#!/usr/bin/env bash
set -euo pipefail

echo "================================================"
echo "Ping Endpoint Manual Verification"
echo "================================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "Step 1: Verify code compiles"
echo "----------------------------"
if go build ./internal/handlers/ 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Handlers package compiles"
else
    echo "✗ Compilation failed"
    exit 1
fi

if go build ./cmd/traveler 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Main application compiles"
    rm -f traveler
else
    echo "✗ Application compilation failed"
    exit 1
fi
echo ""

echo "Step 2: Check file structure"
echo "----------------------------"
files=(
    "internal/handlers/ping.go"
    "internal/handlers/ping_test.go"
    "internal/handlers/routes.go"
    "docs/ping-endpoint.md"
    "docs/ping-quick-reference.md"
    "scripts/test-ping-endpoints.sh"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        lines=$(wc -l < "$file" | tr -d ' ')
        echo -e "${GREEN}✓${NC} $file (${lines} lines)"
    else
        echo "✗ Missing: $file"
    fi
done
echo ""

echo "Step 3: Verify handler functions exist"
echo "---------------------------------------"
if grep -q "func PingHandler" internal/handlers/ping.go; then
    echo -e "${GREEN}✓${NC} PingHandler function found"
fi
if grep -q "func PingHandlerSimple" internal/handlers/ping.go; then
    echo -e "${GREEN}✓${NC} PingHandlerSimple function found"
fi
if grep -q "func RegisterRoutes" internal/handlers/routes.go; then
    echo -e "${GREEN}✓${NC} RegisterRoutes function found"
fi
echo ""

echo "Step 4: Verify test functions exist"
echo "------------------------------------"
if grep -q "func TestPingHandler" internal/handlers/ping_test.go; then
    echo -e "${GREEN}✓${NC} TestPingHandler found"
fi
if grep -q "func BenchmarkPingHandler" internal/handlers/ping_test.go; then
    echo -e "${GREEN}✓${NC} BenchmarkPingHandler found"
fi
echo ""

echo "Step 5: Code statistics"
echo "-----------------------"
echo "Handler code:"
echo "  - ping.go: $(wc -l < internal/handlers/ping.go | tr -d ' ') lines"
echo "  - ping_test.go: $(wc -l < internal/handlers/ping_test.go | tr -d ' ') lines"
echo "  - routes.go: $(wc -l < internal/handlers/routes.go | tr -d ' ') lines"
echo ""
echo "Documentation:"
echo "  - ping-endpoint.md: $(wc -l < docs/ping-endpoint.md | tr -d ' ') lines"
echo "  - ping-quick-reference.md: $(wc -l < docs/ping-quick-reference.md | tr -d ' ') lines"
echo "  - handlers/README.md: $(wc -l < internal/handlers/README.md | tr -d ' ') lines"
echo ""

echo "================================================"
echo -e "${GREEN}✓ Verification Complete${NC}"
echo "================================================"
echo ""
echo "To test the endpoints manually:"
echo "  1. Start server: go run ./cmd/traveler"
echo "  2. In another terminal:"
echo "     curl http://localhost:8080/ping"
echo "     curl http://localhost:8080/ping/simple"
echo ""
echo "To run automated tests:"
echo "  ./scripts/test-ping-endpoints.sh"

