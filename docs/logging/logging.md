# Logging Configuration Guide

## Overview

The traveler application uses structured JSON logging with configurable log levels. The minimum log level is specified in the configuration file, and only messages at or above that level will be logged.

**Logging can be configured to write to:**
- **stdout only** (default) - logs appear in console
- **stdout + file** - logs appear in console AND are saved to a file

## Configuration File

Edit `configs/config.yaml` to set the log level and optional file output:

```yaml
server:
  port: 8080

log:
  level: info  # Options: debug, info, warn, error
  file: ""     # Optional: path to log file (empty = stdout only)
```

### File Logging

To enable file logging, specify a file path in the `log.file` field:

```yaml
log:
  level: info
  file: logs/app.log  # Logs will be written to BOTH stdout and this file
```

**Important notes:**
- If `file` is empty or omitted, logs only go to stdout (console)
- If `file` is specified, logs go to BOTH stdout AND the file
- The directory for the log file must exist (e.g., create `logs/` directory)
- Log files are appended to, not overwritten
- All log entries are in JSON format for easy parsing

## Log Levels (in order of verbosity)

| Level | Description | What Gets Logged |
|-------|-------------|------------------|
| **debug** | Most verbose | Everything (debug, info, warn, error, fatal) |
| **info** | Normal operations | Info, warn, error, fatal |
| **warn** | Warnings and above | Warn, error, fatal |
| **error** | Errors only | Error, fatal |

## Using the Logger in Code

Import the log package and use the convenience methods:

```go
import "traveler/pkg/log"

// Debug level - only appears when log level is "debug"
log.Debug("processing request", "user_id", 123, "action", "fetch")

// Info level - appears when log level is "debug" or "info"
log.Info("request completed", "duration_ms", 45)

// Warn level - appears when log level is "debug", "info", or "warn"
log.Warn("rate limit approaching", "current", 90, "max", 100)

// Error level - appears at any log level except when set to fatal-only
log.Error("failed to process", "error", err)

// Fatal level - always logs and then exits the application
log.Fatal("critical error", "error", err)
```

## Log Output Format

All logs are output in structured JSON format:

```json
{
  "level": "info",
  "ts": "2025-12-05T08:58:39.132+0200",
  "caller": "app/app.go:25",
  "msg": "starting server",
  "address": ":8080"
}
```

## Testing Different Log Levels

We provide example configs for testing:

1. **Debug level** (see everything):
   ```bash
   # Edit configs/config.yaml and set: level: debug
   go run ./cmd/traveler
   ```

2. **Info level** (default, production recommended):
   ```bash
   # Edit configs/config.yaml and set: level: info
   go run ./cmd/traveler
   ```

3. **Error level** (only errors):
   ```bash
   # Edit configs/config.yaml and set: level: error
   go run ./cmd/traveler
   ```

## Testing File Logging

To test that logs are being written to a file:

```bash
# Create logs directory
mkdir -p logs

# Run the file logging test
go run ./cmd/test-file-logging

# View the log file contents
cat logs/app.log

# Or view with pretty printing
jq '.' logs/app.log
```

The test uses `configs/config.with-file.yaml` which has file logging enabled.

## Demo Script

Run the logging demonstration:

```bash
./scripts/demo-logging.sh
```

Or manually test:

```bash
go run ./cmd/test-logging
```

This will show how different log levels filter messages.

## Best Practices

1. **Development**: Use `debug` level to see all application activity
2. **Production**: Use `info` level for normal operations
3. **Troubleshooting**: Temporarily switch to `debug` level
4. **High-volume services**: Consider `warn` or `error` to reduce log volume
5. **Always include context**: Use key-value pairs to add relevant information to log messages

## Example: Updating Log Level

1. Open `configs/config.yaml`
2. Change the `log.level` value:
   ```yaml
   log:
     level: debug  # Change from info to debug
   ```
3. Restart the application
4. You'll now see debug messages that were previously hidden

