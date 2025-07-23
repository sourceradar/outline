package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sourceradar/outline/internal/cli"
	"github.com/sourceradar/outline/internal/server"
)

func main() {
	var mcpMode bool
	var language string
	var help bool
	
	flag.BoolVar(&mcpMode, "mcp", false, "Run in MCP server mode")
	flag.StringVar(&language, "language", "", "Override language detection (go, java, javascript, typescript, python)")
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.BoolVar(&help, "h", false, "Show help message")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `outline - A code analysis tool that generates structured outlines

USAGE:
    outline [OPTIONS] <file>
    outline --mcp

OPTIONS:
    --language <lang>    Override language detection
                        Supported: go, java, javascript, typescript, python
    --mcp               Run in MCP (Model Context Protocol) server mode
    --help, -h          Show this help message

EXAMPLES:
    outline main.go                      # Analyze a Go file
    outline --language go script.txt     # Force Go parsing
    outline --mcp                        # Run as MCP server

For MCP server mode, add to your MCP client configuration:
{
  "mcpServers": {
    "outline": {
      "command": "outline",
      "args": ["--mcp"]
    }
  }
}
`)
	}
	
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if mcpMode {
		if err := server.Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := cli.Run(flag.Args(), language); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}