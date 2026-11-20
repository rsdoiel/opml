[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

# opml

A package to manipulate OPML files in Go. It is based on the spec as described at http://dev.opml.org/spec2.html.
The `opml` package imports only the standard library. It reads and writes OPML (UTF-8, ISO-8859-1, US-ASCII and
windows-1252 input), walks an outline, and can list the pages of a table of contents with `Flatten()` and `Includes()`. The
`cmark` package converts an outline to a CommonMark list and back.

It includes command line programs that demonstrate it:

_opmlsort_
: sorts an OPML outline by text attribute values

_opmlcat_
: Concatenates one or more OPML outlines

_opml2json_
: Converts an OPML file into a JSON document

_opml2urls_
: Lists the feed URLs of an OPML file, or the pages of a table of contents, or writes newsboat's url file

_urls2opml_
: Converts a list of urls into OPML XML

_opml2cmark_
: Writes an OPML outline as a CommonMark list, to edit in a text editor

_cmark2opml_
: Reads that Markdown list back into OPML


## Copying

The opml package is copyright (c) 2021 by R. S. Doiel and [released](https://github.com/rsdoiel/opml/releases/latest) under the GNU Affero General Public License, version 3 or later (see LICENSE).

## References

+ [OPML Spec 2](http://dev.opml.org/spec2.html)


