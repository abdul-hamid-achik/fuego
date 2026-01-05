package nexo

import (
	"bufio"
	"net"
	"net/http"
)

// responseWriter wraps http.ResponseWriter to capture status code and response size.
// This is used by the app-level logger to accurately track response information.
type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int64
	wroteHeader bool
}

// newResponseWriter creates a new responseWriter that wraps the given http.ResponseWriter.
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK, // Default status
	}
}

// WriteHeader captures the status code and delegates to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size and delegates to the underlying ResponseWriter.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += int64(n)
	return n, err
}

// Status returns the HTTP status code of the response.
func (rw *responseWriter) Status() int {
	return rw.status
}

// Size returns the number of bytes written to the response body.
func (rw *responseWriter) Size() int64 {
	return rw.size
}

// Written returns true if the response has been written to.
func (rw *responseWriter) Written() bool {
	return rw.wroteHeader
}

// Hijack implements the http.Hijacker interface for WebSocket support.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Flush implements the http.Flusher interface for streaming support.
func (rw *responseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Push implements the http.Pusher interface for HTTP/2 server push.
func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := rw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
