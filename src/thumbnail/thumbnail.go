// Package thumbnail provides functionality for generating thumbnails from URLs
package thumbnail

import (
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/coderchirag/pdf-thumbnail-generator/types"
	"github.com/gen2brain/go-fitz"
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

func (s *ThumbnailService) GenerateThumbnail(ctx context.Context, quality int) (string, error) {
	filePath, err := s.downloadFile(ctx)
	fmt.Println("filePath", filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if filePath != "" {
			_ = os.Remove(filePath)
		}
	}()
	thumbnailPath, err := generateThumbnail(ctx, filePath, quality)
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
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

func generateThumbnail(
	ctx context.Context,
	filePath string,
	quality int,
) (string, error) {
	thumbnailPathChan := make(chan string, 1)
	errChan := make(chan error, 1)

	doc, err := fitz.New(filePath)
	if err != nil {
		return "", handleError(err)
	}

	fmt.Println("Generating thumbnail")
	go func() {
		thumbnailPath := path.Join(
			path.Dir(filePath),
			strings.ReplaceAll(path.Base(filePath), path.Ext(filePath), ".jpeg"),
		)
		file, err := createFile(thumbnailPath)
		if err != nil {
			errChan <- err
			return
		}
		img, err := doc.Image(0)
		if err != nil {
			errChan <- handleError(err)
			return
		}

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
		if err != nil {
			errChan <- handleError(err)
			return
		}
		thumbnailPathChan <- thumbnailPath
	}()

	select {
	case thumbnailPath := <-thumbnailPathChan:
		fmt.Println("Thumbnail generated")
		return thumbnailPath, nil
	case err := <-errChan:
		fmt.Println("Error generating thumbnail")
		return "", err
	case <-ctx.Done():
		fmt.Println("Context done")
		return "", handleError(ctx.Err())
	}
}
