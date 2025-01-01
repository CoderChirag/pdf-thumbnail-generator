package thumbnail

import (
	"image"
	"image/jpeg"
	"io"

	thumbnail_errors "github.com/coderchirag/pdf-thumbnail-generator/entities/thumbnail/errors"
)

func encodeImageToFile(img *image.RGBA, file io.Writer, quality int) error {
	err := jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return thumbnail_errors.ConstructErrorWithCode(err, thumbnail_errors.EncodeImageError)
	}

	return nil
}
