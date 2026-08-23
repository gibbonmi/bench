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

// skipOwnerDir is the one package allowed to call the testing package's skip methods.
// Every other skip in the suite goes through its two helpers, so a skip always leaves a
// structured line behind. A bare t.Skip stays invisible under non-verbose `go test`, and a
// skipped security assertion that prints nothing looks the same as a passing one.
const skipOwnerDir = "internal/capability"

// skipMethods lists the testing.TB methods that end a test without running it. The check
// matches any receiver, not only the literal `t.`, because a skip stays just as invisible
// through a differently named variable. An allowlist a rename can evade is no allowlist.
var skipMethods = map[string]bool{"Skip": true, "Skipf": true, "SkipNow": true}

// checkSkipOwnership reports every skip call outside skipOwnerDir. It reads the module's own
// Go source through the AST, not by text search. A forbidden call spelled inside a string
// literal is data, never a violation. This covers this check's own diagnostic, its bite
// proof's synthetic sources, and the canary fixture that reintroduces one. Only a real call
// expression counts as a violation.
func checkSkipOwnership(root string) []string {
	var diags []string
	for _, path := range moduleGoFiles(root) {
		rel := slashRel(root, path)
		if rel == skipOwnerDir || strings.HasPrefix(rel, skipOwnerDir+"/") {
			continue
		}
		diags = append(diags, skipOwnershipDiags(rel, path)...)
	}
	return uniqueSorted(diags)
}

// moduleGoFiles lists the Go source the module compiles: the root package's files plus
// everything under cmd and internal. Fixture payloads under tests/canary are inputs other
// canaries copy into their own roots, not source this module builds.
func moduleGoFiles(root string) []string {
	var paths []string
	if entries, err := os.ReadDir(root); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
				paths = append(paths, filepath.Join(root, entry.Name()))
			}
		}
	}
	for _, top := range []string{"cmd", "internal"} {
		_ = filepath.WalkDir(filepath.Join(root, top), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
	}
	return paths
}

func skipOwnershipDiags(rel, path string) []string {
	body := readIfExists(path)
	if !strings.Contains(body, "Skip") {
		return nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for skip ownership: " + err.Error()}
	}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || !skipMethods[selector.Sel.Name] {
			return true
		}
		diags = append(diags, fmt.Sprintf(
			"%s calls %s outside %s (line %d); route the skip through capability.Capability(t, class, reason) for a host capability the test needs, or capability.Environment(t, reason) for a staging fact",
			rel, expressionText(fset, call.Fun), skipOwnerDir, fset.Position(call.Pos()).Line))
		return true
	})
	return diags
}

// TestSkipOwnershipBites is the recorded bite proof for checkSkipOwnership. It runs against
// a synthetic tree, not the repo, and walks three states:
//
//   - the owner package skips freely
//   - one bare skip appears elsewhere
//   - the same forbidden text appears only as a string literal
//
// The last case matters because this check's own source and its canary fixture both carry
// that text as data.
func TestSkipOwnershipBites(t *testing.T) {
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

	write(skipOwnerDir+"/capability.go", "package capability\n\nfunc Environment(t TB, reason string) {\n\tt.Skip(reason)\n}\n")
	if diags := checkSkipOwnership(root); len(diags) != 0 {
		t.Fatalf("owner package alone: want no diagnostics, got %v", diags)
	}

	write("internal/example/example_test.go", "package example\n\nimport \"testing\"\n\nfunc TestX(t *testing.T) {\n\tt.Skip(\"no host support\")\n}\n")
	diags := checkSkipOwnership(root)
	if len(diags) != 1 || !strings.Contains(diags[0], "internal/example/example_test.go calls t.Skip outside "+skipOwnerDir) {
		t.Fatalf("bare skip outside the owner: want one diagnostic naming it, got %v", diags)
	}
	if !strings.Contains(diags[0], "capability.Environment(t, reason)") {
		t.Fatalf("diagnostic does not name the replacement helpers: %q", diags[0])
	}

	// A renamed receiver must not evade the owner allowlist, and Skipf is the same defect.
	write("internal/example/example_test.go", "package example\n\nimport \"testing\"\n\nfunc TestX(tb *testing.T) {\n\ttb.Skipf(\"no %s\", \"host support\")\n}\n")
	if diags := checkSkipOwnership(root); len(diags) != 1 || !strings.Contains(diags[0], "calls tb.Skipf outside") {
		t.Fatalf("renamed receiver calling Skipf: want one diagnostic, got %v", diags)
	}

	// SkipNow ends the test with no message at all. It is the most invisible of the
	// three, so the check must reach it.
	write("internal/example/example_test.go", "package example\n\nimport \"testing\"\n\nfunc TestX(t *testing.T) {\n\tt.SkipNow()\n}\n")
	if diags := checkSkipOwnership(root); len(diags) != 1 || !strings.Contains(diags[0], "calls t.SkipNow outside") {
		t.Fatalf("bare SkipNow outside the owner: want one diagnostic, got %v", diags)
	}

	// The forbidden call appears as data, never as a call. This is the state that would
	// make the check flag its own source and the canary fixture that reintroduces a bare skip.
	write("internal/example/example_test.go", "package example\n\nconst forbidden = \"t.Skip(\\\"reason\\\")\"\n\n// t.Skipf is also named in a comment here.\n")
	if diags := checkSkipOwnership(root); len(diags) != 0 {
		t.Fatalf("forbidden text in a literal and a comment: want no diagnostics, got %v", diags)
	}
}
