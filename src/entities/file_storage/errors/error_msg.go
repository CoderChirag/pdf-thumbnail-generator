package file_storage_errors

import "fmt"

var _ fmt.Stringer = (*ErrorMessageEnum)(nil)

type ErrorMessageEnum int

type errorMessageMap struct {
	Message string
	Code    string
}

const (
	HttpError ErrorMessageEnum = iota
	IOCopyError
	WriterCloseError
	UrlParseError
	CreateDirError
	OperationAlreadyInProgress
	CreateTempFileError
	UnknownError
)

var errorMessagesEnum = map[ErrorMessageEnum]errorMessageMap{
	HttpError: {
		Message: "http client error",
		Code:    "400",
	},
	IOCopyError: {
		Message: "io copy failed",
		Code:    "500",
	},
	WriterCloseError: {
		Message: "failed to close writer",
		Code:    "500",
	},
	UrlParseError: {
		Message: "failed to parse url",
		Code:    "400",
	},
	CreateDirError: {
		Message: "failed to create directory",
		Code:    "500",
	},
	OperationAlreadyInProgress: {
		Message: "operation already in progress",
		Code:    "422",
	},
	CreateTempFileError: {
		Message: "failed to create temp file",
		Code:    "500",
	},
	UnknownError: {
		Message: "unknown error",
		Code:    "500",
	},
}

func (em ErrorMessageEnum) String() string {
	return errorMessagesEnum[em].Message
}

func (em ErrorMessageEnum) Message() string {
	return errorMessagesEnum[em].Message
}

func (em ErrorMessageEnum) Code() string {
	return errorMessagesEnum[em].Code
}
