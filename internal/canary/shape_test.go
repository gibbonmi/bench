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

func TestCanaryMarkerReadersUseDeclaredNames(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	set := token.NewFileSet()
	readers := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(set, entry.Name(), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || len(call.Args) < 2 {
				return true
			}
			function, ok := call.Fun.(*ast.Ident)
			if !ok || function.Name != "readMarker" {
				return true
			}
			readers++
			if _, ok := call.Args[1].(*ast.Ident); !ok {
				t.Errorf("%s passes a marker name that has no declared owner", set.Position(call.Pos()))
			}
			return true
		})
	}
	if readers == 0 {
		t.Fatal("production canary package has no marker reader")
	}
}
