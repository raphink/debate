package validatetopic

import (
	"encoding/json"
	"log"
	"net/http"
)

// handleValidateTopicImpl is the HTTP handler for the validate-topic Cloud Function
func handleValidateTopicImpl(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

	// Set SSE headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only accept POST requests
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:     "Method not allowed. Use POST.",
			Code:      ErrInternalError,
			Retryable: false,
		})
		return
	}

	// Parse request body
	var req TopicValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:     "Invalid request body",
			Code:      ErrInternalError,
			Retryable: false,
		})
		return
	}

	// Validate input
	if err := ValidateTopicInput(req.Topic); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MapErrorToResponse(err))
		return
	}

	// Sanitize the topic
	sanitizedTopic := SanitizeTopic(req.Topic)

	// Limit suggested names to 5 and sanitize
	suggestedNames := make([]string, 0, 5)
	for i, name := range req.SuggestedNames {
		if i >= 5 {
			break
		}
		sanitized := SanitizeTopic(name) // Reuse sanitization logic
		if sanitized != "" {
			suggestedNames = append(suggestedNames, sanitized)
		}
	}
	log.Printf("Received suggested names: %v, sanitized to: %v", req.SuggestedNames, suggestedNames)

	// Create Claude client
	claudeClient, err := NewClaudeClient()
	if err != nil {
		log.Printf("Error creating Claude client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:     "Service configuration error. Please try again later.",
			Code:      ErrInternalError,
			Retryable: true,
		})
		return
	}

	// Validate topic and stream panelist suggestions from Claude
	if err := claudeClient.ValidateTopicAndSuggestPanelists(r.Context(), sanitizedTopic, suggestedNames, w); err != nil {
		log.Printf("Error validating topic with Claude: %v", err)
		// Send error chunk
		errorChunk := map[string]string{
			"type":  "error",
			"error": "The AI service is temporarily unavailable. Please try again in a few moments.",
		}
		json.NewEncoder(w).Encode(errorChunk)
		return
	}
}
