package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	thumbnail_controller "github.com/coderchirag/pdf-thumbnail-generator/controllers/thumbnail"
	thumbnail_usecase "github.com/coderchirag/pdf-thumbnail-generator/usecases/thumbnail"
)

var server = &http.Server{
	Addr:         ":8080",
	Handler:      http.DefaultServeMux,
	ReadTimeout:  5 * time.Minute,
	WriteTimeout: 30 * time.Minute,
	IdleTimeout:  60 * time.Second,
}

func main() {
	fmt.Printf("Process ID: %d\n", os.Getpid())

	go gracefulShutdown()

	startThumbnailGenerationPipeline()

	thumbnail_controller.RegisterRoutes()

	fmt.Println("Starting server on port 8080")
	_ = server.ListenAndServe()
}

func gracefulShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	shutdownThumbnailGenerationPipeline()
	_ = server.Shutdown(context.Background())
	os.Exit(0)
}

func startThumbnailGenerationPipeline() {
	fmt.Println("Starting thumbnail generation pipeline")
	_ = thumbnail_usecase.NewThumbnailPipeline(5)
}

func shutdownThumbnailGenerationPipeline() {
	fmt.Println("Shutting down thumbnail generation pipeline")
	thumbnail_usecase.GetThumbnailPipeline().Shutdown()
}
