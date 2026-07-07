package runtime

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/contract"
	"path/filepath"
	"strings"
	"testing"
)

func TestRuntimeStatusContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench idea/roadmap contract", testRuntimeIdeaRoadmap)
	contract.RunParallel(t, "bench status clean contract", testRuntimeStatusClean)
	contract.RunParallel(t, "bench status drain-row contract", testRuntimeStatusDrainRow)
	contract.RunParallel(t, "bench status stale-gate contract", testRuntimeStatusStaleGate)
	contract.RunParallel(t, "bench status stale-gate drift classification contract", testRuntimeStatusStaleGateDriftClassification)
	contract.RunParallel(t, "bench status fresh-green contract", testRuntimeStatusFreshGreen)
	contract.RunParallel(t, "bench status decisions contract", testRuntimeStatusDecisions)
	contract.RunParallel(t, "bench status unresolved-maps count contract", testRuntimeStatusUnresolvedMapsCount)
	contract.RunParallel(t, "bench status budget contract", testRuntimeStatusBudget)
	contract.RunParallel(t, "bench status warm-pool contract", testRuntimeStatusWarmPool)
	contract.RunParallel(t, "bench status retirement-signal contract", testRuntimeStatusRetirementSignal)
	contract.RunParallel(t, "bench status orphaned-pickup contract", testRuntimeStatusOrphanedPickup)
	contract.RunParallel(t, "bench status learnings-floor contract", testRuntimeStatusLearningsFloor)
}

func testRuntimeIdeaRoadmap(t *testing.T) {
	f := contract.NewFixture(t)

	absent := f.Bench("roadmap")
	absent.RequireExit(0)
	contract.RequireContains(t, absent.Stdout, "no ROADMAP.md")
	contract.RequireContains(t, absent.Stdout, "/bench-what-next")
	for _, args := range [][]string{nil, {""}, {"   ", "\t"}} {
		blank := contract.NewFixture(t)
		p := blank.Bench(append([]string{"idea"}, args...)...)
		p.RequireExit(2)
		contract.RequireContains(t, p.Stdout, "usage: bench idea \"<text>\"")
		if blank.Exists("IDEAS.md") {
			t.Fatalf("args %q: empty idea created IDEAS.md", args)
		}
	}
	f.Bench("idea", "ship dark mode").RequireExit(0)
	contract.RequireFileMatches(t, f, "IDEAS.md", `(?m)^- [0-9]{4}-[0-9]{2}-[0-9]{2}  ship dark mode$`, "idea entry not dated")
	if f.Exists("ROADMAP.md") {
		t.Fatal("idea created ROADMAP.md; capture should write only IDEAS.md")
	}
	contract.Mkdir(t, filepath.Join(f.Root, "sub"))
	contract.RunAt(t, f, filepath.Join(f.Root, "sub"), nil, "bash", benchPath(t), "idea", "from sub").RequireExit(0)
	contract.RequireFileMatches(t, f, "IDEAS.md", `(?m)^- [0-9]{4}-[0-9]{2}-[0-9]{2}  from sub$`, "idea from nested cwd did not append to root IDEAS.md")
	if f.Exists("sub/IDEAS.md") {
		t.Fatal("idea from nested cwd created sub/IDEAS.md")
	}
	before := contract.LineCount(f.ReadFile("IDEAS.md"))
	for _, args := range [][]string{nil, {""}, {"   ", "\t"}} {
		p := f.Bench(append([]string{"idea"}, args...)...)
		p.RequireExit(2)
		contract.RequireContains(t, p.Stdout, "usage: bench idea \"<text>\"")
	}
	contract.RequireIntEqual(t, contract.LineCount(f.ReadFile("IDEAS.md")), before, "empty idea appended a blank entry")
	f.Bench("idea", "capture", "all", "the", "words").RequireExit(0)
	contract.RequireFileMatches(t, f, "IDEAS.md", `(?m)^- [0-9]{4}-[0-9]{2}-[0-9]{2}  capture all the words$`, "idea did not join unquoted multi-word args")
	f.WriteFile("IDEAS.md", "- 2026-06-01  hand added")
	f.Bench("idea", "after handedit").RequireExit(0)
	contract.RequireIntEqual(t, strings.Count(f.ReadFile("IDEAS.md"), "- "), 2, "idea merged onto a newline-less last line")
	f.WriteFile("ROADMAP.md", "")
	zero := f.Bench("roadmap")
	zero.RequireExit(0)
	contract.RequireContains(t, zero.Stdout, "no ROADMAP.md")
	contract.RequireContains(t, zero.Stdout, "/bench-what-next")

	drain := contract.NewFixture(t)
	drain.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n")
	drain.WriteFile("IDEAS.md", "- 2026-07-05  parked idea\n")
	drain.WriteFile(".bench/learnings.md", "## 2026-07-05 — open learning  [open]\n")
	roadmap := drain.Bench("roadmap")
	roadmap.RequireExit(0)
	contract.RequireContains(t, roadmap.Stdout, "# Roadmap")
	contract.RequireContains(t, roadmap.Stdout, "ideas: 1 parked in IDEAS.md")
	contract.RequireContains(t, roadmap.Stdout, "learnings: 1 open in .bench/learnings.md")
	contract.RequireContains(t, roadmap.Stdout, "/bench-what-next")

	next := contract.NewFixture(t)
	next.WriteFile("ROADMAP.md", "# Roadmap\n\n## Context\n\nKeep current.\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n2. Implement next item - /bench-implement-spec\n")
	out := next.Bench("roadmap")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "## Next action\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n2. Implement next item - /bench-implement-spec\n")
	contract.RequireNotContains(t, out.Stdout, "## Drain status")

	headless := contract.NewFixture(t)
	headless.WriteFile("ROADMAP.md", "# Roadmap\n\n## Context\n\nNo sequence here.\n")
	noSection := headless.Bench("roadmap")
	noSection.RequireExit(0)
	contract.RequireContains(t, noSection.Stdout, "no ## Recommended sequence section")
	contract.RequireContains(t, noSection.Stdout, "/bench-what-next")

	malformed := contract.NewFixture(t)
	malformed.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. Only item - /bench-shape-idea\n")
	short := malformed.Bench("roadmap")
	short.RequireExit(0)
	contract.RequireContains(t, short.Stdout, "malformed ## Recommended sequence: 1 numbered item(s)")
	contract.RequireContains(t, short.Stdout, "/bench-what-next")
}

func testRuntimeStatusClean(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	if status := f.Bench("status"); !strings.Contains(status.Stdout, "clean — nothing pending") {
		t.Fatalf("clean repo did not report all-clear\nstdout:\n%s\nstderr:\n%s", status.Stdout, status.Stderr)
	}
}

func testRuntimeStatusDrainRow(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n")
	f.WriteFile("IDEAS.md", "- 2026-07-05  parked idea\n")
	// Template heading + one real open heading: the shared parser counts only the real one.
	f.WriteFile(".bench/learnings.md", "## <date> — <short title>  [open]\n## 2026-07-05 — open learning  [open]\n")
	f.CommitAll("s")
	out := f.Bench("status").Stdout
	contract.RequireContains(t, out, "1 idea(s), 1 open learning(s)")
	contract.RequireContains(t, out, "/bench-what-next")
	contract.RequireNotContains(t, out, "/bench-integrate-learnings")
	contract.RequireNotContains(t, out, "parked — bench roadmap")
	contract.RequireNotContains(t, out, "clean — nothing pending")
	if n := strings.Count(out, "→ /bench-what-next"); n != 1 {
		t.Fatalf("want one combined maintenance row, got %d in:\n%s", n, out)
	}
	if !strings.HasPrefix(out, "▶ /bench-what-next  (drain)") {
		t.Fatalf("drain row as the only signal must lead, got:\n%s", out)
	}

	// A working roadmap alone is not pending capture: no drain row, board stays clean.
	clean := contract.NewFixture(t)
	clean.WriteFile("ROADMAP.md", "# Roadmap\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n")
	clean.CommitAll("s")
	out = clean.Bench("status").Stdout
	contract.RequireContains(t, out, "clean — nothing pending")
	contract.RequireNotContains(t, out, "/bench-what-next")
}

func testRuntimeStatusStaleGate(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"), "green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n")
	out := f.Bench("status").Stdout
	contract.RequireContains(t, out, "re-run the gate")
	contract.RequireContains(t, out, "stale (gated tree")
	contract.RequireNotContains(t, out, "capture-only drift")
	contract.RequireNotContains(t, out, "clean — nothing pending")
}

func testRuntimeStatusFreshGreen(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"), fmt.Sprintf("green %s 2026-06-30T00:00:00Z\n", strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout)))
	f.Bench("status").RequireContains(f.Bench("status").Stdout, "clean — nothing pending")
}

func testRuntimeStatusDecisions(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/x.md", "### Answer\n— (deferred)\n")
	f.CommitAll("s")
	f.Bench("status").RequireContains(f.Bench("status").Stdout, "craft-grill → /bench-write-spec")
}

func testRuntimeStatusUnresolvedMapsCount(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/one.md", "## #1: a?\nType: Grill\n### Answer\n— (open)\n\n## #2: b?\nType: Grill\n### Answer\n— (deferred)\n")
	f.WriteFile("decisions/two.md", "## #1: c?\nType: Grill\n### Answer\n— (open)\n")
	f.Bench("status").RequireContains(f.Bench("status").Stdout, "2 unresolved map(s)")
}

func testRuntimeStatusBudget(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — a  [open]\n")
	f.WriteFile("big.py", strings.Repeat("x = \n", 401))
	f.WriteFile("decisions/x.md", "### Answer\n— (deferred)\n")
	f.CommitAll("s")
	f.WriteFile("dirty.txt", "dirty\n")
	f.Git("worktree", "add", "-q", "--detach", filepath.Join(f.Root, "wt2"), "HEAD")
	tree := strings.TrimSpace(f.Bench("tree-hash").Stdout)
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"), fmt.Sprintf("red %s 2026-06-30T00:00:00Z\n", tree))

	out := f.Bench("status").Stdout
	first := strings.SplitN(out, "\n", 2)[0]
	contract.RequireContains(t, first, "fix before commit")
	contract.RequireContains(t, out, "+1 more")
	contract.RequireContains(t, out, "/bench-what-next")
	contract.RequireContains(t, out, "split (craft-seams)")
	contract.RequireContains(t, out, "commit on green / push")
	contract.RequireContains(t, out, "bench worktree clean")
	contract.RequireNotContains(t, out, "craft-grill → /bench-write-spec")
	if rows := countStatusRows(out); rows > 5 {
		t.Fatalf("budget exceeded five rows (%d):\n%s", rows, out)
	}
}

func testRuntimeStatusWarmPool(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	benchHome := filepath.Join(f.Root, ".bh")
	pool := addRuntimePoolWorktrees(t, f, benchHome)
	outOfPool := filepath.Join(f.Root, "outside pool")
	f.Git("worktree", "add", "-q", "--detach", outOfPool, "HEAD")
	out := f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireContains(t, out, "1 out-of-pool worktree")
	contract.RequireContains(t, out, "bench worktree clean")
	contract.RequireContains(t, out, "1 leased pool worktree")
	requireStatusLineNotContains(t, out, "1 leased pool worktree", "bench worktree clean")
	leasedOut := contract.RunAt(t, f, pool.Leased, map[string]string{"BENCH_HOME": benchHome}, "bash", benchPath(t), "status").Stdout
	contract.RequireContains(t, leasedOut, "1 out-of-pool worktree")
	contract.RequireContains(t, leasedOut, "1 leased pool worktree")
	requireStatusLineNotContains(t, leasedOut, "1 leased pool worktree", "bench worktree clean")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireContains(t, out, "1 out-of-pool worktree")
	contract.RequireContains(t, out, "1 leased pool worktree")
	contract.RequireContains(t, out, "resume leased worktree")
	requireStatusLineNotContains(t, out, "1 leased pool worktree", "bench worktree clean")
	contract.Remove(t, pool.LeaseFile)
	f.Git("branch", "worktree-agent-orphan")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireContains(t, out, "1 out-of-pool worktree")
	contract.RequireContains(t, out, "orphaned worktree branch")
	contract.RequireContains(t, out, "delete scratch branch")
	requireStatusLineNotContains(t, out, "orphaned worktree branch", "bench worktree clean")
}

func testRuntimeStatusRetirementSignal(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	f.Git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	f.WriteFile("specs/done.md", "# Done\n\nStatus: implemented\n")
	f.WriteFile("specs/staged.md", "# Staged\n\nStatus: staged\n")
	f.WriteFile("specs/plain.md", "# Plain\n\nno status line here\n")
	f.WriteFile("specs/fenced.md", "# Fenced\n\nexample:\n\n```\nStatus: implemented\n```\n")
	f.CommitAll("init")
	out := f.Bench("status").Stdout
	contract.RequireContains(t, out, "1 merged spec(s) awaiting retirement")
	contract.RequireContains(t, out, "promote-then-delete (spec-retire)")
	f.WriteFile("scratch.txt", "scratch\n")
	out = f.Bench("status").Stdout
	contract.RequireContains(t, out, "awaiting retirement")
	if strings.Contains(strings.SplitN(out, "\n", 2)[0], "spec-retire") {
		t.Fatal("retirement signal wrongly led over the git signal")
	}
	contract.Remove(t, filepath.Join(f.Root, "scratch.txt"))
	f.Git("checkout", "-q", "-b", "feature")
	contract.RequireNotContains(t, f.Bench("status").Stdout, "awaiting retirement")
	f.Git("checkout", "-q", "main")
	f.Git("rm", "-q", "specs/done.md")
	f.CommitAll("retire")
	contract.RequireNotContains(t, f.Bench("status").Stdout, "awaiting retirement")
}

func testRuntimeStatusOrphanedPickup(t *testing.T) {
	// A review pickup with no matching spec is flagged with a clean-up action.
	f := contract.NewFixture(t)
	f.WriteFile("reviews/x.md", "# review of x\n")
	f.CommitAll("orphan pickup")
	out := f.Bench("status").Stdout
	contract.RequireContains(t, out, "orphaned review pickup")
	contract.RequireContains(t, out, "promote or delete by hand")

	// A paired pickup (its spec still present) is expected state — no row fires.
	paired := contract.NewFixture(t)
	paired.WriteFile("reviews/x.md", "# review of x\n")
	paired.WriteFile("specs/x.md", "# x\n\nStatus: staged\n")
	paired.CommitAll("paired")
	contract.RequireNotContains(t, paired.Bench("status").Stdout, "orphaned review pickup")
}

func testRuntimeStatusLearningsFloor(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — a  [open]\n")
	f.CommitAll("s")
	contract.RequireNotContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "2"}, "status").Stdout, "/bench-what-next")
	contract.RequireContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "1"}, "status").Stdout, "/bench-what-next")
	f.WriteFile(".bench/learnings.md", "## <date> — <short title>  [open]\n")
	contract.RequireNotContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "1"}, "status").Stdout, "/bench-what-next")
	// The floor gates only the learnings component; parked ideas always count, and a
	// real open heading below the floor renders as zero rather than leaking through.
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — a  [open]\n")
	f.Bench("idea", "parked past the floor").RequireExit(0)
	out := f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "2"}, "status").Stdout
	contract.RequireContains(t, out, "1 idea(s), 0 open learning(s)")
	contract.RequireContains(t, out, "/bench-what-next")
}

func requireStatusLineNotContains(t testing.TB, out, needle, forbidden string) {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, needle) {
			contract.RequireNotContains(t, line, forbidden)
			return
		}
	}
	t.Fatalf("status line containing %q not found in:\n%s", needle, out)
}
