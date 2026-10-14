package gate

import (
	"slices"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// repairTest is the Go test file that a test-only repair changes, and repairSource is the
// production Go file beside it.
const (
	repairTest   = "repair_test.go"
	repairSource = "repair.go"
)

func seedRepair(f *recordtest.Fixture) {
	f.Write(repairTest, "package fixture\n\nconst want = 1\n")
	f.Write(repairSource, "package fixture\n\nconst value = 1\n")
}

// Each repair below writes one repair delta after the reviewed tip.
func modifyTest(f *recordtest.Fixture) {
	f.Write(repairTest, "package fixture\n\nconst want = 2\n")
}

func addTest(f *recordtest.Fixture) {
	f.Write("added_test.go", "package fixture\n\nconst added = 1\n")
}

func modifySource(f *recordtest.Fixture) {
	f.Write(repairSource, "package fixture\n\nconst value = 2\n")
}

func addPlanFile(f *recordtest.Fixture) {
	f.Write("specs/example/notes.md", "# Notes\n")
}

// repairedFixture is one delegated chunk whose Standards result fails with F-1 and whose Spec
// and Coverage results pass at the reviewed pair. repair changes the chunk, and the chunk
// refreezes at the repair with no confirmation. It answers the fixture and the reviewed tip.
func repairedFixture(t *testing.T, repair func(*recordtest.Fixture)) (*recordtest.Fixture, string) {
	t.Helper()
	f := attachedCheckpointFixture(t, recordtest.AttachDelegated, seedRepair)
	reviewed := f.Record.Chunks[0].Tip
	standards := &f.Record.Chunks[0].Reviews[0]
	standards.Outcome, standards.FindingIDs = "fail", []string{"F-1"}
	repair(f)
	f.Commit("repair the chunk")
	f.Refreeze()
	f.Save()
	f.Commit("retain the refrozen chunk")
	return f, reviewed
}

// confirmedFixture is the repaired fixture with a confirmation that closes F-1.
func confirmedFixture(t *testing.T, repair func(*recordtest.Fixture)) (*recordtest.Fixture, string) {
	t.Helper()
	f, reviewed := repairedFixture(t, repair)
	f.Confirm("confirm-1", reviewed, "F-1")
	f.Save()
	f.Commit("retain the confirmation")
	return f, reviewed
}

// TestReviewCheckpointPassesAWrittenConfirmation grades the confirmation that the record
// writer derives for the repaired fixture: the chunk checkpoint passes with it.
func TestReviewCheckpointPassesAWrittenConfirmation(t *testing.T) {
	f, _ := repairedFixture(t, modifyTest)
	call := rr.ConfirmationCall{Chunk: "1", Findings: []string{"F-1"},
		Evidence: rr.Evidence{ID: "confirm-1", Performer: recordtest.Orchestrator, Model: "unknown", Effort: "unknown", NativeRef: recordtest.Native("mutation failed; restore passed")},
		Probe:    recordtest.ConfirmationProbe()}
	if _, _, err := rr.RecordConfirmation(f.Root, recordtest.Spec, call); err != nil {
		t.Fatal(err)
	}
	f.Commit("retain the written confirmation")
	if code, out := runCheckpoint(t, f); code != 0 {
		t.Fatalf("checkpoint with the written confirmation = %d %s", code, out)
	}
	if got := outcomeRuns(t, f.Root); got != 1 {
		t.Fatalf("oracle runs = %d, want 1", got)
	}
}

// repairAgain commits a second test-only repair and refreezes the chunk at it, with no new
// confirmation.
func repairAgain(f *recordtest.Fixture) {
	f.Write(repairTest, "package fixture\n\nconst want = 3\n")
	f.Commit("repair the chunk again")
	f.Refreeze()
}

// freshReview is a passing result of axis at the chunk pair that supersedes the result at the
// reviewed pair. The fixture records one result of each axis in axis order.
func freshReview(f *recordtest.Fixture, axis string) rr.Review {
	chunk := f.Record.Chunks[0]
	prior := chunk.Reviews[slices.Index(rr.Axes(), axis)].ID
	return rr.Review{Evidence: f.Evidence("fresh-"+axis, chunk.SourceDigest, "independent-review"), Axis: axis, Base: chunk.Base, Tip: chunk.Tip, Supersedes: []string{prior}}
}

func TestReviewCheckpointCarriesConfirmedResults(t *testing.T) {
	for _, tc := range []struct {
		name     string
		repair   func(*recordtest.Fixture)
		after    func(f *recordtest.Fixture, reviewed string)
		complete bool
	}{
		{name: "modified test", repair: modifyTest},
		{name: "complete", repair: modifyTest, after: func(f *recordtest.Fixture, _ string) { retainCompletion(f) }, complete: true},
		{name: "added test", repair: addTest},
		{name: "second confirmation", repair: modifyTest, after: func(f *recordtest.Fixture, reviewed string) {
			repairAgain(f)
			f.Confirm("confirm-2", reviewed, "F-1")
			f.Save()
			f.Commit("retain the second confirmation")
		}},
		{name: "fresh axes after a production repair", repair: modifyTest, after: func(f *recordtest.Fixture, _ string) {
			f.Write("source.txt", "production repair\n")
			f.Commit("repair the production source")
			f.Refreeze()
			for _, axis := range rr.Axes() {
				f.Record.Chunks[0].Reviews = append(f.Record.Chunks[0].Reviews, freshReview(f, axis))
			}
			f.Save()
			f.Commit("retain the fresh reviews")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, reviewed := confirmedFixture(t, tc.repair)
			if tc.after != nil {
				tc.after(f, reviewed)
			}
			checkpoint := runCheckpoint
			if tc.complete {
				checkpoint = func(t *testing.T, f *recordtest.Fixture) (int, string) { return completeCheckpoint(t, f) }
			}
			if code, out := checkpoint(t, f); code != 0 {
				t.Fatalf("confirmed checkpoint = %d %s", code, out)
			}
			if got := outcomeRuns(t, f.Root); got != 1 {
				t.Fatalf("oracle runs = %d, want 1", got)
			}
		})
	}
}

// TestReviewCheckpointConfirmationRefusals breaks one condition of the confirmed fixture in
// each case. The reason names the chunk and the condition, the route names the record form,
// and the oracle does not run. In a reason, <base> and <reviewed> stand for the chunk base and
// the reviewed tip, and <moved> stands for the value that the change answers.
func TestReviewCheckpointConfirmationRefusals(t *testing.T) {
	latest := func(f *recordtest.Fixture) *rr.Confirmation {
		confirmations := f.Record.Chunks[0].Confirmations
		return &confirmations[len(confirmations)-1]
	}
	// moveReviewed moves the three axis results and the confirmation to the reviewed tip tip.
	moveReviewed := func(f *recordtest.Fixture, tip string) string {
		for i := range f.Record.Chunks[0].Reviews {
			f.Record.Chunks[0].Reviews[i].Tip = tip
		}
		latest(f).ReviewedTip = tip
		return tip
	}
	for _, tc := range []struct {
		name, reason, form string
		repair             func(*recordtest.Fixture)
		change             func(f *recordtest.Fixture, reviewed string) string
	}{
		{"non-test Go file", "chunk 1: stale axis results at the reviewed tip <reviewed>: the delta to the chunk tip is not test-only: repair.go: path is not a Go test file", "review", modifySource, nil},
		{"plan file", "the delta to the chunk tip is not test-only: specs/example/notes.md: path is not a Go test file", "review", addPlanFile, nil},
		{"later production repair", "the delta to the chunk tip is not test-only: source.txt: path is not a Go test file", "review", modifyTest, func(f *recordtest.Fixture, reviewed string) string {
			f.Write("source.txt", "production repair\n")
			f.Commit("repair the production source")
			f.Refreeze()
			f.Confirm("confirm-2", reviewed, "F-1")
			return ""
		}},
		{"silent probe", "chunk 1: confirmation confirm-1 probe did not bite: outcome silent", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).Probe.Outcome = "silent"
			return ""
		}},
		{"failed restore", "chunk 1: confirmation confirm-1 probe restore fail", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).Probe.Restore = "fail"
			return ""
		}},
		{"zero focused exit", "chunk 1: confirmation confirm-1 probe did not bite: outcome bit, focused exit 0", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).Probe.ExitCode = 0
			return ""
		}},
		{"author performer", "chunk 1: confirmation confirm-1 performer <moved> is not the completion reconciler " + recordtest.Orchestrator, "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).Performer = recordtest.Author(f.Plan.Chunks[0].Tickets[0])
			return latest(f).Performer
		}},
		{"omitted finding", "chunk 1: confirmation confirm-1 closes the findings [], and the open findings at the reviewed pair are [F-1]", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).FindingIDs = []string{}
			return ""
		}},
		{"extra finding", "chunk 1: confirmation confirm-1 closes the findings [F-1 F-2], and the open findings at the reviewed pair are [F-1]", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).FindingIDs = []string{"F-1", "F-2"}
			return ""
		}},
		{"side reviewed tip", "chunk 1: reviewed tip <moved> is not an ancestor of the chunk tip", "review", modifyTest, func(f *recordtest.Fixture, reviewed string) string {
			return moveReviewed(f, f.Git("-c", "user.name=bench", "-c", "user.email=bench@local", "commit-tree", reviewed+"^{tree}", "-p", reviewed, "-m", "side"))
		}},
		{"reviewed tip before base", "chunk 1: reviewed tip <moved> does not descend from the chunk base <base>", "review", modifyTest, func(f *recordtest.Fixture, _ string) string {
			return moveReviewed(f, f.Git("rev-parse", f.Record.Chunks[0].Base+"^"))
		}},
		{"empty reviewed pair", "chunk 1: empty reviewed pair <reviewed>..<reviewed>", "review", modifyTest, func(f *recordtest.Fixture, reviewed string) string {
			chunk := &f.Record.Chunks[0]
			chunk.Base, latest(f).Base = reviewed, reviewed
			for i := range chunk.Reviews {
				chunk.Reviews[i].Base = reviewed
			}
			return ""
		}},
		{"author carried result", `chunk 1: invalid Standards performer "<moved>"`, "review", modifyTest, func(f *recordtest.Fixture, _ string) string {
			f.Record.Chunks[0].Reviews[0].Performer = recordtest.Author(f.Plan.Chunks[0].Tickets[0])
			return f.Record.Chunks[0].Reviews[0].Performer
		}},
		{"forged carried digest", "chunk 1: stale Spec; review the current frozen pair", "review", modifyTest, func(f *recordtest.Fixture, _ string) string {
			f.Record.Chunks[0].Reviews[1].SourceDigest = f.Record.Chunks[0].SourceDigest
			return ""
		}},
		{"stale repaired tip", "chunk 1: confirmation confirm-1 covers the stale pair <base>..<moved>", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			repaired := f.Record.Chunks[0].Tip
			repairAgain(f)
			return repaired
		}},
		{"other base", "chunk 1: confirmation confirm-1 covers the stale pair <reviewed>..<moved>; confirm the chunk pair <base>..<moved>", "confirmation", modifyTest, func(f *recordtest.Fixture, reviewed string) string {
			latest(f).Base = reviewed
			return f.Record.Chunks[0].Tip
		}},
		{"failed carried result", "chunk 1: failed Standards", "review", modifyTest, func(f *recordtest.Fixture, _ string) string {
			f.Record.Chunks[0].Reviews[0].State = "failed"
			return ""
		}},
		{"pending confirmation", "chunk 1: confirmation confirm-1 state is pending", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).State = "pending"
			return ""
		}},
		{"fresh Coverage", "chunk 1: mixed review pairs: Standards at <base>..<reviewed>, Spec at <base>..<reviewed>, Coverage at <base>..<moved>", "review", modifyTest, func(f *recordtest.Fixture, _ string) string {
			f.Record.Chunks[0].Reviews = append(f.Record.Chunks[0].Reviews, freshReview(f, "Coverage"))
			return f.Record.Chunks[0].Tip
		}},
		{"missing axis", "chunk 1: missing Coverage; a reviewed pair needs a current result of each axis", "review", modifyTest, func(f *recordtest.Fixture, _ string) string {
			f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
			return ""
		}},
		{"other reviewed tip", "chunk 1: confirmation confirm-1 names the reviewed tip <base>, and the axis results sit at <reviewed>", "confirmation", modifyTest, func(f *recordtest.Fixture, _ string) string {
			latest(f).ReviewedTip = f.Record.Chunks[0].Base
			return ""
		}},
		{"verification at the reviewed digest", "chunk 1: verification tests is incomplete, failed, or stale", "verification", modifyTest, func(f *recordtest.Fixture, _ string) string {
			chunk := &f.Record.Chunks[0]
			for i := range chunk.Verification {
				chunk.Verification[i].SourceDigest = chunk.Reviews[0].SourceDigest
			}
			return ""
		}},
		{"one carried session", "chunk 1: use three distinct review sessions", "review", modifyTest, func(f *recordtest.Fixture, _ string) string {
			for i := range f.Record.Chunks[0].Reviews {
				f.Record.Chunks[0].Reviews[i].Performer = "fixture-reviewer-shared"
			}
			return ""
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, reviewed := confirmedFixture(t, tc.repair)
			moved := ""
			if tc.change != nil {
				moved = tc.change(f, reviewed)
				f.Save()
			}
			reason := strings.NewReplacer("<base>", f.Record.Chunks[0].Base, "<reviewed>", reviewed, "<moved>", moved).Replace(tc.reason)
			code, out := runCheckpoint(t, f)
			if next := "next=bench record '" + tc.form + "' --help\n"; code == 0 || !strings.Contains(out, reason) || !strings.Contains(out, next) {
				t.Fatalf("confirmation refusal = %d %s; want the reason %q and the route %q", code, out, reason, next)
			}
			assertNoOracleRun(t, f)
		})
	}
}
