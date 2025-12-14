// Package sanitize provides HTML sanitization and XSS prevention utilities
package sanitize

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// htmlTagPattern matches HTML tags
	htmlTagPattern = regexp.MustCompile(`<[^>]*>`)
	// scriptPattern matches script tags and javascript: protocols
	scriptPattern = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>|javascript:`)
)

// StripHTML removes all HTML tags from input string
func StripHTML(input string) string {
	// Remove script tags first for safety
	cleaned := scriptPattern.ReplaceAllString(input, "")
	// Remove all HTML tags
	cleaned = htmlTagPattern.ReplaceAllString(cleaned, "")
	// Trim whitespace
	return strings.TrimSpace(cleaned)
}

// ValidateTopicText sanitizes and validates topic text input
// Returns sanitized text and validation error if any
func ValidateTopicText(topic string, minLength, maxLength int) (string, error) {
	// Strip HTML and trim
	cleaned := StripHTML(topic)

	// Validate length
	if len(cleaned) < minLength {
		return "", &ValidationError{
			Field:   "topic",
			Message: fmt.Sprintf("Topic must be at least %d characters long", minLength),
			Code:    "INVALID_TOPIC_LENGTH",
		}
	}

	if len(cleaned) > maxLength {
		return "", &ValidationError{
			Field:   "topic",
			Message: fmt.Sprintf("Topic must not exceed %d characters", maxLength),
			Code:    "INVALID_TOPIC_LENGTH",
		}
	}

	return cleaned, nil
}

// SanitizeTextField sanitizes any text field for safe storage and display
func SanitizeTextField(text string) string {
	return StripHTML(text)
}

// ValidatePanelistID checks if panelist ID is alphanumeric only
func ValidatePanelistID(id string) bool {
	if len(id) < 3 || len(id) > 20 {
		return false
	}
	alphanumericPattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphanumericPattern.MatchString(id)
}

// ValidationError represents a sanitization validation error
type ValidationError struct {
	Field   string
	Message string
	Code    string
}

func (e *ValidationError) Error() string {
	return e.Message
}
