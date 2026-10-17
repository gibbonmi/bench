package ownership_test

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/tickets"
)

// TestKitCustomRootRow grades that a verified kit source locates a planned change of a custom
// grammar root at that root's own registry row, not at the ordinary parser's row.
func TestKitCustomRootRow(t *testing.T) {
	lines, err := tickets.BindingLines(ownershiptest.KitRegistry(t))
	custom := tickets.GrammarRoot{Dir: filepath.ToSlash(filepath.Dir(operationsFile)), Name: "operations"}
	if err != nil || lines[custom.String()] < 1 || lines[custom.String()] == lines[tickets.GrammarParser().String()] {
		t.Fatalf("kit registry locates no distinct row for %s: %v", custom, err)
	}
	files := map[string]string{operationsFile: grammarBase[operationsFile], breakdownFile: grammarBase[breakdownFile]}
	out := mustPropose(t, ownershiptest.ProposalRepo{Writes: []string{breakdownFile}, Files: files, OwnInventory: true, Seal: freshness.Publish,
		Tail: plannedForm(operationsFile, "evidencecmd.operations", "var")}.Seed(t))
	rows := matching(t, out, "closure_evidence", map[string]string{"ticket": "one.md", "path": rootHelp})
	if !slices.ContainsFunc(rows, func(row map[string]any) bool { return cell(row, "line") == fmt.Sprint(lines[custom.String()]) }) {
		t.Fatalf("evidence for %s = %v, want a reason at line %d", rootHelp, rows, lines[custom.String()])
	}
}

// withoutParserRow is the kit's binding registry without the line that holds the ordinary
// parser's row, at the line that the registry's own line owner locates.
func withoutParserRow(t *testing.T) string {
	t.Helper()
	registry := ownershiptest.KitRegistry(t)
	lines, err := tickets.BindingLines(registry)
	line := lines[tickets.GrammarParser().String()]
	if err != nil || line < 1 {
		t.Fatalf("kit registry locates no parser row: %v", err)
	}
	return ownershiptest.Cut(t, string(registry), strings.SplitAfter(string(registry), "\n")[line-1])
}
