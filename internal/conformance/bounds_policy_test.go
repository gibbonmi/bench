package conformance

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func checkBoundsPolicy(root string) []string {
	registryPath := filepath.Join(root, "internal", "bounds", "bounds.go")
	registry := readIfExists(registryPath)
	if registry == "" {
		return []string{"internal/bounds policy registry is absent"}
	}
	required := []string{"ProviderTimeout", "GitRefreshTimeout", "WorktreeListTimeout", "GuardScanTimeout", "GateTimeout", "PackageLoadTimeout", "ModelReadLimit", "OutlineFileLimit", "ControlRecordLimit", "IterationMin", "IterationMax", "MainIterationsDefault", "RefactorIterationsDefault", "MaxWall", "LeaseStale", "AssignmentStale", "PreviewRuneLimit"}
	var diags []string
	for _, name := range required {
		if !strings.Contains(registry, name) {
			diags = append(diags, "internal/bounds policy registry missing "+name)
		}
	}
	owners := map[string][]string{
		"internal/models/models.go":                 {"bounds.ProviderTimeout", "bounds.ModelReadLimit"},
		"internal/sessioninspect/sessioninspect.go": {"bounds.ProviderTimeout"},
		"internal/outline/read.go":                  {"bounds.OutlineFileLimit"},
		"internal/learnings/learnings.go":           {"bounds.ControlRecordLimit"},
		"internal/maps/maps.go":                     {"bounds.ControlRecordLimit"},
		"internal/roadmap/roadmap.go":               {"bounds.ControlRecordLimit"},
		"internal/guards/guards.go":                 {"bounds.GuardScanTimeout"},
		"internal/sanitize/sanitize.go":             {"bounds.PreviewRuneLimit"},
		"internal/gate/gate.go":                     {"bounds.GateTimeout"},
		"internal/coverage/citation_execution.go":   {"bounds.PackageLoadTimeout"},
		"internal/refresh/refresh.go":               {"bounds.GitRefreshTimeout"},
		"internal/git/git.go":                       {"bounds.WorktreeListTimeout"},
		"internal/worktree/lifecycle.go":            {"bounds.LeaseStale"},
		"internal/worktree/classifier.go":           {"bounds.AssignmentStale"},
		"internal/shift/loop.go":                    {"bounds.MainIterationsDefault", "bounds.RefactorIterationsDefault", "bounds.IterationMin", "bounds.IterationMax", "bounds.MaxWall"},
	}
	for rel, tokens := range owners {
		body := readIfExists(filepath.Join(root, filepath.FromSlash(rel)))
		for _, token := range tokens {
			if !strings.Contains(body, token) {
				diags = append(diags, rel+" does not consume "+token)
			}
		}
	}
	diags = append(diags, checkBoundCallers(root, registryPath)...)
	wrapper := readIfExists(filepath.Join(root, "bin", "bench.sh"))
	if !strings.Contains(wrapper, `[[ "${BENCH_OFFLINE:-}" == 1 ]]`) || !strings.Contains(registry, `os.Getenv("BENCH_OFFLINE") == "1"`) {
		diags = append(diags, "wrapper and Go offline checks do not share exact BENCH_OFFLINE=1 semantics")
	}
	return diags
}

func checkBoundCallers(root, registryPath string) []string {
	fset := token.NewFileSet()
	registry, err := parser.ParseFile(fset, registryPath, nil, 0)
	if err != nil {
		return []string{"internal/bounds policy registry is not valid Go: " + err.Error()}
	}
	owned := boundsOwnedExpressions(fset, registry)
	seams, err := boundsReadSeams(filepath.Dir(registryPath))
	if err != nil {
		return []string{"internal/bounds read seams cannot be parsed: " + err.Error()}
	}
	var diags []string
	for _, top := range []string{"cmd", "internal"} {
		base := filepath.Join(root, top)
		_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel := slashRel(root, path)
			if strings.HasPrefix(rel, "internal/bounds/") {
				return nil
			}
			diags = append(diags, checkBoundCaller(fset, path, rel, owned, seams)...)
			return nil
		})
	}
	return uniqueSorted(diags)
}

// boundsOwnedExpressions maps the expression text of each registry value to its entry name.
func boundsOwnedExpressions(fset *token.FileSet, registry *ast.File) map[string]string {
	owned := map[string]string{}
	for _, decl := range registry.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, item := range gen.Specs {
			spec := item.(*ast.ValueSpec)
			for i, name := range spec.Names {
				if i < len(spec.Values) {
					owned[expressionText(fset, spec.Values[i])] = name.Name
				}
			}
		}
	}
	return owned
}

// boundsReadSeams maps each read seam of the bounds package to the position of its limit
// argument. A read seam is an exported function with a limit parameter of type int64; the
// package source is the one list of seams.
func boundsReadSeams(dir string) (map[string]int, error) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	seams := map[string]int{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			return nil, err
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !fn.Name.IsExported() {
				continue
			}
			position := 0
			for _, field := range fn.Type.Params.List {
				kind, _ := field.Type.(*ast.Ident)
				for _, param := range field.Names {
					if param.Name == "limit" && kind != nil && kind.Name == "int64" {
						seams["bounds."+fn.Name.Name] = position
					}
					position++
				}
			}
		}
	}
	return seams, nil
}

func checkBoundCaller(fset *token.FileSet, path, rel string, ownedExpressions map[string]string, seams map[string]int) []string {
	if readIfExists(path) == "" {
		return nil
	}
	caller, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for bounds ownership: " + err.Error()}
	}
	var diags []string
	ast.Inspect(caller, func(node ast.Node) bool {
		switch value := node.(type) {
		case *ast.ValueSpec:
			for i, expr := range value.Values {
				if owner := ownedExpressions[expressionText(fset, expr)]; owner != "" {
					name := ""
					if i < len(value.Names) {
						name = value.Names[i].Name
					}
					if !isBasicLiteral(expr) || boundLikeName(name, owner) {
						diags = append(diags, rel+" redeclares "+owner+" policy value")
					}
				}
			}
		case *ast.CallExpr:
			selector, ok := value.Fun.(*ast.SelectorExpr)
			if !ok {
				break
			}
			pkg, ok := selector.X.(*ast.Ident)
			if !ok {
				break
			}
			call := pkg.Name + "." + selector.Sel.Name
			if position, ok := seams[call]; ok && position < len(value.Args) {
				if owner := ownedExpressions[expressionText(fset, value.Args[position])]; owner != "" {
					diags = append(diags, rel+" restates "+owner+" in the limit of "+call)
				}
			}
			if call == "context.WithTimeout" || call == "context.WithTimeoutCause" || call == "io.LimitReader" {
				for _, arg := range value.Args {
					if expressionOwnsBound(fset, arg, ownedExpressions) {
						diags = append(diags, rel+" reimplements bounded operation with "+call+" instead of internal/bounds")
						break
					}
				}
			}
		}
		return true
	})
	return uniqueSorted(diags)
}

func isBasicLiteral(expr ast.Expr) bool {
	_, ok := expr.(*ast.BasicLit)
	return ok
}

func boundLikeName(name, owner string) bool {
	name = strings.ToLower(name)
	words := map[string][]string{
		"ProviderTimeout":           {"provider", "timeout"},
		"GitRefreshTimeout":         {"refresh", "timeout"},
		"GuardScanTimeout":          {"guard", "timeout"},
		"GateTimeout":               {"gate", "timeout"},
		"ModelReadLimit":            {"model", "limit"},
		"OutlineFileLimit":          {"outline", "file", "limit"},
		"ControlRecordLimit":        {"control", "record", "limit"},
		"IterationMin":              {"iteration", "min"},
		"IterationMax":              {"iteration", "max"},
		"MainIterationsDefault":     {"main", "iteration", "default"},
		"RefactorIterationsDefault": {"refactor", "iteration", "default"},
		"MaxWall":                   {"wall", "max"},
	}[owner]
	for _, word := range words {
		if !strings.Contains(name, word) {
			return false
		}
	}
	return len(words) > 0
}

func expressionOwnsBound(fset *token.FileSet, expr ast.Expr, owned map[string]string) bool {
	if owned[expressionText(fset, expr)] != "" {
		return true
	}
	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok || selector.Sel == nil {
			return true
		}
		pkg, ok := selector.X.(*ast.Ident)
		if ok && pkg.Name == "bounds" {
			found = true
			return false
		}
		return true
	})
	return found
}

func expressionText(fset *token.FileSet, expr ast.Expr) string {
	var out bytes.Buffer
	if err := format.Node(&out, fset, expr); err != nil {
		return ""
	}
	return out.String()
}
