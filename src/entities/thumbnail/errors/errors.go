/*
Package thumbnail_errors provides errors for package thumbnail
*/
package thumbnail_errors

import (
	"fmt"

	"github.com/coderchirag/pdf-thumbnail-generator/types"
)

var _ types.Error = (*ThumbnailError)(nil)

type ThumbnailError struct {
	err  error
	name string
	msg  string
	code string
}

/* New returns a new instance of *ThumbnailError */
func New(err error, msg string, code string) *ThumbnailError {
	return &ThumbnailError{name: "ThumbnailError", msg: msg, code: code, err: err}
}

func NewFromEnum(err error, enum ErrorMessageEnum) *ThumbnailError {
	return &ThumbnailError{
		name: "ThumbnailError",
		msg:  errorMessagesEnum[enum].Message,
		code: errorMessagesEnum[enum].Code,
		err:  err,
	}
}

func (e *ThumbnailError) Error() string {
	return fmt.Sprintf("[%s]: %s - %s", e.name, e.msg, e.err.Error())
}

func (e *ThumbnailError) Unwrap() error {
	return e.err
}

func (e *ThumbnailError) Is(target error) bool {
	_, ok := target.(*ThumbnailError)
	return ok
}

func (e *ThumbnailError) Name() string {
	return e.name
}

func (e *ThumbnailError) Code() string {
	return e.code
}

func (e *ThumbnailError) Message() string {
	return e.msg
}
