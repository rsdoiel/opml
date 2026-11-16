package opml

import (
	"os"
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
