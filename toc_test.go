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
	"reflect"
	"testing"
)

// tocSrc is a table of contents of the shape antennaApp writes and Dave
// Winer publishes: a month outline with a url and children, entries with
// no url, and an include.
const tocSrc = `<opml version="2.0"><head/><body>
<outline text="Pages">
  <outline text="Home" url="https://example.org/index.html"/>
  <outline text="About" url="https://example.org/about.html"/>
</outline>
<outline text="Posts" type="calendarYear">
  <outline text="November 2026" type="calendarMonth" url="https://example.org/2026/11/index.html">
    <outline text="Nov 1" url="https://example.org/2026/11/01/a.html"/>
    <outline text="Nov 2" url="https://example.org/2026/11/02/b.html"/>
  </outline>
</outline>
<outline text="More" type="include" url="https://example.org/more/toc.opml"/>
<outline text="Archive" type="include" url="https://example.org/old/toc.opml"/>
</body></opml>`

func TestFlatten(t *testing.T) {
	o, err := Parse([]byte(tocSrc))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, ol := range o.Flatten() {
		got = append(got, ol.URL)
	}
	want := []string{
		"https://example.org/index.html",
		"https://example.org/about.html",
		"https://example.org/2026/11/index.html",
		"https://example.org/2026/11/01/a.html",
		"https://example.org/2026/11/02/b.html",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flatten urls = %v, want %v", got, want)
	}
}

func TestIncludes(t *testing.T) {
	o, err := Parse([]byte(tocSrc))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"https://example.org/more/toc.opml", "https://example.org/old/toc.opml"}
	if got := o.Includes(); !reflect.DeepEqual(got, want) {
		t.Errorf("Includes = %v, want %v", got, want)
	}
}

func TestFlattenAndIncludesOnEmptyDocuments(t *testing.T) {
	for name, o := range map[string]*OPML{
		"zero value": {},
		"New":        New(),
		"nil body":   {Body: nil},
	} {
		if got := o.Flatten(); len(got) != 0 {
			t.Errorf("%s: Flatten = %v, want none", name, got)
		}
		if got := o.Includes(); len(got) != 0 {
			t.Errorf("%s: Includes = %v, want none", name, got)
		}
	}
}

func TestIncludeWithoutURLIsIgnored(t *testing.T) {
	o, err := Parse([]byte(`<opml version="2.0"><head/><body><outline text="x" type="include"/></body></opml>`))
	if err != nil {
		t.Fatal(err)
	}
	if got := o.Includes(); len(got) != 0 {
		t.Errorf("Includes = %v, want none", got)
	}
}
