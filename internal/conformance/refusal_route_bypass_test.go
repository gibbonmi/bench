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

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// routeBypassPackages are the packages in which a write verb raises or prints a refusal.
// The scan reads each one's own directory and no subdirectory, so a package that this list
// omits goes unscanned.
var routeBypassPackages = []string{
	"internal/commit",
	"internal/commitment",
	"internal/commitment/commitcmd",
	"internal/commitment/repository",
	"internal/gate",
	"internal/gate/authorization",
	"internal/landing",
	"internal/worktree",
	"internal/worktree/landingpolicy",
}

// routeBypassAllowed is the reviewed allowlist: each file of a scanned package whose routes
// serve a non-write verb, which the registry does not own, keyed to that non-write caller.
var routeBypassAllowed = map[string]string{
	"internal/worktree/path.go":        "the target refusal of bench worktree path, exec, show, build, and create",
	"internal/worktree/build.go":       "the success route of bench worktree build",
	"internal/worktree/tree_target.go": "the success route of bench worktree create",
}

// routeBypassNeedles are the route spellings that only the registry writes: the next label
// as a record field and as a table header, the commitment route tail, and the route joiner.
// The spec pins the tail and the joiner, so the check spells them.
var routeBypassNeedles = []string{refusalroute.NextField + "=", refusalroute.NextField + "[", "run bench ", "; then "}

// routeBypassFaults names each string literal in a production Go file of a scanned package
// that composes a route outside the registry: one that holds a needle, or that is the bare
// next label. A comment does not print, so the scan reads literals and not text. A scanned
// package with no production file is a fault too, so a moved package cannot leave the scan.
func routeBypassFaults(root string) []string {
	var faults []string
	for _, pkg := range routeBypassPackages {
		entries, _ := os.ReadDir(filepath.Join(root, filepath.FromSlash(pkg)))
		scanned := 0
		for _, entry := range entries {
			name := entry.Name()
			rel := pkg + "/" + name
			if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
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

// TestRouteBypassCheckBites is RR50 and RR65. Each planted route literal, the commitment
// tail in internal/commitment/repository included, reds with a fault that names its file.
// The same literal in a test file or a comment does not red, and a literal that is not
// line-safe keeps the fault on one line. An empty tree reds once for each scanned package.
func TestRouteBypassCheckBites(t *testing.T) {
	t.Run("planted routes", func(t *testing.T) {
		plants := map[string]string{
			"internal/worktree/planted.go":              "package worktree\n\nvar planted = \"refused{reason=dirty,next=bench doctor}\"\n",
			"internal/commit/planted.go":                "package commit\n\nvar planted = `next[1]:\n  bench doctor`\n",
			"internal/commitment/repository/planted.go": "package repository\n\nimport \"errors\"\n\nvar planted = errors.New(\"no delivery binding; run bench commitment plan --input <file>\")\n",
			"internal/landing/planted.go":               "package landing\n\nvar planted = \"next\"\n",
			"internal/gate/planted.go":                  "package gate\n\nvar planted = \"refused\\x1b[31m\\nnext=bench doctor\"\n",
		}
		quiet := map[string]string{
			"internal/worktree/planted_test.go": plants["internal/worktree/planted.go"],
			"internal/worktree/prose.go":        "package worktree\n\n// A refusal prints its route on a next= line; then the operator runs it.\nfunc prose() {}\n",
		}
		files := map[string]string{}
		for _, set := range []map[string]string{plants, quiet} {
			for path, body := range set {
				files[path] = body
			}
		}
		faults := routeBypassFaults(throwawayRoot{files: files}.build(t))
		for path := range plants {
			if !containsDiagnostic(faults, path+":") {
				t.Errorf("faults = %q, want one that names the planted route in %s", faults, path)
			}
		}
		for path := range quiet {
			if containsDiagnostic(faults, path+":") {
				t.Errorf("faults = %q, want none that names %s, which prints no route", faults, path)
			}
		}
		for _, fault := range faults {
			if !sanitize.LineSafe(fault) {
				t.Errorf("fault %q is not line-safe", fault)
			}
		}
	})
	t.Run("an empty tree", func(t *testing.T) {
		if faults := routeBypassFaults(t.TempDir()); len(faults) != len(routeBypassPackages) {
			t.Fatalf("faults = %q, want one for each of the %d scanned packages", faults, len(routeBypassPackages))
		}
	})
}
