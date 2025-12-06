#!/bin/bash
# Recreate Keycloak with the updated realm configuration

set -e

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║  Recreating Keycloak to Fix 401 Error                        ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo
echo "This will:"
echo "  1. Stop and remove the existing Keycloak container"
echo "  2. Start a fresh Keycloak with the updated realm configuration"
echo "  3. Wait for Keycloak to fully start"
echo "  4. Test the /api/offerings/specials endpoint"
echo
echo "⚠️  This will reset any manual changes made in Keycloak Admin Console"
echo
read -p "Press Enter to continue or Ctrl+C to cancel..."
echo

cd "$(dirname "$0")/.."

echo "→ Stopping and removing Keycloak container..."
docker compose -f docker/docker-compose.yml stop keycloak
docker compose -f docker/docker-compose.yml rm -f keycloak

echo "→ Starting fresh Keycloak with updated realm..."
docker compose -f docker/docker-compose.yml up -d keycloak

echo "→ Waiting for Keycloak to start (this takes 40-60 seconds)..."
for i in {1..60}; do
    if curl -sf http://localhost:8081/realms/traveler-dev > /dev/null 2>&1; then
        echo "   ✓ Keycloak is ready! (took ${i} seconds)"
        break
    fi
    printf "\r   Waiting... %d/60 seconds" $i
    sleep 1
done
echo

if ! curl -sf http://localhost:8081/realms/traveler-dev > /dev/null 2>&1; then
    echo "   ✗ Keycloak did not start properly"
    echo "   Check logs: docker compose -f docker/docker-compose.yml logs keycloak"
    exit 1
fi

echo
echo "→ Obtaining fresh token..."
TOKEN_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser#1!" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token)

if ! echo "$TOKEN_RESPONSE" | jq -e '.access_token' > /dev/null 2>&1; then
    echo "   ✗ Failed to get token"
    echo "$TOKEN_RESPONSE" | jq '.'
    exit 1
fi

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
echo "   ✓ Token obtained"

# Decode and show token claims
PAYLOAD=$(echo "$ACCESS_TOKEN" | cut -d'.' -f2)
DECODED=$(echo "$PAYLOAD=" | base64 -D 2>/dev/null || echo "$PAYLOAD==" | base64 -D 2>/dev/null)
AUD=$(echo "$DECODED" | jq -r '.aud // empty')
AZP=$(echo "$DECODED" | jq -r '.azp // empty')

echo
echo "   Token claims:"
echo "     aud: $AUD"
echo "     azp: $AZP"

if [ -n "$AUD" ] && [ "$AUD" != "null" ]; then
    echo "   ✓ Token has 'aud' claim!"
else
    echo "   ⚠️  Token still missing 'aud' claim (but 'azp' should work)"
fi

echo
echo "→ Testing /api/offerings/specials endpoint..."
HTTP_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  http://localhost:8080/api/offerings/specials)

HTTP_STATUS=$(echo "$HTTP_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
BODY=$(echo "$HTTP_RESPONSE" | sed '/HTTP_STATUS:/d')

echo
if [ "$HTTP_STATUS" = "200" ]; then
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║  ✅ SUCCESS! The 401 error is FIXED!                         ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo
    echo "Response:"
    echo "$BODY" | jq . 2>/dev/null || echo "$BODY"
    echo
    echo "✨ You can now use the HTTP client file:"
    echo "   api/offerings/specials.http"
    echo
    echo "Just run request #1 to get a token, then request #2 to call the endpoint!"
elif [ "$HTTP_STATUS" = "401" ]; then
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║  ❌ Still getting 401 - Manual configuration needed          ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo
    echo "The realm import didn't include the audience mapper."
    echo "You need to add it manually in Keycloak Admin Console:"
    echo
    echo "1. Open: http://localhost:8081"
    echo "2. Login: admin / Admin#1!"
    echo "3. Select realm: traveler-dev"
    echo "4. Go to: Clients → traveler-app → Client scopes"
    echo "5. Click: traveler-app-dedicated"
    echo "6. Click: Add mapper → By configuration → Audience"
    echo "7. Configure:"
    echo "   - Name: audience-mapper"
    echo "   - Included Client Audience: traveler-app"
    echo "   - Add to ID token: OFF"
    echo "   - Add to access token: ON"
    echo "8. Save"
    echo
    echo "Then try running this script again to test."
    exit 1
else
    echo "⚠️  Unexpected status: $HTTP_STATUS"
    echo "$BODY"
    exit 1
fi

