package reviewrecord_test

import (
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// DI11: an author replacement keeps every earlier occurrence under the
// assignment its own frozen plan named. The amendment lands inside the next
// chunk's delta, so the source chain stays unbroken.
func TestDelegatedReplacement(t *testing.T) {
	f := recordtest.NewDelegated(t, 2)
	f.AddChunk()
	before := f.Plan.Digest
	historical := f.Record.Chunks[0].Verification[0].Performer
	if historical != recordtest.Author("1.md") {
		t.Fatalf("the first chunk did not record its own author: %q", historical)
	}

	f.Plan.Execution.Assignments["1.md"] = replaced(f.Plan.Execution.Assignments["1.md"][0], "session-lost")
	f.WritePlan()
	base := f.Tip()
	f.Write("source.txt", "implemented chunk 2\n")
	f.Commit("implement chunk 2 and amend the assignment history")
	f.Reload()
	f.Record.Amendments = append(f.Record.Amendments, rr.Amendment{
		From: before, To: f.Plan.Digest,
		ChunkIDs: map[string][]string{"1": {"1"}, "2": {"2"}},
	})
	f.Record.PlanDigest = f.Plan.Digest
	f.RecordChunk(base)
	f.Complete()
	f.Save()
	f.Commit("retain evidence across the replacement")

	if err := accept(f); err != nil {
		t.Fatalf("replacement rejected the earlier chunk's own frozen assignment: %v", err)
	}
	if got := f.Record.Chunks[0].Verification[0].Performer; got != historical {
		t.Fatalf("the historical occurrence lost its original author: %q", got)
	}
}

// DI36: once the plan names a successor, only the successor satisfies a new
// obligation. Relabelling the predecessor's pass cannot close it.
func TestDelegatedReplacementFreshness(t *testing.T) {
	f := recordtest.NewDelegated(t, 1)
	first := f.Plan.Execution.Assignments["1.md"][0]
	f.Plan.Execution.Assignments["1.md"] = replaced(first, "terminal-failure")
	f.RewritePlan()
	f.Record.PlanDigest = f.Plan.Digest
	f.AddChunk()
	f.Complete()
	f.Save()
	f.Commit("retain the successor's verification")

	successor := f.Record.Chunks[0].Verification[0].Performer
	if successor == first.Session {
		t.Fatal("the fixture recorded the predecessor after the replacement")
	}
	if err := accept(f); err != nil {
		t.Fatalf("the successor's fresh verification refused: %v", err)
	}

	for i := range f.Record.Chunks[0].Verification {
		f.Record.Chunks[0].Verification[i].Performer = first.Session
	}
	f.Save()
	f.Commit("relabel the successor obligation to the predecessor")
	if err := accept(f); err == nil || !strings.Contains(err.Error(), "verification tests") {
		t.Fatalf("a relabelled predecessor pass satisfied a successor obligation: %v", err)
	}
}
