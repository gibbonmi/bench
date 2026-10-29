package recordtest

import (
	"github.com/gibbonmi/bench/internal/gittest"
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

// Refreeze moves the last recorded chunk to the current tip and the current plan, as a record
// call after a repair does. The tip, the plan digest, and the source digest move. Each
// verification entry keeps the commit that it graded, and the reviews keep their pair.
func (f *Fixture) Refreeze() {
	f.T.Helper()
	chunk := &f.Record.Chunks[len(f.Record.Chunks)-1]
	chunk.Tip, chunk.PlanDigest, chunk.SourceDigest = f.Tip(), f.Plan.Digest, f.digestAt(f.Tip())
}

// RecordAtTip appends to the last recorded chunk one entry for each planned requirement, by its
// current owner and bound to the current tip, as the record call of a repair session does.
func (f *Fixture) RecordAtTip() {
	f.T.Helper()
	chunk := &f.Record.Chunks[len(f.Record.Chunks)-1]
	_, verification := f.freeze(chunk.ID, "chunk-"+chunk.ID+"-"+f.Tip())
	chunk.Verification = append(chunk.Verification, verification...)
}

// LaterTicket commits a later ticket that writes later.txt and each path of paths, and moves
// the last recorded chunk and its reviews to that commit and the current plan through Refreeze.
func (f *Fixture) LaterTicket(paths ...string) {
	f.T.Helper()
	for _, path := range append([]string{"later.txt"}, paths...) {
		f.Write(path, "later ticket\n")
	}
	f.Commit("later ticket")
	f.Refreeze()
	chunk := &f.Record.Chunks[len(f.Record.Chunks)-1]
	for i := range chunk.Reviews {
		chunk.Reviews[i].Tip, chunk.Reviews[i].SourceDigest = chunk.Tip, chunk.SourceDigest
	}
}

// Bind binds entry to commit and to the source digest of the tree of commit.
func (f *Fixture) Bind(entry *rr.Verification, commit string) {
	f.T.Helper()
	entry.SourceCommit, entry.SourceDigest = commit, f.digestAt(commit)
}

// OffLine commits the tree of the last chunk tip as a child of the chunk base on no branch, so
// the commit descends from the base and is not an ancestor of the tip.
func (f *Fixture) OffLine() string {
	f.T.Helper()
	chunk := f.Record.Chunks[len(f.Record.Chunks)-1]
	return f.Git(append(gittest.IdentityArgs(), "commit-tree", chunk.Tip+"^{tree}", "-p", chunk.Base, "-m", "off the chunk line")...)
}

// digestAt is the source digest of the tree of commit.
func (f *Fixture) digestAt(commit string) string {
	f.T.Helper()
	digest, err := rr.SourceDigest(f.Root, f.Git("rev-parse", "--verify", commit+"^{tree}"), f.Record.Spec)
	if err != nil {
		f.T.Fatal(err)
	}
	return digest
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
