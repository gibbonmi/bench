package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

type checkpointFixture struct {
	root     string
	service  *Service
	gate     *countingGate
	assigned Assignment
	run      record
	receipt  receipt
}

func newCheckpointFixture(t *testing.T, configure ...func(string)) checkpointFixture {
	t.Helper()
	root := repo(t)
	write(t, filepath.Join(root, ".gitignore"), "dist/\n")
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\nAssumptions: receipt contract\n\n- [ ] [R10-R15] checkpoint receipt\n- [ ] [R54] framing\n")
	applyCheckpointFixtureConfiguration(root, configure)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "receipt ticket")
	gate := &countingGate{}
	service := New(root, gate, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", "checkpoint receipt")
	if err != nil {
		t.Fatalf("Assign: %v", err)
	}
	write(t, filepath.Join(assigned.Path, "internal", "specbuild", "checkpoint-change.go"), "package specbuild\n")
	return checkpointFixture{
		root: root, service: service, gate: gate, assigned: assigned, run: loadRun(t, service),
		receipt: checkpointReceiptFor(t, service, assigned, []string{"internal/specbuild/checkpoint-change.go"}),
	}
}
func writeCheckpointReceipt(t *testing.T, rec receipt, suffix string) string {
	t.Helper()
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "receipt.json")
	write(t, path, string(data)+suffix)
	return path
}

type checkpointSnapshot struct {
	state, refs, candidate, tree, status string
}

func checkpointSnapshotFor(t *testing.T, fixture checkpointFixture) checkpointSnapshot {
	t.Helper()
	path, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	state, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return checkpointSnapshot{state: string(state), refs: git(t, fixture.root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/"), candidate: git(t, fixture.root, "rev-parse", fixture.run.Candidate), tree: git(t, fixture.assigned.Path, "rev-parse", "HEAD^{tree}"), status: git(t, fixture.assigned.Path, "status", "--porcelain", "--untracked-files=all")}
}
func requireCheckpointRefusal(t *testing.T, fixture checkpointFixture, path string, before checkpointSnapshot) {
	t.Helper()
	if _, err := fixture.service.Checkpoint(context.Background(), "build demo", fixture.assigned.ID, path); err == nil {
		t.Fatal("Checkpoint unexpectedly accepted receipt")
	}
	after := checkpointSnapshotFor(t, fixture)
	if after != before {
		t.Fatalf("receipt refusal mutated state: before=%#v after=%#v", before, after)
	}
}
func checkpointPathInside(t *testing.T, fixture checkpointFixture, rec receipt) string {
	t.Helper()
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(fixture.assigned.Path, "receipt.json")
	write(t, path, string(data)+"\n")
	return path
}
func receiptBeforeAssignment(t *testing.T, fixture checkpointFixture) string {
	t.Helper()
	_, stored, ok := assignmentFor(fixture.run, fixture.assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	created, err := time.Parse(time.RFC3339Nano, stored.Created)
	if err != nil {
		t.Fatal(err)
	}
	return created.Add(-time.Second).Format(time.RFC3339Nano)
}
func changedTicket(t *testing.T, fixture checkpointFixture, text string) {
	t.Helper()
	write(t, filepath.Join(fixture.root, "specs", "build demo", "tickets", "one.md"), text)
}

// checkpointReceiptFor builds the receipt assigned's current worktree state earns. A test
// that needs the receipt without spending it shares this one source with the checkpointing
// helper below rather than restating the receipt's shape.
func checkpointReceiptFor(t *testing.T, service *Service, assigned Assignment, ownership []string) receipt {
	t.Helper()
	run := loadRun(t, service)
	_, stored, ok := assignmentFor(run, assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	rows := make([]rowReceipt, len(stored.Rows))
	for i, row := range stored.Rows {
		rows[i] = rowReceipt{Row: row, Outcome: "passed"}
	}
	tree := benchgit.TreeHash(assigned.Path)
	return receipt{Version: receiptVersion, Run: run.Run, Assignment: assigned.ID, Base: assigned.Base, Tree: tree, TicketDigest: stored.TicketDigest, Rows: rows, Checks: []check{{Name: "go test ./internal/specbuild", Passed: true}}, Probe: probe{Producer: "coordinator", Assignment: assigned.ID, Tree: tree, Command: "go test ./internal/specbuild", Exit: 0, OutputDigest: digest("focused pass"), Produced: time.Now().UTC().Format(time.RFC3339Nano)}, Ownership: ownership, Assumptions: assumptionDigests(stored.Assumptions)}
}
func checkpointAssignment(t *testing.T, root string, service *Service, assigned Assignment, ownership []string) record {
	t.Helper()
	rec := checkpointReceiptFor(t, service, assigned, ownership)
	if _, err := service.Checkpoint(t.Context(), "build demo", assigned.ID, writeCheckpointReceipt(t, rec, "\n")); err != nil {
		t.Fatalf("Checkpoint(%s): %v", assigned.ID, err)
	}
	return loadRun(t, service)
}
func siblingCheckpoints(t *testing.T, firstPath, firstContent, secondPath, secondContent string) (string, *Service, Assignment, Assignment, record) {
	t.Helper()
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R17] sibling replay\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	first, _, err := service.Assign(t.Context(), "build demo", "one.md", "first sibling")
	if err != nil {
		t.Fatalf("Assign first: %v", err)
	}
	second, _, err := service.Assign(t.Context(), "build demo", "two.md", "second sibling")
	if err != nil {
		t.Fatalf("Assign second: %v", err)
	}
	for _, change := range []struct {
		assignment    Assignment
		path, content string
	}{{first, firstPath, firstContent}, {second, secondPath, secondContent}} {
		write(t, filepath.Join(change.assignment.Path, change.path), change.content)
		checkpointAssignment(t, root, service, change.assignment, []string{change.path})
	}
	return root, service, first, second, loadRun(t, service)
}
func loadRun(t *testing.T, service *Service) record {
	t.Helper()
	run, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load run: found:%v err:%v", found, err)
	}
	return run
}

func TestSharedPreconditionsRefuseEveryAssignmentOwnershipIdentity(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*testing.T, *preconditionFixture)
	}{
		{"request", func(t *testing.T, f *preconditionFixture) {
			updateRun(t, f, func(run *record) {
				assigned := run.Assignments[f.assignmentKey]
				assigned.OwnerRequest = "changed"
				run.Assignments[f.assignmentKey] = assigned
			})
		}},
		{"id", func(t *testing.T, f *preconditionFixture) {
			updateRun(t, f, func(run *record) {
				assigned := run.Assignments[f.assignmentKey]
				assigned.ID = "changed"
				run.Assignments[f.assignmentKey] = assigned
			})
		}},
		{"path", func(t *testing.T, f *preconditionFixture) {
			updateRun(t, f, func(run *record) {
				assigned := run.Assignments[f.assignmentKey]
				assigned.Path = filepath.Join(t.TempDir(), "moved")
				run.Assignments[f.assignmentKey] = assigned
			})
		}},
		{"common Git dir", func(t *testing.T, f *preconditionFixture) {
			other := repo(t)
			updateRun(t, f, func(run *record) {
				assigned := run.Assignments[f.assignmentKey]
				assigned.Path = other
				run.Assignments[f.assignmentKey] = assigned
			})
			owned, found, err := intent.FindAssignmentByRequest(f.root, f.run.Assignments[f.assignmentKey].OwnerRequest)
			if err != nil || !found {
				t.Fatalf("owned assignment: found:%v err:%v", found, err)
			}
			owned.Worktree = other
			if err := intent.PutAssignment(f.root, owned); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newPreconditionFixture(t, true)
			test.mutate(t, &fixture)
			before := snapshotPrecondition(t, fixture)
			if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assignment.ID, ""); err == nil {
				t.Fatal("Checkpoint accepted ownership drift")
			}
			if after := snapshotPrecondition(t, fixture); after != before {
				t.Fatalf("ownership refusal mutated: before=%#v after=%#v", before, after)
			}
		})
	}
}
func TestSharedPreconditionsRefuseSwappedSiblingOwnershipTuples(t *testing.T) {
	fixture := newPreconditionFixture(t, true)
	second, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "second sibling")
	if err != nil {
		t.Fatalf("Assign second: %v", err)
	}
	run := loadRun(t, fixture.service)
	firstKey, first, firstOK := assignmentFor(run, fixture.assignment.ID)
	secondKey, secondStored, secondOK := assignmentFor(run, second.ID)
	if !firstOK || !secondOK {
		t.Fatal("sibling assignments missing")
	}
	first.ID, first.Path, first.OwnerRequest, secondStored.ID, secondStored.Path, secondStored.OwnerRequest = secondStored.ID, secondStored.Path, secondStored.OwnerRequest, first.ID, first.Path, first.OwnerRequest
	run.Assignments[firstKey], run.Assignments[secondKey] = first, secondStored
	fixture.run = run
	saveRun(t, fixture.service, run)
	before := snapshotPrecondition(t, fixture)
	for _, id := range []string{fixture.assignment.ID, second.ID} {
		if _, err := fixture.service.Checkpoint(t.Context(), "build demo", id, ""); err == nil || !strings.Contains(err.Error(), "ownership") {
			t.Fatalf("Checkpoint(%s) = %v, want ownership refusal", id, err)
		}
	}
	if after := snapshotPrecondition(t, fixture); after != before {
		t.Fatalf("tuple-swap refusal mutated: before=%#v after=%#v", before, after)
	}
}
func TestAssignmentsScopeIdenticalRequestsToTheirRuns(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "second demo", "spec.md"), "# Second\n\nStatus: staged\n")
	write(t, filepath.Join(root, "specs", "second demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\n\n- [ ] [R01] separate request scope\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second spec")
	service := New(root, reuseGreenGate{}, &preconditionOwner{})
	for _, slug := range []string{"build demo", "second demo"} {
		if _, err := service.Start(t.Context(), slug); err != nil {
			t.Fatalf("Start(%s): %v", slug, err)
		}
	}
	first, _, err := service.Assign(t.Context(), "build demo", "one.md", "same request")
	if err != nil {
		t.Fatalf("Assign first: %v", err)
	}
	second, _, err := service.Assign(t.Context(), "second demo", "one.md", "same request")
	if err != nil || first.ID == second.ID {
		t.Fatalf("Assign second = %#v, %v", second, err)
	}
	for _, replay := range []struct{ slug, id string }{{"build demo", first.ID}, {"second demo", second.ID}} {
		assigned, _, err := service.Assign(t.Context(), replay.slug, "one.md", "same request")
		if err != nil || assigned.ID != replay.id {
			t.Fatalf("replay %s = %#v, %v", replay.slug, assigned, err)
		}
	}
	firstRun := loadRun(t, service)
	secondRun, _, _ := service.load("second demo")
	_, firstStored, _ := assignmentFor(firstRun, first.ID)
	_, secondStored, _ := assignmentFor(secondRun, second.ID)
	if firstStored.OwnerRequest == secondStored.OwnerRequest {
		t.Fatal("owner requests collided across runs")
	}
}
func TestCheckpointJournalRecoversRetainedRefAndRejectsDifferentReceipt(t *testing.T) {
	fixture := newCheckpointFixture(t)
	runner := &countingRunner{}
	fixture.service.runner = runner
	receipt := writeCheckpointReceipt(t, fixture.receipt, "\n")
	fixture.service.fault = func(point string) error {
		if point == "checkpoint/ref" {
			return errors.New("injected")
		}
		return nil
	}
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, receipt); err == nil {
		t.Fatal("fault did not interrupt checkpoint")
	}
	run := loadRun(t, fixture.service)
	ref := "refs/bench/specbuild/checkpoint/" + digest(run.Run+"\x00"+fixture.assigned.ID)
	retained := git(t, fixture.root, "rev-parse", ref)
	fixture.service.fault = nil
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, receipt); err != nil {
		t.Fatal(err)
	}
	if runner.calls == 0 {
		t.Fatal("checkpoint did not use runner")
	}
	if got := git(t, fixture.root, "rev-parse", ref); got != retained {
		t.Fatalf("checkpoint replay changed ref: %s != %s", got, retained)
	}
	fixture.receipt.Probe.Command = "different"
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err == nil || !strings.Contains(err.Error(), "conflicts") {
		t.Fatalf("different receipt = %v", err)
	}
}
func TestIntegrateCancellationKeepsPreparedReplayRecoverable(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	pidPath, grandPath := filepath.Join(t.TempDir(), "child"), filepath.Join(t.TempDir(), "grandchild")
	runner := &countingRunner{child: pidPath, grand: grandPath, block: true}
	fixture.service.runner = runner
	before, count := git(t, fixture.root, "rev-parse", fixture.run.Candidate), git(t, fixture.root, "rev-list", "--count", fixture.run.Candidate)
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() {
		_, err := fixture.service.Integrate(ctx, "build demo", fixture.assigned.ID)
		done <- err
	}()
	pids := []int{waitPID(t, pidPath), waitPID(t, grandPath)}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Integrate error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Integrate did not return after cancellation")
	}
	for _, pid := range pids {
		waitProcessExit(t, pid)
	}
	run := loadRun(t, fixture.service)
	op, found := fixture.service.operation(run, "integrate", fixture.assigned.ID)
	if !found || op.State != "prepared" || op.Result != "" || git(t, fixture.root, "rev-parse", run.Candidate) != before {
		t.Fatalf("canceled integration state = %#v", run)
	}
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate retry: %v", err)
	}
	if got := git(t, fixture.root, "rev-list", "--count", fixture.run.Candidate); got == count || runner.commits != 1 {
		t.Fatalf("retry candidate count=%s prior=%s commits=%d", got, count, runner.commits)
	}
}
func waitPID(t *testing.T, path string) int {
	t.Helper()
	for deadline := time.Now().Add(time.Second); time.Now().Before(deadline); {
		b, _ := os.ReadFile(path)
		if pid, err := strconv.Atoi(strings.TrimSpace(string(b))); err == nil && pid > 0 {
			return pid
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("process did not publish a positive PID: %s", path)
	return 0
}
func waitProcessExit(t *testing.T, pid int) {
	t.Helper()
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); {
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("process %d still exists after cancellation", pid)
}

func injectFault(want string) func(string) error {
	return func(got string) error {
		if got == want {
			return errors.New("injected")
		}
		return nil
	}
}

func (*abandonOwner) Create(context.Context, string, string, string, string) (OwnedWorktree, error) {
	return OwnedWorktree{}, errors.New("unexpected create")
}

type abandonSnapshot struct{ state, refs, worktrees string }

// The assign path is where an unanchored crossing costs something: a ticket
// that cannot write the artifact advertising what it changes is refused at
// lease time rather than discovered a review round later.
func TestAssignRefusesAContractCrossingNoFencedPath(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, ".gitignore"), "dist/\n")
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"),
		"# One\n\nOwnership fence: internal/specbuild\n"+
			"Contracts: every registry row's port name crosses registry→derived inventory\n"+
			"Assumptions: the crossing anchors nothing\n\n- [ ] [R10] unanchored crossing\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "unanchored ticket")
	service := New(root, &countingGate{}, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	assigned, _, err := service.Assign(context.Background(), "build demo", "one.md", "unanchored")
	if err == nil {
		t.Fatalf("Assign leased %#v for a crossing naming no fenced path", assigned)
	}
	const want = "spec build ticket one.md declares a contract crossing no path in its ownership fence"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err, want)
	}
}
