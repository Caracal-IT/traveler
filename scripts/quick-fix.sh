#!/bin/bash
# Quick Start - Test the 401 Fix

echo "╔════════════════════════════════════════════════════════╗"
echo "║  Traveler API - 401 Error Fix                         ║"
echo "╚════════════════════════════════════════════════════════╝"
echo
echo "This will restart Keycloak and test the specials endpoint."
echo
echo "Prerequisites:"
echo "  ✓ Docker is running"
echo "  ✓ Traveler API is running (port 8080)"
echo
echo "Press Ctrl+C to cancel, or Enter to continue..."
read

# Navigate to project root
cd "$(dirname "$0")/.." || exit 1

echo
echo "→ Restarting Keycloak with fixed realm configuration..."
docker compose -f docker/docker-compose.yml restart keycloak

echo "→ Waiting 35 seconds for Keycloak to fully start..."
for i in {1..35}; do
    printf "\r   Progress: [%-35s] %d/35s" $(printf '#%.0s' $(seq 1 $i)) $i
    sleep 1
done
echo

echo "→ Testing Keycloak availability..."
if curl -sf http://localhost:8081/realms/traveler-dev > /dev/null 2>&1; then
    echo "   ✓ Keycloak is ready"
else
    echo "   ✗ Keycloak is not responding"
    echo "   Run: docker compose -f docker/docker-compose.yml logs keycloak"
    exit 1
fi

echo
echo "→ Obtaining JWT token..."
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser%231%21" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token | jq -r .access_token)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "   ✗ Failed to get token"
    exit 1
fi
echo "   ✓ Token obtained"

echo
echo "→ Testing /api/offerings/specials..."
RESPONSE=$(curl -s -w "\nSTATUS:%{http_code}" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/offerings/specials)

STATUS=$(echo "$RESPONSE" | grep "STATUS:" | cut -d':' -f2)
BODY=$(echo "$RESPONSE" | sed '/STATUS:/d')

echo
if [ "$STATUS" = "200" ]; then
    echo "╔════════════════════════════════════════════════════════╗"
    echo "║  ✅ SUCCESS! The 401 error is FIXED!                  ║"
    echo "╚════════════════════════════════════════════════════════╝"
    echo
    echo "Response:"
    echo "$BODY" | jq . 2>/dev/null || echo "$BODY"
    echo
    echo "Next steps:"
    echo "  • Use the HTTP client: api/offerings/specials.http"
    echo "  • See documentation: docs/api/offerings/specials.md"
    echo "  • Run diagnostics: ./scripts/diagnose-auth.sh"
else
    echo "╔════════════════════════════════════════════════════════╗"
    echo "║  ❌ Still getting HTTP $STATUS                            ║"
    echo "╚════════════════════════════════════════════════════════╝"
    echo
    echo "Response: $BODY"
    echo
    echo "Troubleshooting:"
    echo "  1. Check logs: tail -f logs/app.log"
    echo "  2. Run diagnostics: ./scripts/fix-auth-401.sh"
    echo "  3. See guide: docs/troubleshooting/README.md"
    exit 1
fi

