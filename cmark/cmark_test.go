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
	"errors"
	"strings"
	"testing"

	"github.com/rsdoiel/opml"
)

func mustParse(t *testing.T, src string) *opml.OPML {
	t.Helper()
	o, err := opml.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	return o
}

// A table of contents: nesting, a month that has both a url and children,
// attributes after the text, and includes (opml DR-0001, decisions 1, 5, 6).
func TestFromOPMLTableOfContents(t *testing.T) {
	o := mustParse(t, `<opml version="2.0"><head><title>Table of contents</title></head><body>
<outline text="Pages">
  <outline text="Home" url="https://example.org/index.html"/>
  <outline text="About" url="https://example.org/about.html"/>
</outline>
<outline text="Posts" type="calendarYear">
  <outline text="November 2026" type="calendarMonth" name="11" url="https://example.org/2026/11/index.html">
    <outline text="Nov 1" url="https://example.org/2026/11/01/a.html"/>
  </outline>
</outline>
<outline text="More" type="include" url="https://example.org/more/toc.opml"/>
</body></opml>`)
	want := `---
title: Table of contents
---

- Pages
  - [Home](https://example.org/index.html)
  - [About](https://example.org/about.html)
- Posts {type="calendarYear"}
  - [November 2026](https://example.org/2026/11/index.html) {type="calendarMonth" name="11"}
    - [Nov 1](https://example.org/2026/11/01/a.html)
- [More](https://example.org/more/toc.opml) {type="include"}
`
	got, err := FromOPML(o)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// A subscription list reads as a list of feed links: the link is xmlUrl for
// type rss or atom, and title is written even when it equals text.
func TestFromOPMLFeedList(t *testing.T) {
	o := mustParse(t, `<opml version="1.0"><head/><body>
<outline text="BBC News" title="BBC News" type="rss" xmlUrl="http://feeds.example/bbc.xml" htmlUrl="http://www.example/bbc"/>
<outline text="Atom one" type="atom" xmlUrl="http://a.example/atom" url="http://a.example/page"/>
<outline text="Rss with no feed" type="rss" url="http://r.example/p"/>
<outline text="Plain with feed" xmlUrl="http://p.example/feed"/>
</body></opml>`)
	want := `---
opmlVersion: "1.0"
---

- [BBC News](http://feeds.example/bbc.xml) {type="rss" title="BBC News" htmlUrl="http://www.example/bbc"}
- [Atom one](http://a.example/atom) {type="atom" url="http://a.example/page"}
- Rss with no feed {type="rss" url="http://r.example/p"}
- Plain with feed {xmlUrl="http://p.example/feed"}
`
	got, err := FromOPML(o)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFromOPMLHead(t *testing.T) {
	o := mustParse(t, `<opml version="2.0"><head><title>Why: a title</title><dateCreated>2016-05-23</dateCreated><dateModified>Fri, 20 Nov 2026 14:06:24 GMT</dateModified><generator>opmlToMarkdown v0.5.26</generator><ownerName>J. Doe</ownerName><windowTop>5</windowTop></head><body><outline text="x"/></body></opml>`)
	want := `---
title: 'Why: a title'
dateCreated: "2016-05-23"
dateModified: Fri, 20 Nov 2026 14:06:24 GMT
generator: opmlToMarkdown v0.5.26
ownerName: J. Doe
windowTop: 5
---

- x
`
	got, err := FromOPML(o)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFromOPMLNoHeadFieldsNoFrontMatter(t *testing.T) {
	got, err := FromOPML(mustParse(t, `<opml version="2.0"><head/><body><outline text="a"><outline text="b"/></outline></body></opml>`))
	if err != nil {
		t.Fatal(err)
	}
	if want := "- a\n  - b\n"; string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFromOPMLEmptyBody(t *testing.T) {
	got, err := FromOPML(opml.New())
	if err != nil || len(got) != 0 {
		t.Errorf("empty document: %q, %v; want nothing", got, err)
	}
	if _, err := FromOPML(nil); err == nil {
		t.Error("a nil document should be an error")
	}
}

func TestEscapeText(t *testing.T) {
	for in, want := range map[string]string{
		"plain text":              "plain text",
		"a*b_c":                   `a\*b\_c`,
		"[x] y":                   `\[x\] y`,
		"1. first":                `1\. first`,
		"2) second":               `2\) second`,
		"10.5 apples":             "10.5 apples",
		"# not a heading":         `\# not a heading`,
		"- dash":                  `\- dash`,
		"+ plus":                  `\+ plus`,
		"> quote":                 `\> quote`,
		"=== setext":              `\=== setext`,
		"~~~ fence":               `\~~~ fence`,
		"x {y} z":                 `x \{y} z`,
		`a\b`:                     `a\\b`,
		"<b>&amp; ok":             `\<b>\&amp; ok`,
		"AT&T":                    "AT&T",
		"`code`":                  "\\`code\\`",
		"BBC | Programmes":        "BBC | Programmes",
		"Vox Product Blog -  All": "Vox Product Blog -  All",
		"___":                     `\_\_\_`,
		"":                        "",
	} {
		if got := escapeText(in); got != want {
			t.Errorf("escapeText(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFormatDestination(t *testing.T) {
	for in, want := range map[string]string{
		"http://a/b?x=1&y=2": "http://a/b?x=1&y=2",
		"http://a/b?x&copy;": `http://a/b?x\&copy;`,
		"http://a/b c":       "<http://a/b c>",
		"http://a/(x)":       "<http://a/(x)>",
		`http://a/b\c`:       `http://a/b\\c`,
		"http://a/<b>":       `<http://a/\<b\>>`,
	} {
		got, err := formatDestination(in)
		if err != nil || got != want {
			t.Errorf("formatDestination(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	if _, err := formatDestination("http://a/b\nc"); !errors.Is(err, ErrUnrepresentable) {
		t.Errorf("a newline in a url: err = %v, want ErrUnrepresentable", err)
	}
}

// Text a list item cannot hold is an error that names the outline, not a
// silent change (opml DR-0001, open points).
func TestFromOPMLUnrepresentable(t *testing.T) {
	for name, text := range map[string]string{
		"newline":        "two&#10;lines",
		"leading space":  " lead",
		"trailing space": "trail ",
		"tab":            "a&#9;b",
	} {
		o := mustParse(t, `<opml version="2.0"><head/><body><outline text="ok"><outline text="`+text+`"/></outline></body></opml>`)
		_, err := FromOPML(o)
		if !errors.Is(err, ErrUnrepresentable) {
			t.Errorf("%s: err = %v, want ErrUnrepresentable", name, err)
			continue
		}
		if !strings.Contains(err.Error(), "outline 1.1") {
			t.Errorf("%s: %q does not name the outline (want \"outline 1.1\")", name, err)
		}
	}
}

func TestAttributeValuesAreEscaped(t *testing.T) {
	o := mustParse(t, `<opml version="2.0"><head/><body><outline text="x" description='say "hi" \ there'/></body></opml>`)
	got, err := FromOPML(o)
	if err != nil {
		t.Fatal(err)
	}
	if want := `- x {description="say \"hi\" \\ there"}` + "\n"; string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
