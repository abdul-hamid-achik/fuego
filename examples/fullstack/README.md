# Fullstack Example

A complete task manager demonstrating Fuego's full-stack capabilities with templ templates, Tailwind CSS, and HTMX.

## Features

- **Templ templates** — Type-safe HTML with Go
- **Tailwind CSS v4** — Standalone binary, no Node.js required
- **HTMX** — Interactive UI without writing JavaScript

## Running

```bash
cd examples/fullstack

# Install dependencies
go mod tidy

# Start the development server
fuego dev
```

Visit http://localhost:3000

## Project Structure

```
fullstack/
├── app/
│   ├── api/
│   │   └── tasks/
│   │       ├── route.go         # GET, POST, DELETE /api/tasks
│   │       └── toggle/
│   │           └── route.go     # POST /api/tasks/toggle
│   ├── dashboard/
│   │   └── page.templ           # Task manager UI
│   ├── layout.templ             # Shared HTML layout
│   └── page.templ               # Home page
├── internal/
│   └── tasks/
│       └── store.go             # Shared in-memory task storage
├── styles/
│   └── input.css                # Tailwind source CSS
├── static/
│   └── css/
│       └── output.css           # Compiled Tailwind (generated)
├── main.go                      # Application entry point
└── go.mod
```

## Key Files

| File | Purpose |
|------|---------|
| `app/layout.templ` | HTML shell with Tailwind and HTMX |
| `app/page.templ` | Landing page with feature highlights |
| `app/dashboard/page.templ` | Interactive task list UI |
| `app/api/tasks/route.go` | Task CRUD endpoints |
| `app/api/tasks/toggle/route.go` | Toggle task completion |
| `internal/tasks/store.go` | Thread-safe task storage |

## How HTMX Works

The task list uses HTMX for seamless updates without page reloads:

### Loading Tasks

```html
<div id="task-list" hx-get="/api/tasks" hx-trigger="load">
  Loading...
</div>
```

On page load, HTMX fetches `/api/tasks` and replaces the div content with the response.

### Adding Tasks

```html
<form hx-post="/api/tasks" hx-target="#task-list">
  <input type="text" name="title" />
  <button type="submit">Add</button>
</form>
```

Form submission POSTs to `/api/tasks`, and the response replaces `#task-list`.

### Toggling Completion

```html
<input type="checkbox" hx-post="/api/tasks/toggle?id=1" hx-target="#task-list" />
```

Clicking the checkbox POSTs to toggle and refreshes the list.

### Deleting Tasks

```html
<button hx-delete="/api/tasks?id=1" hx-target="#task-list">Delete</button>
```

## Tailwind CSS

Tailwind is built automatically when running `fuego dev`:

1. Source CSS is in `styles/input.css`
2. Compiled CSS goes to `static/css/output.css`
3. The layout references `/static/css/output.css`

To manually build CSS:

```bash
fuego tailwind build    # Production build (minified)
fuego tailwind watch    # Development (watch mode)
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/tasks` | List all tasks (HTML) |
| POST | `/api/tasks` | Create a task |
| DELETE | `/api/tasks?id=N` | Delete a task |
| POST | `/api/tasks/toggle?id=N` | Toggle completion |

All endpoints return HTML for HTMX consumption.
