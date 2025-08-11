package detector

// LanguageInfo contains metadata about a supported language
type LanguageInfo struct {
	Name        string
	Extensions  []string
	Description string
}

// SupportedLanguages returns a map of language name to LanguageInfo
// This is the single source of truth for all supported languages
func SupportedLanguages() map[string]LanguageInfo {
	return map[string]LanguageInfo{
		"go": {
			Name:        "go",
			Extensions:  []string{".go"},
			Description: "Go programming language",
		},
		"java": {
			Name:        "java",
			Extensions:  []string{".java"},
			Description: "Java programming language",
		},
		"javascript": {
			Name:        "javascript",
			Extensions:  []string{".js", ".jsx"},
			Description: "JavaScript programming language",
		},
		"typescript": {
			Name:        "typescript",
			Extensions:  []string{".ts"},
			Description: "TypeScript programming language",
		},
		"tsx": {
			Name:        "tsx",
			Extensions:  []string{".tsx"},
			Description: "TypeScript JSX",
		},
		"python": {
			Name:        "python",
			Extensions:  []string{".py"},
			Description: "Python programming language",
		},
		"swift": {
			Name:        "swift",
			Extensions:  []string{".swift"},
			Description: "Swift programming language",
		},
		"c": {
			Name:        "c",
			Extensions:  []string{".c", ".h"},
			Description: "C programming language",
		},
		"cpp": {
			Name:        "cpp",
			Extensions:  []string{".cpp", ".cxx", ".cc", ".hpp", ".hxx", ".hh"},
			Description: "C++ programming language",
		},
	}
}

// GetLanguageNames returns a slice of supported language names
func GetLanguageNames() []string {
	languages := SupportedLanguages()
	names := make([]string, 0, len(languages))
	for name := range languages {
		names = append(names, name)
	}
	return names
}

// GetAllExtensions returns all supported file extensions
func GetAllExtensions() []string {
	languages := SupportedLanguages()
	var extensions []string
	for _, lang := range languages {
		extensions = append(extensions, lang.Extensions...)
	}
	return extensions
}

// GetLanguageDisplayNames returns language names formatted for display
func GetLanguageDisplayNames() []string {
	languages := SupportedLanguages()
	names := make([]string, 0, len(languages))
	for _, lang := range languages {
		names = append(names, lang.Name)
	}
	return names
}
