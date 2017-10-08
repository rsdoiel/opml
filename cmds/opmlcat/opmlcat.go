//
// opmlcat is a command line utility that reads in one or more OMPL files, concatenates them
// at their roots and returns a single file as a result.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2016 all rights reserved.
// Released under the BSD 3-Clause License
// See: https://opensource.org/licenses/BSD-3-Clause
//
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	// My Packages
	"github.com/rsdoiel/opml"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
)

var (
	usage = `USAGE: %s [OPTIONS] OPML_FILE [OPML_FILE ...]`

	description = `
SYNOPSIS

%s concatenates one or more opml files as siblings to standard out.
`

	examples = `
EXAMPLES

This is an example of using %s and opmlsort together to 
create a combined sorted opml file.

    %s file1.opml file1.opml | opmlsort -o combined-sorted.opml
`

	// Standard options
	showHelp    bool
	showVersion bool
	showLicense bool

	// Application options
	prettyPrint bool
)

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// Application Options
	flag.BoolVar(&prettyPrint, "p", false, "pretty print XML output")
	flag.BoolVar(&prettyPrint, "pretty", false, "pretty print XML output")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, "OPML",
		fmt.Sprintf(opml.LicenseText, appName, opml.Version),
		opml.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName)
	cfg.OptionsText = "OPTIONS\n"
	cfg.ExampleText = fmt.Sprintf(examples, appName, appName)

	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	o := opml.New()
	if len(args) == 0 {
		src, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Missing input, %s", err)
			os.Exit(1)
		}
		o, err = opml.Parse(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error on stdin, %s", err)
			os.Exit(1)
		}
	}

	for _, iFName := range args {
		next := opml.New()
		err := next.ReadFile(iFName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't read %s, %s\n", iFName, err)
			os.Exit(1)
		}
		err = o.Append(next)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
	}
	if prettyPrint == true {
		src, _ := xml.MarshalIndent(o, "  ", "    ")
		fmt.Printf("%s", src)
	} else {
		fmt.Println(o.String())
	}
}
