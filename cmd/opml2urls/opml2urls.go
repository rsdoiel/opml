package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

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

{app_name} converts a OPML XML to a list of urls for elements that
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
cat subscriptions.xml | {app_name} > urls.txt
~~~

Convert an OPML XML file to the newsboart
`+"`"+`.newboat/urls`+"`"+` format.

~~~
cat subscriptions.xml | {app_name} -newsboat
~~~


`
)

var (
	// Common options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App options
	xmlurl        bool
	htmlurl       bool
	textAsComment bool
	newsboat      bool
)

func displayHelp(appName string, txt string) string {
	return strings.ReplaceAll(txt, `{app_name}`, appName)
}

func main() {
	appName := path.Base(os.Args[0])
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&xmlurl, "xmlurl", true, "output the xmlUrl attributes, defaults true")
	flag.BoolVar(&htmlurl, "htmlurl", false, "output the htmlUrl attributes, default false")
	flag.BoolVar(&textAsComment, "text-as-comment", true, "output the text attribute as comment, defaults to true")
	flag.BoolVar(&newsboat, "newsboat", false, "output in newsboat's url file format")
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

	src, err := io.ReadAll(in)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	o, err := opml.Parse(src)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	if err := o.Walk(func(elem *opml.Outline) bool {
		if elem == nil {
			return false
		}
		if newsboat {
			if elem.XMLURL != "" {
				fmt.Fprintf(out, "%s", elem.XMLURL)
				if elem.Text != "" {
					fmt.Fprintf(out, " %q", fmt.Sprintf("~%s", elem.Text))
				}
				fmt.Fprintf(out, "\n")
			}

		} else {
			if textAsComment && (elem.Text != "") {
				fmt.Fprintf(out, "# %s\n", elem.Text)
			}
			if xmlurl && (elem.XMLURL != "") {
				fmt.Fprintf(out, "%s\n", elem.XMLURL)
			}
			if htmlurl && (elem.HTMLURL != "") {
				fmt.Fprintf(out, "%s\n", elem.HTMLURL)
			}
		}
		return true
	}); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}

}
