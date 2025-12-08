package handlers

import (
	"database/sql"
	"traveler/internal/handlers/offerings"
	"traveler/pkg/auth"
	"traveler/pkg/config"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all application routes with the Fiber app.
func RegisterRoutes(app *fiber.App, cfg *config.Config, db *sql.DB) {
	app.Get("/", RootHandler)

	api := app.Group("/api")
	api.Get("/ping", PingHandler)
	api.Get("/ping/simple", PingHandlerSimple)

	authMW := auth.JWTMiddleware(cfg)
	offeringsGroup := api.Group("/offerings", authMW)
	offeringsGroup.Get("/specials", offerings.SpecialsHandler(db))
}
