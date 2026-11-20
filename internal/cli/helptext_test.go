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
package cli

import (
	"strings"
	"testing"

	"github.com/rsdoiel/opml"
)

// The help text is the man page source and the website page. It starts with
// a YAML front matter block, not a pandoc "%" title block: pandoc reads both,
// and a Markdown-web reader strips the front matter but would show the "%"
// lines as a paragraph.
func TestHelpPagesStartWithYAMLFrontMatter(t *testing.T) {
	for name, text := range map[string]string{
		"opml2json": Opml2jsonHelpText, "opml2urls": Opml2urlsHelpText, "opmlcat": OpmlcatHelpText,
		"opmlsort": OpmlsortHelpText, "urls2opml": Urls2opmlHelpText,
		"opml2cmark": Opml2cmarkHelpText, "cmark2opml": Cmark2opmlHelpText,
	} {
		page := opml.FmtHelp(text, name, "0.0.12", "2026-11-20", "abc1234")
		lines := strings.Split(page, "\n")
		if lines[0] != "---" {
			t.Errorf("%s: first line is %q, want ---", name, lines[0])
			continue
		}
		end := 0
		for i, l := range lines[1:] {
			if l == "---" {
				end = i + 1
				break
			}
		}
		if end == 0 {
			t.Errorf("%s: the front matter is never closed", name)
			continue
		}
		fm := strings.Join(lines[1:end], "\n")
		for _, want := range []string{
			"title: " + name, "section: 1", "header: user manual",
			"footer: version 0.0.12 abc1234", "date: 2026-11-20", "author: R. S. Doiel",
		} {
			if !strings.Contains(fm+"\n", want+"\n") {
				t.Errorf("%s: front matter lacks %q:\n%s", name, want, fm)
			}
		}
		for i, l := range lines {
			if strings.HasPrefix(l, "%") {
				t.Errorf("%s: line %d starts with %%, a pandoc title block: %q", name, i+1, l)
			}
		}
		if strings.Contains(page, "{app_name}") || strings.Contains(page, "{version}") {
			t.Errorf("%s: a token was left unreplaced", name)
		}
	}
}
