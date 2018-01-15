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
	"fmt"
	"io/ioutil"
	"os"

	// My Packages
	"github.com/rsdoiel/opml"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
)

var (
	description = `
%s concatenates one or more opml files as siblings to standard out.
`

	examples = `
This is an example of using %s and opmlsort together to 
create a combined sorted opml file.

    %s file1.opml file1.opml | opmlsort -o combined-sorted.opml
`

	// Standard options
	showHelp             bool
	showVersion          bool
	showLicense          bool
	showExamples         bool
	inputFName           string
	outputFName          string
	quiet                bool
	newLine              bool
	generateMarkdownDocs bool

	// Application options
	prettyPrint bool
)

func main() {
	app := cli.NewCli(opml.Version)
	appName := app.AppName()

	// Add non-option parameter docs
	app.AddParams("OPML_FILE", "[OPML_FILE ...]")

	// Add Help docs
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName)))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display examples")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&newLine, "nl,newline", false, "add a trailing newline")
	app.StringVar(&inputFName, "i,input", "", "set input filename")
	app.StringVar(&outputFName, "o,output", "", "set output filename")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "generate Markdown documentation")

	// Application Options
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print XML output")

	// Process environment and options
	app.Parse()
	args := app.Args()

	// Setup I/O
	var err error

	app.Eout = os.Stderr
	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle options
	if generateMarkdownDocs {
		app.GenerateMarkdownDocs(os.Stdout)
		os.Exit(0)
	}
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintln(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	o := opml.New()
	if len(args) == 0 {
		src, err := ioutil.ReadAll(app.In)
		cli.ExitOnError(app.Eout, err, quiet)

		o, err = opml.Parse(src)
		cli.ExitOnError(app.Eout, err, quiet)
	}

	for _, inputFName := range args {
		next := opml.New()
		err := next.ReadFile(inputFName)
		cli.ExitOnError(app.Eout, err, quiet)

		err = o.Append(next)
		cli.ExitOnError(app.Eout, err, quiet)
	}

	var src []byte

	if prettyPrint == true {
		src, err = xml.MarshalIndent(o, "", "    ")
		cli.ExitOnError(app.Eout, err, quiet)
	} else {
		src = []byte(o.String())
	}

	fmt.Fprintln(app.Out, `<?xml version="1.0" encoding="UTF-8"?>`)
	if newLine {
		fmt.Fprintf(app.Out, "%s\n", src)
	} else {
		fmt.Fprintf(app.Out, "%s", src)
	}
}
