package autocompletetopics

// PanelistSummary represents minimal panelist info for autocomplete
type PanelistSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// DebateSummary represents a debate summary for autocomplete results
type DebateSummary struct {
	ID            string            `json:"id"`
	Topic         string            `json:"topic"`
	Panelists     []PanelistSummary `json:"panelists"`
	PanelistCount int               `json:"panelistCount"`
	CreatedAt     string            `json:"createdAt"`
}

// AutocompleteResponse represents the HTTP response
type AutocompleteResponse struct {
	Debates []DebateSummary `json:"debates"`
}