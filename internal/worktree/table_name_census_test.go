package worktree

// This file owns the table-name census. A table name is the value of a package string
// constant that the non-test source passes to a toon renderer. The census reports a
// renderer call that names its table any other way. In the test files, it reports a
// string literal that equals a table name as a block argument, and a string literal that
// begins with a table name and "[", which is the rendered header.

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// toonImport is the import path of the package that renders the tables.
const toonImport = "github.com/gibbonmi/bench/internal/toon"

// tableBlockArgs maps each call that takes a table name to the index of that argument. A
// toon renderer matches by its import path and name. A test helper or a document method
// matches by its name.
var tableBlockArgs = map[string]int{
	kindKey(toonImport, "Table"): 0, kindKey(toonImport, "TableTyped"): 0,
	"mustRows": 1, "readVerbRows": 1, "Rows": 0,
}

// tableFinding selects one of the two census messages.
type tableFinding int

const (
	tableNamedInSource tableFinding = iota
	tableSpelledInTest
)

// tableNameReport renders one census message. A source call names the renderer and the
// spelling of its table argument. A test literal names its spelling and the constant that
// names its table.
func tableNameReport(file string, line int, finding tableFinding, spelled, name string) string {
	where := fmt.Sprintf("%s:%d: ", file, line)
	if finding == tableSpelledInTest {
		return where + spelled + " spells a table name; read the constant " + name
	}
	return where + name + " names its table with " + spelled + ", not a package constant"
}

// stringLiteral returns the value of expr when expr is a string literal.
func stringLiteral(expr ast.Expr) (string, bool) {
	literal, ok := expr.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return "", false
	}
	text, err := strconv.Unquote(literal.Value)
	return text, err == nil
}

// tableCallKey keys call for tableBlockArgs: a package call by its import path and name,
// a method call by its selected name, and a bare call by its name.
func tableCallKey(call *ast.CallExpr, imports map[string]string) string {
	if _, key := qualifiedName(call.Fun, imports); key != "" {
		return key
	}
	if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
		return selector.Sel.Name
	}
	return calleeName(call)
}

// tableNames maps each table name in the non-test source of dir to the constant that
// names it. The reports name each renderer call that passes anything other than a package
// string constant, sorted. A source with no renderer call is an error, because a census
// with no table to grade passes vacuously.
func tableNames(dir string) (map[string]string, []string, error) {
	files, fset, _, err := parseGoFiles(dir, isSourceFile)
	if err != nil {
		return nil, nil, err
	}
	constants := map[string]string{}
	for _, file := range files {
		for _, node := range file.Decls {
			if decl, ok := node.(*ast.GenDecl); ok && decl.Tok == token.CONST {
				for _, spec := range decl.Specs {
					value := spec.(*ast.ValueSpec)
					for i, ident := range value.Names {
						if i >= len(value.Values) {
							break
						}
						if text, ok := stringLiteral(value.Values[i]); ok {
							constants[ident.Name] = text
						}
					}
				}
			}
		}
	}
	names, calls := map[string]string{}, 0
	var reports []string
	for _, file := range files {
		imports := fileImportNames(file)
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			spelled, key := qualifiedName(call.Fun, imports)
			index, renderer := tableBlockArgs[key]
			if importPath, _ := splitKind(key); !renderer || importPath != toonImport || index >= len(call.Args) {
				return true
			}
			calls++
			argument := call.Args[index]
			if ident, ok := argument.(*ast.Ident); ok && constants[ident.Name] != "" {
				names[constants[ident.Name]] = ident.Name
				return true
			}
			position := fset.Position(argument.Pos())
			reports = append(reports, tableNameReport(filepath.Base(position.Filename), position.Line, tableNamedInSource, types.ExprString(argument), spelled))
			return true
		})
	}
	if calls == 0 {
		return nil, nil, fmt.Errorf("the non-test source in %s holds no toon.Table or toon.TableTyped call", dir)
	}
	slices.Sort(reports)
	return names, reports, nil
}

// tableNameLiteralCensus reports each string literal in the test files of dir that equals
// a table name as a block argument, or that begins with a table name and "[". The reports
// are sorted. A source that names no table with a constant is an error.
func tableNameLiteralCensus(dir string) ([]string, error) {
	names, _, err := tableNames(dir)
	if err == nil && len(names) == 0 {
		err = fmt.Errorf("the non-test source in %s names no table with a package constant", dir)
	}
	if err != nil {
		return nil, err
	}
	files, fset, _, err := parseTestFiles(dir)
	if err != nil {
		return nil, err
	}
	var reports []string
	for _, file := range files {
		imports := fileImportNames(file)
		// The walk visits a call before its arguments, so a block argument is marked first.
		blockArgs := map[ast.Node]bool{}
		ast.Inspect(file, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if index, ok := tableBlockArgs[tableCallKey(call, imports)]; ok && index < len(call.Args) {
					blockArgs[call.Args[index]] = true
				}
			}
			literal, ok := node.(*ast.BasicLit)
			if !ok {
				return true
			}
			text, _ := stringLiteral(literal)
			for name, constant := range names {
				if (blockArgs[literal] && text == name) || strings.HasPrefix(text, name+"[") {
					position := fset.Position(literal.Pos())
					reports = append(reports, tableNameReport(filepath.Base(position.Filename), position.Line, tableSpelledInTest, literal.Value, constant))
				}
			}
			return true
		})
	}
	slices.Sort(reports)
	return reports, nil
}

// syntheticTableConstant names the one table of a synthetic package.
const syntheticTableConstant = "cleanupTable"

// syntheticRenderFile is the non-test file of a synthetic package. It declares
// syntheticTableConstant with the live cleanup table name, and its one renderer call, on
// line 7, passes argument.
func syntheticRenderFile(argument string) string {
	return "package worktree\n\nimport " + strconv.Quote(toonImport) + "\n\nconst " + syntheticTableConstant + " = " + strconv.Quote(cleanupTable) +
		"\n\nfunc render() { _, _ = toon.Table(" + argument + ", nil, nil) }\n"
}

// wantTableCensus runs the literal census over a synthetic package whose render file
// passes syntheticTableConstant, and whose table_test.go holds body from line 11. The
// reports must equal want.
func wantTableCensus(t *testing.T, body string, want ...string) {
	t.Helper()
	dir := plantTestFiles(t, map[string]string{
		"render.go":     syntheticRenderFile(syntheticTableConstant),
		"table_test.go": "package worktree\n\nimport (\n\t\"path/filepath\"\n\t\"testing\"\n\n\t" + strconv.Quote(toonImport) + "\n)\n\nfunc TestRows(t *testing.T) {\n" + body + "}\n",
	})
	reports, err := tableNameLiteralCensus(dir)
	if err != nil {
		t.Fatalf("table-name census: %v", err)
	}
	if !slices.Equal(reports, want) {
		t.Fatalf("table-name census = %q, want exactly %q", reports, want)
	}
}

// TestTableNamesAreProductionConstants proves a renderer call that passes a literal draws
// a report, a source with no renderer call is an error, and each renderer call in the
// package passes a constant.
func TestTableNamesAreProductionConstants(t *testing.T) {
	t.Parallel()
	if names, _, err := tableNames(plantTestFiles(t, map[string]string{"render.go": "package worktree\n"})); err == nil {
		t.Fatalf("table names of a source with no renderer call = %q, want an error", names)
	}
	literal := strconv.Quote(cleanupTable)
	_, reports, err := tableNames(plantTestFiles(t, map[string]string{"render.go": syntheticRenderFile(literal)}))
	if want := []string{tableNameReport("render.go", 7, tableNamedInSource, literal, "toon.Table")}; err != nil || !slices.Equal(reports, want) {
		t.Fatalf("table names of a literal call = %q, %v, want exactly %q", reports, err, want)
	}
	if _, reports, err = tableNames("."); err != nil || len(reports) != 0 {
		t.Fatalf("the package names %d tables without a constant (%v):\n%s", len(reports), err, strings.Join(reports, "\n"))
	}
}

// TestTableNameLiteralCensusReportsABlockArgument proves each block call reports a table
// name literal, and a path segment that equals a table name is not a block argument.
func TestTableNameLiteralCensusReportsABlockArgument(t *testing.T) {
	t.Parallel()
	literal := strconv.Quote(cleanupTable)
	for _, call := range []string{"r.mustRows(t, %s)", "readVerbRows(out, %s)", "document.Rows(%s)", "toon.Table(%s, nil, nil)", "toon.TableTyped(%s, nil, nil)"} {
		t.Run(call, func(t *testing.T) {
			t.Parallel()
			wantTableCensus(t, "\t_ = "+fmt.Sprintf(call, literal)+"\n\t_ = filepath.Join(\"home\", "+literal+")\n",
				tableNameReport("table_test.go", 11, tableSpelledInTest, literal, syntheticTableConstant))
		})
	}
}

// TestTableNameLiteralCensusReportsARenderedHeader proves an interpreted or a raw literal
// that begins with a table name and "[" draws a report, and a whole name outside a block
// argument does not.
func TestTableNameLiteralCensusReportsARenderedHeader(t *testing.T) {
	t.Parallel()
	header, raw := strconv.Quote(cleanupTable+"[1]{target}:\n"), "`"+cleanupTable+"[0]`"
	wantTableCensus(t, "\t_ = strings.Contains(out, "+header+")\n\t_ = strings.HasPrefix(out, "+raw+")\n\t_ = strings.Contains(out, "+strconv.Quote(cleanupTable)+")\n",
		tableNameReport("table_test.go", 11, tableSpelledInTest, header, syntheticTableConstant),
		tableNameReport("table_test.go", 12, tableSpelledInTest, raw, syntheticTableConstant))
}

// TestTableNameLiteralCensusOnTheLiveTree proves each test reads a table name from its
// constant.
func TestTableNameLiteralCensusOnTheLiveTree(t *testing.T) {
	t.Parallel()
	reports, err := tableNameLiteralCensus(".")
	if err != nil {
		t.Fatalf("table-name census the package: %v", err)
	}
	if len(reports) != 0 {
		t.Fatalf("the table-name census reports %d literals:\n%s", len(reports), strings.Join(reports, "\n"))
	}
}
