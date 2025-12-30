package app

import (
	"strings"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Proxy demonstrates various proxy layer features:
// - URL rewriting (legacy path migration)
// - Access control (blocking admin paths)
// - Request header manipulation
// - A/B testing redirects
func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
	path := c.Path()

	// 1. Legacy URL Rewriting: /v1/* -> /api/*
	if strings.HasPrefix(path, "/v1/") {
		newPath := strings.Replace(path, "/v1/", "/api/", 1)
		return fuego.Rewrite(newPath), nil
	}

	// 2. Access Control: Block /api/admin without auth
	if strings.HasPrefix(path, "/api/admin") {
		authHeader := c.Header("Authorization")
		if authHeader == "" {
			return fuego.ResponseJSON(401, `{"error":"unauthorized","message":"Admin access requires authentication"}`), nil
		}
		// In real app, validate the token here
		if authHeader != "Bearer admin-token" {
			return fuego.ResponseJSON(403, `{"error":"forbidden","message":"Invalid admin credentials"}`), nil
		}
	}

	// 3. Maintenance Mode (toggle with env var or config)
	// Uncomment to enable:
	// if os.Getenv("MAINTENANCE_MODE") == "true" && !strings.HasPrefix(path, "/health") {
	//     return fuego.ResponseHTML(503, "<h1>Under Maintenance</h1><p>We'll be back soon!</p>"), nil
	// }

	// 4. Add tracking headers for all requests
	return fuego.Continue().
		WithHeader("X-Proxy-Version", "1.0").
		WithHeader("X-Request-Path", path), nil
}

// ProxyMatchers returns patterns for selective proxy execution.
// Only these paths will trigger the proxy.
func ProxyMatchers() []string {
	return []string{
		"/v1/*",        // Legacy API paths
		"/api/admin/*", // Admin paths (protected)
		"/api/*",       // API paths (for header injection)
	}
}
