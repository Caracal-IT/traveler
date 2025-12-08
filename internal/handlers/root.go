package handlers

import (
	"traveler/pkg/log"

	"github.com/gofiber/fiber/v2"
)

// RootHandler handles requests to the root path "/".
func RootHandler(c *fiber.Ctx) error {
	log.Debug("handling root request", "method", c.Method(), "path", c.Path())
	return c.SendString("traveler: hello\n")
}
