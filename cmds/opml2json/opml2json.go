//
// opml2json is a command line utility that can read in a OPML file and
// return it in JSON format.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2016 all rights reserved.
// Released under the BSD 3-Clause License
// See: https://opensource.org/licenses/BSD-3-Clause
//
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	// My Packages
	"github.com/rsdoiel/opml"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
)

var (
	description = `
%s is a program that converts OPML's XML to JSON.
`

	examples = `
Convert *myfeeds.ompl* to *myfeeds.json*.

    %s myfeeds.opml myfeeds.json
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

	// Document non-option parameters
	app.AddParams("INPUT_OPML_FILENAME", "[OUTPUT_OPML_FILENAME]")

	// Add Help Docs
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName)))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display examples")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&newLine, "nl,newline", false, "add trailing newline")
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
	if len(inputFName) > 0 {
		err := o.ReadFile(inputFName)
		cli.ExitOnError(app.Eout, err, quiet)
	} else {
		src, err := ioutil.ReadAll(app.In)
		cli.ExitOnError(app.Eout, err, quiet)
		o, err = opml.Parse(src)
		cli.ExitOnError(app.Eout, err, quiet)
	}
	o.Sort()

	var src []byte
	if prettyPrint {
		src, err = json.MarshalIndent(o, "", "    ")
		cli.ExitOnError(app.Eout, err, quiet)
	} else {
		src, err = json.Marshal(o, "", "    ")
		cli.ExitOnError(app.Eout, err, quiet)
	}

	if newLine {
		fmt.Fprintf(app.Out, "%s\n", src)
	} else {
		fmt.Fprintf(app.Out, "%s", src)
	}
}
