package consumers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"
)

// Match is one declaration a query resolved to. Qualified is the exact re-query
// spelling for this match, so an ambiguous answer can hand the agent a literal
// argument instead of a name it must reconstruct.
type Match struct {
	Obj       types.Object
	PkgPath   string
	Qualified string
	Kind      string
	// Name is the declaration's spelling inside its package: a name, or Type.Member.
	Name string
}

// Resolve maps a symbol query onto the declarations it names. The grammar is an
// import-path-suffix qualification (`outline.Command`, `pkg.Type.Method`) or a bare
// identifier. A bare identifier is deliberately allowed to return several matches; the
// caller decides whether that is an answer or a candidates table.
//
// Every returned Obj is an origin object, so a match found through an alias spelling and
// the same match found through the origin spelling are one match, not two.
func Resolve(pkgs []*Package, query string) ([]Match, error) {
	parts := strings.Split(query, ".")
	if slices.Contains(parts, "") {
		return nil, fmt.Errorf("malformed symbol %q", query)
	}
	var out []Match
	var names []string
	seen := map[string]bool{}
	add := func(pkg *Package, obj types.Object, name, kind string) {
		if obj == nil {
			return
		}
		o := origin(obj)
		key := declKey(pkg.Fset, o)
		if seen[key] {
			return
		}
		seen[key] = true
		names = append(names, name)
		out = append(out, Match{Obj: o, PkgPath: pkg.PkgPath, Kind: kind, Name: name})
	}
	for _, pkg := range pkgs {
		if pkg.Types == nil {
			continue
		}
		scope := pkg.Types.Scope()
		for _, local := range readings(pkg.PkgPath, parts) {
			if len(local) == 1 {
				if obj := scope.Lookup(local[0]); obj != nil {
					add(pkg, obj, local[0], objKind(obj))
				}
			} else if m := lookupMember(scope, local[0], local[1]); m != nil {
				add(pkg, m, local[0]+"."+local[1], objKind(m))
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no declaration named %q", query)
	}
	qualify(out, names)
	return out, nil
}

// readings is every declaration spelling that the query parts name inside the package at
// pkgPath. A trailing single name reads as a package-scope declaration, and a trailing pair
// reads as Type.Member, each when the parts before it qualify pkgPath. Both readings are
// kept, because `a.B` is ambiguous between a package-qualified name and a member of a
// package-scope type in an unqualified package. The single name comes first.
func readings(pkgPath string, parts []string) [][]string {
	n := len(parts)
	var out [][]string
	if pkgSuffixMatches(pkgPath, parts[:n-1]) {
		out = append(out, parts[n-1:])
	}
	if n >= 2 && pkgSuffixMatches(pkgPath, parts[:n-2]) {
		out = append(out, parts[n-2:])
	}
	return out
}

// LocalName is the spelling inside the package at pkgPath that a symbol query names, by the
// grammar Resolve reads: a name, or Type.Member. When both readings qualify, the single name
// wins, as the package-qualified reading. A malformed query, or one that names no
// declaration spelling of that package, is an error.
func LocalName(pkgPath, query string) (string, error) {
	parts := strings.Split(query, ".")
	if slices.Contains(parts, "") {
		return "", fmt.Errorf("malformed symbol %q", query)
	}
	local := readings(pkgPath, parts)
	if len(local) == 0 {
		return "", fmt.Errorf("symbol %q names no declaration of package %q", query, pkgPath)
	}
	return strings.Join(local[0], "."), nil
}

// qualify fills in each match's re-query spelling. The spelling is an import-path suffix
// plus the declaration name, and it must resolve back to one match, or the agent re-runs
// the same ambiguous query. So the suffix grows one path segment at a time until every
// spelling is unique. Growing is enough because a suffix spelling matches exactly the
// declarations of that name in packages the suffix qualifies, which is a subset of what
// the original query already matched.
//
// Two limits stop the growth. A segment that contains a dot cannot be spelled, because
// the grammar splits the query on dots; and two matches can share a whole import path,
// as a package and its external test package do not but a repeated load would. At either
// limit the longest distinguishing suffix stands, and the spelling stays ambiguous.
func qualify(out []Match, names []string) {
	for depth := 1; ; depth++ {
		counts := map[string]int{}
		for i := range out {
			out[i].Qualified = pathSuffix(out[i].PkgPath, depth) + "." + names[i]
			counts[out[i].Qualified]++
		}
		unique := true
		for i := range out {
			if counts[out[i].Qualified] > 1 {
				unique = false
			}
		}
		if unique || !spellableDeeper(out, depth+1) {
			return
		}
	}
}

// spellableDeeper reports whether a longer suffix would say something new. It is false
// when no match has another segment left, and false when the next segment of any match
// carries a dot the query grammar cannot express.
func spellableDeeper(out []Match, depth int) bool {
	longer := false
	for _, m := range out {
		next := pathSuffix(m.PkgPath, depth)
		if strings.Contains(next, ".") {
			return false
		}
		if next != pathSuffix(m.PkgPath, depth-1) {
			longer = true
		}
	}
	return longer
}

// lookupMember finds the method or field named member on the package-scope named type
// typeName, including one promoted from an embedded field. It returns nil when either
// half is absent. A method and a field are one lookup, because Type.Name reads the same
// either way at the query surface.
func lookupMember(scope *types.Scope, typeName, member string) types.Object {
	obj := scope.Lookup(typeName)
	if obj == nil {
		return nil
	}
	tn, ok := obj.(*types.TypeName)
	if !ok {
		return nil
	}
	m, _, _ := types.LookupFieldOrMethod(tn.Type(), true, tn.Pkg(), member)
	switch m.(type) {
	case *types.Func, *types.Var:
		return m
	}
	return nil
}

// pkgSuffixMatches reports whether path is qualified by the dotted suffix the query
// carried. An empty suffix matches every package, which is what a bare identifier means.
func pkgSuffixMatches(path string, suffix []string) bool {
	return PathHasSuffix(path, strings.Join(suffix, "/"))
}

// PathHasSuffix reports whether the import path ends in the slash-separated suffix at an
// element boundary. An empty suffix matches every path.
func PathHasSuffix(path, suffix string) bool {
	return suffix == "" || path == suffix || strings.HasSuffix(path, "/"+suffix)
}

// pathSuffix is the import path's last depth elements, the spelling a qualified re-query
// carries. A depth past the path's own length is the whole path.
func pathSuffix(path string, depth int) string {
	parts := strings.Split(path, "/")
	if depth >= len(parts) {
		return path
	}
	return strings.Join(parts[len(parts)-depth:], "/")
}

// objKind names the declaration class for a candidates row.
func objKind(obj types.Object) string {
	switch o := obj.(type) {
	case *types.Func:
		if sig, ok := o.Type().(*types.Signature); ok && sig.Recv() != nil {
			return "method"
		}
		return "func"
	case *types.TypeName:
		return "type"
	case *types.Const:
		return "const"
	case *types.Var:
		if o.IsField() {
			return "field"
		}
		return "var"
	}
	return "decl"
}

// declKey is the one source of typed declaration identity: every comparison of two loaded
// objects goes through it. The base side of a blast pair is parsed without type
// information and has no object to key, so sourceDeclNames keys a package by its directory
// there. Object pointers cannot serve: go/packages loads
// a package that has tests both plainly and as a test variant, so one declaration arrives
// as two objects with no pointer relation. The key names the declaring package and the
// origin declaration's position instead, which is stable across every load of one tree.
//
// fset must be the file set that produced obj's position. Both producers of Package — the
// loader seam and the in-process test fixtures — share one file set across all packages,
// so a caller passes any loaded package's Fset.
func declKey(fset *token.FileSet, obj types.Object) string {
	o := origin(obj)
	path := ""
	if o.Pkg() != nil {
		path = o.Pkg().Path()
	}
	pos := "-"
	if fset != nil && o.Pos().IsValid() {
		pos = fset.Position(o.Pos()).String()
	}
	return path + "@" + pos + "@" + o.Name()
}

// origin is the object-level half of identity: an alias resolves to the origin type's
// name, and a generic instantiation resolves to its generic origin, so two spellings of
// one declaration key to one object.
func origin(obj types.Object) types.Object {
	switch o := obj.(type) {
	case *types.TypeName:
		if named, ok := types.Unalias(o.Type()).(*types.Named); ok {
			return named.Obj()
		}
		return o
	case *types.Func:
		return o.Origin()
	case *types.Var:
		return o.Origin()
	}
	return obj
}

// fieldStruct is the struct type that declares the field target, found among the
// package-scope types of target's package, including an anonymous struct nested in one, and the
// package-scope type whose struct holds it. Both are nil when target is no field.
func fieldStruct(target types.Object) (*types.Struct, *types.TypeName) {
	field, ok := target.(*types.Var)
	if !ok || !field.IsField() || field.Pkg() == nil {
		return nil, nil
	}
	field = field.Origin()
	var declaring func(s *types.Struct) *types.Struct
	declaring = func(s *types.Struct) *types.Struct {
		for i := range s.NumFields() {
			if s.Field(i) == field {
				return s
			}
			if inner, ok := types.Unalias(s.Field(i).Type()).(*types.Struct); ok {
				if found := declaring(inner); found != nil {
					return found
				}
			}
		}
		return nil
	}
	scope := field.Pkg().Scope()
	for _, name := range scope.Names() {
		if tn, ok := scope.Lookup(name).(*types.TypeName); ok {
			if s, ok := tn.Type().Underlying().(*types.Struct); ok {
				if found := declaring(s); found != nil {
					return found, tn
				}
			}
		}
	}
	return nil, nil
}

// literalStructs is every struct type that a composite literal of the recorded type t can
// write. An elided element literal records its hint type, so t can also be a
// pointer to the struct, a named pointer, or a type parameter. A type parameter gives the
// struct of each term of its constraint, alone or under the one elided pointer. The terms of a
// literal that type-checks share one core type, so a term outside the type set only widens
// the answer toward the diagnostic and never hides a write.
func literalStructs(t types.Type) []*types.Struct {
	var out []*types.Struct
	var visit func(t types.Type, deref bool)
	visit = func(t types.Type, deref bool) {
		switch u := types.Unalias(t).Underlying().(type) {
		case *types.Struct:
			out = append(out, u)
		case *types.Pointer:
			if deref {
				visit(u.Elem(), false)
			}
		case *types.Interface:
			for i := range u.NumEmbeddeds() {
				visit(u.EmbeddedType(i), deref)
			}
		case *types.Union:
			for i := range u.Len() {
				visit(u.Term(i).Type(), deref)
			}
		}
	}
	visit(t, true)
	return out
}

// declNames is the one derivation of the declaration-name forms this surface reports: a
// function, a Type.Method, a package-level type, and a package-level var or const spec. A
// node outside those reports false, so the row walker's innermost open name is the answer
// for every row, and a use with no enclosing declaration at all keeps the empty string.
// fileScope says no named declaration is open yet.
//
// A spec declaring several names reports them all, because the blast walk names each
// declaration separately. The row walker reads the first, because a use inside the spec
// sits inside one declaration statement whatever that statement declares.
func declNames(n ast.Node, fileScope bool) ([]string, bool) {
	switch d := n.(type) {
	case *ast.FuncDecl:
		if recv := receiverName(d); recv != "" {
			return []string{recv + "." + d.Name.Name}, true
		}
		return []string{d.Name.Name}, true
	case *ast.TypeSpec:
		if fileScope {
			return []string{d.Name.Name}, true
		}
	case *ast.ValueSpec:
		// Only a package-level spec names a row. A local `var x = ...` inside a function
		// is not the declaration an agent is looking for; the function is.
		if fileScope && len(d.Names) > 0 {
			names := make([]string, 0, len(d.Names))
			for _, name := range d.Names {
				names = append(names, name.Name)
			}
			return names, true
		}
	}
	return nil, false
}

// DeclNames is declNames for a package-level node: the names one top-level declaration or
// spec declares, in the spelling a row's enclosing cell carries.
func DeclNames(n ast.Node) ([]string, bool) {
	return declNames(n, true)
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
