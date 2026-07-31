package specbuild

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

type countingOwner struct{ calls int }

func (o *countingOwner) Create(context.Context, string, string, string, string) (OwnedWorktree, error) {
	o.calls++
	return OwnedWorktree{}, nil
}

type releaseOwner struct {
	realOwner
	calls, released int
	err             error
	inspect         func()
}

func TestReviewRefusesIncompleteOrHostileReceiptsWithoutMutation(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatal(err)
	}
	run := loadRun(t, fixture.service)
	valid := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		path func(*testing.T) string
	}{
		{"missing axis", func(t *testing.T) string {
			return writeReviewReceipt(t, reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: valid.Axes[:2]})
		}},
		{"duplicate axis", func(t *testing.T) string {
			receipt := valid
			receipt.Axes[2].Axis = "Spec"
			return writeReviewReceipt(t, receipt)
		}},
		{"unknown axis", func(t *testing.T) string {
			receipt := valid
			receipt.Axes[2].Axis = "Other"
			return writeReviewReceipt(t, receipt)
		}},
		{"unresolved finding", func(t *testing.T) string {
			receipt := valid
			receipt.Axes[0].Findings = []reviewFinding{{ID: "f"}}
			return writeReviewReceipt(t, receipt)
		}},
		{"wrong candidate", func(t *testing.T) string {
			receipt := valid
			receipt.Candidate = "wrong"
			return writeReviewReceipt(t, receipt)
		}},
		{"wrong run", func(t *testing.T) string {
			receipt := valid
			receipt.Run = "wrong"
			return writeReviewReceipt(t, receipt)
		}},
		{"malformed", func(t *testing.T) string {
			path := filepath.Join(t.TempDir(), "review.json")
			write(t, path, "{\n")
			return path
		}},
		{"fifo", func(t *testing.T) string {
			path := filepath.Join(t.TempDir(), "review.fifo")
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
			return path
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := fixture.service.Review(t.Context(), "build demo", test.path(t)); err == nil {
				t.Fatal("Review unexpectedly accepted receipt")
			}
			after, err := os.ReadFile(statePath)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("review refusal mutated state: %v", err)
			}
		})
	}
}

func writeReviewReceipt(t *testing.T, receipt reviewReceipt) string {
	t.Helper()
	data, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "review.json")
	write(t, path, string(data)+"\n")
	return path
}

func TestCheckpointRejectsHostileReceiptsBeforeMutationOrBlocking(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*testing.T, checkpointFixture) string
	}{
		{"empty", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "empty.json")
			write(t, path, "")
			return path
		}},
		{"oversized", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "oversized.json")
			write(t, path, strings.Repeat("x", int(bounds.ControlRecordLimit)+1))
			return path
		}},
		{"malformed", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "malformed.json")
			write(t, path, "{\n")
			return path
		}},
		{"unreadable", func(t *testing.T, fixture checkpointFixture) string {
			path := writeCheckpointReceipt(t, fixture.receipt, "\n")
			if err := os.Chmod(path, 0); err != nil {
				t.Fatal(err)
			}
			return path
		}},
		{"fifo", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "receipt.fifo")
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
			return path
		}},
		{"device", func(_ *testing.T, _ checkpointFixture) string { return "/dev/null" }},
		{"socket", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "receipt.sock")
			listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = listener.Close() })
			return path
		}},
		{"regular symlink", func(t *testing.T, fixture checkpointFixture) string {
			target := writeCheckpointReceipt(t, fixture.receipt, "\n")
			path := filepath.Join(t.TempDir(), "receipt-link.json")
			if err := os.Symlink(target, path); err != nil {
				t.Fatal(err)
			}
			return path
		}},
		{"dangling symlink", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "dangling-link.json")
			if err := os.Symlink("missing.json", path); err != nil {
				t.Fatal(err)
			}
			return path
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			before := checkpointSnapshotFor(t, fixture)
			path := tc.setup(t, fixture)
			done := make(chan error, 1)
			go func() {
				_, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, path)
				done <- err
			}()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("Checkpoint unexpectedly accepted hostile receipt")
				}
			case <-time.After(time.Second):
				t.Fatal("Checkpoint blocked on hostile receipt")
			}
			if after := checkpointSnapshotFor(t, fixture); after != before {
				t.Fatalf("hostile receipt mutated state: before=%#v after=%#v", before, after)
			}
		})
	}
}

func TestIntegrateReplaysDisjointSiblingPatchByteForByte(t *testing.T) {
	root, service, first, second, run := siblingCheckpoints(t, "internal/specbuild/first.go", "package specbuild\n// first\n", "internal/specbuild/second.go", "package specbuild\n// second\n")
	_, stored, ok := assignmentFor(run, second.ID)
	if !ok {
		t.Fatal("missing second assignment")
	}
	patch, err := checkpointPatch(root, second.Base, stored.Checkpoint)
	if err != nil {
		t.Fatalf("recorded patch: %v", err)
	}
	if _, err := service.Integrate(t.Context(), "build demo", first.ID); err != nil {
		t.Fatalf("Integrate first: %v", err)
	}
	beforeReplay := loadRun(t, service)
	if _, err := service.Integrate(t.Context(), "build demo", second.ID); err != nil {
		t.Fatalf("Integrate older sibling: %v", err)
	}
	afterReplay := loadRun(t, service)
	replayedPatch, err := checkpointPatch(root, beforeReplay.CandidateTip, afterReplay.CandidateTip)
	if err != nil || !bytes.Equal(replayedPatch, patch) {
		t.Fatalf("replayed patch changed: equal:%v err:%v", bytes.Equal(replayedPatch, patch), err)
	}
	for path, want := range map[string]string{"internal/specbuild/first.go": "package specbuild\n// first", "internal/specbuild/second.go": "package specbuild\n// second"} {
		if got := git(t, root, "show", afterReplay.CandidateTip+":"+path); got != want {
			t.Fatalf("candidate %s = %q, want %q", path, got, want)
		}
	}
}

func TestIntegrateRefusesOverlapAndApplyConflictWithoutMovingCandidate(t *testing.T) {
	for _, tc := range []struct {
		name, firstPath, secondPath string
	}{
		{"overlap", "internal/specbuild/shared.go", "internal/specbuild/shared.go"},
		{"apply conflict", "internal/specbuild/collision", "internal/specbuild/collision/nested.go"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, service, first, second, _ := siblingCheckpoints(t, tc.firstPath, "package specbuild\n", tc.secondPath, "package specbuild\n")
			if _, err := service.Integrate(t.Context(), "build demo", first.ID); err != nil {
				t.Fatalf("Integrate first: %v", err)
			}
			before := git(t, root, "rev-parse", "refs/bench/specbuild/candidate/"+digest(filepath.Join(root, "specs", "build demo", "spec.md")))
			if _, err := service.Integrate(t.Context(), "build demo", second.ID); err == nil {
				t.Fatal("Integrate accepted unsafe replay")
			}
			requireDelegatedCandidate(t, root, service, second.ID, before)
		})
	}
}

func TestIntegrateRefusesCheckpointAndTicketDriftWithoutMovingCandidate(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, checkpointFixture, *record)
	}{
		{"patch drift", func(*testing.T, checkpointFixture, *record) {}},
		{"ownership drift", func(t *testing.T, fixture checkpointFixture, _ *record) {
			changedTicket(t, fixture, "# One\n\nOwnership fence: internal/other\nAssumptions: receipt contract\n\n- [ ] [R10-R15] checkpoint receipt\n- [ ] [R54] framing\n")
		}},
		{"changed assumptions", func(t *testing.T, fixture checkpointFixture, _ *record) {
			changedTicket(t, fixture, "# One\n\nOwnership fence: internal/specbuild\nAssumptions: changed contract\n\n- [ ] [R10-R15] checkpoint receipt\n- [ ] [R54] framing\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
				t.Fatalf("Checkpoint: %v", err)
			}
			run := loadRun(t, fixture.service)
			if tc.name == "patch drift" {
				key, assigned, ok := assignmentFor(run, fixture.assigned.ID)
				if !ok {
					t.Fatal("missing assignment")
				}
				assigned.CheckpointPatch = "changed patch"
				run.Assignments[key] = assigned
			}
			tc.mutate(t, fixture, &run)
			saveRun(t, fixture.service, run)
			before := git(t, fixture.root, "rev-parse", run.Candidate)
			beforeState := checkpointSnapshotFor(t, fixture)
			if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil {
				t.Fatal("Integrate accepted drift")
			}
			if tc.name == "patch drift" {
				requireDelegatedCandidate(t, fixture.root, fixture.service, fixture.assigned.ID, before)
			} else if checkpointSnapshotFor(t, fixture) != beforeState {
				t.Fatal("shared precondition mutated a drifted working subject")
			}
		})
	}
}

func TestIntegrateRoutesRetargetedCheckpointReferenceBackToDelegate(t *testing.T) {
	fixture := newCheckpointFixture(t)
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	run := loadRun(t, fixture.service)
	_, assigned, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	before := git(t, fixture.root, "rev-parse", run.Candidate)
	git(t, fixture.root, "update-ref", assigned.CheckpointRef, assigned.Base, assigned.Checkpoint)
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil {
		t.Fatal("Integrate accepted a retargeted checkpoint ref")
	}
	requireDelegatedCandidate(t, fixture.root, fixture.service, fixture.assigned.ID, before)
}

func TestIntegrateRetriesOnlyWhileTheCheckpointContractHolds(t *testing.T) {
	t.Run("candidate change conflicts with prepared result", func(t *testing.T) {
		fixture := newCheckpointFixture(t)
		if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
			t.Fatalf("Checkpoint: %v", err)
		}
		before := git(t, fixture.root, "rev-parse", fixture.run.Candidate)
		moved := ""
		fixture.service.beforeCandidateCAS = func() {
			if moved != "" {
				return
			}
			write(t, filepath.Join(fixture.root, "earlier.go"), "package earlier\n")
			git(t, fixture.root, "add", ".")
			git(t, fixture.root, "commit", "-qm", "concurrent candidate move")
			moved = git(t, fixture.root, "rev-parse", "HEAD")
			git(t, fixture.root, "update-ref", fixture.run.Candidate, moved, before)
		}
		if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil || err.Error() != "spec build prepared integration result conflicts with replay" {
			t.Fatalf("Integrate retry = %v", err)
		}
		run := loadRun(t, fixture.service)
		op, _ := fixture.service.operation(run, "integrate", fixture.assigned.ID)
		if moved == "" || git(t, fixture.root, "rev-parse", run.Candidate) != moved || run.CandidateTip != before || op.Result == "" {
			t.Fatalf("retry run = %#v", run)
		}
	})
	t.Run("drift after candidate move refuses", func(t *testing.T) {
		fixture := newCheckpointFixture(t)
		if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
			t.Fatalf("Checkpoint: %v", err)
		}
		before := git(t, fixture.root, "rev-parse", fixture.run.Candidate)
		moved := ""
		fixture.service.beforeCandidateCAS = func() {
			if moved != "" {
				return
			}
			write(t, filepath.Join(fixture.root, "earlier.go"), "package earlier\n")
			git(t, fixture.root, "add", ".")
			git(t, fixture.root, "commit", "-qm", "concurrent candidate move")
			moved = git(t, fixture.root, "rev-parse", "HEAD")
			git(t, fixture.root, "update-ref", fixture.run.Candidate, moved, before)
			changedTicket(t, fixture, "# One\n\nOwnership fence: internal/specbuild\nAssumptions: changed contract\n\n- [ ] [R10-R15] checkpoint receipt\n- [ ] [R54] framing\n")
		}
		if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil {
			t.Fatal("Integrate retried after changed assumptions")
		}
		if moved == "" {
			t.Fatal("candidate move seam did not run")
		}
		requireDelegatedCandidate(t, fixture.root, fixture.service, fixture.assigned.ID, moved)
	})
}

func requireDelegatedCandidate(t *testing.T, root string, service *Service, assignmentID, want string) {
	t.Helper()
	run := loadRun(t, service)
	if got := git(t, root, "rev-parse", run.Candidate); got != want || run.CandidateTip != want {
		t.Fatalf("candidate changed: ref=%s state=%s want=%s", got, run.CandidateTip, want)
	}
	_, assigned, ok := assignmentFor(run, assignmentID)
	if !ok || !assigned.DelegatePending || !strings.Contains(run.status().Next, assignmentID) {
		t.Fatalf("delegate route = %#v next=%q", assigned, run.status().Next)
	}
}

func TestPromoteRecomposesAWorkingAdvanceBeforeGate(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	before, gateCalls := loadRun(t, fixture.service), fixture.gate.calls
	advanceWorking(t, fixture.root)
	working := git(t, fixture.root, "rev-parse", "HEAD")
	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.Next != "bench spec build review build demo" || fixture.gate.calls != gateCalls {
		t.Fatalf("recomposition = %#v, %v; gate calls=%d", status, err, fixture.gate.calls)
	}
	after := loadRun(t, fixture.service)
	if after.Base != working || after.CandidateTip == before.CandidateTip || after.Review != nil || git(t, fixture.root, "rev-parse", "HEAD") != working || git(t, fixture.root, "rev-parse", after.Candidate) != after.CandidateTip {
		t.Fatalf("recomposed run = %#v", after)
	}
	for _, path := range []string{"advanced.txt", "internal/specbuild/checkpoint-change.go"} {
		if git(t, fixture.root, "show", after.CandidateTip+":"+path) == "" {
			t.Fatalf("recomposed candidate omitted %s", path)
		}
	}
}

func (promotionGate) Bootstrap(context.Context, string, string, string) error { return nil }
