package main

import (
	"log"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func main() {
	// Create a new Fuego app
	app := fuego.New()

	// App-level logger is enabled by default!
	// It captures all requests with timing info.

	// Add recover middleware for graceful error handling
	app.Use(fuego.Recover())

	// Routes are automatically scanned from the app/ directory
	// Alternatively, use: fuego build --generate to create RegisterRoutes()

	// Serve static files (CSS, JS, images)
	app.Static("/static", "static")

	// Start the server
	log.Println("Starting fullstack example on http://localhost:3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
