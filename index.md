[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

# opml

A package to manipulate OPML files in Go. It is based on the spec as described at http://dev.opml.org/spec2.html.
The `opml` package imports only the standard library. It reads and writes OPML (UTF-8, ISO-8859-1, US-ASCII and
windows-1252 input), walks an outline, and can list the pages of a table of contents with `Flatten()` and `Includes()`. The
`cmark` package converts an outline to a CommonMark list and back.

It includes command line programs that demonstrate it (each has a manual page; the [user manual](user-manual.md) lists them):

- [opmlsort](opmlsort.1.md): sorts an OPML outline by text attribute values
- [opmlcat](opmlcat.1.md): concatenates one or more OPML outlines
- [opml2json](opml2json.1.md): converts an OPML file into a JSON document
- [opml2urls](opml2urls.1.md): lists the feed URLs of an OPML file, or the pages of a table of contents, or writes newsboat's url file
- [urls2opml](urls2opml.1.md): converts a list of urls into OPML XML
- [opml2cmark](opml2cmark.1.md): writes an OPML outline as a CommonMark list, to edit in a text editor
- [cmark2opml](cmark2opml.1.md): reads that Markdown list back into OPML

See [INSTALL](INSTALL.md) to install them and [CHANGES](CHANGES.md) for what changed in each release.

## Markdown web

This documentation is also published as Markdown, for [BonBon](https://github.com/scripting/bonBon/) and
[cmarkbrowser](https://github.com/rsdoiel/cmarkbrowser). The site's address returns the HTML pages, because the host does
not look at what a browser asks for, so open the Markdown home page directly:
<https://rsdoiel.github.io/opml/index.md>. Its table of contents is `/.well-known/toc.opml` in the same folder.

## Copying

The opml package is copyright (c) 2021 by R. S. Doiel and [released](https://github.com/rsdoiel/opml/releases/latest) under the GNU Affero General Public License, version 3 or later (see LICENSE).

## References

+ [OPML Spec 2](http://dev.opml.org/spec2.html)


