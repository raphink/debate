package listdebates

import (
	"reflect"
	"testing"
)

func TestNormalizeAndTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "lowercase conversion",
			input:    "Should Animals Have Rights",
			expected: []string{"should", "animals", "have", "rights"},
		},
		{
			name:     "hyphen replacement",
			input:    "climate-change effects",
			expected: []string{"climate", "change", "effects"},
		},
		{
			name:     "slash replacement",
			input:    "AI/ML risks and benefits",
			expected: []string{"risks", "and", "benefits"},
		},
		{
			name:     "punctuation removal",
			input:    "Is AI dangerous? Yes, it might be!",
			expected: []string{"dangerous", "yes", "might"},
		},
		{
			name:     "token length filtering - keeps ≥3 chars",
			input:    "do we ban AI in US",
			expected: []string{"ban"},
		},
		{
			name:     "all short words filtered out",
			input:    "is it ok",
			expected: []string{},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "special characters only",
			input:    "!@#$%^&*()",
			expected: []string{},
		},
		{
			name:     "mixed punctuation and words",
			input:    "What's the climate-change policy?",
			expected: []string{"what", "the", "climate", "change", "policy"},
		},
		{
			name:     "compound words with multiple separators",
			input:    "self-driving/autonomous cars",
			expected: []string{"self", "driving", "autonomous", "cars"},
		},
		{
			name:     "numbers preserved",
			input:    "AI in 2024 and beyond",
			expected: []string{"2024", "and", "beyond"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAndTokenize(tt.input)
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("NormalizeAndTokenize(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountMatchingTokens(t *testing.T) {
	tests := []struct {
		name        string
		queryTokens []string
		topicTokens []string
		expected    int
	}{
		{
			name:        "all query tokens found",
			queryTokens: []string{"animals", "rights"},
			topicTokens: []string{"should", "animals", "have", "rights"},
			expected:    2,
		},
		{
			name:        "bag-of-words - order independent",
			queryTokens: []string{"climate", "change"},
			topicTokens: []string{"change", "the", "climate", "policy"},
			expected:    2,
		},
		{
			name:        "partial match - returns count",
			queryTokens: []string{"animals", "extinction"},
			topicTokens: []string{"should", "animals", "have", "rights"},
			expected:    1,
		},
		{
			name:        "single token match",
			queryTokens: []string{"climate"},
			topicTokens: []string{"climate", "change", "effects"},
			expected:    1,
		},
		{
			name:        "no matches",
			queryTokens: []string{"war", "peace"},
			topicTokens: []string{"climate", "change"},
			expected:    0,
		},
		{
			name:        "empty query tokens",
			queryTokens: []string{},
			topicTokens: []string{"climate", "change"},
			expected:    0,
		},
		{
			name:        "empty topic tokens",
			queryTokens: []string{"climate"},
			topicTokens: []string{},
			expected:    0,
		},
		{
			name:        "both empty",
			queryTokens: []string{},
			topicTokens: []string{},
			expected:    0,
		},
		{
			name:        "duplicate tokens in topic",
			queryTokens: []string{"climate"},
			topicTokens: []string{"climate", "climate", "change"},
			expected:    1,
		},
		{
			name:        "partial match - multiple tokens",
			queryTokens: []string{"climate", "change", "policy"},
			topicTokens: []string{"climate", "change"},
			expected:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountMatchingTokens(tt.queryTokens, tt.topicTokens)
			if result != tt.expected {
				t.Errorf("CountMatchingTokens(%v, %v) = %d, want %d",
					tt.queryTokens, tt.topicTokens, result, tt.expected)
			}
		})
	}
}
