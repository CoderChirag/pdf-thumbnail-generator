/*
Package file_storage provides functions to download or upload files from http.
*/
package file_storage

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"

	file_storage_errors "github.com/coderchirag/pdf-thumbnail-generator/entities/file_storage/errors"
	"github.com/coderchirag/pdf-thumbnail-generator/util/fsutil"
)

func DownloadFileToTempDir(
	ctx context.Context,
	fileUrl string,
	baseDir string,
	ext string,
) (string, error) {
	pipeReader, pipeWriter := io.Pipe()
	defer func() {
		if err := pipeReader.Close(); err != nil {
			_ = pipeReader.CloseWithError(err)
		}
	}()

	//nolint:errcheck
	go downloadFileFromUrl(ctx, fileUrl, pipeWriter)

	file, filePath, err := createTempFileFromUrl(fileUrl, baseDir, ext)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, pipeReader)
	if err != nil {
		_ = fsutil.RemoveFile(filePath)
		if errors.Is(err, &file_storage_errors.FileStorageError{}) {
			return "", err
		}
		return "", file_storage_errors.ConstructErrorWithCode(
			err,
			file_storage_errors.IOCopyError,
		)
	}

	return filePath, nil
}

func createTempFileFromUrl(fileUrl string, baseDir string, ext string) (*os.File, string, error) {
	parsedURL, err := url.Parse(fileUrl)
	if err != nil {
		return nil, "", file_storage_errors.ConstructErrorWithCode(
			err,
			file_storage_errors.UrlParseError,
		)
	}

	filePath, dir := fsutil.ConstructTempFileAndDirPath(
		parsedURL,
		baseDir,
		ext,
	)
	err = fsutil.CreateDirIfNotExists(dir)
	if err != nil {
		return nil, "", file_storage_errors.ConstructErrorWithCode(
			err,
			file_storage_errors.CreateDirError,
		)
	}

	file, err := fsutil.CreateFile(filePath, func() (bool, error) {
		return false, file_storage_errors.ConstructErrorWithCode(
			err,
			file_storage_errors.OperationAlreadyInProgress,
		)
	})
	if err != nil {
		if file != nil {
			_ = fsutil.RemoveFile(filePath)
		}
		if errors.Is(err, &file_storage_errors.FileStorageError{}) {
			return nil, "", err
		}
		return nil, "", file_storage_errors.ConstructErrorWithCode(
			err,
			file_storage_errors.CreateTempFileError,
		)
	}

	return file, filePath, nil
}
