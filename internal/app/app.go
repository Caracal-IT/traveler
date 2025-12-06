package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"traveler/internal/handlers"
	"traveler/pkg/config"
	"traveler/pkg/log"
)

// Run starts the application. It runs a Fiber HTTP server until context is cancelled.
func Run(ctx context.Context, cfg *config.Config) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		log.Debug("handling root request", "method", c.Method(), "path", c.Path())
		return c.SendString("traveler: hello\n")
	})

	// Register additional routes
	handlers.RegisterRoutes(app, cfg)

	// run server in background
	errCh := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("starting server", "address", addr)
		if err := app.Listen(addr); err != nil {
			errCh <- err
		}
	}()

	// wait for context done or server error
	select {
	case <-ctx.Done():
		log.Info("shutting down server gracefully")
		// graceful shutdown with timeout
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return app.ShutdownWithContext(ctxShutdown)
	case err := <-errCh:
		return err
	}
}
