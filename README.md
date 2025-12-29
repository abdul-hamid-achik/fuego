# Fuego

A file-system based Go framework for APIs and websites, inspired by Next.js App Router.

The name "Fuego" (Spanish for "fire") contains "GO" naturally embedded (fue**GO**) and reflects its Mexican origin and blazing fast performance.

## Features

- **File system is the router** - No manual route registration. Drop a file, get a route.
- **Zero boilerplate to start** - `main.go` is 5 lines.
- **Explicit over magic** - Handlers are plain Go functions.
- **Fast iteration** - Hot reload in dev, sub-second rebuilds.
- **Scalable conventions** - Works for a 3-route API or a 300-route app.
- **Templ-native** - First-class support for type-safe HTML templating.

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

## Project Structure

```
myapp/
├── app/
│   ├── layout.templ          # Root layout
│   ├── page.templ            # GET /
│   └── api/
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
| `page.templ` | UI for a route |
| `layout.templ` | Shared UI wrapper |
| `middleware.go` | Middleware for segment and children |
| `loader.go` | Data fetching for pages |
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
        return c.Error(400, "invalid input")
    }
    return c.JSON(201, map[string]string{"created": input.Name})
}
```

## Middleware

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

## Configuration

```yaml
# fuego.yaml
port: 3000
app_dir: "app"
static_dir: "static"
static_path: "/static"

dev:
  hot_reload: true
  watch_extensions: [".go", ".templ"]

middleware:
  logger: true
  recover: true
```

## Documentation

- [Getting Started](docs/getting-started.md)
- [Routing](docs/routing.md)
- [Middleware](docs/middleware.md)
- [Templating](docs/templating.md)
- [API Reference](docs/api-reference.md)

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

**Abdul Hamid Achik** - [@abdul-hamid-achik](https://github.com/abdul-hamid-achik)

## License

MIT License - see [LICENSE](LICENSE) for details.
