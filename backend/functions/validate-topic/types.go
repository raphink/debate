package main

// TopicValidationRequest represents the incoming request to validate a topic
type TopicValidationRequest struct {
	Topic string `json:"topic"`
}

// TopicValidationResponse represents the response after validating a topic
type TopicValidationResponse struct {
	IsRelevant bool   `json:"isRelevant"`
	Message    string `json:"message"`
	Topic      string `json:"topic"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error      string `json:"error"`
	Code       string `json:"code"`
	Retryable  bool   `json:"retryable"`
	RetryAfter *int   `json:"retryAfter,omitempty"`
}

// Error codes
const (
	ErrInvalidTopicLength  = "INVALID_TOPIC_LENGTH"
	ErrInvalidTopicContent = "INVALID_TOPIC_CONTENT"
	ErrRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"
	ErrInternalError       = "INTERNAL_ERROR"
	ErrServiceUnavailable  = "SERVICE_UNAVAILABLE"
)