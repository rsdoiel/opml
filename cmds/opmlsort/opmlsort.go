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
	"flag"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/rsdoiel/opml"
)

const Version = "0.0.0"

var (
	showHelp    bool
	showVersion bool
)

func init() {
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

		fmt.Printf("\nVersion %s\n", Version)
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Printf("Version %s\n", Version)
		os.Exit(0)
	}

	var (
		iFName string
		oFName string
	)

	args := flag.Args()
	fmt.Printf("DEBUG args: %+v\n", args)
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
	sort.Sort(opml.ByText(o.Body.Outline))

	if oFName == "" {
		fmt.Println(o.String())
	} else {
		o.WriteFile(oFName)
	}
}
