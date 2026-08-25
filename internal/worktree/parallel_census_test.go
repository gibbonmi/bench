package worktree

// This file owns the parallel census and the package's static pins. The census
// parses a directory's _test.go files with the Go AST and derives the serial
// set from call edges to the harness helpers bindEnv and chdir. A top-level
// test that reaches neither helper is eligible and must call t.Parallel().
// (Coverage rows WF02-WF05, WF12, WF14.)

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
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
// file that holds it, and the line of its name.
type testFileFunc struct {
	decl *ast.FuncDecl
	file string
	line int
}

// parseTestFiles parses every regular _test.go file in dir. A special file, for
// example a FIFO, is not regular and is skipped, so the walk never blocks on a
// read.
func parseTestFiles(dir string) ([]*ast.File, *token.FileSet, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, nil, err
	}
	fset := token.NewFileSet()
	var files []*ast.File
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, "_test.go") {
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

// callsSerialHelper reports whether the body of decl calls bindEnv or chdir
// anywhere, a subtest closure included.
func callsSerialHelper(decl *ast.FuncDecl) bool {
	found := false
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok && serialHelpers[calleeName(call)] {
			found = true
		}
		return !found
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

// reachesSerialHelper walks the call edges from decl over the test-file
// functions and reports whether any reached function calls bindEnv or chdir.
// The walk is transitive, so a helper that reaches the helper through another
// helper is still a serial edge.
func reachesSerialHelper(decl *ast.FuncDecl, funcs map[string]testFileFunc) bool {
	seen := map[string]bool{}
	var walk func(*ast.FuncDecl) bool
	walk = func(current *ast.FuncDecl) bool {
		if callsSerialHelper(current) {
			return true
		}
		for _, name := range directCallees(current) {
			if seen[name] {
				continue
			}
			seen[name] = true
			if next, ok := funcs[name]; ok && walk(next.decl) {
				return true
			}
		}
		return false
	}
	return walk(decl)
}

// parallelCensus reports every top-level test in dir that breaks the parallel
// rule: an eligible test without t.Parallel(), or a serial test with it. A test
// is serial when its body, or any test-file function it reaches through call
// edges, calls bindEnv or chdir.
func parallelCensus(dir string) ([]string, error) {
	files, fset, _, err := parseTestFiles(dir)
	if err != nil {
		return nil, err
	}
	// funcs indexes every function and method declared in the test files by
	// name. The call graph walks this index, so a function outside the test
	// files is a leaf.
	funcs := map[string]testFileFunc{}
	var tests []testFileFunc
	for _, file := range files {
		for _, node := range file.Decls {
			decl, ok := node.(*ast.FuncDecl)
			if !ok || decl.Body == nil || decl.Name == nil {
				continue
			}
			position := fset.Position(decl.Name.Pos())
			entry := testFileFunc{decl: decl, file: filepath.Base(position.Filename), line: position.Line}
			funcs[decl.Name.Name] = entry
			if isTopLevelTest(decl) {
				tests = append(tests, entry)
			}
		}
	}
	var reports []string
	for _, test := range tests {
		serial := reachesSerialHelper(test.decl, funcs)
		parallel := callsParallelDirectly(test.decl)
		switch {
		case serial && parallel:
			reports = append(reports, fmt.Sprintf("%s:%d: %s is serial and calls t.Parallel()", test.file, test.line, test.decl.Name.Name))
		case !serial && !parallel:
			reports = append(reports, fmt.Sprintf("%s:%d: %s is eligible and does not call t.Parallel()", test.file, test.line, test.decl.Name.Name))
		}
	}
	sort.Strings(reports)
	return reports, nil
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
