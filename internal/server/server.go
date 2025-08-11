package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sourceradar/outline/internal/detector"
)

// Run starts the MCP server
func Run() error {
	// Create server with implementation details
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "outline",
		Version: "1.0.0",
	}, nil)

	// Register the outline tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "outline",
		Description: getToolDescription(),
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"file": {
					Type:        "string",
					Description: "Path to the source code file to analyze",
				},
			},
			Required: []string{"file"},
		},
	}, OutlineToolHandler)

	// Run server using stdio transport
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}

	return nil
}

// getToolDescription generates the tool description with supported languages
func getToolDescription() string {
	languages := detector.SupportedLanguages()
	var langNames []string
	for _, lang := range languages {
		langNames = append(langNames, strings.Title(lang.Name))
	}

	return fmt.Sprintf("Extract a structured, high-level overview of code symbols from source files. Shows function signatures, class definitions, interfaces, types, and documentation comments without implementation details. Ideal for understanding code architecture, APIs, and large codebases quickly. Supports %s. More efficient than reading entire files when you need to understand code structure and available symbols.", strings.Join(langNames, ", "))
}
