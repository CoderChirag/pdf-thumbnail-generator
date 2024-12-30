package main

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/coderchirag/pdf-thumbnail-generator/thumbnail"
	"github.com/coderchirag/pdf-thumbnail-generator/types"
)

func main() {
	svc, err := thumbnail.NewThumbnailService(
		"https://drive.usercontent.google.com/download?id=1Qspoh1gKWl4KS9MPj0LCQE3hKdKZInmb&export=download&authuser=0",
	)
	handleError(err)

	s, err := svc.GenerateThumbnail(context.Background())
	handleError(err)
	fmt.Println(s)
}

func handleError(err error) {
	if err != nil {
		var typeErr types.Error
		var code string
		if stderrors.As(err, &typeErr) {
			code = typeErr.Code()
			fmt.Printf("Error Code: %s\n", code)
		}
		panic(err)
	}
}
