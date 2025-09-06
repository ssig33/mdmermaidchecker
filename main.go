package main

import (
	"fmt"
	"os"

	"github.com/ssig33/mdmermaidchecker/cmd"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <markdown-file>\n", os.Args[0])
		os.Exit(2)
	}

	filepath := os.Args[1]
	exitCode := cmd.ValidateMarkdownFile(filepath)
	os.Exit(exitCode)
}
