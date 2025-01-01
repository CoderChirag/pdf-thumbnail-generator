package main

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coderchirag/pdf-thumbnail-generator/thumbnail"
	"github.com/coderchirag/pdf-thumbnail-generator/types"
)

func pdf(w http.ResponseWriter, req *http.Request) {
	svc, err := thumbnail.NewThumbnailService(
		"https://bfhldatapipeline.blob.core.windows.net/health-vault/625813eabd79a9da06ca4f99%2F0a42fcec-39c4-4f1c-af7e-0c2ff13ba79f.pdf?sv=2021-08-06&spr=https&st=2024-12-30T12%3A21%3A23Z&se=2025-12-30T12%3A21%3A23Z&sr=b&sp=racw&sig=sp6PrWoVK1X3zup%2FTXhLDeJL1WB%2B1Yj7M2Unfc4aCY4%3D",
	)
	handleError(err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s, err := svc.GenerateThumbnail(ctx, 40)
		handleError(err)
		fmt.Println(s)

		_, err = w.Write([]byte(s))
		handleError(err)
	}()
	time.Sleep(200 * time.Millisecond)
	cancel()
	_, err = w.Write([]byte("done"))
	handleError(err)
}

func main() {
	fmt.Printf("%d\n", os.Getpid())
	http.HandleFunc("/pdf", pdf)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 30 * time.Minute,
		IdleTimeout:  60 * time.Second,
	}

	_ = server.ListenAndServe()
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
