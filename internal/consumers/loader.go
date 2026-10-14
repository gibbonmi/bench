package consumers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

// loadMode is the smallest go/packages mode that yields the core's Package contract:
// the import path, the parsed files, the type-checked package, and the use map.
const loadMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
	packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps |
	packages.NeedImports

// load is the package-internal loader seam: the one site that reaches the go tool. It is
// unexported and injectable inside the package only, because the core must stay a pure
// function of typed packages and the command layer must stay testable without a
// subprocess. It joins no audited port registry; the audit's package set is closed.
var load = loadPackages

// Load type-checks every package of the checkout at root, its tests included, in the
// default build context, for a caller outside this package. Resolve and Rows read the
// answer. The packages hold the checkout's bytes, so only a clean checkout answers for
// its committed HEAD, and a dirty one refuses.
func Load(root string) ([]*Package, error) {
	if (citation{root: root}).state() != "clean" {
		return nil, fmt.Errorf("checkout %s is dirty; typed packages answer only for a clean committed HEAD", root)
	}
	return loadPackages(root, "./...")
}

// loadPackages type-checks patterns rooted at dir and returns them in the core's shape.
// It fails closed on the first load error, because an enumeration over a partially typed
// tree under-reports and reads as an absence proof.
func loadPackages(dir string, patterns ...string) ([]*Package, error) {
	cfg := &packages.Config{Mode: loadMode, Dir: dir, Tests: true}
	loaded, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}
	var out []*Package
	for _, p := range loaded {
		if len(p.Errors) > 0 {
			e := p.Errors[0]
			return nil, fmt.Errorf("%s: %s", e.Pos, e.Msg)
		}
		if p.TypesInfo == nil || p.Types == nil {
			continue
		}
		// go/packages synthesizes a test-binary main file in the build cache and delivers
		// it inside the [pkg.test] variant. It references every Test function, so a query
		// of one would otherwise enumerate a file outside the checkout. This is the one
		// place an out-of-root file is first seen, so it is dropped here rather than
		// filtered again in each renderer.
		files := make([]*ast.File, 0, len(p.Syntax))
		for _, f := range p.Syntax {
			if insideRoot(dir, p.Fset.Position(f.Pos()).Filename) {
				files = append(files, f)
			}
		}
		out = append(out, &Package{
			PkgPath: p.PkgPath,
			Fset:    p.Fset,
			Files:   files,
			Types:   p.Types,
			Info:    p.TypesInfo,
		})
	}
	return out, nil
}
