# CloudBridge SDK Makefile
.PHONY: all build test clean install cli help

# Variables
BINARY_NAME=cloudbridge
CLI_BINARY=cmd/cloudbridge/cloudbridge
GO=go
GOFLAGS=-v
VERSION=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

all: test build ## Run tests and build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build SDK library
	@echo "Building CloudBridge SDK..."
	$(GO) build $(GOFLAGS) ./cloudbridge

cli: ## Build CLI tool
	@echo "Building CloudBridge CLI..."
	cd cmd/cloudbridge && $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "CLI built: cmd/cloudbridge/$(BINARY_NAME)"

install: cli ## Install CLI to /usr/local/bin
	@echo "Installing CloudBridge CLI..."
	sudo mv cmd/cloudbridge/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test -cover ./...
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	$(GO) mod tidy

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f cmd/cloudbridge/$(BINARY_NAME)
	rm -f coverage.out coverage.html
	$(GO) clean

run-cli: cli ## Build and run CLI (example: make run-cli ARGS="health")
	@echo "Running CLI..."
	./cmd/cloudbridge/$(BINARY_NAME) $(ARGS)

# Development targets
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev-test: ## Run tests in watch mode (requires entr)
	@echo "Watching for changes..."
	find . -name '*.go' | entr -c make test

# CLI examples
example-health: cli ## Run health check example
	./cmd/cloudbridge/$(BINARY_NAME) health --verbose

example-discover: cli ## Run discover example
	./cmd/cloudbridge/$(BINARY_NAME) discover

example-version: cli ## Show version
	./cmd/cloudbridge/$(BINARY_NAME) version

# Docker targets
docker-build: ## Build Docker image
	docker build -t cloudbridge-cli:$(VERSION) -f Dockerfile .

docker-run: ## Run CLI in Docker
	docker run --rm -e CLOUDBRIDGE_TOKEN=$(CLOUDBRIDGE_TOKEN) cloudbridge-cli:$(VERSION) $(ARGS)

# Release targets
release: test build cli ## Build release artifacts
	@echo "Creating release $(VERSION)..."
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/cloudbridge
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/cloudbridge
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/cloudbridge
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/cloudbridge
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/cloudbridge
	@echo "Release artifacts created in dist/"

checksums: release ## Generate checksums for release
	cd dist && sha256sum * > checksums.txt
	@echo "Checksums generated: dist/checksums.txt"

.DEFAULT_GOAL := help
