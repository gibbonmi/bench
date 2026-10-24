package recordtest

import (
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
)

// bitProbe is a probe of mutation that bit and restored.
func bitProbe(mutation string) rr.Probe {
	return rr.Probe{Mutation: mutation, Outcome: "bit", ExitCode: 1, Restore: "pass", NativeRef: Native("mutation failed at the required assertion; restore passed")}
}

// ConfirmationProbe is the probe that Confirm records.
func ConfirmationProbe() rr.Probe {
	return bitProbe("omit the repaired assertion")
}

// Refreeze moves the last recorded chunk to the current tip and the current plan, as the chunk
// form does after a repair. The tip, the plan digest, and the source digest move, the chunk
// verification is retained again at the new digest, and the reviews keep their pair.
func (f *Fixture) Refreeze() {
	f.T.Helper()
	chunk := &f.Record.Chunks[len(f.Record.Chunks)-1]
	chunk.Tip, chunk.PlanDigest = f.Tip(), f.Plan.Digest
	digest, verification := f.freeze(chunk.ID, "chunk-"+chunk.ID+"-"+chunk.Tip)
	chunk.SourceDigest, chunk.Verification = digest, append(chunk.Verification, verification...)
}

// Reassign commits the assignment-only plan change that dispatches a repair session: through
// the assignment writer, ticket gains a user-directed replacement of its current author. The
// plan reloads, so the fixture record names the new plan digest.
func (f *Fixture) Reassign(ticket string) {
	f.T.Helper()
	execution, history := f.Plan.Execution, f.Plan.Execution.Assignments[ticket]
	repair := Assign("repair-" + ticket)
	repair.Predecessor, repair.Trigger = history[len(history)-1].Session, "user-directed"
	repair.Stopped, repair.Preserved = "fixture:stopped-"+ticket, "fixture-preserved-"+ticket
	call := rr.AssignmentCall{Source: f.Tip(), Ticket: ticket, Assignment: repair, RunID: execution.RunID,
		OrchestratorSession: execution.OrchestratorSession, ReviewMode: execution.ReviewMode, AuthorLimit: execution.AuthorLimit}
	if _, err := rr.RecordAssignment(f.Root, f.spec, call); err != nil {
		f.T.Fatal(err)
	}
	f.Commit("dispatch the repair of " + ticket)
	f.Reload()
}

// Confirm appends to the last recorded chunk the confirmation that the completion
// reconciler writes for a repair from the reviewed tip to the chunk tip. The confirmation
// closes findings, holds a probe that bit, and supersedes the prior confirmation of the chunk.
func (f *Fixture) Confirm(id, reviewed string, findings ...string) {
	f.T.Helper()
	chunk := &f.Record.Chunks[len(f.Record.Chunks)-1]
	reconciler, _ := f.owes(rr.Requirement{})
	supersedes := []string{}
	if prior := len(chunk.Confirmations); prior != 0 {
		supersedes = []string{chunk.Confirmations[prior-1].ID}
	}
	chunk.Confirmations = append(chunk.Confirmations, rr.Confirmation{
		Evidence: f.EvidenceAs(id, chunk.SourceDigest, rr.ConfirmationRole, reconciler),
		Base:     chunk.Base, ReviewedTip: reviewed, RepairedTip: chunk.Tip,
		FindingIDs: append([]string{}, findings...), Supersedes: supersedes, Probe: ConfirmationProbe(),
	})
}
