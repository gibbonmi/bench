package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// childEnvironment runs a real child through the exec path against worktree and returns
// the environment that child actually received, one assignment per line.
func childEnvironment(t *testing.T, worktree string) string {
	t.Helper()
	t.Setenv("BENCH_KIT", "/caller/kit")
	t.Setenv("BENCH_WRAPPER", "/caller/bin/bench")
	t.Setenv(runbinary.Env, "/caller/run/bench")
	t.Setenv("BENCH_EXEC_TEST_CARRIED", "carried-value")
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"sh", "-c", "env"}, worktree, nil, &stdout, &stderr)
	requireTest(t, code == 0, "child exit = %d, stderr %q", code, stderr.String())
	return stdout.String()
}

// wrapperPathIn is the path exec inspects inside a worktree, with its parent created.
func wrapperPathIn(t *testing.T, worktree string) string {
	t.Helper()
	requireTest(t, os.MkdirAll(filepath.Join(worktree, "bin"), 0o755) == nil, "create bin/")
	return filepath.Join(worktree, "bin", "bench.sh")
}

// worktreeWithWrapper is a worktree carrying a regular-file wrapper holding content.
func worktreeWithWrapper(t *testing.T) (string, string) {
	t.Helper()
	worktree := t.TempDir()
	wrapper := wrapperPathIn(t, worktree)
	requireTest(t, os.WriteFile(wrapper, []byte("#!/bin/sh\nexit 0\n"), 0o755) == nil, "write wrapper")
	return worktree, wrapper
}

// assignment returns the value the child received for name, and whether it received it.
func assignment(seen, name string) (string, bool) {
	for _, line := range strings.Split(seen, "\n") {
		if value, ok := strings.CutPrefix(line, name+"="); ok {
			return value, true
		}
	}
	return "", false
}

// TestExecChildIsRootedAtTheWorktreeWrapper covers WX1 and WX2.
func TestExecChildIsRootedAtTheWorktreeWrapper(t *testing.T) {
	worktree, wrapper := worktreeWithWrapper(t)
	seen := childEnvironment(t, worktree)
	value, ok := assignment(seen, "BENCH_WRAPPER")
	requireTest(t, ok, "child received no BENCH_WRAPPER:\n%s", seen)
	requireTest(t, value == wrapper, "child saw BENCH_WRAPPER=%q, want %q", value, wrapper)
	requireTest(t, filepath.IsAbs(value), "child saw a relative BENCH_WRAPPER=%q", value)
}

// TestExecChildOwnsRatherThanInheritsItsRunBinary covers the environment half of WX3:
// the owner lookup returns an owning selection exactly when the wrapper marker is
// non-empty and the run-binary variable is absent. The lookup itself is graded by the
// composed run recorded as WX20 evidence.
func TestExecChildOwnsRatherThanInheritsItsRunBinary(t *testing.T) {
	worktree, _ := worktreeWithWrapper(t)
	seen := childEnvironment(t, worktree)
	marker, ok := assignment(seen, "BENCH_WRAPPER")
	requireTest(t, ok && marker != "", "child received no non-empty wrapper marker:\n%s", seen)
	_, inherited := assignment(seen, runbinary.Env)
	requireTest(t, !inherited, "child received %s, so its gate would inherit a binary:\n%s", runbinary.Env, seen)
}

// TestExecChildDropsWrapperRouting covers WX4 and WX5.
func TestExecChildDropsWrapperRouting(t *testing.T) {
	worktree, _ := worktreeWithWrapper(t)
	seen := childEnvironment(t, worktree)
	for _, name := range []string{"BENCH_KIT", runbinary.Env} {
		_, ok := assignment(seen, name)
		requireTest(t, !ok, "child saw %s in its environment:\n%s", name, seen)
	}
	value, _ := assignment(seen, "BENCH_WRAPPER")
	requireTest(t, value != "/caller/bin/bench", "child inherited the caller's wrapper value:\n%s", seen)
}

// TestExecChildKeepsUnrelatedCallerVariables covers WX6.
func TestExecChildKeepsUnrelatedCallerVariables(t *testing.T) {
	worktree, _ := worktreeWithWrapper(t)
	seen := childEnvironment(t, worktree)
	value, ok := assignment(seen, "BENCH_EXEC_TEST_CARRIED")
	requireTest(t, ok && value == "carried-value", "child lost the caller's unrelated variable:\n%s", seen)
}

// TestExecChildTakesNoMarkerFromANonWrapper covers WX7 through WX10: the marker follows
// one predicate, so every path that is not a regular file leaves the child as it is
// today.
func TestExecChildTakesNoMarkerFromANonWrapper(t *testing.T) {
	for _, tc := range []struct {
		name  string
		build func(t *testing.T, worktree string)
	}{
		{name: "absent", build: func(*testing.T, string) {}},
		{name: "directory", build: func(t *testing.T, worktree string) {
			requireTest(t, os.MkdirAll(wrapperPathIn(t, worktree), 0o755) == nil, "create wrapper directory")
		}},
		{name: "fifo", build: func(t *testing.T, worktree string) {
			requireTest(t, syscall.Mkfifo(wrapperPathIn(t, worktree), 0o644) == nil, "create wrapper fifo")
		}},
		{name: "dangling symlink", build: func(t *testing.T, worktree string) {
			wrapper := wrapperPathIn(t, worktree)
			requireTest(t, os.Symlink(filepath.Join(worktree, "gone"), wrapper) == nil, "create dangling wrapper link")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			worktree := t.TempDir()
			tc.build(t, worktree)
			seen := childEnvironment(t, worktree)
			value, ok := assignment(seen, "BENCH_WRAPPER")
			requireTest(t, !ok, "a %s wrapper path became the marker BENCH_WRAPPER=%q:\n%s", tc.name, value, seen)
		})
	}
}

// TestExecChildTakesTheMarkerFromAnEmptyWrapper covers WX11: the predicate is the file's
// kind, never its content.
func TestExecChildTakesTheMarkerFromAnEmptyWrapper(t *testing.T) {
	worktree := t.TempDir()
	wrapper := wrapperPathIn(t, worktree)
	requireTest(t, os.WriteFile(wrapper, nil, 0o755) == nil, "write empty wrapper")
	seen := childEnvironment(t, worktree)
	value, ok := assignment(seen, "BENCH_WRAPPER")
	requireTest(t, ok && value == wrapper, "an empty wrapper gave BENCH_WRAPPER=%q (present=%v), want %q", value, ok, wrapper)
}

// TestExecChildDiffersByTheWrapperAssignmentAlone covers WX13, and with it the unit half
// of WX12: an unmarked child is exactly today's child, so the whole change to any verb's
// environment is the one assignment. Both children run in one worktree so nothing
// path-derived varies between them.
func TestExecChildDiffersByTheWrapperAssignmentAlone(t *testing.T) {
	worktree := t.TempDir()
	wrapper := wrapperPathIn(t, worktree)
	today := strings.Split(childEnvironment(t, worktree), "\n")
	requireTest(t, os.WriteFile(wrapper, []byte("#!/bin/sh\n"), 0o755) == nil, "write wrapper")
	marked := strings.Split(childEnvironment(t, worktree), "\n")

	added := map[string]int{}
	for _, line := range marked {
		added[line]++
	}
	for _, line := range today {
		added[line]--
	}
	var difference []string
	for line, count := range added {
		for ; count > 0; count-- {
			difference = append(difference, line)
		}
		requireTest(t, count >= 0, "the marked child lost the assignment %q", line)
	}
	requireTest(t, len(difference) == 1 && difference[0] == "BENCH_WRAPPER="+wrapper,
		"the marked child differs from today's by %q, want the wrapper assignment alone", difference)
}
