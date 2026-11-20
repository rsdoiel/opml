---
title: "OPML, a Go package and cli"
---

## Next

+ [ ] `cmark2opml` and `opml2cmark` are built (opml DR-0001, accepted; 2026-11-20). What is left:
    + add a line to DR-0001 on the OPML root `version` attribute (written as `opmlVersion` in the front matter when it is not 2.0; an absent version is written as an empty one), on empty items (an empty first child needs a blank line before it, as CommonMark does not let an empty item interrupt a paragraph), and on goldmark's `Link.Destination` (raw source, so the reader unescapes it)
    + a `-lenient` option that warns and skips foreign Markdown instead of exiting 65 (DR-0001, decision 7, left it for later)
    + text that has a newline, tab or leading or trailing white space, and reference-style links, are still errors
    + head attributes on `<head>` itself and unmodelled head children are not carried
    + `go test -fuzz FuzzRoundTrip ./cmark` ran 514,000 executions with no failure; run it again after any change to the mapping

## From the antennaApp table-of-contents work (2026-11-19)

Found by running the module against Dave Winer's published `toc.opml` files while writing antennaApp DR-0013 (see also the `opml` project in the knowledge base). `Walk` and the `langauge` tag were found on the released v0.0.10 and are already fixed in v0.0.11 (see Completed); these are what is still open in the current source.

+ [ ] After the v0.0.12 release, re-install `opml2urls`: the installed `~/bin/opml2urls` is 0.0.10 (2026-11-16), from before the fix that makes it honour a file argument, so it ignores the file name and waits on standard input
+ [ ] When the above are done, release v0.0.12; antennaApp pins v0.0.10 and should move to v0.0.11 now for the `Walk` fix

## Someday, maybe

+ opmlviewer - a cli/terminal based opml viewer
+ omplfilter
    + might function something like Unix find command
+ ompledit - CRUD operations to individual Outline elements and element lists
+ Review River5 by Dave Winer for insights into what additional functions are needed in package 

## Dropped

+ [x] ~~Support for Frontier's fttb (fat pages) into OPML~~ Dropped 2026-11-20: the `#pageData` payload is Frontier's undocumented internal binary format, there is no Frontier installation to test against, and the 36 sample objects on scripting.com/fatpages are scripts and tables, not outlines. Only the envelope (an HTML comment of `#name value` lines with a Base64 `#pageData`) is specified
+ [x] ~~A tool that reads an OPML file, harvests the current feeds, indexes them for browsing and sends interesting articles to Pocket~~ Dropped 2026-11-20: Pocket no longer exists, and feeds are aggregated from a Markdown list with antennaApp, not from OPML files

## Completed

+ [x] Removed the TypeScript port (`opml.ts`, `opml_test.ts`, `deno.json`, `deno.lock`): it had drifted from the Go package (`ownerID`, no `generator`, `Flatten`, `Includes`) and is not being developed. Release note for v0.0.12 (2026-11-20)
+ [x] Convert an OPML file for the newsboat TUI feed reader: `opml2urls -newsboat` writes newsboat's url file (`URL "~Text"`). The old TODO line said the option was added to `urls2opml`; it was `opml2urls`, and had been done earlier without being checked off (2026-11-20)
+ [x] `Head.Generator` (`<generator>`, JSON `generator`) is read and written, so Dave Winer's `toc.opml` head round-trips; unset writes nothing. Other unknown child elements of `head` are still dropped, by choice (2026-11-20)
+ [x] `opml2urls` follows the exit-code rules for an empty listing: a document with nothing to list exits 0 and prints nothing, where it exited 70 with "outline is empty" (2026-11-20)
+ [x] `opml2urls` suits an outline of pages: `-url` lists the url attribute, and it is on by itself when a document has no xmlUrl or htmlUrl; `type="include"` outlines are left out unless `-includes`; outlines are walked at any depth in document order. Man page updated (2026-11-20)
+ [x] Comments say RFC 822, not 882, in `opml.go` and `opml.ts` (2026-11-20)
+ [x] `Head.Modified` JSON key is `dateModified` (was the typo `dataModified`). Release note: `opml2json` output changes; anything reading `dataModified` must read `dateModified` (2026-11-20)
+ [x] `Head.OwnerID` reads and writes `ownerId` (XML and JSON), as the spec says; before, a spec-correct `<ownerId>` was silently dropped on read. Release note: `opml2json` now writes `ownerId`, and a file written by an earlier version with `<OwnerId>` is no longer read (2026-11-20)
+ [x] `MarshalIndent` and `(*OPML).WriteFileIndent` write a declaration, two-space indent and a final newline; `Marshal`, `String` and `WriteFile` keep the compact form; `WriteFile` returns the marshal error, uses `os`, and no longer panics on a nil `Body` (2026-11-20)
+ [x] `(*OPML).Flatten()` and `(*OPML).Includes()` for reading a table of contents (toc.go); on Dave Winer's `toc.opml`: 3,553 page entries, no includes (2026-11-20)
+ [x] Read ISO-8859-1, US-ASCII and windows-1252 documents in `Parse`, `ReadFile` and `Unmarshal` (stdlib `charset.go`, no new dependency); Dave Winer's `toc.opml` now parses (2026-11-20)
+ [x] Fix `Walk` visiting only the first child at each level, and stopping early from a nested outline (2026-11-18)
+ [x] Fix the `Outline.Language` struct tag typo (`langauge`), so the `language` attribute is read and written (2026-11-18)
+ [x] Add opml2json 
+ [x] Add support for custom attributes
+ [x] Add Bash script to fetch Dave Winer's userland samples at http://scripting.com/misc/userlandSamples.zip

## Modernization (2026-11-16)

To match the workspace standard (knowledge and opml are the models; Go 1.27.3 across repositories). Do one repository at a time; run the tests and `make -n release` after each step.

- [x] Replace `github.com/caltechlibrary/cli`: the programs already used the standard library `flag`; the unused requirement is gone from go.mod and go.sum, and -h, -l and -v now work beside the long spellings. Moving the code into internal/cli and DR-0003 exit codes (every failure exits 1 today) is a follow-up, as done for osf.
- [x] Moved the five programs into `internal/cli` behind thin mains with the DR-0003 exit codes (2, 65, 66, 73, 74, 77), tested through the built binaries (`cli_exit_test.go`, goldens in `testdata/cli`). Found and fixed on the way: `opmlcat a b` and `opml2urls -o` created the input file (emptying it), `opmlsort` exited 0 on a missing input, `opml2urls` ignored positional arguments, `urls2opml` had no `-i`/`-o`. New: `-nl`, `-input`, `-output`, an EXIT STATUS section in each help page. Unreleased.
