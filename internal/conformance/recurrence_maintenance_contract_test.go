package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

func checkOccurrenceLedgerAndMaintenance(root string) []string {
	return append(checkOccurrenceLedgerMigration(root), checkRecurrenceMaintenanceContract(root)...)
}

func checkRecurrenceMaintenanceContract(root string) []string {
	path := filepath.Join(root, ".agents", "commands", "bench-drain.md")
	text := stripHTMLComments(readIfExists(path))
	if strings.TrimSpace(text) == "" {
		return []string{"bench-drain recurrence maintenance command is unavailable"}
	}

	worktrees, worktreeCount := markdownH2Sections(text, "Worktree closure")
	var diags []string
	if worktreeCount != 1 || !strings.Contains(collapseSpace(worktrees), "Run `bench worktree list` before any batch mutation") {
		diags = append(diags, "bench-drain does not run bench worktree list before any batch mutation")
	}
	if !strings.Contains(collapseSpace(text), "After approval and the tracked landing, the coordinator runs `bench capture drain commit <drain-id>`.") {
		diags = append(diags, "bench-drain does not delay sealed-generation retirement until approval and landing")
	}
	return diags
}

func TestRecurrenceMaintenanceContractCheckBites(t *testing.T) {
	h := NewHarness(t)
	command := readIfExists(filepath.Join(h.KitRoot, ".agents", "commands", "bench-drain.md"))
	if command == "" {
		t.Fatal("bench-drain command unavailable")
	}

	mutations := []struct {
		name, old, new, want string
	}{
		{"orphan inventory", "Run `bench worktree list` before any batch mutation", "Run `bench worktree list` after the batch mutation", "bench-drain does not run bench worktree list before any batch mutation"},
		{"sealed retirement after approval", "After approval and the tracked landing, the coordinator runs\n`bench capture drain commit <drain-id>`.", "Before approval, retire the sealed generation.", "bench-drain does not delay sealed-generation retirement until approval and landing"},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			if baseline := checkRecurrenceMaintenanceContract(h.KitRoot); len(baseline) != 0 {
				t.Fatalf("live command has baseline diagnostics before mutation: %s", strings.Join(baseline, "; "))
			}
			if strings.Count(command, mutation.old) != 1 {
				t.Fatalf("mutation anchor %q occurs %d times", mutation.old, strings.Count(command, mutation.old))
			}
			root := t.TempDir()
			path := filepath.Join(root, ".agents", "commands", "bench-drain.md")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(strings.Replace(command, mutation.old, mutation.new, 1)), 0o644); err != nil {
				t.Fatal(err)
			}
			if !containsDiagnostic(checkRecurrenceMaintenanceContract(root), mutation.want) {
				t.Fatalf("mutation did not bite: want %q", mutation.want)
			}
			t.Logf("observed red: %s", mutation.want)
		})
	}
}

func TestOccurrenceLedgerMigrationCheckBitesOnFT98Count(t *testing.T) {
	root := t.TempDir()
	kit := NewHarness(t).KitRoot
	index, err := os.ReadFile(filepath.Join(kit, "ROADMAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	// The mutation drops a key from the ledger, which lives in the row's own detail file.
	// The fixture copies the live board's index beside a mutated roadmap/FT98.md.
	rowFile := filepath.Join(kit, "roadmap", "FT98.md")
	row, err := os.ReadFile(rowFile)
	if err != nil {
		t.Fatal(err)
	}
	mutated := strings.Replace(string(row), ", baseline-02\n", "\n", 1)
	if mutated == string(row) {
		t.Fatal("FT98 migration-count mutation did not change its ledger")
	}
	roadmaptest.WriteSplitBoard(t, root, string(index), map[string]string{"FT98.md": mutated})
	found := false
	for _, diag := range checkOccurrenceLedgerMigration(root) {
		if diag == "ROADMAP.md occurrence-ledger migration count for FT98 is wrong" {
			found = true
		}
	}
	if !found {
		t.Fatal("FT98 count mutation passed occurrence-ledger migration check")
	}
}

// A retired migration row leaves the board with its spec. The check must stop requiring
// its count, or every retirement lands red.
func TestOccurrenceLedgerMigrationAllowsRetiredRows(t *testing.T) {
	root := t.TempDir()
	index, files := "", map[string]string{}
	for _, row := range [][2]string{
		{"FT71", "baseline-01"},
		{"FT128", "baseline-01"},
		{"FT98", "baseline-01, baseline-02, baseline-03"},
		{"FT169", "baseline-01"},
		{"FT141", "baseline-01"},
		{"FT94", "baseline-01"},
		{"FT125", "baseline-01"},
	} {
		heading := "**" + row[0] + " — migration baseline.**"
		index += heading + "\n\n"
		files[row[0]+".md"] = heading + "\nOccurrences: " + row[1] + "\n"
	}
	roadmaptest.WriteSplitBoard(t, root, index, files)
	diags := checkOccurrenceLedgerMigration(root)
	for _, retired := range []string{"FT126", "FT158"} {
		if containsDiagnostic(diags, "ROADMAP.md occurrence-ledger migration count for "+retired+" is wrong") {
			t.Fatalf("retired %s remained required by occurrence-ledger migration", retired)
		}
	}
}
