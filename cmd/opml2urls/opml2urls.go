// opml2urls is part of opml, a Go package for OPML. Run opml2urls -help.
package main

import (
	"os"
	"path"

	"github.com/rsdoiel/opml/internal/cli"
)

func main() {
	os.Exit(cli.Opml2urls(path.Base(os.Args[0]), os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
