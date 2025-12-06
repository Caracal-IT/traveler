# 401 Error Fix - Missing Audience Claim

## Problem Identified
The Keycloak JWT tokens issued for the `traveler-app` client are missing the `aud` (audience) claim. The token only contains the `azp` (authorized party) claim.

While the JWT middleware code does check for `azp` as a fallback, explicitly adding the audience claim is the correct solution.

## Solution Applied
Added an audience protocol mapper to the `traveler-app` client in the Keycloak realm configuration.

## Steps to Apply the Fix

### 1. The realm export has been updated
The file `docker/keycloak/realm-export/traveler-dev-realm.json` now includes an audience mapper for the `traveler-app` client.

### 2. Restart Keycloak to import the updated realm
```bash
cd /Users/ettienemare/GolandProjects/traveler
docker compose -f docker/docker-compose.yml restart keycloak
```

Wait about 30 seconds for Keycloak to fully restart and import the realm.

### 3. Test the fix
```bash
# Get a new token
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser#1!" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token | jq -r .access_token)

# Verify the token now has an aud claim
echo $TOKEN | cut -d'.' -f2 | base64 -D 2>/dev/null | jq '{iss, aud, azp}'

# Test the specials endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/offerings/specials | jq .
```

Expected token claims after the fix:
```json
{
  "iss": "http://localhost:8081/realms/traveler-dev",
  "aud": "traveler-app",
  "azp": "traveler-app"
}
```

### 4. If the issue persists
If you still get 401 errors after restarting Keycloak:

1. **Verify Keycloak imported the realm properly:**
   ```bash
   docker compose -f docker/docker-compose.yml logs keycloak | grep -i "import"
   ```

2. **Manually add the audience mapper via Keycloak Admin Console:**
   - Open http://localhost:8081
   - Login with admin credentials (from `docker/keycloak.user.env`)
   - Select realm: `traveler-dev`
   - Go to: Clients → `traveler-app` → Client scopes → `traveler-app-dedicated` → Add mapper → By configuration
   - Choose: "Audience"
   - Set:
     - Name: `audience-mapper`
     - Included Client Audience: `traveler-app`
     - Add to access token: ON
   - Save

3. **Get a fresh token** (old tokens won't have the new claim)

4. **Check application logs:**
   ```bash
   tail -f logs/app.log
   ```
   Look for any JWT validation error messages.

## Alternative: If you can't restart Keycloak

If the JWT middleware is correctly checking the `azp` claim (which it should be based on the code), you can verify by checking the application logs to see which specific validation is failing.

The middleware checks audience in this order:
1. `aud` claim (string or array)
2. `azp` claim (authorized party)
3. `resource_access` map

If `azp` matching isn't working, there may be a bug in the audience validation logic.

## Manual Workaround (Temporary)

If you need an immediate workaround and can't wait for Keycloak restart:

Update the audience check to be case-insensitive or more lenient in `pkg/auth/jwt.go`:
```go
// Around line 130
if azp, ok := claims["azp"].(string); ok && strings.EqualFold(azp, audience) {
    audOK = true
}
```

This is already implemented, so the issue is likely that Keycloak needs to be restarted with the new realm configuration.

