// Package thumbnail provides functionality for generating thumbnails from URLs
package thumbnail

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/coderchirag/pdf-thumbnail-generator/types"
)

type ThumbnailService struct {
	url string
}

func NewThumbnailService(url string) (*ThumbnailService, error) {
	if err := validateThumbnailURL(url); err != nil {
		return nil, err
	}

	return &ThumbnailService{url: url}, nil
}

func (s *ThumbnailService) GenerateThumbnail(ctx context.Context) (string, error) {
	return s.downloadFile(ctx)
}

func (s *ThumbnailService) downloadFile(ctx context.Context) (string, error) {
	pipeReader, pipeWriter := io.Pipe()
	defer func() {
		if err := pipeReader.Close(); err != nil {
			_ = pipeReader.CloseWithError(err)
		}
	}()

	go s.downloadFileFromUrl(ctx, pipeWriter)

	parsedURL, err := url.Parse(s.url)
	if err != nil {
		return "", handleError(err)
	}

	filePath, dir := getFileAndDirPath(parsedURL)
	err = createDirIfNotExists(dir)
	if err != nil {
		return "", err
	}

	file, err := createFile(filePath)
	if err != nil {
		if file != nil {
			_ = removeFile(filePath)
		}
		return "", err
	}

	_, err = io.Copy(file, pipeReader)
	if err != nil {
		_ = removeFile(filePath)
		return "", handleError(err)
	}

	return filePath, nil
}

func (s *ThumbnailService) downloadFileFromUrl(
	ctx context.Context,
	writer types.WriteCloserWithError,
) {
	defer func() {
		if err := writer.Close(); err != nil {
			// Handle or log the error as needed
			_ = writer.CloseWithError(err)
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		_ = writer.CloseWithError(err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			if cerr := res.Body.Close(); cerr != nil {
				_ = writer.CloseWithError(cerr)
			}
		}
	}()
	if err != nil {
		_ = writer.CloseWithError(err)
		return
	}

	_, err = io.Copy(writer, res.Body)
	if err != nil {
		_ = writer.CloseWithError(err)
	}
}
