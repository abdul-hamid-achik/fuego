package users

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// Sample users data
var users = []map[string]any{
	{"id": 1, "name": "Alice", "email": "alice@example.com"},
	{"id": 2, "name": "Bob", "email": "bob@example.com"},
}

// Get handles GET /api/users - List all users
func Get(c *fuego.Context) error {
	return c.JSON(200, users)
}

// Post handles POST /api/users - Create a user
func Post(c *fuego.Context) error {
	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.Bind(&user); err != nil {
		return fuego.BadRequest("invalid request body")
	}

	if user.Name == "" {
		return fuego.BadRequest("name is required")
	}

	newUser := map[string]any{
		"id":    len(users) + 1,
		"name":  user.Name,
		"email": user.Email,
	}

	users = append(users, newUser)

	return c.JSON(201, newUser)
}
