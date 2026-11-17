// urls2opml is part of opml, a Go package for OPML. Run urls2opml -help.
package main

import (
	"os"
	"path"

	"github.com/rsdoiel/opml/internal/cli"
)

func main() {
	os.Exit(cli.Urls2opml(path.Base(os.Args[0]), os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
