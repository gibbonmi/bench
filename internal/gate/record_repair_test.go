package gate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestCheckpointEntryRepairForms covers NR17 through NR21: each checkpoint refusal for a
// missing entry routes to the help of the record form that writes that entry.
func TestCheckpointEntryRepairForms(t *testing.T) {
	undispatched := func(t *testing.T) *recordtest.Fixture {
		f := attachedCheckpointFixture(t, recordtest.AttachDelegated, func(f *recordtest.Fixture) {
			recordtest.Undispatched(&f.Plan)
			f.WritePlan()
			f.Commit("empty every assignment history")
			f.Reload()
		})
		// An undispatched ticket has no author, so the chunk holds no result of its own.
		f.Record.Chunks[0].Verification = nil
		f.Save()
		return f
	}
	for _, tc := range []struct {
		name, form string
		complete   bool
		fixture    func(*testing.T) *recordtest.Fixture
	}{
		{"missing chunk", "chunk", false, func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Record.Chunks = nil
			f.Save()
			return f
		}},
		{"missing verification", "verification", false, func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Record.Chunks[0].Verification = nil
			f.Save()
			return f
		}},
		{"missing axis", "review", false, func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
			f.Save()
			return f
		}},
		{"missing completion", "completion", true, checkpointFixture},
		{"undispatched ticket", "assignment", false, undispatched},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := tc.fixture(t)
			args := []string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "1"}
			if tc.complete {
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--complete"}
			}
			var out, diagnostic bytes.Buffer
			code := RunCommand(args, &out, &diagnostic)
			if want := "next=bench record '" + tc.form + "' --help\n"; code == 0 || !strings.Contains(diagnostic.String(), want) {
				t.Fatalf("checkpoint = exit %d, stderr %q; want a refusal whose route is %q", code, diagnostic.String(), want)
			}
		})
	}
}
