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
	contract.RunParallel(t, "glob/space path survives whole", testCommitWeirdPath)
	contract.RunParallel(t, "deleted path stages the removal", testCommitStagesDeletion)
	contract.RunParallel(t, "staged deletion commits", testCommitStagesStagedDeletion)
	contract.RunParallel(t, "staged rename commits whole", testCommitStagesStagedRename)
	contract.RunParallel(t, "unknown path fails before the gate", testCommitUnknownPathFailsFast)
	contract.RunParallel(t, "spec flip lands in one commit", testCommitSpecFlip)
	contract.RunParallel(t, "bad --spec fails before the gate", testCommitSpecFailsFast)
	contract.RunParallel(t, "empty commit refused", testCommitEmptyRefused)
	contract.RunParallel(t, "usage errors exit 2", testCommitUsageExitTwo)
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

	f.Bench("commit", "-m", "do work", "work.txt").RequireExit(0)

	if headSha(f) == before {
		t.Fatal("HEAD did not advance on a green gate")
	}
	contract.RequireIntEqual(t, gateRuns(t, f), 1, "commit re-ran the gate despite a fresh green verdict for the identical tree")
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

	p := f.BenchEnv(map[string]string{"GATE_RC": "1"}, "commit", "-m", "do work", "work.txt")
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

func testCommitWeirdPath(t *testing.T) {
	f := commitFixture(t)
	weird := "a b*c.txt"
	f.WriteFile(weird, "weird\n")

	f.Bench("commit", "-m", "weird path", weird).RequireExit(0)
	contract.RequireContains(t, committedNames(f), weird)
}

func testCommitSpecFlip(t *testing.T) {
	f := commitFixture(t)
	// The staged spec already exists and is tracked at build finish; commit it clean
	// first, then the build modifies work.txt and flips the spec in one commit.
	f.WriteFile("specs/feature.md", "# feature\nStatus: staged\n")
	f.Bench("commit", "-m", "add spec", "specs/feature.md").RequireExit(0)

	f.WriteFile("work.txt", "changed\n")
	f.Bench("commit", "-m", "finish", "--spec", "feature", "work.txt").RequireExit(0)

	names := committedNames(f)
	contract.RequireContains(t, names, "work.txt")
	contract.RequireContains(t, names, "specs/feature.md")
	contract.RequireContains(t, f.ReadFile("specs/feature.md"), "Status: implemented")
}

func testCommitSpecFailsFast(t *testing.T) {
	// A --spec with no `Status: staged` line must be rejected before the gate runs, not after
	// a green one burns. Pair the bad spec with a red gate: if the check preceded the gate the
	// flip error surfaces; a check that ran only after the gate would report "gate is red".
	f := commitFixture(t)
	f.WriteFile("specs/done.md", "# done\nStatus: implemented\n")
	f.Bench("commit", "-m", "add spec", "specs/done.md").RequireExit(0)

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
