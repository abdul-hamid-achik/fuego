package protected

import (
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Middleware for protected routes - requires authentication.
func Middleware() fuego.MiddlewareFunc {
	return func(next fuego.HandlerFunc) fuego.HandlerFunc {
		return func(c *fuego.Context) error {
			token := c.Header("Authorization")
			if token == "" {
				return c.JSON(401, map[string]string{
					"error":   "unauthorized",
					"message": "Authorization header required",
				})
			}

			// In a real app, validate the JWT token here
			if token != "Bearer valid-token" {
				return c.JSON(403, map[string]string{
					"error":   "forbidden",
					"message": "Invalid token",
				})
			}

			// Add user info to context
			c.Set("user_id", "user-123")
			c.Set("user_role", "admin")

			return next(c)
		}
	}
}
