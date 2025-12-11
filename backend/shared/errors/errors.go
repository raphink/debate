// Package errors provides user-friendly error types and utilities
package errors

import "fmt"

// AppError represents an application error with user-friendly messaging
type AppError struct {
	Message   string `json:"error"`
	Code      string `json:"code"`
	Retryable bool   `json:"retryable"`
	Internal  error  `json:"-"` // Internal error not exposed to client
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(message, code string, retryable bool, internal error) *AppError {
	return &AppError{
		Message:   message,
		Code:      code,
		Retryable: retryable,
		Internal:  internal,
	}
}

// Common error constructors

// BadRequest creates a 400 Bad Request error
func BadRequest(message, code string) *AppError {
	return &AppError{
		Message:   message,
		Code:      code,
		Retryable: true,
		Internal:  nil,
	}
}

// InternalError creates a 500 Internal Server Error
func InternalError(message string, err error) *AppError {
	return &AppError{
		Message:   "An internal error occurred. Please try again later.",
		Code:      "INTERNAL_ERROR",
		Retryable: true,
		Internal:  err,
	}
}

// APIError creates an error for external API failures
func APIError(service string, err error) *AppError {
	return &AppError{
		Message:   fmt.Sprintf("Failed to communicate with %s. Please try again.", service),
		Code:      "API_ERROR",
		Retryable: true,
		Internal:  err,
	}
}

// ValidationError creates a validation error
func ValidationError(field, message, code string) *AppError {
	return &AppError{
		Message:   message,
		Code:      code,
		Retryable: true,
		Internal:  fmt.Errorf("validation failed for field %s", field),
	}
}

// RateLimitError creates a rate limit exceeded error
func RateLimitError() *AppError {
	return &AppError{
		Message:   "Rate limit exceeded. Please try again in a few moments.",
		Code:      "RATE_LIMIT_EXCEEDED",
		Retryable: true,
		Internal:  nil,
	}
}
