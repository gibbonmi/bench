package specbuild

import (
	"errors"
	"path/filepath"
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
