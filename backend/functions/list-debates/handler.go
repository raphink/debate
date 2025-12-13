package listdebates

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

// initFirestore initializes the Firestore client
func initFirestore() error {
	ctx := context.Background()

	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath == "" {
		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("GCP_PROJECT_ID not set")
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}

	firestoreClient = client
	return nil
}

// ListDebatesHandler handles GET requests to fetch debate history
func ListDebatesHandler(w http.ResponseWriter, r *http.Request) {
	// Handle CORS
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	offset := 0 // default

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			sendError(w, "Invalid limit: must be between 1 and 100", http.StatusBadRequest)
			return
		}
	}

	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			sendError(w, "Invalid offset: must be >= 0", http.StatusBadRequest)
			return
		}
	}

	// Initialize Firestore client if needed
	if firestoreClient == nil {
		if err := initFirestore(); err != nil {
			log.Printf("Failed to initialize Firestore: %v", err)
			sendError(w, "Failed to initialize database connection", http.StatusInternalServerError)
			return
		}
	}

	// Query debates
	ctx := context.Background()
	debates, total, err := queryDebates(ctx, firestoreClient, limit, offset)
	if err != nil {
		log.Printf("Failed to query debates: %v", err)
		sendError(w, "Failed to query debates from Firestore", http.StatusInternalServerError)
		return
	}

	// Calculate hasMore
	hasMore := (offset + len(debates)) < total

	// Send response
	response := ListDebatesResponse{
		Debates: debates,
		Total:   total,
		HasMore: hasMore,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendError sends a JSON error response
func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
