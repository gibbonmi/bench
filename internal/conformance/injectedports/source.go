package injectedports

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// productionRoots are the graded-root directories whose Go source is production source.
// Fixture source outside them is not part of the census.
var productionRoots = []string{"cmd", "internal"}

// sourceTree is the parsed Go source of the production roots, one entry per package
// directory. The variable and setter derivation reads it whole, so a seam in a package
// that no registry row names is still found.
type sourceTree struct {
	fset     *token.FileSet
	packages map[string]*sourcePackage
}

// sourcePackage is one directory's Go source. Production files are the non-test files.
// Test files include both the internal and the external test package.
type sourcePackage struct {
	rel         string
	prod, tests []*sourceFile
}

// sourceFile is one parsed file with its module-internal import bindings, keyed by the
// local name that the file uses for each imported package.
type sourceFile struct {
	rel     string
	clause  string
	file    *ast.File
	imports map[string]string
}

// loadSourceTree parses every Go file under the production roots of root. A root that
// the graded tree does not hold is skipped, so a partial fixture root stays benign. A
// directory or file the walk cannot read or parse is a diagnostic, never silence.
func loadSourceTree(root string) (*sourceTree, []string) {
	tree := &sourceTree{fset: token.NewFileSet(), packages: map[string]*sourcePackage{}}
	module := modulePath(root)
	var diags []string
	for _, top := range productionRoots {
		base := filepath.Join(root, top)
		if absentPackage(base) {
			continue
		}
		err := filepath.WalkDir(base, func(path string, entry fs.DirEntry, err error) error {
			rel := relSlash(root, path)
			if err != nil {
				diags = append(diags, fmt.Sprintf("injected port derivation cannot read %s: %v", rel, err))
				return nil
			}
			if entry.IsDir() {
				if path != base && skippedDir(entry.Name()) {
					return filepath.SkipDir
				}
				return nil
			}
			if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".go") {
				return nil
			}
			file, err := parser.ParseFile(tree.fset, path, nil, 0)
			if err != nil {
				diags = append(diags, fmt.Sprintf("injected port derivation cannot parse %s: %v", rel, err))
				return nil
			}
			dir := relSlash(root, filepath.Dir(path))
			pkg := tree.packages[dir]
			if pkg == nil {
				pkg = &sourcePackage{rel: dir}
				tree.packages[dir] = pkg
			}
			source := &sourceFile{rel: rel, clause: file.Name.Name, file: file, imports: moduleImports(file, module)}
			if strings.HasSuffix(entry.Name(), "_test.go") {
				pkg.tests = append(pkg.tests, source)
			} else {
				pkg.prod = append(pkg.prod, source)
			}
			return nil
		})
		if err != nil {
			diags = append(diags, fmt.Sprintf("injected port derivation cannot read %s: %v", top, err))
		}
	}
	return tree, diags
}

// skippedDir reports the directories that the go tool itself ignores.
func skippedDir(name string) bool {
	return name == "testdata" || name == "vendor" || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

func relSlash(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

// modulePath reads the module line of the graded root's go.mod. A root without one has
// no module-internal imports, so only same-package source resolves.
func modulePath(root string) string {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if fields := strings.Fields(line); len(fields) == 2 && fields[0] == "module" {
			return strings.Trim(fields[1], `"`)
		}
	}
	return ""
}

// moduleImports maps each local binding of a module-internal import to the imported
// package's directory relative to the root, and each standard-library binding to its
// import path. An alias therefore resolves like the default name, so testing.Testing or
// os.Getenv is found under any spelling. A blank or dot import binds no name.
func moduleImports(file *ast.File, module string) map[string]string {
	imports := map[string]string{}
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		dir := path
		switch {
		case module != "" && strings.HasPrefix(path, module+"/"):
			dir = strings.TrimPrefix(path, module+"/")
		case strings.Contains(strings.Split(path, "/")[0], "."):
			continue
		}
		name := dir[strings.LastIndex(dir, "/")+1:]
		if spec.Name != nil {
			name = spec.Name.Name
		}
		if name == "_" || name == "." {
			continue
		}
		imports[name] = dir
	}
	return imports
}

// sortedKeys returns the keys of a set in order, so a derivation reports stably.
func sortedKeys[V any](set map[string]V) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
