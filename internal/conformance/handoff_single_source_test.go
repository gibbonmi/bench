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
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/status"
)

// The Shape text lives in the binary. shapeSourceFile is the file that holds it.
// status.HandoffFile is the one document allowed to carry a rendered copy of it.
const shapeSourceFile = "internal/handoff/text.go"

// checkHandoffShape grades the Shape text's single-source claim two ways. The tracked
// artifact must still be byte-equal to the constant the command renders from. No other
// tracked file may restate it. Untracked scratch and ignored build output are out of
// scope, because a copy only drifts once it is something the repo carries.
//
// The constant compared against is always the kit's own. The check stays silent unless
// the graded root is the kit, because another tree cannot answer for a contract whose
// source it does not contain.
func checkHandoffShape(root string) []string {
	if !exists(filepath.Join(root, filepath.FromSlash(shapeSourceFile))) {
		return nil
	}

	var diags []string
	tracked := trackedPaths(root)
	// The artifact is graded only while the repo carries it: an untracked handoff is
	// local scratch, and a copy only drifts once it is something the repo carries.
	if !contains(tracked, status.HandoffFile) {
		return scanShapeCopies(root, tracked)
	}
	body, found := shapeSectionBody(readIfExists(filepath.Join(root, status.HandoffFile)))
	switch {
	case !found:
		diags = append(diags, fmt.Sprintf("%s carries no %q section; run bench handoff to regenerate it", status.HandoffFile, handoff.ShapeHeading))
	case body != strings.Trim(handoff.ShapeSection, "\n"):
		diags = append(diags, fmt.Sprintf("%s has a %q section that no longer matches the text bench handoff renders, so the artifact has become a second source; run bench handoff to regenerate it", status.HandoffFile, handoff.ShapeHeading))
	}

	diags = append(diags, scanShapeCopies(root, tracked)...)
	return uniqueSorted(diags)
}

// scanShapeCopies reports every tracked file, other than the source and the artifact,
// that restates the Shape text.
func scanShapeCopies(root string, tracked []string) []string {
	var diags []string
	needle := shapeSentence()
	for _, rel := range tracked {
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

// contains reports whether list carries exactly value.
func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// shapeSentence is the Shape text's opening sentence with its wrapping collapsed. It
// reads out of the constant, so this file carries no copy of the text it forbids copying.
// Collapsing whitespace is what makes a re-wrapped restatement match the original.
func shapeSentence() string {
	sentence := collapseSpace(handoff.ShapeSection)
	if end := strings.Index(sentence, ". "); end >= 0 {
		sentence = sentence[:end+1]
	}
	return sentence
}

// shapeSectionBody returns the document's Shape body: everything below the heading, with
// surrounding blank lines trimmed. It locates the heading by the writer's own constant.
// It runs to EOF because the writer places Shape last. This function restates no part of
// the command's splitting rule. A second derivation of where a section ends is exactly
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

// TestHandoffShapeSingleSourcedBites is the recorded bite proof for checkHandoffShape. It
// runs against a synthetic repository rather than the kit tree. It walks the three states
// that matter. The first two states are a derived artifact with no other copy, and a
// second tracked file carrying the text. The third state is an artifact whose Shape body
// has drifted from the constant.
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
	// The source file is a presence sentinel. The text the check grades against comes from
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

	// An untracked artifact is out of scope: a fresh worktree without the local handoff
	// must stay green, drifted body and all.
	runGit(t, root, "rm", "--cached", "-q", status.HandoffFile)
	if diags := checkHandoffShape(root); len(diags) != 0 {
		t.Fatalf("untracked artifact: want no diagnostics, got %v", diags)
	}
}

// TestHandoffShapeNamesTheSectionGrammar grades the Shape text against the grammar
// the leaf package renders. checkHandoffShape holds the artifact to the constant and
// says nothing about what the constant claims, so a Shape that still described the
// one-phase document would pass it while the file grew a section per assignment.
//
// The spellings come from internal/handoffdoc, never from a second copy here.
func TestHandoffShapeNamesTheSectionGrammar(t *testing.T) {
	for _, token := range []string{
		handoffdoc.MainHeading,
		handoffdoc.RequestHeadingPrefix,
		handoffdoc.StateHeading,
		handoffdoc.LabelRequestToken,
		handoffdoc.LabelWorktreeTip,
		handoffdoc.LabelRecordedBase,
		handoffdoc.LabelSpecStatus,
		handoffdoc.LabelNextCommand,
	} {
		if !strings.Contains(handoff.ShapeSection, strings.TrimSpace(token)) {
			t.Errorf("the Shape text names no %q, so it does not describe the section grammar", token)
		}
	}
}

// handoffPkgDir holds the consumer package. prefixTablePkgDir and prefixTableFile hold the
// harness record, which owns the phase invocation forms. prefixTableVar holds its rows
// variable, and prefixTableKey holds the row field the forms sit behind.
const (
	handoffPkgDir     = "internal/handoff"
	prefixTablePkgDir = "internal/harnesses"
	prefixTableFile   = "harnesses.go"
	prefixTableVar    = "Rows"
	prefixTableKey    = "PhaseForm"
)

// checkHarnessPrefix reports any string literal in the handoff package that writes a
// phase invocation form the harness record already owns. Examples include a trailing
// replacement with a hardcoded target, an inline conditional, or a second table. Each is
// a producer the table cannot see. A harness added as a row would leave any of them
// behind.
//
// The literals are read through the AST for the reason packageReachesGrammar states. The
// package's doc comments legitimately discuss these forms in prose. A substring scan
// would either fire on that prose or force it to be mangled. The forbidden forms are the
// table's own values, so a new row is covered the moment it lands. This check restates
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
	forms := prefixTable(tableFile)
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

// prefixTable returns the phase invocation forms the record declares. The rows are struct
// literals nested inside a slice literal, so the reader descends one level and collects the
// string literal behind each prefixTableKey field. A row with an empty form contributes no
// forbidden string, because an empty string matches every literal.
func prefixTable(file *ast.File) []string {
	var forms []string
	if file == nil {
		return nil
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
					row, ok := element.(*ast.CompositeLit)
					if !ok {
						continue
					}
					forms = append(forms, rowPhaseForms(row)...)
				}
			}
		}
	}
	return uniqueSorted(forms)
}

// rowPhaseForms returns the non-empty forms one row's prefixTableKey fields hold.
func rowPhaseForms(row *ast.CompositeLit) []string {
	var forms []string
	for _, field := range row.Elts {
		pair, ok := field.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := pair.Key.(*ast.Ident)
		if !ok || key.Name != prefixTableKey {
			continue
		}
		form, ok := stringLiteral(pair.Value)
		if !ok || form == "" {
			continue
		}
		forms = append(forms, form)
	}
	return forms
}

// TestHarnessPrefixSingleSourcedBites is the recorded bite proof for checkHarnessPrefix.
// It runs against a synthetic record whose rows carry invented forms. The forms come from
// the record, so a fixture needs no copy of the real ones and cannot go stale when a
// harness is added. The fixture takes the record's nested row shape, so the test proves the
// reader against the shape the live tree carries, and it holds a formless row to prove an
// empty form forbids nothing. It walks four states: a table alone, and the same forms named
// only in prose. The other two states are an inline literal beside the table, and a second
// table-shaped literal in another file.
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
	table := "package harnesses\n\nvar " + prefixTableVar + " = []Row{\n" +
		"\t{Harness: \"claude\", " + prefixTableKey + ": \"/synth-\"},\n" +
		"\t{Harness: \"codex\", " + prefixTableKey + ": \"$synth-\"},\n" +
		"\t{Harness: \"opencode\", " + prefixTableKey + ": \"\"},\n}\n"

	if !containsDiagnostic(checkHarnessPrefix(root), prefixTablePkgDir+"/"+prefixTableFile+" declares no "+prefixTableVar+" table") {
		t.Fatalf("missing route owner: want a diagnostic, got %v", checkHarnessPrefix(root))
	}

	write(tableDir, prefixTableFile, table)
	if diags := checkHarnessPrefix(root); len(diags) != 0 {
		t.Fatalf("table alone: want no diagnostics, got %v", diags)
	}

	// A substring scan cannot tell this state from a violation: the forms discussed in a doc
	// comment. Every file in this package legitimately does that.
	write(handoffDir, "sections.go", "package handoff\n\n// Phase invocations render as /synth- or $synth- depending on the harness.\nfunc split() {}\n")
	if diags := checkHarnessPrefix(root); len(diags) != 0 {
		t.Fatalf("forms named only in prose: want no diagnostics, got %v", diags)
	}

	write(handoffDir, "sections.go", "package handoff\n\nfunc translate(name string) string { return \"$synth-\" + name }\n")
	if !containsDiagnostic(checkHarnessPrefix(root), `writes the phase-invocation form "$synth-"`) {
		t.Fatalf("inline literal beside the table: want a diagnostic, got %v", checkHarnessPrefix(root))
	}

	write(handoffDir, "sections.go", "package handoff\n\nvar legacyPrefix = map[string]string{\n\t\"claude\": \"/synth-\",\n}\n")
	if !containsDiagnostic(checkHarnessPrefix(root), `sections.go writes the phase-invocation form "/synth-"`) {
		t.Fatalf("second table-shaped literal: want a diagnostic naming the file, got %v", checkHarnessPrefix(root))
	}

	// The formless row forbids nothing. An empty form would otherwise match every literal.
	write(handoffDir, "sections.go", "package handoff\n\nfunc name() string { return \"handoff\" }\n")
	if diags := checkHarnessPrefix(root); len(diags) != 0 {
		t.Fatalf("formless row: want no diagnostics, got %v", diags)
	}
}
