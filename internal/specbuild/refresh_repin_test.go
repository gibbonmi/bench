package specbuild

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// newDocsRepairFixture stages a docs-repair ticket fenced only to the
// consumer's own ticket file, walks it through the ordinary lifecycle with a
// zero-byte checkpoint — the ticket's job is naming the metadata fix, not
// touching source — and integrates it, leaving the consumer's original
// assignment preserved and uncheckpointed: exactly the shape a rewrite of the
// consumer's ticket text must reach the preserved assignment through.
func newDocsRepairFixture(t *testing.T) (string, *Service, Assignment) {
	t.Helper()
	root := repo(t)
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	consumer, _, err := service.Assign(t.Context(), "build demo", "one.md", consumerRequest)
	if err != nil {
		t.Fatalf("Assign consumer: %v", err)
	}
	repair := "# Repair one metadata\n\n" +
		"Blocked by: none\n" +
		"Ownership fence: `specs/build demo/tickets/one.md`\n" +
		"Integration surfaces: repaired dependency metadata→one.md + R90\n" +
		"Contracts: repaired ticket metadata crosses `specs/build demo/tickets/one.md`→the refreshed consumer assignment, asserted by R90 against the committed ticket\n" +
		"Closure: R90/metadata-fix\n\n" +
		"- [ ] [R90] repair one.md's stale dependency metadata\n\n" +
		"## Red mutations\n\n" +
		"| criterion | mutation | owner | operation sequence |\n" +
		"|---|---|---|---|\n" +
		"| R90/metadata-fix | omit the metadata repair | the metadata regression | apply the omission, run the metadata regression, require red |\n"
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "repair.md"), repair)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "stage repair ticket")
	// The mid-run insertion refuses exactly like the other refresh fixtures'
	// repair ticket does, until promote recomposes onto the new commit.
	if _, _, err := service.Assign(t.Context(), "build demo", "repair.md", "repair request"); !errors.Is(err, errRecompose) {
		t.Fatalf("mid-run repair insertion = %v, want recomposition route", err)
	}
	if _, err := service.Promote(t.Context(), "build demo"); err != nil {
		t.Fatalf("recomposing Promote: %v", err)
	}
	repairAssignment, _, err := service.Assign(t.Context(), "build demo", "repair.md", "repair request")
	if err != nil {
		t.Fatalf("Assign repair: %v", err)
	}
	checkpointAssignment(t, root, service, repairAssignment, nil)
	if _, err := service.Integrate(t.Context(), "build demo", repairAssignment.ID); err != nil {
		t.Fatalf("Integrate repair: %v", err)
	}
	return root, service, consumer
}

// rewriteConsumerTicket commits text as the new tickets/one.md and recomposes
// the run onto it. The rewrite is staged directly to the working tree — the
// same route the fixture above uses for repair.md itself — rather than
// through the repair assignment's own checkpoint: ParseTicket and
// requireCommittedTicket read the ticket from the main working tree, so a
// rewrite that existed only on the integrated candidate would leave that read
// point stale. Committing it directly here keeps the read point valid without
// changing where Refresh or checkpoint/integrate read from.
func rewriteConsumerTicket(t *testing.T, root string, service *Service, text string) {
	t.Helper()
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), text)
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "repair consumer ticket metadata")
	if _, err := service.Promote(t.Context(), "build demo"); err != nil {
		t.Fatalf("recomposing Promote for ticket rewrite: %v", err)
	}
}

func debugReceiptForDocsRepair(run record, consumer Assignment) debugReceipt {
	return debugReceipt{
		Version: debugReceiptVersion, Run: run.Run, Assignment: consumer.ID, Base: consumer.Base,
		Repro:         debugRepro{Command: "go test ./internal/contract/runtime -run TestRuntimeCommitContracts", Exit: 1, OutputDigest: digest("deterministic red"), Produced: time.Now().UTC().Format(time.RFC3339Nano)},
		Cause:         "dependency metadata repair landed mid-run",
		RequiredFence: []string{"specs/build demo/tickets/one.md"},
		Resumable:     true,
	}
}

// TestRefreshRePinsTicketAfterDocsRepairRewrite is the regression for the
// architecture defect: a legitimate docs-repair rewrite of the consumer's
// ticket, naming the reciprocal Blocked by: edge and leaving the ownership
// fence untouched, must reach the preserved assignment — refresh re-pins the
// recorded ticket digest and rows instead of refusing forever on the stale
// one recorded at assign.
func TestRefreshRePinsTicketAfterDocsRepairRewrite(t *testing.T) {
	root, service, consumer := newDocsRepairFixture(t)
	rewritten := "# One\n\n" +
		"Blocked by: repair.md\n" +
		"Ownership fence: `internal/specbuild`\n" +
		"Integration surfaces: repaired dependency metadata consumed→internal/specbuild + R01\n" +
		"Contracts: repaired dependency metadata crosses `internal/specbuild`→the refreshed consumer assignment, asserted by R01 against the real inventory\n" +
		"Closure: R01/metadata-consumed\n\n" +
		"- [ ] [R01] exact start\n\n" +
		"## Red mutations\n\n" +
		"| criterion | mutation | owner | operation sequence |\n" +
		"|---|---|---|---|\n" +
		"| R01/metadata-consumed | ignore the repaired metadata | the start regression | apply the omission, run the start regression, require red |\n"
	rewriteConsumerTicket(t, root, service, rewritten)

	run := loadRun(t, service)
	refreshed, _, err := service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, debugReceiptForDocsRepair(run, consumer)))
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if refreshed.Base != run.CandidateTip {
		t.Fatalf("refreshed base = %s, want repaired candidate %s", refreshed.Base, run.CandidateTip)
	}
	newTicket, err := ParseTicket(run.Spec, "one.md")
	if err != nil {
		t.Fatalf("ParseTicket: %v", err)
	}
	reloaded := loadRun(t, service)
	_, stored, ok := assignmentFor(reloaded, consumer.ID)
	if !ok {
		t.Fatal("missing consumer assignment")
	}
	if stored.TicketDigest != newTicket.Digest {
		t.Fatalf("reloaded assignment ticket digest = %s, want re-pinned %s", stored.TicketDigest, newTicket.Digest)
	}
	if !sameStrings(stored.Rows, newTicket.Rows) {
		t.Fatalf("reloaded assignment rows = %v, want re-pinned %v", stored.Rows, newTicket.Rows)
	}
}

// TestRefreshRefusesTicketRewriteThatChangesOwnershipFence is the boundary on
// the re-pin: a rewrite that still names the reciprocal Blocked by: edge but
// moves the consumer's ownership fence is not a metadata repair refresh may
// carry — the fence is the assignment's write envelope, and refresh does not
// own moving it.
func TestRefreshRefusesTicketRewriteThatChangesOwnershipFence(t *testing.T) {
	root, service, consumer := newDocsRepairFixture(t)
	rewritten := "# One\n\n" +
		"Blocked by: repair.md\n" +
		"Ownership fence: `internal/other`\n" +
		"Integration surfaces: repaired dependency metadata consumed→internal/other + R01\n" +
		"Contracts: none crosses\n\n" +
		"- [ ] [R01] exact start\n"
	rewriteConsumerTicket(t, root, service, rewritten)

	run := loadRun(t, service)
	_, _, err := service.Refresh(t.Context(), "build demo", "one.md", consumerRequest, writeDebugReceipt(t, debugReceiptForDocsRepair(run, consumer)))
	if err == nil {
		t.Fatal("refresh carried a rewrite that moved the assignment's ownership fence")
	}
	if !strings.Contains(err.Error(), "Ownership fence") {
		t.Fatalf("refresh refusal = %v, want the changed Ownership fence named", err)
	}
}
