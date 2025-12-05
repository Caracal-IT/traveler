package handlers

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers all application routes with the Fiber app.
func RegisterRoutes(app *fiber.App) {
	// Health check endpoints
	app.Get("/ping", PingHandler)
	app.Get("/ping/simple", PingHandlerSimple)
}
