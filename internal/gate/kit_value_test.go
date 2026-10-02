package gate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// TestKitValueIsTheOneBenchKitRead pins the one kit read. A second inline read, in
// kitRoot or KitDir, is a reader that a caller cannot replace with an explicit kit.
func TestKitValueIsTheOneBenchKitRead(t *testing.T) {
	if readers := benchKitReaders(t, "."); !slices.Equal(readers, []string{"KitValue"}) {
		t.Fatalf("BENCH_KIT readers = %v, want [KitValue]", readers)
	}
}

// TestKitDirReturnsTheBoundKitValue pins KitDir to KitValue. A KitDir that skips
// KitValue answers the executable's parent under a bound kit.
func TestKitDirReturnsTheBoundKitValue(t *testing.T) {
	kit := t.TempDir()
	t.Setenv("BENCH_KIT", kit)

	if got := KitDir(); got != kit {
		t.Fatalf("KitDir() = %q, want the bound kit %q", got, kit)
	}
}

// TestLaneForCommitAtKitReadsTheRootManifest pins the kit argument. A form that reads the
// process kit answers the kit BENCH_KIT names, or the built-in lane when it names none.
func TestLaneForCommitAtKitReadsTheRootManifest(t *testing.T) {
	root := manifestLaneRoot(t)
	kit := t.TempDir()

	lane, err := LaneForCommitAtKit(root, kit)
	if err != nil {
		t.Fatal(err)
	}
	if lane == nil || lane.Kit != kit || lane.Selective {
		t.Fatalf("LaneForCommitAtKit(root, %q) = %+v, want the manifest lane under that kit", kit, lane)
	}
}

// TestLaneForCommitAtKitFallsBackToTheRoot pins the lane fallback. An empty kit makes the
// root its own kit, so a fallback to the executable's parent answers the manifest lane.
func TestLaneForCommitAtKitFallsBackToTheRoot(t *testing.T) {
	root := manifestLaneRoot(t)

	lane, err := LaneForCommitAtKit(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if lane == nil || !lane.Selective {
		t.Fatalf("LaneForCommitAtKit(root, \"\") = %+v, want a selective lane", lane)
	}
}

// TestKitSourceCheckoutAtKitMatchesAnotherSpelling pins the kit argument. A form that
// ignores its kit compares the root with the executable's parent.
func TestKitSourceCheckoutAtKitMatchesAnotherSpelling(t *testing.T) {
	root := t.TempDir()
	link := filepath.Join(t.TempDir(), "kit-link")
	if err := os.Symlink(root, link); err != nil {
		t.Fatal(err)
	}

	if !KitSourceCheckoutAtKit(root, link) {
		t.Fatalf("KitSourceCheckoutAtKit(%q, %q) = false, want true", root, link)
	}
}

// TestKitSourceCheckoutAtKitRefusesAnotherDirectory pins the compare. A form that answers
// true for each resolvable kit counts a consumer repository as the kit.
func TestKitSourceCheckoutAtKitRefusesAnotherDirectory(t *testing.T) {
	root, other := t.TempDir(), t.TempDir()

	if KitSourceCheckoutAtKit(root, other) {
		t.Fatalf("KitSourceCheckoutAtKit(%q, %q) = true, want false", root, other)
	}
}

// TestKitSourceCheckoutAtKitFallsBackToTheExecutableParent pins the kit-source fallback.
// A temporary root is not the executable's parent, so a fallback to the root answers true.
func TestKitSourceCheckoutAtKitFallsBackToTheExecutableParent(t *testing.T) {
	root := t.TempDir()

	if KitSourceCheckoutAtKit(root, "") {
		t.Fatalf("KitSourceCheckoutAtKit(%q, \"\") = true, want false", root)
	}
}

// manifestLaneRoot answers a temporary root whose phase manifest declares a lane.
func manifestLaneRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeLaneFile(t, filepath.Join(root, ".bench", "phases.json"), `{
	  "phases": [{"name": "build", "argv": ["go", "build", "./..."]}],
	  "lane": [{"name": "fmt", "argv": ["make", "fmt"]}]
	}`)
	return root
}

// benchKitReaders names the declaration that holds each os.Getenv or os.LookupEnv call
// with the literal argument BENCH_KIT, across the non-test files of dir. A package-level
// read has no declaration name, so it reports as the empty name.
func benchKitReaders(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var readers []string
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range file.Decls {
			var declName string
			if fn, ok := decl.(*ast.FuncDecl); ok {
				declName = fn.Name.Name
			}
			ast.Inspect(decl, func(node ast.Node) bool {
				if call, ok := node.(*ast.CallExpr); ok && readsBenchKit(call) {
					readers = append(readers, declName)
				}
				return true
			})
		}
	}
	slices.Sort(readers)
	return readers
}

func readsBenchKit(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || (sel.Sel.Name != "Getenv" && sel.Sel.Name != "LookupEnv") || len(call.Args) != 1 {
		return false
	}
	if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "os" {
		return false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	value, err := strconv.Unquote(lit.Value)
	return err == nil && value == "BENCH_KIT"
}
