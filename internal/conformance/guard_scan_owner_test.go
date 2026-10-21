package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// guardScanOwnerDir holds the production Go source that owns the guard scan, and
// guardScanOwnerFile is the file in it that holds Scan.
const (
	guardScanOwnerDir  = "internal/guards"
	guardScanOwnerFile = guardScanOwnerDir + "/guards.go"
)

// guardScanOwnerDiagPrefix starts each diagnostic of checkGuardScanLoopOwner.
const guardScanOwnerDiagPrefix = "guard scan loop owner "

// checkGuardScanLoopOwner keeps the guard scan one synchronous loop. It parses each production
// Go file in internal/guards, unused function bodies included, and reports each go statement,
// select statement, channel construction, channel send, and channel receive. It reads syntax
// nodes, so a comment or a string literal never matches. It skips test files, so a test
// goroutine stays legal. A root without internal/guards, such as a linked consumer, holds no
// scan source and gets no diagnostic. Inside that directory, a missing or unreadable guards.go
// is a diagnostic, because an absent owner would otherwise read as a clean pass.
func checkGuardScanLoopOwner(root string) []string {
	dir := filepath.Join(root, filepath.FromSlash(guardScanOwnerDir))
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(dir)
	var production []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
			production = append(production, guardScanOwnerDir+"/"+entry.Name())
		}
	}
	if err != nil || !slices.Contains(production, guardScanOwnerFile) {
		return []string{guardScanOwnerDiagPrefix + guardScanOwnerFile + " is missing or unreadable"}
	}
	var diags []string
	fset := token.NewFileSet()
	for _, rel := range production {
		file, err := parser.ParseFile(fset, filepath.Join(root, filepath.FromSlash(rel)), nil, 0)
		if err != nil {
			diags = append(diags, guardScanOwnerDiagPrefix+"cannot parse "+rel)
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if construct := scanPipelineConstruct(node); construct != "" {
				diags = append(diags, fmt.Sprintf("%s%s:%d: %s", guardScanOwnerDiagPrefix, rel, fset.Position(node.Pos()).Line, construct))
			}
			return true
		})
	}
	return diags
}

// scanPipelineConstruct names the concurrent pipeline construct that node is, or gives "".
// Channel construction is a make call whose type argument is a channel type; a type argument
// cannot reach a function a file declares, so the call names the builtin.
func scanPipelineConstruct(node ast.Node) string {
	switch n := node.(type) {
	case *ast.GoStmt:
		return "go statement"
	case *ast.SelectStmt:
		return "select statement"
	case *ast.SendStmt:
		return "channel send"
	case *ast.UnaryExpr:
		if n.Op == token.ARROW {
			return "channel receive"
		}
	case *ast.CallExpr:
		if ident, ok := n.Fun.(*ast.Ident); ok && ident.Name == "make" && len(n.Args) > 0 {
			if _, chanType := ast.Unparen(n.Args[0]).(*ast.ChanType); chanType {
				return "channel construction"
			}
		}
	}
	return ""
}

// guardScanOwnerDiags runs the registered package-core-guard binding on root and keeps only the
// scan loop owner diagnostics, so each test grades the subcheck through its registered owner.
func guardScanOwnerDiags(t *testing.T, root string) []string {
	t.Helper()
	binding, bound := conformanceChecks["package-core-guard"]
	if !bound {
		t.Fatal("package-core-guard conformance owner is not bound")
	}
	var owned []string
	for _, diag := range binding.run(root, root, registry.Dev) {
		if strings.HasPrefix(diag, guardScanOwnerDiagPrefix) {
			owned = append(owned, diag)
		}
	}
	return owned
}

// guardScanPipelineCopy is the retired worker pipeline of Scan under a new name. No caller
// reaches it, so only a check that grades every function body sees it.
const guardScanPipelineCopy = `package guards

import "context"

func scanPipeline(ctx context.Context, root string) ScanResult {
	type enumeration struct {
		candidates []candidate
		err        error
	}
	enumerated := make(chan enumeration, 1)
	go func() {
		candidates, err := enumerateGuards(ctx, root)
		enumerated <- enumeration{candidates: candidates, err: err}
	}()
	var candidates []candidate
	select {
	case result := <-enumerated:
		candidates = result.candidates
	case <-ctx.Done():
		<-enumerated
		return ScanResult{Status: "incomplete"}
	}
	var rows [][]string
	for _, c := range candidates {
		finished := make(chan [][]string, 1)
		go func(candidate candidate) { finished <- inspectGuard(ctx, root, candidate) }(c)
		rows = append(rows, <-finished...)
	}
	return ScanResult{Rows: rows, Status: "complete"}
}
`

// TestGuardScanLoopOwnerBitesOnPrivateCopy plants a renamed, unused copy of the retired
// pipeline beside a direct loop. The registered package-core-guard check must report each
// construct kind in the copy. A goroutine and a channel in a test file, and pipeline text in a
// comment or a string literal, must not get a diagnostic.
func TestGuardScanLoopOwnerBitesOnPrivateCopy(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)), "package guards\n\n// Scan once ran go func() { ch <- rows }() and a select on ctx.Done().\nfunc Scan() string { return \"go func() { <-done }(); select {}\" }\n")
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/pipeline.go")), guardScanPipelineCopy)
	writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/guards_test.go")), "package guards\n\nfunc waitInTest() {\n\tdone := make(chan struct{})\n\tgo func() { done <- struct{}{} }()\n\tselect {\n\tcase <-done:\n\t}\n}\n")
	diags := guardScanOwnerDiags(t, root)
	for _, construct := range []string{"go statement", "select statement", "channel construction", "channel send", "channel receive"} {
		if !containsDiagnostic(diags, ": "+construct) {
			t.Errorf("want a %s diagnostic on the private pipeline copy, got %v", construct, diags)
		}
	}
	for _, diag := range diags {
		if !strings.Contains(diag, guardScanOwnerDir+"/pipeline.go:") {
			t.Errorf("want diagnostics only on the private pipeline copy, got %q", diag)
		}
	}
}

// TestGuardScanLoopOwnerOptionalSurface grades a linked consumer root. It holds Go source with
// a goroutine outside internal/guards and no internal/guards directory, so the subcheck has no
// scan source to grade and reports nothing.
func TestGuardScanLoopOwnerOptionalSurface(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "internal", "worker", "worker.go"), "package worker\n\nfunc Start(f func()) { go f() }\n")
	if diags := guardScanOwnerDiags(t, root); len(diags) != 0 {
		t.Fatalf("want no scan loop owner diagnostic for a root without %s, got %v", guardScanOwnerDir, diags)
	}
}

// TestGuardScanLoopOwnerMissingOwner amputates guards.go from a root that holds internal/guards.
// An absent file, a directory in its place, and a file that does not parse each get a named
// diagnostic, so a missing owner never reads as a clean pass.
func TestGuardScanLoopOwnerMissingOwner(t *testing.T) {
	missing := guardScanOwnerDiagPrefix + guardScanOwnerFile + " is missing or unreadable"
	for _, tc := range []struct {
		name, want string
		plant      func(t *testing.T, path string)
	}{
		{"absent", missing, func(t *testing.T, path string) {}},
		{"directory", missing, func(t *testing.T, path string) {
			if err := os.MkdirAll(path, 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		{"unparseable", guardScanOwnerDiagPrefix + "cannot parse " + guardScanOwnerFile, func(t *testing.T, path string) {
			writeFixtureFile(t, path, "package guards\n\nfunc Scan( {\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerDir+"/guards_test.go")), "package guards\n")
			tc.plant(t, filepath.Join(root, filepath.FromSlash(guardScanOwnerFile)))
			if diags := guardScanOwnerDiags(t, root); !containsDiagnostic(diags, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, diags)
			}
		})
	}
}

// TestGuardScanLoopOwnerHoldsOnTheLiveTree runs the registered package-core-guard binding on the
// kit root, which holds the guard scan source, and requires no scan loop owner diagnostic.
func TestGuardScanLoopOwnerHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	root := h.KitRoot
	if !exists(filepath.Join(root, filepath.FromSlash(guardScanOwnerFile))) {
		t.Fatalf("the kit root holds no %s, so the live-tree check would grade nothing", guardScanOwnerFile)
	}
	if diags := guardScanOwnerDiags(t, root); len(diags) != 0 {
		t.Fatalf("the kit guard scan is not one synchronous loop:\n%s", strings.Join(diags, "\n"))
	}
}
