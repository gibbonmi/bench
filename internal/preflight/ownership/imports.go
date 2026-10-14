package ownership

import (
	"errors"
	"go/ast"
	"path"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
)

// edge is one reference from a grammar unit to a declaration of another package of the
// module: that declaration's selector, and whether the reference calls it.
type edge struct {
	target tickets.GrammarRoot
	called bool
}

// imported is the import path that expr names in file, when expr is an import name there.
func (p *producers) imported(file string, expr ast.Expr) (string, bool) {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return "", false
	}
	imported, ok := p.imports[file][id.Name]
	return imported, ok
}

// inModule is the package directory of one import path inside the module.
func (p *producers) inModule(imported string) (string, bool) {
	if p.module == "" {
		return "", false
	}
	if imported == p.module {
		return ".", true
	}
	return strings.CutPrefix(imported, p.module+"/")
}

// reachedBy reports whether a reference reaches u: a called function reaches its body, and
// every reference reaches a value, const, or type.
func (u unit) reachedBy(called bool) bool {
	_, function := u.node.(*ast.FuncDecl)
	return called || !function
}

// ambiguity is the reason that a walk reaches a name with more than one declaration.
func ambiguity(name string) string {
	return "reaches the ambiguous declaration " + name
}

// revRead is one fact read at one revision, or the fault that stopped it.
type revRead[T any] struct {
	value T
	err   error
}

// module is the module path that go.mod declares at rev, read once. A revision with no
// go.mod has no module, so every import there is outside the module.
func (e *effects) module(rev string) (string, error) {
	if cached, ok := e.modules[rev]; ok {
		return cached.value, cached.err
	}
	module := ""
	listing, err := git.Output("-C", e.root, "ls-tree", "-z", "--name-only", rev, "--", "go.mod")
	if err == nil && listing != "" {
		var data []byte
		if data, err = git.ReadTreeFile(e.root, rev, "go.mod"); err == nil {
			if module = modfile.ModulePath(data); module == "" {
				err = errors.New("go.mod declares no module path at " + rev)
			}
		}
	}
	e.modules[rev] = revRead[string]{module, err}
	return module, err
}

// referred is every unit of dir at rev that a grammar root of another package reaches through
// in-module imports, with each such root as a grammar change that names no unit yet. Only a
// referrer's resolved roots refer: a producer that the referrer cannot resolve is that
// package's own diagnostic. The failure names the first referrer read or reference that the
// static walk cannot resolve.
func (e *effects) referred(rev, dir string) (map[string][]grammarChange, string) {
	referrers, err := e.referrers(rev, dir)
	if err != nil {
		return nil, err.Error()
	}
	found := map[string][]grammarChange{}
	for _, from := range referrers {
		read := e.read(rev, from)
		if read.err != nil {
			return nil, read.err.Error()
		}
		for _, key := range read.producers.roots {
			at := read.producers.units[key]
			origin := grammarChange{root: tickets.GrammarRoot{Dir: from, Name: key}.String(), file: at.file, line: at.line}
			reached, failure := e.follow(rev, read.producers.external[key])
			if failure != "" {
				return nil, origin.root + " " + failure
			}
			for _, target := range reached {
				if target.Dir == dir {
					found[target.Name] = append(found[target.Name], origin)
				}
			}
		}
	}
	return found, ""
}

// follow walks references to other packages of the module to their fixed point at rev: each
// referenced declaration that the reference reaches, and each unit that it reaches in its own
// package by the rule reach states. A reference whose package does not declare it exactly
// once stays unresolved.
func (e *effects) follow(rev string, edges []edge) ([]tickets.GrammarRoot, string) {
	var reached []tickets.GrammarRoot
	seen := map[edge]bool{}
	for queue := slices.Clone(edges); len(queue) > 0; queue = queue[1:] {
		next := queue[0]
		if seen[next] {
			continue
		}
		seen[next] = true
		read := e.read(rev, next.target.Dir)
		declared, ok := read.producers.units[next.target.Name]
		switch {
		case read.err != nil:
			return nil, read.err.Error()
		case !ok:
			return nil, "references " + next.target.String() + ", which its package does not declare"
		case read.producers.ambiguous[next.target.Name]:
			return nil, ambiguity(next.target.String())
		}
		if !declared.reachedBy(next.called) {
			continue
		}
		keys, external, problems := read.producers.reach(next.target.Name)
		if len(problems) > 0 {
			return nil, problems[0]
		}
		for _, key := range keys {
			reached = append(reached, tickets.GrammarRoot{Dir: next.target.Dir, Name: key})
		}
		queue = append(queue, external...)
	}
	return reached, ""
}

// referrers is every grammar package at rev, other than dir, that can reach dir through
// in-module imports: the transitive importers of dir that declare a custom root or name the
// ordinary parser. The fixed string search is a superset of the importing files, so it drops
// no referrer.
func (e *effects) referrers(rev, dir string) ([]string, error) {
	module, err := e.module(rev)
	if err != nil || module == "" {
		return nil, err
	}
	grammar, err := e.grammarDirs(rev)
	if err != nil {
		return nil, err
	}
	importers := map[string]bool{dir: true}
	var found []string
	for frontier := []string{dir}; len(frontier) > 0; {
		needles := make([]string, 0, len(frontier))
		for _, importer := range frontier {
			needles = append(needles, path.Join(module, importer))
		}
		files, err := grepFiles(e.root, rev, needles...)
		if err != nil {
			return nil, err
		}
		frontier = nil
		for _, file := range files {
			if importer := path.Dir(file); producerFile(file) && !importers[importer] {
				importers[importer] = true
				frontier = append(frontier, importer)
				if grammar[importer] {
					found = append(found, importer)
				}
			}
		}
	}
	slices.Sort(found)
	return found, nil
}

// grammarDirs is every package directory at rev that can declare a grammar root, read once:
// each custom root's directory, and each directory whose production sources name the
// ordinary parser's package.
func (e *effects) grammarDirs(rev string) (map[string]bool, error) {
	if cached, ok := e.grammar[rev]; ok {
		return cached.value, cached.err
	}
	dirs := map[string]bool{}
	for _, root := range tickets.GrammarRoots() {
		dirs[root.Dir] = true
	}
	files, err := grepFiles(e.root, rev, tickets.GrammarParser().Dir)
	for _, file := range files {
		if producerFile(file) {
			dirs[path.Dir(file)] = true
		}
	}
	e.grammar[rev] = revRead[map[string]bool]{dirs, err}
	return dirs, err
}
