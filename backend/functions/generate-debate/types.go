package main

// Panelist represents a debate participant
type Panelist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tagline   string `json:"tagline"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatarUrl"`
	Position  string `json:"position"`
}

// DebateRequest represents the incoming request to generate a debate
type DebateRequest struct {
	Topic             string     `json:"topic"`
	SelectedPanelists []Panelist `json:"selectedPanelists"`
}

// StreamChunk represents a single chunk of the streaming response
type StreamChunk struct {
	Type       string `json:"type"`        // "message", "error", "done"
	PanelistID string `json:"panelistId"`  // ID of the speaking panelist
	Text       string `json:"text"`        // Partial or complete text
	Done       bool   `json:"done"`        // Whether streaming is complete
	Error      string `json:"error,omitempty"` // Error message if type="error"
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
	ErrInvalidRequest     = "INVALID_REQUEST"
	ErrInvalidPanelists   = "INVALID_PANELISTS"
	ErrRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
	ErrInternalError      = "INTERNAL_ERROR"
	ErrServiceUnavailable = "SERVICE_UNAVAILABLE"
)
