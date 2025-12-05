package main

import (
	"fmt"
	"os"
	"traveler/pkg/config"
	"traveler/pkg/log"
)

func main() {
	// Test with file logging
	println("=== Testing with FILE logging ===")
	cfg := config.LoadOrDefault("configs/config.with-file.yaml")

	if err := log.Init(cfg.Log.Level, cfg.Log.File); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("This message goes to both stdout and file", "test", "file-logging")
	log.Debug("Debug message (won't appear at info level)")
	log.Warn("Warning message")
	log.Error("Error message")

	_ = log.Sync()

	println("\n✓ Log messages written to stdout and", cfg.Log.File)
	println("Check the log file to verify file logging works!")
}
