package listdebates

import (
	"regexp"
	"strings"
)

// minTokenLength is the minimum length for tokens to be considered significant
const minTokenLength = 3

// punctuationRegex matches all punctuation and special characters except spaces and alphanumeric
var punctuationRegex = regexp.MustCompile(`[^a-z0-9\s]+`)

// NormalizeAndTokenize converts text to lowercase, splits compound words (hyphens/slashes),
// removes punctuation, and returns array of significant tokens (≥3 characters)
func NormalizeAndTokenize(text string) []string {
	// Convert to lowercase
	normalized := strings.ToLower(text)

	// Replace hyphens and slashes with spaces to split compound words
	// "climate-change" → "climate change", "AI/ML" → "AI ML"
	normalized = strings.ReplaceAll(normalized, "-", " ")
	normalized = strings.ReplaceAll(normalized, "/", " ")

	// Remove all punctuation
	normalized = punctuationRegex.ReplaceAllString(normalized, " ")

	// Split on whitespace
	words := strings.Fields(normalized)

	// Filter to keep only tokens ≥ minTokenLength
	var tokens []string
	for _, word := range words {
		if len(word) >= minTokenLength {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

// CountMatchingTokens counts how many query tokens appear in topic tokens (bag-of-words).
// Returns 0 if not ALL query tokens are found (failed match), otherwise returns the count
// of matching tokens.
func CountMatchingTokens(queryTokens, topicTokens []string) int {
	if len(queryTokens) == 0 {
		return 0
	}

	// Create a map for efficient lookup of topic tokens
	topicMap := make(map[string]bool)
	for _, token := range topicTokens {
		topicMap[token] = true
	}

	// Count how many query tokens are found in topic
	matchCount := 0
	for _, queryToken := range queryTokens {
		if topicMap[queryToken] {
			matchCount++
		}
	}

	// Return 0 if not all query tokens were found (failed match)
	if matchCount < len(queryTokens) {
		return 0
	}

	// Return the count of matching tokens
	return matchCount
}
