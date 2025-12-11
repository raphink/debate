package validatetopic
package main





























)	ErrServiceUnavailable  = "SERVICE_UNAVAILABLE"	ErrInternalError       = "INTERNAL_ERROR"	ErrRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"	ErrInvalidTopicContent = "INVALID_TOPIC_CONTENT"	ErrInvalidTopicLength  = "INVALID_TOPIC_LENGTH"const (// Error codes}	RetryAfter *int   `json:"retryAfter,omitempty"`	Retryable  bool   `json:"retryable"`	Code       string `json:"code"`	Error      string `json:"error"`type ErrorResponse struct {// ErrorResponse represents an error response from the API}	Topic      string `json:"topic"`	Message    string `json:"message"`	IsRelevant bool   `json:"isRelevant"`type TopicValidationResponse struct {// TopicValidationResponse represents the response after validating a topic}	Topic string `json:"topic"`type TopicValidationRequest struct {// TopicValidationRequest represents the incoming request to validate a topic