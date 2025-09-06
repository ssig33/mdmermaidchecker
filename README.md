# mdmermaidchecker

A CLI tool written in Go that validates Mermaid diagram syntax within Markdown files.

## Prerequisites

- Go 1.19 or later
- Node.js and npm (for running `npx`)
- `@mermaid-js/mermaid-cli` will be automatically installed on first run

## Installation

```bash
go install github.com/ssig33/mdmermaidchecker@latest
```

Or build from source:

```bash
git clone https://github.com/ssig33/mdmermaidchecker.git
cd mdmermaidchecker
go build
```

## Usage

```bash
mdmermaidchecker <markdown-file>
```

Example:

```bash
mdmermaidchecker README.md
```

## Exit Codes

- `0`: All Mermaid diagrams are valid (or no Mermaid diagrams found)
- `1`: One or more Mermaid diagrams have validation errors
- `2`: System error (file not found, etc.)

## How It Works

1. Parses the Markdown file to extract all Mermaid code blocks (` ```mermaid` blocks)
2. Validates each Mermaid diagram using `@mermaid-js/mermaid-cli`
3. Reports any validation errors with line numbers
4. Returns appropriate exit code

## Error Output

When validation errors are found, the tool outputs the line number and the error message from the Mermaid CLI:

```
Line 42: Mermaid validation failed:
Error: Parse error on line 2:
...
```

## Testing

Run all tests:

```bash
go test ./...
```

Run E2E tests (requires Node.js):

```bash
go test ./test -v
```

## License

MIT