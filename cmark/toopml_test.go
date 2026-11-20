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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rsdoiel/opml"
)

// model returns the document as the bytes the opml package writes, which is
// the comparison "the same OPML model" (opml DR-0001, decision 8).
func model(t *testing.T, o *opml.OPML) string {
	t.Helper()
	b, err := opml.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func wantXML(t *testing.T, xml string) string {
	t.Helper()
	return model(t, mustParse(t, xml))
}

func TestToOPMLCanonical(t *testing.T) {
	md := `---
title: My list
opmlVersion: "1.0"
---

- [BBC](http://f/x.xml) {type="rss" title="BBC" htmlUrl="http://h/"}
- Parent {type="grp" name="p 1"}
  - child \* star
  - [Month](https://e.org/m.html) {type="calendarMonth"}
- [More](https://e.org/toc.opml) {type="include"}
`
	o, err := ToOPML([]byte(md))
	if err != nil {
		t.Fatal(err)
	}
	want := wantXML(t, `<opml version="1.0"><head><title>My list</title></head><body>
<outline text="BBC" type="rss" title="BBC" xmlUrl="http://f/x.xml" htmlUrl="http://h/"/>
<outline text="Parent" type="grp" name="p 1"><outline text="child * star"/><outline text="Month" type="calendarMonth" url="https://e.org/m.html"/></outline>
<outline text="More" type="include" url="https://e.org/toc.opml"/>
</body></opml>`)
	if got := model(t, o); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

// What goldmark lets a hand-edited list get away with (opml DR-0001, decision 7).
func TestToOPMLLenientForms(t *testing.T) {
	for name, c := range map[string]struct{ md, body string }{
		"other bullets":    {"* a\n  * b\n+ c\n", `<outline text="a"><outline text="b"/></outline><outline text="c"/>`},
		"four space":       {"- a\n    - b\n", `<outline text="a"><outline text="b"/></outline>`},
		"ordered":          {"1. a\n2. b\n", `<outline text="a"/><outline text="b"/>`},
		"loose":            {"- a\n\n- b\n\n  - c\n", `<outline text="a"/><outline text="b"><outline text="c"/></outline>`},
		"wrapped":          {"- first line\n  second line\n", `<outline text="first line second line"/>`},
		"crlf":             {"- a\r\n  - b\r\n", `<outline text="a"><outline text="b"/></outline>`},
		"no final newline": {"- a", `<outline text="a"/>`},
		"empty parent":     {"-\n  - b\n", `<outline text=""><outline text="b"/></outline>`},
		"escaped brace":    {"- x \\{y} z\n", `<outline text="x {y} z"/>`},
		"plain braces":     {"- x {y} z\n", `<outline text="x {y} z"/>`},
		"unclosed block":   {"- x {type=\"a\n", `<outline text="x {type=&#34;a"/>`},
		"typed markup":     {"- *emph* here\n", `<outline text="*emph* here"/>`},
		"link no attrs":    {"- [a](http://u/p)\n", `<outline text="a" url="http://u/p"/>`},
		"angle link":       {"- [a](<http://u/a b>)\n", `<outline text="a" url="http://u/a b"/>`},
		"escaped label":    {"- [a \\[1\\]](http://u)\n", `<outline text="a [1]" url="http://u"/>`},
		"empty document":   {"", ``},
		"blank lines only": {"\n\n", ``},
	} {
		o, err := ToOPML([]byte(c.md))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		want := wantXML(t, `<opml version="2.0"><head/><body>`+c.body+`</body></opml>`)
		if got := model(t, o); got != want {
			t.Errorf("%s:\n got  %s\n want %s", name, got, want)
		}
	}
}

func TestToOPMLHeadOnly(t *testing.T) {
	o, err := ToOPML([]byte("---\ntitle: Only a head\nwindowTop: 7\n---\n"))
	if err != nil {
		t.Fatal(err)
	}
	if o.Head.Title != "Only a head" || o.Head.WindowTop != 7 || len(o.Body.Outline) != 0 || o.Version != "2.0" {
		t.Errorf("head %+v version %q outlines %d", o.Head, o.Version, len(o.Body.Outline))
	}
}

// Content with no outline meaning is an error that names the line, never a
// silent drop (opml DR-0001, decision 7).
func TestToOPMLErrors(t *testing.T) {
	for name, c := range map[string]struct{ md, want string }{
		"heading":           {"# Heading\n- a\n", "line 1"},
		"paragraph":         {"- a\n\nSome paragraph\n", "line 3"},
		"code in item":      {"- a\n  ```\n  code\n  ```\n", "line 2"},
		"second paragraph":  {"- a\n\n  more\n", "line 3"},
		"reference link":    {"- [x][ref]\n\n[ref]: http://u\n", "line 1"},
		"text after link":   {"- [x](http://u) trailing words\n", "line 1"},
		"duplicate key":     {"- x {type=\"a\" type=\"b\"}\n", "line 1"},
		"link and url key":  {"- [x](http://u) {url=\"http://v\"}\n", "line 1"},
		"link and xml key":  {"- [x](http://u) {type=\"rss\" xmlUrl=\"http://v\"}\n", "line 1"},
		"text key":          {"- x {text=\"y\"}\n", "line 1"},
		"bad bool":          {"- x {isComment=\"maybe\"}\n", "line 1"},
		"unknown front key": {"---\ntitle: t\nbogus: 1\n---\n\n- a\n", "bogus"},
		"thematic break":    {"- a\n\n---\n", "line 3"},
		"block quote":       {"> quote\n", "line 1"},
		"error after front": {"---\ntitle: t\n---\n\n- a\n\n# H\n", "line 7"},
		"link title":        {"- [x](http://u \"title\")\n", "line 1"},
	} {
		_, err := ToOPML([]byte(c.md))
		if !errors.Is(err, ErrNotOutline) {
			t.Errorf("%s: err = %v, want ErrNotOutline", name, err)
			continue
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %q does not mention %q", name, err, c.want)
		}
	}
}

// OPML to Markdown to OPML gives the same model, and Markdown to OPML to
// Markdown is byte-equal for canonical Markdown (decision 8).
func roundTrip(t *testing.T, name string, o *opml.OPML) {
	t.Helper()
	md, err := FromOPML(o)
	if err != nil {
		t.Errorf("%s: FromOPML: %v", name, err)
		return
	}
	back, err := ToOPML(md)
	if err != nil {
		t.Errorf("%s: ToOPML: %v\n%.300s", name, err, md)
		return
	}
	if got, want := model(t, back), model(t, o); got != want {
		t.Errorf("%s: the model changed in OPML -> Markdown -> OPML", name)
	}
	md2, err := FromOPML(back)
	if err != nil || string(md2) != string(md) {
		t.Errorf("%s: Markdown -> OPML -> Markdown is not byte-equal (err %v)", name, err)
	}
}

func TestRoundTripRealFiles(t *testing.T) {
	files, _ := filepath.Glob("../testdata/*.opml")
	if toc := filepath.Join(os.Getenv("HOME"), "Laboratory", "tmp", "dave-toc.opml"); fileExists(toc) {
		files = append(files, toc)
	}
	if len(files) < 5 {
		t.Fatalf("expected the testdata files, found %v", files)
	}
	for _, f := range files {
		o, err := opml.ReadFile(f)
		if err != nil {
			t.Errorf("%s: %v", f, err)
			continue
		}
		roundTrip(t, filepath.Base(f), o)
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// Awkward text, urls and attribute values from the escaping rules.
func TestRoundTripAwkwardValues(t *testing.T) {
	texts := []string{"plain", "a*b_c", "[x] y", "1. first", "2) second", "10.5 apples", "# h", "- d", "+ p", "> q",
		"=== s", "~~~ f", "x {y} z", `a\b`, "<b>&amp; ok", "AT&T", "`code`", "BBC | P", "___", "{type=\"a\"}", "end {", `\{`, "&copy; 2026",
		"日本語 テキスト", "a  b", "[a](b)", "![i](u)", "<http://x>", "***", "1.", "1)", "#", "-", "*", "}"}
	urls := []string{"http://a/b", "http://a/b?x=1&y=2", "http://a/b?x&copy;", "http://a/b c", "http://a/(x)", `http://a/b\c`, "http://a/<b>", "http://a/é"}
	attrs := []string{"v", `say "hi"`, `back\slash`, `\"`, "a b", "{x}", "}", "=", "é", "&amp;"}
	for _, text := range texts {
		for _, typ := range []string{"", "rss", "atom", "other"} {
			o := opml.New()
			ol := &opml.Outline{Text: text, Type: typ, Outline: []*opml.Outline{{Text: "child " + text}}}
			o.Body.Outline = []*opml.Outline{ol}
			for i, u := range urls {
				if i%2 == 0 {
					ol.URL = u
				} else {
					ol.XMLURL = u
				}
				ol.HTMLURL = urls[(i+1)%len(urls)]
				ol.Description = attrs[i%len(attrs)]
				ol.Category = attrs[(i+3)%len(attrs)]
				ol.Title = text
				roundTrip(t, text+"|"+typ+"|"+u, o)
			}
			ol.URL, ol.XMLURL = "", ""
			roundTrip(t, text+"|"+typ+"|nourl", o)
		}
	}
}

func TestRoundTripAttributesAndHead(t *testing.T) {
	o := mustParse(t, `<opml version="1.0"><head><title>Why: a title</title><dateCreated>2016-05-23</dateCreated><generator>g</generator><ownerName>J</ownerName><expansionState>1,3</expansionState><windowTop>5</windowTop></head><body>
<outline text="a" isComment="true" isBreakpoint="true" created="Mon" category="x,y" language="en" version="v" description="d" custom="2" other-one="1"/>
</body></opml>`)
	roundTrip(t, "attrs", o)
	o2 := mustParse(t, `<opml><head/><body><outline text="x"/></body></opml>`)
	o2.Version = ""
	roundTrip(t, "no version", o2)
}

// FuzzRoundTrip: any outline FromOPML can write, ToOPML reads back to the
// same model. Run with: go test -fuzz FuzzRoundTrip ./cmark
func FuzzRoundTrip(f *testing.F) {
	for _, s := range []string{"plain", "a*b_c", "[x](y)", "1. n", "# h", "{type=\"a\"}", `\`, "&amp;", "<b>", "x {y", "", "é"} {
		f.Add(s, "http://a/b?x=1&y=2", "rss", s)
		f.Add(s, "", "", s)
	}
	f.Fuzz(func(t *testing.T, text, url, typ, attr string) {
		o := opml.New()
		ol := &opml.Outline{Text: text, Type: typ, URL: url, XMLURL: url + "x", Description: attr,
			Outline: []*opml.Outline{{Text: attr, URL: url}}}
		o.Body.Outline = []*opml.Outline{ol}
		md, err := FromOPML(o)
		if errors.Is(err, ErrUnrepresentable) {
			t.Skip()
		}
		if err != nil {
			t.Fatalf("FromOPML: %v", err)
		}
		back, err := ToOPML(md)
		if err != nil {
			t.Fatalf("ToOPML: %v\n%q", err, md)
		}
		if got, want := model(t, back), model(t, o); got != want {
			t.Fatalf("model changed\n md   %q\n got  %s\n want %s", md, got, want)
		}
	})
}
