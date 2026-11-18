%opml2urls(1) user manual | version 0.0.11 a083b9b
% R. S. Doiel
% 2026-11-18

# NAME

opml2urls — list the feed URLs of an OPML document

# SYNOPSIS

opml2urls [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

opml2urls reads an OPML document and writes the URLs of the outlines that have an xmlUrl attribute, one per line.

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

-xmlurl
: write the xmlUrl attribute (the default; -xmlurl=false turns it off)

-htmlurl
: write the htmlUrl attribute (the default is not to)

-text-as-comment
: write the text attribute as a comment line starting with # (the default; -text-as-comment=false turns it off)

-newsboat
: write newsboat's url file format

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

    cat subscriptions.opml | opml2urls > urls.txt

Write newsboat's url file format:

    opml2urls -newsboat subscriptions.opml urls

