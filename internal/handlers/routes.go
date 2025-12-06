package handlers

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers all application routes with the Fiber app.
func RegisterRoutes(app *fiber.App) {
    // Health check endpoints
    api := app.Group("/api")
    api.Get("/ping", PingHandler)
    api.Get("/ping/simple", PingHandlerSimple)
}
