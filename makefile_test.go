package opml

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

/** TestMakefileProgramsMatchManPages checks that the Makefile's PROGRAMS line
 * names every cmd/ directory and every section 1 man page source. The man page
 * list is derived from PROGRAMS, so a program left out would drop its page from
 * `make man` and `make install` while the build still succeeds.
 */
func TestMakefileProgramsMatchManPages(t *testing.T) {
	mk, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?m)^PROGRAMS *= *(.*)$`).FindSubmatch(mk)
	if m == nil {
		t.Fatal("Makefile has no PROGRAMS line")
	}
	got := strings.Fields(string(m[1]))
	var cmds, pages []string
	dirs, _ := filepath.Glob("cmd/*")
	for _, d := range dirs {
		cmds = append(cmds, filepath.Base(d))
	}
	md, _ := filepath.Glob("*.1.md")
	for _, f := range md {
		pages = append(pages, strings.TrimSuffix(f, ".1.md"))
	}
	sort.Strings(got)
	sort.Strings(cmds)
	sort.Strings(pages)
	if strings.Join(got, " ") != strings.Join(cmds, " ") {
		t.Errorf("Makefile PROGRAMS = %v, cmd/ holds %v", got, cmds)
	}
	if strings.Join(got, " ") != strings.Join(pages, " ") {
		t.Errorf("Makefile PROGRAMS = %v, *.1.md pages = %v", got, pages)
	}
}

/** TestCheckCommitted runs `make check-committed` in a scratch repository with
 * a bare origin: it passes when the tree is clean and HEAD is pushed, and fails
 * for an uncommitted change to a tracked file and for an unpushed commit.
 */
func TestCheckCommitted(t *testing.T) {
	for _, tool := range []string{"git", "make"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not available", tool)
		}
	}
	mk, err := filepath.Abs("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	origin, work := filepath.Join(root, "origin.git"), filepath.Join(root, "work")
	run := func(dir string, name string, args ...string) error {
		c := exec.Command(name, args...)
		c.Dir = dir
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.org",
			"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.org")
		if b, err := c.CombinedOutput(); err != nil {
			t.Logf("%s %v: %s", name, args, b)
			return err
		}
		return nil
	}
	must := func(dir, name string, args ...string) {
		t.Helper()
		if err := run(dir, name, args...); err != nil {
			t.Fatalf("%s %v: %v", name, args, err)
		}
	}
	must(root, "git", "init", "-q", "--bare", origin)
	must(root, "git", "init", "-q", "-b", "main", work)
	must(work, "git", "remote", "add", "origin", origin)
	if err := os.WriteFile(filepath.Join(work, "a.txt"), []byte("one\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	must(work, "git", "add", "a.txt")
	must(work, "git", "commit", "-q", "-m", "one")
	must(work, "git", "push", "-q", "origin", "main")
	check := func() error { return run(work, "make", "-s", "-f", mk, "check-committed") }
	if err := check(); err != nil {
		t.Errorf("clean and pushed: want pass, got %v", err)
	}
	if err := os.WriteFile(filepath.Join(work, "a.txt"), []byte("two\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := check(); err == nil {
		t.Error("uncommitted change to a tracked file: want failure")
	}
	must(work, "git", "commit", "-q", "-am", "two")
	if err := check(); err == nil {
		t.Error("committed but unpushed: want failure")
	}
	must(work, "git", "push", "-q", "origin", "main")
	if err := check(); err != nil {
		t.Errorf("pushed again: want pass, got %v", err)
	}
}

/** TestWebsiteMakPages checks that website.mak lists the same PROGRAMS as the
 * Makefile and that every page in SITE_PAGES has a Markdown source, so a
 * renamed or missing file fails here instead of silently dropping a page.
 */
func TestWebsiteMakPages(t *testing.T) {
	read := func(name string) string {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	programs := func(s string) string {
		m := regexp.MustCompile(`(?m)^PROGRAMS *= *(.*)$`).FindStringSubmatch(s)
		if m == nil {
			t.Fatal("no PROGRAMS line")
		}
		return m[1]
	}
	web := read("website.mak")
	if a, b := programs(read("Makefile")), programs(web); a != b {
		t.Errorf("Makefile PROGRAMS = %q, website.mak PROGRAMS = %q", a, b)
	}
	m := regexp.MustCompile(`(?m)^SITE_PAGES *= *((?:.*\\\n)*.*)$`).FindStringSubmatch(web)
	if m == nil {
		t.Fatal("website.mak has no SITE_PAGES")
	}
	list := strings.NewReplacer("\\\n", " ", "\t", " ").Replace(m[1])
	var pages []string
	for _, f := range strings.Fields(list) {
		if f == "$(addsuffix" || f == ".1,$(PROGRAMS))" {
			continue
		}
		pages = append(pages, f)
	}
	for _, p := range strings.Fields(programs(web)) {
		pages = append(pages, p+".1")
	}
	for _, p := range pages {
		if _, err := os.Stat(p + ".md"); err != nil {
			t.Errorf("website.mak lists %s, no %s.md: %v", p, p, err)
		}
	}
	if len(pages) < 8 {
		t.Errorf("parsed only %v from SITE_PAGES", pages)
	}
}
