package getdebate

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/google/uuid"
	"github.com/raphink/debate/shared/firebase"
)

var allowedOrigin string

func init() {
	allowedOrigin = os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}
	log.Printf("ALLOWED_ORIGIN set to: %s", allowedOrigin)

	// Initialize Firestore client
	ctx := context.Background()
	if err := firebase.InitFirestore(ctx); err != nil {
		log.Printf("Failed to initialize Firestore: %v", err)
	}
}

// HandleGetDebate handles GET requests to retrieve a debate by UUID
func HandleGetDebate(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Get UUID from query parameter
	debateID := strings.TrimSpace(r.URL.Query().Get("id"))
	if debateID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing debate ID parameter",
		})
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(debateID); err != nil {
		log.Printf("Invalid UUID format: %s", debateID)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid debate ID format",
		})
		return
	}

	// Retrieve debate from Firestore
	ctx := r.Context()
	debate, err := firebase.GetDebate(ctx, debateID)
	if err != nil {
		// Check if it's a not found error
		if strings.Contains(err.Error(), "not found") {
			log.Printf("Debate not found: %s", debateID)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Debate not found",
			})
			return
		}

		// Internal server error
		log.Printf("Failed to retrieve debate %s: %v", debateID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to load debate",
		})
		return
	}

	// Return debate data
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(debate); err != nil {
		log.Printf("Failed to encode debate response: %v", err)
	}
}
