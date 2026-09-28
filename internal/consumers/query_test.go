package consumers

import (
	"strings"
	"testing"
)

// scopeFixture plants one consumer of the target in a production file and one in a
// _test.go file of the same package, so each filter has exactly one row to keep and one
// row to drop.
func scopeFixture(root string) []fixturePkg {
	use := func(name string) string {
		return "package consumer\n\nimport \"example.com/target\"\n\nfunc " + name + "() { target.Symbol() }\n"
	}
	return []fixturePkg{targetPkg(root), {path: "example.com/consumer", files: map[string]string{
		root + "/consumer/consumer.go":      use("Run"),
		root + "/consumer/consumer_test.go": use("check"),
	}}}
}

// pairFixture declares two target symbols and consumes each once, so a two-symbol query
// has one row per symbol to label.
func pairFixture(root string) []fixturePkg {
	return []fixturePkg{
		declPkg(root, "target", "target", "package target\n\nfunc Symbol() {}\n\nfunc Other() {}\n"),
		declPkg(root, "consumer", "consumer", "package consumer\n\nimport \"example.com/target\"\n\n"+
			"func UseSymbol() { target.Symbol() }\n\nfunc UseOther() { target.Other() }\n"),
	}
}

// With no filter the scope fixture answers both rows, so a filter test that drops one of
// them observes the filter and not a fixture that never planted the row.
func TestUnfilteredQueryKeepsProductionAndTestRows(t *testing.T) {
	stubLoad(t, scopeFixture)
	out, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	want := "consumers[2]{file,line,via,enclosing}:\n" +
		"  consumer/consumer.go,5,call,Run\n" +
		"  consumer/consumer_test.go,5,call,check\n"
	if !strings.Contains(out, want) {
		t.Fatalf("stdout = %q, want both rows %q", out, want)
	}
}

// --test keeps only the rows whose file is a _test.go file, and the meta accounting
// counts the kept rows.
func TestTestFilterKeepsOnlyTestFileRows(t *testing.T) {
	stubLoad(t, scopeFixture)
	out, code := run(t, "target.Symbol", "--test")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	want := "consumers[1]{file,line,via,enclosing}:\n  consumer/consumer_test.go,5,call,check\n" +
		"meta[1]{packages,files,matches,rows,truncated}:\n  2,1,1,1,false\n"
	if !strings.Contains(out, want) {
		t.Fatalf("stdout = %q, want only the test-file row %q", out, want)
	}
}

// --production keeps only the rows whose file is not a _test.go file.
func TestProductionFilterKeepsOnlyProductionRows(t *testing.T) {
	stubLoad(t, scopeFixture)
	out, code := run(t, "target.Symbol", "--production")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	want := "consumers[1]{file,line,via,enclosing}:\n  consumer/consumer.go,5,call,Run\n" +
		"meta[1]{packages,files,matches,rows,truncated}:\n  2,1,1,1,false\n"
	if !strings.Contains(out, want) {
		t.Fatalf("stdout = %q, want only the production row %q", out, want)
	}
}

// The two filters name disjoint row sets, so asking for both is a grammar error, and a
// filter never combines with the --changed form.
func TestFilterGrammarErrorsExitTwo(t *testing.T) {
	for _, args := range [][]string{
		{"target.Symbol", "--production", "--test"},
		{"--changed", "--test"},
	} {
		out, code := run(t, args...)
		if code != 2 || !strings.HasPrefix(out, "usage: bench consumers") {
			t.Errorf("%v = stdout %q exit %d, want a usage line at 2", args, out, code)
		}
	}
}

// An over-cap action replays the query, so it carries the filter the query named.
func TestOverCapActionRepeatsTheFilter(t *testing.T) {
	stubLoad(t, splitFixture(101, 100))
	out, code := run(t, "target.Symbol", "--production")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.HasSuffix(out, "help[1]{cmd,why}:\n  bench consumers target.Symbol --production --full,emit every consumer row\n") {
		t.Fatalf("stdout = %q, want the --full action with the filter", out)
	}
}

// Several symbols answer in one response: each row leads with the operand it answers, in
// operand order, and one meta row and one citation row close the answer.
func TestSeveralSymbolsLabelEachRowAndCiteOnce(t *testing.T) {
	stubLoad(t, pairFixture)
	out, code := run(t, "target.Other", "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	want := "consumers[2]{symbol,file,line,via,enclosing}:\n" +
		"  target.Other,consumer/consumer.go,7,call,UseOther\n" +
		"  target.Symbol,consumer/consumer.go,5,call,UseSymbol\n" +
		"meta[1]{packages,files,matches,rows,truncated}:\n  2,1,2,2,false\n" +
		"citation[1]{sha,state,version,cmd,hash}:\n"
	if !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want the labeled rows, meta, then the citation %q", out, want)
	}
	if n := strings.Count(out, "citation["); n != 1 {
		t.Fatalf("stdout = %q, want one citation block, got %d", out, n)
	}
	if !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
		t.Fatalf("stdout = %q, want the terminal empty help envelope", out)
	}
}

// The row cap applies to each symbol: the over-cap symbol answers its aggregate rows and
// its own --full action, and the other symbol keeps its consumer rows.
func TestSeveralSymbolsCapEachSymbol(t *testing.T) {
	stubLoad(t, func(root string) []fixturePkg {
		other := declPkg(root, "beta", "beta", "package beta\n\nimport \"example.com/target\"\n\nfunc UseOther() { target.Other() }\n")
		return []fixturePkg{pairFixture(root)[0], usePkg(root, "alpha", 201), other}
	})
	out, code := run(t, "target.Symbol", "target.Other")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	for _, want := range []string{
		"consumers[1]{symbol,file,line,via,enclosing}:\n  target.Other,beta/beta.go,5,call,UseOther\n",
		"consumers_packages[1]{symbol,dir,rows}:\n  target.Symbol,alpha,201\n",
		"meta[1]{packages,files,matches,rows,truncated}:\n  3,2,2,202,true\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout = %q, want %q", out, want)
		}
	}
	if !strings.HasSuffix(out, "help[1]{cmd,why}:\n  bench consumers target.Symbol --full,emit every consumer row\n") {
		t.Fatalf("stdout = %q, want one --full action for the over-cap symbol", out)
	}
}

// The candidates form applies to each symbol: an ambiguous bare name answers its labeled
// candidates, and a resolved operand beside it keeps its definitive consumers table.
func TestSeveralSymbolsAnswerCandidatesPerSymbol(t *testing.T) {
	stubLoad(t, ambiguousFixture)
	out, code := run(t, "Symbol", "alpha.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	for _, want := range []string{
		"consumers[0]{symbol,file,line,via,enclosing}:\n",
		"consumers_candidates[2]{symbol,qualified,file,line,kind}:\n" +
			"  Symbol,alpha.Symbol,alpha/alpha.go,3,func\n  Symbol,beta.Symbol,beta/beta.go,3,type\n",
		"meta[1]{packages,files,matches,rows,truncated}:\n  2,0,3,0,false\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout = %q, want %q", out, want)
		}
	}
}

// The help states the several-symbol operand and the filter pair, so the grammar the
// parser enforces is the grammar the usage line advertises.
func TestHelpAdvertisesSeveralSymbolsAndTheFilter(t *testing.T) {
	out, code := run(t, "--help")
	if code != 0 || !strings.HasPrefix(out, "usage: bench consumers <qualified-symbol>... [--production|--test] [--full] |") {
		t.Fatalf("help = %q exit %d, want the several-symbol and filter usage", out, code)
	}
}
