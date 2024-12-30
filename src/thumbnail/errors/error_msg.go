package errors

import "fmt"

var _ fmt.Stringer = (*ErrorMessageEnum)(nil)

type ErrorMessageEnum int

type ErrorMessageMap struct {
	Message string
	Code    string
}

const (
	ValidationError ErrorMessageEnum = iota
	UnknownError
)

var errorMessagesEnum = map[ErrorMessageEnum]ErrorMessageMap{
	ValidationError: {
		Message: "validation error",
		Code:    "400",
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
