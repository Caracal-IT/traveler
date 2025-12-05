# Ping Endpoint Quick Reference

## Endpoints

| Endpoint | Method | Response Type | Use Case |
|----------|--------|---------------|----------|
| `/ping` | GET | JSON | Monitoring systems, detailed health checks |
| `/ping/simple` | GET | Plain text | Load balancers, simple health checks |

## Quick Commands

### Test Locally

```bash
# Start server
go run ./cmd/traveler

# Test JSON endpoint
curl http://localhost:8080/ping

# Test simple endpoint
curl http://localhost:8080/ping/simple
```

### Run Tests

```bash
# Unit tests
go test ./internal/handlers/

# Unit tests with coverage
go test -cover ./internal/handlers/

# Benchmarks
go test -bench=. -benchmem ./internal/handlers/

# Integration test
./scripts/test-ping-endpoints.sh
```

## Response Examples

### /ping Response

```json
{
  "status": "ok",
  "message": "pong",
  "timestamp": "2025-12-05T10:30:45.123456Z",
  "version": "1.0.0"
}
```

### /ping/simple Response

```
pong
```

## Load Balancer Configuration

### AWS ALB

```yaml
HealthCheck:
  Path: /ping/simple
  Protocol: HTTP
  Port: 8080
  HealthyThresholdCount: 2
  UnhealthyThresholdCount: 3
  Interval: 30
  Timeout: 5
```

### Kubernetes

```yaml
livenessProbe:
  httpGet:
    path: /ping/simple
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ping
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### Docker

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/ping/simple || exit 1
```

## Performance

- **Response Time**: < 10ms typical
- **Throughput**: > 100k requests/second
- **Memory**: ~2KB per request (JSON), ~512B (simple)

## Implementation Files

- Handler: `internal/handlers/ping.go`
- Tests: `internal/handlers/ping_test.go`
- Routes: `internal/handlers/routes.go`
- Docs: `docs/ping-endpoint.md`
- API Spec: `api/openapi.yaml`

## Related Documentation

- [Full Ping Endpoint Documentation](ping-endpoint.md)
- [Handlers Package README](../internal/handlers/README.md)
- [API Documentation](../api/openapi.yaml)

