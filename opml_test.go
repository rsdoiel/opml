//
// Package opml provides basic utility functions for working with
// OPML files.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2016 all rights reserved.
// Released under the BSD 3-Clause License
// See: https://opensource.org/licenses/BSD-3-Clause
//
package opml

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
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

	s := o.String()
	if s != `<opml version="2.0"><head></head><body></body></opml>` {
		t.Errorf("Expected a minimal OPML document [%s]", s)
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

/*
func TestToMarkdown() {
	o := New()

	source := `
<opml version="2.0">
<head>
<title>n4</title>
<ownerProfile>http://rsdoiel.github.io</ownerProfile>
<ownerName>R. S. Doiel</ownerName>
<ownerEmail>rsdoiel@gmail.com</ownerEmail>
<dateModified>Sat, 28 May 2016 15:50:19 GMT</dateModified>
<expansionState></expansionState>
<lastCursor>2</lastCursor>
</head>
<body>
<outline text="&lt;a href=&quot;http://jane-doe.example.org&quot;&gt;Jane Doe&lt;/a&gt;"/>
<outline text="&lt;a href=&quot;http://example.org/jordan-doe&quot;&gt;Jordan Doe&lt;/a&gt;" created="Sat, 28 May 2016 15:48:19 GMT"/>
<outline text="Spancer Wallace" created="Sat, 28 May 2016 15:50:13 GMT"/>
</body>
</opml>
`
	expected := `
# Test List Names

+ [Jane Doe](http://jane-doe.example.org)
+ [Jordan Doe](http://example.org/jordan-doe)
+ Spancer Wallace
`

}
*/

func TestParse(t *testing.T) {
	sources := []string{
		`<opml version="2.0">
<head>
	<title>no-children</title>
</head>
<body>
	<outline text="root zero"/>
	<outline text="root one"/>
	<outline text="root two"/>
	<outline text="root three"/>
</body>
</opml>
`,
		`<opml version="2.0">
<head>
	<title>children</title>
</head>
<body>
	<outline text="root zero">
		<outline text="child zero, root zero"/>
		<outline text="child one, root zero"/>
		<outline text="child two, root zero"/>
	</outline>
	<outline text="root one">
		<outline text="child zero, root one"/>
		<outline text="child one, root one"/>
		<outline text="child two, root one"/>
	</outline>
	<outline text="root two">
		<outline text="child zero, root two"/>
		<outline text="child one, root two"/>
		<outline text="child two, root two"/>
	</outline>
	<outline text="root three">
		<outline text="child zero, root three"/>
		<outline text="child one, root three"/>
		<outline text="child two, root three"/>
	</outline>
</body>
</opml>
`,
		`<opml version="2.0">
<head>
	<title>some grand children</title>
</head>
<body>
	<outline text="root zero">
		<outline text="child zero, root zero">
			<outline text="granchild zero, child zero, root zero"/>
		</outline>
		<outline text="child one, root zero"/>
		<outline text="child two, root zero"/>
	</outline>
	<outline text="root one">
		<outline text="child zero, root one"/>
		<outline text="child one, root one">
			<outline text="granchild zero, child one, root one"/>
			<outline text="granchild one, child one, root one"/>
			<outline text="granchild two, child one, root one"/>
		</outline>
		<outline text="child two, root one"/>
	</outline>
	<outline text="root two">
		<outline text="child zero, root two"/>
		<outline text="child one, root two"/>
		<outline text="child two, root two"/>
	</outline>
	<outline text="root three">
		<outline text="child zero, root three"/>
		<outline text="child one, root three"/>
		<outline text="child two, root three">
			<outline text="granchild zero, child two, root three"/>
			<outline text="granchild one, child two, root three"/>
			<outline text="granchild two, child two, root three"/>
			<outline text="granchild three, child two, root three"/>
			<outline text="granchild four, child two, root three"/>
		</outline>
	</outline>
</body>
</opml>
`,
	}

	for i, src := range sources {
		//fmt.Printf("DEBUG Parsing %d, %s\n", i, src)
		_, err := Parse([]byte(src))
		if err != nil {
			t.Errorf("%d %s", i, err)
			t.FailNow()
		}
	}
}

func TestSelect(t *testing.T) {
	sources := []string{
		`<opml version="2.0">
<head>
	<title>no-children</title>
</head>
<body>
	<outline text="root zero"/>
	<outline text="root one"/>
	<outline text="root two"/>
	<outline text="root three"/>
</body>
</opml>
`,
		`<opml version="2.0">
<head>
	<title>children</title>
</head>
<body>
	<outline text="root zero">
		<outline text="child zero, root zero"/>
		<outline text="child one, root zero"/>
		<outline text="child two, root zero"/>
	</outline>
	<outline text="root one">
		<outline text="child zero, root one"/>
		<outline text="child one, root one"/>
		<outline text="child two, root one"/>
	</outline>
	<outline text="root two">
		<outline text="child zero, root two"/>
		<outline text="child one, root two"/>
		<outline text="child two, root two"/>
	</outline>
	<outline text="root three">
		<outline text="child zero, root three"/>
		<outline text="child one, root three"/>
		<outline text="child two, root three"/>
	</outline>
</body>
</opml>
`,
		`<opml version="2.0">
<head>
	<title>some grand children</title>
</head>
<body>
	<outline text="root zero">
		<outline text="child zero, root zero">
			<outline text="granchild zero, child zero, root zero"/>
		</outline>
		<outline text="child one, root zero"/>
		<outline text="child two, root zero"/>
	</outline>
	<outline text="root one">
		<outline text="child zero, root one"/>
		<outline text="child one, root one">
			<outline text="granchild zero, child one, root one"/>
			<outline text="granchild one, child one, root one"/>
			<outline text="granchild two, child one, root one"/>
		</outline>
		<outline text="child two, root one"/>
	</outline>
	<outline text="root two">
		<outline text="child zero, root two"/>
		<outline text="child one, root two"/>
		<outline text="child two, root two"/>
	</outline>
	<outline text="root three">
		<outline text="child zero, root three"/>
		<outline text="child one, root three"/>
		<outline text="child two, root three">
			<outline text="granchild zero, child two, root three"/>
			<outline text="granchild one, child two, root three"/>
			<outline text="granchild two, child two, root three"/>
			<outline text="granchild three, child two, root three"/>
			<outline text="granchild four, child two, root three"/>
		</outline>
	</outline>
</body>
</opml>
`,
	}

	testVals := []string{}

	expected := []string{}

	for i, src := range sources {
	}
}
