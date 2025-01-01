package thumbnail_controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/coderchirag/pdf-thumbnail-generator/types"
	thumbnail_usecase "github.com/coderchirag/pdf-thumbnail-generator/usecases/thumbnail"
)

func generateThumbnail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pdfUrl := r.FormValue("pdfUrl")
	thumbnailPath, err := thumbnail_usecase.GenerateThumbnailSequentially(
		ctx,
		pdfUrl,
	)
	if err != nil {
		handleError(err, w)
		return
	}
	defer os.Remove(thumbnailPath)

	file, err := os.Open(thumbnailPath)
	if err != nil {
		handleError(err, w)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		handleError(err, w)
		return
	}
}

func handleError(err error, w http.ResponseWriter) {
	var typeErr types.Error
	var code string
	if errors.As(err, &typeErr) {
		code = typeErr.Code()
		fmt.Printf("Error Code: %s\n", code)
		fmt.Println(typeErr)
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
