package backgroundstate_test

import (
	"path"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// TestImportClosureHoldsNoOtherGatePackage pins the leaf boundary. The import closure of
// backgroundstate holds no internal/gate package other than itself, so a package that
// internal/gate imports can read and stop a grade without an import cycle.
func TestImportClosureHoldsNoOtherGatePackage(t *testing.T) {
	roots, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps}, ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 1 {
		t.Fatalf("the load of this package answered %d packages, want 1", len(roots))
	}
	if errs := roots[0].Errors; len(errs) != 0 {
		t.Fatalf("the load of this package: %v", errs)
	}
	self := roots[0].PkgPath
	gate := path.Dir(self)
	packages.Visit(roots, nil, func(dep *packages.Package) {
		if dep.PkgPath != self && (dep.PkgPath == gate || strings.HasPrefix(dep.PkgPath, gate+"/")) {
			t.Errorf("the import closure of %s holds %s", self, dep.PkgPath)
		}
	})
}
