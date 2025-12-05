# API Directory

This directory contains API documentation and testing files for the traveler service.

## Contents

- **openapi.yaml** - OpenAPI 3.0 specification for the traveler API
- **ping-endpoints.http** - HTTP request file for testing ping endpoints
- **http-client.env.json** - Environment configuration for HTTP requests

## Testing with HTTP Files

The `ping-endpoints.http` file provides a convenient way to test API endpoints directly from your IDE.

### Supported IDEs

#### JetBrains IDEs (IntelliJ IDEA, GoLand, WebStorm)

Built-in support - no plugins needed.

**How to use:**
1. Open `ping-endpoints.http` in your IDE
2. Click the green play button (▶️) next to any request
3. View results in the "Run" panel at the bottom
4. Tests will automatically validate responses

**Run all requests:**
- Right-click the file → "Run All Requests in File"
- Or use the toolbar button

#### Visual Studio Code

Requires the "REST Client" extension by Huachao Mao.

**Setup:**
1. Install extension: `ext install humao.rest-client`
2. Open `ping-endpoints.http`
3. Click "Send Request" above each request
4. View results in a new panel

**Run all requests:**
- Command Palette → "REST Client: Send All Requests in File"

### Available Requests

The HTTP file includes 15+ test scenarios:

1. **Basic Health Checks**
   - JSON ping endpoint
   - Simple text ping endpoint
   - Root endpoint

2. **Custom Headers**
   - User-Agent testing
   - Request ID tracking

3. **Load Balancer Simulation**
   - AWS ALB health check
   - Kubernetes probe

4. **Error Scenarios**
   - Wrong HTTP method
   - Non-existent endpoints

5. **Performance Testing**
   - Multiple sequential requests
   - Latency measurement

### Environment Configuration

The `http-client.env.json` file defines environment-specific settings:

```json
{
  "dev": {
    "baseUrl": "http://localhost:8080"
  },
  "staging": {
    "baseUrl": "http://staging.traveler.local:8080"
  },
  "prod": {
    "baseUrl": "https://api.traveler.example.com"
  }
}
```

**Switch environments:**
- IntelliJ/GoLand: Use the environment selector in the HTTP request toolbar
- VS Code: Add to settings.json:
  ```json
  "rest-client.environmentVariables": {
    "$shared": {},
    "dev": { "baseUrl": "http://localhost:8080" },
    "prod": { "baseUrl": "https://api.example.com" }
  }
  ```

### Private Environment Variables

For secrets (API keys, tokens), create `http-client.private.env.json`:

```json
{
  "dev": {
    "apiKey": "dev-secret-key",
    "authToken": "Bearer dev-token"
  }
}
```

**Important:** This file is gitignored by default.

## Quick Start

### 1. Start the Server

```bash
go run ./cmd/traveler
```

### 2. Test Endpoints

**Option A: Using HTTP file (recommended)**
- Open `ping-endpoints.http`
- Click "Send Request" for any test

**Option B: Using curl**
```bash
# JSON endpoint
curl http://localhost:8080/ping

# Simple endpoint
curl http://localhost:8080/ping/simple
```

**Option C: Using automated script**
```bash
./scripts/test-ping-endpoints.sh
```

## Response Examples

### JSON Ping Response

```json
{
  "status": "ok",
  "message": "pong",
  "timestamp": "2025-12-05T10:30:45.123456Z",
  "version": "1.0.0"
}
```

### Simple Ping Response

```
pong
```

## OpenAPI Specification

The `openapi.yaml` file provides a complete API specification in OpenAPI 3.0 format.

### View the Specification

**Online:**
1. Copy the content of `openapi.yaml`
2. Visit https://editor.swagger.io/
3. Paste and view the interactive documentation

**Local:**
```bash
# Install swagger UI
npm install -g swagger-ui-watcher

# Serve the spec
swagger-ui-watcher api/openapi.yaml
```

### Generate Client Code

```bash
# Install openapi-generator
npm install -g @openapitools/openapi-generator-cli

# Generate Go client
openapi-generator-cli generate \
  -i api/openapi.yaml \
  -g go \
  -o ./generated/client

# Generate TypeScript client
openapi-generator-cli generate \
  -i api/openapi.yaml \
  -g typescript-axios \
  -o ./generated/ts-client
```

## Integration Testing

### Using the HTTP File

The HTTP file includes automated tests:

```javascript
> {%
    client.test("Status is 200", function() {
        client.assert(response.status === 200);
    });
%}
```

These tests run automatically when you send requests in supported IDEs.

### CI/CD Integration

For automated testing in CI/CD pipelines, use the test script:

```bash
./scripts/test-ping-endpoints.sh
```

Or use Go tests:

```bash
go test ./internal/handlers/ -v
```

## Documentation

- [Ping Endpoint Documentation](../docs/ping-endpoint.md) (if exists)
- [Ping Quick Reference](../docs/ping-quick-reference.md)
- [Handlers README](../internal/handlers/README.md)

## Tips

### Response Time Measurement

Both IntelliJ and VS Code show response times for each request:
- IntelliJ: Shows in the response panel header
- VS Code: Shows in the response status line

### History

IDEs keep a history of all requests:
- IntelliJ: View → Tool Windows → HTTP Client → Show History
- VS Code: Responses are saved in `.vscode/rest-client/` (gitignored)

### Variables

Use variables for reusable values:

```http
@token = my-auth-token

GET {{baseUrl}}/ping
Authorization: Bearer {{token}}
```

### Dynamic Variables

Use built-in dynamic variables:

```http
GET {{baseUrl}}/ping
X-Request-ID: {{$uuid}}
X-Timestamp: {{$timestamp}}
```

## Troubleshooting

### Connection Refused

**Error:** `Connection refused`

**Solution:**
1. Ensure server is running: `go run ./cmd/traveler`
2. Check port: Server should be on port 8080
3. Verify baseUrl in HTTP file matches server address

### 404 Not Found

**Error:** `404 Not Found`

**Solution:**
1. Check endpoint path is correct
2. Verify routes are registered in `internal/handlers/routes.go`
3. Check server logs for routing issues

### IDE Not Recognizing HTTP File

**IntelliJ/GoLand:**
- File → Settings → Languages & Frameworks → HTTP Client
- Ensure plugin is enabled

**VS Code:**
- Install "REST Client" extension
- File association should be automatic for `.http` files

## Related Files

- `../internal/handlers/ping.go` - Endpoint implementation
- `../internal/handlers/ping_test.go` - Go unit tests
- `../scripts/test-ping-endpoints.sh` - Integration test script

