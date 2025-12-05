#!/usr/bin/env bash
set -euo pipefail

echo "========================================"
echo "Ping Endpoint Integration Test"
echo "========================================"
echo ""

# Start the server in background
echo "Starting server..."
go run ./cmd/traveler &
SERVER_PID=$!

# Give server time to start
sleep 2

echo "Testing ping endpoints..."
echo ""

# Test JSON ping endpoint
echo "1. Testing /ping (JSON response):"
PING_RESPONSE=$(curl -s http://localhost:8080/ping)
echo "   Response: $PING_RESPONSE"

if echo "$PING_RESPONSE" | grep -q '"status":"ok"'; then
    echo "   ✓ JSON ping endpoint working"
else
    echo "   ✗ JSON ping endpoint failed"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo ""

# Test simple ping endpoint
echo "2. Testing /ping/simple (plain text response):"
SIMPLE_RESPONSE=$(curl -s http://localhost:8080/ping/simple)
echo "   Response: $SIMPLE_RESPONSE"

if [ "$SIMPLE_RESPONSE" = "pong" ]; then
    echo "   ✓ Simple ping endpoint working"
else
    echo "   ✗ Simple ping endpoint failed"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo ""

# Test response time
echo "3. Testing response time:"
TIME_RESPONSE=$(curl -s -o /dev/null -w "%{time_total}" http://localhost:8080/ping)
echo "   Response time: ${TIME_RESPONSE}s"

if (( $(echo "$TIME_RESPONSE < 0.1" | bc -l) )); then
    echo "   ✓ Response time acceptable"
else
    echo "   ⚠ Response time slower than expected"
fi
echo ""

# Clean up
echo "Stopping server..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

echo ""
echo "========================================"
echo "✓ All ping endpoint tests passed!"
echo "========================================"

