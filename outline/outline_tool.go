package outline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)



// OutlineToolParams defines the parameters for the outline tool
type OutlineToolParams struct {
	File string `json:"file" jsonschema:"description=Path to the file to analyze"`
}

// OutlineToolHandler handles outline tool requests
func OutlineToolHandler(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[OutlineToolParams]) (*mcp.CallToolResultFor[any], error) {
	filePath := params.Arguments.File

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: file not found: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	if fileInfo.IsDir() {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: expected a file, got directory",
				},
			},
			IsError: true,
		}, nil
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error reading file: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Detect language based on file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	var language string

	switch ext {
	case ".go":
		language = "go"
	case ".js":
		language = "javascript"
	case ".jsx":
		language = "javascript"
	case ".ts":
		language = "typescript"
	case ".tsx":
		language = "tsx"
	case ".py":
		language = "python"
	default:
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: unsupported file extension: %s", ext),
				},
			},
			IsError: true,
		}, nil
	}

	// Extract symbols based on language
	outline, err := ExtractOutline(content, language)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error extracting outline: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	result := fmt.Sprintf("Language: %s\n\n%s", language, outline)
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result,
			},
		},
	}, nil
}
