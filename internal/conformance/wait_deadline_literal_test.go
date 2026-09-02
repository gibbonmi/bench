package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// derivedDeadlineHelper is the bounds helper that turns an inner window into an outer test
// deadline. Its own argument is the bound the deadline has to outlast, not the deadline, so
// the scanner stops at the call. A literal inside it names the window under test, and every
// migrated wait in the tree spells its deadline that way.
const derivedDeadlineHelper = "TestDeadline"

// checkWaitDeadlineLiterals fails a test wait whose deadline is a numeric duration. The
// marker-wait check grades one helper's slow leg; this one grades the general shape behind
// it. A wall-clock guess in an outer wait passes on an idle machine and flakes on a loaded
// one, and the value it should carry is derivable: bounds.TestDeadline(inner) returns a
// window strictly greater than the bound the wait contains.
//
// A poll interval is the exception the shape needs. A tick fires many times inside the
// wait, so its literal is a sampling rate, not a deadline, and it stays legal.
//
// The check reads the deadline argument, as the marker-wait check does, and it does not
// resolve a named window back to its declaration. A window a package names is already a
// reviewed policy value, and the defect this catches is the number written at the wait.
func checkWaitDeadlineLiterals(root string) []string {
	var diags []string
	for _, top := range []string{"cmd", "internal"} {
		_ = filepath.WalkDir(filepath.Join(root, top), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || !strings.HasSuffix(path, "_test.go") {
				return nil
			}
			body := readIfExists(path)
			if !strings.Contains(body, "time.After(") && !strings.Contains(body, "time.Now().Add(") {
				return nil
			}
			diags = append(diags, waitDeadlineLiteralDiags(path, slashRel(root, path), body)...)
			return nil
		})
	}
	return uniqueSorted(diags)
}

// waitDeadlineLiteralDiags reports each literal deadline in one test file. Two spellings
// carry a deadline: time.After feeds a select's timeout arm, and time.Now().Add fixes the
// instant a poll loop gives up at. A negative offset from now is a backdated fixture
// timestamp, never a deadline, so it is not a wait at all.
func waitDeadlineLiteralDiags(path, rel, body string) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for wait deadlines: " + err.Error()}
	}
	arms, polls := map[ast.Stmt]bool{}, map[ast.Node]bool{}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		// A select arm holds its receive in an expression statement too, and that arm is the
		// deadline the wait races. Preorder visits the clause before the statement, so the arm
		// is known by the time the statement decides whether it is a poll.
		if clause, ok := node.(*ast.CommClause); ok {
			arms[clause.Comm] = true
			return true
		}
		// A bare `<-time.After(d)` receive statement waits out one tick and continues. It
		// paces a poll loop rather than bounding it, so it carries a sampling rate. The
		// statement is visited before the call it holds, so the exemption is recorded first.
		if statement, ok := node.(*ast.ExprStmt); ok {
			receive, ok := statement.X.(*ast.UnaryExpr)
			if ok && receive.Op == token.ARROW && !arms[statement] {
				polls[receive.X] = true
			}
			return true
		}
		call, ok := node.(*ast.CallExpr)
		if !ok || polls[call] || len(call.Args) != 1 {
			return true
		}
		deadline, spelling := call.Args[0], waitDeadlineSpelling(call)
		if spelling == "" || !literalDeadline(deadline) {
			return true
		}
		diags = append(diags, fmt.Sprintf(
			"%s:%d waits on duration literal %s in %s; derive the window from bounds.TestDeadline(inner) so the wait outlasts the bound it contains",
			rel, fset.Position(call.Pos()).Line, expressionText(fset, deadline), spelling))
		return true
	})
	return diags
}

// waitDeadlineSpelling names the wait a call opens, or "" when the call is neither.
func waitDeadlineSpelling(call *ast.CallExpr) string {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	if pkg, ok := selector.X.(*ast.Ident); ok && pkg.Name == "time" && selector.Sel.Name == "After" {
		return "time.After"
	}
	if selector.Sel.Name != "Add" || !isTimeNow(selector.X) || negativeOffset(call.Args[0]) {
		return ""
	}
	return "time.Now().Add"
}

func isTimeNow(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Now" {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == "time"
}

// negativeOffset reports whether an offset from now points backwards. The sign sits on the
// leftmost leaf, because `-8 * 24 * time.Hour` parses as a product whose first factor is the
// negated number, so the walk descends the left operand rather than reading the top node.
func negativeOffset(expr ast.Expr) bool {
	switch node := expr.(type) {
	case *ast.ParenExpr:
		return negativeOffset(node.X)
	case *ast.BinaryExpr:
		return negativeOffset(node.X)
	case *ast.UnaryExpr:
		return node.Op == token.SUB
	}
	return false
}

// literalDeadline reports whether a deadline expression carries a wall-clock number. It
// reuses containsNumericLiteral, the marker-wait scanner, and adds one stop rule on top of
// it: a TestDeadline call is already derived, and descending into it would read the inner
// bound as the defect it exists to fix.
func literalDeadline(expr ast.Expr) bool {
	switch node := expr.(type) {
	case *ast.ParenExpr:
		return literalDeadline(node.X)
	case *ast.UnaryExpr:
		return literalDeadline(node.X)
	case *ast.BinaryExpr:
		return literalDeadline(node.X) || literalDeadline(node.Y)
	case *ast.CallExpr:
		if derivedDeadline(node) {
			return false
		}
	}
	return containsNumericLiteral(expr)
}

// derivedDeadline recognises the bounds helper through either spelling. The bounds package
// calls it unqualified in its own tests, and every other package reaches it as
// bounds.TestDeadline.
func derivedDeadline(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name == derivedDeadlineHelper
	case *ast.SelectorExpr:
		return fun.Sel.Name == derivedDeadlineHelper
	}
	return false
}

// TestWaitDeadlineLiteralsBites is the recorded bite proof for checkWaitDeadlineLiterals. It
// runs against a synthetic tree, not the repo, and walks the four states that decide the
// check: a literal deadline, a derived deadline paced by a tick, a derived deadline whose
// inner bound is a literal, and a backdated fixture timestamp.
func TestWaitDeadlineLiteralsBites(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	const rel = "internal/example/example_test.go"

	write(rel, "package example\n\nimport (\n\t\"testing\"\n\t\"time\"\n)\n\nfunc TestX(t *testing.T) {\n\tdone := make(chan struct{})\n\tselect {\n\tcase <-done:\n\tcase <-time.After(5 * time.Second):\n\t\tt.Fatal(\"timeout\")\n\t}\n}\n")
	diags := checkWaitDeadlineLiterals(root)
	if len(diags) != 1 || !strings.Contains(diags[0], rel+":12 waits on duration literal 5 * time.Second in time.After") {
		t.Fatalf("literal wait deadline: want one diagnostic naming it, got %v", diags)
	}
	if !strings.Contains(diags[0], "bounds.TestDeadline(inner)") {
		t.Fatalf("diagnostic does not name the derivation: %q", diags[0])
	}

	// A derived deadline paced by a 10ms tick is the shape the check exists to protect. The
	// tick's literal is a sampling rate inside the wait, and only the outer window is graded.
	write(rel, "package example\n\nimport (\n\t\"testing\"\n\t\"time\"\n\n\t\"github.com/gibbonmi/bench/internal/bounds\"\n)\n\nfunc TestX(t *testing.T) {\n\tdone := make(chan struct{})\n\tticker := time.NewTicker(10 * time.Millisecond)\n\tdefer ticker.Stop()\n\tfor {\n\t\tselect {\n\t\tcase <-done:\n\t\t\treturn\n\t\tcase <-time.After(bounds.TestDeadline(0)):\n\t\t\tt.Fatal(bounds.TestTimeoutVerdict(\"done\", bounds.TestDeadline(0)))\n\t\tdefault:\n\t\t}\n\t\t<-time.After(10 * time.Millisecond)\n\t}\n}\n")
	if diags := checkWaitDeadlineLiterals(root); len(diags) != 0 {
		t.Fatalf("derived deadline with a 10ms tick: want no diagnostics, got %v", diags)
	}

	// The inner bound a derivation outlasts is the window under test, not the deadline. Every
	// migrated wait spells its window bounds.TestDeadline(0), so descending would red them all.
	write(rel, "package example\n\nimport (\n\t\"testing\"\n\t\"time\"\n\n\t\"github.com/gibbonmi/bench/internal/bounds\"\n)\n\nfunc TestX(t *testing.T) {\n\tdone := make(chan struct{})\n\tdeadline := time.Now().Add(bounds.TestDeadline(5 * time.Second))\n\tfor time.Now().Before(deadline) {\n\t\tselect {\n\t\tcase <-done:\n\t\t\treturn\n\t\tdefault:\n\t\t}\n\t}\n\tt.Fatal(\"timeout\")\n}\n")
	if diags := checkWaitDeadlineLiterals(root); len(diags) != 0 {
		t.Fatalf("literal inner bound inside a derivation: want no diagnostics, got %v", diags)
	}

	// A negative offset from now is a backdated fixture timestamp. It is already in the past,
	// so no wait can expire at it, and the check must not read it as a deadline.
	write(rel, "package example\n\nimport (\n\t\"testing\"\n\t\"time\"\n)\n\nfunc TestX(t *testing.T) {\n\tstale := time.Now().Add(-8 * 24 * time.Hour)\n\tif stale.IsZero() {\n\t\tt.Fatal(\"zero\")\n\t}\n}\n")
	if diags := checkWaitDeadlineLiterals(root); len(diags) != 0 {
		t.Fatalf("backdated fixture timestamp: want no diagnostics, got %v", diags)
	}

	// A poll loop's own give-up instant is a deadline wherever it is spelled, so the
	// time.Now().Add half bites on a forward literal.
	write(rel, "package example\n\nimport (\n\t\"testing\"\n\t\"time\"\n)\n\nfunc TestX(t *testing.T) {\n\tdeadline := time.Now().Add(30 * time.Second)\n\tif time.Now().Before(deadline) {\n\t\tt.Fatal(\"timeout\")\n\t}\n}\n")
	if diags := checkWaitDeadlineLiterals(root); len(diags) != 1 || !strings.Contains(diags[0], "waits on duration literal 30 * time.Second in time.Now().Add") {
		t.Fatalf("literal poll-loop deadline: want one diagnostic, got %v", diags)
	}
}
