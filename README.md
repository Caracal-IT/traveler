# traveler

Minimal Go project skeleton for the traveler service with structured logging and configuration management.

## Quick start

Build:

```bash
go build ./...
```

Run (from project root):

```bash
go run ./cmd/traveler
```

## Configuration

Configuration is loaded from `configs/config.yaml`. Example:

```yaml
server:
  port: 8080

log:
  level: info  # Options: debug, info, warn, error
  file: ""     # Optional: logs/app.log (empty = stdout only)
```

### File Logging

To enable logging to a file, specify a path in the `log.file` field:

```yaml
log:
  level: info
  file: logs/app.log  # Logs go to BOTH stdout and this file
```

**Note:** Make sure the directory exists before starting the app:
```bash
mkdir -p logs
```

### Log Levels

The application supports the following log levels (in order of verbosity):

- **debug**: Most verbose, logs everything including debug messages
- **info**: Logs informational messages, warnings, and errors
- **warn**: Logs warnings and errors only
- **error**: Logs only error messages

When you set a minimum log level in the config file, only messages at that level or higher will be logged.

Example configs are provided:
- `configs/config.yaml` - Info level (default)
- `configs/config.debug.yaml` - Debug level (all messages)
- `configs/config.error.yaml` - Error level (errors only)

### Using the Logger

The `pkg/log` package provides convenient methods:

```go
import "traveler/pkg/log"

log.Debug("debug message", "key", "value")
log.Info("info message", "key", "value")
log.Warn("warning message", "key", "value")
log.Error("error message", "key", "value")
log.Fatal("fatal message", "key", "value")  // exits the program
```

All logs are output in structured JSON format with timestamps, caller information, and custom key-value pairs.

## Testing

### Test log level filtering
```bash
go run ./cmd/test-logging
```

### Test file logging
```bash
mkdir -p logs
go run ./cmd/test-file-logging
cat logs/app.log  # View the log file
```

## API Endpoints

### Health Check

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

