package consumers

import (
	"go/ast"
	"go/types"
	"sort"
)

// viaReference is the classification every row carries today. A later ticket splits call
// and implements out of this one function; the callers do not change when it does.
const viaReference = "reference"

// classify names the edge one use forms. It is deliberately a single small function, so
// the call and implements classes land here and nowhere else.
func classify(*Package, *ast.Ident) string {
	return viaReference
}

// Rows enumerates every reference to target across pkgs, sorted by file, line, then
// column. root makes the file cells repository-relative.
func Rows(pkgs []*Package, target types.Object, root string) []Row {
	var out []Row
	for _, pkg := range pkgs {
		if pkg.Info == nil {
			continue
		}
		// Identity is resolved per package against that package's file set, so a use in a
		// test variant of the target's own package matches the plain declaration.
		want := declKey(pkg.Fset, target)
		targets := map[*ast.Ident]bool{}
		for id, obj := range pkg.Info.Uses {
			if declKey(pkg.Fset, obj) == want {
				targets[id] = true
			}
		}
		if len(targets) == 0 {
			continue
		}
		for _, file := range pkg.Files {
			out = append(out, fileRows(pkg, file, targets, root)...)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.Column < b.Column
	})
	return dedupe(out)
}

// dedupe drops rows that name one position twice. A loader that reports one file in both
// a package and its test variant otherwise doubles every row in it.
func dedupe(rows []Row) []Row {
	out := rows[:0]
	seen := map[Row]bool{}
	for _, r := range rows {
		if seen[r] {
			continue
		}
		seen[r] = true
		out = append(out, r)
	}
	return out
}

// fileRows walks one file once and emits a row for each target identifier it reaches,
// carrying the innermost enclosing named declaration open at that point.
func fileRows(pkg *Package, file *ast.File, targets map[*ast.Ident]bool, root string) []Row {
	var out []Row
	var names []string
	var pushed []bool
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			if pushed[len(pushed)-1] {
				names = names[:len(names)-1]
			}
			pushed = pushed[:len(pushed)-1]
			return false
		}
		name, ok := declName(n, len(names) == 0)
		if ok {
			names = append(names, name)
		}
		pushed = append(pushed, ok)
		if id, isIdent := n.(*ast.Ident); isIdent && targets[id] {
			pos := pkg.Fset.Position(id.Pos())
			enclosing := ""
			if len(names) > 0 {
				enclosing = names[len(names)-1]
			}
			out = append(out, Row{
				File:      relPath(root, pos.Filename),
				Line:      pos.Line,
				Column:    pos.Column,
				Via:       classify(pkg, id),
				Enclosing: enclosing,
			})
		}
		return true
	})
	return out
}

// declName is the one derivation of the enclosing-declaration name. It reports the four
// forms the surface names: a function, a Type.Method, a package-level type, and a
// package-level var or const spec. A node outside those reports false, so the walker's
// innermost open name is the answer for every row. A use with no enclosing declaration at
// all keeps the empty string. fileScope says no named declaration is open yet.
func declName(n ast.Node, fileScope bool) (string, bool) {
	switch d := n.(type) {
	case *ast.FuncDecl:
		if recv := receiverName(d); recv != "" {
			return recv + "." + d.Name.Name, true
		}
		return d.Name.Name, true
	case *ast.TypeSpec:
		if fileScope {
			return d.Name.Name, true
		}
	case *ast.ValueSpec:
		// Only a package-level spec names a row. A local `var x = ...` inside a function
		// is not the declaration an agent is looking for; the function is.
		if fileScope && len(d.Names) > 0 {
			return d.Names[0].Name, true
		}
	}
	return "", false
}

// receiverName is the base type name a method hangs off, with the pointer star and any
// type parameters stripped. It is the empty string for a plain function.
func receiverName(d *ast.FuncDecl) string {
	if d.Recv == nil || len(d.Recv.List) == 0 {
		return ""
	}
	expr := d.Recv.List[0].Type
	for {
		switch t := expr.(type) {
		case *ast.StarExpr:
			expr = t.X
		case *ast.IndexExpr:
			expr = t.X
		case *ast.IndexListExpr:
			expr = t.X
		case *ast.Ident:
			return t.Name
		default:
			return ""
		}
	}
}
