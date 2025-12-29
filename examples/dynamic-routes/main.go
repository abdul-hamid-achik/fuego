package main

import (
	"log"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func main() {
	app := fuego.New()

	// Global middleware
	app.Use(fuego.Logger())
	app.Use(fuego.Recover())

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
