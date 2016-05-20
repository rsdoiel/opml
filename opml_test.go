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
	"testing"
)

func TestBasic(t *testing.T) {
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
