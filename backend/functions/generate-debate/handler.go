package generatedebate

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// handleGenerateDebateImpl handles debate generation requests with SSE streaming
func handleGenerateDebateImpl(w http.ResponseWriter, r *http.Request) {
	// Enable CORS - allow configured origin or localhost for dev
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		sendError(w, "Method not allowed", ErrInvalidRequest, false, http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req DebateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", ErrInvalidRequest, false, http.StatusBadRequest)
		return
	}

	// Validate request
	if err := ValidateDebateRequest(&req); err != nil {
		sendError(w, err.Error(), ErrInvalidPanelists, false, http.StatusBadRequest)
		return
	}

	// Create Claude client
	claudeClient, err := NewClaudeClient()
	if err != nil {
		log.Printf("Failed to create Claude client: %v", err)
		sendError(w, "Service configuration error", ErrInternalError, true, http.StatusInternalServerError)
		return
	}

	// Set up Server-Sent Events headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

	// Flush headers
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// Stream the debate
	if err := claudeClient.GenerateDebate(r.Context(), &req, w); err != nil {
		log.Printf("Error generating debate: %v", err)
		errorChunk := StreamChunk{
			Type:  "error",
			Error: "Failed to generate debate. Please try again.",
		}
		json.NewEncoder(w).Encode(errorChunk)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		return
	}
}

// sendError sends a JSON error response
func sendError(w http.ResponseWriter, message, code string, retryable bool, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Error:     message,
		Code:      code,
		Retryable: retryable,
	}

	json.NewEncoder(w).Encode(errorResponse)
}
