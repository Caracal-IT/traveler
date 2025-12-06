
**GET /ping** - JSON health check with metadata
```bash
curl http://localhost:8080/ping
```

**GET /ping/simple** - Minimal text health check (for load balancers)
```bash
curl http://localhost:8080/ping/simple
```

### Testing with IDE

Use the HTTP request file for interactive testing:
- Open `api/ping-endpoints.http` in IntelliJ IDEA, GoLand, or VS Code (with REST Client extension)
- Click the play button next to any request to test
- Includes 15+ test scenarios with automatic validation

See [API README](api/README.md) for detailed instructions.

## Testing Endpoints

### Quick Integration Test

```bash
# Start the server
go run ./cmd/traveler

# In another terminal, test the endpoints:
curl http://localhost:8080/ping
curl http://localhost:8080/ping/simple
```

### Automated Integration Test

```bash
./scripts/test-ping-endpoints.sh
```

## Project Structure

- `cmd/traveler` - application entrypoint
- `cmd/test-logging` - logging level demonstration
- `cmd/test-file-logging` - file logging demonstration
- `internal/app` - core app logic (Fiber-based HTTP server)
- `internal/handlers` - HTTP request handlers
  - `ping.go` - health check endpoints
  - `ping_test.go` - comprehensive tests and benchmarks
  - `routes.go` - route registration
- `pkg/config` - configuration loading with viper
- `pkg/log` - structured logging with zap
- `configs` - configuration files (YAML)
- `docs` - comprehensive documentation
- `scripts` - helper scripts
- `api` - OpenAPI placeholder
- `logs` - log files directory (created at runtime)
- `api` - OpenAPI placeholder
- `logs` - log files directory (created at runtime)



### Packages

Protected endpoint (requires a valid Keycloak access token for audience `traveler-app`):

- GET /api/packages/specials — returns sample specials data

Usage example (after starting Keycloak and importing realm):

1. Obtain a token (using the pre-provisioned `api-user`):
```
curl -s -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=traveler-app" \
  -d "username=api-user" \
  -d "password=ApiUser#1!" \
  http://localhost:8081/realms/traveler-dev/protocol/openid-connect/token | jq -r .access_token
```

2. Call the endpoint with the Bearer token:
```
TOKEN="$(# command above to fetch token)"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/packages/specials
```

Configuration for auth is under `auth` in `configs/config.yaml` and defaults to the local Keycloak realm.
