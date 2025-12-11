package validatetopic

import (
	"log"
	"net/http"
)

// HandleValidateTopic is the entry point for the Cloud Function
func HandleValidateTopic(w http.ResponseWriter, r *http.Request) {
	log.Printf("Validate topic request received: %s %s", r.Method, r.URL.Path)
	
	// Delegate to the existing handler
	handleValidateTopicImpl(w, r)
}
