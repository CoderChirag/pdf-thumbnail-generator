/*
Package file_storage_errors provides errors for package file_storage
*/
package file_storage_errors

import (
	"fmt"

	"github.com/coderchirag/pdf-thumbnail-generator/types"
)

var _ types.Error = (*FileStorageError)(nil)

type FileStorageError struct {
	err  error
	name string
	msg  string
	code string
}

/* New returns a new instance of *FileStorageError */
func New(err error, msg string, code string) *FileStorageError {
	return &FileStorageError{name: "FileStorageError", msg: msg, code: code, err: err}
}

func NewFromEnum(err error, enum ErrorMessageEnum) *FileStorageError {
	return &FileStorageError{
		name: "FileStorageError",
		msg:  errorMessagesEnum[enum].Message,
		code: errorMessagesEnum[enum].Code,
		err:  err,
	}
}

func (e *FileStorageError) Error() string {
	return fmt.Sprintf("[%s]: %s - %s", e.name, e.msg, e.err.Error())
}

func (e *FileStorageError) Unwrap() error {
	return e.err
}

func (e *FileStorageError) Is(target error) bool {
	_, ok := target.(*FileStorageError)
	return ok
}

func (e *FileStorageError) Name() string {
	return e.name
}

func (e *FileStorageError) Code() string {
	return e.code
}

func (e *FileStorageError) Message() string {
	return e.msg
}
