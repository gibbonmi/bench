package consumers

import (
	"go/ast"
	"go/token"
	"go/types"
	"maps"
	pathpkg "path"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// The three values the via column carries. reference is the default: an appearance of
// the declaration that is neither a static call nor a method-set satisfaction.
const (
	viaReference  = "reference"
	viaCall       = "call"
	viaImplements = "implements"
)

// ViaCall is the via value of a static call, for a caller outside this package that keeps
// only the call rows.
const ViaCall = viaCall

// classify names the edge one use forms. It is deliberately a single small function, so
// the call and implements classes land here and nowhere else. callee says the walker
// reached this identifier as the static callee of a call expression; the type check then
// separates a real call from a conversion spelled the same way.
func classify(pkg *Package, id *ast.Ident, callee bool) string {
	if !callee || pkg.Info == nil {
		return viaReference
	}
	if _, isFunc := pkg.Info.Uses[id].(*types.Func); isFunc {
		return viaCall
	}
	return viaReference
}

// staticCallee is the identifier a call expression names as its callee, or nil when the
// callee is computed rather than named. It strips the parentheses and the explicit type
// arguments a generic call carries, because neither changes which declaration is called.
func staticCallee(call *ast.CallExpr) *ast.Ident {
	expr := call.Fun
	for {
		switch e := expr.(type) {
		case *ast.ParenExpr:
			expr = e.X
		case *ast.IndexExpr:
			expr = e.X
		case *ast.IndexListExpr:
			expr = e.X
		case *ast.SelectorExpr:
			return e.Sel
		case *ast.Ident:
			return e
		default:
			return nil
		}
	}
}

// Rows enumerates every reference to target across pkgs, sorted by file, line, then
// column. root makes the file cells repository-relative.
func Rows(pkgs []*Package, target types.Object, root string) []Row {
	return sortRows(append(useRows(pkgs, target, root, false), implementsRows(pkgs, target, root)...))
}

// WriteRows is every static write of target across pkgs, in the order Rows uses: an
// assignment or op-assignment target, an increment or decrement operand, and a field key of
// a composite literal. A read is no write.
func WriteRows(pkgs []*Package, target types.Object, root string) []Row {
	return sortRows(useRows(pkgs, target, root, true))
}

// UnkeyedWrites is every unkeyed composite literal across pkgs whose struct type declares the
// field target, in the order Rows uses. Such a literal writes the field by position and names
// no field, so WriteRows cannot list it.
func UnkeyedWrites(pkgs []*Package, target types.Object, root string) []Row {
	var out []Row
	for _, pkg := range pkgs {
		if pkg.Info == nil {
			continue
		}
		want := declKey(pkg.Fset, target)
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.CompositeLit)
				if !ok || len(lit.Elts) == 0 || pkg.Info.TypeOf(lit) == nil {
					return true
				}
				if _, keyed := lit.Elts[0].(*ast.KeyValueExpr); keyed {
					return true
				}
				if fields, isStruct := types.Unalias(pkg.Info.TypeOf(lit)).Underlying().(*types.Struct); isStruct {
					for i := range fields.NumFields() {
						if declKey(pkg.Fset, fields.Field(i)) == want {
							pos := pkg.Fset.Position(lit.Pos())
							out = append(out, Row{File: relPath(root, pos.Filename), Line: pos.Line, Column: pos.Column})
							break
						}
					}
				}
				return true
			})
		}
	}
	return sortRows(out)
}

// useRows is every use of target across pkgs, or only every static write of it.
func useRows(pkgs []*Package, target types.Object, root string, writesOnly bool) []Row {
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
			out = append(out, fileRows(pkg, file, targets, root, writesOnly)...)
		}
	}
	return out
}

// sortRows orders rows by file, line, then column, and drops repeated positions.
func sortRows(out []Row) []Row {
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
// carrying the innermost enclosing named declaration open at that point. writesOnly keeps
// only the identifiers that a statement or a composite literal writes.
func fileRows(pkg *Package, file *ast.File, targets map[*ast.Ident]bool, root string, writesOnly bool) []Row {
	var out []Row
	var names []string
	var pushed []bool
	// Inspect visits a call expression before the identifier inside it, so the callee set
	// is complete for every identifier by the time the walk reaches it. This is the one
	// derivation of call position; classify reads it and never re-walks.
	callees := map[*ast.Ident]bool{}
	// The write set fills the same way: a statement or a literal precedes the identifiers it
	// writes.
	writes := map[*ast.Ident]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			if pushed[len(pushed)-1] {
				names = names[:len(names)-1]
			}
			pushed = pushed[:len(pushed)-1]
			return false
		}
		declared, ok := declNames(n, len(names) == 0)
		if ok {
			names = append(names, declared[0])
		}
		pushed = append(pushed, ok)
		if call, isCall := n.(*ast.CallExpr); isCall {
			if id := staticCallee(call); id != nil {
				callees[id] = true
			}
		}
		markWrites(pkg, n, writes)
		if id, isIdent := n.(*ast.Ident); isIdent && targets[id] && (!writesOnly || writes[id]) {
			pos := pkg.Fset.Position(id.Pos())
			enclosing := ""
			if len(names) > 0 {
				enclosing = names[len(names)-1]
			}
			out = append(out, Row{
				File:      relPath(root, pos.Filename),
				Line:      pos.Line,
				Column:    pos.Column,
				Via:       classify(pkg, id, callees[id]),
				Enclosing: enclosing,
			})
		}
		return true
	})
	return out
}

// markWrites adds each identifier that n writes: each assignment target, an increment or
// decrement operand, and each composite-literal key that names a field. A selector target
// writes its selected name.
func markWrites(pkg *Package, n ast.Node, writes map[*ast.Ident]bool) {
	written := func(expr ast.Expr) {
		for {
			switch e := expr.(type) {
			case *ast.ParenExpr:
				expr = e.X
				continue
			case *ast.SelectorExpr:
				writes[e.Sel] = true
			case *ast.Ident:
				writes[e] = true
			}
			return
		}
	}
	switch s := n.(type) {
	case *ast.AssignStmt:
		for _, lhs := range s.Lhs {
			written(lhs)
		}
	case *ast.IncDecStmt:
		written(s.X)
	case *ast.CompositeLit:
		for _, elt := range s.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if key, ok := kv.Key.(*ast.Ident); ok && pkg.Info != nil {
				if field, ok := pkg.Info.Uses[key].(*types.Var); ok && field.IsField() {
					writes[key] = true
				}
			}
		}
	}
}

// SourceUse is one file of a revision's sources that names a declaration, at the line of
// its first naming.
type SourceUse struct {
	File string
	Line int
}

// SourceUses reads one revision's sources without type information, the way the deleted
// table reads its base side. file declares name at package level in the package whose
// import path is importPath. declared is the first line of that declaration, zero when
// file declares no such name. uses is every other source that names it: a file of the
// declaring package through the bare identifier, and any other file through a selector
// on its import of importPath. The reading is lexical, so a local name that shadows the
// declaration still counts as a use, and a source that does not parse names nothing. A
// method spelling has no receiver type to resolve here, so it can declare but names
// nothing.
func SourceUses(sources map[string]string, file, importPath, name string) (declared int, uses []SourceUse) {
	fset := token.NewFileSet()
	owner, err := parseDecls(fset, file, sources[file])
	if err != nil {
		return 0, nil
	}
	for _, site := range topLevelDecls(fset, owner) {
		if site.Name == name && declared == 0 {
			declared = site.Start
		}
	}
	if declared == 0 || strings.Contains(name, ".") {
		return declared, nil
	}
	for _, other := range slices.Sorted(maps.Keys(sources)) {
		parsed, err := parseDecls(fset, other, sources[other])
		if other == file || err != nil {
			continue
		}
		same := pathpkg.Dir(other) == pathpkg.Dir(file) && parsed.Name.Name == owner.Name.Name
		if line := firstUse(fset, parsed, same, importPath, name); line > 0 {
			uses = append(uses, SourceUse{File: other, Line: line})
		}
	}
	return declared, uses
}

// firstUse is the first line of file that names name: a bare identifier in the declaring
// package, or a selector on the file's import of importPath elsewhere. It is zero when
// the file names it nowhere.
func firstUse(fset *token.FileSet, file *ast.File, same bool, importPath, name string) int {
	qualifier := ""
	for _, spec := range file.Imports {
		if name, imported := ImportName(spec); imported == importPath {
			qualifier = name
		}
	}
	line := 0
	var visit func(ast.Node) bool
	visit = func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.SelectorExpr:
			if id, ok := n.X.(*ast.Ident); ok && !same && qualifier != "" && id.Name == qualifier && n.Sel.Name == name && line == 0 {
				line = fset.Position(n.Sel.Pos()).Line
			}
			ast.Inspect(n.X, visit)
			return false
		case *ast.Ident:
			if same && n.Name == name && line == 0 {
				line = fset.Position(n.Pos()).Line
			}
		}
		return line == 0
	}
	ast.Inspect(file, visit)
	return line
}

// ImportName is the one reading of an import's local name: its explicit name, else the
// last element of its path. path is the unquoted import path.
func ImportName(spec *ast.ImportSpec) (name, path string) {
	imported, _ := strconv.Unquote(spec.Path.Value)
	if spec.Name != nil {
		return spec.Name.Name, imported
	}
	return pathpkg.Base(imported), imported
}

// implementsRows answers the interface half of the via column. A query for an interface
// type name lists every named type in the loaded packages whose method set satisfies it,
// positioned at that type's own declaration. A row is emitted once per declaration, and
// the value and pointer method sets are one question: a pointer-receiver implementer is
// the edge an agent most needs, and it never appears in the value method set.
//
// A non-interface query, and an empty interface that every type would satisfy, emit
// nothing. The interface itself and every other interface are excluded, because "this
// interface implements that interface" is an embedding question, not a consumer edge.
func implementsRows(pkgs []*Package, target types.Object, root string) []Row {
	name, isType := target.(*types.TypeName)
	if !isType {
		return nil
	}
	iface, isIface := types.Unalias(name.Type()).Underlying().(*types.Interface)
	if !isIface || iface.Empty() {
		return nil
	}
	var out []Row
	seen := map[string]bool{}
	for _, pkg := range pkgs {
		if pkg.Types == nil {
			continue
		}
		scope := pkg.Types.Scope()
		for _, member := range scope.Names() {
			named, ok := types.Unalias(scope.Lookup(member).Type()).(*types.Named)
			if !ok || named.Obj().Name() != member || named.TypeParams().Len() > 0 {
				continue
			}
			if _, other := named.Underlying().(*types.Interface); other {
				continue
			}
			if !types.Implements(named, iface) && !types.Implements(types.NewPointer(named), iface) {
				continue
			}
			key := declKey(pkg.Fset, named.Obj())
			if seen[key] {
				continue
			}
			seen[key] = true
			pos := pkg.Fset.Position(named.Obj().Pos())
			out = append(out, Row{
				File:      relPath(root, pos.Filename),
				Line:      pos.Line,
				Column:    pos.Column,
				Via:       viaImplements,
				Enclosing: named.Obj().Name(),
			})
		}
	}
	return out
}
