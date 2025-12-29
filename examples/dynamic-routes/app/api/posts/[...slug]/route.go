package posts

import (
	"strings"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Get handles GET /api/posts/*
// The [...slug] directory creates a catch-all route that matches any path
//
// Examples:
//   - GET /api/posts/hello         -> slug = "hello"
//   - GET /api/posts/2024/01/hello -> slug = "2024/01/hello"
//   - GET /api/posts/a/b/c/d       -> slug = "a/b/c/d"
func Get(c *fuego.Context) error {
	// The catch-all parameter contains the entire remaining path
	slug := c.Param("slug")

	// You can split the slug to get individual segments
	segments := strings.Split(slug, "/")

	return c.JSON(200, map[string]any{
		"slug":     slug,
		"segments": segments,
		"count":    len(segments),
		"message":  "Catch-all route matched",
	})
}
