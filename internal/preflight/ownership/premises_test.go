package ownership_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestMetadataGrammar drives each planning-facts grammar case through the actual proposal.
// A diagnostic case keeps the inventory incomplete, so the premises table names it; a
// valid inventory reports its own items instead.
func TestMetadataGrammar(t *testing.T) {
	cases := []struct {
		name, text string
		id, state  string
		reason     string
	}{
		{"OC55 duplicate fence", fenced(completeEmpty) + fenced(completeEmpty), "metadata", "unknown", "more than once"},
		{"OC56 duplicate key", fenced(strings.Replace(completeEmpty, `"version":1,`, `"version":1,"version":1,`, 1)), "metadata", "unknown", `duplicate object field "version"`},
		{"OC56 nested duplicate key", fenced(withPremises(strings.Replace(linesPremise("p1"), `"expected":{"count":1}`, `"expected":{"count":1,"count":2}`, 1))), "metadata", "unknown", `duplicate object field "count"`},
		{"OC57 unknown field", fenced(strings.Replace(completeEmpty, `"version":1,`, `"version":1,"extra":true,`, 1)), "metadata", "unknown", `unknown field "extra"`},
		{"OC58 duplicate id", fenced(withPremises(linesPremise("p1"), linesPremise("p1"))), "metadata", "unknown", `duplicate item id "p1"`},
		{"OC58 cross-array duplicate id", fenced(strings.Replace(withPremises(linesPremise("p1")), `"scope":[]`, `"scope":[{"id":"p1",`+itemLocation+`,"paths":["internal/example"]}]`, 1)), "metadata", "unknown", `duplicate item id "p1"`},
		{"OC59 unsupported version", fenced(strings.Replace(completeEmpty, `"version":1`, `"version":2`, 1)), "metadata", "unknown", "unsupported metadata version 2"},
		{"OC60 invalid enum", fenced(strings.Replace(completeEmpty, `"declared-complete"`, `"Declared-Complete"`, 1)), "metadata", "unknown", `invalid completeness "Declared-Complete"`},
		{"OC91 opener inside another fence", "\n```text\n```bench-planning-facts\n" + completeEmpty + "\n```\n```\n", "metadata", "unknown", "no bench-planning-facts fence"},
		{"unclosed fence", "\n```bench-planning-facts\n" + completeEmpty + "\n", "metadata", "unknown", "not closed"},
		{"nested fence", "\n```bench-planning-facts\n" + completeEmpty + "\n```json\n```\n", "metadata", "unknown", "another fence marker"},
		{"missing array", fenced(strings.Replace(completeEmpty, `,"manual_claims":[]`, "", 1)), "metadata", "unknown", `missing field "manual_claims"`},
		{"stale current premise stays unproved", fenced(withPremises(strings.Replace(linesPremise("p1"), `"count":1`, `"count":2`, 1))), "p1", "unknown", "lines"},
		{"current-premise manual claim", fenced(manualClaim(`"applicability":"current-premise"`)), "semantic", "manual", "current-premise"},
		{"CRLF fence section", strings.ReplaceAll(fenced(manualClaim(`"applicability":"current-premise"`)), "\n", "\r\n"), "semantic", "manual", "current-premise"},
		{"current-premise rejects native fields", fenced(manualClaim(`"applicability":"current-premise","row":"PF1"`)), "metadata", "unknown", `unknown field "row"`},
		{"native-future requires its descriptor", fenced(manualClaim(`"applicability":"native-future","row":"PF1","chunk":"c1","ticket":"one.md","verification":"tests"`)), "metadata", "unknown", `missing field "phase_owner"`},
		{"native-future stays manual", fenced(manualClaim(`"applicability":"native-future","row":"PF1","chunk":"c1","ticket":"one.md","verification":"tests","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":1,"needle":"# Source","source":"tip"},"observation":"unobserved"`)), "semantic", "manual", "native-future"},
		{"native-future observation is unobserved", fenced(manualClaim(`"applicability":"native-future","row":"PF1","chunk":"c1","ticket":"one.md","verification":"tests","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":1,"needle":"# Source","source":"tip"},"observation":"observed"`)), "metadata", "unknown", `invalid observation "observed"`},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			out := mustPropose(t, emptyTree(test.text).Seed(t))
			row := one(t, out, "premises", map[string]string{"id": test.id})
			if cell(row, "state") != test.state || !strings.Contains(cell(row, "reason"), test.reason) {
				t.Fatalf("premise %s = %v, want state %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
	for name, test := range map[string]struct {
		text     string
		complete bool
	}{
		"declared-complete empty inventory": {fenced(completeEmpty), true},
		"partial inventory":                 {fenced(strings.Replace(completeEmpty, `"declared-complete"`, `"partial"`, 1)), false},
		"manual claim":                      {fenced(manualClaim(`"applicability":"current-premise"`)), false},
	} {
		t.Run("Complete "+name, func(t *testing.T) {
			snapshot, err := ownership.Collect(t.TempDir(), ownership.Source{}, ownership.ParseManifest([]byte(test.text)), nil)
			if err != nil || snapshot.Complete != test.complete {
				t.Fatalf("Collect = Complete %t, %v; want Complete %t", snapshot.Complete, err, test.complete)
			}
		})
	}
}

// TestMetadataGrammarFenceRoles is the parity case of the scanner TestMetadataGrammar reads:
// the fence roles join each scanned line while every other field keeps its value. A
// labeled marker inside another fence closes that fence and opens none.
func TestMetadataGrammarFenceRoles(t *testing.T) {
	scan := maps.FieldScan{
		Table:     []maps.FieldSpec{{Name: "Owner", Syntax: "Owner: "}},
		Duplicate: func(spec maps.FieldSpec, _ string) string { return "duplicate " + spec.Name },
	}
	lines, diagnostics := scan.Scan([]byte("Owner: a\n```text\nOwner: b\n```bench-planning-facts\nOwner: c\n```\n"))
	want := []maps.FieldLine{
		{Text: "Owner: a", Field: "Owner", Value: "a"},
		{Text: "```text", Fenced: true, FenceRole: maps.FenceOpen, FenceLabel: "text"},
		{Text: "Owner: b", Fenced: true},
		{Text: "```bench-planning-facts", Fenced: true, FenceRole: maps.FenceClose, FenceLabel: "bench-planning-facts"},
		{Text: "Owner: c", Field: "Owner", Value: "c", Diagnostic: "duplicate Owner"},
		{Text: "```", Fenced: true, FenceRole: maps.FenceOpen},
		{Text: "", Fenced: true},
	}
	if !reflect.DeepEqual(lines, want) || !reflect.DeepEqual(diagnostics, []string{"duplicate Owner"}) {
		t.Fatalf("Scan = %#v, %#v; want %#v and one duplicate", lines, diagnostics, want)
	}
}

// TestSizingPremises drives the line-count, glob-member, and headroom collectors through
// the actual proposal. The tip adds one physical line to a.go, so a count names its source.
// Each unproved form keeps its premise unknown with the reason the case names.
func TestSizingPremises(t *testing.T) {
	lines := func(id, path, source string, count int) string {
		return premiseItem(id, "lines", `{"path":"`+path+`","source":"`+source+`"}`, fmt.Sprintf(`{"count":%d}`, count))
	}
	const (
		hidden, a1, literal, visible = `"internal/glob/.hidden.go"`, `"internal/glob/a1.go"`, `"internal/glob/a[1].go"`, `"internal/glob/visible.go"`
		plannedLines                 = `{"id":"lines-planned",` + itemLocation + `,"kind":"lines","state":"planned","subject":{"path":"internal/sized/future.go","source":"tip"},"expected":{"count":40}}`
	)
	facts := withItems("premises",
		lines("lines-stale", "internal/sized/a.go", "tip", 1),
		lines("lines-tip", "internal/sized/a.go", "tip", 2),
		lines("lines-base", "internal/sized/a.go", "base", 1),
		lines("lines-unbuilt", "internal/sized/future.go", "tip", 40),
		plannedLines,
		globItem("glob-visible-only", "internal/glob/*.go", a1+","+literal+","+visible),
		globItem("glob-complete", "internal/glob/*.go", hidden+","+a1+","+literal+","+visible),
		globItem("glob-literal", `internal/glob/a\\[1\\].go`, literal),
		globItem("glob-class", "internal/glob/a[1].go", literal),
		globItem("glob-unsorted", "internal/glob/*.go", visible+","+hidden+","+a1+","+literal),
		globItem("glob-recursive", "internal/glob/**/*.go", ""),
		globItem("glob-escaping", "internal/sized/*.go", ""),
		globItem("glob-malformed", "internal/glob/[", ""),
		headroomItem("headroom-limit", "internal/sized", "", "", 3, 3, 0),
		headroomItem("headroom-copied-cap", "internal/sized", "", "", 12, 3, 0),
		headroomItem("headroom-swap", "internal/sized", `"internal/sized/d.go"`, `"internal/sized/b.go"`, 3, 3, 0),
		headroomItem("headroom-ignored-deletion", "internal/sized", `"internal/sized/d.go"`, `"internal/sized/b.go"`, 3, 3, 1),
		headroomItem("headroom-over", "internal/sized", `"internal/sized/d.go"`, "", 3, 3, 1),
		headroomItem("headroom-absent-deletion", "internal/sized", "", `"internal/sized/z.go"`, 3, 3, -1),
		headroomItem("headroom-granted", "internal/granted", `"internal/granted/n.go"`, "", 1, 1, 1))
	root := ownershiptest.ProposalRepo{Writes: []string{"internal/sized/a.go"}, Files: sizingBase, Tail: facts}.Seed(t)
	preflighttest.MustWriteFile(t, "internal/sized/a.go", "package sized\n\n")
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "tip")
	out := mustPropose(t, root)
	for _, test := range []struct{ row, id, state, reason string }{
		{"OC32", "lines-stale", "unknown", "holds 2 line(s) at the source tip"},
		{"OC32", "lines-tip", "green", "holds 2 line(s) at the source tip"},
		{"OC32", "lines-base", "green", "holds 1 line(s) at the source base"},
		{"OC33", "glob-visible-only", "unknown", "omits internal/glob/.hidden.go"},
		{"OC33", "glob-complete", "green", "4 tracked member(s)"},
		{"OC34", "glob-literal", "green", "1 tracked member(s)"},
		{"OC34", "glob-class", "unknown", "omits internal/glob/a1.go; names internal/glob/a[1].go, which is no member"},
		{"OC35", "glob-recursive", "unknown", "recursive"},
		{"OC35", "glob-escaping", "unknown", "escapes its scope"},
		{"OC35", "glob-malformed", "unknown", "syntax"},
		{"OC35", "glob-unsorted", "unknown", "complete sorted"},
		{"OC36", "headroom-limit", "green", "limit 3"},
		{"OC36", "headroom-copied-cap", "unknown", "maximum 12, the canonical limit is 3"},
		{"OC36", "headroom-over", "unknown", "exceeds"},
		{"OC36", "headroom-granted", "green", "one cohesive family"},
		{"OC37", "headroom-swap", "green", "net 0"},
		{"OC37", "headroom-ignored-deletion", "unknown", "net 1, the plan nets 0"},
		{"OC37", "headroom-absent-deletion", "unknown", "internal/sized/z.go is no tracked source"},
		{"OC38", "lines-unbuilt", "unknown", "unobserved estimate"},
		{"OC38", "lines-planned", "planned", "certifies no current source"},
	} {
		t.Run(test.row+" "+test.id, func(t *testing.T) {
			row := one(t, out, "premises", map[string]string{"id": test.id})
			if cell(row, "state") != test.state || !strings.Contains(cell(row, "reason"), test.reason) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}

// TestSymbolPremises drives the symbol, import, anchor-kind, sentinel, error, and bound
// collectors through the actual proposal. Each case names its owner by exact identity, so
// a familiar name, a same-valued literal, or a same-named declaration of another owner
// cannot certify it. Each unproved form keeps its premise unknown with the reason the case
// names.
func TestSymbolPremises(t *testing.T) {
	owner := func(path, symbol, kind, signature string) ownership.Owner {
		return ownership.Owner{Path: path, Package: symbolModule + "/" + filepath.ToSlash(filepath.Dir(path)), Symbol: symbol, Kind: kind, Signature: signature}
	}
	encode := func(value any) string {
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	const rules, other, local, match = "internal/rules/rules.go", "internal/other/other.go", "internal/rules/local.go", "internal/anchors/match.go"
	check := owner(rules, "rules.Check", "func", "func Check(rows int) error")
	write := owner(rules, "rules.write", "func", "func write(v int) int")
	fail := owner(rules, "rules.fail", "func", "func fail(dynamic bool) error")
	missing := owner(rules, "rules.ErrMissing", "var", `var ErrMissing = errors.New("missing")`)
	symbol := func(id string, o ownership.Owner) string {
		return premiseItem(id, "symbol", encode(o), encode(map[string]string{"kind": o.Kind, "signature": o.Signature}))
	}
	edge := func(id, caller string, target ownership.Owner) string {
		return premiseItem(id, "import", encode(map[string]any{"caller": symbolModule + "/" + caller, "target": target.Package, "owner": target}), `{"reachable":true}`)
	}
	anchor := func(id, kind string, line int, needle, source string) string {
		return premiseItem(id, "anchor-kind", encode(map[string]any{"owner": ownership.Location{Path: match, Line: line, Needle: needle, Source: source}, "kind": kind}), `{"supported":true}`)
	}
	rule := func(id, kind string, o ownership.Owner, value string, producers ...ownership.Owner) string {
		if producers == nil {
			producers = []ownership.Owner{}
		}
		return premiseItem(id, kind, encode(map[string]any{"owner": o, "producers": producers}), encode(map[string]string{"value": value}))
	}
	bound := func(id string, o ownership.Owner, value string) string {
		return premiseItem(id, "bound", encode(o), encode(map[string]string{"value": value}))
	}
	misplaced := check
	misplaced.Package, misplaced.Symbol = symbolModule+"/internal/other", "other.Check"
	moved := owner(local, "rules.maxRows", "const", "")
	familiar := missing
	familiar.Package, familiar.Symbol = symbolModule+"/internal/other", "other.ErrMissing"
	facts := withItems("premises",
		symbol("symbol-owner", check),
		symbol("symbol-other-owner", misplaced),
		symbol("symbol-type", owner(rules, "rules.dynamicError", "type", "type dynamicError struct {\n\tmsg string\n}")),
		symbol("symbol-signature", owner(rules, "rules.write", "func", "func write(v string) int")),
		symbol("symbol-kind", owner(rules, "rules.Check", "method", "func Check(rows int) error")),
		edge("import-exported", "internal/other", check),
		edge("import-unexported", "internal/other", write),
		edge("import-internal", "internal/rules", owner("internal/other/internal/secret/secret.go", "secret.Exported", "func", "func Exported()")),
		edge("import-cycle", "internal/rules", owner(other, "other.Check", "func", "func Check() error")),
		anchor("anchor-step", "RequireInStep", 8, "RequireInStep", "tip"),
		anchor("anchor-section-at-step", "RequireInSection", 8, "RequireInStep", "tip"),
		anchor("anchor-invented", "ForbidInStep", 8, "RequireInStep", "tip"),
		anchor("anchor-other-const", "Width", 11, "Width", "tip"),
		rule("sentinel-owner", "sentinel", missing, `errors.New("missing")`, check),
		rule("sentinel-familiar", "sentinel", familiar, `errors.New("missing")`, check),
		rule("sentinel-dynamic", "sentinel", owner(rules, "rules.errDynamic", "var", "var errDynamic = errors.New(reason())"), "errors.New(reason())", check),
		rule("sentinel-unused-producer", "sentinel", missing, `errors.New("missing")`, write),
		rule("sentinel-no-producer", "sentinel", missing, `errors.New("missing")`),
		rule("error-static", "error", owner(rules, "rules.staticError", "type", "type staticError struct{}"), `"static"`, fail),
		rule("error-dynamic", "error", owner(rules, "rules.dynamicError", "type", "type dynamicError struct {\n\tmsg string\n}"), "e.msg", fail),
		bound("bound-owner", owner(rules, "rules.maxRows", "const", ""), "32"),
		bound("bound-local-literal", moved, "32"),
		bound("bound-text", owner(rules, "rules.maxRows", "const", ""), "2 << 4"),
		bound("bound-runtime", owner(rules, "rules.limit", "var", "var limit = len(os.Args)"), "1"))
	out := mustPropose(t, ownershiptest.ProposalRepo{Writes: []string{rules}, Files: symbolBase, Tail: facts}.Seed(t))
	for _, test := range []struct{ row, id, state, reason string }{
		{"OC39", "symbol-owner", "green", "rules.Check resolves in internal/rules/rules.go"},
		{"OC39", "symbol-other-owner", "unknown", "other.Check is declared in internal/other/other.go, not internal/rules/rules.go"},
		{"OC39", "symbol-type", "green", "rules.dynamicError resolves in internal/rules/rules.go"},
		{"OC40", "symbol-signature", "unknown", `declares "func write(v int) int", not "func write(v string) int"`},
		{"OC40", "symbol-kind", "unknown", "rules.Check is a func, not a method"},
		{"OC41", "import-exported", "green", "reaches rules.Check"},
		{"OC41", "import-unexported", "unknown", "rules.write is unexported"},
		{"OC41", "import-internal", "unknown", "outside the tree of " + symbolModule + "/internal/other"},
		{"OC42", "import-cycle", "unknown", "already imports " + symbolModule + "/internal/rules"},
		{"OC43", "anchor-step", "green", "declares the anchor kind RequireInStep at internal/anchors/match.go line 8"},
		{"OC43", "anchor-section-at-step", "unknown", "declares RequireInSection at internal/anchors/match.go line 7, not line 8"},
		{"OC43", "anchor-invented", "unknown", `"ForbidInStep" is no constant of the anchor kind inventory`},
		{"OC43", "anchor-other-const", "unknown", `"Width" is no constant of the anchor kind inventory`},
		{"OC44", "sentinel-owner", "green", `constructs errors.New("missing")`},
		{"OC44", "sentinel-familiar", "unknown", "other.ErrMissing is declared in internal/other/other.go, not internal/rules/rules.go"},
		{"OC44", "sentinel-unused-producer", "unknown", "rules.write does not statically use rules.ErrMissing"},
		{"OC44", "sentinel-no-producer", "unknown", "names no producer"},
		{"OC44", "error-static", "green", `constructs "static"`},
		{"OC45", "sentinel-dynamic", "unknown", "is no static construction"},
		{"OC45", "error-dynamic", "unknown", "the Error method of rules.dynamicError returns no constant"},
		{"OC46", "bound-owner", "green", "rules.maxRows holds the constant value 32"},
		{"OC46", "bound-local-literal", "unknown", "rules.maxRows is declared in internal/rules/rules.go, not internal/rules/local.go"},
		{"OC46", "bound-text", "unknown", "holds the constant value 32, not 2 << 4"},
		{"OC47", "bound-runtime", "unknown", "rules.limit is a var, so its value is no static bound"},
	} {
		t.Run(test.row+" "+test.id, func(t *testing.T) {
			row := one(t, out, "premises", map[string]string{"id": test.id})
			if cell(row, "state") != test.state || !strings.Contains(cell(row, "reason"), test.reason) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}

// TestEvidencePremises drives the reused-test premise collector through the actual
// proposal over committed review-record fixtures. Only a source-bound record whose named
// verification carries a restored, biting behavioral probe of the named test subject with
// red_needle turns green; every unproved form stays unknown with a diagnostic reason.
func TestEvidencePremises(t *testing.T) {
	files := map[string]string{
		evidenceTestPath: "package reused\n\nimport \"testing\"\n\nfunc " + evidenceTestFunc + "(t *testing.T) {}\n",
	}
	root := ownershiptest.ProposalRepo{Writes: []string{"internal/example/after.go"}, Files: files}.Seed(t)
	commits := seedEvidenceCommits(t, root)
	facts := withItems("premises",
		reusedTestItem("evidence-certified", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.valid, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-missing-record", evidenceTestPath, evidenceTestFunc, "specs/unrecorded/spec.md", commits.source, commits.valid, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-stale-source", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.valid, commits.valid, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-confirmation", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.valid, "1", "confirm-1", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-other-subject", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.otherSubject, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-other-mutation", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.valid, "1", "tests", "other mutation", evidenceNeedle),
		reusedTestItem("evidence-missing-needle", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.missingNeedle, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-absent-func", evidenceTestPath, "TestAbsent", evidenceSpec, commits.source, commits.valid, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-unrestored-planned", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.unrestoredPlanned, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-unrestored-probe", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.unrestored, "1", "additional", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-compile-failure", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.compileFailure, "1", "tests", evidenceMutation, evidenceNeedle))
	writeEvidenceSpec(t, root, facts)
	out := mustPropose(t, root)
	for _, test := range []struct{ row, id, state, reason string }{
		{"OC48", "evidence-certified", "green", "is certified by specs/prior/spec.md chunk 1 verification tests"},
		{"OC48", "evidence-missing-record", "unknown", "failed check for source_commit"},
		{"OC48", "evidence-stale-source", "unknown", "binds to"},
		{"OC48", "evidence-confirmation", "unknown", "coordinator confirmation, not a verification"},
		{"OC49", "evidence-other-subject", "unknown", "does not witness test subject TestReused"},
		{"OC49", "evidence-other-mutation", "unknown", `probe mutation is "omit source check", not "other mutation"`},
		{"OC49", "evidence-missing-needle", "unknown", `does not hold red_needle "expected behavioral needle"`},
		{"OC49", "evidence-absent-func", "unknown", "declares no test function TestAbsent"},
		{"OC50", "evidence-unrestored-planned", "unknown", "restore failed"},
		{"OC50", "evidence-unrestored-probe", "unknown", `probe restoration is "fail", not "pass"`},
		{"OC87", "evidence-compile-failure", "unknown", "records a compile failure rather than a behavioral failure of TestReused"},
	} {
		t.Run(test.row+" "+test.id, func(t *testing.T) {
			row := one(t, out, "premises", map[string]string{"id": test.id})
			if cell(row, "state") != test.state || !strings.Contains(cell(row, "reason"), test.reason) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}
