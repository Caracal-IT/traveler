package main

import (
	"traveler/pkg/config"
	"traveler/pkg/log"
)

func main() {
	// Test with info level (default)
	println("\n=== Testing with INFO level ===")
	cfg := config.LoadOrDefault("configs/config.yaml")
	_ = log.Init(cfg.Log.Level, cfg.Log.File)

	log.Debug("This is a DEBUG message - should NOT appear")
	log.Info("This is an INFO message - should appear")
	log.Warn("This is a WARN message - should appear")
	log.Error("This is an ERROR message - should appear")
	_ = log.Sync()

	// Test with debug level
	println("\n=== Testing with DEBUG level ===")
	cfgDebug := config.LoadOrDefault("configs/config.debug.yaml")
	_ = log.Init(cfgDebug.Log.Level, cfgDebug.Log.File)

	log.Debug("This is a DEBUG message - should appear")
	log.Info("This is an INFO message - should appear")
	log.Warn("This is a WARN message - should appear")
	log.Error("This is an ERROR message - should appear")
	_ = log.Sync()

	// Test with error level
	println("\n=== Testing with ERROR level ===")
	cfgError := config.LoadOrDefault("configs/config.error.yaml")
	_ = log.Init(cfgError.Log.Level, cfgError.Log.File)

	log.Debug("This is a DEBUG message - should NOT appear")
	log.Info("This is an INFO message - should NOT appear")
	log.Warn("This is a WARN message - should NOT appear")
	log.Error("This is an ERROR message - should appear")
	_ = log.Sync()
}
