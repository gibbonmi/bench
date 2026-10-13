package probe

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// buildEachRun replaces the fixture's inherited run binary with a build script, so each
// focused run owns its selection. The script copies the subject named by rel into the
// returned records directory, one numbered file for each build, and then copies the
// fixture's sealed executable pair to the output. The freshness verifier accepts the copy,
// because the bytes and the adjacent seal are unchanged, the subject is outside the
// digest closure, and the private directory holds no symlink component.
func buildEachRun(t *testing.T, f *fixture, rel string) string {
	t.Helper()
	sealed := os.Getenv(runbinary.Env)
	records := scratchDir(t)
	script := "set -eu\n" +
		"if [ \"$1\" = --manifest-dir ]; then shift 2; fi\n" +
		"n=$(ls " + sanitize.ShellQuote(records) + " | wc -l)\n" +
		"cp \"$1\"/" + sanitize.ShellQuote(rel) + " " + sanitize.ShellQuote(records) + "/\"$n\"\n" +
		"cp " + sanitize.ShellQuote(sealed) + " \"$2\"\n" +
		"cp " + sanitize.ShellQuote(sealed+".seal") + " \"$2.seal\"\n"
	installBuildScript(t, f, script)
	t.Setenv("TMPDIR", scratchDir(t))
	return records
}

// installBuildScript writes script as the fixture's run-binary build script and removes the
// inherited run binary, so each run builds through the script.
func installBuildScript(t *testing.T, f *fixture, script string) {
	t.Helper()
	writeFixtureFile(t, filepath.Join(f.root, "scripts", "go-build.sh"), script, 0o755)
	// ReuseOrOwn tests the presence of the variable, so an empty value would still inherit.
	t.Setenv(runbinary.Env, "")
	if err := os.Unsetenv(runbinary.Env); err != nil {
		t.Fatal(err)
	}
}

// buildRecords answers the subject bytes that each build read, in build order.
func buildRecords(t *testing.T, records string) []string {
	t.Helper()
	entries, err := os.ReadDir(records)
	if err != nil {
		t.Fatal(err)
	}
	read := make([]string, len(entries))
	for i := range entries {
		body, err := os.ReadFile(filepath.Join(records, strconv.Itoa(i)))
		if err != nil {
			t.Fatal(err)
		}
		read[i] = string(body)
	}
	return read
}

// trackFixtureFile writes rel and adds it to the fixture's Git index. The prose grade of a
// work-tree top grades only the files that the index lists.
func trackFixtureFile(t *testing.T, f *fixture, rel, body string) {
	t.Helper()
	writeFixtureFile(t, filepath.Join(f.root, rel), body, 0o644)
	if out, err := exec.Command("git", "-C", f.root, "add", rel).CombinedOutput(); err != nil {
		t.Fatalf("git add %s: %v\n%s", rel, err, out)
	}
}
