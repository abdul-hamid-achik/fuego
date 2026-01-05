package nexo

import (
	"errors"
	"testing"
)

func TestHTTPError_Error(t *testing.T) {
	err := NewHTTPError(400, "bad request")
	expected := "400: bad request"

	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestHTTPError_ErrorWithCause(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewHTTPErrorWithCause(500, "internal error", cause)

	if !errors.Is(err, cause) {
		t.Error("Expected error to wrap cause")
	}

	if !errors.Is(err.Unwrap(), cause) {
		t.Error("Expected Unwrap to return cause")
	}

	errStr := err.Error()
	if errStr != "500: internal error: underlying error" {
		t.Errorf("Unexpected error string: %s", errStr)
	}
}

func TestIsHTTPError(t *testing.T) {
	httpErr := NewHTTPError(404, "not found")
	regularErr := errors.New("regular error")

	// Test with HTTPError
	result, ok := IsHTTPError(httpErr)
	if !ok {
		t.Error("Expected IsHTTPError to return true for HTTPError")
	}
	if result.Code != 404 {
		t.Errorf("Expected code 404, got %d", result.Code)
	}

	// Test with regular error
	_, ok = IsHTTPError(regularErr)
	if ok {
		t.Error("Expected IsHTTPError to return false for regular error")
	}

	// Test with wrapped HTTPError
	wrapped := WrapError(httpErr, "context")
	result, ok = IsHTTPError(wrapped)
	if !ok {
		t.Error("Expected IsHTTPError to unwrap and find HTTPError")
	}
	if result.Code != 404 {
		t.Errorf("Expected code 404, got %d", result.Code)
	}
}

func TestWrapError(t *testing.T) {
	// Test nil error
	if WrapError(nil, "context") != nil {
		t.Error("Expected nil for nil error")
	}

	// Test wrapping
	original := errors.New("original")
	wrapped := WrapError(original, "context")

	if !errors.Is(wrapped, original) {
		t.Error("Wrapped error should contain original")
	}

	expected := "context: original"
	if wrapped.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, wrapped.Error())
	}
}

func TestCommonHTTPErrors(t *testing.T) {
	tests := []struct {
		name    string
		err     *HTTPError
		code    int
		message string
	}{
		{"BadRequest", ErrBadRequest, 400, "bad request"},
		{"Unauthorized", ErrUnauthorized, 401, "unauthorized"},
		{"Forbidden", ErrForbidden, 403, "forbidden"},
		{"NotFound", ErrHTTPNotFound, 404, "not found"},
		{"Conflict", ErrConflict, 409, "conflict"},
		{"InternalServerError", ErrInternalServerError, 500, "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, tt.err.Code)
			}
			if tt.err.Message != tt.message {
				t.Errorf("Expected message %q, got %q", tt.message, tt.err.Message)
			}
		})
	}
}
