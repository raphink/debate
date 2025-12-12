package autocompletetopics
package autocompletetopics

























}	Debates []DebateSummary `json:"debates"`type AutocompleteResponse struct {// AutocompleteResponse represents the HTTP response}	CreatedAt     string            `json:"createdAt"`	PanelistCount int               `json:"panelistCount"`	Panelists     []PanelistSummary `json:"panelists"`	Topic         string            `json:"topic"`	ID            string            `json:"id"`type DebateSummary struct {// DebateSummary represents a debate summary for autocomplete results}	Slug string `json:"slug"`	Name string `json:"name"`	ID   string `json:"id"`type PanelistSummary struct {// PanelistSummary represents minimal panelist info for autocomplete)	"time"import (