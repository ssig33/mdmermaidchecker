package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestE2E(t *testing.T) {
	// Check if Node.js is available
	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("Skipping E2E test: npx not found in PATH")
	}

	// Build the binary
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	binaryPath := filepath.Join(projectRoot, "mdmermaidchecker")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, projectRoot)
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove(binaryPath)

	tests := []struct {
		name         string
		file         string
		expectedExit int
	}{
		{
			name:         "valid markdown file",
			file:         "testdata/valid.md",
			expectedExit: 0,
		},
		{
			name:         "invalid markdown file",
			file:         "testdata/invalid.md",
			expectedExit: 1,
		},
		{
			name:         "non-existent file",
			file:         "testdata/nonexistent.md",
			expectedExit: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.file)
			cmd.Dir = filepath.Dir(currentFile)

			err := cmd.Run()

			// Check exit code
			exitCode := 0
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else if err != nil {
				t.Fatalf("Unexpected error type: %v", err)
			}

			if exitCode != tt.expectedExit {
				t.Errorf("Expected exit code %d, got %d", tt.expectedExit, exitCode)
			}
		})
	}
}

func TestE2EUsage(t *testing.T) {
	// Build the binary
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	binaryPath := filepath.Join(projectRoot, "mdmermaidchecker")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, projectRoot)
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove(binaryPath)

	// Test with no arguments
	cmd := exec.Command(binaryPath)
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Errorf("Expected exit code 2 for no arguments, got %d", exitErr.ExitCode())
		}
	} else {
		t.Errorf("Expected exit error for no arguments")
	}

	// Test with too many arguments
	cmd = exec.Command(binaryPath, "file1.md", "file2.md")
	err = cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Errorf("Expected exit code 2 for too many arguments, got %d", exitErr.ExitCode())
		}
	} else {
		t.Errorf("Expected exit error for too many arguments")
	}
}
