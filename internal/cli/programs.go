package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/rsdoiel/opml"
	"github.com/rsdoiel/opml/cmark"
)

/** Opml2json is opml2json: it reads an OPML document and writes it as JSON.
 *
 * Parameters:
 *   appName (string) — the program name, for help and error messages
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Opml2json("opml2json", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Opml2json(appName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var o options
	var pretty bool
	fset := newFlagSet(appName, &o, false)
	fset.BoolVar(&pretty, "pretty", false, "pretty print JSON output")
	if _, code, done := start(appName, Opml2jsonHelpText, fset, &o, args, 2, stdout, stderr); done {
		return code
	}
	doc, code := readOPML(appName, &o, o.input, stdin, stderr)
	if code != ExitOK {
		return code
	}
	doc.Sort()
	var src []byte
	var err error
	if pretty {
		src, err = json.MarshalIndent(doc, "", "    ")
	} else {
		src, err = json.Marshal(doc)
	}
	if err != nil {
		return fail(appName, &o, stderr, ExitInternal, "%v", err)
	}
	return write(appName, &o, src, stdout, stderr)
}

/** Cmark2opml is cmark2opml: it reads a Markdown outline, the form opml2cmark
 * writes, and writes it as an OPML document with the XML declaration and two
 * spaces of indentation (opml DR-0001). Markdown that is not an outline list
 * is a data error that names the line.
 *
 * Parameters:
 *   appName (string) — the program name, for help and error messages
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Cmark2opml("cmark2opml", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Cmark2opml(appName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var o options
	fset := newFlagSet(appName, &o, false)
	if _, code, done := start(appName, Cmark2opmlHelpText, fset, &o, args, 2, stdout, stderr); done {
		return code
	}
	src, code := readSource(appName, &o, o.input, stdin, stderr)
	if code != ExitOK {
		return code
	}
	doc, err := cmark.ToOPML(src)
	if err != nil {
		if errors.Is(err, cmark.ErrNotOutline) {
			if namedInput(o.input) {
				return fail(appName, &o, stderr, ExitData, "%s: %v", o.input, err)
			}
			return fail(appName, &o, stderr, ExitData, "%v", err)
		}
		return fail(appName, &o, stderr, ExitInternal, "%v", err)
	}
	out, err := opml.MarshalIndent(doc)
	if err != nil {
		return fail(appName, &o, stderr, ExitInternal, "%v", err)
	}
	return write(appName, &o, out, stdout, stderr)
}

/** Opml2cmark is opml2cmark: it reads an OPML document and writes it as a
 * CommonMark list, with the head as YAML front matter, so the outline can be
 * edited in a text editor and read back with cmark2opml (opml DR-0001). The
 * outlines keep their order.
 *
 * Parameters:
 *   appName (string) — the program name, for help and error messages
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Opml2cmark("opml2cmark", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Opml2cmark(appName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var o options
	fset := newFlagSet(appName, &o, false)
	if _, code, done := start(appName, Opml2cmarkHelpText, fset, &o, args, 2, stdout, stderr); done {
		return code
	}
	doc, code := readOPML(appName, &o, o.input, stdin, stderr)
	if code != ExitOK {
		return code
	}
	src, err := cmark.FromOPML(doc)
	if err != nil {
		if errors.Is(err, cmark.ErrUnrepresentable) {
			return fail(appName, &o, stderr, ExitData, "%v", err)
		}
		return fail(appName, &o, stderr, ExitInternal, "%v", err)
	}
	return write(appName, &o, src, stdout, stderr)
}

/** Opml2urls is opml2urls: it reads an OPML document and writes the feed
 * URLs, as a plain list or in newsboat's url file format.
 *
 * Parameters:
 *   appName (string) — the program name, for help and error messages
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Opml2urls("opml2urls", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Opml2urls(appName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var o options
	var xmlurl, htmlurl, urlAttr, includes, textAsComment, newsboat bool
	fset := newFlagSet(appName, &o, false)
	fset.BoolVar(&xmlurl, "xmlurl", true, "output the xmlUrl attributes")
	fset.BoolVar(&htmlurl, "htmlurl", false, "output the htmlUrl attributes")
	fset.BoolVar(&urlAttr, "url", false, "output the url attributes (on by itself when the document has no xmlUrl or htmlUrl)")
	fset.BoolVar(&includes, "includes", false, "also output the url of include outlines")
	fset.BoolVar(&textAsComment, "text-as-comment", true, "output the text attribute as a comment")
	fset.BoolVar(&newsboat, "newsboat", false, "output in newsboat's url file format")
	if _, code, done := start(appName, Opml2urlsHelpText, fset, &o, args, 2, stdout, stderr); done {
		return code
	}
	doc, code := readOPML(appName, &o, o.input, stdin, stderr)
	if code != ExitOK {
		return code
	}
	// Nothing to list is an empty listing, not an error: exit 0.
	if doc.Body == nil || len(doc.Body.Outline) == 0 {
		return write(appName, &o, nil, stdout, stderr)
	}
	// A table of contents has url attributes and no feed URLs: list the
	// url attribute without being asked.
	if !urlAttr && !newsboat {
		feeds := false
		_ = doc.Walk(func(elem *opml.Outline) bool {
			feeds = feeds || elem.XMLURL != "" || elem.HTMLURL != ""
			return !feeds
		})
		urlAttr = !feeds
	}
	var buf bytes.Buffer
	if err := doc.Walk(func(elem *opml.Outline) bool {
		if elem == nil {
			return false
		}
		// An include's url is another outline, not a page.
		if !newsboat && strings.EqualFold(elem.Type, "include") && !includes {
			return true
		}
		if newsboat {
			if elem.XMLURL != "" {
				fmt.Fprintf(&buf, "%s", elem.XMLURL)
				if elem.Text != "" {
					fmt.Fprintf(&buf, " %q", fmt.Sprintf("~%s", elem.Text))
				}
				fmt.Fprintf(&buf, "\n")
			}
		} else {
			if textAsComment && elem.Text != "" {
				fmt.Fprintf(&buf, "# %s\n", elem.Text)
			}
			if xmlurl && elem.XMLURL != "" {
				fmt.Fprintf(&buf, "%s\n", elem.XMLURL)
			}
			if htmlurl && elem.HTMLURL != "" {
				fmt.Fprintf(&buf, "%s\n", elem.HTMLURL)
			}
			if urlAttr && elem.URL != "" {
				fmt.Fprintf(&buf, "%s\n", elem.URL)
			}
		}
		return true
	}); err != nil {
		return fail(appName, &o, stderr, ExitInternal, "%v", err)
	}
	return write(appName, &o, buf.Bytes(), stdout, stderr)
}

/** Opmlcat is opmlcat: it reads one or more OPML documents, standard input
 * when none is named, and writes them as one document with their outlines
 * side by side.
 *
 * Parameters:
 *   appName (string) — the program name, for help and error messages
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Opmlcat("opmlcat", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Opmlcat(appName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var o options
	var pretty bool
	fset := newFlagSet(appName, &o, false)
	fset.BoolVar(&pretty, "pretty", false, "pretty print XML output")
	inputs, code, done := start(appName, OpmlcatHelpText, fset, &o, args, 0, stdout, stderr)
	if done {
		return code
	}
	if o.input != "" {
		inputs = append([]string{o.input}, inputs...)
	}
	if len(inputs) == 0 {
		inputs = []string{"-"}
	}
	combined := opml.New()
	for i, name := range inputs {
		doc, code := readOPML(appName, &o, name, stdin, stderr)
		if code != ExitOK {
			return code
		}
		if i == 0 {
			combined = doc
			continue
		}
		if err := combined.Append(doc); err != nil {
			return fail(appName, &o, stderr, ExitData, "%s: %v", name, err)
		}
	}
	return writeXML(appName, &o, combined, pretty, stdout, stderr)
}

/** Opmlsort is opmlsort: it reads an OPML document and writes it with its
 * outlines sorted.
 *
 * Parameters:
 *   appName (string) — the program name, for help and error messages
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Opmlsort("opmlsort", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Opmlsort(appName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var o options
	var pretty, caseInsensitive, byTitle bool
	fset := newFlagSet(appName, &o, false)
	fset.BoolVar(&pretty, "pretty", false, "pretty print XML output")
	fset.BoolVar(&caseInsensitive, "case-insensitive", false, "case insensitive sort")
	fset.BoolVar(&byTitle, "title", true, "sort by title")
	if _, code, done := start(appName, OpmlsortHelpText, fset, &o, args, 2, stdout, stderr); done {
		return code
	}
	doc, code := readOPML(appName, &o, o.input, stdin, stderr)
	if code != ExitOK {
		return code
	}
	switch {
	case byTitle && caseInsensitive:
		doc.SortTitleCaseInsensitive()
	case byTitle:
		doc.SortTitle()
	case caseInsensitive:
		doc.SortCaseInsensitive()
	default:
		doc.Sort()
	}
	return writeXML(appName, &o, doc, pretty, stdout, stderr)
}

// writeXML writes doc after the XML declaration, indented when pretty is true.
func writeXML(appName string, o *options, doc *opml.OPML, pretty bool, stdout, stderr io.Writer) int {
	var src []byte
	if pretty {
		var err error
		if src, err = xml.MarshalIndent(doc, "", "    "); err != nil {
			return fail(appName, o, stderr, ExitInternal, "%v", err)
		}
	} else {
		src = []byte(doc.String())
	}
	out := append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"), src...)
	return write(appName, o, out, stdout, stderr)
}

/** Urls2opml is urls2opml: it reads a text file with one URL per line and
 * writes an OPML document with an outline for each http or https URL. Lines
 * that are empty, comments or not usable URLs are skipped with a warning on
 * standard error.
 *
 * Parameters:
 *   appName (string) — the program name, for help and error messages
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Urls2opml("urls2opml", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Urls2opml(appName string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var o options
	fset := newFlagSet(appName, &o, true)
	if _, code, done := start(appName, Urls2opmlHelpText, fset, &o, args, 2, stdout, stderr); done {
		return code
	}
	src, code := readSource(appName, &o, o.input, stdin, stderr)
	if code != ExitOK {
		return code
	}
	doc := opml.New()
	doc.Head.Title = fmt.Sprintf("url list convert with %s %s", appName, opml.Version)
	doc.Head.Created = time.Now().Format(time.RFC822Z)
	doc.Body.Outline = []*opml.Outline{}
	scan := bufio.NewScanner(bytes.NewReader(src))
	lineNo := 0
	for scan.Scan() {
		lineNo++
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			warn(appName, &o, stderr, "line %d is empty", lineNo)
			continue
		}
		if strings.HasPrefix(line, "#") {
			warn(appName, &o, stderr, "line %d, skipping comment %q", lineNo, line)
			continue
		}
		u, err := url.Parse(line)
		if err != nil {
			warn(appName, &o, stderr, "line %d not a url %q, %s", lineNo, line, err)
			continue
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			warn(appName, &o, stderr, "line %d, skipping unsupported url %q", lineNo, u.String())
			continue
		}
		doc.Body.Outline = append(doc.Body.Outline, &opml.Outline{XMLURL: line})
	}
	if err := scan.Err(); err != nil {
		return fail(appName, &o, stderr, ExitIO, "%v", err)
	}
	return writeXML(appName, &o, doc, true, stdout, stderr)
}
