.PHONY: build test install clean

# Build the binary
build:
	go build -o outline ./cmd/outline

# Run tests
test:
	go test ./...

# Install binary to GOPATH/bin
install:
	go install ./cmd/outline

# Clean build artifacts
clean:
	rm -f outline

# Run with MCP inspector for development
dev:
	npx @modelcontextprotocol/inspector go run ./cmd/outline -- --mcp

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy