# Proxy

The `proxy.go` convention allows you to intercept and modify requests before they reach your route handlers. Inspired by Next.js 16's middleware system, Fuego's proxy runs at the edge of your request flow.

## Overview

Proxy runs **before route matching**, giving you the power to:

- **Rewrite URLs** - Change the internal path without changing the browser URL (A/B testing, feature flags)
- **Redirect** - Send users to different URLs (301/302/307/308 redirects)
- **Respond Early** - Return responses without hitting your handlers (auth checks, rate limiting)
- **Modify Requests** - Add headers, set cookies, or transform requests

## Quick Start

Create `app/proxy.go`:

```go
package app

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    // Redirect old URLs
    if c.Path() == "/old-page" {
        return fuego.Redirect("/new-page", 301), nil
    }
    
    // Continue to normal routing
    return fuego.Continue(), nil
}
```

## Request Flow

```
Request → Proxy → Global Middleware → Route Middleware → Handler
```

The proxy runs **first**, before any middleware or route handlers.

## ProxyResult Helpers

### Continue

Proceed with normal routing:

```go
return fuego.Continue(), nil
```

### Redirect

Send an HTTP redirect:

```go
// Permanent redirect (301)
return fuego.Redirect("/new-page", 301), nil

// Temporary redirect (302)
return fuego.Redirect("/temp-page", 302), nil

// Preserve method (307/308)
return fuego.Redirect("/api/v2/users", 308), nil
```

### Rewrite

Change the internal path (URL bar stays the same):

```go
// User sees /products but server handles /catalog/products
return fuego.Rewrite("/catalog" + c.Path()), nil
```

### Response

Return a response directly, bypassing routing:

```go
// JSON response
return fuego.ResponseJSON(403, `{"error":"forbidden"}`), nil

// HTML response
return fuego.ResponseHTML(503, "<h1>Maintenance</h1>"), nil

// Custom response
return fuego.Response(429, []byte("Rate limited"), "text/plain"), nil
```

### Adding Headers

Add headers to redirects or responses:

```go
return fuego.Redirect("/login", 302).WithHeader("X-Reason", "session-expired"), nil

return fuego.ResponseJSON(200, `{}`).WithHeaders(map[string]string{
    "X-Custom-1": "value1",
    "X-Custom-2": "value2",
}), nil
```

## ProxyConfig (Optional)

Limit which paths run through the proxy:

```go
package app

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// ProxyConfig is optional - if not defined, proxy runs on all paths
var ProxyConfig = &fuego.ProxyConfig{
    Matcher: []string{
        "/api/:path*",     // All API routes
        "/admin/*",        // Admin routes
    },
}

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    // This only runs for paths matching the patterns above
    return fuego.Continue(), nil
}
```

### Matcher Patterns

- `/api/users` - Exact match
- `/api/:param` - Single dynamic segment
- `/api/:path*` - Zero or more segments (wildcard)
- `/api/:path+` - One or more segments
- `/api/:param?` - Optional segment
- `/(api|admin)` - Regex group

## Use Cases

### Authentication

```go
func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    // Skip auth for public paths
    if c.Path() == "/" || c.Path() == "/login" {
        return fuego.Continue(), nil
    }
    
    token := c.Header("Authorization")
    if token == "" {
        return fuego.Redirect("/login", 302), nil
    }
    
    if !isValidToken(token) {
        return fuego.ResponseJSON(401, `{"error":"invalid token"}`), nil
    }
    
    return fuego.Continue(), nil
}
```

### A/B Testing

```go
func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    // Check experiment cookie
    variant := c.Cookie("experiment")
    
    if variant == "B" && c.Path() == "/pricing" {
        return fuego.Rewrite("/pricing-variant-b"), nil
    }
    
    return fuego.Continue(), nil
}
```

### URL Migration

```go
func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    path := c.Path()
    
    // Redirect old API version
    if strings.HasPrefix(path, "/api/v1/") {
        newPath := strings.Replace(path, "/api/v1/", "/api/v2/", 1)
        return fuego.Redirect(newPath, 308), nil
    }
    
    return fuego.Continue(), nil
}
```

### Rate Limiting

```go
var limiter = NewRateLimiter(100, time.Minute)

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    ip := c.ClientIP()
    
    if !limiter.Allow(ip) {
        return fuego.ResponseJSON(429, `{"error":"too many requests"}`).
            WithHeader("Retry-After", "60"), nil
    }
    
    return fuego.Continue(), nil
}
```

### Maintenance Mode

```go
var maintenanceMode = false

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    if maintenanceMode && !strings.HasPrefix(c.Path(), "/health") {
        return fuego.ResponseHTML(503, `
            <!DOCTYPE html>
            <html>
                <body>
                    <h1>We'll be back soon!</h1>
                </body>
            </html>
        `), nil
    }
    
    return fuego.Continue(), nil
}
```

### Geolocation Routing

```go
func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    country := c.Header("CF-IPCountry") // Cloudflare header
    
    if country == "DE" && c.Path() == "/" {
        return fuego.Rewrite("/de/home"), nil
    }
    
    return fuego.Continue(), nil
}
```

## Proxy vs Middleware

| Feature | Proxy | Middleware |
|---------|-------|------------|
| Runs | Before routing | After routing |
| Can rewrite URLs | Yes | No |
| Can redirect | Yes | Yes |
| Access to route params | No | Yes |
| Per-route configuration | Via matchers | Per-route |
| File location | `app/proxy.go` | `app/**/middleware.go` |

**Use Proxy when you need to:**
- Modify the URL before routing
- Decide routing based on request properties
- Return early responses for all paths

**Use Middleware when you need to:**
- Access route parameters
- Wrap specific route handlers
- Apply logic after routing is determined

## Error Handling

If your proxy function returns an error, Fuego returns a 500 Internal Server Error:

```go
func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    user, err := validateSession(c)
    if err != nil {
        // Returning error = 500 response
        return nil, err
    }
    
    // Better: handle errors explicitly
    if err != nil {
        return fuego.ResponseJSON(500, `{"error":"internal error"}`), nil
    }
    
    return fuego.Continue(), nil
}
```

## CLI Support

View proxy configuration with the routes command:

```bash
fuego routes
```

Output:
```
  Fuego Routes

  PROXY   Proxy enabled
          Matchers: [/api/:path* /admin/*]
          File: app/proxy.go

  GET     /api/users                    app/api/users/route.go
  POST    /api/users                    app/api/users/route.go
  ...
```

## Best Practices

1. **Keep it fast** - Proxy runs on every request, so avoid slow operations
2. **Handle errors gracefully** - Return explicit error responses rather than returning errors
3. **Use matchers** - Only run proxy on paths that need it
4. **Log sparingly** - Too much logging can slow down requests
5. **Test thoroughly** - Proxy bugs affect all matching requests
