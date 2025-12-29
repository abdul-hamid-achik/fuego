package users

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Get handles GET /api/users
// Returns a list of all users
func Get(c *fuego.Context) error {
	users := []map[string]any{
		{"id": "1", "name": "Alice", "email": "alice@example.com"},
		{"id": "2", "name": "Bob", "email": "bob@example.com"},
		{"id": "3", "name": "Charlie", "email": "charlie@example.com"},
	}

	return c.JSON(200, map[string]any{
		"users": users,
		"count": len(users),
	})
}
