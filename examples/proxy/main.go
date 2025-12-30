package main

import (
	"log"

	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

func main() {
	app := fuego.New()

	// App-level logger captures ALL requests including proxy actions!
	// Example output:
	//   [12:34:56] GET /v1/users → /api/users 200 in 52ms [rewrite]
	//   [12:34:57] GET /api/admin 403 in 1ms [proxy]
	//   [12:34:58] GET /old-page 301 in 2ms [redirect → /new-page]

	// Add global middleware
	app.Use(fuego.Recover())

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
