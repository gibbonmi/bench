package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type waitBinding struct {
	expr ast.Expr
	file *ast.File
}

type waitPolicy struct {
	seams    map[string]int
	bindings map[string]map[string]waitBinding
}

// The bounds signatures own their duration argument positions. Test-only helpers
// use other parameter names and are not production wait seams.
func newWaitPolicy(registry *ast.File) *waitPolicy {
	p := &waitPolicy{seams: map[string]int{}, bindings: map[string]map[string]waitBinding{}}
	for _, decl := range registry.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || !fn.Name.IsExported() {
			continue
		}
		for _, position := range durationPositions(fn.Type, registry, "limit") {
			p.seams["bounds."+fn.Name.Name] = position
		}
	}
	return p
}

func (p *waitPolicy) check(file *ast.File, path, rel string) []string {
	bindings, err := p.packageBindings(filepath.Dir(path))
	if err != nil {
		return []string{rel + " cannot resolve wait windows: " + err.Error()}
	}
	var diags []string
	ast.Inspect(file, func(node ast.Node) bool {
		if comparison, ok := node.(*ast.BinaryExpr); ok {
			switch comparison.Op {
			case token.LSS, token.LEQ, token.GTR, token.GEQ, token.EQL, token.NEQ:
				for _, pair := range [][2]ast.Expr{{comparison.X, comparison.Y}, {comparison.Y, comparison.X}} {
					if elapsedWaitTime(pair[0], file, bindings) && !classifiedWait(pair[1], file, bindings, map[ast.Expr]bool{}) {
						diags = append(diags, rel+" has an unclassified wait in current-time deadline")
					}
				}
			}
		}
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		fun := ast.Unparen(call.Fun)
		if ident, ok := fun.(*ast.Ident); ok && ident.Obj != nil {
			if field, ok := ident.Obj.Decl.(*ast.Field); ok {
				if fn, ok := field.Type.(*ast.FuncType); ok {
					for _, position := range durationPositions(fn, file, "") {
						if position < len(call.Args) && !classifiedWait(call.Args[position], file, bindings, map[ast.Expr]bool{}) {
							diags = append(diags, rel+" has an unclassified wait in injected duration function")
						}
					}
				}
			}
		}
		names := waitFunctionNames(fun, file, bindings, map[ast.Expr]bool{})
		if selector, ok := fun.(*ast.SelectorExpr); ok && selector.Sel.Name == "Add" && currentWaitTime(selector.X, file, bindings, map[ast.Expr]bool{}) {
			names = append(names, "current-time deadline")
		}
		for _, name := range names {
			position, wait := p.seams[name]
			switch name {
			case "context.WithTimeout", "context.WithTimeoutCause", "context.WithDeadline", "context.WithDeadlineCause":
				position, wait = 1, true
			case "time.Sleep", "time.After", "time.AfterFunc", "time.NewTimer", "time.NewTicker", "time.Tick", "current-time deadline":
				position, wait = 0, true
			}
			if wait && position < len(call.Args) && !classifiedWait(call.Args[position], file, bindings, map[ast.Expr]bool{}) {
				diags = append(diags, rel+" has an unclassified wait in "+name)
			}
		}
		return true
	})
	return diags
}

// Package defaults can live beside their callers. Local declarations resolve
// through the parser's lexical objects before a package binding is considered.
func (p *waitPolicy) packageBindings(dir string) (map[string]waitBinding, error) {
	if found, ok := p.bindings[dir]; ok {
		return found, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	bindings := map[string]waitBinding{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(dir, name)
		source, err := readBoundsSource(path)
		if err != nil {
			return nil, err
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, source, 0)
		if err != nil {
			return nil, err
		}
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, item := range gen.Specs {
				value, ok := item.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, name := range value.Names {
					if i < len(value.Values) {
						bindings[name.Name] = waitBinding{value.Values[i], file}
					}
				}
			}
		}
	}
	p.bindings[dir] = bindings
	return bindings, nil
}

func classifiedWait(expr ast.Expr, file *ast.File, bindings map[string]waitBinding, seen map[ast.Expr]bool) bool {
	expr = ast.Unparen(expr)
	if seen[expr] {
		return false
	}
	seen[expr] = true
	switch value := expr.(type) {
	case *ast.CallExpr:
		if selector, ok := ast.Unparen(value.Fun).(*ast.SelectorExpr); ok && selector.Sel.Name == "Add" && len(value.Args) == 1 && currentWaitTime(selector.X, file, bindings, map[ast.Expr]bool{}) {
			return classifiedWait(value.Args[0], file, bindings, seen)
		}
		names := waitFunctionNames(value.Fun, file, bindings, map[ast.Expr]bool{})
		return slices.Contains(names, "bounds.VerdictWindow") || slices.Contains(names, "bounds.FixedWindow")
	case *ast.Ident:
		binding, ok := resolveWaitBinding(value, file, bindings)
		if !ok || !classifiedWait(binding.expr, binding.file, bindings, seen) {
			return false
		}
		for _, assigned := range localWaitAssignments(value, file) {
			if !classifiedWait(assigned, file, bindings, seen) {
				return false
			}
		}
		return true
	}
	return false
}

// Package windows retain raw test setters. Local windows must classify each write.
func localWaitAssignments(value *ast.Ident, file *ast.File) []ast.Expr {
	if value.Obj == nil {
		return nil
	}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			if spec == value.Obj.Decl {
				return nil
			}
		}
	}
	var values []ast.Expr
	ast.Inspect(file, func(node ast.Node) bool {
		assignment, ok := node.(*ast.AssignStmt)
		if !ok || assignment == value.Obj.Decl {
			return true
		}
		for i, target := range assignment.Lhs {
			if ident, ok := ast.Unparen(target).(*ast.Ident); ok && ident.Obj == value.Obj && i < len(assignment.Rhs) {
				values = append(values, assignment.Rhs[i])
			}
		}
		return true
	})
	return values
}

// A dot import has no qualifier that selects its package. Keep every supported
// candidate so import order cannot hide a call from the rules that own its name.
func waitCallNames(expr ast.Expr, file *ast.File) []string {
	var qualifier, member string
	switch value := ast.Unparen(expr).(type) {
	case *ast.SelectorExpr:
		ident, ok := value.X.(*ast.Ident)
		if !ok || ident.Obj != nil {
			return nil
		}
		qualifier, member = ident.Name, value.Sel.Name
	case *ast.Ident:
		if value.Obj != nil {
			return nil
		}
		qualifier, member = ".", value.Name
	default:
		return nil
	}
	var names []string
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := filepath.Base(path)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		if name != qualifier {
			continue
		}
		switch path {
		case "time", "context":
			names = append(names, path+"."+member)
		case "github.com/gibbonmi/bench/internal/bounds":
			names = append(names, "bounds."+member)
		}
	}
	return names
}

func resolveWaitBinding(value *ast.Ident, file *ast.File, bindings map[string]waitBinding) (waitBinding, bool) {
	if value.Obj == nil {
		binding, ok := bindings[value.Name]
		return binding, ok
	}
	switch decl := value.Obj.Decl.(type) {
	case *ast.ValueSpec:
		for i, name := range decl.Names {
			if name.Name == value.Name && i < len(decl.Values) {
				return waitBinding{decl.Values[i], file}, true
			}
		}
	case *ast.AssignStmt:
		for i, lhs := range decl.Lhs {
			if name, ok := lhs.(*ast.Ident); ok && name.Name == value.Name && i < len(decl.Rhs) {
				return waitBinding{decl.Rhs[i], file}, true
			}
		}
	}
	return waitBinding{}, false
}

func waitFunctionNames(expr ast.Expr, file *ast.File, bindings map[string]waitBinding, seen map[ast.Expr]bool) []string {
	expr = ast.Unparen(expr)
	if seen[expr] {
		return nil
	}
	seen[expr] = true
	names := waitCallNames(expr, file)
	if ident, ok := expr.(*ast.Ident); ok {
		if binding, found := resolveWaitBinding(ident, file, bindings); found {
			names = append(names, waitFunctionNames(binding.expr, binding.file, bindings, seen)...)
		}
	}
	return names
}

func currentWaitTime(expr ast.Expr, file *ast.File, bindings map[string]waitBinding, seen map[ast.Expr]bool) bool {
	expr = ast.Unparen(expr)
	if seen[expr] {
		return false
	}
	seen[expr] = true
	switch value := expr.(type) {
	case *ast.Ident:
		binding, ok := resolveWaitBinding(value, file, bindings)
		return ok && currentWaitTime(binding.expr, binding.file, bindings, seen)
	case *ast.CallExpr:
		if slices.Contains(waitFunctionNames(value.Fun, file, bindings, map[ast.Expr]bool{}), "time.Now") {
			return true
		}
		fun := ast.Unparen(value.Fun)
		if selector, ok := fun.(*ast.SelectorExpr); ok {
			switch selector.Sel.Name {
			case "UTC", "Local", "Round", "Truncate", "Add":
				return currentWaitTime(selector.X, file, bindings, seen)
			}
		}
		if ident, ok := fun.(*ast.Ident); ok && ident.Obj != nil {
			if field, ok := ident.Obj.Decl.(*ast.Field); ok {
				if fn, ok := field.Type.(*ast.FuncType); ok && fn.Results != nil && len(fn.Results.List) == 1 {
					return slices.Contains(waitCallNames(fn.Results.List[0].Type, file), "time.Time")
				}
			}
		}
	}
	return false
}

func elapsedWaitTime(expr ast.Expr, file *ast.File, bindings map[string]waitBinding) bool {
	call, ok := ast.Unparen(expr).(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return false
	}
	if slices.Contains(waitFunctionNames(call.Fun, file, bindings, map[ast.Expr]bool{}), "time.Since") {
		return currentWaitTime(call.Args[0], file, bindings, map[ast.Expr]bool{})
	}
	selector, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr)
	return ok && selector.Sel.Name == "Sub" && currentWaitTime(selector.X, file, bindings, map[ast.Expr]bool{}) && currentWaitTime(call.Args[0], file, bindings, map[ast.Expr]bool{})
}

func durationPositions(fn *ast.FuncType, file *ast.File, requiredName string) []int {
	var positions []int
	position := 0
	for _, field := range fn.Params.List {
		count := max(1, len(field.Names))
		if slices.Contains(waitCallNames(field.Type, file), "time.Duration") {
			for i := 0; i < count; i++ {
				if requiredName == "" || len(field.Names) > i && field.Names[i].Name == requiredName {
					positions = append(positions, position+i)
				}
			}
		}
		position += count
	}
	return positions
}
