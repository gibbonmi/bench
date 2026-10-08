package gate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/landing/published"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// citation is the diagnostic of citationGate.
const citation = "decisions/x.md cites missing roadmap/FT1.md"

// citationGate stands in for decision-map integrity: the project gate exits 9 when
// decisions/x.md exists and the roadmap detail file that it cites does not.
func citationGate(_ *testrepo.GateFixture, script string) string {
	return script + `if [ -e decisions/x.md ] && [ ! -e roadmap/FT1.md ]; then echo '` + citation + `' >&2; exit 9; fi
`
}

// completedSource commits one reviewed and completed chunk under the configured project
// gate. The chunk commit carries each file that prepare writes.
func completedSource(t *testing.T, gate func(*testrepo.GateFixture, string) string, prepare func(*recordtest.Fixture)) *recordtest.Fixture {
	t.Helper()
	f := recordtest.Attach(t, outcomeFixture(t, gate), 1)
	prepare(f)
	f.AddChunk()
	f.Save()
	f.Commit("retain review evidence")
	retainCompletion(f)
	return f
}

// closureSource is a completed source whose delivery closes roadmap row FT1, and the
// closure deletes roadmap/FT1.md. With cite, the source also holds a decision map that
// cites that file, so only the closed tree is red under citationGate.
func closureSource(t *testing.T, cite bool) *recordtest.Fixture {
	t.Helper()
	return completedSource(t, citationGate, func(f *recordtest.Fixture) {
		commitmenttest.SeedClosure(t, f.Root, recordtest.Spec)
		if cite {
			f.Write("decisions/x.md", "# Map\n\nSources:\n- Path: `roadmap/FT1.md`\n")
		}
	})
}

func retainCompletion(f *recordtest.Fixture) {
	f.T.Helper()
	f.Complete()
	f.Save()
	f.Commit("retain complete evidence")
}

func completeCheckpoint(t *testing.T, f *recordtest.Fixture, flags ...string) (int, string) {
	t.Helper()
	var out, diagnostic bytes.Buffer
	code := RunCommand(append([]string{f.Root, "--checkpoint", recordtest.Spec, "--complete"}, flags...), &out, &diagnostic)
	return code, out.String() + diagnostic.String()
}

// publishedEvidence inspects the retained evidence for the tree that the landing
// publishes from HEAD, under the landing's completion obligation.
func publishedEvidence(t *testing.T, f *recordtest.Fixture) EvidenceInspection {
	t.Helper()
	tip := f.Tip()
	tree, err := published.Tree(f.Root, f.Tree(), recordtest.Spec, tip)
	if err != nil {
		t.Fatalf("closure transform: %v", err)
	}
	return InspectTreeContext(WithCompletion(context.Background(), recordtest.Spec, tip), f.Root, tree)
}

// annotateRecord adds prose to the review record and leaves the evidence unchanged.
func annotateRecord(t *testing.T, f *recordtest.Fixture) {
	t.Helper()
	path, err := rr.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(path, string(outcomeRead(t, filepath.Join(f.Root, path)))+"\nThe orchestrator reconciled the final evidence.\n")
}

func assertNoOracleRun(t *testing.T, f *recordtest.Fixture) {
	t.Helper()
	if _, err := os.Stat(runCountWitness.path(t, f.Root)); !os.IsNotExist(err) {
		t.Fatalf("oracle ran before the refusal: %v", err)
	}
}

func TestCompleteCheckpointGradesTheClosedTree(t *testing.T) {
	f := closureSource(t, true)
	if result := Execute(context.Background(), f.Root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("ordinary gate on the source checkout = %d, want 0", result.ActionExit)
	}

	code, out := completeCheckpoint(t, f)
	if code != 9 || !strings.Contains(out, citation) {
		t.Fatalf("complete checkpoint on a source whose delivery closure is red = (%d, %q), want exit 9 and %q", code, out, citation)
	}

	// Witness: the landing's own closure of this exact source is red.
	closed, err := published.Tree(f.Root, f.Tree(), recordtest.Spec, f.Tip())
	if err != nil {
		t.Fatalf("closure transform: %v", err)
	}
	f.Git("read-tree", "--reset", "-u", closed)
	var closedErr bytes.Buffer
	if result := Execute(context.Background(), f.Root, &bytes.Buffer{}, &closedErr); result.ActionExit != 9 || !strings.Contains(closedErr.String(), citation) {
		t.Fatalf("witness lost: the closed tree is not red (exit %d): %s", result.ActionExit, closedErr.String())
	}
}

func TestChunkCheckpointGradesTheCheckoutTree(t *testing.T) {
	f := closureSource(t, true)
	if code, out := runCheckpoint(t, f); code != 0 {
		t.Fatalf("chunk checkpoint on the unclosed source = (%d, %q), want 0", code, out)
	}
}

func TestCompleteCheckpointEvidenceNamesThePublishedTree(t *testing.T) {
	f := closureSource(t, false)
	if code, out := completeCheckpoint(t, f); code != 0 {
		t.Fatalf("complete checkpoint: %d %s", code, out)
	}
	if got := publishedEvidence(t, f); !got.ReusableGreen {
		t.Fatalf("published tree of HEAD has no reusable green under the completion obligation: %+v", got)
	}
	if _, err := os.Stat(recordDuringWitness.path(t, f.Root)); err != nil {
		t.Fatalf("the project gate did not read the run record from the common directory: %v", err)
	}

	annotateRecord(t, f)
	f.Commit("annotate the review record")
	if got := publishedEvidence(t, f); got.ReusableGreen {
		t.Fatalf("published tree of the new tip is green before its checkpoint: %+v", got)
	}
	if code, out := completeCheckpoint(t, f); code != 0 {
		t.Fatalf("complete checkpoint on the new tip: %d %s", code, out)
	}
	if got := publishedEvidence(t, f); !got.ReusableGreen {
		t.Fatalf("published tree of the new tip has no reusable green: %+v", got)
	}
}

func TestCompleteCheckpointGradesTheImplementedStatus(t *testing.T) {
	stagedGate := func(g *testrepo.GateFixture, script string) string {
		return script + `if ` + g.Command("grep") + ` -q '^Status: staged$' '` + recordtest.Spec + `'; then exit 8; fi
`
	}
	f := completedSource(t, stagedGate, func(*recordtest.Fixture) {})
	if code, out := completeCheckpoint(t, f); code != 0 {
		t.Fatalf("complete checkpoint graded the staged status: %d %s", code, out)
	}
}

func TestCompleteCheckpointReuseAndFresh(t *testing.T) {
	f := closureSource(t, false)
	if code, out := completeCheckpoint(t, f); code != 0 {
		t.Fatalf("complete checkpoint: %d %s", code, out)
	}
	code, out := completeCheckpoint(t, f)
	if code != 0 || !strings.Contains(out, "gate: green (fresh verdict reused for this tree)") {
		t.Fatalf("repeated complete checkpoint = (%d, %q), want the reuse line", code, out)
	}
	if got := outcomeRuns(t, f.Root); got != 1 {
		t.Fatalf("repeated complete checkpoint ran the oracle: %d runs, want 1", got)
	}
	if code, out := completeCheckpoint(t, f, "--fresh"); code != 0 {
		t.Fatalf("fresh complete checkpoint: %d %s", code, out)
	}
	if got := outcomeRuns(t, f.Root); got != 2 {
		t.Fatalf("fresh complete checkpoint: %d runs, want 2", got)
	}
}

func TestCompleteCheckpointLeavesTheCheckout(t *testing.T) {
	f := closureSource(t, true)
	head := f.Tip()
	if code, out := completeCheckpoint(t, f); code == 0 {
		t.Fatalf("complete checkpoint on a red closure passed: %s", out)
	}
	if got := f.Tip(); got != head {
		t.Fatalf("complete checkpoint moved HEAD from %s to %s", head, got)
	}
	if status := f.Git("status", "--porcelain"); status != "" {
		t.Fatalf("complete checkpoint changed the checkout:\n%s", status)
	}
}

func TestCompleteCheckpointRefusesADirtyCheckout(t *testing.T) {
	completed := func(t *testing.T) *recordtest.Fixture {
		f := checkpointFixture(t)
		retainCompletion(f)
		return f
	}
	t.Run("ignored file", func(t *testing.T) {
		f := completed(t)
		f.Write(".gate-local", "ignored\n")
		if code, out := completeCheckpoint(t, f); code != 0 || strings.Contains(out, cleanCheckoutRefusal) {
			t.Errorf("complete checkpoint with only an ignored file = (%d, %q), want exit 0", code, out)
		}
	})
	for _, tc := range []struct {
		name  string
		dirty func(*testing.T) *recordtest.Fixture
	}{
		{"tracked edit", func(t *testing.T) *recordtest.Fixture {
			f := completed(t)
			f.Write("tracked.txt", "edited\n")
			return f
		}},
		{"untracked file", func(t *testing.T) *recordtest.Fixture {
			f := completed(t)
			f.Write("untracked.txt", "stray\n")
			return f
		}},
		{"uncommitted record", func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Complete()
			f.Save()
			return f
		}},
		{"stale evidence", func(t *testing.T) *recordtest.Fixture {
			f := attachedCheckpointFixture(t, recordtest.Attach, seedCommentGap)
			retainCompletion(f)
			correctCommentGap(f)
			f.Write("tracked.txt", "edited\n")
			return f
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := tc.dirty(t)
			code, out := completeCheckpoint(t, f)
			if code != 1 || !strings.Contains(out, cleanCheckoutRefusal) || strings.Contains(out, "completion is incomplete or stale") {
				t.Errorf("dirty complete checkpoint = (%d, %q), want exit 1 and only %q", code, out, cleanCheckoutRefusal)
			}
			// RR41: the refusal names the commit that cleans the checkout.
			if !strings.Contains(out, "\nnext=bench commit --in ") {
				t.Errorf("dirty complete checkpoint = %q, want a next= route that commits the checkout", out)
			}
			assertNoOracleRun(t, f)
		})
	}
}

func TestCompleteCheckpointRefusesAnUntransformableSpec(t *testing.T) {
	f := attachedCheckpointFixture(t, func(t testing.TB, root string, count int) *recordtest.Fixture {
		return recordtest.Prepare(t, root, count, recordtest.Spec, "# Example\n\n")
	})
	retainCompletion(f)
	code, out := completeCheckpoint(t, f)
	if code != 1 || !strings.Contains(out, "spec has no Status: staged line") {
		t.Errorf("complete checkpoint on a spec with no staged status = (%d, %q), want exit 1 and the transform reason", code, out)
	}
	// RR66: the spec status is the reviewer's, so the route hands back.
	if !strings.Contains(out, "\nnext=reviewer: ") {
		t.Errorf("complete checkpoint on a spec with no staged status = %q, want a next= reviewer route", out)
	}
	assertNoOracleRun(t, f)
}
