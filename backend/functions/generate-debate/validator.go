package main

import (
	"errors"
	"strings"
)

// ValidateDebateRequest validates the debate generation request
func ValidateDebateRequest(req *DebateRequest) error {
	if req == nil {
		return errors.New("request body is required")
	}

	// Validate topic
	topic := strings.TrimSpace(req.Topic)
	if len(topic) < 10 {
		return errors.New("topic must be at least 10 characters")
	}
	if len(topic) > 500 {
		return errors.New("topic must not exceed 500 characters")
	}

	// Validate panelists
	if len(req.SelectedPanelists) < 2 {
		return errors.New("at least 2 panelists are required")
	}
	if len(req.SelectedPanelists) > 5 {
		return errors.New("maximum 5 panelists allowed")
	}

	// Validate each panelist
	for _, panelist := range req.SelectedPanelists {
		if panelist.ID == "" || panelist.Name == "" {
			return errors.New("all panelists must have id and name")
		}
	}

	return nil
}
