package main

import (
	"log"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func main() {
	app := fuego.New()

	// Add global middleware
	app.Use(fuego.Logger())
	app.Use(fuego.Recover())
	app.Use(fuego.RequestID())

	// Serve static files
	app.Static("/static", "static")

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
