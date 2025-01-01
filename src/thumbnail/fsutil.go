package thumbnail

import (
	stderrors "errors"
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"

	"github.com/coderchirag/pdf-thumbnail-generator/thumbnail/errors"
)

func getFileAndDirPath(url *url.URL) (string, string) {
	var filePath string
	var fileDirPath string
	if path.Ext(url.Path) == ".pdf" {
		filePath = path.Join(os.TempDir(), "pdf-parsing", path.Base(url.Path))
	} else {
		filePath = path.Join(os.TempDir(), "pdf-parsing", uuid.NewString()+".pdf")
	}
	fileDirPath = path.Dir(filePath)

	return filePath, fileDirPath
}

func createDirIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0o755); err != nil {
			return handleError(err)
		}
	} else if err != nil {
		return handleError(err)
	}

	return nil
}

func createFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	if err == nil {
		// File already exists
		return nil, errors.NewFromEnum(
			stderrors.New("Operation already in progress."),
			errors.ValidationError,
		)
	} else if !os.IsNotExist(err) {
		return nil, handleError(err)
	}
	file, err := os.Create(path)
	if err != nil {
		if file != nil {
			_ = removeFile(path)
		}
		return nil, handleError(err)
	}

	return file, nil
}

func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return handleError(err)
	}
	return nil
}
