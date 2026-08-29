package consumers

import (
	"fmt"
	"strings"
	"testing"
)

// The two help clauses are pinned here as literals rather than read from the package
// constants. An expectation that reads the same constant cannot go red when the clause is
// reworded, and the reworded clause is exactly the failure CS16 and CS24 exist to catch.
const (
	wantPromiseClause   = "consumers identifies resolved reference edges; it does not bless a seam — projects/<name>.md owns that."
	wantSoundnessClause = "sound for static Go references only: reflection, go:linkname, plugin, and exec edges are invisible; the default build context is the graded one."
)

// stubLoad swaps the package-internal loader seam for in-process typed fixtures, so a
// command test spawns no subprocess. build receives the repository root the command
// resolved, and names its files under it, so the rendered file cells stay repo-relative.
func stubLoad(t *testing.T, build func(root string) []fixturePkg) {
	t.Helper()
	original := load
	load = func(dir string, _ ...string) ([]*Package, error) {
		return typecheckFixture(t, build(dir)), nil
	}
	t.Cleanup(func() { load = original })
}

// targetPkg is the queried declaration every command fixture resolves to.
func targetPkg(root string) fixturePkg {
	return fixturePkg{path: "example.com/target", files: map[string]string{
		root + "/target/target.go": "package target\n\nfunc Symbol() {}\n"}}
}

// usePkg plants count references to the target symbol, one per line, in one package.
func usePkg(root, name string, count int) fixturePkg {
	var b strings.Builder
	fmt.Fprintf(&b, "package %s\n\nimport \"example.com/target\"\n", name)
	for i := 0; i < count; i++ {
		fmt.Fprintf(&b, "\nfunc Use%d() { target.Symbol() }\n", i)
	}
	return fixturePkg{path: "example.com/" + name, files: map[string]string{
		root + "/" + name + "/" + name + ".go": b.String()}}
}

// referencesFixture is the command fixture shape: the target plus one consumer package
// holding count references.
func referencesFixture(count int) func(string) []fixturePkg {
	return func(root string) []fixturePkg {
		if count == 0 {
			return []fixturePkg{targetPkg(root)}
		}
		return []fixturePkg{targetPkg(root), usePkg(root, "consumer", count)}
	}
}

// splitFixture spreads references over two consumer packages, so an aggregate row per
// package is observably a grouping rather than one total.
func splitFixture(first, second int) func(string) []fixturePkg {
	return func(root string) []fixturePkg {
		return []fixturePkg{targetPkg(root), usePkg(root, "alpha", first), usePkg(root, "beta", second)}
	}
}

func run(t *testing.T, args ...string) (string, int) {
	t.Helper()
	return Command(args)
}

// CS6 (story 5): a matched symbol with zero references answers with the definitive empty
// table rather than an error, so an absence claim is warranted.
func TestZeroReferenceSymbolEmitsDefinitiveEmptyTable(t *testing.T) {
	stubLoad(t, referencesFixture(0))
	out, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "consumers[0]{file,line,via,enclosing}:\n") {
		t.Fatalf("stdout = %q, want the definitive empty consumers table", out)
	}
}

// CS14 (story 32): the common case is complete, and the meta accounting states the count.
func TestThreeReferencesEmitThreeRowsAndMetaRowCount(t *testing.T) {
	stubLoad(t, referencesFixture(3))
	out, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "consumers[3]{file,line,via,enclosing}:\n") {
		t.Fatalf("stdout = %q, want three consumer rows", out)
	}
	if !strings.Contains(out, "meta[1]{packages,files,matches,rows,truncated}:\n  2,1,1,3,false\n") {
		t.Fatalf("stdout = %q, want meta rows=3", out)
	}
}

// CS15 (story 14): an unknown flag is a usage error on stdout at exit 2.
func TestUnknownFlagExitsTwoWithUsageOnStdout(t *testing.T) {
	out, code := run(t, "target.Symbol", "--unknown-probe")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "usage: bench consumers") {
		t.Fatalf("stdout = %q, want a usage line", out)
	}
}

// CS16 (story 15): the help text states the static-analysis limit verbatim.
func TestHelpCarriesSoundnessClause(t *testing.T) {
	out, code := run(t, "--help")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, wantSoundnessClause) {
		t.Fatalf("help = %q, want the soundness clause %q", out, wantSoundnessClause)
	}
}

// CS24 (story 15): the help text states the identifies-edges promise verbatim.
func TestHelpCarriesPromiseClause(t *testing.T) {
	out, code := run(t, "--help")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, wantPromiseClause) {
		t.Fatalf("help = %q, want the promise clause %q", out, wantPromiseClause)
	}
}

// CS19 (story 14): a symbol result is a terminal read, so its envelope is honestly empty.
func TestTerminalSymbolResultEndsWithEmptyHelpEnvelope(t *testing.T) {
	stubLoad(t, referencesFixture(3))
	out, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
		t.Fatalf("stdout = %q, want a terminal empty help envelope", out)
	}
}

// CS20 (story 32): an over-cap default aggregates per package, says it truncated, and
// offers exactly the one action that returns the complete set.
func TestOverCapDefaultAggregatesAndOffersFull(t *testing.T) {
	stubLoad(t, splitFixture(101, 100))
	out, code := run(t, "target.Symbol")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "consumers_packages[2]{dir,rows}:\n  alpha,101\n  beta,100\n") {
		t.Fatalf("stdout = %q, want per-package aggregate rows", out)
	}
	if strings.Contains(out, "consumers[") {
		t.Fatalf("stdout = %q, want no per-reference rows over the cap", out)
	}
	if !strings.Contains(out, ",201,true\n") {
		t.Fatalf("stdout = %q, want rows=201 and truncated=true in meta", out)
	}
	if !strings.HasSuffix(out, "help[1]{cmd,why}:\n  bench consumers target.Symbol --full,emit every consumer row\n") {
		t.Fatalf("stdout = %q, want exactly one --full action", out)
	}
}

// CS21 (story 13): --full is the complete set, whatever the cap.
func TestFullEmitsEveryRowPastTheCap(t *testing.T) {
	stubLoad(t, splitFixture(101, 100))
	out, code := run(t, "target.Symbol", "--full")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; out=%q", code, out)
	}
	if !strings.Contains(out, "consumers[201]{file,line,via,enclosing}:\n") {
		t.Fatalf("stdout = %q, want every planted row", out)
	}
	if !strings.Contains(out, ",201,false\n") {
		t.Fatalf("stdout = %q, want truncated=false under --full", out)
	}
}

// CS23 (story 14): every help spelling prints usage on stdout at exit 0.
func TestHelpSpellingsPrintUsageOnStdout(t *testing.T) {
	for _, spelling := range []string{"--help", "-h", "help"} {
		out, code := run(t, spelling)
		if code != 0 || !strings.HasPrefix(out, "usage: bench consumers") {
			t.Errorf("%s = stdout %q exit %d, want a usage line at 0", spelling, out, code)
		}
	}
}
