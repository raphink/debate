package main

import (
	"log"
	"net/http"
	"os"

	generatedebate "github.com/raphink/debate/functions/generate-debate"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", generatedebate.HandleGenerateDebate)

	log.Printf("Starting generate-debate server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
