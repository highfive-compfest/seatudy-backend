package apierror

import (
	"errors"
	"net/http"
)

// ApiErrorBuilder is used to build an ApiError step by step.
type ApiErrorBuilder struct {
	httpStatus int
	message    string
	payload    any
}

// NewApiErrorBuilder initializes a new ApiErrorBuilder with default values.
func NewApiErrorBuilder() *ApiErrorBuilder {
	return &ApiErrorBuilder{
		httpStatus: http.StatusInternalServerError, // default status
	}
}

// WithHttpStatus sets the HTTP status code for the error.
func (b *ApiErrorBuilder) WithHttpStatus(status int) *ApiErrorBuilder {
	b.httpStatus = status
	return b
}

// WithMessage sets the message for the error.
func (b *ApiErrorBuilder) WithMessage(message string) *ApiErrorBuilder {
	b.message = message
	return b
}

// WithPayload sets the payload for the error.
func (b *ApiErrorBuilder) WithPayload(payload any) *ApiErrorBuilder {
	b.payload = payload
	return b
}

// Build constructs the ApiError with the configured parameters.
func (b *ApiErrorBuilder) Build() *ApiError {
	return &ApiError{
		HttpStatus: b.httpStatus,
		Message:    b.message,
		Payload:    b.payload,
	}
}

// ApiError represents a structured API error.
type ApiError struct {
	HttpStatus int
	Message    string
	Payload    any
}

// Error returns the error message.
func (e *ApiError) Error() string {
	return e.Message
}

// GetHttpStatus retrieves the HTTP status code from the error.
func GetHttpStatus(err error) int {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr.HttpStatus
	}
	return http.StatusInternalServerError
}

// GetPayload retrieves the payload from the error.
func GetPayload(err error) any {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr.Payload
	}
	return nil
}
