package runtime

import (
	"bytes"
	"fmt"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	contract.RunParallel(t, "bench status branch-prefix neutrality contract", testRuntimeStatusBranchPrefixNeutrality)
	contract.RunParallel(t, "bench status retirement-signal contract", testRuntimeStatusRetirementSignal)
	contract.RunParallel(t, "bench status orphaned-pickup contract", testRuntimeStatusOrphanedPickup)
	contract.RunParallel(t, "bench status roadmap-reconcile contract", testRuntimeStatusRoadmapReconcile)
	contract.RunParallel(t, "bench status learnings-floor contract", testRuntimeStatusLearningsFloor)
	contract.RunParallel(t, "bench status guards-signal contract", testRuntimeStatusGuardsSignal)
	contract.RunParallel(t, "bench status guards primary-checkout contract", testRuntimeStatusGuardsPrimaryOnly)
	contract.RunParallel(t, "bench status landed-state contract", testRuntimeStatusLandedState)
	contract.RunParallel(t, "bench status intent common-dir contract", testRuntimeStatusIntentCommonDir)
}

func testRuntimeStatusIntentCommonDir(t *testing.T) {
	contract.NoteContractFailure(t, "intent common-dir contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	f.Git("branch", "worktree-agent-live")
	linked := filepath.Join(t.TempDir(), "linked intent reader")
	f.Git("worktree", "add", "-q", "--detach", linked, "HEAD")
	ledger := `{"schema":1,"entries":[{"key":"shared","kind":"claude-agent","objective":"shared objective","created_at":"2026-07-11T00:00:00Z"}]}` + "\n"
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "bench-intent.json"), ledger)
	out := contract.RunAt(t, f, linked, nil, "bash", benchPath(t), "status", "--all")
	out.RequireExit(0)
	out.RequireContains(out.Stdout, "shared objective")
}

func testRuntimeStatusLandedState(t *testing.T) {
	contract.NoteContractFailure(t, "status landed-state contract failed")
	f := contract.NewFixture(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	f.Git("branch", "worktree-agent-live")
	f.WriteFile("tracked.txt", "dirty\n")
	ledger := `{"schema":1,"entries":[` +
		`{"key":"a","kind":"claude-agent","objective":"old` + "\\n\\u001b" + `","created_at":"2026-07-11T00:00:00Z"},` +
		`{"key":"b","kind":"claude-agent","objective":"two","created_at":"2026-07-11T00:00:01Z"},` +
		`{"key":"c","kind":"claude-agent","objective":"three","created_at":"2026-07-11T00:00:02Z"},` +
		`{"key":"d","kind":"claude-agent","objective":"four","created_at":"2026-07-11T00:00:03Z"},` +
		`{"key":"e","kind":"claude-agent","objective":"five","created_at":"2026-07-11T00:00:04Z"},` +
		`{"key":"f","kind":"claude-agent","objective":"six","created_at":"2026-07-11T00:00:05Z"}]}` + "\n"
	path := filepath.Join(gitDir(t, f), "bench-intent.json")
	contract.WriteFileAbs(t, path, ledger)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	compact := f.Bench("status").Stdout
	if strings.Count(compact, "  intent") != 1 || !strings.Contains(compact, "6 uncorrelated") || strings.ContainsRune(compact, '\x1b') {
		t.Fatalf("compact intent rendering unsafe or unbounded:\n%s", compact)
	}
	if countStatusRows(compact) > 5 {
		t.Fatalf("default status exceeded five rows:\n%s", compact)
	}
	all := f.Bench("status", "--all").Stdout
	if strings.Count(all, "  intent") != 6 || strings.ContainsRune(all, '\x1b') {
		t.Fatalf("expanded intent rendering =\n%s", all)
	}
	if again := f.Bench("status").Stdout; again != compact {
		t.Fatalf("status bytes changed\nfirst=%s\nsecond=%s", compact, again)
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("status mutated ledger: %v", err)
	}
}

// testRuntimeStatusGuardsSignal pins story 7: in a routed repo whose pre-push is missing, a
// low-noise guards row fires with the bench link remedy, ranked worktree < guards < drain so
// it never crowds the gate/git rows.
func testRuntimeStatusGuardsSignal(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/lines.env", "BENCH_TIER_TOP=t\nBENCH_TIER_MID=m\nBENCH_TIER_CHEAP=c\n")
	f.WriteFile("IDEAS.md", "- 2026-07-05  parked idea\n")
	f.CommitAll("routed base") // commit so the git signal is quiet and only the ladder rows remain
	// An out-of-pool worktree adds the worktree signal (sev just above guards).
	f.Git("worktree", "add", "-q", "--detach", filepath.Join(f.Root, "outside pool"), "HEAD")

	out := f.Bench("status").Stdout
	requireStatusLineContains(t, out, "guards", "bench link")
	contract.RequireContains(t, out, "pre-push")

	worktree := strings.Index(out, "out-of-pool")
	guards := strings.Index(out, "guards")
	drain := strings.Index(out, "→ /bench-what-next")
	if worktree < 0 || guards < 0 || drain < 0 {
		t.Fatalf("ladder fixture missing a row (worktree=%d guards=%d drain=%d):\n%s", worktree, guards, drain, out)
	}
	if !(worktree < guards && guards < drain) {
		t.Fatalf("severity ladder broken (worktree=%d guards=%d drain=%d):\n%s", worktree, guards, drain, out)
	}

	// A managed pre-push clears the gap: no guards row.
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "hooks", "pre-push"), "#!/usr/bin/env bash\n# bench:managed-pre-push\nexit 0\n")
	managed := f.Bench("status").Stdout
	contract.RequireNotContains(t, managed, "guards")

	// An unrouted repo (no .bench/lines.env) never fires the signal, even with no pre-push.
	unrouted := contract.NewFixture(t)
	contract.RequireNotContains(t, unrouted.Bench("status").Stdout, "guards")
}

// testRuntimeStatusGuardsPrimaryOnly pins story 8: the signal fires only on the primary
// checkout. A pool/linked worktree shares the main .git, so running status from it must not
// double-report the same hook.
func testRuntimeStatusGuardsPrimaryOnly(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/lines.env", "BENCH_TIER_TOP=t\nBENCH_TIER_MID=m\nBENCH_TIER_CHEAP=c\n")
	f.CommitAll("routed base")
	linked := filepath.Join(f.Root, "linked wt")
	f.Git("worktree", "add", "-q", "--detach", linked, "HEAD")

	primary := f.Bench("status").Stdout
	if n := strings.Count(primary, "guards"); n != 1 {
		t.Fatalf("want exactly one guards row from the primary checkout, got %d:\n%s", n, primary)
	}
	fromLinked := contract.RunAt(t, f, linked, nil, "bash", benchPath(t), "status").Stdout
	contract.RequireNotContains(t, fromLinked, "guards")
}

func testRuntimeIdeaRoadmap(t *testing.T) {
	f := contract.NewFixture(t)

	absent := f.Bench("roadmap")
	absent.RequireExit(0)
	contract.RequireContains(t, absent.Stdout, "no ROADMAP.md")
	contract.RequireContains(t, absent.Stdout, "/bench-what-next")

	// Outside a git repo, `roadmap` takes the same structured not-in-repo posture
	// as its sibling `idea` — exit 1, not the in-repo maintenance prompt above.
	noRepo := contract.NewFixture(t, contract.WithNoRepo())
	outsideRepo := noRepo.Bench("roadmap")
	outsideRepo.RequireExit(1)
	contract.RequireContains(t, outsideRepo.Stdout, "error: not in a git repository")
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
	writeGateCache(t, f, fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":"%s","oracle":"%s","recorded_at":%q}`+"\n", strings.Repeat("d", 40), strings.Repeat("0", 64), time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)))
	out := f.Bench("status").Stdout
	contract.RequireContains(t, out, "re-run the gate")
	contract.RequireContains(t, out, "stale (gated tree")
	contract.RequireNotContains(t, out, "capture-only drift")
	contract.RequireNotContains(t, out, "clean — nothing pending")
}

func testRuntimeStatusFreshGreen(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.CommitAll("init")
	f.Bench("gate").RequireExit(0)
	f.Bench("status").RequireContains(f.Bench("status").Stdout, "clean — nothing pending")
}

func testRuntimeStatusDecisions(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/x.md", "### Answer\n— (deferred)\n")
	f.CommitAll("s")
	out := f.Bench("status").Stdout
	contract.RequireContains(t, out, "/bench-shape-idea")
	contract.RequireNotContains(t, out, "craft-grill")
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
	writeGateCache(t, f, fmt.Sprintf(`{"schema":1,"state":"ready","status":"red","tree":%q,"oracle":"%s","recorded_at":%q}`+"\n", tree, strings.Repeat("0", 64), time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)))

	out := f.Bench("status").Stdout
	first := strings.SplitN(out, "\n", 2)[0]
	contract.RequireContains(t, first, "fix before commit")
	contract.RequireContains(t, out, "+1 more (bench status --all)")
	contract.RequireContains(t, out, "/bench-what-next")
	contract.RequireContains(t, out, "split (craft-seams)")
	contract.RequireContains(t, out, "commit on green")
	contract.RequireContains(t, out, "bench worktree clean <path>")
	// The 6th signal (decisions, sev 6) is truncated off the default budget board.
	contract.RequireNotContains(t, out, "/bench-shape-idea")
	if rows := countStatusRows(out); rows > 5 {
		t.Fatalf("budget exceeded five rows (%d):\n%s", rows, out)
	}

	// --all lifts the budget: the 6th signal's action prints and no overflow line remains.
	all := f.Bench("status", "--all").Stdout
	contract.RequireContains(t, all, "/bench-shape-idea")
	contract.RequireNotContains(t, all, "+1 more")
	contract.RequireNotContains(t, all, "(bench status --all)")
	if rows := countStatusRows(all); rows != 6 {
		t.Fatalf("--all should print all six rows, got %d:\n%s", rows, all)
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
	contract.RequireContains(t, out, "bench worktree clean <path>")
	contract.RequireContains(t, out, "1 leased pool worktree")
	requireStatusLineNotContains(t, out, "1 leased pool worktree", "bench worktree clean <path>")
	leasedOut := contract.RunAt(t, f, pool.Leased, map[string]string{"BENCH_HOME": benchHome}, "bash", benchPath(t), "status").Stdout
	contract.RequireContains(t, leasedOut, "1 out-of-pool worktree")
	contract.RequireContains(t, leasedOut, "1 leased pool worktree")
	requireStatusLineNotContains(t, leasedOut, "1 leased pool worktree", "bench worktree clean <path>")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireContains(t, out, "1 out-of-pool worktree")
	contract.RequireContains(t, out, "1 leased pool worktree")
	contract.RequireContains(t, out, "resume leased worktree")
	requireStatusLineNotContains(t, out, "1 leased pool worktree", "bench worktree clean <path>")
	contract.Remove(t, pool.LeaseFile)
	f.Git("branch", "worktree-agent-orphan")
	out = f.BenchEnv(map[string]string{"BENCH_HOME": benchHome}, "status").Stdout
	contract.RequireContains(t, out, "1 out-of-pool worktree")
	contract.RequireContains(t, out, "bench worktree clean <path>")
	contract.RequireNotContains(t, out, "orphaned worktree branch")
}

// A scratch-looking branch name is not ownership evidence and cannot create a status
// cleanup decision. Exact registered worktrees remain the only lifecycle signal.
func testRuntimeStatusBranchPrefixNeutrality(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	commitAllowEmpty(t, f, "init")
	f.Git("branch", "worktree-agent-kept")
	f.Git("branch", "worktree-agent-landed")

	out := f.Bench("status").Stdout
	contract.RequireNotContains(t, out, "orphaned worktree branch")
	contract.RequireNotContains(t, out, "un-landed salvage branch")
	contract.RequireNotContains(t, out, "bench worktree clean")
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
	contract.RequireContains(t, out, "bench spec retire <slug>")
	contract.RequireNotContains(t, out, "promote-then-delete")
	f.WriteFile("scratch.txt", "scratch\n")
	out = f.Bench("status").Stdout
	contract.RequireContains(t, out, "awaiting retirement")
	if strings.Contains(strings.SplitN(out, "\n", 2)[0], "bench spec retire") {
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

func testRuntimeStatusRoadmapReconcile(t *testing.T) {
	// The roadmap-reconcile signal is branch-gated to the default branch, mirroring the
	// retirement signal, so every fixture pins HEAD and origin/HEAD to main.
	onMain := func(f contract.Fixture) {
		f.Git("symbolic-ref", "HEAD", "refs/heads/main")
		f.Git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	}

	// Story 1: a row naming a merged-implemented spec — wrapped in the bold/backtick
	// decoration the real roadmap uses — fires the shipped-row signal with its action.
	merged := contract.NewFixture(t)
	onMain(merged)
	merged.WriteFile("specs/shipped.md", "# Shipped\n\nStatus: implemented\n")
	merged.WriteFile("ROADMAP.md", "# Roadmap\n\n**FT1 — shipped feature.** Staged: `specs/shipped.md` — done.\n")
	merged.CommitAll("merged row")
	out := merged.Bench("status").Stdout
	requireStatusLineContains(t, out, "1 row for merged work", "roadmap", "→ /bench-what-next")
	contract.RequireNotContains(t, out, "retired spec")

	// Story 2: a row naming a spec file that no longer exists fires the dangling detail.
	dangling := contract.NewFixture(t)
	onMain(dangling)
	dangling.WriteFile("ROADMAP.md", "# Roadmap\n\n**FT2 — gone.** Staged: `specs/gone.md` — retired but row survived.\n")
	dangling.CommitAll("dangling row")
	out = dangling.Bench("status").Stdout
	requireStatusLineContains(t, out, "1 row names a retired spec", "roadmap", "→ /bench-what-next")
	contract.RequireNotContains(t, out, "for merged work")

	// Guard: a row naming a still-staged spec is normal open work — nothing fires. A naive
	// "names any spec" implementation would flag it, so this pins the classification.
	staged := contract.NewFixture(t)
	onMain(staged)
	staged.WriteFile("specs/open.md", "# Open\n\nStatus: staged\n")
	staged.WriteFile("ROADMAP.md", "# Roadmap\n\n**FT3 — open.** Staged: `specs/open.md` — in flight.\n")
	staged.CommitAll("staged row")
	out = staged.Bench("status").Stdout
	contract.RequireNotContains(t, out, "for merged work")
	contract.RequireNotContains(t, out, "retired spec")
	if strings.Contains(out, "roadmap") {
		t.Fatalf("staged-spec row wrongly fired the roadmap signal:\n%s", out)
	}

	// Ladder: red gate + dirty tree + shipped row — gate (0) < git (1) < roadmap (9). The
	// new severity must never displace the gate or git rows.
	ladder := contract.NewFixture(t)
	onMain(ladder)
	ladder.WriteFile("specs/shipped.md", "# Shipped\n\nStatus: implemented\n")
	ladder.WriteFile("ROADMAP.md", "# Roadmap\n\n**FT1 — shipped.** `specs/shipped.md`\n")
	ladder.CommitAll("ladder base")
	ladder.WriteFile("dirty.txt", "dirty\n")
	tree := strings.TrimSpace(ladder.Bench("tree-hash").Stdout)
	writeGateCache(t, ladder, fmt.Sprintf(`{"schema":1,"state":"ready","status":"red","tree":%q,"oracle":"%s","recorded_at":%q}`+"\n", tree, strings.Repeat("0", 64), time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)))
	out = ladder.Bench("status").Stdout
	gate := strings.Index(out, "fix before commit")
	gitRow := strings.Index(out, "commit on green")
	road := strings.Index(out, "row for merged work")
	if gate < 0 || gitRow < 0 || road < 0 {
		t.Fatalf("ladder fixture missing a row (gate=%d git=%d roadmap=%d):\n%s", gate, gitRow, road, out)
	}
	if !(gate < gitRow && gitRow < road) {
		t.Fatalf("severity ladder broken (gate=%d git=%d roadmap=%d):\n%s", gate, gitRow, road, out)
	}

	// Boundary: two shipped rows collapse into one signal row carrying the count of 2.
	two := contract.NewFixture(t)
	onMain(two)
	two.WriteFile("specs/one.md", "# One\n\nStatus: implemented\n")
	two.WriteFile("specs/two.md", "# Two\n\nStatus: implemented\n")
	two.WriteFile("ROADMAP.md", "# Roadmap\n\n**A** `specs/one.md`\n\n**B** `specs/two.md`\n")
	two.CommitAll("two rows")
	out = two.Bench("status").Stdout
	contract.RequireContains(t, out, "2 rows for merged work")
	if n := strings.Count(out, "roadmap"); n != 1 {
		t.Fatalf("want exactly one roadmap signal row, got %d:\n%s", n, out)
	}
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

// requireStatusLineContains asserts the single status row carrying needle also carries every
// want — the positive twin of requireStatusLineNotContains, so a row's signal, detail, and
// action are checked as one line rather than merely co-present in the board.
func requireStatusLineContains(t testing.TB, out, needle string, wants ...string) {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, needle) {
			for _, w := range wants {
				contract.RequireContains(t, line, w)
			}
			return
		}
	}
	t.Fatalf("status line containing %q not found in:\n%s", needle, out)
}
