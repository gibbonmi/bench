package runtime

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func benchPath(t testing.TB) string {
	t.Helper()
	return filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
}

func gitDir(t testing.TB, f contract.Fixture) string {
	t.Helper()
	return strings.TrimSpace(f.Git("rev-parse", "--absolute-git-dir").Stdout)
}

func commitAllowEmpty(t testing.TB, f contract.Fixture, message string) {
	t.Helper()
	f.Run("git", "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", message).RequireExit(0)
}

func cksum(t testing.TB, value string) string {
	t.Helper()
	cmd := exec.Command("cksum")
	cmd.Stdin = strings.NewReader(value + "\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("cksum %q: %v", value, err)
	}
	return strings.Fields(string(out))[0]
}

type runtimePoolWorktrees struct {
	Pool, Warm, Leased, LeaseFile string
}

func addRuntimePoolWorktrees(t testing.TB, f contract.Fixture, benchHome string) runtimePoolWorktrees {
	t.Helper()
	repoRoot := strings.TrimSpace(f.Git("rev-parse", "--show-toplevel").Stdout)
	pool := filepath.Join(benchHome, "worktrees", filepath.Base(repoRoot)+"-"+cksum(t, repoRoot))
	warm := filepath.Join(pool, "warm")
	leased := filepath.Join(pool, "leased")
	contract.Mkdir(t, pool)
	f.Git("worktree", "add", "-q", "--detach", warm, "HEAD")
	f.Git("worktree", "add", "-q", "--detach", leased, "HEAD")
	lease := strings.TrimSpace(contract.RunAt(t, f, leased, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, lease, "")
	return runtimePoolWorktrees{Pool: pool, Warm: warm, Leased: leased, LeaseFile: lease}
}

func repeatLines(n int, line string) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(line)
	}
	return b.String()
}

func copyRuntimeFile(t testing.TB, src, dst string, mode os.FileMode) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, data, mode); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
}

func countStatusRows(out string) int {
	n := 0
	for _, line := range strings.Split(out, "\n") {
		if regexp.MustCompile(`^  [a-z]`).MatchString(line) {
			n++
		}
	}
	return n
}
