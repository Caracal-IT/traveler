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

