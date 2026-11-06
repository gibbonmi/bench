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
}

// deriveSeams resolves every function of the tree, then derives the setters and the
// variable seams. A test file of an external test package resolves its own names to that
// test package, never to the production package of the same directory.
func deriveSeams(tree *sourceTree) *seamCensus {
	census := &seamCensus{
		vars: map[string]bool{}, prod: map[string]*funcEntry{}, readers: map[string]map[string]bool{},
		tests: map[string]*funcEntry{}, testFiles: map[string]string{}, setters: map[string][]string{},
		variables: map[string]string{},
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
	return census
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
	for _, file := range pkg.prod {
		clause = file.clause
		for _, decl := range file.file.Decls {
			switch typed := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range typed.Specs {
					value, ok := spec.(*ast.ValueSpec)
					if !ok || typed.Tok != token.VAR {
						continue
					}
					for i, name := range value.Names {
						c.vars[pkg.rel+"."+name.Name] = true
						if i < len(value.Values) {
							read(name.Name, exprFacts(fset, file, pkg.rel, value.Values[i]))
						}
					}
				}
			case *ast.FuncDecl:
				facts := factsOf(fset, file, pkg.rel, typed)
				read(ownerName(typed), facts)
				if typed.Recv == nil && typed.Name.Name != "init" {
					c.prod[pkg.rel+"."+typed.Name.Name] = &funcEntry{file: file, facts: facts}
				}
			}
		}
	}
	for _, file := range pkg.tests {
		scope := pkg.rel
		if file.clause != clause {
			scope = pkg.rel + "_test"
		}
		c.testFiles[file.rel] = scope
		for _, decl := range file.file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil {
				c.tests[scope+"."+fn.Name.Name] = &funcEntry{file: file, facts: factsOf(fset, file, scope, fn)}
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
