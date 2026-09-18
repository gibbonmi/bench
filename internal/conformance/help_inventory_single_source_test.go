package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func checkRoadmapHelpLine(helpInventorySource string) []string {
	if strings.Contains(helpInventorySource, "show the top 10 roadmap rows + drain state") {
		return nil
	}
	return []string{"cmd/bench/main.go help inventory does not describe the top-10 board"}
}

const helpInventoryTitle = "bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants."

// helpRowProjections names the one command whose help rows derive from its own operation
// registry, and the projection call that must supply them. Every other public command keeps
// literal helpRow metadata.
var helpRowProjections = map[string]string{"preflight": "preflightHelpRows"}

// isHelpRowProjection reports whether a publicInventory call spreads exactly one call of
// the named projection and nothing else.
func isHelpRowProjection(args []ast.Expr, spread bool, projection string) bool {
	if len(args) != 1 || !spread {
		return false
	}
	call, ok := args[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	name, ok := call.Fun.(*ast.Ident)
	return ok && name.Name == projection
}

func checkHelpInventorySingleSource(root string) []string {
	var diags []string
	occurrences := 0
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := filepath.ToSlash(strings.TrimPrefix(path, root+string(filepath.Separator)))
		if entry.IsDir() {
			if rel == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		productionGo := strings.HasSuffix(rel, ".go") && !strings.HasSuffix(rel, "_test.go")
		publicSurface := rel == "bin/bench.sh" || rel == ".bench/BENCH.md" || rel == ".bench/BENCH-reference.md"
		if productionGo || publicSurface {
			if !entry.Type().IsRegular() {
				diags = append(diags, rel+" is not a regular help-inventory source")
				return nil
			}
			occurrences += strings.Count(readIfExists(path), helpInventoryTitle)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, "scan help inventory sources: "+err.Error())
	}
	if occurrences != 1 {
		diags = append(diags, fmt.Sprintf("help inventory title appears %d times on production surfaces, want exactly 1", occurrences))
	}

	path := filepath.Join(root, "cmd", "bench", "main.go")
	entries, err := parseCommandRegistry(path, readIfExists(path))
	if err != nil {
		return append(diags, "help inventory commandRegistry: "+err.Error())
	}
	helpFound := false
	publicCommands := map[string]bool{}
	for _, entry := range entries {
		fields := entry.fields["Inventory"]
		if len(fields) != 1 {
			diags = append(diags, fmt.Sprintf("help inventory command %q has %d classifications, want exactly 1", entry.name, len(fields)))
			continue
		}
		switch value := fields[0].(type) {
		case *ast.Ident:
			if value.Name != "internalInventory" {
				diags = append(diags, fmt.Sprintf("help inventory command %q has unknown classification %q", entry.name, value.Name))
			}
		case *ast.CallExpr:
			owner, ok := value.Fun.(*ast.Ident)
			if !ok || owner.Name != "publicInventory" {
				diags = append(diags, fmt.Sprintf("help inventory command %q must use publicInventory or internalInventory", entry.name))
				continue
			}
			publicCommands[entry.name] = true
			if entry.name == "help" {
				helpFound = true
				if len(value.Args) != 0 {
					diags = append(diags, "help command must not own a standalone command-row catalog")
				}
				continue
			}
			if projection, projected := helpRowProjections[entry.name]; projected {
				if !isHelpRowProjection(value.Args, value.Ellipsis.IsValid(), projection) {
					diags = append(diags, fmt.Sprintf("public help command %q must derive its rows from exactly one %s(...) projection", entry.name, projection))
				}
				continue
			}
			if len(value.Args) == 0 {
				diags = append(diags, fmt.Sprintf("public help command %q owns no inventory rows", entry.name))
			}
			for _, argument := range value.Args {
				row, ok := argument.(*ast.CompositeLit)
				if !ok {
					diags = append(diags, fmt.Sprintf("public help command %q has a row outside its literal helpRow metadata", entry.name))
					continue
				}
				rowType, typed := row.Type.(*ast.Ident)
				if !typed || rowType.Name != "helpRow" {
					diags = append(diags, fmt.Sprintf("public help command %q has a row outside its literal helpRow metadata", entry.name))
				}
			}
		default:
			diags = append(diags, fmt.Sprintf("help inventory command %q has malformed classification", entry.name))
		}
	}
	if !helpFound {
		diags = append(diags, "help command is not classified as public inventory")
	}
	if !publicCommands["repair"] {
		diags = append(diags, "CLI inventory omits bench repair")
	}
	for _, rel := range []string{"cmd/bench/main.go", "cmd/bench/command_registry.go"} {
		sourcePath := filepath.Join(root, filepath.FromSlash(rel))
		file, parseErr := parser.ParseFile(token.NewFileSet(), sourcePath, readIfExists(sourcePath), 0)
		if parseErr != nil {
			diags = append(diags, rel+" cannot be parsed for standalone help catalogs: "+parseErr.Error())
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			text, unquoteErr := strconv.Unquote(literal.Value)
			if unquoteErr == nil && (strings.HasPrefix(text, "  bench ") || strings.HasPrefix(text, "  bash bin/bench.sh ")) {
				diags = append(diags, rel+" carries a standalone rendered help command row")
			}
			return true
		})
	}
	registryOwner := readIfExists(filepath.Join(root, "cmd", "bench", "command_registry.go"))
	if !strings.Contains(registryOwner, "definition.Name + row.Suffix") {
		diags = append(diags, "help renderer does not derive each command token from commandRegistry Name")
	}
	return diags
}

func TestHelpInventorySingleSourceBites(t *testing.T) {
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
	title := "bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants."
	registry := "package main\n\nvar commandRegistry = []commandDefinition{\n\t{Name: \"help\", Inventory: publicInventory()},\n\t{Name: \"repair\", WrapperOnly: true, Inventory: publicInventory(helpRow{Order: 1, Description: \"repair\"})},\n\t{Name: \"status\", Inventory: publicInventory(helpRow{Order: 2, Description: \"state\"})},\n\t{Name: \"plumbing\", Inventory: internalInventory},\n}\n"
	write("cmd/bench/main.go", registry)
	write("cmd/bench/command_registry.go", "package main\nconst title = \""+title+"\"\nfunc render(definition commandDefinition, row helpRow) { _ = definition.Name + row.Suffix }\n")
	write("bin/bench.sh", "#!/bin/sh\n")
	write(".bench/BENCH.md", "Run `bench help`.\n")
	if diags := checkHelpInventorySingleSource(root); len(diags) != 0 {
		t.Fatalf("registry-owned inventory = %v, want no diagnostics", diags)
	}

	write("bin/bench.sh", "#!/bin/sh\n# "+title+"\n")
	if diags := checkHelpInventorySingleSource(root); !containsDiagnostic(diags, "appears 2 times") {
		t.Fatalf("duplicate production inventory = %v, want duplicate-source diagnostic", diags)
	}

	write("bin/bench.sh", "#!/bin/sh\n")
	write("cmd/bench/main.go", registry+"\nvar helpInventory = []string{\"  bench status  state\"}\n")
	if diags := checkHelpInventorySingleSource(root); !containsDiagnostic(diags, "standalone rendered help command row") {
		t.Fatalf("inventory moved beside registry = %v, want standalone-catalog diagnostic", diags)
	}

	write("cmd/bench/main.go", strings.Replace(registry, "{Name: \"status\", Inventory: publicInventory(helpRow{Order: 2, Description: \"state\"})}", "{Name: \"status\"}", 1))
	if diags := checkHelpInventorySingleSource(root); !containsDiagnostic(diags, "status\" has 0 classifications") {
		t.Fatalf("unclassified public route = %v, want classification diagnostic", diags)
	}

	// The preflight rows derive from exactly one named projection. A second row beside it,
	// another unnamed expression, a missing spread, and literal rows in its place each red,
	// and the projection name opens no exception for any other command.
	const projected = "{Name: \"preflight\", Inventory: publicInventory(preflightHelpRows(21)...)},\n}\n"
	withPreflight := func(row string) string { return strings.TrimSuffix(registry, "}\n") + row }
	write("cmd/bench/main.go", withPreflight(projected))
	if diags := checkHelpInventorySingleSource(root); len(diags) != 0 {
		t.Fatalf("named preflight projection = %v, want no diagnostics", diags)
	}
	for name, row := range map[string]string{
		"second row":         "{Name: \"preflight\", Inventory: publicInventory(preflightHelpRows(21), extraRow)},\n}\n",
		"unnamed projection": "{Name: \"preflight\", Inventory: publicInventory(otherRows(21)...)},\n}\n",
		"missing spread":     "{Name: \"preflight\", Inventory: publicInventory(preflightHelpRows(21))},\n}\n",
		"projection removed": "{Name: \"preflight\", Inventory: publicInventory(helpRow{Order: 21, Description: \"checks\"})},\n}\n",
	} {
		write("cmd/bench/main.go", withPreflight(row))
		if diags := checkHelpInventorySingleSource(root); !containsDiagnostic(diags, "\"preflight\" must derive its rows from exactly one preflightHelpRows(...) projection") {
			t.Fatalf("%s = %v, want projection diagnostic", name, diags)
		}
	}
	write("cmd/bench/main.go", withPreflight("{Name: \"other\", Inventory: publicInventory(preflightHelpRows(21)...)},\n}\n"))
	if diags := checkHelpInventorySingleSource(root); !containsDiagnostic(diags, "\"other\" has a row outside its literal helpRow metadata") {
		t.Fatalf("projection on another command = %v, want literal-row diagnostic", diags)
	}
}
