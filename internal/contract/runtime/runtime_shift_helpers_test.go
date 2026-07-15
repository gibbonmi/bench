package runtime

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func shiftFixture(t *testing.T, gate string) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Env["GIT_AUTHOR_NAME"] = "Bench"
	f.Env["GIT_AUTHOR_EMAIL"] = "bench@local"
	f.Env["GIT_COMMITTER_NAME"] = "Bench"
	f.Env["GIT_COMMITTER_EMAIL"] = "bench@local"
	f.WriteExecutable(".bench/gate.sh", gate)
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":["BENCH_TEST_PROMPTS","BENCH_TEST_STATE"],"paths":[],"tools":[]}`+"\n")
	f.CommitAll("init")
	return f
}

func shiftBranch(t *testing.T, output string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^■ shift done: ([^,]+),`)
	m := re.FindStringSubmatch(output)
	if len(m) != 2 {
		t.Fatalf("shift summary did not name the branch:\n%s", output)
	}
	return m[1]
}

func shiftWorktree(t *testing.T, output string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^  worktree: (.+)$`)
	m := re.FindStringSubmatch(output)
	if len(m) != 2 {
		t.Fatalf("shift did not report its worktree:\n%s", output)
	}
	return m[1]
}

func requireNoWorktreeBranch(t *testing.T, f contract.Fixture, branch string) {
	t.Helper()
	out := f.Git("worktree", "list", "--porcelain").Stdout
	if strings.Contains(out, "branch refs/heads/"+branch) {
		t.Fatalf("released worktree still holds the shift branch:\n%s", out)
	}
}

func requireNoLease(t *testing.T, home string) {
	t.Helper()
	var found []string
	err := filepath.WalkDir(home, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "bench-lease" {
			found = append(found, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk BENCH_HOME for leases: %v", err)
	}
	if len(found) > 0 {
		t.Fatalf("shift worktree lease was not released: %s", strings.Join(found, ", "))
	}
}

func requireEqual(t *testing.T, got, want, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %q, want %q", msg, got, want)
	}
}

func runGitAt(t *testing.T, root string, args ...string) string {
	t.Helper()
	f := contract.NewFixtureAt(t, root, contract.IsolatedEnv(t, root))
	return f.Git(args...).Stdout
}

func waitSeconds(t *testing.T, seconds int) {
	t.Helper()
	time.Sleep(time.Duration(seconds) * time.Second)
}

func strPtr(s string) *string {
	return &s
}

func requireFileContains(t *testing.T, path, needle, msg string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(data), needle) {
		t.Fatal(msg)
	}
}

func requireWorktreeLifecycleSharedResolver(t *testing.T, f contract.Fixture) {
	t.Helper()
	hook := f.ReadFile(".bench/hooks/worktree-lifecycle.sh")
	for _, needle := range []string{"../lib/resolve-bench.sh", `. "$lib"`, "bench_resolve_wrapper"} {
		if !strings.Contains(hook, needle) {
			t.Fatalf("linked worktree lifecycle hook does not use the shared wrapper resolver: missing %q", needle)
		}
	}
	for _, duplicate := range []string{"for candidate in", `"$root/.bench/bin/bench.sh"`, `"$root/bin/bench.sh"`} {
		if strings.Contains(hook, duplicate) {
			t.Fatalf("linked worktree lifecycle hook duplicates wrapper resolution instead of using the shared resolver: found %q", duplicate)
		}
	}
}

// requireNoShiftBranch asserts no branch matching bench/shift-* exists — the signal
// that a usage failure exited before the loop created anything.
func requireNoShiftBranch(t *testing.T, f contract.Fixture) {
	t.Helper()
	if out := f.Git("for-each-ref", "refs/heads/bench/shift-*").Stdout; strings.TrimSpace(out) != "" {
		t.Fatalf("usage failure created a shift branch:\n%s", out)
	}
}

// requireShiftResult asserts the shift_result TOON block's row equals wantRow exactly —
// header plus the one data row, in the pinned field order.
func requireShiftResult(t *testing.T, stdout, wantRow string) {
	t.Helper()
	header := "shift_result[1]{outcome,exit,branch,committed,iterations_used,recovery,detail}:"
	if !strings.Contains(stdout, header) {
		t.Fatalf("shift_result block missing or malformed header:\n%s", stdout)
	}
	if !strings.Contains(stdout, wantRow) {
		t.Fatalf("shift_result row = %q not found in:\n%s", wantRow, stdout)
	}
}
