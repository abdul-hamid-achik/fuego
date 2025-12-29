# Middleware Example

This example demonstrates Fuego's middleware system, including global middleware, route-level middleware, and middleware inheritance.

## Structure

```
middleware/
├── main.go                     # App entry point with global middleware
├── go.mod
└── app/
    └── api/
        ├── middleware.go       # API-level middleware (applies to all /api/* routes)
        ├── public/
        │   └── route.go        # Public route (no auth required)
        └── protected/
            ├── middleware.go   # Auth middleware (applies to /api/protected/*)
            └── route.go        # Protected route (requires auth)
```

## Middleware Inheritance

Fuego supports hierarchical middleware that follows the file system structure:

1. **Global middleware** - Applied in `main.go` using `app.Use()`
2. **Route-level middleware** - Defined in `middleware.go` files within route directories

When a request hits `/api/protected/resource`:
1. Global middleware runs first (Logger, Recover, RequestID, etc.)
2. API middleware runs next (`app/api/middleware.go`)
3. Protected middleware runs last (`app/api/protected/middleware.go`)
4. Finally, the route handler executes

## Built-in Middleware

Fuego provides several built-in middleware:

```go
app.Use(fuego.Logger())        // Request/response logging
app.Use(fuego.Recover())       // Panic recovery
app.Use(fuego.RequestID())     // X-Request-ID header
app.Use(fuego.SecureHeaders()) // Security headers (XSS, etc.)
app.Use(fuego.CORS())          // CORS with defaults
```

### CORS Configuration

```go
app.Use(fuego.CORSWithConfig(fuego.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

## Custom Middleware

Create custom middleware by implementing the signature:

```go
func Middleware(next fuego.HandlerFunc) fuego.HandlerFunc {
    return func(c *fuego.Context) error {
        // Before handler
        start := time.Now()
        
        // Call next handler
        err := next(c)
        
        // After handler
        c.SetHeader("X-Response-Time", time.Since(start).String())
        
        return err
    }
}
```

## Running the Example

```bash
cd examples/middleware
go run main.go
```

## Testing the Routes

### Public route (no auth required):
```bash
curl http://localhost:3000/api/public
# Response: {"message": "This is a public endpoint", "timestamp": "..."}
```

### Protected route (auth required):
```bash
# Without token - returns 401
curl http://localhost:3000/api/protected
# Response: {"error": "unauthorized", "message": "Authorization header required"}

# With valid token - returns 200
curl -H "Authorization: Bearer valid-token" http://localhost:3000/api/protected
# Response: {"message": "Welcome to the protected area", "user_id": "user-123", ...}
```

## Response Headers

All API routes include:
- `X-API-Version: 1.0` (from API middleware)
- `X-Response-Time: 123µs` (from API middleware)
- `X-Request-ID: uuid` (from global RequestID middleware)
