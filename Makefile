.PHONY: build test install clean

# Default values if not set
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Build the binary
build:
	mkdir -p dist/$(GOOS)-$(GOARCH)
	CGO_ENABLED=1 go build -o dist/$(GOOS)-$(GOARCH)/outline ./cmd/outline

# Run tests
test:
	go test ./...

# Install binary to GOPATH/bin
install:
	go install ./cmd/outline

# Clean build artifacts
clean:
	rm -f outline
	rm -rf dist

# Run with MCP inspector for development
dev:
	npx @modelcontextprotocol/inspector go run ./cmd/outline -- --mcp

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy