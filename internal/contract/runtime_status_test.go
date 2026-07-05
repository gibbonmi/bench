package contract

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRuntimeStatusContracts(t *testing.T) {
	skipIfSubjectBenchMissing(t)
	t.Parallel()
	runParallel(t, "bench idea/roadmap contract", testRuntimeIdeaRoadmap)
	runParallel(t, "bench status clean contract", testRuntimeStatusClean)
	runParallel(t, "bench status footer contract", testRuntimeStatusFooter)
	runParallel(t, "bench status stale-gate contract", testRuntimeStatusStaleGate)
	runParallel(t, "bench status fresh-green contract", testRuntimeStatusFreshGreen)
	runParallel(t, "bench status decisions contract", testRuntimeStatusDecisions)
	runParallel(t, "bench status unresolved-maps count contract", testRuntimeStatusUnresolvedMapsCount)
	runParallel(t, "bench status budget contract", testRuntimeStatusBudget)
	runParallel(t, "bench status warm-pool contract", testRuntimeStatusWarmPool)
	runParallel(t, "bench status retirement-signal contract", testRuntimeStatusRetirementSignal)
	runParallel(t, "bench status learnings-floor contract", testRuntimeStatusLearningsFloor)
}

func testRuntimeIdeaRoadmap(t *testing.T) {
	f := NewFixture(t)

	if roadmap := f.Bench("roadmap"); !strings.Contains(strings.ToLower(roadmap.Stdout+roadmap.Stderr), "empty") {
		t.Fatalf("roadmap on absent file did not report empty\nstdout:\n%s\nstderr:\n%s", roadmap.Stdout, roadmap.Stderr)
	}
	f.Bench("idea", "ship dark mode").RequireExit(0)
	requireFileMatches(t, f, "ROADMAP.md", `(?m)^- [0-9]{4}-[0-9]{2}-[0-9]{2}  ship dark mode$`, "idea entry not dated")
	before := lineCount(f.ReadFile("ROADMAP.md"))
	if p := f.Bench("idea"); p.ExitCode == 0 {
		t.Fatalf("no-arg idea succeeded; should error\nstdout:\n%s\nstderr:\n%s", p.Stdout, p.Stderr)
	}
	requireIntEqual(t, lineCount(f.ReadFile("ROADMAP.md")), before, "no-arg idea appended a blank entry")
	f.Bench("roadmap").RequireContains(f.Bench("roadmap").Stdout, "ship dark mode")
	f.Bench("idea", "capture", "all", "the", "words").RequireExit(0)
	requireFileMatches(t, f, "ROADMAP.md", `(?m)^- [0-9]{4}-[0-9]{2}-[0-9]{2}  capture all the words$`, "idea did not join unquoted multi-word args")
	f.WriteFile("ROADMAP.md", "- 2026-06-01  hand added")
	f.Bench("idea", "after handedit").RequireExit(0)
	requireIntEqual(t, strings.Count(f.ReadFile("ROADMAP.md"), "- "), 2, "idea merged onto a newline-less last line")
	f.WriteFile("ROADMAP.md", "")
	f.Bench("roadmap").RequireContains(f.Bench("roadmap").Stdout, "empty")
}

func testRuntimeStatusClean(t *testing.T) {
	f := NewFixture(t)
	commitAllowEmpty(t, f, "init")
	if status := f.Bench("status"); !strings.Contains(status.Stdout, "clean — nothing pending") {
		t.Fatalf("clean repo did not report all-clear\nstdout:\n%s\nstderr:\n%s", status.Stdout, status.Stderr)
	}
}

func testRuntimeStatusFooter(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("ROADMAP.md", "- 2026-06-30  an idea\n")
	f.CommitAll("s")
	out := f.Bench("status").Stdout
	requireContains(t, out, "clean — nothing pending")
	requireContains(t, out, "parked — bench roadmap")
	if regexp.MustCompile(`(?m)^▶.*bench roadmap`).MatchString(out) {
		t.Fatal("roadmap footer became the lead")
	}
}

func testRuntimeStatusStaleGate(t *testing.T) {
	f := NewFixture(t)
	commitAllowEmpty(t, f, "init")
	writeFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"), "green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n")
	out := f.Bench("status").Stdout
	requireContains(t, out, "re-run the gate")
	requireNotContains(t, out, "clean — nothing pending")
}

func testRuntimeStatusFreshGreen(t *testing.T) {
	f := NewFixture(t)
	commitAllowEmpty(t, f, "init")
	writeFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"), fmt.Sprintf("green %s 2026-06-30T00:00:00Z\n", strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout)))
	f.Bench("status").RequireContains(f.Bench("status").Stdout, "clean — nothing pending")
}

func testRuntimeStatusDecisions(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("decisions/x.md", "### Answer\n— (deferred)\n")
	f.CommitAll("s")
	f.Bench("status").RequireContains(f.Bench("status").Stdout, "craft-grill → /bench-write-spec")
}

func testRuntimeStatusUnresolvedMapsCount(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("decisions/one.md", "## #1: a?\nType: Grill\n### Answer\n— (open)\n\n## #2: b?\nType: Grill\n### Answer\n— (deferred)\n")
	f.WriteFile("decisions/two.md", "## #1: c?\nType: Grill\n### Answer\n— (open)\n")
	f.Bench("status").RequireContains(f.Bench("status").Stdout, "2 unresolved map(s)")
}

func testRuntimeStatusBudget(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — a  [open]\n")
	f.WriteFile("big.py", strings.Repeat("x = \n", 401))
	f.WriteFile("decisions/x.md", "### Answer\n— (deferred)\n")
	f.CommitAll("s")
	f.WriteFile("dirty.txt", "dirty\n")
	f.Git("worktree", "add", "-q", "--detach", filepath.Join(f.Root, "wt2"), "HEAD")
	tree := strings.TrimSpace(f.Bench("tree-hash").Stdout)
	writeFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"), fmt.Sprintf("red %s 2026-06-30T00:00:00Z\n", tree))

	out := f.Bench("status").Stdout
	first := strings.SplitN(out, "\n", 2)[0]
	requireContains(t, first, "fix before commit")
	requireContains(t, out, "+1 more")
	requireContains(t, out, "/bench-integrate-learnings")
	requireContains(t, out, "split (craft-seams)")
	requireContains(t, out, "commit on green / push")
	requireContains(t, out, "resume or clean up")
	requireNotContains(t, out, "craft-grill → /bench-write-spec")
	if rows := countStatusRows(out); rows > 5 {
		t.Fatalf("budget exceeded five rows (%d):\n%s", rows, out)
	}
}

func testRuntimeStatusWarmPool(t *testing.T) {
	f := NewFixture(t)
	commitAllowEmpty(t, f, "init")
	benchHome := filepath.Join(f.Root, ".bh")
	repoRoot := strings.TrimSpace(f.Git("rev-parse", "--show-toplevel").Stdout)
	pool := filepath.Join(benchHome, "worktrees", filepath.Base(repoRoot)+"-"+cksum(t, repoRoot))
	warm := filepath.Join(pool, "warm")
	mkdir(t, pool)
	f.Git("worktree", "add", "-q", "--detach", warm, "HEAD")
	out := f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	requireNotContains(t, out, "resume or clean up")
	lease := strings.TrimSpace(runAt(t, f, warm, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	writeFileAbs(t, lease, "")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	requireContains(t, out, "1 active worktree")
	remove(t, lease)
	f.Git("branch", "worktree-agent-orphan")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	requireContains(t, out, "orphaned worktree branch")
}

func testRuntimeStatusRetirementSignal(t *testing.T) {
	f := NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	f.Git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	f.WriteFile("specs/done.md", "# Done\n\nStatus: implemented\n")
	f.WriteFile("specs/staged.md", "# Staged\n\nStatus: staged\n")
	f.WriteFile("specs/plain.md", "# Plain\n\nno status line here\n")
	f.WriteFile("specs/fenced.md", "# Fenced\n\nexample:\n\n```\nStatus: implemented\n```\n")
	f.CommitAll("init")
	out := f.Bench("status").Stdout
	requireContains(t, out, "1 merged spec(s) awaiting retirement")
	requireContains(t, out, "promote-then-delete (spec-retire)")
	f.WriteFile("scratch.txt", "scratch\n")
	out = f.Bench("status").Stdout
	requireContains(t, out, "awaiting retirement")
	if strings.Contains(strings.SplitN(out, "\n", 2)[0], "spec-retire") {
		t.Fatal("retirement signal wrongly led over the git signal")
	}
	remove(t, filepath.Join(f.Root, "scratch.txt"))
	f.Git("checkout", "-q", "-b", "feature")
	requireNotContains(t, f.Bench("status").Stdout, "awaiting retirement")
	f.Git("checkout", "-q", "main")
	f.Git("rm", "-q", "specs/done.md")
	f.CommitAll("retire")
	requireNotContains(t, f.Bench("status").Stdout, "awaiting retirement")
}

func testRuntimeStatusLearningsFloor(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — a  [open]\n")
	f.CommitAll("s")
	requireNotContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "2"}, "status").Stdout, "/bench-integrate-learnings")
	requireContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "1"}, "status").Stdout, "/bench-integrate-learnings")
	f.WriteFile(".bench/learnings.md", "## <date> — <short title>  [open]\n")
	requireNotContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "1"}, "status").Stdout, "/bench-integrate-learnings")
}
