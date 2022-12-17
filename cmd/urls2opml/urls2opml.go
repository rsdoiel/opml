package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	// My packages
	"github.com/rsdoiel/opml"
)

const (
	HelpText = `% {app_name}(1) user manual
% R. S. Doiel
% 2022-12-16

# NAME

{app_name}

# SYNOPSIS

{app_name} converts a text file, one url per line, to OPML
XML.  It reads from standard input and writes to standard output.

# OPTIONS

help
: Display this help page

version
: Display version

-license
: Display license

# EXAMPLE

Convert a newsboat "url" file to OPML.

~~~
cat .newsboat/url | cut -d \" -f 1 |\
    grep -v '#' | {app_name} \
	>subscriptions.opml
~~~

`
)

var (
	showHelp    bool
	showLicense bool
	showVersion bool
)

func displayHelp(appName string, txt string) string {
	return strings.ReplaceAll(txt, `{app_name}`, appName)
}

func main() {
	appName := path.Base(os.Args[0])
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.Parse()

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr
	if showHelp {
		fmt.Fprintf(out, "%s", displayHelp(appName, HelpText))
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s\n", appName, opml.Version)
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", opml.LicenseText)
		os.Exit(0)
	}

	//FIXME: Should allow for creation/pick of outline to append to.
	// E.g. mimik a "subscripts" outline element that has URLs as children.
	label := fmt.Sprintf("url list convert with %s %s", appName, opml.Version)
	o := opml.New()
	o.Head.Title = label
	o.Head.Created = time.Now().Format(time.RFC822Z)
	o.Body.Outline = []*opml.Outline{}
	scan := bufio.NewScanner(in)
	i := 1
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			fmt.Fprintf(eout, "line %d is empty\n", i)
			continue
		}
		if strings.HasPrefix(line, "#") {
			fmt.Fprintf(eout, "line %d, skipping comment %q\n", i, line)
			continue
		}
		u, err := url.Parse(line)
		if err != nil {
			fmt.Fprintf(eout, "line %d not a url %q, %s\n", i, line, err)
			continue
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			fmt.Fprintf(eout, "line %d, skipping unsupported url %q\n", i, u.String())
			continue
		}
		// FIXME: there should be an option to verify the link
		// is still available by executing a GET, it could
		// then populate the Outline element appropriately
			elem := new(opml.Outline)
			elem.XMLURL = line
			o.Body.Outline = append(o.Body.Outline, elem)
		i++
	}
	if err := scan.Err(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	src, err := xml.MarshalIndent(o, "", "    ")
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(out, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintf(out, "%s\n", src)
}
