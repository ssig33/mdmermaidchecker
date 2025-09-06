package parser

import (
	"strings"
)

// MermaidBlock represents a mermaid code block found in markdown
type MermaidBlock struct {
	Content    string
	LineNumber int
}

// ExtractMermaidBlocks extracts all mermaid code blocks from markdown content
func ExtractMermaidBlocks(content string) []MermaidBlock {
	var blocks []MermaidBlock
	lines := strings.Split(content, "\n")

	inMermaidBlock := false
	var currentBlock strings.Builder
	var blockStartLine int

	for i, line := range lines {
		lineNumber := i + 1

		if !inMermaidBlock {
			// Check for mermaid block start
			if isMermaidStart(line) {
				inMermaidBlock = true
				blockStartLine = lineNumber
				currentBlock.Reset()
			}
		} else {
			// Check for code block end
			if isCodeBlockEnd(line) {
				blocks = append(blocks, MermaidBlock{
					Content:    currentBlock.String(),
					LineNumber: blockStartLine,
				})
				inMermaidBlock = false
			} else {
				if currentBlock.Len() > 0 {
					currentBlock.WriteString("\n")
				}
				currentBlock.WriteString(line)
			}
		}
	}

	return blocks
}

// isMermaidStart checks if a line starts a mermaid code block
func isMermaidStart(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "```mermaid"
}

// isCodeBlockEnd checks if a line ends a code block
func isCodeBlockEnd(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "```"
}
