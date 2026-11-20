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
	"errors"
	"io"
	"os"
	"path"
	"strings"
	"testing"
)

// legacyEncodingSrc is an OPML document in a single-byte encoding. 0xe9 is
// é in ISO-8859-1 and windows-1252; 0x93 and 0x94 are curly quotes in
// windows-1252 only.
func legacyEncodingSrc(enc string, text []byte) []byte {
	src := []byte(`<?xml version="1.0" encoding="` + enc + `"?><opml version="2.0"><head/><body><outline text="`)
	src = append(src, text...)
	return append(src, []byte(`"/></body></opml>`)...)
}

func TestLegacyEncodings(t *testing.T) {
	tests := []struct {
		enc  string
		in   []byte
		want string
	}{
		{"ISO-8859-1", []byte("caf\xe9"), "café"},
		{"iso-8859-1", []byte("caf\xe9"), "café"},
		{"US-ASCII", []byte("plain"), "plain"},
		{"windows-1252", []byte("\x93caf\xe9\x94 \x80"), "“café” €"},
		{"UTF-8", []byte("café"), "café"},
	}
	dir := t.TempDir()
	for _, tc := range tests {
		src := legacyEncodingSrc(tc.enc, tc.in)
		check := func(how string, o *OPML, err error) {
			t.Helper()
			if err != nil {
				t.Errorf("%s %s: %v", tc.enc, how, err)
				return
			}
			if got := o.Body.Outline[0].Text; got != tc.want {
				t.Errorf("%s %s: Text = %q, want %q", tc.enc, how, got, tc.want)
			}
		}
		o, err := Parse(src)
		check("Parse", o, err)
		u := &OPML{}
		err = Unmarshal(src, u)
		check("Unmarshal", u, err)
		fname := path.Join(dir, "legacy.opml")
		if err := os.WriteFile(fname, src, 0o644); err != nil {
			t.Fatal(err)
		}
		r, err := ReadFile(fname)
		check("ReadFile", r, err)
		m := New()
		err = m.ReadFile(fname)
		check("(*OPML).ReadFile", m, err)
	}
}

func TestUnknownEncodingIsAnError(t *testing.T) {
	if _, err := Parse(legacyEncodingSrc("EBCDIC-XYZ", []byte("x"))); err == nil {
		t.Error("expected an error for an unsupported encoding")
	}
}

func TestCharsetReaderLabels(t *testing.T) {
	for _, label := range []string{"ISO-8859-1", "latin1", "Latin-1", "US-ASCII", "ascii", " windows-1252 ", "CP1252"} {
		r, err := charsetReader(label, strings.NewReader("x"))
		if err != nil || r == nil {
			t.Errorf("label %q: reader %v, err %v", label, r, err)
		}
	}
	if _, err := charsetReader("shift_jis", strings.NewReader("x")); err == nil {
		t.Error("shift_jis should be unsupported")
	}
}

func TestWindows1252UndefinedBytesKeepTheirValue(t *testing.T) {
	r, _ := charsetReader("windows-1252", strings.NewReader("\x81\x8d\x8f\x90\x9d"))
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if want := "\u0081\u008d\u008f\u0090\u009d"; string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("disk gone") }

func TestCharsetReaderReportsReadErrors(t *testing.T) {
	r, err := charsetReader("ISO-8859-1", failingReader{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(r); err == nil || !strings.Contains(err.Error(), "disk gone") {
		t.Errorf("err = %v, want the read error", err)
	}
	if err := Unmarshal(legacyEncodingSrc("ISO-8859-1", []byte("x")), &OPML{}); err != nil {
		t.Errorf("a good document should still parse: %v", err)
	}
}
