package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/handoff"
	"github.com/gibbonmi/bench/internal/status"
)

// The Shape text lives in the binary; shapeSourceFile is the file that holds it and
// status.HandoffFile is the one document allowed to carry a rendered copy of it.
const shapeSourceFile = "internal/handoff/text.go"

// checkHandoffShape grades the Shape text's single-source claim two ways: the tracked
// artifact must still be byte-equal to the constant the command renders from, and no other
// tracked file may restate it. Untracked scratch and ignored build output are out of scope —
// a copy only drifts once it is something the repo carries.
//
// The constant compared against is always the kit's own, which is why the check stays silent
// unless the graded root is the kit: another tree cannot answer for a contract whose source
// it does not contain.
func checkHandoffShape(root string) []string {
	if !exists(filepath.Join(root, filepath.FromSlash(shapeSourceFile))) {
		return nil
	}

	var diags []string
	body, found := shapeSectionBody(readIfExists(filepath.Join(root, status.HandoffFile)))
	switch {
	case !found:
		diags = append(diags, fmt.Sprintf("%s carries no %q section; run bench handoff to regenerate it", status.HandoffFile, handoff.ShapeHeading))
	case body != strings.Trim(handoff.ShapeSection, "\n"):
		diags = append(diags, fmt.Sprintf("%s has a %q section that no longer matches the text bench handoff renders, so the artifact has become a second source; run bench handoff to regenerate it", status.HandoffFile, handoff.ShapeHeading))
	}

	needle := shapeSentence()
	for _, rel := range trackedPaths(root) {
		if rel == shapeSourceFile || rel == status.HandoffFile {
			continue
		}
		text := readIfExists(filepath.Join(root, filepath.FromSlash(rel)))
		if text == "" || strings.IndexByte(text, 0) >= 0 {
			continue
		}
		if strings.Contains(collapseSpace(text), needle) {
			diags = append(diags, fmt.Sprintf("%s restates the handoff's Shape text, which %s owns and %s alone derives; delete the copy", rel, shapeSourceFile, status.HandoffFile))
		}
	}
	return uniqueSorted(diags)
}

// shapeSentence is the Shape text's opening sentence with its wrapping collapsed, read out
// of the constant so this file carries no copy of the text it forbids copying. Collapsing
// whitespace is what makes a re-wrapped restatement match the original.
func shapeSentence() string {
	sentence := collapseSpace(handoff.ShapeSection)
	if end := strings.Index(sentence, ". "); end >= 0 {
		sentence = sentence[:end+1]
	}
	return sentence
}

// shapeSectionBody returns the document's Shape body: everything below the heading, with
// surrounding blank lines trimmed. It locates the heading by the writer's own constant and
// runs to EOF because the writer places Shape last, so this restates no part of the
// command's splitting rule — a second derivation of "where does a section end" is exactly
// what would let the check and the emitter drift.
func shapeSectionBody(doc string) (string, bool) {
	lines := strings.Split(doc, "\n")
	for i, line := range lines {
		if strings.TrimRight(line, " \t\r") == handoff.ShapeHeading {
			return strings.Trim(strings.Join(lines[i+1:], "\n"), "\n"), true
		}
	}
	return "", false
}

func trackedPaths(root string) []string {
	out, err := exec.Command("git", "-C", root, "ls-files", "-z").Output()
	if err != nil {
		return nil
	}
	listing := strings.TrimRight(string(out), "\x00")
	if listing == "" {
		return nil
	}
	return strings.Split(listing, "\x00")
}

// TestHandoffShapeSingleSourcedBites is the recorded bite proof for checkHandoffShape (per
// craft-gate). It runs against a synthetic repository rather than the kit tree and walks the
// three states that matter: a derived artifact and no other copy, a second tracked file
// carrying the text, and an artifact whose Shape body has drifted from the constant.
func TestHandoffShapeSingleSourcedBites(t *testing.T) {
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
		runGit(t, root, "add", "-A")
	}
	runGit(t, root, "init")
	// The source file is a presence sentinel: the text the check grades against comes from
	// the imported constant, never from this tree.
	write(shapeSourceFile, "package handoff\n")
	write(status.HandoffFile, "# Session handoff\n\n"+handoff.ShapeHeading+"\n\n"+handoff.ShapeSection)

	if diags := checkHandoffShape(root); len(diags) != 0 {
		t.Fatalf("derived artifact and no other copy: want no diagnostics, got %v", diags)
	}

	write("docs/pinned-shape.md", "Notes for a cold session.\n\n"+handoff.ShapeSection)
	if !containsDiagnostic(checkHandoffShape(root), "docs/pinned-shape.md restates the handoff's Shape text") {
		t.Fatalf("second tracked copy: want a diagnostic naming the copy, got %v", checkHandoffShape(root))
	}
	if err := os.Remove(filepath.Join(root, "docs", "pinned-shape.md")); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "-A")

	write(status.HandoffFile, "# Session handoff\n\n"+handoff.ShapeHeading+"\n\nRewritten whenever somebody feels like it.\n")
	if !containsDiagnostic(checkHandoffShape(root), "no longer matches the text bench handoff renders") {
		t.Fatalf("drifted Shape body: want a drift diagnostic, got %v", checkHandoffShape(root))
	}
}

// handoffPkgDir holds the consumer package, prefixTablePkgDir/File hold the route owner's
// harness table, and prefixTableVar is its name.
const (
	handoffPkgDir     = "internal/handoff"
	prefixTablePkgDir = "internal/status"
	prefixTableFile   = "route.go"
	prefixTableVar    = "harnessPrefix"
)

// checkHarnessPrefix reports any string literal in the handoff package that writes a phase
// invocation form the status route owner already owns — a trailing replacement with a
// hardcoded target, an inline conditional, a second table. Each is a producer the table
// cannot see, so a harness added as a row would leave it behind.
//
// The literals are read through the AST for the reason packageReachesGrammar states: the
// package's doc comments legitimately discuss these forms in prose, and a substring scan
// would either fire on that prose or force it to be mangled. The forbidden forms are the
// table's own values, so a new row is covered the moment it lands and this check restates
// nothing the table already says.
func checkHarnessPrefix(root string) []string {
	fset := token.NewFileSet()
	tableDir := filepath.Join(root, filepath.FromSlash(prefixTablePkgDir))
	if !exists(tableDir) {
		return nil
	}
	tablePath := filepath.Join(tableDir, prefixTableFile)
	tableSource := readIfExists(tablePath)
	if tableSource == "" {
		return []string{fmt.Sprintf("%s/%s declares no %s table, so no single owner of the phase-invocation forms remains", prefixTablePkgDir, prefixTableFile, prefixTableVar)}
	}
	tableFile, err := parser.ParseFile(fset, tablePath, tableSource, 0)
	if err != nil {
		return []string{fmt.Sprintf("%s/%s cannot be parsed for harness prefix literals: %v", prefixTablePkgDir, prefixTableFile, err)}
	}
	forms, _ := prefixTable(tableFile)
	if len(forms) == 0 {
		return []string{fmt.Sprintf("%s/%s declares no %s table, so no single owner of the phase-invocation forms remains", prefixTablePkgDir, prefixTableFile, prefixTableVar)}
	}

	dir := filepath.Join(root, filepath.FromSlash(handoffPkgDir))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	parsed := map[string]*ast.File{}
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(dir, name)
		file, err := parser.ParseFile(fset, path, readIfExists(path), 0)
		if err != nil {
			return []string{fmt.Sprintf("%s/%s cannot be parsed for harness prefix literals: %v", handoffPkgDir, name, err)}
		}
		parsed[name] = file
		names = append(names, name)
	}

	var diags []string
	for _, name := range names {
		ast.Inspect(parsed[name], func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if !ok {
				return true
			}
			value, ok := stringLiteral(literal)
			if !ok {
				return true
			}
			for _, form := range forms {
				if strings.Contains(value, form) {
					diags = append(diags, fmt.Sprintf("%s/%s writes the phase-invocation form %q in the literal %q; only the %s table in %s/%s may produce it", handoffPkgDir, name, form, value, prefixTableVar, prefixTablePkgDir, prefixTableFile))
				}
			}
			return true
		})
	}
	return uniqueSorted(diags)
}

// prefixTable returns the harness table's values and the positions of the literals holding
// them, so the table is both what the rest of the package is measured against and the only
// place its own values are exempt.
func prefixTable(file *ast.File) ([]string, map[token.Pos]bool) {
	var forms []string
	owned := map[token.Pos]bool{}
	if file == nil {
		return nil, owned
	}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok || len(value.Names) != 1 || value.Names[0].Name != prefixTableVar {
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
					form, ok := stringLiteral(pair.Value)
					if !ok || form == "" {
						continue
					}
					forms = append(forms, form)
					owned[pair.Value.Pos()] = true
				}
			}
		}
	}
	return uniqueSorted(forms), owned
}

// TestHarnessPrefixSingleSourcedBites is the recorded bite proof for checkHarnessPrefix (per
// craft-gate). It runs against a synthetic package whose table carries invented forms, which
// is the whole mechanism under proof: the forms come from the table, so a fixture needs no
// copy of the real ones and cannot go stale when a harness is added. It walks four states —
// a table alone, the same forms named only in prose, an inline literal beside the table, and
// a second table-shaped literal in another file.
func TestHarnessPrefixSingleSourcedBites(t *testing.T) {
	root := t.TempDir()
	handoffDir := filepath.Join(root, filepath.FromSlash(handoffPkgDir))
	tableDir := filepath.Join(root, filepath.FromSlash(prefixTablePkgDir))
	for _, dir := range []string{handoffDir, tableDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(dir, name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	table := "package status\n\nvar " + prefixTableVar + " = map[string]string{\n\t\"claude\": \"/synth-\",\n\t\"codex\":  \"$synth-\",\n}\n"

	if !containsDiagnostic(checkHarnessPrefix(root), prefixTablePkgDir+"/"+prefixTableFile+" declares no "+prefixTableVar+" table") {
		t.Fatalf("missing route owner: want a diagnostic, got %v", checkHarnessPrefix(root))
	}

	write(tableDir, prefixTableFile, table)
	if diags := checkHarnessPrefix(root); len(diags) != 0 {
		t.Fatalf("table alone: want no diagnostics, got %v", diags)
	}

	// The state a substring scan cannot tell from a violation: the forms discussed in a doc
	// comment, which every file in this package legitimately does.
	write(handoffDir, "sections.go", "package handoff\n\n// Phase invocations render as /synth- or $synth- depending on the harness.\nfunc split() {}\n")
	if diags := checkHarnessPrefix(root); len(diags) != 0 {
		t.Fatalf("forms named only in prose: want no diagnostics, got %v", diags)
	}

	write(handoffDir, "sections.go", "package handoff\n\nfunc translate(name string) string { return \"$synth-\" + name }\n")
	if !containsDiagnostic(checkHarnessPrefix(root), `writes the phase-invocation form "$synth-"`) {
		t.Fatalf("inline literal beside the table: want a diagnostic, got %v", checkHarnessPrefix(root))
	}

	write(handoffDir, "sections.go", "package handoff\n\nvar legacyPrefix = map[string]string{\n\t\"claude\": \"/synth-\",\n}\n")
	if !containsDiagnostic(checkHarnessPrefix(root), "sections.go writes the phase-invocation form") {
		t.Fatalf("second table-shaped literal: want a diagnostic naming the file, got %v", checkHarnessPrefix(root))
	}
}
