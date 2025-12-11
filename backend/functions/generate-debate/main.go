package main

import (
	"log"
	"net/http"
)

// HandleGenerateDebate is the entry point for the Cloud Function
func HandleGenerateDebate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Generate debate request received: %s %s", r.Method, r.URL.Path)
	
	// Delegate to the existing handler
	handleGenerateDebateImpl(w, r)
}
