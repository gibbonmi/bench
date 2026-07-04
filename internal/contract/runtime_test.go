package contract

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestRuntimeContracts(t *testing.T) {
	skipIfSubjectBenchMissing(t)
	t.Run("bench idea/roadmap contract", testRuntimeIdeaRoadmap)
	t.Run("bench gate repo-root cwd contract", testRuntimeGateRepoRootCWD)
	t.Run("bench gate BENCH_GATE cwd contract", testRuntimeGateBenchGateCWD)
	t.Run("bench gate resolution-order contract", testRuntimeGateResolutionOrder)
	t.Run("bench status clean contract", testRuntimeStatusClean)
	t.Run("bench status footer contract", testRuntimeStatusFooter)
	t.Run("bench status stale-gate contract", testRuntimeStatusStaleGate)
	t.Run("bench status fresh-green contract", testRuntimeStatusFreshGreen)
	t.Run("bench status decisions contract", testRuntimeStatusDecisions)
	t.Run("bench status unresolved-maps count contract", testRuntimeStatusUnresolvedMapsCount)
	t.Run("bench status budget contract", testRuntimeStatusBudget)
	t.Run("bench status warm-pool contract", testRuntimeStatusWarmPool)
	t.Run("bench status gate-cache write contract", testRuntimeStopHookGateCacheWrite)
	t.Run("stop hook no-gate no-cache contract", testRuntimeStopHookNoGateNoCache)
	t.Run("stop hook missing-core-binary fail-safe contract", testRuntimeStopHookMissingCoreBinary)
	t.Run("bench gate missing-core-binary fail-safe contract", testRuntimeGateMissingCoreBinary)
	t.Run("bench status retirement-signal contract", testRuntimeStatusRetirementSignal)
	t.Run("bench gate verdict-record contract", testRuntimeGateVerdictRecord)
	t.Run("bench status learnings-floor contract", testRuntimeStatusLearningsFloor)
	t.Run("bench structure shell-file contract", testRuntimeStructureShellFile)
	t.Run("bench structure budgets contract", testRuntimeStructureBudgets)
	t.Run("bench structure path-with-spaces contract", testRuntimeStructurePathWithSpaces)
	t.Run("bench worktree lease/reuse contract", testRuntimeWorktreeLeaseReuse)
	t.Run("bench worktree lease hardening contract", testRuntimeWorktreeLeaseHardening)
	t.Run("bench worktree concurrent-acquire contract", testRuntimeWorktreeConcurrentAcquire)
	t.Run("bench symlinked kit-dir portability contract", testRuntimeSymlinkedKitDir)
	t.Run("stop hook stop_hook_active contract", testRuntimeStopHookActive)
	t.Run("stop hook missing-bench fail-open contract", testRuntimeStopHookMissingBenchFailOpen)
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

func testRuntimeGateRepoRootCWD(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("package.json", "{\"ok\":true}\n")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\n[ -f package.json ]\n")
	mkdir(t, filepath.Join(f.Root, "sub"))

	runAt(t, f, filepath.Join(f.Root, "sub"), nil, "bash", benchPath(t), "gate").RequireExit(0)
}

func testRuntimeGateBenchGateCWD(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("package.json", "{\"ok\":true}\n")
	f.WriteExecutable("gate-root.sh", "#!/usr/bin/env bash\n[ -f package.json ]\n")
	mkdir(t, filepath.Join(f.Root, "sub"))

	runAt(t, f, filepath.Join(f.Root, "sub"), map[string]string{"BENCH_GATE": "./gate-root.sh"}, "bash", benchPath(t), "gate").RequireExit(0)
}

func testRuntimeGateResolutionOrder(t *testing.T) {
	f := NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.BenchEnv(map[string]string{"BENCH_GATE": "exit 1"}, "gate").RequireExit(0)

	remove(t, filepath.Join(f.Root, ".bench", "gate.sh"))
	f.WriteFile("package.json", "{\"private\":true}\n")
	f.BenchEnv(map[string]string{"BENCH_GATE": "exit 0"}, "gate").RequireExit(0)

	auto := f.Bench("gate")
	auto.RequireNotContains(auto.Stdout+auto.Stderr, "no gate found")
	if auto.ExitCode == 3 {
		t.Fatalf("package.json resolved to no-gate exit 3\nstdout:\n%s\nstderr:\n%s", auto.Stdout, auto.Stderr)
	}

	remove(t, filepath.Join(f.Root, "package.json"))
	commitAllowEmpty(t, f, "init")
	remove(t, filepath.Join(gitDir(t, f), "bench-last-gate"))
	noGate := f.Bench("gate")
	noGate.RequireExit(3)
	noGate.RequireContains(noGate.Stdout+noGate.Stderr, "no gate found")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("no-gate case recorded a verdict")
	}
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

func testRuntimeStopHookGateCacheWrite(t *testing.T) {
	f := copiedCLIHookFixture(t, true)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")

	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n")
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	data := readFileAbs(t, cache)
	if !regexp.MustCompile(`^(green|red) [0-9a-f]+ [0-9T:Z-]+$`).MatchString(strings.TrimSpace(data)) {
		t.Fatalf("gate cache not <status> <tree> <iso8601>: %q", data)
	}
}

func testRuntimeStopHookNoGateNoCache(t *testing.T) {
	f := copiedCLIHookFixture(t, true)
	f.CommitAll("init")

	probe := runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n")
	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout+probe.Stderr, "no gate found")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("armed no-gate stop recorded a gate cache")
	}
}

func testRuntimeStopHookMissingCoreBinary(t *testing.T) {
	f := copiedCLIHookFixture(t, false)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")

	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("missing core binary forged a gate verdict")
	}
}

func testRuntimeGateMissingCoreBinary(t *testing.T) {
	f := copiedCLIHookFixture(t, false)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")

	f.Run("bash", filepath.Join(f.Root, "bin", "bench.sh"), "gate")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("gate_record forged a verdict with no core binary")
	}
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

func testRuntimeGateVerdictRecord(t *testing.T) {
	f := NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("init")
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	writeFileAbs(t, cache, "green deadbeefdeadbeefdeadbeefdeadbeefdeadbeef 2026-06-30T00:00:00Z\n")
	f.Bench("gate").RequireExit(0)
	requireContains(t, readFileAbs(t, cache), "green "+strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout))
	requireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	commitAllowEmpty(t, f, "same-tree")
	requireNotContains(t, f.Bench("status").Stdout, "re-run the gate")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
	if p := f.Bench("gate"); p.ExitCode == 0 {
		t.Fatal("red gate run exited zero")
	}
	requireContains(t, readFileAbs(t, cache), "red "+strings.TrimSpace(f.Bench("tree-hash").Stdout))
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

func testRuntimeStructureShellFile(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("big.sh", repeatLines(401, "x=\n"))
	f.Git("add", "big.sh")
	probe := f.Bench("structure")
	if probe.ExitCode == 0 {
		t.Fatalf("shell source over the line budget did not fail structure\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	probe.RequireContains(probe.Stdout+probe.Stderr, "FILE TOO LONG")
}

func testRuntimeStructureBudgets(t *testing.T) {
	f := NewFixture(t)
	f.WriteFile("big.sh", repeatLines(401, "x=\n"))
	f.WriteFile("mid.sh", repeatLines(200, "y=\n"))
	for i := 1; i <= 13; i++ {
		f.WriteFile(fmt.Sprintf("sub/f%d.sh", i), "z=1\n")
	}
	f.WriteFile(".bench/structure.budgets", "# reviewer grants\nbig.sh 500\nsub/ 20\nweird abc\nmid.sh 100")
	f.CommitAll("s")
	probe := f.Bench("structure")
	if probe.ExitCode == 0 {
		t.Fatal("override below the global cap did not fail structure")
	}
	out := probe.Stdout + probe.Stderr
	requireContains(t, out, "ignoring malformed line")
	requireNotContains(t, out, "big.sh")
	requireNotContains(t, out, "DIR CROWDED")
	requireContains(t, out, "200 lines (max 100)   mid.sh")
}

func testRuntimeStructurePathWithSpaces(t *testing.T) {
	f := NewFixture(t)
	for i := 1; i <= 13; i++ {
		f.WriteFile(fmt.Sprintf("space dir/file%d.sh", i), fmt.Sprintf("x=%d\n", i))
	}
	f.Git("add", "space dir")
	probe := f.BenchEnv(map[string]string{"BENCH_MAX_DIR_FILES": "12"}, "structure")
	if probe.ExitCode == 0 {
		t.Fatalf("crowded path-with-spaces directory did not fail structure\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
	out := probe.Stdout + probe.Stderr
	requireContains(t, out, "space dir/")
	requireNotContains(t, out, "   ./")
	requireNotContains(t, out, "   dir/")
}

func testRuntimeWorktreeLeaseReuse(t *testing.T) {
	f := NewFixture(t)
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
	paths := nonEmptyLines(readFileAbs(t, record))
	requireIntEqual(t, len(paths), 2, "worktree shell did not run twice")
	requireEqual(t, paths[0], paths[1], "worktree pool did not reuse a clean released path")
	if _, err := os.Stat(strings.TrimSpace(runAt(t, f, paths[1], nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)); err == nil {
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
	f := NewFixture(t)
	commitAllowEmpty(t, f, "init")
	f.WriteExecutable("rec-shell", "#!/usr/bin/env bash\n: \"${BENCH_WT_RECORD:?}\"\npwd >> \"$BENCH_WT_RECORD\"\n")
	record := filepath.Join(f.Root, "paths")
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bh"), "BENCH_WT_RECORD": record, "SHELL": filepath.Join(f.Root, "rec-shell")}
	runWT := func() { f.BenchEnv(env, "worktree").RequireExit(0) }
	runWT()
	p := strings.TrimSpace(readFileAbs(t, record))
	lease := strings.TrimSpace(runAt(t, f, p, nil, "git", "rev-parse", "--git-path", "bench-lease").Stdout)
	writeFileAbs(t, lease, "4194300 2020-01-01T00:00:00Z\n")
	runWT()
	writeFileAbs(t, lease, fmt.Sprintf("%d 2026-01-01T00:00:00Z\n", os.Getpid()))
	runWT()
	if _, err := os.Stat(lease); err != nil {
		t.Fatal("live foreign lease was removed by a foreign release")
	}
	writeFileAbs(t, lease, "")
	old := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	if err := os.Chtimes(lease, old, old); err != nil {
		t.Fatalf("age empty lease: %v", err)
	}
	runWT()
	writeFileAbs(t, lease, "")
	runWT()
	remove(t, lease)
	paths := nonEmptyLines(readFileAbs(t, record))
	requireIntEqual(t, len(paths), 5, "expected five worktree runs")
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
	f := NewFixture(t)
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
	done := make(chan Probe, 2)
	for i := 0; i < 2; i++ {
		go func() { done <- f.BenchEnv(env, "worktree") }()
	}
	for i := 0; i < 2; i++ {
		(<-done).RequireExit(0)
	}
	paths := nonEmptyLines(readFileAbs(t, record))
	sort.Strings(paths)
	requireIntEqual(t, len(paths), 2, "concurrent worktree runs did not both complete")
	if paths[0] == paths[1] {
		t.Fatal("concurrent acquires shared a worktree")
	}
}

func testRuntimeSymlinkedKitDir(t *testing.T) {
	tmp := t.TempDir()
	repo := filepath.Join(tmp, "repo")
	binDir := filepath.Join(tmp, "bin")
	shim := filepath.Join(tmp, "shim")
	mkdir(t, repo)
	mkdir(t, binDir)
	mkdir(t, shim)
	if err := os.Symlink(benchPath(t), filepath.Join(binDir, "bench")); err != nil {
		t.Fatalf("symlink bench: %v", err)
	}
	writeExecutableAbs(t, filepath.Join(shim, "readlink"), "#!/usr/bin/env bash\nif [ \"${1:-}\" = \"-f\" ]; then exit 1; fi\n/usr/bin/readlink \"$@\"\n")
	f := Fixture{t: t, Root: repo, Env: isolatedEnv(t, repo)}
	f.Git("init", "-q")
	runAt(t, f, repo, map[string]string{"PATH": shim + ":/usr/bin:/bin"}, filepath.Join(binDir, "bench"), "link").RequireExit(0)
	if _, err := os.Stat(filepath.Join(repo, ".bench", "BENCH.md")); err != nil {
		t.Fatal("symlinked bench did not resolve kit dir without readlink -f")
	}
}

func testRuntimeStopHookActive(t *testing.T) {
	f := copiedCLIHookFixture(t, true)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 1\n")
	f.CommitAll("init")
	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{\"stop_hook_active\":true}\n").RequireExit(0)
	runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1"}, "{}\n").RequireExit(2)
}

func testRuntimeStopHookMissingBenchFailOpen(t *testing.T) {
	f := NewFixture(t)
	commitAllowEmpty(t, f, "init")
	probe := runStopHook(t, f, map[string]string{"BENCH_SHIFT": "1", "PATH": "/usr/bin:/bin"}, "{}\n")
	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout+probe.Stderr, "bench")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-last-gate")); err == nil {
		t.Fatal("missing bench forged a gate cache")
	}
}

func copiedCLIHookFixture(t *testing.T, withCore bool) Fixture {
	t.Helper()
	f := NewFixture(t)
	mkdir(t, filepath.Join(f.Root, ".bench"))
	mkdir(t, filepath.Join(f.Root, "bin"))
	if withCore {
		mkdir(t, filepath.Join(f.Root, "dist"))
		copyRuntimeFile(t, filepath.Join(SubjectRoot(t), "dist", "bench"), filepath.Join(f.Root, "dist", "bench"), 0o755)
	}
	matches, err := filepath.Glob(filepath.Join(SubjectRoot(t), "bin", "*.sh"))
	if err != nil {
		t.Fatalf("glob bin scripts: %v", err)
	}
	for _, src := range matches {
		copyRuntimeFile(t, src, filepath.Join(f.Root, "bin", filepath.Base(src)), 0o755)
	}
	return f
}

func runStopHook(t *testing.T, f Fixture, env map[string]string, stdin string) Probe {
	t.Helper()
	cmd := exec.Command("bash", filepath.Join(SubjectRoot(t), ".bench", "hooks", "stop.sh"))
	cmd.Dir = f.Root
	cmd.Env = mergeEnv(f.Env, envToSpec(env))
	cmd.Stdin = strings.NewReader(stdin)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return Probe{t: t, ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String()}
}

func runAt(t testing.TB, f Fixture, dir string, env map[string]string, name string, args ...string) Probe {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnv(f.Env, envToSpec(env))
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return Probe{t: t, ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String()}
}

func envToSpec(env map[string]string) Env {
	out := make(Env, len(env))
	for k, v := range env {
		value := v
		out[k] = &value
	}
	return out
}

func benchPath(t testing.TB) string {
	t.Helper()
	return filepath.Join(SubjectRoot(t), "bin", "bench.sh")
}

func gitDir(t testing.TB, f Fixture) string {
	t.Helper()
	return strings.TrimSpace(f.Git("rev-parse", "--absolute-git-dir").Stdout)
}

func commitAllowEmpty(t testing.TB, f Fixture, message string) {
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

func writeExecutableAbs(t testing.TB, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o755); err != nil {
		t.Fatalf("write executable %s: %v", path, err)
	}
}

func writeFileAbs(t testing.TB, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readFileAbs(t testing.TB, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func mkdir(t testing.TB, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func remove(t testing.TB, path string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("remove %s: %v", path, err)
	}
}

func requireFileMatches(t testing.TB, f Fixture, path, pattern, msg string) {
	t.Helper()
	if !regexp.MustCompile(pattern).MatchString(f.ReadFile(path)) {
		t.Fatalf("%s:\n%s", msg, f.ReadFile(path))
	}
}

func requireContains(t testing.TB, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("missing %q in:\n%s", needle, haystack)
	}
}

func requireNotContains(t testing.TB, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Fatalf("unexpected %q in:\n%s", needle, haystack)
	}
}

func requireIntEqual(t testing.TB, got, want int, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %d, want %d", msg, got, want)
	}
}

func lineCount(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n")
}

func repeatLines(n int, line string) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(line)
	}
	return b.String()
}

func nonEmptyLines(s string) []string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}
	return out
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
