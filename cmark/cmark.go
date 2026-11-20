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
// Package cmark converts an OPML document to a CommonMark list and back, so
// an outline can be edited in any text editor (opml DR-0001). The opml
// package itself imports only the standard library; this package is the one
// that depends on a YAML library (and, for reading, goldmark).
package cmark

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/rsdoiel/opml"
	"gopkg.in/yaml.v3"
)

/** ErrUnrepresentable is wrapped by the error FromOPML returns for content a
 * Markdown list item cannot hold: a newline or tab in text, leading or
 * trailing white space, a newline in a URL or attribute value, an attribute
 * name that is not a plain name. The error names the outline.
 *
 * Example:
 *   _, err := cmark.FromOPML(doc)
 *   if errors.Is(err, cmark.ErrUnrepresentable) {
 *       fmt.Println(err) // names the outline, such as "outline 2.1: ..."
 *   }
 */
var ErrUnrepresentable = errors.New("cannot be written as a Markdown list item")

// frontMatter holds the head fields, in the order they are written. An empty
// field is left out; opmlVersion is left out when it is the usual 2.0.
type frontMatter struct {
	OPMLVersion     *string `yaml:"opmlVersion,omitempty"`
	Title           string  `yaml:"title,omitempty"`
	DateCreated     string  `yaml:"dateCreated,omitempty"`
	DateModified    string  `yaml:"dateModified,omitempty"`
	Generator       string  `yaml:"generator,omitempty"`
	OwnerName       string  `yaml:"ownerName,omitempty"`
	OwnerEmail      string  `yaml:"ownerEmail,omitempty"`
	OwnerID         string  `yaml:"ownerId,omitempty"`
	Docs            string  `yaml:"docs,omitempty"`
	ExpansionState  string  `yaml:"expansionState,omitempty"`
	VertScrollState int     `yaml:"vertScrollState,omitempty"`
	WindowTop       int     `yaml:"windowTop,omitempty"`
	WindowLeft      int     `yaml:"windowLeft,omitempty"`
	WindowBottom    int     `yaml:"windowBottom,omitempty"`
	WindowRight     int     `yaml:"windowRight,omitempty"`
}

/** FromOPML writes an OPML document as CommonMark: the head fields as YAML
 * front matter (when there are any), then one list item per outline with
 * children as a nested list. The output is the canonical form of opml
 * DR-0001: "[text](url)" when the outline has a URL for the link (xmlUrl for
 * type rss or atom, url otherwise), plain text when it has none, and the
 * remaining attributes in a trailing {key="value"} block.
 *
 * Parameters:
 *   o (*opml.OPML) — the document
 *
 * Returns:
 *   []byte — the Markdown, ending in a newline; empty for a document with no
 *     head fields and no outlines
 *   error — nil, or an error wrapping ErrUnrepresentable that names the outline
 *
 * Example:
 *   doc, _ := opml.ReadFile("feeds.opml")
 *   src, err := cmark.FromOPML(doc)
 */
func FromOPML(o *opml.OPML) ([]byte, error) {
	if o == nil {
		return nil, errors.New("no OPML document")
	}
	var out strings.Builder
	if o.Head != nil || o.Version != "2.0" {
		fm, err := headFrontMatter(o)
		if err != nil {
			return nil, err
		}
		out.WriteString(fm)
	}
	if o.Body != nil {
		if err := writeOutlines(&out, o.Body.Outline, 0, nil); err != nil {
			return nil, err
		}
	}
	return []byte(out.String()), nil
}

// headFrontMatter returns the front matter block and the blank line after it,
// or "" when the head has nothing to write.
func headFrontMatter(o *opml.OPML) (string, error) {
	var fm frontMatter
	if o.Version != "2.0" {
		v := o.Version
		fm.OPMLVersion = &v
	}
	if h := o.Head; h != nil {
		fm.Title, fm.DateCreated, fm.DateModified = h.Title, h.Created, h.Modified
		fm.Generator, fm.OwnerName, fm.OwnerEmail, fm.OwnerID = h.Generator, h.OwnerName, h.OwnerEmail, h.OwnerID
		fm.Docs, fm.ExpansionState = h.Docs, h.ExpansionState
		fm.VertScrollState, fm.WindowTop, fm.WindowLeft = h.VertScrollState, h.WindowTop, h.WindowLeft
		fm.WindowBottom, fm.WindowRight = h.WindowBottom, h.WindowRight
		for _, s := range []string{h.Title, h.Created, h.Modified, h.Generator, h.OwnerName, h.OwnerEmail, h.OwnerID, h.Docs, h.ExpansionState} {
			if strings.ContainsAny(s, "\r\n") {
				return "", fmt.Errorf("head: a value has a newline: %w", ErrUnrepresentable)
			}
		}
	}
	src, err := yaml.Marshal(&fm)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(string(src)) == "{}" {
		return "", nil
	}
	return "---\n" + string(src) + "---\n\n", nil
}

// writeOutlines writes one list level; path is the 1-based position of the
// parent, for error messages.
func writeOutlines(out *strings.Builder, outlines []*opml.Outline, depth int, path []int) error {
	for i, ol := range outlines {
		if ol == nil {
			continue
		}
		here := append(append([]int{}, path...), i+1)
		line, err := itemLine(ol)
		if err != nil {
			return fmt.Errorf("outline %s: %w", pathString(here), err)
		}
		out.WriteString(strings.Repeat("  ", depth))
		out.WriteString("-")
		if line != "" {
			out.WriteString(" " + line)
		}
		out.WriteString("\n")
		if err := writeOutlines(out, ol.Outline, depth+1, here); err != nil {
			return err
		}
	}
	return nil
}

func pathString(p []int) string {
	parts := make([]string, len(p))
	for i, n := range p {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ".")
}

// isFeedType reports whether the link of an outline of this type is its
// xmlUrl (opml DR-0001, decision 5).
func isFeedType(t string) bool {
	t = strings.ToLower(t)
	return t == "rss" || t == "atom"
}

// itemLine is the text of a list item, after "- ": the text or link, then
// the attribute block.
func itemLine(ol *opml.Outline) (string, error) {
	if err := checkText(ol.Text); err != nil {
		return "", fmt.Errorf("text %w", err)
	}
	linkAttr, linkURL := "url", ol.URL
	if isFeedType(ol.Type) {
		linkAttr, linkURL = "xmlUrl", ol.XMLURL
	}
	var parts []string
	if linkURL != "" {
		dest, err := formatDestination(linkURL)
		if err != nil {
			return "", fmt.Errorf("%s %w", linkAttr, err)
		}
		parts = append(parts, "["+escape(ol.Text, false)+"]("+dest+")")
	} else if ol.Text != "" {
		parts = append(parts, escape(ol.Text, true))
	}
	block, err := attributeBlock(ol, linkAttr)
	if err != nil {
		return "", err
	}
	if block != "" {
		parts = append(parts, block)
	}
	return strings.Join(parts, " "), nil
}

// checkText rejects text a one-line list item would change.
func checkText(s string) error {
	if s == "" {
		return nil
	}
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' || unicode.IsControl(r) {
			return fmt.Errorf("has a control character such as a newline or tab: %w", ErrUnrepresentable)
		}
	}
	if first := []rune(s)[0]; unicode.IsSpace(first) {
		return fmt.Errorf("starts with white space: %w", ErrUnrepresentable)
	}
	rs := []rune(s)
	if unicode.IsSpace(rs[len(rs)-1]) {
		return fmt.Errorf("ends with white space: %w", ErrUnrepresentable)
	}
	return nil
}

var (
	entityLike   = regexp.MustCompile(`^&(#[0-9]{1,7}|#[xX][0-9a-fA-F]{1,6}|[A-Za-z][A-Za-z0-9]{1,31});`)
	leadingOList = regexp.MustCompile(`^[0-9]{1,9}[.)]( |$)`)
)

// escape backslash-escapes the characters of s that a reader of the list
// would take for Markdown. With leading set, s is at the start of an item
// and a block marker there (a heading, quote, bullet, number, rule, fence)
// is escaped too.
func escape(s string, leading bool) string {
	var b strings.Builder
	if leading && s != "" {
		if m := leadingOList.FindString(s); m != "" {
			i := strings.IndexAny(m, ".)")
			b.WriteString(s[:i] + `\`)
			s = s[i:]
		} else if strings.ContainsRune("#>+-=~", rune(s[0])) {
			b.WriteByte('\\')
		}
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\', '`', '*', '_', '[', ']', '<', '{':
			b.WriteByte('\\')
		case '&':
			if entityLike.MatchString(s[i:]) {
				b.WriteByte('\\')
			}
		}
		b.WriteByte(c)
	}
	return b.String()
}

// escapeText escapes the text of an item that has no link.
func escapeText(s string) string { return escape(s, true) }

// formatDestination writes a URL as the destination of a Markdown link: plain
// when it can be, in angle brackets when it holds a space or a parenthesis.
func formatDestination(u string) (string, error) {
	for _, r := range u {
		if r == '\n' || r == '\r' || (unicode.IsControl(r) && r != ' ') {
			return "", fmt.Errorf("has a control character such as a newline: %w", ErrUnrepresentable)
		}
	}
	angle := strings.ContainsAny(u, " ()<>")
	var b strings.Builder
	if angle {
		b.WriteByte('<')
	}
	for i := 0; i < len(u); i++ {
		c := u[i]
		switch {
		case c == '\\', angle && (c == '<' || c == '>'):
			b.WriteByte('\\')
		case c == '&' && entityLike.MatchString(u[i:]):
			b.WriteByte('\\')
		}
		b.WriteByte(c)
	}
	if angle {
		b.WriteByte('>')
	}
	return b.String(), nil
}

var attrName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.:-]*$`)

// attributeBlock writes {key="value" ...} for every attribute the list syntax
// does not carry, in the order of opml.Outline, then the custom ones. skip is
// the URL attribute already written as the link.
func attributeBlock(ol *opml.Outline, skip string) (string, error) {
	var pairs []string
	add := func(name, value string) error {
		if strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("attribute %s has a newline: %w", name, ErrUnrepresentable)
		}
		v := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(value)
		pairs = append(pairs, name+`="`+v+`"`)
		return nil
	}
	type kv struct{ name, value string }
	known := []kv{
		{"type", ol.Type}, {"title", ol.Title},
		{"isComment", boolString(ol.IsComment)}, {"isBreakpoint", boolString(ol.IsBreakpoint)},
		{"created", ol.Created}, {"category", ol.Category},
		{"xmlUrl", ol.XMLURL}, {"htmlUrl", ol.HTMLURL},
		{"language", ol.Language}, {"description", ol.Description},
		{"version", ol.Version}, {"url", ol.URL},
	}
	for _, a := range known {
		if a.value == "" || a.name == skip {
			continue
		}
		if err := add(a.name, a.value); err != nil {
			return "", err
		}
	}
	for _, a := range ol.OtherAttr {
		name := a.Name.Local
		if a.Name.Space != "" {
			name = a.Name.Space + ":" + name
		}
		if !attrName.MatchString(name) {
			return "", fmt.Errorf("attribute name %q: %w", name, ErrUnrepresentable)
		}
		if err := add(name, a.Value); err != nil {
			return "", err
		}
	}
	if len(pairs) == 0 {
		return "", nil
	}
	return "{" + strings.Join(pairs, " ") + "}", nil
}

func boolString(b bool) string {
	if b {
		return "true"
	}
	return ""
}
