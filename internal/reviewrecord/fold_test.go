package reviewrecord_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// foldedAfterChunk1 is a linked fixture with count planned chunks whose chunk 1 entry is
// committed. A fold of the default branch follows, and the record holds it with delta.
func foldedAfterChunk1(t *testing.T, count int, delta rr.FoldDelta) (*recordtest.Fixture, rr.Fold) {
	t.Helper()
	f, _ := recordtest.NewLinked(t, count)
	f.AddChunk()
	f.Save()
	f.Commit("retain chunk 1")
	f.HandFold(f.Mainline(map[string]string{"main.txt": "main\n"}), nil)
	return f, f.RecordFold(delta)
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
// after one or two recorded folds serves as the chunk 2 base. Without the fold entries, or
// off the fold commit, that base is a chain gap that the fold form repairs.
func TestFoldChainAcceptsARecordedFoldBetweenChunks(t *testing.T) {
	for _, folds := range []int{1, 2} {
		t.Run(fmt.Sprintf("%d folds", folds), func(t *testing.T) {
			f, fold := foldedAfterChunk1(t, 2, rr.FoldDelta{})
			for n := 2; n <= folds; n++ {
				f.Commit("retain fold " + fold.ID)
				f.HandFold(f.Mainline(map[string]string{fmt.Sprintf("main%d.txt", n): "main\n"}), nil)
				fold = f.RecordFold(rr.FoldDelta{})
			}
			f.Commit("retain the folds")
			f.AddChunk()
			if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), f.Record, "2", false); err != nil {
				t.Errorf("chunk 2 on the record commit after %d recorded folds: %v", folds, err)
			}
			unrecorded := f.Record
			unrecorded.Folds = nil
			foldRepair(t, rr.CheckSource(f.Root, f.Tree(), f.Tip(), unrecorded, "2", false), "stale review chain gap")

			// The same base tree and chunk 2 tree on a parent chain that skips the fold commit.
			offFold := f.Record
			offFold.Chunks = slices.Clone(f.Record.Chunks)
			chunk := &offFold.Chunks[1]
			chunk.Base = f.Git("commit-tree", chunk.Base+"^{tree}", "-p", offFold.Chunks[0].Tip, "-m", "fold digest off the fold")
			chunk.Tip = f.Git("commit-tree", chunk.Tip+"^{tree}", "-p", chunk.Base, "-m", "chunk 2 off the fold")
			chunk.Reviews = slices.Clone(chunk.Reviews)
			for i := range chunk.Reviews {
				chunk.Reviews[i].Base, chunk.Reviews[i].Tip = chunk.Base, chunk.Tip
			}
			foldRepair(t, rr.CheckSource(f.Root, f.Tree(), chunk.Tip, offFold, "2", false), "stale review chain gap")
		})
	}
}

// TestChainRefusesAFailedAncestryQuery proves that an ancestry query Git cannot run
// refuses and names the query. The commit is absent, so merge-base exits 128, and no
// record form repairs it.
func TestChainRefusesAFailedAncestryQuery(t *testing.T) {
	const absent = "0123456789abcdef0123456789abcdef01234567"
	for _, tc := range []struct {
		name    string
		corrupt func(record *rr.Record)
	}{
		{"chunk base", func(record *rr.Record) { record.Chunks[1].Base = absent }},
		{"fold commit", func(record *rr.Record) { record.Folds[0].Commit = absent }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, _ := foldedAfterChunk1(t, 2, rr.FoldDelta{})
			f.Commit("retain the fold")
			f.AddChunk()
			broken := f.Record
			broken.Chunks = slices.Clone(f.Record.Chunks)
			broken.Folds = slices.Clone(f.Record.Folds)
			tc.corrupt(&broken)
			err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), broken, "2", false)
			var repair rr.Repair
			if err == nil || errors.As(err, &repair) || !strings.Contains(err.Error(), "ancestry query of "+absent) {
				t.Errorf("chunk 2 over an absent %s = %v, want a failed ancestry query with no repair form", tc.name, err)
			}
		})
	}
}

// TestFoldCoversACompletionSourcePastTheLastChunk covers BW37 and BW38: a recorded fold
// covers a completion source past the last chunk tip, and without the fold entry the
// completion refuses the stale source as a repair of the fold form.
func TestFoldCoversACompletionSourcePastTheLastChunk(t *testing.T) {
	f, fold := foldedAfterChunk1(t, 1, rr.FoldDelta{})
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

// TestFoldResolutionNeedsReview covers BW39 and BW40: a clean fold with no review passes the
// checkpoint. A fold with a resolution refuses the checkpoint until its third axis passes.
func TestFoldResolutionNeedsReview(t *testing.T) {
	clean, _ := foldedAfterChunk1(t, 1, rr.FoldDelta{})
	if err := rr.CheckSource(clean.Root, clean.Tree(), clean.Tip(), clean.Record, "1", false); err != nil {
		t.Errorf("checkpoint over a clean fold with no review: %v", err)
	}
	f, fold := foldedAfterChunk1(t, 1, rr.FoldDelta{Resolution: []string{"source.txt"}})
	axes := rr.Axes()
	for i, axis := range axes {
		call := rr.ReviewCall{Fold: fold.ID, Axis: axis, Evidence: f.Evidence("fold-"+axis, "", "independent-review")}
		if _, err := rr.RecordReview(f.Root, recordtest.Spec, call); err != nil {
			t.Fatal(err)
		}
		record, err := rr.Read(f.Root, recordtest.Spec)
		if err != nil {
			t.Fatal(err)
		}
		err = rr.CheckSource(f.Root, f.Tree(), f.Tip(), record, "1", false)
		var repair rr.Repair
		switch {
		case i == len(axes)-1 && err != nil:
			t.Errorf("checkpoint after every fold axis passed: %v", err)
		case i < len(axes)-1 && (!errors.As(err, &repair) || repair.Form != rr.FormReview || repair.Fold != fold.ID || repair.Axis != axes[i+1]):
			t.Errorf("checkpoint after %d fold axes = %v, want the review form for fold %s axis %s", i+1, err, fold.ID, axes[i+1])
		}
	}
}

// TestFoldDerivedFinalRequirements covers BW42: the complete checkpoint refuses until a
// passing final result exists for the merged-test requirement of a fold package.
func TestFoldDerivedFinalRequirements(t *testing.T) {
	f, fold := foldedAfterChunk1(t, 1, rr.FoldDelta{TestPackages: []string{"./pkg/a"}})
	f.Record.Completion.Verification = f.Verification("final", fold.SourceDigest, f.Plan.FinalVerification)
	f.Save()
	f.Commit("retain final verification")
	merged := rr.MergedTestRequirement(fold.TestPackages[0])
	_, err := rr.RecordCompletion(f.Root, recordtest.Spec, f.Tip())
	var repair rr.Repair
	if !errors.As(err, &repair) || repair.Form != rr.FormVerification || repair.Requirement != merged.ID || !strings.Contains(err.Error(), merged.Command) {
		t.Fatalf("completion without a merged-test result = %v, want the verification form for %s", err, merged.ID)
	}
	call := rr.VerificationCall{Source: f.Tip(), Requirement: merged.ID, Evidence: f.Evidence("final-"+merged.ID, "", "author-verification")}
	if _, err := rr.RecordVerification(f.Root, recordtest.Spec, call); err != nil {
		t.Fatalf("final verification of %s: %v", merged.ID, err)
	}
	if _, err := rr.RecordCompletion(f.Root, recordtest.Spec, f.Tip()); err != nil {
		t.Errorf("completion with a passing merged-test result: %v", err)
	}
}
