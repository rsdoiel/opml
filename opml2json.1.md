%s
%opml2json(1) | version 0.0.9 db50e8d
% R. S. Doiel
% 2024-05-20

# NAME

opml2json

# SYNOPSIS

opml2json [OPTIONS] [INPUT_FILENAME] [OUTPUT_FILENAME]

# DESCRIPTION

opml2json is a program that converts OPML's XML to JSON.

# OPTIONS

-help
: display help

-license
: display license

-version
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

