package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// childOutput runs a real child through the exec path against worktree, with the caller's
// routing variables set, and returns what that child wrote.
func childOutput(t *testing.T, worktree, script string) string {
	t.Helper()
	return childOutputAtHome(t, worktree, script, Home())
}

// childOutputAtHome is childOutput with the home the verb resolved passed explicitly, so
// a row can separate the caller's home from the one the process carries.
func childOutputAtHome(t *testing.T, worktree, script, home string) string {
	t.Helper()
	bindEnv(t, "BENCH_KIT", "/caller/kit")
	bindEnv(t, "BENCH_WRAPPER", "/caller/bin/bench")
	bindEnv(t, runbinary.Env, "/caller/run/bench")
	bindEnv(t, "BENCH_EXEC_TEST_CARRIED", "carried-value")
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"sh", "-c", script}, worktree, home, nil, nil, &stdout, &stderr)
	requireTest(t, code == 0, "child exit = %d, stderr %q", code, stderr.String())
	return stdout.String()
}

// childEnvironment returns the environment the child actually received, one assignment
// per line. That format cannot frame a value holding a newline, so it grades presence and
// absence over known-safe values only; childWrapper grades the marker's exact value.
func childEnvironment(t *testing.T, worktree string) string {
	t.Helper()
	return childOutput(t, worktree, "env")
}

// childWrapper returns the exact wrapper marker the child received, and whether it
// received one. The child emits the presence flag and the value as NUL-terminated fields.
// No environment value can forge that frame. A path holding a newline, a space, or a
// glob character arrives byte-for-byte instead of being re-split by the reader.
func childWrapper(t *testing.T, worktree string) (string, bool) {
	t.Helper()
	emitted := childOutput(t, worktree, `printf '%s\0%s\0' "${BENCH_WRAPPER+set}" "${BENCH_WRAPPER-}"`)
	fields := strings.Split(emitted, "\x00")
	requireTest(t, len(fields) == 3 && fields[2] == "", "child emitted %q, want two NUL-terminated fields", emitted)
	return fields[1], fields[0] == "set"
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

// TestExecChildTakesTheExplicitHome covers WF16. The child reads the home the verb was
// given, not the one the caller's process carries, so it resolves the same pool the verb
// resolved. The two homes differ here, so an inherited value fails the comparison. The
// child emits the value as a NUL-terminated field, which no environment value can forge.
// The bind keeps this test serial.
func TestExecChildTakesTheExplicitHome(t *testing.T) {
	worktree, _ := worktreeWithWrapper(t)
	explicit := t.TempDir()
	bindEnv(t, homeEnv, t.TempDir())
	emitted := childOutputAtHome(t, worktree, `printf '%s\0' "${BENCH_HOME-}"`, explicit)
	fields := strings.Split(emitted, "\x00")
	requireTest(t, len(fields) == 2 && fields[1] == "", "child emitted %q, want one NUL-terminated field", emitted)
	requireTest(t, fields[0] == explicit, "child saw BENCH_HOME=%q, want the explicit home %q", fields[0], explicit)
}

// TestExecChildIsRootedAtTheWorktreeWrapper covers WX1 and WX2.
func TestExecChildIsRootedAtTheWorktreeWrapper(t *testing.T) {
	worktree, wrapper := worktreeWithWrapper(t)
	value, ok := childWrapper(t, worktree)
	requireTest(t, ok, "child received no BENCH_WRAPPER")
	requireTest(t, value == wrapper, "child saw BENCH_WRAPPER=%q, want %q", value, wrapper)
	requireTest(t, filepath.IsAbs(value), "child saw a relative BENCH_WRAPPER=%q", value)
}

// TestExecChildOwnsRatherThanInheritsItsRunBinary covers the environment half of WX3.
// The owner lookup returns an owning selection exactly when the wrapper marker is
// non-empty and the run-binary variable is absent. The lookup itself is graded by the
// composed run recorded as WX20 evidence.
func TestExecChildOwnsRatherThanInheritsItsRunBinary(t *testing.T) {
	worktree, _ := worktreeWithWrapper(t)
	seen := childEnvironment(t, worktree)
	marker, ok := childWrapper(t, worktree)
	requireTest(t, ok && marker != "", "child received no non-empty wrapper marker")
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
	value, _ := childWrapper(t, worktree)
	requireTest(t, value != "/caller/bin/bench", "child inherited the caller's wrapper value: %q", value)
}

// TestExecChildKeepsUnrelatedCallerVariables covers WX6.
func TestExecChildKeepsUnrelatedCallerVariables(t *testing.T) {
	worktree, _ := worktreeWithWrapper(t)
	seen := childEnvironment(t, worktree)
	value, ok := assignment(seen, "BENCH_EXEC_TEST_CARRIED")
	requireTest(t, ok && value == "carried-value", "child lost the caller's unrelated variable:\n%s", seen)
}

// TestExecChildTakesNoMarkerFromANonWrapper covers WX7 through WX10. The marker follows
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
			value, ok := childWrapper(t, worktree)
			requireTest(t, !ok, "a %s wrapper path became the marker BENCH_WRAPPER=%q", tc.name, value)
		})
	}
}

// TestExecChildTakesTheMarkerFromAnEmptyWrapper covers WX11: the predicate is the file's
// kind, never its content.
func TestExecChildTakesTheMarkerFromAnEmptyWrapper(t *testing.T) {
	worktree := t.TempDir()
	wrapper := wrapperPathIn(t, worktree)
	requireTest(t, os.WriteFile(wrapper, nil, 0o755) == nil, "write empty wrapper")
	value, ok := childWrapper(t, worktree)
	requireTest(t, ok && value == wrapper, "an empty wrapper gave BENCH_WRAPPER=%q (present=%v), want %q", value, ok, wrapper)
}

// TestExecChildDiffersByTheWrapperAssignmentAlone covers WX13, and with it the unit half
// of WX12. An unmarked child is exactly today's child, so the whole change to any verb's
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

// TestExecChildTakesTheMarkerThroughALiveWrapperLink covers WX21. The predicate follows
// links, so a link to a regular file is a wrapper. The marker names the link path the
// child's own wrapper resolution will walk rather than the target behind it.
func TestExecChildTakesTheMarkerThroughALiveWrapperLink(t *testing.T) {
	worktree := t.TempDir()
	target := filepath.Join(worktree, "linked-wrapper.sh")
	requireTest(t, os.WriteFile(target, []byte("#!/bin/sh\n"), 0o755) == nil, "write link target")
	wrapper := wrapperPathIn(t, worktree)
	requireTest(t, os.Symlink(target, wrapper) == nil, "create live wrapper link")
	value, ok := childWrapper(t, worktree)
	requireTest(t, ok, "a live wrapper link took no marker, want BENCH_WRAPPER=%q", wrapper)
	requireTest(t, value == wrapper, "a live wrapper link gave BENCH_WRAPPER=%q, want the link path %q", value, wrapper)
}

// TestExecChildTakesTheExactPathFromAnAwkwardWorktreeName covers WX22. The marker is a
// joined path handed to the child as one value, so a space cannot split it and a glob
// character cannot expand.
func TestExecChildTakesTheExactPathFromAnAwkwardWorktreeName(t *testing.T) {
	worktree := filepath.Join(t.TempDir(), "nested [*] path")
	requireTest(t, os.MkdirAll(worktree, 0o755) == nil, "create awkward worktree")
	wrapper := wrapperPathIn(t, worktree)
	requireTest(t, os.WriteFile(wrapper, []byte("#!/bin/sh\n"), 0o755) == nil, "write wrapper")
	value, ok := childWrapper(t, worktree)
	requireTest(t, ok, "an awkwardly named worktree took no marker, want BENCH_WRAPPER=%q", wrapper)
	requireTest(t, value == wrapper, "child saw BENCH_WRAPPER=%q, want the exact path %q", value, wrapper)
}

// execAtOwnedTarget runs ExecCommand against a real owned assignment, so a row grades the
// public seam an agent types. It returns the home the verb resolved, the child's stdout,
// the refusal stream, and the exit code.
func execAtOwnedTarget(t *testing.T, request string, args ...string) (string, string, string, int) {
	t.Helper()
	root, creation, home := newOwnedAssignment(t, request)
	var stdout, stderr bytes.Buffer
	code := ExecCommand(root, home, append([]string{creation.Assignment.Label}, args...), nil, &stdout, &stderr)
	return home, stdout.String(), stderr.String(), code
}

// TestExecHelpCarriesStdinAndTheExitRule covers X3 and X6. The help is three lines: the
// grammar, a stdin heredoc example, and the rule that tells a grammar refusal from a
// child's own exit 2.
func TestExecHelpCarriesStdinAndTheExitRule(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	code := ExecCommand(t.TempDir(), t.TempDir(), []string{"--help"}, nil, &stdout, &stderr)
	requireTest(t, code == 0, "help exited %d", code)
	lines := strings.Split(strings.TrimRight(stderr.String(), "\n"), "\n")
	requireTest(t, len(lines) == 3, "help printed %d lines, want 3:\n%s", len(lines), stderr.String())
	requireTest(t, lines[0] == "usage: "+usage.WorktreeExec, "help line 1 = %q, want the grammar line", lines[0])
	requireTest(t, strings.Contains(lines[1], "<<'EOF'") && strings.Contains(lines[1], "-- python3 -"),
		"help line 2 = %q, want a heredoc feeding -- python3 -", lines[1])
	requireTest(t, strings.Contains(lines[2], "usage: bench worktree exec") && strings.Contains(lines[2], "exit 2"),
		"help line 3 = %q, want the exit-2 rule naming the refusal prefix", lines[2])
}

// TestExecChildReadsTheCallersStdin covers X4. The caller's reader reaches the child byte
// for byte, so a heredoc arrives whole. The NUL byte proves no line-oriented rewrite.
func TestExecChildReadsTheCallersStdin(t *testing.T) {
	t.Parallel()
	fed := "first\nsecond\nthird\x00"
	var stdout, stderr bytes.Buffer
	code := runWorktreeChild([]string{"cat"}, t.TempDir(), t.TempDir(), nil, strings.NewReader(fed), &stdout, &stderr)
	requireTest(t, code == 0, "cat exited %d, stderr %q", code, stderr.String())
	requireTest(t, stdout.String() == fed, "child emitted %q, want the fed bytes %q", stdout.String(), fed)
}

// TestExecEnvReachesTheChildOnly covers X7. The child reads the assignment the caller
// never set, so the flag the parser accepts is the flag the child receives.
func TestExecEnvReachesTheChildOnly(t *testing.T) {
	t.Parallel()
	requireTest(t, os.Getenv("FOO") == "", "the test process already carries FOO")
	_, stdout, stderr, code := execAtOwnedTarget(t, "env-child-only", "--env", "FOO=bar", "--", "sh", "-c", "echo $FOO")
	requireTest(t, code == 0, "exec exited %d: %s", code, stderr)
	requireTest(t, strings.TrimSpace(stdout) == "bar", "child printed %q, want bar", stdout)
	requireTest(t, os.Getenv("FOO") == "", "the caller's process gained FOO")
}

// TestExecRefusesAMalformedEnvValue covers X10. A value with no "=" and a value with a
// bad KEY each refuse at exit 2 naming the value, and neither starts a child. The marker
// file the child would write stays absent.
func TestExecRefusesAMalformedEnvValue(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"FOO", "1X=y"} {
		t.Run(value, func(t *testing.T) {
			marker := filepath.Join(t.TempDir(), "started")
			_, stdout, stderr, code := execAtOwnedTarget(t, "env-malformed",
				"--env", value, "--", "sh", "-c", "touch '"+marker+"'")
			requireTest(t, code == 2, "--env %q exited %d, want 2", value, code)
			requireTest(t, strings.TrimSpace(stderr) == toon.Usage("bench worktree exec", value),
				"--env %q printed %q, want the usage line naming the value", value, stderr)
			_, err := os.Stat(marker)
			requireTest(t, os.IsNotExist(err), "--env %q started a child: %s exists (stdout %q)", value, marker, stdout)
		})
	}
}

// TestExecEnvCannotRepointTheChildsPool covers X11. The verb's own routing values apply
// after the --env values, so the child reads the home the verb resolved.
func TestExecEnvCannotRepointTheChildsPool(t *testing.T) {
	t.Parallel()
	home, stdout, stderr, code := execAtOwnedTarget(t, "env-home-loses",
		"--env", "BENCH_HOME=/x", "--", "sh", "-c", "echo $BENCH_HOME")
	requireTest(t, code == 0, "exec exited %d: %s", code, stderr)
	requireTest(t, strings.TrimSpace(stdout) == home, "child saw BENCH_HOME=%q, want the resolved home %q", strings.TrimSpace(stdout), home)
}

// TestExecEnvCannotRestoreTheStrippedKit covers X13. The routing strip runs after the
// --env values, so a named kit leaves the child with no assignment at all.
func TestExecEnvCannotRestoreTheStrippedKit(t *testing.T) {
	t.Parallel()
	_, stdout, stderr, code := execAtOwnedTarget(t, "env-kit-stripped",
		"--env", "BENCH_KIT=/x", "--", "sh", "-c", "echo ${BENCH_KIT-unset}")
	requireTest(t, code == 0, "exec exited %d: %s", code, stderr)
	requireTest(t, strings.TrimSpace(stdout) == "unset", "child saw BENCH_KIT=%q, want it stripped", strings.TrimSpace(stdout))
}
