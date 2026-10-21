// The worktree merge verb's built-in lane: a merge from any checkout of the kit repository
// grades the composed tree with the target's selective lane, and an adopted kit keeps the
// whole gate.
package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// goLaneMergeFixture is the whole-gate merge fixture over a repository whose base commit
// already holds a Go module, so the built-in lane of each target finds one trivial package.
// The whole gate appends to the set's tally and then runs gateExit.
func goLaneMergeFixture(t *testing.T, gateExit string) mergeSet {
	t.Helper()
	root := newWorktreeRepo(t)
	writeGoMainFixture(t, root)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-q", "-m", "declare the go module")
	tally := filepath.Join(t.TempDir(), "gate-tally")
	f := wholeGateMergeSetAt(t, root, "printf w >> "+sanitize.ShellQuote(tally)+"; "+gateExit)
	f.tally = tally
	return f
}

// mergeGoModLine is the merge of a default-branch commit that appends line to the
// fixture's go.mod into the set's target. A go.mod change selects no Bench-owned check of
// the built-in lane, so the lane runs with no Bench build.
func mergeGoModLine(t *testing.T, f mergeSet, line string) verbCall {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(f.root, "go.mod"))
	mustNoError(t, err)
	return f.merge(mergeTargetArgs(f, commitOnDefault(t, f.root, "go.mod", string(body)+line+"\n"))...)
}

// requireSelectiveLanePass fails t unless the merge exited 0 with the selective lane's pass
// line and the whole gate left no tally.
func requireSelectiveLanePass(t *testing.T, f mergeSet, r verbResult) {
	t.Helper()
	if r.exit != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	var lane string
	for _, line := range strings.Split(r.stdout, "\n") {
		if strings.HasPrefix(line, "lane{outcome=pass,") {
			lane = line
		}
	}
	if !strings.Contains(lane, "classes=") {
		t.Errorf("stdout = %q, want a selective lane pass line", r.stdout)
	}
	requireNoWholeGate(t, f)
}

// requireNoWholeGate marks t failed when the whole gate wrote its tally, and lets the
// caller's later checks still run.
func requireNoWholeGate(t *testing.T, f mergeSet) {
	t.Helper()
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Errorf("whole-gate tally = %v, want no whole-gate run", err)
	}
}

// A merge whose kit is the primary checkout grades with the built-in lane of the target, as
// a merge from inside the target does. The caller's kit declares no lane for the target, so
// a merge that keeps it runs the whole gate.
func TestMergeFromThePrimaryGradesWithTheTargetLane(t *testing.T) {
	t.Parallel()
	f := goLaneMergeFixture(t, "exit 0")
	f.kit = f.root
	requireSelectiveLanePass(t, f, runVerb(t, verbMerge, mergeGoModLine(t, f, "// lane fixture")))
}

// The fold's lane runs go vet on an unknown go.mod directive and fails. The target tip alone
// passes that lane, so the face names the fold. A grade of the target alone under the
// caller's kit runs the whole gate, which fails and names the target. The two faces share a
// sentence, so the printed route must equal the route that the registry renders for the
// fold face over the merge's own facts.
func TestMergeFromThePrimaryGradesTheTargetAloneWithTheFoldLane(t *testing.T) {
	t.Parallel()
	produced := produceMergeFace(t, mergeRefusalFixture{face: faceMergeFoldRed, exit: 1,
		build: func(t *testing.T) (mergeSet, []string) {
			f := goLaneMergeFixture(t, "exit 1")
			f.kit = f.root
			return f, mergeGoModLine(t, f, "unknowndirective x").args
		}})
	requireNoWholeGate(t, produced.f)
	from, target := produced.args[1], produced.f.created[0].Assignment
	facts := conflictFacts(from, target.ID, target.Worktree)
	facts[refusalroute.FactLabel], facts[refusalroute.FactFrom] = target.Label, from
	if want := landingRoute(faceMergeFoldRed, mergeRerun(from, target.ID), facts); produced.next != want {
		t.Fatalf("fold red next = %q, want the %s route %q", produced.next, faceMergeFoldRed, want)
	}
}

// A merge whose kit is another linked worktree of the target's repository grades with
// the target's built-in lane too.
func TestMergeFromASiblingWorktreeGradesWithTheTargetLane(t *testing.T) {
	t.Parallel()
	f := goLaneMergeFixture(t, "exit 0")
	f.kit = filepath.Join(t.TempDir(), "sibling")
	gitRun(t, f.root, "worktree", "add", "-q", "--detach", f.kit, "HEAD")
	requireSelectiveLanePass(t, f, runVerb(t, verbMerge, mergeGoModLine(t, f, "// lane fixture")))
}

// An adopted kit at <root>/.bench is no top level of a checkout, so it stays the kit.
// The target declares no lane, so the whole gate grades the merge.
func TestMergeUnderAnAdoptedKitKeepsTheWholeGate(t *testing.T) {
	t.Parallel()
	f := goLaneMergeFixture(t, "exit 0")
	f.kit = filepath.Join(f.root, ".bench")
	r := runVerb(t, verbMerge, mergeGoModLine(t, f, "// lane fixture"))
	if r.exit != 0 || strings.Contains(r.stdout, "lane{") {
		t.Fatalf("merge = (%d, %q, %q), want exit 0 and no lane line", r.exit, r.stdout, r.stderr)
	}
	if _, err := os.Stat(f.tally); err != nil {
		t.Fatalf("whole-gate tally = %v, want one whole-gate run", err)
	}
}
