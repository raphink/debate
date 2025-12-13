package listdebates

import (
	"time"
)

// PanelistInfo represents minimal panelist information
type PanelistInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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





}	Error string `json:"error"`type ErrorResponse struct {// ErrorResponse is the error response structure