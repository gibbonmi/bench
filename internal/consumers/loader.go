package consumers

import (
	"fmt"
	"go/ast"
	"go/token"
	"maps"
	pathpkg "path"
	"slices"
	"strconv"
	"strings"

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

// SourceUse is one file of a revision's sources that names a declaration, at the line of
// its first naming.
type SourceUse struct {
	File string
	Line int
}

// SourceUses reads one revision's sources without type information, the way the deleted
// table reads its base side. file declares name at package level in the package whose
// import path is importPath. declared is the first line of that declaration, zero when
// file declares no such name. uses is every other source that names it: a file of the
// declaring package through the bare identifier, and any other file through a selector
// on its import of importPath. The reading is lexical, so a local name that shadows the
// declaration still counts as a use, and a source that does not parse names nothing. A
// method spelling has no receiver type to resolve here, so it can declare but names
// nothing.
func SourceUses(sources map[string]string, file, importPath, name string) (declared int, uses []SourceUse) {
	fset := token.NewFileSet()
	owner, err := parseDecls(fset, file, sources[file])
	if err != nil {
		return 0, nil
	}
	for _, site := range topLevelDecls(fset, owner) {
		if site.Name == name && declared == 0 {
			declared = site.Start
		}
	}
	if declared == 0 || strings.Contains(name, ".") {
		return declared, nil
	}
	for _, other := range slices.Sorted(maps.Keys(sources)) {
		parsed, err := parseDecls(fset, other, sources[other])
		if other == file || err != nil {
			continue
		}
		same := pathpkg.Dir(other) == pathpkg.Dir(file) && parsed.Name.Name == owner.Name.Name
		if line := firstUse(fset, parsed, same, importPath, name); line > 0 {
			uses = append(uses, SourceUse{File: other, Line: line})
		}
	}
	return declared, uses
}

// firstUse is the first line of file that names name: a bare identifier in the declaring
// package, or a selector on the file's import of importPath elsewhere. It is zero when
// the file names it nowhere.
func firstUse(fset *token.FileSet, file *ast.File, same bool, importPath, name string) int {
	qualifier := ""
	for _, spec := range file.Imports {
		if name, imported := ImportName(spec); imported == importPath {
			qualifier = name
		}
	}
	line := 0
	var visit func(ast.Node) bool
	visit = func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.SelectorExpr:
			if id, ok := n.X.(*ast.Ident); ok && !same && qualifier != "" && id.Name == qualifier && n.Sel.Name == name && line == 0 {
				line = fset.Position(n.Sel.Pos()).Line
			}
			ast.Inspect(n.X, visit)
			return false
		case *ast.Ident:
			if same && n.Name == name && line == 0 {
				line = fset.Position(n.Pos()).Line
			}
		}
		return line == 0
	}
	ast.Inspect(file, visit)
	return line
}

// ImportName is the one reading of an import's local name: its explicit name, else the
// last element of its path. path is the unquoted import path.
func ImportName(spec *ast.ImportSpec) (name, path string) {
	imported, _ := strconv.Unquote(spec.Path.Value)
	if spec.Name != nil {
		return spec.Name.Name, imported
	}
	return pathpkg.Base(imported), imported
}
