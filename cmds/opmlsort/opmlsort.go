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

	"github.com/rsdoiel/opml"
)

var (
	showHelp    bool
	showVersion bool
	prettyPrint bool
)

func init() {
	flag.BoolVar(&prettyPrint, "p", false, "pretty print XML output")
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	if showHelp == true {
		fmt.Printf(` USAGE: %s [OPTIONS] INPUT_OPML_FILENAME [OUTPUT_OPML_FILENAME]
  This program sorts the outline in an OPML document.

  OPTIONS
`, appName)
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("\t-%s\t(defaults to %s) %s\n", f.Name, f.Value, f.Usage)
		})

		fmt.Printf("\nVersion %s\n", opml.Version)
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Printf("Version %s\n", opml.Version)
		os.Exit(0)
	}

	var (
		iFName string
		oFName string
	)

	args := flag.Args()
	if len(args) > 0 {
		iFName = args[0]
	}
	if len(args) > 1 {
		oFName = args[1]
	}

	if iFName == "" {
		fmt.Println("Missing input opml filename")
		os.Exit(1)
	}
	o := opml.New()
	err := o.ReadFile(iFName)
	if err != nil {
		fmt.Printf("%s, %s\n", iFName, err)
		os.Exit(1)
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
