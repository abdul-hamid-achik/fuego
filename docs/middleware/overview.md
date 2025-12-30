# Middleware

Fuego provides a powerful middleware system with built-in middleware and support for custom middleware.

## How Middleware Works

Middleware wraps handlers to add functionality:

```go
func MyMiddleware() fuego.MiddlewareFunc {
    return func(next fuego.HandlerFunc) fuego.HandlerFunc {
        return func(c *fuego.Context) error {
            // Before handler
            fmt.Println("Request started")
            
            err := next(c) // Call the handler
            
            // After handler
            fmt.Println("Request finished")
            
            return err
        }
    }
}
```

## Middleware Execution Order

```
Request
    → App-Level Logger (captures ALL requests)
    → Proxy (if configured)
    → Global Middleware (in order added)
        → Path Middleware (inherited from parent)
            → Route Middleware
                → Handler
            ← Route Middleware
        ← Path Middleware
    ← Global Middleware
Response
```

## Global Middleware

Apply to all routes via `app.Use()`:

```go
app := fuego.New()
app.Use(fuego.Recover())
app.Use(fuego.RequestID())
```

## File-Based Middleware

Create `middleware.go` in any app directory:

```
app/
├── middleware.go           # Applies to all routes
├── api/
│   ├── middleware.go      # Applies to /api/*
│   └── users/
│       └── route.go
```

Example `app/api/middleware.go`:

```go
package api

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Middleware() fuego.MiddlewareFunc {
    return func(next fuego.HandlerFunc) fuego.HandlerFunc {
        return func(c *fuego.Context) error {
            c.SetHeader("X-API-Version", "1.0")
            return next(c)
        }
    }
}
```

## Built-in Middleware

### Request Logger (App-Level)

Fuego includes an app-level request logger that captures **all** requests, including those handled by the proxy layer. The logger is enabled by default.

```go
app := fuego.New() // Logger enabled by default!
```

Output:
```
[12:34:56] GET /api/users 200 in 45ms (1.2KB)
[12:34:57] POST /api/tasks 201 in 123ms (256B)
[12:34:58] GET /v1/users → /api/users 200 in 52ms [rewrite]
[12:34:59] GET /api/admin 403 in 1ms [proxy]
```

#### Configuration

```go
app.SetLogger(fuego.RequestLoggerConfig{
    ShowIP:        true,   // Show client IP
    ShowUserAgent: true,   // Show user agent
    ShowSize:      true,   // Show response size (default: true)
    SkipStatic:    true,   // Don't log static files
    SkipPaths:     []string{"/health", "/metrics"},
    Level:         fuego.LogLevelInfo, // debug, info, warn, error
    TimeUnit:      "auto", // "ms" (default), "us", or "auto"
})
```

Detailed output:
```
[12:34:56] GET /api/users 200 in 45ms (1.2KB) [192.168.1.100] [Mozilla/5.0...]
```

#### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Compact` | `bool` | `true` | Use compact Next.js-style format |
| `ShowTimestamp` | `bool` | `true` | Show `[HH:MM:SS]` timestamp |
| `ShowIP` | `bool` | `false` | Show client IP address |
| `ShowUserAgent` | `bool` | `false` | Show user agent (truncated) |
| `ShowSize` | `bool` | `true` | Show response size |
| `ShowErrors` | `bool` | `true` | Show error details inline |
| `ShowProxyAction` | `bool` | `true` | Show `[proxy]`, `[rewrite]`, `[redirect]` tags |
| `TimeUnit` | `string` | `"ms"` | Time unit: `"ms"`, `"us"`, or `"auto"` |
| `Level` | `LogLevel` | `LogLevelInfo` | Log level filtering |
| `SkipPaths` | `[]string` | `[]` | Paths to skip entirely |
| `SkipStatic` | `bool` | `false` | Skip static file requests |
| `StaticPaths` | `[]string` | `["/static"]` | Paths considered static |
| `DisableColors` | `bool` | `false` | Disable color output |

#### Log Levels

| Level | What's Logged |
|-------|---------------|
| `LogLevelDebug` | Everything + internal details |
| `LogLevelInfo` | All requests (default) |
| `LogLevelWarn` | 4xx + 5xx only |
| `LogLevelError` | 5xx only |
| `LogLevelOff` | Nothing |

#### Environment Variables

- `FUEGO_LOG_LEVEL` - Set log level (`debug`, `info`, `warn`, `error`, `off`)
- `FUEGO_DEV=true` - Automatically sets debug level
- `GO_ENV=production` - Automatically sets warn level

#### Disable/Enable Logger

```go
app.DisableLogger() // Disable logging
app.EnableLogger()  // Re-enable with default config
```

### Middleware-Level Logger (Legacy)

For backward compatibility or fine-grained control, you can use the middleware logger:

```go
app := fuego.New()
app.DisableLogger() // Disable app-level logger
app.Use(fuego.Logger()) // Use middleware logger instead
```

Note: The middleware logger does NOT capture requests handled by the proxy layer.

### Recover

Recovers from panics and returns 500 error:

```go
app.Use(fuego.Recover())
```

### RequestID

Adds a unique request ID header:

```go
app.Use(fuego.RequestID())
```

Header: `X-Request-Id: abc123...`

### CORS

Handle Cross-Origin Resource Sharing:

```go
// Simple - allow all origins
app.Use(fuego.CORS())

// Configured
app.Use(fuego.CORSWithConfig(fuego.CORSConfig{
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
    MaxAge:           86400,
}))
```

### Timeout

Set request timeout:

```go
app.Use(fuego.Timeout(30 * time.Second))
```

### BasicAuth

Simple username/password authentication:

```go
app.Use(fuego.BasicAuth(func(user, pass string) bool {
    return user == "admin" && pass == "secret"
}))
```

### SecureHeaders

Add security-related headers:

```go
app.Use(fuego.SecureHeaders())
```

Adds:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: SAMEORIGIN`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`

### RateLimiter

Limit requests per IP:

```go
app.Use(fuego.RateLimiter(100, time.Minute)) // 100 requests per minute
```

## Custom Middleware

Create your own middleware using the factory pattern:

```go
func AuthMiddleware() fuego.MiddlewareFunc {
    return func(next fuego.HandlerFunc) fuego.HandlerFunc {
        return func(c *fuego.Context) error {
            token := c.Header("Authorization")
            if token == "" {
                return fuego.Unauthorized("missing authorization header")
            }
            
            user, err := validateToken(token)
            if err != nil {
                return fuego.Unauthorized("invalid token")
            }
            
            // Store user in context for handlers
            c.Set("user", user)
            
            return next(c)
        }
    }
}
```

## Middleware vs Proxy

| Feature | Middleware | Proxy |
|---------|------------|-------|
| Runs | After routing | Before routing |
| URL rewriting | No | Yes |
| Access route params | Yes | No |
| Per-route control | Yes | Via matchers |
| Location | `middleware.go` | `app/proxy.go` |

See [Proxy Documentation](./proxy.md) for pre-routing manipulation.

## Best Practices

1. **Order matters** - Add recover first, then other middleware
2. **Keep it focused** - Each middleware should do one thing
3. **Handle errors** - Return proper errors, don't panic
4. **Use context** - Store shared data with `c.Set()`/`c.Get()`
5. **Be careful with state** - Middleware runs concurrently
6. **Use app-level logger** - It captures all requests including proxy actions
