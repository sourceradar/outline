package server

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sourceradar/outline/internal/detector"
	"github.com/sourceradar/outline/pkg/outline"
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
	language, ok := detector.DetectLanguage(filePath)
	if !ok {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: unsupported file extension"),
				},
			},
			IsError: true,
		}, nil
	}

	// Extract symbols based on language
	result, err := outline.ExtractOutline(content, language)
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

	formattedResult := fmt.Sprintf("Language: %s\n\n%s", language, result)
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: formattedResult,
			},
		},
	}, nil
}
