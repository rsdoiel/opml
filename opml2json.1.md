%opml2json(1) user manual | version 0.0.10 f7001d7
% R. S. Doiel
% 2025-10-01

# NAME

opml2json — convert an OPML document to JSON

# SYNOPSIS

opml2json [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

opml2json reads an OPML document, sorts its outlines and writes it as JSON.

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

-pretty
: indent the JSON

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

    opml2json myfeeds.opml myfeeds.json

