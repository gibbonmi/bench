package recordtest

import (
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
)

// bitProbe is a probe of mutation that bit and restored.
func bitProbe(mutation string) rr.Probe {
	return rr.Probe{Mutation: mutation, Outcome: "bit", ExitCode: 1, Restore: "pass", NativeRef: Native("mutation failed at the required assertion; restore passed")}
}

// Refreeze moves the last recorded chunk to the current tip, as the chunk form does after a
// repair. The tip and the source digest move, the chunk verification is retained again at
// the new digest, and the reviews keep their pair.
func (f *Fixture) Refreeze() {
	f.T.Helper()
	chunk := &f.Record.Chunks[len(f.Record.Chunks)-1]
	digest, err := rr.SourceDigest(f.Root, f.Tree(), f.Record.Spec)
	if err != nil {
		f.T.Fatal(err)
	}
	chunk.Tip, chunk.SourceDigest = f.Tip(), digest
	chunk.Verification = append(chunk.Verification, f.Verification("chunk-"+chunk.ID+"-"+chunk.Tip, digest, f.Plan.Chunk(chunk.ID).Verification)...)
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
		FindingIDs: append([]string{}, findings...), Supersedes: supersedes, Probe: bitProbe("omit the repaired assertion"),
	})
}
