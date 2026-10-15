package consumers

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"sort"
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
// a composite literal. A read is no write. UnnamedWrites owns every write that names no field.
func WriteRows(pkgs []*Package, target types.Object, root string) []Row {
	return sortRows(useRows(pkgs, target, root, true))
}

// The three Via values of an unnamed write.
const (
	viaUnkeyed   = "an unkeyed literal"
	viaGeneric   = "a generic struct"
	viaIdentical = "a type identical to its struct"
)

// UnnamedWrites is every write across pkgs that can reach the field target and names no field
// of its own struct type, in the order Rows uses. Via names the form. A generic struct gives one
// row at its declaration, because an instantiation is a struct of its own. Otherwise an unkeyed
// composite literal of the field's struct writes the field by position, and each appearance of
// another type identical to that struct gives a row, because a value of that type reaches the
// field through an assignment, a conversion, an argument, or a return. literalStructs owns which
// struct a literal writes, and sameStruct owns struct identity.
func UnnamedWrites(pkgs []*Package, target types.Object, root string) []Row {
	s, own := fieldStruct(target)
	if s == nil {
		return nil
	}
	var out []Row
	for _, pkg := range pkgs {
		if pkg.Info == nil {
			continue
		}
		add := func(at token.Pos, via string) {
			pos := pkg.Fset.Position(at)
			out = append(out, Row{File: relPath(root, pos.Filename), Line: pos.Line, Column: pos.Column, Via: via})
		}
		if genericDecl(own.Type()) {
			add(own.Pos(), viaGeneric)
			return out
		}
		foreign := identicalType(pkg, target, own, s)
		for expr, tv := range pkg.Info.Types {
			if foreign(tv.Type) {
				add(expr.Pos(), viaIdentical)
			}
		}
		for id, obj := range pkg.Info.Defs {
			if obj != nil && foreign(obj.Type()) {
				add(id.Pos(), viaIdentical)
			}
		}
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if lit, ok := n.(*ast.CompositeLit); ok && len(lit.Elts) > 0 && writesStruct(pkg.Info.TypeOf(lit), s) {
					if _, keyed := lit.Elts[0].(*ast.KeyValueExpr); !keyed {
						add(lit.Pos(), viaUnkeyed)
					}
				}
				return true
			})
		}
	}
	return sortRows(out)
}

// genericDecl reports whether the declared type t has type parameters, as a generic named type or
// a generic alias.
func genericDecl(t types.Type) bool {
	if alias, ok := t.(*types.Alias); ok && alias.TypeParams().Len() > 0 {
		return true
	}
	named, ok := types.Unalias(t).(*types.Named)
	return ok && named.TypeParams().Len() > 0
}

// sameStruct is the one source of struct identity. Go converts between two structs that
// types.IdenticalIgnoreTags compares as one, whatever their tags and alias spellings. A package
// and its test variant hold distinct objects for one named field type, so two structs that print
// as one type are one struct too. Either test can join two structs Go keeps apart; that only
// widens a census toward its diagnostic.
func sameStruct(a, b *types.Struct) bool {
	return types.IdenticalIgnoreTags(a, b) || types.TypeString(a, nil) == types.TypeString(b, nil)
}

// identicalType reports whether a type recorded in pkg, or an element of a recorded tuple, is a
// type other than own whose struct, alone or under one pointer, is the same struct as want. The
// struct that declares target is own's struct, named or nested, so it is no other type.
func identicalType(pkg *Package, target types.Object, own *types.TypeName, want *types.Struct) func(types.Type) bool {
	ownKey, fieldKey := declKey(pkg.Fset, own), declKey(pkg.Fset, target)
	var foreign func(t types.Type) bool
	foreign = func(t types.Type) bool {
		if tuple, ok := t.(*types.Tuple); ok {
			for i := range tuple.Len() {
				if foreign(tuple.At(i).Type()) {
					return true
				}
			}
			return false
		}
		if t == nil {
			return false
		}
		if p, ok := types.Unalias(t).(*types.Pointer); ok {
			t = p.Elem()
		}
		t = types.Unalias(t)
		u, ok := t.Underlying().(*types.Struct)
		if !ok || !sameStruct(u, want) {
			return false
		}
		if named, ok := t.(*types.Named); ok {
			return declKey(pkg.Fset, named.Obj()) != ownKey
		}
		for i := range u.NumFields() {
			if declKey(pkg.Fset, u.Field(i)) == fieldKey {
				return false
			}
		}
		return true
	}
	return foreign
}

// writesStruct reports whether a literal of the recorded type t writes the same struct as want.
func writesStruct(t types.Type, want *types.Struct) bool {
	if t == nil {
		return false
	}
	return slices.ContainsFunc(literalStructs(t), func(s *types.Struct) bool { return sameStruct(s, want) })
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
