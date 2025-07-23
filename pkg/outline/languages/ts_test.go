package languages

import (
	"strings"
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

func TestTypeScriptOutlineWithImportsAndRequires(t *testing.T) {
	tsCode := `import * as React from 'react';
import { Component } from 'react';
const lodash = require('lodash');
const utils = require('./utils');

interface User {
    name: string;
    age?: number;
}

type Status = 'active' | 'inactive';

/**
 * Sample function with type annotations
 */
function processUser(user: User): Status {
    return 'active';
}

const processData = (data: any[]): User[] => {
    return data;
};

class UserManager implements Manager {
    private users: User[] = [];

    addUser(user: User): void {
        this.users.push(user);
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(typescript.LanguageTypescript())); err != nil {
		t.Fatalf("Failed to set TypeScript language: %v", err)
	}

	tree := parser.Parse([]byte(tsCode), nil)
	defer tree.Close()

	result := ExtractTSOutline(tree.RootNode(), []byte(tsCode))

	// Check that imports are included
	if !strings.Contains(result, "import * as React from 'react'") {
		t.Error("Expected namespace import to be included")
	}
	if !strings.Contains(result, "import { Component } from 'react'") {
		t.Error("Expected named import to be included")
	}

	// Check that require statements are included
	if !strings.Contains(result, "const lodash = require('lodash')") {
		t.Error("Expected require statement to be included")
	}
	if !strings.Contains(result, "const utils = require('./utils')") {
		t.Error("Expected local require statement to be included")
	}

	// Check that interface is included
	if !strings.Contains(result, "interface User {") {
		t.Error("Expected interface declaration to be included")
	}

	// Check that type alias is included
	if !strings.Contains(result, "type Status = 'active' | 'inactive'") {
		t.Error("Expected type alias to be included")
	}

	// Check that function with types is included
	if !strings.Contains(result, "function processUser(user: User): Status") {
		t.Error("Expected typed function to be included")
	}

	// Check that arrow function with types is included
	if !strings.Contains(result, "const processData = (data: any[]): User[]") {
		t.Error("Expected typed arrow function to be included")
	}

	// Check that class with implements is included
	if !strings.Contains(result, "class UserManager implements Manager") {
		t.Error("Expected class with implements to be included")
	}
}
