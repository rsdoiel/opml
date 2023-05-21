% opml2urls(1) user manual
% R. S. Doiel
% 2022-12-16

# NAME

opml2urls

# SYNOPSIS

opml2urls converts a OPML XML to a list of urls for elements that
have the xmlUrl attribute set.

# OPTIONS

help
: Display this help page

version
: Display version

-license
: Display license

-xmlurl
: output the xmlUrl attribute (default is true)

-htmlurl
: output the htmlUrl attribute (default is false)

-text-as-comment
: output the text attribute as a hash prefix comment

-newsboat
: convert the OPML into a newsboat url file format

# EXAMPLES

Convert OPML XML to a list of plain text URLs 
from the xmlUrl attributes on the OPML.

~~~
cat subscriptions.xml | opml2urls > urls.txt
~~~

Convert an OPML XML file to the newsboart
`.newboat/urls` format.

~~~
cat subscriptions.xml | opml2urls -newsboat
~~~



