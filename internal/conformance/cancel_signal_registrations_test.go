package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// cancelSignalOwner is the one package that spells the termination-signal set. Its own
// registration reads the set from the variable beside it, so the walk exempts it.
const cancelSignalOwner = "internal/subprocess"

// cancelSignalSet is the expression every production registration spreads.
const cancelSignalSet = "subprocess.CancelSignals"

// checkCancelSignalRegistrations grades the single-source rule for the signals a Bench
// command traps. A production owner that names its own signals drifts from
// subprocess.CancelSignals, and the drift is silent: the command still runs, and the
// signal it declines to trap leaks the detached process group it owns.
//
// The check inspects call sites rather than file bytes, so the token in a comment or a
// string is not a registration. Test files are exempt, because a fixture registers
// whatever signal its own scenario sends.
func checkCancelSignalRegistrations(root string) []string {
	var diags []string
	for _, top := range []string{"cmd", "internal"} {
		base := filepath.Join(root, top)
		_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel := slashRel(root, path)
			if strings.HasPrefix(rel, cancelSignalOwner+"/") {
				return nil
			}
			diags = append(diags, cancelSignalRegistrationDiags(path, rel)...)
			return nil
		})
	}
	return uniqueSorted(diags)
}

// cancelSignalRegistrationDiags reports each non-conforming registration in one file. The
// substring prefilter keeps the whole-tree cost at one scan per file, because both
// spellings share the "signal.Notify" prefix.
func cancelSignalRegistrationDiags(path, rel string) []string {
	body := readIfExists(path)
	if !strings.Contains(body, "signal.Notify") {
		return nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for cancel-signal registrations: " + err.Error()}
	}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		spelling := cancelSignalSpelling(call)
		if spelling == "" || registersCancelSignals(call, fset) {
			return true
		}
		diags = append(diags, rel+":"+strconv.Itoa(fset.Position(call.Pos()).Line)+
			" registers "+spelling+" with "+registeredSignalText(call, fset)+
			" instead of "+cancelSignalSet)
		return true
	})
	return diags
}

// cancelSignalSpellings is the closed set of calls that register a termination handler.
// signal.Notify takes a channel first and signal.NotifyContext takes a parent context
// first, so the signals begin at the same index in both.
var cancelSignalSpellings = []string{"Notify", "NotifyContext"}

// cancelSignalSpelling names the registration a call opens, or "" when the call opens
// none.
func cancelSignalSpelling(call *ast.CallExpr) string {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok || pkg.Name != "signal" || len(call.Args) == 0 {
		return ""
	}
	for _, name := range cancelSignalSpellings {
		if selector.Sel.Name == name {
			return "signal." + name
		}
	}
	return ""
}

// registersCancelSignals reports whether a call spreads the owner's set and nothing else.
// A registration that spreads the set and then adds a signal of its own is still a second
// source, so the shape is exactly one spread argument.
func registersCancelSignals(call *ast.CallExpr, fset *token.FileSet) bool {
	return len(call.Args) == 2 && call.Ellipsis.IsValid() && expressionText(fset, call.Args[1]) == cancelSignalSet
}

// registeredSignalText renders the signals a call names, for the diagnostic.
func registeredSignalText(call *ast.CallExpr, fset *token.FileSet) string {
	if len(call.Args) == 1 {
		return "no signals"
	}
	var parts []string
	for _, arg := range call.Args[1:] {
		parts = append(parts, expressionText(fset, arg))
	}
	text := strings.Join(parts, ", ")
	if call.Ellipsis.IsValid() {
		text += "..."
	}
	return text
}

// TestCancelSignalRegistrationsBites is the recorded bite proof for
// checkCancelSignalRegistrations. It runs against a synthetic tree, not the repo, and
// walks the states that decide the check: a production registration of its own signal
// set, the conforming spread, the same token in a comment and a string, and a test file
// that registers anything it likes.
func TestCancelSignalRegistrationsBites(t *testing.T) {
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
	const production = "internal/example/example.go"

	write(production, "package example\n\nimport (\n\t\"os\"\n\t\"os/signal\"\n)\n\nfunc Trap() chan os.Signal {\n\tch := make(chan os.Signal, 1)\n\tsignal.Notify(ch, os.Interrupt)\n\treturn ch\n}\n")
	diags := checkCancelSignalRegistrations(root)
	if len(diags) != 1 || !strings.Contains(diags[0], production+":10 registers signal.Notify with os.Interrupt instead of subprocess.CancelSignals") {
		t.Fatalf("own signal set: want one diagnostic naming it, got %v", diags)
	}

	// The spread of the owner's set is the shape the check exists to protect.
	write(production, "package example\n\nimport (\n\t\"os\"\n\t\"os/signal\"\n\n\t\"github.com/gibbonmi/bench/internal/subprocess\"\n)\n\nfunc Trap() chan os.Signal {\n\tch := make(chan os.Signal, 1)\n\tsignal.Notify(ch, subprocess.CancelSignals...)\n\treturn ch\n}\n")
	if diags := checkCancelSignalRegistrations(root); len(diags) != 0 {
		t.Fatalf("spread of the owner set: want no diagnostics, got %v", diags)
	}

	// The context spelling registers the same set through a different call, so it carries
	// its own red.
	write(production, "package example\n\nimport (\n\t\"context\"\n\t\"os\"\n\t\"os/signal\"\n)\n\nfunc Trap() (context.Context, context.CancelFunc) {\n\treturn signal.NotifyContext(context.Background(), os.Interrupt)\n}\n")
	diags = checkCancelSignalRegistrations(root)
	if len(diags) != 1 || !strings.Contains(diags[0], "registers signal.NotifyContext with os.Interrupt instead of subprocess.CancelSignals") {
		t.Fatalf("context spelling: want one diagnostic, got %v", diags)
	}

	// A byte grep reds the comment and the string. The walk reads call sites, so neither is
	// a registration.
	write(production, "package example\n\n// Trap documents signal.Notify(ch, os.Interrupt) without calling it.\nfunc Trap() string {\n\treturn \"signal.Notify(ch, os.Interrupt)\"\n}\n")
	if diags := checkCancelSignalRegistrations(root); len(diags) != 0 {
		t.Fatalf("token in a comment and a string: want no diagnostics, got %v", diags)
	}

	// A test file registers whatever signal its own scenario sends, so the walk skips it.
	write(production, "package example\n")
	write("internal/example/example_test.go", "package example\n\nimport (\n\t\"os\"\n\t\"os/signal\"\n\t\"testing\"\n)\n\nfunc TestTrap(t *testing.T) {\n\tch := make(chan os.Signal, 1)\n\tsignal.Notify(ch, os.Interrupt)\n\tdefer signal.Stop(ch)\n}\n")
	if diags := checkCancelSignalRegistrations(root); len(diags) != 0 {
		t.Fatalf("test-file registration: want no diagnostics, got %v", diags)
	}
}
