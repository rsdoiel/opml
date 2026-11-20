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

import "strings"

// isInclude reports whether an outline is an include: its url names
// another outline document, not a page.
func (ol *Outline) isInclude() bool {
	return strings.EqualFold(ol.Type, "include")
}

/** Flatten returns the outlines of a table of contents that name a page,
 * in document order: every outline, at any depth, that has a url and is
 * not an include. A parent that has a url comes before its children. The
 * outlines are the document's own, not copies.
 *
 * Returns:
 *   []*Outline — the page outlines, empty when there are none
 *
 * Example:
 *   o, _ := opml.ReadFile(".well-known/toc.opml")
 *   for _, page := range o.Flatten() {
 *       fmt.Println(page.Text, page.URL)
 *   }
 */
func (o *OPML) Flatten() []*Outline {
	var pages []*Outline
	if o == nil || o.Body == nil {
		return pages
	}
	for _, ol := range o.Body.Outline {
		walk(ol, func(ol *Outline) bool {
			if ol.URL != "" && !ol.isInclude() {
				pages = append(pages, ol)
			}
			return true
		})
	}
	return pages
}

/** Includes returns the addresses of the include outlines of a table of
 * contents, in document order. An include is an outline with
 * type="include"; its url is another outline document to merge in. An
 * include with no url is ignored.
 *
 * Returns:
 *   []string — the include addresses, empty when there are none
 *
 * Example:
 *   for _, addr := range o.Includes() {
 *       // fetch addr, opml.Parse it, and merge it in
 *   }
 */
func (o *OPML) Includes() []string {
	var addrs []string
	if o == nil || o.Body == nil {
		return addrs
	}
	for _, ol := range o.Body.Outline {
		walk(ol, func(ol *Outline) bool {
			if ol.URL != "" && ol.isInclude() {
				addrs = append(addrs, ol.URL)
			}
			return true
		})
	}
	return addrs
}
