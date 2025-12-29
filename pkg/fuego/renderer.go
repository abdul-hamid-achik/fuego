package fuego

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
)

// Renderer handles templ component rendering with layout support.
type Renderer struct {
	// layouts stores layout components by path prefix
	layouts map[string]LayoutFunc

	// errorComponents stores error boundary components by path prefix
	errorComponents map[string]ErrorComponent

	// notFoundComponent is the global 404 component
	notFoundComponent templ.Component

	// loadingComponents stores loading skeleton components by path prefix
	loadingComponents map[string]templ.Component
}

// LayoutFunc is a function that wraps content with a layout.
// It receives the page title and returns a component that wraps children.
type LayoutFunc func(title string, children templ.Component) templ.Component

// ErrorComponent is a function that renders an error page.
type ErrorComponent func(err error) templ.Component

// LoaderFunc is a function that fetches data for a page.
type LoaderFunc func(c *Context) (any, error)

// PageHandler combines a loader and page component.
type PageHandler struct {
	Loader    LoaderFunc
	Component func(data any) templ.Component
	Title     string
}

// NewRenderer creates a new Renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		layouts:           make(map[string]LayoutFunc),
		errorComponents:   make(map[string]ErrorComponent),
		loadingComponents: make(map[string]templ.Component),
	}
}

// SetLayout registers a layout for a path prefix.
func (r *Renderer) SetLayout(pathPrefix string, layout LayoutFunc) {
	r.layouts[pathPrefix] = layout
}

// SetErrorComponent registers an error component for a path prefix.
func (r *Renderer) SetErrorComponent(pathPrefix string, errComp ErrorComponent) {
	r.errorComponents[pathPrefix] = errComp
}

// SetNotFoundComponent sets the global 404 component.
func (r *Renderer) SetNotFoundComponent(comp templ.Component) {
	r.notFoundComponent = comp
}

// SetLoadingComponent registers a loading component for a path prefix.
func (r *Renderer) SetLoadingComponent(pathPrefix string, comp templ.Component) {
	r.loadingComponents[pathPrefix] = comp
}

// GetLayout returns the most specific layout for a path.
func (r *Renderer) GetLayout(path string) LayoutFunc {
	// Find the most specific matching layout
	var bestMatch string
	var bestLayout LayoutFunc

	for prefix, layout := range r.layouts {
		if len(prefix) > len(bestMatch) && matchesPrefix(path, prefix) {
			bestMatch = prefix
			bestLayout = layout
		}
	}

	return bestLayout
}

// GetErrorComponent returns the most specific error component for a path.
func (r *Renderer) GetErrorComponent(path string) ErrorComponent {
	var bestMatch string
	var bestComp ErrorComponent

	for prefix, comp := range r.errorComponents {
		if len(prefix) > len(bestMatch) && matchesPrefix(path, prefix) {
			bestMatch = prefix
			bestComp = comp
		}
	}

	return bestComp
}

// matchesPrefix checks if path starts with prefix (with proper path boundary handling).
func matchesPrefix(path, prefix string) bool {
	if prefix == "/" || prefix == "" {
		return true
	}
	if len(path) < len(prefix) {
		return false
	}
	if path[:len(prefix)] != prefix {
		return false
	}
	// Ensure we match at a path boundary
	if len(path) > len(prefix) && path[len(prefix)] != '/' {
		return false
	}
	return true
}

// Render renders a templ component as the response.
func (r *Renderer) Render(c *Context, status int, comp templ.Component) error {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(status)
	return comp.Render(c.Context(), c.Response)
}

// RenderWithLayout renders a component wrapped in the appropriate layout.
func (r *Renderer) RenderWithLayout(c *Context, status int, title string, comp templ.Component) error {
	layout := r.GetLayout(c.Path())

	var finalComp templ.Component
	if layout != nil {
		finalComp = layout(title, comp)
	} else {
		finalComp = comp
	}

	return r.Render(c, status, finalComp)
}

// RenderError renders an error using the appropriate error component.
func (r *Renderer) RenderError(c *Context, err error) error {
	status := http.StatusInternalServerError

	// Check if it's an HTTP error
	if httpErr, ok := IsHTTPError(err); ok {
		status = httpErr.Code
	}

	errComp := r.GetErrorComponent(c.Path())
	if errComp != nil {
		return r.Render(c, status, errComp(err))
	}

	// Default error response
	return c.Error(status, err.Error())
}

// RenderNotFound renders the 404 page.
func (r *Renderer) RenderNotFound(c *Context) error {
	if r.notFoundComponent != nil {
		return r.Render(c, http.StatusNotFound, r.notFoundComponent)
	}

	return c.Error(http.StatusNotFound, "page not found")
}

// TemplComponent is a helper to render templ components directly from handlers.
func TemplComponent(c *Context, status int, comp templ.Component) error {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(status)
	return comp.Render(c.Context(), c.Response)
}

// TemplWithLayout renders a component with the given layout.
func TemplWithLayout(c *Context, status int, layout LayoutFunc, title string, comp templ.Component) error {
	var finalComp templ.Component
	if layout != nil {
		finalComp = layout(title, comp)
	} else {
		finalComp = comp
	}

	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(status)
	return finalComp.Render(c.Context(), c.Response)
}

// WrapLayout is a helper to create a layout wrapper component.
// This is useful for creating layout functions that work with templ's { children... } pattern.
type WrapLayout struct {
	Title    string
	Layout   func(title string) templ.Component
	Children templ.Component
}

// Render implements templ.Component.
func (w WrapLayout) Render(ctx context.Context, wr io.Writer) error {
	return w.Layout(w.Title).Render(ctx, wr)
}

// StreamingRenderer provides support for streaming HTML responses.
type StreamingRenderer struct {
	*Renderer
}

// NewStreamingRenderer creates a streaming-capable renderer.
func NewStreamingRenderer() *StreamingRenderer {
	return &StreamingRenderer{
		Renderer: NewRenderer(),
	}
}

// RenderStreaming renders a component with streaming support (chunked transfer).
func (sr *StreamingRenderer) RenderStreaming(c *Context, comp templ.Component) error {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.SetHeader("Transfer-Encoding", "chunked")
	c.Response.WriteHeader(http.StatusOK)

	// Flush after rendering
	if flusher, ok := c.Response.(http.Flusher); ok {
		defer flusher.Flush()
	}

	return comp.Render(c.Context(), c.Response)
}
