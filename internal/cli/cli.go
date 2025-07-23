package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/sourceradar/outline/internal/detector"
	"github.com/sourceradar/outline/pkg/outline"
)

// Run executes the CLI application
func Run(args []string, languageOverride string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: outline [--language <lang>] <file>")
	}

	filePath := args[0]

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %v", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("expected a file, got directory")
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Determine language
	var language string
	if languageOverride != "" {
		language = languageOverride
	} else {
		var ok bool
		language, ok = detector.DetectLanguage(filePath)
		if !ok {
			supportedExts := strings.Join(detector.SupportedExtensions(), ", ")
			return fmt.Errorf("unsupported file extension. Supported extensions: %s\nOr use --language flag to override", supportedExts)
		}
	}

	// Extract outline
	result, err := outline.ExtractOutline(content, language)
	if err != nil {
		return fmt.Errorf("error extracting outline: %v", err)
	}

	fmt.Printf("Language: %s\n\n%s", language, result)
	return nil
}