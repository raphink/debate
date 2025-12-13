package listdebates

import (
	"time"
)

// PanelistInfo represents minimal panelist information
type PanelistInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	Tagline   string `json:"tagline"`
	Bio       string `json:"bio"`
}

// DebateSummary represents a debate summary for the list view
type DebateSummary struct {
	ID        string         `json:"id"`
	Topic     string         `json:"topic"`
	Panelists []PanelistInfo `json:"panelists"`
	StartedAt time.Time      `json:"startedAt"`
}

// ListDebatesResponse is the response structure for the list endpoint
type ListDebatesResponse struct {
	Debates []DebateSummary `json:"debates"`
	Total   int             `json:"total"`
	HasMore bool            `json:"hasMore"`
}

// ErrorResponse is the error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}
