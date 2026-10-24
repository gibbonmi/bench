package gate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
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

// TestCheckpointPrepareRunsTheGateWithoutEvidence covers RC34. The record commit of a frozen
// chunk holds no review and a sentence over the prose bound, so the chunk checkpoint refuses
// it. The prepared run grades its subject alone: the oracle runs once, and the evidence store
// holds a green for the HEAD tree without the record.
func TestCheckpointPrepareRunsTheGateWithoutEvidence(t *testing.T) {
	f, _ := unreviewedFixture(t, prosePolicy)
	commitOverBoundRecord(t, f)
	requirePrepared(t, f)
	stripped, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if got := InspectTree(f.Root, stripped); !got.ReusableGreen {
		t.Fatalf("evidence of the HEAD tree without the record = %+v; want the prepared green", got)
	}
}

// TestChunkCheckpointReusesThePreparedGreen covers RC36. The record commit that retains the
// reviews follows the prepared run, and a second call reloads the evidence store, so the
// chunk checkpoint reuses the prepared green and runs no phase.
func TestChunkCheckpointReusesThePreparedGreen(t *testing.T) {
	f, reviews := unreviewedFixture(t)
	requirePrepared(t, f)
	f.Record.Chunks[0].Reviews = reviews
	resave(t, f, "retain the reviews")
	requireReuse(t, f, "after the prepared run and a record-only commit")
}

// TestChunkCheckpointRunsAfterASpecChange covers RC37: the spec stays in the key.
func TestChunkCheckpointRunsAfterASpecChange(t *testing.T) {
	requireRunAfterAChange(t, recordtest.Spec)
}

// TestChunkCheckpointRunsAfterASourceChange covers RC38: a source path other than the record
// stays in the key.
func TestChunkCheckpointRunsAfterASourceChange(t *testing.T) {
	requireRunAfterAChange(t, "source.txt")
}

// TestCheckpointPrepareSurvivesARecordCommit covers RC42. After the prepared run takes the
// gate lock, the worktree retains the reviews in a record write and a record commit. The run
// resolves no tip and captures its private checkout, so it stays green, and the chunk
// checkpoint after it reuses that green.
func TestCheckpointPrepareSurvivesARecordCommit(t *testing.T) {
	f, reviews := unreviewedFixture(t)
	retain := func() {
		f.Record.Chunks[0].Reviews = reviews
		resave(t, f, "retain the reviews during the prepared run")
	}
	var stderr bytes.Buffer
	if code := RunCheckpointArmed(t, f.Root, []string{"--checkpoint", recordtest.Spec, "--chunk", "1", flagPrepare}, retain, &stderr); code != 0 {
		t.Fatalf("prepared run across a record commit = %d %s; want its green", code, stderr.String())
	}
	requireReuse(t, f, "after a record commit during the prepared run")
}

// TestGateUsageNamesPrepare covers RC49: `bench gate --help` names the prepare flag.
func TestGateUsageNamesPrepare(t *testing.T) {
	var stdout bytes.Buffer
	if code := Command([]string{"--help"}, nil, &stdout, io.Discard); code != 0 || !strings.Contains(stdout.String(), flagPrepare) {
		t.Fatalf("gate --help = %d %q; want a usage line that names %s", code, stdout.String(), flagPrepare)
	}
}

// TestChunkCheckpointRunsAfterThePreparedGreenExpires covers RC52. A prepared green past the
// freshness window serves no chunk checkpoint, so the checkpoint runs the gate again.
func TestChunkCheckpointRunsAfterThePreparedGreenExpires(t *testing.T) {
	f := checkpointFixture(t)
	requirePrepared(t, f)
	expireEvidence(t, f.Root)
	requireRun(t, f, "after the prepared green expires")
}

// TestCheckpointPrepareRefusesADirtyCheckout grades a prepared run over a dirty checkout. The
// run grades the HEAD tree that the chunk checkpoint grades, so it refuses the checkout before
// the oracle, and the rerun of its route prepares again.
func TestCheckpointPrepareRefusesADirtyCheckout(t *testing.T) {
	f := checkpointFixture(t)
	f.Write("tracked.txt", "edited\n")
	code, out := runCheckpointWith(t, f, flagPrepare)
	requireCleanCheckoutRefusal(t, f, code, out, "stale reviewed source")
	if next, one := routetest.Next(out); !one || !strings.HasSuffix(next, " "+flagPrepare) {
		t.Errorf("dirty prepared run = %q; want a next= rerun that keeps %s", out, flagPrepare)
	}
}

// unreviewedFixture is a checkpoint fixture whose last record commit holds the frozen chunk
// with no review, as at a freeze before the review axes run. It answers the reviews that the
// record leaves out.
func unreviewedFixture(t *testing.T, prepare ...func(*recordtest.Fixture)) (*recordtest.Fixture, []rr.Review) {
	t.Helper()
	f := attachedCheckpointFixture(t, recordtest.Attach, prepare...)
	reviews := f.Record.Chunks[0].Reviews
	f.Record.Chunks[0].Reviews = []rr.Review{}
	resave(t, f, "retain the frozen chunk before its reviews")
	return f, reviews
}

// requireRunAfterAChange prepares chunk 1 of a checkpoint fixture over a commit that changes
// path. It then commits the reviewed bytes of path again, so the evidence passes. The chunk
// checkpoint must then run the gate on its own subject, not reuse the prepared green.
func requireRunAfterAChange(t *testing.T, path string) {
	t.Helper()
	f := checkpointFixture(t)
	reviewed := string(outcomeRead(t, filepath.Join(f.Root, path)))
	f.Write(path, reviewed+"\nA change after the review.\n")
	f.Commit("change " + path)
	requirePrepared(t, f)
	f.Write(path, reviewed)
	f.Commit("restore " + path)
	requireRun(t, f, "after a change of "+path)
}

// requireRun requires that a chunk checkpoint of f exits 0 after a second oracle run.
func requireRun(t *testing.T, f *recordtest.Fixture, when string) {
	t.Helper()
	if code, out := runCheckpoint(t, f); code != 0 || outcomeRuns(t, f.Root) != 2 {
		t.Fatalf("checkpoint %s = %d %s after %d runs; want a second run", when, code, out, outcomeRuns(t, f.Root))
	}
}

// recordedAt matches the recorded time of one evidence-store record.
var recordedAt = regexp.MustCompile(`"recorded_at":"[^"]*"`)

// expireEvidence moves the recorded time of each timed record in the evidence store of root
// back by the freshness window, so no retained green is fresh.
func expireEvidence(t *testing.T, root string) {
	t.Helper()
	common, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	store := filepath.Dir(evidencePath(common, subject{}))
	entries, err := os.ReadDir(store)
	if err != nil {
		t.Fatal(err)
	}
	expired := []byte(`"recorded_at":"` + time.Now().UTC().Add(-freshness).Truncate(time.Second).Format(time.RFC3339) + `"`)
	aged := 0
	for _, entry := range entries {
		path := filepath.Join(store, entry.Name())
		if data := outcomeRead(t, path); recordedAt.Match(data) {
			if err := os.WriteFile(path, recordedAt.ReplaceAll(data, expired), 0o600); err != nil {
				t.Fatal(err)
			}
			aged++
		}
	}
	if aged == 0 {
		t.Fatalf("evidence store %s holds no timed record to expire", store)
	}
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
