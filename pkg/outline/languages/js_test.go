package languages

import (
	"strings"
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

func TestJavaScriptOutlineWithImportsAndRequires(t *testing.T) {
	jsCode := `import React from 'react';
import { useState } from 'react';
const fs = require('fs');
const path = require('path');

/**
 * A sample function
 */
function myFunction(param) {
    return param * 2;
}

const arrowFunc = (x) => {
    return x + 1;
};

class MyClass extends Component {
    constructor(props) {
        super(props);
    }

    render() {
        return null;
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(javascript.Language())); err != nil {
		t.Fatalf("Failed to set JavaScript language: %v", err)
	}

	tree := parser.Parse([]byte(jsCode), nil)
	defer tree.Close()

	result := ExtractJSOutline(tree.RootNode(), []byte(jsCode))

	// Check that ES6 imports are included
	if !strings.Contains(result, "import React from 'react'") {
		t.Error("Expected default import to be included")
	}
	if !strings.Contains(result, "import { useState } from 'react'") {
		t.Error("Expected named import to be included")
	}

	// Check that require statements are included
	if !strings.Contains(result, "const fs = require('fs')") {
		t.Error("Expected require statement to be included")
	}
	if !strings.Contains(result, "const path = require('path')") {
		t.Error("Expected require statement to be included")
	}

	// Check that function is included
	if !strings.Contains(result, "function myFunction(param)") {
		t.Error("Expected function declaration to be included")
	}

	// Check that arrow function is included
	if !strings.Contains(result, "const arrowFunc = (x) =>") {
		t.Error("Expected arrow function to be included")
	}

	// Check that class is included
	if !strings.Contains(result, "class MyClass extends Component") {
		t.Error("Expected class declaration to be included")
	}
}
