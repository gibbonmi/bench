package runtime

import (
	"fmt"
	"os"
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
	contract.RunParallel(t, "bench worktree clean removes out-of-pool contract", testRuntimeWorktreeCleanRemovesOutOfPool)
	contract.RunParallel(t, "bench worktree clean salvages dirty onto branch contract", testRuntimeWorktreeCleanSalvagesDirtyOntoBranch)
	contract.RunParallel(t, "bench worktree clean refuses dirty detached contract", testRuntimeWorktreeCleanRefusesDirtyDetached)
	contract.RunParallel(t, "bench worktree clean stale registration prune contract", testRuntimeWorktreeCleanPrunesMissingRegistration)
	contract.RunParallel(t, "bench worktree clean pool cwd classification contract", testRuntimeWorktreeCleanFromPoolCwd)
	contract.RunParallel(t, "bench worktree usage contract", testRuntimeWorktreeRejectsUnknownArgs)
	contract.RunParallel(t, "bench worktree lease/reuse contract", testRuntimeWorktreeLeaseReuse)
	contract.RunParallel(t, "bench worktree lease hardening contract", testRuntimeWorktreeLeaseHardening)
	contract.RunParallel(t, "bench worktree concurrent-acquire contract", testRuntimeWorktreeConcurrentAcquire)
	contract.RunParallel(t, "bench worktree clean sweeps merged orphan (non-tty) contract", testRuntimeWorktreeSweepDeletesMerged)
	contract.RunParallel(t, "bench worktree clean sweeps content-landed orphan contract", testRuntimeWorktreeSweepDeletesContentLanded)
	contract.RunParallel(t, "bench worktree clean keeps unmerged orphan contract", testRuntimeWorktreeSweepKeepsUnmerged)
	contract.RunParallel(t, "bench worktree clean keeps unique-patch orphan contract", testRuntimeWorktreeSweepKeepsUniquePatch)
	contract.RunParallel(t, "bench worktree clean keeps evil-merge orphan contract", testRuntimeWorktreeSweepKeepsEvilMerge)
	contract.RunParallel(t, "bench worktree clean spares active and non-scratch contract", testRuntimeWorktreeSweepSparesProtected)
	contract.RunParallel(t, "bench worktree clean deletes slashed unicode orphan contract", testRuntimeWorktreeSweepDeletesSlashUnicode)
	contract.RunParallel(t, "bench worktree clean refuses unresolvable default contract", testRuntimeWorktreeSweepRefusesUnresolvableDefault)
}

// onMainFixture returns a fixture whose default branch resolves to a real `main` commit — the
// sweep's happy-path precondition. git.DefaultBranch falls back to "main", but a bare `git init`
// fixture is born on "master", so an explicit HEAD symref lands the first commit on main and makes
// the default branch resolve.
func onMainFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	commitAllowEmpty(t, f, "init")
	return f
}

// headExists reports whether refs/heads/<name> resolves — an exit-code probe so a slashed or
// unicode branch name is checked verbatim, never through a substring or glob match.
func headExists(f contract.Fixture, name string) bool {
	return f.GitAllow("show-ref", "--verify", "--quiet", "refs/heads/"+name).ExitCode == 0
}

// Story 1 + 5: a merged worktree-* orphan is deleted with its line on stdout, run non-interactively
// (no PTY) with no out-of-pool worktrees — proving the sweep is not coupled to the TTY confirmation.
func testRuntimeWorktreeSweepDeletesMerged(t *testing.T) {
	f := onMainFixture(t)
	f.Git("branch", "worktree-agent-merged")

	out := f.Bench("worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "deleted branch worktree-agent-merged")
	if headExists(f, "worktree-agent-merged") {
		t.Fatalf("merged scratch orphan survived the sweep:\n%s", f.Git("for-each-ref", "--format=%(refname:short)", "refs/heads/").Stdout)
	}
}

// FT44 story 3: a worktree-* orphan whose commit was cherry-picked into the default branch is a
// non-ancestor whose every patch already landed — `git cherry` all `-`. The sweep deletes it and
// says which proof landed it, so a hand-inspection of a provably-empty branch never recurs.
func testRuntimeWorktreeSweepDeletesContentLanded(t *testing.T) {
	f := onMainFixture(t)
	f.Git("checkout", "-q", "-b", "worktree-agent-landed")
	f.WriteFile("work.txt", "landed work\n")
	f.CommitAll("real work")
	f.Git("checkout", "-q", "main")
	// A divergent commit before the cherry-pick guarantees the branch tip is not an
	// ancestor of main, so only the patch-containment proof can land it.
	commitAllowEmpty(t, f, "diverge")
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "cherry-pick", "worktree-agent-landed")

	out := f.Bench("worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "deleted branch worktree-agent-landed (landed by content)")
	if headExists(f, "worktree-agent-landed") {
		t.Fatalf("content-landed scratch orphan survived the sweep:\n%s", f.Git("for-each-ref", "--format=%(refname:short)", "refs/heads/").Stdout)
	}
}

// Story 2 + 4: an unmerged worktree-* orphan survives with the verbatim kept line, and keeping it
// exits 0 — a branch carrying unique commits is never destroyed and keeping is not a failure.
func testRuntimeWorktreeSweepKeepsUnmerged(t *testing.T) {
	f := onMainFixture(t)
	f.Git("checkout", "-q", "-b", "worktree-agent-unmerged")
	commitAllowEmpty(t, f, "unique work")
	f.Git("checkout", "-q", "main")

	out := f.Bench("worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "kept branch worktree-agent-unmerged (unique commits — inspect or delete by hand)")
	if !headExists(f, "worktree-agent-unmerged") {
		t.Fatal("unmerged scratch orphan was destroyed by the conservative default")
	}
}

// FT44 story 3, destructive direction: a worktree-* orphan carrying a real unique patch (not just
// an empty commit) reports `+` under `git cherry` and must survive the content-landed proof.
func testRuntimeWorktreeSweepKeepsUniquePatch(t *testing.T) {
	f := onMainFixture(t)
	f.Git("checkout", "-q", "-b", "worktree-agent-unique")
	f.WriteFile("unique.txt", "unique work\n")
	f.CommitAll("unique work")
	f.Git("checkout", "-q", "main")

	out := f.Bench("worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "kept branch worktree-agent-unique (unique commits — inspect or delete by hand)")
	if !headExists(f, "worktree-agent-unique") {
		t.Fatal("unique-patch scratch orphan was destroyed by the content-landed proof")
	}
}

// FT44 review fix (merge blindness): `git cherry` enumerates only non-merge commits, so a branch
// whose ordinary commits all landed but whose merge commit carries unique content (a conflict
// resolution) is invisible to the patch proof. Such a branch must survive the sweep: the ordering
// matters — the branch merges main first, then its feature commit squash-lands, so cherry reports
// the feature commit `-` while the merge-only file exists nowhere in main.
func testRuntimeWorktreeSweepKeepsEvilMerge(t *testing.T) {
	f := onMainFixture(t)
	f.Git("checkout", "-q", "-b", "worktree-agent-evil")
	f.WriteFile("feat.txt", "feature\n")
	f.CommitAll("feature work")
	f.Git("checkout", "-q", "main")
	f.WriteFile("other.txt", "other\n")
	f.CommitAll("mainline work")
	f.Git("checkout", "-q", "worktree-agent-evil")
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "merge", "--no-commit", "--no-ff", "main")
	f.WriteFile("evil.txt", "merge-only resolution\n")
	f.CommitAll("merge main (evil)")
	f.Git("checkout", "-q", "main")
	// Squash-land the feature commit (the merge commit's first parent) after the merge.
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "cherry-pick", "worktree-agent-evil~1")

	out := f.Bench("worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "kept branch worktree-agent-evil (unique commits — inspect or delete by hand)")
	if !headExists(f, "worktree-agent-evil") {
		t.Fatal("orphan carrying merge-only content was deleted by the content-landed proof")
	}
}

// Story 3: a scratch branch checked out in a live worktree, a bench/shift-* review branch, a plain
// branch, and main all survive the sweep — the active-worktree and non-scratch filters hold. The
// live out-of-pool worktree makes the worktree-removal phase refuse under no TTY (exit 1), which is
// orthogonal to the sweep; this row asserts only branch survival.
func testRuntimeWorktreeSweepSparesProtected(t *testing.T) {
	f := onMainFixture(t)
	active := filepath.Join(f.Root, "..", "active-wt")
	f.Git("worktree", "add", "-q", "-b", "worktree-agent-active", active, "HEAD")
	f.Git("branch", "bench/shift-review")
	f.Git("branch", "plain-branch")

	f.Bench("worktree", "clean")
	for _, name := range []string{"worktree-agent-active", "bench/shift-review", "plain-branch", "main"} {
		if !headExists(f, name) {
			t.Fatalf("sweep deleted protected branch %s:\n%s", name, f.Git("for-each-ref", "--format=%(refname:short)", "refs/heads/").Stdout)
		}
	}
}

// Story 1 edge: a merged scratch name carrying a slash and a non-ASCII character inside the
// worktree- prefix is detected and deleted — the ref name is neither mangled nor lost to a glob that
// stops at a slash.
func testRuntimeWorktreeSweepDeletesSlashUnicode(t *testing.T) {
	f := onMainFixture(t)
	const name = "worktree-agent-café/x"
	f.Git("branch", name)

	out := f.Bench("worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "deleted branch "+name)
	if headExists(f, name) {
		t.Fatalf("merged slashed/unicode scratch orphan survived the sweep:\n%s", f.Git("for-each-ref", "--format=%(refname:short)", "refs/heads/").Stdout)
	}
}

// Story 7: when the resolved default branch does not resolve to a commit, the sweep refuses loudly on
// stderr, deletes nothing, and exits 1 — the false-empty guard. The fixture is born on "master" with
// no origin/HEAD, so git.DefaultBranch's "main" fallback names a ref that never resolves.
func testRuntimeWorktreeSweepRefusesUnresolvableDefault(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	f.Git("branch", "worktree-agent-x")

	out := f.Bench("worktree", "clean")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stderr, "cannot resolve the default branch")
	contract.RequireNotContains(t, out.Stdout, "deleted branch")
	if !headExists(f, "worktree-agent-x") {
		t.Fatal("sweep deleted a branch despite an unresolvable default branch")
	}
}

func testRuntimeWorktreeCleanRemovesOutOfPool(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	candidate := filepath.Join(f.Root, "..", "outside wt [one]")
	f.Git("worktree", "add", "-q", "--detach", candidate, "HEAD")
	nested := filepath.Join(f.Root, "sub", "dir")
	contract.Mkdir(t, nested)

	out := contract.RunAt(t, f, nested, nil, "bash", benchPath(t), "worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, candidate)
	contract.RequireContains(t, out.Stdout, "removed")
	if strings.Contains(out.Stdout+out.Stderr, "--force") {
		t.Fatalf("cleanup mentioned forced removal:\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
	if strings.Contains(f.Git("worktree", "list", "--porcelain").Stdout, candidate) {
		t.Fatalf("cleanup left worktree registered:\n%s", f.Git("worktree", "list", "--porcelain").Stdout)
	}

	rerun := contract.RunAt(t, f, nested, nil, "bash", benchPath(t), "worktree", "clean")
	rerun.RequireExit(0)
	contract.RequireContains(t, rerun.Stdout, "nothing to clean")
}

func testRuntimeWorktreeCleanSalvagesDirtyOntoBranch(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	dirty := filepath.Join(f.Root, "..", "dirty on branch")
	f.Git("worktree", "add", "-q", "-b", "salvaged-wip", dirty, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(dirty, "dirty.txt"), "dirty\n")

	out := contract.RunAt(t, f, f.Root, nil, "bash", benchPath(t), "worktree", "clean")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "salvaged")
	contract.RequireContains(t, out.Stdout, "salvaged-wip")
	contract.RequireContains(t, out.Stdout, "removed")
	if strings.Contains(f.Git("worktree", "list", "--porcelain").Stdout, dirty) {
		t.Fatalf("salvage left worktree registered:\n%s", f.Git("worktree", "list", "--porcelain").Stdout)
	}
	show := f.Git("show", "salvaged-wip:dirty.txt")
	contract.RequireContains(t, show.Stdout, "dirty")
}

func testRuntimeWorktreeCleanRefusesDirtyDetached(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	dirty := filepath.Join(f.Root, "..", "dirty outside")
	f.Git("worktree", "add", "-q", "--detach", dirty, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(dirty, "dirty.txt"), "dirty\n")

	out := contract.RunAt(t, f, f.Root, nil, "bash", benchPath(t), "worktree", "clean")
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
	// Overlap is guaranteed by a test-owned barrier, not a timed poll: each shell records
	// its worktree and then holds until the test drops the go-file, which the test creates
	// only after seeing both records. A capped self-poll here is a schedule race — under
	// full-gate load the second spawn can outlive the first run's window, the pool hands
	// the released worktree back, and by-design reuse reads as a shared acquire (the FT37
	// flake). The shell's own loop is only a leak backstop, and it exits loud.
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	f.WriteExecutable("rv-shell", `#!/usr/bin/env bash
: "${BENCH_WT_RECORD:?}" "${BENCH_WT_GO:?}"
pwd >> "$BENCH_WT_RECORD"
for _ in $(seq 600); do
  [ -e "$BENCH_WT_GO" ] && exit 0
  sleep 0.1
done
exit 1
`)
	record := filepath.Join(f.Root, "paths")
	goFile := filepath.Join(f.Root, "go-file")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "BENCH_WT_GO": goFile, "SHELL": filepath.Join(f.Root, "rv-shell")}
	done := make(chan contract.Probe, 2)
	for i := 0; i < 2; i++ {
		go func() { done <- f.BenchEnv(env, "worktree") }()
	}
	backstop := time.Now().Add(2 * time.Minute)
	var overlapDeadline time.Time
	finished := make([]contract.Probe, 0, 2)
	for {
		raw, _ := os.ReadFile(record)
		lines := contract.NonEmptyLines(string(raw))
		if len(lines) >= 2 {
			break
		}
		if overlapDeadline.IsZero() && len(lines) >= 1 {
			// First record confirms one run has acquired and is holding — arm the overlap
			// window from this event rather than from spawn, so spawn latency under load
			// never counts against the assertion this window actually carries.
			overlapDeadline = time.Now().Add(time.Minute)
		}
	drainDone:
		for len(finished) < 2 {
			select {
			case p := <-done:
				finished = append(finished, p)
			default:
				break drainDone
			}
		}
		if len(finished) == 2 {
			// An exited run can never record (the shell records before it holds), but the
			// second record may land between the read above and the exits — re-read once.
			raw, _ := os.ReadFile(record)
			if len(contract.NonEmptyLines(string(raw))) < 2 {
				t.Fatal("both worktree runs exited with fewer than two acquires recorded — the runs never overlapped")
			}
			break
		}
		if overlapDeadline.IsZero() {
			if time.Now().After(backstop) {
				contract.WriteFileAbs(t, goFile, "") // release the straggler before failing
				t.Fatal("no bench worktree run recorded within 2m — acquire appears wedged")
			}
		} else if time.Now().After(overlapDeadline) {
			contract.WriteFileAbs(t, goFile, "") // release the straggler before failing
			t.Fatal("second acquire did not record within 60s of the first — the runs never overlapped")
		}
		time.Sleep(50 * time.Millisecond)
	}
	contract.WriteFileAbs(t, goFile, "")
	for len(finished) < 2 {
		finished = append(finished, <-done)
	}
	for _, p := range finished {
		p.RequireExit(0)
	}
	paths := contract.NonEmptyLines(contract.ReadFileAbs(t, record))
	sort.Strings(paths)
	contract.RequireIntEqual(t, len(paths), 2, "concurrent worktree runs did not both complete")
	if paths[0] == paths[1] {
		t.Fatal("concurrent acquires shared a worktree")
	}
}
