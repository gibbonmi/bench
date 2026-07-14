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
