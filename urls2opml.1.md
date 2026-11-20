---
title: urls2opml
section: 1
header: user manual
footer: version 0.0.12 c15bc0e
date: 2026-11-20
author: R. S. Doiel
---

# NAME

urls2opml — convert a list of URLs to OPML

# SYNOPSIS

urls2opml [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

urls2opml reads a text file with one URL per line and writes an OPML document with an outline for each http or https URL. Empty lines, comments and lines that are not usable URLs are skipped, with a warning on standard error.

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
: add a trailing newline (the default; -nl=false turns it off)

-quiet
: suppress error messages; the exit status still says what happened

The input and output may also be given as the first and second argument.

# EXIT STATUS

0
: success

2
: bad command line: an unknown option, a surplus argument

65
: the input cannot be read as text lines

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

Convert a newsboat url file to OPML.

    cut -d ' ' -f 1 ~/.newsboat/urls | grep -v '#' | urls2opml > subscriptions.opml

