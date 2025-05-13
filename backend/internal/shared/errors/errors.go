package errors

import "net/http"

// ErrorResponse defines the standard error structure for API responses.
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`    // Optional app-level error code
	Details any    `json:"details,omitempty"` // Optional additional context
}

// New creates a standard ErrorResponse with just a message and HTTP status code.
func New(message string, status int) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Code:    status,
	}
}

// WithDetails creates an ErrorResponse with additional context.
func WithDetails(message string, status int, details any) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Code:    status,
		Details: details,
	}
}

// Predefined common error templates
var (
	ErrBadRequest          = New("Bad request", http.StatusBadRequest)
	ErrUnauthorized        = New("Unauthorized", http.StatusUnauthorized)
	ErrForbidden           = New("Forbidden", http.StatusForbidden)
	ErrNotFound            = New("Not found", http.StatusNotFound)
	ErrInternalServerError = New("Internal server error", http.StatusInternalServerError)
)
