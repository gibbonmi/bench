package runtime

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// commitFixture builds a repo on branch main with a committed green gate (whose verdict
// is `exit ${GATE_RC:-0}`, so a red run needs no untracked marker) and a configured git
// identity, so `bench commit`'s own `git commit` has an author. The gate tallies each
// run in `.git/gate-runs` — inside the git dir so the tally never reads as an
// unexplained working-tree file — so verdict-reuse tests can count real gate runs.
func commitFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\necho run >> .git/gate-runs\nexit ${GATE_RC:-0}\n")
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":["GATE_RC"],"paths":[],"tools":[]}`+"\n")
	// The passlist requires every declared value to exist. Green controls declare
	// the value explicitly instead of relying on the script's shell default.
	f.Env["GATE_RC"] = "0"
	f.WriteFile("seed.txt", "seed\n")
	f.CommitAll("seed")
	return f
}

func gateRuns(t *testing.T, f contract.Fixture) int {
	t.Helper()
	return len(contract.NonEmptyLines(contract.ReadFileAbs(t, filepath.Join(f.Root, ".git", "gate-runs"))))
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
	contract.RunParallel(t, "stale verdict re-runs the gate", testCommitStaleVerdictRerunsGate)
	contract.RunParallel(t, "red gate refuses commit", testCommitRedRefuses)
	contract.RunParallel(t, "unexplained file blocks before gate", testCommitUnexplainedBlocks)
	contract.RunParallel(t, "glob/space path survives whole", testCommitHostileFilenames)
	contract.RunParallel(t, "named directory commits its changed children", testCommitDirectoryCommitsChildren)
	contract.RunParallel(t, "named directory stops at its own segment", testCommitDirectoryStopsAtSegment)
	contract.RunParallel(t, "named directory leaves an outside file blocking", testCommitDirectoryLeavesOutsideFileBlocking)
	contract.RunParallel(t, "deleted named directory commits its removals", testCommitDeletedDirectoryCommitsRemovals)
	contract.RunParallel(t, "deleted path stages the removal", testCommitStagesDeletion)
	contract.RunParallel(t, "staged deletion commits", testCommitStagesStagedDeletion)
	contract.RunParallel(t, "staged rename commits whole", testCommitStagesStagedRename)
	contract.RunParallel(t, "unknown path fails before the gate", testCommitUnknownPathFailsFast)
	contract.RunParallel(t, "spec flip lands in one commit", testCommitSpecFlip)
	contract.RunParallel(t, "bad --spec fails before the gate", testCommitSpecFailsFast)
	contract.RunParallel(t, "empty commit refused", testCommitEmptyRefused)
	contract.RunParallel(t, "usage errors exit 2", testCommitUsageExitTwo)
	contract.RunParallel(t, "-- makes a leading-dash path expressible", testCommitDoubleDashPath)
	contract.RunParallel(t, "empty positional is usage, not a path", testCommitEmptyPositionalIsUsage)
}

// testCommitDoubleDashPath drives the grammar's `--` rule end to end: a file whose name
// begins with a dash is inexpressible without it, and the rule has to survive every layer
// between argv and the git pathspec — the shell router, the parse, the block-check's
// allow-set, and staging — so only the real binary against a real repo proves it.
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
	contract.RequireIntEqual(t, gateRuns(t, f), 1, "commit re-ran the gate despite a fresh green verdict for the identical tree")
	// The reuse is silent to the operator unless the gate says so: a skipped run that
	// reports nothing reads as a gate that never ran. The line comes from the gate's own
	// emitter, the single home of verdict-reuse policy.
	contract.RequireContains(t, res.Stdout, "gate: green (fresh verdict reused for this tree)")
}

func testCommitStaleVerdictRerunsGate(t *testing.T) {
	// A verdict recorded for a different tree proves nothing about this diff: any edit
	// after the gate run must send commit back through the full gate.
	f := commitFixture(t)
	f.Bench("gate").RequireExit(0)
	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)

	f.Bench("commit", "-m", "do work", "work.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	contract.RequireIntEqual(t, gateRuns(t, f), 2, "commit trusted a verdict recorded for a different tree")
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
	// spec-retire always deletes a spec, so the sanctioned commit path must stage a removed
	// file. A staging step that skips named paths no longer on disk records no deletion, so
	// the commit finds nothing staged and exits 1 — this test goes red on that.
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
	// from both worktree and index, where a per-path `git add` pathspec matches nothing
	// and dies with git's exit 128. The staging step must recognize the change as already
	// in the index and commit it.
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
	// A named path absent from worktree, index, and HEAD is a naming error, not a staging
	// step to attempt: it must be reported before the gate runs (pair with a red gate — a
	// post-gate check would report "gate is red" instead) with a real message, not git's
	// raw exit 128.
	f := commitFixture(t)
	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)

	p := f.BenchEnv(map[string]string{"GATE_RC": "1"}, "commit", "-m", "ghost", "ghost.txt", "work.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "not found in worktree, index, or HEAD")
	contract.RequireNotContains(t, p.Stderr, "gate is red")
	if headSha(f) != before {
		t.Fatal("HEAD advanced despite an unknown named path")
	}
}

func testCommitRedRefuses(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)

	p := f.BenchEnv(map[string]string{"GATE_RC": "23"}, "commit", "-m", "do work", "work.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "gate is red")
	if headSha(f) != before {
		t.Fatal("commit landed on a red gate")
	}
}

func testCommitUnexplainedBlocks(t *testing.T) {
	// Both an untracked and a tracked-modified file outside the named set must block,
	// and — with the gate left green — the block must fire regardless of the verdict.
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

			p := f.Bench("commit", "-m", "do work", "work.txt")
			p.RequireExit(1)
			p.RequireContains(p.Stderr, "outside the named set")
			if headSha(f) != before {
				t.Fatal("HEAD advanced despite an unexplained working-tree file")
			}
		})
	}
}

// testCommitHostileFilenames pins both hostile filename shapes a shell CLI actually meets:
// an embedded space, which word-splits if any layer between argv and the git pathspec drops
// its quoting, and a `*`, which selects other files if the pathspec loses its `:(literal)`
// prefix. Naming only the glob-shaped file is the discriminator — as a glob that name also
// matches the space-shaped one, so a non-literal chain would explain a file the reviewer
// never named.
func testCommitHostileFilenames(t *testing.T) {
	f := commitFixture(t)
	space, glob := "a b.txt", "a*b.txt"
	f.WriteFile(space, "space\n")
	f.WriteFile(glob, "glob\n")

	p := f.Bench("commit", "-m", "glob-shaped name", glob)
	p.RequireExit(1)
	p.RequireContains(p.Stderr, space)

	f.Bench("commit", "-m", "hostile names", space, glob).RequireExit(0)
	names := committedNames(f)
	contract.RequireContains(t, names, space)
	contract.RequireContains(t, names, glob)
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
// is a sibling of `sub`, not a path beneath it, so its changed file stays unexplained. A
// string-prefix implementation silently widens the commit to it.
func testCommitDirectoryStopsAtSegment(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("sub/a.txt", "a\n")
	f.WriteFile("subdir/x.txt", "x\n")
	before := headSha(f)

	p := f.Bench("commit", "-m", "span a directory", "sub")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "subdir/x.txt")
	if headSha(f) != before {
		t.Fatal("HEAD advanced despite a changed file under a sibling directory")
	}
}

// testCommitDirectoryLeavesOutsideFileBlocking pins the safety property the block-check
// exists for: a named directory explains what is beneath it and nothing else, so a change
// anywhere outside it still refuses the commit by name.
func testCommitDirectoryLeavesOutsideFileBlocking(t *testing.T) {
	f := commitFixture(t)
	f.WriteFile("sub/a.txt", "a\n")
	f.WriteFile("stray.txt", "x\n")
	before := headSha(f)

	p := f.Bench("commit", "-m", "span a directory", "sub")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "outside the named set")
	p.RequireContains(p.Stderr, "stray.txt")
	if headSha(f) != before {
		t.Fatal("HEAD advanced despite a changed file outside the named directory")
	}
}

// testCommitDeletedDirectoryCommitsRemovals pins the directory rule to the side of the
// removal: `rm -r sub` then naming `sub` is how a directory is retired, and the block-check
// has to explain the deletions the removal produced. An attribution that asks only the
// working tree finds nothing to widen from — the directory is gone — and reports every one
// of its own children as an unexplained offender. Two children, because a directory holding
// one is satisfied by an implementation that never widens past a single path.
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

func testCommitSpecFlip(t *testing.T) {
	f := commitFixture(t)
	// The staged spec already exists and is tracked at build finish; commit it clean
	// first, then the build modifies work.txt and flips the spec in one commit.
	f.WriteFile("specs/feature/spec.md", "# feature\nStatus: staged\n")
	f.Bench("commit", "-m", "add spec", "specs/feature").RequireExit(0)

	f.WriteFile("work.txt", "changed\n")
	f.Bench("commit", "-m", "finish", "--spec", "feature", "work.txt").RequireExit(0)

	names := committedNames(f)
	contract.RequireContains(t, names, "work.txt")
	contract.RequireContains(t, names, "specs/feature/spec.md")
	contract.RequireContains(t, f.ReadFile("specs/feature/spec.md"), "Status: implemented")
}

func testCommitSpecFailsFast(t *testing.T) {
	// A --spec with no `Status: staged` line must be rejected before the gate runs, not after
	// a green one burns. Pair the bad spec with a red gate: if the check preceded the gate the
	// flip error surfaces; a check that ran only after the gate would report "gate is red".
	f := commitFixture(t)
	f.WriteFile("specs/done/spec.md", "# done\nStatus: implemented\n")
	f.Bench("commit", "-m", "add spec", "specs/done/spec.md").RequireExit(0)

	f.WriteFile("work.txt", "changed\n")
	before := headSha(f)
	p := f.BenchEnv(map[string]string{"GATE_RC": "1"}, "commit", "-m", "finish", "--spec", "done", "work.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "no `Status: staged`")
	contract.RequireNotContains(t, p.Stderr, "gate is red")
	if headSha(f) != before {
		t.Fatal("HEAD advanced despite a bad --spec")
	}
}

func testCommitEmptyRefused(t *testing.T) {
	f := commitFixture(t)
	before := headSha(f)
	// seed.txt is committed and clean; naming it produces no staged change.
	p := f.Bench("commit", "-m", "nothing", "seed.txt")
	p.RequireExit(1)
	p.RequireContains(p.Stderr, "nothing to commit")
	if headSha(f) != before {
		t.Fatal("HEAD advanced on an empty commit")
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

	p := contract.RunAt(t, f, filepath.Join(f.Root, "sub"), nil, "bash", benchPath(t), "commit", "-m", "msg", "")

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
