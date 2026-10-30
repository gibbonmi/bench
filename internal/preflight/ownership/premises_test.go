package ownership_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestMetadataGrammar drives each planning-facts grammar case through the actual
// collection. A diagnostic case keeps the inventory incomplete, so a premise row names it;
// a valid inventory reports its own items instead.
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
		{"absent manifest stays diagnostic", "", "metadata", "unknown", "no bench-planning-facts fence"},
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
			row := onePremise(t, premises(t, emptyTree(test.text).Seed(t)), test.id)
			if string(row.State) != test.state || !strings.Contains(row.Reason, test.reason) {
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
			snapshot := ownership.Collect(t.TempDir(), ownership.Source{}, ownership.ParseManifest([]byte(test.text)), nil)
			if snapshot.Complete != test.complete {
				t.Fatalf("Collect = Complete %t; want Complete %t", snapshot.Complete, test.complete)
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
	// A longer run closes the facts fence, and its label is the text after a run of the
	// opener's length, so the facts reader still refuses the nested marker.
	nested := "```" + ownership.MetadataLabel + "\n" + completeEmpty + "\n````\n"
	if lines, _ := scan.Scan([]byte(nested)); lines[len(lines)-2].FenceLabel != "`" {
		t.Fatalf("longer closer = %#v; want the label `", lines[len(lines)-2])
	}
	if got := ownership.ParseManifest([]byte(nested)).Diagnostics; len(got) != 1 || !strings.Contains(got[0].Reason, "another fence marker") {
		t.Fatalf("ParseManifest diagnostics = %v; want the nested-marker refusal", got)
	}
}

// TestSizingPremises drives the line-count, glob-member, and headroom collectors through
// ownership.Collect over the seeded spec. The tip adds one physical line to a.go, so a count names its source.
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
		headroomItem("headroom-granted", "internal/granted", `"internal/granted/n.go"`, "", 1, 1, 1),
		headroomItem("headroom-empty-directory", "", "", "", 3, 3, 0),
		headroomItem("headroom-dot-directory", ".", "", "", 3, 3, 0))
	root := ownershiptest.Repo{Writes: []string{"internal/sized/a.go"}, Files: sizingBase, OwnInventory: true, Tail: facts}.Seed(t)
	preflighttest.MustWriteFile(t, "internal/sized/a.go", "package sized\n\n")
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "tip")
	rows := premises(t, root)
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
		{"OC36", "headroom-empty-directory", "unknown", `subject: directory "" is no clean repository directory`},
		{"OC36", "headroom-dot-directory", "unknown", `subject: directory "." is no clean repository directory`},
		{"OC37", "headroom-swap", "green", "net 0"},
		{"OC37", "headroom-ignored-deletion", "unknown", "net 1, the plan nets 0"},
		{"OC37", "headroom-absent-deletion", "unknown", "internal/sized/z.go is no tracked source"},
		{"OC38", "lines-unbuilt", "unknown", "unobserved estimate"},
		{"OC38", "lines-planned", "planned", "certifies no current source"},
	} {
		t.Run(test.row+" "+test.id, func(t *testing.T) {
			row := onePremise(t, rows, test.id)
			if string(row.State) != test.state || !strings.Contains(row.Reason, test.reason) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}

// TestSymbolPremises drives the symbol, import, anchor-kind, sentinel, error, and bound
// collectors through ownership.Collect over the seeded spec. Each case names its owner by exact identity, so
// a familiar name, a same-valued literal, or a same-named declaration of another owner
// cannot certify it. Each unproved form keeps its premise unknown with the reason the case
// names.
func TestSymbolPremises(t *testing.T) {
	const rules, other, local, match = "internal/rules/rules.go", "internal/other/other.go", "internal/rules/local.go", "internal/anchors/match.go"
	check := symbolOwner(rules, "rules.Check", "func", "func Check(rows int) error")
	write := symbolOwner(rules, "rules.write", "func", "func write(v int) int")
	fail := symbolOwner(rules, "rules.fail", "func", "func fail(dynamic bool) error")
	missing := symbolOwner(rules, "rules.ErrMissing", "var", `var ErrMissing = errors.New("missing")`)
	misplaced := check
	misplaced.Package, misplaced.Symbol = symbolModule+"/internal/other", "other.Check"
	moved := symbolOwner(local, "rules.maxRows", "const", "")
	familiar := missing
	familiar.Package, familiar.Symbol = symbolModule+"/internal/other", "other.ErrMissing"
	facts := withItems("premises",
		symbolItem("symbol-owner", check),
		symbolItem("symbol-other-owner", misplaced),
		symbolItem("symbol-type", symbolOwner(rules, "rules.dynamicError", "type", "type dynamicError struct {\n\tmsg string\n}")),
		symbolItem("symbol-signature", symbolOwner(rules, "rules.write", "func", "func write(v string) int")),
		symbolItem("symbol-kind", symbolOwner(rules, "rules.Check", "method", "func Check(rows int) error")),
		edgeItem("import-exported", "internal/other", check),
		edgeItem("import-unexported", "internal/other", write),
		edgeItem("import-internal", "internal/rules", symbolOwner("internal/other/internal/secret/secret.go", "secret.Exported", "func", "func Exported()")),
		edgeItem("import-cycle", "internal/rules", symbolOwner(other, "other.Check", "func", "func Check() error")),
		anchorItem(match, "anchor-step", "RequireInStep", 8, "RequireInStep", "tip"),
		anchorItem(match, "anchor-section-at-step", "RequireInSection", 8, "RequireInStep", "tip"),
		anchorItem(match, "anchor-invented", "ForbidInStep", 8, "RequireInStep", "tip"),
		anchorItem(match, "anchor-other-const", "Width", 11, "Width", "tip"),
		ruleItem("sentinel-owner", "sentinel", missing, `errors.New("missing")`, check),
		ruleItem("sentinel-familiar", "sentinel", familiar, `errors.New("missing")`, check),
		ruleItem("sentinel-dynamic", "sentinel", symbolOwner(rules, "rules.errDynamic", "var", "var errDynamic = errors.New(reason())"), "errors.New(reason())", check),
		ruleItem("sentinel-unused-producer", "sentinel", missing, `errors.New("missing")`, write),
		ruleItem("sentinel-no-producer", "sentinel", missing, `errors.New("missing")`),
		ruleItem("error-static", "error", symbolOwner(rules, "rules.staticError", "type", "type staticError struct{}"), `"static"`, fail),
		ruleItem("error-dynamic", "error", symbolOwner(rules, "rules.dynamicError", "type", "type dynamicError struct {\n\tmsg string\n}"), "e.msg", fail),
		boundItem("bound-owner", symbolOwner(rules, "rules.maxRows", "const", ""), "32"),
		boundItem("bound-local-literal", moved, "32"),
		boundItem("bound-text", symbolOwner(rules, "rules.maxRows", "const", ""), "2 << 4"),
		boundItem("bound-runtime", symbolOwner(rules, "rules.limit", "var", "var limit = len(os.Args)"), "1"))
	rows := premises(t, ownershiptest.Repo{Writes: []string{rules}, Files: symbolBase, OwnInventory: true, Tail: facts}.Seed(t))
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
			row := onePremise(t, rows, test.id)
			if string(row.State) != test.state || !strings.Contains(row.Reason, test.reason) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}

// TestEvidencePremises drives the reused-test premise collector through ownership.Collect
// over committed review-record fixtures. Only a source-bound record whose named
// verification carries a restored, biting behavioral probe of the named test subject with
// red_needle turns green; every unproved form stays unknown with a diagnostic reason.
func TestEvidencePremises(t *testing.T) {
	root := seedEvidenceRepo(t)
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
		reusedTestItem("evidence-compile-failure", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.compileFailure, "1", "tests", evidenceMutation, evidenceNeedle),
		reusedTestItem("evidence-forged-compile", evidenceTestPath, evidenceTestFunc, evidenceSpec, commits.source, commits.forgedCompile, "1", "tests", evidenceMutation, evidenceNeedle))
	writeEvidenceSpec(t, root, facts)
	rows := premises(t, root)
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
		{"OC87", "evidence-compile-failure", "unknown", "verification tests requires a successful diagnostic probe"},
		{"OC87", "evidence-forged-compile", "unknown", `does not hold red_needle "expected behavioral needle"`},
	} {
		t.Run(test.row+" "+test.id, func(t *testing.T) {
			row := onePremise(t, rows, test.id)
			if string(row.State) != test.state || !strings.Contains(row.Reason, test.reason) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", test.id, row, test.state, test.reason)
			}
		})
	}
}

// TestReusedTestPremiseAcceptsAnEarlierTicketCommit covers RC20: the premise names the chunk
// tip, and its probe entry stays at the commit of the first ticket.
func TestReusedTestPremiseAcceptsAnEarlierTicketCommit(t *testing.T) {
	if row := laterTicketPremise(t, "tests", nil); string(row.State) != "green" {
		t.Fatalf("premise = %v, want green", row)
	}
}

// TestReusedTestPremiseRefusesAnEntryOffTheChunkLine covers RC21 through an entry that only the
// currency test of the premise grades.
func TestReusedTestPremiseRefusesAnEntryOffTheChunkLine(t *testing.T) {
	if row := laterTicketPremise(t, "extra", offLineExtra); string(row.State) != "unknown" || !strings.Contains(row.Reason, "is not current for chunk 1") {
		t.Fatalf("premise = %v, want unknown and a reason that names the entry off the chunk line", row)
	}
}

// TestPredicatePremises drives decision predicate-to-row mappings and seam citation checks
// through ownership.Collect over the seeded spec.
func TestPredicatePremises(t *testing.T) {
	dupPredSpec := strings.Replace(preflighttest.SpecBody(slug), "| PF2 | 1 | does y |", "| PF2 | 1 | does x |", 1)
	for _, test := range []struct{ name, spec, facts, id, state, reason string }{
		{"OC51 exact predicate mapping", "", manifestJSON("declared-complete", nil, plannedSeam, rowDecisions, nil), "dec-1", "green", "maps predicate to coverage row PF1"},
		{"OC52 changed predicate bytes", "", manifestJSON("declared-complete", nil, plannedSeam, []string{decisionItem("dec-stale", "PF1", "does stale", 11, "does x"), rowDecisions[1]}, nil), "dec-stale", "unknown", "does x"},
		{"OC51 missing coverage row", "", manifestJSON("declared-complete", nil, plannedSeam, []string{decisionItem("dec-missing", "PF9", "does x", 11, "does x")}, nil), "metadata", "unknown", "PF9"},
		{"OC51 duplicate row mapping", "", manifestJSON("declared-complete", nil, plannedSeam, []string{decisionItem("dec-1", "PF1", "does x", 11, "does x"), decisionItem("dec-dup", "PF1", "does x", 11, "does x"), rowDecisions[1]}, nil), "dec-1", "unknown", "multiple decisions"},
		{"OC51 multi-row predicate", dupPredSpec, manifestJSON("declared-complete", nil, plannedSeam, []string{decisionItem("dec-m1", "PF1", "does x", 11, "does x"), decisionItem("dec-m2", "PF2", "does x", 12, "does x")}, nil), "dec-m1", "unknown", "multiple coverage rows"},
		{"OC88 omitted second seam citation", citedSpec("`TestFoo`, `TestBar`"), manifestJSON("declared-complete", nil, plannedSeam, nil, nil), "metadata", "unknown", "TestBar"},
		{"OC88 omitted single seam citation", citedSpec("`TestFoo`"), manifestJSON("declared-complete", nil, nil, nil, nil), "metadata", "unknown", "TestFoo"},
		{"OC89 manual claim stays manual", "", manifestJSON("declared-complete", nil, plannedSeam, nil, []string{currentManualJSON("manual-sem", "semantic equivalence")}), "manual-sem", "manual", "semantic equivalence"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := emptyTree(fenced(test.facts)).Seed(t)
			if test.spec != "" {
				writePredicateSpec(t, root, test.spec, test.facts)
			}
			row := onePremise(t, premises(t, root), test.id)
			if string(row.State) != test.state || !strings.Contains(row.Reason, test.reason) {
				t.Fatalf("premise %s = %v, want state %s and reason %q", test.id, row, test.state, test.reason)
			}
		})
	}
}

// TestChargeReadiness grades the charge readiness of proved and unproved premise inventories
// through the actual collection. The audit reads the readiness blocker, so each unmet case
// names its own fault.
func TestChargeReadiness(t *testing.T) {
	for name, facts := range map[string]string{
		"complete decision inventory":           manifestJSON("declared-complete", nil, plannedSeam, rowDecisions, nil),
		"OC51 decisions need not map every row": manifestJSON("declared-complete", nil, plannedSeam, rowDecisions[:1], nil),
	} {
		t.Run(name, func(t *testing.T) {
			root, _ := seedChargeRepo(t, "", facts)
			if ready := collectChargeSnapshot(t, root).ChargeReadiness; !ready.InventoryComplete || ready.Blocker != "" {
				t.Fatalf("charge readiness = %+v, want met", ready)
			}
		})
	}
	for _, test := range []struct {
		name, spec, facts, want string
	}{
		{"OC54 partial inventory", "", manifestJSON("partial", nil, nil, nil, nil), "partial"},
		{"OC54 stale current premise", "", withPremises(stalePremise), "p1"},
		{"OC52 changed predicate bytes", "", manifestJSON("declared-complete", nil, plannedSeam, []string{decisionItem("dec-stale", "PF1", "does stale", 11, "does x"), rowDecisions[1]}, nil), "dec-stale"},
		{"OC51 missing coverage row", "", manifestJSON("declared-complete", nil, plannedSeam, []string{decisionItem("dec-missing", "PF9", "does x", 11, "does x")}, nil), "PF9"},
		{"OC88 omitted second seam citation", citedSpec("`TestFoo`, `TestBar`"), manifestJSON("declared-complete", nil, plannedSeam, nil, nil), "TestBar"},
		{"OC88 omitted single seam citation", citedSpec("`TestFoo`"), manifestJSON("declared-complete", nil, nil, nil, nil), "TestFoo"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, _ := seedChargeRepo(t, test.spec, test.facts)
			assertNotReady(t, root, test.want)
		})
	}
}

// TestChargeManualApplicability grades manual claims through the actual collection. A valid
// native-future claim keeps Snapshot.Complete false, and the charge and current-evidence
// routes still serve it. Each manualRefusals case keeps the charge readiness unmet.
func TestChargeManualApplicability(t *testing.T) {
	valid := validNative().json()
	t.Run("OC89 OC93 native-future claim allows charge and current evidence", func(t *testing.T) {
		root, args := seedChargeRepo(t, "", manifestJSON("declared-complete", nil, nil, nil, []string{valid}))
		snap := collectChargeSnapshot(t, root)
		if snap.Complete || !snap.ChargeReadiness.InventoryComplete || snap.ChargeReadiness.Blocker != "" {
			t.Fatalf("native-future snapshot = Complete %t, readiness %+v, premises %+v, diag %+v", snap.Complete, snap.ChargeReadiness, snap.Premises, snap.Diagnostics)
		}
		if out, code := checkCurrent(assertChargePublished(t, root, args)); code != 0 || !strings.Contains(out, "current[1]") {
			t.Fatalf("current evidence = (%d):\n%s", code, out)
		}
	})
	for _, test := range manualRefusals() {
		t.Run(test.name, func(t *testing.T) {
			root, _ := test.repo.seed(t)
			if test.pending {
				if ready := collectChargeSnapshot(t, root).ChargeReadiness; strings.Contains(ready.Blocker, validNative().id) {
					t.Fatalf("charge readiness = %+v, want the valid native claim pending, not blocking", ready)
				}
			}
			assertNotReady(t, root, test.want)
		})
	}
}
