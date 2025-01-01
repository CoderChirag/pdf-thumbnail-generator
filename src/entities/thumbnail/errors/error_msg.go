package thumbnail_errors

import "fmt"

var _ fmt.Stringer = (*ErrorMessageEnum)(nil)

type ErrorMessageEnum int

type ErrorMessageMap struct {
	Message string
	Code    string
}

const (
	PDFParseError ErrorMessageEnum = iota
	CreateImageFileError
	ImageCreationError
	EncodeImageError
	ContextError
	UnknownError
)

var errorMessagesEnum = map[ErrorMessageEnum]ErrorMessageMap{
	PDFParseError: {
		Message: "Error parsing PDF",
		Code:    "400",
	},
	CreateImageFileError: {
		Message: "Error creating image file",
		Code:    "400",
	},
	ImageCreationError: {
		Message: "Error creating image from PDF",
		Code:    "500",
	},
	EncodeImageError: {
		Message: "Error encoding image",
		Code:    "500",
	},
	ContextError: {
		Message: "Context error",
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
