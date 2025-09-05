# Alfred API Makefile

.PHONY: test test-verbose test-coverage test-clean build run dev help

# Default target
help:
	@echo "Available commands:"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-clean     - Clean test artifacts and run tests"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run tests and start the application"
	@echo "  clean          - Clean build artifacts"

# Run all tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean test artifacts and run tests
test-clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html
	@echo "Running clean tests..."
	go clean -testcache
	go test ./...

# Build the application
build:
	@echo "Building application..."
	go build -o bin/alfred ./cmd/api

# Run the application
run: build
	@echo "Starting Alfred API..."
	./bin/alfred

# Development workflow: test then run
dev: test-clean
	@echo "Development mode: tests passed, starting server..."
	go run ./cmd/api

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/alfred
	rm -f coverage.out coverage.html
	go clean