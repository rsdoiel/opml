package opml

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

/** TestCommandsInformationalOptions builds each command and runs it with
 * -help, -license and -version, the way a user does, and checks that the
 * output holds the real text. A Fprintln or Fprint given a "%s\n" format
 * prints the directive itself, not the value.
 */
func TestCommandsInformationalOptions(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"opml2json", "opml2urls", "opmlcat", "opmlsort", "urls2opml"} {
		bin := filepath.Join(dir, name)
		if out, err := exec.Command("go", "build", "-o", bin, "./cmd/"+name).CombinedOutput(); err != nil {
			t.Fatalf("%s: build failed: %v\n%s", name, err, out)
		}
		for opt, want := range map[string]string{
			"-help":    name,
			"-license": "GNU AFFERO GENERAL PUBLIC LICENSE",
			"-version": name + " " + Version,
		} {
			out, err := exec.Command(bin, opt).Output()
			if err != nil {
				t.Errorf("%s %s: %v", name, opt, err)
				continue
			}
			s := string(out)
			if strings.Contains(s, "%s") {
				t.Errorf("%s %s printed a literal %%s directive:\n%.200s", name, opt, s)
			}
			if !strings.Contains(s, want) {
				t.Errorf("%s %s: output lacks %q:\n%.200s", name, opt, want, s)
			}
		}
	}
}
