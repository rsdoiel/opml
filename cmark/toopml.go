/*
opml is a Go package for working with OPML XML files.
Copyright (C) 2021 R. S. Doiel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package cmark

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rsdoiel/opml"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

/** ErrNotOutline is wrapped by the error ToOPML returns for Markdown that is
 * not an outline list: a heading, paragraph, code block, quote or rule
 * outside the list, a second paragraph or a code block inside an item, a
 * reference-style link, text after a link, a repeated or unknown attribute,
 * an unknown front matter key. The message names the line.
 *
 * Example:
 *   _, err := cmark.ToOPML(src)
 *   if errors.Is(err, cmark.ErrNotOutline) {
 *       fmt.Println(err) // "line 7: a heading is not part of an outline: ..."
 *   }
 */
var ErrNotOutline = errors.New("not an outline list")

/** ToOPML reads a Markdown outline, the form FromOPML writes and the forms a
 * person editing it by hand may leave: front matter for the head, then lists.
 * It is lenient about the list (any bullet or number, any indent, blank
 * lines between items, an item wrapped over lines, sibling lists, CRLF) and
 * strict about everything else, so no text is dropped silently (opml
 * DR-0001, decision 7). The Markdown is read with the goldmark AST.
 *
 * Parameters:
 *   src ([]byte) — the Markdown
 *
 * Returns:
 *   *opml.OPML — the document; an empty one for empty input
 *   error — nil, or an error wrapping ErrNotOutline that names the line
 *
 * Example:
 *   src, _ := os.ReadFile("feeds.md")
 *   doc, err := cmark.ToOPML(src)
 */
func ToOPML(src []byte) (*opml.OPML, error) {
	o := opml.New()
	body, base, err := readFrontMatter(src, o)
	if err != nil {
		return nil, err
	}
	r := &reader{body: body, base: base}
	root := goldmark.New().Parser().Parse(text.NewReader(body))
	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		list, ok := n.(*ast.List)
		if !ok {
			return nil, r.foreign(n)
		}
		outlines, err := r.list(list)
		if err != nil {
			return nil, err
		}
		o.Body.Outline = append(o.Body.Outline, outlines...)
	}
	return o, nil
}

// readFrontMatter splits off a leading --- block, applies it to the head and
// returns the rest and the number of lines that came before it.
func readFrontMatter(src []byte, o *opml.OPML) (rest []byte, lines int, err error) {
	first, after, _ := bytes.Cut(src, []byte("\n"))
	if string(bytes.TrimRight(first, " \t\r")) != "---" {
		return src, 0, nil
	}
	pos, n := len(first)+1, 1
	for pos <= len(src) {
		line, _, _ := bytes.Cut(src[pos:], []byte("\n"))
		n++
		pos += len(line) + 1
		if s := string(bytes.TrimRight(line, " \t\r")); s == "---" || s == "..." {
			fm := after[:max(0, pos-len(line)-1-len(first)-1)]
			if pos > len(src) {
				pos = len(src)
			}
			if err := applyFrontMatter(fm, o, n); err != nil {
				return nil, 0, err
			}
			return src[pos:], n, nil
		}
	}
	return src, 0, nil // never closed: not front matter
}

func applyFrontMatter(fm []byte, o *opml.OPML, endLine int) error {
	var m frontMatter
	dec := yaml.NewDecoder(bytes.NewReader(fm))
	dec.KnownFields(true)
	if err := dec.Decode(&m); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("front matter (lines 1 to %d): %v: %w", endLine, strings.TrimPrefix(err.Error(), "yaml: "), ErrNotOutline)
	}
	if m.OPMLVersion != nil {
		o.Version = *m.OPMLVersion
	}
	h := o.Head
	h.Title, h.Created, h.Modified, h.Generator = m.Title, m.DateCreated, m.DateModified, m.Generator
	h.OwnerName, h.OwnerEmail, h.OwnerID, h.Docs, h.ExpansionState = m.OwnerName, m.OwnerEmail, m.OwnerID, m.Docs, m.ExpansionState
	h.VertScrollState, h.WindowTop, h.WindowLeft = m.VertScrollState, m.WindowTop, m.WindowLeft
	h.WindowBottom, h.WindowRight = m.WindowBottom, m.WindowRight
	return nil
}

// reader holds the Markdown body and the number of file lines before it, so
// an error names the line in the file.
type reader struct {
	body []byte
	base int
}

func (r *reader) lineAt(offset int) int {
	offset = min(max(offset, 0), len(r.body))
	return r.base + 1 + bytes.Count(r.body[:offset], []byte("\n"))
}

func (r *reader) fail(line int, format string, a ...any) error {
	return fmt.Errorf("line %d: %s: %w", line, fmt.Sprintf(format, a...), ErrNotOutline)
}

// position finds where a block starts: the first line it holds, or for a
// fenced code block the fence before its first line. ok is false for a block
// with no line of its own (a thematic break).
func (r *reader) position(n ast.Node) (offset int, ok bool) {
	if n.Type() == ast.TypeBlock {
		if lines := n.Lines(); lines.Len() > 0 {
			start := lines.At(0).Start
			if _, fenced := n.(*ast.FencedCodeBlock); fenced {
				// the opening fence is the line before the first content line
				contentLine := bytes.LastIndexByte(r.body[:start], '\n') + 1
				return bytes.LastIndexByte(r.body[:max(contentLine-1, 0)], '\n') + 1, true
			}
			return start, true
		}
	}
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, isText := c.(*ast.Text); isText {
			return t.Segment.Start, true
		}
		if off, ok := r.position(c); ok {
			return off, true
		}
	}
	return 0, false
}

// end is the offset after the last line a node holds, or -1.
func (r *reader) end(n ast.Node) int {
	best := -1
	if n.Type() == ast.TypeBlock {
		if lines := n.Lines(); lines.Len() > 0 {
			best = lines.At(lines.Len() - 1).Stop
		}
	}
	for c := n.LastChild(); c != nil; c = c.PreviousSibling() {
		if e := r.end(c); e > best {
			best = e
		}
		if t, ok := c.(*ast.Text); ok && t.Segment.Stop > best {
			best = t.Segment.Stop
		}
	}
	return best
}

// lineOf is the best line number for a block: its own, or for a block with
// no line the first non-blank line after the previous block.
func (r *reader) lineOf(n ast.Node) int {
	if off, ok := r.position(n); ok {
		return r.lineAt(off)
	}
	from := 0
	for p := n.PreviousSibling(); p != nil; p = p.PreviousSibling() {
		if e := r.end(p); e >= 0 {
			from = e
			break
		}
	}
	for from < len(r.body) {
		i := bytes.IndexByte(r.body[from:], '\n')
		line := r.body[from:]
		if i >= 0 {
			line = r.body[from : from+i]
		}
		if len(bytes.TrimSpace(line)) > 0 || i < 0 {
			break
		}
		from += i + 1
	}
	return r.lineAt(from)
}

func (r *reader) foreign(n ast.Node) error {
	what := strings.ToLower(n.Kind().String())
	switch n.(type) {
	case *ast.Heading:
		what = "a heading"
	case *ast.Paragraph, *ast.TextBlock:
		what = "a paragraph"
	case *ast.FencedCodeBlock, *ast.CodeBlock:
		what = "a code block"
	case *ast.Blockquote:
		what = "a block quote"
	case *ast.ThematicBreak:
		what = "a thematic break"
	case *ast.HTMLBlock:
		what = "an HTML block"
	}
	return r.fail(r.lineOf(n), "%s is not part of an outline (only a list is)", what)
}

func (r *reader) list(l *ast.List) ([]*opml.Outline, error) {
	var out []*opml.Outline
	for it := l.FirstChild(); it != nil; it = it.NextSibling() {
		item, ok := it.(*ast.ListItem)
		if !ok {
			return nil, r.foreign(it)
		}
		ol, err := r.item(item)
		if err != nil {
			return nil, err
		}
		out = append(out, ol)
	}
	return out, nil
}

func (r *reader) item(li *ast.ListItem) (*opml.Outline, error) {
	ol := &opml.Outline{}
	c := li.FirstChild()
	switch c.(type) {
	case *ast.Paragraph, *ast.TextBlock:
		if err := r.content(c, ol); err != nil {
			return nil, err
		}
		c = c.NextSibling()
	}
	for ; c != nil; c = c.NextSibling() {
		sub, ok := c.(*ast.List)
		if !ok {
			return nil, r.foreign(c)
		}
		children, err := r.list(sub)
		if err != nil {
			return nil, err
		}
		ol.Outline = append(ol.Outline, children...)
	}
	return ol, nil
}

// content reads the first paragraph of an item: text or a link, then an
// optional attribute block.
func (r *reader) content(block ast.Node, ol *opml.Outline) error {
	lines := block.Lines()
	parts := make([]string, 0, lines.Len())
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		parts = append(parts, strings.TrimSpace(string(seg.Value(r.body))))
	}
	raw := strings.Join(parts, " ")
	line := r.lineOf(block)

	var dest string
	var tail string
	if link, ok := block.FirstChild().(*ast.Link); ok && strings.HasPrefix(raw, "[") {
		label, after, err := scanLink(raw)
		if err != nil {
			return r.fail(line, "%v", err)
		}
		ol.Text = unescape(label)
		// goldmark keeps the destination as written, minus any angle
		// brackets; resolving backslash escapes is left to the renderer.
		dest, tail = unescape(string(link.Destination)), strings.TrimSpace(after)
		if tail != "" {
			if !strings.HasPrefix(tail, "{") {
				return r.fail(line, "text after a link (%q): put attributes in a {key=\"value\"} block", tail)
			}
			if _, ok := parseBlock(tail); !ok {
				return r.fail(line, "text after a link (%q) is not a {key=\"value\"} block", tail)
			}
		}
	} else {
		text, block := splitBlock(raw)
		ol.Text, tail = unescape(text), block
	}
	if tail != "" {
		pairs, _ := parseBlock(tail)
		if err := assign(ol, pairs); err != nil {
			return r.fail(line, "%v", err)
		}
	}
	if dest != "" {
		attr, field := "url", &ol.URL
		if isFeedType(ol.Type) {
			attr, field = "xmlUrl", &ol.XMLURL
		}
		if *field != "" {
			return r.fail(line, "the link is the %s, which the attribute block also has", attr)
		}
		*field = dest
	}
	return nil
}

// scanLink reads "[label](destination)" at the start of s and returns the
// raw label and what follows the closing parenthesis. The destination itself
// comes from the AST, already decoded; this only finds where it ends.
func scanLink(s string) (label, after string, err error) {
	depth, i := 1, 1
	for ; i < len(s) && depth > 0; i++ {
		switch s[i] {
		case '\\':
			i++
		case '[':
			depth++
		case ']':
			depth--
		}
	}
	if depth != 0 {
		return "", "", errors.New("the link label is not closed")
	}
	label = s[1 : i-1]
	if i >= len(s) || s[i] != '(' {
		return "", "", errors.New("reference-style links are not supported; write [text](address)")
	}
	i++
	for i < len(s) && s[i] == ' ' {
		i++
	}
	if i < len(s) && s[i] == '<' {
		for i++; i < len(s) && s[i] != '>'; i++ {
			if s[i] == '\\' {
				i++
			}
		}
		i++
	} else {
		parens := 0
		for ; i < len(s); i++ {
			if s[i] == '\\' {
				i++
				continue
			}
			if s[i] == ' ' || (s[i] == ')' && parens == 0) {
				break
			}
			if s[i] == '(' {
				parens++
			} else if s[i] == ')' {
				parens--
			}
		}
	}
	for i < len(s) && s[i] == ' ' {
		i++
	}
	if i >= len(s) || s[i] != ')' {
		return "", "", errors.New("a link title is not supported; write [text](address)")
	}
	return label, s[i+1:], nil
}

// splitBlock separates the trailing {key="value"} block from the text. A
// block is the last { that is not escaped and that starts a block running to
// the end of the line; text may hold other braces.
func splitBlock(s string) (text, block string) {
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == '\\':
			i++
		case s[i] == '{' && (i == 0 || s[i-1] == ' '):
			if _, ok := parseBlock(s[i:]); ok {
				return strings.TrimSpace(s[:i]), s[i:]
			}
		}
	}
	return s, ""
}

type pair struct{ key, value string }

// parseBlock reads {key="value" key="value"} that runs to the end of s.
func parseBlock(s string) ([]pair, bool) {
	if !strings.HasPrefix(s, "{") {
		return nil, false
	}
	var pairs []pair
	i := 1
	for {
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i < len(s) && s[i] == '}' {
			return pairs, len(pairs) > 0 && i == len(s)-1
		}
		start := i
		for i < len(s) && s[i] != '=' && s[i] != ' ' && s[i] != '}' {
			i++
		}
		key := s[start:i]
		if !attrName.MatchString(key) || i+1 >= len(s) || s[i] != '=' || s[i+1] != '"' {
			return nil, false
		}
		i += 2
		var v strings.Builder
		for ; i < len(s) && s[i] != '"'; i++ {
			if s[i] == '\\' && i+1 < len(s) && (s[i+1] == '\\' || s[i+1] == '"') {
				i++
			}
			v.WriteByte(s[i])
		}
		if i >= len(s) {
			return nil, false
		}
		i++
		pairs = append(pairs, pair{key, v.String()})
	}
}

// assign sets the attributes of an outline from a block, in order; the ones
// the opml.Outline struct has no field for are kept as custom attributes.
func assign(ol *opml.Outline, pairs []pair) error {
	seen := map[string]bool{}
	for _, p := range pairs {
		if seen[p.key] {
			return fmt.Errorf("the attribute %s is given twice", p.key)
		}
		seen[p.key] = true
		switch p.key {
		case "text":
			return errors.New("text is the item text, not an attribute")
		case "type":
			ol.Type = p.value
		case "title":
			ol.Title = p.value
		case "isComment", "isBreakpoint":
			if p.value != "true" && p.value != "false" {
				return fmt.Errorf("%s is \"true\" or \"false\", not %q", p.key, p.value)
			}
			if p.key == "isComment" {
				ol.IsComment = p.value == "true"
			} else {
				ol.IsBreakpoint = p.value == "true"
			}
		case "created":
			ol.Created = p.value
		case "category":
			ol.Category = p.value
		case "xmlUrl":
			ol.XMLURL = p.value
		case "htmlUrl":
			ol.HTMLURL = p.value
		case "language":
			ol.Language = p.value
		case "description":
			ol.Description = p.value
		case "version":
			ol.Version = p.value
		case "url":
			ol.URL = p.value
		default:
			name := xml.Name{Local: p.key}
			if i := strings.LastIndex(p.key, ":"); i > 0 {
				name = xml.Name{Space: p.key[:i], Local: p.key[i+1:]}
			}
			ol.OtherAttr = append(ol.OtherAttr, xml.Attr{Name: name, Value: p.value})
		}
	}
	return nil
}

// unescape resolves CommonMark backslash escapes: a backslash before an ASCII
// punctuation character is dropped. Other markup stays as typed.
func unescape(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && strings.IndexByte("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", s[i+1]) >= 0 {
			i++
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
