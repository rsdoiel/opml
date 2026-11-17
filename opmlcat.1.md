%opmlcat(1) user manual | version 0.0.10 f7001d7
% R. S. Doiel
% 2025-10-01

# NAME

opmlcat — concatenate OPML documents

# SYNOPSIS

opmlcat [OPTIONS] [INPUT_FILENAME ...]

# DESCRIPTION

opmlcat reads one or more OPML documents and writes one document with their outlines side by side. With no input file it reads standard input. Every argument is an input file; use -o for the output.

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
: indent the XML

Every argument is an input file; the output is -o.

# EXIT STATUS

0
: success

2
: bad command line: an unknown option, a surplus argument

65
: an input is not valid OPML, or cannot be combined

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

    opmlcat file1.opml file2.opml | opmlsort -o combined-sorted.opml

