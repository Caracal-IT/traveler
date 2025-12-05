package handlers

import (
	"time"

	"traveler/pkg/log"

	"github.com/gofiber/fiber/v2"
)

// PingResponse represents the response structure for the ping endpoint.
type PingResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// PingHandler handles the /ping endpoint for health checks.
// It returns a simple response indicating the service is alive and responsive.
//
// @Summary      Health check endpoint
// @Description  Returns service health status
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  PingResponse
// @Router       /ping [get]
func PingHandler(c *fiber.Ctx) error {
	log.Debug("ping endpoint called", "ip", c.IP(), "user_agent", c.Get("User-Agent"))

	response := PingResponse{
		Status:    "ok",
		Message:   "pong",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// PingHandlerSimple is a minimal ping handler that returns plain text.
// Useful for load balancers that expect a simple text response.
//
// @Summary      Simple health check
// @Description  Returns simple text response
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "pong"
// @Router       /ping/simple [get]
func PingHandlerSimple(c *fiber.Ctx) error {
	log.Debug("simple ping endpoint called", "ip", c.IP())
	return c.Status(fiber.StatusOK).SendString("pong")
}
