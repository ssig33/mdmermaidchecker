package validator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ssig33/mdmermaidchecker/internal/parser"
)

// Validator provides mermaid validation functionality
type Validator struct {
	// CommandExecutor can be overridden for testing
	CommandExecutor func(name string, arg ...string) *exec.Cmd
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		CommandExecutor: exec.Command,
	}
}

// ValidateMermaid validates a single mermaid block
func (v *Validator) ValidateMermaid(block parser.MermaidBlock) error {
	// Create temporary files
	tempDir := os.TempDir()
	tempMmdFile := filepath.Join(tempDir, fmt.Sprintf("mermaid_%d.mmd", block.LineNumber))
	tempSvgFile := filepath.Join(tempDir, fmt.Sprintf("mermaid_%d.svg", block.LineNumber))

	// Clean up temporary files
	defer func() {
		os.Remove(tempMmdFile)
		os.Remove(tempSvgFile)
	}()

	// Write mermaid content to temporary file
	if err := os.WriteFile(tempMmdFile, []byte(block.Content), 0644); err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Execute mmdc command
	args := []string{"--yes", "@mermaid-js/mermaid-cli", "-i", tempMmdFile, "-o", tempSvgFile}
	
	// Check for puppeteer config from environment variable (for CI)
	if puppeteerConfig := os.Getenv("PUPPETEER_CONFIG"); puppeteerConfig != "" {
		args = append(args, "-p", puppeteerConfig)
	}
	
	cmd := v.CommandExecutor("npx", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// mmdc validation failed
		return fmt.Errorf("Line %d: Mermaid validation failed:\n%s", block.LineNumber, stderr.String())
	}

	return nil
}

// ValidateAll validates all mermaid blocks and returns all errors
func (v *Validator) ValidateAll(blocks []parser.MermaidBlock) []error {
	var errors []error

	for _, block := range blocks {
		if err := v.ValidateMermaid(block); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
