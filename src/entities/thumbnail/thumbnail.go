/*
Package thumbnail provides functionality for generating thumbnails
*/
package thumbnail

import (
	"context"
	"io"
	"path"
	"strings"

	thumbnail_errors "github.com/coderchirag/pdf-thumbnail-generator/entities/thumbnail/errors"
	"github.com/coderchirag/pdf-thumbnail-generator/util/fsutil"
)

func GenerateThumbnailFromPdf(ctx context.Context, pdfPath string, quality int) (string, error) {
	doneChan := make(chan bool, 1)
	errChan := make(chan error, 1)

	thumbnailFile, thumbnailPath, err := createThumbnailFile(pdfPath)
	if err != nil {
		return "", err
	}
	go generateThumbnailFromPdfConcurrently(thumbnailFile, pdfPath, quality, doneChan, errChan)

	select {
	case <-doneChan:
		return thumbnailPath, nil
	case err := <-errChan:
		_ = fsutil.RemoveFile(thumbnailPath)
		return "", err
	case <-ctx.Done():
		_ = fsutil.RemoveFile(thumbnailPath)
		return "", thumbnail_errors.ConstructErrorWithCode(ctx.Err(), thumbnail_errors.ContextError)
	}
}

func createThumbnailFile(pdfPath string) (io.Writer, string, error) {
	thumbnailPath := path.Join(
		path.Dir(pdfPath),
		strings.ReplaceAll(path.Base(pdfPath), path.Ext(pdfPath), ".jpeg"),
	)
	file, err := fsutil.CreateFile(thumbnailPath, func() (bool, error) {
		return true, nil
	})
	if err != nil {
		if file != nil {
			_ = fsutil.RemoveFile(thumbnailPath)
		}
		return nil, "", thumbnail_errors.ConstructErrorWithCode(
			err,
			thumbnail_errors.CreateImageFileError,
		)
	}

	return file, thumbnailPath, nil
}

func generateThumbnailFromPdfConcurrently(
	file io.Writer,
	pdfPath string,
	quality int,
	doneChan chan bool,
	errChan chan error,
) {
	img, err := generateThumbnailImageFromPdf(pdfPath)
	if err != nil {
		errChan <- err
		return
	}

	err = encodeImageToFile(img, file, quality)
	if err != nil {
		errChan <- err
		return
	}

	doneChan <- true
}
