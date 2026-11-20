package cli

/** Opml2jsonHelpText is the opml2json(1) page, shown by `opml2json -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to opml2json.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.Opml2jsonHelpText, "opml2json", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const Opml2jsonHelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} — convert an OPML document to JSON

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

{app_name} reads an OPML document, sorts its outlines and writes it as JSON.

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

    {app_name} myfeeds.opml myfeeds.json
`

/** Opml2urlsHelpText is the opml2urls(1) page, shown by `opml2urls -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to opml2urls.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.Opml2urlsHelpText, "opml2urls", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const Opml2urlsHelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} — list the feed URLs, or the page URLs, of an OPML document

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

{app_name} reads an OPML document and writes the URLs of the outlines that have an xmlUrl attribute, one per line. Outlines are walked in document order, at any depth.

A table of contents (for example .well-known/toc.opml) names pages with the url attribute and has no feed URLs. When a document has no xmlUrl or htmlUrl attribute at all, {app_name} writes the url attributes, as if -url were given. An outline with type="include" has a url that is another outline, not a page, so it is left out unless -includes is given.

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

-url
: write the url attribute (the default only for a document with no xmlUrl or htmlUrl attribute)

-includes
: also write the url of outlines with type="include"

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

    cat subscriptions.opml | {app_name} > urls.txt

Write newsboat's url file format:

    {app_name} -newsboat subscriptions.opml urls
`

/** OpmlcatHelpText is the opmlcat(1) page, shown by `opmlcat -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to opmlcat.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.OpmlcatHelpText, "opmlcat", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const OpmlcatHelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} — concatenate OPML documents

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME ...]

# DESCRIPTION

{app_name} reads one or more OPML documents and writes one document with their outlines side by side. With no input file it reads standard input. Every argument is an input file; use -o for the output.

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

    {app_name} file1.opml file2.opml | opmlsort -o combined-sorted.opml
`

/** OpmlsortHelpText is the opmlsort(1) page, shown by `opmlsort -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to opmlsort.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.OpmlsortHelpText, "opmlsort", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const OpmlsortHelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} — sort the outlines of an OPML document

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

{app_name} reads an OPML document and writes it with its outlines sorted.

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

-case-insensitive
: sort without regard to case

-title
: sort by title (the default; -title=false sorts by text)

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

    {app_name} myfeeds.opml sorted-feeds.opml
`

/** Urls2opmlHelpText is the urls2opml(1) page, shown by `urls2opml -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to urls2opml.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.Urls2opmlHelpText, "urls2opml", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const Urls2opmlHelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} — convert a list of URLs to OPML

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

{app_name} reads a text file with one URL per line and writes an OPML document with an outline for each http or https URL. Empty lines, comments and lines that are not usable URLs are skipped, with a warning on standard error.

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
: add a trailing newline (the default; -nl=false turns it off)

-quiet
: suppress error messages; the exit status still says what happened

The input and output may also be given as the first and second argument.

# EXIT STATUS

0
: success

2
: bad command line: an unknown option, a surplus argument

65
: the input cannot be read as text lines

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

Convert a newsboat url file to OPML.

    cut -d ' ' -f 1 ~/.newsboat/urls | grep -v '#' | {app_name} > subscriptions.opml
`
