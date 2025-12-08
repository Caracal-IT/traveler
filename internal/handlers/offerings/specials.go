package offerings

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"

	repo "traveler/internal/db/offerings"
	"traveler/pkg/log"
)

// SpecialsHandler returns a Fiber handler that serves the authenticated specials endpoint.
// Route: GET /api/offerings/specials
func SpecialsHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		if ctx == nil {
			ctx = context.Background()
		}

		items, err := repo.GetActiveSpecials(ctx, db)
		if err != nil {
			log.Error("failed to list specials", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch specials",
			})
		}

		// Keep response shape stable: { "items": [ ... ] }
		// Directly marshal the repo.Special values.
		return c.JSON(fiber.Map{"items": items})
	}
}
