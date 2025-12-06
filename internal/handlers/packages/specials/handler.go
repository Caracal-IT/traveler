package specials

import (
	"github.com/gofiber/fiber/v2"
)

// Handler serves the authenticated specials endpoint.
// Route: GET /api/packages/specials
func Handler(c *fiber.Ctx) error {
	// Example payload; in real use this would be fetched from a service/db
	data := fiber.Map{
		"items": []fiber.Map{
			{"id": "sp-1001", "name": "Winter Escape", "price": 799.0, "currency": "USD"},
			{"id": "sp-1002", "name": "City Break Deluxe", "price": 499.0, "currency": "USD"},
		},
	}

	return c.JSON(data)
}
