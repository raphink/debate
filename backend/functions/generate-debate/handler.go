package generatedebate

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/raphink/debate/shared/firebase"
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

	// Generate UUID for this debate
	debateID := uuid.New().String()
	log.Printf("Generated debate ID: %s", debateID)

	// Initialize Firestore (if not already initialized)
	ctx := r.Context()
	if firebase.GetClient() == nil {
		if err := firebase.InitFirestore(ctx); err != nil {
			log.Printf("Failed to initialize Firestore (non-blocking): %v", err)
			// Continue anyway - debate will work, just won't be saved
		}
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
	w.Header().Set("X-Debate-Id", debateID)   // Send debate ID to frontend

	// Flush headers
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// Create accumulator for debate messages
	accumulator := NewDebateAccumulator(debateID, req.Topic, req.SelectedPanelists)

	// Wrap writer to accumulate messages
	wrappedWriter := &AccumulatingWriter{
		writer:      w,
		accumulator: accumulator,
	}

	// Stream the debate
	if err := claudeClient.GenerateDebate(ctx, &req, wrappedWriter); err != nil {
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

	// Save completed debate to Firestore (non-blocking)
	go saveDebateToFirestore(context.Background(), accumulator, r.Header.Get("User-Agent"))
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
