# Fuego

A file-system based Go framework for APIs and websites, inspired by Next.js App Router.

The name "Fuego" (Spanish for "fire") contains "GO" naturally embedded (fue**GO**) and reflects its Mexican origin and blazing fast performance.

[![Go Reference](https://pkg.go.dev/badge/github.com/abdul-hamid-achik/fuego.svg)](https://pkg.go.dev/github.com/abdul-hamid-achik/fuego)
[![Go Report Card](https://goreportcard.com/badge/github.com/abdul-hamid-achik/fuego)](https://goreportcard.com/report/github.com/abdul-hamid-achik/fuego)

## Features

- **File system is the router** - No manual route registration. Drop a file, get a route.
- **Zero boilerplate to start** - `main.go` is 5 lines.
- **Explicit over magic** - Handlers are plain Go functions.
- **Fast iteration** - Hot reload in dev, sub-second rebuilds.
- **Scalable conventions** - Works for a 3-route API or a 300-route app.
- **Templ-native** - First-class support for type-safe HTML templating.
- **Proxy layer** - Intercept requests for rewrites, redirects, and early responses.
- **Built-in middleware** - Logger, CORS, rate limiting, auth, and more.

## Quick Start

```bash
# Install Fuego CLI
go install github.com/abdul-hamid-achik/fuego/cmd/fuego@latest

# Create a new project
fuego new myapp

# Start development server
cd myapp
fuego dev
```

Visit http://localhost:3000

## Project Structure

```
myapp/
├── app/
│   ├── proxy.go              # Request interception (optional)
│   ├── middleware.go         # Global middleware
│   ├── layout.templ          # Root layout
│   ├── page.templ            # GET /
│   └── api/
│       ├── middleware.go     # API middleware
│       └── health/
│           └── route.go      # GET /api/health
├── static/
├── main.go
├── fuego.yaml
└── go.mod
```

## File Conventions

| File | Purpose |
|------|---------|
| `route.go` | API endpoint (exports Get, Post, Put, Patch, Delete, etc.) |
| `proxy.go` | Request interception before routing (app root only) |
| `middleware.go` | Middleware for segment and children |
| `page.templ` | UI for a route |
| `layout.templ` | Shared UI wrapper |
| `error.templ` | Error boundary UI |
| `loading.templ` | Loading skeleton |
| `notfound.templ` | Not found UI |

## Dynamic Routes

| Pattern | Example | Matches |
|---------|---------|---------|
| `[param]` | `app/users/[id]/` | `/users/123` |
| `[...param]` | `app/docs/[...slug]/` | `/docs/a/b/c` |
| `[[...param]]` | `app/shop/[[...categories]]/` | `/shop`, `/shop/clothes` |

## Example: API Route

```go
// app/api/users/route.go
package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// GET /api/users
func Get(c *fuego.Context) error {
    return c.JSON(200, map[string]any{
        "users": []string{"alice", "bob"},
    })
}

// POST /api/users
func Post(c *fuego.Context) error {
    var input struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&input); err != nil {
        return fuego.BadRequest("invalid input")
    }
    return c.JSON(201, map[string]string{"created": input.Name})
}
```

## Example: Dynamic Route

```go
// app/api/users/[id]/route.go
package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// GET /api/users/:id
func Get(c *fuego.Context) error {
    id := c.Param("id")
    return c.JSON(200, map[string]any{
        "id": id,
        "name": "User " + id,
    })
}
```

## Proxy (Request Interception)

Intercept requests before routing for rewrites, redirects, and early responses:

```go
// app/proxy.go
package app

import (
    "strings"
    "github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    path := c.Path()
    
    // Redirect old URLs
    if strings.HasPrefix(path, "/api/v1/") {
        newPath := strings.Replace(path, "/api/v1/", "/api/v2/", 1)
        return fuego.Redirect(newPath, 301), nil
    }
    
    // Block unauthorized access
    if strings.HasPrefix(path, "/admin") && !isAdmin(c) {
        return fuego.ResponseJSON(403, `{"error":"forbidden"}`), nil
    }
    
    // Rewrite for A/B testing
    if c.Cookie("experiment") == "variant-b" {
        return fuego.Rewrite("/variant-b" + path), nil
    }
    
    // Continue to normal routing
    return fuego.Continue(), nil
}
```

## Middleware

### File-based Middleware

```go
// app/api/middleware.go
package api

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Middleware() fuego.MiddlewareFunc {
    return func(next fuego.HandlerFunc) fuego.HandlerFunc {
        return func(c *fuego.Context) error {
            // Applied to all routes under /api
            c.SetHeader("X-API-Version", "1.0")
            return next(c)
        }
    }
}
```

### Built-in Middleware

```go
app := fuego.New()

// Logging
app.Use(fuego.Logger())

// Panic recovery
app.Use(fuego.Recover())

// Request ID
app.Use(fuego.RequestID())

// CORS
app.Use(fuego.CORSWithConfig(fuego.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
}))

// Rate limiting
app.Use(fuego.RateLimiter(100, time.Minute))

// Timeout
app.Use(fuego.Timeout(30 * time.Second))

// Basic auth
app.Use(fuego.BasicAuth(map[string]string{
    "admin": "secret",
}))

// Security headers
app.Use(fuego.SecureHeaders())
```

## Context API

```go
func Get(c *fuego.Context) error {
    // URL parameters
    id := c.Param("id")
    idInt, _ := c.ParamInt("id")
    
    // Query parameters
    page := c.Query("page")
    limit := c.QueryDefault("limit", "10")
    
    // Headers
    auth := c.Header("Authorization")
    c.SetHeader("X-Custom", "value")
    
    // Request body
    var body MyStruct
    c.Bind(&body)
    
    // Cookies
    token := c.Cookie("session")
    c.SetCookie("session", "abc123", 3600)
    
    // Context store
    c.Set("user", user)
    user := c.Get("user")
    
    // Response
    return c.JSON(200, data)
    return c.String(200, "Hello")
    return c.HTML(200, "<h1>Hello</h1>")
    return c.Redirect("/login", 302)
    return c.NoContent()
    return c.Blob(200, "application/pdf", pdfBytes)
}
```

## Error Handling

```go
func Get(c *fuego.Context) error {
    if notFound {
        return fuego.NotFound("resource not found")
    }
    
    if unauthorized {
        return fuego.Unauthorized("invalid token")
    }
    
    if badInput {
        return fuego.BadRequest("invalid input")
    }
    
    return fuego.InternalServerError("something went wrong")
}
```

## CLI Commands

```bash
# Create new project
fuego new myapp
fuego new myapp --api-only      # Without templ templates
fuego new myapp --with-proxy    # Include proxy.go example

# Development server with hot reload
fuego dev

# Build for production
fuego build

# List all routes
fuego routes
```

## Configuration

```yaml
# fuego.yaml
port: 3000
host: "0.0.0.0"
app_dir: "app"
static_dir: "static"
static_path: "/static"

dev:
  hot_reload: true
  watch_extensions: [".go", ".templ"]
  exclude_dirs: ["node_modules", ".git"]

middleware:
  logger: true
  recover: true
```

## Development

We use [Task](https://taskfile.dev) for development commands:

```bash
# Build
task build

# Run tests
task test

# Format code
task fmt

# Run all checks
task check

# Install globally
task install
```

## Documentation

- [Quick Start](docs/getting-started/quickstart.md)
- [File-based Routing](docs/routing/file-based.md)
- [Middleware](docs/middleware/overview.md)
- [Proxy](docs/middleware/proxy.md)

## Examples

See the [examples](examples/) directory for complete examples:

- [Basic](examples/basic/) - Simple API with file-based routing

## Acknowledgments

Fuego stands on the shoulders of giants. We're grateful to the authors and maintainers of:

- **[chi](https://github.com/go-chi/chi)** by Peter Kieltyka - The lightweight router that powers Fuego
- **[templ](https://github.com/a-h/templ)** by Adrian Hesketh - Type-safe HTML templating for Go
- **[fsnotify](https://github.com/fsnotify/fsnotify)** - Cross-platform file watching
- **[cobra](https://github.com/spf13/cobra)** by Steve Francia - CLI framework
- **[viper](https://github.com/spf13/viper)** by Steve Francia - Configuration management
- **[Next.js](https://nextjs.org)** by Vercel - The inspiration for our file-system routing

Thank you for making the Go ecosystem amazing!

## Author

**Abdul Hamid Achik** ([@abdulachik](https://x.com/abdulachik))

A Syrian-Mexican software engineer based in Guadalajara, Mexico.

- GitHub: [@abdul-hamid-achik](https://github.com/abdul-hamid-achik)
- Twitter/X: [@abdulachik](https://x.com/abdulachik)

## License

MIT License - see [LICENSE](LICENSE) for details.
