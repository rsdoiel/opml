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
	if s != `<opml version="2.0"><head></head><body><outline text=""></outline></body></opml>` {
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

func TestFilterForTypes(t *testing.T) {
	t.Errorf("Filter for types not implemented")
	t.FailNow()
}

func TestWalk(t *testing.T) {
	txt := `<?xml version="1.0" encoding="ISO-8859-1"?>
<opml version="2.0">
    <head> 
	    <title>scriptingNewsDirectory.opml</title> 
		<dateCreated>Thu, 13 Oct 2005 15:34:07 GMT</dateCreated> 
		<dateModified>Tue, 25 Oct 2005 21:33:57 GMT</dateModified> 
		<ownerName>Dave Winer</ownerName> 
		<ownerEmail>dwiner@yahoo.com</ownerEmail> 
		<expansionState></expansionState> 
		<vertScrollState>1</vertScrollState> 
		<windowTop>105</windowTop> 
		<windowLeft>466</windowLeft> 
		<windowBottom>386</windowBottom> 
		<windowRight>964</windowRight> 
	</head> 
	<body> 
	    <outline text="Scripting News sites" created="Sun, 16 Oct 2005 05:56:10 GMT" type="link" url="http://hosting.opml.org/dave/mySites.opml"/> 
		<outline text="News.Com top 100 OPML" created="Tue, 25 Oct 2005 21:33:28 GMT" type="link" url="http://news.com.com/html/ne/blogs/CNETNewsBlog100.opml"/> 
		<outline text="BloggerCon III Blogroll" created="Mon, 24 Oct 2005 05:23:52 GMT" type="link" url="http://static.bloggercon.org/iii/blogroll.opml"/> 
		<outline text="TechCrunch reviews" type="link" url="http://hosting.opml.org/techcrunch.opml.org/TechCrunch.opml"/> 
		<outline text="Tod Maffin's directory of Public Radio podcasts" type="link" url="http://todmaffin.com/radio.opml"/> 
		<outline text="Adam Curry's iPodder.org director type="link" url="http://homepage.mac.com/dailysourcecode/DSC/ipodderDirectory.opml"/> 
		<outline text="Memeorandum" created="Thu, 13 Oct 2005 15:19:05 GMT" type="link" url="http://tech.memeorandum.com/index.opml"/> 
		<outline text="DaveNet archive" created="Wed, 12 Oct 2005 01:39:56 GMT" type="link" url="http://davenet.opml.org/index.opml"/> 
	</body> 
</opml>
`
	t.Errorf("Walk not implemented")
	t.FailNow()
}
