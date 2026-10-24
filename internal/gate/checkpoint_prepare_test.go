package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
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
	commitOverBoundRecord(t, f)
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
	commitOverBoundRecord(t, f)
	code, out := runCheckpointWith(t, f, flagFresh)
	requireRecordProseRefusal(t, code, out)
}

// TestChunkCheckpointSkipsTheRecordGradeWithNoPolicy grades a HEAD tree with no prose policy.
// The prose package holds that tree ungradeable, so a record sentence over the prose bound
// passes the chunk checkpoint there.
func TestChunkCheckpointSkipsTheRecordGradeWithNoPolicy(t *testing.T) {
	f := checkpointFixture(t)
	commitOverBoundRecord(t, f)
	if code, out := runCheckpoint(t, f); code != 0 {
		t.Fatalf("checkpoint over record prose with no policy = %d %s; want the record grade skipped", code, out)
	}
}

// TestChunkCheckpointReadsTheRecordProsePolicy grades the policy that the record grade reads
// from the HEAD tree. A policy row with a real target holds: a row that excludes the record
// keeps the record out of the grade, and a row that names another file leaves it graded. A
// policy or a record that does not read, and a malformed policy row, refuse before the oracle.
func TestChunkCheckpointReadsTheRecordProsePolicy(t *testing.T) {
	unreadablePolicy := prose.UnreadableExclusionDiagnostic("")
	for _, tc := range []struct {
		name, policy, text string
		// want is a clause of the refusal; it is empty for a pass.
		want string
	}{
		{"row that excludes the record", "reviews/ the record stays out of the grade\n", overBoundProse, ""},
		{"row that names another file", "tracked.txt the fixture file\n", overBoundProse, boundFinding},
		{"invalid UTF-8 policy", "tracked.txt \xff\n", "", unreadablePolicy},
		{"oversized policy", strings.Repeat("a", int(bounds.ControlRecordLimit)+1), "", unreadablePolicy},
		{"malformed policy row", "tracked.txt\n", "", "malformed exclusion row"},
		{"invalid UTF-8 record", "", "Short \xff prose.", prose.UnreadableSubjectDiagnostic(fixtureRecord(t), "")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := attachedCheckpointFixture(t, recordtest.Attach, func(f *recordtest.Fixture) { f.Write(prose.ExclusionFile, tc.policy) })
			if tc.text != "" {
				appendRecordProse(t, f, tc.text)
				f.Commit("retain the record prose")
			}
			code, out := runCheckpoint(t, f)
			if tc.want == "" {
				if code != 0 {
					t.Fatalf("checkpoint = %d %s; want a pass", code, out)
				}
				return
			}
			if code == 0 || !strings.Contains(out, tc.want) {
				t.Fatalf("checkpoint = %d %s; want a refusal that names %q", code, out, tc.want)
			}
			assertNoOracleRun(t, f)
		})
	}
}

// TestChunkCheckpointRefusesALinkedRecord grades a symbolic link at the record path of the
// HEAD tree. The evidence pass reads the record before the record grade, so the refusal of
// the evidence pass names the link, and the oracle does not run.
func TestChunkCheckpointRefusesALinkedRecord(t *testing.T) {
	f := proseCheckpointFixture(t)
	record := filepath.Join(f.Root, fixtureRecord(t))
	if err := os.Rename(record, record+".target"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Base(record)+".target", record); err != nil {
		t.Fatal(err)
	}
	f.Commit("link the review record")
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "completion evidence: ") || !strings.Contains(out, "nonregular tree file "+fixtureRecord(t)) {
		t.Fatalf("checkpoint over a linked record = %d %s; want the evidence refusal of the link", code, out)
	}
	assertNoOracleRun(t, f)
}

// TestChunkCheckpointRunsTheProspectiveGate grades a tree that holds the prospective gate. The
// chunk checkpoint runs its subject through the prospective owner, so the prospective gate
// that its key hashes runs, and the ordinary gate does not.
func TestChunkCheckpointRunsTheProspectiveGate(t *testing.T) {
	f := attachedCheckpointFixture(t, recordtest.Attach, func(f *recordtest.Fixture) {
		outcomeWrite(t, f.Root, prospectiveGatePath, "#!/bin/sh\nexit 6\n", 0o755)
	})
	if code, out := runCheckpoint(t, f); code != 6 {
		t.Fatalf("checkpoint under a prospective gate that exits 6 = %d %s; want 6", code, out)
	}
	assertNoOracleRun(t, f)
}

// proseCheckpointFixture is a checkpoint fixture whose committed tree holds an empty prose
// policy, so the record grade excludes nothing.
func proseCheckpointFixture(t *testing.T) *recordtest.Fixture {
	return attachedCheckpointFixture(t, recordtest.Attach, prosePolicy)
}

// prosePolicy writes an empty prose policy into the checkout of f.
func prosePolicy(f *recordtest.Fixture) { f.Write(prose.ExclusionFile, "") }

// overBoundProse is one sentence of one word over the prose bound.
var overBoundProse = strings.Repeat("word ", prose.MaxSentenceWords) + "word."

// boundFinding is the clause of the prose finding for a sentence over the bound.
var boundFinding = fmt.Sprintf("over the %d-word bound", prose.MaxSentenceWords)

// commitOverBoundRecord commits a sentence over the prose bound in the review record of f, and
// answers the record bytes before that sentence.
func commitOverBoundRecord(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	before := outcomeRead(t, filepath.Join(f.Root, fixtureRecord(t)))
	appendRecordProse(t, f, overBoundProse)
	f.Commit("retain record prose over the bound")
	return string(before)
}

// fixtureRecord is the path of the review record of the fixture spec.
func fixtureRecord(t *testing.T) string {
	t.Helper()
	path, err := rr.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return path
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
	if path := fixtureRecord(t); code == 0 || !strings.Contains(out, strconv.Quote(path)) || !strings.Contains(out, boundFinding) {
		t.Fatalf("checkpoint over record prose = %d %s; want a red that names %s and %q", code, out, path, boundFinding)
	}
}
