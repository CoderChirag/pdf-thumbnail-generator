package thumbnail

import (
	"image"

	"github.com/gen2brain/go-fitz"

	thumbnail_errors "github.com/coderchirag/pdf-thumbnail-generator/entities/thumbnail/errors"
)

func generateThumbnailImageFromPdf(filePath string) (*image.RGBA, error) {
	doc, err := fitz.New(filePath)
	if err != nil {
		return nil, thumbnail_errors.ConstructErrorWithCode(err, thumbnail_errors.PDFParseError)
	}
	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		return nil, thumbnail_errors.ConstructErrorWithCode(
			err,
			thumbnail_errors.ImageCreationError,
		)
	}
	return img, nil
}
