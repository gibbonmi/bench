package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate"
)

func checkPublishedExecutablePath(root string) []string {
	const owner = "internal/freshness/freshness_verify.go"
	if !exists(filepath.Join(root, "internal/freshness")) {
		return nil
	}
	want, err := publishedExecutableSpelling(filepath.Join(root, owner))
	if err != nil {
		return []string{owner + ": published executable path unavailable: " + err.Error()}
	}
	var diags []string
	for _, site := range []struct{ path, function, pattern string }{
		{".bench/lib/resolve-bench.sh", "bench_rebuild_action", `(?m)^\s*printf [^\n]*bench_shell_quote "\$1/([^"\n]+)"`},
		{"bin/bench.sh", "bench_binary_path", `(?m)^\s*c="\$k/([^"\n]+)"`},
	} {
		body, err := readBuildContractSource(filepath.Join(root, filepath.FromSlash(site.path)))
		if err != nil {
			diags = append(diags, site.path+": published executable path unavailable: "+err.Error())
			continue
		}
		function := regexp.MustCompile(`(?ms)^` + regexp.QuoteMeta(site.function) + `\(\) \{\n(.*?)^\}`).FindSubmatch(body)
		if len(function) != 2 {
			diags = append(diags, site.path+": published executable path function unavailable")
			continue
		}
		matches := regexp.MustCompile(site.pattern).FindAllSubmatch(function[1], -1)
		if len(matches) != 1 {
			diags = append(diags, site.path+": published executable path must have one shell derivation")
			continue
		}
		if got := string(matches[0][1]); got != want {
			diags = append(diags, fmt.Sprintf("%s: published executable path %q differs from freshness.PublishedExecutable %q", site.path, got, want))
		}
	}
	return diags
}

func publishedExecutableSpelling(path string) (string, error) {
	body, err := readBuildContractSource(path)
	if err != nil {
		return "", err
	}
	file, err := parser.ParseFile(token.NewFileSet(), path, body, 0)
	if err != nil {
		return "", err
	}
	for _, declaration := range file.Decls {
		fn, ok := declaration.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "PublishedExecutable" {
			continue
		}
		if fn.Body == nil || len(fn.Body.List) != 1 || len(fn.Type.Params.List) != 1 || len(fn.Type.Params.List[0].Names) != 1 {
			break
		}
		result, ok := fn.Body.List[0].(*ast.ReturnStmt)
		if !ok || len(result.Results) != 1 {
			break
		}
		call, ok := result.Results[0].(*ast.CallExpr)
		if !ok || len(call.Args) < 2 {
			break
		}
		name, qualifier := callName(call.Fun)
		if name != "Join" || importedPackage(file, qualifier) != "path/filepath" {
			break
		}
		root, ok := call.Args[0].(*ast.Ident)
		if !ok || root.Name != fn.Type.Params.List[0].Names[0].Name {
			break
		}
		parts := make([]string, 0, len(call.Args)-1)
		for _, arg := range call.Args[1:] {
			part, ok := stringLiteralValue(arg)
			if !ok {
				return "", fmt.Errorf("PublishedExecutable has a nonliteral path component")
			}
			parts = append(parts, part)
		}
		return filepath.ToSlash(filepath.Join(parts...)), nil
	}
	return "", fmt.Errorf("PublishedExecutable has no readable filepath.Join derivation")
}

func importedPackage(file *ast.File, qualifier string) string {
	for _, imp := range file.Imports {
		path, ok := stringLiteralValue(imp.Path)
		if !ok {
			continue
		}
		name := filepath.Base(path)
		if imp.Name != nil {
			name = imp.Name.Name
		}
		if name == qualifier || name == "." && qualifier == "" {
			return path
		}
	}
	return ""
}

func TestPublishedExecutablePathRejectsDrift(t *testing.T) {
	root := throwawayRoot{files: map[string]string{
		"internal/freshness/freshness_verify.go": `package freshness
import "path/filepath"
func PublishedExecutable(root string) string { return filepath.Join(root, "dist", "bench") }
`,
		".bench/lib/resolve-bench.sh": `bench_rebuild_action() {
 printf '%s' "$(bench_shell_quote "$1/wrong/bench")"
}
`,
		"bin/bench.sh": `bench_binary_path() {
 c="$k/dist/bench"
}
`,
	}}.build(t)
	diags := checkPublishedExecutablePath(root)
	if len(diags) != 1 || !strings.Contains(diags[0], ".bench/lib/resolve-bench.sh") {
		t.Fatalf("path drift diagnostics = %v, want resolver mismatch", diags)
	}
}

func checkGoBuildVCS(root string) []string {
	var diags []string
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
		if filepath.Ext(path) != ".go" {
			return nil
		}
		body, err := readBuildContractSource(path)
		if err != nil {
			return err
		}
		rel := slashRel(root, path)
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, body, 0)
		if err != nil {
			diags = append(diags, rel+": Go build VCS scan cannot parse source: "+err.Error())
			return nil
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			name, qualifier := callName(call.Fun)
			if importedPackage(file, qualifier) != "os/exec" || name != "Command" && name != "CommandContext" {
				return true
			}
			args := call.Args
			if name == "CommandContext" && len(args) > 0 {
				args = args[1:]
			}
			if len(args) < 2 {
				return true
			}
			tool, _ := stringLiteralValue(args[0])
			operation, _ := stringLiteralValue(args[1])
			if tool != "go" || operation != "build" && operation != "list" {
				return true
			}
			protected, conflicting := false, false
			for _, arg := range args[2:] {
				if value, ok := stringLiteralValue(arg); ok {
					if value == gate.DisableBuildVCS {
						protected = true
					}
					if strings.HasPrefix(value, "-buildvcs=") && value != gate.DisableBuildVCS {
						conflicting = true
					}
				}
				symbol, pkg := callName(arg)
				if symbol == "DisableBuildVCS" && (importedPackage(file, pkg) == "github.com/gibbonmi/bench/internal/gate" || pkg == "" && file.Name.Name == "gate") {
					protected = true
				}
			}
			if !protected || conflicting {
				diags = append(diags, fmt.Sprintf("%s:%d: go %s must carry -buildvcs=false or gate.DisableBuildVCS without a conflicting VCS flag", rel, fset.Position(call.Pos()).Line, operation))
			}
			return true
		})
		return nil
	})
	if err != nil {
		diags = append(diags, "Go build VCS scan unavailable: "+err.Error())
	}
	return diags
}

func TestGoBuildVCSRejectsMissingFlag(t *testing.T) {
	root := throwawayRoot{files: map[string]string{
		"internal/example/build_test.go": `package example
import "os/exec"
func build() { _ = exec.Command("go", "build", "./cmd/bench") }
`,
	}}.build(t)
	diags := checkGoBuildVCS(root)
	if len(diags) != 1 || !strings.Contains(diags[0], "build_test.go") || !strings.Contains(diags[0], "-buildvcs=false") {
		t.Fatalf("VCS flag diagnostics = %v, want missing flag at build_test.go", diags)
	}
}

func TestPublishedExecutablePathSourceFamily(t *testing.T) {
	files := map[string]string{
		"internal/freshness/freshness_verify.go": `package freshness
import "path/filepath"
func PublishedExecutable(root string) string { return filepath.Join(root, "out", "tool") }
`,
		".bench/lib/resolve-bench.sh": "bench_rebuild_action() {\n printf '%s' \"$(bench_shell_quote \"$1/out/tool\")\"\n}\n",
		"bin/bench.sh":                "bench_binary_path() {\n c=\"$k/out/tool\"\n}\n",
	}
	root := throwawayRoot{files: files}.build(t)
	if diags := checkPublishedExecutablePath(root); len(diags) != 0 {
		t.Fatalf("matching paths: %v", diags)
	}
	for path, body := range files {
		t.Run(path, func(t *testing.T) {
			replacement := strings.ReplaceAll(body, "out", "changed")
			target := filepath.Join(root, filepath.FromSlash(path))
			if err := os.WriteFile(target, []byte(replacement), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkPublishedExecutablePath(root); len(diags) == 0 {
				t.Fatal("changed source path passed")
			}
			if err := os.WriteFile(target, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkPublishedExecutablePath(root); len(diags) != 0 {
				t.Fatalf("restored source: %v", diags)
			}
		})
	}
	if diags := checkPublishedExecutablePath(t.TempDir()); len(diags) != 0 {
		t.Fatalf("absent surface: %v", diags)
	}
}

func TestGoBuildVCSCallFamily(t *testing.T) {
	for _, tc := range []struct {
		name, imports, call string
		red                 bool
	}{
		{"build literal", `"os/exec"`, `exec.Command("go", "build", "-buildvcs=false", "./...")`, false},
		{"list missing", `"os/exec"`, `exec.Command("go", "list", "./...")`, true},
		{"build alias", `process "os/exec"`, `process.Command("go", "build", "./...")`, true},
		{"context list", `"os/exec"`, `exec.CommandContext(ctx, "go", "list", "./...")`, true},
		{"gate constant", `"os/exec"; oracle "github.com/gibbonmi/bench/internal/gate"`, `exec.Command("go", "list", oracle.DisableBuildVCS, "./...")`, false},
		{"wrong gate", `"os/exec"; gate "example.com/gate"`, `exec.Command("go", "list", gate.DisableBuildVCS, "./...")`, true},
		{"dot exec", `. "os/exec"`, `Command("go", "build", "./...")`, true},
		{"conflict", `"os/exec"`, `exec.Command("go", "build", "-buildvcs=false", "-buildvcs=true", "./...")`, true},
		{"other Go verb", `"os/exec"`, `exec.Command("go", "version")`, false},
		{"other executable", `"os/exec"`, `exec.Command("git", "build")`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := throwawayRoot{files: map[string]string{
				"hidden/child/build_test.go": "package sample\nimport (" + tc.imports + ")\nfunc run() { _ = " + tc.call + " }\n",
			}}.build(t)
			diags := checkGoBuildVCS(root)
			if got := len(diags) != 0; got != tc.red {
				t.Fatalf("diagnostics = %v, want red=%t", diags, tc.red)
			}
		})
	}
}

func readBuildContractSource(path string) ([]byte, error) {
	classified := bounds.ClassifyNoFollow(path)
	if classified.State != bounds.StateParsed && classified.State != bounds.StateEmpty {
		return nil, fmt.Errorf("%s: %s: %s", path, classified.State, classified.Reason)
	}
	return classified.Data, nil
}

func TestBuildContractChecksRefuseSpecialFiles(t *testing.T) {
	for _, tc := range []struct {
		path  string
		check func(string) []string
	}{
		{"internal/freshness/freshness_verify.go", checkPublishedExecutablePath},
		{".bench/lib/resolve-bench.sh", checkPublishedExecutablePath},
		{"bin/bench.sh", checkPublishedExecutablePath},
		{"internal/example/probe.go", checkGoBuildVCS},
	} {
		t.Run(tc.path, func(t *testing.T) {
			files := map[string]string{
				"internal/freshness/freshness_verify.go": `package freshness
import "path/filepath"
func PublishedExecutable(root string) string { return filepath.Join(root, "dist", "bench") }
`,
				".bench/lib/resolve-bench.sh": "bench_rebuild_action() {\n printf '%s' \"$(bench_shell_quote \"$1/dist/bench\")\"\n}\n",
				"bin/bench.sh":                "bench_binary_path() {\n c=\"$k/dist/bench\"\n}\n",
			}
			delete(files, tc.path)
			root := throwawayRoot{files: files, plants: map[string]func(*testing.T, string){tc.path: hostileSkillPlanters["fifo"]}}.build(t)
			diags := strings.Join(tc.check(root), "\n")
			if !strings.Contains(diags, tc.path) || !strings.Contains(diags, "wrong-type") {
				t.Fatalf("special-file refusal = %q", diags)
			}
		})
	}
}
