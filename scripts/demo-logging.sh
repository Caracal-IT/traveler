#!/usr/bin/env bash
set -euo pipefail

echo "====================================="
echo "Traveler Logging Level Demonstration"
echo "====================================="
echo ""

echo "1. Testing with INFO level (default config)"
echo "   - Debug messages will NOT appear"
echo "   - Info, Warn, Error will appear"
echo ""
go run ./cmd/test-logging 2>&1 | head -n 20

echo ""
echo "2. To test different log levels, update configs/config.yaml"
echo "   Available levels: debug, info, warn, error"
echo ""
echo "3. Run the server with: go run ./cmd/traveler"
echo "   The server will use the log level from configs/config.yaml"
echo ""

