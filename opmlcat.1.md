%opmlcat(1) | version 0.0.8 7dd10d4
% R. S. Doiel
% 2023-05-20

# NAME

opmlcat

# SYNOPSIS

opmlcat [OPTIONS]

# DESCRIPTION

opmlcat concatenates one or more opml files as siblings to standard out.

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
: add trailing newline

-pretty
: pretty print JSON output

# EXAMPLES

This is an example of using opmlcat and opmlsort together to 
create a combined sorted opml file.

~~~
    opmlcat file1.opml file1.opml | opmlsort -o combined-sorted.opml
~~~


