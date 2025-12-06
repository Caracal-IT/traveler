#!/bin/bash
# Quick fix script for 401 errors on the specials endpoint

set -e

echo "=== Fixing 401 Error for /api/offerings/specials ==="
echo

# Configuration
KC_URL="http://localhost:8081"
REALM="traveler-dev"
CLIENT_ID="traveler-app"
USERNAME="api-user"
PASSWORD="ApiUser#1!"
API_URL="http://localhost:8080"

# Check and start services
echo "Step 1: Checking services..."

# Check Keycloak
if ! curl -sf "${KC_URL}/realms/${REALM}" > /dev/null 2>&1; then
    echo "Keycloak is not running. Starting it..."
    docker compose -f docker/docker-compose.yml up -d keycloak
    echo "Waiting for Keycloak to start (30 seconds)..."
    sleep 30

    if curl -sf "${KC_URL}/realms/${REALM}" > /dev/null 2>&1; then
        echo "✓ Keycloak is now running"
    else
        echo "✗ Failed to start Keycloak. Please check docker logs:"
        echo "  docker compose -f docker/docker-compose.yml logs keycloak"
        exit 1
    fi
else
    echo "✓ Keycloak is already running"
fi

# Check Traveler API
if ! curl -sf "${API_URL}/api/ping" > /dev/null 2>&1; then
    echo "Traveler API is not responding. You may need to start it:"
    echo "  docker compose -f docker/docker-compose.yml up -d traveler"
    echo "  OR: go run ./cmd/traveler"
    exit 1
else
    echo "✓ Traveler API is running"
fi

echo

echo "Step 2: Testing token acquisition..."
TOKEN_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=${CLIENT_ID}" \
  -d "username=${USERNAME}" \
  -d "password=${PASSWORD}" \
  "${KC_URL}/realms/${REALM}/protocol/openid-connect/token")

if ! echo "$TOKEN_RESPONSE" | jq -e '.access_token' > /dev/null 2>&1; then
    echo "✗ Failed to get token from Keycloak"
    echo "Response:"
    echo "$TOKEN_RESPONSE" | jq '.'
    echo
    echo "Possible issues:"
    echo "  - User 'api-user' doesn't exist in Keycloak"
    echo "  - Wrong password"
    echo "  - Client 'traveler-app' not configured properly"
    echo
    echo "Try accessing Keycloak admin console:"
    echo "  http://localhost:8081"
    exit 1
fi

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
echo "✓ Successfully obtained access token"

# Decode and display token claims
echo
echo "Token claims:"
CLAIMS=$(echo "$ACCESS_TOKEN" | cut -d'.' -f2 | base64 -d 2>/dev/null | jq '{iss, aud, azp, exp}')
echo "$CLAIMS"

# Extract and check claims
ISS=$(echo "$CLAIMS" | jq -r '.iss')
AUD=$(echo "$CLAIMS" | jq -r '.aud')
AZP=$(echo "$CLAIMS" | jq -r '.azp')

echo

echo "Step 3: Validating configuration..."

# Check config file
if [ ! -f configs/config.yaml ]; then
    echo "✗ Config file not found: configs/config.yaml"
    exit 1
fi

CFG_ISSUER=$(grep "issuer:" configs/config.yaml | awk '{print $2}' | tr -d '"')
CFG_AUDIENCE=$(grep "audience:" configs/config.yaml | awk '{print $2}' | tr -d '"')

echo "Config issuer:   $CFG_ISSUER"
echo "Token issuer:    $ISS"

if [ "$CFG_ISSUER" != "$ISS" ]; then
    echo "✗ ISSUER MISMATCH!"
    echo "  Update configs/config.yaml to set issuer to: $ISS"
    exit 1
else
    echo "✓ Issuer matches"
fi

echo
echo "Config audience: $CFG_AUDIENCE"
echo "Token aud:       $AUD"
echo "Token azp:       $AZP"

# Check if audience matches
AUDIENCE_MATCH=false
if [ "$CFG_AUDIENCE" = "$AUD" ] || [ "$CFG_AUDIENCE" = "$AZP" ]; then
    AUDIENCE_MATCH=true
fi

# Also check if aud is an array
if echo "$TOKEN_RESPONSE" | jq -e ".access_token" > /dev/null 2>&1; then
    FULL_AUD=$(echo "$ACCESS_TOKEN" | cut -d'.' -f2 | base64 -d 2>/dev/null | jq -r '.aud')
    if echo "$FULL_AUD" | jq -e "index(\"$CFG_AUDIENCE\")" > /dev/null 2>&1; then
        AUDIENCE_MATCH=true
    fi
fi

if [ "$AUDIENCE_MATCH" = "false" ]; then
    echo "✗ AUDIENCE MISMATCH!"
    echo "  The token doesn't contain the expected audience: $CFG_AUDIENCE"
    echo "  You may need to add an audience mapper in Keycloak for client '$CLIENT_ID'"
    exit 1
else
    echo "✓ Audience matches"
fi

echo

echo "Step 4: Testing the specials endpoint..."
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  "${API_URL}/api/offerings/specials")

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS:/d')

if [ "$HTTP_STATUS" = "200" ]; then
    echo "✓ SUCCESS! Endpoint returned 200 OK"
    echo
    echo "Response:"
    echo "$BODY" | jq '.'
    echo
    echo "The issue is resolved! 🎉"
elif [ "$HTTP_STATUS" = "401" ]; then
    echo "✗ Still getting 401 Unauthorized"
    echo
    echo "Response:"
    echo "$BODY"
    echo
    echo "Check application logs for more details:"
    echo "  tail -f logs/app.log"
    echo
    echo "Look for messages like:"
    echo "  - 'token validation failed'"
    echo "  - 'token issuer mismatch'"
    echo "  - 'token audience/azp mismatch'"
    echo "  - 'failed to get JWKS'"
    exit 1
else
    echo "Unexpected status: $HTTP_STATUS"
    echo "Response:"
    echo "$BODY"
    exit 1
fi

