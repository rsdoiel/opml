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
	"io/ioutil"
)

type OPML struct {
	XMLName xml.Name `xml:"opml" json:"-"`
	Version string   `xml:"version,attr" json:"version"`
	Head    *Head    `xml:"head" json:"head"`
	Body    *Body    `xml:"body" json:"body"`
}

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

type Body struct {
	XMLName xml.Name  `json:"-"`
	Outline []Outline `xml:"outline" json:"outline"`
}

type Outline struct {
	XMLName      xml.Name  `json:"-"`
	Text         string    `xml:"text,attr" json:"text"`
	Type         string    `xml:"type,attr,omitempty" json:"type,omitempty"`
	IsComment    bool      `xml:"isComment,attr,omitempty" json:"isComment,omitempty"`
	IsBreakpoint bool      `xml:"isBreakpoint,attr,omitempty" json:"isBreakpoint,omitempty"`
	Created      string    `xml:"created,attr,omitempty" json:"created,omitempty"` // RFC 882 date and time
	Category     string    `xml:"category,attr,omitempty" json:"category,omitempty"`
	XmlURL       string    `xml:"xmlUrl,attr,omitempty" json:"xmlUrl,omitempty"`   // url
	HtmlURL      string    `xml:"htmlUrl,attr,omitempty" json:"htmlUrl,omitempty"` // url
	Language     string    `xml:"langauge,attr,omitempty" json:"language,omitempty"`
	Description  string    `xml:"description,attr,omitempty" json:"description,omitempty"`
	Version      string    `xml:"version,attr,omitempty" json:"version,omitempty"`
	URL          string    `xml:"url,attr,omitempty" json:"url,omitempty"` // url
	Outline      []Outline `xml:"outline,omitempty" json:"outline,omitempty"`
}

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

func (o *OPML) String() string {
	if len(o.Body.Outline) == 0 {
		o.Body.Outline = append(o.Body.Outline, Outline{
			Text: "",
		})
	}
	s, _ := xml.Marshal(o)
	return string(s)
}

func (o *OPML) ReadFile(s string) error {
	src, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}
	return xml.Unmarshal(src, &o)
}
