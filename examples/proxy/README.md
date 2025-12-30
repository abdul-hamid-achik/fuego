# Proxy Example

Demonstrates the Fuego proxy layer for request interception before routing.

## Features Demonstrated

1. **URL Rewriting** - Migrate legacy `/v1/*` paths to `/api/*`
2. **Access Control** - Block `/api/admin/*` without valid authorization
3. **Header Injection** - Add tracking headers to all requests
4. **Selective Matching** - Only run proxy for specific path patterns

## Structure

```
proxy/
├── app/
│   ├── proxy.go              # Proxy configuration
│   └── api/
│       ├── users/
│       │   └── route.go      # GET /api/users
│       └── admin/
│           └── route.go      # GET /api/admin (protected)
├── main.go
└── go.mod
```

## Running

```bash
cd examples/proxy
go mod tidy
go run .
```

## Testing

### 1. Normal API Request
```bash
curl http://localhost:3000/api/users
```

Response includes proxy headers:
```json
{
  "users": [...],
  "proxy": {
    "version": "1.0",
    "path": "/api/users"
  }
}
```

### 2. Legacy URL Rewriting
```bash
# This gets rewritten from /v1/users to /api/users
curl http://localhost:3000/v1/users
```

### 3. Admin Access Control

Without auth (401):
```bash
curl http://localhost:3000/api/admin
```

With invalid token (403):
```bash
curl -H "Authorization: Bearer wrong-token" http://localhost:3000/api/admin
```

With valid token (200):
```bash
curl -H "Authorization: Bearer admin-token" http://localhost:3000/api/admin
```

## Proxy Configuration

The proxy is defined in `app/proxy.go`:

```go
func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
    path := c.Path()

    // URL Rewriting
    if strings.HasPrefix(path, "/v1/") {
        return fuego.Rewrite(strings.Replace(path, "/v1/", "/api/", 1)), nil
    }

    // Access Control
    if strings.HasPrefix(path, "/api/admin") {
        if c.Header("Authorization") == "" {
            return fuego.ResponseJSON(401, `{"error":"unauthorized"}`), nil
        }
    }

    // Continue with headers
    return fuego.Continue().WithHeader("X-Proxy-Version", "1.0"), nil
}
```

## Available Proxy Results

- `fuego.Continue()` - Continue to normal routing
- `fuego.Rewrite(path)` - Rewrite URL internally
- `fuego.Redirect(url, status)` - HTTP redirect
- `fuego.Response(status, body, contentType)` - Direct response
- `fuego.ResponseJSON(status, json)` - JSON response
- `fuego.ResponseHTML(status, html)` - HTML response

All results support `.WithHeader(key, value)` and `.WithHeaders(map)` for adding response headers.
