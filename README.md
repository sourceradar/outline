# Outline

A command-line code analysis tool that generates structured outlines for multiple programming languages using tree-sitter parsers. The tool analyzes source code files and generates clean outlines showing functions, classes, types, and other symbols.

**Can also be used as an MCP (Model Context Protocol) server** for integration with Claude Desktop and other MCP clients.

## Features

- **Multi-language support**: Go, Java, JavaScript, TypeScript, Python
- **Comprehensive symbol extraction**: Functions, classes, methods, types, interfaces, constants
- **Documentation extraction**: JSDoc, Go doc comments, Python docstrings, Javadoc
- **Fast and accurate**: Tree-sitter powered parsing
- **Dual mode**: CLI tool and optional MCP server

## Supported Languages

| Language   | File Extensions | Symbols Extracted |
|------------|-----------------|-------------------|
| Go         | `.go`           | Functions, methods, types, constants, variables, structs, interfaces |
| Java       | `.java`         | Classes, interfaces, enums, methods, constructors, fields, with modifiers and inheritance |
| JavaScript | `.js`, `.jsx`   | Functions, classes, arrow functions |
| TypeScript | `.ts`, `.tsx`   | Functions, classes, interfaces, types, with type annotations |
| Python     | `.py`           | Functions, classes (public symbols only) |

## Installation

```bash
go install github.com/sourceradar/outline@latest
```

Add to Claude Code:
```bash
claude mcp add -s user outline -- outline --mcp
```

Or build from source:

```bash
git clone https://github.com/sourceradar/outline.git
cd outline
go build ./cmd/outline
```

## Usage

### CLI Tool (Primary Usage)

Analyze a single file:

```bash
outline path/to/file.go
```

Override language detection:

```bash
outline --language go path/to/file.txt
```

### MCP Server Mode (Optional)

Run as MCP server:

```bash
outline --mcp
```

#### Claude Code Integration

After installing outline, add it to Claude Code:

```bash
claude mcp add -s user outline -- outline --mcp
```

#### Other MCP Clients

The server is compatible with any MCP client that supports `STDIO` communication. You can run it with:

```bash
outline --mcp
```

#### Development with MCP Inspector

Test the MCP server during development:

```bash
npx @modelcontextprotocol/inspector go run ./cmd/outline -- --mcp
```

This launches a web interface to test MCP tool calls interactively.

#### MCP Tool Usage

The server provides a single `outline` tool that accepts a file path parameter:

**Example Usage:**
```json
{
  "name": "outline",
  "arguments": {
    "file": "/path/to/your/source/file.go"
  }
}
```

**Response Format:**
The tool returns a text response containing the structured outline with language detection and symbol extraction.

## Example Output

For a Go file:

```
Language: go

package main

func main() {
    // Entry point of the application
}

type Config struct {
    Port     int    `json:"port"`
    Database string `json:"database"`
}

func (c *Config) Validate() error {
    // Validates the configuration
}
```

## Development

### Requirements

- Go 1.24.5 or later

### Commands

```bash
# Build the project
go build ./cmd/outline

# Run CLI directly
go run ./cmd/outline file.go

# Run as MCP server
go run ./cmd/outline --mcp

# Run tests
go test ./...

# Format code
go fmt ./...

# Clean dependencies
go mod tidy
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- General contribution process
- Adding support for new programming languages
- Testing requirements
- Code style guidelines

## License

MIT License - see the [LICENSE](LICENSE) file for details.

## Dependencies

- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official MCP Go SDK (for MCP mode)
- [Tree-sitter Go bindings](https://github.com/tree-sitter/go-tree-sitter) - Core tree-sitter functionality
- Language-specific tree-sitter grammars for parsing support
