// Package cli holds the command line behavior of opml2json, opml2urls, opmlcat,
// opmlsort and urls2opml so that plain go test covers it. The programs under
// cmd/ only call Opml2json, Opml2urls, Opmlcat, Opmlsort and Urls2opml.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/rsdoiel/opml"
)

// The exit codes of every command line tool in the workspace (DR-0003).
const (
	ExitOK         = 0
	ExitUsage      = 2
	ExitData       = 65
	ExitNoInput    = 66
	ExitInternal   = 70
	ExitCantCreate = 73
	ExitIO         = 74
	ExitNoPerm     = 77
)

// options are the options every program takes.
type options struct {
	help, license, version bool
	newline, quiet         bool
	input, output          string
}

// newFlagSet makes a FlagSet that reports its errors on stderr and registers
// the standard options, each in a short and a long spelling. newlineDefault is
// whether a trailing newline is written unless -newline=false is given.
func newFlagSet(appName string, o *options, newlineDefault bool) *flag.FlagSet {
	fset := flag.NewFlagSet(appName, flag.ContinueOnError)
	fset.SetOutput(io.Discard)
	fset.BoolVar(&o.help, "h", false, "display help")
	fset.BoolVar(&o.help, "help", false, "display help")
	fset.BoolVar(&o.license, "l", false, "display license")
	fset.BoolVar(&o.license, "license", false, "display license")
	fset.BoolVar(&o.version, "v", false, "display version")
	fset.BoolVar(&o.version, "version", false, "display version")
	fset.BoolVar(&o.newline, "nl", newlineDefault, "add a trailing newline")
	fset.BoolVar(&o.newline, "newline", newlineDefault, "add a trailing newline")
	fset.BoolVar(&o.quiet, "quiet", false, "suppress error messages")
	fset.StringVar(&o.input, "i", "", "input filename")
	fset.StringVar(&o.input, "input", "", "input filename")
	fset.StringVar(&o.output, "o", "", "output filename")
	fset.StringVar(&o.output, "output", "", "output filename")
	return fset
}

// fail writes one plain line, prefixed with the program name, to stderr unless
// the program was run with -quiet, and returns code (DR-0011).
func fail(appName string, o *options, stderr io.Writer, code int, format string, a ...any) int {
	if !o.quiet {
		fmt.Fprintf(stderr, "%s: %s\n", appName, fmt.Sprintf(format, a...))
	}
	return code
}

// warn writes one plain warning line to stderr unless -quiet was given.
func warn(appName string, o *options, stderr io.Writer, format string, a ...any) {
	if !o.quiet {
		fmt.Fprintf(stderr, "%s: %s\n", appName, fmt.Sprintf(format, a...))
	}
}

// classify turns an error from reading a named input into an exit code.
func classify(err error) int {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return ExitNoInput
	case errors.Is(err, fs.ErrPermission):
		return ExitNoPerm
	default:
		return ExitIO
	}
}

// start parses args, answers -help, -license and -version, and works out the
// input and output names. With maxPositional 2 they are the options or the
// first and second positional arguments; with 0 every positional argument is
// left in the returned slice for the program (opmlcat names many inputs). It
// returns done true with the exit code when the program is to stop there.
func start(appName, helpText string, fset *flag.FlagSet, o *options, args []string, maxPositional int, stdout, stderr io.Writer) (rest []string, code int, done bool) {
	if err := fset.Parse(args); err != nil {
		// The message is the flag package's, without a trailing newline or usage.
		return nil, fail(appName, o, stderr, ExitUsage, "%v (try %s -help)", err, appName), true
	}
	switch {
	case o.help:
		fmt.Fprintln(stdout, opml.FmtHelp(helpText, appName, opml.Version, opml.ReleaseDate, opml.ReleaseHash))
		return nil, ExitOK, true
	case o.license:
		fmt.Fprintln(stdout, strings.TrimSpace(opml.LicenseText))
		return nil, ExitOK, true
	case o.version:
		fmt.Fprintf(stdout, "%s %s %s\n", appName, opml.Version, opml.ReleaseHash)
		return nil, ExitOK, true
	}
	rest = fset.Args()
	if maxPositional == 0 {
		return rest, ExitOK, false
	}
	if len(rest) > maxPositional {
		return nil, fail(appName, o, stderr, ExitUsage, "surplus argument %q (try %s -help)", rest[maxPositional], appName), true
	}
	if o.input == "" && len(rest) > 0 {
		o.input = rest[0]
	}
	if o.output == "" && len(rest) > 1 {
		o.output = rest[1]
	}
	return nil, ExitOK, false
}

// namedInput says whether an input name is a file, not standard input.
func namedInput(name string) bool {
	return name != "" && name != "-"
}

// readSource returns the bytes of the named file, or of stdin when name is
// empty or "-". On failure it returns the exit code.
func readSource(appName string, o *options, name string, stdin io.Reader, stderr io.Writer) ([]byte, int) {
	if namedInput(name) {
		src, err := os.ReadFile(name)
		if err != nil {
			return nil, fail(appName, o, stderr, classify(err), "%v", err)
		}
		return src, ExitOK
	}
	src, err := io.ReadAll(stdin)
	if err != nil {
		return nil, fail(appName, o, stderr, ExitIO, "%v", err)
	}
	return src, ExitOK
}

// readOPML reads and parses one OPML document from the named file or stdin.
// A file that is not OPML is a data error (65), a missing one 66.
func readOPML(appName string, o *options, name string, stdin io.Reader, stderr io.Writer) (*opml.OPML, int) {
	src, code := readSource(appName, o, name, stdin, stderr)
	if code != ExitOK {
		return nil, code
	}
	doc, err := opml.Parse(src)
	if err != nil {
		if namedInput(name) {
			return nil, fail(appName, o, stderr, ExitData, "%s: %v", name, err)
		}
		return nil, fail(appName, o, stderr, ExitData, "%v", err)
	}
	return doc, ExitOK
}

// write sends the result to the named output file, or to stdout when there is
// none or it is "-". The file is created only now, after the input has been
// read and converted, so a run that fails leaves an existing file alone. A
// trailing newline is added when the options ask for one.
func write(appName string, o *options, result []byte, stdout, stderr io.Writer) int {
	if o.newline {
		result = append(result, '\n')
	}
	if o.output == "" || o.output == "-" {
		if _, err := stdout.Write(result); err != nil {
			return fail(appName, o, stderr, ExitIO, "%v", err)
		}
		return ExitOK
	}
	f, err := os.Create(o.output)
	if err != nil {
		return fail(appName, o, stderr, ExitCantCreate, "%v", err)
	}
	if _, err := f.Write(result); err != nil {
		f.Close()
		return fail(appName, o, stderr, ExitIO, "%v", err)
	}
	if err := f.Close(); err != nil {
		return fail(appName, o, stderr, ExitIO, "%v", err)
	}
	return ExitOK
}
