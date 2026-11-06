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
// identity that a production reader reads and that test source sets, through Setenv or
// through a "NAME=value" entry. Test-support packages contribute setters, never readers
// or branches.
type controlCensus struct {
	runtime map[string]string              // branch identity -> the probe or predicate it calls
	readers map[string]map[string]string   // hook name -> reader body identity -> site
	sets    map[*funcFacts]map[string]bool // test body -> environment names it sets
	hooks   map[string]bool
	diags   []string
}

func (c *seamCensus) deriveControls(families map[string][]string) *controlCensus {
	c.families = families
	cc := &controlCensus{runtime: map[string]string{}, readers: map[string]map[string]string{}, sets: map[*funcFacts]map[string]bool{}, hooks: map[string]bool{}}
	set := map[string]bool{}
	for _, facts := range c.testBodies {
		cc.sets[facts] = c.namesSet(facts)
		for name := range cc.sets[facts] {
			set[name] = true
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
			for name := range c.namesSet(facts) {
				set[name] = true
			}
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
	for name := range cc.readers {
		if set[name] {
			cc.hooks[name] = true
		}
	}
	return cc
}

// namesSet returns the environment names that one body sets.
func (c *seamCensus) namesSet(facts *funcFacts) map[string]bool {
	names := map[string]bool{}
	for _, op := range facts.envSets {
		values, _ := c.values(facts.ctx, op.scope, op.expr, 0)
		for _, value := range values {
			if envName.MatchString(value) {
				names[value] = true
			}
		}
	}
	for _, op := range facts.entries {
		for _, name := range c.entryNames(facts.ctx, op) {
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
// the test bodies that assign a test variable it reads, such as a check table. When
// throughProduction holds, it also follows the production calls along the way.
func (c *seamCensus) reachesControl(file, test string, throughProduction bool, hit func(*funcFacts) bool) (declared, reached bool) {
	scope, known := c.testFiles[file]
	entry := c.tests[scope+"."+test]
	if !known || entry == nil || entry.file.rel != file || !strings.HasPrefix(test, "Test") {
		return false, false
	}
	visited := map[*funcFacts]bool{entry.facts: true}
	queue := []*funcFacts{entry.facts}
	enqueue := func(facts *funcFacts) {
		if facts != nil && !visited[facts] {
			visited[facts] = true
			queue = append(queue, facts)
		}
	}
	for len(queue) > 0 {
		facts := queue[0]
		queue = queue[1:]
		if hit(facts) {
			return true, true
		}
		for callee := range facts.calls {
			if next := c.tests[callee]; next != nil {
				enqueue(next.facts)
			} else if throughProduction {
				enqueue(c.bodies[callee])
			}
		}
		if c.bodies[facts.ctx.pkg+"."+facts.ctx.owner] == facts {
			continue
		}
		for read := range facts.reads {
			if next := c.tests[read]; next != nil {
				enqueue(next.facts)
			}
			for _, assigner := range c.assigners[read] {
				enqueue(assigner)
			}
		}
	}
	return true, false
}
