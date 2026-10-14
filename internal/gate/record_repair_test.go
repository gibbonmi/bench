package gate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestCheckpointEntryRepairForms covers NR17 through NR21: each checkpoint refusal for a
// missing entry routes to the help of the record form that writes that entry. A changed
// plan routes to the amendment form, and an unplanned chunk, which no form writes, keeps
// the preflight review route. RC43 and RC44: axis results at an earlier pair with a failed
// Standards result route to the confirmation form after a test-only repair, and to the review
// form after a repair of a non-test Go file.
func TestCheckpointEntryRepairForms(t *testing.T) {
	repaired := func(repair func(*recordtest.Fixture)) func(*testing.T) *recordtest.Fixture {
		return func(t *testing.T) *recordtest.Fixture {
			f, _ := repairedFixture(t, repair)
			return f
		}
	}
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
	// An empty chunk selects the complete checkpoint.
	for _, tc := range []struct {
		name, chunk, next string
		fixture           func(*testing.T) *recordtest.Fixture
	}{
		{"missing chunk", "1", "bench record 'chunk' --help", func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Record.Chunks = nil
			f.Save()
			return f
		}},
		{"missing verification", "1", "bench record 'verification' --help", func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Record.Chunks[0].Verification = nil
			f.Save()
			return f
		}},
		{"missing axis", "1", "bench record 'review' --help", func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
			f.Save()
			return f
		}},
		{"missing completion", "", "bench record 'completion' --help", checkpointFixture},
		{"undispatched ticket", "1", "bench record 'assignment' --help", undispatched},
		{"changed plan", "1", "bench record 'amendment' --help", func(t *testing.T) *recordtest.Fixture {
			f := checkpointFixture(t)
			f.Plan.FinalVerification[0].Command = "go test -count=1 ./..."
			f.WritePlan()
			f.Commit("change the plan")
			return f
		}},
		{"unplanned chunk", "9", "bench preflight review 'example'", checkpointFixture},
		{"test-only repair", "1", "bench record 'confirmation' --help", repaired(modifyTest)},
		{"non-test repair", "1", "bench record 'review' --help", repaired(modifySource)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := tc.fixture(t)
			args := []string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", tc.chunk}
			if tc.chunk == "" {
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--complete"}
			}
			var out, diagnostic bytes.Buffer
			code := RunCommand(args, &out, &diagnostic)
			if want := "next=" + tc.next + "\n"; code == 0 || !strings.Contains(diagnostic.String(), want) {
				t.Fatalf("checkpoint = exit %d, stderr %q; want a refusal whose route is %q", code, diagnostic.String(), want)
			}
		})
	}
}
