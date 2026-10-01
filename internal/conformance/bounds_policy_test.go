package conformance

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

func checkBoundsPolicy(root string) []string {
	registryPath := filepath.Join(root, "internal", "bounds", "bounds.go")
	registry, err := readBoundsSource(registryPath)
	if err != nil {
		return []string{err.Error()}
	}
	if registry == "" {
		return []string{"internal/bounds policy registry is absent"}
	}
	required := []string{"ModelReadLimit", "OutlineFileLimit", "ControlRecordLimit", "IterationMin", "IterationMax", "MainIterationsDefault", "RefactorIterationsDefault", "MaxWall", "LeaseStale", "AssignmentStale", "PreviewRuneLimit"}
	var diags []string
	owners := map[string][]string{
		"internal/models/models.go":                 {"bounds.VerdictWindow(bounds.ProviderTimeout)", "bounds.ModelReadLimit"},
		"internal/sessioninspect/sessioninspect.go": {"bounds.VerdictWindow(bounds.ProviderTimeout)", "bounds.VerdictWindow(bounds.EnvironmentDiscoveryTimeout)"},
		"internal/outline/read.go":                  {"bounds.OutlineFileLimit"},
		"internal/learnings/learnings.go":           {"bounds.ControlRecordLimit"},
		"internal/maps/maps.go":                     {"bounds.ControlRecordLimit"},
		"internal/roadmap/roadmap.go":               {"bounds.ControlRecordLimit"},
		"internal/guards/guards.go":                 {"bounds.VerdictWindow(bounds.GuardScanTimeout)"},
		"internal/sanitize/sanitize.go":             {"bounds.PreviewRuneLimit"},
		"internal/gate/gate.go":                     {"bounds.VerdictWindow(bounds.GateTimeout)"},
		"internal/coverage/citation_execution.go":   {"bounds.VerdictWindow(bounds.PackageLoadTimeout)"},
		"internal/refresh/refresh.go":               {"bounds.VerdictWindow(bounds.GitRefreshTimeout)"},
		"internal/git/git.go":                       {"bounds.VerdictWindow(bounds.WorktreeListTimeout)", "bounds.VerdictWindow(bounds.RefCheckTimeout)"},
		"internal/intent/intent.go":                 {"bounds.VerdictWindow(bounds.IntentLockTimeout)"},
		"internal/handoffdoc/store.go":              {"bounds.VerdictWindow(bounds.HandoffLockTimeout)"},
		"internal/capturetx/transaction.go":         {"bounds.VerdictWindow(bounds.CaptureLockTimeout)"},
		"internal/worktree/lifecycle.go":            {"bounds.LeaseStale"},
		"internal/worktree/classifier.go":           {"bounds.AssignmentStale"},
		"internal/shift/loop.go":                    {"bounds.MainIterationsDefault", "bounds.RefactorIterationsDefault", "bounds.IterationMin", "bounds.IterationMax", "bounds.MaxWall"},
	}
	for _, tokens := range owners {
		for _, consumer := range tokens {
			if name, verdict := strings.CutPrefix(consumer, "bounds.VerdictWindow(bounds."); verdict {
				required = append(required, strings.TrimSuffix(name, ")"))
			}
		}
	}
	for _, name := range required {
		if !strings.Contains(registry, name) {
			diags = append(diags, "internal/bounds policy registry missing "+name)
		}
	}
	for rel, tokens := range owners {
		body, err := readBoundsSource(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			diags = append(diags, err.Error())
			continue
		}
		for _, token := range tokens {
			if !strings.Contains(body, token) {
				diags = append(diags, rel+" does not consume "+token)
			}
		}
	}
	diags = append(diags, checkBoundCallers(root, registryPath)...)
	wrapper, err := readBoundsSource(filepath.Join(root, "bin", "bench.sh"))
	if err != nil {
		diags = append(diags, err.Error())
	}
	if !strings.Contains(wrapper, `[[ "${BENCH_OFFLINE:-}" == 1 ]]`) || !strings.Contains(registry, `os.Getenv("BENCH_OFFLINE") == "1"`) {
		diags = append(diags, "wrapper and Go offline checks do not share exact BENCH_OFFLINE=1 semantics")
	}
	return diags
}

func checkBoundCallers(root, registryPath string) []string {
	fset := token.NewFileSet()
	registryBytes, err := readBoundsSource(registryPath)
	if err != nil {
		return []string{err.Error()}
	}
	registry, err := parser.ParseFile(fset, registryPath, registryBytes, 0)
	if err != nil {
		return []string{"internal/bounds policy registry is not valid Go: " + err.Error()}
	}
	owned := boundsOwnedExpressions(fset, registry)
	seams, err := boundsReadSeams(filepath.Dir(registryPath))
	if err != nil {
		return []string{"internal/bounds read seams cannot be parsed: " + err.Error()}
	}
	waits := newWaitPolicy(registry)
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
			diags = append(diags, checkBoundCaller(fset, path, rel, owned, seams, waits)...)
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
		path := filepath.Join(dir, name)
		source, err := readBoundsSource(path)
		if err != nil {
			return nil, err
		}
		file, err := parser.ParseFile(fset, path, source, 0)
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

func checkBoundCaller(fset *token.FileSet, path, rel string, ownedExpressions map[string]string, seams map[string]int, waits *waitPolicy) []string {
	source, err := readBoundsSource(path)
	if err != nil {
		return []string{err.Error()}
	}
	if source == "" {
		return nil
	}
	caller, err := parser.ParseFile(fset, path, source, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for bounds ownership: " + err.Error()}
	}
	diags := waits.check(caller, path, rel)
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

func TestBoundsPolicyRejectsSpecialSources(t *testing.T) {
	for _, kind := range []string{"symlink", "fifo"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "internal", "sample")
			requireFixtureNoError(t, os.MkdirAll(dir, 0o755))
			owner := filepath.Join(root, "internal", "bounds")
			requireFixtureNoError(t, os.MkdirAll(owner, 0o755))
			requireFixtureNoError(t, os.WriteFile(filepath.Join(owner, "bounds.go"), []byte("package bounds\n"), 0o600))
			requireFixtureNoError(t, os.WriteFile(filepath.Join(dir, "a.go"), []byte("package sample\n"), 0o600))
			path := filepath.Join(dir, "blocked.go")
			requireFixtureNoError(t, os.WriteFile(path, []byte("package sample\n"), 0o600))
			check := conformanceChecks["bounds-policy"]
			if diagnostics := check.run(root, root, check.tier); containsDiagnostic(diagnostics, "blocked.go") {
				t.Fatalf("regular source refused: %v", diagnostics)
			}
			requireFixtureNoError(t, os.Remove(path))
			if kind == "symlink" {
				target := filepath.Join(t.TempDir(), "target")
				requireFixtureNoError(t, os.WriteFile(target, []byte("package sample\n"), 0o600))
				if err := os.Symlink(target, path); err != nil {
					capability.Capability(t, capability.Symlink, err.Error())
				}
			} else if err := syscall.Mkfifo(path, 0o600); err != nil {
				capability.Capability(t, capability.Fifo, err.Error())
			}
			if diagnostics := check.run(root, root, check.tier); !containsDiagnostic(diagnostics, "blocked.go") {
				t.Fatalf("special source lacks a named refusal: %v", diagnostics)
			}
		})
	}
}

// The classified shapes are the ones a two-leg marker wait used: a fixed window passed to an
// injected sleep, and a fixed window on the right of an injected-clock elapsed compare.
func TestBoundsPolicyClassifiesInjectedAndElapsedWaits(t *testing.T) {
	for _, tc := range []struct{ name, body, diagnostic string }{
		{"injected classified", "sleep(bounds.FixedWindow(10 * time.Millisecond))", ""},
		{"injected raw", "sleep(10 * time.Millisecond)", "unclassified wait in injected duration function"},
		{"elapsed classified", "for now().Sub(started) < bounds.FixedWindow(deadline) {\n\t}", ""},
		{"elapsed raw", "for now().Sub(started) < deadline {\n\t}", "unclassified wait in current-time deadline"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := throwawayRoot{files: map[string]string{
				"internal/bounds/bounds.go": "package bounds\n",
				"internal/sample/wait.go": "package sample\n\nimport (\n\t\"time\"\n\n\t\"github.com/gibbonmi/bench/internal/bounds\"\n)\n\n" +
					"func wait(deadline time.Duration, now func() time.Time, sleep func(time.Duration)) {\n\tstarted := now()\n\t" + tc.body + "\n}\n",
			}}.build(t)
			diagnostics := checkBoundCallers(root, filepath.Join(root, "internal", "bounds", "bounds.go"))
			if tc.diagnostic == "" && len(diagnostics) != 0 {
				t.Fatalf("classified wait read red: %v", diagnostics)
			}
			if tc.diagnostic != "" && !containsDiagnostic(diagnostics, "internal/sample/wait.go has an "+tc.diagnostic) {
				t.Fatalf("raw wait lacks %q: %v", tc.diagnostic, diagnostics)
			}
		})
	}
}

func readBoundsSource(path string) (string, error) {
	file := bounds.ClassifyNoFollow(path)
	if file.State.Failed() {
		return "", fmt.Errorf("%s: %s", path, file.Reason)
	}
	return string(file.Data), nil
}
