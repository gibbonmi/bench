package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// contractSuiteDir is the source tree the contract phase compiles and runs from the kit
// checkout, and so the only tree this check sweeps. The narrowness is deliberate: the
// stripped subject construction closes every other package's reads, and widening the
// sweep would turn a checkable assertion into a prohibition the rest of the kit never
// agreed to.
const contractSuiteDir = "internal/contract"

// subjectResolverName and kitResolverName are the contract package's two root-resolution
// helpers. A test resolves the subject through SubjectRoot; KitRoot names the checkout the
// suite compiles from, which stripping never reaches.
const (
	subjectResolverName = "SubjectRoot"
	kitResolverName     = "KitRoot"
)

// checkContractCaptureReads reports every contract test that resolves a reduced-scope
// path from the kit checkout. The contract phase runs the suite from the kit checkout
// against a separate subject root, and story 4's stripped construction reaches only the
// subject — so a kit-relative read of an allowlisted path would keep reading the real
// tree with nothing able to red it. It reads the suite's Go source through the AST rather
// than by text search, so the forbidden call spelled inside a string literal — this
// check's own bite proof, a diagnostic — is data and never a violation. Only a
// filepath.Join whose root is KitRoot (directly or through a variable assigned from it)
// and whose literal segments name a declared path is one; a path built from dynamic
// segments is invisible to the sweep, which is the honest boundary of a static check.
func checkContractCaptureReads(root string) []string {
	var diags []string
	scope := gate.ReducedScope()
	_ = filepath.WalkDir(filepath.Join(root, filepath.FromSlash(contractSuiteDir)), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		diags = append(diags, contractCaptureReadDiags(scope, slashRel(root, path), path)...)
		return nil
	})
	return uniqueSorted(diags)
}

func contractCaptureReadDiags(scope gate.Scope, rel, path string) []string {
	body := readIfExists(path)
	if !strings.Contains(body, kitResolverName) {
		return nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for kit-relative capture reads: " + err.Error()}
	}
	var diags []string
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		kitVars := kitRootVariables(fn.Body)
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			joined, line, kitRooted := kitRootedJoin(fset, node, kitVars)
			if kitRooted && reducedScopePath(scope, joined) {
				diags = append(diags, fmt.Sprintf(
					"%s resolves reduced-scope path %q from the kit checkout (line %d); resolve it through contract.SubjectRoot so the stripped subject construction can red the dependency",
					rel, joined, line))
			}
			return true
		})
	}
	return diags
}

// kitRootVariables collects the names a function binds to a KitRoot call, so a read
// laundered through `kit := contract.KitRoot(t)` is graded the same as a direct one.
// Tracking is per function: the suite reuses names like `root` across tests for both
// resolvers, and a file-wide set would let one function's binding misgrade another's.
func kitRootVariables(body *ast.BlockStmt) map[string]bool {
	vars := map[string]bool{}
	ast.Inspect(body, func(node ast.Node) bool {
		assign, ok := node.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != len(assign.Rhs) {
			return true
		}
		for i, rhs := range assign.Rhs {
			ident, ok := assign.Lhs[i].(*ast.Ident)
			if !ok {
				continue
			}
			// A rebinding away from KitRoot drops the name: only its current
			// resolver decides how a later join is graded.
			vars[ident.Name] = isResolverCall(rhs, kitResolverName)
		}
		return true
	})
	return vars
}

// kitRootedJoin reports whether node is a filepath.Join resolving from the kit checkout,
// returning the slash-joined literal segments and the call's line when it is.
func kitRootedJoin(fset *token.FileSet, node ast.Node, kitVars map[string]bool) (string, int, bool) {
	call, ok := node.(*ast.CallExpr)
	if !ok || !isPackageCall(call, "filepath", "Join") || len(call.Args) < 2 {
		return "", 0, false
	}
	head := call.Args[0]
	ident, isIdent := head.(*ast.Ident)
	if !isResolverCall(head, kitResolverName) && !(isIdent && kitVars[ident.Name]) {
		return "", 0, false
	}
	var segments []string
	for _, arg := range call.Args[1:] {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			// A dynamic segment ends what the sweep can know about the path.
			break
		}
		segments = append(segments, strings.Trim(lit.Value, "`\""))
	}
	return strings.Join(segments, "/"), fset.Position(call.Pos()).Line, true
}

// isResolverCall reports whether expr calls the named root resolver, through any package
// qualifier or none: inside the contract package the call is bare, and a qualifier
// allowlist a rename or alias evades is no allowlist.
func isResolverCall(expr ast.Expr, name string) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name == name
	case *ast.SelectorExpr:
		return fun.Sel.Name == name
	}
	return false
}

func isPackageCall(call *ast.CallExpr, pkg, name string) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != name {
		return false
	}
	ident, ok := selector.X.(*ast.Ident)
	return ok && ident.Name == pkg
}

// reducedScopePath reports whether a slash-joined literal path is declared: a member of
// the reduced scope, or a declared directory itself, whose listing is as much a read of
// the capture surface as any file under it.
func reducedScopePath(scope gate.Scope, path string) bool {
	if scope.Member(path) {
		return true
	}
	for _, dir := range scope.Directories() {
		if path+"/" == dir {
			return true
		}
	}
	return false
}

// TestContractCaptureReadsCheckBites is the recorded bite proof for
// checkContractCaptureReads (per craft-gate). It runs against a synthetic tree,
// and walks the states that matter: a planted kit-relative read of an
// allowlisted path fires, the same read resolved from the subject root does
// not, a kit-relative read of an unlisted path does not, and the forbidden
// call spelled inside a string literal is data, never a violation.
func TestContractCaptureReadsCheckBites(t *testing.T) {
	write := func(t *testing.T, body string) string {
		t.Helper()
		root := t.TempDir()
		path := filepath.Join(root, "internal", "contract", "runtime", "planted_test.go")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}
	source := func(read string) string {
		return "package runtime\n\nimport (\n\t\"os\"\n\t\"path/filepath\"\n\t\"testing\"\n\n\t\"github.com/gibbonmi/bench/internal/contract\"\n)\n\nfunc TestPlanted(t *testing.T) {\n\tdata, _ := os.ReadFile(" + read + ")\n\t_ = data\n}\n"
	}

	t.Run("kit-relative read of an allowlisted file fires", func(t *testing.T) {
		diags := checkContractCaptureReads(write(t, source(`filepath.Join(contract.KitRoot(t), "ROADMAP.md")`)))
		if len(diags) != 1 || !strings.Contains(diags[0], `"ROADMAP.md"`) ||
			!strings.Contains(diags[0], "internal/contract/runtime/planted_test.go") {
			t.Fatalf("planted kit-relative ROADMAP.md read: want one diagnostic naming it, got %v", diags)
		}
	})

	t.Run("kit-root variable does not launder the read", func(t *testing.T) {
		body := "package runtime\n\nimport (\n\t\"os\"\n\t\"path/filepath\"\n\t\"testing\"\n\n\t\"github.com/gibbonmi/bench/internal/contract\"\n)\n\nfunc TestPlanted(t *testing.T) {\n\tkit := contract.KitRoot(t)\n\tdata, _ := os.ReadFile(filepath.Join(kit, \"capture\", \"learnings.md\"))\n\t_ = data\n}\n"
		diags := checkContractCaptureReads(write(t, body))
		if len(diags) != 1 || !strings.Contains(diags[0], `"capture/learnings.md"`) {
			t.Fatalf("kit-root variable joined with a capture path: want one diagnostic, got %v", diags)
		}
	})

	t.Run("subject-root read of the same path is clean", func(t *testing.T) {
		if diags := checkContractCaptureReads(write(t, source(`filepath.Join(contract.SubjectRoot(t), "ROADMAP.md")`))); len(diags) != 0 {
			t.Fatalf("subject-root ROADMAP.md read: want no diagnostics, got %v", diags)
		}
	})

	t.Run("subject-root variable read is clean", func(t *testing.T) {
		body := "package runtime\n\nimport (\n\t\"os\"\n\t\"path/filepath\"\n\t\"testing\"\n\n\t\"github.com/gibbonmi/bench/internal/contract\"\n)\n\nfunc TestPlanted(t *testing.T) {\n\troot := contract.SubjectRoot(t)\n\tdata, _ := os.ReadFile(filepath.Join(root, \"capture\", \"learnings.md\"))\n\t_ = data\n}\n"
		if diags := checkContractCaptureReads(write(t, body)); len(diags) != 0 {
			t.Fatalf("subject-root variable read: want no diagnostics, got %v", diags)
		}
	})

	t.Run("kit-relative read of an unlisted path is clean", func(t *testing.T) {
		if diags := checkContractCaptureReads(write(t, source(`filepath.Join(contract.KitRoot(t), "bin", "bench.sh")`))); len(diags) != 0 {
			t.Fatalf("kit-relative unlisted read: want no diagnostics, got %v", diags)
		}
	})

	t.Run("forbidden call inside a string literal is data", func(t *testing.T) {
		body := "package runtime\n\n// filepath.Join(contract.KitRoot(t), \"ROADMAP.md\") is also named in a comment.\nconst forbidden = \"filepath.Join(contract.KitRoot(t), \\\"ROADMAP.md\\\")\"\n"
		if diags := checkContractCaptureReads(write(t, body)); len(diags) != 0 {
			t.Fatalf("forbidden text in a literal and a comment: want no diagnostics, got %v", diags)
		}
	})

	t.Run("unparsable contract source is reported", func(t *testing.T) {
		diags := checkContractCaptureReads(write(t, "package runtime\n\nfunc KitRoot broken {\n"))
		if len(diags) != 1 || !strings.Contains(diags[0], "cannot be parsed") {
			t.Fatalf("unparsable source: want one parse diagnostic, got %v", diags)
		}
	})
}
