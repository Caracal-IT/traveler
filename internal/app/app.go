package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	appdb "traveler/internal/db"
	"traveler/internal/handlers"
	"traveler/pkg/config"
	"traveler/pkg/log"
)

// initDatabase initializes the SQLite database and applies the schema.
func initDatabase(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	dbPath := cfg.Database.Path
	schemaPath := "db/schema.sql"
	sqlDb, err := appdb.Init(ctx, dbPath, schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	log.Info("database initialized", "db", dbPath)
	return sqlDb, nil
}

// Run starts the application. It runs a Fiber HTTP server until context is cancelled.
func Run(ctx context.Context, cfg *config.Config) error {
	// Initialize local SQLite database and apply schema
	sqlDb, err := initDatabase(ctx, cfg)
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		log.Debug("handling root request", "method", c.Method(), "path", c.Path())
		return c.SendString("traveler: hello\n")
	})

	// Register additional routes
	handlers.RegisterRoutes(app, cfg, sqlDb)

	// run server in the background
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
		if err := app.ShutdownWithContext(ctxShutdown); err != nil {
			_ = sqlDb.Close()
			return err
		}
		return sqlDb.Close()
	case err := <-errCh:
		_ = sqlDb.Close()
		return err
	}
}
