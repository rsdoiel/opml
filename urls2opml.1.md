% urls2opml(1) user manual | 0.0.9 fdad103
% R. S. Doiel
% 2024-05-20

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

