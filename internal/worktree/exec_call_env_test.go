package worktree

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestExecChildTakesTheCallEnvironment proves that an exec child starts from the base
// environment of the call. The child sees the call's own value and no value that only
// the test process holds. A BENCH_ variable, PATH, and PWD are no witness, because the
// child environment strips, sets, or shares them. (Coverage row LGP19.)
func TestExecChildTakesTheCallEnvironment(t *testing.T) {
	t.Parallel()
	processOnly := ""
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if name != "PATH" && name != "PWD" && !strings.HasPrefix(name, "BENCH_") && !strings.Contains(entry, "\n") {
			processOnly = entry
			break
		}
	}
	requireTest(t, processOnly != "", "the test process holds no variable outside the call environment")
	call := []string{"PATH=" + os.Getenv("PATH"), "EXEC_CALL_ENV=from-call"}
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"env"}, t.TempDir(), t.TempDir(), call, nil, nil, &stdout, &stderr)
	requireTest(t, code == 0, "child exit = %d, stderr %q", code, stderr.String())
	seen := "\n" + stdout.String()
	requireTest(t, strings.Contains(seen, "\nEXEC_CALL_ENV=from-call\n"), "the child environment lacks the call value:\n%s", stdout.String())
	requireTest(t, !strings.Contains(seen, "\n"+processOnly+"\n"), "the child environment holds the process-only value %q", processOnly)
}
