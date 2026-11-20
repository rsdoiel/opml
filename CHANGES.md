# Changes

Reconstructed for v0.0.10 and v0.0.11 from each tag's `codemeta.json` release notes (v0.0.3 to v0.0.9 have none; see
the git history); maintained going forward.

## v0.0.12 — 2026-11-20

### Changed

- **`opml2json` writes `ownerId` and `dateModified`.** The head's JSON key `OwnerId` is now `ownerId`, and the typo
  `dataModified` is now `dateModified`. **A script that reads the old keys must read the new ones.** In the XML, `Head.OwnerID`
  is now read and written as `<ownerId>`, which is what the OPML spec says; before, a spec-correct `<ownerId>` was silently
  dropped on read (`testdata/example2.opml` lost its owner on a round trip). A file an earlier version wrote with
  `<OwnerId>` is no longer read.
- **`opml2urls` lists the pages of a table of contents.** `-url` writes the `url` attribute, and it is on by itself when a
  document has no `xmlUrl` or `htmlUrl`. An outline with `type="include"` is left out unless `-includes` is given. Outlines
  are walked at any depth in document order.
- **`opml2urls` exits 0 for a document with nothing to list** and prints nothing, where it exited 70 with "outline is
  empty" (workspace DR-0003: an empty listing is 0).
- **`WriteFile` returns the marshal error**, uses `os` in place of `ioutil`, and no longer panics on an OPML with a nil
  `Body`. `Marshal`, `String` and `WriteFile` keep the compact form, byte for byte.
- The comments say RFC 822 where they said 882.

### Added

- **`opml2cmark` and `cmark2opml`** write an OPML outline as a CommonMark list and read it back, so an outline can be
  edited in any text editor (opml DR-0001). The head is YAML front matter, each outline is a list item, a link carries the
  address (`xmlUrl` for type rss or atom, `url` otherwise), and other attributes follow in a `{key="value"}` block.
  `cmark2opml` accepts any bullet, indentation, blank lines and wrapped items, and stops with an error that names the line
  for anything that is not part of the list. OPML to Markdown to OPML gives the same model for every file in `testdata` and
  for Dave Winer's table of contents.
- **The `cmark` package** (`FromOPML`, `ToOPML`, `ErrUnrepresentable`, `ErrNotOutline`). The `opml` package itself still
  imports only the standard library; `cmark` is the one that needs goldmark and yaml.v3.
- **Encodings.** `Parse`, `ReadFile` and `Unmarshal` read ISO-8859-1, US-ASCII and windows-1252 documents (standard library
  only). Dave Winer's `toc.opml`, which declares ISO-8859-1, could not be read before.
- **`Flatten()` and `Includes()`** for reading a table of contents: the page outlines in document order, and the addresses of
  the include outlines.
- **`MarshalIndent` and `WriteFileIndent`** write the XML declaration, two spaces per level and a final newline.
- **`Head.Generator`** reads and writes the `<generator>` element, so a head like Dave Winer's round trips.

### Removed

- **The TypeScript port** (`opml.ts`, `opml_test.ts`, `deno.json`, `deno.lock`). It had drifted from the Go package and is
  not being developed.

### Requirements

- `github.com/yuin/goldmark` v1.8.7 and `gopkg.in/yaml.v3` v3.0.1, imported by the `cmark` package and the two new programs
  only.

## v0.0.11 — 2026-11-18

### Fixed

- `Walk()` visited only the first child at each level and did not stop when the callback returned false at depth.
- The `Outline.Language` tag typo, so the `language` attribute is read and written.
- `opmlcat` and `opml2urls` no longer empty an input file named by mistake; `opmlsort` exits non-zero on a missing input.
- `-help`, `-license` and `-version` output in `opml2json` and `opmlsort`.

### Changed

- The programs moved into `internal/cli` with documented exit codes (2, 65, 66, 73, 74, 77) and an EXIT STATUS section in
  each help page.
- Added `-h`, `-l` and `-v`; new `-nl`, `-input` and `-output` options.
- Dropped the unused `caltechlibrary/cli` requirement.

## v0.0.10 — 2025-10-01

### Changed

- Normalized the license to AGPL >= 3.
- Updated the codemeta status to inactive (it is stable).
