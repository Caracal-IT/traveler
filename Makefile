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

