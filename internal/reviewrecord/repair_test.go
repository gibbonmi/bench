package reviewrecord_test

import (
	"errors"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// repairOf is the repair that err carries, and it fails the test when err carries none.
func repairOf(t *testing.T, err error) rr.Repair {
	t.Helper()
	var repair rr.Repair
	if !errors.As(err, &repair) {
		t.Fatalf("error %v carries no repair", err)
	}
	return repair
}

// repairValues is the form and the known values of repair, without its error.
func repairValues(repair rr.Repair) [6]string {
	return [6]string{repair.Form, repair.Chunk, repair.Requirement, repair.Ticket, repair.Axis, repair.Source}
}

// verifyChunk asks for one chunk verification result of requirement tests.
func verifyChunk(f *recordtest.Fixture, chunk string) error {
	_, err := rr.RecordVerification(f.Root, f.Record.Spec, rr.VerificationCall{Chunk: chunk, Requirement: "tests",
		Evidence: rr.Evidence{ID: "tests-" + chunk, Performer: "fixture-author", Model: "unknown", Effort: "unknown", NativeRef: recordtest.Native("ok")}})
	return err
}

// TestRecordRepairPreservesErrorIdentity covers NR33: a missing record and a missing chunk
// entry name the chunk form, and the existing sentinel stays reachable through the repair.
func TestRecordRepairPreservesErrorIdentity(t *testing.T) {
	t.Run("missing record", func(t *testing.T) {
		err := verifyChunk(recordtest.New(t, 1), "1")
		if repair := repairOf(t, err); !errors.Is(err, rr.ErrNoChunkEntry) || repair.Form != rr.FormChunk {
			t.Fatalf("missing record = %v with repair %+v, want ErrNoChunkEntry and the chunk form", err, repair)
		}
	})
	t.Run("missing chunk", func(t *testing.T) {
		f := recordtest.New(t, 2)
		f.AddChunk()
		f.Save()
		err := verifyChunk(f, "2")
		if repair := repairOf(t, err); !errors.Is(err, rr.ErrNoChunkEntry) || repair.Form != rr.FormChunk || repair.Chunk != "2" {
			t.Fatalf("missing chunk = %v with repair %+v, want ErrNoChunkEntry and the chunk form for chunk 2", err, repair)
		}
	})
}

// TestRecordRepairNamesOwnerValues grades the values that each refusing owner supplies:
// record.go names the chunk and the missing axis, the verifier names the undispatched
// requirement and its ticket, the checkpoint names the requested source, and a result
// writer names the chunk of its call when the record is absent.
func TestRecordRepairNamesOwnerValues(t *testing.T) {
	recorded := func(t *testing.T, options ...recordtest.Option) *recordtest.Fixture {
		f := recordtest.Prepare(t, gittest.RepoOnBranch(t, "main"), 1, recordtest.Spec, fixtureBody, options...)
		f.AddChunk()
		return f
	}
	check := func(f *recordtest.Fixture, chunk string) error {
		f.Save()
		f.Commit("retain the evidence under test")
		return rr.Check(f.Root, f.Tree(), f.Tip(), f.Record.Spec, chunk, chunk == "")
	}
	for _, tc := range []struct {
		name string
		err  func(t *testing.T) (error, string)
		want rr.Repair
	}{
		{name: "missing axis", want: rr.Repair{Form: rr.FormReview, Chunk: "1", Axis: "Coverage"}, err: func(t *testing.T) (error, string) {
			chunk := recorded(t).Record.Chunks[0]
			chunk.Reviews = chunk.Reviews[:2]
			return rr.CheckReviews(chunk, nil, false), ""
		}},
		{name: "missing verification", want: rr.Repair{Form: rr.FormVerification, Chunk: "1", Requirement: "tests"}, err: func(t *testing.T) (error, string) {
			f := recorded(t)
			f.Record.Chunks[0].Verification = nil
			return check(f, "1"), ""
		}},
		{name: "undispatched ticket", want: rr.Repair{Form: rr.FormAssignment, Requirement: "tests", Ticket: "1.md"}, err: func(t *testing.T) (error, string) {
			f := recorded(t, recordtest.Delegate, recordtest.Undispatched)
			// An undispatched ticket has no author, so the record holds no result of its own.
			f.Record.Chunks[0].Verification = nil
			f.Save()
			return verifyChunk(f, "1"), ""
		}},
		{name: "missing completion", want: rr.Repair{Form: rr.FormCompletion}, err: func(t *testing.T) (error, string) {
			f := recorded(t)
			err := check(f, "")
			return err, f.Tip()
		}},
		{name: "changed plan", want: rr.Repair{Form: rr.FormAmendment}, err: func(t *testing.T) (error, string) {
			f := recorded(t)
			f.Save()
			f.Plan.FinalVerification[0].Command = "go test -count=1 ./..."
			f.WritePlan()
			f.Commit("change the plan")
			return rr.Check(f.Root, f.Tree(), f.Tip(), f.Record.Spec, "1", false), f.Tip()
		}},
		{name: "absent record verification", want: rr.Repair{Form: rr.FormChunk, Chunk: "1"}, err: func(t *testing.T) (error, string) {
			return verifyChunk(recordtest.New(t, 1), "1"), ""
		}},
		{name: "absent record review", want: rr.Repair{Form: rr.FormChunk, Chunk: "1"}, err: func(t *testing.T) (error, string) {
			f := recordtest.New(t, 1)
			_, err := rr.RecordReview(f.Root, f.Record.Spec, rr.ReviewCall{Chunk: "1", Axis: "Standards",
				Evidence: rr.Evidence{ID: "standards-1", Performer: "fixture-reviewer", Model: "unknown", Effort: "unknown", NativeRef: recordtest.Native("ok")}})
			return err, ""
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err, source := tc.err(t)
			got := repairOf(t, err)
			tc.want.Source = source
			if repairValues(got) != repairValues(tc.want) {
				t.Fatalf("repair form, chunk, requirement, ticket, axis, and source = %q, want %q", repairValues(got), repairValues(tc.want))
			}
		})
	}
}
