package canary

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProductionCanaryDeclsHaveNoFunctionTypedParameters(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	set := token.NewFileSet()
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(set, entry.Name(), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(node ast.Node) bool {
			if spec, ok := node.(*ast.TypeSpec); ok && spec.Name.Name == "DispatchResult" {
				t.Errorf("retired DispatchResult declaration remains")
			}
			if field, ok := node.(*ast.Field); ok {
				for _, name := range field.Names {
					if name.Name == "Dispatched" {
						t.Errorf("retired Dispatched field remains")
					}
				}
			}
			decl, ok := node.(*ast.FuncDecl)
			if !ok || decl.Type.Params == nil {
				return true
			}
			for _, field := range decl.Type.Params.List {
				if _, ok := field.Type.(*ast.FuncType); ok {
					t.Errorf("%s declares function-typed parameter", decl.Name.Name)
				}
			}
			return true
		})
	}
}
