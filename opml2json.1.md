%opml2json(1) | version 0.0.10 f7001d7
% R. S. Doiel
% 2025-10-01

# NAME

opml2json

# SYNOPSIS

opml2json [OPTIONS] [INPUT_FILENAME] [OUTPUT_FILENAME]

# DESCRIPTION

opml2json is a program that converts OPML's XML to JSON.

# OPTIONS

-h, -help
: display help

-l, -license
: display license

-v, -version
: display version

-i
: read from filename

-o
: write to filename

-newline
: add a trailing newline

-pretty
: pretty print JSON output


# EXAMPLES

Convert *myfeeds.ompl* to *myfeeds.json*.

~~~
opml2json myfeeds.opml myfeeds.json
~~~


