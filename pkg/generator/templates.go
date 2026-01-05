package generator

// Template data structures

type routeTemplateData struct {
	Package string
	Methods []methodInfo
	Params  []ParamInfo
	Pattern string
}

type methodInfo struct {
	Method   string // HTTP method (GET, POST, etc.)
	FuncName string // Go function name (Get, Post, etc.)
}

type middlewareTemplateData struct {
	Package string
	Name    string
	Path    string
}

type pageTemplateData struct {
	Package  string
	Title    string
	FilePath string
}

// Route template
// Note: Methods should be title case (Get, Post, Put, Delete) to match scanner expectations
var routeTemplate = `package {{.Package}}

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"
{{range .Methods}}
// {{.FuncName}} handles {{.Method}} /api/{{$.Pattern}}
func {{.FuncName}}(c *nexo.Context) error {
{{- range $.Params}}
	{{.Name}} := c.Param("{{.Name}}")
	_ = {{.Name}} // TODO: use this parameter
{{- end}}
	return c.JSON(200, map[string]any{
{{- range $.Params}}
		"{{.Name}}": {{.Name}},
{{- end}}
		// TODO: Implement {{.FuncName}} handler
	})
}
{{end}}`

// Middleware templates
var middlewareTemplates = map[string]string{
	"blank": `package {{.Package}}

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// Middleware runs before all routes in {{.Path}}
func Middleware(next nexo.HandlerFunc) nexo.HandlerFunc {
	return func(c *nexo.Context) error {
		// TODO: Add middleware logic here
		return next(c)
	}
}
`,
	"auth": `package {{.Package}}

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// Middleware provides authentication for routes in {{.Path}}
func Middleware(next nexo.HandlerFunc) nexo.HandlerFunc {
	return func(c *nexo.Context) error {
		token := c.Header("Authorization")
		if token == "" {
			return c.JSON(401, map[string]string{
				"error":   "unauthorized",
				"message": "Authorization header required",
			})
		}

		// TODO: Validate the token
		// Example: Verify JWT, check database, etc.
		// if !isValidToken(token) {
		//     return c.JSON(403, map[string]string{
		//         "error": "forbidden",
		//         "message": "Invalid or expired token",
		//     })
		// }

		// Optionally set user info in context
		// c.Set("user_id", extractUserID(token))

		return next(c)
	}
}
`,
	"logging": `package {{.Package}}

import (
	"log"
	"time"

	"github.com/abdul-hamid-achik/nexo/pkg/nexo"
)

// Middleware provides request logging for routes in {{.Path}}
func Middleware(next nexo.HandlerFunc) nexo.HandlerFunc {
	return func(c *nexo.Context) error {
		start := time.Now()

		// Log request
		log.Printf("[REQUEST] %s %s", c.Method(), c.Path())

		// Call next handler
		err := next(c)

		// Log response
		duration := time.Since(start)
		if err != nil {
			log.Printf("[RESPONSE] %s %s - ERROR: %v (%s)", c.Method(), c.Path(), err, duration)
		} else {
			log.Printf("[RESPONSE] %s %s - OK (%s)", c.Method(), c.Path(), duration)
		}

		return err
	}
}
`,
	"timing": `package {{.Package}}

import (
	"time"

	"github.com/abdul-hamid-achik/nexo/pkg/nexo"
)

// Middleware adds timing headers for routes in {{.Path}}
func Middleware(next nexo.HandlerFunc) nexo.HandlerFunc {
	return func(c *nexo.Context) error {
		start := time.Now()

		// Call next handler
		err := next(c)

		// Add timing header
		duration := time.Since(start)
		c.SetHeader("X-Response-Time", duration.String())
		c.SetHeader("Server-Timing", "total;dur="+duration.String())

		return err
	}
}
`,
	"cors": `package {{.Package}}

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// CORS configuration
var (
	allowedOrigins = []string{"*"} // TODO: Configure allowed origins
	allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	allowedHeaders = []string{"Content-Type", "Authorization", "X-Requested-With"}
)

// Middleware provides CORS support for routes in {{.Path}}
func Middleware(next nexo.HandlerFunc) nexo.HandlerFunc {
	return func(c *nexo.Context) error {
		origin := c.Header("Origin")

		// Check if origin is allowed
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.SetHeader("Access-Control-Allow-Origin", origin)
			c.SetHeader("Access-Control-Allow-Methods", joinStrings(allowedMethods))
			c.SetHeader("Access-Control-Allow-Headers", joinStrings(allowedHeaders))
			c.SetHeader("Access-Control-Allow-Credentials", "true")
			c.SetHeader("Access-Control-Max-Age", "86400")
		}

		// Handle preflight
		if c.Method() == "OPTIONS" {
			return c.NoContent()
		}

		return next(c)
	}
}

func joinStrings(s []string) string {
	result := ""
	for i, str := range s {
		if i > 0 {
			result += ", "
		}
		result += str
	}
	return result
}
`,
}

// Proxy templates
var proxyTemplates = map[string]string{
	"blank": `package app

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// Proxy runs before route matching.
// Use it for request interception, URL rewriting, or early responses.
func Proxy(c *nexo.Context) (*nexo.ProxyResult, error) {
	// Continue with normal routing
	return nexo.Continue(), nil
}
`,
	"auth-check": `package app

import (
	"strings"

	"github.com/abdul-hamid-achik/nexo/pkg/nexo"
)

// Public paths that don't require authentication
var publicPaths = []string{
	"/",
	"/api/health",
	"/api/public",
	"/login",
	"/register",
}

// Proxy runs before route matching to check authentication.
func Proxy(c *nexo.Context) (*nexo.ProxyResult, error) {
	path := c.Path()

	// Skip auth for public paths
	for _, p := range publicPaths {
		if path == p || strings.HasPrefix(path, p+"/") {
			return nexo.Continue(), nil
		}
	}

	// Skip auth for static files
	if strings.HasPrefix(path, "/static/") {
		return nexo.Continue(), nil
	}

	// Check for auth token
	token := c.Header("Authorization")
	if token == "" {
		return nexo.ResponseJSON(401, map[string]string{
			"error":   "unauthorized",
			"message": "Authorization header required",
		}), nil
	}

	// TODO: Validate token
	// if !isValidToken(token) {
	//     return nexo.ResponseJSON(403, map[string]string{
	//         "error": "forbidden",
	//         "message": "Invalid or expired token",
	//     }), nil
	// }

	// Add header to indicate proxy processed the request
	c.SetHeader("X-Auth-Checked", "true")

	return nexo.Continue(), nil
}
`,
	"rate-limit": `package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/abdul-hamid-achik/nexo/pkg/nexo"
)

// Rate limit configuration
var (
	rateLimitMu sync.Mutex
	requests    = make(map[string][]time.Time)
	maxRequests = 100           // Maximum requests per window
	window      = time.Minute   // Time window
)

// Proxy implements simple IP-based rate limiting.
func Proxy(c *nexo.Context) (*nexo.ProxyResult, error) {
	ip := c.ClientIP()

	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()

	now := time.Now()
	windowStart := now.Add(-window)

	// Clean old requests and count recent ones
	var recent []time.Time
	for _, t := range requests[ip] {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}

	// Check if rate limit exceeded
	if len(recent) >= maxRequests {
		retryAfter := recent[0].Add(window).Sub(now)
		c.SetHeader("Retry-After", retryAfter.String())
		c.SetHeader("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.SetHeader("X-RateLimit-Remaining", "0")
		
		return nexo.ResponseJSON(429, map[string]string{
			"error":   "too_many_requests",
			"message": "Rate limit exceeded. Please try again later.",
		}), nil
	}

	// Record this request
	requests[ip] = append(recent, now)

	// Add rate limit headers
	c.SetHeader("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
	c.SetHeader("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests-len(recent)-1))

	return nexo.Continue(), nil
}
`,
	"maintenance": `package app

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// Set to true to enable maintenance mode
var maintenanceMode = false

// Allowed IPs during maintenance (e.g., admin IPs)
var allowedIPs = []string{
	// "192.168.1.1",
	// "10.0.0.1",
}

// Proxy returns 503 for all requests when in maintenance mode.
func Proxy(c *nexo.Context) (*nexo.ProxyResult, error) {
	if !maintenanceMode {
		return nexo.Continue(), nil
	}

	// Check if IP is allowed during maintenance
	clientIP := c.ClientIP()
	for _, ip := range allowedIPs {
		if ip == clientIP {
			c.SetHeader("X-Maintenance-Bypass", "true")
			return nexo.Continue(), nil
		}
	}

	// Return maintenance response
	c.SetHeader("Retry-After", "3600") // Suggest retry in 1 hour

	return nexo.ResponseJSON(503, map[string]string{
		"error":   "service_unavailable",
		"message": "Service is under maintenance. Please try again later.",
	}), nil
}
`,
	"redirect-www": `package app

import (
	"strings"

	"github.com/abdul-hamid-achik/nexo/pkg/nexo"
)

// Configuration:
// - redirectToWWW = true:  example.com -> www.example.com
// - redirectToWWW = false: www.example.com -> example.com
var redirectToWWW = false

// Proxy handles www/non-www redirects.
func Proxy(c *nexo.Context) (*nexo.ProxyResult, error) {
	host := c.Request.Host

	// Skip for localhost/IP addresses
	if strings.HasPrefix(host, "localhost") || 
	   strings.HasPrefix(host, "127.0.0.1") ||
	   strings.HasPrefix(host, "[::1]") {
		return nexo.Continue(), nil
	}

	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}

	if redirectToWWW {
		// Redirect non-www to www
		if !strings.HasPrefix(host, "www.") {
			newURL := scheme + "://www." + host + c.Request.RequestURI
			return nexo.Redirect(newURL, 301), nil
		}
	} else {
		// Redirect www to non-www
		if strings.HasPrefix(host, "www.") {
			newHost := strings.TrimPrefix(host, "www.")
			newURL := scheme + "://" + newHost + c.Request.RequestURI
			return nexo.Redirect(newURL, 301), nil
		}
	}

	return nexo.Continue(), nil
}
`,
}

// Loader template
var loaderTemplate = `package {{.Package}}

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// {{.DataType}} holds the data for this page.
// Add your data fields here.
type {{.DataType}} struct {
	// TODO: Add your data fields
	// Example:
	// UserName string
	// Items    []Item
}

// Loader loads data for the page.
// This function is automatically called before rendering the page.
func Loader(c *nexo.Context) ({{.DataType}}, error) {
	// TODO: Load your data here
	// Example:
	// - Fetch from database
	// - Call external API
	// - Read from cache
	//
	// Return an error to stop page rendering:
	// if notFound {
	//     return {{.DataType}}{}, nexo.NotFound("Resource not found")
	// }

	return {{.DataType}}{}, nil
}
`

// Page templates
var pageTemplate = `package {{.Package}}

templ Page() {
	@Layout("{{.Title}}") {
		<main style="max-width: 800px; margin: 0 auto; padding: 2rem;">
			<h1>{{.Title}}</h1>
			<p>Edit this page at {{.FilePath}}</p>
		</main>
	}
}
`

var layoutTemplate = `package {{.Package}}

templ Layout(title string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title }</title>
			<style>
				* { box-sizing: border-box; margin: 0; padding: 0; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
					line-height: 1.6;
					color: #333;
				}
			</style>
		</head>
		<body>
			{ children... }
		</body>
	</html>
}
`

// Routes generation templates

var emptyRoutesTemplate = `// Code generated by nexo. DO NOT EDIT.
// This file is automatically regenerated when routes change.

package main

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// RegisterRoutes registers all file-based routes with the app.
// This file is generated because no routes were found in the app directory.
func RegisterRoutes(app *nexo.App) {
	// No routes found. Add route.go files in the app/api directory.
	// Example: app/api/health/route.go with a Get function
}
`

var routesGenTemplate = `// Code generated by nexo. DO NOT EDIT.
// This file is automatically regenerated when routes change.
// Generator schema version: 1

package main

import (
	"github.com/abdul-hamid-achik/nexo/pkg/nexo"
{{range .Imports}}
	{{.Alias}} "{{.Path}}"
{{- end}}
)

// RegisterRoutes registers all file-based routes with the app.
func RegisterRoutes(app *nexo.App) {
{{- if .Proxy}}
	// Register proxy (from {{.Proxy.FilePath}})
	{{- if .Proxy.HasConfig}}
	_ = app.SetProxy({{.Proxy.ImportAlias}}.Proxy, {{.Proxy.ImportAlias}}.ProxyConfig)
	{{- else}}
	_ = app.SetProxy({{.Proxy.ImportAlias}}.Proxy, nil)
	{{- end}}
{{end}}
{{- range .Middlewares}}
	// Middleware for {{.PathPrefix}} (from {{.FilePath}})
	app.RouteTree().AddMiddleware("{{.PathPrefix}}", {{.ImportAlias}}.Middleware)
{{- end}}
{{range .Routes}}
	// {{.Method}} {{.Pattern}} (from {{.FilePath}})
	app.RegisterRoute("{{.Method}}", "{{.Pattern}}", {{.ImportAlias}}.{{.Handler}})
{{- end}}
{{- range .Pages}}
{{- if .HasLoader}}
	// Page: {{.Pattern}} (from {{.FilePath}})
	// Data loaded by: {{.LoaderPackage}}.Loader()
	app.Get("{{.Pattern}}", func(c *nexo.Context) error {
		data, err := {{.ImportAlias}}.Loader(c)
		if err != nil {
			return err
		}
		return nexo.TemplComponent(c, 200, {{.ImportAlias}}.Page(data))
	})
{{- else if .HasParams}}
	// Page: {{.Pattern}} (from {{.FilePath}})
	// Dynamic page with signature: {{.ParamSignature}}
	app.Get("{{.Pattern}}", func(c *nexo.Context) error {
		{{- range .Params}}
		{{- if .FromPath}}
		{{.Name}} := c.Param("{{.Name}}")
		{{- end}}
		{{- end}}
		return nexo.TemplComponent(c, 200, {{.ImportAlias}}.Page({{paramArgs .Params}}))
	})
{{- else}}
	// Page: {{.Pattern}} (from {{.FilePath}})
	app.Get("{{.Pattern}}", func(c *nexo.Context) error {
		return nexo.TemplComponent(c, 200, {{.ImportAlias}}.Page())
	})
{{- end}}
{{- end}}
}
`
