package specbuild

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

type countingRunner struct {
	calls, commits int
	child, grand   string
	block          bool
}

func (r *countingRunner) Output(ctx context.Context, program string, args ...string) (string, error) {
	r.calls++
	if len(args) > 2 && args[2] == "commit-tree" {
		r.commits++
	}
	return (processRunner{}).Output(ctx, program, args...)
}
func (r *countingRunner) Run(ctx context.Context, command Command) (string, error) {
	for _, arg := range command.Args {
		if r.block && arg == "apply" {
			r.block = false
			return (processRunner{}).Output(ctx, "sh", "-c", "trap 'exit 0' TERM; sh -c 'trap \"\" TERM HUP; sleep 30 & echo $! > "+r.grand+"; wait' >/dev/null 2>&1 & echo $! > "+r.child+"; wait")
		}
	}
	r.calls++
	return (processRunner{}).Run(ctx, command)
}
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
	bind := func(t *testing.T, fixture checkpointFixture, rec *receipt) {
		tree := benchgit.TreeHash(fixture.assigned.Path)
		if tree == "none" {
			t.Fatal("live assignment tree is unavailable")
		}
		rec.Tree, rec.Probe.Tree = tree, tree
	}
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
			if err := os.Symlink("internal/specbuild/checkpoint-change.go", filepath.Join(fixture.assigned.Path, "outside-link")); err != nil {
				t.Fatal(err)
			}
			bind(t, fixture, rec)
			rec.Ownership = []string{"outside-link", "internal/specbuild/checkpoint-change.go"}
		}},
		{"unexplained path", func(t *testing.T, fixture checkpointFixture, rec *receipt) {
			write(t, filepath.Join(fixture.assigned.Path, "internal", "specbuild", "unexplained\nname.go"), "package specbuild\n")
			bind(t, fixture, rec)
		}},
		{"ticket file drift", func(t *testing.T, fixture checkpointFixture, _ *receipt) {
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
func TestCheckpointCreatesCommitFromVerifiedLiveAssignmentTree(t *testing.T) {
	fixture := newCheckpointFixture(t)
	runner := &countingRunner{}
	fixture.service.runner = runner
	before := checkpointSnapshotFor(t, fixture)
	if _, err := fixture.service.Checkpoint(context.Background(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	after := loadRun(t, fixture.service)
	_, stored, ok := assignmentFor(after, fixture.assigned.ID)
	candidate, ref := git(t, fixture.root, "rev-parse", after.Candidate), git(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/checkpoint/")
	if !ok || stored.Checkpoint == "" || stored.CheckpointRef == "" || stored.CheckpointTree != fixture.receipt.Tree || stored.ReceiptDigest == "" || candidate != before.candidate || fixture.gate.calls != 1 || ref != stored.CheckpointRef || runner.commits != 1 {
		t.Fatalf("checkpoint effects: attribution=%#v candidate=%s gate_calls=%d ref=%s commits=%d", stored, candidate, fixture.gate.calls, ref, runner.commits)
	}
	_, err := fixture.service.Checkpoint(context.Background(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n"))
	if replayRef := git(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/checkpoint/"); err != nil || replayRef != stored.CheckpointRef {
		t.Fatalf("checkpoint replay: error=%v ref=%q", err, replayRef)
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
	after := loadRun(t, fixture.service)
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
		changedTicket(t, fixture, "ticket dirt after durable integration\n")
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

func TestIntegrateRefusesTicketIndexDriftWithoutMutation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		apply func(*testing.T, checkpointFixture, string, []byte)
	}{
		{"staged mismatch", func(t *testing.T, fixture checkpointFixture, path string, original []byte) {
			write(t, path, string(original)+"\n")
			git(t, fixture.root, "add", filepath.ToSlash(path[len(fixture.root)+1:]))
			write(t, path, string(original))
		}},
		{"missing index", func(t *testing.T, fixture checkpointFixture, path string, _ []byte) {
			git(t, fixture.root, "rm", "--cached", "-q", filepath.ToSlash(path[len(fixture.root)+1:]))
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := checkpointedReleaseFixture(t)
			path := filepath.Join(fixture.root, "specs", "build demo", "tickets", "one.md")
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			tc.apply(t, fixture, path, original)
			before := checkpointSnapshotFor(t, fixture)
			if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil || !strings.Contains(err.Error(), "ticket no longer matches committed subject") {
				t.Fatalf("Integrate with %s = %v", tc.name, err)
			}
			if after := checkpointSnapshotFor(t, fixture); after != before {
				t.Fatalf("Integrate with %s mutated: before=%#v after=%#v", tc.name, before, after)
			}
		})
	}
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

func (o *releaseOwner) Release(context.Context, string, string, string, ReleaseEvidence) error {
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
func TestIntegrateJournalRecoversCandidateCASBeforeState(t *testing.T) {
	for _, point := range []string{"integrate/commit", "integrate/candidate-cas", "integrate/state", "integrate/release"} {
		t.Run(point, func(t *testing.T) {
			fixture, runner := checkpointedReleaseFixture(t), &countingRunner{}
			fixture.service.runner, fixture.service.fault = runner, func(got string) error {
				if got == point {
					return errors.New("injected")
				}
				return nil
			}
			if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil {
				t.Fatal("fault did not interrupt")
			}
			moved := git(t, fixture.root, "rev-parse", fixture.run.Candidate)
			fixture.service.fault = nil
			if point == "integrate/candidate-cas" {
				changedTicket(t, fixture, "ticket dirt after candidate CAS\n")
			}
			if point == "integrate/commit" {
				run := loadRun(t, fixture.service)
				op, _ := fixture.service.operation(run, "integrate", fixture.assigned.ID)
				original := op.Result
				op.Result = git(t, fixture.root, "commit-tree", git(t, fixture.assigned.Path, "rev-parse", "HEAD^{tree}"), "-p", fixture.run.Candidate, "-m", "swapped prepared result")
				run.Operations[operationID("integrate", fixture.assigned.ID)] = op
				saveRun(t, fixture.service, run)
				if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil || err.Error() != "spec build prepared integration result conflicts with replay" {
					t.Fatalf("tampered prepared result = %v", err)
				}
				after := loadRun(t, fixture.service)
				if git(t, fixture.root, "rev-parse", after.Candidate) != moved || after.Operations[operationID("integrate", fixture.assigned.ID)].Result != op.Result || runner.commits != 1 {
					t.Fatal("tampered result mutated integration")
				}
				op.Result = original
				run.Operations[operationID("integrate", fixture.assigned.ID)] = op
				saveRun(t, fixture.service, run)
			}
			if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
				t.Fatal(err)
			}
			run := loadRun(t, fixture.service)
			if git(t, fixture.root, "rev-parse", run.Candidate) != run.CandidateTip || (point != "integrate/commit" && run.CandidateTip != moved) || runner.calls == 0 || (point == "integrate/commit" && runner.commits != 1) {
				t.Fatal("integration replayed or bypassed runner")
			}
			if point == "integrate/candidate-cas" {
				requireReleased(t, fixture)
			}
		})
	}
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
	run := loadRun(t, fixture.service)
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
	run := loadRun(t, fixture.service)
	_, assigned, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok || assigned.CleanupPending || !assigned.Released {
		t.Fatalf("release state = %#v", assigned)
	}
}

func setAssignmentsReleased(t *testing.T, service *Service, released bool) record {
	run := loadRun(t, service)
	for key, assigned := range run.Assignments {
		assigned.Released = released
		run.Assignments[key] = assigned
	}
	saveRun(t, service, run)
	return run
}

type promotionGate struct {
	tree                       string
	red, accept, contradictory bool
	executions, validations    int
	disposition                GateDisposition
	err                        error
	inspect                    func(string, string)
}
