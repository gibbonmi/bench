package injectedports

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"
)

// funcEntry is one top-level function with its body facts.
type funcEntry struct {
	file  *sourceFile
	facts *funcFacts
}

// seamCensus is the resolved variable and setter census of one graded root.
//
// A variable seam is a production package variable that a production path reads and that
// test source or a setter replaces. A setter is a production function that assigns a
// package variable a value built from its own parameter, or that forwards to a setter.
// Its canonical seams are the variables it finally replaces. Ordinary mutable data that
// production code maintains from its own state is neither.
type seamCensus struct {
	vars      map[string]bool            // production variable identities
	prod      map[string]*funcEntry      // production function identities (pkg.name)
	readers   map[string]map[string]bool // variable identity -> production owners reading it
	tests     map[string]*funcEntry      // test function identities (test pkg.name)
	testFiles map[string]string          // test file path -> its test package identity
	setters   map[string][]string        // setter identity -> canonical variable identities
	variables map[string]string          // variable seam identity -> first replacement site
	// The runtime and environment census reads these: every production body, methods and
	// package initializers included, the constants of each package scope, and the imports
	// that decide which packages are test support.
	bodies  map[string]*funcFacts      // production body identity (pkg.owner) -> facts
	consts  map[string]constant        // constant identity -> its declaration
	imports map[string]map[string]bool // package -> packages its production files import
	tested  map[string]bool            // packages that some test file imports
	mains   map[string]bool            // packages that declare a command entry
	support map[string]bool            // packages that only tests and test support reach
	// families are the policy's members of a local name, keyed "pkg.owner.name".
	families map[string][]string
	// assigners maps a test-package variable to the test bodies, init included, that
	// assign it, so a test that reads a table reaches the functions the table names.
	assigners  map[string][]*funcFacts
	testBodies []*funcFacts
	// testMethods maps a method name of a test package to the method bodies of its test
	// files, so a consumer and an environment entry can reach through a fixture's method.
	testMethods map[string][]*funcFacts
	// methods maps a method name of a package to its production method bodies.
	methods map[string][]*funcFacts
}

// constant is one constant declaration with the context its value resolves in.
type constant struct {
	expr ast.Expr
	ctx  *funcContext
}

// deriveSeams resolves every function of the tree, then derives the setters and the
// variable seams. A test file of an external test package resolves its own names to that
// test package, never to the production package of the same directory.
func deriveSeams(tree *sourceTree) *seamCensus {
	census := &seamCensus{
		vars: map[string]bool{}, prod: map[string]*funcEntry{}, readers: map[string]map[string]bool{},
		tests: map[string]*funcEntry{}, testFiles: map[string]string{}, setters: map[string][]string{},
		variables: map[string]string{}, bodies: map[string]*funcFacts{}, consts: map[string]constant{},
		imports: map[string]map[string]bool{}, tested: map[string]bool{}, mains: map[string]bool{},
		assigners: map[string][]*funcFacts{}, methods: map[string][]*funcFacts{}, testMethods: map[string][]*funcFacts{},
	}
	for _, dir := range sortedKeys(tree.packages) {
		census.indexPackage(tree.fset, tree.packages[dir])
	}
	for _, id := range sortedKeys(census.prod) {
		census.resolveSetter(id, map[string]bool{})
	}
	// A setter reads the variable it saves and restores, but it is not the production path
	// that the replacement controls, so it cannot own the seam.
	for _, id := range sortedKeys(census.setters) {
		pkg, name := seamPackage(id)
		for _, seam := range census.setters[id] {
			if pkg != packageOf(seam) {
				continue
			}
			delete(census.readers[seam], name)
			if site, direct := census.prod[id].facts.assigns[seam]; direct {
				census.replaced(seam, site)
			}
		}
	}
	for _, id := range sortedKeys(census.tests) {
		for seam, site := range census.tests[id].facts.assigns {
			census.replaced(seam, site)
		}
	}
	for seam := range census.variables {
		if len(census.readers[seam]) == 0 {
			delete(census.variables, seam)
		}
	}
	census.support = census.testSupport()
	return census
}

// testSupport returns the packages that test source reaches and no production package
// reaches. A package is test support by its importers, not by its name or its testing
// import: a test imports it, it declares no command entry, and every production file that
// imports it belongs to test support too.
func (c *seamCensus) testSupport() map[string]bool {
	support := map[string]bool{}
	for pkg := range c.tested {
		if !c.mains[pkg] {
			support[pkg] = true
		}
	}
	for changed := true; changed; {
		changed = false
		for importer, imported := range c.imports {
			if support[importer] {
				continue
			}
			for pkg := range imported {
				if support[pkg] {
					delete(support, pkg)
					changed = true
				}
			}
		}
	}
	return support
}

// declareConsts indexes the constants of one declaration under the package scope pkg.
func (c *seamCensus) declareConsts(file *sourceFile, pkg string, gen *ast.GenDecl) {
	if gen.Tok != token.CONST {
		return
	}
	ctx := &funcContext{file: file, pkg: pkg, top: &bodyScope{names: map[string]bool{}}}
	for _, spec := range gen.Specs {
		value, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for i, name := range value.Names {
			if i < len(value.Values) {
				c.consts[pkg+"."+name.Name] = constant{expr: value.Values[i], ctx: ctx}
			}
		}
	}
}

// importsOf records the module packages that file imports.
func (c *seamCensus) importsOf(file *sourceFile, into map[string]bool) {
	for _, dir := range file.imports {
		if strings.HasPrefix(dir, "cmd/") || strings.HasPrefix(dir, "internal/") {
			into[dir] = true
		}
	}
}

func (c *seamCensus) indexPackage(fset *token.FileSet, pkg *sourcePackage) {
	clause := ""
	read := func(owner string, facts *funcFacts) {
		for seam := range facts.reads {
			if c.readers[seam] == nil {
				c.readers[seam] = map[string]bool{}
			}
			c.readers[seam][owner] = true
		}
	}
	prodImports := map[string]bool{}
	c.imports[pkg.rel] = prodImports
	for _, file := range pkg.prod {
		clause = file.clause
		c.mains[pkg.rel] = c.mains[pkg.rel] || clause == "main"
		c.importsOf(file, prodImports)
		for _, decl := range file.file.Decls {
			switch typed := decl.(type) {
			case *ast.GenDecl:
				c.declareConsts(file, pkg.rel, typed)
				for _, spec := range typed.Specs {
					value, ok := spec.(*ast.ValueSpec)
					if !ok || typed.Tok != token.VAR {
						continue
					}
					for i, name := range value.Names {
						c.vars[pkg.rel+"."+name.Name] = true
						if i < len(value.Values) {
							facts := exprFacts(fset, file, pkg.rel, value.Values[i])
							read(name.Name, facts)
							c.bodies[pkg.rel+"."+name.Name] = facts
						}
					}
				}
			case *ast.FuncDecl:
				facts := factsOf(fset, file, pkg.rel, typed)
				read(ownerName(typed), facts)
				if typed.Name.Name != "init" {
					c.bodies[pkg.rel+"."+ownerName(typed)] = facts
				}
				if typed.Recv != nil {
					key := methodKey(pkg.rel, typed.Name.Name)
					c.methods[key] = append(c.methods[key], facts)
				}
				if typed.Recv == nil && typed.Name.Name != "init" {
					c.prod[pkg.rel+"."+typed.Name.Name] = &funcEntry{file: file, facts: facts}
				}
			}
		}
	}
	for _, file := range pkg.tests {
		scope := pkg.rel
		if file.clause != clause && (clause != "" || strings.HasSuffix(file.clause, "_test")) {
			scope = pkg.rel + "_test"
		}
		c.testFiles[file.rel] = scope
		tested := map[string]bool{}
		c.importsOf(file, tested)
		delete(tested, pkg.rel)
		for dir := range tested {
			c.tested[dir] = true
		}
		for _, decl := range file.file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok {
				c.declareConsts(file, scope, gen)
			}
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv != nil {
				key := methodKey(scope, fn.Name.Name)
				c.testMethods[key] = append(c.testMethods[key], factsOf(fset, file, scope, fn))
			} else if ok {
				facts := factsOf(fset, file, scope, fn)
				c.tests[scope+"."+fn.Name.Name] = &funcEntry{file: file, facts: facts}
				c.testBodies = append(c.testBodies, facts)
				for variable := range facts.assigns {
					c.assigners[variable] = append(c.assigners[variable], facts)
				}
			}
		}
	}
}

// resolveSetter records id as a setter when it injects a parameter into a package
// variable, or when it forwards to a setter. A forwarding cycle resolves to nothing.
func (c *seamCensus) resolveSetter(id string, visiting map[string]bool) []string {
	if canonical, done := c.setters[id]; done {
		return canonical
	}
	entry := c.prod[id]
	if entry == nil || visiting[id] {
		return nil
	}
	visiting[id] = true
	var canonical []string
	if entry.facts.forward != "" {
		canonical = c.resolveSetter(entry.facts.forward, visiting)
	} else {
		for seam := range entry.facts.injects {
			if c.vars[seam] {
				canonical = append(canonical, seam)
			}
		}
		slices.Sort(canonical)
	}
	if len(canonical) > 0 {
		c.setters[id] = canonical
	}
	return canonical
}

// replaced records a replacement of a production variable. The first recorded site wins,
// and the callers record in a fixed order, so a diagnostic names a stable place.
func (c *seamCensus) replaced(seam, site string) {
	if _, seen := c.variables[seam]; c.vars[seam] && !seen {
		c.variables[seam] = site
	}
}

// reaches reports whether the top-level test in file installs a replacement of seam:
// it assigns the variable or calls a setter of it, directly or through a test helper of
// its own test package. declared is false when file declares no such test.
func (c *seamCensus) reaches(file, test, seam string) (declared, reached bool) {
	scope, known := c.testFiles[file]
	start := scope + "." + test
	if entry := c.tests[start]; !known || entry == nil || entry.file.rel != file || !strings.HasPrefix(test, "Test") {
		return false, false
	}
	visited := map[string]bool{start: true}
	queue := []string{start}
	for len(queue) > 0 {
		facts := c.tests[queue[0]].facts
		queue = queue[1:]
		if _, assigns := facts.assigns[seam]; assigns {
			return true, true
		}
		for callee := range facts.calls {
			if slices.Contains(c.setters[callee], seam) {
				return true, true
			}
			if c.tests[callee] != nil && !visited[callee] {
				visited[callee] = true
				queue = append(queue, callee)
			}
		}
	}
	return true, false
}

// ownerName names a production function or method the way an admission row cites it.
func ownerName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return fn.Name.Name
	}
	recv := fn.Recv.List[0].Type
	for {
		switch typed := recv.(type) {
		case *ast.StarExpr:
			recv = typed.X
		case *ast.IndexExpr:
			recv = typed.X
		case *ast.IndexListExpr:
			recv = typed.X
		case *ast.Ident:
			return typed.Name + "." + fn.Name.Name
		default:
			return fn.Name.Name
		}
	}
}

// seamPackage splits a seam identity into its package directory and its name.
func seamPackage(identity string) (string, string) {
	cut := strings.LastIndex(identity, ".")
	return identity[:cut], identity[cut+1:]
}

func packageOf(identity string) string {
	pkg, _ := seamPackage(identity)
	return pkg
}
