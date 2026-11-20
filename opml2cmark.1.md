%opml2cmark(1) user manual | version 0.0.12 c15bc0e
% R. S. Doiel
% 2026-11-20

# NAME

opml2cmark — write an OPML document as a CommonMark list

# SYNOPSIS

opml2cmark [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

opml2cmark reads an OPML document and writes it as Markdown, so the outline can be edited in any text editor. The result is read back with cmark2opml.

The head fields (title, dates, owner, generator and so on) are YAML front matter. Each outline is one list item and its sub-outlines are a nested list, in the order of the document; nothing is sorted. An outline that has a URL is written as a link: the text is the link text, and the address is the xmlUrl attribute for a type of rss or atom, the url attribute otherwise. Every other attribute is written after the item as a block of key="value" pairs in braces:

    - Pages
      - [Home](https://example.org/index.html)
    - [BBC News](http://feeds.example/bbc.xml) {type="rss" title="BBC News" htmlUrl="http://www.example/bbc"}

A document that is not OPML 2.0 has an opmlVersion key in its front matter. A document with no head fields and no outlines gives no output.

Some text cannot be written on one list line: a newline or tab, or white space at the start or end of the text. opml2cmark stops with an error that names the outline, such as outline 2.1, and writes nothing.

# OPTIONS

-h, -help
: display this help page

-l, -license
: display the license

-v, -version
: display the version

-i, -input FILENAME
: read from FILENAME; the default, and -, is standard input

-o, -output FILENAME
: write to FILENAME; the default, and -, is standard output

-nl, -newline
: add a trailing newline (the default is none; the output already ends in one)

-quiet
: suppress error messages; the exit status still says what happened

The input and output may also be given as the first and second argument.

# EXIT STATUS

0
: success, including a document with nothing to write

2
: bad command line: an unknown option, a surplus argument

65
: the input is not valid OPML, or has text a list item cannot hold

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

    opml2cmark subscriptions.opml subscriptions.md

Edit the list, then convert it back:

    cmark2opml subscriptions.md subscriptions.opml

