package specbuild

import (
	"bytes"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

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
	beforeReplay, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load before replay: found:%v err:%v", found, err)
	}
	if _, err := service.Integrate(t.Context(), "build demo", second.ID); err != nil {
		t.Fatalf("Integrate older sibling: %v", err)
	}
	afterReplay, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load after replay: found:%v err:%v", found, err)
	}
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
			run, found, err := fixture.service.load("build demo")
			if err != nil || !found {
				t.Fatalf("load checkpoint: found:%v err:%v", found, err)
			}
			if tc.name == "patch drift" {
				key, assigned, ok := assignmentFor(run, fixture.assigned.ID)
				if !ok {
					t.Fatal("missing assignment")
				}
				assigned.CheckpointPatch = "changed patch"
				run.Assignments[key] = assigned
			}
			tc.mutate(t, fixture, &run)
			if err := fixture.service.save(run); err != nil {
				t.Fatalf("save drift: %v", err)
			}
			before := git(t, fixture.root, "rev-parse", run.Candidate)
			if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err == nil {
				t.Fatal("Integrate accepted drift")
			}
			requireDelegatedCandidate(t, fixture.root, fixture.service, fixture.assigned.ID, before)
		})
	}
}

func TestIntegrateRoutesRetargetedCheckpointReferenceBackToDelegate(t *testing.T) {
	fixture := newCheckpointFixture(t)
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	run, found, err := fixture.service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load checkpoint: found:%v err:%v", found, err)
	}
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
	t.Run("unchanged retry succeeds", func(t *testing.T) {
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
		if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
			t.Fatalf("Integrate retry: %v", err)
		}
		run, found, err := fixture.service.load("build demo")
		if err != nil || !found || moved == "" {
			t.Fatalf("retry run = %#v, found:%v err:%v", run, found, err)
		}
		if parent := git(t, fixture.root, "rev-parse", run.CandidateTip+"^"); parent != moved {
			t.Fatalf("retry parent = %s, want injected candidate %s", parent, moved)
		}
		for path := range map[string]bool{"earlier.go": true, "internal/specbuild/checkpoint-change.go": true} {
			if git(t, fixture.root, "show", run.CandidateTip+":"+path) == "" {
				t.Fatalf("retry lost %s", path)
			}
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
	run, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load refused integration: found:%v err:%v", found, err)
	}
	if got := git(t, root, "rev-parse", run.Candidate); got != want || run.CandidateTip != want {
		t.Fatalf("candidate changed: ref=%s state=%s want=%s", got, run.CandidateTip, want)
	}
	_, assigned, ok := assignmentFor(run, assignmentID)
	if !ok || !assigned.DelegatePending || !strings.Contains(run.status().Next, assignmentID) {
		t.Fatalf("delegate route = %#v next=%q", assigned, run.status().Next)
	}
}
