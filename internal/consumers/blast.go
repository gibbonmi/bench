package consumers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"sort"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/toon"
)

// blastFields is the blast schema: the declaration the diff touched, one consumer
// position, and whether that consumer file also sits inside the diff. The base and the
// tip are constant across the table, so neither gets a column.
var blastFields = []string{"changed_symbol", "file", "line", "touched"}

// blastDeletedFields is the deleted-declaration schema. The row carries a base-side
// position because the declaration no longer exists at the tip, and it carries no
// consumer columns: base-side enumeration is a stated non-goal, and a green tip has
// already edited every consumer the deletion had.
var blastDeletedFields = []string{"changed_symbol", "base_file", "base_line"}

// blastRow is one consumer of one touched declaration. Column is carried for the sort
// order only, the way a consumers Row carries it, and is never printed.
type blastRow struct {
	Symbol  string
	File    string
	Line    int
	Column  int
	Touched bool
}

// deletedRow is one declaration the base declared and the tip does not.
type deletedRow struct {
	Symbol   string
	BaseFile string
	BaseLine int
}

// touchedDecl is one tip declaration whose span meets an added line run.
type touchedDecl struct {
	Pkg  *Package
	Obj  types.Object
	Name string
}

// declSite is one package-level declaration as the source names it: its spelling, the
// lines it spans, and its own position. Both the tip walk and the base-side walk read
// this one derivation, so a method spells the same `Type.Method` on both sides.
type declSite struct {
	Name       string
	Start, End int
	Pos        token.Pos
}

// topLevelDecls names every package-level declaration in one file and the lines it spans.
// declNames owns which nodes name a declaration and how each one is spelled, so this walk
// contributes the spans only. A parenthesized group spans per spec, so one edited constant
// in a long block does not touch the whole block.
func topLevelDecls(fset *token.FileSet, file *ast.File) []declSite {
	var out []declSite
	add := func(names []string, node ast.Node) {
		for _, name := range names {
			if name == "_" {
				continue
			}
			out = append(out, declSite{
				Name:  name,
				Start: fset.Position(node.Pos()).Line,
				End:   fset.Position(node.End()).Line,
				Pos:   node.Pos(),
			})
		}
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if names, ok := declNames(d, true); ok {
				add(names, d)
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				node := ast.Node(spec)
				if !d.Lparen.IsValid() {
					node = d
				}
				if names, ok := declNames(spec, true); ok {
					add(names, node)
				}
			}
		}
	}
	return out
}

// lookupDecl resolves a declaration spelling inside one package. Resolve owns the query
// grammar, so the spelling is read there rather than split a second time here, and the
// answer is kept only when the declaration belongs to this package.
func lookupDecl(pkg *Package, name string) types.Object {
	if pkg.Types == nil {
		return nil
	}
	matches, err := Resolve([]*Package{pkg}, name)
	if err != nil {
		return nil
	}
	for _, m := range matches {
		if m.Obj.Pkg() == pkg.Types {
			return m.Obj
		}
	}
	return nil
}

// touchedDecls names every tip declaration whose span meets an added line run in its own
// file. added is keyed by repository-relative path. Identity is declKey, so a package
// loaded both plainly and as a test variant contributes one declaration.
func touchedDecls(pkgs []*Package, root string, added map[string][]lineSpan) []touchedDecl {
	var out []touchedDecl
	seen := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			rel := relPath(root, pkg.Fset.Position(file.Pos()).Filename)
			spans := added[rel]
			if len(spans) == 0 {
				continue
			}
			for _, site := range topLevelDecls(pkg.Fset, file) {
				if !intersects(spans, site.Start, site.End) {
					continue
				}
				obj := lookupDecl(pkg, site.Name)
				if obj == nil {
					continue
				}
				key := declKey(pkg.Fset, obj)
				if seen[key] {
					continue
				}
				seen[key] = true
				out = append(out, touchedDecl{Pkg: pkg, Obj: obj, Name: site.Name})
			}
		}
	}
	return out
}

// qualifiedSpelling is the re-query argument a blast row carries. It comes from Resolve's
// own qualify pass over every declaration of that name, so the `--full` action the
// envelope offers resolves back to this one declaration rather than to a fresh ambiguity.
func qualifiedSpelling(pkgs []*Package, decl touchedDecl) string {
	if matches, err := Resolve(pkgs, decl.Name); err == nil {
		want := declKey(decl.Pkg.Fset, decl.Obj)
		for _, m := range matches {
			if declKey(decl.Pkg.Fset, m.Obj) == want {
				return m.Qualified
			}
		}
	}
	return pathSuffix(decl.Pkg.PkgPath, 1) + "." + decl.Name
}

// blastRows enumerates the consumers of every touched declaration. It is a pure function
// of the typed packages, the touched set, and the diff's file set: the enumeration itself
// is the existing Rows core, so the blast mode adds a marking pass and no second walker.
// A declaration's own definition site is not a consumer of it, so it is dropped.
func blastRows(pkgs []*Package, root string, decls []touchedDecl, changed map[string]bool) []blastRow {
	var out []blastRow
	for _, decl := range decls {
		symbol := qualifiedSpelling(pkgs, decl)
		pos := decl.Pkg.Fset.Position(decl.Obj.Pos())
		self := relPath(root, pos.Filename)
		for _, r := range Rows(pkgs, decl.Obj, root) {
			if r.File == self && r.Line == pos.Line && r.Column == pos.Column {
				continue
			}
			out = append(out, blastRow{
				Symbol: symbol, File: r.File, Line: r.Line, Column: r.Column,
				Touched: changed[r.File],
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		switch {
		case a.Symbol != b.Symbol:
			return a.Symbol < b.Symbol
		case a.File != b.File:
			return a.File < b.File
		case a.Line != b.Line:
			return a.Line < b.Line
		}
		return a.Column < b.Column
	})
	return out
}

// deletedRows names every declaration the base declared inside a removed line run whose
// name the tip's package no longer declares. A package is a directory, so a declaration
// moved between files of one package is not a deletion. baseSources maps a
// repository-relative path to that path's base-side bytes, so this function parses rather
// than reads: the `git show` calls stay at the command's rim.
func deletedRows(pkgs []*Package, root string, hunks []fileHunks, baseSources map[string]string) []deletedRow {
	live := tipDeclNames(pkgs, root)
	var out []deletedRow
	for _, fh := range hunks {
		source, ok := baseSources[fh.BasePath]
		if !ok || len(fh.Removed) == 0 {
			continue
		}
		fset := token.NewFileSet()
		file, err := parseBaseFile(fset, fh.BasePath, source)
		if err != nil {
			continue
		}
		dir := path.Dir(fh.BasePath)
		for _, site := range topLevelDecls(fset, file) {
			if !intersects(fh.Removed, site.Start, site.End) || live[dir+"\x00"+site.Name] {
				continue
			}
			out = append(out, deletedRow{
				Symbol:   path.Base(dir) + "." + site.Name,
				BaseFile: fh.BasePath,
				BaseLine: site.Start,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		switch {
		case a.Symbol != b.Symbol:
			return a.Symbol < b.Symbol
		case a.BaseFile != b.BaseFile:
			return a.BaseFile < b.BaseFile
		}
		return a.BaseLine < b.BaseLine
	})
	return out
}

// parseBaseFile parses one base-side file for its declarations only. The base revision is
// not type-checked: a blast answer never enumerates base-side consumers, so the parse
// needs the declaration names and their spans and nothing else.
func parseBaseFile(fset *token.FileSet, name, source string) (*ast.File, error) {
	return parser.ParseFile(fset, name, source, parser.SkipObjectResolution)
}

// tipDeclNames keys every tip declaration by its package directory and its spelling. The
// directory is the key rather than the import path because the base side is parsed
// without type information and knows only where the file sits.
func tipDeclNames(pkgs []*Package, root string) map[string]bool {
	live := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			rel := relPath(root, pkg.Fset.Position(file.Pos()).Filename)
			dir := path.Dir(rel)
			for _, site := range topLevelDecls(pkg.Fset, file) {
				live[dir+"\x00"+site.Name] = true
			}
		}
	}
	return live
}

// blastResponse renders the whole blast answer: the consumer table, the deleted table
// when it has rows, the meta accounting, the citation, and the per-symbol help envelope.
// It composes envelope, so the blast mode cannot grow its own accounting or lose its
// disclosure.
func blastResponse(source citation, pkgCount, matchCount int, rows []blastRow, deleted []deletedRow, full bool, replay []string) (string, int) {
	rows, files, dropped := representableFiles(rows, func(r blastRow) string { return r.File })
	deleted, _, deletedDropped := representableFiles(deleted, func(r deletedRow) string { return r.BaseFile })
	overCap := !full && len(rows) > rowCap
	var block string
	var err error
	if overCap {
		block, err = toon.TableTyped("consumers_packages", aggregateFields, aggregate(files))
	} else {
		cells := make([][]any, len(rows))
		for i, r := range rows {
			cells[i] = []any{r.Symbol, r.File, r.Line, r.Touched}
		}
		block, err = toon.TableTyped("blast", blastFields, cells)
	}
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	if len(deleted) > 0 {
		cells := make([][]any, len(deleted))
		for i, r := range deleted {
			cells[i] = []any{r.Symbol, r.BaseFile, r.BaseLine}
		}
		deletedBlock, err := toon.TableTyped("blast_deleted", blastDeletedFields, cells)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		block += deletedBlock
	}
	return envelope(source, block, pkgCount, countFiles(files), matchCount, len(rows),
		overCap || dropped || deletedDropped, blastActions(rows, overCap, replay))
}

// blastActions is the blast envelope. An over-cap default offers the one invocation that
// returns every row; otherwise each changed symbol with a consumer outside the diff earns
// one per-symbol `--full` action, because that outside-diff set is what the review walks.
// A blast with no untouched consumer is a terminal read.
func blastActions(rows []blastRow, truncated bool, replay []string) []axi.Action {
	if truncated {
		args := []axi.InvocationArgument{axi.KnownArgument("consumers")}
		for _, arg := range replay {
			args = append(args, axi.KnownArgument(arg))
		}
		return []axi.Action{axi.ExecutableInvocation("emit every consumer row", append(args, axi.KnownArgument("--full"))...)}
	}
	seen := map[string]bool{}
	var symbols []string
	for _, r := range rows {
		if r.Touched || seen[r.Symbol] {
			continue
		}
		seen[r.Symbol] = true
		symbols = append(symbols, r.Symbol)
	}
	sort.Strings(symbols)
	actions := make([]axi.Action, 0, len(symbols))
	for _, symbol := range symbols {
		actions = append(actions, axi.ExecutableInvocation("walk the consumers outside the diff",
			axi.KnownArgument("consumers"), axi.KnownArgument(symbol), axi.KnownArgument("--full")))
	}
	return actions
}
