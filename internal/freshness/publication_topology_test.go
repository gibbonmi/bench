package freshness

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/visibletree"
)

const sealedPublicationEntry = "cmd/bench/freshness_publish.go:freshnessPublish"

const publicationDispatchFile = "cmd/bench/main.go"
const builderFile = "scripts/go-build.sh"

func TestFreshnessPublicationTopology(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	if diags := checkFreshnessPublicationTopology(root); len(diags) != 0 {
		t.Fatalf("freshness publication topology: %v", diags)
	}
}

func TestFreshnessPublicationTopologyBites(t *testing.T) {
	root := t.TempDir()
	writePublicationTopologyFile(t, root, "cmd/bench/freshness_publish.go", `package main
import "github.com/gibbonmi/bench/internal/freshness"
const usage = "freshness-publish"
func freshnessPublish() { _ = freshness.Publish("", "", "", "") }
`)
	writePublicationTopologyFile(t, root, "cmd/bench/main.go", "package main\nvar _ = \"freshness-publish\"\n")
	writePublicationTopologyFile(t, root, "scripts/go-build.sh", "staged freshness-publish root output\n")
	writePublicationTopologyFile(t, root, "internal/artifactpublisher/main.go", `package main
import "github.com/gibbonmi/bench/internal/freshness"
func main() { _ = freshness.Publish("", "", "", "") }
`)
	writePublicationTopologyFile(t, root, "internal/artifactpublisher/library.go", `package artifactpublisher
import "github.com/gibbonmi/bench/internal/freshness"
func publishQuietly() { _ = freshness.Publish("", "", "", "") }
`)
	writePublicationTopologyFile(t, root, "scripts/artifact-publish.sh", "staged freshness-publish root output\n")

	diags := strings.Join(checkFreshnessPublicationTopology(root), "\n")
	for _, want := range []string{
		"unexpected sealed-publication call internal/artifactpublisher/main.go:main",
		"unexpected sealed-publication call internal/artifactpublisher/library.go:publishQuietly",
		"unexpected freshness-publish token owner",
	} {
		if !strings.Contains(diags, want) {
			t.Fatalf("artifact helper diagnostics = %q, want %q", diags, want)
		}
	}
}

// TestFreshnessPublicationTopologyGradesOnlyTheVisibleTree is the bite proof for the
// topology's scope. A token owner under an ignored path is a local copy no commit carries,
// so it earns no diagnostic. The same owner at a visible untracked path earns one.
func TestFreshnessPublicationTopologyGradesOnlyTheVisibleTree(t *testing.T) {
	root := gittest.RepoIgnoring(t, ".logs/")
	owner := "package extra\nvar _ = \"freshness-publish\"\n"

	writePublicationTopologyFile(t, root, ".logs/copy/internal/extra/extra.go", owner)
	if diags := strings.Join(checkFreshnessPublicationTopology(root), "\n"); strings.Contains(diags, ".logs/") {
		t.Fatalf("ignored token owner: diagnostics %q name the ignored copy", diags)
	}
	writePublicationTopologyFile(t, root, "internal/extra/extra.go", owner)
	want := "unexpected freshness-publish token owner internal/extra/extra.go"
	if diags := strings.Join(checkFreshnessPublicationTopology(root), "\n"); !strings.Contains(diags, want) {
		t.Fatalf("visible token owner: diagnostics %q, want %q", diags, want)
	}
}

func checkFreshnessPublicationTopology(root string) []string {
	var diags []string
	calls := map[string]bool{}
	tokenOwners := map[string]bool{}
	expectedTokenOwners := freshnessPublishTokenOwners()
	visible, err := visibletree.Files(root)
	if err != nil {
		return []string{"inspect freshness publication topology: " + err.Error()}
	}
	for _, rel := range visible {
		if skippedPublicationTopologyDir(rel) || strings.HasSuffix(rel, "_test.go") || (!strings.HasSuffix(rel, ".go") && !strings.HasSuffix(rel, ".sh")) {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(rel))
		body, err := os.ReadFile(path)
		if err != nil {
			return []string{"inspect freshness publication topology: " + err.Error()}
		}
		if strings.Contains(string(body), "freshness-publish") {
			tokenOwners[rel] = true
			if !expectedTokenOwners[rel] {
				diags = append(diags, fmt.Sprintf("unexpected freshness-publish token owner %s", rel))
			}
		}
		if strings.HasSuffix(rel, ".sh") && strings.Contains(string(body), "go run") && strings.Contains(string(body), "publish") {
			diags = append(diags, fmt.Sprintf("standalone freshness publisher invocation remains in %s", rel))
		}
		if strings.HasSuffix(rel, ".go") {
			fileCalls, pkg, parseErr := freshnessPublishCalls(path, rel, body)
			if parseErr != nil {
				return []string{"inspect freshness publication topology: " + parseErr.Error()}
			}
			for _, call := range fileCalls {
				calls[call] = true
			}
			if pkg == "main" && strings.HasPrefix(rel, "internal/freshness/") && filepath.ToSlash(filepath.Dir(rel)) != "internal/freshness/check" {
				diags = append(diags, fmt.Sprintf("standalone freshness publisher package remains at %s", filepath.ToSlash(filepath.Dir(rel))))
			}
		}
	}
	expectedCalls := map[string]bool{sealedPublicationEntry: true}
	for call := range calls {
		if !expectedCalls[call] {
			diags = append(diags, "unexpected sealed-publication call "+call)
		}
	}
	for call := range expectedCalls {
		if !calls[call] {
			diags = append(diags, "missing sealed-publication call "+call)
		}
	}
	for owner := range expectedTokenOwners {
		if !tokenOwners[owner] {
			diags = append(diags, "missing freshness-publish token owner "+owner)
		}
	}
	sort.Strings(diags)
	return diags
}

func freshnessPublishTokenOwners() map[string]bool {
	entryFile, _, _ := strings.Cut(sealedPublicationEntry, ":")
	return map[string]bool{entryFile: true, publicationDispatchFile: true, builderFile: true}
}

func skippedPublicationTopologyDir(name string) bool {
	top := strings.Split(name, "/")[0]
	switch top {
	case ".git", ".agents", "capture", "decisions", "dist", "node_modules", "specs", "vendor":
		return true
	default:
		return false
	}
}

func freshnessPublishCalls(path, rel string, body []byte) ([]string, string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, body, 0)
	if err != nil {
		return nil, "", err
	}
	aliases := map[string]bool{}
	for _, imported := range file.Imports {
		path, err := strconv.Unquote(imported.Path.Value)
		if err != nil || path != "github.com/gibbonmi/bench/internal/freshness" {
			continue
		}
		name := "freshness"
		if imported.Name != nil {
			name = imported.Name.Name
		}
		aliases[name] = true
	}
	var calls []string
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || selector.Sel.Name != "Publish" {
				return true
			}
			pkg, ok := selector.X.(*ast.Ident)
			if ok && aliases[pkg.Name] {
				calls = append(calls, rel+":"+fn.Name.Name)
			}
			return true
		})
	}
	return calls, file.Name.Name, nil
}

func writePublicationTopologyFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
