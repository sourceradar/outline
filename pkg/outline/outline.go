package outline

import (
	"fmt"
	sitter "github.com/tree-sitter/go-tree-sitter"
	golang "github.com/tree-sitter/tree-sitter-go/bindings/go"
	java "github.com/tree-sitter/tree-sitter-java/bindings/go"
	javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	python "github.com/tree-sitter/tree-sitter-python/bindings/go"
	swift "github.com/alex-pinkus/tree-sitter-swift/bindings/go"
	typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"

	"github.com/sourceradar/outline/pkg/outline/languages"
)

// SymbolInfo represents information about a code symbol (for internal use)
type SymbolInfo struct {
	Type          string       `json:"type"`
	Name          string       `json:"name"`
	Signature     string       `json:"signature,omitempty"`
	Documentation string       `json:"documentation,omitempty"`
	Line          int          `json:"line"`
	Column        int          `json:"column"`
	EndLine       int          `json:"endLine"`
	EndColumn     int          `json:"endColumn"`
	IsPublic      bool         `json:"isPublic"`
	Children      []SymbolInfo `json:"children,omitempty"`
}

// ExtractOutline analyzes the syntax tree to generate a compact outline
func ExtractOutline(content []byte, language string) (string, error) {
	// Parse content
	parser, err := createParserForLanguage(language)

	if err != nil {
		return "", fmt.Errorf("error creating parser: %v", err)
	}

	tree := parser.Parse(content, nil)
	root := tree.RootNode()

	switch language {
	case "go":
		return languages.ExtractGoOutline(root, content), nil
	case "java":
		return languages.ExtractJavaOutline(root, content), nil
	case "javascript":
		return languages.ExtractJSOutline(root, content), nil
	case "swift":
		return languages.ExtractSwiftOutline(root, content), nil
	case "typescript":
		return languages.ExtractTSOutline(root, content), nil
	case "python":
		return languages.ExtractPythonOutline(root, content), nil
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}

func createParserForLanguage(language string) (*sitter.Parser, error) {
	var err error
	parser := sitter.NewParser()

	switch language {
	case "go":
		err = parser.SetLanguage(sitter.NewLanguage(golang.Language()))
	case "java":
		err = parser.SetLanguage(sitter.NewLanguage(java.Language()))
	case "javascript":
		err = parser.SetLanguage(sitter.NewLanguage(javascript.Language()))
	case "swift":
		err = parser.SetLanguage(sitter.NewLanguage(swift.Language()))
	case "typescript":
		err = parser.SetLanguage(sitter.NewLanguage(typescript.LanguageTypescript()))
	case "tsx":
		err = parser.SetLanguage(sitter.NewLanguage(typescript.LanguageTSX()))
	case "python":
		language = "python"
		err = parser.SetLanguage(sitter.NewLanguage(python.Language()))
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	if err != nil {
		return nil, fmt.Errorf("error setting language parser: %v", err)
	}

	return parser, nil
}
