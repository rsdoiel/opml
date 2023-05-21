//
// opml2json is a command line utility that can read in a OPML file and
// return it in JSON format.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
//
// Copyright (c) 2021, R. S. Doiel
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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	// My Packages
	"github.com/rsdoiel/opml"
)

var (
	helpText = `%{app_name}(1) | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] [INPUT_FILENAME] [OUTPUT_FILENAME]

# DESCRIPTION

{app_name} is a program that converts OPML's XML to JSON.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-i
: read from filename

-o
: write to filename

-newline
: add a trailing newline

-pretty
: pretty print JSON output


# EXAMPLES

Convert *myfeeds.ompl* to *myfeeds.json*.

~~~
{app_name} myfeeds.opml myfeeds.json
~~~

`

	// Standard options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	inputFName       string
	outputFName      string
	quiet            bool
	newLine          bool

	// Application options
	prettyPrint bool
)

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: the following are set when version.go is generated
	version := opml.Version
	releaseDate := opml.ReleaseDate
	releaseHash := opml.ReleaseHash
	fmtHelp := opml.FmtHelp

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&quiet, "quiet", false, "suppress error messages")
	flag.BoolVar(&newLine, "newline", false, "add trailing newline")
	flag.StringVar(&inputFName, "i", "", "set input filename")
	flag.StringVar(&outputFName, "o", "", "set output filename")

	// Application Options
	flag.BoolVar(&prettyPrint, "pretty", false, "pretty print XML output")

	// Process environment and options
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		inputFName = args[0]
	}
	if len(args) > 1 {
		outputFName = args[1]
	}

	// Setup I/O
	var err error

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	// Handle options
	if showHelp {
		fmt.Fprint(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(out, "%s\n", opml.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(out, "%s %s %s\n", appName, version, releaseHash)
		os.Exit(0)
	}

	if inputFName != "" {
		in, err = os.Open(inputFName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		defer in.Close()
	}

	if outputFName != "" {
		out, err = os.Create(outputFName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}


	o := opml.New()
	if len(inputFName) > 0 {
		if err := o.ReadFile(inputFName); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
	} else {
		src, err := ioutil.ReadAll(in)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		o, err = opml.Parse(src)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
	}
	o.Sort()

	var src []byte
	if prettyPrint {
		src, err = json.MarshalIndent(o, "", "    ")
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
	} else {
		src, err = json.Marshal(o)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
	}

	fmt.Fprintf(out, "%s", src)
	if newLine {
		fmt.Fprintln(out)
	}
}
