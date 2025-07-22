# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an MCP (Model Context Protocol) server that provides code outline generation for multiple programming languages using tree-sitter parsers. The server analyzes source code files and generates structured outlines showing functions, classes, types, and other symbols.

## Supported Languages

- **Go** (.go files) - Functions, methods, types, constants, variables
- **JavaScript** (.js, .jsx files) - Functions, classes, arrow functions  
- **TypeScript** (.ts, .tsx files) - Functions, classes, interfaces, types
- **Python** (.py files) - Functions, classes (public symbols only)

## Architecture

### Core Components

- `main.go` - MCP server entry point with tool registration and stdio transport
- `outline/outline.go` - Main outline extraction logic with language detection and parser creation
- `outline/outline_tool.go` - MCP tool handler implementing the outline functionality
- `outline/util.go` - Utility functions for tree-sitter node processing
- `outline/languages/` - Language-specific outline extractors:
  - `go.go` - Go language parser with struct/interface/method handling
  - `js.go` - JavaScript parser with class and function extraction  
  - `ts.go` - TypeScript parser with type annotations and interfaces
  - `python.go` - Python parser filtering private symbols (underscore prefix)

### Key Functions

- `ExtractOutline(content []byte, language string)` - Main entry point in `outline/outline.go:27`
- `createParserForLanguage(language string)` - Parser factory in `outline/outline.go:52`
- `OutlineToolHandler()` - MCP tool handler in `outline/outline_tool.go:12`
- `getNodeText()` and `findDocComment()` - Utility functions in `outline/util.go`
- Language-specific extractors: `extractGoOutline()`, `extractJSOutline()`, `extractTSOutline()`, `extractPythonOutline()`

## Development Commands

```bash
# Build the project
go build .

# Run the MCP server (for development with MCP inspector)
npx @modelcontextprotocol/inspector go run .

# Run tests (if any exist)
go test ./...

# Install dependencies and clean up
go mod tidy

# Format code
go fmt ./...
```

## Code Patterns

- Each language parser follows a recursive tree traversal pattern
- Documentation comments are extracted when available (JSDoc, Go doc comments, Python docstrings)
- Go parser handles methods with receivers, struct fields, and interface methods
- TypeScript parser includes type annotations and extends/implements clauses
- Python parser filters out private symbols (names starting with underscore)
- All parsers generate readable outline format with indentation and structure

## MCP Integration

This server implements the Model Context Protocol (MCP) specification:

- Uses stdio transport for communication with MCP clients
- Registers a single "outline" tool that accepts a file path parameter
- Returns structured text outlines of code symbols
- Handles errors gracefully with proper MCP error responses

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - Official MCP Go SDK for server implementation
- `github.com/tree-sitter/go-tree-sitter` - Core tree-sitter Go bindings
- Language-specific tree-sitter grammars for Go, JavaScript, TypeScript, Python