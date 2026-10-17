package ownership_test

import (
	"cmp"
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/tickets"
)

// Grammar effect fixture paths. command.go passes a local grammar to the parser, and that
// grammar reads its help from help.go. helper.go is an unrelated helper of the same package.
// cross.go passes a grammar whose help reads names.Suffix, and whose handler literal calls
// helper and names.Count. shared.go passes names.Grammar, which grammar.go declares. called.go
// passes a grammar whose help calls names.Help, which help.go of names declares.
// gauge.go, in a package with no Grammar row, passes a literal grammar to the parser.
// dynamic.go, in another bound package, passes a parameter to the parser. operations.go
// declares the custom root operations beside a case-variant Operations.
const (
	commandFile    = "internal/harnesses/command.go"
	helpFile       = "internal/harnesses/help.go"
	helperFile     = "internal/harnesses/helper.go"
	crossFile      = "internal/harnesses/cross.go"
	namesFile      = "internal/names/names.go"
	sharedFile     = "internal/harnesses/shared.go"
	namesGrammar   = "internal/names/grammar.go"
	calledFile     = "internal/harnesses/called.go"
	namesHelp      = "internal/names/help.go"
	gaugeFile      = "internal/gauge/gauge.go"
	dynamicFile    = "internal/roadmap/dynamic.go"
	terminalFile   = "internal/terminal/terminal.go"
	operationsFile = "internal/preflight/evidencecmd/operations.go"
	breakdownFile  = "internal/preflight/evidencecmd/breakdown.go"
	usageFile      = "internal/usage/usage.go"
	rootHelp       = "cmd/bench/help_inventory_test.go"
	routingCensus  = "internal/conformance/subcommand_routing_table_test.go"
)

// grammarBase is the base source of every grammar effect case. It has no go.mod, so every
// import is outside the module.
var grammarBase = map[string]string{
	commandFile:    "package harnesses\n\nimport \"github.com/gibbonmi/bench/internal/usage\"\n\nfunc Command(args []string) int {\n\tgrammar := usage.Grammar{Cmd: \"bench harnesses\", Help: harnessesHelp}\n\t_, _, code := usage.Parse(grammar, args)\n\treturn code + helper()\n}\n",
	helpFile:       "package harnesses\n\nconst harnessesHelp = \"usage: bench harnesses\"\n",
	helperFile:     "package harnesses\n\nfunc helper() int { return 0 }\n",
	crossFile:      "package harnesses\n\nimport (\n\t\"github.com/gibbonmi/bench/internal/names\"\n\t\"github.com/gibbonmi/bench/internal/usage\"\n)\n\nfunc Cross(args []string) int {\n\tgrammar := usage.Grammar{Cmd: \"bench cross\", Help: harnessesHelp + names.Suffix, Run: func() int { return helper() + names.Count() }}\n\t_, _, code := usage.Parse(grammar, args)\n\treturn code\n}\n",
	namesFile:      "package names\n\nconst Suffix = \" [--all]\"\n\nfunc Count() int { return 0 }\n",
	gaugeFile:      "package gauge\n\nimport \"github.com/gibbonmi/bench/internal/usage\"\n\nfunc Command(args []string) int {\n\t_, _, code := usage.Parse(usage.Grammar{Cmd: \"bench gauge\"}, args)\n\treturn code\n}\n",
	dynamicFile:    "package roadmap\n\nimport \"github.com/gibbonmi/bench/internal/usage\"\n\nfunc run(g usage.Grammar, args []string) int {\n\t_, _, code := usage.Parse(g, args)\n\treturn code\n}\n",
	terminalFile:   "package terminal\n\nfunc Run() int { return 0 }\n",
	operationsFile: "package evidencecmd\n\nvar operations = []string{\"build\"}\n\nvar Operations = []string{\"x\"}\n",
	breakdownFile:  "package evidencecmd\n",
}

// crossBase is grammarBase inside the module that its imports name, with the parser package
// the module's package index needs and a grammar value that another package passes.
var crossBase = merged(grammarBase, map[string]string{
	"go.mod":     "module github.com/gibbonmi/bench\n\ngo 1.25\n",
	usageFile:    "package usage\n\ntype Grammar struct {\n\tCmd, Help string\n\tRun       func() int\n}\n\nfunc Parse(g Grammar, args []string) (any, any, int) { return nil, nil, 0 }\n",
	sharedFile:   "package harnesses\n\nimport (\n\t\"github.com/gibbonmi/bench/internal/names\"\n\t\"github.com/gibbonmi/bench/internal/usage\"\n)\n\nfunc Shared(args []string) int {\n\t_, _, code := usage.Parse(names.Grammar, args)\n\treturn code\n}\n",
	namesGrammar: "package names\n\nimport \"github.com/gibbonmi/bench/internal/usage\"\n\nvar Grammar = usage.Grammar{Cmd: \"bench shared\"}\n",
})

// edited is one base file with its first old replaced by new.
func edited(path, old, new string) map[string]string {
	return map[string]string{path: strings.Replace(grammarBase[path], old, new, 1)}
}

// plannedForm is a partial inventory that declares one planned grammar change of symbol in
// operations.go or command.go.
func plannedForm(path, symbol, kind string) string {
	return fenced(strings.Replace(strings.Replace(completeEmpty, `"declared-complete"`, `"partial"`, 1), `"changes":[]`,
		`"changes":[{"id":"new-form",`+itemLocation+`,"kind":"grammar","owner":{"path":"`+path+`","package":"github.com/gibbonmi/bench/`+filepath.ToSlash(filepath.Dir(path))+`","symbol":"`+symbol+`","kind":"`+kind+`","signature":""},"caller_scope":"production-and-tests","state":"planned","destination":null,"dispositions":[]}]`, 1))
}

// effectCase is one actual proposal over a committed base and a tip that edits the case's
// files; an empty tip body deletes its file. want maps each path to whether the closure
// requires it, and evidence maps a required path to text that one of its reasons holds. A
// later hop of the fixed point can require a path again, so only its reason shows which
// rule required it. premise names a premise row whose reason must hold reason. A nonempty
// refusal is text that the proposal's refusal must name instead. A non-nil seal publishes a
// verified kit binary, so the registry source is the authenticated kit context, and located
// requires that the verified source locates every registry reason: no kit-registry premise.
type effectCase struct {
	name                          string
	writes                        []string
	tip                           map[string]string
	facts                         string
	want                          map[string]bool
	evidence                      map[string]string
	premise, premiseState, reason string
	refusal                       string
	files                         map[string]string
	fence                         []string
	seal                          ownershiptest.Publisher
	located                       bool
}

// run seeds the case's base, commits its tip, and grades the proposal. A case with no facts
// declares the complete-empty inventory.
func (test effectCase) run(t *testing.T) {
	root := ownershiptest.ProposalRepo{Writes: test.writes, Files: test.files, Fence: test.fence, OwnInventory: true, Seal: test.seal, Tail: cmp.Or(test.facts, fenced(completeEmpty))}.Seed(t)
	if len(test.tip) > 0 {
		for path, body := range test.tip {
			if body == "" {
				preflighttest.RunGit(t, "rm", "-q", path)
				continue
			}
			preflighttest.MustWriteFile(t, path, body)
		}
		preflighttest.RunGit(t, "add", "-A")
		preflighttest.RunGit(t, "commit", "-q", "-m", "tip")
	}
	if test.refusal != "" {
		out, code := propose(t, root)
		refused(t, out, code, test.refusal)
		return
	}
	out := mustPropose(t, root)
	for path, want := range test.want {
		if got := len(matching(t, out, "closure", map[string]string{"ticket": "one.md", "requires": path})) > 0; got != want {
			t.Errorf("closure requires %s = %t, want %t:\n%s", path, got, want, out)
		}
	}
	for path, reason := range test.evidence {
		rows := matching(t, out, "closure_evidence", map[string]string{"ticket": "one.md", "path": path})
		if !slices.ContainsFunc(rows, func(row map[string]any) bool { return strings.Contains(cell(row, "reason"), reason) }) {
			t.Errorf("evidence for %s = %v, want a reason naming %q", path, rows, reason)
		}
	}
	if rows := matching(t, out, "premises", map[string]string{"id": "kit-registry"}); test.located && len(rows) > 0 {
		t.Errorf("kit-registry premise = %v, want every registry reason located", rows)
	}
	if test.premise != "" {
		state := cmp.Or(test.premiseState, "unknown")
		row := one(t, out, "premises", map[string]string{"id": test.premise})
		if cell(row, "state") != state || !strings.Contains(cell(row, "reason"), test.reason) {
			t.Errorf("premise %s = %v, want %s and a reason naming %q", test.premise, row, state, test.reason)
		}
	}
}

// TestGrammarEffects drives each grammar effect through the actual proposal.
func TestGrammarEffects(t *testing.T) {
	cases := []effectCase{
		{name: "OC14 body-only edit adds no grammar obligation", writes: []string{commandFile, helperFile},
			tip:  merged(edited(helperFile, "return 0", "return 1"), edited(commandFile, "code + helper()", "code + helper() + 1")),
			want: map[string]bool{rootHelp: false, routingCensus: false, helpFile: false}},
		{name: "OC15 body-only terminal edit keeps its Core obligation", writes: []string{terminalFile},
			tip:  edited(terminalFile, "return 0", "return 1"),
			want: map[string]bool{"internal/adopt/setup_test.go": true, "internal/gate/command_test.go": true, rootHelp: false}},
		{name: "OC16 changed producer requires its grammar source", writes: []string{helpFile},
			tip:  edited(helpFile, "usage: bench harnesses", "usage: bench harnesses [--all]"),
			want: map[string]bool{commandFile: true, rootHelp: true, routingCensus: true}},
		{name: "OC17 planned unbound form requires the root help inventory", writes: []string{breakdownFile},
			facts:    plannedForm(operationsFile, "evidencecmd.operations", "var"),
			want:     map[string]bool{operationsFile: true, rootHelp: true, routingCensus: true},
			evidence: map[string]string{rootHelp: "planned grammar change new-form requires " + rootHelp}},
		{name: "changed custom root requires the shared inventories", writes: []string{operationsFile},
			tip:  edited(operationsFile, `"build"`, `"build", "breakdown"`),
			want: map[string]bool{rootHelp: true, routingCensus: true}},
		{name: "OC18 Command-named registration stays diagnostic", writes: []string{helperFile},
			facts: plannedForm(commandFile, "harnesses.Command", "func"),
			want:  map[string]bool{commandFile: false}, premise: "new-form", reason: "no registered grammar producer"},
		{name: "OC18 dynamic parser argument stays diagnostic", writes: []string{dynamicFile},
			tip:  edited(dynamicFile, "return code", "return code + 1"),
			want: map[string]bool{rootHelp: true}, premise: "grammar-effects", reason: "no static producer"},
		{name: "OC19 case-variant selector stays distinct", writes: []string{breakdownFile},
			facts: plannedForm(operationsFile, "evidencecmd.Operations", "var"),
			want:  map[string]bool{operationsFile: false, rootHelp: false}, premise: "new-form", reason: "no registered grammar producer"},
		{name: "OC19 case-variant declaration is no root", writes: []string{operationsFile},
			tip:  edited(operationsFile, `"x"`, `"y"`),
			want: map[string]bool{rootHelp: false}},
		{name: "cross-package producer requires its root and the shared inventories", writes: []string{namesFile}, files: crossBase,
			tip:      edited(namesFile, "[--all]", "[--all] [--json]"),
			want:     map[string]bool{crossFile: true, rootHelp: true, routingCensus: true},
			evidence: map[string]string{crossFile: "internal/names.Suffix changes the grammar that internal/harnesses.Cross.grammar produces"}},
		{name: "cross-package called producer requires its root and the shared inventories", writes: []string{namesHelp},
			files: merged(crossBase, map[string]string{calledFile: "package harnesses\n\nimport (\n\t\"github.com/gibbonmi/bench/internal/names\"\n\t\"github.com/gibbonmi/bench/internal/usage\"\n)\n\nfunc Called(args []string) int {\n\tgrammar := usage.Grammar{Cmd: \"bench called\", Help: names.Help()}\n\t_, _, code := usage.Parse(grammar, args)\n\treturn code\n}\n",
				namesHelp: "package names\n\nfunc Help() string { return \"usage: bench called\" }\n"}),
			tip:      map[string]string{namesHelp: "package names\n\nfunc Help() string { return \"usage: bench called [--all]\" }\n"},
			want:     map[string]bool{calledFile: true, rootHelp: true, routingCensus: true},
			evidence: map[string]string{calledFile: "internal/names.Help changes the grammar that internal/harnesses.Called.grammar produces"}},
		{name: "a referrer's unresolved producer stays its own diagnostic", writes: []string{namesFile},
			files: merged(crossBase, map[string]string{dynamicFile: "package roadmap\n\nimport (\n\t\"github.com/gibbonmi/bench/internal/names\"\n\t\"github.com/gibbonmi/bench/internal/usage\"\n)\n\nfunc run(g usage.Grammar, args []string) int {\n\t_, _, code := usage.Parse(g, args)\n\treturn code + names.Count()\n}\n"}),
			tip:   edited(namesFile, "[--all]", "[--all] [--json]"),
			want:  map[string]bool{crossFile: true, rootHelp: true}},
		{name: "another package's grammar value requires the passing root", writes: []string{namesGrammar}, files: crossBase,
			tip:      map[string]string{namesGrammar: strings.Replace(crossBase[namesGrammar], "bench shared", "bench shared [--all]", 1)},
			want:     map[string]bool{sharedFile: true, rootHelp: true, routingCensus: true},
			evidence: map[string]string{sharedFile: "internal/names.Grammar changes the grammar that internal/harnesses.Shared#1 produces"}},
		{name: "handler literal body adds no grammar obligation", writes: []string{helperFile, namesFile}, files: crossBase,
			tip:  merged(edited(helperFile, "return 0", "return 1"), edited(namesFile, "return 0", "return 1")),
			want: map[string]bool{crossFile: false, rootHelp: false, routingCensus: false}},
		{name: "unbound package grammar requires the shared inventories", writes: []string{gaugeFile},
			tip:      edited(gaugeFile, "bench gauge", "bench gauge [--all]"),
			want:     map[string]bool{rootHelp: true, routingCensus: true},
			evidence: map[string]string{rootHelp: "grammar root internal/gauge.Command#1 requires " + rootHelp}},
		{name: "a verified kit source locates an unbound package's parser root", writes: []string{gaugeFile}, seal: freshness.Publish, files: kitGauge, tip: kitGaugeTip,
			want:     map[string]bool{rootHelp: true},
			evidence: map[string]string{rootHelp: "grammar root internal/gauge.Command#1 requires " + rootHelp + " beside " + gaugeFile},
			located:  true},
		{name: "a kit source with no parser row stays diagnostic", writes: []string{gaugeFile}, seal: freshness.Publish, tip: kitGaugeTip,
			files: merged(kitGauge, map[string]string{tickets.BindingRegistry: withoutParserRow(t)}),
			want:  map[string]bool{rootHelp: true}, premise: "kit-registry", reason: "holds no row for grammar root internal/gauge.Command#1"},
		{name: "ambiguous producer declaration stays diagnostic", writes: []string{helpFile},
			files: merged(grammarBase, map[string]string{"internal/harnesses/help_other.go": "//go:build ignore\n\npackage harnesses\n\nconst harnessesHelp = \"other\"\n"}),
			tip:   edited(helpFile, "usage: bench harnesses", "usage: bench harnesses [--all]"),
			want:  map[string]bool{rootHelp: true}, premise: "grammar-effects", reason: "reaches the ambiguous declaration harnessesHelp"},
		{name: "malformed tip source stays diagnostic", writes: []string{helperFile},
			tip:  map[string]string{helperFile: malformedHelper},
			want: map[string]bool{rootHelp: true}, premise: "grammar-effects", reason: "tip " + helperFile + ":"},
		{name: "malformed base source stays diagnostic", writes: []string{helperFile},
			files: merged(grammarBase, map[string]string{helperFile: malformedHelper}),
			tip:   map[string]string{helperFile: grammarBase[helperFile]},
			want:  map[string]bool{rootHelp: true}, premise: "grammar-effects", reason: "base " + helperFile + ":"},
	}
	for _, test := range cases {
		if test.files == nil {
			test.files = grammarBase
		}
		t.Run(test.name, test.run)
	}
}

// malformedHelper is helper.go with a syntax error.
const malformedHelper = "package harnesses\n\nfunc helper( int {\n"

// merged is the union of several edit sets.
func merged(sets ...map[string]string) map[string]string {
	all := map[string]string{}
	for _, set := range sets {
		maps.Copy(all, set)
	}
	return all
}

// TestHelperCallers drives each declared helper change through the actual proposal. The
// ticket owns the helper and its first caller, so each other caller is a proposed write.
// The third caller sits outside every seed path, so a sampled sweep misses it. After a
// relocation the third caller stops using the helper, so only the before-state census can
// still name it.
func TestHelperCallers(t *testing.T) {
	signature := withItems("changes", changeItem("signature", "existing", mergeHelper, "merge.mergeFixture", ""))
	cases := []effectCase{
		{name: "OC20 signature change proposes each static test caller", writes: []string{mergeHelper, callerA}, facts: signature,
			want:     map[string]bool{callerB: true, callerC: true},
			evidence: map[string]string{callerC: "declared in " + mergeHelper + " line 3 at the source tip"}},
		{name: "OC21 fixture-helper change proposes its external call site", writes: []string{seedFile},
			facts:    withItems("changes", changeItem("fixture-helper", "existing", seedFile, "mergetest.Seed", "")),
			want:     map[string]bool{external: true},
			evidence: map[string]string{external: "uses mergetest.Seed"}},
		{name: "OC22 relocation retains the before-state caller", writes: []string{mergeHelper + " (deleted)", movedFile + " (new)", callerA},
			facts: withItems("changes", changeItem("signature", "existing", mergeHelper, "merge.mergeFixture", movedFile)),
			tip:   map[string]string{mergeHelper: "", movedFile: helperBase[mergeHelper], callerC: caller("C", false)},
			want:  map[string]bool{callerB: true, callerC: true},
			evidence: map[string]string{callerB: "declared in " + movedFile + " line 3 at the source tip",
				callerC: "declared in " + mergeHelper + " line 3 at the source base"}},
		{name: "OC23 deletion retains the before-state callers", writes: []string{mergeHelper + " (deleted)", callerA}, facts: signature,
			tip:      map[string]string{mergeHelper: "", callerA: caller("A", false), callerB: caller("B", false), callerC: caller("C", false)},
			want:     map[string]bool{callerB: true, callerC: true},
			evidence: map[string]string{callerB: "at the source base", callerC: "at the source base"}},
		{name: "OC23 deleted fixture helper retains its external before-state caller", writes: []string{seedFile + " (deleted)"},
			facts:    withItems("changes", changeItem("fixture-helper", "existing", seedFile, "mergetest.Seed", "")),
			tip:      map[string]string{seedFile: "", external: "package other\n\nimport \"testing\"\n\nfunc TestOther(t *testing.T) {}\n"},
			want:     map[string]bool{external: true},
			evidence: map[string]string{external: "declared in " + seedFile + " line 3 at the source base"}},
		{name: "planned helper has no invented callers", writes: []string{mergeHelper, callerA},
			facts: withItems("changes", changeItem("signature", "planned", mergeHelper, "merge.mergeFixture", "")),
			want:  map[string]bool{callerB: false, callerC: false}},
		{name: "a present method's former caller stays diagnostic", writes: []string{teeFile, teeCallerD},
			facts:   withItems("changes", changeItem("signature", "existing", teeFile, "merge.T.Merge", "")),
			tip:     map[string]string{teeFile: strings.Replace(helperBase[teeFile], "return a", `return a + "!"`, 1), teeCallerE: teeCaller("E", false)},
			premise: "caller-census", reason: "merge.T.Merge is a method, whose former callers need the types"},
		{name: "a liveness change grades its dispositions against the retirement census", writes: []string{policyFile},
			facts: withItems("changes", `{"id":"retire",`+itemLocation+`,"kind":"liveness","owner":`+ownerJSON(policyFile, "policy.retired", "func")+
				`,"caller_scope":"production-and-tests","state":"existing","destination":null,"dispositions":[{"symbol":"policy.keepA","action":"remove"}]}`),
			premise: "liveness-census", reason: "internal/policy/consumers.go keepB has 0 dispositions, want 1"},
		{name: "OC53 dynamic caller claim stays diagnostic", writes: []string{mergeHelper, callerA}, facts: signature,
			tip:     map[string]string{callerA: "package merge\n\nimport (\n\t\"reflect\"\n\t\"testing\"\n)\n\nfunc TestA(t *testing.T) { _ = reflect.ValueOf(mergeFixture) }\n"},
			premise: "caller-census", reason: "dynamic"},
	}
	for _, test := range cases {
		test.files, test.fence = helperBase, helperFence
		t.Run(test.name, test.run)
	}
}

// TestDeletedEntry drives each deletion marker shape through the actual proposal: a file
// tracked at the source base resolves whether the tip still holds it or not, and no other
// shape gains authority from the marker. The legacy unmarked committed deletion stays
// accepted, while the same unmarked spelling of a never-tracked path stays red.
func TestDeletedEntry(t *testing.T) {
	never := "internal/merge/never.go"
	cases := []effectCase{
		{name: "OC24 deleted marker resolves a removed before path", writes: []string{goneFile + " (deleted)"}, tip: map[string]string{goneFile: ""}},
		{name: "OC26 planned deleted marker accepts a present before path", writes: []string{goneFile + " (deleted)"}},
		{name: "OC25 never-tracked deleted path refuses", writes: []string{never + " (deleted)"}, refusal: "no tracked file at the review base: one.md: " + never},
		{name: "OC25 destination cannot repair a never-tracked deletion", writes: []string{never + " (deleted)", movedFile + " (new)"},
			tip: map[string]string{movedFile: "package merge\n"}, refusal: "no tracked file at the review base: one.md: " + never},
		{name: "OC25 directory before path refuses", writes: []string{"internal/merge (deleted)"}, refusal: "no tracked file at the review base"},
		{name: "OC25 deleted directory shorthand refuses", writes: []string{"internal/merge/ (deleted)"}, refusal: "names one exact file, not a directory"},
		{name: "OC27 new then deleted markers refuse parsing", writes: []string{goneFile + " (new) (deleted)"}, refusal: "both the (new) and the (deleted) marker"},
		{name: "OC27 deleted then new markers refuse parsing", writes: []string{goneFile + " (deleted) (new)"}, refusal: "both the (new) and the (deleted) marker"},
		{name: "OC82 committed D-status stays accepted", writes: []string{goneFile}, tip: map[string]string{goneFile: ""}},
		{name: "OC82 unmarked never-tracked path stays red", writes: []string{never}, refusal: "names no tree path and carries no (new) marker: one.md: " + never},
	}
	for _, test := range cases {
		test.files, test.fence = helperBase, helperFence
		t.Run(test.name, test.run)
	}
}

// TestOwnerPremises grades each owner premise of one inventory through the actual
// proposal. Each partial census omits a second real member: a second lifecycle consumer, an
// independent writer site, and a second renderer branch. The printer's read of the field is
// no producer.
func TestOwnerPremises(t *testing.T) {
	consumer := func(symbol string) string { return ownerJSON("internal/policy/consumers.go", symbol, "func") }
	home := func(line int) string {
		return `{"owner":{"path":"` + policyFile + `","line":` + fmt.Sprint(line) + `,"needle":"func home() bool","source":"tip"}}`
	}
	list := func(path string, symbols ...string) string {
		owners := make([]string, len(symbols))
		for i, symbol := range symbols {
			owners[i] = ownerJSON(path, symbol, "func")
		}
		return "[" + strings.Join(owners, ",") + "]"
	}
	probe := `{"path":"internal/policy/policy_test.go","function":"TestWrite","record_spec":"specs/policy/spec.md","source_commit":"a","evidence_commit":"b","chunk":"C1","verification":"policy-1","mutation":"drop-sync"`
	posture := func(callers, extra string) string {
		return `{"owner":` + ownerJSON(policyFile, "policy.write", "func") + `,"callers":` + callers + `,"writer_probe":` + probe + extra + `}}`
	}
	field := func(producers string) string {
		return `{"field":` + ownerJSON(policyFile, "policy.row.Reason", "field") + `,"producers":` + producers + `}`
	}
	writers := list("internal/policy/writers.go", "policy.writerA")
	both := `[` + ownerJSON("internal/policy/writers.go", "policy.writerA", "func") + `,` + ownerJSON("internal/policy/sync.go", "policy.syncWrite", "func") + `]`
	renderers := list("internal/policy/render.go", "policy.renderA", "policy.renderB")
	facts := withItems("premises",
		premiseItem("home-owner", "reused-owner", home(3), `{"required":true}`),
		premiseItem("moved-owner", "reused-owner", home(4), `{"required":true}`),
		premiseItem("partial-retirement", "retirement", `{"owner":`+ownerJSON(policyFile, "policy.retired", "func")+`,"consumers":[`+consumer("policy.keepA")+`]}`, `{"dispositions":[{"symbol":"policy.keepA","action":"remove"}]}`),
		premiseItem("complete-retirement", "retirement", `{"owner":`+ownerJSON(policyFile, "policy.retired", "func")+`,"consumers":[`+consumer("policy.keepA")+`,`+consumer("policy.keepB")+`]}`,
			`{"dispositions":[{"symbol":"policy.keepA","action":"remove"},{"symbol":"policy.keepB","action":"retain"}]}`),
		premiseItem("partial-posture", "posture", posture(writers, ""), `{"callers":`+writers+`}`),
		premiseItem("complete-posture", "posture", posture(both, ""), `{"callers":`+both+`}`),
		premiseItem("loose-probe", "posture", posture(both, `,"extra":true`), `{"callers":`+both+`}`),
		premiseItem("partial-field", "printed-field", field(list("internal/policy/render.go", "policy.renderA")), `{"producers":`+list("internal/policy/render.go", "policy.renderA")+`}`),
		premiseItem("complete-field", "printed-field", field(renderers), `{"producers":`+renderers+`}`),
		premiseItem("mismatched-field", "printed-field", field(renderers), `{"producers":`+list("internal/policy/render.go", "policy.renderA")+`}`))
	out := mustPropose(t, ownershiptest.ProposalRepo{Writes: []string{"internal/policy/render.go"}, Files: helperBase, Fence: helperFence, OwnInventory: true, Tail: facts}.Seed(t))
	t.Run("OC28 reused owner adds its rule source", func(t *testing.T) {
		one(t, out, "closure", map[string]string{"ticket": "one.md", "path": "internal/policy/render.go", "requires": policyFile})
		row := one(t, out, "closure_evidence", map[string]string{"ticket": "one.md", "path": policyFile, "kind": "owner"})
		if cell(row, "source") != policyFile || cell(row, "line") != "3" || len(matching(t, out, "writes_proposal", map[string]string{"path": policyFile})) != 0 {
			t.Fatalf("owner evidence = %v, want %s line 3 and no proposed write", row, policyFile)
		}
	})
	for id, want := range map[string][2]string{
		"home-owner":          {"green", "holds its needle"},
		"moved-owner":         {"unknown", "line 4 does not hold its needle exactly once"},
		"partial-retirement":  {"unknown", "omits internal/policy/consumers.go keepB"},
		"complete-retirement": {"green", "complete at the source tip: 2 member(s)"},
		"partial-posture":     {"unknown", "omits internal/policy/sync.go syncWrite"},
		"complete-posture":    {"unknown", "no collector certifies its writer_probe evidence"},
		"loose-probe":         {"unknown", `subject: writer_probe: unknown field "extra"`},
		"partial-field":       {"unknown", "omits internal/policy/render.go renderB"},
		"complete-field":      {"green", "complete at the source tip: 2 member(s)"},
		"mismatched-field":    {"unknown", "the expected producers differ from the subject producers"},
	} {
		t.Run(id, func(t *testing.T) {
			row := one(t, out, "premises", map[string]string{"id": id})
			if cell(row, "state") != want[0] || !strings.Contains(cell(row, "reason"), want[1]) {
				t.Fatalf("premise %s = %v, want %s and a reason naming %q", id, row, want[0], want[1])
			}
		})
	}
	printedFieldWrites(t, field(renderers), renderers)
}
