package reviewrecord_test

import (
	"errors"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// foldedAfterChunk1 is a linked fixture with count planned chunks whose chunk 1 entry is
// committed, followed by a clean fold of the default branch that the record holds.
func foldedAfterChunk1(t *testing.T, count int) (*recordtest.Fixture, rr.Fold) {
	t.Helper()
	f, _ := recordtest.NewLinked(t, count)
	f.AddChunk()
	f.Save()
	f.Commit("retain chunk 1")
	f.HandFold(f.Mainline(map[string]string{"main.txt": "main\n"}), nil)
	return f, f.RecordCleanFold()
}

// foldRepair requires err to be a refusal that the fold form repairs and that holds want.
func foldRepair(t *testing.T, err error, want string) {
	t.Helper()
	var repair rr.Repair
	if !errors.As(err, &repair) || repair.Form != rr.FormFold || !strings.Contains(err.Error(), want) {
		t.Errorf("refusal = %v, want %q repaired by the fold form", err, want)
	}
}

// TestFoldChainAcceptsARecordedFoldBetweenChunks covers BW34 and BW35: the record commit
// after a recorded fold serves as the chunk 2 base, and without the fold entry that base
// is a chain gap that the fold form repairs.
func TestFoldChainAcceptsARecordedFoldBetweenChunks(t *testing.T) {
	f, _ := foldedAfterChunk1(t, 2)
	f.Commit("retain the fold")
	f.AddChunk()
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), f.Record, "2", false); err != nil {
		t.Errorf("chunk 2 on the record commit after a recorded fold: %v", err)
	}
	unrecorded := f.Record
	unrecorded.Folds = nil
	foldRepair(t, rr.CheckSource(f.Root, f.Tree(), f.Tip(), unrecorded, "2", false), "stale review chain gap")
}

// TestFoldCoversACompletionSourcePastTheLastChunk covers BW37 and BW38: a recorded fold
// covers a completion source past the last chunk tip, and without the fold entry the
// completion refuses the stale source as a repair of the fold form.
func TestFoldCoversACompletionSourcePastTheLastChunk(t *testing.T) {
	f, fold := foldedAfterChunk1(t, 1)
	f.Record.Completion.Verification = f.Verification("final", fold.SourceDigest, f.Plan.FinalVerification)
	f.Save()
	f.Commit("retain final verification")
	recorded, unrecorded := f.Record, f.Record
	unrecorded.Folds = nil
	f.ReplaceRecord(unrecorded)
	_, err := rr.RecordCompletion(f.Root, recordtest.Spec, f.Tip())
	foldRepair(t, err, "stale reviewed source: no chunk review covers")
	f.ReplaceRecord(recorded)
	if _, err := rr.RecordCompletion(f.Root, recordtest.Spec, f.Tip()); err != nil {
		t.Errorf("completion past the last chunk tip under a recorded fold: %v", err)
	}
}
