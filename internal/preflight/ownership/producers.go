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
	"strconv"
	"strings"

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
// declaration and every producer expression, roots names the grammar roots, reached maps
// each unit a root reaches to those roots, external maps each root to the declarations of
// other module packages that its units reference, and unresolved names each producer that
// the static walk cannot resolve.
type producers struct {
	units      map[string]unit
	roots      []string
	reached    map[string][]string
	external   map[string][]edge
	unresolved []string

	fset      *token.FileSet
	module    string
	methods   map[string][]string
	ambiguous map[string]bool
	// imports maps each file to its import table: each local import name to its path.
	imports map[string]map[string]string
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
	p := producers{units: map[string]unit{}, reached: map[string][]string{}, external: map[string][]edge{}, fset: token.NewFileSet(),
		module: module, methods: map[string][]string{}, ambiguous: map[string]bool{}, imports: map[string]map[string]string{}}
	var files []*ast.File
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
		files = append(files, file)
		p.index(file)
	}
	for _, custom := range tickets.GrammarRoots() {
		if _, declared := p.units[custom.Name]; custom.Dir == dir && declared {
			p.root(custom.Name)
		}
	}
	for _, file := range files {
		p.parserCalls(file)
	}
	return p, nil
}

// index adds each package-level declaration of one file. A const group is one unit, because
// each member's value can depend on its place in the group.
func (p *producers) index(file *ast.File) {
	imports := map[string]string{}
	for _, spec := range file.Imports {
		imported, _ := strconv.Unquote(spec.Path.Value)
		name := path.Base(imported)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		imports[name] = imported
	}
	p.imports[p.fset.Position(file.Package).Filename] = imports
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				p.add(d.Name.Name, d)
				continue
			}
			key := receiverName(d.Recv.List[0].Type) + "." + d.Name.Name
			p.add(key, d)
			p.methods[d.Name.Name] = append(p.methods[d.Name.Name], key)
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					node := ast.Node(s)
					if d.Tok == token.CONST {
						node = d
					}
					for _, name := range s.Names {
						p.add(name.Name, node)
					}
				case *ast.TypeSpec:
					p.add(s.Name.Name, s)
				}
			}
		}
	}
}

// receiverName is the base type name of one method receiver.
func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverName(t.X)
	case *ast.IndexExpr:
		return receiverName(t.X)
	case *ast.IndexListExpr:
		return receiverName(t.X)
	case *ast.Ident:
		return t.Name
	}
	return ""
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
		if imported, _ := strconv.Unquote(spec.Path.Value); imported == parse.Dir || strings.HasSuffix(imported, "/"+parse.Dir) {
			local = path.Base(imported)
			if spec.Name != nil {
				local = spec.Name.Name
			}
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
				p.producer(fn, call, calls)
			}
			return true
		})
	}
}

// producer roots the grammar argument of one parser call in fn. A package declaration, a
// single-definition local, a call of a package function or method, and a literal resolve;
// a parameter, a reassigned local, and another package's value stay unresolved.
func (p *producers) producer(fn *ast.FuncDecl, call *ast.CallExpr, n int) {
	where := fmt.Sprintf("%s parser call %d", fn.Name.Name, n)
	if len(call.Args) == 0 {
		p.unresolved = append(p.unresolved, where+" passes no grammar")
		return
	}
	switch arg := call.Args[0].(type) {
	case *ast.Ident:
		value, local, ok := localValue(fn, arg.Name)
		switch {
		case local && ok:
			p.rootUnit(fn.Name.Name+"."+arg.Name, value)
			return
		case !local:
			if _, declared := p.units[arg.Name]; declared {
				p.root(arg.Name)
				return
			}
		}
	case *ast.CallExpr:
		if p.resolvedCall(p.fset.Position(call.Pos()).Filename, arg) {
			p.rootUnit(fmt.Sprintf("%s#%d", fn.Name.Name, n), arg)
			return
		}
	case *ast.CompositeLit, *ast.UnaryExpr:
		p.rootUnit(fmt.Sprintf("%s#%d", fn.Name.Name, n), arg)
		return
	}
	p.unresolved = append(p.unresolved, where+" passes a grammar that no static producer declares")
}

// rootUnit roots one producer expression under key. Two parser calls that pass one local
// share its unit.
func (p *producers) rootUnit(key string, node ast.Node) {
	if _, declared := p.units[key]; !declared {
		p.add(key, node)
	}
	p.root(key)
}

// resolvedCall reports whether call, in file, invokes a package function or a method this
// package declares.
func (p *producers) resolvedCall(file string, call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		_, declared := p.units[fun.Name]
		return declared
	case *ast.SelectorExpr:
		_, imported := p.imported(file, fun.X)
		return !imported && len(p.methods[fun.Sel.Name]) > 0
	}
	return false
}

// localValue is the one value fn assigns to a local name. local reports that fn binds the
// name at all; ok is false for a parameter or for a local with no single defining value.
func localValue(fn *ast.FuncDecl, name string) (value ast.Expr, local, ok bool) {
	for _, list := range []*ast.FieldList{fn.Recv, fn.Type.Params, fn.Type.Results} {
		if list != nil && slices.ContainsFunc(list.List, func(field *ast.Field) bool {
			return slices.ContainsFunc(field.Names, func(id *ast.Ident) bool { return id.Name == name })
		}) {
			return nil, true, false
		}
	}
	var values []ast.Expr
	bind := func(names []ast.Expr, rhs []ast.Expr) {
		for i, lhs := range names {
			if isIdent(lhs, name) {
				local = true
				if len(names) == len(rhs) {
					values = append(values, rhs[i])
				} else {
					values = append(values, nil)
				}
			}
		}
	}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		switch s := node.(type) {
		case *ast.AssignStmt:
			bind(s.Lhs, s.Rhs)
		case *ast.ValueSpec:
			names := make([]ast.Expr, len(s.Names))
			for i, id := range s.Names {
				names[i] = id
			}
			bind(names, s.Values)
		case *ast.RangeStmt:
			bind([]ast.Expr{s.Key, s.Value}, nil)
		case *ast.FuncLit:
			for _, field := range s.Type.Params.List {
				for _, id := range field.Names {
					bind([]ast.Expr{id}, nil)
				}
			}
		}
		return true
	})
	if len(values) == 1 && values[0] != nil {
		return values[0], true, true
	}
	return nil, local, false
}

func isIdent(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}

// root walks one grammar root to every unit it reaches, and records each reference to
// another package of the module that those units hold.
func (p *producers) root(key string) {
	if slices.Contains(p.roots, key) {
		return
	}
	p.roots = append(p.roots, key)
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
