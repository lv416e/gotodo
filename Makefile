.PHONY: help build run test clean docker-build docker-up docker-down docker-logs dev

# Default target
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application locally"
	@echo "  make build        - Build the application binary"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-logs  - View Docker logs"
	@echo "  make dev          - Run in development mode with hot reload"

# Build the application
build:
	go build -o gotodo main.go

# Run the application locally
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f gotodo
	rm -rf tmp/

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Development mode with hot reload
dev:
	docker-compose -f docker-compose.dev.yml up