package admin

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// GET /api/admin - Protected by proxy
func Get(c *fuego.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "Welcome to admin panel",
		"user":    "admin",
	})
}
