package docs

import (
	"strings"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Get handles GET /api/docs and GET /api/docs/*
// The [[...path]] directory creates an OPTIONAL catch-all route
//
// Examples:
//   - GET /api/docs           -> path = "" (matches root)
//   - GET /api/docs/intro     -> path = "intro"
//   - GET /api/docs/api/users -> path = "api/users"
func Get(c *fuego.Context) error {
	// The optional catch-all may be empty
	path := c.Param("path")

	if path == "" {
		// Root documentation page
		return c.JSON(200, map[string]any{
			"title":   "Documentation Home",
			"message": "Welcome to the API documentation",
			"sections": []string{
				"intro",
				"getting-started",
				"api/users",
				"api/posts",
			},
		})
	}

	// Specific documentation page
	segments := strings.Split(path, "/")

	return c.JSON(200, map[string]any{
		"path":     path,
		"segments": segments,
		"title":    "Documentation: " + path,
		"message":  "Optional catch-all route matched",
	})
}
