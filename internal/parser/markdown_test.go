package parser

import (
	"testing"
)

func TestExtractMermaidBlocks(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []MermaidBlock
	}{
		{
			name: "single mermaid block",
			content: `# Test Document

Some text here.

` + "```mermaid" + `
graph TD
    A[Start] --> B[End]
` + "```" + `

More text.`,
			expected: []MermaidBlock{
				{
					Content: `graph TD
    A[Start] --> B[End]`,
					LineNumber: 5,
				},
			},
		},
		{
			name: "multiple mermaid blocks",
			content: `# Document

` + "```mermaid" + `
graph LR
    A --> B
` + "```" + `

Some text.

` + "```mermaid" + `
sequenceDiagram
    Alice->>Bob: Hello
` + "```",
			expected: []MermaidBlock{
				{
					Content: `graph LR
    A --> B`,
					LineNumber: 3,
				},
				{
					Content: `sequenceDiagram
    Alice->>Bob: Hello`,
					LineNumber: 10,
				},
			},
		},
		{
			name: "no mermaid blocks",
			content: `# Document

Just regular text.

` + "```javascript" + `
console.log("hello");
` + "```",
			expected: []MermaidBlock{},
		},
		{
			name:     "empty document",
			content:  "",
			expected: []MermaidBlock{},
		},
		{
			name: "mermaid block with indented content",
			content: "```mermaid" + `
    graph TD
        A[Start] --> B{Decision}
        B -->|Yes| C[OK]
        B -->|No| D[End]
` + "```",
			expected: []MermaidBlock{
				{
					Content: `    graph TD
        A[Start] --> B{Decision}
        B -->|Yes| C[OK]
        B -->|No| D[End]`,
					LineNumber: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := ExtractMermaidBlocks(tt.content)

			if len(blocks) != len(tt.expected) {
				t.Errorf("Expected %d blocks, got %d", len(tt.expected), len(blocks))
				return
			}

			for i, block := range blocks {
				if block.Content != tt.expected[i].Content {
					t.Errorf("Block %d content mismatch.\nExpected:\n%s\nGot:\n%s",
						i, tt.expected[i].Content, block.Content)
				}
				if block.LineNumber != tt.expected[i].LineNumber {
					t.Errorf("Block %d line number mismatch. Expected %d, got %d",
						i, tt.expected[i].LineNumber, block.LineNumber)
				}
			}
		})
	}
}
