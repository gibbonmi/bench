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
	contract.RunParallel(t, "bench status learnings-floor contract", testRuntimeStatusLearningsFloor)
}

func testRuntimeIdeaRoadmap(t *testing.T) {
	f := contract.NewFixture(t)

	if roadmap := f.Bench("roadmap"); !strings.Contains(strings.ToLower(roadmap.Stdout+roadmap.Stderr), "empty") {
		t.Fatalf("roadmap on absent file did not report empty\nstdout:\n%s\nstderr:\n%s", roadmap.Stdout, roadmap.Stderr)
	}
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
	f.Bench("roadmap").RequireContains(f.Bench("roadmap").Stdout, "empty")

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
	f.WriteFile(".bench/learnings.md", "## 2026-07-05 — open learning  [open]\n")
	f.CommitAll("s")
	out := f.Bench("status").Stdout
	contract.RequireContains(t, out, "1 idea(s), 1 open learning(s)")
	contract.RequireContains(t, out, "/bench-what-next")
	contract.RequireNotContains(t, out, "/bench-integrate-learnings")
	contract.RequireNotContains(t, out, "parked — bench roadmap")
	if n := strings.Count(out, "→ /bench-what-next"); n != 1 {
		t.Fatalf("want one combined maintenance row, got %d in:\n%s", n, out)
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

func testRuntimeStatusStaleGateDriftClassification(t *testing.T) {
	benign := statusRowExpectation{
		contains:    []string{"stale (capture-only drift)", "re-run when convenient"},
		notContains: []string{"stale (gated tree"},
	}
	strong := statusRowExpectation{
		contains:    []string{"stale (gated tree", "re-run the gate"},
		notContains: []string{"capture-only drift"},
	}
	cases := []staleGateStatusCase{
		{name: "added ROADMAP.md is capture-only", mutate: writeRuntimeFile("ROADMAP.md", "- 2026-07-05  parked idea\n"), want: benign},
		{name: "modified ROADMAP.md is capture-only", seed: writeRuntimeFile("ROADMAP.md", "- 2026-07-04  old idea\n"), mutate: writeRuntimeFile("ROADMAP.md", "- 2026-07-05  parked idea\n"), want: benign},
		{name: "deleted ROADMAP.md is capture-only", seed: writeRuntimeFile("ROADMAP.md", "- 2026-07-04  old idea\n"), mutate: removeRuntimePath("ROADMAP.md"), want: benign},
		{name: "added .bench-notes.md is capture-only", mutate: writeRuntimeFile(".bench-notes.md", "scratch\n"), want: benign},
		{name: "modified .bench-notes.md is capture-only", seed: writeRuntimeFile(".bench-notes.md", "old\n"), mutate: writeRuntimeFile(".bench-notes.md", "new\n"), want: benign},
		{name: "deleted .bench-notes.md is capture-only", seed: writeRuntimeFile(".bench-notes.md", "old\n"), mutate: removeRuntimePath(".bench-notes.md"), want: benign},
		{name: "docs ROADMAP lookalike is strong stale", mutate: writeRuntimeFile("docs/ROADMAP.md", "doc drift\n"), want: strong},
		{name: "ROADMAP backup lookalike is strong stale", mutate: writeRuntimeFile("ROADMAP.md.bak", "doc drift\n"), want: strong},
		{name: ".bench notes lookalike is strong stale", mutate: writeRuntimeFile(".bench/notes.md", "doc drift\n"), want: strong},
		{name: "nested ROADMAP lookalike is strong stale", mutate: writeRuntimeFile("notes/ROADMAP.md", "doc drift\n"), want: strong},
		{name: "mixed capture-only and real drift is strong stale", mutate: func(_ testing.TB, f contract.Fixture) {
			f.WriteFile("ROADMAP.md", "- 2026-07-05  parked idea\n")
			f.WriteFile("docs/x.md", "doc drift\n")
		}, want: strong},
		{name: "cache missing tree is strong stale", cache: literalGateCache("green\n"), want: strong},
		{name: "cache tree none is strong stale", cache: literalGateCache("green none 2026-06-30T00:00:00Z\n"), want: strong},
		{name: "cache missing timestamp is strong stale", cache: func(t testing.TB, f contract.Fixture) string {
			t.Helper()
			return fmt.Sprintf("green %s\n", strings.TrimSpace(f.Bench("tree-hash").Stdout))
		}, want: strong},
		{name: "untrusted cache status is strong stale", cache: func(t testing.TB, f contract.Fixture) string {
			t.Helper()
			return gateCacheLine(t, f, "yellow")
		}, want: strong},
		{name: "deep cwd ROADMAP drift is capture-only", seed: writeRuntimeFile("sub/.keep", ""), mutate: writeRuntimeFile("ROADMAP.md", "- 2026-07-05  parked idea\n"), runDir: "sub", want: benign},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertStaleGateStatusCase(t, tc)
		})
	}
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
	contract.RequireContains(t, out, "resume or clean up")
	contract.RequireNotContains(t, out, "craft-grill → /bench-write-spec")
	if rows := countStatusRows(out); rows > 5 {
		t.Fatalf("budget exceeded five rows (%d):\n%s", rows, out)
	}
}

func testRuntimeStatusWarmPool(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	benchHome := filepath.Join(f.Root, ".bh")
	repoRoot := strings.TrimSpace(f.Git("rev-parse", "--show-toplevel").Stdout)
	pool := filepath.Join(benchHome, "worktrees", filepath.Base(repoRoot)+"-"+cksum(t, repoRoot))
	warm := filepath.Join(pool, "warm")
	contract.Mkdir(t, pool)
	f.Git("worktree", "add", "-q", "--detach", warm, "HEAD")
	out := f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireNotContains(t, out, "resume or clean up")
	lease := strings.TrimSpace(contract.RunAt(t, f, warm, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	contract.WriteFileAbs(t, lease, "")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireContains(t, out, "1 active worktree")
	contract.Remove(t, lease)
	f.Git("branch", "worktree-agent-orphan")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireContains(t, out, "orphaned worktree branch")
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

func testRuntimeStatusLearningsFloor(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — a  [open]\n")
	f.CommitAll("s")
	contract.RequireNotContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "2"}, "status").Stdout, "/bench-what-next")
	contract.RequireContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "1"}, "status").Stdout, "/bench-what-next")
	f.WriteFile(".bench/learnings.md", "## <date> — <short title>  [open]\n")
	contract.RequireNotContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "1"}, "status").Stdout, "/bench-what-next")
	// The floor gates only the learnings component; parked ideas always count.
	f.Bench("idea", "parked past the floor").RequireExit(0)
	contract.RequireContains(t, f.BenchEnv(map[string]string{"BENCH_LEARNINGS_FLOOR": "2"}, "status").Stdout, "/bench-what-next")
}

type staleGateStatusCase struct {
	name   string
	seed   runtimeFixtureEdit
	mutate runtimeFixtureEdit
	cache  func(testing.TB, contract.Fixture) string
	runDir string
	want   statusRowExpectation
}

type runtimeFixtureEdit func(testing.TB, contract.Fixture)

type statusRowExpectation struct {
	contains    []string
	notContains []string
}

func assertStaleGateStatusCase(t testing.TB, tc staleGateStatusCase) {
	t.Helper()
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	if tc.seed != nil {
		tc.seed(t, f)
		f.CommitAll("seed")
	}
	cache := gateCacheLine(t, f, "green")
	if tc.cache != nil {
		cache = tc.cache(t, f)
	}
	writeGateCache(t, f, cache)
	if tc.mutate != nil {
		tc.mutate(t, f)
	}
	out := f.Bench("status").Stdout
	if tc.runDir != "" {
		out = contract.RunAt(t, f, filepath.Join(f.Root, filepath.FromSlash(tc.runDir)), nil, "bash", benchPath(t), "status").Stdout
	}
	requireStatusRow(t, out, tc.want)
}

func requireStatusRow(t testing.TB, out string, want statusRowExpectation) {
	t.Helper()
	for _, needle := range want.contains {
		contract.RequireContains(t, out, needle)
	}
	for _, needle := range want.notContains {
		contract.RequireNotContains(t, out, needle)
	}
}

func writeRuntimeFile(path, contents string) runtimeFixtureEdit {
	return func(_ testing.TB, f contract.Fixture) {
		f.WriteFile(path, contents)
	}
}

func removeRuntimePath(path string) runtimeFixtureEdit {
	return func(t testing.TB, f contract.Fixture) {
		t.Helper()
		contract.Remove(t, filepath.Join(f.Root, filepath.FromSlash(path)))
	}
}

func literalGateCache(line string) func(testing.TB, contract.Fixture) string {
	return func(t testing.TB, _ contract.Fixture) string {
		t.Helper()
		return line
	}
}

func writeGateCache(t testing.TB, f contract.Fixture, line string) {
	t.Helper()
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "bench-last-gate"), line)
}

func gateCacheLine(t testing.TB, f contract.Fixture, status string) string {
	t.Helper()
	tree := strings.TrimSpace(f.Bench("tree-hash").Stdout)
	return fmt.Sprintf("%s %s 2026-06-30T00:00:00Z\n", status, tree)
}
