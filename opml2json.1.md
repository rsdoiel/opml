% opml2json(1) user manual
% R. S. Doiel
% 2022-12-16

# NAME

opml2json

# SYNOPSIS

opml2json [OPTIONS] INPUT_OPML_FILENAME [OUTPUT_OPML_FILENAME]

# DESCRIPTION

opml2json is a program that converts OPML's XML to JSON.

# OPTIONS

-examples
: display examples

-generate-manpage
: generate man page

-generate-markdown
: generate Markdown documentation

-h, -help
: display help

-i, -input
: set input filename

-l, -license
: display license

-nl, -newline
: add trailing newline

-o, -output
: set output filename

-p, -pretty
: pretty print XML output

-quiet
: suppress error messages

-v, -version
: display version


# EXAMPLES

Convert *myfeeds.ompl* to *myfeeds.json*.

~~~
    opml2json myfeeds.opml myfeeds.json
~~~
