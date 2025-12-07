package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"traveler/internal/app"
	"traveler/pkg/config"
	"traveler/pkg/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load configuration from file
	cfg := config.LoadOrDefault("configs/config.yaml")

	// Initialize logger with configured level, file path, and optional Elasticsearch sink
	if err := log.Init(cfg.Log.Level, cfg.Log.File, &cfg.Log.Elasticsearch); err != nil {
		log.Fatal("failed to init logger", "error", err)
	}
	defer func() { _ = log.Sync() }()

	logMsg := "starting application"
	logFields := []interface{}{"port", cfg.Server.Port, "log_level", cfg.Log.Level}
	if cfg.Log.File != "" {
		logFields = append(logFields, "log_file", cfg.Log.File)
	}
	log.Info(logMsg, logFields...)

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal("application error", "error", err)
	}

	// allow logger flush if any
	time.Sleep(50 * time.Millisecond)
}
