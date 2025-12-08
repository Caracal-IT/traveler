# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o traveler ./cmd/traveler/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/traveler .

# Copy configs if they exist
COPY --from=builder /app/configs ./configs

# Copy database assets (schema) for runtime initialization
COPY --from=builder /app/db ./db

EXPOSE 8080

CMD ["./traveler"]
