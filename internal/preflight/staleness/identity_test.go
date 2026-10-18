package staleness_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// symbolModule is the module of the cited-symbol fixture.
const symbolModule = "example.com/sym"

// overBound is a Go comment line longer than the control-record bound, so a committed object
// that holds it does not read.
func overBound() string { return "// " + strings.Repeat("x", int(bounds.ControlRecordLimit)) + "\n\n" }

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
// signature, loses its row. The kind words come from consumers, the one source that both
// the typed and the syntax readings spell. The breakdown over the same pair gives each
// premise verdict, which an ambiguous bare spelling does not fail.
func TestSymbolIdentity(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	const sym, moved, broken = "internal/example/sym.go", "internal/example/moved.go", "internal/example/broken.go"
	premise := func(id, path, symbol, kind, signature string) string {
		return symbolPremise(t, slug, id, path, "internal/example", symbol, kind, signature)
	}
	premises := []string{
		premise("gone", sym, "example.Gone", consumers.KindFunc, "func Gone()"),
		premise("change", sym, "example.Change", consumers.KindFunc, "func Change(a int) int"),
		premise("receiver", sym, "example.Origin.Method", consumers.KindMethod, "func (Origin) Method()"),
		premise("kept", sym, "example.Origin.Kept", consumers.KindMethod, "func (Origin) Kept()"),
		premise("steady", sym, "example.Steady", consumers.KindFunc, "func Steady()"),
		premise("move", sym, "example.Move", consumers.KindFunc, "func Move()"),
		premise("rule", sym, "example.Rule", consumers.KindConst, ""),
		premise("fresh", sym, "example.Fresh", consumers.KindFunc, "func Fresh()"),
		premise("bare", sym, "Dup", consumers.KindFunc, "func Dup()"),
		premise("alias", sym, "example.Alias", consumers.KindType, "type Origin struct{ N int }"),
		premise("broken", broken, "example.Broken", consumers.KindFunc, "func Broken()"),
		premise("tagged", "internal/example/tagged.go", "example.Tagged", consumers.KindFunc, "func Tagged()"),
		symbolPremise(t, slug, "prose", "docs/notes.md", "docs", "docs.Notes", consumers.KindFunc, "func Notes()"),
	}
	write(t, map[string]string{
		specPath(slug):                inventorySpec(slug, premises, nil),
		"go.mod":                      "module " + symbolModule + "\n\ngo 1.25\n",
		"internal/example/sub/sub.go": "package sub\n\nfunc Dup() {}\n",
		"internal/example/tagged.go":  "//go:build benchnever\n\npackage example\n\nfunc Tagged() {}\n",
		"docs/notes.md":               "# Notes\n\nfunc Notes()\n",
		broken:                        "package example\n\nfunc Broken( {\n",
		sym: "package example\n\ntype Origin struct{ N int }\n\ntype Alias = Origin\n\nfunc (Origin) Method() {}\n\nfunc (Origin) Kept() {}\n\nfunc Steady() {}\n\n" +
			"func Gone() {}\n\nfunc Change(a int) int { return a }\n\nfunc Move() {}\n\nconst Rule = 1\n\nfunc Dup() {}\n",
	})
	base := commitAll(t, "cited symbols")
	write(t, map[string]string{
		moved:  "package example\n\nfunc Move() {}\n",
		broken: "package example\n\nfunc Broken() {}\n",
		sym: "package example\n\ntype Origin struct{ N int }\n\ntype Alias = Origin\n\nfunc (*Origin) Method() {}\n\nfunc (Origin) Kept() {}\n\nfunc Steady() {}\n\n" +
			"func Change(a string) int { return len(a) }\n\nvar Rule = 1\n\nfunc Dup() {}\n\nfunc Fresh() {}\n",
	})
	tip := commitAll(t, "cited symbols change")
	breakdown, code := preflight.CommandWithVersion("staleness-test")([]string{"build", slug, "--breakdown", "--base", base, "--source-tip", tip})
	if code != 0 {
		t.Fatalf("breakdown = (%d), want exit 0:\n%s", code, breakdown)
	}
	verdicts := map[string]string{}
	for _, row := range table(t, breakdown, "premises") {
		verdicts[fmt.Sprint(row["id"])] = fmt.Sprint(row["state"])
	}

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
	for _, test := range []struct{ row, premise, verdict, drift, diagnostic string }{
		{"SS13", "gone", "unknown", sym + "|" + absent + "|symbol example.Gone absent|func Gone()|" + absent, ""},
		{"SS14", "change", "unknown", sym + "|" + sym + "|symbol example.Change signature|func Change(a int) int|func Change(a string) int", ""},
		{"SS14", "receiver", "unknown", sym + "|" + sym + "|symbol example.Origin.Method signature|func (Origin) Method()|func (*Origin) Method()", ""},
		// An unchanged method or function reads as the same kind on both sides, so it adds
		// no row.
		{"SS16", "kept", "green", "", ""},
		{"SS16", "steady", "green", "", ""},
		{"SS15", "move", "unknown", sym + "|" + moved + "|symbol example.Move moved|func Move()|func Move()", ""},
		{"SS16", "rule", "unknown", sym + "|" + sym + "|symbol example.Rule kind|" + consumers.KindConst + "|var Rule = 1", ""},
		{"SS13", "fresh", "green", absent + "|" + sym + "|symbol example.Fresh absent|" + absent + "|func Fresh()", ""},
		// The package filter selects the cited owner, so the premise passes; the bare
		// spelling alone proves no owner, so the audit keeps it diagnostic.
		{"SS17", "bare", "green", "", "cited symbol Dup is ambiguous"},
		{"SS18", "alias", "green", "", ""},
		// A baseline owner that does not parse stays unread, never absent.
		{"SS19", "broken", "green", "", "cited symbol example.Broken is unread"},
		{"SS19", "tagged", "unknown", "", "cited symbol example.Tagged is unsupported-context"},
		{"SS19", "prose", "unknown", "", "cited symbol docs.Notes is unread"},
	} {
		if test.drift != "" {
			rows++
		}
		t.Run(test.row+" "+test.premise, func(t *testing.T) {
			if got := verdicts[test.premise]; got != test.verdict {
				t.Errorf("premise %s = %q, want %q", test.premise, got, test.verdict)
			}
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

	// A baseline owner object over the control-record bound does not read, so the audit
	// refuses with no staleness result.
	const big = "internal/example/big.go"
	write(t, map[string]string{
		specPath(slug): inventorySpec(slug, append(premises, premise("big", big, "example.Big", consumers.KindFunc, "func Big()")), nil),
		big:            "package example\n\n" + overBound() + "func Big() {}\n",
	})
	oversized := commitAll(t, "oversized owner")
	write(t, map[string]string{big: "package example\n\nfunc Big() {}\n"})
	if out, code := audit(t, slug, oversized, commitAll(t, "owner shrinks")); code != 1 || !strings.Contains(out, "cited symbol source refused") || strings.Contains(out, "staleness[") {
		t.Errorf("oversized baseline owner = (%d), want the cited symbol source refusal and no staleness result:\n%s", code, out)
	}
}

// TestTestIdentity covers SS20, SS21, SS22, SS23, SS24, SS25, SS26, and SS27. One committed
// pair changes each test seam in its own way, and the actual audit-spec response grades
// them all. Each expected drift row and diagnostic is stated here by hand, with object IDs
// from git rev-parse. The ordinary preflight over the same pair is the continuation, which
// must not repeat the initial planned-test absence. A test object that does not read at
// either identity refuses the audit.
func TestTestIdentity(t *testing.T) {
	_, slug := preflighttest.SeedAdmitted(t)
	const (
		kept, gone, removed = "internal/example/kept_test.go", "internal/example/gone_test.go", "internal/example/removed_test.go"
		future, computed    = "internal/example/future_test.go", "internal/example/computed_test.go"
		tagged, untagged    = "internal/example/tagged_test.go", "internal/example/untagged_test.go"
	)
	seam := func(id, path, function, state string) string {
		return seamItem(slug, id, path, "example.com/internal/example", function, state, "PF1")
	}
	const header = "package example\n\nimport \"testing\"\n\n"
	seams := []string{
		seam("kept", kept, "TestKept", "existing"),
		seam("gone", gone, "TestGone", "existing"),
		seam("removed", removed, "TestRemoved", "existing"),
		seam("future", future, "TestFuture", "planned"),
		seam("early", future, "TestEarly", "planned"),
		seam("subtest", kept, "TestSub/case one", "existing"),
		seam("nosub", kept, "TestSub/case two", "existing"),
		seam("deep", kept, "TestSub/case one/inner", "existing"),
		seam("computed", computed, "TestComputed/x", "existing"),
		seam("tagged", tagged, "TestTagged", "existing"),
		seam("untagged", untagged, "TestUntagged", "existing"),
	}
	write(t, map[string]string{
		specPath(slug): inventorySpec(slug, nil, seams),
		kept:           header + "func TestKept(t *testing.T) {}\n\nfunc TestSub(t *testing.T) {\n\tt.Run(\"case one\", func(t *testing.T) {})\n}\n",
		gone:           header + "func TestGone(t *testing.T) {}\n",
		removed:        header + "func TestRemoved(t *testing.T) {}\n\nfunc TestStays(t *testing.T) {}\n",
		future:         header + "func TestOther(t *testing.T) {}\n",
		computed:       header + "func TestComputed(t *testing.T) {\n\tname := \"x\"\n\tt.Run(name, func(t *testing.T) {})\n}\n",
		tagged:         "//go:build benchnever\n\n" + header + "func TestTagged(t *testing.T) {}\n",
		untagged:       "//go:build benchnever\n\n" + header + "func TestUntagged(t *testing.T) {}\n",
	})
	base := commitAll(t, "cited tests")
	preflighttest.RunGit(t, "rm", "-q", gone)
	write(t, map[string]string{
		removed:  header + "func TestStays(t *testing.T) {}\n",
		future:   header + "func TestOther(t *testing.T) {}\n\nfunc TestEarly(t *testing.T) {}\n",
		untagged: header + "func TestStays(t *testing.T) {}\n",
	})
	tip := commitAll(t, "cited tests change")

	out, code := audit(t, slug, base, tip)
	if code != 1 || !strings.Contains(out, "cited-tests,red") {
		t.Fatalf("audit = (%d), want exit 1 and a red cited-tests row:\n%s", code, out)
	}
	testRows := map[string]bool{}
	for _, row := range driftRows(t, out) {
		if strings.Contains(row, "|test ") {
			testRows[row] = true
		}
	}
	diagnostics := map[string]string{}
	for _, row := range table(t, out, "audit_diagnostics") {
		diagnostics[fmt.Sprint(row["subject"])] = fmt.Sprint(row["reason"])
	}
	same := func(path, kind string) string {
		return path + "|" + path + "|" + kind + "|" + object(t, base, path) + "|" + object(t, tip, path)
	}
	rows := 0
	for _, test := range []struct{ row, seam, drift, diagnostic string }{
		{"SS20", "gone", gone + "|" + absent + "|test TestGone file-absent|" + object(t, base, gone) + "|" + absent, ""},
		{"SS21", "removed", same(removed, "test TestRemoved function-absent"), ""},
		{"SS22", "future", "", ""},
		{"SS23", "early", same(future, "test TestEarly already-shipped"), ""},
		{"SS24", "kept", "", ""},
		{"SS25", "subtest", "", ""},
		{"SS25", "nosub", same(kept, "test TestSub/case two function-absent"), ""},
		{"SS26", "computed", "", "cited test TestComputed/x is unsupported"},
		// A declared first segment cannot vouch for a deeper one, which needs the call graph.
		{"SS26", "deep", "", "cited test TestSub/case one/inner is unsupported"},
		{"SS27", "tagged", "", "cited test TestTagged is unsupported-context"},
		// A wrong context at the base cannot mask the removal at the tip.
		{"SS27", "untagged", same(untagged, "test TestUntagged function-absent"), ""},
	} {
		if test.drift != "" {
			rows++
		}
		t.Run(test.row+" "+test.seam, func(t *testing.T) {
			if test.drift != "" && !testRows[test.drift] {
				t.Errorf("drift lacks %q; test rows %v", test.drift, testRows)
			}
			if got := diagnostics[test.seam]; test.diagnostic == "" && got != "" || !strings.HasPrefix(got, test.diagnostic) {
				t.Errorf("diagnostic of %s = %q, want prefix %q", test.seam, got, test.diagnostic)
			}
		})
	}
	// A planned function in a file that holds other tests, and a declared subtest, add no row.
	if len(testRows) != rows {
		t.Errorf("test drift rows = %d, want %d: %v", len(testRows), rows, testRows)
	}

	// The continuation: the ordinary preflight over the pair keeps its own verdict, so a
	// predecessor's legitimate planned test leaves it usable.
	if out, code := preflight.CommandWithVersion("staleness-test")([]string{"build", slug, "--base", base, "--source-tip", tip}); code != 0 || strings.Contains(out, "already-shipped") {
		t.Errorf("ordinary preflight = (%d), want exit 0 and no initial planned-test absence:\n%s", code, out)
	}

	// A test object over the control-record bound at either identity does not read, so the
	// audit refuses with no staleness result. Each pair ends at the head, which the audit
	// grades as its tip, so only the named side is oversized.
	const big = "internal/example/big_test.go"
	small, large := header+"func TestBig(t *testing.T) {}\n", header+overBound()+"func TestBig(t *testing.T) {}\n"
	write(t, map[string]string{specPath(slug): inventorySpec(slug, nil, append(seams, seam("big", big, "TestBig", "existing"))), big: large})
	before := commitAll(t, "oversized test")
	for _, step := range []struct{ side, body string }{{"baseline", small}, {"tip", large}} {
		write(t, map[string]string{big: step.body})
		after := commitAll(t, "test changes size")
		if out, code := audit(t, slug, before, after); code != 1 || !strings.Contains(out, "cited test source refused") || strings.Contains(out, "staleness[") {
			t.Errorf("oversized %s test = (%d), want the cited test source refusal and no staleness result:\n%s", step.side, code, out)
		}
		before = after
	}
}
