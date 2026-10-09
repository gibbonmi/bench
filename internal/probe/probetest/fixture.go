// Package probetest installs the canned `go` that probe tests run under. A probe starts two
// `go test` children, the baseline and then the mutated run, so a test supplies one answer
// for each start. Only tests import this package; it holds no verdict, diagnostic, row, or
// record schema, and the caller authors every byte and exit it supplies.
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

// InstallStarts puts a canned `go` ahead of the real one for the rest of the test and
// records every start in marker. Non-test calls reach the real toolchain for settings and
// module discovery. The marker is removed first, so each install restarts the count: an
// odd test start answers baseline and the even start after it answers mutated. Both hooks
// read the start's ordinal in `$n`.
func InstallStarts(t testing.TB, marker string, baseline, mutated Start) {
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
		"n=$(grep -c '^test$' " + quoted + ")\n" +
		"if [ $((n % 2)) = 1 ]; then\n" + answer(baseline) + "fi\n" + answer(mutated)
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func answer(start Start) string {
	return start.Hook + "cat <<'PROBEEOF'\n" + start.Events + "\nPROBEEOF\nexit " + strconv.Itoa(start.Exit) + "\n"
}
