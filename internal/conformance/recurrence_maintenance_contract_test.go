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
	if entryCount != 1 || !strings.Contains(entry, "`context.schema = 4`") ||
		!strings.Contains(entry, "every other schema stops the phase before any batch mutation") ||
		!strings.Contains(entry, "Do not guess recurrence facts from an older schema") {
		diags = append(diags, "bench-what-next does not require schema 4 before recurrence maintenance")
	}
	if entryCount != 1 || !strings.Contains(entry, "schema-4 index is the complete local inventory") ||
		!strings.Contains(entry, "The index proves what exists; targeted fetches and named body reads prove content") ||
		!strings.Contains(entry, "Fetch complete roadmap detail only for rows the reconcile touches") ||
		!strings.Contains(entry, "Read idea, learning, and retro bodies from the paths the index names") {
		diags = append(diags, "bench-what-next does not use index-first evidence for recurrence maintenance")
	}
	schemaAt := strings.Index(entry, "`context.schema = 4`")
	trustAt := strings.Index(entry, "Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that snapshot.")
	trustStop := "snapshot. When it is false, stop before any batch mutation"
	if entryCount != 1 || trustAt <= schemaAt || !strings.Contains(entry, trustStop) {
		diags = append(diags, "bench-what-next does not stop untrusted recurrence evidence before reconciliation")
	}
	if entryCount != 1 || trustAt <= schemaAt || !strings.Contains(entry, "`occurrence_discrepancies` row together with the complete index snapshot") {
		diags = append(diags, "bench-what-next does not preserve complete index evidence before reconciliation")
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
		{"schema", "`context.schema = 4`", "`context.schema = 3`", "bench-what-next does not require schema 4 before recurrence maintenance"},
		{"schema guessing", "Do not guess recurrence facts from an older schema.", "Infer recurrence facts from an older schema.", "bench-what-next does not require schema 4 before recurrence maintenance"},
		{"index inventory", "schema-4\nindex is the complete local inventory", "schema-4\nindex is a partial inventory", "bench-what-next does not use index-first evidence for recurrence maintenance"},
		{"targeted detail", "Fetch complete roadmap detail only for rows the reconcile touches", "Fetch complete roadmap detail for every row", "bench-what-next does not use index-first evidence for recurrence maintenance"},
		{"named capture bodies", "Read idea, learning, and retro bodies\nfrom the paths the index names", "Read capture bodies from a directory sweep", "bench-what-next does not use index-first evidence for recurrence maintenance"},
		{"trust after reconcile", "Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that\nsnapshot.", "After `## 1. Reconcile first`, read `context.sequence_trusted` from that\nsnapshot.", "bench-what-next does not stop untrusted recurrence evidence before reconciliation"},
		{"trust polarity", "snapshot. When it is false, stop before any batch mutation", "snapshot. When it is true, stop before any batch mutation", "bench-what-next does not stop untrusted recurrence evidence before reconciliation"},
		{"trust stop", "stop before any batch mutation", "continue into batch mutation", "bench-what-next does not stop untrusted recurrence evidence before reconciliation"},
		{"context evidence", "`occurrence_discrepancies` row together with the complete index snapshot", "`occurrence_discrepancies` row", "bench-what-next does not preserve complete index evidence before reconciliation"},
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

func TestOccurrenceLedgerMigrationCheckBitesOnFT158Count(t *testing.T) {
	root := t.TempDir()
	data, err := os.ReadFile(filepath.Join(NewHarness(t).KitRoot, "ROADMAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	mutated := strings.Replace(string(data), "Occurrences: baseline-01, baseline-02, baseline-03", "Occurrences: baseline-01, baseline-02", 1)
	if mutated == string(data) {
		t.Fatal("FT158 migration-count mutation did not change its ledger")
	}
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte(mutated), 0o644); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, diag := range checkOccurrenceLedgerMigration(root) {
		if diag == "ROADMAP.md occurrence-ledger migration count for FT158 is wrong" {
			found = true
		}
	}
	if !found {
		t.Fatal("FT158 count mutation passed occurrence-ledger migration check")
	}
}

func TestOccurrenceLedgerMigrationAllowsRetiredFT126(t *testing.T) {
	root := t.TempDir()
	const migrated = "**FT71 — migration baseline.**\nOccurrences: baseline-01\n\n" +
		"**FT158 — migration baseline.**\nOccurrences: baseline-01, baseline-02, baseline-03\n\n" +
		"**FT128 — migration baseline.**\nOccurrences: baseline-01\n\n" +
		"**FT98 — migration baseline.**\nOccurrences: baseline-01, baseline-02, baseline-03\n\n" +
		"**FT169 — migration baseline.**\nOccurrences: baseline-01\n\n" +
		"**FT141 — migration baseline.**\nOccurrences: baseline-01\n\n" +
		"**FT94 — migration baseline.**\nOccurrences: baseline-01\n\n" +
		"**FT125 — migration baseline.**\nOccurrences: baseline-01\n"
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte(migrated), 0o644); err != nil {
		t.Fatal(err)
	}
	if containsDiagnostic(checkOccurrenceLedgerMigration(root), "ROADMAP.md occurrence-ledger migration count for FT126 is wrong") {
		t.Fatal("retired FT126 remained required by occurrence-ledger migration")
	}
}
