package settings

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Get handles GET /api/settings
// The (admin) directory is a route GROUP - it doesn't affect the URL path
// This allows organizing related routes without adding URL segments
//
// File path:  app/api/(admin)/settings/route.go
// URL path:   /api/settings (NOT /api/admin/settings)
func Get(c *fuego.Context) error {
	return c.JSON(200, map[string]any{
		"settings": map[string]any{
			"theme":         "dark",
			"notifications": true,
			"language":      "en",
		},
		"message": "Route group example - (admin) folder doesn't appear in URL",
	})
}

// Put handles PUT /api/settings
func Put(c *fuego.Context) error {
	return c.JSON(200, map[string]string{
		"message": "Settings updated",
	})
}
