package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func checkOccurrenceLedgerAndMaintenance(root string) []string {
	return append(checkOccurrenceLedgerMigration(root), checkRecurrenceMaintenanceContract(root)...)
}

func checkRecurrenceMaintenanceContract(root string) []string {
	path := filepath.Join(root, ".agents", "commands", "bench-what-next.md")
	text := stripHTMLComments(readIfExists(path))
	if strings.TrimSpace(text) == "" {
		return []string{"bench-what-next recurrence maintenance command is unavailable"}
	}

	entry, entryCount := markdownH2Sections(text, "Entry orientation")
	occurrences, occurrenceCount := markdownH2Sections(text, "2. Drain occurrence evidence")
	sequence, sequenceCount := markdownH2Sections(text, "7. Refresh the sequence")
	entry = collapseSpace(entry)
	occurrences = collapseSpace(occurrences)
	sequence = collapseSpace(sequence)
	var diags []string
	if entryCount != 1 || !strings.Contains(entry, "`context.schema = 3`") ||
		!strings.Contains(entry, "every other schema stops the phase before any batch mutation") ||
		!strings.Contains(entry, "Do not guess recurrence facts from an older schema") {
		diags = append(diags, "bench-what-next does not require schema 3 before recurrence maintenance")
	}
	schemaAt := strings.Index(entry, "`context.schema = 3`")
	trustAt := strings.Index(entry, "Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that snapshot.")
	trustStop := "snapshot. When it is false, stop before any batch mutation"
	if entryCount != 1 || trustAt <= schemaAt || !strings.Contains(entry, trustStop) {
		diags = append(diags, "bench-what-next does not stop untrusted recurrence evidence before reconciliation")
	}
	if entryCount != 1 || trustAt <= schemaAt || !strings.Contains(entry, "`occurrence_discrepancies` row together with the complete context snapshot") {
		diags = append(diags, "bench-what-next does not preserve complete context evidence before reconciliation")
	}
	pending := "add its incident key to that owner's `Occurrences:` line in `ROADMAP.md` before removing any source unit"
	if occurrenceCount != 1 || !strings.Contains(occurrences, "For every `pending` owner/incident pair") || !strings.Contains(occurrences, pending) {
		diags = append(diags, "bench-what-next does not ledger every pending occurrence before source removal")
	}
	if occurrenceCount != 1 || !strings.Contains(occurrences, "Every `already-recorded` source already has that key: remove its source unit without adding another key") {
		diags = append(diags, "bench-what-next does not remove already-recorded occurrence sources without another key")
	}

	precedence := []string{
		"Rank rows by severity.",
		"choose actionable work over blocked work",
		"apply literal dependencies, then explicit reviewer pricing.",
		"rank by descending occurrence count.",
		"existing reproduced defect-over-feature rule, then cheapest-first cost rule.",
	}
	previous := -1
	if sequenceCount != 1 {
		diags = append(diags, "bench-what-next recurrence sequence precedence is unavailable")
	} else {
		for _, anchor := range precedence {
			at := strings.Index(sequence, anchor)
			if at < 0 || at <= previous {
				diags = append(diags, "bench-what-next recurrence sequence precedence is incomplete or out of order")
				break
			}
			previous = at
		}
	}
	return diags
}

func TestRecurrenceMaintenanceContractCheckBites(t *testing.T) {
	h := NewHarness(t)
	command := readIfExists(filepath.Join(h.KitRoot, ".agents", "commands", "bench-what-next.md"))
	if command == "" {
		t.Fatal("bench-what-next command unavailable")
	}

	mutations := []struct {
		name, old, new, want string
	}{
		{"schema", "`context.schema = 3`", "`context.schema = 2`", "bench-what-next does not require schema 3 before recurrence maintenance"},
		{"schema guessing", "Do not guess recurrence facts from an older schema.", "Infer recurrence facts from an older schema.", "bench-what-next does not require schema 3 before recurrence maintenance"},
		{"trust after reconcile", "Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that\nsnapshot.", "After `## 1. Reconcile first`, read `context.sequence_trusted` from that\nsnapshot.", "bench-what-next does not stop untrusted recurrence evidence before reconciliation"},
		{"trust polarity", "snapshot. When it is false, stop before any batch mutation", "snapshot. When it is true, stop before any batch mutation", "bench-what-next does not stop untrusted recurrence evidence before reconciliation"},
		{"trust stop", "stop before any batch mutation", "continue into batch mutation", "bench-what-next does not stop untrusted recurrence evidence before reconciliation"},
		{"context evidence", "`occurrence_discrepancies` row together with the complete context snapshot", "`occurrence_discrepancies` row", "bench-what-next does not preserve complete context evidence before reconciliation"},
		{"pending ledger", "add its incident\nkey to that owner's `Occurrences:` line in `ROADMAP.md` before removing any source\nunit", "record the pending pair for later", "bench-what-next does not ledger every pending occurrence before source removal"},
		{"already recorded", "remove its\nsource unit without adding another key", "remove its source unit after adding another key", "bench-what-next does not remove already-recorded occurrence sources without another key"},
		{"equal class", "apply literal dependencies, then explicit\nreviewer pricing. Only when all four stronger inputs tie, rank by descending\noccurrence count.", "rank by descending occurrence count before explicit reviewer pricing.", "bench-what-next recurrence sequence precedence is incomplete or out of order"},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			if strings.Count(command, mutation.old) != 1 {
				t.Fatalf("mutation anchor %q occurs %d times", mutation.old, strings.Count(command, mutation.old))
			}
			root := t.TempDir()
			path := filepath.Join(root, ".agents", "commands", "bench-what-next.md")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(strings.Replace(command, mutation.old, mutation.new, 1)), 0o644); err != nil {
				t.Fatal(err)
			}
			if !containsDiagnostic(checkRecurrenceMaintenanceContract(root), mutation.want) {
				t.Fatalf("mutation did not bite: want %q", mutation.want)
			}
		})
	}
}
