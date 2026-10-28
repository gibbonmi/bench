package backgroundstate_test

import (
	"os/exec"
	"path"
	"strings"
	"testing"
)

// TestImportClosureHoldsNoOtherGatePackage pins the leaf boundary. The import closure of
// backgroundstate holds no internal/gate package other than itself, so a package that
// internal/gate imports can read and stop a grade without an import cycle.
func TestImportClosureHoldsNoOtherGatePackage(t *testing.T) {
	out, err := exec.Command("go", "list", "-buildvcs=false", "-deps", "-f", "{{.ImportPath}}", ".").Output()
	if err != nil {
		t.Fatalf("go list -deps: %v", err)
	}
	// go list -deps prints the closure in dependency order, so the package itself comes last.
	deps := strings.Fields(string(out))
	if len(deps) == 0 {
		t.Fatal("go list -deps printed no package")
	}
	self := deps[len(deps)-1]
	gate := path.Dir(self)
	for _, dep := range deps {
		if dep != self && (dep == gate || strings.HasPrefix(dep, gate+"/")) {
			t.Errorf("the import closure of %s holds %s", self, dep)
		}
	}
}
