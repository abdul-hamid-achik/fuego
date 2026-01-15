# Nexo

**File-based routing for Go. Fast to write. Faster to run.**

Your file structure *is* your router. Build APIs and full-stack web apps with conventions that feel natural — if you've used modern meta-frameworks, you'll be productive in minutes.

[![Go Reference](https://pkg.go.dev/badge/github.com/abdul-hamid-achik/nexo.svg)](https://pkg.go.dev/github.com/abdul-hamid-achik/nexo)
[![Go Report Card](https://goreportcard.com/badge/github.com/abdul-hamid-achik/nexo)](https://goreportcard.com/report/github.com/abdul-hamid-achik/nexo)

> **fue**GO — Spanish for "fire", with Go built right in. Born in Mexico, built for speed.

## Why Nexo?

Traditional Go routing requires manual registration:

```go
// The old way
r.HandleFunc("/api/users", usersHandler)
r.HandleFunc("/api/users/{id}", userHandler)
r.HandleFunc("/api/posts", postsHandler)
// ...repeat for every route
```

With Nexo, your file structure is your router:

```
app/api/users/route.go      → GET/POST /api/users
app/api/users/_id/route.go  → GET/PUT/DELETE /api/users/:id
app/api/posts/route.go      → GET/POST /api/posts
```

**No registration. No configuration. Just files.**

## Features

- **File system routing** — Your directory structure defines your routes. No manual registration.
- **Zero-config start** — A working API in 5 lines of code.
- **Convention over configuration** — Sensible defaults, full control when you need it.
- **Type-safe templates** — First-class [templ](https://templ.guide) support with compile-time HTML validation.
- **HTMX-ready** — Build interactive UIs without client-side JavaScript frameworks.
- **Standalone Tailwind** — Built-in Tailwind CSS v4 binary. No Node.js required.
- **Request interception** — Proxy layer for auth checks, rewrites, and early responses.
- **Hot reload** — Sub-second rebuilds during development.
- **Production ready** — Single binary deployment, minimal dependencies.

## Installation

### Using Homebrew (macOS/Linux)

```bash
brew install abdul-hamid-achik/tap/nexo-cli
```

### Using Go

```bash
go install github.com/abdul-hamid-achik/nexo/cmd/nexo@latest
```

## Quick Start

```bash
# Create a new project
nexo new myapp

# Start development server
cd myapp
nexo dev
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
├── nexo.yaml
└── go.mod
```

## Familiar Conventions

Nexo uses file-based routing patterns found in modern web frameworks:

| Pattern | File | Route |
|---------|------|-------|
| Static | `app/api/users/route.go` | `/api/users` |
| Dynamic | `app/api/users/_id/route.go` | `/api/users/:id` |
| Catch-all | `app/docs/__slug/route.go` | `/docs/*` |
| Optional catch-all | `app/shop/___categories/route.go` | `/shop`, `/shop/*` |
| Middleware | `app/api/middleware.go` | Applies to `/api/*` |
| Pages | `app/dashboard/page.templ` | `/dashboard` |
| Layouts | `app/layout.templ` | Wraps child pages |

If you've used Next.js, Nuxt, SvelteKit, or similar frameworks, these patterns will feel familiar.

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

## Example: API Route

```go
// app/api/users/route.go
package users

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// GET /api/users
func Get(c *nexo.Context) error {
    return c.JSON(200, map[string]any{
        "users": []string{"alice", "bob"},
    })
}

// POST /api/users
func Post(c *nexo.Context) error {
    var input struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&input); err != nil {
        return nexo.BadRequest("invalid input")
    }
    return c.JSON(201, map[string]string{"created": input.Name})
}
```

## Example: Dynamic Route

```go
// app/api/users/_id/route.go
package id

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

// GET /api/users/:id
func Get(c *nexo.Context) error {
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
    "github.com/abdul-hamid-achik/nexo/pkg/nexo"
)

func Proxy(c *nexo.Context) (*nexo.ProxyResult, error) {
    path := c.Path()
    
    // Redirect old URLs
    if strings.HasPrefix(path, "/api/v1/") {
        newPath := strings.Replace(path, "/api/v1/", "/api/v2/", 1)
        return nexo.Redirect(newPath, 301), nil
    }
    
    // Block unauthorized access
    if strings.HasPrefix(path, "/admin") && !isAdmin(c) {
        return nexo.ResponseJSON(403, `{"error":"forbidden"}`), nil
    }
    
    // Rewrite for A/B testing
    if c.Cookie("experiment") == "variant-b" {
        return nexo.Rewrite("/variant-b" + path), nil
    }
    
    // Continue to normal routing
    return nexo.Continue(), nil
}
```

## Middleware

### File-based Middleware

```go
// app/api/middleware.go
package api

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

func Middleware() nexo.MiddlewareFunc {
    return func(next nexo.HandlerFunc) nexo.HandlerFunc {
        return func(c *nexo.Context) error {
            // Applied to all routes under /api
            c.SetHeader("X-API-Version", "1.0")
            return next(c)
        }
    }
}
```

### Built-in Middleware

```go
app := nexo.New()

// Request logging is enabled by default!
// Output: [12:34:56] GET /api/users 200 in 45ms (1.2KB)
// Customize with: app.SetLogger(nexo.RequestLoggerConfig{...})

// Panic recovery
app.Use(nexo.Recover())

// Request ID
app.Use(nexo.RequestID())

// CORS
app.Use(nexo.CORSWithConfig(nexo.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
}))

// Rate limiting
app.Use(nexo.RateLimiter(100, time.Minute))

// Timeout
app.Use(nexo.Timeout(30 * time.Second))

// Basic auth
app.Use(nexo.BasicAuth(func(user, pass string) bool {
    return user == "admin" && pass == "secret"
}))

// Security headers
app.Use(nexo.SecureHeaders())
```

## Context API

```go
func Get(c *nexo.Context) error {
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
func Get(c *nexo.Context) error {
    if notFound {
        return nexo.NotFound("resource not found")
    }
    
    if unauthorized {
        return nexo.Unauthorized("invalid token")
    }
    
    if badInput {
        return nexo.BadRequest("invalid input")
    }
    
    return nexo.InternalServerError("something went wrong")
}
```

## CLI Commands

```bash
# Create new project
nexo new myapp
nexo new myapp --api-only      # Without templ templates
nexo new myapp --with-proxy    # Include proxy.go example

# Development server with hot reload
nexo dev

# Build for production
nexo build

# List all routes
nexo routes
```

## Configuration

```yaml
# nexo.yaml
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

## Documentation

Full documentation is available at [nexo.build](https://nexo.build).

**Getting Started:**
- [Quick Start](https://nexo.build/docs/getting-started/quickstart)
- [Familiar Patterns](https://nexo.build/docs/getting-started/familiar-patterns)

**Core Concepts:**
- [File-based Routing](https://nexo.build/docs/routing/file-based)
- [Middleware](https://nexo.build/docs/middleware/overview)
- [Proxy](https://nexo.build/docs/middleware/proxy)
- [Templates](https://nexo.build/docs/core-concepts/templates)
- [Static Files](https://nexo.build/docs/core-concepts/static-files)

**Frontend:**
- [HTMX Integration](https://nexo.build/docs/frontend/htmx)
- [Tailwind CSS](https://nexo.build/docs/frontend/tailwind)
- [Forms](https://nexo.build/docs/frontend/forms)

**Guides:**
- [Examples](https://nexo.build/docs/guides/examples) - Working code examples for common patterns
- [Authentication](https://nexo.build/docs/guides/authentication)
- [Database](https://nexo.build/docs/guides/database)
- [Deployment](https://nexo.build/docs/guides/deployment)

**API Reference:**
- [Overview](https://nexo.build/docs/api/overview) - Quick reference tables for all types
- [App](https://nexo.build/docs/api/app) - Application lifecycle and routing
- [Context](https://nexo.build/docs/api/context) - Request/response methods
- [Config](https://nexo.build/docs/api/config) - Configuration options
- [Middleware](https://nexo.build/docs/api/middleware) - Built-in middleware reference
- [Proxy](https://nexo.build/docs/api/proxy) - Request interception
- [Errors](https://nexo.build/docs/api/errors) - Error handling helpers
- [CLI](https://nexo.build/docs/api/cli) - Command-line interface

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

## Acknowledgments

Nexo stands on the shoulders of giants:

- **[chi](https://github.com/go-chi/chi)** by Peter Kieltyka — The lightweight router that powers Nexo
- **[templ](https://github.com/a-h/templ)** by Adrian Hesketh — Type-safe HTML templating for Go
- **[fsnotify](https://github.com/fsnotify/fsnotify)** — Cross-platform file watching
- **[cobra](https://github.com/spf13/cobra)** by Steve Francia — CLI framework
- **[viper](https://github.com/spf13/viper)** by Steve Francia — Configuration management

## Author

**Abdul Hamid Achik** ([@abdulachik](https://x.com/abdulachik))

A Syrian-Mexican software engineer based in Guadalajara, Mexico.

- GitHub: [@abdul-hamid-achik](https://github.com/abdul-hamid-achik)
- Twitter/X: [@abdulachik](https://x.com/abdulachik)

## License

MIT License - see [LICENSE](LICENSE) for details.
