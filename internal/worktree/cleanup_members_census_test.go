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

	"github.com/gibbonmi/bench/internal/intent"
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
// on that literal as a top-level result on the live path, or a loop whose body skips, stops,
// or appends as a visit loop does, directly or through a function literal it holds or calls.
func traversalCopyProblems(fset *token.FileSet, decl *ast.FuncDecl) []string {
	var problems []string
	bound := map[string]bool{}
	closures := map[string]*ast.FuncLit{}
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		if n, ok := node.(*ast.AssignStmt); ok {
			for i, rhs := range n.Rhs {
				name := ""
				if i < len(n.Lhs) {
					if ident, ok := n.Lhs[i].(*ast.Ident); ok {
						name = ident.Name
					}
				}
				switch value := rhs.(type) {
				case *ast.CompositeLit:
					if isCleanupTraversalType(value.Type) && name != "" {
						bound[name] = true
						problems = append(problems, traversalHookProblems(value)...)
					}
				case *ast.FuncLit:
					if name != "" {
						closures[name] = value
					}
				}
			}
		}
		return true
	})
	runs := 0
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CallExpr:
			if traversalRun(n, bound) {
				runs++
			}
		case *ast.ForStmt:
			problems = append(problems, visitLoopProblem(fset, n.Pos(), n.Body, closures, bound)...)
		case *ast.RangeStmt:
			problems = append(problems, visitLoopProblem(fset, n.Pos(), n.Body, closures, bound)...)
		}
		return true
	})
	if len(bound) != 1 {
		problems = append(problems, fmt.Sprintf("binds %d cleanupTraversal literals, want 1", len(bound)))
	}
	if runs != 1 {
		problems = append(problems, fmt.Sprintf("runs its cleanupTraversal %d times, want 1", runs))
	}
	if !liveTopLevelRun(decl.Body, bound) {
		problems = append(problems, "does not return or assign its cleanupTraversal run as a top-level statement on the live path")
	}
	return problems
}

// traversalRun reports whether call is the run method of a bound cleanupTraversal.
func traversalRun(call *ast.CallExpr, bound map[string]bool) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "run" {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	return ok && bound[ident.Name]
}

// liveTopLevelRun reports whether a statement of body itself returns or assigns the run of a
// bound cleanupTraversal, with no unconditional return ahead of it. A run inside a branch,
// a loop, or a closure is one a condition can skip.
func liveTopLevelRun(body *ast.BlockStmt, bound map[string]bool) bool {
	for _, stmt := range body.List {
		var results []ast.Expr
		switch s := stmt.(type) {
		case *ast.ReturnStmt:
			results = s.Results
		case *ast.AssignStmt:
			results = s.Rhs
		}
		for _, result := range results {
			if call, ok := result.(*ast.CallExpr); ok && traversalRun(call, bound) {
				return true
			}
		}
		if _, ok := stmt.(*ast.ReturnStmt); ok {
			return false
		}
	}
	return false
}

// visitLoopProblem reports a loop at pos whose body skips a member, stops the loop, or appends
// a row: the work the shared traversal owns. A function literal in the body, a call to a
// closure the adapter binds, and a call on the bound traversal all count as that body.
func visitLoopProblem(fset *token.FileSet, pos token.Pos, body *ast.BlockStmt, closures map[string]*ast.FuncLit, bound map[string]bool) []string {
	found := ""
	var visit func(ast.Node) bool
	visit = func(node ast.Node) bool {
		if found != "" {
			return false
		}
		switch n := node.(type) {
		case *ast.ReturnStmt:
			found = "returns"
		case *ast.BranchStmt:
			found = n.Tok.String() + "s"
		case *ast.CallExpr:
			switch fun := n.Fun.(type) {
			case *ast.Ident:
				if fun.Name == "append" {
					found = "appends"
				} else if lit := closures[fun.Name]; lit != nil {
					delete(closures, fun.Name)
					ast.Inspect(lit.Body, visit)
					closures[fun.Name] = lit
				}
			case *ast.SelectorExpr:
				if ident, ok := fun.X.(*ast.Ident); ok && bound[ident.Name] {
					found = "calls its cleanupTraversal"
				}
			}
		}
		return true
	}
	ast.Inspect(body, visit)
	if found == "" {
		return nil
	}
	return []string{fmt.Sprintf("member-visit loop at line %d %s inside its body", fset.Position(pos).Line, found)}
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

// TestCleanupTraversalFingerprintKnownAnswer pins each set version's digest of one fixed member list
// through the production input builder and hash. The literal digests are independent of the
// hash route on purpose: a uniform change of the hashed bytes, such as a dropped version,
// leaves every rendered record green, and only a known answer sees it.
func TestCleanupTraversalFingerprintKnownAnswer(t *testing.T) {
	t.Parallel()
	options := CleanupOptions{DiscardBranch: true, Unclaimed: true}
	verdict := refVerdict{oid: "0123abcd", class: refClass("unique"), holder: "none"}
	branch := unclaimedAssignmentBranch{ref: "refs/heads/bench/a1", refVerdict: verdict}
	explicit := explicitCleanupSet{
		rows:       []explicitCleanupRow{{assignment: intent.Assignment{ID: "a1"}, plan: CleanupPlan{Target: "/pool/a1", Action: ActionRemoved}, targetFingerprint: "t1"}},
		unrecorded: []unrecordedCleanupRow{{unclaimedAssignmentBranch: branch, discarded: "refs/bench/discarded/a1"}},
	}
	landed := []landedCleanupRow{{assignment: intent.Assignment{ID: "a1", Worktree: "/pool/a1"}, plan: CleanupPlan{Action: ActionRemoved, Tracked: "clean"}, headOID: "4567cdef", lease: "l1"}}
	for _, c := range []struct{ version, got, want string }{
		{explicitSetFingerprintVersion, fingerprintInputs(explicitSetFingerprintVersion, explicitFingerprintInputs(explicit, options)), "7db38ed1ed843957cd6fa7bb21340f52dc030b6960a81c70eb5cc2daf566d9fc"},
		{landedSetFingerprintVersion, fingerprintInputs(landedSetFingerprintVersion, landedFingerprintInputs(landed, options)), "68da89b4ebc201d2ac3079589026eca6fcd2f2aa178a1904959033a47554c02c"},
		{unclaimedAssignmentFingerprintVersion, fingerprintInputs(unclaimedAssignmentFingerprintVersion, unclaimedFingerprintInputs([]unclaimedAssignmentBranch{branch}, options)), "3c8dab601e3b3dcefed187d755e607837c831a7f11720ff523cd2f715e8c600a"},
	} {
		if c.got != c.want {
			t.Errorf("%s fingerprint = %s, want the known answer %s", c.version, c.got, c.want)
		}
	}
}
