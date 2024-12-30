// Package errors provides errors for package thumbnail
package errors

import (
	"fmt"

	"github.com/coderchirag/pdf-thumbnail-generator/types"
)

var _ types.Error = (*ThumbnailServiceError)(nil)

type ThumbnailServiceError struct {
	err  error
	name string
	msg  string
	code string
}

/* New returns a new instance of *ThumbnailServiceError */
func New(err error, msg string, code string) *ThumbnailServiceError {
	return &ThumbnailServiceError{name: "ThumbnailServiceError", msg: msg, code: code, err: err}
}

func NewFromEnum(err error, enum ErrorMessageEnum) *ThumbnailServiceError {
	return &ThumbnailServiceError{
		name: "ThumbnailServiceError",
		msg:  errorMessagesEnum[enum].Message,
		code: errorMessagesEnum[enum].Code,
		err:  err,
	}
}

func (e *ThumbnailServiceError) Error() string {
	return fmt.Sprintf("[%s]: %s - %s", e.name, e.msg, e.err.Error())
}

func (e *ThumbnailServiceError) Unwrap() error {
	return e.err
}

func (e *ThumbnailServiceError) Name() string {
	return e.name
}

func (e *ThumbnailServiceError) Code() string {
	return e.code
}

func (e *ThumbnailServiceError) Message() string {
	return e.msg
}
