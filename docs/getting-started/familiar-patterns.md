# Familiar Patterns

If you've used modern meta-frameworks like Next.js, Nuxt, SvelteKit, or Remix, Fuego's conventions will feel natural.

## Quick Comparison

| Concept | Other Frameworks | Fuego |
|---------|------------------|-------|
| Route file | `route.ts` / `+server.ts` | `route.go` |
| Page file | `page.tsx` / `+page.svelte` | `page.templ` |
| Layout | `layout.tsx` / `+layout.svelte` | `layout.templ` |
| Middleware | `middleware.ts` | `middleware.go` |
| Dynamic segment | `[id]` | `[id]` |
| Catch-all | `[...slug]` | `[...slug]` |
| Optional catch-all | `[[...slug]]` | `[[...slug]]` |
| Route groups | `(group)` | `(group)` |

## Key Differences

### Handlers Are Named After HTTP Methods

In JavaScript frameworks, you export named functions like `GET`, `POST`, etc. In Fuego, you define functions with capitalized method names:

```go
// app/api/users/route.go
package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
    return c.JSON(200, []string{"alice", "bob"})
}

func Post(c *fuego.Context) error {
    var user User
    if err := c.Bind(&user); err != nil {
        return fuego.BadRequest("invalid input")
    }
    return c.JSON(201, user)
}

func Put(c *fuego.Context) error { ... }
func Patch(c *fuego.Context) error { ... }
func Delete(c *fuego.Context) error { ... }
```

### Templates Use templ, Not JSX

Instead of JSX/TSX, Fuego uses [templ](https://templ.guide) for type-safe HTML:

```go
// app/dashboard/page.templ
package dashboard

templ Page() {
    <div class="container">
        <h1>Dashboard</h1>
        <p>Welcome, { user.Name }!</p>
    </div>
}
```

Templ provides:
- Compile-time type checking
- No runtime template parsing
- Go code completion in templates
- Smaller binary size than text/template

### No Build Step for Go Code

Unlike JavaScript frameworks that require bundling, Fuego compiles to a single binary:

```bash
# Development
fuego dev

# Production
fuego build
./myapp
```

The binary includes everything — no need for `node_modules`, npm, or a separate build step for your Go code.

### HTMX Instead of Client-Side Frameworks

For interactivity, Fuego encourages HTMX over React/Vue/Svelte:

```html
<!-- Load data on page load -->
<div hx-get="/api/users" hx-trigger="load">
    Loading...
</div>

<!-- Submit form without page reload -->
<form hx-post="/api/users" hx-target="#user-list">
    <input name="name" />
    <button type="submit">Add User</button>
</form>
```

Benefits:
- No JavaScript to write
- Server renders HTML
- Smaller page weight
- Works with any backend

## Migration Patterns

### From Next.js App Router

| Next.js | Fuego |
|---------|-------|
| `app/api/users/route.ts` | `app/api/users/route.go` |
| `export async function GET()` | `func Get(c *fuego.Context) error` |
| `NextRequest` | `*fuego.Context` |
| `NextResponse.json()` | `c.JSON(status, data)` |
| `NextResponse.redirect()` | `c.Redirect(status, url)` |
| `cookies().get()` | `c.Cookie("name")` |
| `headers().get()` | `c.Header("name")` |

### From Nuxt

| Nuxt | Fuego |
|------|-------|
| `server/api/users.ts` | `app/api/users/route.go` |
| `defineEventHandler()` | `func Get(c *fuego.Context) error` |
| `getQuery(event)` | `c.Query("key")` |
| `readBody(event)` | `c.Bind(&body)` |
| `setResponseStatus()` | `c.JSON(status, data)` |

### From SvelteKit

| SvelteKit | Fuego |
|-----------|-------|
| `+server.ts` | `route.go` |
| `+page.svelte` | `page.templ` |
| `+layout.svelte` | `layout.templ` |
| `RequestHandler` | `func Get(c *fuego.Context) error` |
| `json()` helper | `c.JSON(status, data)` |

## File Structure Comparison

### JavaScript Framework (Next.js style)
```
app/
├── api/
│   └── users/
│       └── [id]/
│           └── route.ts
├── dashboard/
│   ├── page.tsx
│   └── layout.tsx
└── layout.tsx
```

### Fuego
```
app/
├── api/
│   └── users/
│       └── [id]/
│           └── route.go
├── dashboard/
│   ├── page.templ
│   └── layout.templ
└── layout.templ
```

The structure is nearly identical — just swap the file extensions.

## Why Go?

If you're comfortable with JavaScript frameworks, why use Fuego?

1. **Performance** — Go is significantly faster than Node.js for CPU-bound tasks
2. **Memory** — Lower memory footprint, important for containers
3. **Type safety** — Compile-time checks catch errors early
4. **Deployment** — Single binary, no runtime dependencies
5. **Concurrency** — Goroutines handle concurrent requests efficiently
6. **Simplicity** — No package manager drama, stable stdlib

## Next Steps

- [Quick Start](./quickstart.md) — Get up and running
- [File-based Routing](../routing/file-based.md) — Deep dive into routing
- [Context API](../api/context.md) — Request and response handling
