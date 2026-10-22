// Package testcensus owns the walk that the package parallel censuses share. It parses
// a directory's Go files with the Go AST, finds the top-level tests, reads a test's own
// t.Parallel() call, resolves a write target to the identifier it writes, and plants the
// synthetic file sets the census tests grade. Each census keeps its own rules. Only tests
// import this package.
package testcensus

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ParseTestFiles parses every regular _test.go file in dir. A special file, for
// example a FIFO, is not regular and is skipped, so the walk never blocks on a
// read.
func ParseTestFiles(dir string) ([]*ast.File, *token.FileSet, []string, error) {
	return ParseGoFiles(dir, func(name string) bool { return strings.HasSuffix(name, "_test.go") })
}

// ParseGoFiles parses every regular file in dir whose name want accepts. It returns the
// files, their file set, and the base name of each file.
func ParseGoFiles(dir string, want func(string) bool) ([]*ast.File, *token.FileSet, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, nil, err
	}
	fset := token.NewFileSet()
	var files []*ast.File
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if !want(name) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, nil, nil, err
		}
		if !info.Mode().IsRegular() {
			continue
		}
		parsed, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, nil, nil, err
		}
		files = append(files, parsed)
		names = append(names, name)
	}
	return files, fset, names, nil
}

// IsTopLevelTest reports whether decl is a TestXxx(t *testing.T) function.
// TestMain is neither eligible nor serial, so it is not a top-level test.
func IsTopLevelTest(decl *ast.FuncDecl) bool {
	if decl.Recv != nil || decl.Body == nil || decl.Name == nil {
		return false
	}
	name := decl.Name.Name
	if !strings.HasPrefix(name, "Test") || name == "TestMain" || len(name) == len("Test") {
		return false
	}
	if next := rune(name[len("Test")]); next >= 'a' && next <= 'z' {
		return false
	}
	params := decl.Type.Params.List
	if len(params) != 1 {
		return false
	}
	star, ok := params[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	selector, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == "testing" && selector.Sel.Name == "T"
}

// CallsParallelDirectly reports whether the body of decl calls t.Parallel() on
// its own testing parameter, outside every nested function literal. A subtest's
// t.Parallel() inside a closure is the subtest's call, not the parent's.
func CallsParallelDirectly(decl *ast.FuncDecl) bool {
	params := decl.Type.Params.List
	if len(params) != 1 || len(params[0].Names) != 1 {
		return false
	}
	receiver := params[0].Names[0].Name
	found := false
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncLit); ok {
			return false // a closure's call belongs to the closure
		}
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "Parallel" {
			return true
		}
		if ident, ok := selector.X.(*ast.Ident); ok && ident.Name == receiver {
			found = true
		}
		return true
	})
	return found
}

// WriteTargets returns the expressions node writes: the left side of an assignment
// that is not a := declaration, and the operand of ++ or --. Any other node writes
// nothing.
func WriteTargets(node ast.Node) []ast.Expr {
	switch write := node.(type) {
	case *ast.AssignStmt:
		if write.Tok != token.DEFINE {
			return write.Lhs
		}
	case *ast.IncDecStmt:
		return []ast.Expr{write.X}
	}
	return nil
}

// AssignRoot resolves a write target to the identifier it ultimately writes, and
// to the selector name that qualifies that identifier. The walk strips the index,
// pointer, and parenthesis wrappers and the outer selectors, so pkg.Slice[0],
// pkg.Var.Field, and (*pkg.Ptr).Field all resolve to the root pkg and the selector
// the root qualifies. The selector name is empty when the target is a bare
// identifier or a wrapper around one, and the root name is empty when the target
// roots in no identifier at all. A write through any of these shapes reaches the
// same storage as a write to the bare name.
func AssignRoot(target ast.Expr) (root string, sel string) {
	for {
		switch node := target.(type) {
		case *ast.ParenExpr:
			target = node.X
		case *ast.StarExpr:
			target = node.X
		case *ast.IndexExpr:
			target = node.X
		case *ast.SelectorExpr:
			sel = node.Sel.Name
			target = node.X
		case *ast.Ident:
			return node.Name, sel
		default:
			return "", ""
		}
	}
}

// Plant writes one synthetic file set into a temporary directory and returns the
// directory, so a census test grades planted sources instead of the live package.
func Plant(t testing.TB, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatalf("plant %s: %v", name, err)
		}
	}
	return dir
}
