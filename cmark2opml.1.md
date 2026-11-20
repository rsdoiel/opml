%cmark2opml(1) user manual | version 0.0.12 c15bc0e
% R. S. Doiel
% 2026-11-20

# NAME

cmark2opml — read a Markdown outline and write it as OPML

# SYNOPSIS

cmark2opml [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

cmark2opml reads the Markdown that opml2cmark writes, after you have edited it in a text editor, and writes an OPML document with the XML declaration and two spaces of indentation.

The Markdown is a list. Each list item is an outline and a nested list holds its sub-outlines. An item is text, or a link whose text is the outline text; a link address is the xmlUrl attribute when the outline's type is rss or atom, and the url attribute otherwise. Other attributes follow the item in braces:

    - [BBC News](http://feeds.example/bbc.xml) {type="rss" title="BBC News"}
    - Pages
      - [Home](https://example.org/index.html)

A block of YAML front matter at the top holds the head: title, dateCreated, dateModified, generator, ownerName, ownerEmail, ownerId, docs, expansionState, vertScrollState, windowTop, windowLeft, windowBottom, windowRight, and opmlVersion for a document that is not OPML 2.0. An unknown key is an error.

The list may be written the way a person writes one: any bullet (-, * or +) or numbers, any indentation, blank lines between items, an item wrapped over several lines (the lines are joined with a space), Windows line endings. Markdown typed in text, such as emphasis, is kept as typed; a backslash before punctuation is removed.

Anything that is not part of the list stops the run with an error that names the line, so no text is dropped without your knowing it: a heading, paragraph, code block, quote or rule, a second paragraph or code block inside an item, a reference-style link, a link title, text after a link, an attribute given twice, or the same address given as both the link and an attribute.

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
: success; empty input gives an empty outline

2
: bad command line: an unknown option, a surplus argument

65
: the input is Markdown that is not an outline list; the message names the line

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

Edit an outline as Markdown and convert it back:

    opml2cmark subscriptions.opml subscriptions.md
    cmark2opml subscriptions.md subscriptions.opml

