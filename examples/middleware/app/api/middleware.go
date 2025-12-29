package api

import (
	"time"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Middleware applies to all routes under /api/*
// This demonstrates route-level middleware inheritance
func Middleware(next fuego.HandlerFunc) fuego.HandlerFunc {
	return func(c *fuego.Context) error {
		// Add API version header
		c.SetHeader("X-API-Version", "1.0")

		// Add timing
		start := time.Now()
		err := next(c)
		c.SetHeader("X-Response-Time", time.Since(start).String())

		return err
	}
}
