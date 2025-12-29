package user

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Get handles GET /api/users/:id
// The [id] directory creates a dynamic route segment
func Get(c *fuego.Context) error {
	// Access the dynamic parameter
	id := c.Param("id")

	// Simulated user lookup
	users := map[string]map[string]string{
		"1": {"id": "1", "name": "Alice", "email": "alice@example.com"},
		"2": {"id": "2", "name": "Bob", "email": "bob@example.com"},
		"3": {"id": "3", "name": "Charlie", "email": "charlie@example.com"},
	}

	user, exists := users[id]
	if !exists {
		return c.JSON(404, map[string]string{
			"error":   "not_found",
			"message": "User not found",
		})
	}

	return c.JSON(200, user)
}

// Put handles PUT /api/users/:id
// Updates a specific user
func Put(c *fuego.Context) error {
	id := c.Param("id")

	return c.JSON(200, map[string]string{
		"message": "User updated",
		"id":      id,
	})
}

// Delete handles DELETE /api/users/:id
// Deletes a specific user
func Delete(c *fuego.Context) error {
	id := c.Param("id")

	return c.JSON(200, map[string]string{
		"message": "User deleted",
		"id":      id,
	})
}
