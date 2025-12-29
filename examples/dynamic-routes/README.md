# Dynamic Routes Example

This example demonstrates Fuego's file-system based routing with dynamic segments, inspired by Next.js App Router.

## Structure

```
dynamic-routes/
├── main.go
├── go.mod
└── app/
    └── api/
        ├── users/
        │   ├── route.go          # GET /api/users
        │   └── [id]/
        │       └── route.go      # GET/PUT/DELETE /api/users/:id
        ├── posts/
        │   └── [...slug]/
        │       └── route.go      # GET /api/posts/* (catch-all)
        ├── docs/
        │   └── [[...path]]/
        │       └── route.go      # GET /api/docs and /api/docs/* (optional catch-all)
        └── (admin)/              # Route group (doesn't affect URL)
            └── settings/
                └── route.go      # GET/PUT /api/settings
```

## Route Segment Types

### 1. Dynamic Segments `[param]`

Directory names wrapped in square brackets create dynamic route parameters:

```
app/api/users/[id]/route.go  ->  /api/users/:id
```

Access the parameter in your handler:
```go
func Get(c *fuego.Context) error {
    id := c.Param("id")
    // id = "123" for /api/users/123
}
```

### 2. Catch-All Segments `[...param]`

Three dots before the parameter name match any number of path segments:

```
app/api/posts/[...slug]/route.go  ->  /api/posts/*
```

Examples:
- `/api/posts/hello` -> slug = `"hello"`
- `/api/posts/2024/01/my-post` -> slug = `"2024/01/my-post"`

```go
func Get(c *fuego.Context) error {
    slug := c.Param("slug")
    segments := strings.Split(slug, "/")
}
```

### 3. Optional Catch-All Segments `[[...param]]`

Double brackets make the catch-all optional, also matching the root:

```
app/api/docs/[[...path]]/route.go  ->  /api/docs and /api/docs/*
```

Examples:
- `/api/docs` -> path = `""` (empty, matches root)
- `/api/docs/intro` -> path = `"intro"`
- `/api/docs/api/users` -> path = `"api/users"`

### 4. Route Groups `(name)`

Parentheses create route groups for organization without affecting the URL:

```
app/api/(admin)/settings/route.go  ->  /api/settings (NOT /api/admin/settings)
```

Use cases:
- Organize related routes together
- Apply shared middleware to a group
- Keep code organized without URL changes

## Running the Example

```bash
cd examples/dynamic-routes
go run main.go
```

## Testing the Routes

### Static route - List users:
```bash
curl http://localhost:3000/api/users
```

### Dynamic segment - Get user by ID:
```bash
curl http://localhost:3000/api/users/1
curl http://localhost:3000/api/users/2
curl http://localhost:3000/api/users/999  # Returns 404
```

### Catch-all - Posts with any path:
```bash
curl http://localhost:3000/api/posts/hello
curl http://localhost:3000/api/posts/2024/01/my-first-post
curl http://localhost:3000/api/posts/category/tech/article
```

### Optional catch-all - Docs root and sub-pages:
```bash
curl http://localhost:3000/api/docs           # Root (shows index)
curl http://localhost:3000/api/docs/intro     # Specific page
curl http://localhost:3000/api/docs/api/users # Nested page
```

### Route group - Settings (no /admin in URL):
```bash
curl http://localhost:3000/api/settings
```

## Route Priority

When multiple routes could match, Fuego uses this priority:
1. Static routes (exact match)
2. Dynamic segments `[param]`
3. Catch-all segments `[...param]`
4. Optional catch-all `[[...param]]`

For example, given:
- `/api/users/route.go`
- `/api/users/[id]/route.go`
- `/api/users/[...slug]/route.go`

A request to `/api/users` matches the static route, `/api/users/123` matches `[id]`, and `/api/users/123/profile` matches `[...slug]`.
