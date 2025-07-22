# MCP Outline

An MCP (Model Context Protocol) server that provides code outline generation for multiple programming languages using tree-sitter parsers. The server analyzes source code files and generates structured outlines showing functions, classes, types, and other symbols.

## Features

- **Multi-language support**: Go, JavaScript, TypeScript, Python
- **Comprehensive symbol extraction**: Functions, classes, methods, types, interfaces, constants
- **Documentation extraction**: JSDoc, Go doc comments, Python docstrings
- **MCP integration**: Compatible with any MCP client
- **Tree-sitter powered**: Fast and accurate parsing

## Supported Languages

| Language   | File Extensions | Symbols Extracted |
|------------|-----------------|-------------------|
| Go         | `.go`           | Functions, methods, types, constants, variables, structs, interfaces |
| JavaScript | `.js`, `.jsx`   | Functions, classes, arrow functions |
| TypeScript | `.ts`, `.tsx`   | Functions, classes, interfaces, types, with type annotations |
| Python     | `.py`           | Functions, classes (public symbols only) |

## Installation

```bash
go install github.com/sourceradar/mcp-outline@latest
```

Or build from source:

```bash
git clone https://github.com/sourceradar/mcp-outline.git
cd mcp-outline
go build .
```

## Usage

### As an MCP Server

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "outline": {
      "command": "mcp-outline"
    }
  }
}
```

### Development with MCP Inspector

```bash
npx @modelcontextprotocol/inspector go run .
```

### Tool Usage

The server provides a single `outline` tool that accepts a file path:

```json
{
  "name": "outline",
  "arguments": {
    "file_path": "/path/to/your/source/file.go"
  }
}
```

## Example Output

For a Go file:

```
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
go build .

# Run tests
go test ./...

# Format code
go fmt ./...

# Clean dependencies
go mod tidy
```

### Project Structure

```
├── main.go                    # MCP server entry point
├── outline/
│   ├── outline.go            # Main outline extraction logic
│   ├── outline_tool.go       # MCP tool handler
│   ├── util.go               # Tree-sitter utilities
│   └── languages/            # Language-specific parsers
│       ├── go.go
│       ├── js.go
│       ├── ts.go
│       └── python.go
└── go.mod                    # Go module definition
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

- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official MCP Go SDK
- [Tree-sitter Go bindings](https://github.com/tree-sitter/go-tree-sitter) - Core tree-sitter functionality
- Language-specific tree-sitter grammars for parsing support