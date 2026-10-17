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

// kitFiles is each named crossBase file inside a sealed kit's module, so each of its
// in-module imports names a package of the kit.
func kitFiles(paths ...string) map[string]string {
	files := map[string]string{}
	for _, path := range paths {
		files[path] = strings.ReplaceAll(crossBase[path], "github.com/gibbonmi/bench", ownershiptest.KitModule)
	}
	return files
}

// kitGauge is the gauge package inside a sealed kit's module, beside the parser package that
// the module's package index needs, and kitGaugeTip changes its grammar.
var (
	kitGauge    = kitFiles(gaugeFile, usageFile)
	kitGaugeTip = map[string]string{gaugeFile: strings.Replace(kitGauge[gaugeFile], "bench gauge", "bench gauge [--all]", 1)}
)

// TestKitRootRowSides grades that a verified kit source locates each grammar change at the row
// that registered its root on the side of the range that declares that root: a root that only
// the base declares, and a root of another package that reaches the changed declaration.
func TestKitRootRowSides(t *testing.T) {
	kitShared := kitFiles(sharedFile, namesGrammar, usageFile)
	cases := []effectCase{
		{name: "a verified kit source locates a root that only the base declares", writes: []string{gaugeFile}, seal: freshness.Publish, files: kitGauge,
			tip:      map[string]string{gaugeFile: "package gauge\n\nfunc Command(args []string) int { return 0 }\n"},
			want:     map[string]bool{rootHelp: true},
			evidence: map[string]string{rootHelp: "grammar root internal/gauge.Command#1 requires " + rootHelp + " beside " + gaugeFile},
			located:  true},
		{name: "a verified kit source locates another package's passing root", writes: []string{namesGrammar}, seal: freshness.Publish, files: kitShared,
			tip:      map[string]string{namesGrammar: strings.Replace(kitShared[namesGrammar], "bench shared", "bench shared [--all]", 1)},
			want:     map[string]bool{sharedFile: true, rootHelp: true},
			evidence: map[string]string{rootHelp: "grammar root internal/harnesses.Shared#1 requires " + rootHelp + " beside " + namesGrammar},
			located:  true},
	}
	for _, test := range cases {
		t.Run(test.name, test.run)
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
