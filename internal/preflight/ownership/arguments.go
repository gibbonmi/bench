package ownership

import (
	"fmt"
	"go/ast"
	"slices"
)

// grammarUnit is one producer that a grammar argument resolves to: a package declaration
// when node is nil, otherwise a producer expression that rootUnit adds under key.
type grammarUnit struct {
	key  string
	node ast.Node
}

// resolve is each producer of one grammar argument in fn; key names an expression producer.
// A package declaration, a single-definition local, a call of a package function or method, a
// literal, and a declaration of another package of the module resolve. A parameter resolves
// when callers resolves it. A reassigned local and a value from outside the module stay
// unresolved. visiting holds each function whose parameter the resolution is inside.
func (p *producers) resolve(fn *ast.FuncDecl, arg ast.Expr, key string, visiting map[*ast.FuncDecl]bool) ([]grammarUnit, bool) {
	file := p.fset.Position(arg.Pos()).Filename
	switch arg := arg.(type) {
	case *ast.Ident:
		value, local, ok := localValue(fn, arg.Name)
		switch {
		case local && ok:
			return []grammarUnit{{fn.Name.Name + "." + arg.Name, value}}, true
		case local:
			return p.callers(fn, arg.Name, visiting)
		}
		if _, declared := p.units[arg.Name]; declared {
			return []grammarUnit{{arg.Name, nil}}, true
		}
	case *ast.CallExpr:
		if p.resolvedCall(file, arg) {
			return []grammarUnit{{key, arg}}, true
		}
	case *ast.CompositeLit, *ast.UnaryExpr:
		return []grammarUnit{{key, arg}}, true
	case *ast.SelectorExpr:
		if imported, ok := p.imported(file, arg.X); ok {
			if _, in := p.inModule(imported); in {
				return []grammarUnit{{key, arg}}, true
			}
		}
	}
	return nil, false
}

// callers resolves parameter name of fn through each call site of fn in the package: the
// argument at the parameter's place resolves as resolve states. fn must be an unexported
// function, not a method, with an unambiguous name that the package only calls, so the
// package holds every call site. A name that is not a parameter of fn, a recursive
// resolution, a use of fn as a value, a call site outside a function or with too few
// arguments, a call site that does not resolve, or no call site leaves it unresolved.
func (p *producers) callers(fn *ast.FuncDecl, name string, visiting map[*ast.FuncDecl]bool) ([]grammarUnit, bool) {
	at := paramIndex(fn, name)
	if at < 0 || fn.Recv != nil || ast.IsExported(fn.Name.Name) || p.ambiguous[fn.Name.Name] || visiting[fn] {
		return nil, false
	}
	visiting[fn] = true
	defer delete(visiting, fn)
	var found []grammarUnit
	sites, ok := 0, true
	for _, file := range p.files {
		for _, decl := range file.Decls {
			caller, _ := decl.(*ast.FuncDecl)
			called := map[*ast.Ident]bool{fn.Name: true}
			ast.Inspect(decl, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.CallExpr:
					id, isCall := n.Fun.(*ast.Ident)
					if !isCall || id.Name != fn.Name.Name {
						return true
					}
					called[id], sites = true, sites+1
					if caller == nil || at >= len(n.Args) {
						ok = false
						return true
					}
					units, resolved := p.resolve(caller, n.Args[at], fmt.Sprintf("%s#%s%d", caller.Name.Name, fn.Name.Name, sites), visiting)
					found, ok = append(found, units...), ok && resolved
				case *ast.Ident:
					ok = ok && (n.Name != fn.Name.Name || called[n])
				}
				return true
			})
		}
	}
	return found, ok && sites > 0
}

// paramIndex is the place of parameter name in fn's parameter list, or -1.
func paramIndex(fn *ast.FuncDecl, name string) int {
	at := 0
	for _, field := range fn.Type.Params.List {
		if len(field.Names) == 0 {
			at++
		}
		for _, id := range field.Names {
			if id.Name == name {
				return at
			}
			at++
		}
	}
	return -1
}

// resolvedCall reports whether call, in file, invokes a package function or a method this
// package declares.
func (p *producers) resolvedCall(file string, call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		_, declared := p.units[fun.Name]
		return declared
	case *ast.SelectorExpr:
		_, imported := p.imported(file, fun.X)
		return !imported && len(p.methods[fun.Sel.Name]) > 0
	}
	return false
}

// localValue is the one value fn assigns to a local name. local reports that fn binds the
// name at all; ok is false for a parameter or for a local with no single defining value.
func localValue(fn *ast.FuncDecl, name string) (value ast.Expr, local, ok bool) {
	for _, list := range []*ast.FieldList{fn.Recv, fn.Type.Params, fn.Type.Results} {
		if list != nil && slices.ContainsFunc(list.List, func(field *ast.Field) bool {
			return slices.ContainsFunc(field.Names, func(id *ast.Ident) bool { return id.Name == name })
		}) {
			return nil, true, false
		}
	}
	var values []ast.Expr
	bind := func(names []ast.Expr, rhs []ast.Expr) {
		for i, lhs := range names {
			if isIdent(lhs, name) {
				local = true
				if len(names) == len(rhs) {
					values = append(values, rhs[i])
				} else {
					values = append(values, nil)
				}
			}
		}
	}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		switch s := node.(type) {
		case *ast.AssignStmt:
			bind(s.Lhs, s.Rhs)
		case *ast.ValueSpec:
			names := make([]ast.Expr, len(s.Names))
			for i, id := range s.Names {
				names[i] = id
			}
			bind(names, s.Values)
		case *ast.RangeStmt:
			bind([]ast.Expr{s.Key, s.Value}, nil)
		case *ast.FuncLit:
			for _, field := range s.Type.Params.List {
				for _, id := range field.Names {
					bind([]ast.Expr{id}, nil)
				}
			}
		}
		return true
	})
	if len(values) == 1 && values[0] != nil {
		return values[0], true, true
	}
	return nil, local, false
}

func isIdent(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}
