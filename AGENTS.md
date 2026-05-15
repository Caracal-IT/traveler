# AGENTS: How to be productive in this repository

Purpose: give an AI coding agent the minimal, actionable orientation to become productive quickly in this Go microservice.

Quick start (commands)
- Run locally: `go run ./cmd/traveler` (uses `configs/config.yaml` or defaults)
- Build: `go build -o bin/traveler ./cmd/traveler`
- Unit tests: `go test ./...`
- Run handler tests only: `go test ./internal/handlers/`
- Use the interactive HTTP file: open `api/ping-endpoints.http` in GoLand/IDE or VS Code REST Client to exercise endpoints.

Big picture architecture
- Single Go module `module traveler` (see `go.mod`).
- HTTP server: Fiber (`internal/app/app.go`) — `Run` initializes DB, registers routes (`internal/handlers/routes.go`) and starts Fiber.
- Routing and handlers live in `internal/handlers` (examples: `ping.go`, `routes.go`). Tests and benchmarks are in the same package and exercise handlers via `fiber.App.Test`.
- Config: `pkg/config` wraps Viper and provides `LoadOrDefault` used by `cmd/traveler/main.go`.
- Logging: structured JSON logging via `pkg/log` (zap). Optional Elasticsearch sink is supported via config.
- Auth: Keycloak/OpenID integration implemented in `pkg/auth/jwt.go` (JWKS caching/refresh and Keycloak-compatible audience checks).
- Database: local SQLite using `modernc.org/sqlite` (pure-Go) in `internal/db/init.go`. Schema applied from `db/schema.sql` on startup.

Key repository conventions & patterns (examples you should follow)
- Config-first: load config at startup with `config.LoadOrDefault("configs/config.yaml")` (see `cmd/traveler/main.go`). Defaults are defined in `pkg/config/config.go`.
- Single DB connection for embedded DB: `Init` sets `SetMaxOpenConns(1)` and WAL mode in the DSN. When modifying schema, update `db/schema.sql` (Init re-applies only on empty schema/run).
- Structured logging via `traveler/pkg/log` helpers (e.g. `log.Info("msg", "key", val)`). Do not construct your own global loggers — reuse `pkg/log`.
- Route registration: add endpoints in `internal/handlers` and register them from `routes.go`. Protect groups using `auth.JWTMiddleware(cfg)`; handlers should use `c.Locals("claims")` to access validated token claims.
- JWT/audience handling: `pkg/auth/jwt.go` does non-standard audience checks (supports `aud`, `azp`, and `resource_access`) and allows JWKS URL override via config (`auth.jwks_url`). If you need a different JWKS host (containers), use that override.

Developer workflows & debugging hints
- IDE-driven HTTP testing: `api/ping-endpoints.http` is the canonical interactive test suite for endpoints. Use `http-client.env.json` to swap environments.
- Local Keycloak for integration testing is provided under `docker/keycloak` (env files + realm export in `docker/keycloak/realm-export/`). README shows how to fetch a token and call `/api/offerings/specials`.
- When iterating on handlers, unit tests use `fiber.App.Test` (see `internal/handlers/ping_test.go`). Run `go test -run TestPingHandler ./internal/handlers/` for focused iterations.
- Graceful shutdown: `internal/app` uses `signal.NotifyContext` and `app.ShutdownWithContext` (5s timeout) — ensure handlers finish quickly for tests relying on shutdown.

Integration points & external dependencies
- Keycloak (OpenID Connect) — issuer configured under `auth.issuer` in `configs/config.yaml`; audience defaults to `traveler-app`. Example token fetch in `README.md`.
- Elasticsearch (optional) — `log.elasticsearch` section controls shipping logs; `pkg/log` creates an ES syncer when enabled.
- SQLite (embedded) — `modernc.org/sqlite` driver; DB path default `db/traveler.db` (see `pkg/config` defaults).

Where to look for common tasks
- Add a new HTTP endpoint: `internal/handlers/` + register in `internal/handlers/routes.go` + update `api/openapi.yaml` and `api/ping-endpoints.http` if you want IDE-driven tests.
- Change config defaults: `pkg/config/config.go` (viper defaults and ReadInConfig).
- Modify logging behavior: `pkg/log/log.go` (encoder config, multi-writer, ES syncer hook).
- Work with tokens: `pkg/auth/jwt.go` (JWKS caching, leeway, audience compatibility logic).
- Database/schema changes: edit `db/schema.sql` (Init applies it on startup). Keep scripts idempotent (schema uses IF NOT EXISTS guards).

Gotchas & quick tips for agents
- Tests and HTTP file expect server on port 8080 by default (see `pkg/config` defaults and `api/openapi.yaml`).
- The JWT middleware tolerates issuer scheme differences (http/https) and trailing slashes — tests and integration may run Keycloak on a different host; use `auth.jwks_url` to override JWKS URL when needed.
- The internal DB is opened with `SetMaxOpenConns(1)`; concurrent transactional tests that expect multiple connections will fail — design tests accordingly.
- README references helper scripts (e.g. `./scripts/test-ping-endpoints.sh`) that may not be present in the tree — prefer `go run` + `api/*.http` or `go test` for reproducible runs.

Minimal checklist an agent should run before making changes
1. Run `go test ./...` and inspect failures.
2. Start server `go run ./cmd/traveler` and exercise `api/ping-endpoints.http`.
3. If touching auth, ensure Keycloak is running (use `docker/keycloak` assets) and fetch a token as in README.

Files referenced above (entry points)
- `cmd/traveler/main.go`
- `internal/app/app.go`
- `internal/db/init.go`, `db/schema.sql`
- `internal/handlers/*` (handlers, tests, routes.go)
- `pkg/config/config.go`, `pkg/log/log.go`, `pkg/auth/jwt.go`
- `api/openapi.yaml`, `api/ping-endpoints.http`

