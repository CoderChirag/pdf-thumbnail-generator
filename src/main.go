package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	thumbnail_controller "github.com/coderchirag/pdf-thumbnail-generator/controllers/thumbnail"
)

func main() {
	fmt.Printf("Process ID: %d\n", os.Getpid())
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 30 * time.Minute,
		IdleTimeout:  60 * time.Second,
	}

	thumbnail_controller.RegisterRoutes()

	fmt.Println("Starting server on port 8080")
	_ = server.ListenAndServe()
}
