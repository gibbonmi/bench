package conformance

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// capabilityImport is the import path of the skip owner.
const capabilityImport = benchModule + skipOwnerDir

// unavailableClasses are the capability classes whose skips route through
// capability.Unavailable, so a low-level Capability call of one outside the owner fails.
var unavailableClasses = []string{"Fifo", "Symlink"}

// unavailableOwnershipDiags reports each low-level Capability call in rel whose class is a
// FIFO or symlink class, under any import name. Other classes keep the low-level helper.
func unavailableOwnershipDiags(rel, path string) []string {
	body := readIfExists(path)
	if !strings.Contains(body, "Capability(") {
		return nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for unavailable-capability ownership: " + err.Error()}
	}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || guardProjectionImported(file, call, capabilityImport) != "Capability" || len(call.Args) < 2 {
			return true
		}
		if class := importedName(file, call.Args[1], capabilityImport); slices.Contains(unavailableClasses, class) {
			diags = append(diags, fmt.Sprintf("%s calls the low-level Capability for the %s class (line %d); route the skip through capability.Unavailable(t, class, detail)", rel, class, fset.Position(call.Pos()).Line))
		}
		return true
	})
	return diags
}

// unavailableSite is one capability.Unavailable call: its file, its enclosing declaration,
// its ordinal among the calls of that declaration, and its class as the source names it.
type unavailableSite struct {
	Path   string `json:"path"`
	Symbol string `json:"symbol"`
	Case   int    `json:"case"`
	Class  string `json:"class"`
}

// unavailableCallerSites lists each capability.Unavailable call outside the owner.
func unavailableCallerSites(root string) ([]unavailableSite, error) {
	var sites []unavailableSite
	for _, path := range moduleGoFiles(root) {
		rel := slashRel(root, path)
		body := readIfExists(path)
		if strings.HasPrefix(rel, skipOwnerDir+"/") || !strings.Contains(body, "Unavailable(") {
			continue
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, body, 0)
		if err != nil {
			return nil, err
		}
		cases := map[string]int{}
		inspectDecls(file, func(symbol string, node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || guardProjectionImported(file, call, capabilityImport) != "Unavailable" || len(call.Args) < 2 {
				return true
			}
			class := importedName(file, call.Args[1], capabilityImport)
			if class == "" {
				class = expressionText(fset, call.Args[1])
			}
			cases[symbol]++
			sites = append(sites, unavailableSite{Path: rel, Symbol: symbol, Case: cases[symbol], Class: class})
			return true
		})
	}
	return sites, nil
}

// TestUnavailableOwnershipBites is GF26. A direct FIFO or symlink Capability call outside
// the owner fails the registered skip check, also under an import alias. The owner, the
// Unavailable route, and another class keep the low-level helper.
func TestUnavailableOwnershipBites(t *testing.T) {
	const header = "package example\n\nimport (\n\t\"testing\"\n\n\t%s \"" + capabilityImport + "\"\n)\n\nfunc TestX(t *testing.T) {\n\t%s\n}\n"
	for _, tc := range []struct {
		name, rel, alias, call, want string
	}{
		{"fifo", "internal/example/example_test.go", "capability", `capability.Capability(t, capability.Fifo, "no fifo")`, "Fifo class"},
		{"aliased symlink", "internal/example/example_test.go", "caps", `caps.Capability(t, caps.Symlink, "no symlink")`, "Symlink class"},
		{"owner", skipOwnerDir + "/owner_test.go", "capability", `capability.Capability(t, capability.Fifo, "no fifo")`, ""},
		{"unavailable", "internal/example/example_test.go", "capability", `capability.Unavailable(t, capability.Fifo, "no fifo")`, ""},
		{"other class", "internal/example/example_test.go", "capability", `capability.Capability(t, capability.Privilege, "no device")`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := throwawayRoot{files: map[string]string{tc.rel: fmt.Sprintf(header, tc.alias, tc.call)}}.build(t)
			diags := checkSkipOwnership(root)
			if tc.want == "" && len(diags) != 0 || tc.want != "" && (len(diags) != 1 || !strings.Contains(diags[0], tc.want) || !strings.Contains(diags[0], "capability.Unavailable")) {
				t.Fatalf("diagnostics = %v, want %q", diags, tc.want)
			}
		})
	}
}

// unavailableClassFixture is the independent record of each migrated FIFO and symlink caller
// and the class it held before the migration.
const unavailableClassFixture = "testdata/fixture-skip-classes.json"

// TestUnavailableCallerClasses is GF31. Each capability.Unavailable call in the tree matches
// one row of the independent class fixture by file, declaration, and case, with the same
// class. A missing row, an unrecorded call, and a changed class each fail.
//
// The fixture is independent of the call sites under AGENTS.md's independent-expectation
// exception. A fixture read from the tree would change with the tree, so a caller that moved
// from the FIFO class to the symlink class would pass. The named mutation is that move in
// the testrepo socket helper, and this test then fails on that row.
func TestUnavailableCallerClasses(t *testing.T) {
	data, err := os.ReadFile(filepath.FromSlash(unavailableClassFixture))
	if err != nil {
		t.Fatal(err)
	}
	var recorded []unavailableSite
	if err := json.Unmarshal(data, &recorded); err != nil {
		t.Fatal(err)
	}
	found, err := unavailableCallerSites(NewHarness(t).KitRoot)
	if err != nil {
		t.Fatal(err)
	}
	type key struct {
		path, symbol string
		ordinal      int
	}
	want := map[key]string{}
	for _, site := range recorded {
		want[key{site.Path, site.Symbol, site.Case}] = site.Class
	}
	for _, site := range found {
		k := key{site.Path, site.Symbol, site.Case}
		class, ok := want[k]
		switch {
		case !ok:
			t.Errorf("%s %s case %d calls Unavailable with class %s, and the class fixture has no row for it", site.Path, site.Symbol, site.Case, site.Class)
		case class != site.Class:
			t.Errorf("%s %s case %d changed class from %s to %s", site.Path, site.Symbol, site.Case, class, site.Class)
		}
		delete(want, k)
	}
	for k, class := range want {
		t.Errorf("%s %s case %d (class %s) is in the class fixture, and the tree has no such Unavailable call", k.path, k.symbol, k.ordinal, class)
	}
}
