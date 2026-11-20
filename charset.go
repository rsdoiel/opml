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
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// windows1252 maps bytes 0x80-0x9F of windows-1252 to runes. Bytes the
// encoding leaves undefined map to the C1 control of the same value, as
// browsers do.
var windows1252 = [32]rune{
	0x20AC, 0x81, 0x201A, 0x0192, 0x201E, 0x2026, 0x2020, 0x2021,
	0x02C6, 0x2030, 0x0160, 0x2039, 0x0152, 0x8D, 0x017D, 0x8F,
	0x90, 0x2018, 0x2019, 0x201C, 0x201D, 0x2022, 0x2013, 0x2014,
	0x02DC, 0x2122, 0x0161, 0x203A, 0x0153, 0x9D, 0x017E, 0x0178,
}

/** charsetReader is an xml.Decoder CharsetReader for the single-byte
 * encodings found in the wild in OPML files.
 *
 * Parameters:
 *   label (string) — the encoding named in the XML declaration
 *   input (io.Reader) — the document bytes
 *
 * Returns:
 *   io.Reader — UTF-8 bytes
 *   error — when the encoding is not supported
 *
 * Example:
 *   d := xml.NewDecoder(r)
 *   d.CharsetReader = charsetReader
 */
func charsetReader(label string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(strings.TrimSpace(label)) {
	case "iso-8859-1", "iso8859-1", "latin1", "latin-1", "us-ascii", "ascii":
		return recode(input, false), nil
	case "windows-1252", "cp1252":
		return recode(input, true), nil
	}
	return nil, fmt.Errorf("unsupported XML encoding %q", label)
}

// recode converts single-byte text to UTF-8.
func recode(input io.Reader, cp1252 bool) io.Reader {
	src, err := io.ReadAll(input)
	if err != nil {
		return &errReader{err}
	}
	var buf bytes.Buffer
	buf.Grow(len(src) + len(src)/8)
	for _, b := range src {
		r := rune(b)
		if cp1252 && b >= 0x80 && b <= 0x9F {
			r = windows1252[b-0x80]
		}
		buf.WriteRune(r)
	}
	return &buf
}

type errReader struct{ err error }

func (e *errReader) Read([]byte) (int, error) { return 0, e.err }

// unmarshal decodes an OPML document in any supported encoding.
func unmarshal(src []byte, o *OPML) error {
	d := xml.NewDecoder(bytes.NewReader(src))
	d.CharsetReader = charsetReader
	return d.Decode(o)
}
