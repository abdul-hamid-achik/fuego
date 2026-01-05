package nexo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// LogLevel represents the logging level.
type LogLevel int

const (
	// LogLevelDebug logs everything including internal details.
	LogLevelDebug LogLevel = iota
	// LogLevelInfo logs all requests (default).
	LogLevelInfo
	// LogLevelWarn logs only 4xx and 5xx responses.
	LogLevelWarn
	// LogLevelError logs only 5xx responses.
	LogLevelError
	// LogLevelOff disables logging entirely.
	LogLevelOff
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	case LogLevelOff:
		return "off"
	default:
		return "info"
	}
}

// ParseLogLevel parses a string into a LogLevel.
func ParseLogLevel(s string) LogLevel {
	switch strings.ToLower(s) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "off", "none", "disabled":
		return LogLevelOff
	default:
		return LogLevelInfo
	}
}

// RequestLoggerConfig holds configuration for the request logger.
type RequestLoggerConfig struct {
	// Display Options
	Compact         bool // Use compact Next.js-style format (default: true)
	ShowTimestamp   bool // Show [HH:MM:SS] timestamp (default: true)
	ShowIP          bool // Show client IP (default: false)
	ShowUserAgent   bool // Show user agent (default: false)
	ShowErrors      bool // Show error details inline (default: true)
	ShowProxyAction bool // Show proxy action tags (default: true)
	ShowSize        bool // Show response size (default: true)

	// Formatting
	TimeUnit        string // "ms" (default), "us", or "auto"
	TimestampFormat string // "15:04:05" (default)

	// Filtering
	Level       LogLevel // Log level (default: LogLevelInfo)
	SkipPaths   []string // Paths to skip entirely
	SkipStatic  bool     // Skip static files (default: false)
	StaticPaths []string // Paths considered static

	// Colors
	DisableColors bool // Force disable colors (default: false, auto-detected)

	// MaxErrorLength is the maximum length for error messages in logs.
	// Messages longer than this are truncated. Default: 100.
	MaxErrorLength int
}

// DefaultRequestLoggerConfig returns sensible defaults for the request logger.
func DefaultRequestLoggerConfig() RequestLoggerConfig {
	level := LogLevelInfo

	// Check environment variable for log level
	if envLevel := os.Getenv("NEXO_LOG_LEVEL"); envLevel != "" {
		level = ParseLogLevel(envLevel)
	} else {
		// Auto-detect dev vs prod mode
		if os.Getenv("NEXO_DEV") == "true" || os.Getenv("GO_ENV") == "development" {
			level = LogLevelDebug
		} else if os.Getenv("GO_ENV") == "production" {
			level = LogLevelWarn
		}
	}

	return RequestLoggerConfig{
		Compact:         true,
		ShowTimestamp:   true,
		ShowErrors:      true,
		ShowProxyAction: true,
		ShowSize:        true,
		TimeUnit:        "ms",
		TimestampFormat: "15:04:05",
		Level:           level,
		StaticPaths:     []string{"/static", "/assets", "/public", "/_next"},
		MaxErrorLength:  100,
	}
}

// RequestLogger handles request logging with configurable output.
type RequestLogger struct {
	config RequestLoggerConfig

	// Color functions
	methodColors map[string]func(a ...interface{}) string
	statusColors map[int]func(a ...interface{}) string
	dim          func(a ...interface{}) string
	cyan         func(a ...interface{}) string
	yellow       func(a ...interface{}) string
}

// NewRequestLogger creates a new request logger with the given configuration.
func NewRequestLogger(config RequestLoggerConfig) *RequestLogger {
	// Auto-detect TTY for color support
	if !config.DisableColors {
		config.DisableColors = !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())
	}

	if config.DisableColors {
		color.NoColor = true
	}

	rl := &RequestLogger{
		config:       config,
		methodColors: make(map[string]func(a ...interface{}) string),
	}

	// Set up method colors
	rl.methodColors[http.MethodGet] = color.New(color.FgBlue).SprintFunc()
	rl.methodColors[http.MethodPost] = color.New(color.FgGreen).SprintFunc()
	rl.methodColors[http.MethodPut] = color.New(color.FgYellow).SprintFunc()
	rl.methodColors[http.MethodPatch] = color.New(color.FgMagenta).SprintFunc()
	rl.methodColors[http.MethodDelete] = color.New(color.FgRed).SprintFunc()
	rl.methodColors[http.MethodHead] = color.New(color.FgCyan).SprintFunc()
	rl.methodColors[http.MethodOptions] = color.New(color.FgWhite).SprintFunc()

	// Set up status color ranges
	rl.statusColors = make(map[int]func(a ...interface{}) string)

	// Helper colors
	rl.dim = color.New(color.Faint).SprintFunc()
	rl.cyan = color.New(color.FgCyan).SprintFunc()
	rl.yellow = color.New(color.FgYellow).SprintFunc()

	return rl
}

// getMethodColor returns the color function for a given HTTP method.
func (rl *RequestLogger) getMethodColor(method string) func(a ...interface{}) string {
	if colorFunc, ok := rl.methodColors[method]; ok {
		return colorFunc
	}
	return color.New(color.FgWhite).SprintFunc()
}

// getStatusColor returns the color function for a given status code.
func (rl *RequestLogger) getStatusColor(status int) func(a ...interface{}) string {
	switch {
	case status >= 500:
		return color.New(color.FgRed).SprintFunc()
	case status >= 400:
		return color.New(color.FgYellow).SprintFunc()
	case status >= 300:
		return color.New(color.FgCyan).SprintFunc()
	default:
		return color.New(color.FgGreen).SprintFunc()
	}
}

// ShouldLog determines if a request should be logged based on configuration.
func (rl *RequestLogger) ShouldLog(path string, status int) bool {
	// Check level
	switch rl.config.Level {
	case LogLevelOff:
		return false
	case LogLevelError:
		if status < 500 {
			return false
		}
	case LogLevelWarn:
		if status < 400 {
			return false
		}
	case LogLevelDebug, LogLevelInfo:
		// Log everything
	}

	// Check skip paths
	for _, skipPath := range rl.config.SkipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath+"/") {
			return false
		}
	}

	// Check static files
	if rl.config.SkipStatic && rl.isStaticPath(path) {
		return false
	}

	return true
}

// isStaticPath checks if a path is a static file path.
func (rl *RequestLogger) isStaticPath(path string) bool {
	// Check configured static paths
	for _, prefix := range rl.config.StaticPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Check file extensions
	ext := strings.ToLower(filepath.Ext(path))
	staticExts := []string{
		".css", ".js", ".map",
		".jpg", ".jpeg", ".png", ".gif", ".svg", ".ico", ".webp",
		".woff", ".woff2", ".ttf", ".eot", ".otf",
		".mp4", ".webm", ".mp3", ".wav",
		".pdf", ".zip", ".tar", ".gz",
	}
	for _, staticExt := range staticExts {
		if ext == staticExt {
			return true
		}
	}

	return false
}

// formatLatency formats the duration based on configuration.
func (rl *RequestLogger) formatLatency(d time.Duration) string {
	switch rl.config.TimeUnit {
	case "us":
		return fmt.Sprintf("%dµs", d.Microseconds())
	case "auto":
		if d < time.Millisecond {
			return fmt.Sprintf("%dµs", d.Microseconds())
		} else if d < time.Second {
			return fmt.Sprintf("%dms", d.Milliseconds())
		} else {
			return fmt.Sprintf("%.2fs", d.Seconds())
		}
	case "ms":
		fallthrough
	default:
		ms := d.Milliseconds()
		if ms == 0 && d > 0 {
			return "<1ms"
		}
		return fmt.Sprintf("%dms", ms)
	}
}

// formatSize formats the response size in a human-readable format.
func (rl *RequestLogger) formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	} else {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	}
}

// looksLikeBody detects if content looks like HTML/JSON body content
// that shouldn't appear in request logs.
func looksLikeBody(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return false
	}

	// HTML indicators
	if strings.HasPrefix(s, "<!") ||
		strings.HasPrefix(s, "<html") ||
		strings.HasPrefix(s, "<HTML") ||
		strings.HasPrefix(s, "<head") ||
		strings.HasPrefix(s, "<HEAD") ||
		strings.HasPrefix(s, "<body") ||
		strings.HasPrefix(s, "<BODY") {
		return true
	}

	// Large JSON objects/arrays (likely response bodies)
	if (strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[")) && len(s) > 200 {
		return true
	}

	return false
}

// formatError extracts a clean, sanitized error message suitable for logging.
// It never returns body content (HTML, large JSON) - only concise error messages.
func (rl *RequestLogger) formatError(err error) string {
	if err == nil {
		return ""
	}

	var msg string

	// Check for HTTPError - use only the semantic message
	if httpErr, ok := IsHTTPError(err); ok {
		msg = httpErr.Message
	} else {
		msg = err.Error()
	}

	// Skip if it looks like body content
	if looksLikeBody(msg) {
		return ""
	}

	// Truncate if too long
	maxLen := rl.config.MaxErrorLength
	if maxLen <= 0 {
		maxLen = 100
	}
	if len(msg) > maxLen {
		return msg[:maxLen-3] + "..."
	}

	return msg
}

// ProxyAction represents the action taken by the proxy.
type ProxyAction struct {
	Type   string // "continue", "rewrite", "redirect", "response"
	Target string // URL for rewrite/redirect
}

// Log logs a request with the given parameters.
func (rl *RequestLogger) Log(r *http.Request, status int, size int64, latency time.Duration, proxyAction *ProxyAction, err error) {
	path := r.URL.Path

	// Check if we should log this request
	if !rl.ShouldLog(path, status) {
		return
	}

	// Build the log message
	var msg strings.Builder

	// Timestamp
	if rl.config.ShowTimestamp {
		timestamp := time.Now().Format(rl.config.TimestampFormat)
		msg.WriteString(rl.dim(fmt.Sprintf("[%s] ", timestamp)))
	}

	// Method (color-coded)
	methodColor := rl.getMethodColor(r.Method)
	msg.WriteString(methodColor(r.Method))
	msg.WriteString(" ")

	// Path (with optional rewrite indicator)
	if proxyAction != nil && proxyAction.Type == "rewrite" && proxyAction.Target != "" {
		// Show original path → rewritten path
		msg.WriteString(path)
		msg.WriteString(" ")
		msg.WriteString(rl.dim("→"))
		msg.WriteString(" ")
		msg.WriteString(proxyAction.Target)
	} else {
		msg.WriteString(path)
	}
	msg.WriteString(" ")

	// Status (color-coded)
	statusColor := rl.getStatusColor(status)
	msg.WriteString(statusColor(fmt.Sprintf("%d", status)))
	msg.WriteString(" ")

	// Latency
	msg.WriteString(rl.dim("in "))
	msg.WriteString(rl.formatLatency(latency))

	// Size (optional)
	if rl.config.ShowSize && size > 0 {
		msg.WriteString(" ")
		msg.WriteString(rl.dim(fmt.Sprintf("(%s)", rl.formatSize(size))))
	}

	// Proxy action tag (optional)
	if rl.config.ShowProxyAction && proxyAction != nil {
		switch proxyAction.Type {
		case "redirect":
			msg.WriteString(" ")
			msg.WriteString(rl.cyan(fmt.Sprintf("[redirect → %s]", proxyAction.Target)))
		case "response":
			msg.WriteString(" ")
			msg.WriteString(rl.cyan("[proxy]"))
		case "rewrite":
			msg.WriteString(" ")
			msg.WriteString(rl.cyan("[rewrite]"))
		}
	}

	// Client IP (optional)
	if rl.config.ShowIP {
		ip := getClientIP(r)
		msg.WriteString(" ")
		msg.WriteString(rl.dim(fmt.Sprintf("[%s]", ip)))
	}

	// User agent (optional)
	if rl.config.ShowUserAgent {
		ua := r.UserAgent()
		if len(ua) > 50 {
			ua = ua[:47] + "..."
		}
		msg.WriteString(" ")
		msg.WriteString(rl.dim(fmt.Sprintf("[%s]", ua)))
	}

	// Error (optional)
	if rl.config.ShowErrors && err != nil {
		errMsg := rl.formatError(err)
		if errMsg != "" {
			msg.WriteString(" ")
			msg.WriteString(rl.yellow(fmt.Sprintf("[%s]", errMsg)))
		}
	}

	// Print the log message
	log.Println(msg.String())
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	// Check X-Real-IP header
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if colonIdx := strings.LastIndex(ip, ":"); colonIdx != -1 {
		ip = ip[:colonIdx]
	}
	return ip
}
