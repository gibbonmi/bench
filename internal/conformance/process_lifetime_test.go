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
	functions, err := processLifetimePackageFunctions(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return []string{"process-lifetime census entry " + entry.caller + ": " + err.Error()}
	}
	start := slices.IndexFunc(functions[entry.function], func(candidate processLifetimeDeclaration) bool {
		return slashRel(root, candidate.path) == entry.file
	})
	if start < 0 {
		return []string{fmt.Sprintf("process-lifetime census entry %s: %s declares no %s", entry.caller, entry.file, entry.function)}
	}
	owner := false
	queue := []processLifetimeDeclaration{functions[entry.function][start]}
	seen := map[*ast.FuncDecl]bool{queue[0].decl: true}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		reached[dir+"#"+current.decl.Name.Name] = true
		ast.Inspect(current.decl.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			if slices.Contains(processLifetimeOwnerCalls, guardProjectionImported(current.file, call, processLifetimeOwnerImport)) {
				owner = true
			}
			for _, next := range functions[processLifetimeCallee(current.file, call)] {
				if !seen[next.decl] {
					seen[next.decl] = true
					queue = append(queue, next)
				}
			}
			return true
		})
	}
	if !owner {
		return []string{fmt.Sprintf("process-lifetime census entry %s: %s %s does not reach subprocess.NewLifetime or subprocess.RunInherited", entry.caller, entry.file, entry.function)}
	}
	return nil
}

// processLifetimeCallee names the same-package function or method that a call can select. A
// selection through an import qualifier names another package, so it selects nothing here.
func processLifetimeCallee(file *ast.File, call *ast.CallExpr) string {
	switch fun := ast.Unparen(call.Fun).(type) {
	case *ast.Ident:
		return fun.Name
	case *ast.SelectorExpr:
		if ident, ok := fun.X.(*ast.Ident); ok && ident.Obj == nil && importedPackage(file, ident.Name) != "" {
			return ""
		}
		return fun.Sel.Name
	}
	return ""
}

// processLifetimeDeclaration is one function declaration with the file that holds it.
type processLifetimeDeclaration struct {
	path string
	file *ast.File
	decl *ast.FuncDecl
}

// processLifetimePackageFunctions parses every production file of dir, under every build
// constraint, and indexes its function and method declarations by name.
func processLifetimePackageFunctions(dir string) (map[string][]processLifetimeDeclaration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	functions := map[string][]processLifetimeDeclaration{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(dir, name)
		file, err := parser.ParseFile(token.NewFileSet(), path, readIfExists(path), 0)
		if err != nil {
			return nil, err
		}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil {
				functions[fn.Name.Name] = append(functions[fn.Name.Name], processLifetimeDeclaration{path, file, fn})
			}
		}
	}
	return functions, nil
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
