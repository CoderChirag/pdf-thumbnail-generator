// Package thumbnail provides functionality for generating thumbnails from URLs
package thumbnail

import (
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"

	thumbnailErrors "github.com/coderchirag/pdf-thumbnail-generator/thumbnail/errors"
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
	_, pipeWriter := io.Pipe()

	go func() {
		defer func() {
			if err := pipeWriter.Close(); err != nil {
				// Handle or log the error as needed
				_ = pipeWriter.CloseWithError(err)
			}
		}()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			pipeWriter.CloseWithError(err)
			return
		}

		defer func() {
			if cerr := res.Body.Close(); cerr != nil && err == nil {
				pipeWriter.CloseWithError(cerr)
			}
		}()

		_, err = io.Copy(pipeWriter, res.Body)
		if err != nil {
			pipeWriter.CloseWithError(err)
		}
	}()

	parsedURL, err := url.Parse(s.url)
	if err != nil {
		return "", handleError(err)
	}
	var filePath string
	if path.Ext(parsedURL.Path) == ".pdf" {
		filePath = path.Join(os.TempDir(), "pdf-parsing", path.Base(parsedURL.Path))
	} else {
		filePath = path.Join(os.TempDir(), "pdf-parsing", uuid.NewString()+".pdf")
	}

	return filePath, nil
}

func validateThumbnailURL(fileURL string) error {
	if fileURL == "" {
		err := thumbnailErrors.New(
			stderrors.New("Url should not be empty."),
			fmt.Sprint(thumbnailErrors.ValidationError),
			thumbnailErrors.ValidationError.Code(),
		)
		return err
	}

	_, err := url.Parse(fileURL)
	if err != nil {
		err := thumbnailErrors.New(
			stderrors.New("Invalid url."),
			fmt.Sprint(thumbnailErrors.ValidationError),
			thumbnailErrors.ValidationError.Code(),
		)
		return err
	}

	return nil
}

func handleError(e error) error {
	return thumbnailErrors.New(
		e,
		fmt.Sprint(thumbnailErrors.UnknownError),
		thumbnailErrors.UnknownError.Code(),
	)
}
