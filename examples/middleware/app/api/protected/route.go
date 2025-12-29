package protected

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// GET /api/protected - requires valid auth token
func Get(c *fuego.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "This is protected data",
		"user": map[string]string{
			"id":   c.GetString("user_id"),
			"role": c.GetString("user_role"),
		},
	})
}
