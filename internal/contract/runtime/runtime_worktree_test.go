package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeWorktreeContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench worktree clean confirmed cleanup contract", testRuntimeWorktreeCleanRemovesConfirmedOutOfPool)
	contract.RunParallel(t, "bench worktree clean non-tty refusal contract", testRuntimeWorktreeCleanRejectsNonTTY)
	contract.RunParallel(t, "bench worktree clean dirty refusal contract", testRuntimeWorktreeCleanLeavesDirtyOutOfPool)
	contract.RunParallel(t, "bench worktree clean stale registration prune contract", testRuntimeWorktreeCleanPrunesMissingRegistration)
	contract.RunParallel(t, "bench worktree clean pool cwd classification contract", testRuntimeWorktreeCleanFromPoolCwd)
	contract.RunParallel(t, "bench worktree usage contract", testRuntimeWorktreeRejectsUnknownArgs)
	contract.RunParallel(t, "bench worktree lease/reuse contract", testRuntimeWorktreeLeaseReuse)
	contract.RunParallel(t, "bench worktree lease hardening contract", testRuntimeWorktreeLeaseHardening)
	contract.RunParallel(t, "bench worktree concurrent-acquire contract", testRuntimeWorktreeConcurrentAcquire)
}

func testRuntimeWorktreeCleanRemovesConfirmedOutOfPool(t *testing.T) {
	if _, err := exec.LookPath("script"); err != nil {
		t.Skip("script command unavailable for PTY-backed confirmation")
	}
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	candidate := filepath.Join(f.Root, "..", "outside wt [one]")
	f.Git("worktree", "add", "-q", "--detach", candidate, "HEAD")
	nested := filepath.Join(f.Root, "sub", "dir")
	contract.Mkdir(t, nested)

	command := "bash " + shellQuote(benchPath(t)) + " worktree clean"
	out := contract.RunAtWithInput(t, f, nested, nil, "clean worktrees\n", "script", "-q", "-e", "-c", command, "/dev/null")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, candidate)
	contract.RequireContains(t, out.Stdout, "removed")
	if strings.Contains(out.Stdout+out.Stderr, "--force") {
		t.Fatalf("cleanup mentioned forced removal:\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
	if strings.Contains(f.Git("worktree", "list", "--porcelain").Stdout, candidate) {
		t.Fatalf("confirmed cleanup left worktree registered:\n%s", f.Git("worktree", "list", "--porcelain").Stdout)
	}

	rerun := contract.RunAtWithInput(t, f, nested, nil, "clean worktrees\n", "script", "-q", "-e", "-c", command, "/dev/null")
	rerun.RequireExit(0)
	contract.RequireContains(t, rerun.Stdout, "nothing to clean")
}

func testRuntimeWorktreeCleanRejectsNonTTY(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	candidate := filepath.Join(f.Root, "..", "outside")
	f.Git("worktree", "add", "-q", "--detach", candidate, "HEAD")

	out := contract.RunAtWithInput(t, f, f.Root, nil, "clean worktrees\n", "bash", benchPath(t), "worktree", "clean")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stderr, "requires an interactive TTY")
	contract.RequireContains(t, f.Git("worktree", "list", "--porcelain").Stdout, candidate)
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func testRuntimeWorktreeCleanLeavesDirtyOutOfPool(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	dirty := filepath.Join(f.Root, "..", "dirty outside")
	f.Git("worktree", "add", "-q", "--detach", dirty, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(dirty, "dirty.txt"), "dirty\n")

	out := contract.RunAtWithInput(t, f, f.Root, nil, "clean worktrees\n", "bash", benchPath(t), "worktree", "clean")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stdout, "refused")
	contract.RequireContains(t, out.Stdout, dirty)
	contract.RequireContains(t, f.Git("worktree", "list", "--porcelain").Stdout, dirty)
}

func testRuntimeWorktreeCleanPrunesMissingRegistration(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	missing := filepath.Join(f.Root, "..", "missing outside")
	f.Git("worktree", "add", "-q", "--detach", missing, "HEAD")
	contract.Remove(t, missing)

	out := f.Bench("worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "nothing to clean")
	contract.RequireNotContains(t, f.Git("worktree", "list", "--porcelain").Stdout, missing)
}

func testRuntimeWorktreeCleanFromPoolCwd(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	benchHome := filepath.Join(f.Root, ".bh")
	pool := addRuntimePoolWorktrees(t, f, benchHome)

	out := contract.RunAt(t, f, pool.Leased, map[string]string{"BENCH_HOME": benchHome}, "bash", benchPath(t), "worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "nothing to clean")
	contract.RequireNotContains(t, out.Stdout, f.Root)
	contract.RequireNotContains(t, out.Stdout, pool.Warm)
	contract.RequireNotContains(t, out.Stdout, pool.Leased)
}

func testRuntimeWorktreeRejectsUnknownArgs(t *testing.T) {
	f := contract.NewFixture(t)
	for _, args := range [][]string{
		{"worktree", "badverb"},
		{"worktree", "clean", "extra"},
		{"worktree", "bad", "verb"},
	} {
		out := f.Bench(args...)
		out.RequireExit(2)
		contract.RequireContains(t, out.Stderr, "usage: bench worktree")
		contract.RequireContains(t, out.Stderr, "bench worktree clean")
	}
}

func testRuntimeWorktreeLeaseReuse(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	f.WriteFile(".gitignore", "ignored/\n")
	f.CommitAll("ignore")
	f.WriteExecutable("wt-shell", `#!/usr/bin/env bash
: "${BENCH_WT_RECORD:?}"
pwd >> "$BENCH_WT_RECORD"
lease="$(git rev-parse --git-path bench-lease)"
[ -f "$lease" ] || { echo "lease missing"; exit 7; }
[ ! -e dirty.txt ] || { echo "dirty file carried into reused worktree"; exit 8; }
[ ! -e ignored/leak.txt ] || { echo "ignored artifact carried into reused worktree"; exit 9; }
echo dirty > dirty.txt
mkdir -p ignored; echo leak > ignored/leak.txt
`)
	record := filepath.Join(f.Root, "paths")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "SHELL": filepath.Join(f.Root, "wt-shell")}
	f.BenchEnv(env, "worktree").RequireExit(0)
	f.BenchEnv(env, "worktree").RequireExit(0)
	paths := contract.NonEmptyLines(contract.ReadFileAbs(t, record))
	contract.RequireIntEqual(t, len(paths), 2, "worktree shell did not run twice")
	requireEqual(t, paths[0], paths[1], "worktree pool did not reuse a clean released path")
	if _, err := os.Stat(strings.TrimSpace(contract.RunAt(t, f, paths[1], nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)); err == nil {
		t.Fatal("worktree lease was not removed on release")
	}
	if _, err := os.Stat(filepath.Join(paths[1], "dirty.txt")); err == nil {
		t.Fatal("worktree release did not clean dirty files")
	}
	if _, err := os.Stat(filepath.Join(paths[1], "ignored", "leak.txt")); err == nil {
		t.Fatal("worktree release did not clean ignored artifacts")
	}
}

func testRuntimeWorktreeLeaseHardening(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	f.WriteExecutable("rec-shell", "#!/usr/bin/env bash\n: \"${BENCH_WT_RECORD:?}\"\npwd >> \"$BENCH_WT_RECORD\"\n")
	record := filepath.Join(f.Root, "paths")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "SHELL": filepath.Join(f.Root, "rec-shell")}
	runWT := func() { f.BenchEnv(env, "worktree").RequireExit(0) }
	runWT()
	p := strings.TrimSpace(contract.ReadFileAbs(t, record))
	lease := strings.TrimSpace(contract.RunAt(t, f, p, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, lease, "4194300 2020-01-01T00:00:00Z\n")
	runWT()
	contract.WriteFileAbs(t, lease, fmt.Sprintf("%d 2026-01-01T00:00:00Z\n", os.Getpid()))
	runWT()
	if _, err := os.Stat(lease); err != nil {
		t.Fatal("live foreign lease was removed by a foreign release")
	}
	contract.WriteFileAbs(t, lease, "")
	old := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	if err := os.Chtimes(lease, old, old); err != nil {
		t.Fatalf("age empty lease: %v", err)
	}
	runWT()
	contract.WriteFileAbs(t, lease, "")
	runWT()
	contract.Remove(t, lease)
	paths := contract.NonEmptyLines(contract.ReadFileAbs(t, record))
	contract.RequireIntEqual(t, len(paths), 5, "expected five worktree runs")
	requireEqual(t, paths[1], p, "dead-pid lease was not reclaimed")
	if paths[2] == p {
		t.Fatal("live foreign lease was stolen")
	}
	requireEqual(t, paths[3], p, "aged-out empty lease was not reclaimed")
	if paths[4] == p {
		t.Fatal("fresh empty lease was stolen")
	}
}

func testRuntimeWorktreeConcurrentAcquire(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	f.WriteExecutable("rv-shell", `#!/usr/bin/env bash
: "${BENCH_WT_RECORD:?}"
pwd >> "$BENCH_WT_RECORD"
for _ in $(seq 100); do
  [ "$(grep -c . "$BENCH_WT_RECORD" 2>/dev/null)" -ge 2 ] && exit 0
  sleep 0.1
done
exit 0
`)
	record := filepath.Join(f.Root, "paths")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "SHELL": filepath.Join(f.Root, "rv-shell")}
	done := make(chan contract.Probe, 2)
	for i := 0; i < 2; i++ {
		go func() { done <- f.BenchEnv(env, "worktree") }()
	}
	for i := 0; i < 2; i++ {
		(<-done).RequireExit(0)
	}
	paths := contract.NonEmptyLines(contract.ReadFileAbs(t, record))
	sort.Strings(paths)
	contract.RequireIntEqual(t, len(paths), 2, "concurrent worktree runs did not both complete")
	if paths[0] == paths[1] {
		t.Fatal("concurrent acquires shared a worktree")
	}
}
