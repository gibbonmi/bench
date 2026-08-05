package specbuild

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestOperationJournalAdmitsSixtyFifthEntry(t *testing.T) {
	fixture := newCheckpointFixture(t)
	run := loadRun(t, fixture.service)
	fillOperationJournal(&run, 64)
	mustSaveRun(t, fixture.service, run)

	op, completed, err := fixture.service.beginOperation(&run, "journal-probe", "request-64", "input-64")
	if err != nil || completed || op.State != "prepared" || len(run.Operations) != 65 {
		t.Fatalf("65th operation = %#v, completed=%v, operations=%d, err=%v", op, completed, len(run.Operations), err)
	}
	stored := loadRun(t, fixture.service)
	if len(stored.Operations) != 65 || stored.Operations[operationID("journal-probe", "request-64")].State != "prepared" {
		t.Fatalf("stored 65th operation = %#v", stored.Operations)
	}
}

func TestOperationJournalRetainsFiniteBoundAndExistingReplay(t *testing.T) {
	fixture := newCheckpointFixture(t)
	run := loadRun(t, fixture.service)
	fillOperationJournal(&run, operationLimit)
	mustSaveRun(t, fixture.service, run)
	statePath, err := fixture.service.statePath(run.Slug)
	if err != nil {
		t.Fatalf("state path: %v", err)
	}
	before, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("read state: %v", err)
	}

	if _, _, err := fixture.service.beginOperation(&run, "journal-probe", "overflow", "overflow-input"); err == nil || !strings.Contains(err.Error(), "journal is full") {
		t.Fatalf("operation beyond bound = %v, want full-journal refusal", err)
	}
	after, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("read state after refusal: %v", err)
	}
	if string(after) != string(before) || len(run.Operations) != operationLimit {
		t.Fatal("full-journal refusal mutated durable or in-memory state")
	}

	op, completed, err := fixture.service.beginOperation(&run, "journal-probe", "request-0", "input-0")
	if err != nil || !completed || op.State != "completed" {
		t.Fatalf("existing operation replay = %#v, completed=%v, err=%v", op, completed, err)
	}
}

func fillOperationJournal(run *record, count int) {
	run.Operations = make(map[string]operation, count)
	for i := 0; i < count; i++ {
		request := fmt.Sprintf("request-%d", i)
		op := operation{Command: "journal-probe", Request: request, Input: digest(fmt.Sprintf("input-%d", i)), State: "completed"}
		run.Operations[operationID(op.Command, op.Request)] = op
	}
}

func mustSaveRun(t *testing.T, service *Service, run record) {
	t.Helper()
	if err := service.save(run); err != nil {
		t.Fatalf("save run: %v", err)
	}
}
