package nexo

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// ---------- App Construction Tests ----------

func TestApp_New(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("expected app to be non-nil")
	}
	if app.router == nil {
		t.Error("expected router to be initialized")
	}
	if app.config == nil {
		t.Error("expected config to be initialized")
	}
	if app.routeTree == nil {
		t.Error("expected routeTree to be initialized")
	}
	if app.scanner == nil {
		t.Error("expected scanner to be initialized")
	}
}

func TestApp_NewWithOptions(t *testing.T) {
	app := New(
		WithPort("8080"),
		WithHost("127.0.0.1"),
	)

	if app.config.Port != "8080" {
		t.Errorf("expected port 8080, got %s", app.config.Port)
	}
	if app.config.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", app.config.Host)
	}
}

// ---------- Accessor Method Tests ----------

func TestApp_Use(t *testing.T) {
	app := New()

	middlewareCalled := false
	mw := func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			middlewareCalled = true
			return next(c)
		}
	}

	app.Use(mw)

	if len(app.middlewares) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(app.middlewares))
	}

	// Verify middleware works when mounted
	app.Get("/test", func(c *Context) error {
		return c.String(200, "ok")
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	app.ServeHTTP(w, r)

	if !middlewareCalled {
		t.Error("expected middleware to be called")
	}
}

func TestApp_Router(t *testing.T) {
	app := New()

	router := app.Router()
	if router == nil {
		t.Error("expected router to be non-nil")
	}
}

func TestApp_Config(t *testing.T) {
	app := New(WithPort("9999"))

	config := app.Config()
	if config == nil {
		t.Fatal("expected config to be non-nil")
	}
	if config.Port != "9999" {
		t.Errorf("expected port 9999, got %s", config.Port)
	}
}

func TestApp_RouteTree(t *testing.T) {
	app := New()

	tree := app.RouteTree()
	if tree == nil {
		t.Error("expected route tree to be non-nil")
	}
}

func TestApp_Scan(t *testing.T) {
	// Create temp directory with a route file
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app", "api")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	routeContent := `package api

import "github.com/abdul-hamid-achik/nexo/pkg/nexo"

func Get(c *nexo.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "route.go"), []byte(routeContent), 0644); err != nil {
		t.Fatalf("failed to write route.go: %v", err)
	}

	app := New(WithAppDir(filepath.Join(tmpDir, "app")))

	if err := app.Scan(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	routes := app.RouteTree().Routes()
	if len(routes) != 1 {
		t.Errorf("expected 1 route, got %d", len(routes))
	}
}

func TestApp_Scan_NonExistentDir(t *testing.T) {
	app := New(WithAppDir("/nonexistent/path"))

	// Should not error, just find no routes
	if err := app.Scan(); err != nil {
		t.Errorf("expected no error for non-existent dir, got %v", err)
	}
}

// ---------- HTTP Method Tests ----------

func TestApp_Post(t *testing.T) {
	app := New()
	app.Post("/users", func(c *Context) error {
		return c.JSON(201, map[string]string{"created": "true"})
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", nil)
	app.ServeHTTP(w, r)

	if w.Code != 201 {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestApp_Put(t *testing.T) {
	app := New()
	app.Put("/users/1", func(c *Context) error {
		return c.JSON(200, map[string]string{"updated": "true"})
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/users/1", nil)
	app.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestApp_Patch(t *testing.T) {
	app := New()
	app.Patch("/users/1", func(c *Context) error {
		return c.JSON(200, map[string]string{"patched": "true"})
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PATCH", "/users/1", nil)
	app.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestApp_Delete(t *testing.T) {
	app := New()
	app.Delete("/users/1", func(c *Context) error {
		return c.NoContent()
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users/1", nil)
	app.ServeHTTP(w, r)

	if w.Code != 204 {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestApp_Head(t *testing.T) {
	app := New()
	app.Head("/users", func(c *Context) error {
		c.SetHeader("X-Total-Count", "100")
		return c.NoContent()
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("HEAD", "/users", nil)
	app.ServeHTTP(w, r)

	if w.Code != 204 {
		t.Errorf("expected 204, got %d", w.Code)
	}
	if w.Header().Get("X-Total-Count") != "100" {
		t.Error("expected X-Total-Count header")
	}
}

func TestApp_Options(t *testing.T) {
	app := New()
	app.Options("/users", func(c *Context) error {
		c.SetHeader("Allow", "GET, POST, PUT, DELETE")
		return c.NoContent()
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("OPTIONS", "/users", nil)
	app.ServeHTTP(w, r)

	if w.Code != 204 {
		t.Errorf("expected 204, got %d", w.Code)
	}
	if w.Header().Get("Allow") != "GET, POST, PUT, DELETE" {
		t.Error("expected Allow header")
	}
}

func TestApp_AllHTTPMethods(t *testing.T) {
	tests := []struct {
		method     string
		register   func(*App, string, HandlerFunc)
		statusCode int
	}{
		{"GET", func(a *App, p string, h HandlerFunc) { a.Get(p, h) }, 200},
		{"POST", func(a *App, p string, h HandlerFunc) { a.Post(p, h) }, 200},
		{"PUT", func(a *App, p string, h HandlerFunc) { a.Put(p, h) }, 200},
		{"PATCH", func(a *App, p string, h HandlerFunc) { a.Patch(p, h) }, 200},
		{"DELETE", func(a *App, p string, h HandlerFunc) { a.Delete(p, h) }, 200},
		{"HEAD", func(a *App, p string, h HandlerFunc) { a.Head(p, h) }, 200},
		{"OPTIONS", func(a *App, p string, h HandlerFunc) { a.Options(p, h) }, 200},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			app := New()
			tt.register(app, "/test", func(c *Context) error {
				return c.String(200, "ok")
			})
			app.Mount()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "/test", nil)
			app.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("expected %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

// ---------- Static File Tests ----------

func TestApp_Static(t *testing.T) {
	// Create temp directory with a file
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	app := New()
	app.Static("/static", tmpDir)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/static/test.txt", nil)
	app.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "hello" {
		t.Errorf("expected body 'hello', got %q", w.Body.String())
	}
}

func TestApp_Static_WithoutLeadingSlash(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	app := New()
	app.Static("assets", tmpDir) // No leading slash

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/assets/file.txt", nil)
	app.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestApp_Static_EmptyPath(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	app := New()
	app.Static("", tmpDir) // Empty path defaults to "/"

	// Test that the path was normalized to "/" by checking route is registered
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test.txt", nil)
	app.ServeHTTP(w, r)

	// 200 for direct file access or 301 for directory redirect are both valid
	// The important thing is that the route is registered (not 404)
	if w.Code == 404 {
		t.Errorf("expected route to be registered, got 404")
	}
}

// ---------- Route Group Tests ----------

func TestApp_Group(t *testing.T) {
	app := New()

	app.Group("/api", func(g *RouteGroup) {
		g.Get("/users", func(c *Context) error {
			return c.JSON(200, map[string]string{"users": "list"})
		})
		g.Post("/users", func(c *Context) error {
			return c.JSON(201, map[string]string{"user": "created"})
		})
	})

	app.Mount()

	// Test GET /api/users
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/users", nil)
	app.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("GET /api/users: expected 200, got %d", w.Code)
	}

	// Test POST /api/users
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/users", nil)
	app.ServeHTTP(w, r)

	if w.Code != 201 {
		t.Errorf("POST /api/users: expected 201, got %d", w.Code)
	}
}

func TestApp_Group_WithMiddleware(t *testing.T) {
	app := New()

	app.Group("/api", func(g *RouteGroup) {
		g.Use(func(next HandlerFunc) HandlerFunc {
			return func(c *Context) error {
				c.SetHeader("X-Group-MW", "true")
				return next(c)
			}
		})

		g.Get("/test", func(c *Context) error {
			return c.String(200, "ok")
		})
	})

	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/test", nil)
	app.ServeHTTP(w, r)

	if w.Header().Get("X-Group-MW") != "true" {
		t.Error("expected group middleware to set header")
	}
}

func TestRouteGroup_Use(t *testing.T) {
	app := New()

	var order []string

	app.Group("/v1", func(g *RouteGroup) {
		g.Use(func(next HandlerFunc) HandlerFunc {
			return func(c *Context) error {
				order = append(order, "mw1")
				return next(c)
			}
		})
		g.Use(func(next HandlerFunc) HandlerFunc {
			return func(c *Context) error {
				order = append(order, "mw2")
				return next(c)
			}
		})

		g.Get("/test", func(c *Context) error {
			order = append(order, "handler")
			return c.String(200, "ok")
		})
	})

	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/v1/test", nil)
	app.ServeHTTP(w, r)

	expected := []string{"mw1", "mw2", "handler"}
	if len(order) != len(expected) {
		t.Errorf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if i < len(order) && order[i] != v {
			t.Errorf("expected order[%d] = %q, got %q", i, v, order[i])
		}
	}
}

func TestRouteGroup_AllMethods(t *testing.T) {
	app := New()

	app.Group("/api", func(g *RouteGroup) {
		g.Get("/resource", func(c *Context) error { return c.String(200, "get") })
		g.Post("/resource", func(c *Context) error { return c.String(201, "post") })
		g.Put("/resource", func(c *Context) error { return c.String(200, "put") })
		g.Patch("/resource", func(c *Context) error { return c.String(200, "patch") })
		g.Delete("/resource", func(c *Context) error { return c.String(200, "delete") })
	})

	app.Mount()

	tests := []struct {
		method string
		status int
		body   string
	}{
		{"GET", 200, "get"},
		{"POST", 201, "post"},
		{"PUT", 200, "put"},
		{"PATCH", 200, "patch"},
		{"DELETE", 200, "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "/api/resource", nil)
			app.ServeHTTP(w, r)

			if w.Code != tt.status {
				t.Errorf("expected %d, got %d", tt.status, w.Code)
			}
			if w.Body.String() != tt.body {
				t.Errorf("expected body %q, got %q", tt.body, w.Body.String())
			}
		})
	}
}

// ---------- Proxy Error Handling Tests ----------

func TestApp_ServeHTTP_ProxyError(t *testing.T) {
	app := New()

	// Set proxy that returns an error
	_ = app.SetProxy(func(c *Context) (*ProxyResult, error) {
		return nil, NewHTTPError(500, "proxy error")
	}, nil)

	app.Get("/test", func(c *Context) error {
		return c.String(200, "should not reach here")
	})
	app.Mount()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	app.ServeHTTP(w, r)

	if w.Code != 500 {
		t.Errorf("expected 500 for proxy error, got %d", w.Code)
	}
}

// ---------- Addr Tests ----------

func TestApp_Addr_BeforeStart(t *testing.T) {
	app := New()

	addr := app.Addr()
	if addr != "" {
		t.Errorf("expected empty addr before start, got %q", addr)
	}
}
