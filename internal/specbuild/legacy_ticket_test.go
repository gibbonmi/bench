package specbuild

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// requireLegacyTicketLine fails unless the fixture ticket still carries the
// retired grammar's line. Both tolerance tests below are vacuous the moment the
// fixture is rewritten to the current grammar, and neither would say so.
func requireLegacyTicketLine(t *testing.T, root string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, "specs", "build demo", "tickets", "one.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "\nAssumptions:") {
		t.Fatal("fixture ticket no longer carries the retired Assumptions line")
	}
}

// A ticket staged under an earlier grammar carries a build from lease to
// integration with the retired line skipped rather than refused.
func TestLegacyAssumptionsLineParsesAndIntegrates(t *testing.T) {
	fixture := newCheckpointFixture(t)
	requireLegacyTicketLine(t, fixture.root)

	ticket, err := ParseTicket(filepath.Join(fixture.root, "specs", "build demo", "spec.md"), "one.md")
	if err != nil {
		t.Fatalf("ParseTicket over a legacy ticket: %v", err)
	}
	if want := []string{"internal/specbuild"}; !slices.Equal(ticket.Fence, want) {
		t.Errorf("Fence = %q, want %q", ticket.Fence, want)
	}
	if want := []string{"R10", "R11", "R12", "R13", "R14", "R15", "R54"}; !slices.Equal(ticket.Rows, want) {
		t.Errorf("Rows = %q, want %q", ticket.Rows, want)
	}
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate: %v", err)
	}
}

// A record written before the retirement carries assumption data no current
// structure names, and the lifecycle continues from it across a fresh reload.
func TestPreRetirementRecordReloadContinuesTheLifecycle(t *testing.T) {
	fixture := newCheckpointFixture(t)
	requireLegacyTicketLine(t, fixture.root)
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	path, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	const anchor = `"Fence":["internal/specbuild"]`
	if !strings.Contains(string(data), anchor) {
		t.Fatalf("persisted assignment does not carry %s: %s", anchor, data)
	}
	legacy := strings.Replace(string(data), anchor, anchor+`,"Assumptions":["receipt contract"]`, 1)
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}

	reloaded := New(fixture.root, fixture.gate, realOwner{})
	if _, err := reloaded.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate after reloading a pre-retirement record: %v", err)
	}
}
