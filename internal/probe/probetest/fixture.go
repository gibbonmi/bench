// Package probetest installs the canned `go` that probe tests run under. A probe starts one
// baseline `go test` child and then one mutated child for each mutant, so a test supplies one
// answer for each start. Only tests import this package; it holds no verdict, diagnostic, row,
// or record schema, and the caller authors every byte and exit it supplies.
package probetest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/sanitize"
)

// Start is what the canned `go` does on one `go test` start: the shell it runs first, the
// event stream it prints, and the exit it takes.
type Start struct {
	Hook   string
	Events string
	Exit   int
}

// unansweredExit is the exit of a `go test` start that InstallSequence has no answer for.
const unansweredExit = 97

// InstallStarts puts a canned `go` ahead of the real one for the single form. An odd test
// start answers baseline, and the even start after it answers mutated.
func InstallStarts(t testing.TB, marker string, baseline, mutated Start) {
	t.Helper()
	install(t, marker, "if [ $((n % 2)) = 1 ]; then\n"+answer(baseline)+"fi\n"+answer(mutated))
}

// InstallSequence puts a canned `go` ahead of the real one for a mutant set. The first test
// start answers baseline, and the k-th later start answers mutated[k-1]. A start past the list
// prints no event and exits unansweredExit.
func InstallSequence(t testing.TB, marker string, baseline Start, mutated ...Start) {
	t.Helper()
	arms := "1)\n" + answer(baseline) + ";;\n"
	for k, start := range mutated {
		arms += strconv.Itoa(k+2) + ")\n" + answer(start) + ";;\n"
	}
	install(t, marker, "case $n in\n"+arms+"esac\nexit "+strconv.Itoa(unansweredExit)+"\n")
}

// install writes the canned `go` for the rest of the test. Its prelude records every start in
// marker and passes each non-test call to the real toolchain. It reads the test start's
// ordinal into `$n` for body and its hooks. The marker is removed first, so each install
// restarts the count.
func install(t testing.TB, marker, body string) {
	t.Helper()
	real, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(marker); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	quoted := sanitize.ShellQuote(marker)
	script := "#!/usr/bin/env bash\n" +
		"printf '%s\\n' \"$1\" >> " + quoted + "\n" +
		"if [ \"$1\" != test ]; then exec " + sanitize.ShellQuote(real) + " \"$@\"; fi\n" +
		"n=$(grep -c '^test$' " + quoted + ")\n" + body
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func answer(start Start) string {
	return start.Hook + "cat <<'PROBEEOF'\n" + start.Events + "\nPROBEEOF\nexit " + strconv.Itoa(start.Exit) + "\n"
}
