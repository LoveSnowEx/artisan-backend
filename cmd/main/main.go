package main

import (
	"artisan-backend/internal/service"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	src := service.New()

	app.Get("/api/forward/:px", src.Forward)
	app.Get("/api/left/:deg", src.Left)
	app.Get("/api/right/:deg", src.Right)
	app.Get("/api/reverse", src.Reverse)
	app.Get("/api/inversion", src.Inversion)
	app.Get("/api/circulars", src.Circulars)
	app.Get("/api/img", src.Img)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
