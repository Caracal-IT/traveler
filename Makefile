build:
	go build ./...

run:
	go run ./cmd/traveler

test:
	go test ./... -v

fmt:
	gofmt -w .

lint:
	golangci-lint run || true

# Docker targets
docker-build:
	docker build -t traveler:latest .

docker-run:
	docker run -p 8080:8080 traveler:latest

docker-stop:
	docker stop $$(docker ps -q --filter ancestor=traveler:latest)

docker-clean:
	docker rmi traveler:latest

# Local SQLite database targets
DB_FILE ?= db/traveler.db
DB_SCHEMA ?= db/schema.sql

db-init:
	@if ! command -v sqlite3 >/dev/null 2>&1; then \
		echo "Error: sqlite3 CLI not found. Please install SQLite (e.g., brew install sqlite)."; \
		exit 1; \
	fi
	@mkdir -p db
	@echo "Initializing $(DB_FILE) from $(DB_SCHEMA)" && sqlite3 $(DB_FILE) < $(DB_SCHEMA)
	@echo "Done. Created: $(DB_FILE)"

db-clean:
	@rm -f $(DB_FILE)
	@echo "Removed: $(DB_FILE)"

