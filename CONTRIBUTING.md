# Contributing to Outline

Thank you for your interest in contributing to Outline! This document provides guidelines for contributing to the project, including how to add support for new programming languages.

## General Contributing Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass (`go test ./...`)
6. Format your code (`go fmt ./...`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## Adding New Language Support

This section explains how to add support for new programming languages to the outline tool.

## Overview

The outline tool uses tree-sitter parsers to analyze source code and generate structured outlines. Each language requires:

1. A tree-sitter grammar/parser
2. Language-specific outline extraction logic
3. Integration with the main parser factory
4. Comprehensive tests

## Quick Start

To add support for a new language (e.g., Rust):

1. **Add tree-sitter dependency** to `go.mod`
2. **Create extractor** in `pkg/outline/languages/rust.go`
3. **Update parser factory** in `pkg/outline/outline.go`
4. **Add file extension mapping** in `internal/detector/`
5. **Write tests** in `pkg/outline/languages/rust_test.go`

## Step-by-Step Guide

### 1. Add Tree-Sitter Dependency

Add the tree-sitter grammar for your language to `go.mod`:

```bash
go get github.com/tree-sitter/tree-sitter-rust/bindings/go
```

Import it in `pkg/outline/outline.go`:

```go
import (
    // ... existing imports
    rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)
```

### 2. Create Language Extractor

Create a new file `pkg/outline/languages/{language}.go` with the following structure:

```go
package languages

import (
    "fmt"
    "strings"
    sitter "github.com/tree-sitter/go-tree-sitter"
)

// Extract{Language}Outline extracts {Language} outline directly from the code
func Extract{Language}Outline(root *sitter.Node, content []byte) string {
    var result strings.Builder

    // Function to process a node and its children
    var processNode func(node *sitter.Node, indentLevel int)
    processNode = func(node *sitter.Node, indentLevel int) {
        indent := strings.Repeat("    ", indentLevel) // Adjust indentation as needed

        // Process based on node type
        switch node.Kind() {
        case "source_file", "program", "module":
            // Root node - process all children
            for i := 0; i < int(node.NamedChildCount()); i++ {
                child := node.NamedChild(uint(i))
                processNode(child, indentLevel)
            }

        case "import_statement", "use_declaration":
            // Handle import statements
            importText := getNodeText(node, content)
            result.WriteString(fmt.Sprintf("%s\n", importText))

        case "function_declaration", "function_item":
            // Handle function declarations
            processFunctionDeclaration(node, content, &result, indent)

        case "struct_declaration", "struct_item":
            // Handle struct/class declarations
            processStructDeclaration(node, content, &result, indent, indentLevel)

        // Add more cases for language-specific constructs
        }
    }

    processNode(root, 0)
    return result.String()
}

func processFunctionDeclaration(node *sitter.Node, content []byte, result *strings.Builder, indent string) {
    nameNode := node.ChildByFieldName("name")
    if nameNode == nil {
        return
    }

    name := getNodeText(nameNode, content)
    
    // Get parameters
    paramNode := node.ChildByFieldName("parameters")
    paramText := ""
    if paramNode != nil {
        paramText = getNodeText(paramNode, content)
    }

    // Get return type if applicable
    returnNode := node.ChildByFieldName("return_type")
    returnText := ""
    if returnNode != nil {
        returnText = getNodeText(returnNode, content)
    }

    // Get documentation comment if present
    doc := findDocComment(node, content, "rust") // Replace with your language
    if doc != "" {
        docLines := strings.Split(doc, "\n")
        for _, line := range docLines {
            result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
        }
    }

    // Write function declaration with placeholder body
    result.WriteString(fmt.Sprintf("%sfn %s%s%s {\n", indent, name, paramText, returnText))
    result.WriteString(fmt.Sprintf("%s    // ...\n", indent))
    result.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

// Add more helper functions as needed
```

### 3. Update Parser Factory

In `pkg/outline/outline.go`, add your language to the `createParserForLanguage` function:

```go
func createParserForLanguage(language string) (*sitter.Parser, error) {
    var err error
    parser := sitter.NewParser()

    switch language {
    // ... existing cases
    case "rust":
        err = parser.SetLanguage(sitter.NewLanguage(rust.Language()))
    default:
        return nil, fmt.Errorf("unsupported language: %s", language)
    }

    if err != nil {
        return nil, fmt.Errorf("error setting language parser: %v", err)
    }

    return parser, nil
}
```

And add the extraction case in `ExtractOutline`:

```go
func ExtractOutline(content []byte, language string) (string, error) {
    // ... parser creation logic

    switch language {
    // ... existing cases
    case "rust":
        return languages.ExtractRustOutline(root, content), nil
    default:
        return "", fmt.Errorf("unsupported language: %s", language)
    }
}
```

### 4. Add File Extension Mapping

In `internal/detector/detector.go`, add file extension detection:

```go
func DetectLanguage(filePath string) (string, bool) {
    ext := strings.ToLower(filepath.Ext(filePath))
    switch ext {
    // ... existing cases
    case ".rs":
        return "rust", true
    default:
        return "", false
    }
}
```

### 5. Write Comprehensive Tests

Create `pkg/outline/languages/{language}_test.go`:

```go
package languages

import (
    "strings"
    "testing"
    
    sitter "github.com/tree-sitter/go-tree-sitter"
    rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)

func TestRustOutlineWithImports(t *testing.T) {
    rustCode := `use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Debug)]
struct User {
    name: String,
    age: u32,
}

impl User {
    fn new(name: String, age: u32) -> Self {
        User { name, age }
    }
    
    fn greet(&self) -> String {
        format!("Hello, {}", self.name)
    }
}

fn main() {
    let user = User::new("Alice".to_string(), 30);
    println!("{}", user.greet());
}
`

    parser := sitter.NewParser()
    defer parser.Close()

    if err := parser.SetLanguage(sitter.NewLanguage(rust.Language())); err != nil {
        t.Fatalf("Failed to set Rust language: %v", err)
    }

    tree := parser.Parse([]byte(rustCode), nil)
    defer tree.Close()

    result := ExtractRustOutline(tree.RootNode(), []byte(rustCode))

    // Test that imports are included
    if !strings.Contains(result, "use std::collections::HashMap") {
        t.Error("Expected use statement to be included")
    }
    if !strings.Contains(result, "use serde::{Deserialize, Serialize}") {
        t.Error("Expected grouped use statement to be included")
    }

    // Test that struct is included
    if !strings.Contains(result, "struct User {") {
        t.Error("Expected struct declaration to be included")
    }

    // Test that impl block methods are included
    if !strings.Contains(result, "fn new(") {
        t.Error("Expected constructor method to be included")
    }
    if !strings.Contains(result, "fn greet(&self) -> String") {
        t.Error("Expected instance method with return type to be included")
    }

    // Test that main function is included
    if !strings.Contains(result, "fn main()") {
        t.Error("Expected main function to be included")
    }
}
```

## Tree-Sitter Integration Patterns

### Common Node Types

Most tree-sitter grammars follow similar patterns:

- **Root nodes**: `source_file`, `program`, `module`
- **Import statements**: `import_statement`, `import_declaration`, `use_declaration`
- **Functions**: `function_declaration`, `function_definition`, `function_item`
- **Classes/Structs**: `class_declaration`, `struct_declaration`, `struct_item`
- **Variables**: `variable_declaration`, `let_declaration`
- **Comments**: `comment`, `line_comment`, `block_comment`

### Field Access Patterns

Use `node.ChildByFieldName()` to access semantic parts:

```go
nameNode := node.ChildByFieldName("name")
paramsNode := node.ChildByFieldName("parameters")
bodyNode := node.ChildByFieldName("body")
typeNode := node.ChildByFieldName("type")
returnTypeNode := node.ChildByFieldName("return_type")
```

### Heritage/Inheritance Patterns

For class inheritance, look for:

- `class_heritage` (JavaScript/TypeScript)
- `superclass` field
- `extends_clause`, `implements_clause`
- `base_clause` (Python)

### Debugging Tree-Sitter Parsing

Use this helper function to explore node structure:

```go
func debugNode(node *sitter.Node, content []byte, depth int) {
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%sNode: %s, Text: '%s'\n", indent, node.Kind(), getNodeText(node, content))
    
    for i := 0; i < int(node.ChildCount()); i++ {
        child := node.Child(uint(i))
        if child != nil {
            debugNode(child, content, depth+1)
        }
    }
}
```

## Language-Specific Considerations

### Access Modifiers

Handle visibility appropriately:

**Python**: Skip names starting with `_`
```go
if strings.HasPrefix(name, "_") {
    return // Skip private symbols
}
```

**Java/C#**: Check for `public`, `private`, `protected` modifiers

**Rust**: Check for `pub` keyword

### Documentation Comments

Different languages use different doc comment styles:

- **Go**: `//` comments above declarations
- **Rust**: `///` and `//!` comments
- **Java**: `/** */` JavaDoc comments
- **Python**: Triple-quoted strings as first statement in function/class body
- **JavaScript/TypeScript**: `/** */` JSDoc comments

### Type Annotations

Handle type information appropriately:

- Some languages have explicit type annotations in return_type fields
- Some include the `:` in the type annotation, others don't
- Generic types may need special handling

### Module Systems

Each language handles imports differently:

- **Go**: `import "package"` and `import ("multi", "line")`
- **Rust**: `use path::to::item;` and `use path::{item1, item2};`
- **Python**: `import module` and `from module import item`
- **JavaScript**: `import item from 'module'` and `const item = require('module')`

## Testing Requirements

### Minimum Test Coverage

Each language should have tests covering:

1. **Import statements** - All import/use/require variants
2. **Function declarations** - With and without parameters, return types
3. **Class/struct declarations** - With inheritance/implementations
4. **Method declarations** - Instance and static methods
5. **Variable declarations** - Constants and variables
6. **Documentation comments** - Language-specific doc comment formats
7. **Access modifiers** - Public/private visibility filtering
8. **Generic types** - If the language supports them
9. **Error handling** - Malformed code should not crash

### Test Structure

Follow this pattern:

```go
func Test{Language}OutlineWith{Feature}(t *testing.T) {
    // 1. Define test code string
    code := `...`
    
    // 2. Set up parser
    parser := sitter.NewParser()
    defer parser.Close()
    
    // 3. Parse code
    if err := parser.SetLanguage(sitter.NewLanguage(lang.Language())); err != nil {
        t.Fatalf("Failed to set language: %v", err)
    }
    tree := parser.Parse([]byte(code), nil)
    defer tree.Close()
    
    // 4. Extract outline
    result := Extract{Language}Outline(tree.RootNode(), []byte(code))
    
    // 5. Assert expected elements are present
    if !strings.Contains(result, "expected_content") {
        t.Error("Expected content to be included")
    }
    
    // 6. Assert unwanted elements are absent (e.g., private symbols)
    if strings.Contains(result, "private_content") {
        t.Error("Private content should not be included")
    }
}
```

## Performance Considerations

- **Lazy imports**: Only import tree-sitter grammars when needed
- **Parser reuse**: Consider caching parsers for frequently used languages
- **Memory management**: Always call `defer parser.Close()` and `defer tree.Close()`
- **Large files**: Consider streaming or chunking for very large source files

## Common Pitfalls

1. **Double colons in type annotations**: Check if tree-sitter already includes `:` in return_type
2. **Missing inheritance clauses**: Use `class_heritage` instead of `superclass` for JS/TS
3. **Wrong indent patterns**: Match the language's conventional indentation
4. **Forgetting to handle named vs unnamed children**: Use `NamedChild()` for semantic nodes
5. **Not handling edge cases**: Empty files, syntax errors, missing optional fields

## Contributing Guidelines

When adding a new language:

1. **Follow existing patterns** - Look at similar languages for guidance
2. **Write comprehensive tests** - Cover all major language features
3. **Update documentation** - Add your language to README.md and CLAUDE.md
4. **Handle errors gracefully** - Don't crash on malformed input
5. **Consider the end user** - Generate clean, readable outlines

## Examples

See the existing language implementations for reference:

- **Go** (`pkg/outline/languages/go.go`) - Complex example with packages, imports, functions, methods, structs, interfaces
- **Java** (`pkg/outline/languages/java.go`) - Classes, interfaces, enums, modifiers, inheritance, Javadoc
- **TypeScript** (`pkg/outline/languages/ts.go`) - Type annotations, interfaces, classes, heritage clauses
- **JavaScript** (`pkg/outline/languages/js.go`) - Classes, arrow functions, CommonJS requires
- **Python** (`pkg/outline/languages/python.go`) - Docstrings, private symbol filtering

Each implementation demonstrates different tree-sitter integration patterns and language-specific considerations.