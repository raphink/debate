package validatetopic

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MinTopicLength = 10
	MaxTopicLength = 500
)

var (
	ErrTopicTooShort = errors.New("topic must be at least 10 characters long")
	ErrTopicTooLong  = errors.New("topic must not exceed 500 characters")
	ErrTopicInvalid  = errors.New("topic contains invalid characters or HTML content")
)

// ValidateTopicInput validates the topic input from the request
func ValidateTopicInput(topic string) error {
	// Trim whitespace
	topic = strings.TrimSpace(topic)

	// Check length
	if len(topic) < MinTopicLength {
		return ErrTopicTooShort
	}
	if len(topic) > MaxTopicLength {
		return ErrTopicTooLong
	}

	// Check for HTML tags (basic XSS prevention)
	if strings.Contains(topic, "<") || strings.Contains(topic, ">") {
		return ErrTopicInvalid
	}

	// Check for common HTML entities
	if strings.Contains(topic, "&lt;") || strings.Contains(topic, "&gt;") ||
		strings.Contains(topic, "&#") || strings.Contains(topic, "&amp;") {
		return ErrTopicInvalid
	}

	return nil
}

// SanitizeTopic removes potentially dangerous characters from the topic
func SanitizeTopic(topic string) string {
	// Trim whitespace
	topic = strings.TrimSpace(topic)

	// Remove any HTML tags (just in case they slipped through)
	topic = strings.ReplaceAll(topic, "<", "")
	topic = strings.ReplaceAll(topic, ">", "")

	// Normalize whitespace
	topic = strings.Join(strings.Fields(topic), " ")

	return topic
}

// MapErrorToResponse maps validation errors to ErrorResponse
func MapErrorToResponse(err error) ErrorResponse {
	switch {
	case errors.Is(err, ErrTopicTooShort):
		return ErrorResponse{
			Error:     fmt.Sprintf("Topic must be at least %d characters long", MinTopicLength),
			Code:      ErrInvalidTopicLength,
			Retryable: true,
		}
	case errors.Is(err, ErrTopicTooLong):
		return ErrorResponse{
			Error:     fmt.Sprintf("Topic must not exceed %d characters", MaxTopicLength),
			Code:      ErrInvalidTopicLength,
			Retryable: true,
		}
	case errors.Is(err, ErrTopicInvalid):
		return ErrorResponse{
			Error:     "Topic contains invalid characters or HTML content",
			Code:      ErrInvalidTopicContent,
			Retryable: true,
		}
	default:
		return ErrorResponse{
			Error:     "An error occurred while validating the topic. Please try again.",
			Code:      ErrInternalError,
			Retryable: true,
		}
	}
}
