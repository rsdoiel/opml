---
title: opmlsort
section: 1
header: user manual
footer: version 0.0.12 b8a4071
date: 2026-11-20
author: R. S. Doiel
---

# NAME

opmlsort — sort the outlines of an OPML document

# SYNOPSIS

opmlsort [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

opmlsort reads an OPML document and writes it with its outlines sorted.

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
: add a trailing newline (the default is none)

-quiet
: suppress error messages; the exit status still says what happened

-pretty
: indent the XML

-case-insensitive
: sort without regard to case

-title
: sort by title (the default; -title=false sorts by text)

The input and output may also be given as the first and second argument.

# EXIT STATUS

0
: success

2
: bad command line: an unknown option, a surplus argument

65
: the input is not valid OPML

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

    opmlsort myfeeds.opml sorted-feeds.opml

