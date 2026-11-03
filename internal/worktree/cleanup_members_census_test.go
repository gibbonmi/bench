// The copy-survival census of the shared cleanup traversal. Each cleanup apply mode adapter
// must hand its members to cleanupTraversal and run it, and no adapter may keep a loop of its
// own that visits members: one that skips, stops, or accumulates outcome rows.
package worktree

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"testing"
)

// cleanupTraversalAdapters is every mode adapter that visits cleanup members.
var cleanupTraversalAdapters = []string{"applyExplicitSet", "applyUnrecordedRows", "applyLandedSet", "applyUnclaimedAssignmentSet"}

// TestCleanupTraversalCopySurvival is R25. It reads the package sources and fails when an
// adapter no longer runs a complete cleanupTraversal or still holds a member-visit loop.
func TestCleanupTraversalCopySurvival(t *testing.T) {
	t.Parallel()
	files, fset, _, err := parseGoFiles(".", isSourceFile)
	mustNoError(t, err)
	funcs := topLevelFuncs(files)
	for _, name := range cleanupTraversalAdapters {
		decl, ok := funcs[name]
		if !ok {
			t.Errorf("adapter %s is not declared in the package sources", name)
			continue
		}
		for _, problem := range traversalCopyProblems(fset, decl) {
			t.Errorf("%s: %s", name, problem)
		}
	}
}

// traversalCopyProblems reports what keeps decl from delegating every visit to the shared
// traversal: no traversal literal with all four mode hooks and a visit closure, no run call
// on that literal, or a loop whose body skips, stops, or appends as a visit loop does.
func traversalCopyProblems(fset *token.FileSet, decl *ast.FuncDecl) []string {
	var problems []string
	bound := map[string]bool{}
	runs := 0
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.AssignStmt:
			for i, rhs := range n.Rhs {
				if lit, ok := rhs.(*ast.CompositeLit); ok && isCleanupTraversalType(lit.Type) && i < len(n.Lhs) {
					if ident, ok := n.Lhs[i].(*ast.Ident); ok {
						bound[ident.Name] = true
					}
					problems = append(problems, traversalHookProblems(lit)...)
				}
			}
		case *ast.CallExpr:
			if sel, ok := n.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "run" {
				if ident, ok := sel.X.(*ast.Ident); ok && bound[ident.Name] {
					runs++
				}
			}
		case *ast.ForStmt:
			problems = append(problems, visitLoopProblem(fset, n.Pos(), n.Body)...)
		case *ast.RangeStmt:
			problems = append(problems, visitLoopProblem(fset, n.Pos(), n.Body)...)
		}
		return true
	})
	if len(bound) != 1 {
		problems = append(problems, fmt.Sprintf("binds %d cleanupTraversal literals, want 1", len(bound)))
	}
	if runs != 1 {
		problems = append(problems, fmt.Sprintf("runs its cleanupTraversal %d times, want 1", runs))
	}
	return problems
}

// isCleanupTraversalType reports whether expr is an instance of cleanupTraversal.
func isCleanupTraversalType(expr ast.Expr) bool {
	index, ok := expr.(*ast.IndexExpr)
	if !ok {
		return false
	}
	ident, ok := index.X.(*ast.Ident)
	return ok && ident.Name == "cleanupTraversal"
}

// traversalHookProblems reports each cleanupTraversal field the literal leaves unset, and a visit hook that
// is not a closure the adapter writes in place.
func traversalHookProblems(lit *ast.CompositeLit) []string {
	set := map[string]ast.Expr{}
	for _, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.Ident); ok {
				set[key.Name] = kv.Value
			}
		}
	}
	var problems []string
	hooks := reflect.TypeFor[cleanupTraversal[struct{}]]()
	for i := range hooks.NumField() {
		if hook := hooks.Field(i).Name; set[hook] == nil {
			problems = append(problems, "cleanupTraversal leaves its "+hook+" hook unset")
		}
	}
	if _, ok := set["visit"].(*ast.FuncLit); set["visit"] != nil && !ok {
		problems = append(problems, "cleanupTraversal visit hook is not a closure")
	}
	return problems
}

// visitLoopProblem reports a loop at pos whose body, outside any nested closure, skips a
// member, stops the loop, or appends a row: the work the shared traversal owns.
func visitLoopProblem(fset *token.FileSet, pos token.Pos, body *ast.BlockStmt) []string {
	found := ""
	ast.Inspect(body, func(node ast.Node) bool {
		if found != "" {
			return false
		}
		switch n := node.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt:
			found = "returns"
		case *ast.BranchStmt:
			found = n.Tok.String() + "s"
		case *ast.CallExpr:
			if calleeName(n) == "append" {
				found = "appends"
			}
		}
		return true
	})
	if found == "" {
		return nil
	}
	return []string{fmt.Sprintf("member-visit loop at line %d %s inside its body", fset.Position(pos).Line, found)}
}
