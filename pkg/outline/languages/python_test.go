package languages

import (
	"strings"
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

func TestPythonOutlineWithImports(t *testing.T) {
	pythonCode := `import os
import sys
from typing import List, Optional
from .local_module import helper

def public_function(name: str) -> str:
    """A public function."""
    return f"Hello {name}"

def _private_function():
    """This should not appear."""
    pass

class PublicClass:
    """A public class."""
    
    def __init__(self, value: int):
        self.value = value
    
    def public_method(self) -> int:
        """A public method."""
        return self.value

class _PrivateClass:
    """This should not appear."""
    pass
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(python.Language())); err != nil {
		t.Fatalf("Failed to set Python language: %v", err)
	}

	tree := parser.Parse([]byte(pythonCode), nil)
	defer tree.Close()

	result := ExtractPythonOutline(tree.RootNode(), []byte(pythonCode))

	// Check that imports are included
	if !strings.Contains(result, "import os") {
		t.Error("Expected simple import to be included")
	}
	if !strings.Contains(result, "from typing import List, Optional") {
		t.Error("Expected from import to be included")
	}
	if !strings.Contains(result, "from .local_module import helper") {
		t.Error("Expected relative import to be included")
	}

	// Check that public function is included
	if !strings.Contains(result, "def public_function(name: str) -> str:") {
		t.Error("Expected public function to be included")
	}

	// Check that private function is NOT included
	if strings.Contains(result, "_private_function") {
		t.Error("Private function should not be included")
	}

	// Check that public class is included
	if !strings.Contains(result, "class PublicClass:") {
		t.Error("Expected public class to be included")
	}

	// Check that public method is included
	if !strings.Contains(result, "def public_method(self) -> int:") {
		t.Error("Expected public method to be included")
	}

	// Check that private class is NOT included
	if strings.Contains(result, "_PrivateClass") {
		t.Error("Private class should not be included")
	}
}