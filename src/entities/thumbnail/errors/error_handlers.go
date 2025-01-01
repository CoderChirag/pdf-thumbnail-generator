package thumbnail_errors

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
