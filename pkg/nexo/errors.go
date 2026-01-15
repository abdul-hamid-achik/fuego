// Package nexo provides a file-system based Go framework for APIs and websites.
package nexo

import (
	"errors"
	"fmt"
	"net/http"
)

// Common errors returned by the framework.
var (
	ErrNotFound         = errors.New("not found")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrInvalidHandler   = errors.New("invalid handler signature")
	ErrScanFailed       = errors.New("failed to scan routes")
	ErrNoAppDir         = errors.New("app directory not found")
)

// HTTPError represents an HTTP error with a status code and message.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// NewHTTPError creates a new HTTPError.
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// NewHTTPErrorWithCause creates a new HTTPError with an underlying cause.
func NewHTTPErrorWithCause(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common HTTP errors for convenience.
var (
	ErrBadRequest          = NewHTTPError(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = NewHTTPError(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = NewHTTPError(http.StatusForbidden, "forbidden")
	ErrHTTPNotFound        = NewHTTPError(http.StatusNotFound, "not found")
	ErrConflict            = NewHTTPError(http.StatusConflict, "conflict")
	ErrInternalServerError = NewHTTPError(http.StatusInternalServerError, "internal server error")
)

// WrapError wraps an error with additional context.
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsHTTPError checks if an error is an HTTPError and returns it.
func IsHTTPError(err error) (*HTTPError, bool) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr, true
	}
	return nil, false
}

// ---------- Error Helper Functions ----------

// BadRequest creates a 400 Bad Request error with a custom message.
func BadRequest(message string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, message)
}

// Unauthorized creates a 401 Unauthorized error with a custom message.
func Unauthorized(message string) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, message)
}

// Forbidden creates a 403 Forbidden error with a custom message.
func Forbidden(message string) *HTTPError {
	return NewHTTPError(http.StatusForbidden, message)
}

// NotFound creates a 404 Not Found error with a custom message.
func NotFound(message string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, message)
}

// Conflict creates a 409 Conflict error with a custom message.
func Conflict(message string) *HTTPError {
	return NewHTTPError(http.StatusConflict, message)
}

// InternalServerError creates a 500 Internal Server Error with a custom message.
func InternalServerError(message string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, message)
}
