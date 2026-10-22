//go:build system

package systemtest

// This file owns the system parallel census. The census parses the package's test files
// through the shared testcensus walk. A top-level test is parallel when it calls
// t.Parallel(), and serial when systemSerialTests holds a reason for it. Parallel tests
// share the one owner, so a parallel test must not reach the owner's repositories or name
// a ledger field outside an owner writer. An owner writer writes a ledger field, and the
// owner mutex holds each of its ledger-field references.

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/testcensus"
)

// ledgerFields name the owner fields that record the suite's observations. TestMain
// verifies them after an unfiltered run, while parallel tests write them.
var ledgerFields = map[string]bool{"starts": true, "seen": true, "terminal": true}

// lockExempt names the one function whose ledger writes skip the owner mutex. It plants
// the verification's negative controls in a child process that runs no other test.
const lockExempt = "TestOwnerSelectionChild"

// systemFunc is one function declared in a test file, with the base name of that file.
type systemFunc struct {
	decl *ast.FuncDecl
	file string
}

// systemPackage is one parse of a directory's test files. funcs indexes a function by
// its name and a method by a dot and its name, because a bare identifier resolves to a
// function and a selector resolves to a method.
type systemPackage struct {
	fset  *token.FileSet
	all   []systemFunc
	funcs map[string][]systemFunc
}

// parseSystemPackage parses dir's test files once for one census.
func parseSystemPackage(dir string) (systemPackage, error) {
	files, fset, names, err := testcensus.ParseTestFiles(dir)
	if err != nil {
		return systemPackage{}, err
	}
	pkg := systemPackage{fset: fset, funcs: map[string][]systemFunc{}}
	for i, file := range files {
		for _, node := range file.Decls {
			decl, ok := node.(*ast.FuncDecl)
			if !ok || decl.Body == nil {
				continue
			}
			fn := systemFunc{decl: decl, file: names[i]}
			key := decl.Name.Name
			if decl.Recv != nil {
				key = "." + key
			}
			pkg.funcs[key] = append(pkg.funcs[key], fn)
			pkg.all = append(pkg.all, fn)
		}
	}
	return pkg, nil
}

func (p systemPackage) where(fn systemFunc, pos token.Pos) string {
	return fmt.Sprintf("%s:%d", fn.file, p.fset.Position(pos).Line)
}

// ownerValue reports whether name spells the owner value: owner for the package's owner,
// and o for a systemOwner receiver.
func ownerValue(name string) bool { return name == "owner" || name == "o" }

// ownerSelector splits a selector on the owner value into that value and the name it
// selects. Any other node selects nothing.
func ownerSelector(node ast.Node) (value, name string) {
	selector, ok := node.(*ast.SelectorExpr)
	if !ok {
		return "", ""
	}
	ident, ok := selector.X.(*ast.Ident)
	if !ok || !ownerValue(ident.Name) {
		return "", ""
	}
	return ident.Name, selector.Sel.Name
}

// ownerMutexCall reports whether call locks or unlocks the owner mutex, and which.
func ownerMutexCall(call *ast.CallExpr) (lock, ok bool) {
	method, isSelector := call.Fun.(*ast.SelectorExpr)
	if !isSelector || (method.Sel.Name != "Lock" && method.Sel.Name != "Unlock") {
		return false, false
	}
	if value, name := ownerSelector(method.X); value == "" || name != "mu" {
		return false, false
	}
	return method.Sel.Name == "Lock", true
}

// ledgerRef is one ledger-field reference on the owner value, such as o.starts.
type ledgerRef struct {
	text   string
	pos    token.Pos
	write  bool
	locked bool
}

// ledgerRefs returns every ledger-field reference in decl in source order. The owner
// mutex holds a reference when the last owner-mutex call before it is a Lock. A deferred
// Unlock runs at return, so it releases nothing inside the body.
func ledgerRefs(decl *ast.FuncDecl) []ledgerRef {
	var refs []ledgerRef
	held := false
	deferred := map[*ast.CallExpr]bool{}
	ast.Inspect(decl.Body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.DeferStmt:
			deferred[n.Call] = true
		case *ast.CallExpr:
			if lock, ok := ownerMutexCall(n); ok && !deferred[n] {
				held = lock
			}
		}
		for _, target := range testcensus.WriteTargets(node) {
			if root, field := testcensus.AssignRoot(target); ownerValue(root) && ledgerFields[field] {
				refs = append(refs, ledgerRef{text: root + "." + field, pos: target.Pos(), write: true, locked: held})
			}
		}
		if value, name := ownerSelector(node); ledgerFields[name] {
			refs = append(refs, ledgerRef{text: value + "." + name, pos: node.Pos(), locked: held})
		}
		return true
	})
	return refs
}

// ownerWriter reports whether decl writes a ledger field while the owner mutex holds
// each of its ledger-field references.
func ownerWriter(decl *ast.FuncDecl) bool {
	refs := ledgerRefs(decl)
	wrote := false
	for _, ref := range refs {
		if !ref.locked {
			return false
		}
		wrote = wrote || ref.write
	}
	return wrote
}

// callees returns each test-file function that body can run, as a call or as a value. A
// bare identifier resolves to the function of its name, and a selector to each method of
// its name on any receiver, because the census reads no types. A selector's own name and
// a composite-literal key name no function.
func (p systemPackage) callees(body *ast.BlockStmt) []systemFunc {
	var found []systemFunc
	skip := map[ast.Expr]bool{}
	ast.Inspect(body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.SelectorExpr:
			skip[n.Sel] = true
			found = append(found, p.funcs["."+n.Sel.Name]...)
		case *ast.KeyValueExpr:
			skip[n.Key] = true
		case *ast.Ident:
			if !skip[n] {
				found = append(found, p.funcs[n.Name]...)
			}
		}
		return true
	})
	return found
}

// reach names the first owner-repository reference and the first ledger-field reference
// that test reaches through its body and the functions it can run, each with the helper
// it reaches the reference through. A ledger field in the test body counts even under the
// owner mutex. Inside an owner writer the walk drops only the ledger rule. The visited set
// keys a function by the ledger rule, so a function that a writer runs still grades its
// ledger references when the test also reaches it outside every writer.
func (p systemPackage) reach(test systemFunc) (repos, ledger string) {
	type visit struct {
		decl   *ast.FuncDecl
		graded bool
	}
	visited := map[visit]bool{{test.decl, true}: true}
	var walk func(fn systemFunc, through string, graded bool)
	walk = func(fn systemFunc, through string, graded bool) {
		ast.Inspect(fn.decl.Body, func(node ast.Node) bool {
			if value, name := ownerSelector(node); name == "repos" && repos == "" {
				repos = value + ".repos" + through
			} else if graded && ledgerFields[name] && ledger == "" {
				ledger = value + "." + name + through
			}
			return true
		})
		for _, next := range p.callees(fn.decl.Body) {
			step := visit{next.decl, graded && !ownerWriter(next.decl)}
			if visited[step] {
				continue
			}
			visited[step] = true
			hop := through
			if hop == "" {
				hop = " through " + next.decl.Name.Name
			}
			walk(next, hop, step.graded)
		}
	}
	walk(test, "", true)
	return repos, ledger
}

// systemParallelCensus reports every top-level test in dir that calls no t.Parallel()
// and holds no serial reason, and every parallel test that reaches the owner's
// repositories or a ledger field.
func systemParallelCensus(dir string, serial map[string]string) ([]string, error) {
	pkg, err := parseSystemPackage(dir)
	if err != nil {
		return nil, err
	}
	var reports []string
	for _, fn := range pkg.all {
		if !testcensus.IsTopLevelTest(fn.decl) {
			continue
		}
		name, where := fn.decl.Name.Name, pkg.where(fn, fn.decl.Name.Pos())
		if !testcensus.CallsParallelDirectly(fn.decl) {
			if serial[name] == "" {
				reports = append(reports, fmt.Sprintf("%s: %s calls no t.Parallel() and holds no serial reason", where, name))
			}
			continue
		}
		repos, ledger := pkg.reach(fn)
		if repos != "" {
			reports = append(reports, fmt.Sprintf("%s: %s calls t.Parallel() and reaches %s", where, name, repos))
		}
		if ledger != "" {
			reports = append(reports, fmt.Sprintf("%s: %s calls t.Parallel() and names %s", where, name, ledger))
		}
	}
	sort.Strings(reports)
	return reports, nil
}

// ledgerLockCensus reports every ledger-field write in dir outside the owner mutex,
// except in lockExempt.
func ledgerLockCensus(dir string) ([]string, error) {
	pkg, err := parseSystemPackage(dir)
	if err != nil {
		return nil, err
	}
	var reports []string
	for _, fn := range pkg.all {
		if fn.decl.Recv == nil && fn.decl.Name.Name == lockExempt {
			continue
		}
		for _, ref := range ledgerRefs(fn.decl) {
			if ref.write && !ref.locked {
				reports = append(reports, fmt.Sprintf("%s: %s writes %s outside the owner mutex", pkg.where(fn, ref.pos), fn.decl.Name.Name, ref.text))
			}
		}
	}
	sort.Strings(reports)
	return reports, nil
}

// serialTableCensus reports every serial row that names no top-level test in dir, names
// a test that calls t.Parallel(), or holds no reason.
func serialTableCensus(dir string, serial map[string]string) ([]string, error) {
	pkg, err := parseSystemPackage(dir)
	if err != nil {
		return nil, err
	}
	parallel := map[string]bool{}
	for _, fn := range pkg.all {
		if testcensus.IsTopLevelTest(fn.decl) {
			parallel[fn.decl.Name.Name] = testcensus.CallsParallelDirectly(fn.decl)
		}
	}
	var reports []string
	for name, reason := range serial {
		if isParallel, declared := parallel[name]; !declared {
			reports = append(reports, fmt.Sprintf("serial row %s names no top-level test", name))
		} else if isParallel {
			reports = append(reports, fmt.Sprintf("serial row %s names a test that calls t.Parallel()", name))
		}
		if reason == "" {
			reports = append(reports, fmt.Sprintf("serial row %s holds no reason", name))
		}
	}
	sort.Strings(reports)
	return reports, nil
}

// requireCensus runs census over one planted test file, whose body starts on line 5,
// and requires exactly want.
func requireCensus(t *testing.T, census func(string) ([]string, error), body string, want []string) {
	t.Helper()
	dir := testcensus.Plant(t, map[string]string{"planted_test.go": "package systemtest\n\nimport \"testing\"\n\n" + body})
	requireReports(t, census, dir, want)
}

// requireReports runs census over dir and requires exactly want.
func requireReports(t *testing.T, census func(string) ([]string, error), dir string, want []string) {
	t.Helper()
	reports, err := census(dir)
	if err != nil {
		t.Fatalf("census %s: %v", dir, err)
	}
	if !slices.Equal(reports, want) {
		t.Fatalf("census %s =\n%s\nwant\n%s", dir, strings.Join(reports, "\n"), strings.Join(want, "\n"))
	}
}

// TestSystemParallelCensus proves the census refuses a test with no t.Parallel() call and
// no serial reason, and a parallel test that reaches the owner's repositories or names a
// ledger field outside an owner writer, and that the live package passes. (Coverage rows
// LGP25, LGP29.)
func TestSystemParallelCensus(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, body string
		serial     map[string]string
		want       []string
	}{
		{"no call and no reason", "func TestQuiet(t *testing.T) {}\n", nil,
			[]string{"planted_test.go:5: TestQuiet calls no t.Parallel() and holds no serial reason"}},
		{"serial reason", "func TestQuiet(t *testing.T) {}\n", map[string]string{"TestQuiet": pendingMigration}, nil},
		{"owner repository", "func TestShare(t *testing.T) { t.Parallel(); _ = owner.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestShare calls t.Parallel() and reaches owner.repos"}},
		{"owner repository through a method", "func TestShare(t *testing.T) { t.Parallel(); first() }\n\nfunc first() string { return owner.pick() }\n\nfunc (o *systemOwner) pick() string { return o.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestShare calls t.Parallel() and reaches o.repos through first"}},
		{"ledger field", "func TestBump(t *testing.T) { t.Parallel(); owner.starts++ }\n", nil,
			[]string{"planted_test.go:5: TestBump calls t.Parallel() and names owner.starts"}},
		{"ledger field through a helper", "func TestPeek(t *testing.T) { t.Parallel(); peek() }\n\nfunc peek() int { return len(owner.seen) }\n", nil,
			[]string{"planted_test.go:5: TestPeek calls t.Parallel() and names owner.seen through peek"}},
		{"owner writer", "func TestRun(t *testing.T) { t.Parallel(); owner.bump() }\n\nfunc (o *systemOwner) bump() { o.mu.Lock(); o.starts++; o.mu.Unlock() }\n", nil, nil},
		{"owner repository through an owner writer", "func TestStore(t *testing.T) { t.Parallel(); owner.store() }\n\nfunc (o *systemOwner) store() { o.mu.Lock(); o.starts++; _ = o.repos[0]; o.mu.Unlock() }\n", nil,
			[]string{"planted_test.go:5: TestStore calls t.Parallel() and reaches o.repos through store"}},
		{"owner repository through a method value", "func TestPick(t *testing.T) { t.Parallel(); f := owner.pick; f() }\n\nfunc (o *systemOwner) pick() string { return o.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestPick calls t.Parallel() and reaches o.repos through pick"}},
		{"owner repository through a function value", "func TestFirst(t *testing.T) { t.Parallel(); f := first; f() }\n\nfunc first() string { return owner.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestFirst calls t.Parallel() and reaches owner.repos through first"}},
		{"owner repository through a non-owner method", "func TestHold(t *testing.T) { t.Parallel(); var h holder; h.first() }\n\ntype holder struct{}\n\nfunc (holder) first() string { return owner.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestHold calls t.Parallel() and reaches owner.repos through first"}},
		{"local names", "func TestLocal(t *testing.T) { t.Parallel(); seen := map[string]bool{\"starts\": true}; seen[\"terminal\"] = true }\n", nil, nil},
	} {
		t.Run(row.name, func(t *testing.T) {
			requireCensus(t, func(dir string) ([]string, error) { return systemParallelCensus(dir, row.serial) }, row.body, row.want)
		})
	}
	requireReports(t, func(dir string) ([]string, error) { return systemParallelCensus(dir, systemSerialTests) }, ".", nil)
}

// TestLedgerWritesHoldTheOwnerLock proves every ledger write outside lockExempt holds the
// owner mutex: an unguarded write reds by name, and the live package passes. (Coverage
// row LGP28.)
func TestLedgerWritesHoldTheOwnerLock(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, body string
		want       []string
	}{
		{"unguarded", "func (o *systemOwner) bump() { o.starts++ }\n",
			[]string{"planted_test.go:5: bump writes o.starts outside the owner mutex"}},
		{"after the unlock", "func record(x int) {\n\towner.mu.Lock()\n\towner.mu.Unlock()\n\towner.seen = append(owner.seen, x)\n}\n",
			[]string{"planted_test.go:8: record writes owner.seen outside the owner mutex"}},
		{"guarded", "func (o *systemOwner) mark(k string) {\n\to.mu.Lock()\n\tdefer o.mu.Unlock()\n\to.terminal[k] = true\n}\n", nil},
		{"exempt", "func TestOwnerSelectionChild(t *testing.T) { owner.starts = 0 }\n", nil},
	} {
		t.Run(row.name, func(t *testing.T) {
			requireCensus(t, ledgerLockCensus, row.body, row.want)
		})
	}
	requireReports(t, ledgerLockCensus, ".", nil)
}

// TestSystemSerialTableIsExact proves a serial row that names a missing test or a
// parallel test, or that holds no reason, reds, and that the live table is exact.
// (Coverage row LGP31.)
func TestSystemSerialTableIsExact(t *testing.T) {
	t.Parallel()
	table := map[string]string{"TestSerial": "", "TestGone": pendingMigration, "TestMain": pendingMigration, "TestQuick": pendingMigration}
	requireCensus(t, func(dir string) ([]string, error) { return serialTableCensus(dir, table) },
		"func TestMain(m *testing.M) {}\n\nfunc TestSerial(t *testing.T) {}\n\nfunc TestQuick(t *testing.T) { t.Parallel() }\n",
		[]string{
			"serial row TestGone names no top-level test",
			"serial row TestMain names no top-level test",
			"serial row TestQuick names a test that calls t.Parallel()",
			"serial row TestSerial holds no reason",
		})
	requireReports(t, func(dir string) ([]string, error) { return serialTableCensus(dir, systemSerialTests) }, ".", nil)
}
