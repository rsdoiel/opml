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
	err := o.WriteFile(fname)
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
}
