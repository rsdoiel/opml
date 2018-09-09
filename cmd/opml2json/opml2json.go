//
// opml2json is a command line utility that can read in a OPML file and
// return it in JSON format.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
//
// Copyright (c) 2018, R. S. Doiel
// All rights not granted herein are expressly reserved by R. S. Doiel.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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

	license = `
%s %s

Copyright (c) 2018, R. S. Doiel
All rights not granted herein are expressly reserved by R. S. Doiel.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	// Standard options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	showExamples     bool
	inputFName       string
	outputFName      string
	quiet            bool
	newLine          bool
	generateMarkdown bool
	generateManPage  bool

	// Application options
	prettyPrint bool
)

func main() {
	app := cli.NewCli(opml.Version)
	appName := app.AppName()

	// Document non-option parameters
	app.SetParams("INPUT_OPML_FILENAME", "[OUTPUT_OPML_FILENAME]")

	// Add Help Docs
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName)))
	app.AddHelp("license", []byte(fmt.Sprintf(license, appName, opml.Version)))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display examples")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&newLine, "nl,newline", false, "add trailing newline")
	app.StringVar(&inputFName, "i,input", "", "set input filename")
	app.StringVar(&outputFName, "o,output", "", "set output filename")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// Application Options
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print XML output")

	// Process environment and options
	app.Parse()
	args := app.Args()

	if len(args) > 0 {
		inputFName = args[0]
	}
	if len(args) > 1 {
		outputFName = args[1]
	}

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
	if generateMarkdown {
		app.GenerateMarkdown(os.Stdout)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(os.Stdout)
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
		src, err = json.Marshal(o)
		cli.ExitOnError(app.Eout, err, quiet)
	}

	if newLine {
		fmt.Fprintf(app.Out, "%s\n", src)
	} else {
		fmt.Fprintf(app.Out, "%s", src)
	}
}
