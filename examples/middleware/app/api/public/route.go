package public

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// GET /api/public - no auth required
func Get(c *fuego.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "This is public data",
		"info":    "No authentication required",
	})
}
