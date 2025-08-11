package detector

import (
	"path/filepath"
	"strings"
)

// DetectLanguage determines the programming language based on file extension
func DetectLanguage(filePath string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filePath))

	languages := SupportedLanguages()
	for langName, langInfo := range languages {
		for _, supportedExt := range langInfo.Extensions {
			if ext == supportedExt {
				return langName, true
			}
		}
	}

	return "", false
}

// SupportedExtensions returns a list of supported file extensions
func SupportedExtensions() []string {
	return GetAllExtensions()
}
