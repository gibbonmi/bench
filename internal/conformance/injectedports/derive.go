package injectedports

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// derivedPort is one port the source sweep found, carrying the site that injected it. A
// red then says where the port entered the package, not only that it did.
type derivedPort struct {
	pkg, name, site string
}

// derivedInjectedPorts reports the injected ports one package's non-test source declares.
//
// A port is a named type whose shape is port-shaped: a non-empty interface, a func type,
// or a struct whose every field is a func. The package must also inject it, by passing it
// as a function parameter or by narrowing to it through a type assertion or type switch.
//
// The assertion arm is not decoration. A capability widened out of an already-injected
// owner is never a parameter at all. A parameter-only rule would therefore leave exactly
// the silently-downgradable capabilities out of the inventory. Exportedness is not part
// of the rule either. A package main declares nothing exported, and its ports are ports
// all the same.
func derivedInjectedPorts(dir, pkg string) ([]derivedPort, []string) {
	fset := token.NewFileSet()
	shaped := map[string]bool{}
	sites := map[string]string{}
	var diags []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []string{fmt.Sprintf("injected port derivation cannot read %s: %v", pkg, err)}
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		path := filepath.Join(dir, name)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			diags = append(diags, fmt.Sprintf("injected port derivation cannot parse %s/%s: %v", pkg, name, err))
			continue
		}
		rel := pkg + "/" + name
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.TypeSpec:
				if portShapedType(typed.Type) {
					shaped[typed.Name.Name] = true
				}
			case *ast.FuncDecl:
				if typed.Type.Params == nil {
					return true
				}
				for _, param := range typed.Type.Params.List {
					recordInjectionSite(sites, param.Type, fmt.Sprintf("a parameter of %s in %s:%d", typed.Name.Name, rel, fset.Position(param.Pos()).Line))
				}
			case *ast.TypeAssertExpr:
				recordInjectionSite(sites, typed.Type, fmt.Sprintf("a type assertion in %s:%d", rel, fset.Position(typed.Pos()).Line))
			case *ast.CaseClause:
				for _, expr := range typed.List {
					recordInjectionSite(sites, expr, fmt.Sprintf("a type switch case in %s:%d", rel, fset.Position(expr.Pos()).Line))
				}
			}
			return true
		})
	}

	var ports []derivedPort
	for name := range shaped {
		if site, injected := sites[name]; injected {
			ports = append(ports, derivedPort{pkg: pkg, name: name, site: site})
		}
	}
	sort.Slice(ports, func(i, j int) bool { return ports[i].name < ports[j].name })
	return ports, diags
}

// recordInjectionSite keeps the first site a name is injected at. A port injected twice
// is still one port, and the first site is the stable one to name in a diagnostic.
func recordInjectionSite(sites map[string]string, expr ast.Expr, site string) {
	name := typeIdentName(expr)
	if name == "" {
		return
	}
	if _, seen := sites[name]; !seen {
		sites[name] = site
	}
}

// typeIdentName reduces a type expression to the package-local name it names, seeing
// through a pointer. A qualified type belongs to another package and is not this
// package's port to register.
func typeIdentName(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.StarExpr:
		return typeIdentName(typed.X)
	}
	return ""
}

// portShapedType reports whether a type declaration has the shape of a port. An empty
// interface constrains nothing and names no producer, so it is not one.
func portShapedType(expr ast.Expr) bool {
	switch typed := expr.(type) {
	case *ast.InterfaceType:
		return typed.Methods != nil && len(typed.Methods.List) > 0
	case *ast.FuncType:
		return true
	case *ast.StructType:
		if typed.Fields == nil || len(typed.Fields.List) == 0 {
			return false
		}
		for _, field := range typed.Fields.List {
			if _, isFunc := field.Type.(*ast.FuncType); !isFunc {
				return false
			}
		}
		return true
	}
	return false
}
