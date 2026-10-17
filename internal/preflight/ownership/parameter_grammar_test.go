package ownership_test

import (
	"cmp"
	"strings"
	"testing"
)

// TestParameterGrammarResolvesThroughCallers drives a parameter grammar of the unbound gauge
// package through the actual proposal. run.go passes its parameter to the parser, and
// callers.go calls run with a package grammar and with a local grammar, so the parameter
// roots to both. One more caller that passes its own parameter, or a use of run as a value,
// hides a call site, so the parser call stays unresolved.
func TestParameterGrammarResolvesThroughCallers(t *testing.T) {
	run := strings.Replace(grammarBase[dynamicFile], "package roadmap", "package gauge", 1)
	callers := "package gauge\n\nimport \"github.com/gibbonmi/bench/internal/usage\"\n\nvar listGrammar = usage.Grammar{Cmd: \"bench gauge list\"}\n\nfunc List(args []string) int { return run(listGrammar, args) }\n\nfunc Show(args []string) int {\n\tgrammar := usage.Grammar{Cmd: \"bench gauge show\"}\n\treturn run(grammar, args)\n}\n"
	callersFile, tip := "internal/gauge/callers.go", strings.NewReplacer("gauge list", "gauge list [--all]")
	cases := []effectCase{
		{name: "AP49 declared callers root the parameter",
			want:     map[string]bool{rootHelp: true, routingCensus: true},
			evidence: map[string]string{rootHelp: "grammar root internal/gauge.listGrammar requires " + rootHelp}},
		{name: "AP49 one undeclared caller keeps the parameter unresolved", premise: "grammar-effects", reason: "no static producer",
			files: map[string]string{callersFile: callers + "\nfunc relay(g usage.Grammar, args []string) int { return run(g, args) }\n"}},
		{name: "AP49 a value use of the function keeps the parameter unresolved", premise: "grammar-effects", reason: "no static producer",
			files: map[string]string{callersFile: callers + "\nvar handler = run\n"}},
	}
	for _, test := range cases {
		body := cmp.Or(test.files[callersFile], callers)
		test.writes, test.files, test.tip = []string{callersFile}, merged(grammarBase, map[string]string{"internal/gauge/run.go": run, callersFile: body}), map[string]string{callersFile: tip.Replace(body)}
		t.Run(test.name, test.run)
	}
}
