package conformance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// The guard owner checks put their own prefix before these diagnostic pieces.
const (
	guardOwnerUnreadable = " is missing or unreadable"
	guardOwnerIrregular  = " is not a regular file"
	guardOwnerUnparsed   = "cannot parse "
)

// guardOwnerSources lists the production Go files beneath dir of root as slash paths relative to root.
// It skips test files and testdata directories. It names a Go file that is not a regular file before
// any open, because an open of a FIFO can block the gate. An entry that the walk cannot read gets a
// diagnostic, because an unlisted subdirectory would hide its production files from the check.
func guardOwnerSources(root, dir, prefix string) (production, diags []string) {
	_ = filepath.WalkDir(filepath.Join(root, filepath.FromSlash(dir)), func(path string, entry fs.DirEntry, err error) error {
		relPath, relErr := filepath.Rel(root, path)
		rel := filepath.ToSlash(relPath)
		switch {
		case err != nil || relErr != nil:
			diags = append(diags, prefix+rel+guardOwnerUnreadable)
		case entry.IsDir() && entry.Name() == "testdata":
			return fs.SkipDir
		case entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go"):
		case !entry.Type().IsRegular():
			diags = append(diags, prefix+rel+guardOwnerIrregular)
		default:
			production = append(production, rel)
		}
		return nil
	})
	return production, diags
}

// guardOwnerParse parses each file of rels beneath root, in order, and collects the diagnostics that
// inspect gives for it. A file that does not parse gets a diagnostic instead.
func guardOwnerParse(root, prefix string, rels []string, inspect func(fset *token.FileSet, rel string, file *ast.File) []string) []string {
	var diags []string
	fset := token.NewFileSet()
	for _, rel := range rels {
		file, err := parser.ParseFile(fset, filepath.Join(root, filepath.FromSlash(rel)), nil, 0)
		if err != nil {
			diags = append(diags, prefix+guardOwnerUnparsed+rel)
			continue
		}
		diags = append(diags, inspect(fset, rel, file)...)
	}
	return diags
}

// guardOwnerDiags runs the registered package-core-guard binding on root and keeps only the
// diagnostics that start with prefix, so each test grades one subcheck through its registered owner.
func guardOwnerDiags(t *testing.T, root, prefix string) []string {
	t.Helper()
	binding, bound := conformanceChecks["package-core-guard"]
	if !bound {
		t.Fatal("package-core-guard conformance owner is not bound")
	}
	var owned []string
	for _, diag := range binding.run(root, root, registry.Dev) {
		if strings.HasPrefix(diag, prefix) {
			owned = append(owned, diag)
		}
	}
	return owned
}
