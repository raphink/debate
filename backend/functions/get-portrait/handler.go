package getportrait

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	// Valid panelist ID pattern: alphanumeric with hyphens
	panelistIDPattern = regexp.MustCompile(`^[a-zA-Z0-9\-]{1,50}$`)
)

// HandleGetPortrait is the HTTP handler for the get-portrait Cloud Function
func HandleGetPortrait(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers - allow configured origin or localhost for dev
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Only accept POST requests
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed", ErrInvalidInput, false)
		return
	}

	// Parse request body
	var req PortraitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request format", ErrInvalidInput, false)
		return
	}

	// Validate input
	if req.PanelistID == "" || req.PanelistName == "" {
		respondWithError(w, http.StatusBadRequest, "Panelist ID and name are required", ErrInvalidInput, false)
		return
	}

	if !panelistIDPattern.MatchString(req.PanelistID) {
		respondWithError(w, http.StatusBadRequest, "Invalid panelist ID format", ErrInvalidInput, false)
		return
	}

	if len(req.PanelistName) > 100 {
		respondWithError(w, http.StatusBadRequest, "Panelist name too long (max 100 characters)", ErrInvalidInput, false)
		return
	}

	// Sanitize name
	sanitizedName := strings.TrimSpace(req.PanelistName)
	if sanitizedName == "" {
		respondWithError(w, http.StatusBadRequest, "Panelist name cannot be empty", ErrInvalidInput, false)
		return
	}

	// Check cache first
	if cachedURL, found := portraitCache.Get(req.PanelistID); found {
		log.Printf("Cache hit for %s", req.PanelistID)
		respondWithSuccess(w, PortraitResponse{
			PanelistID:  req.PanelistID,
			PortraitURL: cachedURL,
			Cached:      true,
		})
		return
	}

	// Fetch from Wikimedia
	wiki := NewWikimediaAPI()
	portraitURL := wiki.FetchPortraitURL(sanitizedName)

	// Use placeholder if not found
	if portraitURL == "" {
		portraitURL = "placeholder-avatar.svg"
	}

	// Cache the result
	portraitCache.Set(req.PanelistID, portraitURL)

	// Return response
	respondWithSuccess(w, PortraitResponse{
		PanelistID:  req.PanelistID,
		PortraitURL: portraitURL,
		Cached:      false,
	})
}

func respondWithError(w http.ResponseWriter, status int, message, code string, retryable bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:     message,
		Code:      code,
		Retryable: retryable,
	})
}

func respondWithSuccess(w http.ResponseWriter, response PortraitResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
