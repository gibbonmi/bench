//go:build system

package systemtest

// This file owns the process-state rules of the system parallel census. Parallel tests
// share one process, so a parallel test changes no process environment and no working
// directory, and no function binds TMPDIR in the process.

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

// processStateCalls name the os functions that change state that every test in the
// process shares.
var processStateCalls = map[string]bool{"Setenv": true, "Unsetenv": true, "Chdir": true}

// processStateCall returns the name of the os function that call runs when that function
// is in processStateCalls, and an empty string otherwise.
func processStateCall(call *ast.CallExpr) string {
	method, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || !processStateCalls[method.Sel.Name] {
		return ""
	}
	if pkg, ok := method.X.(*ast.Ident); !ok || pkg.Name != "os" {
		return ""
	}
	return "os." + method.Sel.Name
}

// tempRootBindings reports every function in the package that binds TMPDIR in the process.
// Parallel tests share the process environment, so a test gives each child its private
// temporary root as an explicit environment value instead.
func (p systemPackage) tempRootBindings() []string {
	var reports []string
	for _, fns := range p.funcs {
		for _, fn := range fns {
			ast.Inspect(fn.decl.Body, func(node ast.Node) bool {
				if call, ok := node.(*ast.CallExpr); ok && bindsTempRoot(call) {
					reports = append(reports, fmt.Sprintf("%s: %s binds TMPDIR in the process", p.where(fn, call.Pos()), fn.decl.Name.Name))
				}
				return true
			})
		}
	}
	return reports
}

// bindsTempRoot reports whether call is a Setenv call on any receiver whose name argument
// can be TMPDIR. The census reads no values, so a name argument that is not a string
// literal can be TMPDIR.
func bindsTempRoot(call *ast.CallExpr) bool {
	method, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || method.Sel.Name != "Setenv" || len(call.Args) != 2 {
		return false
	}
	name, ok := call.Args[0].(*ast.BasicLit)
	if !ok || name.Kind != token.STRING {
		return true
	}
	value, err := strconv.Unquote(name.Value)
	return err == nil && value == "TMPDIR"
}
