.PHONY: build test clean install snapshot help

# Build variables
BINARY_NAME=dnaspec
BUILD_DIR=build
MAIN_PATH=./cmd/dnaspec

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# LDFLAGS to inject version information
LDFLAGS=-ldflags "-s -w \
	-X github.com/aviator5/dnaspec/internal/cli.Version=$(VERSION) \
	-X github.com/aviator5/dnaspec/internal/cli.Commit=$(COMMIT) \
	-X github.com/aviator5/dnaspec/internal/cli.Date=$(DATE)"

## help: Show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary for the current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -rf dist

## install: Install the binary to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PATH)

## snapshot: Create a snapshot release using goreleaser
snapshot:
	@echo "Creating snapshot release..."
	goreleaser release --snapshot --clean

## tidy: Run go mod tidy
tidy:
	@echo "Running go mod tidy..."
	go mod tidy

## lint: Run linters (requires golangci-lint)
lint:
	@echo "Running linters..."
	golangci-lint run

## fmt: Format all Go files
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## check: Run fmt, vet, and test
check: fmt vet test

.DEFAULT_GOAL := help
