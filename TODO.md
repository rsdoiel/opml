{
    "title": "OPML, a Go package and cli"
}

## Next

+ [ ] Added -newsboat option to urls2opml, need to add the url description tas text attribute and optionally use a sub-list to import nicely into NetNewsWire
+ [ ] create a command line tool that reads an OPML and appends a element to the list (e.g. adds a feed URL to an OPML list of feeds)
    + basic verbs would be insert (insert a new list element), append (a new list element), replace (replace a list element) delete (a list element), and find (return the path to an element by name or attribute value)
        + append, insert, replace takes a path and the value to update with
        + delete takes a path to the item to be deleted, if you try to delete the root list you get back an empty list
        + find takes an attribute name and value and returns a path
    + the tree the location for the verb to operate would be assumed to be the root or the index number of the items as a path
        + append "URL" would append the url to the root of the list
        + append "URL" /3 would append the UR" to the with item in the root list
        + append "URL" /3/2 would apend the URL the third items' second entry 
        + find ITEN_NAME append NEW_ITEM would find the item in the OPML tree, then append the content
+ [ ] create a tool that can read an OPML file, harvest the current feeds, index for browsing and support an option to send interesting articles to Pocket
+ [ ] Add support to process Frontier's fttb into OPML
    + See http://scripting.com/fatpages/about.html, http://scripting.com/fatpages/faq.html and http://scripting.com/fatpages/outline.html
    + fttb is a "fatpages" document, it is a Base 64 encoded document like is done with email.

## From the antennaApp table-of-contents work (2026-11-19)

Found by running the module against Dave Winer's published `toc.opml` files while writing antennaApp DR-0013 (see also the `opml` project in the knowledge base). `Walk` and the `langauge` tag were found on the released v0.0.10 and are already fixed in v0.0.11 (see Completed); these are what is still open in the current source.

+ [ ] `Parse`, `ReadFile` and `Unmarshal` fail on any document that declares an encoding other than UTF-8: `xml: encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil`. Dave Winer's table of contents (https://dave.writersweb.org/scripting/.well-known/toc.opml, 456 KB, 3,562 outlines) declares ISO-8859-1, so the module cannot read the one real `toc.opml` found so far.
    + decode with `xml.NewDecoder` and a `CharsetReader`
    + ISO-8859-1 and US-ASCII need only a byte-to-rune reader; windows-1252 needs a table (`golang.org/x/text/encoding/charmap`)
    + add a small ISO-8859-1 file to `testdata`, with a test, before the fix
+ [ ] `Marshal`, `String` and `WriteFile` write no XML declaration and no indentation, so a file is a single line
    + add an indented form that starts with `<?xml version="1.0" encoding="UTF-8"?>` (keep today's output available; the order is deterministic, which callers rely on)
    + `WriteFile` ignores the marshal error and uses the deprecated `ioutil`: return the error and use `os`
+ [ ] `Head.OwnerID` is written as `OwnerId`; the OPML spec says `ownerId`
+ [ ] `Head.Modified` has the JSON tag `dataModified` (should be `dateModified`; this changes `opml2json` output, so note it in the release notes)
+ [ ] Comments say "RFC 882" where they mean RFC 822 (`Head.Created`, `Head.Modified`, `Outline.Created`)
+ [ ] Helpers for reading a table of contents (antennaApp, cmarkbrowser's previous and next from `/.well-known/toc.opml`):
    + `Flatten()`: the outlines that have a `url` and are not `type="include"`, in document order
    + `Includes()`: the addresses of the `include` outlines
+ [ ] `opml2urls` was written when URL lists were flat and does not suit an outline of pages (found 2026-11-19 trying it on antennaApp's `toc.opml`, which it lists as empty)
    + it prints only the `xmlUrl` and `htmlUrl` attributes; a table of contents uses `url` on its entries. Add a `-url` option for the `url` attribute, and print `url` when a document has no `xmlUrl` or `htmlUrl`
    + decide what to do with `type="include"` outlines (their `url` is another outline, not a page): list them apart, or only with an option
    + depth: it should walk nested outlines in document order (the `Walk` fix in v0.0.11 is needed)
    + the installed `~/bin/opml2urls` is 0.0.10 (2026-11-16), from before the fix that makes it honour a file argument (2026-11-17): it ignores the file name and waits on standard input. Re-install after the next release
+ [ ] `Head` cannot hold a `generator` child element, which Dave's `toc.opml` has: decide whether to keep unknown child elements or leave them out
+ [ ] When the above are done, release v0.0.12; antennaApp pins v0.0.10 and should move to v0.0.11 now for the `Walk` fix

## Someday, maybe

+ opmlviewer - a cli/terminal based opml viewer
+ omplfilter
    + might function something like Unix find command
+ ompledit - CRUD operations to individual Outline elements and element lists
+ Review River5 by Dave Winer for insights into what additional functions are needed in package 

## Completed

+ [x] Fix `Walk` visiting only the first child at each level, and stopping early from a nested outline (2026-11-18)
+ [x] Fix the `Outline.Language` struct tag typo (`langauge`), so the `language` attribute is read and written (2026-11-18)
+ [x] Add opml2json 
+ [x] Add support for custom attributes
+ [x] Add Bash script to fetch Dave Winer's userland samples at http://scripting.com/misc/userlandSamples.zip

## Modernization (2026-11-16)

To match the workspace standard (knowledge and opml are the models; Go 1.27.3 across repositories). Do one repository at a time; run the tests and `make -n release` after each step.

- [x] Replace `github.com/caltechlibrary/cli`: the programs already used the standard library `flag`; the unused requirement is gone from go.mod and go.sum, and -h, -l and -v now work beside the long spellings. Moving the code into internal/cli and DR-0003 exit codes (every failure exits 1 today) is a follow-up, as done for osf.
- [x] Moved the five programs into `internal/cli` behind thin mains with the DR-0003 exit codes (2, 65, 66, 73, 74, 77), tested through the built binaries (`cli_exit_test.go`, goldens in `testdata/cli`). Found and fixed on the way: `opmlcat a b` and `opml2urls -o` created the input file (emptying it), `opmlsort` exited 0 on a missing input, `opml2urls` ignored positional arguments, `urls2opml` had no `-i`/`-o`. New: `-nl`, `-input`, `-output`, an EXIT STATUS section in each help page. Unreleased.
