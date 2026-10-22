package worktree

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// childOutput runs a real child through the exec path against worktree, with the caller's
// routing variables in the call's base environment, and returns what that child wrote.
func childOutput(t *testing.T, worktree, script string) string {
	t.Helper()
	return childOutputAtHome(t, worktree, script, Home())
}

// childOutputAtHome is childOutput with the home the verb resolved passed explicitly. The
// call's base environment carries a different caller home, so a row can separate the two.
// The caller values travel in the call and not in the process, so a test can run in parallel.
func childOutputAtHome(t *testing.T, worktree, script, home string) string {
	t.Helper()
	caller := []string{"PATH=" + os.Getenv("PATH"), homeEnv + "=/caller/home", "BENCH_KIT=/caller/kit",
		"BENCH_WRAPPER=/caller/bin/bench", runbinary.Env + "=/caller/run/bench", "BENCH_EXEC_TEST_CARRIED=carried-value"}
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"sh", "-c", script}, worktree, home, caller, nil, nil, &stdout, &stderr)
	requireTest(t, code == 0, "child exit = %d, stderr %q", code, stderr.String())
	return stdout.String()
}

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
