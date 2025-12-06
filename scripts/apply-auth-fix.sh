#!/bin/bash
# Script to apply the Keycloak realm fix and verify the 401 error is resolved

set -e

echo "=== Applying Fix for 401 Error ==="
echo
echo "The issue: Keycloak tokens are missing the 'aud' (audience) claim."
echo "The fix: Added an audience mapper to the traveler-app client."
echo
echo "This script will:"
echo "  1. Restart Keycloak with the updated realm configuration"
echo "  2. Wait for Keycloak to be ready"
echo "  3. Get a new token with the aud claim"
echo "  4. Test the /api/offerings/specials endpoint"
echo
read -p "Press Enter to continue or Ctrl+C to cancel..."
echo

# Configuration
KC_URL="http://localhost:8081"
REALM="traveler-dev"
CLIENT_ID="traveler-app"
USERNAME="api-user"
PASSWORD="ApiUser#1!"
API_URL="http://localhost:8080"

echo "Step 1: Restarting Keycloak..."
docker compose -f docker/docker-compose.yml restart keycloak

echo "Waiting for Keycloak to start (this takes about 30 seconds)..."
sleep 10

# Wait for Keycloak to be ready
for i in {1..30}; do
    if curl -sf "${KC_URL}/realms/${REALM}" > /dev/null 2>&1; then
        echo "✓ Keycloak is ready!"
        break
    fi
    echo "  Still waiting... ($i/30)"
    sleep 2
done

if ! curl -sf "${KC_URL}/realms/${REALM}" > /dev/null 2>&1; then
    echo "✗ Keycloak did not start properly"
    echo "Check logs: docker compose -f docker/docker-compose.yml logs keycloak"
    exit 1
fi

echo

echo "Step 2: Obtaining a new token with the aud claim..."
TOKEN_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=${CLIENT_ID}" \
  -d "username=${USERNAME}" \
  -d "password=${PASSWORD}" \
  "${KC_URL}/realms/${REALM}/protocol/openid-connect/token")

if ! echo "$TOKEN_RESPONSE" | jq -e '.access_token' > /dev/null 2>&1; then
    echo "✗ Failed to get token"
    echo "$TOKEN_RESPONSE" | jq '.'
    exit 1
fi

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
echo "✓ Token obtained"

echo

echo "Step 3: Verifying token claims..."
# Extract payload (2nd part of JWT)
PAYLOAD=$(echo "$ACCESS_TOKEN" | cut -d'.' -f2)
# Decode base64 (add padding if needed)
DECODED=$(echo "$PAYLOAD" | base64 -D 2>/dev/null || echo "$PAYLOAD=" | base64 -D 2>/dev/null || echo "$PAYLOAD==" | base64 -D 2>/dev/null)

ISS=$(echo "$DECODED" | jq -r '.iss')
AUD=$(echo "$DECODED" | jq -r '.aud')
AZP=$(echo "$DECODED" | jq -r '.azp')

echo "Token claims:"
echo "  iss: $ISS"
echo "  aud: $AUD"
echo "  azp: $AZP"

if [ "$AUD" = "null" ] || [ "$AUD" = "" ]; then
    echo
    echo "⚠️  WARNING: Token still doesn't have an 'aud' claim!"
    echo "The realm may not have imported properly."
    echo
    echo "Manually add the audience mapper:"
    echo "  1. Open http://localhost:8081 and login as admin"
    echo "  2. Select realm: traveler-dev"
    echo "  3. Go to: Clients → traveler-app → Client scopes"
    echo "  4. Click on 'traveler-app-dedicated'"
    echo "  5. Add mapper → By configuration → Audience"
    echo "  6. Set: Name=audience-mapper, Included Client Audience=traveler-app"
    echo
    echo "However, the 'azp' claim should still work with the JWT middleware..."
fi

echo

echo "Step 4: Testing /api/offerings/specials endpoint..."
HTTP_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  "${API_URL}/api/offerings/specials")

HTTP_STATUS=$(echo "$HTTP_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
BODY=$(echo "$HTTP_RESPONSE" | sed '/HTTP_STATUS:/d')

echo

if [ "$HTTP_STATUS" = "200" ]; then
    echo "✅ SUCCESS! The endpoint returned 200 OK"
    echo
    echo "Response:"
    echo "$BODY" | jq '.'
    echo
    echo "🎉 The 401 error is FIXED!"
elif [ "$HTTP_STATUS" = "401" ]; then
    echo "❌ Still getting 401 Unauthorized"
    echo
    echo "Response:"
    echo "$BODY"
    echo
    echo "Troubleshooting steps:"
    echo "  1. Check application logs: tail -f logs/app.log"
    echo "  2. Verify Keycloak realm imported: docker compose -f docker/docker-compose.yml logs keycloak | grep -i import"
    echo "  3. Manually add the audience mapper in Keycloak Admin Console"
    echo "  4. See: docs/troubleshooting/FIX-401-missing-audience.md"
else
    echo "⚠️  Unexpected status: $HTTP_STATUS"
    echo "$BODY"
fi

