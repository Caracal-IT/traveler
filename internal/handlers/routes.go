package handlers

import (
	offerings "traveler/internal/handlers/offerings"
	"traveler/pkg/auth"
	"traveler/pkg/config"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all application routes with the Fiber app.
func RegisterRoutes(app *fiber.App, cfg *config.Config) {
	// Health check endpoints
	api := app.Group("/api")
	api.Get("/ping", PingHandler)
	api.Get("/ping/simple", PingHandlerSimple)

	// Authenticated routes
	authMW := auth.JWTMiddleware(cfg)
	offeringsGroup := api.Group("/offerings", authMW)
	offeringsGroup.Get("/specials", offerings.SpecialsHandler)
}
