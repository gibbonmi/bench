package injectedports

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"sort"
)

// FaultProducerMessage names a fault-port row whose owner production supplies a value.
// Such a port has a real producer, so its row names the test that drives that producer.
const FaultProducerMessage = "injected fault port registry row names an owner that a production caller supplies"

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

// faultDiags grades a fault-port row: a named port that every production caller of its
// owner passes as nil. Owner is the production function whose parameter takes the port.
// Every production call of the owner passes nil there, a variable of the port type that
// production declares without a value and never assigns, or a parameter of its own caller
// that is nil by the same rule. A consumer supplies a non-nil value at an injection site:
// it, or a test helper it reaches, calls a function of the package with a non-nil argument
// for a parameter of the port type, or assigns a package variable of that type.
func (c *seamCensus) faultDiags(row Row) []string {
	identity := row.Package + "." + row.Port
	owner := row.Package + "." + row.Owner
	index := slices.Index(c.params[owner], row.Port)
	diags := rowEvidence(row, identity, index >= 0, func(file, test string) (bool, bool) {
		return c.reachesControl(file, test, false, func(facts *funcFacts) bool { return c.suppliesFault(facts, row.Package, row.Port) })
	})
	if index >= 0 && !c.nilAt(owner, index, map[string]bool{}) {
		diags = append(diags, fmt.Sprintf("%s: %s is supplied a value by a production caller of %s — name the test that drives its real producer instead", FaultProducerMessage, identity, row.Owner))
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

// nilAt reports whether every production call of callee passes nil at parameter index. A
// call that a test-support package makes is not production. A forwarding cycle is nil.
func (c *seamCensus) nilAt(callee string, index int, visiting map[string]bool) bool {
	if visiting[callee] {
		return true
	}
	visiting[callee] = true
	for _, id := range sortedKeys(c.bodies) {
		facts := c.bodies[id]
		if c.support[facts.ctx.pkg] {
			continue
		}
		for _, site := range facts.sites {
			if site.callee == callee && (index >= len(site.args) || !c.nilArgument(facts, site.scope, site.args[index], visiting)) {
				return false
			}
		}
	}
	return true
}

// nilArgument reports whether one production argument is nil: the nil literal, a package
// variable declared without a value that no production body assigns, or a parameter of the
// calling function that every production call of it passes as nil.
func (c *seamCensus) nilArgument(facts *funcFacts, scope *bodyScope, arg ast.Expr, visiting map[string]bool) bool {
	ident, ok := arg.(*ast.Ident)
	if !ok {
		return false
	}
	switch owner := scope.owner(ident.Name); {
	case isNil(arg, scope):
		return true
	case owner == nil:
		variable := facts.ctx.pkg + "." + ident.Name
		if _, typed := c.varTypes[variable]; !typed {
			return false
		}
		for _, body := range c.bodies {
			if _, assigned := body.assigns[variable]; assigned {
				return false
			}
		}
		return true
	case owner == facts.ctx.top:
		caller := facts.ctx.pkg + "." + facts.ctx.owner
		param := slices.Index(facts.ctx.params, ident.Name)
		return param >= 0 && c.prod[caller] != nil && c.nilAt(caller, param, visiting)
	}
	return false
}

// isNil reports whether expr is the predeclared nil in scope.
func isNil(expr ast.Expr, scope *bodyScope) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil" && scope.owner("nil") == nil
}
