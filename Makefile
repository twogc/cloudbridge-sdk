.PHONY: help test lint build clean install-go install-python install-js

help:
	@echo "CloudBridge SDK - Build Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make test          - Run all tests"
	@echo "  make test-go       - Run Go tests"
	@echo "  make lint          - Run linters"
	@echo "  make lint-go       - Run Go linter"
	@echo "  make build         - Build all SDKs"
	@echo "  make build-go      - Build Go SDK"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install-go    - Install Go dependencies"
	@echo "  make install-python - Install Python dependencies"
	@echo "  make install-js    - Install JavaScript dependencies"
	@echo "  make fmt           - Format all code"
	@echo "  make fmt-go        - Format Go code"

# Go SDK targets
test-go:
	cd go && go test -v -race -cover ./...

lint-go:
	cd go && golangci-lint run ./...

build-go:
	cd go && go build ./...

fmt-go:
	cd go && go fmt ./...
	cd go && goimports -w .

install-go:
	cd go && go mod download
	cd go && go mod tidy

# Combined targets
test: test-go

lint: lint-go

build: build-go

fmt: fmt-go

clean:
	find . -name "*.test" -delete
	find . -name "coverage.txt" -delete
	find . -name "coverage.html" -delete
	rm -rf go/bin
	rm -rf python/dist
	rm -rf python/build
	rm -rf javascript/dist
	rm -rf javascript/node_modules

# Development
dev-setup:
	@echo "Setting up development environment..."
	make install-go
	@echo "Development setup complete!"

# CI targets
ci-test: test

ci-lint: lint

ci: ci-lint ci-test
