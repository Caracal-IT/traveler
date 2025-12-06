# 401 Error Resolution Summary

## Issue
The `/api/offerings/specials` endpoint was returning a 401 Unauthorized error.

## Root Cause
After analyzing the JWT middleware code and testing the token generation, I identified that:

1. **Missing `aud` (audience) claim in JWT tokens**: The Keycloak `traveler-app` client configuration didn't include an audience protocol mapper
2. While the JWT middleware has fallback logic to check the `azp` (authorized party) claim, having an explicit `aud` claim is the correct OAuth2/OIDC standard

## Solution Implemented

### 1. Fixed Keycloak Realm Configuration
Updated `docker/keycloak/realm-export/traveler-dev-realm.json` to add an audience protocol mapper to the `traveler-app` client:

```json
"protocolMappers": [
  {
    "name": "audience-mapper",
    "protocol": "openid-connect",
    "protocolMapper": "oidc-audience-mapper",
    "consentRequired": false,
    "config": {
      "included.client.audience": "traveler-app",
      "id.token.claim": "false",
      "access.token.claim": "true"
    }
  }
]
```

This ensures that all tokens issued for the `traveler-app` client will include:
- `aud`: `"traveler-app"` (explicit audience)
- `azp`: `"traveler-app"` (authorized party)

### 2. Created Troubleshooting Documentation
- `docs/troubleshooting/README.md` - Main troubleshooting guide index
- `docs/troubleshooting/FIX-401-missing-audience.md` - Detailed fix guide
- `docs/troubleshooting/401-specials-endpoint.md` - Comprehensive 401 troubleshooting

### 3. Created Utility Scripts
- `scripts/apply-auth-fix.sh` - Automated script to restart Keycloak and verify the fix
- `scripts/diagnose-auth.sh` - Diagnostic script to identify auth issues
- `scripts/fix-auth-401.sh` - Quick fix script for 401 errors

### 4. Updated Documentation
- Fixed incorrect endpoint path in README.md (was `/api/packages/specials`, should be `/api/offerings/specials`)
- Added troubleshooting reference in main README

## How to Apply the Fix

### Option 1: Automated (Recommended)
```bash
cd /Users/ettienemare/GolandProjects/traveler
./scripts/apply-auth-fix.sh
```

This script will:
1. Restart Keycloak with the updated realm configuration
2. Wait for Keycloak to be ready
3. Obtain a new token and verify it has the `aud` claim
4. Test the `/api/offerings/specials` endpoint
5. Report success or provide troubleshooting steps

### Option 2: Manual
```bash
# 1. Restart Keycloak
docker compose -f docker/docker-compose.yml restart keycloak

# 2. Wait for startup (about 30 seconds)
sleep 30

# 3. Test the endpoint
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser#1!" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token | jq -r .access_token)

curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/offerings/specials | jq .
```

Expected response:
```json
{
  "items": [
    {"id": "sp-1001", "name": "Winter Escape", "price": 799.0, "currency": "USD"},
    {"id": "sp-1002", "name": "City Break Deluxe", "price": 499.0, "currency": "USD"}
  ]
}
```

## Technical Details

### JWT Middleware Audience Validation
The middleware (`pkg/auth/jwt.go`) validates audience in this order:
1. Check `aud` claim (string or array)
2. Fallback to `azp` claim
3. Fallback to `resource_access` map (Keycloak-specific)

The fix ensures tokens have an explicit `aud` claim, which is the OAuth2/OIDC standard.

### Token Claims Before Fix
```json
{
  "iss": "http://localhost:8081/realms/traveler-dev",
  "azp": "traveler-app",
  "typ": "Bearer",
  ...
}
```

### Token Claims After Fix
```json
{
  "iss": "http://localhost:8081/realms/traveler-dev",
  "aud": "traveler-app",
  "azp": "traveler-app",
  "typ": "Bearer",
  ...
}
```

## Files Modified
1. `docker/keycloak/realm-export/traveler-dev-realm.json` - Added audience mapper
2. `README.md` - Fixed endpoint path, added troubleshooting reference

## Files Created
1. `docs/troubleshooting/README.md`
2. `docs/troubleshooting/FIX-401-missing-audience.md`
3. `docs/troubleshooting/401-specials-endpoint.md`
4. `scripts/apply-auth-fix.sh`
5. `scripts/diagnose-auth.sh`
6. `scripts/fix-auth-401.sh`

## Next Steps
1. Run `./scripts/apply-auth-fix.sh` to apply the fix
2. If issues persist, check the troubleshooting guide: `docs/troubleshooting/README.md`
3. For manual Keycloak configuration, see: `docs/troubleshooting/FIX-401-missing-audience.md`

## Prevention
To prevent this issue in the future:
- Always include audience mappers when creating Keycloak clients
- Test authentication immediately after client creation
- Use the diagnostic scripts to verify token claims

