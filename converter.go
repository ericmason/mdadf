package mdadf

import (
	"encoding/json"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Convert transforms markdown text to ADF JSON bytes
func Convert(markdown string) ([]byte, error) {
	doc, err := ConvertToDoc(markdown)
	if err != nil {
		return nil, err
	}
	return json.Marshal(doc)
}

// ConvertToDoc transforms markdown to an ADF Document struct
func ConvertToDoc(markdown string) (*Document, error) {
	source := []byte(markdown)

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown (tables, strikethrough, etc.)
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	reader := text.NewReader(source)
	mdDoc := md.Parser().Parse(reader)

	adfDoc := NewDocument()
	converter := &adfConverter{source: source}
	adfDoc.Content = converter.convertChildren(mdDoc)

	return adfDoc, nil
}

// adfConverter holds state during conversion
type adfConverter struct {
	source []byte
}

// convertChildren converts all children of an AST node to ADF nodes
func (c *adfConverter) convertChildren(n ast.Node) []Node {
	var nodes []Node
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if converted := c.convertNode(child); converted != nil {
			nodes = append(nodes, *converted)
		}
	}
	return nodes
}

// convertNode converts a single AST node to an ADF node
func (c *adfConverter) convertNode(n ast.Node) *Node {
	switch node := n.(type) {
	case *ast.Document:
		return nil // Document is handled at the top level

	case *ast.Paragraph:
		return c.convertParagraph(node)

	case *ast.Heading:
		return c.convertHeading(node)

	case *ast.List:
		return c.convertList(node)

	case *ast.ListItem:
		return c.convertListItem(node)

	case *ast.FencedCodeBlock:
		return c.convertFencedCodeBlock(node)

	case *ast.CodeBlock:
		return c.convertCodeBlock(node)

	case *ast.Blockquote:
		return c.convertBlockquote(node)

	case *ast.ThematicBreak:
		rule := RuleNode()
		return &rule

	case *ast.Text:
		return c.convertText(node)

	case *ast.String:
		text := TextNode(string(node.Value))
		return &text

	case *ast.CodeSpan:
		return c.convertCodeSpan(node)

	case *ast.Emphasis:
		return c.convertEmphasis(node)

	case *ast.Link:
		return c.convertLink(node)

	case *ast.AutoLink:
		return c.convertAutoLink(node)

	case *ast.Image:
		return c.convertImage(node)

	case *ast.RawHTML:
		return c.convertRawHTML(node)

	case *ast.HTMLBlock:
		return c.convertHTMLBlock(node)

	case *ast.TextBlock:
		return c.convertTextBlock(node)

	case *east.Table:
		return c.convertTable(node)

	case *east.TableRow:
		return c.convertTableRow(node)

	case *east.TableHeader:
		return c.convertTableHeader(node)

	case *east.TableCell:
		return c.convertTableCell(node)

	case *east.Strikethrough:
		return c.convertStrikethrough(node)

	default:
		// For unknown nodes, try to convert children
		children := c.convertChildren(n)
		if len(children) > 0 {
			para := ParagraphNode(children...)
			return &para
		}
		return nil
	}
}

// convertParagraph converts a markdown paragraph to ADF
func (c *adfConverter) convertParagraph(node *ast.Paragraph) *Node {
	content := c.convertInlineChildren(node)
	if len(content) == 0 {
		return nil
	}
	para := ParagraphNode(content...)
	return &para
}

// convertHeading converts a markdown heading to ADF
func (c *adfConverter) convertHeading(node *ast.Heading) *Node {
	content := c.convertInlineChildren(node)
	heading := HeadingNode(node.Level, content...)
	return &heading
}

// convertList converts a markdown list to ADF
func (c *adfConverter) convertList(node *ast.List) *Node {
	items := c.convertChildren(node)
	if node.IsOrdered() {
		list := OrderedListNode(items...)
		return &list
	}
	list := BulletListNode(items...)
	return &list
}

// convertListItem converts a markdown list item to ADF
func (c *adfConverter) convertListItem(node *ast.ListItem) *Node {
	var content []Node
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch child := child.(type) {
		case *ast.TextBlock:
			// TextBlock in list item should become paragraph content
			inlineContent := c.convertInlineChildren(child)
			if len(inlineContent) > 0 {
				content = append(content, ParagraphNode(inlineContent...))
			}
		case *ast.Paragraph:
			if converted := c.convertParagraph(child); converted != nil {
				content = append(content, *converted)
			}
		case *ast.List:
			// Nested list
			if converted := c.convertList(child); converted != nil {
				content = append(content, *converted)
			}
		default:
			if converted := c.convertNode(child); converted != nil {
				content = append(content, *converted)
			}
		}
	}
	item := ListItemNode(content...)
	return &item
}

// convertFencedCodeBlock converts a fenced code block to ADF
func (c *adfConverter) convertFencedCodeBlock(node *ast.FencedCodeBlock) *Node {
	language := string(node.Language(c.source))
	var codeText string
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		codeText += string(line.Value(c.source))
	}
	// Remove trailing newline if present
	if len(codeText) > 0 && codeText[len(codeText)-1] == '\n' {
		codeText = codeText[:len(codeText)-1]
	}
	codeBlock := CodeBlockNode(language, TextNode(codeText))
	return &codeBlock
}

// convertCodeBlock converts an indented code block to ADF
func (c *adfConverter) convertCodeBlock(node *ast.CodeBlock) *Node {
	var codeText string
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		codeText += string(line.Value(c.source))
	}
	if len(codeText) > 0 && codeText[len(codeText)-1] == '\n' {
		codeText = codeText[:len(codeText)-1]
	}
	codeBlock := CodeBlockNode("", TextNode(codeText))
	return &codeBlock
}

// convertBlockquote converts a blockquote to ADF
func (c *adfConverter) convertBlockquote(node *ast.Blockquote) *Node {
	content := c.convertChildren(node)
	blockquote := BlockquoteNode(content...)
	return &blockquote
}

// convertText converts a text node to ADF
func (c *adfConverter) convertText(node *ast.Text) *Node {
	segment := node.Segment
	value := string(segment.Value(c.source))
	if value == "" {
		return nil
	}
	text := TextNode(value)
	return &text
}

// convertTextBlock converts a text block to ADF paragraph
func (c *adfConverter) convertTextBlock(node *ast.TextBlock) *Node {
	content := c.convertInlineChildren(node)
	if len(content) == 0 {
		return nil
	}
	para := ParagraphNode(content...)
	return &para
}

// convertCodeSpan converts inline code to ADF
func (c *adfConverter) convertCodeSpan(node *ast.CodeSpan) *Node {
	var codeText string
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			codeText += string(textNode.Segment.Value(c.source))
		}
	}
	text := TextNode(codeText, CodeMark())
	return &text
}

// convertEmphasis converts emphasis (bold/italic) to ADF
func (c *adfConverter) convertEmphasis(node *ast.Emphasis) *Node {
	var mark Mark
	if node.Level == 2 {
		mark = StrongMark()
	} else {
		mark = EmMark()
	}

	// Get all inline content and apply the mark
	content := c.collectTextWithMark(node, mark)
	if len(content) == 1 {
		return &content[0]
	}

	// If there are multiple nodes, we need to handle them differently
	// For now, return the first one (this handles simple cases)
	if len(content) > 0 {
		return &content[0]
	}
	return nil
}

// collectTextWithMark collects all text content and applies a mark
func (c *adfConverter) collectTextWithMark(n ast.Node, mark Mark) []Node {
	var nodes []Node
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch node := child.(type) {
		case *ast.Text:
			value := string(node.Segment.Value(c.source))
			if value != "" {
				nodes = append(nodes, TextNode(value, mark))
			}
			// Handle line breaks that follow text
			if node.SoftLineBreak() {
				nodes = append(nodes, TextNode(" ", mark))
			}
		case *ast.Emphasis:
			// Nested emphasis - combine marks
			var innerMark Mark
			if node.Level == 2 {
				innerMark = StrongMark()
			} else {
				innerMark = EmMark()
			}
			for innerChild := node.FirstChild(); innerChild != nil; innerChild = innerChild.NextSibling() {
				if textNode, ok := innerChild.(*ast.Text); ok {
					value := string(textNode.Segment.Value(c.source))
					if value != "" {
						nodes = append(nodes, TextNode(value, mark, innerMark))
					}
				}
			}
		case *ast.CodeSpan:
			var codeText string
			for innerChild := node.FirstChild(); innerChild != nil; innerChild = innerChild.NextSibling() {
				if textNode, ok := innerChild.(*ast.Text); ok {
					codeText += string(textNode.Segment.Value(c.source))
				}
			}
			if codeText != "" {
				nodes = append(nodes, TextNode(codeText, mark, CodeMark()))
			}
		case *ast.Link:
			linkText := c.getLinkText(node)
			href := string(node.Destination)
			title := string(node.Title)
			nodes = append(nodes, TextNode(linkText, mark, LinkMark(href, title)))
		default:
			// Recursively handle other nodes
			innerNodes := c.collectTextWithMark(child, mark)
			nodes = append(nodes, innerNodes...)
		}
	}
	return nodes
}

// convertLink converts a link to ADF
func (c *adfConverter) convertLink(node *ast.Link) *Node {
	linkText := c.getLinkText(node)
	href := string(node.Destination)
	title := string(node.Title)
	text := TextNode(linkText, LinkMark(href, title))
	return &text
}

// getLinkText extracts the text content from a link
func (c *adfConverter) getLinkText(node *ast.Link) string {
	var text string
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			text += string(textNode.Segment.Value(c.source))
		}
	}
	return text
}

// convertAutoLink converts an autolink to ADF
func (c *adfConverter) convertAutoLink(node *ast.AutoLink) *Node {
	url := string(node.URL(c.source))
	label := string(node.Label(c.source))
	text := TextNode(label, LinkMark(url, ""))
	return &text
}

// convertImage converts an image to ADF (as a link since we don't have media IDs)
func (c *adfConverter) convertImage(node *ast.Image) *Node {
	// ADF media requires Atlassian-specific IDs, so we convert images to links
	altText := string(node.Text(c.source))
	if altText == "" {
		altText = "image"
	}
	href := string(node.Destination)
	text := TextNode(altText, LinkMark(href, string(node.Title)))
	return &text
}

// convertRawHTML converts raw HTML to ADF (as plain text)
func (c *adfConverter) convertRawHTML(node *ast.RawHTML) *Node {
	// Strip HTML tags and convert to text
	segments := node.Segments
	var htmlContent string
	for i := 0; i < segments.Len(); i++ {
		seg := segments.At(i)
		htmlContent += string(seg.Value(c.source))
	}
	if htmlContent == "" {
		return nil
	}
	text := TextNode(htmlContent)
	return &text
}

// convertHTMLBlock converts an HTML block to ADF (as plain text)
func (c *adfConverter) convertHTMLBlock(node *ast.HTMLBlock) *Node {
	lines := node.Lines()
	var htmlContent string
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		htmlContent += string(line.Value(c.source))
	}
	if htmlContent == "" {
		return nil
	}
	para := ParagraphNode(TextNode(htmlContent))
	return &para
}

// convertTable converts a GFM table to ADF
func (c *adfConverter) convertTable(node *east.Table) *Node {
	rows := c.convertChildren(node)
	table := TableNode(rows...)
	return &table
}

// convertTableRow converts a table row to ADF
func (c *adfConverter) convertTableRow(node *east.TableRow) *Node {
	cells := c.convertChildren(node)
	row := TableRowNode(cells...)
	return &row
}

// convertTableHeader converts a table header to ADF
func (c *adfConverter) convertTableHeader(node *east.TableHeader) *Node {
	cells := c.convertChildren(node)
	row := TableRowNode(cells...)
	return &row
}

// convertTableCell converts a table cell to ADF
func (c *adfConverter) convertTableCell(node *east.TableCell) *Node {
	content := c.convertInlineChildren(node)
	para := ParagraphNode(content...)

	// Check if this is a header cell (parent is TableHeader)
	if _, isHeader := node.Parent().(*east.TableHeader); isHeader {
		cell := TableHeaderNode(para)
		return &cell
	}
	cell := TableCellNode(para)
	return &cell
}

// convertStrikethrough converts strikethrough text to ADF
func (c *adfConverter) convertStrikethrough(node *east.Strikethrough) *Node {
	content := c.collectTextWithMark(node, StrikeMark())
	if len(content) == 1 {
		return &content[0]
	}
	if len(content) > 0 {
		return &content[0]
	}
	return nil
}

// convertInlineChildren converts inline children (text, emphasis, links, etc.)
func (c *adfConverter) convertInlineChildren(n ast.Node) []Node {
	var nodes []Node
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch node := child.(type) {
		case *ast.Text:
			value := string(node.Segment.Value(c.source))
			if value != "" {
				nodes = append(nodes, TextNode(value))
			}
			// Handle line breaks that follow text
			if node.HardLineBreak() {
				nodes = append(nodes, HardBreakNode())
			} else if node.SoftLineBreak() {
				nodes = append(nodes, TextNode(" "))
			}
		case *ast.Emphasis:
			emphNodes := c.convertEmphasisInline(node)
			nodes = append(nodes, emphNodes...)
		case *ast.CodeSpan:
			if converted := c.convertCodeSpan(node); converted != nil {
				nodes = append(nodes, *converted)
			}
		case *ast.Link:
			if converted := c.convertLink(node); converted != nil {
				nodes = append(nodes, *converted)
			}
		case *ast.AutoLink:
			if converted := c.convertAutoLink(node); converted != nil {
				nodes = append(nodes, *converted)
			}
		case *ast.Image:
			if converted := c.convertImage(node); converted != nil {
				nodes = append(nodes, *converted)
			}
		case *ast.RawHTML:
			if converted := c.convertRawHTML(node); converted != nil {
				nodes = append(nodes, *converted)
			}
		case *east.Strikethrough:
			strikeNodes := c.convertStrikethroughInline(node)
			nodes = append(nodes, strikeNodes...)
		default:
			// Try to get any text content
			if converted := c.convertNode(child); converted != nil {
				nodes = append(nodes, *converted)
			}
		}
	}
	return nodes
}

// convertEmphasisInline handles emphasis and returns multiple nodes if needed
func (c *adfConverter) convertEmphasisInline(node *ast.Emphasis) []Node {
	var mark Mark
	if node.Level == 2 {
		mark = StrongMark()
	} else {
		mark = EmMark()
	}
	return c.collectTextWithMark(node, mark)
}

// convertStrikethroughInline handles strikethrough and returns multiple nodes if needed
func (c *adfConverter) convertStrikethroughInline(node *east.Strikethrough) []Node {
	return c.collectTextWithMark(node, StrikeMark())
}
