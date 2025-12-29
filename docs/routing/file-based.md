# File-Based Routing

Fuego uses a file-system based router inspired by Next.js App Router. Your file structure defines your routes.

## Basic Routing

Create `route.go` files to define API endpoints:

```
app/
├── api/
│   ├── users/
│   │   └── route.go    → /api/users
│   └── posts/
│       └── route.go    → /api/posts
└── route.go            → /
```

## Handler Functions

Export functions named after HTTP methods:

```go
package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// GET /api/users
func Get(c *fuego.Context) error {
    return c.JSON(200, []string{"alice", "bob"})
}

// POST /api/users
func Post(c *fuego.Context) error {
    return c.JSON(201, map[string]string{"status": "created"})
}

// PUT /api/users
func Put(c *fuego.Context) error {
    return c.JSON(200, map[string]string{"status": "updated"})
}

// PATCH /api/users
func Patch(c *fuego.Context) error {
    return c.JSON(200, map[string]string{"status": "patched"})
}

// DELETE /api/users
func Delete(c *fuego.Context) error {
    return c.NoContent()
}

// HEAD /api/users
func Head(c *fuego.Context) error {
    return c.NoContent()
}

// OPTIONS /api/users
func Options(c *fuego.Context) error {
    c.SetHeader("Allow", "GET, POST, PUT, PATCH, DELETE")
    return c.NoContent()
}
```

## Dynamic Routes

Use `[param]` folders for dynamic segments:

```
app/api/users/[id]/route.go    → /api/users/:id
```

```go
package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
    id := c.Param("id")
    return c.JSON(200, map[string]string{"id": id})
}
```

### Multiple Parameters

```
app/api/orgs/[orgId]/teams/[teamId]/route.go
→ /api/orgs/:orgId/teams/:teamId
```

```go
func Get(c *fuego.Context) error {
    orgId := c.Param("orgId")
    teamId := c.Param("teamId")
    return c.JSON(200, map[string]string{
        "orgId": orgId,
        "teamId": teamId,
    })
}
```

## Catch-All Routes

Use `[...param]` for catch-all routes:

```
app/docs/[...slug]/route.go    → /docs/*
```

```go
func Get(c *fuego.Context) error {
    // /docs/api/users/create → slug = "api/users/create"
    slug := c.Param("slug")
    return c.JSON(200, map[string]string{"slug": slug})
}
```

## Optional Catch-All

Use `[[...param]]` for optional catch-all (matches with or without segments):

```
app/shop/[[...categories]]/route.go
```

Matches:
- `/shop` → categories = ""
- `/shop/electronics` → categories = "electronics"
- `/shop/electronics/phones` → categories = "electronics/phones"

## Route Groups

Use `(groupname)` folders to organize without affecting URLs:

```
app/
├── (marketing)/
│   ├── about/route.go     → /about
│   └── blog/route.go      → /blog
├── (shop)/
│   ├── products/route.go  → /products
│   └── cart/route.go      → /cart
```

Groups help organize code without adding URL segments.

## Private Folders

Folders starting with `_` are ignored:

```
app/
├── _components/           # Ignored - shared components
│   └── button.go
├── _utils/               # Ignored - utilities
│   └── helpers.go
└── api/
    └── route.go          # Routable
```

## Route Priority

Routes are matched in order of specificity:

1. **Static routes** (highest priority)
   - `/api/users/me` matches before `/api/users/:id`

2. **Dynamic routes**
   - `/api/users/:id` matches after static routes

3. **Catch-all routes** (lowest priority)
   - `/docs/*` matches last

Example:
```
GET /api/users/me → matches /api/users/me (static)
GET /api/users/123 → matches /api/users/:id (dynamic)
GET /docs/anything/here → matches /docs/* (catch-all)
```

## Viewing Routes

Use the CLI to list all routes:

```bash
fuego routes
```

Output:
```
  Fuego Routes

  GET     /                             app/route.go
  GET     /api/health                   app/api/health/route.go
  GET     /api/users                    app/api/users/route.go
  POST    /api/users                    app/api/users/route.go
  GET     /api/users/{id}               app/api/users/[id]/route.go

  Total: 5 routes
```

## Handler Signature

All handlers must have this signature:

```go
func HandlerName(c *fuego.Context) error
```

Invalid signatures are skipped with a warning (in verbose mode).
