package getportrait

// PortraitRequest represents a request to fetch a panelist's portrait
type PortraitRequest struct {
	PanelistID   string `json:"panelistId"`
	PanelistName string `json:"panelistName"`
}

// PortraitResponse represents the response with portrait URL
type PortraitResponse struct {
	PanelistID  string `json:"panelistId"`
	PortraitURL string `json:"portraitUrl"`
	Cached      bool   `json:"cached"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code"`
	Retryable bool   `json:"retryable"`
}

// Error codes
const (
	ErrInvalidInput   = "INVALID_INPUT"
	ErrWikimediaError = "WIKIMEDIA_ERROR"
	ErrInternalError  = "INTERNAL_ERROR"
)
