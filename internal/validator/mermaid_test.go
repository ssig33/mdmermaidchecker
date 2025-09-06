package validator

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ssig33/mdmermaidchecker/internal/parser"
)

// mockCommand creates a mock exec.Cmd for testing
func mockCommand(exitCode int, stderr string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{
			"GO_WANT_HELPER_PROCESS=1",
			"EXIT_CODE=" + string(rune(exitCode+'0')),
			"STDERR_OUTPUT=" + stderr,
		}
		return cmd
	}
}

// TestHelperProcess is a helper process for testing command execution
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	exitCode := 0
	if code := os.Getenv("EXIT_CODE"); code != "" {
		exitCode = int(code[0] - '0')
	}

	if stderr := os.Getenv("STDERR_OUTPUT"); stderr != "" {
		os.Stderr.WriteString(stderr)
	}

	os.Exit(exitCode)
}

func TestValidateMermaid(t *testing.T) {
	tests := []struct {
		name          string
		block         parser.MermaidBlock
		exitCode      int
		stderr        string
		expectError   bool
		errorContains string
	}{
		{
			name: "valid mermaid diagram",
			block: parser.MermaidBlock{
				Content:    "graph TD\n    A[Start] --> B[End]",
				LineNumber: 10,
			},
			exitCode:    0,
			stderr:      "",
			expectError: false,
		},
		{
			name: "invalid mermaid diagram",
			block: parser.MermaidBlock{
				Content:    "graph TD\n    A[Start --> B[End]",
				LineNumber: 20,
			},
			exitCode:      1,
			stderr:        "Error: Parse error on line 2",
			expectError:   true,
			errorContains: "Line 20: Mermaid validation failed",
		},
		{
			name: "empty mermaid block",
			block: parser.MermaidBlock{
				Content:    "",
				LineNumber: 5,
			},
			exitCode:      1,
			stderr:        "Error: No diagram definition found",
			expectError:   true,
			errorContains: "Line 5: Mermaid validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				CommandExecutor: mockCommand(tt.exitCode, tt.stderr),
			}

			err := v.ValidateMermaid(tt.block)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error should contain '%s', got '%s'", tt.errorContains, err.Error())
				}
				if !strings.Contains(err.Error(), tt.stderr) {
					t.Errorf("Error should contain stderr '%s', got '%s'", tt.stderr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateAll(t *testing.T) {
	blocks := []parser.MermaidBlock{
		{Content: "graph TD\n    A --> B", LineNumber: 10},
		{Content: "invalid", LineNumber: 20},
		{Content: "sequenceDiagram\n    Alice->>Bob: Hello", LineNumber: 30},
	}

	v := &Validator{
		CommandExecutor: func(name string, args ...string) *exec.Cmd {
			// Simulate: first block passes, second fails, third passes
			content := args[3] // The input file path (after --yes @mermaid-js/mermaid-cli -i)
			if strings.Contains(content, "_20.mmd") {
				return mockCommand(1, "Error: Invalid syntax")(name, args...)
			}
			return mockCommand(0, "")(name, args...)
		},
	}

	errors := v.ValidateAll(blocks)

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	if len(errors) > 0 && !strings.Contains(errors[0].Error(), "Line 20") {
		t.Errorf("Error should reference Line 20, got: %v", errors[0])
	}
}
