package fuego

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_PathToRoute(t *testing.T) {
	tests := []struct {
		name     string
		appDir   string
		filePath string
		want     string
	}{
		{
			name:     "root route",
			appDir:   "app",
			filePath: "app/route.go",
			want:     "/",
		},
		{
			name:     "simple nested route",
			appDir:   "app",
			filePath: "app/users/route.go",
			want:     "/users",
		},
		{
			name:     "deeply nested route",
			appDir:   "app",
			filePath: "app/api/users/profile/route.go",
			want:     "/api/users/profile",
		},
		{
			name:     "dynamic segment",
			appDir:   "app",
			filePath: "app/users/[id]/route.go",
			want:     "/users/{id}",
		},
		{
			name:     "multiple dynamic segments",
			appDir:   "app",
			filePath: "app/orgs/[orgId]/teams/[teamId]/route.go",
			want:     "/orgs/{orgId}/teams/{teamId}",
		},
		{
			name:     "catch-all segment",
			appDir:   "app",
			filePath: "app/docs/[...slug]/route.go",
			want:     "/docs/*",
		},
		{
			name:     "optional catch-all",
			appDir:   "app",
			filePath: "app/shop/[[...categories]]/route.go",
			want:     "/shop/*",
		},
		{
			name:     "route group",
			appDir:   "app",
			filePath: "app/(auth)/login/route.go",
			want:     "/login",
		},
		{
			name:     "multiple route groups",
			appDir:   "app",
			filePath: "app/(marketing)/(landing)/about/route.go",
			want:     "/about",
		},
		{
			name:     "route group with dynamic segment",
			appDir:   "app",
			filePath: "app/(api)/users/[id]/route.go",
			want:     "/users/{id}",
		},
		{
			name:     "complex nested path",
			appDir:   "app",
			filePath: "app/(admin)/dashboard/users/[userId]/posts/[postId]/route.go",
			want:     "/dashboard/users/{userId}/posts/{postId}",
		},
		{
			name:     "api route",
			appDir:   "app",
			filePath: "app/api/health/route.go",
			want:     "/api/health",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.appDir)
			got := s.pathToRoute(tt.filePath)
			if got != tt.want {
				t.Errorf("pathToRoute() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestScanner_Scan_BasicRoute(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	healthDir := filepath.Join(appDir, "api", "health")

	if err := os.MkdirAll(healthDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	// Create a valid route.go file
	routeContent := `package health

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}

func Post(c *fuego.Context) error {
	return c.JSON(201, nil)
}
`
	routePath := filepath.Join(healthDir, "route.go")
	if err := os.WriteFile(routePath, []byte(routeContent), 0644); err != nil {
		t.Fatalf("Failed to write route.go: %v", err)
	}

	// Scan
	scanner := NewScanner(appDir)
	tree := NewRouteTree()

	if err := scanner.Scan(tree); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	routes := tree.Routes()
	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}

	// Check route patterns and methods
	foundGet := false
	foundPost := false
	for _, r := range routes {
		if r.Pattern == "/api/health" {
			if r.Method == "GET" {
				foundGet = true
			}
			if r.Method == "POST" {
				foundPost = true
			}
		}
	}

	if !foundGet {
		t.Error("Expected GET /api/health route")
	}
	if !foundPost {
		t.Error("Expected POST /api/health route")
	}
}

func TestScanner_Scan_DynamicRoute(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	usersDir := filepath.Join(appDir, "users", "[id]")

	if err := os.MkdirAll(usersDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	routeContent := `package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
	return nil
}
`
	routePath := filepath.Join(usersDir, "route.go")
	if err := os.WriteFile(routePath, []byte(routeContent), 0644); err != nil {
		t.Fatalf("Failed to write route.go: %v", err)
	}

	scanner := NewScanner(appDir)
	tree := NewRouteTree()

	if err := scanner.Scan(tree); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	routes := tree.Routes()
	if len(routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(routes))
	}

	if routes[0].Pattern != "/users/{id}" {
		t.Errorf("Expected pattern '/users/{id}', got '%s'", routes[0].Pattern)
	}
}

func TestScanner_Scan_SkipsPrivateFolders(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	privateDir := filepath.Join(appDir, "_components")
	publicDir := filepath.Join(appDir, "public")

	if err := os.MkdirAll(privateDir, 0755); err != nil {
		t.Fatalf("Failed to create private dir: %v", err)
	}
	if err := os.MkdirAll(publicDir, 0755); err != nil {
		t.Fatalf("Failed to create public dir: %v", err)
	}

	// Route in private folder (should be ignored)
	privateRoute := `package components

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(privateDir, "route.go"), []byte(privateRoute), 0644); err != nil {
		t.Fatalf("Failed to write private route.go: %v", err)
	}

	// Route in public folder (should be found)
	publicRoute := `package public

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(publicDir, "route.go"), []byte(publicRoute), 0644); err != nil {
		t.Fatalf("Failed to write public route.go: %v", err)
	}

	scanner := NewScanner(appDir)
	tree := NewRouteTree()

	if err := scanner.Scan(tree); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	routes := tree.Routes()
	if len(routes) != 1 {
		t.Errorf("Expected 1 route (private folder should be skipped), got %d", len(routes))
	}

	if routes[0].Pattern != "/public" {
		t.Errorf("Expected pattern '/public', got '%s'", routes[0].Pattern)
	}
}

func TestScanner_Scan_RouteGroup(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	authDir := filepath.Join(appDir, "(auth)", "login")

	if err := os.MkdirAll(authDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	routeContent := `package login

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
	return nil
}

func Post(c *fuego.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(authDir, "route.go"), []byte(routeContent), 0644); err != nil {
		t.Fatalf("Failed to write route.go: %v", err)
	}

	scanner := NewScanner(appDir)
	tree := NewRouteTree()

	if err := scanner.Scan(tree); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	routes := tree.Routes()
	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}

	// Route group should not appear in the pattern
	for _, r := range routes {
		if r.Pattern != "/login" {
			t.Errorf("Expected pattern '/login' (group stripped), got '%s'", r.Pattern)
		}
	}
}

func TestScanner_Scan_SkipsInvalidSignatures(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	testDir := filepath.Join(appDir, "test")

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	// Route with invalid signatures
	routeContent := `package test

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// Valid handler
func Get(c *fuego.Context) error {
	return nil
}

// Invalid: wrong parameter type
func Post(w http.ResponseWriter, r *http.Request) {
}

// Invalid: wrong return type
func Put(c *fuego.Context) string {
	return ""
}

// Invalid: too many parameters
func Patch(c *fuego.Context, extra string) error {
	return nil
}

// Invalid: unexported
func delete(c *fuego.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(testDir, "route.go"), []byte(routeContent), 0644); err != nil {
		t.Fatalf("Failed to write route.go: %v", err)
	}

	scanner := NewScanner(appDir)
	tree := NewRouteTree()

	if err := scanner.Scan(tree); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	routes := tree.Routes()
	// Only the valid Get handler should be registered
	if len(routes) != 1 {
		t.Errorf("Expected 1 valid route, got %d", len(routes))
	}

	if len(routes) > 0 && routes[0].Method != "GET" {
		t.Errorf("Expected GET method, got %s", routes[0].Method)
	}
}

func TestScanner_Scan_NonExistentDir(t *testing.T) {
	scanner := NewScanner("/nonexistent/path")
	tree := NewRouteTree()

	// Should not return an error, just no routes
	if err := scanner.Scan(tree); err != nil {
		t.Errorf("Expected no error for non-existent dir, got: %v", err)
	}

	if len(tree.Routes()) != 0 {
		t.Errorf("Expected 0 routes, got %d", len(tree.Routes()))
	}
}

func TestScanner_ScanRouteInfo(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	usersDir := filepath.Join(appDir, "users")

	if err := os.MkdirAll(usersDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	routeContent := `package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
	return nil
}

func Post(c *fuego.Context) error {
	return nil
}

func Delete(c *fuego.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(usersDir, "route.go"), []byte(routeContent), 0644); err != nil {
		t.Fatalf("Failed to write route.go: %v", err)
	}

	scanner := NewScanner(appDir)
	routes, err := scanner.ScanRouteInfo()
	if err != nil {
		t.Fatalf("ScanRouteInfo failed: %v", err)
	}

	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	methods := make(map[string]bool)
	for _, r := range routes {
		methods[r.Method] = true
		if r.Pattern != "/users" {
			t.Errorf("Expected pattern '/users', got '%s'", r.Pattern)
		}
	}

	if !methods["GET"] || !methods["POST"] || !methods["DELETE"] {
		t.Error("Missing expected HTTP methods")
	}
}

func TestCalculatePriority(t *testing.T) {
	tests := []struct {
		pattern  string
		expected int
	}{
		{"/", 100},
		{"/users", 100},
		{"/api/health", 100},
		{"/users/{id}", 50},
		{"/orgs/{orgId}/teams/{teamId}", 50},
		{"/docs/*", 5},
		{"/*", 5},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			priority := CalculatePriority(tt.pattern)
			if priority != tt.expected {
				t.Errorf("CalculatePriority(%q) = %d, want %d", tt.pattern, priority, tt.expected)
			}
		})
	}
}

// ---------- Proxy Scanning Tests ----------

func TestScanner_ScanProxyInfo_ValidProxy(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	proxyContent := `package app

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
	return fuego.Continue(), nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if !info.HasProxy {
		t.Error("expected HasProxy to be true")
	}
	if info.FilePath == "" {
		t.Error("expected FilePath to be set")
	}
}

func TestScanner_ScanProxyInfo_ValidProxyWithContext(t *testing.T) {
	// Test with just "Context" (same package) instead of "fuego.Context"
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	proxyContent := `package app

func Proxy(c *Context) (*ProxyResult, error) {
	return Continue(), nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if !info.HasProxy {
		t.Error("expected HasProxy to be true for same-package types")
	}
}

func TestScanner_ScanProxyInfo_NoProxy(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// No proxy.go file exists

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if info.HasProxy {
		t.Error("expected HasProxy to be false")
	}
}

func TestScanner_ScanProxyInfo_InvalidSignature_WrongParams(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// Proxy with wrong parameter count
	proxyContent := `package app

func Proxy() (*ProxyResult, error) {
	return nil, nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if info.HasProxy {
		t.Error("expected HasProxy to be false for invalid signature")
	}
}

func TestScanner_ScanProxyInfo_InvalidSignature_WrongReturn(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// Proxy with wrong return type
	proxyContent := `package app

func Proxy(c *Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if info.HasProxy {
		t.Error("expected HasProxy to be false for wrong return type")
	}
}

func TestScanner_ScanProxyInfo_InvalidSignature_NotPointer(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// Proxy with non-pointer parameter
	proxyContent := `package app

func Proxy(c Context) (*ProxyResult, error) {
	return nil, nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if info.HasProxy {
		t.Error("expected HasProxy to be false for non-pointer param")
	}
}

func TestScanner_ScanProxyInfo_WithMatchers(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	proxyContent := `package app

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

var ProxyConfig = fuego.ProxyConfig{
	Matcher: []string{
		"/api/*",
		"/admin/*",
	},
}

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
	return fuego.Continue(), nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if !info.HasProxy {
		t.Error("expected HasProxy to be true")
	}

	if len(info.Matchers) != 2 {
		t.Errorf("expected 2 matchers, got %d", len(info.Matchers))
	}

	expectedMatchers := []string{"/api/*", "/admin/*"}
	for i, expected := range expectedMatchers {
		if i < len(info.Matchers) && info.Matchers[i] != expected {
			t.Errorf("expected matcher[%d] = %q, got %q", i, expected, info.Matchers[i])
		}
	}
}

func TestScanner_ScanProxyInfo_WithMatchersPointer(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// ProxyConfig with pointer syntax (&fuego.ProxyConfig{})
	proxyContent := `package app

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

var ProxyConfig = &fuego.ProxyConfig{
	Matcher: []string{
		"/v1/*",
	},
}

func Proxy(c *fuego.Context) (*fuego.ProxyResult, error) {
	return fuego.Continue(), nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	info, err := scanner.ScanProxyInfo()
	if err != nil {
		t.Fatalf("ScanProxyInfo failed: %v", err)
	}

	if len(info.Matchers) != 1 {
		t.Errorf("expected 1 matcher, got %d", len(info.Matchers))
	}

	if len(info.Matchers) > 0 && info.Matchers[0] != "/v1/*" {
		t.Errorf("expected matcher /v1/*, got %q", info.Matchers[0])
	}
}

func TestScanner_ScanProxyInfo_ParseError(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// Invalid Go syntax
	proxyContent := `package app

func Proxy(c *Context) {
	this is not valid go code
}
`
	if err := os.WriteFile(filepath.Join(appDir, "proxy.go"), []byte(proxyContent), 0644); err != nil {
		t.Fatalf("failed to write proxy.go: %v", err)
	}

	scanner := NewScanner(appDir)
	_, err := scanner.ScanProxyInfo()
	if err == nil {
		t.Error("expected error for invalid Go syntax")
	}
}

// ---------- Middleware Scanning Tests ----------

func TestScanner_ScanMiddlewareInfo_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	apiDir := filepath.Join(appDir, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		t.Fatalf("failed to create api dir: %v", err)
	}

	middlewareContent := `package api

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Middleware() fuego.MiddlewareFunc {
	return func(next fuego.HandlerFunc) fuego.HandlerFunc {
		return func(c *fuego.Context) error {
			return next(c)
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(apiDir, "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		t.Fatalf("failed to write middleware.go: %v", err)
	}

	scanner := NewScanner(appDir)
	middlewares, err := scanner.ScanMiddlewareInfo()
	if err != nil {
		t.Fatalf("ScanMiddlewareInfo failed: %v", err)
	}

	if len(middlewares) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(middlewares))
	}

	if len(middlewares) > 0 {
		if middlewares[0].Path != "/api" {
			t.Errorf("expected path /api, got %s", middlewares[0].Path)
		}
	}
}

func TestScanner_ScanMiddlewareInfo_ValidSamePackage(t *testing.T) {
	// Test with just "MiddlewareFunc" instead of "fuego.MiddlewareFunc"
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	middlewareContent := `package app

func Middleware() MiddlewareFunc {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		t.Fatalf("failed to write middleware.go: %v", err)
	}

	scanner := NewScanner(appDir)
	middlewares, err := scanner.ScanMiddlewareInfo()
	if err != nil {
		t.Fatalf("ScanMiddlewareInfo failed: %v", err)
	}

	if len(middlewares) != 1 {
		t.Errorf("expected 1 middleware for same-package type, got %d", len(middlewares))
	}
}

func TestScanner_ScanMiddlewareInfo_NoMiddleware(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// No middleware.go files

	scanner := NewScanner(appDir)
	middlewares, err := scanner.ScanMiddlewareInfo()
	if err != nil {
		t.Fatalf("ScanMiddlewareInfo failed: %v", err)
	}

	if len(middlewares) != 0 {
		t.Errorf("expected 0 middlewares, got %d", len(middlewares))
	}
}

func TestScanner_ScanMiddlewareInfo_NonExistentDir(t *testing.T) {
	scanner := NewScanner("/nonexistent/path")
	middlewares, err := scanner.ScanMiddlewareInfo()
	if err != nil {
		t.Fatalf("expected no error for non-existent dir, got: %v", err)
	}

	if len(middlewares) != 0 {
		t.Errorf("expected 0 middlewares, got %d", len(middlewares))
	}
}

func TestScanner_ScanMiddlewareInfo_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	apiDir := filepath.Join(appDir, "api")
	usersDir := filepath.Join(apiDir, "users")
	if err := os.MkdirAll(usersDir, 0755); err != nil {
		t.Fatalf("failed to create users dir: %v", err)
	}

	middlewareContent := `package placeholder

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Middleware() fuego.MiddlewareFunc {
	return nil
}
`
	// Middleware at /api
	if err := os.WriteFile(filepath.Join(apiDir, "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		t.Fatalf("failed to write api middleware.go: %v", err)
	}

	// Middleware at /api/users
	if err := os.WriteFile(filepath.Join(usersDir, "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		t.Fatalf("failed to write users middleware.go: %v", err)
	}

	scanner := NewScanner(appDir)
	middlewares, err := scanner.ScanMiddlewareInfo()
	if err != nil {
		t.Fatalf("ScanMiddlewareInfo failed: %v", err)
	}

	if len(middlewares) != 2 {
		t.Errorf("expected 2 middlewares, got %d", len(middlewares))
	}
}

func TestScanner_ScanMiddlewareInfo_InvalidSignature(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// Middleware with wrong signature (has parameters)
	middlewareContent := `package app

func Middleware(name string) MiddlewareFunc {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(appDir, "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		t.Fatalf("failed to write middleware.go: %v", err)
	}

	scanner := NewScanner(appDir)
	middlewares, err := scanner.ScanMiddlewareInfo()
	if err != nil {
		t.Fatalf("ScanMiddlewareInfo failed: %v", err)
	}

	if len(middlewares) != 0 {
		t.Errorf("expected 0 middlewares for invalid signature, got %d", len(middlewares))
	}
}

func TestScanner_ScanMiddlewareInfo_WrongReturnType(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("failed to create app dir: %v", err)
	}

	// Middleware with wrong return type
	middlewareContent := `package app

func Middleware() string {
	return ""
}
`
	if err := os.WriteFile(filepath.Join(appDir, "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		t.Fatalf("failed to write middleware.go: %v", err)
	}

	scanner := NewScanner(appDir)
	middlewares, err := scanner.ScanMiddlewareInfo()
	if err != nil {
		t.Fatalf("ScanMiddlewareInfo failed: %v", err)
	}

	if len(middlewares) != 0 {
		t.Errorf("expected 0 middlewares for wrong return type, got %d", len(middlewares))
	}
}

func TestScanner_Scan_RegistersMiddleware(t *testing.T) {
	tmpDir := t.TempDir()
	appDir := filepath.Join(tmpDir, "app")
	apiDir := filepath.Join(appDir, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		t.Fatalf("failed to create api dir: %v", err)
	}

	middlewareContent := `package api

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Middleware() fuego.MiddlewareFunc {
	return func(next fuego.HandlerFunc) fuego.HandlerFunc {
		return func(c *fuego.Context) error {
			c.SetHeader("X-Middleware", "true")
			return next(c)
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(apiDir, "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		t.Fatalf("failed to write middleware.go: %v", err)
	}

	routeContent := `package api

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

func Get(c *fuego.Context) error {
	return nil
}
`
	if err := os.WriteFile(filepath.Join(apiDir, "route.go"), []byte(routeContent), 0644); err != nil {
		t.Fatalf("failed to write route.go: %v", err)
	}

	scanner := NewScanner(appDir)
	tree := NewRouteTree()

	if err := scanner.Scan(tree); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Verify middleware was registered
	chain := tree.GetMiddlewareChain("/api")
	if len(chain) != 1 {
		t.Errorf("expected 1 middleware in chain, got %d", len(chain))
	}
}

func TestScanner_SetVerbose(t *testing.T) {
	scanner := NewScanner("app")
	scanner.SetVerbose(true)

	if !scanner.verbose {
		t.Error("expected verbose to be true")
	}

	scanner.SetVerbose(false)
	if scanner.verbose {
		t.Error("expected verbose to be false")
	}
}
