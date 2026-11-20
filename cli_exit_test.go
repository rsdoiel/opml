package opml

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// These tests run the five programs as built binaries, the way a user does,
// with real files from testdata. The golden files in testdata/cli are what the
// programs wrote before they moved into internal/cli.

var programNames = []string{"opml2json", "opml2urls", "opmlcat", "opmlsort", "urls2opml"}

var builtDir string

// binaries builds the five programs once per test run.
func binaries(t *testing.T) string {
	t.Helper()
	if builtDir != "" {
		return builtDir
	}
	dir, err := os.MkdirTemp("", "opml-cli-")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range programNames {
		if out, err := exec.Command("go", "build", "-o", filepath.Join(dir, name), "./cmd/"+name).CombinedOutput(); err != nil {
			t.Fatalf("%s: build failed: %v\n%s", name, err, out)
		}
	}
	builtDir = dir
	return dir
}

func TestMain(m *testing.M) {
	code := m.Run()
	if builtDir != "" {
		os.RemoveAll(builtDir)
	}
	os.Exit(code)
}

// run runs a built program with stdin and returns stdout, stderr and the exit code.
func run(t *testing.T, dir, name, stdin string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(filepath.Join(binaries(t), name), args...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(stdin)
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	code := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("%s: %v", name, err)
		}
		code = ee.ExitCode()
	}
	return out.String(), errOut.String(), code
}

func readFile(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

/** TestOutputMatchesGolden checks each program against the output it wrote
 * before the move to internal/cli, with the input named as an argument and on
 * standard input.
 */
func TestOutputMatchesGolden(t *testing.T) {
	ex1, simple1, toc1 := "testdata/example1.opml", "testdata/simple1.opml", "testdata/toc1.opml"
	cases := []struct {
		prog, gold, input string
		args              []string
	}{
		{"opml2json", "opml2json-plain", ex1, nil},
		{"opml2json", "opml2json-pretty", ex1, []string{"-pretty", "-newline"}},
		{"opml2json", "opml2json-pretty", ex1, []string{"-pretty", "-nl"}},
		{"opml2urls", "opml2urls-default", ex1, nil},
		{"opml2urls", "opml2urls-newsboat", ex1, []string{"-newsboat"}},
		{"opml2urls", "opml2urls-htmlurl", ex1, []string{"-htmlurl"}},
		// A table of contents has url attributes, not xmlUrl: with no feed
		// URLs in the document the url attribute is listed without asking.
		{"opml2urls", "opml2urls-toc-default", toc1, nil},
		{"opml2urls", "opml2urls-toc-urls", toc1, []string{"-text-as-comment=false"}},
		{"opml2urls", "opml2urls-toc-urls", toc1, []string{"-url", "-text-as-comment=false"}},
		{"opml2urls", "opml2urls-toc-includes", toc1, []string{"-includes", "-text-as-comment=false"}},
		{"opmlsort", "opmlsort-plain", ex1, nil},
		{"opmlsort", "opmlsort-pretty", ex1, []string{"-pretty", "-case-insensitive", "-newline"}},
		{"opmlcat", "opmlcat-stdin", simple1, nil},
		{"opmlcat", "opmlcat-pretty", simple1, []string{"-pretty"}},
	}
	for _, c := range cases {
		want := readFile(t, "testdata/cli/"+c.gold+".golden")
		abs, _ := filepath.Abs(c.input)
		// The input as an argument (opmlcat names every input).
		out, errOut, code := run(t, t.TempDir(), c.prog, "", append(append([]string{}, c.args...), abs)...)
		if code != 0 || out != want {
			t.Errorf("%s %v %s: exit %d, stderr %q, output differs from %s (got %d bytes, want %d)", c.prog, c.args, c.input, code, errOut, c.gold, len(out), len(want))
		}
		// The input on standard input.
		out, errOut, code = run(t, t.TempDir(), c.prog, readFile(t, c.input), c.args...)
		if code != 0 || out != want {
			t.Errorf("%s %v from standard input: exit %d, stderr %q, output differs from %s", c.prog, c.args, code, errOut, c.gold)
		}
		// -i and -o, or for opmlcat the positional input and -o.
		dir := t.TempDir()
		outFile := filepath.Join(dir, "out.txt")
		args := append(append([]string{}, c.args...), "-o", outFile)
		if c.prog == "opmlcat" {
			args = append(args, abs)
		} else {
			args = append(args, "-i", abs)
		}
		if out, errOut, code = run(t, dir, c.prog, "", args...); code != 0 || out != "" || readFile(t, outFile) != want {
			t.Errorf("%s %v -i -o: exit %d, stdout %q, stderr %q, or the file differs from %s", c.prog, c.args, code, out, errOut, c.gold)
		}
	}
}

func TestUrls2opmlMatchesGolden(t *testing.T) {
	want := readFile(t, "testdata/cli/urls2opml.golden")
	filter := func(s string) string {
		var keep []string
		for _, line := range strings.Split(s, "\n") {
			if !strings.Contains(line, "<title>") && !strings.Contains(line, "dateCreated") {
				keep = append(keep, line)
			}
		}
		return strings.Join(keep, "\n")
	}
	urls := readFile(t, "testdata/cli/urls.txt")
	abs, _ := filepath.Abs("testdata/cli/urls.txt")
	for name, args := range map[string][]string{"stdin": nil, "positional": {abs}, "-i": {"-i", abs}} {
		stdin := ""
		if name == "stdin" {
			stdin = urls
		}
		out, errOut, code := run(t, t.TempDir(), "urls2opml", stdin, args...)
		if code != 0 || filter(out) != want {
			t.Errorf("urls2opml %s: exit %d, stderr %q, output differs", name, code, errOut)
		}
	}
}

/** TestOpmlcatKeepsItsInputs is the bug found while converting: with two
 * input files opmlcat treated the second as an output name and created the
 * first one, which emptied it.
 */
func TestOpmlcatKeepsItsInputs(t *testing.T) {
	dir := t.TempDir()
	a, b := filepath.Join(dir, "a.opml"), filepath.Join(dir, "b.opml")
	for _, f := range [][2]string{{a, "testdata/simple1.opml"}, {b, "testdata/simple2.opml"}} {
		if err := os.WriteFile(f[0], []byte(readFile(t, f[1])), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out, errOut, code := run(t, dir, "opmlcat", "", a, b)
	if code != 0 {
		t.Fatalf("opmlcat a b: exit %d, stderr %q", code, errOut)
	}
	if readFile(t, a) != readFile(t, "testdata/simple1.opml") || readFile(t, b) != readFile(t, "testdata/simple2.opml") {
		t.Errorf("opmlcat a b changed an input file")
	}
	one, _, _ := run(t, dir, "opmlcat", "", a)
	if len(out) <= len(one) {
		t.Errorf("opmlcat a b wrote %d bytes, no more than opmlcat a (%d): the second input is missing", len(out), len(one))
	}
}

/** TestExitCodes covers DR-0003 for every program: a missing input is 66, an
 * input that is not OPML is 65, an output that cannot be created is 73, a bad
 * command line is 2, and the output file is not created when the input fails.
 */
func TestExitCodes(t *testing.T) {
	for _, name := range programNames {
		dir := t.TempDir()
		missing := filepath.Join(dir, "missing.opml")
		bad := filepath.Join(dir, "bad.opml")
		if err := os.WriteFile(bad, []byte("this is not xml <"), 0o644); err != nil {
			t.Fatal(err)
		}
		good, _ := filepath.Abs("testdata/simple1.opml")
		if name == "urls2opml" {
			good, _ = filepath.Abs("testdata/cli/urls.txt")
		}
		noDir := filepath.Join(dir, "no", "such", "dir", "out")
		cases := []struct {
			what  string
			stdin string
			args  []string
			want  int
		}{
			{"unknown option", "", []string{"-nosuchoption"}, 2},
			{"surplus argument", "", []string{good, filepath.Join(dir, "o1"), filepath.Join(dir, "o2")}, 2},
			{"missing input file", "", []string{missing}, 66},
			{"missing input with -i", "", []string{"-i", missing}, 66},
			{"output cannot be created", "", []string{"-i", good, "-o", noDir}, 73},
		}
		if name == "opmlcat" {
			cases = cases[:1]
			cases = append(cases,
				struct {
					what  string
					stdin string
					args  []string
					want  int
				}{"missing input file", "", []string{missing}, 66},
				struct {
					what  string
					stdin string
					args  []string
					want  int
				}{"output cannot be created", "", []string{"-o", noDir, good}, 73})
		}
		if name != "urls2opml" {
			cases = append(cases,
				struct {
					what  string
					stdin string
					args  []string
					want  int
				}{"input is not OPML", "", []string{bad}, 65},
				struct {
					what  string
					stdin string
					args  []string
					want  int
				}{"standard input is not OPML", "not xml <", nil, 65})
		}
		for _, c := range cases {
			out, errOut, code := run(t, dir, name, c.stdin, c.args...)
			if code != c.want {
				t.Errorf("%s, %s: exit %d, want %d (stdout %.60q, stderr %q)", name, c.what, code, c.want, out, errOut)
			}
			if c.want != 0 && !strings.HasPrefix(errOut, name+": ") {
				t.Errorf("%s, %s: stderr %q does not start with the program name", name, c.what, errOut)
			}
		}
		// A failing input does not create or empty the output file.
		if name != "urls2opml" {
			outFile := filepath.Join(dir, "keep.txt")
			if err := os.WriteFile(outFile, []byte("keep"), 0o644); err != nil {
				t.Fatal(err)
			}
			args := []string{"-o", outFile, bad}
			if name != "opmlcat" {
				args = []string{"-i", bad, "-o", outFile}
			}
			run(t, dir, name, "", args...)
			if readFile(t, outFile) != "keep" {
				t.Errorf("%s: a bad input changed the existing output file", name)
			}
		}
		// -quiet drops the message, not the exit code.
		out, errOut, code := run(t, dir, name, "", "-quiet", "-i", missing)
		if name == "opmlcat" {
			out, errOut, code = run(t, dir, name, "", "-quiet", missing)
		}
		if code != 66 || errOut != "" || out != "" {
			t.Errorf("%s -quiet: exit %d, stdout %q, stderr %q; want 66 and nothing printed", name, code, out, errOut)
		}
	}
}
