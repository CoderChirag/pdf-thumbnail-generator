package thumbnail_usecase

import (
	"context"
	"os"

	"github.com/coderchirag/pdf-thumbnail-generator/entities/file_storage"
	"github.com/coderchirag/pdf-thumbnail-generator/entities/thumbnail"
)

func GenerateThumbnailSequentially(
	ctx context.Context,
	pdfUrl string,
) (string, error) {
	pdfPath, err := file_storage.DownloadFileToTempDir(ctx, pdfUrl, baseDir, pdfExt)
	if err != nil {
		return "", err
	}
	defer os.Remove(pdfPath)

	thumbnailPath, err := thumbnail.GenerateThumbnailFromPdf(ctx, pdfPath, quality)
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}
