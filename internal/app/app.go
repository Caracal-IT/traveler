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

// Run starts the application. It runs a Fiber HTTP server until context is cancelled.
func Run(ctx context.Context, cfg *config.Config) error {
	sqlDb, err := initDatabase(ctx, cfg)
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	handlers.RegisterRoutes(app, cfg, sqlDb)

	errCh := make(chan error, 1)
	go startServer(app, cfg, errCh)

	select {
	case <-ctx.Done():
		return gracefulShutdown(app, sqlDb)
	case err := <-errCh:
		_ = sqlDb.Close()
		return err
	}
}

// startServer starts the Fiber HTTP server in the background and reports errors via errCh.
func startServer(app *fiber.App, cfg *config.Config, errCh chan<- error) {
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info("starting server", "address", addr)

	if err := app.Listen(addr); err != nil {
		errCh <- err
	}
}

// gracefulShutdown performs a graceful shutdown of the server and closes the database.
func gracefulShutdown(app *fiber.App, sqlDb *sql.DB) error {
	log.Info("shutting down server gracefully")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		_ = sqlDb.Close()
		return err
	}

	return sqlDb.Close()
}

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
