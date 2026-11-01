package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// processLifetimeOwner is the package that owns process lifetime mechanics. Its tree,
// including the protection leaf and the process test fixtures, is exempt from the closure.
const processLifetimeOwner = "internal/subprocess"

// processLifetimeOwnerImport is the import path that the owner calls resolve through.
const processLifetimeOwnerImport = benchModule + processLifetimeOwner

// processLifetimeOwnerCalls are the shared owner entries that a census caller must reach.
var processLifetimeOwnerCalls = []string{"NewLifetime", "RunInherited"}

// processLifetimeSignalPackages are the packages whose Kill sends a signal and whose Signal
// type converts a signal number.
var processLifetimeSignalPackages = []string{"syscall", "golang.org/x/sys/unix"}

// processLifetimeSetpgidAdapters are the platform adapter files that keep the builder's own
// group attributes. The closure permits Setpgid there and nowhere else outside the owner.
const processLifetimeSetpgidAdapters = "internal/runbinary/sysprocattr_*.go"

// processLifetimeEntry is one accepted census caller: its name, the production file, and the
// function that must reach the shared lifetime owner through same-package calls.
type processLifetimeEntry struct {
	caller, file, function string
}

// processLifetimeCensus is the complete accepted caller census.
var processLifetimeCensus = []processLifetimeEntry{
	{"gate runner", "internal/gate/runner.go", "runProcessGroupCommand"},
	{"runbinary builder", "internal/runbinary/runbinary.go", "runBuildScript"},
	{"bounds runner", "internal/bounds/bounds.go", "run"},
	{"prep-release step", "internal/preprelease/preprelease.go", "runStep"},
	{"external preflight", "internal/releasepreflight/command.go", "runExternal"},
	{"vulnerability scanner", "internal/releasepreflight/vulnerability.go", "runVulnerability"},
	{"focused test runner", "internal/testreport/command.go", "runGoTest"},
	{"package-list runner", "internal/testreport/selection.go", "currentPackagesWithEnvironment"},
	{"worktree exec", "internal/worktree/exec.go", "runChild"},
	{"worktree shell", "internal/worktree/subshell.go", "subshellAt"},
	{"shift adapter", "internal/shift/session.go", "runAdapter"},
	{"background launch", "internal/gate/background_launch.go", "startGrade"},
}

// checkProcessLifetimeClosure grades the final migration closure of the live tree: every
// census entry reaches the shared owner, every owner call sits under a census entry, and no
// production file outside the owner keeps a private lifetime loop.
func checkProcessLifetimeClosure(root string) []string {
	return processLifetimeClosureDiags(root, processLifetimeCensus)
}

// processLifetimeClosureDiags applies the closure to root under census.
func processLifetimeClosureDiags(root string, census []processLifetimeEntry) []string {
	reached := map[string]bool{}
	var diags []string
	for _, entry := range census {
		diags = append(diags, processLifetimeEntryDiags(root, entry, reached)...)
	}
	for _, top := range []string{"cmd", "internal"} {
		_ = filepath.WalkDir(filepath.Join(root, top), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel := slashRel(root, path)
			if strings.HasPrefix(rel, processLifetimeOwner+"/") {
				return nil
			}
			diags = append(diags, processLifetimeFileDiags(path, rel, reached)...)
			return nil
		})
	}
	return uniqueSorted(diags)
}

// processLifetimeEntryDiags reports a census entry whose function is absent or reaches neither
// owner call. It marks every function the entry reaches, keyed by directory and name, so the
// reverse direction can find an owner call that no entry reaches.
func processLifetimeEntryDiags(root string, entry processLifetimeEntry, reached map[string]bool) []string {
	dir := path.Dir(entry.file)
	functions, err := processLifetimePackageFunctions(root, dir)
	if err != nil {
		return []string{"process-lifetime census entry " + entry.caller + ": " + err.Error()}
	}
	candidates := functions.named(entry.function)
	start := slices.IndexFunc(candidates, func(candidate *packageFunc) bool { return candidate.rel == entry.file })
	if start < 0 {
		return []string{fmt.Sprintf("process-lifetime census entry %s: %s declares no %s", entry.caller, entry.file, entry.function)}
	}
	owner := false
	visit := func(current *packageFunc) bool {
		reached[dir+"#"+current.decl.Name.Name] = true
		ast.Inspect(current.body, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok && slices.Contains(processLifetimeOwnerCalls, guardProjectionImported(current.file, call, processLifetimeOwnerImport)) {
				owner = true
			}
			return true
		})
		return false
	}
	visit(candidates[start])
	functions.reaches(candidates[start], visit)
	if !owner {
		return []string{fmt.Sprintf("process-lifetime census entry %s: %s %s does not reach subprocess.NewLifetime or subprocess.RunInherited", entry.caller, entry.file, entry.function)}
	}
	return nil
}

// processLifetimePackageFunctions parses every production file of dir beneath root, under
// every build constraint, and indexes its functions.
func processLifetimePackageFunctions(root, dir string) (packageFuncs, error) {
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return nil, err
	}
	functions := packageFuncs{}
	fset := token.NewFileSet()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		rel := path.Join(dir, name)
		abs := filepath.Join(root, filepath.FromSlash(rel))
		file, err := parser.ParseFile(fset, abs, readIfExists(abs), 0)
		if err != nil {
			return nil, err
		}
		functions.add(fset, rel, file)
	}
	return functions, nil
}

// packageCallee is a key of a package function index: a function name, a method name, or a
// closure variable object. The object keeps two closures of one name apart.
type packageCallee struct {
	name    string
	method  bool
	closure *ast.Object
}

// packageFunc is one function, method, or assigned closure of a package. decl is the
// declaration that holds it, and callees lists each callee that its body can select.
type packageFunc struct {
	name, rel string
	line      int
	file      *ast.File
	decl      *ast.FuncDecl
	body      *ast.BlockStmt
	callees   []packageCallee
}

// packageFuncs indexes the functions, methods, and assigned closures of one package by callee.
type packageFuncs map[packageCallee][]*packageFunc

// add indexes each function, method, and assigned closure of file, which sits at rel.
func (index packageFuncs) add(fset *token.FileSet, rel string, file *ast.File) {
	put := func(key packageCallee, name string, pos token.Pos, decl *ast.FuncDecl, body *ast.BlockStmt) {
		index[key] = append(index[key], &packageFunc{name, rel, fset.Position(pos).Line, file, decl, body, packageCallees(file, body)})
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		put(packageCallee{name: fn.Name.Name, method: fn.Recv != nil}, fn.Name.Name, fn.Pos(), fn, fn.Body)
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			if assign, ok := node.(*ast.AssignStmt); ok && len(assign.Lhs) == len(assign.Rhs) {
				for i, lhs := range assign.Lhs {
					ident, named := lhs.(*ast.Ident)
					literal, closure := assign.Rhs[i].(*ast.FuncLit)
					if named && closure && ident.Obj != nil {
						put(packageCallee{closure: ident.Obj}, ident.Name, literal.Pos(), fn, literal.Body)
					}
				}
			}
			return true
		})
	}
}

// packageCallees lists, once each, the callee that each call or function value of body can select.
// A selection through an import qualifier selects another package. A composite literal key names a field.
func packageCallees(file *ast.File, body ast.Node) []packageCallee {
	var callees []packageCallee
	skip := map[*ast.Ident]bool{}
	add := func(key packageCallee) {
		if !slices.Contains(callees, key) {
			callees = append(callees, key)
		}
	}
	ast.Inspect(body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.KeyValueExpr:
			if key, ok := n.Key.(*ast.Ident); ok {
				skip[key] = true
			}
		case *ast.SelectorExpr:
			skip[n.Sel] = true
			if qualifier, ok := n.X.(*ast.Ident); ok && qualifier.Obj == nil && importedPackage(file, qualifier.Name) != "" {
				return false
			}
			add(packageCallee{name: n.Sel.Name, method: true})
		case *ast.Ident:
			switch {
			case skip[n]:
			case n.Obj == nil || n.Obj.Kind == ast.Fun:
				add(packageCallee{name: n.Name})
			case n.Obj.Kind == ast.Var:
				add(packageCallee{closure: n.Obj})
			}
		}
		return true
	})
	return callees
}

// named lists the functions and then the methods of index that a declared name selects.
func (index packageFuncs) named(name string) []*packageFunc {
	return slices.Concat(index[packageCallee{name: name}], index[packageCallee{name: name, method: true}])
}

// reaches reports whether a callee that from selects, directly or through other callees, is done.
// It offers each reached function to done once.
func (index packageFuncs) reaches(from *packageFunc, done func(*packageFunc) bool) bool {
	seen := map[*packageFunc]bool{}
	for stack := []*packageFunc{from}; len(stack) > 0; {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, key := range node.callees {
			for _, next := range index[key] {
				if seen[next] {
					continue
				}
				if seen[next] = true; done(next) {
					return true
				}
				stack = append(stack, next)
			}
		}
	}
	return false
}

// processLifetimeFixtureCensus is the one-entry census of the synthetic trees.
var processLifetimeFixtureCensus = []processLifetimeEntry{{"example caller", "internal/example/run.go", "run"}}

// processLifetimeOwned is a census caller that reaches the owner through a same-package method.
const processLifetimeOwned = `package example

import (
	"context"
	"os/exec"

	"github.com/gibbonmi/bench/internal/subprocess"
)

type runner struct{}

func run(ctx context.Context, cmd *exec.Cmd) error { return runner{}.launch(ctx, cmd) }

func (runner) launch(ctx context.Context, cmd *exec.Cmd) error {
	_, err := subprocess.RunInherited(ctx, cmd, subprocess.Policy{}, subprocess.Streams{}, subprocess.Interrupt(15))
	return err
}
`

// TestProcessLifetimeOwnerEntries is PL1. It grades the census half on synthetic trees: an
// entry that reaches the owner through a same-package call passes; an entry whose caller
// keeps a private launch, an absent entry function, and an owner call that no entry reaches
// each turn red.
func TestProcessLifetimeOwnerEntries(t *testing.T) {
	root := throwawayRoot{files: map[string]string{"internal/example/run.go": processLifetimeOwned}}.build(t)
	if diags := processLifetimeClosureDiags(root, processLifetimeFixtureCensus); len(diags) != 0 {
		t.Fatalf("owner reached through a same-package method: want no diagnostics, got %v", diags)
	}

	// The caller keeps its own launch: the owner call survives in a helper the entry no longer calls.
	bypass := strings.Replace(processLifetimeOwned, "return runner{}.launch(ctx, cmd)", "return cmd.Run()", 1)
	root = throwawayRoot{files: map[string]string{"internal/example/run.go": bypass}}.build(t)
	diags := processLifetimeClosureDiags(root, processLifetimeFixtureCensus)
	want := []string{
		"internal/example/run.go:15 calls subprocess.RunInherited from launch, which no process-lifetime census entry reaches",
		"process-lifetime census entry example caller: internal/example/run.go run does not reach subprocess.NewLifetime or subprocess.RunInherited",
	}
	if !slices.Equal(diags, want) {
		t.Fatalf("bypassed caller: got %q, want %q", diags, want)
	}

	// A renamed entry function leaves the census naming a caller that no longer exists.
	root = throwawayRoot{files: map[string]string{"internal/example/run.go": strings.Replace(processLifetimeOwned, "func run(", "func runRenamed(", 1)}}.build(t)
	diags = processLifetimeClosureDiags(root, processLifetimeFixtureCensus)
	if !slices.Contains(diags, "process-lifetime census entry example caller: internal/example/run.go declares no run") {
		t.Fatalf("absent entry function: got %q", diags)
	}

	// An aliased owner import is the same owner.
	aliased := strings.NewReplacer("\t\"github.com/gibbonmi/bench/internal/subprocess\"", "\tlife \"github.com/gibbonmi/bench/internal/subprocess\"", "subprocess.", "life.").Replace(processLifetimeOwned)
	root = throwawayRoot{files: map[string]string{"internal/example/run.go": aliased}}.build(t)
	if diags := processLifetimeClosureDiags(root, processLifetimeFixtureCensus); len(diags) != 0 {
		t.Fatalf("aliased owner import: want no diagnostics, got %v", diags)
	}
}
