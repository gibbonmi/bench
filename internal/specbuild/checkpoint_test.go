package specbuild

import (
	"context"
	"testing"
)

func TestCheckpointCreatesOneAttributedCommitWithoutCandidateOrGateMutation(t *testing.T) {
	fixture := newCheckpointFixture(t)
	before := checkpointSnapshotFor(t, fixture)
	if _, err := fixture.service.Checkpoint(context.Background(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	after, found, err := fixture.service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load checkpoint: found:%v err:%v", found, err)
	}
	_, stored, ok := assignmentFor(after, fixture.assigned.ID)
	if !ok || stored.Checkpoint == "" || stored.CheckpointRef == "" || stored.ReceiptDigest == "" {
		t.Fatalf("checkpoint attribution = %#v", stored)
	}
	if got := git(t, fixture.root, "rev-parse", after.Candidate); got != before.candidate {
		t.Fatalf("checkpoint moved candidate from %s to %s", before.candidate, got)
	}
	if fixture.gate.calls != 1 {
		t.Fatalf("gate calls = %d, want 1 bootstrap call", fixture.gate.calls)
	}
	if got := git(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/checkpoint/"); got != stored.CheckpointRef {
		t.Fatalf("checkpoint refs = %q, want one %q", got, stored.CheckpointRef)
	}
	if _, err := fixture.service.Checkpoint(context.Background(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint replay: %v", err)
	}
	if got := git(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/checkpoint/"); got != stored.CheckpointRef {
		t.Fatalf("checkpoint replay created another ref: %q", got)
	}
}
