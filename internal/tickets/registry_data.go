package tickets

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"sort"
	"strconv"
)

// BindingRow binds one package prefix to the files a ticket must co-name when
// it writes into that package. The prefix matches a `Writes:` path through Covers.
type BindingRow struct {
	// Prefix is the bound package, repo-relative and slash-spelled.
	Prefix string
	// Files are the registries the ticket must name beside its own edit.
	Files []string
}

// commandRegistries is the file set every command-package ticket co-names: the
// help projection and its assertion, the envelope cases, the approved AXI query
// list, and the subcommand routing census. A verb that ships without one of
// these reaches the gate with a registry nobody updated.
var commandRegistries = []string{
	"cmd/bench/command_registry.go",
	"cmd/bench/command_registry_test.go",
	"cmd/bench/help_inventory_test.go",
	"internal/conformance/axi_query_registry_test.go",
	"internal/conformance/subcommand_routing_table_test.go",
}

// seedOwners are the package prefixes the table binds beyond the command rows:
// the dispatcher, the renderer, and the terminal-lifecycle owner. Each one owns
// pinned test files a ticket must co-name, so the table is the one place their
// names are stated.
var seedOwners = []string{"cmd/bench", "internal/toon", "internal/terminal"}

// SeedOwners returns the owner prefixes the binding registry must carry.
func SeedOwners() []string {
	return append([]string(nil), seedOwners...)
}

// bindings is the ordered binding registry, sorted by prefix. The command rows
// carry the shared registry set. The rows for the seed owners bind their pinned
// test files instead.
var bindings = []BindingRow{
	{Prefix: "cmd/bench", Files: commandRegistries},
	{Prefix: "internal/anchors", Files: commandRegistries},
	{Prefix: "internal/assessment", Files: commandRegistries},
	{Prefix: "internal/commitment", Files: commandRegistries},
	{Prefix: "internal/consumers", Files: commandRegistries},
	{Prefix: "internal/coverage", Files: commandRegistries},
	{Prefix: "internal/diff", Files: commandRegistries},
	{Prefix: "internal/guards", Files: commandRegistries},
	{Prefix: "internal/harnesses", Files: commandRegistries},
	{Prefix: "internal/learnings", Files: commandRegistries},
	{Prefix: "internal/maps", Files: commandRegistries},
	{Prefix: "internal/probe", Files: commandRegistries},
	{Prefix: "internal/roadmap", Files: commandRegistries},
	{Prefix: "internal/terminal", Files: []string{
		"internal/adopt/setup_prompt_test.go",
		"internal/adopt/setup_test.go",
		"internal/gate/command_test.go",
	}},
	{Prefix: "internal/toon", Files: []string{
		"internal/conformance/data_handling_test.go",
		"internal/toon/toon_test.go",
	}},
	{Prefix: "internal/worktree", Files: commandRegistries},
}

// Bindings returns the ordered binding registry.
func Bindings() []BindingRow {
	return append([]BindingRow(nil), bindings...)
}

// CommandPackage reports whether path sits under a bound command package: a binding row
// whose files are the command registry set.
func CommandPackage(path string) bool {
	for _, row := range bindings {
		if slices.Equal(row.Files, commandRegistries) && Covers(row.Prefix, path) {
			return true
		}
	}
	return false
}

// BindingRegistry is the repository path of this registry's source. A closure reason cites
// it as the rule owner, so no consumer spells the location a second time.
const BindingRegistry = "internal/tickets/registry_data.go"

// BindingRule is one binding that path falls under: the bound prefix and one file the
// prefix requires beside the path.
type BindingRule struct {
	Prefix, File string
}

// BindingRules returns every rule whose prefix covers path, in registry order. A row never
// requires path itself.
func BindingRules(path string) []BindingRule {
	var rules []BindingRule
	for _, row := range bindings {
		if !Covers(row.Prefix, path) {
			continue
		}
		for _, file := range row.Files {
			if file != path {
				rules = append(rules, BindingRule{Prefix: row.Prefix, File: file})
			}
		}
	}
	return rules
}

// BoundFiles returns every file bound to a row whose prefix covers path, sorted
// and deduplicated. A path under no bound package answers none. It projects
// BindingRules, so the closure and this answer read one table.
func BoundFiles(path string) []string {
	var files []string
	for _, rule := range BindingRules(path) {
		if !slices.Contains(files, rule.File) {
			files = append(files, rule.File)
		}
	}
	sort.Strings(files)
	return files
}

// BindingLines maps each binding row's prefix to its line in src, the source bytes of
// this registry. A caller that holds verified kit source names the row's actual location
// this way. Source that does not parse, or that declares no bindings table, is an error.
func BindingLines(src []byte) (map[string]int, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, BindingRegistry, src, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	lines := map[string]int{}
	found := false
	ast.Inspect(file, func(node ast.Node) bool {
		spec, ok := node.(*ast.ValueSpec)
		if !ok || len(spec.Names) != 1 || spec.Names[0].Name != "bindings" || len(spec.Values) != 1 {
			return true
		}
		table, ok := spec.Values[0].(*ast.CompositeLit)
		if !ok {
			return true
		}
		found = true
		for _, element := range table.Elts {
			if row, ok := element.(*ast.CompositeLit); ok {
				if prefix, ok := rowPrefix(row); ok {
					lines[prefix] = fset.Position(row.Pos()).Line
				}
			}
		}
		return false
	})
	if !found {
		return nil, fmt.Errorf("%s declares no bindings table", BindingRegistry)
	}
	return lines, nil
}

// rowPrefix is the Prefix string literal of one binding row literal.
func rowPrefix(row *ast.CompositeLit) (string, bool) {
	for _, element := range row.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, keyOK := field.Key.(*ast.Ident)
		value, valueOK := field.Value.(*ast.BasicLit)
		if keyOK && valueOK && key.Name == "Prefix" && value.Kind == token.STRING {
			prefix, err := strconv.Unquote(value.Value)
			return prefix, err == nil
		}
	}
	return "", false
}
