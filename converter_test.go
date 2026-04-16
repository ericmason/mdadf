package mdadf

import (
	"encoding/json"
	"testing"
)

func TestConvert_Paragraph(t *testing.T) {
	md := "Hello world"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc.Version != 1 {
		t.Errorf("expected version 1, got %d", doc.Version)
	}
	if doc.Type != "doc" {
		t.Errorf("expected type 'doc', got %s", doc.Type)
	}
	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 content node, got %d", len(doc.Content))
	}

	para := doc.Content[0]
	if para.Type != "paragraph" {
		t.Errorf("expected paragraph, got %s", para.Type)
	}
	if len(para.Content) < 1 {
		t.Fatalf("expected at least 1 text node, got %d", len(para.Content))
	}
	// Combine all text content
	var text string
	for _, node := range para.Content {
		text += node.Text
	}
	if text != "Hello world" {
		t.Errorf("expected 'Hello world', got %s", text)
	}
}

func TestConvert_Headings(t *testing.T) {
	tests := []struct {
		md    string
		level int
		text  string
	}{
		{"# Heading 1", 1, "Heading 1"},
		{"## Heading 2", 2, "Heading 2"},
		{"### Heading 3", 3, "Heading 3"},
		{"#### Heading 4", 4, "Heading 4"},
		{"##### Heading 5", 5, "Heading 5"},
		{"###### Heading 6", 6, "Heading 6"},
	}

	for _, tt := range tests {
		t.Run(tt.md, func(t *testing.T) {
			doc, err := ConvertToDoc(tt.md)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(doc.Content) != 1 {
				t.Fatalf("expected 1 content node, got %d", len(doc.Content))
			}

			heading := doc.Content[0]
			if heading.Type != "heading" {
				t.Errorf("expected heading, got %s", heading.Type)
			}

			var attrs HeadingAttrs
			if err := json.Unmarshal(heading.Attrs, &attrs); err != nil {
				t.Fatalf("failed to unmarshal attrs: %v", err)
			}
			if attrs.Level != tt.level {
				t.Errorf("expected level %d, got %d", tt.level, attrs.Level)
			}

			if len(heading.Content) < 1 {
				t.Fatalf("expected at least 1 text node, got %d", len(heading.Content))
			}
			// Combine all text content
			var text string
			for _, node := range heading.Content {
				text += node.Text
			}
			if text != tt.text {
				t.Errorf("expected '%s', got '%s'", tt.text, text)
			}
		})
	}
}

func TestConvert_Bold(t *testing.T) {
	md := "This is **bold** text"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]
	if len(para.Content) != 3 {
		t.Fatalf("expected 3 text nodes, got %d", len(para.Content))
	}

	// Check "This is "
	if para.Content[0].Text != "This is " {
		t.Errorf("expected 'This is ', got '%s'", para.Content[0].Text)
	}

	// Check "bold" with strong mark
	boldNode := para.Content[1]
	if boldNode.Text != "bold" {
		t.Errorf("expected 'bold', got '%s'", boldNode.Text)
	}
	if len(boldNode.Marks) != 1 || boldNode.Marks[0].Type != "strong" {
		t.Errorf("expected strong mark, got %v", boldNode.Marks)
	}

	// Check " text"
	if para.Content[2].Text != " text" {
		t.Errorf("expected ' text', got '%s'", para.Content[2].Text)
	}
}

func TestConvert_Italic(t *testing.T) {
	md := "This is *italic* text"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]
	if len(para.Content) != 3 {
		t.Fatalf("expected 3 text nodes, got %d", len(para.Content))
	}

	italicNode := para.Content[1]
	if italicNode.Text != "italic" {
		t.Errorf("expected 'italic', got '%s'", italicNode.Text)
	}
	if len(italicNode.Marks) != 1 || italicNode.Marks[0].Type != "em" {
		t.Errorf("expected em mark, got %v", italicNode.Marks)
	}
}

func TestConvert_InlineCode(t *testing.T) {
	md := "Use `code` here"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]
	if len(para.Content) != 3 {
		t.Fatalf("expected 3 text nodes, got %d", len(para.Content))
	}

	codeNode := para.Content[1]
	if codeNode.Text != "code" {
		t.Errorf("expected 'code', got '%s'", codeNode.Text)
	}
	if len(codeNode.Marks) != 1 || codeNode.Marks[0].Type != "code" {
		t.Errorf("expected code mark, got %v", codeNode.Marks)
	}
}

func TestConvert_Link(t *testing.T) {
	md := "Check [this link](https://example.com)"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]
	if len(para.Content) != 2 {
		t.Fatalf("expected 2 text nodes, got %d", len(para.Content))
	}

	linkNode := para.Content[1]
	if linkNode.Text != "this link" {
		t.Errorf("expected 'this link', got '%s'", linkNode.Text)
	}
	if len(linkNode.Marks) != 1 || linkNode.Marks[0].Type != "link" {
		t.Errorf("expected link mark, got %v", linkNode.Marks)
	}

	var attrs LinkAttrs
	if err := json.Unmarshal(linkNode.Marks[0].Attrs, &attrs); err != nil {
		t.Fatalf("failed to unmarshal attrs: %v", err)
	}
	if attrs.Href != "https://example.com" {
		t.Errorf("expected 'https://example.com', got '%s'", attrs.Href)
	}
}

func TestConvert_BulletList(t *testing.T) {
	md := `- Item 1
- Item 2
- Item 3`
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 content node, got %d", len(doc.Content))
	}

	list := doc.Content[0]
	if list.Type != "bulletList" {
		t.Errorf("expected bulletList, got %s", list.Type)
	}
	if len(list.Content) != 3 {
		t.Fatalf("expected 3 list items, got %d", len(list.Content))
	}

	for i, item := range list.Content {
		if item.Type != "listItem" {
			t.Errorf("item %d: expected listItem, got %s", i, item.Type)
		}
	}
}

func TestConvert_OrderedList(t *testing.T) {
	md := `1. First
2. Second
3. Third`
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := doc.Content[0]
	if list.Type != "orderedList" {
		t.Errorf("expected orderedList, got %s", list.Type)
	}
	if len(list.Content) != 3 {
		t.Fatalf("expected 3 list items, got %d", len(list.Content))
	}
}

func TestConvert_NestedList(t *testing.T) {
	md := `- Item 1
  - Nested 1
  - Nested 2
- Item 2`
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := doc.Content[0]
	if list.Type != "bulletList" {
		t.Errorf("expected bulletList, got %s", list.Type)
	}

	// First item should have nested list
	firstItem := list.Content[0]
	if firstItem.Type != "listItem" {
		t.Errorf("expected listItem, got %s", firstItem.Type)
	}

	// Check for nested list in first item's content
	hasNestedList := false
	for _, child := range firstItem.Content {
		if child.Type == "bulletList" {
			hasNestedList = true
			if len(child.Content) != 2 {
				t.Errorf("expected 2 nested items, got %d", len(child.Content))
			}
		}
	}
	if !hasNestedList {
		t.Error("expected nested list in first item")
	}
}

func TestConvert_TaskList(t *testing.T) {
	md := `- [ ] First task
- [x] Second task`

	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 content node, got %d", len(doc.Content))
	}

	list := doc.Content[0]
	if list.Type != "taskList" {
		t.Fatalf("expected taskList, got %s", list.Type)
	}

	if len(list.Content) != 2 {
		t.Fatalf("expected 2 task items, got %d", len(list.Content))
	}

	tests := []struct {
		index         int
		expectedState string
		expectedText  string
	}{
		{0, "TODO", "First task"},
		{1, "DONE", "Second task"},
	}

	for _, tt := range tests {
		item := list.Content[tt.index]
		if item.Type != "taskItem" {
			t.Fatalf("item %d: expected taskItem, got %s", tt.index, item.Type)
		}

		var attrs TaskItemAttrs
		if err := json.Unmarshal(item.Attrs, &attrs); err != nil {
			t.Fatalf("item %d: failed to unmarshal attrs: %v", tt.index, err)
		}
		if attrs.LocalID == "" {
			t.Fatalf("item %d: expected localId to be set", tt.index)
		}
		if attrs.State != tt.expectedState {
			t.Fatalf("item %d: expected state %s, got %s", tt.index, tt.expectedState, attrs.State)
		}

		var text string
		for _, node := range item.Content {
			text += node.Text
		}
		if text != tt.expectedText {
			t.Fatalf("item %d: expected text %q, got %q", tt.index, tt.expectedText, text)
		}
	}
}

func TestConvert_NestedTaskList(t *testing.T) {
	md := `- [x] Parent
  - [ ] Child
- [ ] Sibling`

	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := doc.Content[0]
	if list.Type != "taskList" {
		t.Fatalf("expected taskList, got %s", list.Type)
	}
	var listAttrs TaskListAttrs
	if err := json.Unmarshal(list.Attrs, &listAttrs); err != nil {
		t.Fatalf("failed to unmarshal taskList attrs: %v", err)
	}
	if listAttrs.LocalID == "" {
		t.Fatal("expected taskList localId to be set")
	}

	parent := list.Content[0]
	if parent.Type != "taskItem" {
		t.Fatalf("expected parent taskItem, got %s", parent.Type)
	}

	if len(parent.Content) < 2 {
		t.Fatalf("expected parent to contain text and nested taskList, got %d nodes", len(parent.Content))
	}

	nested := parent.Content[len(parent.Content)-1]
	if nested.Type != "taskList" {
		t.Fatalf("expected nested taskList, got %s", nested.Type)
	}

	if len(nested.Content) != 1 || nested.Content[0].Type != "taskItem" {
		t.Fatalf("expected one nested taskItem, got %+v", nested.Content)
	}

	var attrs TaskItemAttrs
	if err := json.Unmarshal(nested.Content[0].Attrs, &attrs); err != nil {
		t.Fatalf("failed to unmarshal nested task attrs: %v", err)
	}
	if attrs.LocalID == "" {
		t.Fatal("expected nested task localId to be set")
	}
	if attrs.State != "TODO" {
		t.Fatalf("expected nested task state TODO, got %s", attrs.State)
	}
}

func TestConvert_CodeBlock(t *testing.T) {
	md := "```javascript\nconst x = 1;\n```"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	codeBlock := doc.Content[0]
	if codeBlock.Type != "codeBlock" {
		t.Errorf("expected codeBlock, got %s", codeBlock.Type)
	}

	var attrs CodeBlockAttrs
	if err := json.Unmarshal(codeBlock.Attrs, &attrs); err != nil {
		t.Fatalf("failed to unmarshal attrs: %v", err)
	}
	if attrs.Language != "javascript" {
		t.Errorf("expected language 'javascript', got '%s'", attrs.Language)
	}

	if len(codeBlock.Content) != 1 {
		t.Fatalf("expected 1 text node, got %d", len(codeBlock.Content))
	}
	if codeBlock.Content[0].Text != "const x = 1;" {
		t.Errorf("expected 'const x = 1;', got '%s'", codeBlock.Content[0].Text)
	}
}

func TestConvert_Blockquote(t *testing.T) {
	md := "> This is a quote"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	blockquote := doc.Content[0]
	if blockquote.Type != "blockquote" {
		t.Errorf("expected blockquote, got %s", blockquote.Type)
	}

	if len(blockquote.Content) != 1 {
		t.Fatalf("expected 1 paragraph, got %d", len(blockquote.Content))
	}
	if blockquote.Content[0].Type != "paragraph" {
		t.Errorf("expected paragraph in blockquote, got %s", blockquote.Content[0].Type)
	}
}

func TestConvert_HorizontalRule(t *testing.T) {
	md := "Above\n\n---\n\nBelow"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Content) != 3 {
		t.Fatalf("expected 3 content nodes, got %d", len(doc.Content))
	}

	rule := doc.Content[1]
	if rule.Type != "rule" {
		t.Errorf("expected rule, got %s", rule.Type)
	}
}

func TestConvert_Table(t *testing.T) {
	md := `| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |
| Cell 3   | Cell 4   |`

	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	table := doc.Content[0]
	if table.Type != "table" {
		t.Errorf("expected table, got %s", table.Type)
	}

	// Should have 3 rows (1 header + 2 data)
	if len(table.Content) != 3 {
		t.Fatalf("expected 3 table rows, got %d", len(table.Content))
	}

	// First row should be header row
	headerRow := table.Content[0]
	if headerRow.Type != "tableRow" {
		t.Errorf("expected tableRow, got %s", headerRow.Type)
	}

	// Check that header cells are tableHeader type
	for _, cell := range headerRow.Content {
		if cell.Type != "tableHeader" {
			t.Errorf("expected tableHeader, got %s", cell.Type)
		}
	}

	// Data rows should have tableCell
	for i := 1; i < len(table.Content); i++ {
		row := table.Content[i]
		for _, cell := range row.Content {
			if cell.Type != "tableCell" {
				t.Errorf("expected tableCell, got %s", cell.Type)
			}
		}
	}
}

func TestConvert_Strikethrough(t *testing.T) {
	md := "This is ~~deleted~~ text"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]
	if len(para.Content) < 3 {
		t.Fatalf("expected at least 3 text nodes, got %d", len(para.Content))
	}

	// Find the strikethrough node
	found := false
	for _, node := range para.Content {
		if node.Text == "deleted" {
			found = true
			if len(node.Marks) != 1 || node.Marks[0].Type != "strike" {
				t.Errorf("expected strike mark, got %v", node.Marks)
			}
		}
	}
	if !found {
		t.Error("expected to find 'deleted' text with strike mark")
	}
}

func TestConvert_BoldAndItalic(t *testing.T) {
	md := "This is ***bold and italic*** text"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]

	// Find the bold+italic node
	found := false
	for _, node := range para.Content {
		if node.Text == "bold and italic" {
			found = true
			if len(node.Marks) != 2 {
				t.Errorf("expected 2 marks, got %d", len(node.Marks))
			}
			// Should have both strong and em marks
			hasStrong := false
			hasEm := false
			for _, mark := range node.Marks {
				if mark.Type == "strong" {
					hasStrong = true
				}
				if mark.Type == "em" {
					hasEm = true
				}
			}
			if !hasStrong || !hasEm {
				t.Errorf("expected both strong and em marks, got %v", node.Marks)
			}
		}
	}
	if !found {
		t.Error("expected to find 'bold and italic' text")
	}
}

func TestConvert_MultipleMarks(t *testing.T) {
	md := "Check **[bold link](https://example.com)**"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]

	// Find the node with both marks
	found := false
	for _, node := range para.Content {
		if node.Text == "bold link" {
			found = true
			if len(node.Marks) != 2 {
				t.Errorf("expected 2 marks, got %d: %v", len(node.Marks), node.Marks)
			}
		}
	}
	if !found {
		t.Error("expected to find 'bold link' text")
	}
}

func TestConvert_EmptyDocument(t *testing.T) {
	md := ""
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if doc.Version != 1 {
		t.Errorf("expected version 1, got %d", doc.Version)
	}
	if doc.Type != "doc" {
		t.Errorf("expected type 'doc', got %s", doc.Type)
	}
	if len(doc.Content) != 0 {
		t.Errorf("expected 0 content nodes, got %d", len(doc.Content))
	}
}

func TestConvert_JSON(t *testing.T) {
	md := "# Hello\n\nWorld"
	jsonBytes, err := Convert(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var doc Document
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if doc.Version != 1 {
		t.Errorf("expected version 1, got %d", doc.Version)
	}
	if len(doc.Content) != 2 {
		t.Fatalf("expected 2 content nodes, got %d", len(doc.Content))
	}
}

func TestConvert_Image(t *testing.T) {
	md := "![Alt text](https://example.com/image.png)"
	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	para := doc.Content[0]
	if len(para.Content) != 1 {
		t.Fatalf("expected 1 node, got %d", len(para.Content))
	}

	// Image should be converted to a link (since we don't have media IDs)
	imgNode := para.Content[0]
	if imgNode.Text != "Alt text" {
		t.Errorf("expected 'Alt text', got '%s'", imgNode.Text)
	}
	if len(imgNode.Marks) != 1 || imgNode.Marks[0].Type != "link" {
		t.Errorf("expected link mark, got %v", imgNode.Marks)
	}

	var attrs LinkAttrs
	if err := json.Unmarshal(imgNode.Marks[0].Attrs, &attrs); err != nil {
		t.Fatalf("failed to unmarshal attrs: %v", err)
	}
	if attrs.Href != "https://example.com/image.png" {
		t.Errorf("expected 'https://example.com/image.png', got '%s'", attrs.Href)
	}
}

func TestConvert_ComplexDocument(t *testing.T) {
	md := `# Main Title

This is a paragraph with **bold**, *italic*, and ` + "`code`" + `.

## Sub Heading

- List item 1
- List item 2
  - Nested item

` + "```go\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```" + `

> A quote

---

| Col 1 | Col 2 |
|-------|-------|
| A     | B     |
`

	doc, err := ConvertToDoc(md)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify we got a valid document
	if doc.Version != 1 {
		t.Errorf("expected version 1, got %d", doc.Version)
	}
	if len(doc.Content) < 5 {
		t.Errorf("expected at least 5 content nodes, got %d", len(doc.Content))
	}

	// Verify JSON output is valid
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}
	if len(jsonBytes) == 0 {
		t.Error("expected non-empty JSON output")
	}
}
