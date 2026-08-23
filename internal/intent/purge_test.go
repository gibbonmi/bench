package intent

import (
	"os"
	"path/filepath"
	"testing"
)

// writeLedgerBody puts body at the ledger address verbatim. This is how a record no
// current writer would produce, such as one an older binary wrote, reaches the purge.
func writeLedgerBody(t *testing.T, root, body string) string {
	t.Helper()
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func keepAll(Assignment) bool { return true }

// TestPurgeAssignmentsDropsRecordsReadRefuses holds the purge to the one thing Read
// cannot do. A record this build's strict decoder rejects makes the whole ledger
// unreadable. The reconcile that clears it has to see past the rejection rather than
// inherit it. Everything the decoder does accept is left for the caller's predicate.
func TestPurgeAssignmentsDropsRecordsReadRefuses(t *testing.T) {
	root := newRepo(t)
	valid := activeAssignment()
	if err := PutAssignment(root, valid); err != nil {
		t.Fatal(err)
	}
	body := `{"schema":2,"entries":[],"assignments":[` +
		`{"schema":"bench-assignment/v1","id":"` + valid.ID + `","owner_id":"` + valid.OwnerID + `","request":"` + valid.Request +
		`","label":"delegate","start":"` + valid.Start + `","branch":"` + valid.Branch + `","worktree":"/pool/delegate","state":"active","recovery":[]},` +
		`{"schema":"bench-assignment/v1","id":"cafe0000000000000000000000000001","state":"provisional","ticket":"specs/removed/tickets/one.md"}` +
		`],"cleanup_receipts":[]}`
	path := writeLedgerBody(t, root, body)
	if _, err := Read(root); err == nil {
		t.Fatal("Read accepted the legacy record, so the purge has nothing to see past")
	}

	dropped, err := PurgeAssignments(root, keepAll)
	if err != nil || dropped != 1 {
		t.Fatalf("PurgeAssignments = %d, %v; want 1 dropped", dropped, err)
	}
	ledger, err := Read(root)
	if err != nil || len(ledger.Assignments) != 1 || ledger.Assignments[0].ID != valid.ID {
		t.Fatalf("ledger after purge = %#v, %v", ledger.Assignments, err)
	}

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if dropped, err := PurgeAssignments(root, keepAll); err != nil || dropped != 0 {
		t.Fatalf("second purge = %d, %v; want 0 dropped", dropped, err)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != string(before) {
		t.Fatalf("a purge that dropped nothing rewrote the ledger: %v", err)
	}
}

// TestPurgeAssignmentsDropsEveryLegacySchemaRecord pins the schema rule the strict reader
// already states: the legacy schema authorizes no lifecycle record at all. The purge
// drops them without asking the caller's predicate, which would otherwise keep a record
// the reader refuses to admit.
func TestPurgeAssignmentsDropsEveryLegacySchemaRecord(t *testing.T) {
	root := newRepo(t)
	valid := activeAssignment()
	body := `{"schema":1,"entries":[],"assignments":[` +
		`{"schema":"bench-assignment/v1","id":"` + valid.ID + `","owner_id":"` + valid.OwnerID + `","request":"` + valid.Request +
		`","label":"delegate","start":"` + valid.Start + `","branch":"` + valid.Branch + `","worktree":"/pool/delegate","state":"active","recovery":[]}` +
		`]}`
	writeLedgerBody(t, root, body)

	dropped, err := PurgeAssignments(root, keepAll)
	if err != nil || dropped != 1 {
		t.Fatalf("PurgeAssignments = %d, %v; want 1 dropped", dropped, err)
	}
	ledger, err := Read(root)
	if err != nil || len(ledger.Assignments) != 0 || ledger.Schema != Schema {
		t.Fatalf("ledger after legacy purge = %#v, %v", ledger, err)
	}
}

// TestPurgeAssignmentsIsANoOpWithoutALedger keeps the reconcile silent in a repository
// that has never written one. It runs at every session start, and creating a ledger
// there would manufacture the state it exists to clean up.
func TestPurgeAssignmentsIsANoOpWithoutALedger(t *testing.T) {
	root := newRepo(t)
	dropped, err := PurgeAssignments(root, keepAll)
	if err != nil || dropped != 0 {
		t.Fatalf("PurgeAssignments = %d, %v; want 0 dropped", dropped, err)
	}
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("purge created a ledger at %s: %v", filepath.Base(path), err)
	}
}
