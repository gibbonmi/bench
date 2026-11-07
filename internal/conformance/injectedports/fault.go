package injectedports

import (
	"fmt"
	"go/ast"
	"slices"
	"strings"
)

// FaultProducerMessage names a fault-port row whose port production supplies a value.
// Such a port has a real producer, so its row names the test that drives that producer.
const FaultProducerMessage = "injected fault port registry row names a port that production supplies"

// fieldSet is one production store into a struct field: an element of a literal of a
// resolved struct, or an assignment to a selector whose struct the walk does not resolve.
// A keyed element and an assignment name the field; a positional element gives its index.
type fieldSet struct {
	strct, field string
	index        int
	value        operand
}

// faultDiags grades a fault-port row: a named port that every production caller passes as
// nil. Owner is a production function whose parameter takes the port. A consumer supplies
// a non-nil value at an injection site: it, or a test helper it reaches, calls a function
// of the package with a non-nil argument for a parameter of the port type, or assigns a
// package variable of that type.
//
// The port has a production producer when any production place can carry a value of it:
// a parameter of the port type that a production caller passes a value, a struct field of
// that type that production sets, or a package variable of that type that production sets.
func (c *seamCensus) faultDiags(row Row) []string {
	identity := row.Package + "." + row.Port
	owns := slices.Contains(c.params[row.Package+"."+row.Owner], row.Port)
	diags := rowEvidence(row, identity, owns, func(file, test string) (bool, bool) {
		return c.reachesControl(file, test, false, func(facts *funcFacts) bool { return c.suppliesFault(facts, row.Package, row.Port) })
	})
	trace := &faultTrace{census: c, pkg: row.Package, port: row.Port, forwarders: map[string]bool{}}
	if producers := trace.producers(); len(producers) > 0 {
		diags = append(diags, fmt.Sprintf("%s: %s is supplied a value by %s — name the test that drives its real producer instead", FaultProducerMessage, identity, strings.Join(producers, "; ")))
	}
	return diags
}

// suppliesFault reports whether one body supplies a non-nil value of port at an injection
// site of pkg.
func (c *seamCensus) suppliesFault(facts *funcFacts, pkg, port string) bool {
	for _, site := range facts.sites {
		if packageOf(site.callee) != pkg {
			continue
		}
		for index, name := range c.params[site.callee] {
			if name == port && index < len(site.args) && !isNil(site.args[index], site.scope) {
				return true
			}
		}
	}
	for variable := range facts.assigns {
		if packageOf(variable) == pkg && c.varTypes[variable] == port {
			return true
		}
	}
	return false
}

// isNil reports whether expr is the predeclared nil in scope. A nil scope reads the name
// alone, for a pass that does not track scopes.
func isNil(expr ast.Expr, scope *bodyScope) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil" && scope.owner("nil") == nil
}

// faultTrace decides, for the port of one row, whether production supplies a value of it.
// A production value is nil when it is the nil literal, a variable or field of the port
// type that production never sets, a parameter that every production caller passes as
// nil, a local whose every assignment is nil, or a function literal that composes only nil
// faults. A call that a test-support package makes is not production. A cycle is nil.
type faultTrace struct {
	census     *seamCensus
	pkg, port  string
	visiting   map[string]bool
	forwarders map[string]bool
}

// producers names each production place of the port's package that carries a value of it.
// A method parameter is named too, because the audit does not resolve the callers of a method.
func (t *faultTrace) producers() []string {
	c := t.census
	var found []string
	for _, id := range sortedKeys(c.params) {
		for index, name := range c.params[id] {
			if name != t.port || !strings.HasPrefix(id, t.pkg+".") {
				continue
			}
			if c.prod[id] == nil {
				found = append(found, "the method "+id+", whose callers the audit does not resolve")
				break
			}
			if t.visiting = map[string]bool{}; !t.nilAt(id, index) {
				found = append(found, "a production caller of "+id)
				break
			}
		}
	}
	for _, field := range sortedKeys(c.fieldTypes) {
		if t.visiting = map[string]bool{}; c.fieldTypes[field] == t.port && strings.HasPrefix(field, t.pkg+".") && !t.unsetField(field) {
			found = append(found, "a production store into the field "+field)
		}
	}
	for _, variable := range sortedKeys(c.varTypes) {
		if c.varTypes[variable] == t.port && packageOf(variable) == t.pkg && !t.unset(variable) {
			found = append(found, "a production value of the variable "+variable)
		}
	}
	return found
}

// nilAt reports whether every production call of callee passes nil at parameter index.
func (t *faultTrace) nilAt(callee string, index int) bool {
	if t.visiting[callee] {
		return true
	}
	t.visiting[callee] = true
	for _, facts := range t.census.production {
		for _, site := range facts.sites {
			if site.callee == callee && (index >= len(site.args) || !t.nilValue(facts, site.scope, site.args[index])) {
				return false
			}
		}
	}
	return true
}

// nilValue reports whether one production expression, written in scope of a body, is nil.
func (t *faultTrace) nilValue(facts *funcFacts, scope *bodyScope, expr ast.Expr) bool {
	switch typed := ast.Unparen(expr).(type) {
	case *ast.Ident:
		owner := scope.owner(typed.Name)
		switch {
		case isNil(typed, scope):
			return true
		case owner == nil:
			return t.unset(facts.ctx.pkg + "." + typed.Name)
		case owner == facts.ctx.top:
			caller := facts.ctx.pkg + "." + facts.ctx.owner
			param := slices.Index(facts.ctx.params, typed.Name)
			return param >= 0 && t.census.prod[caller] != nil && t.nilAt(caller, param)
		}
		return t.nilLocal(facts, owner, typed.Name)
	case *ast.SelectorExpr:
		if variable := resolveName(facts.ctx.file, facts.ctx.pkg, scope, typed); variable != "" {
			return t.unset(variable)
		}
		fields := 0
		for _, field := range sortedKeys(t.census.fieldTypes) {
			if t.census.fieldTypes[field] != t.port || !strings.HasPrefix(field, t.pkg+".") || !strings.HasSuffix(field, "."+typed.Sel.Name) {
				continue
			}
			if !t.unsetField(field) {
				return false
			}
			fields++
		}
		return fields > 0
	case *ast.FuncLit:
		only, relayed := returnsOnly(typed.Type, typed.Body, func(call *ast.CallExpr) bool {
			return t.relays(facts, call, func(scope *bodyScope, arg ast.Expr) bool { return t.nilValue(facts, scope, arg) })
		})
		return only && relayed
	}
	return false
}

// nilLocal reports whether every value assigned to a local of the body is nil. A local
// with no recorded value, such as a range variable, is not.
func (t *faultTrace) nilLocal(facts *funcFacts, owner *bodyScope, name string) bool {
	key := fmt.Sprintf("%p.%s", owner, name)
	if t.visiting[key] {
		return true
	}
	t.visiting[key] = true
	values := facts.locals[owner][name]
	for _, value := range values {
		if !t.nilValue(facts, value.scope, value.expr) {
			return false
		}
	}
	return len(values) > 0
}

// unset reports whether a package variable is declared with a type and no value, and no
// production body assigns it.
func (t *faultTrace) unset(variable string) bool {
	if _, typed := t.census.varTypes[variable]; !typed || t.census.initialized[variable] {
		return false
	}
	for _, facts := range t.census.production {
		if _, assigned := facts.assigns[variable]; assigned {
			return false
		}
	}
	return true
}

// unsetField reports whether production stores only nil into a field, written
// "pkg.Struct.field". A keyed or positional element of a literal of its struct stores into
// it, and so does an assignment to any selector of the same name.
func (t *faultTrace) unsetField(field string) bool {
	if t.visiting[field] {
		return true
	}
	t.visiting[field] = true
	strct, name := seamPackage(field)
	for _, facts := range t.census.production {
		for _, set := range facts.fieldSets {
			stored := set.field
			if fields := t.census.structFields[set.strct]; stored == "" && set.index < len(fields) {
				stored = fields[set.index]
			}
			if stored == name && (set.strct == "" || set.strct == strct) && !t.nilValue(facts, set.value.scope, set.value.expr) {
				return false
			}
		}
	}
	return true
}

// relays reports whether a call returns only a fault that it is given: it calls a value
// that carries no fault, or it calls a forwarder of the port's package and every port
// argument of that call carries no fault.
func (t *faultTrace) relays(facts *funcFacts, call *ast.CallExpr, carriesNil func(*bodyScope, ast.Expr) bool) bool {
	scope := facts.callScopes[call]
	if scope == nil {
		return false
	}
	if carriesNil(scope, call.Fun) {
		return true
	}
	callee := resolveName(facts.ctx.file, facts.ctx.pkg, scope, call.Fun)
	if callee == "" || packageOf(callee) != t.pkg || !t.forwards(callee) {
		return false
	}
	ports := 0
	for index, name := range t.census.params[callee] {
		if name != t.port {
			continue
		}
		if index >= len(call.Args) || !carriesNil(scope, call.Args[index]) {
			return false
		}
		ports++
	}
	return ports > 0
}

// forwards reports whether a production function returns only the faults that its port
// parameters carry, so a call of it carries no fault when those arguments carry none.
func (t *faultTrace) forwards(id string) bool {
	if done, seen := t.forwarders[id]; seen {
		return done
	}
	entry := t.census.prod[id]
	if entry == nil {
		return false
	}
	t.forwarders[id] = true
	facts, params := entry.facts, t.census.params[id]
	own := func(scope *bodyScope, arg ast.Expr) bool {
		ident, ok := ast.Unparen(arg).(*ast.Ident)
		if !ok {
			return false
		}
		index := slices.Index(facts.ctx.params, ident.Name)
		return isNil(ident, scope) || scope.owner(ident.Name) == facts.ctx.top && index >= 0 && index < len(params) && params[index] == t.port
	}
	done, _ := returnsOnly(entry.decl.Type, entry.decl.Body, func(call *ast.CallExpr) bool { return t.relays(facts, call, own) })
	t.forwarders[id] = done
	return done
}

// returnsOnly reports whether every value that a function body returns is nil or a call
// that carries no fault, and whether any returned value is such a call. A returned local
// counts when every assignment to it in the body is such a call. A function with named
// results can return them bare, so it never counts.
//
// A function literal composes nil faults only when it relays at least one of them. A
// literal that returns only nil is still a value that production supplies.
func returnsOnly(fn *ast.FuncType, body *ast.BlockStmt, carries func(*ast.CallExpr) bool) (only, relayed bool) {
	if body == nil || fn.Results != nil && len(fn.Results.List) > 0 && len(fn.Results.List[0].Names) > 0 {
		return false, false
	}
	values := map[string][]ast.Expr{}
	var results []ast.Expr
	ast.Inspect(body, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.FuncLit:
			return false
		case *ast.AssignStmt:
			for i, target := range typed.Lhs {
				if ident, ok := target.(*ast.Ident); ok {
					values[ident.Name] = append(values[ident.Name], assignedValue(typed, i))
				}
			}
		case *ast.ReturnStmt:
			results = append(results, typed.Results...)
		}
		return true
	})
	for _, result := range results {
		returned := []ast.Expr{result}
		if ident, ok := result.(*ast.Ident); ok && !isNil(ident, nil) {
			returned = values[ident.Name]
		}
		if len(returned) == 0 {
			return false, false
		}
		for _, value := range returned {
			call, called := ast.Unparen(value).(*ast.CallExpr)
			switch {
			case isNil(value, nil):
			case called && carries(call):
				relayed = true
			default:
				return false, false
			}
		}
	}
	return true, relayed
}
