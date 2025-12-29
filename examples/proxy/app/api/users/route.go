package users

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// GET /api/users
func Get(c *fuego.Context) error {
	users := []map[string]interface{}{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
	}
	return c.JSON(200, map[string]interface{}{
		"users": users,
		"proxy": map[string]string{
			"version": c.Header("X-Proxy-Version"),
			"path":    c.Header("X-Request-Path"),
		},
	})
}
