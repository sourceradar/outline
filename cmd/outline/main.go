package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sourceradar/outline/internal/cli"
	"github.com/sourceradar/outline/internal/detector"
	"github.com/sourceradar/outline/internal/server"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var mcpMode bool
	var language string
	var help bool
	var showVersion bool

	flag.BoolVar(&mcpMode, "mcp", false, "Run in MCP server mode")
	flag.StringVar(&language, "language", "", fmt.Sprintf("Override language detection (%s)", strings.Join(detector.GetLanguageNames(), ", ")))
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.BoolVar(&help, "h", false, "Show help message")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")

	flag.Usage = func() {
		supportedLangs := strings.Join(detector.GetLanguageNames(), ", ")
		fmt.Fprintf(os.Stderr, `outline - A code analysis tool that generates structured outlines

USAGE:
    outline [OPTIONS] <file>
    outline --mcp

OPTIONS:
    --language <lang>   Override language detection
                        Supported: %s
    --mcp               Run in MCP (Model Context Protocol) server mode
    --version, -v       Show version information
    --help, -h          Show this help message

EXAMPLES:
    outline main.go                      # Analyze a Go file
    outline --language go script.txt     # Force Go parsing
    outline --mcp                        # Run as MCP server
    outline --version                    # Show version

For MCP server mode, add to your MCP client configuration:
{
  "mcpServers": {
    "outline": {
      "command": "outline",
      "args": ["--mcp"]
    }
  }
}
`, supportedLangs)
	}

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Printf("outline version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", date)
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
