package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// canonicalPathOwner is the one package allowed to derive a path's canonical spelling.
const canonicalPathOwner = "internal/canonicalpath"

// checkCanonicalPathOwner grades the single-source rule for the canonical-path derivation:
// absolute, then symlink-resolved, then cleaned. A production function that calls both
// filepath.Abs and filepath.EvalSymlinks re-derives that spelling beside the owner, and the
// two copies drift on the not-yet-existing path. Test files are exempt, so a fixture can
// plant either call.
//
// The unit is the function, not the file, because unrelated link work shares a file with an
// absolute-path derivation in internal/gate/subject.go.
func checkCanonicalPathOwner(root string) []string {
	var diags []string
	for _, top := range []string{"cmd", "internal"} {
		base := filepath.Join(root, top)
		_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel := slashRel(root, path)
			if strings.HasPrefix(rel, canonicalPathOwner+"/") {
				return nil
			}
			diags = append(diags, canonicalPathDerivations(path, rel)...)
			return nil
		})
	}
	return uniqueSorted(diags)
}

func canonicalPathDerivations(path, rel string) []string {
	if readIfExists(path) == "" {
		return nil
	}
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return []string{rel + " cannot be parsed for canonical-path ownership: " + err.Error()}
	}
	var diags []string
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		calls := map[string]bool{}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := selector.X.(*ast.Ident); ok && pkg.Name == "filepath" {
					calls[selector.Sel.Name] = true
				}
			}
			return true
		})
		if calls["Abs"] && calls["EvalSymlinks"] {
			diags = append(diags, rel+" derives the canonical path in "+fn.Name.Name+" with filepath.Abs and filepath.EvalSymlinks instead of "+canonicalPathOwner)
		}
	}
	return diags
}
