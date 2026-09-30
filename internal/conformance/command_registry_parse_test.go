package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// worktreeLeafFile is the one file that declares the `bench worktree` family, and
// worktreeLeafTable is the identifier of the composite literal that carries its leaves.
const (
	worktreeLeafFile  = "cmd/bench/worktree_leaves.go"
	worktreeLeafTable = "worktreeLeaves"
)

type commandRegistryEntry struct {
	name   string
	fields map[string][]ast.Expr
}

// parseCommandRegistry owns the syntax shared by registry-backed conformance checks. It
// locates the named variable exactly once and preserves entry and repeated-field order.
// Each check stays responsible for the meaning of the fields it consumes.
func parseCommandRegistry(path, body string) ([]commandRegistryEntry, error) {
	return parseRegistryTable(path, body, dispatchRegistry)
}

// parseWorktreeLeaves reads the `bench worktree` leaf table with the same rules as the
// command registry, so a leaf row and a command row have one syntax.
func parseWorktreeLeaves(path, body string) ([]commandRegistryEntry, error) {
	return parseRegistryTable(path, body, worktreeLeafTable)
}

func parseRegistryTable(path, body, table string) ([]commandRegistryEntry, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, body, 0)
	if err != nil {
		return nil, err
	}
	type declaration struct {
		tok  token.Token
		expr ast.Expr
	}
	var declarations []declaration
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range value.Names {
				if name.Name != table {
					continue
				}
				if len(value.Names) != 1 || len(value.Values) != 1 || i >= len(value.Values) {
					return nil, fmt.Errorf("%s must be one named value with one literal", table)
				}
				declarations = append(declarations, declaration{tok: genDecl.Tok, expr: value.Values[i]})
			}
		}
	}
	if len(declarations) != 1 {
		return nil, fmt.Errorf("found %d %s declarations, want exactly 1", len(declarations), table)
	}
	if declarations[0].tok != token.VAR {
		return nil, fmt.Errorf("%s must be declared with var", table)
	}
	literal, ok := declarations[0].expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("%s is not a composite literal", table)
	}
	entries := make([]commandRegistryEntry, 0, len(literal.Elts))
	seen := make(map[string]bool, len(literal.Elts))
	for i, element := range literal.Elts {
		entry, ok := element.(*ast.CompositeLit)
		if !ok {
			return nil, fmt.Errorf("registry entry %d is not a composite literal", i+1)
		}
		fields := make(map[string][]ast.Expr)
		for _, field := range entry.Elts {
			pair, ok := field.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := pair.Key.(*ast.Ident)
			if !ok {
				continue
			}
			fields[key.Name] = append(fields[key.Name], pair.Value)
		}
		names := fields["Name"]
		if len(names) != 1 {
			return nil, fmt.Errorf("registry entry %d has %d Name fields, want exactly 1", i+1, len(names))
		}
		name, ok := stringLiteral(names[0])
		if !ok || name == "" {
			return nil, fmt.Errorf("registry entry %d has malformed or empty Name", i+1)
		}
		if seen[name] {
			return nil, fmt.Errorf("%s repeats command %q", table, name)
		}
		seen[name] = true
		entries = append(entries, commandRegistryEntry{name: name, fields: fields})
	}
	return entries, nil
}
