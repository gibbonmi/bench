package tickets

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// BindingKind is the obligation class of one binding row, in its exact case-sensitive
// spelling. A Core row binds its files to every write under its prefix. A Grammar row binds
// its files only to a write that changes a grammar producer: a declaration that shapes the
// arguments, the help, the public inventory, or the routing of a command.
type BindingKind string

const (
	BindingCore    BindingKind = "Core"
	BindingGrammar BindingKind = "Grammar"
)

// BindingRow binds one package prefix to the files a ticket must co-name when
// it writes into that package. The prefix matches a `Writes:` path through Covers.
type BindingRow struct {
	// Prefix is the bound package, repo-relative and slash-spelled.
	Prefix string
	// Kind tells when the row's files apply to a write under Prefix.
	Kind BindingKind
	// Files are the registries the ticket must name beside its own edit.
	Files []string
}

// commandRegistries is the file set every grammar change co-names: the help
// projection and its assertion, the envelope cases, the approved AXI query list,
// and the subcommand routing census. A verb that ships without one of these
// reaches the gate with a registry nobody updated.
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
// are Grammar rows over the shared registry set. The terminal and renderer rows
// are Core rows over their pinned test files.
var bindings = []BindingRow{
	{Prefix: "cmd/bench", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/anchors", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/assessment", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/commitment", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/consumers", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/coverage", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/diff", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/guards", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/harnesses", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/learnings", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/maps", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/probe", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/roadmap", Kind: BindingGrammar, Files: commandRegistries},
	{Prefix: "internal/terminal", Kind: BindingCore, Files: []string{
		"internal/adopt/setup_prompt_test.go",
		"internal/adopt/setup_test.go",
		"internal/gate/command_test.go",
	}},
	{Prefix: "internal/toon", Kind: BindingCore, Files: []string{
		"internal/conformance/data_handling_test.go",
		"internal/toon/toon_test.go",
	}},
	{Prefix: "internal/worktree", Kind: BindingGrammar, Files: commandRegistries},
}

// Bindings returns the ordered binding registry.
func Bindings() []BindingRow {
	return append([]BindingRow(nil), bindings...)
}

// CommandPackage reports whether path sits under a bound command package: a Grammar
// binding row.
func CommandPackage(path string) bool {
	for _, row := range bindings {
		if row.Kind == BindingGrammar && Covers(row.Prefix, path) {
			return true
		}
	}
	return false
}

// GrammarRoot is one grammar producer selector: the package directory that declares it and
// the exact, case-sensitive name of its package-level declaration.
type GrammarRoot struct {
	Dir, Name string
}

// String is the selector's spelling: its package directory, a dot, and its name.
func (r GrammarRoot) String() string { return r.Dir + "." + r.Name }

// grammarParser is the ordinary producer selector: each grammar value passed to it is a
// grammar producer.
var grammarParser = GrammarRoot{Dir: "internal/usage", Name: "Parse"}

// grammarRoots are the custom producer selectors. Each declaration shapes command arguments,
// help, public inventory, or routing without being passed to grammarParser itself.
var grammarRoots = []GrammarRoot{
	{Dir: "cmd/bench", Name: "commandRegistry"},
	{Dir: "internal/preflight/evidencecmd", Name: "flagTable"},
	{Dir: "internal/preflight/evidencecmd", Name: "operations"},
	{Dir: "internal/reviewrecord/recordcmd", Name: "forms"},
}

// GrammarParser returns the ordinary producer selector.
func GrammarParser() GrammarRoot { return grammarParser }

// GrammarRoots returns the custom producer selectors in registry order.
func GrammarRoots() []GrammarRoot {
	return append([]GrammarRoot(nil), grammarRoots...)
}

// GrammarFiles returns the public inventory files that every grammar change requires. A
// custom root outside every Grammar row requires the same files, so a new public form
// cannot omit a shared inventory because its package has no prefix binding.
func GrammarFiles() []string {
	return append([]string(nil), commandRegistries...)
}

// BindingRegistry is the repository path of this registry's source. A closure reason cites
// it as the rule owner, so no consumer spells the location a second time.
const BindingRegistry = "internal/tickets/registry_data.go"

// BindingRule is one binding that path falls under: the bound prefix, its kind, and one
// file the prefix requires beside the path.
type BindingRule struct {
	Prefix, File string
	Kind         BindingKind
}

// BindingRules returns every rule whose prefix covers path, in registry order, whatever its
// kind. A row never requires path itself.
func BindingRules(path string) []BindingRule {
	var rules []BindingRule
	for _, row := range bindings {
		if !Covers(row.Prefix, path) {
			continue
		}
		for _, file := range row.Files {
			if file != path {
				rules = append(rules, BindingRule{Prefix: row.Prefix, File: file, Kind: row.Kind})
			}
		}
	}
	return rules
}

// BoundFiles returns every file bound to a row whose prefix covers path, sorted
// and deduplicated, as if every write changed grammar. A path under no bound
// package answers none. It projects BindingRules, so the closure and this answer
// read one table.
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

// BindingLines maps each binding row's prefix, and each grammar root's selector, to its line
// in src, the source bytes of this registry. A caller that holds verified kit source names
// the row's actual location this way. Source that does not parse, or that declares no
// bindings table, is an error.
func BindingLines(src []byte) (map[string]int, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, BindingRegistry, src, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	// Each table names the row fields whose values spell the row's key.
	tables := map[string][]string{"bindings": {"Prefix"}, "grammarRoots": {"Dir", "Name"}}
	lines := map[string]int{}
	found := false
	ast.Inspect(file, func(node ast.Node) bool {
		spec, ok := node.(*ast.ValueSpec)
		if !ok || len(spec.Names) != 1 || tables[spec.Names[0].Name] == nil || len(spec.Values) != 1 {
			return true
		}
		table, ok := spec.Values[0].(*ast.CompositeLit)
		if !ok {
			return true
		}
		found = found || spec.Names[0].Name == "bindings"
		for _, element := range table.Elts {
			if row, ok := element.(*ast.CompositeLit); ok {
				if key, ok := rowKey(row, tables[spec.Names[0].Name]); ok {
					lines[key] = fset.Position(row.Pos()).Line
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

// rowKey joins the named string fields of one row literal with dots, in field order.
func rowKey(row *ast.CompositeLit, names []string) (string, bool) {
	values := make([]string, 0, len(names))
	for _, name := range names {
		value, ok := stringField(row, name)
		if !ok {
			return "", false
		}
		values = append(values, value)
	}
	return strings.Join(values, "."), true
}

// stringField is the string literal value of one keyed field of a row literal.
func stringField(row *ast.CompositeLit, name string) (string, bool) {
	for _, element := range row.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, keyOK := field.Key.(*ast.Ident)
		value, valueOK := field.Value.(*ast.BasicLit)
		if keyOK && valueOK && key.Name == name && value.Kind == token.STRING {
			unquoted, err := strconv.Unquote(value.Value)
			return unquoted, err == nil
		}
	}
	return "", false
}
