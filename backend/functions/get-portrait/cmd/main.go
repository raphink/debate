package main

import (
	"log"
	"net/http"
	"os"

	getportrait "github.com/raphink/debate/backend/functions/get-portrait"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	http.HandleFunc("/", getportrait.HandleGetPortrait)

	log.Printf("Starting get-portrait server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
