# Quick Start

Get up and running with Fuego in under 5 minutes.

## Installation

```bash
go install github.com/abdul-hamid-achik/fuego/cmd/fuego@latest
```

## Create a Project

```bash
fuego new myapp
cd myapp
```

This creates:
```
myapp/
├── app/
│   ├── api/
│   │   └── health/
│   │       └── route.go    # GET /api/health
│   ├── layout.templ        # HTML layout
│   └── page.templ          # Home page
├── static/                  # Static files
├── main.go                  # Entry point
├── fuego.yaml              # Configuration
└── go.mod
```

## Run Development Server

```bash
fuego dev
```

Visit http://localhost:3000

## Add an API Route

Create `app/api/users/route.go`:

```go
package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
    users := []map[string]any{
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
    }
    return c.JSON(200, users)
}

func Post(c *fuego.Context) error {
    var user struct {
        Name string `json:"name"`
    }
    if err := c.Bind(&user); err != nil {
        return fuego.BadRequest("invalid request body")
    }
    return c.JSON(201, map[string]any{
        "id": 3,
        "name": user.Name,
    })
}
```

Now you have:
- `GET /api/users` - List users
- `POST /api/users` - Create user

## Add Middleware

Create `app/api/middleware.go`:

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

This middleware applies to all routes under `/api/`.

## Add Dynamic Routes

Create `app/api/users/[id]/route.go`:

```go
package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
    id := c.Param("id")
    return c.JSON(200, map[string]any{
        "id": id,
        "name": "User " + id,
    })
}
```

Now `GET /api/users/123` returns `{"id": "123", "name": "User 123"}`.

## Build for Production

```bash
fuego build
./myapp
```

## Next Steps

- [Routing Guide](../routing/file-based.md)
- [Middleware](../middleware/overview.md)
- [Context API](../api/context.md)
