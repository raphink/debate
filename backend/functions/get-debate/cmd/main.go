package main

import (
	"log"
	"net/http"
	"os"

	getdebate "github.com/raphink/debate/functions/get-debate"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", getdebate.HandleGetDebate)

	log.Printf("Starting get-debate server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
