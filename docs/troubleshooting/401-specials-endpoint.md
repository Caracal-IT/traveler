# Troubleshooting 401 Unauthorized Error for /api/offerings/specials

## Problem
The `/api/offerings/specials` endpoint returns a 401 Unauthorized error even with what appears to be a valid token.

## Root Causes and Solutions

### 1. Keycloak Not Running
**Symptom**: Cannot obtain a token from Keycloak  
**Solution**:
```bash
cd /Users/ettienemare/GolandProjects/traveler
docker compose -f docker/docker-compose.yml up -d keycloak
```

Wait for Keycloak to fully start (takes ~30-60 seconds). Verify with:
```bash
curl http://localhost:8081/realms/traveler-dev
```

### 2. Traveler API Not Running
**Symptom**: No response from the API  
**Solution**:
```bash
# If running via docker:
docker compose -f docker/docker-compose.yml up -d traveler

# OR if running directly:
go run ./cmd/traveler
```

### 3. Token Issuer Mismatch
**Symptom**: 401 error with log message "token issuer mismatch"  
**Root Cause**: The `iss` claim in the JWT doesn't match `configs/config.yaml` auth.issuer

**Check token issuer**:
```bash
# Get a token
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password&client_id=traveler-app&username=api-user&password=ApiUser%231%21" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token | jq -r .access_token)

# Decode and view claims
echo $TOKEN | cut -d'.' -f2 | base64 -d 2>/dev/null | jq '{iss, aud, azp}'
```

**Expected issuer**: `http://localhost:8081/realms/traveler-dev`

**Solution**: Ensure `configs/config.yaml` has:
```yaml
auth:
  issuer: http://localhost:8081/realms/traveler-dev
  audience: traveler-app
```

### 4. Token Audience Mismatch
**Symptom**: 401 error with log message "token audience/azp mismatch"  
**Root Cause**: The token doesn't contain the expected audience

The JWT middleware checks for the audience in this order:
1. `aud` claim (string or array)
2. `azp` claim (authorized party)
3. `resource_access` map (Keycloak-specific)

**Check token audience**:
```bash
echo $TOKEN | cut -d'.' -f2 | base64 -d 2>/dev/null | jq '{aud, azp, resource_access}'
```

**Expected**: One of these should match `traveler-app`:
- `aud`: "traveler-app" or ["traveler-app", ...]
- `azp`: "traveler-app"
- `resource_access.traveler-app`: {...}

**Solution**: The client configuration in Keycloak realm export looks correct. If the audience is still wrong, you may need to:

1. Add audience mapper to the client in Keycloak UI
2. Or update `configs/config.yaml` to match the actual audience in tokens

### 5. JWKS Endpoint Not Accessible
**Symptom**: 401 error with log message "failed to get JWKS"  
**Root Cause**: The API cannot fetch public keys from Keycloak to verify token signatures

The API tries to fetch keys from: `http://localhost:8081/realms/traveler-dev/protocol/openid-connect/certs`

**Test JWKS endpoint**:
```bash
curl http://localhost:8081/realms/traveler-dev/protocol/openid-connect/certs
```

**Solution**: Ensure Keycloak is running and accessible from the API container/process.

### 6. Token Expired
**Symptom**: 401 error shortly after obtaining token  
**Root Cause**: Tokens have a short lifespan (default 5 minutes in Keycloak)

**Solution**: Get a fresh token before each request, or implement token refresh logic.

### 7. Wrong Keycloak User Credentials
**Symptom**: Cannot obtain token - error response from Keycloak  
**Root Cause**: Using wrong username/password

**Solution**: Use credentials from the realm export:
- Username: `api-user`
- Password: `ApiUser#1!`

## Quick Test Script

Use the diagnostic script to identify the issue:
```bash
./scripts/diagnose-auth.sh
```

## Manual Test

1. **Get a token**:
```bash
ACCESS_TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser#1!" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token | jq -r .access_token)

echo "Token: $ACCESS_TOKEN"
```

2. **Call the endpoint**:
```bash
curl -v -H "Authorization: Bearer $ACCESS_TOKEN" \
  http://localhost:8080/api/offerings/specials
```

3. **Check application logs** for specific error messages:
```bash
tail -f logs/app.log
```

Look for messages like:
- "token validation failed"
- "token issuer mismatch"
- "token audience/azp mismatch"
- "failed to get JWKS"

## Using the HTTP Client

If using the HTTP request file `api/offerings/specials.http`:

1. Ensure environment is set to `dev` in your IDE
2. Run request #1 to get the token (it will be saved automatically)
3. Run request #2 to call the specials endpoint

If request #1 fails:
- Check that Keycloak is running
- Verify credentials in `api/http-client.env.json`

If request #2 fails with 401:
- Check application logs
- Verify token claims match config
- Ensure JWKS endpoint is accessible

## Common Configuration Issues

### Issue: Docker container cannot reach host Keycloak
If running the API in Docker and Keycloak on host:
- Change `localhost` to `host.docker.internal` in configs
- Or run both in the same Docker network

### Issue: HTTP vs HTTPS mismatch  
The middleware is flexible with http/https in development, but ensure consistency:
- Keycloak issuer URL scheme matches token `iss` claim
- API config matches Keycloak realm URL

## Still Not Working?

1. Restart both services:
```bash
docker compose -f docker/docker-compose.yml restart
```

2. Check if services are actually listening:
```bash
lsof -i :8080  # Traveler API
lsof -i :8081  # Keycloak
```

3. Enable debug logging in `configs/config.yaml`:
```yaml
log:
  level: debug
```

4. Check Keycloak admin console:
- Open: http://localhost:8081
- Login with admin credentials (from `docker/keycloak.user.env`)
- Verify realm `traveler-dev` exists
- Verify client `traveler-app` is configured
- Verify user `api-user` exists and is enabled

