package consumers

import (
	"go/ast"
	"go/types"
	"sort"
)

// The three values the via column carries. reference is the default: an appearance of
// the declaration that is neither a static call nor a method-set satisfaction.
const (
	viaReference  = "reference"
	viaCall       = "call"
	viaImplements = "implements"
)

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
	out = append(out, implementsRows(pkgs, target, root)...)
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
	// Inspect visits a call expression before the identifier inside it, so the callee set
	// is complete for every identifier by the time the walk reaches it. This is the one
	// derivation of call position; classify reads it and never re-walks.
	callees := map[*ast.Ident]bool{}
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
		if call, isCall := n.(*ast.CallExpr); isCall {
			if id := staticCallee(call); id != nil {
				callees[id] = true
			}
		}
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
				Via:       classify(pkg, id, callees[id]),
				Enclosing: enclosing,
			})
		}
		return true
	})
	return out
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
