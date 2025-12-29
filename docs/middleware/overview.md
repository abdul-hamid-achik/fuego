# Middleware

Fuego provides a powerful middleware system with built-in middleware and support for custom middleware.

## How Middleware Works

Middleware wraps handlers to add functionality:

```go
func MyMiddleware(next fuego.HandlerFunc) fuego.HandlerFunc {
    return func(c *fuego.Context) error {
        // Before handler
        fmt.Println("Request started")
        
        err := next(c) // Call the handler
        
        // After handler
        fmt.Println("Request finished")
        
        return err
    }
}
```

## Middleware Execution Order

```
Request
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
app.Use(fuego.Logger())
app.Use(fuego.Recover())
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

### Logger

Logs requests with method, path, status, and duration:

```go
app.Use(fuego.Logger())

// With config
app.Use(fuego.LoggerWithConfig(fuego.LoggerConfig{
    SkipPaths: []string{"/health", "/metrics"},
}))
```

Output:
```
2024/01/15 10:30:45 200 GET     /api/users 1.234ms <nil>
```

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
app.Use(fuego.BasicAuth(map[string]string{
    "admin": "secret",
    "user": "password",
}))
```

### SecureHeaders

Add security-related headers:

```go
app.Use(fuego.SecureHeaders())
```

Adds:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`

### RateLimiter

Limit requests per IP:

```go
app.Use(fuego.RateLimiter(100, time.Minute)) // 100 requests per minute
```

## Custom Middleware

Create your own middleware:

```go
func AuthMiddleware(next fuego.HandlerFunc) fuego.HandlerFunc {
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

1. **Order matters** - Add logging/recover first, auth after
2. **Keep it focused** - Each middleware should do one thing
3. **Handle errors** - Return proper errors, don't panic
4. **Use context** - Store shared data with `c.Set()`/`c.Get()`
5. **Be careful with state** - Middleware runs concurrently
