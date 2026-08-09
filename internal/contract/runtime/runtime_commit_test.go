package runtime

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

// commitFixture builds a repo on branch main with a committed green gate (whose verdict
// is `exit ${GATE_RC:-0}`, so a red run needs no untracked marker) and a configured Git
// identity for landing commits. The gate tallies in the common Git directory so both the
// invoking checkout and prospective checkouts report through one observable counter.
func commitFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\ngitdir=\"$(git rev-parse --git-common-dir)\"\necho run >> \"$gitdir/gate-runs\"\nexit ${GATE_RC:-0}\n")
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":["GATE_RC"],"paths":[],"tools":[]}`+"\n")
	// The passlist requires every declared value to exist. Green controls declare
	// the value explicitly instead of relying on the script's shell default.
	f.Env["GATE_RC"] = "0"
	f.WriteFile("seed.txt", "seed\n")
	f.CommitAll("seed")
	f.WriteFile(".git/gate-runs", "")
	return f
}

func headSha(f contract.Fixture) string {
	return strings.TrimSpace(f.Git("rev-parse", "HEAD").Stdout)
}

func committedNames(f contract.Fixture) string {
	return f.Git("show", "--name-only", "--pretty=format:", "HEAD").Stdout
}

func TestRuntimeCommitContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "green gate commits named path", testCommitGreenCommits)
	contract.RunParallel(t, "fresh green verdict is reused, gate not re-run", testCommitFreshVerdictReused)
	contract.RunParallel(t, "fresh green is reused while another run holds the gate lock", testCommitReusesFreshVerdictUnderHeldLock)
	contract.RunParallel(t, "stale verdict re-runs the gate", testCommitStaleVerdictRerunsGate)
	contract.RunParallel(t, "red gate refuses commit", testCommitRedRefuses)
	contract.RunParallel(t, "unnamed work is preserved", testCommitPreservesUnnamedWork)
	contract.RunParallel(t, "prospective gate sees exact tree", testCommitGateSeesProspectiveTree)
	contract.RunParallel(t, "CAS loser recomposes in fresh process", testCommitCASLossRecomposes)
	contract.RunParallel(t, "detached HEAD lands without moving branch", testCommitDetachedHEAD)
	contract.RunParallel(t, "glob/space path survives whole", testCommitHostileFilenames)
	contract.RunParallel(t, "named directory commits its changed children", testCommitDirectoryCommitsChildren)
	contract.RunParallel(t, "named directory stops at its own segment", testCommitDirectoryStopsAtSegment)
	contract.RunParallel(t, "named directory preserves an outside file", testCommitDirectoryPreservesOutsideFile)
	contract.RunParallel(t, "deleted named directory commits its removals", testCommitDeletedDirectoryCommitsRemovals)
	contract.RunParallel(t, "deleted path stages the removal", testCommitStagesDeletion)
	contract.RunParallel(t, "reconciliation failure is bounded", testCommitReconciliationFailureIsBounded)
	contract.RunParallel(t, "staged deletion commits", testCommitStagesStagedDeletion)
	contract.RunParallel(t, "staged rename commits whole", testCommitStagesStagedRename)
	contract.RunParallel(t, "unknown path fails before the gate", testCommitUnknownPathFailsFast)
	contract.RunParallel(t, "spec transition is authorized and lands unchanged", testCommitSpecFlip)
	contract.RunParallel(t, "red gate preserves staged spec", testCommitSpecRedPreservesStaged)
	contract.RunParallel(t, "bad --spec fails before the gate", testCommitSpecFailsFast)
	contract.RunParallel(t, "empty named delta ignores foreign staged work", testCommitEmptyRefused)
	contract.RunParallel(t, "help remains a successful stdout response", testCommitHelp)
	contract.RunParallel(t, "usage errors exit 2", testCommitUsageExitTwo)
	contract.RunParallel(t, "-- makes a leading-dash path expressible", testCommitDoubleDashPath)
	contract.RunParallel(t, "empty positional is usage, not a path", testCommitEmptyPositionalIsUsage)
}

// testCommitDoubleDashPath drives the grammar's `--` rule end to end: a file whose name
// begins with a dash is inexpressible without it, and the rule has to survive every layer
// between argv and literal prospective composition.
func testCommitDoubleDashPath(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("-weird.txt", "dash-led\n")
	before := headSha(f)

	f.Bench("commit", "-m", "dash path", "--", "-weird.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	contract.RequireContains(t, committedNames(f), "-weird.txt")
}

func testCommitFreshVerdictReused(t *testing.T) {
	// The verdict cache is keyed to the content hash of the tested tree, so a green
	// `bench gate` on the exact tree being committed already proves this diff. Commit
	// reuses that verdict instead of paying the full gate a second time — the
	// five-minute-docs-commit defect.
	f := commitFixture(t)
	f.WriteFile("work.txt", "changed\n")
	f.Bench("gate").RequireExit(0)
	before := headSha(f)

	res := f.Bench("commit", "-m", "do work", "work.txt")
	res.RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	contract.RequireIntEqual(t, gateRunTally(t, f), 1, "commit re-ran the gate despite a fresh green verdict for the identical tree")
	// The reuse is silent to the operator unless the gate says so: a skipped run that
	// reports nothing reads as a gate that never ran. The line comes from the gate's own
	// emitter, the single home of verdict-reuse policy.
	contract.RequireContains(t, res.Stdout, "gate: green (fresh verdict reused for this tree)")
}

// testCommitReusesFreshVerdictUnderHeldLock pins reuse against contention, the state a
// concurrent gate run leaves behind: the verdict for this exact tree is already green, so a
// commit must be answered from it without queueing for the execution lock. A reuse check
// that sits under the lock instead refuses the commit and — worse — demotes the green it
// could have reused to `pending` on the way out, which is why the record's bytes are graded
// beside the exit code.
func testCommitReusesFreshVerdictUnderHeldLock(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("work.txt", "changed\n")
	f.Bench("gate").RequireExit(0)
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	seeded := contract.ReadFileAbs(t, cache)
	contract.RequireContains(t, seeded, `"state":"ready"`)
	contract.RequireContains(t, seeded, `"status":"green"`)
	before := headSha(f)

	held, err := os.OpenFile(filepath.Join(gitDir(t, f), "bench-gate.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()
	acquireTestGateLock(t, held)
	defer releaseTestGateLock(held)

	res := f.Bench("commit", "-m", "do work", "work.txt")
	res.RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a reusable green verdict")
	}
	contract.RequireContains(t, res.Stdout, "gate: green (fresh verdict reused for this tree)")
	contract.RequireIntEqual(t, gateRunTally(t, f), 1, "commit reached the oracle despite a reusable green verdict")
	if got := contract.ReadFileAbs(t, cache); got != seeded {
		t.Fatalf("contended commit rewrote the verdict record:\nseeded %q\ngot    %q", seeded, got)
	}
}

func testCommitStaleVerdictRerunsGate(t *testing.T) {
	// A verdict recorded for a different tree proves nothing about this diff: any edit
	// after the gate run must send commit back through the full gate.
	f := commitFixture(t)
	f.Bench("gate").RequireExit(0)
	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)

	f.Bench("commit", "-m", "do work", "work.txt", "work.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	contract.RequireIntEqual(t, gateRunTally(t, f), 2, "commit trusted a verdict recorded for a different tree")
}

func testCommitGreenCommits(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)

	f.Bench("commit", "-m", "do work", "work.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	names := committedNames(f)
	contract.RequireContains(t, names, "work.txt")
	// The no-bare-`git add -A` guarantee: the unchanged seed never rides into the commit.
	contract.RequireNotContains(t, names, "seed.txt")
}

func testCommitStagesDeletion(t *testing.T) {
	// spec-retire always deletes a spec, so prospective composition must preserve a named
	// removal even though the file is no longer present to read from the working tree.
	f := commitFixture(t)
	contract.Remove(t, filepath.Join(f.Root, "seed.txt"))
	before := headSha(f)

	f.Bench("commit", "-m", "remove seed", "seed.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	// The deletion is what landed: the commit records seed.txt and the tree no longer tracks it.
	contract.RequireContains(t, committedNames(f), "seed.txt")
	contract.RequireNotContains(t, f.Git("ls-files").Stdout, "seed.txt")
}

func testCommitStagesStagedDeletion(t *testing.T) {
	// A removal already staged (`git rm`, or a delegate's index) leaves the path absent
	// from both worktree and index. Prospective composition must still recover the deletion
	// from the expected base and land it.
	f := commitFixture(t)
	f.Git("rm", "-q", "seed.txt")
	before := headSha(f)

	f.Bench("commit", "-m", "remove seed", "seed.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	contract.RequireContains(t, committedNames(f), "seed.txt")
	contract.RequireNotContains(t, f.Git("ls-files").Stdout, "seed.txt")
}

func testCommitStagesStagedRename(t *testing.T) {
	// The recorded FT38 shape: `git mv` stages the rename, so the old path is absent from
	// worktree and index while the new path is index-only. Naming both halves must commit
	// the whole rename.
	f := commitFixture(t)
	f.Git("mv", "seed.txt", "renamed.txt")
	before := headSha(f)

	f.Bench("commit", "-m", "rename seed", "seed.txt", "renamed.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	// Rename detection collapses the pair to its destination in `--name-only`; list with
	// --no-renames so the assertion sees both halves of what landed.
	names := f.Git("show", "--name-only", "--no-renames", "--pretty=format:", "HEAD").Stdout
	contract.RequireContains(t, names, "seed.txt")
	contract.RequireContains(t, names, "renamed.txt")
	tracked := f.Git("ls-files").Stdout
	contract.RequireContains(t, tracked, "renamed.txt")
	contract.RequireNotContains(t, tracked, "seed.txt")
}

func testCommitUnknownPathFailsFast(t *testing.T) {
	// A named path absent from worktree, index, and the expected base is an attribution
	// error. Pairing it with a red gate proves the refusal happens before authorization.
	f := commitFixture(t)
	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)

	p := f.BenchEnv(map[string]string{"GATE_RC": "1"}, "commit", "-m", "ghost", "ghost.txt", "work.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "not found in worktree, index, or expected base")
	contract.RequireNotContains(t, p.Stderr, "prospective authorization refused")
	if headSha(f) != before {
		t.Fatal("HEAD advanced despite an unknown named path")
	}
}

func testCommitRedRefuses(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("work.txt", "staged\n")
	f.Git("add", "work.txt")
	f.WriteFile("work.txt", "staged plus unstaged\n")
	f.WriteFile("foreign.txt", "foreign\n")
	before := headSha(f)
	beforeStatus := f.Git("status", "--porcelain", "-uall").Stdout
	beforeIndex := f.Git("ls-files", "--stage").Stdout
	beforeWork := f.ReadFile("work.txt")
	beforeForeign := f.ReadFile("foreign.txt")

	p := f.BenchEnv(map[string]string{"GATE_RC": "23"}, "commit", "-m", "do work", "work.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "prospective authorization refused")
	if headSha(f) != before {
		t.Fatal("commit landed on a red gate")
	}
	if got := f.Git("status", "--porcelain", "-uall").Stdout; got != beforeStatus {
		t.Fatalf("red gate changed classifications: got %q, want %q", got, beforeStatus)
	}
	if got := f.Git("ls-files", "--stage").Stdout; got != beforeIndex {
		t.Fatalf("red gate changed index: got %q, want %q", got, beforeIndex)
	}
	if got := f.ReadFile("work.txt"); got != beforeWork {
		t.Fatalf("red gate changed named bytes: got %q, want %q", got, beforeWork)
	}
	if got := f.ReadFile("foreign.txt"); got != beforeForeign {
		t.Fatalf("red gate changed foreign bytes: got %q, want %q", got, beforeForeign)
	}
	contract.RequireNotContains(t, p.Stderr, "v1:")
}

func testCommitPreservesUnnamedWork(t *testing.T) {
	// A successful attributed landing must leave unrelated work available to its owner.
	for _, tc := range []struct {
		name  string
		setup func(f contract.Fixture)
	}{
		{"untracked", func(f contract.Fixture) { f.WriteFile("stray.txt", "x\n") }},
		{"tracked-modified", func(f contract.Fixture) { f.WriteFile("seed.txt", "tampered\n") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := commitFixture(t)
			f.WriteFile("work.txt", "changed\n")
			tc.setup(f)
			before := headSha(f)

			f.Bench("commit", "-m", "do work", "work.txt").RequireExit(0)
			if headSha(f) == before {
				t.Fatal("HEAD did not advance for the named change")
			}
			contract.RequireNotContains(t, committedNames(f), "stray.txt")
			contract.RequireNotContains(t, committedNames(f), "seed.txt")
		})
	}
}

func testCommitGateSeesProspectiveTree(t *testing.T) {
	f := commitFixture(t)
	f.WriteExecutable(".bench/gate.sh", `#!/usr/bin/env bash
set -eu
gitdir="$(git rev-parse --git-common-dir)"
echo run >> "$gitdir/gate-runs"
test "$(cat seed.txt)" = seed
test "$(cat named.txt)" = "named working bytes"
test "$(cat foreign-staged.txt)" = "foreign staged base"
test "$(cat foreign-unstaged.txt)" = "foreign unstaged base"
test ! -e foreign-untracked.txt
test ! -e ignored.txt
printf inspected > "$gitdir/prospective-inspected"
`)
	f.WriteFile(".gitignore", "ignored.txt\n")
	f.WriteFile("named.txt", "named base\n")
	f.WriteFile("foreign-staged.txt", "foreign staged base\n")
	f.WriteFile("foreign-unstaged.txt", "foreign unstaged base\n")
	f.CommitAll("inspection base")

	f.WriteFile("named.txt", "named staged bytes\n")
	f.Git("add", "named.txt")
	f.WriteFile("named.txt", "named working bytes\n")
	f.WriteFile("foreign-staged.txt", "foreign staged bytes\n")
	f.Git("add", "foreign-staged.txt")
	f.WriteFile("foreign-unstaged.txt", "foreign unstaged bytes\n")
	f.WriteFile("foreign-untracked.txt", "foreign untracked bytes\n")
	f.WriteFile("ignored.txt", "foreign ignored bytes\n")
	foreignStatus := f.Git("status", "--porcelain", "--ignored", "-uall", "--", "foreign-staged.txt", "foreign-unstaged.txt", "foreign-untracked.txt", "ignored.txt").Stdout
	foreignIndex := f.Git("show", ":foreign-staged.txt").Stdout

	f.Bench("commit", "-m", "inspect prospective tree", "named.txt").RequireExit(0)
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "prospective-inspected")); err != nil {
		t.Fatalf("public gate did not record prospective inspection: %v", err)
	}
	if got := f.Git("show", "HEAD:named.txt").Stdout; got != "named working bytes\n" {
		t.Fatalf("published named bytes = %q", got)
	}
	if got := f.Git("show", "HEAD:foreign-staged.txt").Stdout; got != "foreign staged base\n" {
		t.Fatalf("published foreign staged bytes = %q", got)
	}
	if got := f.Git("show", ":foreign-staged.txt").Stdout; got != foreignIndex {
		t.Fatalf("foreign index bytes changed: got %q, want %q", got, foreignIndex)
	}
	for path, want := range map[string]string{
		"foreign-staged.txt":    "foreign staged bytes\n",
		"foreign-unstaged.txt":  "foreign unstaged bytes\n",
		"foreign-untracked.txt": "foreign untracked bytes\n",
		"ignored.txt":           "foreign ignored bytes\n",
	} {
		if got := f.ReadFile(path); got != want {
			t.Fatalf("foreign worktree bytes for %s = %q, want %q", path, got, want)
		}
	}
	if got := f.Git("status", "--porcelain", "--ignored", "-uall", "--", "foreign-staged.txt", "foreign-unstaged.txt", "foreign-untracked.txt", "ignored.txt").Stdout; got != foreignStatus {
		t.Fatalf("foreign classifications changed:\ngot  %q\nwant %q", got, foreignStatus)
	}
	if got := f.Git("status", "--porcelain", "--", "named.txt").Stdout; got != "" {
		t.Fatalf("named path is not clean after landing: %q", got)
	}
}

func testCommitCASLossRecomposes(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("specs/concurrent/spec.md", "# concurrent\nStatus: staged\n")
	f.CommitAll("staged concurrent spec")
	f.WriteFile("winner.txt", "winner bytes\n")
	f.WriteFile("loser.txt", "loser bytes\n")
	base := headSha(f)

	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	shim := t.TempDir()
	contract.WriteExecutableAbs(t, filepath.Join(shim, "git"), `#!/usr/bin/env bash
if [ "$#" -eq 4 ] && [ "$1" = -C ] && [ "$2" = "${BENCH_TEST_CAS_ROOT:-}" ] && [ "$3" = read-tree ] && [ "$4" = "${BENCH_TEST_CAS_BASE:-}" ]; then
  printf x >&3
  IFS= read -r -n 1 _ <&4 || exit 97
fi
exec "$BENCH_TEST_REAL_GIT" "$@"
`)
	readyRead, readyWrite, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	releaseRead, releaseWrite, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer readyRead.Close()
	defer releaseWrite.Close()

	var loserStdout, loserStderr bytes.Buffer
	loser := exec.Command(benchPath(t), "commit", "-m", "loser", "--spec", "concurrent", "loser.txt")
	loser.Dir = f.Root
	loser.Env = selectedSurfaceEnv(t, f, map[string]string{
		"PATH": shim + string(os.PathListSeparator) + os.Getenv("PATH"), "GATE_RC": "0",
		"BENCH_TEST_CAS_ROOT": f.Root, "BENCH_TEST_CAS_BASE": base, "BENCH_TEST_REAL_GIT": realGit,
		"GIT_AUTHOR_NAME": "loser", "GIT_AUTHOR_EMAIL": "loser@local",
	})
	loser.ExtraFiles = []*os.File{readyWrite, releaseRead}
	loser.Stdout, loser.Stderr = &loserStdout, &loserStderr
	if err := loser.Start(); err != nil {
		t.Fatal(err)
	}
	readyWrite.Close()
	releaseRead.Close()
	waitForCommitPipe(t, readyRead, loser)

	f.BenchEnv(map[string]string{"GIT_AUTHOR_NAME": "winner", "GIT_AUTHOR_EMAIL": "winner@local"}, "commit", "-m", "winner", "winner.txt").RequireExit(0)
	winner := headSha(f)
	if _, err := releaseWrite.Write([]byte{'x'}); err != nil {
		t.Fatal(err)
	}
	releaseWrite.Close()
	if err := loser.Wait(); err == nil || loser.ProcessState.ExitCode() != 1 {
		t.Fatalf("CAS loser exit = %v/%d, want 1\nstdout:\n%s\nstderr:\n%s", err, loser.ProcessState.ExitCode(), loserStdout.String(), loserStderr.String())
	}
	contract.RequireContains(t, loserStderr.String(), "destination compare-and-swap refused")
	contract.RequireNotContains(t, loserStderr.String(), "v1:")
	if got := headSha(f); got != winner {
		t.Fatalf("CAS loser replaced winner: got %s, want %s", got, winner)
	}
	contract.RequireContains(t, f.ReadFile("specs/concurrent/spec.md"), "Status: staged")

	f.BenchEnv(map[string]string{"GIT_AUTHOR_NAME": "loser", "GIT_AUTHOR_EMAIL": "loser@local"}, "commit", "-m", "loser", "--spec", "concurrent", "loser.txt").RequireExit(0)
	if got := strings.TrimSpace(f.Git("rev-parse", "HEAD^").Stdout); got != winner {
		t.Fatalf("fresh loser parent = %s, want winner %s", got, winner)
	}
	if got := strings.TrimSpace(f.Git("rev-parse", "HEAD^^").Stdout); got != base {
		t.Fatalf("winner parent = %s, want base %s", got, base)
	}
	if got := f.Git("diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD^").Stdout; got != "winner.txt\n" {
		t.Fatalf("winner attribution = %q", got)
	}
	if got := f.Git("diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD").Stdout; got != "loser.txt\nspecs/concurrent/spec.md\n" {
		t.Fatalf("loser attribution = %q", got)
	}
	if got := strings.TrimSpace(f.Git("show", "-s", "--format=%an <%ae>", "HEAD^").Stdout); got != "winner <winner@local>" {
		t.Fatalf("winner author = %q", got)
	}
	if got := strings.TrimSpace(f.Git("show", "-s", "--format=%an <%ae>", "HEAD").Stdout); got != "loser <loser@local>" {
		t.Fatalf("loser author = %q", got)
	}
	contract.RequireContains(t, f.ReadFile("specs/concurrent/spec.md"), "Status: implemented")
}

func waitForCommitPipe(t *testing.T, pipe *os.File, cmd *exec.Cmd) {
	t.Helper()
	ready := make(chan error, 1)
	go func() {
		var signal [1]byte
		_, err := io.ReadFull(pipe, signal[:])
		ready <- err
	}()
	select {
	case err := <-ready:
		if err != nil {
			_ = cmd.Process.Kill()
			t.Fatalf("wait for CAS synchronization: %v", err)
		}
	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("CAS loser did not reach the captured-base composition seam")
	}
}

func testCommitDetachedHEAD(t *testing.T) {
	f := commitFixture(t)
	base := headSha(f)
	f.Git("checkout", "--detach", "-q")
	f.WriteFile("work.txt", "detached bytes\n")

	f.Bench("commit", "-m", "detached", "work.txt").RequireExit(0)
	landed := headSha(f)
	if landed == base {
		t.Fatal("detached HEAD did not advance")
	}
	if got := strings.TrimSpace(f.Git("rev-parse", "HEAD^").Stdout); got != base {
		t.Fatalf("detached landing parent = %s, want %s", got, base)
	}
	if got := strings.TrimSpace(f.Git("rev-parse", "refs/heads/main").Stdout); got != base {
		t.Fatalf("detached landing moved main to %s, want %s", got, base)
	}
	if got := f.Git("status", "--porcelain", "--", "work.txt").Stdout; got != "" {
		t.Fatalf("detached named path is not clean: %q", got)
	}
}

// testCommitHostileFilenames pins both hostile filename shapes a shell CLI actually meets:
// an embedded space, which word-splits if any layer between argv and the git pathspec drops
// its quoting, and a `*`, which selects other files if the pathspec loses its `:(literal)`
// prefix. Naming only the glob-shaped file is the discriminator — as a glob that name also
// matches the space-shaped one, so a non-literal composition would include bytes the
// reviewer never named.
func testCommitHostileFilenames(t *testing.T) {
	f := commitFixture(t)
	space, glob := "a b.txt", "a*b.txt"
	f.WriteFile(space, "space\n")
	f.WriteFile(glob, "glob\n")

	f.Bench("commit", "-m", "glob-shaped name", glob).RequireExit(0)
	contract.RequireContains(t, committedNames(f), glob)
	contract.RequireContains(t, f.ReadFile(space), "space")

	f.Bench("commit", "-m", "hostile names", space, glob).RequireExit(0)
	names := committedNames(f)
	contract.RequireContains(t, names, space)
}

// testCommitDirectoryCommitsChildren drives the conventional path grammar: a commit that
// spans a directory is expressed by naming the directory, and every changed path beneath it
// is thereby explained. Two changed children, because a directory holding exactly one is
// satisfied by an implementation that never widens beyond a single path.
func testCommitDirectoryCommitsChildren(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("sub/a.txt", "a\n")
	f.WriteFile("sub/b.txt", "b\n")
	before := headSha(f)

	f.Bench("commit", "-m", "span a directory", "sub").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	names := committedNames(f)
	contract.RequireContains(t, names, "sub/a.txt")
	contract.RequireContains(t, names, "sub/b.txt")
}

// testCommitDirectoryStopsAtSegment pins the directory rule to whole path segments: `subdir`
// is a sibling of `sub`, not a path beneath it, so its changed file stays unattributed. A
// string-prefix implementation silently widens the commit to it.
func testCommitDirectoryStopsAtSegment(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("sub/a.txt", "a\n")
	f.WriteFile("subdir/x.txt", "x\n")
	before := headSha(f)

	f.Bench("commit", "-m", "span a directory", "sub").RequireExit(0)
	if headSha(f) == before {
		t.Fatal("HEAD did not advance for the named directory")
	}
	contract.RequireNotContains(t, committedNames(f), "subdir/x.txt")
	contract.RequireContains(t, f.ReadFile("subdir/x.txt"), "x")
}

// testCommitDirectoryPreservesOutsideFile pins exact attribution: a named directory lands
// what is beneath it while an unrelated path stays in the invoking checkout.
func testCommitDirectoryPreservesOutsideFile(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("sub/a.txt", "a\n")
	f.WriteFile("stray.txt", "x\n")
	before := headSha(f)

	f.Bench("commit", "-m", "span a directory", "sub").RequireExit(0)
	if headSha(f) == before {
		t.Fatal("HEAD did not advance for the named directory")
	}
	contract.RequireNotContains(t, committedNames(f), "stray.txt")
	contract.RequireContains(t, f.ReadFile("stray.txt"), "x")
}

// testCommitDeletedDirectoryCommitsRemovals pins the directory rule to the expected-base
// side: after `rm -r sub`, composition must recover every tracked descendant even though
// the working directory no longer exists. Two children prevent a one-entry special case.
func testCommitDeletedDirectoryCommitsRemovals(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("sub/a.txt", "a\n")
	f.WriteFile("sub/b.txt", "b\n")
	f.Bench("commit", "-m", "add a directory", "sub").RequireExit(0)
	contract.Remove(t, filepath.Join(f.Root, "sub"))
	before := headSha(f)

	f.Bench("commit", "-m", "retire a directory", "sub").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	names := committedNames(f)
	contract.RequireContains(t, names, "sub/a.txt")
	contract.RequireContains(t, names, "sub/b.txt")
	contract.RequireNotContains(t, f.Git("ls-files").Stdout, "sub/")
}

// A held real-index lock fails only post-publication reconciliation. The command must say
// that the commit landed so a retry cannot be mistaken for the original publication.
func testCommitReconciliationFailureIsBounded(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("work.txt", "delta\n")
	if err := os.WriteFile(filepath.Join(f.Root, ".git", "index.lock"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	before := headSha(f)

	p := f.Bench("commit", "-m", "blocked staging", "work.txt")
	p.RequireExit(1)
	contract.RequireContains(t, p.Stderr, "landed-but-checkout-incomplete")
	contract.RequireNotContains(t, p.Stderr, "v1:")
	if headSha(f) == before {
		t.Fatal("landing commit was not published before reconciliation failed")
	}
}

func testCommitSpecFlip(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("specs/feature/spec.md", "# feature\nStatus: staged\n")
	f.WriteExecutable(".bench/gate.sh", `#!/usr/bin/env bash
set -eu
gitdir="$(git rev-parse --git-common-dir)"
echo run >> "$gitdir/gate-runs"
test "$(awk '/^Status:/ { print $2 }' specs/feature/spec.md)" = implemented
cp specs/feature/spec.md "$gitdir/spec-seen"
`)
	f.CommitAll("staged feature spec and observing gate")

	f.WriteFile("work.txt", "changed\n")
	f.Bench("commit", "-m", "finish", "--spec", "feature", "work.txt").RequireExit(0)

	names := committedNames(f)
	contract.RequireContains(t, names, "work.txt")
	contract.RequireContains(t, names, "specs/feature/spec.md")
	checkout := f.ReadFile("specs/feature/spec.md")
	committed := f.Git("show", "HEAD:specs/feature/spec.md").Stdout
	seen := contract.ReadFileAbs(t, filepath.Join(gitDir(t, f), "spec-seen"))
	if checkout != committed || checkout != seen {
		t.Fatalf("implemented spec diverged: gate=%q commit=%q checkout=%q", seen, committed, checkout)
	}
	contract.RequireContains(t, checkout, "Status: implemented")
}

func testCommitSpecRedPreservesStaged(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("specs/feature/spec.md", "# feature\nStatus: staged\n")
	f.CommitAll("staged feature spec")
	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)
	beforeStatus := f.Git("status", "--porcelain", "-uall").Stdout

	p := f.BenchEnv(map[string]string{"GATE_RC": "23"}, "commit", "-m", "finish", "--spec", "feature", "work.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "prospective authorization refused")
	if headSha(f) != before || f.Git("status", "--porcelain", "-uall").Stdout != beforeStatus {
		t.Fatal("red spec authorization changed destination or checkout state")
	}
	contract.RequireContains(t, f.ReadFile("specs/feature/spec.md"), "Status: staged")
}

func testCommitSpecFailsFast(t *testing.T) {
	for _, tc := range []struct {
		name, slug, want string
		setup            func(contract.Fixture)
	}{
		{"implemented", "done", "no `Status: staged`", func(f contract.Fixture) {
			f.WriteFile("specs/done/spec.md", "# done\nStatus: implemented\n")
			f.CommitAll("implemented spec")
		}},
		{"missing", "missing", "spec not found", func(contract.Fixture) {}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := commitFixture(t)
			tc.setup(f)
			f.WriteFile("work.txt", "changed\n")
			before := headSha(f)
			beforeRuns := gateRunTally(t, f)

			p := f.BenchEnv(map[string]string{"GATE_RC": "1"}, "commit", "-m", "finish", "--spec", tc.slug, "work.txt")
			p.RequireExit(1)
			p.RequireContains(p.Stderr, tc.want)
			contract.RequireNotContains(t, p.Stderr, "prospective authorization refused")
			if headSha(f) != before || gateRunTally(t, f) != beforeRuns {
				t.Fatal("invalid --spec advanced HEAD or ran the gate")
			}
		})
	}
}

func testCommitEmptyRefused(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("foreign.txt", "base\n")
	f.CommitAll("foreign base")
	f.WriteFile("foreign.txt", "staged\n")
	f.Git("add", "foreign.txt")
	before := headSha(f)
	beforeStatus := f.Git("status", "--porcelain").Stdout
	beforeIndex := f.Git("show", ":foreign.txt").Stdout
	beforeRuns := gateRunTally(t, f)
	// seed.txt is committed and clean; the unrelated index entry cannot make it non-empty.
	p := f.Bench("commit", "-m", "nothing", "seed.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "nothing to commit")
	if headSha(f) != before {
		t.Fatal("HEAD advanced on an empty commit")
	}
	if got := gateRunTally(t, f); got != beforeRuns {
		t.Fatalf("empty prospective tree ran the gate: got %d runs, want %d", got, beforeRuns)
	}
	if got := f.Git("status", "--porcelain").Stdout; got != beforeStatus {
		t.Fatalf("empty refusal changed status: got %q, want %q", got, beforeStatus)
	}
	if got := f.Git("show", ":foreign.txt").Stdout; got != beforeIndex {
		t.Fatalf("empty refusal changed foreign index bytes: got %q, want %q", got, beforeIndex)
	}
}

func testCommitHelp(t *testing.T) {
	f := commitFixture(t)
	for _, arg := range []string{"help", "--help", "-h"} {
		p := f.Bench("commit", arg)
		p.RequireExit(0)
		p.RequireContains(p.Stdout, "usage: bench commit -m <msg> [--spec <slug>] [--] <path>...")
		if p.Stderr != "" {
			t.Fatalf("commit %s stderr = %q, want empty", arg, p.Stderr)
		}
	}
}

// TestCommitResolvesValidPathsFromNestedCWD uses distinct root and nested files because
// resolving against the repository root still succeeds, but commits the wrong a.txt.
func TestCommitResolvesValidPathsFromNestedCWD(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	f := commitFixture(t)
	f.WriteFile("a.txt", "root bytes\n")
	f.WriteFile("sub/a.txt", "nested bytes\n")
	before := headSha(f)

	contract.RunAt(t, f, filepath.Join(f.Root, "sub"), selectedBenchEnv(t, nil), benchPath(t), "commit", "-m", "nested path", "a.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance for the nested named path")
	}
	if got := strings.TrimSpace(committedNames(f)); got != "sub/a.txt" {
		t.Fatalf("committed paths = %q, want only sub/a.txt", got)
	}
	if got := f.Git("show", "HEAD:sub/a.txt").Stdout; got != "nested bytes\n" {
		t.Fatalf("published nested bytes = %q", got)
	}
	if got := f.ReadFile("a.txt"); got != "root bytes\n" {
		t.Fatalf("root a.txt bytes = %q, want preserved root bytes", got)
	}
	if got := f.Git("status", "--porcelain", "--", "sub/a.txt").Stdout; got != "" {
		t.Fatalf("nested named path is not clean: %q", got)
	}
}

// testCommitEmptyPositionalIsUsage pins the empty-positional guard from a subdirectory
// cwd, the shape that actually triggered the bug: an unset shell variable inside quotes
// expands to "", and a subcommand resolving positionals against the filesystem widens
// silently to the cwd instead of naming nothing. Two changed children under the cwd, not
// one, so a fix that merely rejects a bare "" positional but still lets it resolve to a
// directory would still be caught landing a commit nobody named.
func testCommitEmptyPositionalIsUsage(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("sub/a.txt", "a\n")
	f.WriteFile("sub/b.txt", "b\n")
	before := headSha(f)

	p := contract.RunAt(t, f, filepath.Join(f.Root, "sub"), selectedBenchEnv(t, nil), benchPath(t), "commit", "-m", "msg", "")

	// Checked ahead of the exit code and the usage line: the historical bug was a
	// commit landing files nobody named, not a wrong exit code, so this is the
	// assertion the fix actually has to satisfy.
	if headSha(f) != before {
		t.Fatal("HEAD advanced on an empty positional")
	}
	status := strings.TrimSpace(f.Git("status", "--porcelain", "-uall").Stdout)
	if status != "?? sub/a.txt\n?? sub/b.txt" && status != "?? sub/b.txt\n?? sub/a.txt" {
		t.Fatalf("index was touched by the rejected empty positional: status = %q", status)
	}
	p.RequireExit(2)
	p.RequireContains(p.Stderr, `unknown argument: ""`)
}

func testCommitUsageExitTwo(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("work.txt", "changed\n")
	for _, args := range [][]string{
		{"commit", "-m", "msg"},                       // no paths
		{"commit", "work.txt"},                        // no message
		{"commit", "-m", "msg", "--nope", "work.txt"}, // unknown flag
		{"commit", "-m"},                              // dangling -m
	} {
		p := f.Bench(args...)
		p.RequireExit(2)
		p.RequireContains(p.Stderr, "--spec marks the spec implemented")
	}
}
