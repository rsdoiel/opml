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
	"time"
)

type OPML struct {
	XMLName xml.Name `json:"-"`
	Version string   `xml:"version,attr" json:"version"`
	Head    *Head    `xml:"head" json:"head"`
	Body    *Body    `xml:"body" json:"body"`
}

type Head struct {
	XMLName         xml.Name  `json:"-"`
	Title           string    `xml:"title,omitempty" json:"title,omitempty"`
	Created         time.Time `xml:"dateCreated,omitempty" json:"dateCreated,omitempty"`
	Modified        time.Time `xml:"dateModified,omitempty" json:"dataModified,omitempty"`
	OwnerName       string    `xml:"ownerName,omitempty" json:"ownerName,omitempty"`
	OwnerEmail      string    `xml:"ownerEmail,omitempty" json:"ownerEmail,omitempty"`
	OwnerID         url.URL   `xml:"OwnerId,omitempty" json:"OwnerId,omitempty"`
	Docs            url.URL   `xml:"docs,omitempty" json:"docs,omitempty"`
	ExpansionState  []int     `xml:"expansionState,omitempty" json:"expansionState,omitempty"`
	VertScrollState int       `xml:"vertScrollState,omitempty" json:"vertScrollState,omitempty"`
	WindowTop       int       `xml:"windowTop,omitempty" json:"windowTop,omitempty"`
	WindowLeft      int       `xml:"windowLeft,omitempty" json:"windowLeft,omitempty"`
	WindowBottom    int       `xml:"windowBottom,omitempty" json:"windowBottom,omitempty"`
	WindowRight     int       `xml:"windowRight,omitempty" json:"windowRight,omitempty"`
}

type Body struct {
	XMLName xml.Name   `json:"-"`
	Outline []*Outline `xml:"outline" json:"outline"`
}

type Outline struct {
	XMLName      xml.Name   `json:"-"`
	Text         string     `xml:"text,attr" json:"text"`
	Type         string     `xml:"type,attr" json:"type,omitempty"`
	IsComment    bool       `xml:"isComment,attr" json:"isComment,omitempty"`
	IsBreakpoint bool       `xml:"isBreakpoint,attr" json:"isBreakpoint,omitempty"`
	Created      time.Time  `xml:"created,attr" json:"created,omitempty"`
	Category     string     `xml:"category,attr" json:"category,omitempty"`
	XmlURL       url.URL    `xml:"xmlUrl,attr" json:"xmlUrl,omitempty"`
	HtmlURL      url.URL    `xml:"htmlUrl,attr" json:"htmlUrl,omitempty"`
	Language     string     `xml:"langauge,attr" json:"language,omitempty"`
	Description  string     `xml:"description,attr" json:"description,omitempty"`
	Version      string     `xml:"version,attr" json:"version,omitempty"`
	URL          url.URL    `xml:"url,attr" json:"url,omitempty"`
	Outline      []*Outline `xml:"outline,omitempty" json:"outline,omitempty"`
}
