package main

import (
	"log"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func main() {
	app := fuego.New()

	// App-level logger is enabled by default!
	// No need to call app.Use(fuego.Logger())

	// Global middleware
	app.Use(fuego.Recover())

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
