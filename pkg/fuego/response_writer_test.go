package fuego

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriter_DefaultStatus(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Default status should be 200
	if rw.Status() != http.StatusOK {
		t.Errorf("Expected default status 200, got %d", rw.Status())
	}

	// Should not be marked as written yet
	if rw.Written() {
		t.Error("Expected Written() to be false before writing")
	}
}

func TestResponseWriter_CapturesStatus(t *testing.T) {
	testCases := []struct {
		name   string
		status int
	}{
		{"OK", http.StatusOK},
		{"Created", http.StatusCreated},
		{"BadRequest", http.StatusBadRequest},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			rw := newResponseWriter(w)

			rw.WriteHeader(tc.status)

			if rw.Status() != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, rw.Status())
			}

			if !rw.Written() {
				t.Error("Expected Written() to be true after WriteHeader")
			}
		})
	}
}

func TestResponseWriter_CapturesSize(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Write some data
	data := []byte("Hello, World!")
	n, err := rw.Write(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	if rw.Size() != int64(len(data)) {
		t.Errorf("Expected size %d, got %d", len(data), rw.Size())
	}
}

func TestResponseWriter_MultipleWrites(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Write multiple times
	data1 := []byte("Hello, ")
	data2 := []byte("World!")

	_, _ = rw.Write(data1)
	_, _ = rw.Write(data2)

	expectedSize := int64(len(data1) + len(data2))
	if rw.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, rw.Size())
	}
}

func TestResponseWriter_WriteHeaderOnce(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Write header multiple times
	rw.WriteHeader(http.StatusCreated)
	rw.WriteHeader(http.StatusNotFound) // Should be ignored

	// First status should be preserved
	if rw.Status() != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rw.Status())
	}
}

func TestResponseWriter_WriteImpliesOK(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Write without calling WriteHeader first
	_, _ = rw.Write([]byte("test"))

	// Status should be 200 OK
	if rw.Status() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Status())
	}

	if !rw.Written() {
		t.Error("Expected Written() to be true after Write")
	}
}

func TestResponseWriter_Flush(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Flush should not panic
	rw.Flush()
}

func TestResponseWriter_DelegatesToUnderlying(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// Set a header
	rw.Header().Set("X-Test", "value")

	// Write response
	rw.WriteHeader(http.StatusCreated)
	_, _ = rw.Write([]byte("test body"))

	// Check underlying recorder
	if w.Code != http.StatusCreated {
		t.Errorf("Expected underlying status %d, got %d", http.StatusCreated, w.Code)
	}

	if w.Header().Get("X-Test") != "value" {
		t.Error("Expected header to be passed to underlying writer")
	}

	if w.Body.String() != "test body" {
		t.Errorf("Expected body 'test body', got '%s'", w.Body.String())
	}
}
