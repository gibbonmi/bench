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

// dispatchFile is the one file that decides which subcommand names exist, and
// dispatchFunc is the function holding the second of its two dispatch surfaces.
const (
	dispatchFile = "cmd/bench/main.go"
	dispatchMap  = "commands"
	dispatchFunc = "run"
)

// grammarPkg and grammarSel name the single owner of arity, flag recognition, `--`, and
// help. A subcommand recorded as routed must reference it; anything else is hand-rolling a
// second grammar. grammarHelper is the spelling the diagnostics use, derived from the two
// halves the check matches on so the message can never name something else.
const (
	grammarPkg    = "usage"
	grammarSel    = "Parse"
	grammarHelper = grammarPkg + "." + grammarSel
)

// routingEntry records how one dispatch name answers for its argument grammar: Pkg names
// the directory whose entry point must reach the helper, or Exempt states why this name is
// not required to.
type routingEntry struct {
	Pkg    string
	Exempt string
}

func routed(pkg string) routingEntry { return routingEntry{Pkg: pkg} }
func exempt(why string) routingEntry { return routingEntry{Exempt: why} }

// Reasons shared by a whole class of dispatch names, so the class is stated once.
const (
	whyPlumbing = "hook- and adapter-driven plumbing: its argv is produced by the kit, never typed by an agent, so there is no misuse for a grammar to report"
	whyNested   = "dispatches a subcommand tree rather than a flat argv; each leaf owns its own grammar"
)

// subcommandRouting is the explicit registry the routing check grades cmd/bench/main.go
// against. It is deliberately exhaustive rather than a list of the interesting cases: a
// name reaching either dispatch surface with no row here is red, which is what makes the
// check fail closed against the next subcommand somebody adds.
var subcommandRouting = map[string]routingEntry{
	"anchors":   routed("cmd/bench"),
	"commands":  routed("cmd/bench"),
	"commit":    routed("internal/commit"),
	"coverage":  routed("internal/coverage"),
	"dashboard": routed("internal/dashboard"),
	"diff":      routed("internal/diff"),
	"guards":    routed("internal/guards"),
	"handoff":   routed("internal/handoff"),
	"idea":      routed("internal/roadmap"),
	"learnings": routed("internal/learnings"),
	"maps":      routed("internal/maps"),
	"models":    routed("internal/models"),
	"outline":   routed("internal/outline"),
	// prep-release takes a flat argv with no subcommand tree, so it is routed rather
	// than exempt like the release commands beside it in the dispatch switch.
	"prep-release": routed("internal/preprelease"),
	"roadmap":      routed("internal/roadmap"),
	"status":       routed("internal/status"),
	"structure":    routed("internal/structure"),
	"test":         routed("internal/testreport"),

	"check-agent-line":    exempt(whyPlumbing),
	"freshness-check":     exempt(whyPlumbing),
	"freshness-publish":   exempt(whyPlumbing),
	"gate-go":             exempt(whyPlumbing),
	"gate-phases":         exempt(whyPlumbing),
	"guard-git":           exempt(whyPlumbing),
	"resolve-model":       exempt(whyPlumbing),
	"stop-verdict":        exempt(whyPlumbing),
	"tree-hash":           exempt(whyPlumbing),
	"worktree-hook":       exempt(whyPlumbing),
	"worktree-lease-file": exempt(whyPlumbing),
	"worktree-pool":       exempt(whyPlumbing),

	"canary":            exempt(whyNested),
	"doctor":            exempt(whyNested),
	"gate":              exempt(whyNested),
	"gate-pin":          exempt(whyNested),
	"gate-run":          exempt(whyNested),
	"init":              exempt(whyNested),
	"link":              exempt(whyNested),
	"release":           exempt(whyNested),
	"release-preflight": exempt(whyNested),
	"resume-clean":      exempt(whyNested),
	"session-inspect":   exempt(whyNested),
	"setup":             exempt(whyNested),
	"shift":             exempt(whyNested),
	"spec":              exempt(whyNested),
	"unlink":            exempt(whyNested),
	"upgrade":           exempt(whyNested),
	"worktree":          exempt(whyNested),

	"version": exempt("takes no arguments: the dispatch case prints the build-time version line and returns"),
}

// checkSubcommandRouting grades every dispatched subcommand name against the registry
// above. A check that inspected only the names it already knows would pass vacuously on
// exactly the case it exists to catch — the subcommand somebody just added — so the names
// come from the dispatch file itself, read from both of its surfaces, and an unlisted one
// is a violation.
func checkSubcommandRouting(root string) []string {
	path := filepath.Join(root, filepath.FromSlash(dispatchFile))
	body := readIfExists(path)
	if body == "" {
		return nil
	}
	names, err := dispatchNames(path, body)
	if err != nil {
		return []string{dispatchFile + " cannot be parsed for subcommand dispatch: " + err.Error()}
	}

	var diags []string
	dispatched := map[string]bool{}
	for _, name := range names {
		dispatched[name] = true
		entry, ok := subcommandRouting[name]
		if !ok {
			diags = append(diags, fmt.Sprintf("%s dispatches %q with no entry in the subcommand argument-routing registry; record it as routed through %s or as an exemption with its reason", dispatchFile, name, grammarHelper))
			continue
		}
		if entry.Exempt != "" {
			continue
		}
		// The routed claim is verified only where the package is present: a canary
		// fixture materializes the dispatch file without the tree behind it, and a
		// partial tree must not manufacture violations it cannot answer for.
		if pkg := filepath.Join(root, filepath.FromSlash(entry.Pkg)); exists(pkg) && !packageReachesGrammar(pkg) {
			diags = append(diags, fmt.Sprintf("%s is recorded as routing %q through %s but no file there reaches it", entry.Pkg, name, grammarHelper))
		}
	}
	for name := range subcommandRouting {
		if !dispatched[name] {
			diags = append(diags, fmt.Sprintf("the subcommand argument-routing registry names %q, which %s no longer dispatches", name, dispatchFile))
		}
	}
	return uniqueSorted(diags)
}

// dispatchNames reads both dispatch surfaces: the keys of the `commands` map and the case
// labels of the switch in run(). Either one alone leaves half the CLI unexamined.
func dispatchNames(path, body string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, decl := range file.Decls {
		switch node := decl.(type) {
		case *ast.GenDecl:
			names = append(names, mapDispatchNames(node)...)
		case *ast.FuncDecl:
			if node.Name != nil && node.Name.Name == dispatchFunc {
				names = append(names, switchDispatchNames(node)...)
			}
		}
	}
	return names, nil
}

func mapDispatchNames(decl *ast.GenDecl) []string {
	var names []string
	for _, spec := range decl.Specs {
		value, ok := spec.(*ast.ValueSpec)
		if !ok || len(value.Names) != 1 || value.Names[0].Name != dispatchMap {
			continue
		}
		for _, expr := range value.Values {
			literal, ok := expr.(*ast.CompositeLit)
			if !ok {
				continue
			}
			for _, element := range literal.Elts {
				pair, ok := element.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				if name, ok := stringLiteral(pair.Key); ok {
					names = append(names, name)
				}
			}
		}
	}
	return names
}

func switchDispatchNames(fn *ast.FuncDecl) []string {
	var names []string
	ast.Inspect(fn, func(node ast.Node) bool {
		clause, ok := node.(*ast.CaseClause)
		if !ok {
			return true
		}
		for _, expr := range clause.List {
			if name, ok := stringLiteral(expr); ok {
				names = append(names, name)
			}
		}
		return true
	})
	return names
}

func stringLiteral(expr ast.Expr) (string, bool) {
	literal, ok := expr.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return "", false
	}
	value, err := strconv.Unquote(literal.Value)
	if err != nil {
		return "", false
	}
	return value, true
}

// packageReachesGrammar reports whether any non-test Go source directly in dir carries a
// real reference to the grammar helper, read through the AST. Every routed package's
// grammar doc comment names the helper in prose, so a text search would stay green on
// exactly the state this check exists to catch: the call deleted and the comment left
// behind. The entry point and its grammar declaration live in the same package, so a
// directory-level read is enough and no import graph has to be resolved.
func packageReachesGrammar(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		if fileReachesGrammar(filepath.Join(dir, name)) {
			return true
		}
	}
	return false
}

func fileReachesGrammar(path string) bool {
	body := readIfExists(path)
	if body == "" {
		return false
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return false
	}
	found := false
	ast.Inspect(file, func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != grammarSel {
			return true
		}
		if ident, ok := selector.X.(*ast.Ident); ok && ident.Name == grammarPkg {
			found = true
		}
		return true
	})
	return found
}

// TestSubcommandRoutingRegistryBites is the recorded bite proof for
// checkSubcommandRouting (per craft-gate). It runs against a synthetic dispatch file, not
// the repo tree, and walks the three states that matter: every dispatched name registered,
// an unregistered name added to the map surface, and an unregistered name added to the
// switch surface — because a check that read only one surface would leave the other half
// of the CLI unexamined.
func TestSubcommandRoutingRegistryBites(t *testing.T) {
	root := t.TempDir()
	write := func(body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(root, "cmd", "bench"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(dispatchFile)), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// One registered name per surface: "maps" is a routed row and "version" an exempt one,
	// so a clean run proves both dispositions pass rather than only the exempt shortcut.
	// The registry rows for every other real name are absent from this synthetic file, so
	// the check's own "registry names a name no longer dispatched" arm fires for them; the
	// assertions below therefore look for the added name rather than counting diagnostics.
	clean := "package main\n\nvar commands = map[string]int{\n\t\"maps\": 1,\n}\n\nfunc run(args []string) int {\n\tswitch args[0] {\n\tcase \"version\":\n\t\treturn 0\n\t}\n\treturn 2\n}\n"

	write(clean)
	if containsDiagnostic(checkSubcommandRouting(root), "with no entry in the subcommand argument-routing registry") {
		t.Fatalf("registered names alone: want no unregistered-name diagnostic, got %v", checkSubcommandRouting(root))
	}

	write(strings.Replace(clean, "\t\"maps\": 1,\n", "\t\"maps\": 1,\n\t\"newmap\": 1,\n", 1))
	if !containsDiagnostic(checkSubcommandRouting(root), `dispatches "newmap" with no entry`) {
		t.Fatalf("unregistered map key: want a diagnostic naming newmap, got %v", checkSubcommandRouting(root))
	}

	write(strings.Replace(clean, "\tcase \"version\":\n", "\tcase \"version\", \"newcase\":\n", 1))
	if !containsDiagnostic(checkSubcommandRouting(root), `dispatches "newcase" with no entry`) {
		t.Fatalf("unregistered switch label: want a diagnostic naming newcase, got %v", checkSubcommandRouting(root))
	}
}

// TestSubcommandRoutingRoutedClaimBites proves the routed disposition is verified rather
// than merely asserted: a package recorded as routed that no longer reaches the grammar
// helper is reported.
func TestSubcommandRoutingRoutedClaimBites(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(dispatchFile)), []byte("package main\n\nvar commands = map[string]int{\n\t\"maps\": 1,\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	pkg := filepath.Join(root, filepath.FromSlash(subcommandRouting["maps"].Pkg))
	if err := os.MkdirAll(pkg, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(pkg, "maps.go"), []byte("package maps\n\nvar _ = usage.Parse\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if containsDiagnostic(checkSubcommandRouting(root), "but no file there reaches it") {
		t.Fatalf("package reaching the helper: want no routed-claim diagnostic, got %v", checkSubcommandRouting(root))
	}

	if err := os.WriteFile(filepath.Join(pkg, "maps.go"), []byte("package maps\n\nvar _ = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !containsDiagnostic(checkSubcommandRouting(root), "but no file there reaches it") {
		t.Fatalf("package no longer reaching the helper: want a routed-claim diagnostic, got %v", checkSubcommandRouting(root))
	}

	// The state a textual search cannot see: the call deleted, the doc comment that
	// names the helper left behind. Every routed package carries such a comment, so a
	// substring check would pass this exact case vacuously.
	mentionOnly := "package maps\n\n// The declared argument shape usage.Parse enforces lives here.\nconst helper = \"usage.Parse\"\n"
	if err := os.WriteFile(filepath.Join(pkg, "maps.go"), []byte(mentionOnly), 0o644); err != nil {
		t.Fatal(err)
	}
	if !containsDiagnostic(checkSubcommandRouting(root), "but no file there reaches it") {
		t.Fatalf("helper named only in a comment and a literal: want a routed-claim diagnostic, got %v", checkSubcommandRouting(root))
	}
}
