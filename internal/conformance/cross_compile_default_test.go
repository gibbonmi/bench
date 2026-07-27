//go:build !stress

package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preprelease"
)

// crossCompileMatrix is a no-op on the every-commit gate: the four-platform cross-compile
// is portability gold-plating that runs only under `go test -tags stress` (see the stress
// variant), with the release workflow as the standing backstop. This keeps the default
// conformance run — the one the gate invokes — free of the ~8s cross-compile cost.
func crossCompileMatrix(root, buildHelper string) []string { return nil }

// residualCheckFunc and crossCompileFunc are the call site the tripwire below reads.
const (
	residualCheckFunc = "checkGoToolchain"
	crossCompileFunc  = "crossCompileMatrix"
)

// TestResidualCheckCallsCrossCompileMatrix is the dev-tier tripwire for a call the dev
// tier cannot otherwise see: here the matrix returns nil, so a residual check that
// stopped calling it behaves identically and every other assertion in the package stays
// green while ship silently loses the four-platform build. The source is the only
// evidence available where the callee does nothing.
//
// The stress-tagged TestResidualCheckKeepsCrossCompile grades a different fact — that
// the matrix, when it is real, reports a refused target — so neither test stands in for
// the other and neither restates the other's knowledge.
func TestResidualCheckCallsCrossCompileMatrix(t *testing.T) {
	h := NewHarness(t)
	decls := packageTestFuncs(t, filepath.Join(h.KitRoot, filepath.FromSlash("internal/conformance")))
	residual, declared := decls[residualCheckFunc]
	if !declared {
		t.Fatalf("%s is no longer declared in this package; the residual check is what carries cross-compile to ship", residualCheckFunc)
	}
	calls := false
	ast.Inspect(residual, func(node ast.Node) bool {
		if ident, isIdent := node.(*ast.Ident); isIdent && ident.Name == crossCompileFunc {
			calls = true
		}
		return !calls
	})
	if !calls {
		t.Fatalf("%s no longer reaches %s; ship would lose the cross-compile matrix with no test observing it", residualCheckFunc, crossCompileFunc)
	}
}

// TestShipConformanceRunNamesDeclaredTests closes the gap a `-run` filter opens on its
// own: naming a test the package does not declare is not an error anywhere, so the step
// runs whatever is left, exits 0, and reports ship green over an assertion that no
// longer exists.
func TestShipConformanceRunNamesDeclaredTests(t *testing.T) {
	h := NewHarness(t)
	decls := packageTestFuncs(t, filepath.Join(h.KitRoot, filepath.FromSlash("internal/conformance")))
	for _, name := range preprelease.ShipConformanceTests {
		if _, declared := decls[name]; !declared {
			t.Errorf("prep-release runs %q at ship, which this package declares nowhere", name)
		}
	}
}

// packageTestFuncs maps every top-level function this package's test files declare to
// its declaration. It parses rather than builds, so the stress-tagged files the dev tier
// excludes are still in the answer — those are exactly the ones no dev run can observe.
func packageTestFuncs(t *testing.T, dir string) map[string]*ast.FuncDecl {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	funcs := map[string]*ast.FuncDecl{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		for _, decl := range file.Decls {
			if fn, isFunc := decl.(*ast.FuncDecl); isFunc && fn.Recv == nil {
				funcs[fn.Name.Name] = fn
			}
		}
	}
	return funcs
}
