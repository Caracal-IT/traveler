# Handlers Package

This package contains HTTP request handlers for the traveler service.

## Overview

The handlers package provides modular, testable HTTP handlers for the Fiber web framework. Each handler is responsible for processing specific types of requests.

## Handlers

### Ping Handlers

Health check endpoints for monitoring and load balancing.

- **PingHandler** (`/ping`) - JSON health check with metadata
- **PingHandlerSimple** (`/ping/simple`) - Minimal text health check

See [Ping Endpoint Documentation](../../docs/ping-endpoint.md) for details.

## Structure

```
internal/handlers/
├── ping.go           # Ping endpoint handlers
├── ping_test.go      # Comprehensive tests and benchmarks
└── routes.go         # Route registration
```

## Usage

### Registering Routes

Routes are automatically registered in `internal/app/app.go`:

```go
import "traveler/internal/handlers"

// In your app setup:
handlers.RegisterRoutes(app)
```

### Creating New Handlers

1. Create a new file: `internal/handlers/myhandler.go`
2. Implement your handler function:
```go
package handlers

import "github.com/gofiber/fiber/v2"

func MyHandler(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "ok"})
}
```

3. Add tests: `internal/handlers/myhandler_test.go`
4. Register in `routes.go`:
```go
func RegisterRoutes(app *fiber.App) {
    app.Get("/ping", PingHandler)
    app.Get("/myroute", MyHandler)  // Add your route
}
```

## Testing

### Run All Tests

```bash
go test ./internal/handlers/
```

### Run Tests with Coverage

```bash
go test -cover ./internal/handlers/
```

### Run Verbose Tests

```bash
go test -v ./internal/handlers/
```

### Run Specific Test

```bash
go test -run TestPingHandler ./internal/handlers/
```

## Benchmarks

### Run All Benchmarks

```bash
go test -bench=. -benchmem ./internal/handlers/
```

### Run Specific Benchmark

```bash
go test -bench=BenchmarkPingHandler -benchmem ./internal/handlers/
```

### Run Parallel Benchmarks

```bash
go test -bench=Parallel -benchmem ./internal/handlers/
```

## Best Practices

### 1. Handler Function Signature

Always use the Fiber handler signature:

```go
func HandlerName(c *fiber.Ctx) error
```

### 2. Error Handling

Return errors properly:

```go
func MyHandler(c *fiber.Ctx) error {
    if err := someOperation(); err != nil {
        log.Error("operation failed", "error", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "internal server error",
        })
    }
    return c.JSON(fiber.Map{"status": "ok"})
}
```

### 3. Logging

Use structured logging:

```go
import "traveler/pkg/log"

func MyHandler(c *fiber.Ctx) error {
    log.Debug("handling request", "path", c.Path(), "method", c.Method())
    // ... handler logic
}
```

### 4. Response Types

Define response structs for clarity:

```go
type MyResponse struct {
    Status  string `json:"status"`
    Data    string `json:"data"`
}

func MyHandler(c *fiber.Ctx) error {
    resp := MyResponse{
        Status: "ok",
        Data:   "example",
    }
    return c.JSON(resp)
}
```

### 5. Testing

Write comprehensive tests:

```go
func TestMyHandler(t *testing.T) {
    app := fiber.New()
    app.Get("/myroute", MyHandler)

    t.Run("returns 200", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/myroute", nil)
        resp, err := app.Test(req)
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
    })
}
```

### 6. Benchmarking

Add benchmarks for performance-critical handlers:

```go
func BenchmarkMyHandler(b *testing.B) {
    app := fiber.New()
    app.Get("/myroute", MyHandler)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("GET", "/myroute", nil)
        resp, _ := app.Test(req, -1)
        resp.Body.Close()
    }
}
```

## Performance Tips

1. **Minimize Allocations** - Reuse objects when possible
2. **Use Sync.Pool** - For frequently allocated objects
3. **Avoid Reflection** - Direct struct access is faster
4. **Stream Large Responses** - Use `c.SendStream()` for large data
5. **Cache Static Data** - Pre-compute constant responses

## Middleware

Handlers can use middleware for cross-cutting concerns:

```go
import "github.com/gofiber/fiber/v2/middleware/logger"

func RegisterRoutes(app *fiber.App) {
    // Apply middleware to specific routes
    app.Get("/myroute", 
        logger.New(),
        MyHandler,
    )
}
```

## Documentation

- [Ping Endpoint Documentation](../../docs/ping-endpoint.md)
- [API Documentation](../../api/openapi.yaml)
- [Fiber Framework Docs](https://docs.gofiber.io/)

## Contributing

When adding new handlers:

1. ✅ Create handler file with clear naming
2. ✅ Add comprehensive tests (>80% coverage)
3. ✅ Add benchmarks for performance-critical code
4. ✅ Document public APIs with comments
5. ✅ Update OpenAPI specification
6. ✅ Update this README

