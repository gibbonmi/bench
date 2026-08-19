package worktree

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// childEnvironment runs a real child through the exec path and returns the environment
// that child actually received, one assignment per line.
func childEnvironment(t *testing.T) string {
	t.Helper()
	t.Setenv("BENCH_KIT", "/caller/kit")
	t.Setenv("BENCH_WRAPPER", "/caller/bin/bench")
	t.Setenv(runbinary.Env, "/caller/run/bench")
	t.Setenv("BENCH_EXEC_TEST_CARRIED", "carried-value")
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"sh", "-c", "env"}, t.TempDir(), nil, &stdout, &stderr)
	requireTest(t, code == 0, "child exit = %d, stderr %q", code, stderr.String())
	return stdout.String()
}

func TestExecChildDropsWrapperRouting(t *testing.T) {
	seen := childEnvironment(t)
	for _, name := range []string{"BENCH_KIT", "BENCH_WRAPPER", runbinary.Env} {
		requireTest(t, !strings.Contains(seen, name+"="), "child saw %s in its environment:\n%s", name, seen)
	}
}

func TestExecChildKeepsUnrelatedCallerVariables(t *testing.T) {
	seen := childEnvironment(t)
	requireTest(t, strings.Contains(seen, "BENCH_EXEC_TEST_CARRIED=carried-value"), "child lost the caller's unrelated variable:\n%s", seen)
}
