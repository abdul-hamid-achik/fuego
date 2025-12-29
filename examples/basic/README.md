# Basic Example

A simple Fuego application demonstrating the basics of file-based routing.

## Structure

```
basic/
├── app/
│   └── api/
│       ├── middleware.go      # API middleware
│       ├── health/
│       │   └── route.go       # GET /api/health
│       └── users/
│           └── route.go       # GET/POST /api/users
├── static/                     # Static files
├── main.go                     # Entry point
└── go.mod
```

## Running

```bash
cd examples/basic
go run .
```

## Endpoints

- `GET /api/health` - Health check
- `GET /api/users` - List users
- `POST /api/users` - Create user

## Testing

```bash
# Health check
curl http://localhost:3000/api/health

# List users
curl http://localhost:3000/api/users

# Create user
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie","email":"charlie@example.com"}'
```
