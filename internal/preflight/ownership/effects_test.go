package ownership_test

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// Grammar effect fixture paths. command.go passes a local grammar to the parser, and that
// grammar reads its help from help.go. helper.go is an unrelated helper of the same package.
// dynamic.go, in another bound package, passes a parameter to the parser.
// operations.go declares the custom root operations beside a case-variant Operations.
const (
	commandFile    = "internal/harnesses/command.go"
	helpFile       = "internal/harnesses/help.go"
	helperFile     = "internal/harnesses/helper.go"
	dynamicFile    = "internal/roadmap/dynamic.go"
	terminalFile   = "internal/terminal/terminal.go"
	operationsFile = "internal/preflight/evidencecmd/operations.go"
	breakdownFile  = "internal/preflight/evidencecmd/breakdown.go"
	rootHelp       = "cmd/bench/help_inventory_test.go"
	routingCensus  = "internal/conformance/subcommand_routing_table_test.go"
)

// grammarBase is the base source of every grammar effect case.
var grammarBase = map[string]string{
	commandFile:    "package harnesses\n\nimport \"github.com/gibbonmi/bench/internal/usage\"\n\nfunc Command(args []string) int {\n\tgrammar := usage.Grammar{Cmd: \"bench harnesses\", Help: harnessesHelp}\n\t_, _, code := usage.Parse(grammar, args)\n\treturn code + helper()\n}\n",
	helpFile:       "package harnesses\n\nconst harnessesHelp = \"usage: bench harnesses\"\n",
	helperFile:     "package harnesses\n\nfunc helper() int { return 0 }\n",
	dynamicFile:    "package roadmap\n\nimport \"github.com/gibbonmi/bench/internal/usage\"\n\nfunc run(g usage.Grammar, args []string) int {\n\t_, _, code := usage.Parse(g, args)\n\treturn code\n}\n",
	terminalFile:   "package terminal\n\nfunc Run() int { return 0 }\n",
	operationsFile: "package evidencecmd\n\nvar operations = []string{\"build\"}\n\nvar Operations = []string{\"x\"}\n",
	breakdownFile:  "package evidencecmd\n",
}

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

// TestGrammarEffects drives each grammar effect through the actual proposal over a committed
// base and a tip that edits the case's files. want maps each path to whether the closure
// requires it, and evidence maps a required path to text that one of its reasons holds. A
// later hop of the fixed point can require a path again, so only its reason shows which
// rule required it. premise names a premise row whose reason must hold reason.
func TestGrammarEffects(t *testing.T) {
	cases := []struct {
		name            string
		writes          []string
		tip             map[string]string
		facts           string
		want            map[string]bool
		evidence        map[string]string
		premise, reason string
	}{
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
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root := ownershiptest.ProposalRepo{Writes: test.writes, Files: grammarBase, Tail: test.facts}.Seed(t)
			if len(test.tip) > 0 {
				for path, body := range test.tip {
					preflighttest.MustWriteFile(t, path, body)
				}
				preflighttest.RunGit(t, "add", "-A")
				preflighttest.RunGit(t, "commit", "-q", "-m", "tip")
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
			if test.premise != "" {
				row := one(t, out, "premises", map[string]string{"id": test.premise})
				if cell(row, "state") != "unknown" || !strings.Contains(cell(row, "reason"), test.reason) {
					t.Errorf("premise %s = %v, want unknown and a reason naming %q", test.premise, row, test.reason)
				}
			}
		})
	}
}

// merged is the union of several edit sets.
func merged(sets ...map[string]string) map[string]string {
	all := map[string]string{}
	for _, set := range sets {
		maps.Copy(all, set)
	}
	return all
}
