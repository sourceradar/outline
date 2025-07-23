package server

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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
		Description: "Generate an outline of symbols in a file (functions, classes, types, etc.)",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"file": {
					Type:        "string",
					Description: "Path to the file to analyze",
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