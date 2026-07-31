package specbuild

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckpointRejectsProbeProvenance(t *testing.T) {
	tests := []struct {
		name   string
		change func(*testing.T, checkpointFixture, *receipt) string
	}{
		{"delegate-produced", func(_ *testing.T, _ checkpointFixture, rec *receipt) string {
			rec.Probe.Producer = "delegate"
			return ""
		}},
		{"stale", func(t *testing.T, fixture checkpointFixture, rec *receipt) string {
			rec.Probe.Produced = receiptBeforeAssignment(t, fixture)
			return ""
		}},
		{"other-assignment", func(_ *testing.T, _ checkpointFixture, rec *receipt) string {
			rec.Probe.Assignment = "other-assignment"
			return ""
		}},
		{"other-tree", func(_ *testing.T, fixture checkpointFixture, rec *receipt) string {
			rec.Probe.Tree = fixture.run.CandidateTip
			return ""
		}},
		{"inside-worktree", func(t *testing.T, fixture checkpointFixture, rec *receipt) string {
			return checkpointPathInside(t, fixture, *rec)
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			rec := fixture.receipt
			path := tc.change(t, fixture, &rec)
			if path == "" {
				path = writeCheckpointReceipt(t, rec, "\n")
			}
			requireCheckpointRefusal(t, fixture, path, checkpointSnapshotFor(t, fixture))
		})
	}
}

func TestCheckpointAdmitsHonestRowOutcomes(t *testing.T) {
	for _, outcome := range []string{"passed", "already-covered", "not-tdd-able"} {
		t.Run(outcome, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			rec := fixture.receipt
			for i := range rec.Rows {
				rec.Rows[i].Outcome = outcome
			}
			if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, rec, "\n")); err != nil {
				t.Fatalf("Checkpoint(%s): %v", outcome, err)
			}
		})
	}
}

func TestCheckpointRefusesReceiptFieldFailuresWithoutMutation(t *testing.T) {
	tests := []struct {
		name   string
		change func(*receipt)
	}{
		{"charged row omitted", func(rec *receipt) { rec.Rows = nil }},
		{"charged row failed", func(rec *receipt) { rec.Rows[0].Outcome = "failed" }},
		{"focused check omitted", func(rec *receipt) { rec.Checks = nil }},
		{"focused check failed", func(rec *receipt) { rec.Checks[0].Passed = false }},
		{"probe producer omitted", func(rec *receipt) { rec.Probe.Producer = "" }},
		{"probe assignment omitted", func(rec *receipt) { rec.Probe.Assignment = "" }},
		{"probe tree omitted", func(rec *receipt) { rec.Probe.Tree = "" }},
		{"probe command omitted", func(rec *receipt) { rec.Probe.Command = "" }},
		{"probe exit failed", func(rec *receipt) { rec.Probe.Exit = 1 }},
		{"probe output omitted", func(rec *receipt) { rec.Probe.OutputDigest = "" }},
		{"probe timestamp omitted", func(rec *receipt) { rec.Probe.Produced = "" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			rec := fixture.receipt
			tc.change(&rec)
			requireCheckpointRefusal(t, fixture, writeCheckpointReceipt(t, rec, "\n"), checkpointSnapshotFor(t, fixture))
		})
	}
}

func TestCheckpointRereadsEveryLiveFact(t *testing.T) {
	tests := []struct {
		name   string
		change func(*testing.T, checkpointFixture, *receipt)
	}{
		{"base", func(_ *testing.T, _ checkpointFixture, rec *receipt) { rec.Base = "changed-base" }},
		{"tree", func(_ *testing.T, fixture checkpointFixture, rec *receipt) {
			rec.Tree, rec.Probe.Tree = fixture.run.CandidateTip, fixture.run.CandidateTip
		}},
		{"ticket digest", func(_ *testing.T, _ checkpointFixture, rec *receipt) { rec.TicketDigest = "changed-ticket" }},
		{"outside fence", func(t *testing.T, fixture checkpointFixture, rec *receipt) {
			write(t, filepath.Join(fixture.assigned.Path, "outside.go"), "package outside\n")
			git(t, fixture.assigned.Path, "add", ".")
			git(t, fixture.assigned.Path, "commit", "-qm", "outside fence")
			rec.Tree, rec.Probe.Tree = git(t, fixture.assigned.Path, "rev-parse", "HEAD^{tree}"), git(t, fixture.assigned.Path, "rev-parse", "HEAD^{tree}")
			rec.Ownership = []string{"outside.go", "internal/specbuild/checkpoint-change.go"}
		}},
		{"unexplained path", func(t *testing.T, fixture checkpointFixture, rec *receipt) {
			write(t, filepath.Join(fixture.assigned.Path, "internal", "specbuild", "unexplained.go"), "package specbuild\n")
			git(t, fixture.assigned.Path, "add", ".")
			git(t, fixture.assigned.Path, "commit", "-qm", "unexplained path")
			rec.Tree, rec.Probe.Tree = git(t, fixture.assigned.Path, "rev-parse", "HEAD^{tree}"), git(t, fixture.assigned.Path, "rev-parse", "HEAD^{tree}")
		}},
		{"assumption drift", func(t *testing.T, fixture checkpointFixture, _ *receipt) {
			changedTicket(t, fixture, "# One\n\nOwnership fence: internal/specbuild\nAssumptions: changed contract\n\n- [ ] [R10-R15] checkpoint receipt\n- [ ] [R54] framing\n")
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			rec := fixture.receipt
			tc.change(t, fixture, &rec)
			requireCheckpointRefusal(t, fixture, writeCheckpointReceipt(t, rec, "\n"), checkpointSnapshotFor(t, fixture))
		})
	}
}

func TestCheckpointRequiresOneFinalNewlineFraming(t *testing.T) {
	fixture := newCheckpointFixture(t)
	before := checkpointSnapshotFor(t, fixture)
	noNewline := writeCheckpointReceipt(t, fixture.receipt, "")
	_, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, noNewline)
	if !errors.Is(err, errMalformedReceipt) {
		t.Fatalf("no-final-newline error = %v", err)
	}
	if after := checkpointSnapshotFor(t, fixture); after != before {
		t.Fatalf("no-final-newline mutated state: before=%#v after=%#v", before, after)
	}
	malformed := filepath.Join(t.TempDir(), "malformed.json")
	write(t, malformed, "{\n")
	_, err = fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, malformed)
	if !errors.Is(err, errMalformedReceipt) {
		t.Fatalf("malformed error = %v", err)
	}
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint with final newline: %v", err)
	}
}

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

func TestIntegrateRequiresVerifiedCheckpointAndAdvancesOneAttributedCandidate(t *testing.T) {
	fixture := newCheckpointFixture(t)
	before := git(t, fixture.root, "rev-parse", fixture.run.Candidate)
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil {
		t.Fatal("Integrate accepted an unverified assignment")
	}
	if got := git(t, fixture.root, "rev-parse", fixture.run.Candidate); got != before {
		t.Fatalf("unverified integration moved candidate from %s to %s", before, got)
	}
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate: %v", err)
	}
	after, found, err := fixture.service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load integrated run: found:%v err:%v", found, err)
	}
	_, assigned, ok := assignmentFor(after, fixture.assigned.ID)
	if !ok || assigned.Integrated != after.CandidateTip {
		t.Fatalf("integration attribution = %#v", assigned)
	}
	if parent := git(t, fixture.root, "rev-parse", after.CandidateTip+"^"); parent != before {
		t.Fatalf("candidate parent = %s, want %s", parent, before)
	}
	if subject := git(t, fixture.root, "show", "-s", "--format=%B", after.CandidateTip); !strings.Contains(subject, "run="+after.Run) || !strings.Contains(subject, "assignment="+fixture.assigned.ID) || !strings.Contains(subject, "checkpoint="+assigned.Checkpoint) {
		t.Fatalf("candidate attribution = %q", subject)
	}
	if tree := git(t, fixture.root, "rev-parse", after.CandidateTip+"^{tree}"); tree != assigned.CheckpointTree {
		t.Fatalf("candidate tree = %s, want checkpoint tree %s", tree, assigned.CheckpointTree)
	}
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate replay: %v", err)
	}
	if got := git(t, fixture.root, "rev-parse", after.Candidate); got != after.CandidateTip {
		t.Fatalf("integration replay changed candidate from %s to %s", after.CandidateTip, got)
	}
}

func TestIntegrateReleasesOnlyAfterDurableProvenance(t *testing.T) {
	t.Run("unavailable owner stays pending", func(t *testing.T) {
		fixture := checkpointedReleaseFixture(t)
		fixture.service.worktrees = &countingOwner{}

		status, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID)
		if err == nil {
			t.Fatal("Integrate accepted an unavailable release owner")
		}
		requirePendingRelease(t, fixture, status)
	})

	t.Run("failed release resumes without a second candidate commit", func(t *testing.T) {
		fixture := checkpointedReleaseFixture(t)
		owner := &releaseOwner{err: errors.New("release unavailable")}
		fixture.service.worktrees = owner
		owner.inspect = func() { requireReleaseProvenance(t, fixture) }

		status, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID)
		if err == nil {
			t.Fatal("Integrate accepted a failed release")
		}
		requirePendingRelease(t, fixture, status)
		candidate := status.Subject
		commits := git(t, fixture.root, "rev-list", "--count", fixture.run.Candidate)

		owner.err = nil
		status, err = fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID)
		if err != nil {
			t.Fatalf("Integrate resume: %v", err)
		}
		if status.Subject != candidate {
			t.Fatalf("resume advanced candidate from %s to %s", candidate, status.Subject)
		}
		if got := git(t, fixture.root, "rev-list", "--count", fixture.run.Candidate); got != commits {
			t.Fatalf("resume candidate commits = %s, want %s", got, commits)
		}
		if owner.calls != 2 || owner.released != 1 {
			t.Fatalf("release calls=%d successful=%d, want 2 and 1", owner.calls, owner.released)
		}
		requireReleased(t, fixture)
	})

	t.Run("successful release clears pending", func(t *testing.T) {
		fixture := checkpointedReleaseFixture(t)
		owner := &releaseOwner{}
		fixture.service.worktrees = owner
		owner.inspect = func() { requireReleaseProvenance(t, fixture) }

		if _, err := fixture.service.Integrate(context.Background(), "build demo", fixture.assigned.ID); err != nil {
			t.Fatalf("Integrate: %v", err)
		}
		if owner.calls != 1 || owner.released != 1 {
			t.Fatalf("release calls=%d successful=%d, want 1 and 1", owner.calls, owner.released)
		}
		requireReleased(t, fixture)
	})
}

func TestTypedNilReleaseOwnerStaysPending(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	var owner *releaseOwner
	fixture.service.worktrees = owner

	status, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID)
	if err == nil || err.Error() != "spec build integrate requires a release-capable worktree owner" {
		t.Fatalf("Integrate error = %v", err)
	}
	requirePendingRelease(t, fixture, status)
}

type releaseOwner struct {
	realOwner
	calls, released int
	err             error
	inspect         func()
}

func (o *releaseOwner) Release(context.Context, string, string, string) error {
	o.calls++
	if o.inspect != nil {
		o.inspect()
	}
	if o.err != nil {
		return o.err
	}
	o.released++
	return nil
}

func checkpointedReleaseFixture(t *testing.T) checkpointFixture {
	t.Helper()
	fixture := newCheckpointFixture(t)
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	return fixture
}

func requireReleaseProvenance(t *testing.T, fixture checkpointFixture) {
	t.Helper()
	run, found, err := fixture.service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load at release: found:%v err:%v", found, err)
	}
	_, assigned, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok || assigned.Checkpoint == "" || assigned.CheckpointRef == "" || assigned.CheckpointTree == "" || assigned.ReceiptDigest == "" || assigned.CheckpointPatch == "" || assigned.Integrated == "" || !assigned.CleanupPending || assigned.Released {
		t.Fatalf("release provenance = %#v", assigned)
	}
	if assigned.Integrated != run.CandidateTip || !refAt(fixture.root, run.Candidate, assigned.Integrated) {
		t.Fatalf("candidate provenance = %#v", run)
	}
}

func requirePendingRelease(t *testing.T, fixture checkpointFixture, status Status) {
	t.Helper()
	if status.Next != "release assignment "+fixture.assigned.ID {
		t.Fatalf("next = %q", status.Next)
	}
	requireReleaseProvenance(t, fixture)
}

func requireReleased(t *testing.T, fixture checkpointFixture) {
	t.Helper()
	run, found, err := fixture.service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load released: found:%v err:%v", found, err)
	}
	_, assigned, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok || assigned.CleanupPending || !assigned.Released {
		t.Fatalf("release state = %#v", assigned)
	}
}
