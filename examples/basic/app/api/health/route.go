package health

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// Get handles GET /api/health
func Get(c *fuego.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}
