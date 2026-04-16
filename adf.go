package mdadf

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Document represents the root ADF document
type Document struct {
	Version int    `json:"version"`
	Type    string `json:"type"`
	Content []Node `json:"content"`
}

// NewDocument creates a new ADF document with version 1
func NewDocument() *Document {
	return &Document{
		Version: 1,
		Type:    "doc",
		Content: []Node{},
	}
}

// Node represents any ADF node (block or inline)
type Node struct {
	Type    string          `json:"type"`
	Text    string          `json:"text,omitempty"`
	Content []Node          `json:"content,omitempty"`
	Marks   []Mark          `json:"marks,omitempty"`
	Attrs   json.RawMessage `json:"attrs,omitempty"`
}

// Mark represents text formatting (bold, italic, code, link, etc.)
type Mark struct {
	Type  string          `json:"type"`
	Attrs json.RawMessage `json:"attrs,omitempty"`
}

// HeadingAttrs holds attributes for heading nodes
type HeadingAttrs struct {
	Level int `json:"level"`
}

// CodeBlockAttrs holds attributes for code block nodes
type CodeBlockAttrs struct {
	Language string `json:"language,omitempty"`
}

// LinkAttrs holds attributes for link marks
type LinkAttrs struct {
	Href  string `json:"href"`
	Title string `json:"title,omitempty"`
}

// TableAttrs holds attributes for table nodes
type TableAttrs struct {
	IsNumberColumnEnabled bool   `json:"isNumberColumnEnabled,omitempty"`
	Layout                string `json:"layout,omitempty"`
}

// TableCellAttrs holds attributes for table cell/header nodes
type TableCellAttrs struct {
	ColSpan int `json:"colspan,omitempty"`
	RowSpan int `json:"rowspan,omitempty"`
}

// MediaAttrs holds attributes for media nodes
type MediaAttrs struct {
	Type       string `json:"type"`
	ID         string `json:"id,omitempty"`
	Collection string `json:"collection,omitempty"`
	URL        string `json:"url,omitempty"`
	Alt        string `json:"alt,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
}

// MediaSingleAttrs holds attributes for mediaSingle nodes
type MediaSingleAttrs struct {
	Layout string `json:"layout,omitempty"`
}

// OrderedListAttrs holds attributes for ordered list nodes
type OrderedListAttrs struct {
	Order int `json:"order,omitempty"`
}

// TaskListAttrs holds attributes for task list nodes
type TaskListAttrs struct {
	LocalID string `json:"localId"`
}

// TaskItemAttrs holds attributes for task item nodes
type TaskItemAttrs struct {
	LocalID string `json:"localId"`
	State   string `json:"state"`
}

// Helper functions to create common nodes

// TextNode creates a text node with optional marks
func TextNode(text string, marks ...Mark) Node {
	node := Node{
		Type: "text",
		Text: text,
	}
	if len(marks) > 0 {
		node.Marks = marks
	}
	return node
}

// ParagraphNode creates a paragraph node with content
func ParagraphNode(content ...Node) Node {
	return Node{
		Type:    "paragraph",
		Content: content,
	}
}

// HeadingNode creates a heading node with the specified level
func HeadingNode(level int, content ...Node) Node {
	attrs, _ := json.Marshal(HeadingAttrs{Level: level})
	return Node{
		Type:    "heading",
		Attrs:   attrs,
		Content: content,
	}
}

// BulletListNode creates a bullet list node
func BulletListNode(items ...Node) Node {
	return Node{
		Type:    "bulletList",
		Content: items,
	}
}

// OrderedListNode creates an ordered list node
func OrderedListNode(items ...Node) Node {
	return Node{
		Type:    "orderedList",
		Content: items,
	}
}

// ListItemNode creates a list item node
func ListItemNode(content ...Node) Node {
	return Node{
		Type:    "listItem",
		Content: content,
	}
}

// TaskListNode creates a task list node
func TaskListNode(items ...Node) Node {
	attrs, _ := json.Marshal(TaskListAttrs{LocalID: uuid.NewString()})
	return Node{
		Type:    "taskList",
		Attrs:   attrs,
		Content: items,
	}
}

// TaskItemNode creates a task item node with a Jira-compatible state
func TaskItemNode(checked bool, content ...Node) Node {
	state := "TODO"
	if checked {
		state = "DONE"
	}
	attrs, _ := json.Marshal(TaskItemAttrs{
		LocalID: uuid.NewString(),
		State:   state,
	})
	return Node{
		Type:    "taskItem",
		Attrs:   attrs,
		Content: content,
	}
}

// CodeBlockNode creates a code block node
func CodeBlockNode(language string, content ...Node) Node {
	var attrs json.RawMessage
	if language != "" {
		attrs, _ = json.Marshal(CodeBlockAttrs{Language: language})
	}
	return Node{
		Type:    "codeBlock",
		Attrs:   attrs,
		Content: content,
	}
}

// BlockquoteNode creates a blockquote node
func BlockquoteNode(content ...Node) Node {
	return Node{
		Type:    "blockquote",
		Content: content,
	}
}

// RuleNode creates a horizontal rule node
func RuleNode() Node {
	return Node{
		Type: "rule",
	}
}

// HardBreakNode creates a hard break (line break) node
func HardBreakNode() Node {
	return Node{
		Type: "hardBreak",
	}
}

// TableNode creates a table node
func TableNode(rows ...Node) Node {
	return Node{
		Type:    "table",
		Content: rows,
	}
}

// TableRowNode creates a table row node
func TableRowNode(cells ...Node) Node {
	return Node{
		Type:    "tableRow",
		Content: cells,
	}
}

// TableHeaderNode creates a table header cell node
func TableHeaderNode(content ...Node) Node {
	return Node{
		Type:    "tableHeader",
		Content: content,
	}
}

// TableCellNode creates a table cell node
func TableCellNode(content ...Node) Node {
	return Node{
		Type:    "tableCell",
		Content: content,
	}
}

// Mark helper functions

// StrongMark creates a bold/strong mark
func StrongMark() Mark {
	return Mark{Type: "strong"}
}

// EmMark creates an italic/emphasis mark
func EmMark() Mark {
	return Mark{Type: "em"}
}

// CodeMark creates an inline code mark
func CodeMark() Mark {
	return Mark{Type: "code"}
}

// StrikeMark creates a strikethrough mark
func StrikeMark() Mark {
	return Mark{Type: "strike"}
}

// UnderlineMark creates an underline mark
func UnderlineMark() Mark {
	return Mark{Type: "underline"}
}

// LinkMark creates a link mark with href and optional title
func LinkMark(href, title string) Mark {
	attrs := LinkAttrs{Href: href}
	if title != "" {
		attrs.Title = title
	}
	attrsJSON, _ := json.Marshal(attrs)
	return Mark{
		Type:  "link",
		Attrs: attrsJSON,
	}
}
