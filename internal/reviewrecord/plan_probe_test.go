package reviewrecord_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// probeOnly keeps only the probe requirements of each chunk of a delegated plan. The final
// list keeps its plain requirements.
func probeOnly(p *rr.Plan) {
	for i := range p.Chunks {
		p.Chunks[i].Verification = slices.DeleteFunc(p.Chunks[i].Verification, func(item rr.Requirement) bool { return item.Probe == "" })
	}
}

// emptySecond empties the requirement list of the second chunk.
func emptySecond(p *rr.Plan) { p.Chunks[1].Verification = []rr.Requirement{} }

// absentKey is the encoded verification key of a chunk whose list is nil.
const absentKey = `,"verification":null`

// TestPlanProbeOnlyChunks covers BC24, BC25, BC26, BC27, BC28, and BC30. Each plan case starts
// from a two-chunk delegated plan with only its probe requirements, and the reader grades the
// committed plan.
func TestPlanProbeOnlyChunks(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		mutate     func(*rr.Plan)
		absent     bool
	}{
		// BC24 and BC28: an empty chunk list passes beside the plain final list.
		{"empty chunk list", "", emptySecond, false},
		// BC25: the encoded plan holds no verification key for the second chunk.
		{"absent chunk list", "", func(p *rr.Plan) { p.Chunks[1].Verification = nil }, true},
		// BC26: 2.md joins the chunk of 1.md, whose requirements all name 1.md, so 2.md owes none.
		{"ticket that owes nothing", "", mergeTickets, false},
		// BC27: a chunk requirement still names a ticket of its own chunk.
		{"requirement outside its chunk", `"2.md" outside this chunk`, func(p *rr.Plan) { p.Chunks[0].Verification[0].Ticket = "2.md" }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := recordtest.NewDelegated(t, 2)
			tc.mutate(&f.Plan)
			probeOnly(&f.Plan)
			f.WritePlan()
			if tc.absent {
				path := filepath.Join(f.Root, f.Record.Spec)
				data, err := os.ReadFile(path)
				if err != nil || !strings.Contains(string(data), absentKey) {
					t.Fatalf("plan document %v lacks %s:\n%s", err, absentKey, data)
				}
				f.Write(f.Record.Spec, strings.Replace(string(data), absentKey, "", 1))
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
		f := recordtest.Prepare(t, gittest.RepoOnBranch(t, "main"), 2, recordtest.Spec, recordtest.Body, recordtest.Delegate, probeOnly, emptySecond)
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
