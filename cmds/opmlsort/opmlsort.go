//
// opmlsort is a command line utility that can read in a OPML file, sort the outline
// and return the results.
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

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
)

var (
	usage = `USAGE: %s [OPTIONS] INPUT_OPML_FILENAME [OUTPUT_OPML_FILENAME]`

	description = `
SYNOPSIS

%s is a program that sorts the outline in an OPML document.
`

	examples = `
EXAMPLES

    %s myfeeds.opml sorted-feeds.opml
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
	cfg.ExampleText = fmt.Sprintf(examples, appName)

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

	var (
		iFName string
		oFName string
	)
	if len(args) > 0 {
		iFName = args[0]
	}
	if len(args) > 1 {
		oFName = args[1]
	}

	o := opml.New()
	if len(iFName) > 0 {
		err := o.ReadFile(iFName)
		if err != nil {
			fmt.Printf("%s, %s\n", iFName, err)
			os.Exit(1)
		}
	} else {
		src, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Missing OPML input, %s", err)
			os.Exit(1)
		}
		o, err = opml.Parse(src)
	}
	o.Sort()

	if prettyPrint == true {
		src, _ := xml.MarshalIndent(o, "  ", "    ")
		if oFName == "" {
			fmt.Printf("%s\n", src)
		} else {
			ioutil.WriteFile(oFName, src, 0664)
		}
	} else {
		if oFName == "" {
			fmt.Println(o.String())
		} else {
			o.WriteFile(oFName, 0664)
		}
	}
}
