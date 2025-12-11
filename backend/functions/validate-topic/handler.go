package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleValidateTopic is the HTTP handler for the validate-topic Cloud Function
func HandleValidateTopic(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
	w.Header().Set("Content-Type", "application/json")

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

	// Validate topic with Claude
	ctx := r.Context()
	isRelevant, message, err := claudeClient.ValidateTopicRelevance(ctx, sanitizedTopic)
	if err != nil {
		log.Printf("Error validating topic with Claude: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		retryAfter := 30
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:      "The AI service is temporarily unavailable. Please try again in a few moments.",
			Code:       ErrServiceUnavailable,
			Retryable:  true,
			RetryAfter: &retryAfter,
		})
		return
	}

	// Return successful response
	response := TopicValidationResponse{
		IsRelevant: isRelevant,
		Message:    message,
		Topic:      sanitizedTopic,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// HealthCheck is a simple health check handler
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"healthy"}`)
}
