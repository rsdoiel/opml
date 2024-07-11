%opmlsort(1) | version 0.0.9 fdad103
% R. S. Doiel
% 2024-05-20

# NAME

opmlsort

# SYNOPSIS

opmlsort [OPTIONS] [INPUT_FILENAME] [OUTPUT_FILENAME]

# DESCRIPTION

opmlsort is a program that sorts the outline in an OPML document.

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
: add a tailing newline

-pretty
: pretty print JSON output

-case-insensitive
: case insensitive sort

-title
: sort by title


# EXAMPLES

~~~
    opmlsort myfeeds.opml sorted-feeds.opml
~~~


