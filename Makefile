.PHONY: build install test lint dev clean fmt vet

# Version info
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/abdul-hamid-achik/fuego/internal/version.Version=$(VERSION)"

# Binary name
BINARY := fuego

# Build the CLI
build:
	@echo "Building $(BINARY)..."
	@go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/fuego
	@echo "Done: bin/$(BINARY)"

# Install the CLI globally
install:
	@echo "Installing $(BINARY)..."
	@go install $(LDFLAGS) ./cmd/fuego
	@echo "Done: $(BINARY) installed"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Run development server (for testing the framework itself)
dev:
	@go run ./cmd/fuego dev

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@echo "Done"

# Run all checks
check: fmt vet test

# Build for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 ./cmd/fuego
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 ./cmd/fuego
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 ./cmd/fuego
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 ./cmd/fuego
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/fuego
	@echo "Done: binaries in dist/"

# Show help
help:
	@echo "Fuego Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the CLI binary"
	@echo "  install     Install the CLI globally"
	@echo "  test        Run tests"
	@echo "  test-cover  Run tests with coverage report"
	@echo "  lint        Run linter"
	@echo "  fmt         Format code"
	@echo "  vet         Vet code"
	@echo "  check       Run fmt, vet, and test"
	@echo "  clean       Remove build artifacts"
	@echo "  build-all   Build for all platforms"
	@echo "  help        Show this help"
