package staleness_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// symbolModule is the module of the cited-symbol fixture.
const symbolModule = "example.com/sym"

// symbolPremise is one current symbol premise of the fixture inventory, whose expected kind
// and signature repeat its owner identity.
func symbolPremise(t *testing.T, slug, id, path, pkg, symbol, kind, signature string) string {
	t.Helper()
	owner, err := json.Marshal(map[string]string{"path": path, "package": symbolModule + "/" + pkg, "symbol": symbol, "kind": kind, "signature": signature})
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(`{"id":%q,%s,"kind":"symbol","state":"current","subject":%s,"expected":{"kind":%q,"signature":%q}}`, id, itemLocation(slug), owner, kind, signature)
}

// TestSymbolIdentity covers SS13, SS14, SS15, SS16, SS17, SS18, and SS19. One committed
// pair changes each cited symbol in its own way, and the actual audit-spec response grades
// them all. Each expected drift row states the declaration syntax of both sides by hand, so
// a collector that resolves a name to the wrong owner, or that compares no receiver or
// signature, loses its row.
func TestSymbolIdentity(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	const sym, moved = "internal/example/sym.go", "internal/example/moved.go"
	premise := func(id, path, symbol, kind, signature string) string {
		return symbolPremise(t, slug, id, path, "internal/example", symbol, kind, signature)
	}
	premises := []string{
		premise("gone", sym, "example.Gone", "func", "func Gone()"),
		premise("change", sym, "example.Change", "func", "func Change(a int) int"),
		premise("receiver", sym, "example.Origin.Method", "method", "func (Origin) Method()"),
		premise("move", sym, "example.Move", "func", "func Move()"),
		premise("rule", sym, "example.Rule", "const", ""),
		premise("fresh", sym, "example.Fresh", "func", "func Fresh()"),
		premise("bare", sym, "Dup", "func", "func Dup()"),
		premise("alias", sym, "example.Alias", "type", "type Origin struct{ N int }"),
		premise("tagged", "internal/example/tagged.go", "example.Tagged", "func", "func Tagged()"),
		symbolPremise(t, slug, "prose", "docs/notes.md", "docs", "docs.Notes", "func", "func Notes()"),
	}
	write(t, map[string]string{
		specPath(slug):                inventorySpec(slug, premises, nil),
		"go.mod":                      "module " + symbolModule + "\n\ngo 1.25\n",
		"internal/example/sub/sub.go": "package sub\n\nfunc Dup() {}\n",
		"internal/example/tagged.go":  "//go:build benchnever\n\npackage example\n\nfunc Tagged() {}\n",
		"docs/notes.md":               "# Notes\n\nfunc Notes()\n",
		sym: "package example\n\ntype Origin struct{ N int }\n\ntype Alias = Origin\n\nfunc (Origin) Method() {}\n\n" +
			"func Gone() {}\n\nfunc Change(a int) int { return a }\n\nfunc Move() {}\n\nconst Rule = 1\n\nfunc Dup() {}\n",
	})
	base := commitAll(t, "cited symbols")
	write(t, map[string]string{
		moved: "package example\n\nfunc Move() {}\n",
		sym: "package example\n\ntype Origin struct{ N int }\n\ntype Alias = Origin\n\nfunc (*Origin) Method() {}\n\n" +
			"func Change(a string) int { return len(a) }\n\nvar Rule = 1\n\nfunc Dup() {}\n\nfunc Fresh() {}\n",
	})
	tip := commitAll(t, "cited symbols change")

	out, code := audit(t, slug, base, tip)
	if code != 1 || !strings.Contains(out, "cited-symbols,red") {
		t.Fatalf("audit = (%d), want exit 1 and a red cited-symbols row:\n%s", code, out)
	}
	symbolRows := map[string]bool{}
	for _, row := range driftRows(t, out) {
		if strings.Contains(row, "|symbol ") {
			symbolRows[row] = true
		}
	}
	diagnostics := map[string]string{}
	for _, row := range table(t, out, "audit_diagnostics") {
		diagnostics[fmt.Sprint(row["subject"])] = fmt.Sprint(row["reason"])
	}
	rows := 0
	for _, test := range []struct{ row, premise, drift, diagnostic string }{
		{"SS13", "gone", sym + "|" + absent + "|symbol example.Gone absent|func Gone()|" + absent, ""},
		{"SS14", "change", sym + "|" + sym + "|symbol example.Change signature|func Change(a int) int|func Change(a string) int", ""},
		{"SS14", "receiver", sym + "|" + sym + "|symbol example.Origin.Method signature|func (Origin) Method()|func (*Origin) Method()", ""},
		{"SS15", "move", sym + "|" + moved + "|symbol example.Move moved|func Move()|func Move()", ""},
		{"SS16", "rule", sym + "|" + sym + "|symbol example.Rule kind|const|var Rule = 1", ""},
		{"SS13", "fresh", absent + "|" + sym + "|symbol example.Fresh absent|" + absent + "|func Fresh()", ""},
		{"SS17", "bare", "", "cited symbol Dup is ambiguous"},
		{"SS18", "alias", "", ""},
		{"SS19", "tagged", "", "cited symbol example.Tagged is unsupported-context"},
		{"SS19", "prose", "", "cited symbol docs.Notes is unread"},
	} {
		if test.drift != "" {
			rows++
		}
		t.Run(test.row+" "+test.premise, func(t *testing.T) {
			if test.drift != "" && !symbolRows[test.drift] {
				t.Errorf("drift lacks %q; symbol rows %v", test.drift, symbolRows)
			}
			if got := diagnostics[test.premise]; test.diagnostic == "" && got != "" || !strings.HasPrefix(got, test.diagnostic) {
				t.Errorf("diagnostic of %s = %q, want prefix %q", test.premise, got, test.diagnostic)
			}
		})
	}
	// The alias names its origin's declaration, so a lexical reading of it adds a row.
	if len(symbolRows) != rows {
		t.Errorf("symbol drift rows = %d, want %d: %v", len(symbolRows), rows, symbolRows)
	}
}
