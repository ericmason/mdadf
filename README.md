# mdadf

A Go library that converts Markdown to Atlassian Document Format (ADF) for use with Jira and Confluence APIs.

## Installation

**Install the CLI tool:**
```bash
go install github.com/ericmason/mdadf/cmd/mdadf@latest
```

**Or install as a library:**
```bash
go get github.com/ericmason/mdadf
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	adf "github.com/ericmason/mdadf"
)

func main() {
	markdown := `# Hello World

This is **bold** and *italic* text with a [link](https://example.com).

- Item 1
- Item 2

` + "```go\nfmt.Println(\"Hello\")\n```"

	// Convert to JSON bytes
	jsonBytes, err := adf.Convert(markdown)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBytes))

	// Or convert to Document struct
	doc, err := adf.ConvertToDoc(markdown)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Document has %d top-level nodes\n", len(doc.Content))
}
```

## Supported Markdown Elements

| Markdown | ADF Output |
|----------|------------|
| `# Heading` | `heading` (levels 1-6) |
| Paragraph | `paragraph` |
| `**bold**` | `text` with `strong` mark |
| `*italic*` | `text` with `em` mark |
| `` `code` `` | `text` with `code` mark |
| `~~strike~~` | `text` with `strike` mark |
| `[link](url)` | `text` with `link` mark |
| `- item` | `bulletList` > `listItem` |
| `- [ ] item` | `taskList` > `taskItem` (`TODO`) |
| `- [x] item` | `taskList` > `taskItem` (`DONE`) |
| `1. item` | `orderedList` > `listItem` |
| ` ```lang ` | `codeBlock` with language |
| `> quote` | `blockquote` |
| `---` | `rule` |
| `\| table \|` | `table` > `tableRow` > `tableCell` |
| `![img](url)` | `text` with `link` mark (fallback) |

## ADF Output Example

Input:
```markdown
# Hello

This is **bold** text.
```

Output:
```json
{
  "version": 1,
  "type": "doc",
  "content": [
    {
      "type": "heading",
      "attrs": {"level": 1},
      "content": [{"type": "text", "text": "Hello"}]
    },
    {
      "type": "paragraph",
      "content": [
        {"type": "text", "text": "This is "},
        {"type": "text", "text": "bold", "marks": [{"type": "strong"}]},
        {"type": "text", "text": " text."}
      ]
    }
  ]
}
```

## Helper Functions

The library provides helper functions to programmatically build ADF documents:

```go
// Create nodes
para := adf.ParagraphNode(adf.TextNode("Hello"))
heading := adf.HeadingNode(1, adf.TextNode("Title"))
list := adf.BulletListNode(
	adf.ListItemNode(adf.ParagraphNode(adf.TextNode("Item 1"))),
	adf.ListItemNode(adf.ParagraphNode(adf.TextNode("Item 2"))),
)

// Create text with marks
bold := adf.TextNode("bold", adf.StrongMark())
italic := adf.TextNode("italic", adf.EmMark())
link := adf.TextNode("click here", adf.LinkMark("https://example.com", ""))
code := adf.TextNode("code", adf.CodeMark())

// Multiple marks
boldItalic := adf.TextNode("bold and italic", adf.StrongMark(), adf.EmMark())
```

## Command Line Tool

The repository includes a command-line tool `mdadf` for converting markdown files:

```bash
# Build the tool
go build -o mdadf ./cmd/mdadf

# Convert a markdown file
./mdadf input.md > output.json

# Compact output (no indentation)
./mdadf -c input.md > output.json

# Read from stdin
echo "# Hello" | ./mdadf
```

Pre-built binaries are available in the [Releases](https://github.com/ericmason/mdadf/releases) section.

## Notes

- Images are converted to links since ADF media requires Atlassian-specific IDs
- HTML blocks are preserved as plain text
- GFM extensions (tables, strikethrough) are supported via goldmark
- GitHub task list items are converted to Jira-compatible `taskList` / `taskItem` nodes

## Release Process

Releases are created automatically via GitHub Actions when a version tag is pushed.

### Creating a Release

1. **Update version**: Make any necessary code changes
2. **Commit changes**: 
   ```bash
   git add .
   git commit -m "Your commit message"
   git push
   ```
3. **Create and push tag**:
   ```bash
   git tag v0.1.2
   git push origin v0.1.2
   ```

The release workflow will automatically:
- Build binaries for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64, arm64)
- Generate SHA256 checksums for all binaries
- Create a GitHub release with all assets
- Generate release notes from commits

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR.MINOR.PATCH** (e.g., `v1.0.0`)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

Pre-release versions can use tags like `v1.0.0-alpha.1` or `v1.0.0-beta.1`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.
