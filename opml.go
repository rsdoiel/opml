//
// Package opml provides basic utility functions for working with OPML files.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2016 all rights reserved.
// Released under the BSD 3-Clause License
// See: https://opensource.org/licenses/BSD-3-Clause
//
package opml

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

const (
	Version = "0.0.2"
)

// OPML is the root structure for holding an OPML document
type OPML struct {
	XMLName xml.Name `xml:"opml" json:"-"`
	Version string   `xml:"version,attr" json:"version"`
	Head    *Head    `xml:"head" json:"head"`
	Body    *Body    `xml:"body" json:"body"`
}

// Head holds the metadata for an OPML document
type Head struct {
	XMLName         xml.Name `json:"-"`
	Title           string   `xml:"title,omitempty" json:"title,omitempty"`
	Created         string   `xml:"dateCreated,omitempty" json:"dateCreated,omitempty"`   // RFC 882 date and time
	Modified        string   `xml:"dateModified,omitempty" json:"dataModified,omitempty"` // RFC 882 date and time
	OwnerName       string   `xml:"ownerName,omitempty" json:"ownerName,omitempty"`
	OwnerEmail      string   `xml:"ownerEmail,omitempty" json:"ownerEmail,omitempty"`
	OwnerID         string   `xml:"OwnerId,omitempty" json:"OwnerId,omitempty"`               // url
	Docs            string   `xml:"docs,omitempty" json:"docs,omitempty"`                     // url
	ExpansionState  string   `xml:"expansionState,omitempty" json:"expansionState,omitempty"` // array of numbers
	VertScrollState int      `xml:"vertScrollState,omitempty" json:"vertScrollState,omitempty"`
	WindowTop       int      `xml:"windowTop,omitempty" json:"windowTop,omitempty"`
	WindowLeft      int      `xml:"windowLeft,omitempty" json:"windowLeft,omitempty"`
	WindowBottom    int      `xml:"windowBottom,omitempty" json:"windowBottom,omitempty"`
	WindowRight     int      `xml:"windowRight,omitempty" json:"windowRight,omitempty"`
}

// Body holds the outline for an OPML document
type Body struct {
	XMLName xml.Name    `json:"-"`
	Outline OutlineList `xml:"outline" json:"outline"`
}

// Outline is the primary element of an OPML document, may hold sub-Outlines
type Outline struct {
	XMLName      xml.Name    `json:"-"`
	Text         string      `xml:"text,attr" json:"text"`
	Type         string      `xml:"type,attr,omitempty" json:"type,omitempty"`
	IsComment    bool        `xml:"isComment,attr,omitempty" json:"isComment,omitempty"`
	IsBreakpoint bool        `xml:"isBreakpoint,attr,omitempty" json:"isBreakpoint,omitempty"`
	Created      string      `xml:"created,attr,omitempty" json:"created,omitempty"` // RFC 882 date and time
	Category     string      `xml:"category,attr,omitempty" json:"category,omitempty"`
	XMLURL       string      `xml:"xmlUrl,attr,omitempty" json:"xmlUrl,omitempty"`   // url
	HTMLURL      string      `xml:"htmlUrl,attr,omitempty" json:"htmlUrl,omitempty"` // url
	Language     string      `xml:"langauge,attr,omitempty" json:"language,omitempty"`
	Description  string      `xml:"description,attr,omitempty" json:"description,omitempty"`
	Version      string      `xml:"version,attr,omitempty" json:"version,omitempty"`
	URL          string      `xml:"url,attr,omitempty" json:"url,omitempty"` // url
	Outline      OutlineList `xml:"outline,omitempty" json:"outline,omitempty"`
}

type OutlineList []Outline
type ByText []Outline
type ByType []Outline

// New creates an empty OPML structure
func New() *OPML {
	o := new(OPML)
	o.Version = "2.0"

	o.Head = new(Head)
	o.Body = new(Body)
	return o
}

func (h *Head) String() string {
	s, _ := xml.Marshal(h)
	return string(s)
}

func (b *Body) String() string {
	s, _ := xml.Marshal(b)
	return string(s)
}

func (ol *Outline) String() string {
	s, _ := xml.Marshal(ol)
	return string(s)
}

// HasChildren return true if the outline element has a populated child outline
func (ol *Outline) HasChildren() bool {
	if len(ol.Outline) > 0 {
		return true
	}
	return false
}

func (ol OutlineList) Append(elem *Outline) error {
	i := len(ol)
	ol = append(ol, *elem)
	if len(ol) != (i+1) || ol[i].Text != elem.Text {
		return fmt.Errorf("failed to append element")
	}
	return nil
}

func (ol *Outline) AppendChild(elem *Outline) error {
	return ol.Outline.Append(elem)
}

func (o *OPML) String() string {
	if o.Body != nil {
		if o.Body.Outline == nil {
			o.Body.Outline = make(OutlineList, 1)
			o.Body.Outline.Append(&Outline{
				Text: "",
			})
		} else if len(o.Body.Outline) == 0 {
			o.Body.Outline.Append(&Outline{
				Text: "",
			})
		}
	}
	s, _ := xml.Marshal(o)
	return string(s)
}

// Len for ByText sort of Outline
func (a ByText) Len() int {
	return len(a)
}

// Swap for ByText sort of Outline
func (a ByText) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less for ByText sort of Outline
func (a ByText) Less(i, j int) bool {
	return strings.Compare(a[i].Text, a[j].Text) == -1
}

// Sort do a recursive sort over an outline
func (a ByText) Sort() {
	if len(a) > 0 {
		for _, item := range a {
			if len(item.Outline) > 0 {
				ol := ByText(item.Outline)
				ol.Sort()
			}
		}
		sort.Sort(ByText(a))
	}
}

// Sort do a recursive ByText sort of outline elements starting at the OPML struct.
func (o *OPML) Sort() {
	if o.Body != nil && len(o.Body.Outline) > 0 {
		ol := ByText(o.Body.Outline)
		ol.Sort()
	}
}

// Len for ByType sort of Outline
func (a ByType) Len() int {
	return len(a)
}

// Swap for ByType sort of Outline
func (a ByType) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less for ByType sort of Outline
func (a ByType) Less(i, j int) bool {
	return strings.Compare(a[i].Type, a[j].Type) == -1
}

// Sort do a recursive sort over an outline
func (a ByType) Sort() {
	if len(a) > 0 {
		for _, item := range a {
			if len(item.Outline) > 0 {
				ol := ByType(item.Outline)
				ol.Sort()
			}
		}
		sort.Sort(ByType(a))
	}
}

// SortTypes do a recursive ByText sort of outline elements starting at the OPML struct.
func (o *OPML) SortTypes() {
	if o.Body != nil && len(o.Body.Outline) > 0 {
		ol := ByType(o.Body.Outline)
		ol.Sort()
	}
}

// ReadFile reads an OPML file and populates the OPML object appropriately
func (o *OPML) ReadFile(s string) error {
	src, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}
	return xml.Unmarshal(src, &o)
}

// WriteFile writes the contents of a OPML struct to a file
func (o *OPML) WriteFile(s string, perm os.FileMode) error {
	if len(o.Body.Outline) == 0 {
		o.Body.Outline = append(o.Body.Outline, Outline{
			Text: "",
		})
	}
	b, _ := xml.Marshal(o)
	return ioutil.WriteFile(s, b, perm)
}
