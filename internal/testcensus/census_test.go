package testcensus

import (
	"go/ast"
	"maps"
	"slices"
	"testing"
)

// TestWalkFindsTopLevelTestsAndTheirOwnParallelCall proves the walk parses only the
// _test.go files, keeps TestMain, a lowercase name, and a helper out of the top-level
// tests, and credits a test with t.Parallel() only for a call outside every closure.
func TestWalkFindsTopLevelTestsAndTheirOwnParallelCall(t *testing.T) {
	t.Parallel()
	dir := Plant(t, map[string]string{
		"source.go": "package planted\n\nimport \"testing\"\n\nfunc TestInSource(t *testing.T) { t.Parallel() }\n",
		"planted_test.go": `package planted

import "testing"

func TestMain(m *testing.M) {}

func TestDirect(t *testing.T) { t.Parallel() }

func TestClosure(t *testing.T) { t.Run("x", func(t *testing.T) { t.Parallel() }) }

func Testlower(t *testing.T) { t.Parallel() }

func helper(t *testing.T) { t.Parallel() }
`,
	})
	files, _, names, err := ParseTestFiles(dir)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !slices.Equal(names, []string{"planted_test.go"}) {
		t.Fatalf("parsed files = %q, want the test file alone", names)
	}
	parallel := map[string]bool{}
	for _, node := range files[0].Decls {
		if decl, ok := node.(*ast.FuncDecl); ok && IsTopLevelTest(decl) {
			parallel[decl.Name.Name] = CallsParallelDirectly(decl)
		}
	}
	if want := map[string]bool{"TestDirect": true, "TestClosure": false}; !maps.Equal(parallel, want) {
		t.Fatalf("top-level tests and their own t.Parallel() call = %v, want %v", parallel, want)
	}
}
