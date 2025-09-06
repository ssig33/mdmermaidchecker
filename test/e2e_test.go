package test

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ssig33/mdmermaidchecker/cmd"
)

func TestE2E(t *testing.T) {
	// Check if Node.js is available
	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("Skipping E2E test: npx not found in PATH")
	}

	_, currentFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(currentFile)

	tests := []struct {
		name         string
		file         string
		expectedExit int
	}{
		{
			name:         "valid markdown file",
			file:         filepath.Join(testDir, "testdata/valid.md"),
			expectedExit: 0,
		},
		{
			name:         "invalid markdown file",
			file:         filepath.Join(testDir, "testdata/invalid.md"),
			expectedExit: 1,
		},
		{
			name:         "non-existent file",
			file:         filepath.Join(testDir, "testdata/nonexistent.md"),
			expectedExit: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := cmd.ValidateMarkdownFile(tt.file)

			if exitCode != tt.expectedExit {
				t.Errorf("Expected exit code %d, got %d", tt.expectedExit, exitCode)
			}
		})
	}
}
