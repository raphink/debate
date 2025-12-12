package main

import (
	"context"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	autocompletetopics "github.com/raphink/debate/backend/functions/autocomplete-topics"
)

func main() {
	// Register HTTP function
	if err := funcframework.RegisterHTTPFunctionContext(
		context.Background(),
		"/",
		autocompletetopics.AutocompleteTopicsHandler,
	); err != nil {
		log.Fatalf("funcframework.RegisterHTTPFunctionContext: %v\n", err)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
