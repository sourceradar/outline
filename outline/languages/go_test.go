package languages

import (
	"strings"
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	golang "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

func TestGoOutlineWithImports(t *testing.T) {
	goCode := `package main

import (
	"fmt"
	"os"
	"github.com/example/pkg"
)

// MyFunction does something
func MyFunction(name string) string {
	return "Hello " + name
}

type MyStruct struct {
	Name string
	Age  int
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(golang.Language())); err != nil {
		t.Fatalf("Failed to set Go language: %v", err)
	}

	tree := parser.Parse([]byte(goCode), nil)
	defer tree.Close()

	result := ExtractGoOutline(tree.RootNode(), []byte(goCode))

	// Check that package is included
	if !strings.Contains(result, "package main") {
		t.Error("Expected package declaration to be included")
	}

	// Check that imports are included
	if !strings.Contains(result, "import (") {
		t.Error("Expected import block to be included")
	}
	if !strings.Contains(result, `"fmt"`) {
		t.Error("Expected fmt import to be included")
	}
	if !strings.Contains(result, `"github.com/example/pkg"`) {
		t.Error("Expected external import to be included")
	}

	// Check that function is included
	if !strings.Contains(result, "func MyFunction(name string) string") {
		t.Error("Expected function declaration to be included")
	}

	// Check that struct is included
	if !strings.Contains(result, "type MyStruct struct") {
		t.Error("Expected struct declaration to be included")
	}
}