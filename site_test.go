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
package opml

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The documentation is a Markdown site published as HTML and as Markdown for
// BonBon browsers and cmarkbrowser. These tests keep it that way: every page is
// in the page order, the pages link to each other as Markdown, the table of
// contents lists them, and the built HTML points at its Markdown.

// noPublish are the Markdown files that are not website pages.
var noPublish = map[string]bool{"TODO.md": true, "index.md": true}

// websiteVar returns the value of a variable set on one line in website.mak.
func websiteVar(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("website.mak")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?m)^` + name + ` *= *(.*)$`).FindSubmatch(b)
	if m == nil {
		t.Fatalf("website.mak has no %s", name)
	}
	return strings.TrimSpace(string(m[1]))
}

func pageOrder(t *testing.T) []string { return strings.Fields(websiteVar(t, "PAGE_ORDER")) }

func publishedPages(t *testing.T) []string {
	t.Helper()
	all, _ := filepath.Glob("*.md")
	var out []string
	for _, f := range all {
		if !noPublish[f] {
			out = append(out, f)
		}
	}
	sort.Strings(out)
	return out
}

func readText(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestEveryPublishedPageIsInThePageOrder(t *testing.T) {
	order := append([]string(nil), pageOrder(t)...)
	sort.Strings(order)
	if got, want := strings.Join(order, " "), strings.Join(publishedPages(t), " "); got != want {
		t.Errorf("PAGE_ORDER in website.mak\n  %s\nbut the published pages are\n  %s", got, want)
	}
	if o := pageOrder(t); o[0] != "README.md" {
		t.Errorf("the home page must come first: %v", o)
	}
	for _, p := range strings.Fields(websiteVar(t, "PROGRAMS")) {
		found := false
		for _, o := range pageOrder(t) {
			found = found || o == p+".1.md"
		}
		if !found {
			t.Errorf("PAGE_ORDER lacks %s.1.md, the page of a program in PROGRAMS", p)
		}
	}
}

// The Markdown is the truth: a page links to another page as .md, and the
// HTML build turns that into .html (links-to-html.lua).
func TestPagesLinkToEachOtherAsMarkdown(t *testing.T) {
	link := regexp.MustCompile(`\]\(([^)\s]+)\)`)
	published := map[string]bool{}
	for _, p := range pageOrder(t) {
		published[p] = true
	}
	for _, page := range publishedPages(t) {
		for _, m := range link.FindAllStringSubmatch(readText(t, page), -1) {
			dest := m[1]
			if strings.Contains(dest, "://") || strings.HasPrefix(dest, "#") || strings.HasPrefix(dest, "mailto:") {
				continue
			}
			file := strings.SplitN(dest, "#", 2)[0]
			switch {
			case strings.HasSuffix(file, ".html"):
				t.Errorf("%s links to %s: link to the .md page, the HTML build rewrites it", page, dest)
			case strings.HasSuffix(file, ".md") && !published[file]:
				t.Errorf("%s links to %s, which is not a published page", page, dest)
			}
		}
	}
}

// A pandoc "%" title block and JSON front matter show as text in a reader
// that follows CommonMark; YAML front matter is stripped.
func TestPagesStartLikeCommonMark(t *testing.T) {
	for _, page := range publishedPages(t) {
		first := strings.SplitN(readText(t, page), "\n", 2)[0]
		if strings.HasPrefix(first, "%") || strings.HasPrefix(first, "{") {
			t.Errorf("%s starts with %q, not YAML front matter or text", page, first)
		}
	}
}

// pageTitle is the table of contents title the build gives a page.
func runToc(t *testing.T) []byte {
	t.Helper()
	args := append([]string{"./toc.bash", websiteVar(t, "BASE_URL"), websiteVar(t, "PROJECT")}, pageOrder(t)...)
	out, err := exec.Command("bash", args...).Output()
	if err != nil {
		t.Fatalf("toc.bash: %v", err)
	}
	return out
}

func TestTocListsThePagesInOrder(t *testing.T) {
	doc, err := Parse(runToc(t))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	titles := map[string]string{}
	for _, ol := range doc.Flatten() {
		got = append(got, ol.URL)
		titles[ol.URL] = ol.Text
	}
	var want []string
	for _, p := range pageOrder(t) {
		if p == "README.md" {
			p = "index.md"
		}
		want = append(want, websiteVar(t, "BASE_URL")+p)
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("toc urls\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	for url, title := range titles {
		if title == "" || title == "NAME" || strings.Contains(title, "*") || strings.HasSuffix(title, ".md") {
			t.Errorf("%s has the toc title %q: give the page a YAML title or a heading", url, title)
		}
	}
	if len(doc.Includes()) != 0 {
		t.Errorf("the toc has includes: %v", doc.Includes())
	}
}

// The head's dateModified is the date of the last commit that touched a page,
// so committing the toc changes what a rebuild writes: that one line is left
// out of the comparison, or the test could never pass after a commit.
func TestCommittedTocIsCurrent(t *testing.T) {
	withoutDate := func(s string) string {
		var keep []string
		for _, line := range strings.Split(s, "\n") {
			if !strings.Contains(line, "<dateModified>") {
				keep = append(keep, line)
			}
		}
		return strings.Join(keep, "\n")
	}
	if got, want := withoutDate(readText(t, ".well-known/toc.opml")), withoutDate(string(runToc(t))); got != want {
		t.Errorf(".well-known/toc.opml is out of date: run `make website`")
	}
}

func TestIndexMdIsReadmeAndSiteIsNotJekyll(t *testing.T) {
	if readText(t, "index.md") != readText(t, "README.md") {
		t.Error("index.md differs from README.md: run `make website`")
	}
	if _, err := os.Stat(".nojekyll"); err != nil {
		t.Error("no .nojekyll: GitHub Pages would hide /.well-known/")
	}
}

// Each HTML page names its Markdown, for a reader that finds the HTML first.
func TestHTMLPagesPointAtTheirMarkdown(t *testing.T) {
	for _, md := range pageOrder(t) {
		html := strings.TrimSuffix(md, ".md") + ".html"
		want := md
		if md == "README.md" {
			html, want = "index.html", "index.md"
		}
		page := readText(t, html)
		if !strings.Contains(page, `<link rel="alternate" type="text/markdown" href="`+want+`">`) {
			t.Errorf("%s has no alternate link to %s: run `make website`", html, want)
		}
	}
}

func TestLinksToHTMLFilter(t *testing.T) {
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("no pandoc")
	}
	src := "[a](user-manual.md#keys) [b](README.md) [c](https://example.org/x.md) [d](mailto:a@b.c) [e](#top) [f](other.html) [g](/abs.md)\n"
	cmd := exec.Command(pandoc, "--from", "markdown", "--to", "markdown", "--wrap=none", "--lua-filter=links-to-html.lua")
	cmd.Stdin = strings.NewReader(src)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	want := "[a](user-manual.html#keys) [b](index.html) [c](https://example.org/x.md) [d](mailto:a@b.c) [e](#top) [f](other.html) [g](/abs.md)"
	if got := strings.TrimSpace(out.String()); got != want {
		t.Errorf("links-to-html.lua\n got %s\nwant %s", got, want)
	}
}
