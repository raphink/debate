package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Register HTTP handlers
	http.HandleFunc("/", HandleValidateTopic)
	http.HandleFunc("/health", HealthCheck)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting validate-topic function on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(fmt.Errorf("failed to start server: %w", err))
	}
}
