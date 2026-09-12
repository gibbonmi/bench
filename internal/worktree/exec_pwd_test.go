package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

const execPWDHelperEnv = "BENCH_EXEC_PWD_HELPER"

// A helper process controls inherited PWD without changing the test process.
// Direct env and pwd children leave a stale PWD visible; a shell can repair it.
func TestExecPWDMatchesChildDirectory(t *testing.T) {
	t.Parallel()
	if mode := os.Getenv(execPWDHelperEnv); mode != "" {
		checkExecPWD(t, mode)
		return
	}
	for _, mode := range []string{"inherited", "absent", "repeated_overrides"} {
		t.Run(mode, func(t *testing.T) {
			cmd := descendant(t, os.Args[0], "-test.run=^TestExecPWDMatchesChildDirectory$", "-test.v")
			cmd.Env = append(capability.WithoutEnvironment(os.Environ(), "PWD"),
				execPWDHelperEnv+"="+mode, "PWD_EXTRA=inherited-near-match")
			if mode != "absent" {
				cmd.Env = append(cmd.Env, "PWD=/caller/pwd")
			}
			out, err := cmd.CombinedOutput()
			requireTest(t, err == nil, "exec PWD helper: %v; %s", err, strings.ReplaceAll(string(out), "\n", " | "))
		})
	}
}

func checkExecPWD(t *testing.T, mode string) {
	t.Helper()
	for _, wrapper := range []bool{false, true} {
		name := "without_wrapper"
		if wrapper {
			name = "with_wrapper"
		}
		t.Run(name, func(t *testing.T) {
			root, creation, home := newOwnedAssignment(t, "pwd-"+mode+"-"+name)
			dir := creation.Assignment.Worktree
			if wrapper {
				path := wrapperPathIn(t, dir)
				requireTest(t, os.WriteFile(path, nil, 0o644) == nil, "write wrapper")
			}
			args := []string{creation.Assignment.Label}
			if mode == "repeated_overrides" {
				args = append(args, "--env", "PWD=/first/override", "--env", "PWD=/last/override")
			}
			args = append(args, "--env", "BENCH_EXEC_PWD_CARRIED=explicit-value")
			var stdout, stderr bytes.Buffer
			code := ExecCommand(root, home, append(args, "--", "pwd", "-P"), nil, &stdout, &stderr)
			requireTest(t, code == 0, "pwd exited %d: %s", code, stderr.String())
			requireTest(t, stdout.String() == dir+"\n", "actual cwd = %q, want %q", stdout.String(), dir)
			stdout.Reset()
			stderr.Reset()
			code = ExecCommand(root, home, append(args, "--", "env"), nil, &stdout, &stderr)
			requireTest(t, code == 0, "env exited %d: %s", code, stderr.String())
			entries := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
			requireExecPWD(t, entries, dir)
			for key, want := range map[string]string{
				"PWD_EXTRA":              "inherited-near-match",
				"BENCH_EXEC_PWD_CARRIED": "explicit-value",
				"BENCH_HOME":             home,
			} {
				got, present := assignment(stdout.String(), key)
				requireTest(t, present && got == want, "%s = %q, present %t; want %q", key, got, present, want)
			}
			marker, present := assignment(stdout.String(), "BENCH_WRAPPER")
			requireTest(t, present == wrapper, "wrapper present = %t, want %t", present, wrapper)
			if wrapper {
				requireTest(t, marker == filepath.Join(dir, "bin", "bench.sh"), "wrapper = %q", marker)
			}
		})
	}
}

// os/exec deduplicates environment keys before the child starts. Inspect the
// composed slice too, so that deduplication cannot hide multiple PWD entries.
func TestExecEnvironmentContainsOneCanonicalPWD(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	entries := execEnv(dir, t.TempDir(), []string{"PWD=/first", "PWD=/last", "PWD_EXTRA=kept"})
	requireExecPWD(t, entries, dir)
	seen := "\n" + strings.Join(entries, "\n") + "\n"
	requireTest(t, strings.Contains(seen, "\nPWD_EXTRA=kept\n"), "composed environment lost PWD_EXTRA")
}

func requireExecPWD(t *testing.T, entries []string, want string) {
	t.Helper()
	var values []string
	for _, entry := range entries {
		if value, ok := strings.CutPrefix(entry, "PWD="); ok {
			values = append(values, value)
		}
	}
	requireTest(t, len(values) == 1 && values[0] == want, "PWD entries = %q, want exactly [%q]", values, want)
}
