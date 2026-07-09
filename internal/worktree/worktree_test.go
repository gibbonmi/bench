package worktree

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// cksumGolden pins the Go cksum against values produced by the coreutils `cksum`
// tool. Each was derived once with, e.g.:
//
//	printf '%s\n' "/home/mgibs/workspace/bench" | cksum   -> 2826441890 28
//	printf '%s\n' "/tmp/a b/c"                  | cksum   -> 889650394  11
//
// The `\n` is intentional: Pool checksums `root + "\n"` because the shell used
// `echo "$root" | cksum`. The second vector carries a space to exercise a path the
// shell would otherwise word-split.
var cksumGolden = []struct {
	root string
	sum  uint32
}{
	{"/home/mgibs/workspace/bench", 2826441890},
	{"/tmp/a b/c", 889650394},
}

func TestCksumMatchesGolden(t *testing.T) {
	for _, g := range cksumGolden {
		got := cksum([]byte(g.root + "\n"))
		if got != g.sum {
			t.Errorf("cksum(%q+NL) = %d, want %d", g.root, got, g.sum)
		}
	}
}

// TestCksumMatchesSystemTool cross-checks against the live `cksum` when it is on
// PATH, so the pinned goldens can never silently drift from the real tool. It is
// skipped where `cksum` is unavailable, keeping the suite hermetic there.
func TestCksumMatchesSystemTool(t *testing.T) {
	if _, err := exec.LookPath("cksum"); err != nil {
		t.Skip("cksum not available")
	}
	for _, g := range cksumGolden {
		// printf '%s\n' "<root>" | cksum
		printf := exec.Command("printf", "%s\n", g.root)
		ck := exec.Command("cksum")
		pipe, err := printf.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		ck.Stdin = pipe
		out, err := ck.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := ck.Start(); err != nil {
			t.Fatal(err)
		}
		if err := printf.Run(); err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 64)
		n, _ := out.Read(buf)
		if err := ck.Wait(); err != nil {
			t.Fatal(err)
		}
		field := strings.Fields(string(buf[:n]))
		if len(field) == 0 {
			t.Fatalf("empty cksum output for %q", g.root)
		}
		want, err := strconv.ParseUint(field[0], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		if got := cksum([]byte(g.root + "\n")); uint64(got) != want {
			t.Errorf("cksum(%q) = %d, system tool = %d", g.root, got, want)
		}
	}
}

func TestPool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := "/home/mgibs/workspace/bench"
	want := filepath.Join(home, "worktrees", "bench-2826441890")
	if got := Pool(root); got != want {
		t.Errorf("Pool(%q) = %q, want %q", root, got, want)
	}
}

func TestPoolDefaultBenchHome(t *testing.T) {
	// With BENCH_HOME unset, Pool falls back to <home>/.bench.
	t.Setenv("BENCH_HOME", "")
	root := "/tmp/a b/c"
	got := Pool(root)
	suffix := filepath.Join(".bench", "worktrees", "c-889650394")
	if !strings.HasSuffix(got, suffix) {
		t.Errorf("Pool(%q) = %q, want suffix %q", root, got, suffix)
	}
}

func TestClassifyRegisteredWorktrees(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)

	pool := Pool(root)
	warm := filepath.Join(pool, "warm")
	leased := filepath.Join(pool, "leased")
	outOfPool := filepath.Join(filepath.Dir(root), "outside pool")
	if err := os.MkdirAll(pool, 0o755); err != nil {
		t.Fatalf("mkdir pool: %v", err)
	}
	gitRun(t, root, "worktree", "add", "-q", "--detach", warm, "HEAD")
	gitRun(t, root, "worktree", "add", "-q", "--detach", leased, "HEAD")
	gitRun(t, root, "worktree", "add", "-q", "--detach", outOfPool, "HEAD")
	lease, err := LeaseFile(leased)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	if err := os.WriteFile(lease, []byte("123 2026-07-06T00:00:00Z\n"), 0o644); err != nil {
		t.Fatalf("write lease: %v", err)
	}

	entries, err := ClassifyRegisteredWorktrees(root)
	if err != nil {
		t.Fatalf("ClassifyRegisteredWorktrees: %v", err)
	}
	got := map[string]Class{}
	for _, entry := range entries {
		got[entry.Path] = entry.Class
	}
	want := map[string]Class{
		root:      ClassRoot,
		warm:      ClassPoolWarm,
		leased:    ClassPoolLease,
		outOfPool: ClassOutOfPool,
	}
	for path, class := range want {
		if got[path] != class {
			t.Errorf("class %q = %q, want %q", path, got[path], class)
		}
	}
	linkedEntries, err := ClassifyRegisteredWorktrees(leased)
	if err != nil {
		t.Fatalf("ClassifyRegisteredWorktrees from linked worktree: %v", err)
	}
	got = map[string]Class{}
	for _, entry := range linkedEntries {
		got[entry.Path] = entry.Class
	}
	for path, class := range want {
		if got[path] != class {
			t.Errorf("class from linked cwd %q = %q, want %q", path, got[path], class)
		}
	}
}

func TestCleanCommandRemovesOutOfPool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	candidate := filepath.Join(filepath.Dir(root), "outside wt [one]")
	gitRun(t, root, "worktree", "add", "-q", "--detach", candidate, "HEAD")
	nested := filepath.Join(root, "sub", "dir")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested cwd: %v", err)
	}
	chdir(t, nested)

	var stdout, stderr bytes.Buffer
	code := cleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("cleanCommand exit = %d\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), candidate) || !strings.Contains(stdout.String(), "removed") {
		t.Fatalf("cleanup did not report removed candidate:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String()+stderr.String(), "--force") {
		t.Fatalf("cleanup mentioned forced removal:\nstdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
	if out := gitOutput(t, root, "worktree", "list", "--porcelain"); strings.Contains(out, candidate) {
		t.Fatalf("cleanup left worktree registered:\n%s", out)
	}

	stdout.Reset()
	stderr.Reset()
	code = cleanCommand(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("second cleanCommand exit = %d\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "nothing to clean") {
		t.Fatalf("second cleanup did not report idempotent empty state:\n%s", stdout.String())
	}
}

// TestCleanOutOfPoolWorktreesGitFailureNeverReadsAsNothingToClean is the FT29 false-empty
// regression guard for the classifier's last swallowing caller: a `git worktree list`
// failure (deterministically induced by making .git unreadable, the gitOpError-style
// injection FT29 used in structure_test.go) must exit non-zero with the git error on
// stderr, never fall through to the "nothing to clean" exit-0 report that reads
// identically to a genuinely empty, healthy repo.
func TestCleanOutOfPoolWorktreesGitFailureNeverReadsAsNothingToClean(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)

	gitDir := filepath.Join(root, ".git")
	if err := os.Chmod(gitDir, 0o000); err != nil {
		t.Fatalf("chmod .git unreadable: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(gitDir, 0o755) })

	var stdout, stderr bytes.Buffer
	code := cleanOutOfPoolWorktrees(root, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("cleanOutOfPoolWorktrees exit 0 on a git worktree-list failure\nstdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
	if strings.Contains(stdout.String(), "nothing to clean") {
		t.Fatalf("git failure read as the empty-repo all-clear:\nstdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
	if stderr.String() == "" {
		t.Fatal("cleanOutOfPoolWorktrees printed no error to stderr on a git failure")
	}
}

func TestLeaseFile(t *testing.T) {
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	lease, err := LeaseFile(dir)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	if !strings.HasSuffix(lease, "bench-lease") {
		t.Errorf("LeaseFile = %q, want suffix bench-lease", lease)
	}
	if !filepath.IsAbs(lease) {
		t.Errorf("LeaseFile = %q, want absolute — a relative path resolves against the caller's CWD, not the worktree", lease)
	}
}

func TestLeaseFileCommandMissingArg(t *testing.T) {
	out, code := LeaseFileCommand(nil)
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if !strings.HasPrefix(out, "usage:") {
		t.Errorf("out = %q, want usage line", out)
	}
}

func TestPoolCommandExplicitRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	out, code := PoolCommand([]string{"/home/mgibs/workspace/bench"})
	if code != 0 {
		t.Errorf("exit = %d, want 0", code)
	}
	want := filepath.Join(home, "worktrees", "bench-2826441890") + "\n"
	if out != want {
		t.Errorf("out = %q, want %q", out, want)
	}
}

func newWorktreeRepo(t testing.TB) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init")
	gitRun(t, root, "config", "user.email", "bench@local")
	gitRun(t, root, "config", "user.name", "bench")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write tracked.txt: %v", err)
	}
	gitRun(t, root, "add", "tracked.txt")
	gitRun(t, root, "commit", "-q", "-m", "base")
	return root
}

func chdir(t testing.TB, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
}

func gitOutput(t testing.TB, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

func gitRun(t testing.TB, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}
