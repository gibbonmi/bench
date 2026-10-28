package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// routeBypassAllowed is the reviewed allowlist: each file of a scanned package whose routes
// the registry does not own, keyed to the read or shared caller that its routes serve.
var routeBypassAllowed = map[string]string{
	"internal/worktree/path.go":        "the shared --in tree-target refusal of every tree-scoped verb, read and write",
	"internal/worktree/build.go":       "the success routes of bench worktree build and its install form, and the install refusal routes",
	"internal/worktree/tree_target.go": "the success route of bench worktree create",
}

// routeBypassNeedles are the route spellings that only the registry writes: the next label
// as a record field and as a table header, the commitment route tail, and the route joiner.
// The spec pins the tail, so the check spells it; the registry exports the rest.
var routeBypassNeedles = []string{refusalroute.NextField + "=", refusalroute.NextField + "[", "run bench ", refusalroute.StepJoiner}

// routeBypassFaults names each string literal in a production Go file of a scanned package
// that composes a route outside the registry: one that holds a needle, or that is the bare
// next label. A comment does not print, so the scan reads literals and not text. A scanned
// package with no production file is a fault too, so a moved package cannot leave the scan.
func routeBypassFaults(root string) []string {
	var faults []string
	for _, pkg := range registry.RouteBypassPackages {
		entries, _ := os.ReadDir(filepath.Join(root, filepath.FromSlash(pkg)))
		scanned := 0
		for _, entry := range entries {
			rel := pkg + "/" + entry.Name()
			if entry.IsDir() || !registry.RouteBypassScanned(rel) {
				continue
			}
			scanned++
			if _, allowed := routeBypassAllowed[rel]; !allowed {
				faults = append(faults, routeLiteralFaults(filepath.Join(root, filepath.FromSlash(rel)), rel)...)
			}
		}
		if scanned == 0 {
			faults = append(faults, pkg+" holds no production Go file, so the route bypass check scans nothing there")
		}
	}
	return faults
}

// routeLiteralFaults grades the string literals of one production file.
func routeLiteralFaults(path, rel string) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, readIfExists(path), 0)
	if err != nil {
		return []string{rel + " cannot be parsed for route literals: " + err.Error()}
	}
	var faults []string
	ast.Inspect(file, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if !ok {
			return true
		}
		value, ok := stringLiteral(literal)
		if !ok {
			return true
		}
		if value == refusalroute.NextField || containsAnyNeedle(value) {
			// %q keeps a literal that is not line-safe on the fault's one line.
			faults = append(faults, fmt.Sprintf("%s:%d composes a route outside the registry in the literal %q; print it through internal/refusalroute", rel, fset.Position(literal.Pos()).Line, value))
		}
		return true
	})
	return faults
}

func containsAnyNeedle(value string) bool {
	for _, needle := range routeBypassNeedles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

// TestNoWriteVerbComposesARouteOutsideTheRegistry is RR49. After every verb prints its
// routes through the registry, no write-verb literal spells a route.
func TestNoWriteVerbComposesARouteOutsideTheRegistry(t *testing.T) {
	for _, fault := range routeBypassFaults(NewHarness(t).KitRoot) {
		t.Error(fault)
	}
}

// TestRouteBypassLaneNameIsTheRR49Test guards the name that the commit lane runs. A run
// pattern that matches no test passes, so a renamed RR49 test would leave the lane check
// green and empty. The test reads this package's own source, which is what the lane runs.
func TestRouteBypassLaneNameIsTheRR49Test(t *testing.T) {
	const path = "refusal_route_bypass_test.go"
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, declaration := range file.Decls {
		if function, ok := declaration.(*ast.FuncDecl); ok && function.Name.Name == registry.RouteBypassTest {
			return
		}
	}
	t.Fatalf("%s declares no %s, so the commit lane runs no route bypass test", path, registry.RouteBypassTest)
}

// TestRouteBypassCheckBites is RR50 and RR65. Each scanned package holds one planted route
// literal, the commitment tail in internal/commitment/repository included, and each plant
// reds with a fault that names its file. Each literal holds one needle, and each needle has
// a plant, so a needle or a package that leaves the check reds here. The same literal in a
// test file or a comment does not red, and a literal that is not line-safe keeps the fault
// on one line. An empty tree reds once for each scanned package.
func TestRouteBypassCheckBites(t *testing.T) {
	t.Run("planted routes", func(t *testing.T) {
		const (
			nextField   = `"refused{reason=dirty,next=bench doctor}"`
			nextTable   = "`next[1]:\n  bench doctor`"
			tail        = `"no delivery binding; run bench commitment plan --input <file>"`
			bareLabel   = `"next"`
			notLineSafe = `"refused\x1b[31m\nnext=bench doctor"`
		)
		joined := strconv.Quote("bench doctor" + refusalroute.StepJoiner + "bench gate")
		// The planted packages are spelled apart from registry.RouteBypassPackages, so a package that
		// leaves the list leaves its plant unscanned and reds. Each scanned package must hold
		// a plant, so the plants follow the list.
		planted := map[string]string{
			"internal/commit":                 nextTable,
			"internal/commitment":             joined,
			"internal/commitment/commitcmd":   tail,
			"internal/commitment/repository":  tail,
			"internal/gate":                   notLineSafe,
			"internal/gate/authorization":     joined,
			"internal/landing":                bareLabel,
			"internal/reviewrecord/recordcmd": nextField,
			"internal/worktree":               nextField,
			"internal/worktree/landingpolicy": nextTable,
		}
		for _, pkg := range registry.RouteBypassPackages {
			if _, ok := planted[pkg]; !ok {
				t.Errorf("scanned package %s holds no planted route", pkg)
			}
		}
		source := func(pkg, literal string) string {
			return fmt.Sprintf("package %s\n\nimport \"errors\"\n\nvar planted = errors.New(%s)\n", path.Base(pkg), literal)
		}
		files := map[string]string{}
		var plants []string
		for pkg, literal := range planted {
			rel := pkg + "/planted.go"
			files[rel] = source(pkg, literal)
			plants = append(plants, rel)
		}
		quiet := map[string]string{
			"internal/worktree/planted_test.go": source("internal/worktree", nextField),
			"internal/worktree/prose.go":        "package worktree\n\n// A refusal prints its route on a next= line; then the operator runs it.\nfunc prose() {}\n",
		}
		maps.Copy(files, quiet)
		faults := routeBypassFaults(throwawayRoot{files: files}.build(t))
		for _, rel := range plants {
			if !containsDiagnostic(faults, rel+":") {
				t.Errorf("faults = %q, want one that names the planted route in %s", faults, rel)
			}
		}
		for rel := range quiet {
			if containsDiagnostic(faults, rel+":") {
				t.Errorf("faults = %q, want none that names %s, which prints no route", faults, rel)
			}
		}
		for _, fault := range faults {
			if !sanitize.LineSafe(fault) {
				t.Errorf("fault %q is not line-safe", fault)
			}
		}
	})
	t.Run("an empty tree", func(t *testing.T) {
		if faults := routeBypassFaults(t.TempDir()); len(faults) != len(registry.RouteBypassPackages) {
			t.Fatalf("faults = %q, want one for each of the %d scanned packages", faults, len(registry.RouteBypassPackages))
		}
	})
}
