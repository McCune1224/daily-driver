# Smart Apartment Dashboard Makefile

.PHONY: build run clean test deps dev

# Default target
all: deps build

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	go build -o bin/apartment-dashboard .

# Run the application in development mode
dev:
	go run .

# Run the application
run: build
	./bin/apartment-dashboard

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Run with live reload (requires air)
live:
	air

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/apartment-dashboard .

# Docker build
docker-build:
	docker build -t apartment-dashboard .

# Docker run
docker-run:
	docker run -p 8080:8080 apartment-dashboard

# Install air for live reload
install-air:
	go install github.com/cosmtrek/air@latest

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  deps        - Install Go dependencies"
	@echo "  build       - Build the application"
	@echo "  dev         - Run in development mode"
	@echo "  run         - Build and run"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  live        - Run with live reload (requires air)"
	@echo "  build-prod  - Build for production"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  help        - Show this help"
