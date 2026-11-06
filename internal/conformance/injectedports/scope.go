package injectedports

import (
	"fmt"
	"go/ast"
	"go/token"
)

// funcFacts is what one function body does to package state. Each fact is resolved
// through the body's lexical scopes, so an assignment to a shadowing local is not an
// assignment to the package variable of the same name.
type funcFacts struct {
	assigns map[string]string // seam identity (pkg.name) -> first assignment site
	injects map[string]bool   // seam identities assigned a value built from a parameter
	reads   map[string]bool   // identities the body references
	calls   map[string]bool   // callee identities the body calls
	// forward is the callee that the body forwards to, when the whole body is a single
	// call or a single return of a single call.
	forward string
}

// bodyScope is one lexical scope of a function body.
type bodyScope struct {
	names  map[string]bool
	parent *bodyScope
}

// owner returns the nearest scope that declares name, or nil for a package name.
func (s *bodyScope) owner(name string) *bodyScope {
	for ; s != nil; s = s.parent {
		if s.names[name] {
			return s
		}
	}
	return nil
}

// walker resolves the names of one body against its scopes and the file's imports.
type walker struct {
	fset  *token.FileSet
	file  *sourceFile
	pkg   string
	top   *bodyScope
	scope *bodyScope
	facts *funcFacts
}

// factsOf walks one function of package pkg. A name resolves to the package when no
// enclosing scope declares it. A selector on an import binding resolves to that package.
func factsOf(fset *token.FileSet, file *sourceFile, pkg string, fn *ast.FuncDecl) *funcFacts {
	w := newWalker(fset, file, pkg)
	if fn.Body == nil {
		return w.facts
	}
	declareFields(w.top, fn.Recv)
	declareFields(w.top, fn.Type.Params)
	declareFields(w.top, fn.Type.Results)
	if fun := forwardedCall(fn.Body); fun != nil {
		w.facts.forward = w.resolve(fun)
	}
	w.walk(fn.Body)
	return w.facts
}

// exprFacts walks one package-level initializer, such as a command table whose entries
// call a package variable.
func exprFacts(fset *token.FileSet, file *sourceFile, pkg string, expr ast.Expr) *funcFacts {
	w := newWalker(fset, file, pkg)
	w.walk(expr)
	return w.facts
}

func newWalker(fset *token.FileSet, file *sourceFile, pkg string) *walker {
	top := &bodyScope{names: map[string]bool{}}
	return &walker{fset: fset, file: file, pkg: pkg, top: top, scope: top, facts: &funcFacts{
		assigns: map[string]string{}, injects: map[string]bool{}, reads: map[string]bool{}, calls: map[string]bool{},
	}}
}

func (w *walker) resolve(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		if typed.Name != "_" && w.scope.owner(typed.Name) == nil {
			return w.pkg + "." + typed.Name
		}
	case *ast.SelectorExpr:
		if binding, ok := typed.X.(*ast.Ident); ok && w.scope.owner(binding.Name) == nil {
			if dir, imported := w.file.imports[binding.Name]; imported {
				return dir + "." + typed.Sel.Name
			}
		}
	}
	return ""
}

// usesParameter reports whether expr references a parameter of the walked function.
func (w *walker) usesParameter(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		if ident, ok := node.(*ast.Ident); ok && w.scope.owner(ident.Name) == w.top && w.top.names[ident.Name] {
			found = true
		}
		return !found
	})
	return found
}

func (w *walker) walk(root ast.Node) {
	var stack []ast.Node
	pushed := map[ast.Node]bool{}
	selectors := map[*ast.Ident]bool{}
	enter := func(node ast.Node) {
		w.scope = &bodyScope{names: map[string]bool{}, parent: w.scope}
		pushed[node] = true
	}
	ast.Inspect(root, func(node ast.Node) bool {
		if node == nil {
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if pushed[last] {
				w.scope = w.scope.parent
			}
			return true
		}
		stack = append(stack, node)
		switch typed := node.(type) {
		case *ast.BlockStmt, *ast.IfStmt, *ast.ForStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt,
			*ast.SelectStmt, *ast.CaseClause, *ast.CommClause:
			enter(node)
		case *ast.RangeStmt:
			enter(node)
			if typed.Tok == token.DEFINE {
				declareExprs(w.scope, typed.Key, typed.Value)
			}
		case *ast.FuncLit:
			enter(node)
			declareFields(w.scope, typed.Type.Params)
			declareFields(w.scope, typed.Type.Results)
		case *ast.AssignStmt:
			w.assign(typed)
		case *ast.DeclStmt:
			declareGen(w.scope, typed.Decl)
		case *ast.CallExpr:
			if callee := w.resolve(typed.Fun); callee != "" {
				w.facts.calls[callee] = true
			}
		case *ast.SelectorExpr:
			selectors[typed.Sel] = true
			if seam := w.resolve(typed); seam != "" {
				w.facts.reads[seam] = true
			}
		case *ast.Ident:
			if seam := w.resolve(typed); seam != "" && !selectors[typed] {
				w.facts.reads[seam] = true
			}
		}
		return true
	})
}

// assign records every target of one assignment, so the second name of a combined
// assignment resolves as surely as the first. A definition declares instead.
func (w *walker) assign(stmt *ast.AssignStmt) {
	if stmt.Tok == token.DEFINE {
		declareExprs(w.scope, stmt.Lhs...)
		return
	}
	for i, target := range stmt.Lhs {
		seam := w.resolve(assignedRoot(target))
		if seam == "" {
			continue
		}
		if _, seen := w.facts.assigns[seam]; !seen {
			w.facts.assigns[seam] = fmt.Sprintf("%s:%d", w.file.rel, w.fset.Position(target.Pos()).Line)
		}
		value := stmt.Rhs[0]
		if len(stmt.Rhs) == len(stmt.Lhs) {
			value = stmt.Rhs[i]
		}
		if w.usesParameter(value) {
			w.facts.injects[seam] = true
		}
	}
}

// assignedRoot reduces an assignment target to the variable that holds the replaced
// value. A field of a package collaborator struct belongs to that package variable.
func assignedRoot(expr ast.Expr) ast.Expr {
	for {
		switch typed := expr.(type) {
		case *ast.ParenExpr:
			expr = typed.X
		case *ast.SelectorExpr:
			if _, qualified := typed.X.(*ast.Ident); qualified {
				return typed
			}
			expr = typed.X
		default:
			return expr
		}
	}
}

// forwardedCall returns the callee expression of a body that is one call statement or
// one return of one call. Such a body forwards to its callee.
func forwardedCall(body *ast.BlockStmt) ast.Expr {
	if len(body.List) != 1 {
		return nil
	}
	var call ast.Expr
	switch stmt := body.List[0].(type) {
	case *ast.ExprStmt:
		call = stmt.X
	case *ast.ReturnStmt:
		if len(stmt.Results) == 1 {
			call = stmt.Results[0]
		}
	}
	if typed, ok := call.(*ast.CallExpr); ok {
		return typed.Fun
	}
	return nil
}

func declareFields(scope *bodyScope, fields *ast.FieldList) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		for _, name := range field.Names {
			scope.names[name.Name] = true
		}
	}
}

func declareExprs(scope *bodyScope, exprs ...ast.Expr) {
	for _, expr := range exprs {
		if ident, ok := expr.(*ast.Ident); ok {
			scope.names[ident.Name] = true
		}
	}
}

func declareGen(scope *bodyScope, decl ast.Decl) {
	gen, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}
	for _, spec := range gen.Specs {
		switch named := spec.(type) {
		case *ast.ValueSpec:
			for _, name := range named.Names {
				scope.names[name.Name] = true
			}
		case *ast.TypeSpec:
			scope.names[named.Name.Name] = true
		}
	}
}
