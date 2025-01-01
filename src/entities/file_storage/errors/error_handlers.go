package file_storage_errors

func ConstructError(e error) error {
	return NewFromEnum(
		e,
		UnknownError,
	)
}

func ConstructErrorWithCode(e error, code ErrorMessageEnum) error {
	return NewFromEnum(
		e,
		code,
	)
}

func CloseIo[T interface {
	Close() error
	CloseWithError(error) error
}](c T) {
	if err := c.Close(); err != nil {
		_ = c.CloseWithError(ConstructErrorWithCode(err, WriterCloseError))
	}
}

func CloseIoWithErrorAndCode[T interface {
	Close() error
	CloseWithError(error) error
}](c T, err error, code ErrorMessageEnum) error {
	e := ConstructErrorWithCode(err, code)
	_ = c.CloseWithError(e)
	return e
}
