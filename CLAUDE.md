# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a command-line code analysis tool that generates structured outlines for multiple programming languages using tree-sitter parsers. The tool analyzes source code files and generates clean outlines showing functions, classes, types, and other symbols. It can also optionally run as an MCP (Model Context Protocol) server.

## Supported Languages

- **Go** (.go files) - Functions, methods, types, constants, variables, structs, interfaces
- **Java** (.java files) - Classes, interfaces, enums, methods, constructors, fields, with modifiers and inheritance
- **JavaScript** (.js, .jsx files) - Functions, classes, arrow functions
- **TypeScript** (.ts, .tsx files) - Functions, classes, interfaces, types, with type annotations
- **Python** (.py files) - Functions, classes (public symbols only)

## Development Commands

```bash
# Build the project
go build ./cmd/outline

# Run CLI tool directly
go run ./cmd/outline file.go

# Run as MCP server
go run ./cmd/outline --mcp

# Run with MCP inspector for development/testing
npx @modelcontextprotocol/inspector go run ./cmd/outline -- --mcp

# Run all tests
go test ./...

# Run tests for specific language
go test ./pkg/outline/languages/

# Format code
go fmt ./...

# Clean dependencies
go mod tidy

# Install as binary
go install ./cmd/outline
```

## Architecture

### Core Components

- `cmd/outline/main.go` - Application entry point with CLI and MCP mode handling
- `pkg/outline/outline.go` - Main outline extraction logic with language detection and parser creation
- `internal/server/tool.go` - MCP tool handler implementing the outline functionality
- `internal/cli/cli.go` - CLI implementation for standalone usage
- `internal/detector/` - Language detection from file extensions
- `pkg/outline/languages/` - Language-specific outline extractors:
  - `go.go` - Go language parser with struct/interface/method handling
  - `java.go` - Java language parser with class/interface/enum/method handling and modifiers
  - `js.go` - JavaScript parser with class and function extraction
  - `ts.go` - TypeScript parser with type annotations and interfaces
  - `python.go` - Python parser filtering private symbols (underscore prefix)
  - `util.go` - Shared utilities for tree-sitter node processing

### Key Functions

- `ExtractOutline(content []byte, language string)` - Main entry point in `pkg/outline/outline.go`
- `createParserForLanguage(language string)` - Parser factory in `pkg/outline/outline.go`
- `OutlineToolHandler()` - MCP tool handler in `internal/server/tool.go`
- `DetectLanguage(filePath string)` - File extension to language mapping in `internal/detector/`
- `getNodeText()` and `findDocComment()` - Utility functions in `pkg/outline/languages/util.go`

### Adding New Language Support

Follow the pattern in CONTRIBUTING.md:

1. Add tree-sitter dependency to `go.mod`
2. Create `pkg/outline/languages/{lang}.go` with `Extract{Lang}Outline()` function
3. Add language case to `createParserForLanguage()` in `pkg/outline/outline.go`
4. Add extraction case to `ExtractOutline()` in `pkg/outline/outline.go`
5. Add file extension mapping to `DetectLanguage()` in `internal/detector/`
6. Write comprehensive tests in `pkg/outline/languages/{lang}_test.go`

## Code Patterns

- Each language parser follows recursive tree traversal using `processNode` functions
- Documentation comments extracted when available (JSDoc, Go doc comments, Python docstrings, Javadoc)
- Go parser handles methods with receivers, struct fields, and interface methods
- Java parser extracts classes, interfaces, enums with modifiers, inheritance, and member visibility
- TypeScript parser includes type annotations and extends/implements clauses
- Python parser filters out private symbols (names starting with underscore)
- All parsers generate readable outline format with proper indentation
- Memory management: Always use `defer parser.Close()` and `defer tree.Close()`

## CLI Usage

Primary usage as a command-line tool:

```bash
# Analyze a file
outline path/to/file.go

# Override language detection
outline --language go path/to/file.txt
```

## MCP Integration (Optional)

Can optionally run as MCP server:

- Uses stdio transport for communication with MCP clients
- Registers single "outline" tool accepting `file` parameter
- Returns structured text outlines of code symbols
- Handles errors gracefully with proper MCP error responses
- Compatible with Claude Desktop, MCP Inspector, and other MCP clients

## Testing

Each language has comprehensive test coverage in `pkg/outline/languages/*_test.go`:

- Import statements and module systems
- Function/method declarations with parameters and return types
- Class/struct declarations with inheritance
- Documentation comment extraction
- Access modifier filtering (public/private visibility)
- Error handling for malformed code

Run language-specific tests: `go test ./pkg/outline/languages/ -v`

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - Official MCP Go SDK (for MCP mode only)
- `github.com/tree-sitter/go-tree-sitter` - Core tree-sitter Go bindings  
- Language-specific tree-sitter grammars:
  - `github.com/tree-sitter/tree-sitter-go`
  - `github.com/tree-sitter/tree-sitter-java`
  - `github.com/tree-sitter/tree-sitter-javascript` 
  - `github.com/tree-sitter/tree-sitter-typescript`
  - `github.com/tree-sitter/tree-sitter-python`