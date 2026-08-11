package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
)

type architectureSite struct {
	path string
	line int
	kind string
}

func checkOrdinaryBuildCensus(root string) []string {
	if !exists(filepath.Join(root, "internal", "gate", "phases.go")) {
		return nil
	}
	var diags []string
	phases := gate.BenchkitPhases(root, root)
	wantNames := []string{"gofmt", "vet", "test", "race", "system", "shellcheck"}
	gotNames := make([]string, 0, len(phases))
	for _, phase := range phases {
		gotNames = append(gotNames, phase.Name)
		switch phase.Name {
		case "test":
			if !slices.Equal(phase.Argv, []string{"go", "test", "-count=1", "./..."}) {
				diags = append(diags, fmt.Sprintf("ordinary test argv = %q", phase.Argv))
			}
		case "race":
			joined := strings.Join(phase.Argv, " ")
			if !strings.HasPrefix(joined, "go test -race -count=1 -v ") || strings.Contains(joined, "internal/systemtest") {
				diags = append(diags, fmt.Sprintf("race argv is not the registry-only driver: %q", phase.Argv))
			}
		case "system":
			if !slices.Equal(phase.Argv, []string{"go", "test", "-count=1", "-tags=system", "./internal/systemtest"}) {
				diags = append(diags, fmt.Sprintf("system argv = %q", phase.Argv))
			}
		}
	}
	if !slices.Equal(gotNames, wantNames) {
		diags = append(diags, fmt.Sprintf("dev phase names = %v, want %v", gotNames, wantNames))
	}

	sites, err := scanArchitecture(root)
	if err != nil {
		return append(diags, "branch-native architecture census unavailable: "+err.Error())
	}
	gitRepos, gateProcesses := 0, 0
	for _, site := range sites {
		if site.kind == "retired-entry" {
			diags = append(diags, site.diagnostic())
			continue
		}
		if strings.HasPrefix(site.path, "internal/systemtest/") {
			continue
		}
		if (site.kind == "repository" || site.kind == "process") && strings.HasPrefix(site.path, "internal/git/") && strings.HasSuffix(site.path, "_test.go") {
			gitRepos++
			continue
		}
		if site.kind == "process" && site.path == "internal/gate/process_group_adapter_test.go" {
			gateProcesses++
			continue
		}
		if architectureOwnedTest(site.path) && (site.kind == "process" || site.kind == "repository" || strings.HasPrefix(site.kind, "nested-go-") || site.kind == "inner-gate") {
			diags = append(diags, site.diagnostic())
		}
	}
	if gitRepos > 1 {
		diags = append(diags, fmt.Sprintf("internal/git repository constructors = %d, want at most 1", gitRepos))
	}
	if gateProcesses > 1 {
		diags = append(diags, fmt.Sprintf("internal/gate controlled process constructors = %d, want at most 1", gateProcesses))
	}
	owner := readIfExists(filepath.Join(root, "internal", "systemtest", "owner_test.go"))
	if got := strings.Count(owner, "strippedJourneyMarker"); got != 2 {
		// One declaration and one assertion are the only source references to the singular journey.
		diags = append(diags, fmt.Sprintf("stripped system journey marker references = %d, want 2", got))
	}
	for _, required := range []string{"func TestMain(", "for range 3", "len(o.repos) != 3", "strippedJourneyMarker"} {
		if !strings.Contains(owner, required) {
			diags = append(diags, fmt.Sprintf("system owner is missing budget assertion %q", required))
		}
	}
	for _, releaseOnly := range registry.ReleaseOnlyPackages {
		if !hasGoFile(filepath.Join(root, filepath.FromSlash(releaseOnly))) {
			diags = append(diags, fmt.Sprintf("ReleaseOnlyPackages names %q, which has no Go source in the tree", releaseOnly))
		}
	}
	for key := range directArchitectureTests {
		if !exists(filepath.Join(root, filepath.FromSlash(key))) {
			diags = append(diags, fmt.Sprintf("directArchitectureTests names %q, which has no file in the tree", key))
		}
	}
	decisionTests, err := filepath.Glob(filepath.Join(root, "internal", "*", "decision_test.go"))
	if err != nil {
		return append(diags, "decision-domain census unavailable: "+err.Error())
	}
	for _, path := range decisionTests {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return append(diags, err.Error())
		}
		rel = filepath.ToSlash(rel)
		if !directArchitectureTests[rel] {
			diags = append(diags, fmt.Sprintf("decision-domain test %q has no directArchitectureTests entry", rel))
		}
	}
	return diags
}

// hasGoFile reports whether path is a directory containing at least one .go
// file, so a stale package rename (an empty leftover directory) still reds.
func hasGoFile(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			return true
		}
	}
	return false
}

var directArchitectureTests = map[string]bool{
	"cmd/bench/command_registry_test.go":         true,
	"internal/adopt/decision_test.go":            true,
	"internal/canary/decision_test.go":           true,
	"internal/freshness/decision_test.go":        true,
	"internal/gate/decision_test.go":             true,
	"internal/preflight/decision_test.go":        true,
	"internal/releasepreflight/decision_test.go": true,
}

func architectureOwnedTest(path string) bool {
	if !strings.HasSuffix(path, "_test.go") {
		return false
	}
	return directArchitectureTests[path] || strings.HasPrefix(path, "internal/contract/") || strings.HasPrefix(path, "internal/canary/") || strings.HasPrefix(path, "internal/git/") || strings.HasPrefix(path, "internal/gate/")
}

func (s architectureSite) diagnostic() string {
	return fmt.Sprintf("branch-native census forbids %s at %s:%d", s.kind, s.path, s.line)
}

func scanArchitecture(root string) ([]architectureSite, error) {
	var sites []architectureSite
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".logs", "dist", "node_modules", "vendor":
				if path != root {
					return filepath.SkipDir
				}
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), "//go:build ship") {
			return nil
		}
		found, err := scanArchitectureGo(rel, data)
		if err != nil {
			return err
		}
		sites = append(sites, found...)
		return nil
	})
	return sites, err
}

func scanArchitectureGo(rel string, data []byte) ([]architectureSite, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, rel, data, 0)
	if err != nil {
		return nil, err
	}
	execAliases := map[string]bool{}
	for _, spec := range file.Imports {
		path, _ := strconv.Unquote(spec.Path.Value)
		if path != "os/exec" {
			continue
		}
		name := "exec"
		if spec.Name != nil {
			name = spec.Name.Name
		}
		execAliases[name] = true
	}
	var sites []architectureSite
	for _, declaration := range file.Decls {
		fn, ok := declaration.(*ast.FuncDecl)
		if !ok {
			continue
		}
		retired := strings.HasPrefix(rel, "internal/contract/") && (fn.Name.Name == "NewFixture" || fn.Name.Name == "CommitAll")
		retired = retired || fn.Name.Name == "buildStrippedSubjectForGeneration" || fn.Name.Name == "strippedWorktree"
		if fn.Recv != nil && (fn.Name.Name == "Bench" || fn.Name.Name == "BenchWrapper") {
			retired = true
		}
		if retired {
			sites = append(sites, architectureSite{path: rel, line: fset.Position(fn.Pos()).Line, kind: "retired-entry"})
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		line := fset.Position(call.Pos()).Line
		name, qualifier := callName(call.Fun)
		if execAliases[qualifier] && (name == "Command" || name == "CommandContext") {
			kind := "process"
			literals := callLiterals(call)
			if containsSequence(literals, "go", "test") {
				kind = "nested-go-test"
			} else if containsSequence(literals, "go", "run") {
				kind = "nested-go-run"
			} else if containsSequence(literals, "git", "init") || containsSequence(literals, "git", "clone") || containsSequence(literals, "git", "worktree") {
				kind = "repository"
			}
			sites = append(sites, architectureSite{path: rel, line: line, kind: kind})
			return true
		}
		if isRepositoryCall(name, callLiterals(call)) {
			sites = append(sites, architectureSite{path: rel, line: line, kind: "repository"})
		}
		return true
	})
	if architectureOwnedTest(rel) && strings.Contains(string(data), "BENCH_CANARY_INNER") {
		sites = append(sites, architectureSite{path: rel, line: 1, kind: "inner-gate"})
	}
	return sites, nil
}

func callName(expr ast.Expr) (name, qualifier string) {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name, ""
	case *ast.SelectorExpr:
		if id, ok := value.X.(*ast.Ident); ok {
			return value.Sel.Name, id.Name
		}
		return value.Sel.Name, ""
	default:
		return "", ""
	}
}

func callLiterals(call *ast.CallExpr) []string {
	var literals []string
	for _, arg := range call.Args {
		ast.Inspect(arg, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if ok && literal.Kind == token.STRING {
				value, _ := strconv.Unquote(literal.Value)
				literals = append(literals, value)
			}
			return true
		})
	}
	return literals
}

func containsSequence(values []string, want ...string) bool {
	position := 0
	for _, value := range values {
		if position < len(want) && value == want[position] {
			position++
		}
	}
	return position == len(want)
}

func isRepositoryCall(name string, literals []string) bool {
	switch name {
	case "CommitAll", "InitRepo", "NewFixture", "NewRepository", "newRepository":
		return true
	case "Run", "Output":
		return containsSequence(literals, "init") || containsSequence(literals, "clone") || containsSequence(literals, "worktree")
	default:
		return false
	}
}

func TestBranchNativeArchitectureCensus(t *testing.T) {
	if diags := checkOrdinaryBuildCensus(NewHarness(t).KitRoot); len(diags) != 0 {
		t.Fatalf("branch-native architecture census diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	t.Run("mutation entries are classified", func(t *testing.T) {
		source := []byte(`package mutation
import "os/exec"
func hidden() {
	_ = exec.Command("bench", "version")
	_ = exec.Command("go", "test", "./...")
	CommitAll()
}
`)
		sites, err := scanArchitectureGo("internal/contract/mutation_test.go", source)
		if err != nil {
			t.Fatal(err)
		}
		got := make([]string, 0, len(sites))
		for _, site := range sites {
			got = append(got, site.kind)
		}
		want := []string{"process", "nested-go-test", "repository"}
		if !slices.Equal(got, want) {
			t.Fatalf("mutation classifications = %v, want %v", got, want)
		}
	})
}
