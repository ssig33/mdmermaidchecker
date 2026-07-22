package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ssig33/mdmermaidchecker/internal/parser"
	"github.com/ssig33/mdmermaidchecker/internal/validator"
)

// ValidateMarkdownFile validates all mermaid blocks in a markdown file
func ValidateMarkdownFile(filepath string) int {
	// Read file
	content, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: File not found: %s\n", filepath)
			return 2
		}
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return 2
	}

	// Extract mermaid blocks
	blocks := parser.ExtractMermaidBlocks(string(content))

	if len(blocks) == 0 {
		// No mermaid blocks found - this is considered success
		return 0
	}

	// mermaid-cli depends on puppeteer, which needs unzip to extract Chrome
	if _, err := exec.LookPath("unzip"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: unzip is required but not found in PATH. Please install unzip to use mermaid-cli.")
		return 1
	}

	// Validate all blocks
	v := validator.NewValidator()
	errors := v.ValidateAll(blocks)

	// Print errors
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}

	// Return appropriate exit code
	if len(errors) > 0 {
		return 1
	}

	return 0
}
