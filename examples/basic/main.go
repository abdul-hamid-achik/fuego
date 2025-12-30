package main

import (
	"log"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func main() {
	app := fuego.New()

	// App-level logger is enabled by default and captures ALL requests,
	// including those handled by the proxy layer.
	// Customize with: app.SetLogger(fuego.RequestLoggerConfig{...})

	// Add global middleware
	app.Use(fuego.Recover())
	app.Use(fuego.RequestID())

	// Serve static files
	app.Static("/static", "static")

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
