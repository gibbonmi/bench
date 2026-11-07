package injectedports

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
)

// derivedPort is one port the source sweep found, carrying the site that injected it. A
// red then says where the port entered the package, not only that it did.
type derivedPort struct {
	pkg, name, site string
}

// namedPorts derives the injected ports of every production package of tree, keyed by
// package. A test-support package declares no production port, so it contributes none.
func (c *seamCensus) namedPorts(tree *sourceTree) map[string][]derivedPort {
	ports := map[string][]derivedPort{}
	for _, dir := range sortedKeys(tree.packages) {
		if c.support[dir] {
			continue
		}
		if found := packagePorts(tree.fset, tree.packages[dir]); len(found) > 0 {
			ports[dir] = found
		}
	}
	return ports
}

// packagePorts reports the injected ports that one package's production source declares.
//
// A port is a named type whose shape is port-shaped: a non-empty interface, a func type,
// or a struct whose every field is a func. The package must also inject it, by passing it
// as a function parameter or by narrowing to it through a type assertion or type switch.
//
// The assertion arm is not decoration. A capability widened out of an already-injected
// owner is never a parameter at all. A parameter-only rule would therefore leave exactly
// the silently-downgradable capabilities out of the inventory. Exportedness is not part
// of the rule either. A package main declares nothing exported, and its ports are ports
// all the same.
func packagePorts(fset *token.FileSet, pkg *sourcePackage) []derivedPort {
	shaped := map[string]bool{}
	sites := map[string]string{}
	for _, file := range pkg.prod {
		rel := file.rel
		ast.Inspect(file.file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.TypeSpec:
				if portShapedType(typed.Type) {
					shaped[typed.Name.Name] = true
				}
			case *ast.FuncDecl:
				if typed.Type.Params == nil {
					return true
				}
				for _, param := range typed.Type.Params.List {
					recordInjectionSite(sites, param.Type, fmt.Sprintf("a parameter of %s in %s:%d", typed.Name.Name, rel, fset.Position(param.Pos()).Line))
				}
			case *ast.TypeAssertExpr:
				recordInjectionSite(sites, typed.Type, fmt.Sprintf("a type assertion in %s:%d", rel, fset.Position(typed.Pos()).Line))
			case *ast.CaseClause:
				for _, expr := range typed.List {
					recordInjectionSite(sites, expr, fmt.Sprintf("a type switch case in %s:%d", rel, fset.Position(expr.Pos()).Line))
				}
			}
			return true
		})
	}
	var ports []derivedPort
	for name := range shaped {
		if site, injected := sites[name]; injected {
			ports = append(ports, derivedPort{pkg: pkg.rel, name: name, site: site})
		}
	}
	sort.Slice(ports, func(i, j int) bool { return ports[i].name < ports[j].name })
	return ports
}

// recordInjectionSite keeps the first site a name is injected at. A port injected twice
// is still one port, and the first site is the stable one to name in a diagnostic.
func recordInjectionSite(sites map[string]string, expr ast.Expr, site string) {
	name := typeIdentName(expr)
	if name == "" {
		return
	}
	if _, seen := sites[name]; !seen {
		sites[name] = site
	}
}

// typeIdentName reduces a type expression to the package-local name it names, seeing
// through a pointer. A qualified type belongs to another package and is not this
// package's port to register.
func typeIdentName(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.StarExpr:
		return typeIdentName(typed.X)
	}
	return ""
}

// portShapedType reports whether a type declaration has the shape of a port. An empty
// interface constrains nothing and names no producer, so it is not one.
func portShapedType(expr ast.Expr) bool {
	switch typed := expr.(type) {
	case *ast.InterfaceType:
		return typed.Methods != nil && len(typed.Methods.List) > 0
	case *ast.FuncType:
		return true
	case *ast.StructType:
		if typed.Fields == nil || len(typed.Fields.List) == 0 {
			return false
		}
		for _, field := range typed.Fields.List {
			if _, isFunc := field.Type.(*ast.FuncType); !isFunc {
				return false
			}
		}
		return true
	}
	return false
}

// declareFields indexes the fields of one production struct declaration of pkg in order. An
// embedded field is named by its type.
func (c *seamCensus) declareFields(pkg string, spec ast.Spec) {
	named, ok := spec.(*ast.TypeSpec)
	if !ok {
		return
	}
	strct, ok := named.Type.(*ast.StructType)
	if !ok || strct.Fields == nil {
		return
	}
	id := pkg + "." + named.Name.Name
	for _, field := range strct.Fields.List {
		names := []string{typeIdentName(field.Type)}
		if len(field.Names) > 0 {
			names = names[:0]
			for _, name := range field.Names {
				names = append(names, name.Name)
			}
		}
		for _, name := range names {
			c.structFields[id] = append(c.structFields[id], name)
			c.fieldTypes[id+"."+name] = typeIdentName(field.Type)
		}
	}
}

// assignedValue is the value that one assignment gives its target at index i. A
// multi-value call gives every target the call.
func assignedValue(stmt *ast.AssignStmt, i int) ast.Expr {
	if len(stmt.Rhs) == len(stmt.Lhs) {
		return stmt.Rhs[i]
	}
	return stmt.Rhs[0]
}

// store records one assignment for the fault trace: a value of a local, or a store into a
// field through a selector that names no package variable.
func (w *walker) store(target, value ast.Expr) {
	switch typed := target.(type) {
	case *ast.Ident:
		if owner := w.scope.owner(typed.Name); owner != nil {
			w.local(owner, typed.Name, value)
		}
	case *ast.SelectorExpr:
		if w.resolve(typed) == "" {
			w.facts.fieldSets = append(w.facts.fieldSets, fieldSet{field: typed.Sel.Name, value: w.operand(value)})
		}
	}
}

func (w *walker) local(owner *bodyScope, name string, value ast.Expr) {
	if w.facts.locals[owner] == nil {
		w.facts.locals[owner] = map[string][]operand{}
	}
	w.facts.locals[owner][name] = append(w.facts.locals[owner][name], w.operand(value))
}

// declareValues records the values of a local var declaration.
func (w *walker) declareValues(decl ast.Decl) {
	gen, ok := decl.(*ast.GenDecl)
	if !ok || gen.Tok != token.VAR {
		return
	}
	for _, spec := range gen.Specs {
		value, ok := spec.(*ast.ValueSpec)
		for i := 0; ok && i < len(value.Names); i++ {
			w.local(w.scope, value.Names[i].Name, declaredValue(value, i))
		}
	}
}

// declaredValue is the value that one var declaration gives its name at index i. A name
// declared without a value holds its zero value, which for the port type is nil. A
// multi-value call gives every name the call.
func declaredValue(spec *ast.ValueSpec, i int) ast.Expr {
	switch {
	case len(spec.Values) == 0:
		return ast.NewIdent("nil")
	case len(spec.Values) == len(spec.Names):
		return spec.Values[i]
	}
	return spec.Values[0]
}

// composite records each element of a literal of a resolved struct as a store into a field.
func (w *walker) composite(lit *ast.CompositeLit) {
	strct := w.resolve(lit.Type)
	if strct == "" {
		return
	}
	for index, elt := range lit.Elts {
		set := fieldSet{strct: strct, index: index, value: w.operand(elt)}
		if pair, keyed := elt.(*ast.KeyValueExpr); keyed {
			key, named := pair.Key.(*ast.Ident)
			if !named {
				continue
			}
			set.field, set.value = key.Name, w.operand(pair.Value)
		}
		w.facts.fieldSets = append(w.facts.fieldSets, set)
	}
}
