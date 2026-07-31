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
