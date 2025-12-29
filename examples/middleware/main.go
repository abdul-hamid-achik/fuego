package main

import (
	"log"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func main() {
	app := fuego.New()

	// Global middleware - applies to all routes
	app.Use(fuego.Logger())
	app.Use(fuego.Recover())
	app.Use(fuego.RequestID())
	app.Use(fuego.SecureHeaders())

	// CORS for API routes
	app.Use(fuego.CORSWithConfig(fuego.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
