package main

import (
	"log"
	"net/http"
	"os"

	validatetopic "github.com/raphink/debate/functions/validate-topic"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", validatetopic.HandleValidateTopic)

	log.Printf("Starting validate-topic server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
