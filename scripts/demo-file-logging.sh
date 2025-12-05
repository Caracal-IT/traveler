#!/usr/bin/env bash
set -euo pipefail

echo "========================================"
echo "Traveler File Logging Demonstration"
echo "========================================"
echo ""

# Create logs directory if it doesn't exist
mkdir -p logs

# Clean up any existing log file
rm -f logs/app.log

echo "1. Testing console-only logging (default config)"
echo "   Config: configs/config.yaml (no file specified)"
echo ""
echo "Starting server for 2 seconds..."
timeout 2 go run ./cmd/traveler 2>&1 || true
echo ""

echo "2. Testing file + console logging"
echo "   Config: configs/config.with-file.yaml"
echo ""
echo "Running file logging test..."
go run ./cmd/test-file-logging
echo ""

if [ -f logs/app.log ]; then
    echo "✓ Log file created successfully!"
    echo ""
    echo "Contents of logs/app.log:"
    echo "----------------------------------------"
    cat logs/app.log
    echo "----------------------------------------"
    echo ""
    echo "✓ File logging is working!"
else
    echo "✗ Log file was not created"
    exit 1
fi

echo ""
echo "You can now configure file logging in configs/config.yaml by setting:"
echo "  log:"
echo "    file: logs/app.log"
echo ""

