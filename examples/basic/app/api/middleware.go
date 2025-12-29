package api

import "github.com/abdul-hamid-achik/fuego/pkg/fuego"

// Middleware applies to all /api/* routes
func Middleware() fuego.MiddlewareFunc {
	return func(next fuego.HandlerFunc) fuego.HandlerFunc {
		return func(c *fuego.Context) error {
			// Add API version header
			c.SetHeader("X-API-Version", "1.0")
			return next(c)
		}
	}
}
