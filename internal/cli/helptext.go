package cli

/** Opml2jsonHelpText is the opml2json(1) page, shown by `opml2json -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to opml2json.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.Opml2jsonHelpText, "opml2json", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const Opml2jsonHelpText = `---
title: {app_name}
section: 1
header: user manual
footer: version {version} {release_hash}
date: {release_date}
author: R. S. Doiel
---

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

/** Cmark2opmlHelpText is the cmark2opml(1) page, shown by `cmark2opml -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to cmark2opml.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.Cmark2opmlHelpText, "cmark2opml", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const Cmark2opmlHelpText = `---
title: {app_name}
section: 1
header: user manual
footer: version {version} {release_hash}
date: {release_date}
author: R. S. Doiel
---

# NAME

{app_name} — read a Markdown outline and write it as OPML

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

{app_name} reads the Markdown that opml2cmark writes, after you have edited it in a text editor, and writes an OPML document with the XML declaration and two spaces of indentation.

The Markdown is a list. Each list item is an outline and a nested list holds its sub-outlines. An item is text, or a link whose text is the outline text; a link address is the xmlUrl attribute when the outline's type is rss or atom, and the url attribute otherwise. Other attributes follow the item in braces:

    - [BBC News](http://feeds.example/bbc.xml) {type="rss" title="BBC News"}
    - Pages
      - [Home](https://example.org/index.html)

A block of YAML front matter at the top holds the head: title, dateCreated, dateModified, generator, ownerName, ownerEmail, ownerId, docs, expansionState, vertScrollState, windowTop, windowLeft, windowBottom, windowRight, and opmlVersion for a document that is not OPML 2.0. An unknown key is an error.

The list may be written the way a person writes one: any bullet (-, * or +) or numbers, any indentation, blank lines between items, an item wrapped over several lines (the lines are joined with a space), Windows line endings. Markdown typed in text, such as emphasis, is kept as typed; a backslash before punctuation is removed.

Anything that is not part of the list stops the run with an error that names the line, so no text is dropped without your knowing it: a heading, paragraph, code block, quote or rule, a second paragraph or code block inside an item, a reference-style link, a link title, text after a link, an attribute given twice, or the same address given as both the link and an attribute.

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
: add a trailing newline (the default is none; the output already ends in one)

-quiet
: suppress error messages; the exit status still says what happened

The input and output may also be given as the first and second argument.

# EXIT STATUS

0
: success; empty input gives an empty outline

2
: bad command line: an unknown option, a surplus argument

65
: the input is Markdown that is not an outline list; the message names the line

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

Edit an outline as Markdown and convert it back:

    opml2cmark subscriptions.opml subscriptions.md
    cmark2opml subscriptions.md subscriptions.opml
`

/** Opml2cmarkHelpText is the opml2cmark(1) page, shown by `opml2cmark -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to opml2cmark.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.Opml2cmarkHelpText, "opml2cmark", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const Opml2cmarkHelpText = `---
title: {app_name}
section: 1
header: user manual
footer: version {version} {release_hash}
date: {release_date}
author: R. S. Doiel
---

# NAME

{app_name} — write an OPML document as a CommonMark list

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME [OUTPUT_FILENAME]]

# DESCRIPTION

{app_name} reads an OPML document and writes it as Markdown, so the outline can be edited in any text editor. The result is read back with cmark2opml.

The head fields (title, dates, owner, generator and so on) are YAML front matter. Each outline is one list item and its sub-outlines are a nested list, in the order of the document; nothing is sorted. An outline that has a URL is written as a link: the text is the link text, and the address is the xmlUrl attribute for a type of rss or atom, the url attribute otherwise. Every other attribute is written after the item as a block of key="value" pairs in braces:

    - Pages
      - [Home](https://example.org/index.html)
    - [BBC News](http://feeds.example/bbc.xml) {type="rss" title="BBC News" htmlUrl="http://www.example/bbc"}

A document that is not OPML 2.0 has an opmlVersion key in its front matter. A document with no head fields and no outlines gives no output.

Some text cannot be written on one list line: a newline or tab, or white space at the start or end of the text. {app_name} stops with an error that names the outline, such as outline 2.1, and writes nothing.

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
: add a trailing newline (the default is none; the output already ends in one)

-quiet
: suppress error messages; the exit status still says what happened

The input and output may also be given as the first and second argument.

# EXIT STATUS

0
: success, including a document with nothing to write

2
: bad command line: an unknown option, a surplus argument

65
: the input is not valid OPML, or has text a list item cannot hold

66
: an input file does not exist

73
: the output file cannot be created

74
: reading or writing failed part way

77
: the operating system refused access to a file

# EXAMPLES

    opml2cmark subscriptions.opml subscriptions.md

Edit the list, then convert it back:

    cmark2opml subscriptions.md subscriptions.opml
`

/** Opml2urlsHelpText is the opml2urls(1) page, shown by `opml2urls -help`. The {app_name},
 * {version}, {release_hash} and {release_date} tokens are replaced by opml.FmtHelp.
 * The Makefile writes the page to opml2urls.1.md.
 *
 * Example:
 *   fmt.Println(opml.FmtHelp(cli.Opml2urlsHelpText, "opml2urls", opml.Version, opml.ReleaseDate, opml.ReleaseHash))
 */
const Opml2urlsHelpText = `---
title: {app_name}
section: 1
header: user manual
footer: version {version} {release_hash}
date: {release_date}
author: R. S. Doiel
---

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
: success, including a document that has nothing to list (nothing is written)

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
const OpmlcatHelpText = `---
title: {app_name}
section: 1
header: user manual
footer: version {version} {release_hash}
date: {release_date}
author: R. S. Doiel
---

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
const OpmlsortHelpText = `---
title: {app_name}
section: 1
header: user manual
footer: version {version} {release_hash}
date: {release_date}
author: R. S. Doiel
---

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
const Urls2opmlHelpText = `---
title: {app_name}
section: 1
header: user manual
footer: version {version} {release_hash}
date: {release_date}
author: R. S. Doiel
---

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
