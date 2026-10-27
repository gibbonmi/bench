package reviewrecord

import (
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// TestCarryAndCheckpointAgree grades each entry variant through the carry candidate and
// through the chunk verification grade under the plan performer. The two give one verdict,
// because they read one pass predicate. The grade alone accepts an entry with no source
// commit, because EntryCurrent then compares only the digest.
func TestCarryAndCheckpointAgree(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	commit := func(message string) string {
		gittest.Output(t, root, "commit", "-q", "--allow-empty", "-m", message)
		return gittest.Output(t, root, "rev-parse", "HEAD")
	}
	base, source, tip := commit("base"), commit("source"), commit("tip")
	digest, err := commitDigest(root, renderSpec, tip)
	if err != nil {
		t.Fatal(err)
	}
	zero, one := 0, 1
	requirement := Requirement{ID: "tests", Command: "bench test --package ./internal/x", Probe: "omit the check"}
	record := Record{Version: 1, ImplementationSession: "author"}
	passing := func() Verification {
		probe := Probe{Mutation: requirement.Probe, Outcome: "bit", ExitCode: 1, Restore: "pass"}
		return Verification{Evidence: Evidence{Performer: "author", Role: "author-verification", SourceDigest: digest, State: "completed", Outcome: "pass"}, Requirement: requirement.ID, Command: requirement.Command, ExitCode: &zero, Probe: &probe, SourceCommit: source}
	}
	variant := func(change func(*Verification)) Verification {
		entry := passing()
		change(&entry)
		return entry
	}
	for _, row := range []struct {
		name    string
		entries []Verification
		carries bool
	}{
		{"a pass", []Verification{passing()}, true},
		{"a failed outcome", []Verification{variant(func(e *Verification) { e.Outcome = "fail"; e.ExitCode = &one })}, false},
		{"a missing exit code", []Verification{variant(func(e *Verification) { e.ExitCode = nil })}, false},
		{"a changed command", []Verification{variant(func(e *Verification) { e.Command += " --run X" })}, false},
		{"a probe that did not bite", []Verification{variant(func(e *Verification) { e.Probe.Outcome = "silent" })}, false},
		{"a probe whose restore failed", []Verification{variant(func(e *Verification) { e.Probe.Restore = "fail" })}, false},
		{"a source commit at the chunk base", []Verification{variant(func(e *Verification) { e.SourceCommit = base })}, false},
		{"no source commit", []Verification{variant(func(e *Verification) { e.SourceCommit = "" })}, false},
		{"an earlier pass and a failed last entry", []Verification{passing(), variant(func(e *Verification) { e.Outcome = "fail" })}, false},
	} {
		chunk := Chunk{ID: "1", Base: base, Tip: tip, SourceDigest: digest, Verification: row.entries}
		candidate, carries, err := CarryCandidate(root, renderSpec, chunk, requirement)
		if err != nil {
			t.Fatalf("%s: carry candidate: %v", row.name, err)
		}
		current := func(entry Verification) (bool, error) { return EntryCurrent(root, renderSpec, chunk, entry) }
		graded := checkVerification(chunk.Verification, []Requirement{requirement}, Plan{}, record, current, Repair{Chunk: chunk.ID})
		last := row.entries[len(row.entries)-1]
		if carries != row.carries || carries && candidate.SourceCommit != last.SourceCommit {
			t.Errorf("%s: carry candidate %v (%+v); want %v for the last entry", row.name, carries, candidate, row.carries)
		}
		if (graded == nil) != (carries || last.SourceCommit == "") {
			t.Errorf("%s: carry candidate %v, chunk grade %v; want one verdict", row.name, carries, graded)
		}
	}
}
