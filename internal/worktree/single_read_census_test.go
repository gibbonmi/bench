package worktree

// This file owns the single-read census. A census entry is an exported function or
// method. The read set is each effects.go function, each package-qualified function that
// effects.go calls, and each exported gate function that reaches such a gate function. A
// kind is one package-qualified call that a read-set name reaches. An entry can read
// each kind once in its own body, outside each function literal and loop body.

import (
	"fmt"
	"go/ast"
	"go/token"
	"maps"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// effectsFile is the effect boundary. The census reads it for the read set only.
const effectsFile = "effects.go"

// isSourceFile reports whether name is a non-test .go file.
func isSourceFile(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

// readFinding selects one of the three census messages.
type readFinding int

const (
	readBelowEntry readFinding = iota
	readSecondTime
	readThroughEntry
)

// singleReadReport renders one census message. A read below an entry names the source
// spelling, a second read names its kind, and a call names the entry and one kind.
func singleReadReport(file string, line int, declaration string, finding readFinding, name, kind string) string {
	where := fmt.Sprintf("%s:%d: %s", file, line, declaration)
	switch finding {
	case readSecondTime:
		return fmt.Sprintf("%s reads %s a second time", where, kind)
	case readThroughEntry:
		return fmt.Sprintf("%s calls %s, which reads %s", where, name, kind)
	}
	return fmt.Sprintf("%s reads %s below a census entry", where, name)
}

// qualifiedName spells expr as package.Name and keys it as path.Name when expr selects
// from a package that imports maps to its path. Both are empty otherwise. A kind is
// keyed by its import path, so an import alias cannot hide a read.
func qualifiedName(expr ast.Expr, imports map[string]string) (spelled, key string) {
	if selector, ok := expr.(*ast.SelectorExpr); ok {
		if pkg, ok := selector.X.(*ast.Ident); ok && imports[pkg.Name] != "" {
			return pkg.Name + "." + selector.Sel.Name, imports[pkg.Name] + "." + selector.Sel.Name
		}
	}
	return "", ""
}

// topLevelFuncs indexes each function that files declare by name. A method is not
// indexed, because the census follows calls by bare identifier.
func topLevelFuncs(files []*ast.File) map[string]*ast.FuncDecl {
	funcs := map[string]*ast.FuncDecl{}
	for _, file := range files {
		for _, node := range file.Decls {
			if decl, ok := node.(*ast.FuncDecl); ok && decl.Recv == nil && decl.Body != nil {
				funcs[decl.Name.Name] = decl
			}
		}
	}
	return funcs
}

// reachedFuncs adds name and each function of funcs that name reaches by bare calls.
func reachedFuncs(name string, funcs map[string]*ast.FuncDecl, reached map[string]bool) map[string]bool {
	if decl, ok := funcs[name]; ok && !reached[name] {
		reached[name] = true
		for _, callee := range directCallees(decl) {
			reachedFuncs(callee, funcs, reached)
		}
	}
	return reached
}

// readSet maps each read-set name to the sorted kinds that it reaches. A gate function
// takes the import path of the gate package, whose last element is the package name. An
// empty read set is an error, because a census with no read to grade passes vacuously.
func readSet(dir, gateDir string) (map[string][]string, error) {
	files, _, _, err := parseGoFiles(dir, func(name string) bool { return name == effectsFile })
	if err != nil {
		return nil, err
	}
	reads, effects := map[string][]string{}, topLevelFuncs(files)
	for _, file := range files {
		imports := fileImportNames(file)
		for name := range effects {
			reads[name] = []string{}
			for reached := range reachedFuncs(name, effects, map[string]bool{}) {
				ast.Inspect(effects[reached].Body, func(node ast.Node) bool {
					if call, ok := node.(*ast.CallExpr); ok {
						if _, kind := qualifiedName(call.Fun, imports); kind != "" {
							reads[name], reads[kind] = append(reads[name], kind), []string{kind}
						}
					}
					return true
				})
			}
		}
	}
	if len(reads) == 0 {
		return nil, fmt.Errorf("the read set is empty: %s in %s declares no function", effectsFile, dir)
	}
	gateFiles, err := parseSourceFiles(gateDir)
	if err != nil {
		return nil, err
	}
	gate := topLevelFuncs(gateFiles)
	for _, kind := range slices.Collect(maps.Keys(reads)) {
		dot := strings.LastIndex(kind, ".")
		if dot < 0 || gate[kind[dot+1:]] == nil || path.Base(kind[:dot]) != gateFiles[0].Name.Name {
			continue
		}
		for name := range gate {
			if ast.IsExported(name) && reachedFuncs(name, gate, map[string]bool{})[kind[dot+1:]] {
				reads[kind[:dot]+"."+name] = append(reads[kind[:dot]+"."+name], kind)
			}
		}
	}
	for name, kinds := range reads {
		slices.Sort(kinds)
		reads[name] = slices.Compact(kinds)
	}
	return reads, nil
}

// readRef is one reference to a read-set name, with the kinds that it reaches, or to an
// entry. nested marks a reference in a function literal or a loop, which can run twice.
type readRef struct {
	pos    token.Pos
	name   string
	kinds  []string
	entry  bool
	nested bool
}

// declRefs returns each reference in root to a read-set name or to an entry, in source
// order. A selector names a package only when the file imports it, and otherwise it
// names an entry by its selected name, as a method call does. A field, and a struct
// literal key, named like a read is not the read.
func declRefs(root ast.Node, imports map[string]string, reads map[string][]string, entries map[string]bool) []readRef {
	var refs []readRef
	var walk func(ast.Node, bool)
	walk = func(root ast.Node, nested bool) {
		if root == nil {
			return
		}
		ast.Inspect(root, func(node ast.Node) bool {
			name, key := "", ""
			switch node := node.(type) {
			case *ast.FuncLit:
				walk(node.Body, true)
				return false
			case *ast.ForStmt:
				walk(node.Init, nested)
				walk(node.Cond, true)
				walk(node.Post, true)
				walk(node.Body, true)
				return false
			case *ast.RangeStmt:
				walk(node.X, nested)
				walk(node.Body, true)
				return false
			case *ast.KeyValueExpr:
				if _, ok := node.Key.(*ast.Ident); ok {
					walk(node.Value, nested)
					return false
				}
			case *ast.SelectorExpr:
				if name, key = qualifiedName(node, imports); name == "" {
					walk(node.X, nested)
					if entries[node.Sel.Name] {
						refs = append(refs, readRef{pos: node.Sel.Pos(), name: node.Sel.Name, entry: true, nested: nested})
					}
					return false
				}
			case *ast.Ident:
				name, key = node.Name, node.Name
			}
			if kinds, read := reads[key]; read || entries[key] {
				refs = append(refs, readRef{pos: node.Pos(), name: name, kinds: kinds, entry: !read, nested: nested})
			}
			return name == ""
		})
	}
	walk(root, false)
	return refs
}

// singleReadCensus reports each read in the non-test files of dir, other than
// effectsFile, that is not the first read of its kind in an entry's own body. It also
// reports each reference from another declaration to an entry whose own body reads.
// gateDir is the gate package directory. The reports are sorted.
func singleReadCensus(dir, gateDir string) ([]string, error) {
	reads, err := readSet(dir, gateDir)
	if err != nil {
		return nil, err
	}
	files, fset, _, err := parseGoFiles(dir, func(name string) bool { return name != effectsFile && isSourceFile(name) })
	if err != nil {
		return nil, err
	}
	entries := map[string]bool{}
	for _, file := range files {
		for _, node := range file.Decls {
			if decl, ok := node.(*ast.FuncDecl); ok && decl.Name.IsExported() {
				entries[decl.Name.Name] = true
			}
		}
	}
	type graded struct {
		name  string
		entry bool
		refs  []readRef
	}
	var decls []graded
	entryKinds := map[string][]string{}
	for _, file := range files {
		imports := fileImportNames(file)
		for _, node := range file.Decls {
			if decl, ok := node.(*ast.FuncDecl); ok && decl.Body != nil {
				refs := declRefs(decl.Body, imports, reads, entries)
				decls = append(decls, graded{decl.Name.Name, decl.Name.IsExported(), refs})
				for _, ref := range refs {
					entryKinds[decl.Name.Name] = append(entryKinds[decl.Name.Name], ref.kinds...)
				}
			} else if decl, ok := node.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					if value, ok := spec.(*ast.ValueSpec); ok {
						for _, expr := range value.Values {
							decls = append(decls, graded{specName(spec), false, declRefs(expr, imports, reads, entries)})
						}
					}
				}
			}
		}
	}
	var reports []string
	for _, decl := range decls {
		seen := map[string]bool{}
		for _, ref := range decl.refs {
			position := fset.Position(ref.pos)
			report := func(finding readFinding, name, kind string) {
				reports = append(reports, singleReadReport(filepath.Base(position.Filename), position.Line, decl.name, finding, name, kind))
			}
			switch {
			case ref.entry && !decl.entry:
				for _, kind := range entryKinds[ref.name] {
					report(readThroughEntry, ref.name, kind)
				}
			case ref.entry:
				// The called entry's reads are its own, so an entry may call it.
			case ref.nested || !decl.entry:
				report(readBelowEntry, ref.name, "")
			default:
				for _, kind := range ref.kinds {
					if seen[kind] {
						report(readSecondTime, "", kind)
					}
					seen[kind] = true
				}
			}
		}
	}
	// An entry that reads one kind twice yields one helper-call report for that kind.
	sort.Strings(reports)
	return slices.Compact(reports), nil
}
