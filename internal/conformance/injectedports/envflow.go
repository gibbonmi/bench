package injectedports

import (
	"go/ast"
	"strings"
)

// A "NAME=value" string sets an environment hook only where its value flows into an
// environment: the Env field of a process, such as exec.Cmd or os.ProcAttr. The walker
// records, for each string, local, field read, and call result, the one place its value
// flows to. The census then follows those flows across locals, struct fields, helper
// parameters, and helper results to the Env field. A string in an assertion, a comparison,
// or any other position flows nowhere and sets nothing.

// flowKind names one kind of place that a value can flow to.
type flowKind byte

const (
	flowEnv    flowKind = iota + 1 // the Env field of a process
	flowLocal                      // a local variable or parameter, by its declaring scope
	flowParam                      // a parameter of a callee, by its position
	flowReturn                     // the result of a function body
	flowField                      // a struct field of a package, by its name
)

// flowNode is one place a value can flow to.
type flowNode struct {
	kind  flowKind
	scope *bodyScope
	name  string
	index int
}

// flow is one recorded edge: the value at from flows to to.
type flow struct{ from, to flowNode }

// envField is the process field whose value is the environment of a child process.
const envField = "Env"

// callee names the function that a call expression calls. A method call on a value names
// every method of that name in the walked package, since the walker does not know types.
func (w *walker) callee(fun ast.Expr) string {
	if callee := w.resolve(fun); callee != "" {
		return callee
	}
	if sel, ok := fun.(*ast.SelectorExpr); ok {
		return methodKey(w.pkg, sel.Sel.Name)
	}
	return ""
}

func methodKey(pkg, name string) string { return pkg + ".(*)." + name }

// bodyKey names a body the way a call of it names its callee.
func bodyKey(ctx *funcContext) string {
	if _, method, found := strings.Cut(ctx.owner, "."); found {
		return methodKey(ctx.pkg, method)
	}
	return ctx.pkg + "." + ctx.owner
}

// field names the struct field name of the walked package. Every Env field is a process
// environment.
func (w *walker) field(name string) flowNode {
	if name == envField {
		return flowNode{kind: flowEnv}
	}
	return flowNode{kind: flowField, name: w.pkg + "." + name}
}

// flowFrom records that the value of the last node of stack, read from from, flows to
// the place its enclosing expression puts it.
func (w *walker) flowFrom(from flowNode, stack []ast.Node) {
	if to, ok := w.carrier(stack); ok && from.kind != flowEnv {
		w.facts.flows = append(w.facts.flows, flow{from: from, to: to})
	}
}

// carrier returns the place that the value of the last node of stack flows to. A
// composite literal, an append, a slice, or parentheses pass the value on to their own
// place. A keyed element, an assignment, a declaration, a range, a call argument, or a
// return stops at its place. Every other position is no place.
func (w *walker) carrier(stack []ast.Node) (flowNode, bool) {
	for i := len(stack) - 1; i > 0; i-- {
		node := stack[i]
		switch parent := stack[i-1].(type) {
		case *ast.ParenExpr, *ast.CompositeLit:
			continue
		case *ast.SliceExpr:
			if parent.X == node {
				continue
			}
		case *ast.KeyValueExpr:
			if key, ok := parent.Key.(*ast.Ident); ok && parent.Value == node {
				return w.field(key.Name), true
			}
		case *ast.CallExpr:
			if parent.Fun == node {
				return flowNode{}, false
			}
			if fun, ok := parent.Fun.(*ast.Ident); ok && fun.Name == "append" && w.scope.owner("append") == nil {
				continue
			}
			if callee := w.callee(parent.Fun); callee != "" {
				for index, arg := range parent.Args {
					if arg == node {
						return flowNode{kind: flowParam, name: callee, index: index}, true
					}
				}
			}
		case *ast.AssignStmt:
			for index, value := range parent.Rhs {
				if value == node && len(parent.Lhs) == len(parent.Rhs) {
					return w.target(parent.Lhs[index])
				}
			}
		case *ast.ValueSpec:
			for index, value := range parent.Values {
				if value == node && len(parent.Names) == len(parent.Values) {
					return w.target(parent.Names[index])
				}
			}
		case *ast.RangeStmt:
			if parent.X == node && parent.Value != nil {
				return w.target(parent.Value)
			}
		case *ast.ReturnStmt:
			for _, outer := range stack[:i-1] {
				if _, closure := outer.(*ast.FuncLit); closure {
					return flowNode{}, false
				}
			}
			return flowNode{kind: flowReturn, name: bodyKey(w.facts.ctx)}, true
		}
		return flowNode{}, false
	}
	return flowNode{}, false
}

// target returns the place that an assignment or a declaration stores into: a local, or
// a struct field of a value that is not a package-qualified name.
func (w *walker) target(expr ast.Expr) (flowNode, bool) {
	switch typed := expr.(type) {
	case *ast.Ident:
		if owner := w.scope.owner(typed.Name); owner != nil {
			return flowNode{kind: flowLocal, scope: owner, name: typed.Name}, true
		}
	case *ast.SelectorExpr:
		if w.resolve(typed) == "" {
			return w.field(typed.Sel.Name), true
		}
	}
	return flowNode{}, false
}

// allBodies returns every production body, test function, and test method of the tree.
func (c *seamCensus) allBodies() []*funcFacts {
	var all []*funcFacts
	for _, id := range sortedKeys(c.bodies) {
		all = append(all, c.bodies[id])
	}
	for _, id := range sortedKeys(c.tests) {
		all = append(all, c.tests[id].facts)
	}
	for _, key := range sortedKeys(c.testMethods) {
		all = append(all, c.testMethods[key]...)
	}
	return all
}

// envReach returns whether the value at a place reaches a process environment. A callee
// parameter is the callee's own parameter local, and a variadic tail argument is its last
// parameter.
func (c *seamCensus) envReach() func(flowNode) bool {
	all := c.allBodies()
	// A parameter slot flows into the parameter local of each body of that callee name.
	// A variadic body takes every tail slot into its last parameter.
	sources := map[flowNode][]flowNode{}
	tails := map[string]int{}
	for _, facts := range all {
		key := bodyKey(facts.ctx)
		for index, name := range facts.ctx.params {
			local := flowNode{kind: flowLocal, scope: facts.ctx.top, name: name}
			sources[local] = append(sources[local], flowNode{kind: flowParam, name: key, index: index})
		}
		if facts.ctx.variadic && len(facts.ctx.params) > 0 {
			tails[key] = len(facts.ctx.params) - 1
		}
	}
	normal := func(node flowNode) flowNode {
		if tail, variadic := tails[node.name]; node.kind == flowParam && variadic && node.index > tail {
			node.index = tail
		}
		return node
	}
	for _, facts := range all {
		for _, edge := range facts.flows {
			to := normal(edge.to)
			sources[to] = append(sources[to], normal(edge.from))
		}
	}
	reach := map[flowNode]bool{{kind: flowEnv}: true}
	queue := []flowNode{{kind: flowEnv}}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, from := range sources[node] {
			if !reach[from] {
				reach[from] = true
				queue = append(queue, from)
			}
		}
	}
	return func(node flowNode) bool { return reach[normal(node)] }
}
