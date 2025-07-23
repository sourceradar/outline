package detector

import (
	"path/filepath"
	"strings"
)

// DetectLanguage determines the programming language based on file extension
func DetectLanguage(filePath string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".go":
		return "go", true
	case ".java":
		return "java", true
	case ".js", ".jsx":
		return "javascript", true
	case ".ts":
		return "typescript", true
	case ".tsx":
		return "tsx", true
	case ".py":
		return "python", true
	default:
		return "", false
	}
}

// SupportedExtensions returns a list of supported file extensions
func SupportedExtensions() []string {
	return []string{".go", ".java", ".js", ".jsx", ".ts", ".tsx", ".py"}
}
