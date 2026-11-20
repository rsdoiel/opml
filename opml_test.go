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
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNewAndString(t *testing.T) {
	o := New()
	if o == nil {
		t.Errorf("Can't create an opml structure")
	}

	head := o.Head
	if head == nil {
		t.Errorf("Can't find Head in opml structure")
	}

	body := o.Body
	if body == nil {
		t.Errorf("Can't find Body in opml structure")
	}
	o.Body.Outline = []*Outline{}

	expected := `<opml version="2.0"><head></head><body></body></opml>`
	s := o.String()
	if strings.Compare(expected, s) != 0 {
		t.Errorf("expected %q, got %q", expected, s)
	}
}

func TestRead(t *testing.T) {
	o := New()
	err := o.ReadFile("testdata/example1.opml")
	if err != nil {
		t.Errorf(`ReadFile should return an OPML structure an nil error, %s`, err)
		t.FailNow()
	}
	if o.Version != "1.0" {
		t.Errorf(`Expected version 1.0, got %s`, o.Version)
	}
	if o.Head.Title != "johndoe@example.com subscriptions in Go Read" {
		t.Errorf(`Expected "johndoe@example.com subscriptions in Go Read", found -> %s`, o.Head.Title)
	}
	i := 64
	if len(o.Body.Outline) != i {
		t.Errorf(`expected %d outline elements, found %d`, i, len(o.Body.Outline))
	}
	s := o.String()
	if strings.Contains(s, `<outline text=""></outline>`) == true {
		t.Errorf("an empty outline is included in string: %s", s)
	}
	//fmt.Printf("o: %s\n", o)
	o = New()
	err = o.ReadFile("testdata/example2.opml")
	if err != nil {
		t.Errorf(`ReadFile should return an OPML structure and a nil error, %s`, err)
		t.FailNow()
	}
	if o.Version != "2.0" {
		t.Errorf(`Expected version 2.0, got %s`, o.Version)
	}

}

func TestWrite(t *testing.T) {
	fname := path.Join("testdata", "test1.opml")
	o := New()
	err := o.WriteFile(fname, 0664)
	if err != nil {
		t.Errorf("%s", err)
	}

	s := []byte(o.String())
	src, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Errorf("%s", err)
	}
	if bytes.Equal(s, src) != true {
		t.Errorf(`%s != %s`, s, src)
	}
	// cleanup the temp file
	os.Remove(fname)
}

func TestSort(t *testing.T) {
	// Simple unnested list
	fname := path.Join("testdata", "example3.opml")
	o := New()
	err := o.ReadFile(fname)
	if err != nil {
		t.Errorf("can't read %s, %s", fname, err)
		t.FailNow()
	}
	sort.Sort(ByText(o.Body.Outline))
	expected := `<opml version="2.0"><head><title>Unsorted list</title><dateCreated>Mon, 23 May 2016 08:33:00 GMT</dateCreated></head><body><outline text="Alexandrina"></outline><outline text="Tomasa"></outline><outline text="Zachary"></outline></body></opml>`
	result := o.String()
	if strings.Compare(expected, result) != 0 {
		t.Errorf("\n%s\n!=\n%s\n", expected, result)
	}

	// List with two levels of nesting
	fname = path.Join("testdata", "example4.opml")
	o = New()
	err = o.ReadFile(fname)
	if err != nil {
		t.Errorf("can't read %s, %s", fname, err)
		t.FailNow()
	}
	o.Sort()
	expected = `<opml version="2.0"><head><title>test sort</title></head><body><outline text="Places of interest"><outline text="Bay Area"><outline text="Los Gatos"></outline><outline text="Mountain View"></outline><outline text="Palo Alto"></outline><outline text="Woodside"></outline></outline><outline text="Boston"><outline text="Cambridge"></outline><outline text="West Newton"></outline></outline><outline text="New Orleans"><outline text="Congo Square"></outline><outline text="Metairie"></outline><outline text="Uptown"></outline></outline><outline text="New York"><outline text="Midtown"></outline><outline text="Upper Eastside"></outline></outline><outline text="Victoria, BC"></outline></outline></body></opml>`
	result = o.String()
	if strings.Compare(expected, result) != 0 {
		t.Errorf("\n%s\n!=\n%s\n", expected, result)
	}

}

func TestAppend(t *testing.T) {
	o1, err := ReadFile("testdata/simple1.opml")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o2, err := ReadFile("testdata/simple2.opml")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o3 := New()
	err = o3.Append(o1, o2)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	i := len(o1.Body.Outline) + len(o2.Body.Outline)
	if len(o3.Body.Outline) != i {
		t.Errorf("Count wrong o1 %d, o2 %d, o3 %d", len(o1.Body.Outline), len(o2.Body.Outline), len(o3.Body.Outline))
	}
}

/*
func TestFilterForTypes(t *testing.T) {
	t.Errorf("Filter for types not implemented")
	t.FailNow()
}
*/

const walkSrc = `<?xml version="1.0"?>
<opml version="2.0"><head><title>t</title></head><body>
<outline text="a"><outline text="a1"/><outline text="a2"><outline text="a2x"/></outline><outline text="a3"/></outline>
<outline text="b"><outline text="b1"/></outline>
</body></opml>`

// TestWalkVisitsEveryOutline: Walk must visit siblings and nested children,
// depth first, not only the first child at each level.
func TestWalkVisitsEveryOutline(t *testing.T) {
	o := &OPML{}
	if err := Unmarshal([]byte(walkSrc), o); err != nil {
		t.Fatal(err)
	}
	var got []string
	if err := o.Walk(func(ol *Outline) bool {
		got = append(got, ol.Text)
		return true
	}); err != nil {
		t.Fatal(err)
	}
	want := "a a1 a2 a2x a3 b b1"
	if s := strings.Join(got, " "); s != want {
		t.Errorf("visited %q, want %q", s, want)
	}
}

// TestWalkStopsWhenFuncReturnsFalse: false from the callback ends the whole
// walk, including from a nested outline.
func TestWalkStopsWhenFuncReturnsFalse(t *testing.T) {
	o := &OPML{}
	if err := Unmarshal([]byte(walkSrc), o); err != nil {
		t.Fatal(err)
	}
	var got []string
	_ = o.Walk(func(ol *Outline) bool {
		got = append(got, ol.Text)
		return ol.Text != "a2"
	})
	if s := strings.Join(got, " "); s != "a a1 a2" {
		t.Errorf("visited %q, want %q", s, "a a1 a2")
	}
}

// TestLanguageAttribute: the OPML language attribute is read and written.
func TestLanguageAttribute(t *testing.T) {
	src := `<opml version="2.0"><head/><body><outline text="x" language="en"/></body></opml>`
	o := &OPML{}
	if err := Unmarshal([]byte(src), o); err != nil {
		t.Fatal(err)
	}
	if got := o.Body.Outline[0].Language; got != "en" {
		t.Errorf("Language = %q, want %q", got, "en")
	}
	out, err := Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `language="en"`) {
		t.Errorf("marshalled output lacks language attribute: %s", out)
	}
}

func TestMarshalIndent(t *testing.T) {
	o, err := Parse([]byte(tocSrc))
	if err != nil {
		t.Fatal(err)
	}
	out, err := MarshalIndent(o)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.HasPrefix(s, `<?xml version="1.0" encoding="UTF-8"?>`+"\n") {
		t.Errorf("no XML declaration first: %.80q", s)
	}
	if !strings.HasSuffix(s, "\n") {
		t.Errorf("no trailing newline")
	}
	if !strings.Contains(s, "\n  <body>\n    <outline") {
		t.Errorf("not indented two spaces per level: %s", s)
	}
	back, err := Parse(out)
	if err != nil {
		t.Fatalf("indented output does not parse: %v", err)
	}
	if got, want := len(back.Flatten()), len(o.Flatten()); got != want {
		t.Errorf("round trip: %d pages, want %d", got, want)
	}
	if got, want := back.Includes(), o.Includes(); !reflect.DeepEqual(got, want) {
		t.Errorf("round trip includes %v, want %v", got, want)
	}
}

func TestMarshalAndStringStayCompact(t *testing.T) {
	o, _ := Parse([]byte(tocSrc))
	out, err := Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "<?xml") || strings.Contains(string(out), "\n") {
		t.Errorf("Marshal output changed: %.100q", out)
	}
	if o.String() != string(out) {
		t.Error("String and Marshal disagree")
	}
}

func TestWriteFileReturnsTheError(t *testing.T) {
	o := New()
	if err := o.WriteFile(path.Join(t.TempDir(), "no", "such", "dir", "x.opml"), 0o664); err == nil {
		t.Error("WriteFile into a missing directory returned nil")
	}
	if err := o.WriteFileIndent(path.Join(t.TempDir(), "no", "dir", "x.opml"), 0o664); err == nil {
		t.Error("WriteFileIndent into a missing directory returned nil")
	}
}

func TestWriteFileWithNilBody(t *testing.T) {
	o := &OPML{Version: "2.0"}
	fname := path.Join(t.TempDir(), "x.opml")
	if err := o.WriteFile(fname, 0o664); err != nil {
		t.Errorf("WriteFile with nil Body: %v", err)
	}
}

func TestWriteFileIndent(t *testing.T) {
	o, _ := Parse([]byte(tocSrc))
	fname := path.Join(t.TempDir(), "toc.opml")
	if err := o.WriteFileIndent(fname, 0o664); err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := MarshalIndent(o)
	if !bytes.Equal(src, want) {
		t.Errorf("file differs from MarshalIndent:\n%s", src)
	}
}

func TestOwnerIDUsesTheSpecName(t *testing.T) {
	o, err := ReadFile(path.Join("testdata", "example2.opml"))
	if err != nil {
		t.Fatal(err)
	}
	const want = "http://www.example.org/profiles/johndoe"
	if o.Head.OwnerID != want {
		t.Fatalf("OwnerID = %q, want %q (<ownerId> in the file was dropped)", o.Head.OwnerID, want)
	}
	out, err := Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "<ownerId>"+want+"</ownerId>") {
		t.Errorf("Marshal does not write <ownerId>: %s", out)
	}
	if strings.Contains(string(out), "OwnerId") {
		t.Errorf("Marshal still writes OwnerId: %s", out)
	}
	j, err := json.Marshal(o.Head)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(j), `"ownerId":"`+want+`"`) {
		t.Errorf("JSON key is not ownerId: %s", j)
	}
}

func TestModifiedJSONKeyIsDateModified(t *testing.T) {
	src := `<opml version="2.0"><head><dateCreated>Fri, 20 Nov 2026 14:06:24 GMT</dateCreated><dateModified>Sat, 21 Nov 2026 09:00:00 GMT</dateModified></head><body><outline text="x"/></body></opml>`
	o, err := Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	j, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(j), `"dateModified":"Sat, 21 Nov 2026 09:00:00 GMT"`) {
		t.Errorf("JSON lacks dateModified: %s", j)
	}
	if strings.Contains(string(j), "dataModified") {
		t.Errorf("JSON still has the dataModified typo: %s", j)
	}
}
