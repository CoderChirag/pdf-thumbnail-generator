/*
Package fsutil provides utility functions for file system operations.
*/
package fsutil

import (
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"
)

/*
ConstructTempFileAndDirPath returns the file path and the directory path inside the given base directory inside the temp directory for the given url.

If the url has the given extension, the file path is returned as the url path.

If the url does not have the given extension, a new file path is generated using the uuid.
*/
func ConstructTempFileAndDirPath(url *url.URL, baseDir string, ext string) (string, string) {
	var filePath string
	var fileDirPath string
	if path.Ext(url.Path) == ext {
		filePath = path.Join(os.TempDir(), baseDir, path.Base(url.Path))
	} else {
		filePath = path.Join(os.TempDir(), baseDir, uuid.NewString()+ext)
	}
	fileDirPath = path.Dir(filePath)

	return filePath, fileDirPath
}

func CreateDirIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0o755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

/*
CreateFile creates a file at the given path. If the file already exists, it calls the onAlreadyExists() function.

If onAlreadyExists() returns false or an non nil error, the file is not truncated and the error is returned.

If onAlreadyExists() returns true and an nil error, the existing file is truncated.
*/
func CreateFile(path string, onAlreadyExists func() (bool, error)) (*os.File, error) {
	_, err := os.Stat(path)
	if err == nil {
		// File already exists
		ifContinue, e := onAlreadyExists()
		if !ifContinue || e != nil {
			return nil, e
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		if file != nil {
			_ = RemoveFile(path)
		}
		return nil, err
	}

	return file, nil
}

func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
