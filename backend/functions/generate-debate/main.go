package generatedebate

import (
	"log"
	"net/http"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

// HandleGenerateDebate is the entry point for the Cloud Function
func HandleGenerateDebate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Generate debate request received: %s %s", r.Method, r.URL.Path)

	// Delegate to the existing handler
	handleGenerateDebateImpl(w, r)
}
