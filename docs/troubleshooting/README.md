# Troubleshooting Guide

This directory contains troubleshooting guides for common issues with the Traveler API.

## Authentication Issues

### 401 Unauthorized Error on /api/offerings/specials

**Problem**: The `/api/offerings/specials` endpoint returns 401 Unauthorized even with what appears to be a valid token.

**Root Cause**: Keycloak tokens are missing the `aud` (audience) claim. The token only contains `azp` (authorized party) claim.

**Solution**: 
1. The Keycloak realm configuration has been updated to include an audience mapper
2. Restart Keycloak to apply the fix:
   ```bash
   ./scripts/apply-auth-fix.sh
   ```

**Documentation**: 
- [Detailed Fix Guide](./FIX-401-missing-audience.md)
- [Comprehensive 401 Troubleshooting](./401-specials-endpoint.md)

**Quick Fix**:
```bash
# Restart Keycloak with updated realm configuration
cd /Users/ettienemare/GolandProjects/traveler
docker compose -f docker/docker-compose.yml restart keycloak

# Wait 30 seconds for startup, then test
sleep 30
./scripts/fix-auth-401.sh
```

## Common Issues

### Keycloak Not Running
If you get connection errors when trying to obtain a token:
```bash
docker compose -f docker/docker-compose.yml up -d keycloak
```

### Traveler API Not Running
If the API endpoints don't respond:
```bash
# Via Docker:
docker compose -f docker/docker-compose.yml up -d traveler

# Or directly:
go run ./cmd/traveler
```

### Token Expired
Tokens expire after 5 minutes (default). Get a fresh token before each test.

### Wrong Environment in HTTP Client
If using the HTTP request files in your IDE:
1. Ensure you've selected the correct environment (e.g., `dev`)
2. Check `api/http-client.env.json` for correct URLs and credentials

## Useful Scripts

### Diagnostic Scripts
- `scripts/diagnose-auth.sh` - Comprehensive authentication diagnostics
- `scripts/fix-auth-401.sh` - Quick fix script for 401 errors
- `scripts/apply-auth-fix.sh` - Apply the Keycloak realm fix and verify

### Testing Scripts
- `scripts/test-ping-endpoints.sh` - Test health check endpoints
- `scripts/demo-logging.sh` - Demonstrate logging functionality
- `scripts/demo-file-logging.sh` - Demonstrate file logging

## Checking Logs

### Application Logs
```bash
tail -f logs/app.log
```

### Keycloak Logs
```bash
docker compose -f docker/docker-compose.yml logs -f keycloak
```

### Traveler API Logs (Docker)
```bash
docker compose -f docker/docker-compose.yml logs -f traveler
```

## Enable Debug Logging

Edit `configs/config.yaml`:
```yaml
log:
  level: debug  # Change from 'info' to 'debug'
```

Then restart the application to see detailed JWT validation messages.

## Further Help

- [Main README](../../README.md)
- [API Documentation](../api/README.md)
- [Keycloak Getting Started](../tools/keycloak-getting-started.md)
- [Logging Documentation](../logging/logging.md)

