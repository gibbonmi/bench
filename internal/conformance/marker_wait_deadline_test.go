package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// markerWaitHelper is the two-leg marker wait every cross-package caller reaches through
// its package qualifier; the helper's own package calls it unqualified against a fake
// clock, where a literal is the subject under test rather than a wall-clock deadline.
const markerWaitHelper = "WaitForTwoLegMarkers"

// slowDeadlineArg is the position of the slow leg's deadline. That leg is the one that
// has to outlast a window the test already contains, so it is the argument a numeric
// literal makes a coin flip; the fast leg is a startup handshake bounded by nothing else.
const slowDeadlineArg = 3

// checkMarkerWaitDeadlines fails a cross-package marker wait whose slow deadline is spelled
// as a numeric duration. A literal there silently ties the outer wait to a window it must
// outlast — the exact regression that made an outer deadline equal to its inner one — and
// inspecting the helper's own definition would never see it, because the defect lives at
// the call site.
func checkMarkerWaitDeadlines(root string) []string {
	var diags []string
	for _, top := range []string{"cmd", "internal"} {
		_ = filepath.WalkDir(filepath.Join(root, top), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}
			body := readIfExists(path)
			if !strings.Contains(body, markerWaitHelper) {
				return nil
			}
			diags = append(diags, markerWaitDeadlineDiags(path, slashRel(root, path), body)...)
			return nil
		})
	}
	return uniqueSorted(diags)
}

func markerWaitDeadlineDiags(path, rel, body string) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for marker-wait deadlines: " + err.Error()}
	}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != markerWaitHelper {
			return true
		}
		if _, ok := selector.X.(*ast.Ident); !ok || len(call.Args) <= slowDeadlineArg {
			return true
		}
		deadline := call.Args[slowDeadlineArg]
		if !containsNumericLiteral(deadline) {
			return true
		}
		diags = append(diags, rel+" passes duration literal "+expressionText(fset, deadline)+" as the "+markerWaitHelper+" slow deadline; derive it from the bound it must outlast (bounds.TestDeadline)")
		return true
	})
	return diags
}

func containsNumericLiteral(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if ok && (literal.Kind == token.INT || literal.Kind == token.FLOAT) {
			found = true
			return false
		}
		return !found
	})
	return found
}
