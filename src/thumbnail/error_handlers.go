package thumbnail

import (
	stderrors "errors"
	"net/url"

	"github.com/coderchirag/pdf-thumbnail-generator/thumbnail/errors"
)

func validateThumbnailURL(fileURL string) error {
	if fileURL == "" {
		err := errors.NewFromEnum(
			stderrors.New("Url should not be empty."),
			errors.ValidationError,
		)
		return err
	}

	_, err := url.Parse(fileURL)
	if err != nil {
		err := errors.NewFromEnum(
			stderrors.New("Invalid url."),
			errors.ValidationError,
		)
		return err
	}

	return nil
}

func handleError(e error) error {
	return errors.NewFromEnum(
		e,
		errors.UnknownError,
	)
}
