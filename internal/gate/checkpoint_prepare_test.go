package gate

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/prose"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestChunkCheckpointRefusesADirtyCheckout covers RC44: the chunk checkpoint refuses a working
// tree that differs from HEAD before the evidence pass and before the oracle. The review
// record has no exemption, so a valid record change alone also refuses.
func TestChunkCheckpointRefusesADirtyCheckout(t *testing.T) {
	for _, tc := range []struct {
		name  string
		dirty func(*testing.T, *recordtest.Fixture)
	}{
		{"tracked edit", func(_ *testing.T, f *recordtest.Fixture) { f.Write("tracked.txt", "edited\n") }},
		{"uncommitted record", annotateRecord},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := checkpointFixture(t)
			tc.dirty(t, f)
			code, out := runCheckpoint(t, f)
			requireCleanCheckoutRefusal(t, f, code, out, "stale reviewed source")
		})
	}
}

// TestChunkCheckpointReuseKeepsTheEvidencePass covers RC39. After a record-only commit, a
// second checkpoint call reads the evidence store again and reuses the chunk green with no
// run. A stale verification entry then refuses on that same reusable green through the
// evidence pass.
func TestChunkCheckpointReuseKeepsTheEvidencePass(t *testing.T) {
	f := checkpointFixture(t)
	if code, out := runCheckpoint(t, f); code != 0 {
		t.Fatalf("first checkpoint = %d %s", code, out)
	}
	annotateRecord(t, f)
	f.Commit("annotate the review record")
	requireReuse(t, f, "after a record-only commit")
	f.Record.Chunks[0].Verification[0].SourceDigest = f.Record.Chunks[0].Base
	resave(t, f, "retain the stale entry")
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "chunk 1: verification tests is incomplete, failed, or stale") || outcomeRuns(t, f.Root) != 1 {
		t.Fatalf("checkpoint over a stale entry = %d %s after %d runs; want the stale entry refused and one run", code, out, outcomeRuns(t, f.Root))
	}
}

// TestChunkCheckpointGradesTheRecordProseOnReuse covers RC40. A record-only commit leaves the
// chunk green reusable, so a sentence over the prose bound in that commit reds the checkpoint
// with no second run. The corrected record then reuses the green.
func TestChunkCheckpointGradesTheRecordProseOnReuse(t *testing.T) {
	f := proseCheckpointFixture(t)
	if code, out := runCheckpoint(t, f); code != 0 {
		t.Fatalf("first checkpoint = %d %s", code, out)
	}
	overBoundRecord(t, f)
	f.Commit("retain record prose over the bound")
	code, out := runCheckpoint(t, f)
	requireRecordProseRefusal(t, code, out)
	resave(t, f, "correct the record prose")
	requireReuse(t, f, "after the record prose is corrected")
}

// TestChunkCheckpointGradesTheRecordProseOnARun covers RC41. The run grades the tree without
// the record, so a fresh run reds on a sentence over the prose bound only through the
// checkpoint's own record grade.
func TestChunkCheckpointGradesTheRecordProseOnARun(t *testing.T) {
	f := proseCheckpointFixture(t)
	overBoundRecord(t, f)
	f.Commit("retain record prose over the bound")
	code, out := runCheckpointWith(t, f, flagFresh)
	requireRecordProseRefusal(t, code, out)
}

// proseCheckpointFixture is a checkpoint fixture whose committed tree holds an empty prose
// policy, so the record grade excludes nothing.
func proseCheckpointFixture(t *testing.T) *recordtest.Fixture {
	return attachedCheckpointFixture(t, recordtest.Attach, func(f *recordtest.Fixture) { f.Write(prose.ExclusionFile, "") })
}

// overBoundRecord appends to the review record of f one sentence of one word over the prose
// bound.
func overBoundRecord(t *testing.T, f *recordtest.Fixture) {
	appendRecordProse(t, f, strings.Repeat("word ", prose.MaxSentenceWords)+"word.")
}

// requireReuse requires that a chunk checkpoint of f answers with the reused green and leaves
// the oracle at its one run.
func requireReuse(t *testing.T, f *recordtest.Fixture, when string) {
	t.Helper()
	if code, out := runCheckpoint(t, f); code != 0 || !strings.Contains(out, "reused") || outcomeRuns(t, f.Root) != 1 {
		t.Fatalf("checkpoint %s = %d %s after %d runs; want the reused green and one run", when, code, out, outcomeRuns(t, f.Root))
	}
}

// requireRecordProseRefusal requires a red checkpoint that names the sentence over the bound in
// the review record.
func requireRecordProseRefusal(t *testing.T, code int, out string) {
	t.Helper()
	path, err := rr.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if bound := fmt.Sprintf("over the %d-word bound", prose.MaxSentenceWords); code == 0 || !strings.Contains(out, strconv.Quote(path)) || !strings.Contains(out, bound) {
		t.Fatalf("checkpoint over record prose = %d %s; want a red that names %s and %q", code, out, path, bound)
	}
}
