#!/bin/bash
# Diagnostic script to test the specials endpoint authentication

set -e

echo "=== Traveler Auth Diagnostics ==="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
KC_URL="http://localhost:8081"
REALM="traveler-dev"
CLIENT_ID="traveler-app"
USERNAME="api-user"
PASSWORD="ApiUser#1!"
API_URL="http://localhost:8080"

echo "1. Testing Keycloak accessibility..."
if curl -sf "${KC_URL}/realms/${REALM}" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Keycloak is accessible${NC}"
else
    echo -e "${RED}✗ Keycloak is NOT accessible at ${KC_URL}${NC}"
    echo "  Please start Keycloak: docker compose -f docker/docker-compose.yml up -d keycloak"
    exit 1
fi

echo

echo "2. Testing Traveler API accessibility..."
if curl -sf "${API_URL}/api/ping" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Traveler API is accessible${NC}"
else
    echo -e "${RED}✗ Traveler API is NOT accessible at ${API_URL}${NC}"
    echo "  Please start the API: go run ./cmd/traveler"
    exit 1
fi

echo

echo "3. Requesting access token from Keycloak..."
TOKEN_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=${CLIENT_ID}" \
  -d "username=${USERNAME}" \
  -d "password=${PASSWORD}" \
  "${KC_URL}/realms/${REALM}/protocol/openid-connect/token")

if echo "$TOKEN_RESPONSE" | jq -e '.access_token' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Successfully obtained access token${NC}"
    ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

    # Decode token to show claims (header and payload only)
    echo
    echo "Token claims:"
    echo "$ACCESS_TOKEN" | cut -d'.' -f2 | base64 -d 2>/dev/null | jq '.'
else
    echo -e "${RED}✗ Failed to obtain access token${NC}"
    echo "Response from Keycloak:"
    echo "$TOKEN_RESPONSE" | jq '.'
    exit 1
fi

echo

echo "4. Testing /api/offerings/specials endpoint WITH token..."
SPECIALS_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  "${API_URL}/api/offerings/specials")

HTTP_STATUS=$(echo "$SPECIALS_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
BODY=$(echo "$SPECIALS_RESPONSE" | sed '/HTTP_STATUS:/d')

if [ "$HTTP_STATUS" = "200" ]; then
    echo -e "${GREEN}✓ Success! Endpoint returned 200${NC}"
    echo "Response:"
    echo "$BODY" | jq '.'
elif [ "$HTTP_STATUS" = "401" ]; then
    echo -e "${RED}✗ Endpoint returned 401 Unauthorized${NC}"
    echo "Response:"
    echo "$BODY"
    echo
    echo "Possible causes:"
    echo "  - Token issuer mismatch (check config.yaml auth.issuer)"
    echo "  - Token audience mismatch (check config.yaml auth.audience)"
    echo "  - Keycloak JWKS endpoint not accessible from API"
    echo "  - Token expired"
else
    echo -e "${YELLOW}! Endpoint returned status ${HTTP_STATUS}${NC}"
    echo "Response:"
    echo "$BODY"
fi

echo

echo "5. Testing /api/offerings/specials endpoint WITHOUT token..."
NO_TOKEN_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  "${API_URL}/api/offerings/specials")

HTTP_STATUS=$(echo "$NO_TOKEN_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)

if [ "$HTTP_STATUS" = "401" ]; then
    echo -e "${GREEN}✓ Correctly returns 401 without token${NC}"
else
    echo -e "${YELLOW}! Expected 401 but got ${HTTP_STATUS}${NC}"
fi

echo
echo "=== Diagnostics Complete ==="

