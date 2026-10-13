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
// fixture's sealed executable pair to the output.
func buildEachRun(t *testing.T, f *fixture, rel string) string {
	t.Helper()
	records := scratchDir(t)
	record := "n=$(ls " + sanitize.ShellQuote(records) + " | wc -l)\n" +
		"cp \"$1\"/" + sanitize.ShellQuote(rel) + " " + sanitize.ShellQuote(records) + "/\"$n\"\n"
	installSealedBuild(t, f, record)
	return records
}

// parkSecondBuild replaces the fixture's inherited run binary with a build script. Its first
// call copies the sealed pair, and its second call writes the returned marker and parks.
func parkSecondBuild(t *testing.T, f *fixture) string {
	t.Helper()
	dir := scratchDir(t)
	first, parked := sanitize.ShellQuote(filepath.Join(dir, "first")), filepath.Join(dir, "parked")
	installSealedBuild(t, f, "if [ -e "+first+" ]; then : > "+sanitize.ShellQuote(parked)+"; sleep 60; exit 1; fi\n: > "+first+"\n")
	return parked
}

// installSealedBuild installs a build script that runs step and then copies the fixture's
// sealed executable pair to the output. The freshness verifier accepts the copy while the
// digest closure keeps its start bytes. TMPDIR moves, so each private directory is observable.
func installSealedBuild(t *testing.T, f *fixture, step string) {
	t.Helper()
	sealed := os.Getenv(runbinary.Env)
	script := "set -eu\n" +
		"if [ \"$1\" = --manifest-dir ]; then shift 2; fi\n" + step +
		"cp " + sanitize.ShellQuote(sealed) + " \"$2\"\n" +
		"cp " + sanitize.ShellQuote(sealed+".seal") + " \"$2.seal\"\n"
	installBuildScript(t, f, script)
	t.Setenv("TMPDIR", scratchDir(t))
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
