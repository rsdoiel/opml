% urls2opml(1) user manual | 0.0.8 db50e8d
% R. S. Doiel
% 2023-06-20

# NAME

urls2opml

# SYNOPSIS

urls2opml converts a text file, one url per line, to OPML
XML.  It reads from standard input and writes to standard output.

# OPTIONS

-help
: Display this help page

-version
: Display version

-license
: Display license


# EXAMPLE

Convert a newsboat "url" file to OPML.

~~~
cat .newsboat/url | cut -d \" -f 1 |\
    grep -v '#' | urls2opml \
	>subscriptions.opml
~~~

