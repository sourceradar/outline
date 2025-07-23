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
	
	flag.BoolVar(&mcpMode, "mcp", false, "Run in MCP server mode")
	flag.StringVar(&language, "language", "", "Override language detection (go, java, javascript, typescript, python)")
	flag.Parse()

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