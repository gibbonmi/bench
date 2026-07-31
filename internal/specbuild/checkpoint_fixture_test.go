package specbuild

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type checkpointFixture struct {
	root     string
	service  *Service
	gate     *countingGate
	assigned Assignment
	run      record
	receipt  receipt
}

func newCheckpointFixture(t *testing.T) checkpointFixture {
	t.Helper()
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\nAssumptions: receipt contract\n\n- [ ] [R10-R15] checkpoint receipt\n- [ ] [R54] framing\n")
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
	git(t, assigned.Path, "add", ".")
	git(t, assigned.Path, "commit", "-qm", "checkpoint change")
	run, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load: found:%v err:%v", found, err)
	}
	_, stored, ok := assignmentFor(run, assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	rows := make([]rowReceipt, len(stored.Rows))
	for i, row := range stored.Rows {
		rows[i] = rowReceipt{Row: row, Outcome: "passed"}
	}
	tree := git(t, assigned.Path, "rev-parse", "HEAD^{tree}")
	return checkpointFixture{
		root: root, service: service, gate: gate, assigned: assigned, run: run,
		receipt: receipt{Version: receiptVersion, Run: run.Run, Assignment: assigned.ID, Base: assigned.Base, Tree: tree, TicketDigest: stored.TicketDigest, Rows: rows, Checks: []check{{Name: "go test ./internal/specbuild", Passed: true}}, Probe: probe{Producer: "coordinator", Assignment: assigned.ID, Tree: tree, Command: "go test ./internal/specbuild", Exit: 0, OutputDigest: digest("focused pass"), Produced: time.Now().UTC().Format(time.RFC3339Nano)}, Ownership: []string{"internal/specbuild/checkpoint-change.go"}, Assumptions: assumptionDigests(stored.Assumptions)},
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

func checkpointAssignment(t *testing.T, root string, service *Service, assigned Assignment, ownership []string) record {
	t.Helper()
	run, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load before checkpoint: found:%v err:%v", found, err)
	}
	_, stored, ok := assignmentFor(run, assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	rows := make([]rowReceipt, len(stored.Rows))
	for i, row := range stored.Rows {
		rows[i] = rowReceipt{Row: row, Outcome: "passed"}
	}
	tree := git(t, assigned.Path, "rev-parse", "HEAD^{tree}")
	rec := receipt{Version: receiptVersion, Run: run.Run, Assignment: assigned.ID, Base: assigned.Base, Tree: tree, TicketDigest: stored.TicketDigest, Rows: rows, Checks: []check{{Name: "go test ./internal/specbuild", Passed: true}}, Probe: probe{Producer: "coordinator", Assignment: assigned.ID, Tree: tree, Command: "go test ./internal/specbuild", Exit: 0, OutputDigest: digest("focused pass"), Produced: time.Now().UTC().Format(time.RFC3339Nano)}, Ownership: ownership, Assumptions: assumptionDigests(stored.Assumptions)}
	if _, err := service.Checkpoint(t.Context(), "build demo", assigned.ID, writeCheckpointReceipt(t, rec, "\n")); err != nil {
		t.Fatalf("Checkpoint(%s): %v", assigned.ID, err)
	}
	run, found, err = service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load after checkpoint: found:%v err:%v", found, err)
	}
	return run
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
		git(t, change.assignment.Path, "add", ".")
		git(t, change.assignment.Path, "commit", "-qm", change.path)
		checkpointAssignment(t, root, service, change.assignment, []string{change.path})
	}
	run, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load sibling checkpoints: found:%v err:%v", found, err)
	}
	return root, service, first, second, run
}
