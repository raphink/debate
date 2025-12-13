package listdebates

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/raphink/debate/shared/firebase"
	"github.com/raphink/debate/shared/sanitize"
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

func HandleListDebates(w http.ResponseWriter, r *http.Request) {
	ListDebatesHandler(w, r)
}

// ListDebatesHandler handles GET requests to fetch debate history
func ListDebatesHandler(w http.ResponseWriter, r *http.Request) {
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
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	queryParam := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Initialize Firestore client if needed
	ctx := r.Context()
	client := firebase.GetClient()
	if client == nil {
		if err := firebase.InitFirestore(ctx); err != nil {
			log.Printf("Failed to initialize Firestore: %v", err)
			sendError(w, "Failed to initialize database connection", http.StatusInternalServerError)
			return
		}
		client = firebase.GetClient()
	}

	// Autocomplete mode: if q parameter is provided
	if queryParam != "" {
		// Validate query length
		queryParam = strings.TrimSpace(queryParam)
		if len(queryParam) < 3 {
			sendError(w, "Query must be at least 3 characters", http.StatusBadRequest)
			return
		}

		if len(queryParam) > 500 {
			sendError(w, "Query must be less than 500 characters", http.StatusBadRequest)
			return
		}

		// Sanitize query input
		sanitizedQuery := sanitize.StripHTML(queryParam)
		sanitizedQuery = strings.ToLower(sanitizedQuery)

		// Query autocomplete debates
		debates, err := autocompleteDebates(ctx, client, sanitizedQuery)
		if err != nil {
			log.Printf("Autocomplete query failed: query=%q, sanitized=%q, error=%v", queryParam, sanitizedQuery, err)
			sendError(w, "Failed to query autocomplete results", http.StatusInternalServerError)
			return
		}

		log.Printf("Autocomplete query successful: query=%q, results=%d", sanitizedQuery, len(debates))

		// Send autocomplete response (no pagination for autocomplete)
		response := ListDebatesResponse{
			Debates: debates,
			Total:   len(debates),
			HasMore: false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Normal list mode: pagination
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

	// Query debates
	debates, total, err := queryDebates(ctx, client, limit, offset)
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
