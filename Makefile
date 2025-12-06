.PHONY: build build-api build-cli run-api test fmt lint clean deps dev install help

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-s -w -X github.com/numtide/banner-generator/internal/version.Version=$(VERSION) -X github.com/numtide/banner-generator/internal/version.Commit=$(COMMIT)"

# Default target
all: fmt lint test build

# Build all binaries
build: build-api build-cli

# Build API server
build-api:
	@echo "Building API server..."
	go build $(LDFLAGS) -o bin/banner-api ./cmd/banner-api

# Build CLI tool
build-cli:
	@echo "Building CLI tool..."
	go build $(LDFLAGS) -o bin/banner-cli ./cmd/banner-cli

# Install binaries to GOPATH
install:
	go install ./cmd/banner-api
	go install ./cmd/banner-cli

# Run API server
run-api:
	go run ./cmd/banner-api

# Run API server with custom flags
run-api-dev:
	go run ./cmd/banner-api -port 8080 -allow "numtide,nixos"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	find . -name "*.go" -not -path "./.direnv/*" -not -path "./vendor/*" -not -path "./bin/*" | xargs gofmt -w

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run --no-config

# Run quick lint (faster, fewer checks)
lint-fast:
	golangci-lint run --no-config --fast

# Tidy and verify dependencies
mod:
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run API server with live reload (requires gorefresh)
dev:
	@command -v gorefresh >/dev/null 2>&1 || { echo "gorefresh not found. Install with: go install github.com/draganm/gorefresh@latest"; exit 1; }
	gorefresh ./cmd/banner-api

# Generate mocks (if needed in future)
mocks:
	@echo "No mocks to generate yet"

# Run security checks
sec:
	@echo "Running security checks..."
	go list -json -deps ./... | nancy sleuth

# Create release build
release:
	@echo "Building release binaries..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/banner-api-linux-amd64 ./cmd/banner-api
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/banner-cli-linux-amd64 ./cmd/banner-cli
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/banner-api-darwin-amd64 ./cmd/banner-api
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/banner-cli-darwin-amd64 ./cmd/banner-cli
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/banner-api-darwin-arm64 ./cmd/banner-api
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/banner-cli-darwin-arm64 ./cmd/banner-cli

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Format, lint, test, and build"
	@echo "  build         - Build all binaries"
	@echo "  build-api     - Build API server only"
	@echo "  build-cli     - Build CLI tool only"
	@echo "  install       - Install binaries to GOPATH"
	@echo "  run-api       - Run API server"
	@echo "  run-api-dev   - Run API server with dev settings"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linters"
	@echo "  lint-fast     - Run fast lint checks"
	@echo "  mod           - Tidy and verify go modules"
	@echo "  deps          - Download dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  dev           - Run with live reload (requires gorefresh)"
	@echo "  sec           - Run security checks"
	@echo "  release       - Build release binaries for multiple platforms"
	@echo "  help          - Show this help"