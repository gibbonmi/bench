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
	// ctx is the resolution context of the body. envReads are the name operands of its
	// environment reads and envSets the name operands of its Setenv calls. entries are its
	// outermost string expressions that flow to some place, any of which can be a
	// "NAME=value" environment entry, and flows are the places its other values flow to.
	// sites are its calls with their arguments, so a parameter resolves at its callers.
	ctx               *funcContext
	envReads, envSets []operand
	entries           []envEntry
	flows             []flow
	sites             []callSite
	// locals are the values assigned to each local of a scope, fieldSets the stores into
	// struct fields, and callScopes the scope of each call. The fault trace reads them.
	locals     map[*bodyScope]map[string][]operand
	fieldSets  []fieldSet
	callScopes map[*ast.CallExpr]*bodyScope
}

// envEntry is one string expression with the place its value flows to.
type envEntry struct {
	operand
	to flowNode
}

// funcContext is where a body's names resolve: its file, its package, its owner name, and
// its parameters in order. A package-level initializer has no owner and no parameters.
type funcContext struct {
	file   *sourceFile
	pkg    string
	owner  string
	top    *bodyScope
	params []string
	// variadic holds when the last parameter takes the tail arguments of a call.
	variadic bool
}

// operand is one expression with the lexical scope it was written in.
type operand struct {
	expr  ast.Expr
	scope *bodyScope
	site  string
}

// callSite is one call of a resolved callee with its arguments.
type callSite struct {
	callee string
	args   []ast.Expr
	scope  *bodyScope
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
	w.facts.ctx.owner = ownerName(fn)
	if fn.Body == nil {
		return w.facts
	}
	declareFields(w.top, fn.Recv)
	declareFields(w.top, fn.Type.Params)
	declareFields(w.top, fn.Type.Results)
	for _, field := range fn.Type.Params.List {
		for _, name := range field.Names {
			w.facts.ctx.params = append(w.facts.ctx.params, name.Name)
		}
		_, w.facts.ctx.variadic = field.Type.(*ast.Ellipsis)
	}
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
		locals: map[*bodyScope]map[string][]operand{}, callScopes: map[*ast.CallExpr]*bodyScope{},
		ctx: &funcContext{file: file, pkg: pkg, top: top},
	}}
}

func (w *walker) resolve(expr ast.Expr) string { return resolveName(w.file, w.pkg, w.scope, expr) }

// resolveName resolves a package-level name or an import-qualified selector to its
// identity, written "pkg.name". A name that a scope declares is local and resolves to "".
func resolveName(file *sourceFile, pkg string, scope *bodyScope, expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		if typed.Name != "_" && scope.owner(typed.Name) == nil {
			return pkg + "." + typed.Name
		}
	case *ast.SelectorExpr:
		if binding, ok := typed.X.(*ast.Ident); ok && scope.owner(binding.Name) == nil {
			if dir, imported := file.imports[binding.Name]; imported {
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
			w.declareValues(typed.Decl)
		case *ast.CompositeLit:
			w.composite(typed)
		case *ast.CallExpr:
			w.call(typed)
			if callee := w.callee(typed.Fun); callee != "" {
				w.flowFrom(flowNode{kind: flowReturn, name: callee}, stack)
			}
		case *ast.BasicLit, *ast.BinaryExpr:
			if to, ok := w.carrier(stack); ok && !concatOperand(stack, node) {
				w.facts.entries = append(w.facts.entries, envEntry{operand: w.operand(node.(ast.Expr)), to: to})
			}
		case *ast.SelectorExpr:
			selectors[typed.Sel] = true
			if seam := w.resolve(typed); seam != "" {
				w.facts.reads[seam] = true
			} else {
				w.flowFrom(w.field(typed.Sel.Name), stack)
			}
		case *ast.Ident:
			if seam := w.resolve(typed); seam != "" && !selectors[typed] {
				w.facts.reads[seam] = true
			} else if owner := w.scope.owner(typed.Name); owner != nil && !selectors[typed] {
				w.flowFrom(flowNode{kind: flowLocal, scope: owner, name: typed.Name}, stack)
			}
		}
		return true
	})
}

// call records one call: its callee, its arguments for caller-side resolution, an
// environment read, or a Setenv of any receiver. A method call records the method name of
// the walked package.
func (w *walker) call(call *ast.CallExpr) {
	w.facts.callScopes[call] = w.scope
	callee := w.resolve(call.Fun)
	if callee != "" {
		w.facts.calls[callee] = true
		w.facts.sites = append(w.facts.sites, callSite{callee: callee, args: call.Args, scope: w.scope})
	} else if method := w.callee(call.Fun); method != "" {
		w.facts.calls[method] = true
	}
	if len(call.Args) == 0 {
		return
	}
	if callee == "os.Getenv" || callee == "os.LookupEnv" {
		w.facts.envReads = append(w.facts.envReads, w.operand(call.Args[0]))
	}
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Setenv" {
		w.facts.envSets = append(w.facts.envSets, w.operand(call.Args[0]))
	}
}

func (w *walker) operand(expr ast.Expr) operand {
	return operand{expr: expr, scope: w.scope, site: fmt.Sprintf("%s:%d", w.file.rel, w.fset.Position(expr.Pos()).Line)}
}

// concatOperand reports whether node is an operand of an enclosing string concatenation,
// so only the outermost expression of a "NAME=" + value entry is recorded.
func concatOperand(stack []ast.Node, node ast.Node) bool {
	if len(stack) < 2 {
		return false
	}
	parent, ok := stack[len(stack)-2].(*ast.BinaryExpr)
	return ok && parent.Op == token.ADD && (parent.X == node || parent.Y == node)
}

// assign records every target of one assignment, so the second name of a combined
// assignment resolves as surely as the first. A definition declares instead. Every
// target's value is also stored for the fault trace.
func (w *walker) assign(stmt *ast.AssignStmt) {
	if stmt.Tok == token.DEFINE {
		declareExprs(w.scope, stmt.Lhs...)
	}
	for i, target := range stmt.Lhs {
		value := assignedValue(stmt, i)
		w.store(target, value)
		seam := ""
		if stmt.Tok != token.DEFINE {
			seam = w.resolve(w.assignedRoot(target))
		}
		if seam == "" {
			continue
		}
		if _, seen := w.facts.assigns[seam]; !seen {
			w.facts.assigns[seam] = fmt.Sprintf("%s:%d", w.file.rel, w.fset.Position(target.Pos()).Line)
		}
		if w.usesParameter(value) {
			w.facts.injects[seam] = true
		}
	}
}

// assignedRoot reduces an assignment target to the variable that holds the replaced
// value. A field, an index, or a pointer target belongs to the variable at its root. A
// selector on an import binding is that package's variable, and any other selector
// belongs to its operand, so a field of a same-package collaborator struct resolves to
// that package variable.
func (w *walker) assignedRoot(expr ast.Expr) ast.Expr {
	for {
		switch typed := expr.(type) {
		case *ast.ParenExpr:
			expr = typed.X
		case *ast.StarExpr:
			expr = typed.X
		case *ast.IndexExpr:
			expr = typed.X
		case *ast.SelectorExpr:
			if w.resolve(typed) != "" {
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
