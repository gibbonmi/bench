package consumers

import (
	"fmt"
	"go/token"
	"go/types"
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
	for _, p := range parts {
		if p == "" {
			return nil, fmt.Errorf("malformed symbol %q", query)
		}
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
		out = append(out, Match{Obj: o, PkgPath: pkg.PkgPath, Kind: kind})
	}
	n := len(parts)
	for _, pkg := range pkgs {
		if pkg.Types == nil {
			continue
		}
		scope := pkg.Types.Scope()
		// A trailing single name reads as a package-scope declaration.
		if pkgSuffixMatches(pkg.PkgPath, parts[:n-1]) {
			name := parts[n-1]
			if obj := scope.Lookup(name); obj != nil {
				add(pkg, obj, name, objKind(obj))
			}
		}
		// A trailing pair reads as Type.Method. Both readings are attempted, because
		// `a.B` is ambiguous between a package-qualified name and a method on a
		// package-scope type in the current package.
		if n >= 2 && pkgSuffixMatches(pkg.PkgPath, parts[:n-2]) {
			if m := lookupMember(scope, parts[n-2], parts[n-1]); m != nil {
				add(pkg, m, parts[n-2]+"."+parts[n-1], objKind(m))
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no declaration named %q", query)
	}
	qualify(out, names)
	return out, nil
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
	if len(suffix) == 0 {
		return true
	}
	want := strings.Join(suffix, "/")
	return path == want || strings.HasSuffix(path, "/"+want)
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

// declKey is the one source of declaration identity, and every comparison of "is this
// the same declaration" goes through it. Object pointers cannot serve: go/packages loads
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
