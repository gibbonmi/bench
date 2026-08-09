package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

type buildConstructor struct {
	path string
	kind string
}

var expectedBuildConstructors = map[buildConstructor]int{
	{path: "internal/runbinary/runbinary.go", kind: "subject-builder"}:                                     1,
	{path: "internal/conformance/cross_compile_stress_test.go", kind: "subject-builder"}:                   1,
	{path: "internal/contract/freshness_subject_test.go", kind: "subject-builder"}:                         1,
	{path: "internal/contract/runtime/runtime_gate_freshness_routes_test.go", kind: "subject-builder"}:     2,
	{path: "internal/contract/surface/artifact/posture/builder_contract_test.go", kind: "subject-builder"}: 6,
	{path: "internal/contract/surface/artifact/posture/mode_test.go", kind: "subject-builder"}:             10,
	{path: "internal/contract/surface/artifact/posture/reproducibility_test.go", kind: "subject-builder"}:  1,
	{path: "internal/contract/surface/artifact/posture/trace_test.go", kind: "subject-builder"}:            1,
	{path: "scripts/build-artifacts.sh", kind: "subject-builder"}:                                          2,
	{path: "scripts/native-proof.sh", kind: "subject-builder"}:                                             2,
	{path: "scripts/release-preflight.sh", kind: "subject-builder"}:                                        1,
	{path: "internal/gate/build_attestation_test.go", kind: "go-build"}:                                    1,
	{path: "internal/freshness/freshness_test.go", kind: "subject-builder"}:                                1,
	{path: "internal/contract/runtime/runtime_gate_component_boundary_test.go", kind: "go-build"}:          1,
	{path: "scripts/go-build.sh", kind: "go-build"}:                                                        2,
	{path: "internal/gate/gate_go.go", kind: "go-run-bench"}:                                               1,
	{path: "internal/canary/canary.go", kind: "go-test-c"}:                                                 1,
	{path: "internal/contract/runtime/runtime_gate_partial_proof_test.go", kind: "go-test-c"}:              1,
}

func checkOrdinaryBuildCensus(root string) []string {
	if !exists(filepath.Join(root, "internal", "gate", "phases.go")) {
		return nil
	}
	actual, err := scanBuildConstructors(root)
	if err != nil {
		return []string{"ordinary build census unavailable: " + err.Error()}
	}
	var diags []string
	keys := make(map[buildConstructor]bool, len(actual)+len(expectedBuildConstructors))
	for site := range actual {
		keys[site] = true
	}
	for site := range expectedBuildConstructors {
		keys[site] = true
	}
	ordered := make([]buildConstructor, 0, len(keys))
	for site := range keys {
		ordered = append(ordered, site)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].path == ordered[j].path {
			return ordered[i].kind < ordered[j].kind
		}
		return ordered[i].path < ordered[j].path
	})
	for _, site := range ordered {
		if got, want := actual[site], expectedBuildConstructors[site]; got != want {
			diags = append(diags, fmt.Sprintf("ordinary build census %s %s constructors = %d, want %d", site.path, site.kind, got, want))
		}
	}

	phases := gate.BenchkitPhases(root, root)
	for _, phase := range phases {
		if phase.Name == "build" {
			diags = append(diags, "ordinary phase table contains a build phase")
		}
		joined := strings.Join(phase.Argv, " ")
		for _, forbidden := range []string{"go build", "go run ./cmd/bench", "go run ./internal/freshness/check", "scripts/go-build.sh"} {
			if strings.Contains(joined, forbidden) {
				diags = append(diags, fmt.Sprintf("ordinary phase %s contains forbidden constructor %q", phase.Name, forbidden))
			}
		}
	}
	return diags
}

func scanBuildConstructors(root string) (map[buildConstructor]int, error) {
	sites := map[buildConstructor]int{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "dist", "node_modules", "vendor":
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
		if strings.HasPrefix(rel, "tests/canary/") && strings.Contains(rel, "/files/") {
			return nil
		}
		switch filepath.Ext(path) {
		case ".go":
			return scanGoBuildConstructors(path, rel, sites)
		case ".sh":
			return scanShellBuildConstructors(path, rel, sites)
		default:
			return nil
		}
	})
	return sites, err
}

func scanGoBuildConstructors(path, rel string, sites map[buildConstructor]int) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return err
	}
	for _, declaration := range file.Decls {
		fn, ok := declaration.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		builders := map[string]bool{}
		if fn.Type.Params != nil {
			for _, field := range fn.Type.Params.List {
				for _, name := range field.Names {
					if strings.Contains(strings.ToLower(name.Name), "buildhelper") {
						builders[name.Name] = true
					}
				}
			}
		}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.AssignStmt:
				for i, rhs := range n.Rhs {
					if expressionContains(rhs, "go-build.sh", builders) && i < len(n.Lhs) {
						if id, ok := n.Lhs[i].(*ast.Ident); ok {
							builders[id.Name] = true
						}
					}
				}
			case *ast.ValueSpec:
				for i, value := range n.Values {
					if expressionContains(value, "go-build.sh", builders) && i < len(n.Names) {
						builders[n.Names[i].Name] = true
					}
				}
			}
			return true
		})
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || !isExecutionCall(call.Fun) {
				return true
			}
			switch {
			case expressionContains(call, "go-build.sh", builders):
				sites[buildConstructor{path: rel, kind: "subject-builder"}]++
			case callContainsSequence(call, "go", "build"):
				sites[buildConstructor{path: rel, kind: "go-build"}]++
			case callContainsSequence(call, "go", "test", "-c"):
				sites[buildConstructor{path: rel, kind: "go-test-c"}]++
			case callContainsSequence(call, "go", "run", "./cmd/bench"):
				sites[buildConstructor{path: rel, kind: "go-run-bench"}]++
			case callContainsSequence(call, "go", "run", "./internal/freshness/check"):
				sites[buildConstructor{path: rel, kind: "go-run-freshness"}]++
			}
			return true
		})
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if rel == "internal/gate/gate_go.go" && strings.Contains(string(data), `return append(argv, "run", disableBuildVCS, "./cmd/bench"`) {
		sites[buildConstructor{path: rel, kind: "go-run-bench"}]++
	}
	return nil
}

func scanShellBuildConstructors(path, rel string, sites map[buildConstructor]int) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		switch {
		case strings.Contains(line, "scripts/go-build.sh"):
			sites[buildConstructor{path: rel, kind: "subject-builder"}]++
		case strings.Contains(line, "go run ./cmd/bench"):
			sites[buildConstructor{path: rel, kind: "go-run-bench"}]++
		case strings.Contains(line, "go run ./internal/freshness/check"):
			sites[buildConstructor{path: rel, kind: "go-run-freshness"}]++
		case strings.Contains(line, "go build "):
			sites[buildConstructor{path: rel, kind: "go-build"}]++
		}
	}
	return nil
}

func isExecutionCall(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		switch ident.Name {
		case "runAt", "runAtEnv", "runAtWithInput", "runAtWithTimeout":
			return true
		default:
			return false
		}
	}
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	switch selector.Sel.Name {
	case "Command", "CommandContext", "Run", "RunEnv", "RunEnvSpec", "RunAt", "RunAtWithInput", "RunAtWithTimeout":
		return true
	default:
		return false
	}
}

func expressionContains(node ast.Node, value string, identifiers map[string]bool) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		switch item := n.(type) {
		case *ast.BasicLit:
			if item.Kind == token.STRING && strings.Contains(item.Value, value) {
				found = true
				return false
			}
		case *ast.Ident:
			if identifiers[item.Name] {
				found = true
				return false
			}
		}
		return !found
	})
	return found
}

func callContainsSequence(call *ast.CallExpr, values ...string) bool {
	var literals []string
	ast.Inspect(call, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if ok && literal.Kind == token.STRING {
			literals = append(literals, strings.Trim(literal.Value, "`\""))
		}
		return true
	})
	position := 0
	for _, literal := range literals {
		if position < len(values) && literal == values[position] {
			position++
		}
	}
	return position == len(values)
}

func TestOrdinaryBuildCensusMatchesClosedExceptionSet(t *testing.T) {
	if diags := checkOrdinaryBuildCensus(NewHarness(t).KitRoot); len(diags) != 0 {
		t.Fatalf("ordinary build census diagnostics:\n%s", strings.Join(diags, "\n"))
	}
}

func TestOrdinaryBuildCensusRejectsNewConstructor(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "internal", "ordinary.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	source := `package ordinary
import "os/exec"
func hidden() { _ = exec.Command("go", "build", "./cmd/bench") }
`
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	sites, err := scanBuildConstructors(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := sites[buildConstructor{path: "internal/ordinary.go", kind: "go-build"}]; got != 1 {
		t.Fatalf("new ordinary constructor count = %d, want one", got)
	}
}
