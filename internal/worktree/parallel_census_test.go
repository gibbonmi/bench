package worktree

// This file owns the parallel census and the package's static pins. The census
// parses a directory's _test.go files with the Go AST and derives the serial
// set from two edges: a call to a serialHelpers harness helper, and an
// write to an imported package's variable through a selector. Both edges
// reach past the test process, so the test that holds one is serial. A
// top-level test that reaches no such edge is eligible and must call
// t.Parallel().
//
// A write to a package-level variable of a non-test file is a refusal, not
// a serial edge. Every seam this package injects travels in a per-call joins
// value, so a test that reaches for the package variable instead is reported.
// (Coverage rows WF02-WF05, WF12, WF14.)

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// serialHelpers name the harness helpers that bind the process environment or
// change the working directory. The census decides by bare identifier, not by
// type resolution, because a synthetic file set has no harness to resolve the
// name against. The census parses; it never builds.
var serialHelpers = map[string]bool{"bindEnv": true, "chdir": true}

// testFileFunc is one function declared in a test file: its declaration, the
// file that holds it, the line of its name, and the package names that file
// imports. The imports travel with the function because a selector on the left
// of a write names a package only in the file that imports it.
type testFileFunc struct {
	decl    *ast.FuncDecl
	file    string
	line    int
	imports map[string]bool
}

// parseTestFiles parses every regular _test.go file in dir. A special file, for
// example a FIFO, is not regular and is skipped, so the walk never blocks on a
// read.
func parseTestFiles(dir string) ([]*ast.File, *token.FileSet, []string, error) {
	return parseGoFiles(dir, func(name string) bool { return strings.HasSuffix(name, "_test.go") })
}

// parseSourceFiles parses every regular non-test .go file in dir. The census
// reads these files for their package-level variable declarations only.
func parseSourceFiles(dir string) ([]*ast.File, error) {
	files, _, _, err := parseGoFiles(dir, func(name string) bool {
		return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	})
	return files, err
}

// parseGoFiles parses every regular file in dir whose name want accepts.
func parseGoFiles(dir string, want func(string) bool) ([]*ast.File, *token.FileSet, []string, error) {
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

// isTopLevelTest reports whether decl is a TestXxx(t *testing.T) function.
// TestMain is neither eligible nor serial, so it is not a top-level test.
func isTopLevelTest(decl *ast.FuncDecl) bool {
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

// calleeName returns the identifier a call expression names, for a bare call
// such as helper(t). A method or package call, for example t.Run, names no
// test-file function and returns the empty string.
func calleeName(call *ast.CallExpr) string {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// serialHelperCalled names the serialHelpers helper the body of decl calls, a
// subtest closure included, and is empty when it calls none.
func serialHelperCalled(decl *ast.FuncDecl) string {
	found := ""
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok && serialHelpers[calleeName(call)] {
			found = calleeName(call)
		}
		return found == ""
	})
	return found
}

// packageVarNames collects every package-level identifier declared with var in
// the non-test files of dir, a grouped declaration included. A swap of one of
// these names is a process-wide effect, so the test that makes it is serial.
func packageVarNames(dir string) (map[string]bool, error) {
	files, err := parseSourceFiles(dir)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for _, file := range files {
		for _, node := range file.Decls {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.VAR {
				continue
			}
			for _, spec := range decl.Specs {
				value, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, ident := range value.Names {
					names[ident.Name] = true
				}
			}
		}
	}
	return names, nil
}

// assignRoot resolves a write target to the identifier it ultimately
// writes, and to the selector name that qualifies that identifier. The census
// strips the index, pointer, and parenthesis wrappers and the outer selectors,
// so pkg.Slice[0], pkg.Var.Field, and (*pkg.Ptr).Field all resolve to the root
// pkg and the selector the root qualifies. The selector name is empty when the
// target is a bare identifier or a wrapper around one, and the root name is
// empty when the target roots in no identifier at all. A write through any of
// these shapes reaches the same storage as a write to the bare name, so the
// serial edge and the package-variable refusal both read the root.
func assignRoot(target ast.Expr) (root string, sel string) {
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

func writeTargets(node ast.Node) []ast.Expr {
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

// assignedPackageVar names the first package-level variable the body of decl
// writes anywhere, a closure included, and is empty when it writes none. A
// restore inside t.Cleanup is such a closure. The census skips a := shadow
// because it declares a local instead of writing the package variable.
func assignedPackageVar(decl *ast.FuncDecl, vars map[string]bool) string {
	found := ""
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		for _, target := range writeTargets(node) {
			if root, _ := assignRoot(target); vars[root] {
				found = root
			}
		}
		return found == ""
	})
	return found
}

// fileImportNames returns the package names file imports. The name is the
// explicit alias when the import declares one, and the last element of the
// import path otherwise. A blank or a dot import names no package, so the
// census skips it.
func fileImportNames(file *ast.File) map[string]bool {
	names := map[string]bool{}
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := path[strings.LastIndex(path, "/")+1:]
		if spec.Name != nil {
			name = spec.Name.Name
		}
		if name == "" || name == "_" || name == "." {
			continue
		}
		names[name] = true
	}
	return names
}

// assignedImportedVar names the first imported package's variable the body of
// decl writes through a selector, a closure included, and is empty when it
// writes none. Such a write reaches past the test process, so it is a serial
// edge in the class of bindEnv, not a refusal. The census skips a := shadow of
// the package name because it declares a local instead of writing the import.
func assignedImportedVar(decl *ast.FuncDecl, imports map[string]bool) string {
	found := ""
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		for _, target := range writeTargets(node) {
			root, sel := assignRoot(target)
			if sel != "" && imports[root] {
				found = "assigns " + root + "." + sel
			}
		}
		return found == ""
	})
	return found
}

// directCallees returns the names of the test-file functions the body of decl
// calls, a subtest closure included.
func directCallees(decl *ast.FuncDecl) []string {
	var names []string
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if name := calleeName(call); name != "" {
				names = append(names, name)
			}
		}
		return true
	})
	return names
}

// callsParallelDirectly reports whether the body of decl calls t.Parallel() on
// its own testing parameter, outside every nested function literal. A subtest's
// t.Parallel() inside a closure is the subtest's call, not the parent's.
func callsParallelDirectly(decl *ast.FuncDecl) bool {
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

// reachedEffects walks the call edges from decl over the test-file functions and
// reports both effects any reached function holds: the serialHelpers helper a
// reached function calls, and the first package-level variable an assignment
// names. The walk is transitive, so a helper that reaches an effect through
// another helper counts. The serial reason names the helper alone when the test
// calls it, and names the helper the test reaches through a named helper
// otherwise.
func reachedEffects(entry testFileFunc, funcs map[string]testFileFunc, vars map[string]bool) (string, string) {
	seen := map[string]bool{}
	reason, assigned := "", ""
	var walk func(testFileFunc, string)
	walk = func(current testFileFunc, via string) {
		edge := serialHelperCalled(current.decl)
		if edge == "" {
			edge = assignedImportedVar(current.decl, current.imports)
		}
		if edge != "" && reason == "" {
			reason = edge
			if via != "" {
				reason = edge + " through " + via
			}
		}
		if assigned == "" {
			assigned = assignedPackageVar(current.decl, vars)
		}
		for _, name := range directCallees(current.decl) {
			if seen[name] {
				continue
			}
			seen[name] = true
			if next, ok := funcs[name]; ok {
				hop := via
				if hop == "" {
					hop = name
				}
				walk(next, hop)
			}
		}
	}
	walk(entry, "")
	return reason, assigned
}

// testFact is the census verdict for one top-level test: where it is declared,
// why it is serial, the package-level variable it assigns, and whether it calls
// t.Parallel() itself. An empty serialReason means the test is eligible.
type testFact struct {
	name         string
	file         string
	line         int
	serialReason string
	assigned     string
	parallel     bool
}

// censusFacts walks dir once and returns one fact for every top-level test. The
// parallel census and the serial set both read this walk, so the package
// classifies a test in one place. A test is serial when its body, or any
// test-file function it reaches through call edges, calls a serialHelpers
// helper.
func censusFacts(dir string) ([]testFact, error) {
	files, fset, _, err := parseTestFiles(dir)
	if err != nil {
		return nil, err
	}
	vars, err := packageVarNames(dir)
	if err != nil {
		return nil, err
	}
	// funcs indexes every function and method declared in the test files by
	// name. The call graph walks this index, so a function outside the test
	// files is a leaf.
	funcs := map[string]testFileFunc{}
	var tests []testFileFunc
	for _, file := range files {
		imports := fileImportNames(file)
		for _, node := range file.Decls {
			decl, ok := node.(*ast.FuncDecl)
			if !ok || decl.Body == nil || decl.Name == nil {
				continue
			}
			position := fset.Position(decl.Name.Pos())
			entry := testFileFunc{decl: decl, file: filepath.Base(position.Filename), line: position.Line, imports: imports}
			funcs[decl.Name.Name] = entry
			if isTopLevelTest(decl) {
				tests = append(tests, entry)
			}
		}
	}
	var facts []testFact
	for _, test := range tests {
		reason, assigned := reachedEffects(test, funcs, vars)
		facts = append(facts, testFact{
			name:         test.decl.Name.Name,
			file:         test.file,
			line:         test.line,
			serialReason: reason,
			assigned:     assigned,
			parallel:     callsParallelDirectly(test.decl),
		})
	}
	return facts, nil
}

// parallelCensus reports every top-level test in dir that breaks the parallel
// rule: an eligible test without t.Parallel(), or a serial test with it. A test
// that reaches an assignment to a package-level variable of a non-test file is
// reported for that assignment instead; the seam it wants belongs in the joins
// value it passes.
func parallelCensus(dir string) ([]string, error) {
	facts, err := censusFacts(dir)
	if err != nil {
		return nil, err
	}
	var reports []string
	for _, fact := range facts {
		serial := fact.serialReason != ""
		if fact.assigned != "" {
			reports = append(reports, fmt.Sprintf("%s:%d: %s assigns package variable %s", fact.file, fact.line, fact.name, fact.assigned))
			continue
		}
		switch {
		case serial && fact.parallel:
			reports = append(reports, fmt.Sprintf("%s:%d: %s is serial and calls t.Parallel()", fact.file, fact.line, fact.name))
		case !serial && !fact.parallel:
			reports = append(reports, fmt.Sprintf("%s:%d: %s is eligible and does not call t.Parallel()", fact.file, fact.line, fact.name))
		}
	}
	sort.Strings(reports)
	return reports, nil
}

// serialSet returns one line for every serial test in facts, sorted, with the
// reason the census classified it serial.
func serialSet(facts []testFact) []string {
	var serial []string
	for _, fact := range facts {
		if fact.serialReason == "" {
			continue
		}
		serial = append(serial, fmt.Sprintf("%s:%d: %s (%s)", fact.file, fact.line, fact.name, fact.serialReason))
	}
	sort.Strings(serial)
	return serial
}

// serialCeilingBreach returns the refusal for a serial set above ceiling, and is
// empty when the set fits. The refusal lists the whole set with each reason, so
// the reader sees which test is new.
func serialCeilingBreach(facts []testFact, ceiling int) string {
	serial := serialSet(facts)
	if len(serial) <= ceiling {
		return ""
	}
	return fmt.Sprintf("the package holds %d serial tests, above the ceiling of %d:\n%s",
		len(serial), ceiling, strings.Join(serial, "\n"))
}

// --- census unit tests over synthetic file sets ---

// plantTestFiles writes one synthetic file set into a temporary directory and
// returns the directory.
func plantTestFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatalf("plant %s: %v", name, err)
		}
	}
	return dir
}

func censusOf(t *testing.T, dir string) []string {
	t.Helper()
	reports, err := parallelCensus(dir)
	if err != nil {
		t.Fatalf("census: %v", err)
	}
	return reports
}

func containsReport(reports []string, want string) bool {
	for _, report := range reports {
		if strings.Contains(report, want) {
			return true
		}
	}
	return false
}

// TestCensusReportsEligibleTestWithoutParallel proves the census names the file
// and the line of an eligible test that omits t.Parallel(). (Coverage row WF02.)
func TestCensusReportsEligibleTestWithoutParallel(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"plain_test.go": `package worktree

import "testing"

func TestQuiet(t *testing.T) {
	_ = t
}
`})
	reports := censusOf(t, dir)
	want := "plain_test.go:5: TestQuiet is eligible and does not call t.Parallel()"
	if len(reports) != 1 || reports[0] != want {
		t.Fatalf("census = %q, want exactly [%q]", reports, want)
	}
}

// TestCensusFollowsHelperToBindEnv proves a test that calls a helper chain
// ending in bindEnv is serial, so the census does not report it for the missing
// call. The chain is two hops, so a depth-one walk misses it. (Coverage row WF03.)
func TestCensusFollowsHelperToBindEnv(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"helper_test.go": `package worktree

import "testing"

func TestBound(t *testing.T) {
	outerHelper(t)
}

func outerHelper(t *testing.T) {
	innerHelper(t)
}

func innerHelper(t *testing.T) {
	bindEnv(t, "BENCH_HOME", "value")
}
`})
	if reports := censusOf(t, dir); len(reports) != 0 {
		t.Fatalf("census = %q, want no report for a serial test without t.Parallel()", reports)
	}
}

// TestCensusFollowsSubtestClosureToChdir proves a chdir inside a subtest
// closure makes the parent serial. (Coverage row WF04.)
func TestCensusFollowsSubtestClosureToChdir(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"closure_test.go": `package worktree

import "testing"

func TestClosure(t *testing.T) {
	dir := t.TempDir()
	t.Run("x", func(t *testing.T) {
		chdir(t, dir)
	})
}
`})
	if reports := censusOf(t, dir); len(reports) != 0 {
		t.Fatalf("census = %q, want no report for a serial test without t.Parallel()", reports)
	}
}

// TestCensusReportsSerialTestWithParallel proves the census names a serial test
// that calls t.Parallel(), which Go would otherwise catch at run time.
// (Coverage row WF05.)
func TestCensusReportsSerialTestWithParallel(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"pair_test.go": `package worktree

import "testing"

func TestPair(t *testing.T) {
	t.Parallel()
	bindEnv(t, "BENCH_HOME", "value")
}
`})
	reports := censusOf(t, dir)
	want := "pair_test.go:5: TestPair is serial and calls t.Parallel()"
	if len(reports) != 1 || reports[0] != want {
		t.Fatalf("census = %q, want exactly [%q]", reports, want)
	}
}

// TestCensusRefusesAssignmentThroughAHelper proves the refusal follows the call
// edges: a test whose helper chain assigns the package-level variable is
// reported, and the report names the test, not the helper.
func TestCensusRefusesAssignmentThroughAHelper(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{
		"hooks.go": syntheticHooksFile,
		"chain_test.go": `package worktree

import "testing"

func TestChain(t *testing.T) {
	t.Parallel()
	swapHook(t)
}

func swapHook(t *testing.T) {
	hook = func() {}
}
`,
	})
	reports := censusOf(t, dir)
	want := "chain_test.go:5: TestChain assigns package variable hook"
	if len(reports) != 1 || reports[0] != want {
		t.Fatalf("census = %q, want exactly [%q]", reports, want)
	}
}

// TestCensusRefusesAssignmentInsideASubtestClosure proves a subtest closure is
// no shelter for the assignment: the parent carries the refusal.
func TestCensusRefusesAssignmentInsideASubtestClosure(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{
		"hooks.go": syntheticHooksFile,
		"subtest_test.go": `package worktree

import "testing"

func TestSubtest(t *testing.T) {
	t.Parallel()
	t.Run("x", func(t *testing.T) {
		hook = func() {}
	})
}
`,
	})
	reports := censusOf(t, dir)
	want := "subtest_test.go:5: TestSubtest assigns package variable hook"
	if len(reports) != 1 || reports[0] != want {
		t.Fatalf("census = %q, want exactly [%q]", reports, want)
	}
}

// syntheticHooksFile is the non-test file of a synthetic file set. It declares
// the package-level variable the assignment edge reads.
const syntheticHooksFile = `package worktree

var hook = func() {}

var hooks []func()

var hookBox struct{ Field int }

var hookPtr *struct{ Field int }
`

// TestCensusRefusesAssignmentToPackageVariable proves a test that assigns a
// package-level variable of a non-test file is refused, in the body and inside a
// t.Cleanup closure alike, and that a read or a := shadow of the same name
// leaves the test eligible.
func TestCensusRefusesAssignmentToPackageVariable(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name string
		body string
		want string
	}{
		{
			name: "body-assignment",
			body: `	hook = func() {}
`,
			want: "stub_test.go:5: TestStub assigns package variable hook",
		},
		{
			name: "increment",
			body: `	hookBox.Field++
`,
			want: "stub_test.go:5: TestStub assigns package variable hookBox",
		},
		{
			name: "decrement",
			body: `	hookBox.Field--
`,
			want: "stub_test.go:5: TestStub assigns package variable hookBox",
		},
		{
			name: "cleanup-closure-assignment",
			body: `	old := hook
	t.Cleanup(func() { hook = old })
`,
			want: "stub_test.go:5: TestStub assigns package variable hook",
		},
		{
			name: "index-assignment",
			body: `	hooks[0] = nil
`,
			want: "stub_test.go:5: TestStub assigns package variable hooks",
		},
		{
			name: "field-assignment",
			body: `	hookBox.Field = 1
`,
			want: "stub_test.go:5: TestStub assigns package variable hookBox",
		},
		{
			name: "pointer-field-assignment",
			body: `	(*hookPtr).Field = 1
`,
			want: "stub_test.go:5: TestStub assigns package variable hookPtr",
		},
		{
			name: "read-only",
			body: `	_ = hook
`,
			want: "stub_test.go:5: TestStub is eligible and does not call t.Parallel()",
		},
		{
			name: "shadow-declaration",
			body: `	hook := func() {}
	_ = hook
`,
			want: "stub_test.go:5: TestStub is eligible and does not call t.Parallel()",
		},
	} {
		t.Run(row.name, func(t *testing.T) {
			dir := plantTestFiles(t, map[string]string{
				"hooks.go": syntheticHooksFile,
				"stub_test.go": `package worktree

import "testing"

func TestStub(t *testing.T) {
` + row.body + `}
`,
			})
			reports := censusOf(t, dir)
			if len(reports) != 1 || reports[0] != row.want {
				t.Fatalf("census = %q, want exactly [%q]", reports, row.want)
			}
		})
	}
}

// TestCensusReportsPackageVariableAssignmentWithParallel proves the assignment
// refusal, not the serial verdict, is what a test that also calls t.Parallel()
// carries.
func TestCensusReportsPackageVariableAssignmentWithParallel(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{
		"hooks.go": syntheticHooksFile,
		"stub_pair_test.go": `package worktree

import "testing"

func TestStubPair(t *testing.T) {
	t.Parallel()
	hook = func() {}
}
`,
	})
	reports := censusOf(t, dir)
	want := "stub_pair_test.go:5: TestStubPair assigns package variable hook"
	if len(reports) != 1 || reports[0] != want {
		t.Fatalf("census = %q, want exactly [%q]", reports, want)
	}
}

// TestCensusTreatsImportedPackageVariableAssignmentAsSerial proves an assignment
// to an imported package's variable through a selector is a serial edge in the
// class of bindEnv: the test is serial and silent, and a := shadow or a read of
// the same selector leaves the test eligible.
func TestCensusTreatsImportedPackageVariableAssignmentAsSerial(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name       string
		importSpec string
		body       string
		want       string
	}{
		{
			name: "body-assignment",
			body: `	rand.Reader = nil
`,
			want: "",
		},
		{
			name: "increment",
			body: `	rand.Reader++
`,
			want: "",
		},
		{
			name: "decrement",
			body: `	rand.Reader--
`,
			want: "",
		},
		{
			name: "assignment-through-a-helper",
			body: `	swapReader()
}

func swapReader() {
	rand.Reader = nil
`,
			want: "",
		},
		{
			name:       "aliased-import",
			importSpec: `cryptorand "crypto/rand"`,
			body: `	cryptorand.Reader = nil
`,
			want: "",
		},
		{
			name: "index-assignment",
			body: `	rand.Reader[0] = nil
`,
			want: "",
		},
		{
			name: "field-assignment",
			body: `	rand.Reader.Field = nil
`,
			want: "",
		},
		{
			name: "pointer-field-assignment",
			body: `	(*rand.Reader).Field = nil
`,
			want: "",
		},
		{
			name: "read-only",
			body: `	_ = rand.Reader
`,
			want: "reader_test.go:8: TestReader is eligible and does not call t.Parallel()",
		},
		{
			name: "shadow-declaration",
			body: `	rand := struct{ Reader int }{}
	_ = rand.Reader
`,
			want: "reader_test.go:8: TestReader is eligible and does not call t.Parallel()",
		},
	} {
		t.Run(row.name, func(t *testing.T) {
			dir := plantTestFiles(t, map[string]string{"reader_test.go": syntheticReaderFile(row.importSpec, row.body)})
			reports := censusOf(t, dir)
			if row.want == "" {
				if len(reports) != 0 {
					t.Fatalf("census = %q, want no report for a serial test without t.Parallel()", reports)
				}
				return
			}
			if len(reports) != 1 || reports[0] != row.want {
				t.Fatalf("census = %q, want exactly [%q]", reports, row.want)
			}
		})
	}
}

// TestCensusReportsImportedPackageVariableAssignmentWithParallel proves the pair
// report names the selector the test assigns, so a serial test that also calls
// t.Parallel() is caught before Go runs it.
func TestCensusReportsImportedPackageVariableAssignmentWithParallel(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"reader_test.go": syntheticReaderFile("", `	t.Parallel()
	rand.Reader = nil
`)})
	reports := censusOf(t, dir)
	want := "reader_test.go:8: TestReader is serial and calls t.Parallel()"
	if len(reports) != 1 || reports[0] != want {
		t.Fatalf("census = %q, want exactly [%q]", reports, want)
	}
}

// syntheticReaderFile is a synthetic test file that imports crypto/rand with
// importSpec, an empty importSpec being the plain path, and holds body as the
// whole body of its one top-level test.
func syntheticReaderFile(importSpec, body string) string {
	if importSpec == "" {
		importSpec = `"crypto/rand"`
	}
	return `package worktree

import (
	` + importSpec + `
	"testing"
)

func TestReader(t *testing.T) {
` + body + `}
`
}

// TestCensusWalksBuildTaggedFileAndSkipsTestMain proves a build-tagged file is
// walked like any other and that TestMain is neither eligible nor serial.
func TestCensusWalksBuildTaggedFileAndSkipsTestMain(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{
		"tagged_test.go": `//go:build linux

package worktree

import "testing"

func TestTagged(t *testing.T) {
	_ = t
}
`,
		"main_test.go": `package worktree

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
`,
	})
	reports := censusOf(t, dir)
	if len(reports) != 1 || !containsReport(reports, "TestTagged is eligible") {
		t.Fatalf("census = %q, want the build-tagged test alone", reports)
	}
}

// TestCensusAcceptsTableParentWithParallelSubtests proves a table parent that
// binds no environment and calls t.Parallel() itself is eligible and silent,
// even when its subtests call t.Parallel() in their closures.
func TestCensusAcceptsTableParentWithParallelSubtests(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"table_test.go": `package worktree

import "testing"

func TestTable(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"a", "b"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_ = name
		})
	}
}
`})
	if reports := censusOf(t, dir); len(reports) != 0 {
		t.Fatalf("census = %q, want no report", reports)
	}
}

// TestCensusReportsTableParentWithoutOwnParallel proves a subtest's t.Parallel()
// inside a closure is not the parent's call: an eligible parent must call
// t.Parallel() too. (Edge inventory: "the parent calls t.Parallel() too".)
func TestCensusReportsTableParentWithoutOwnParallel(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"closure_only_test.go": `package worktree

import "testing"

func TestClosureOnly(t *testing.T) {
	for _, name := range []string{"a", "b"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_ = name
		})
	}
}
`})
	reports := censusOf(t, dir)
	want := "closure_only_test.go:5: TestClosureOnly is eligible and does not call t.Parallel()"
	if len(reports) != 1 || reports[0] != want {
		t.Fatalf("census = %q, want exactly [%q]", reports, want)
	}
}

// TestCensusSkipsSpecialFile proves the walk skips a FIFO named like a test
// file, so a read never blocks.
func TestCensusSkipsSpecialFile(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"plain_test.go": `package worktree

import "testing"

func TestQuiet(t *testing.T) {
	t.Parallel()
}
`})
	if err := syscall.Mkfifo(filepath.Join(dir, "fifo_test.go"), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	if reports := censusOf(t, dir); len(reports) != 0 {
		t.Fatalf("census = %q, want no report", reports)
	}
}

// TestCensusRefusesASerialSetAboveTheCeiling proves the ceiling check counts the
// serial set and names every member with its reason, so one test above the
// ceiling turns the check red. (Coverage row WF18.)
func TestCensusRefusesASerialSetAboveTheCeiling(t *testing.T) {
	t.Parallel()
	dir := plantTestFiles(t, map[string]string{"serial_test.go": `package worktree

import "testing"

func TestBoundOne(t *testing.T) {
	bindEnv(t, "BENCH_HOME", "value")
}

func TestBoundTwo(t *testing.T) {
	moveAway(t)
}

func moveAway(t *testing.T) {
	chdir(t, t.TempDir())
}
`})
	facts, err := censusFacts(dir)
	if err != nil {
		t.Fatalf("census: %v", err)
	}
	if breach := serialCeilingBreach(facts, 2); breach != "" {
		t.Fatalf("a set of two under a ceiling of two = %q, want no refusal", breach)
	}
	breach := serialCeilingBreach(facts, 1)
	for _, want := range []string{
		"the package holds 2 serial tests, above the ceiling of 1",
		"serial_test.go:5: TestBoundOne (bindEnv)",
		"serial_test.go:9: TestBoundTwo (chdir through moveAway)",
	} {
		if !strings.Contains(breach, want) {
			t.Errorf("refusal = %q, want it to contain %q", breach, want)
		}
	}
}

// --- live-tree pins ---

// TestParallelCensusOnTheLiveTree proves every eligible test in the package
// calls t.Parallel() and no serial test calls it. (Coverage row WF01.)
func TestParallelCensusOnTheLiveTree(t *testing.T) {
	t.Parallel()
	reports, err := parallelCensus(".")
	if err != nil {
		t.Fatalf("census the package: %v", err)
	}
	if len(reports) != 0 {
		t.Fatalf("the parallel census reports %d breaks:\n%s", len(reports), strings.Join(reports, "\n"))
	}
}

// worktreeSerialCeiling is the package's pinned serial-test count. Four classes
// hold the whole set. The boundary-default graders bind the process environment
// because the default Home() resolves is what they grade. The exec child tests
// bind the caller environment because the child's inherited environment is what
// they grade. The stub tests bind PATH because the child under test reads it.
// The operand tests change the working directory because a relative or a
// prefixed operand is what they grade. A fixture that falls back to a process
// bind instead of the home it owns raises the count above this ceiling.
const worktreeSerialCeiling = 45

// TestSerialSetStaysBelowTheCeiling proves no new test joins the serial set.
// (Coverage row WF18.)
func TestSerialSetStaysBelowTheCeiling(t *testing.T) {
	t.Parallel()
	facts, err := censusFacts(".")
	if err != nil {
		t.Fatalf("census the package: %v", err)
	}
	if breach := serialCeilingBreach(facts, worktreeSerialCeiling); breach != "" {
		t.Fatal(breach)
	}
	t.Logf("the package holds %d serial tests, at or below the ceiling of %d", len(serialSet(facts)), worktreeSerialCeiling)
}

// worktreeTestFloor is the package's pinned top-level test count. A new test
// raises the count; a removal below the pin turns the gate red.
const worktreeTestFloor = 334

// TestPackageTestCountPin proves no test is removed or merged for wall-clock.
// (Coverage row WF12.)
func TestPackageTestCountPin(t *testing.T) {
	t.Parallel()
	files, _, _, err := parseTestFiles(".")
	if err != nil {
		t.Fatalf("parse the package test files: %v", err)
	}
	count := 0
	for _, file := range files {
		for _, node := range file.Decls {
			if decl, ok := node.(*ast.FuncDecl); ok && isTopLevelTest(decl) {
				count++
			}
		}
	}
	if count < worktreeTestFloor {
		t.Fatalf("the package declares %d top-level tests, below the pin of %d", count, worktreeTestFloor)
	}
	t.Logf("the package declares %d top-level tests", count)
}

// TestPackageClausePin proves every test file stays in package worktree, so no
// sub-package split arrives by another name. (Coverage row WF14.)
func TestPackageClausePin(t *testing.T) {
	t.Parallel()
	files, _, names, err := parseTestFiles(".")
	if err != nil {
		t.Fatalf("parse the package test files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("the walk found no test file")
	}
	for i, file := range files {
		if file.Name.Name != "worktree" {
			t.Errorf("%s declares package %s, want package worktree", names[i], file.Name.Name)
		}
	}
}
