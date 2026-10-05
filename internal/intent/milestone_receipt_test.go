package intent

import (
	"strings"
	"testing"
)

// A milestone receipt needs one line-safe identity and a payload, and its identity is
// unique. The strict read refuses the whole ledger over one invalid receipt.
func TestMilestoneReceiptValidation(t *testing.T) {
	valid := `{"id":"sha256:one","payload":"{}"}`
	for _, row := range []struct{ name, receipts string }{
		{"empty-identity", `{"id":"","payload":"{}"}`},
		{"empty-payload", `{"id":"sha256:one","payload":""}`},
		{"control-identity", `{"id":"sha256:one\u0007","payload":"{}"}`},
		{"duplicate-identity", valid + "," + valid},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := newRepo(t)
			writeLedgerBody(t, root, `{"schema":2,"entries":[],"milestone_receipts":[`+row.receipts+`]}`)
			if _, err := Read(root); err == nil || !strings.Contains(err.Error(), "invalid milestone receipt") {
				t.Fatalf("Read = %v, want the invalid milestone receipt refusal", err)
			}
		})
	}
	root := newRepo(t)
	writeLedgerBody(t, root, `{"schema":2,"entries":[],"milestone_receipts":[`+valid+`]}`)
	ledger, err := Read(root)
	if err != nil || len(ledger.MilestoneReceipts) != 1 || ledger.MilestoneReceipts[0].ID != "sha256:one" {
		t.Fatalf("Read = %+v, %v; want the one valid receipt", ledger.MilestoneReceipts, err)
	}
}

// The purge carries milestone receipts through untouched while it drops a debris record.
func TestPurgeAssignmentsKeepsMilestoneReceipts(t *testing.T) {
	root := newRepo(t)
	writeLedgerBody(t, root, `{"schema":2,"entries":[],"assignments":[{"schema":"bench-assignment/v1","id":"cafe0000000000000000000000000001","state":"provisional"}],"milestone_receipts":[{"id":"sha256:one","payload":"{}"}]}`)
	if dropped, err := PurgeAssignments(root, keepAll); err != nil || dropped != 1 {
		t.Fatalf("PurgeAssignments = %d, %v; want 1 dropped", dropped, err)
	}
	ledger, err := Read(root)
	if err != nil || len(ledger.MilestoneReceipts) != 1 || ledger.MilestoneReceipts[0].Payload != "{}" {
		t.Fatalf("receipts after purge = %+v, %v; want the receipt kept", ledger.MilestoneReceipts, err)
	}
}
