package ownership

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
)

// unit is one producer declaration: the file and line that declare it, its printed syntax,
// and its syntax tree.
type unit struct {
	file string
	line int
	text string
	node ast.Node
}

// producers is one package's grammar producers at one revision. units holds every package
// declaration and every producer expression, roots names the grammar roots, rows maps each
// root to the selector of the registry row that registered it, reached maps each unit a root
// reaches to those roots, external maps each root to the declarations of other module
// packages that its units reference, and unresolved names each producer that the static walk
// cannot resolve.
type producers struct {
	units      map[string]unit
	roots      []string
	rows       map[string]string
	reached    map[string][]string
	external   map[string][]edge
	unresolved []string

	fset      *token.FileSet
	module    string
	methods   map[string][]string
	ambiguous map[string]bool
	// imports maps each file to its import table: each local import name to its path.
	imports map[string]map[string]string
	// files holds the package's parsed production files, where callers finds call sites.
	files []*ast.File
}

// readProducers reads the production Go files of dir at rev and walks its grammar roots: each
// value passed to the registry's ordinary parser, and each custom root the registry names
// in dir. module is the module path at rev, which places an import inside the module; an
// empty module places every import outside it.
func readProducers(root, rev, dir, module string) (producers, error) {
	listing, err := git.Output("-C", root, "ls-tree", "-z", "--name-only", rev, "--", dir+"/")
	if err != nil {
		return producers{}, fmt.Errorf("%s listing at %s: %w", dir, rev, err)
	}
	p := producers{units: map[string]unit{}, rows: map[string]string{}, reached: map[string][]string{}, external: map[string][]edge{}, fset: token.NewFileSet(),
		module: module, methods: map[string][]string{}, ambiguous: map[string]bool{}, imports: map[string]map[string]string{}}
	for _, name := range strings.Split(listing, "\x00") {
		if path.Dir(name) != dir || !producerFile(name) {
			continue
		}
		data, err := git.ReadTreeFile(root, rev, name)
		if err != nil {
			return producers{}, err
		}
		file, err := parser.ParseFile(p.fset, name, data, parser.SkipObjectResolution)
		if err != nil {
			return producers{}, err
		}
		p.files = append(p.files, file)
		p.index(file)
	}
	for _, custom := range tickets.GrammarRoots() {
		if _, declared := p.units[custom.Name]; custom.Dir == dir && declared {
			p.root(custom.Name, custom)
		}
	}
	for _, file := range p.files {
		p.parserCalls(file)
	}
	return p, nil
}

// index adds each package-level declaration of one file under the name that the consumers
// declaration reader gives it. A const group is one unit, because each member's value can
// depend on its place in the group.
func (p *producers) index(file *ast.File) {
	imports := map[string]string{}
	for _, spec := range file.Imports {
		name, imported := consumers.ImportName(spec)
		imports[name] = imported
	}
	p.imports[p.fset.Position(file.Package).Filename] = imports
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			names, _ := consumers.DeclNames(d)
			p.add(names[0], d)
			if d.Recv != nil {
				p.methods[d.Name.Name] = append(p.methods[d.Name.Name], names[0])
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				names, ok := consumers.DeclNames(spec)
				if !ok {
					continue
				}
				node := ast.Node(spec)
				if d.Tok == token.CONST {
					node = d
				}
				for _, name := range names {
					p.add(name, node)
				}
			}
		}
	}
}

// add records one unit. A second declaration of one name is ambiguous, so a walk that
// reaches it cannot classify the package.
func (p *producers) add(key string, node ast.Node) {
	if key == "_" || key == "init" {
		return
	}
	if _, seen := p.units[key]; seen {
		p.ambiguous[key] = true
		return
	}
	var text bytes.Buffer
	if err := printer.Fprint(&text, p.fset, node); err != nil {
		p.unresolved = append(p.unresolved, key+" does not print: "+err.Error())
	}
	position := p.fset.Position(node.Pos())
	p.units[key] = unit{file: position.Filename, line: position.Line, text: text.String(), node: node}
}

// parserCalls roots each value that one file passes to the registry's ordinary parser.
func (p *producers) parserCalls(file *ast.File) {
	parse := tickets.GrammarParser()
	local := ""
	for _, spec := range file.Imports {
		if name, imported := consumers.ImportName(spec); consumers.PathHasSuffix(imported, parse.Dir) {
			local = name
		}
	}
	if local == "" {
		return
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		calls := 0
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == parse.Name && isIdent(sel.X, local) {
				calls++
				p.producer(fn, call, calls, parse)
			}
			return true
		})
	}
}

// producer roots the grammar argument of one parser call in fn, which parse registers.
func (p *producers) producer(fn *ast.FuncDecl, call *ast.CallExpr, n int, parse tickets.GrammarRoot) {
	where := fmt.Sprintf("%s parser call %d", fn.Name.Name, n)
	if len(call.Args) == 0 {
		p.unresolved = append(p.unresolved, where+" passes no grammar")
		return
	}
	found, ok := p.resolve(fn, call.Args[0], fmt.Sprintf("%s#%d", fn.Name.Name, n), map[*ast.FuncDecl]bool{})
	if !ok {
		p.unresolved = append(p.unresolved, where+" passes a grammar that no static producer declares")
		return
	}
	for _, producer := range found {
		if producer.node == nil {
			p.root(producer.key, parse)
		} else {
			p.rootUnit(producer.key, producer.node, parse)
		}
	}
}

// rootUnit roots one producer expression under key, which row registers. Two parser calls
// that pass one local share its unit.
func (p *producers) rootUnit(key string, node ast.Node, row tickets.GrammarRoot) {
	if _, declared := p.units[key]; !declared {
		p.add(key, node)
	}
	p.root(key, row)
}

// root walks one grammar root to every unit it reaches, and records each reference to
// another package of the module that those units hold. The first row that registers a root
// keeps it.
func (p *producers) root(key string, row tickets.GrammarRoot) {
	if slices.Contains(p.roots, key) {
		return
	}
	p.roots, p.rows[key] = append(p.roots, key), row.String()
	keys, external, problems := p.reach(key)
	for _, current := range keys {
		p.reached[current] = append(p.reached[current], key)
	}
	p.external[key] = external
	p.unresolved = append(p.unresolved, problems...)
}

// reach is every unit that key reaches, the references to other packages of the module that
// those units hold, and each reach that the static walk cannot resolve. A value, const, or
// type reference reaches its declaration. A called function or method reaches its body. A
// function used only as a value contributes its name, not its body, and a function literal
// contributes nothing, so a handler's body is no grammar.
func (p *producers) reach(key string) (keys []string, external []edge, problems []string) {
	seen := map[string]bool{key: true}
	for queue := []string{key}; len(queue) > 0; queue = queue[1:] {
		current := queue[0]
		keys = append(keys, current)
		local, out, bad := p.deps(p.units[current])
		external, problems = append(external, out...), append(problems, bad...)
		for _, dep := range local {
			if p.ambiguous[dep] {
				problems = append(problems, key+" "+ambiguity(dep))
			}
			if !seen[dep] {
				seen[dep] = true
				queue = append(queue, dep)
			}
		}
	}
	return keys, external, problems
}

// deps is every package unit that one unit references, by the rule reach states, and every
// declaration of another package of the module that it references through its file's
// import table. An import outside the module adds nothing. A dot import hides which package
// an identifier names, so it leaves the unit unresolved.
func (p *producers) deps(u unit) (deps []string, external []edge, problems []string) {
	if u.node == nil {
		return nil, nil, nil
	}
	if imported, dot := p.imports[u.file]["."]; dot {
		problems = append(problems, u.file+" dot-imports "+imported)
	}
	called := map[ast.Expr]bool{}
	var visit func(ast.Node) bool
	visit = func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.FuncLit:
			return false
		case *ast.CallExpr:
			called[n.Fun] = true
			if fun, ok := n.Fun.(*ast.SelectorExpr); ok {
				if _, imported := p.imported(u.file, fun.X); !imported {
					deps = append(deps, p.methods[fun.Sel.Name]...)
				}
			}
		case *ast.SelectorExpr:
			imported, isImport := p.imported(u.file, n.X)
			if !isImport {
				ast.Inspect(n.X, visit)
			} else if dir, in := p.inModule(imported); in {
				external = append(external, edge{target: tickets.GrammarRoot{Dir: dir, Name: n.Sel.Name}, called: called[n]})
			}
			return false
		case *ast.Ident:
			if declared, ok := p.units[n.Name]; ok && declared.reachedBy(called[n]) {
				deps = append(deps, n.Name)
			}
		}
		return true
	}
	ast.Inspect(u.node, visit)
	return deps, external, problems
}
