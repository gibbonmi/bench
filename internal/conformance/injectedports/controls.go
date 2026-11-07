package injectedports

import (
	"go/ast"
	"go/token"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// These prefixes name the runtime-branch and environment-hook failures.
const (
	UnregisteredRuntimeMessage = "test runtime branch has no registry row"
	UnregisteredHookMessage    = "environment hook has no registry row"
	UnresolvedHookMessage      = "environment reader identity does not resolve"
)

// testingProbe is the call that reports a running test binary.
const testingProbe = "testing.Testing"

// maxValues bounds a constructed identity set, and maxDepth bounds the resolution of
// one operand through constants and callers. A construction past either is unresolved.
const (
	maxValues = 64
	maxDepth  = 8
)

var envName = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// controlCensus is the runtime-branch and environment-hook census of one graded root.
//
// A runtime branch is a production body that calls testing.Testing through any import
// binding, or calls a production predicate that forwards to it. An environment hook is an
// identity that test source controls at a production reader: a "NAME=value" entry that
// flows into a process environment controls every reader of the name, and a Setenv in a
// test's own process controls each reader that the same test reaches through production
// calls. Test-support packages contribute setters, never readers or branches.
type controlCensus struct {
	runtime  map[string]string              // branch identity -> the probe or predicate it calls
	readers  map[string]map[string]string   // hook name -> controlled reader body identity -> site
	setenvs  map[*funcFacts]map[string]bool // setter body -> names it sets in its own process
	children map[*funcFacts]map[string]bool // setter body -> names it sets for a child process
	hooks    map[string]bool
	diags    []string
}

func (c *seamCensus) deriveControls(families map[string][]string) *controlCensus {
	c.families = families
	cc := &controlCensus{runtime: map[string]string{}, readers: map[string]map[string]string{}, setenvs: map[*funcFacts]map[string]bool{}, children: map[*funcFacts]map[string]bool{}, hooks: map[string]bool{}}
	reach := c.envReach()
	child := map[string]bool{}
	setters := append([]*funcFacts{}, c.testBodies...)
	for _, key := range sortedKeys(c.testMethods) {
		setters = append(setters, c.testMethods[key]...)
	}
	for _, id := range sortedKeys(c.bodies) {
		if c.support[c.bodies[id].ctx.pkg] {
			setters = append(setters, c.bodies[id])
		}
	}
	for _, facts := range setters {
		cc.setenvs[facts], cc.children[facts] = c.setenvNames(facts), c.childNames(facts, reach)
		for name := range cc.children[facts] {
			child[name] = true
		}
	}
	predicates := map[string]bool{}
	for changed := true; changed; {
		changed = false
		for id, facts := range c.bodies {
			if !predicates[id] && (facts.forward == testingProbe || predicates[facts.forward]) {
				predicates[id], changed = true, true
			}
		}
	}
	for _, id := range sortedKeys(c.bodies) {
		facts := c.bodies[id]
		if c.support[facts.ctx.pkg] {
			continue
		}
		for _, callee := range sortedKeys(facts.calls) {
			if _, found := cc.runtime[id]; !found && !predicates[id] && (callee == testingProbe || predicates[callee]) {
				cc.runtime[id] = callee
			}
		}
		for _, read := range facts.envReads {
			names, ok := c.values(facts.ctx, read.scope, read.expr, 0)
			if !ok {
				cc.diags = append(cc.diags, UnresolvedHookMessage+": "+id+" reads the environment at "+read.site+
					" — name it with constants, bounded string construction, constant caller arguments, or a family the policy supplies")
				continue
			}
			for _, name := range names {
				if cc.readers[name] == nil {
					cc.readers[name] = map[string]string{}
				}
				if _, seen := cc.readers[name][id]; !seen {
					cc.readers[name][id] = read.site
				}
			}
		}
	}
	controlled := c.inProcessReaders(cc)
	for name, byReader := range cc.readers {
		for id := range byReader {
			if !child[name] && !controlled[name][id] {
				delete(byReader, id)
			}
		}
		if len(byReader) > 0 {
			cc.hooks[name] = true
		}
	}
	return cc
}

// inProcessReaders returns, for each name, the readers that some top-level test both sets
// in its own process, itself or through its test helpers and test-support calls, and
// reaches through production calls.
func (c *seamCensus) inProcessReaders(cc *controlCensus) map[string]map[string]bool {
	callers := map[*funcFacts][]*funcFacts{}
	always := func(*funcFacts) bool { return true }
	for _, facts := range c.allBodies() {
		for _, next := range c.successors(facts, always) {
			callers[next] = append(callers[next], facts)
		}
	}
	reachers := map[string]map[*funcFacts]bool{}
	reachersOf := func(reader string) map[*funcFacts]bool {
		if found, done := reachers[reader]; done {
			return found
		}
		found := map[*funcFacts]bool{c.bodies[reader]: true}
		queue := []*funcFacts{c.bodies[reader]}
		for len(queue) > 0 {
			facts := queue[0]
			queue = queue[1:]
			for _, caller := range callers[facts] {
				if !found[caller] {
					found[caller] = true
					queue = append(queue, caller)
				}
			}
		}
		reachers[reader] = found
		return found
	}
	support := func(facts *funcFacts) bool { return c.support[facts.ctx.pkg] }
	controlled := map[string]map[string]bool{}
	for _, id := range sortedKeys(c.tests) {
		entry := c.tests[id]
		if !strings.HasPrefix(entry.facts.ctx.owner, "Test") {
			continue
		}
		visited := map[*funcFacts]bool{entry.facts: true}
		queue := []*funcFacts{entry.facts}
		for len(queue) > 0 {
			facts := queue[0]
			queue = queue[1:]
			for name := range cc.setenvs[facts] {
				for reader := range cc.readers[name] {
					if reachersOf(reader)[entry.facts] {
						if controlled[name] == nil {
							controlled[name] = map[string]bool{}
						}
						controlled[name][reader] = true
					}
				}
			}
			for _, next := range c.successors(facts, support) {
				if !visited[next] {
					visited[next] = true
					queue = append(queue, next)
				}
			}
		}
	}
	return controlled
}

// setenvNames returns the environment names that one body sets through Setenv.
func (c *seamCensus) setenvNames(facts *funcFacts) map[string]bool {
	names := map[string]bool{}
	for _, op := range facts.envSets {
		values, _ := c.values(facts.ctx, op.scope, op.expr, 0)
		for _, value := range values {
			if envName.MatchString(value) {
				names[value] = true
			}
		}
	}
	return names
}

// childNames returns the environment names of the "NAME=value" entries of one body whose
// value reaches a process environment.
func (c *seamCensus) childNames(facts *funcFacts, reach func(flowNode) bool) map[string]bool {
	names := map[string]bool{}
	for _, entry := range facts.entries {
		if !reach(entry.to) {
			continue
		}
		for _, name := range c.entryNames(facts.ctx, entry.operand) {
			names[name] = true
		}
	}
	return names
}

// entryNames resolves the leading operands of a concatenation until each value holds
// "=", and returns the environment names before it. A value part may stay dynamic.
func (c *seamCensus) entryNames(ctx *funcContext, op operand) []string {
	prefixes := []string{""}
	for _, part := range concatParts(op.expr) {
		values, ok := c.values(ctx, op.scope, part, 0)
		if !ok || len(prefixes)*len(values) > maxValues {
			return nil
		}
		prefixes = product(prefixes, values)
		if !slices.ContainsFunc(prefixes, func(value string) bool { return !strings.Contains(value, "=") }) {
			break
		}
	}
	var names []string
	for _, value := range prefixes {
		name, _, found := strings.Cut(value, "=")
		if !found || !envName.MatchString(name) {
			return nil
		}
		names = append(names, name)
	}
	return names
}

func concatParts(expr ast.Expr) []ast.Expr {
	if binary, ok := expr.(*ast.BinaryExpr); ok && binary.Op == token.ADD {
		return append(concatParts(binary.X), concatParts(binary.Y)...)
	}
	return []ast.Expr{expr}
}

func product(left, right []string) []string {
	out := make([]string, 0, len(left)*len(right))
	for _, l := range left {
		for _, r := range right {
			out = append(out, l+r)
		}
	}
	return out
}

// values resolves a string operand to the set of values it can take: a literal, a
// constant, a concatenation, a case fold, a parameter through every production caller,
// or a local name whose family the policy supplies. A supplied empty family resolves to
// no values; a name with no supplied family and no other source is unresolved.
func (c *seamCensus) values(ctx *funcContext, scope *bodyScope, expr ast.Expr, depth int) ([]string, bool) {
	if depth > maxDepth {
		return nil, false
	}
	switch typed := expr.(type) {
	case *ast.BasicLit:
		value, err := strconv.Unquote(typed.Value)
		return []string{value}, typed.Kind == token.STRING && err == nil
	case *ast.ParenExpr:
		return c.values(ctx, scope, typed.X, depth+1)
	case *ast.BinaryExpr:
		left, okLeft := c.values(ctx, scope, typed.X, depth+1)
		right, okRight := c.values(ctx, scope, typed.Y, depth+1)
		if typed.Op != token.ADD || !okLeft || !okRight || len(left)*len(right) > maxValues {
			return nil, false
		}
		return product(left, right), true
	case *ast.CallExpr:
		fold := map[string]func(string) string{"strings.ToUpper": strings.ToUpper, "strings.ToLower": strings.ToLower}[resolveName(ctx.file, ctx.pkg, scope, typed.Fun)]
		if fold == nil || len(typed.Args) != 1 {
			return nil, false
		}
		values, ok := c.values(ctx, scope, typed.Args[0], depth+1)
		folded := make([]string, len(values))
		for i, value := range values {
			folded[i] = fold(value)
		}
		return folded, ok
	case *ast.Ident:
		if owner := scope.owner(typed.Name); owner != nil {
			if family, supplied := c.families[ctx.pkg+"."+ctx.owner+"."+typed.Name]; supplied {
				return slices.Clone(family), true
			}
			if index := slices.Index(ctx.params, typed.Name); owner == ctx.top && index >= 0 {
				return c.argumentValues(ctx, index, depth+1)
			}
			return nil, false
		}
		return c.constValues(ctx.pkg+"."+typed.Name, depth+1)
	case *ast.SelectorExpr:
		if id := resolveName(ctx.file, ctx.pkg, scope, typed); id != "" {
			return c.constValues(id, depth+1)
		}
	}
	return nil, false
}

func (c *seamCensus) constValues(id string, depth int) ([]string, bool) {
	declared, ok := c.consts[id]
	if !ok {
		return nil, false
	}
	return c.values(declared.ctx, declared.ctx.top, declared.expr, depth)
}

// argumentValues resolves parameter index of the body in ctx through every production
// call of it. A body that no production call names is unresolved, never empty.
func (c *seamCensus) argumentValues(ctx *funcContext, index, depth int) ([]string, bool) {
	target := ctx.pkg + "." + ctx.owner
	var out []string
	found := false
	for _, id := range sortedKeys(c.bodies) {
		facts := c.bodies[id]
		for _, site := range facts.sites {
			if site.callee != target {
				continue
			}
			if index >= len(site.args) {
				return nil, false
			}
			values, ok := c.values(facts.ctx, site.scope, site.args[index], depth)
			if !ok || len(out)+len(values) > maxValues {
				return nil, false
			}
			out, found = append(out, values...), true
		}
	}
	return out, found
}

// reachesControl reports whether the top-level test in file reaches a body that hit
// accepts. It follows the test's own helpers, the test functions it names as values, and
// the test bodies that assign a test variable it reads, such as a check table, and the
// test methods it calls. When throughProduction holds, it also follows the production
// calls along the way.
func (c *seamCensus) reachesControl(file, test string, throughProduction bool, hit func(*funcFacts) bool) (declared, reached bool) {
	scope, known := c.testFiles[file]
	entry := c.tests[scope+"."+test]
	if !known || entry == nil || entry.file.rel != file || !strings.HasPrefix(test, "Test") {
		return false, false
	}
	visited := map[*funcFacts]bool{entry.facts: true}
	queue := []*funcFacts{entry.facts}
	production := func(*funcFacts) bool { return throughProduction }
	for len(queue) > 0 {
		facts := queue[0]
		queue = queue[1:]
		if hit(facts) {
			return true, true
		}
		for _, next := range c.successors(facts, production) {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return true, false
}

// successors returns the bodies that one body hands control to: the test functions and
// test methods it calls, the production bodies it calls that production admits, and,
// for a test body, the test functions it names and the test bodies that assign a test
// variable it reads. A method call names each method of that name in the caller's package.
func (c *seamCensus) successors(facts *funcFacts, production func(*funcFacts) bool) []*funcFacts {
	var next []*funcFacts
	for callee := range facts.calls {
		if entry := c.tests[callee]; entry != nil {
			next = append(next, entry.facts)
			continue
		}
		next = append(next, c.testMethods[callee]...)
		for _, body := range append([]*funcFacts{c.bodies[callee]}, c.methods[callee]...) {
			if body != nil && production(body) {
				next = append(next, body)
			}
		}
	}
	if c.bodies[facts.ctx.pkg+"."+facts.ctx.owner] == facts {
		return next
	}
	for read := range facts.reads {
		if entry := c.tests[read]; entry != nil {
			next = append(next, entry.facts)
		}
		next = append(next, c.assigners[read]...)
	}
	return next
}
