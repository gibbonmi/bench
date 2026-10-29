package gate

// Each whole-gate verb of this package looks up a green under the grade key before it picks a
// reduced schedule. The fixtures count the runs of the fixture gate, so a reuse adds none.

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// GR58: the running executable is in the grade key. A complete checkpoint green answers a
// prospective run of its tree under the same executable, and a run under another executable
// runs the whole gate.
func TestCheckpointUnderAnotherExecutableRunsTheGate(t *testing.T) {
	f := closureSource(t, false)
	setRunningExecutable(t, strings.Repeat("a", 64))
	if code, out := completeCheckpoint(t, f); code != 0 {
		t.Fatalf("complete checkpoint = (%d, %q), want green", code, out)
	}
	tree := publishedTree(t, f)
	gradeKeyRun(t, f.Root, tree, reuseFreshGreen, 0)
	requireRuns(t, f.Root, 1, "after a prospective run under the checkpoint's executable")
	setRunningExecutable(t, strings.Repeat("b", 64))
	gradeKeyRun(t, f.Root, tree, reuseFreshGreen, 0)
	requireRuns(t, f.Root, 2, "after a prospective run under another executable")
}

// GR48: a plain `bench gate` in the primary checkout reuses the green of a prospective run of
// the same tree in a linked worktree.
func TestGateReusesASiblingGreen(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	var stdout, stderr bytes.Buffer
	if code := RunCommand([]string{primary}, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), reuseLine) {
		t.Fatalf("plain gate = (%d, %q, %q), want the reuse line", code, stdout.String(), stderr.String())
	}
	requireRuns(t, primary, 1, "after a plain gate of the linked worktree's green tree")
}

// GR52: a kit landing that is eligible for the reader set reuses a whole green of its tree. The
// lookup comes before the reader plan, so the landing writes no reader record.
func TestReaderSetLandingReusesAWholeGreen(t *testing.T) {
	f := newLandingFixture(t, nil)
	f.chain(1, -1)
	docsChange(f)
	f.commit("change")
	if result := Execute(context.Background(), f.root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("whole green = %+v", result)
	}
	before := outcomeRuns(t, f.root)
	selection := stubSelection(t, f.root)
	owner := func(context.Context, string) (*runbinary.Selection, error) { return selection, nil }
	var stdout, stderr bytes.Buffer
	result := executeSubjectReplacing(f.landing(), f.root, f.root, &stdout, &stderr, nil, reuseFreshGreen, newGateEvaluation(f.root), owner, "", durablefile.Native())
	tree := gittest.Output(t, f.root, "rev-parse", "HEAD^{tree}")
	if recorded, err := ReaderRecorded(f.root, tree); result.ActionExit != 0 || stdout.String() != reuseLine || recorded || err != nil {
		t.Fatalf("landing = %+v, recorded %t %v\nstdout:\n%s\nstderr:\n%s\nwant the reuse line and no reader record", result, recorded, err, stdout.String(), stderr.String())
	}
	requireRuns(t, f.root, before, "after a landing of a whole green tree")
}
