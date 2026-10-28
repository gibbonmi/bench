package reviewrecord_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// emptySecond empties the requirement list of the second chunk.
func emptySecond(p *rr.Plan) { p.Chunks[1].Verification = []rr.Requirement{} }

// absentKey and absentFinalKey are the encoded keys of a chunk list and of a final list
// that are nil.
const (
	absentKey      = `,"verification":null`
	absentFinalKey = `,"final_verification":null`
)

// missingFinal is the refusal of an empty or absent final list.
const missingFinal = "final verification: missing required verification inventory"

// TestPlanProbeOnlyChunks covers BC23, BC24, BC25, BC26, BC27, BC28, and BC30. Each plan case
// starts from a two-chunk delegated plan whose chunk requirements are all probes, beside the
// plain final list, and the reader grades the committed plan.
func TestPlanProbeOnlyChunks(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		mutate     func(*rr.Plan)
		// drop is the encoded key that the case removes from the plan document, or empty.
		drop string
	}{
		// BC24 and BC28: an empty chunk list passes beside the plain final list.
		{"empty chunk list", "", emptySecond, ""},
		// BC25: the encoded plan holds no verification key for the second chunk.
		{"absent chunk list", "", func(p *rr.Plan) { p.Chunks[1].Verification = nil }, absentKey},
		// BC26: 2.md joins the chunk of 1.md, whose requirements all name 1.md, so 2.md owes none.
		{"ticket that owes nothing", "", mergeTickets, ""},
		// BC27: a chunk requirement still names a ticket of its own chunk.
		{"requirement outside its chunk", `"2.md" outside this chunk`, func(p *rr.Plan) { p.Chunks[0].Verification[0].Ticket = "2.md" }, ""},
		// BC23: a chunk requirement without a probe refuses, and the refusal names the chunk
		// and the requirement.
		{"plain chunk requirement", "chunk 2: verification plain has no probe", func(p *rr.Plan) {
			p.Chunks[1].Verification = append(p.Chunks[1].Verification, rr.Requirement{ID: "plain", Command: "go vet ./...", Ticket: "2.md"})
		}, ""},
		// The final list must still be non-empty, when it is empty and when its key is absent.
		{"empty final list", missingFinal, func(p *rr.Plan) { p.FinalVerification = []rr.Requirement{} }, ""},
		{"absent final list", missingFinal, func(p *rr.Plan) { p.FinalVerification = nil }, absentFinalKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := recordtest.NewDelegated(t, 2)
			tc.mutate(&f.Plan)
			f.WritePlan()
			if tc.drop != "" {
				path := filepath.Join(f.Root, f.Record.Spec)
				data, err := os.ReadFile(path)
				if err != nil || !strings.Contains(string(data), tc.drop) {
					t.Fatalf("plan document %v lacks %s:\n%s", err, tc.drop, data)
				}
				f.Write(f.Record.Spec, strings.Replace(string(data), tc.drop, "", 1))
			}
			f.Commit("stage the plan under test")
			_, err := rr.ReadPlan(f.Root, f.Tree(), f.Record.Spec)
			if tc.want == "" && err != nil {
				t.Fatalf("plan grade refused the probe-only plan: %v", err)
			}
			if tc.want != "" && (err == nil || !strings.Contains(err.Error(), tc.want)) {
				t.Fatalf("plan grade = %v, want a refusal naming %q", err, tc.want)
			}
		})
	}
	// BC30: the checkpoint still grades each planned probe, so a missing probe entry refuses
	// through the verification form.
	t.Run("missing probe entry", func(t *testing.T) {
		f := recordtest.Prepare(t, gittest.RepoOnBranch(t, "main"), 2, recordtest.Spec, recordtest.Body, recordtest.Delegate, emptySecond)
		f.AddChunk()
		f.AddChunk()
		f.Complete()
		f.Record.Chunks[0].Verification = nil
		f.Save()
		f.Commit("retain the chunks without the probe entry")
		var repair rr.Repair
		err := accept(f)
		if !errors.As(err, &repair) || repair.Form != rr.FormVerification || repair.Chunk != "1" || repair.Requirement != "tests" {
			t.Fatalf("checkpoint = %v (%+v), want the verification repair of chunk 1 requirement tests", err, repair)
		}
	})
}
