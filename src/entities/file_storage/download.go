package file_storage

import (
	"context"
	"fmt"
	"io"
	"net/http"

	file_storage_errors "github.com/coderchirag/pdf-thumbnail-generator/entities/file_storage/errors"
	"github.com/coderchirag/pdf-thumbnail-generator/types"
)

func downloadFileFromUrl(
	ctx context.Context,
	url string,
	writer types.WriteCloserWithError,
) error {
	defer file_storage_errors.CloseIo(writer)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return file_storage_errors.CloseIoWithErrorAndCode(
			writer,
			err,
			file_storage_errors.HttpError,
		)
	}

	res, err := http.DefaultClient.Do(req)
	defer closeResponseBody(res)
	if err != nil {
		return file_storage_errors.CloseIoWithErrorAndCode(
			writer,
			err,
			file_storage_errors.HttpError,
		)
	}

	if res.StatusCode != http.StatusOK {
		return file_storage_errors.CloseIoWithErrorAndCode(
			writer,
			fmt.Errorf("unexpected status code: %d", res.StatusCode),
			file_storage_errors.HttpError,
		)
	}

	_, err = io.Copy(writer, res.Body)
	if err != nil {
		return file_storage_errors.CloseIoWithErrorAndCode(
			writer,
			err,
			file_storage_errors.IOCopyError,
		)
	}

	return nil
}

func closeResponseBody(res *http.Response) {
	if res != nil && res.Body != nil {
		_ = res.Body.Close()
	}
}
