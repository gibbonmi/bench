package runtime

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// commitFixture builds a repo on branch main with a committed green gate (whose verdict
// is `exit ${GATE_RC:-0}`, so a red run needs no untracked marker) and a configured git
// identity, so `bench commit`'s own `git commit` has an author.
func commitFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("symbolic-ref", "HEAD", "refs/heads/main")
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit ${GATE_RC:-0}\n")
	f.WriteFile("seed.txt", "seed\n")
	f.CommitAll("seed")
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
	contract.RunParallel(t, "red gate refuses commit", testCommitRedRefuses)
	contract.RunParallel(t, "unexplained file blocks before gate", testCommitUnexplainedBlocks)
	contract.RunParallel(t, "glob/space path survives whole", testCommitWeirdPath)
	contract.RunParallel(t, "spec flip lands in one commit", testCommitSpecFlip)
	contract.RunParallel(t, "bad --spec fails before the gate", testCommitSpecFailsFast)
	contract.RunParallel(t, "empty commit refused", testCommitEmptyRefused)
	contract.RunParallel(t, "usage errors exit 2", testCommitUsageExitTwo)
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
	// The no-`git add -A` guarantee: the unchanged seed never rides into the commit.
	contract.RequireNotContains(t, names, "seed.txt")
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
	}
}
