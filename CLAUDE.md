# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server that provides code outline generation for multiple programming languages using tree-sitter parsers. The server analyzes source code files and generates structured outlines showing functions, classes, types, and other symbols.

## Supported Languages

- **Go** (.go files) - Functions, methods, types, constants, variables, structs, interfaces
- **JavaScript** (.js, .jsx files) - Functions, classes, arrow functions
- **TypeScript** (.ts, .tsx files) - Functions, classes, interfaces, types, with type annotations
- **Python** (.py files) - Functions, classes (public symbols only)

## Development Commands

```bash
# Build the project
go build .

# Run the MCP server directly
./mcp-outline

# Run with MCP inspector for development/testing
npx @modelcontextprotocol/inspector go run .

# Run all tests
go test ./...

# Run tests for specific language
go test ./outline/languages/

# Format code
go fmt ./...

# Clean dependencies
go mod tidy

# Install as binary
go install
```

## Architecture

### Core Components

- `main.go` - MCP server entry point with tool registration and stdio transport
- `outline/outline.go` - Main outline extraction logic with language detection and parser creation
- `outline/outline_tool.go` - MCP tool handler implementing the outline functionality and file extension detection
- `outline/languages/` - Language-specific outline extractors:
  - `go.go` - Go language parser with struct/interface/method handling
  - `js.go` - JavaScript parser with class and function extraction
  - `ts.go` - TypeScript parser with type annotations and interfaces
  - `python.go` - Python parser filtering private symbols (underscore prefix)
  - `util.go` - Shared utilities for tree-sitter node processing

### Key Functions

- `ExtractOutline(content []byte, language string)` - Main entry point in `outline/outline.go`
- `createParserForLanguage(language string)` - Parser factory in `outline/outline.go`
- `OutlineToolHandler()` - MCP tool handler in `outline/outline_tool.go`
- `detectLanguage(filePath string)` - File extension to language mapping in `outline/outline_tool.go`
- `getNodeText()` and `findDocComment()` - Utility functions in `outline/languages/util.go`

### Adding New Language Support

Follow the pattern in CONTRIBUTING.md:

1. Add tree-sitter dependency to `go.mod`
2. Create `outline/languages/{lang}.go` with `Extract{Lang}Outline()` function
3. Add language case to `createParserForLanguage()` in `outline/outline.go`
4. Add extraction case to `ExtractOutline()` in `outline/outline.go`
5. Add file extension mapping to `detectLanguage()` in `outline/outline_tool.go`
6. Write comprehensive tests in `outline/languages/{lang}_test.go`

## Code Patterns

- Each language parser follows recursive tree traversal using `processNode` functions
- Documentation comments extracted when available (JSDoc, Go doc comments, Python docstrings)
- Go parser handles methods with receivers, struct fields, and interface methods
- TypeScript parser includes type annotations and extends/implements clauses
- Python parser filters out private symbols (names starting with underscore)
- All parsers generate readable outline format with proper indentation
- Memory management: Always use `defer parser.Close()` and `defer tree.Close()`

## MCP Integration

Implements Model Context Protocol (MCP) specification:

- Uses stdio transport for communication with MCP clients
- Registers single "outline" tool accepting `file_path` parameter
- Returns structured text outlines of code symbols
- Handles errors gracefully with proper MCP error responses
- Compatible with Claude Desktop, MCP Inspector, and other MCP clients

## Testing

Each language has comprehensive test coverage in `outline/languages/*_test.go`:

- Import statements and module systems
- Function/method declarations with parameters and return types
- Class/struct declarations with inheritance
- Documentation comment extraction
- Access modifier filtering (public/private visibility)
- Error handling for malformed code

Run language-specific tests: `go test ./outline/languages/ -v`

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - Official MCP Go SDK
- `github.com/tree-sitter/go-tree-sitter` - Core tree-sitter Go bindings  
- Language-specific tree-sitter grammars:
  - `github.com/tree-sitter/tree-sitter-go`
  - `github.com/tree-sitter/tree-sitter-javascript` 
  - `github.com/tree-sitter/tree-sitter-typescript`
  - `github.com/tree-sitter/tree-sitter-python`