package testreport

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/subprocess"
)

const currentPackagePattern = "./..."

type listedPackage struct {
	Dir             string
	ImportPath      string
	ForTest         string
	Match           []string
	Imports         []string
	TestImports     []string
	XTestImports    []string
	EmbedFiles      []string
	TestEmbedFiles  []string
	XTestEmbedFiles []string
}

// packageIndex is the one path-to-package map of the current package graph. The
// `--changed` selection and the ticket floor both read it.
type packageIndex struct {
	// byDirectory maps a slash-spelled directory, relative to the root, to its import path.
	byDirectory map[string]string
	// byEmbed maps a slash-spelled embed file of any class, relative to the root, to the
	// import path of its owner.
	byEmbed map[string]string
	// byImport maps an import path to its package.
	byImport map[string]listedPackage
	// system is the import path of the system-suite package, or empty when the index does
	// not hold it.
	system string
}

func newPackageIndex(root string, packages []listedPackage) (packageIndex, error) {
	index := packageIndex{
		byDirectory: make(map[string]string, len(packages)),
		byEmbed:     make(map[string]string),
		byImport:    make(map[string]listedPackage, len(packages)),
		system:      systemSuiteImport(packages),
	}
	for _, pkg := range packages {
		directory, err := filepath.Rel(root, pkg.Dir)
		if err != nil || directory == ".." || strings.HasPrefix(directory, ".."+string(filepath.Separator)) {
			return packageIndex{}, fmt.Errorf("go list returned a package outside the repository")
		}
		directory = filepath.ToSlash(directory)
		index.byDirectory[directory] = pkg.ImportPath
		index.byImport[pkg.ImportPath] = pkg
		for _, name := range append(append(append([]string{}, pkg.EmbedFiles...), pkg.TestEmbedFiles...), pkg.XTestEmbedFiles...) {
			index.byEmbed[filepath.ToSlash(filepath.Join(directory, name))] = pkg.ImportPath
		}
	}
	return index, nil
}

// owners returns the direct packages of one slash-spelled path: the package whose
// directory holds the path as a Go file, and the package that embeds the path. An empty
// import path means that no package has that role.
func (index packageIndex) owners(file string) (goPackage, embedOwner string) {
	if isGoFile(file) {
		goPackage = index.byDirectory[path.Dir(file)]
	}
	return goPackage, index.byEmbed[file]
}

var listCurrentPackages = currentPackages

// changedSelection maps each selected import path to the one cause that selected it.
type changedSelection map[string]string

const (
	causeGoMetadata = "go-metadata"
	causeChanged    = "changed"
	causeEmbed      = "embed"
	// causeImports prefixes the import path of the dependency that selected the package.
	causeImports = "imports "
)

// causePrecedence lists the direct causes from the strongest to the weakest. A package
// that two inputs select keeps the stronger cause.
var causePrecedence = []string{causeGoMetadata, causeChanged, causeEmbed}

// packages returns the selected import paths in sorted order.
func (s changedSelection) packages() []string {
	result := make([]string, 0, len(s))
	for importPath := range s {
		result = append(result, importPath)
	}
	sort.Strings(result)
	return result
}

func resolveChangedPackages(ctx context.Context, root string, paths []string) (changedSelection, error) {
	return resolveChangedPackagesWithLoader(ctx, root, paths, listCurrentPackages)
}

// resolveChangedRun resolves the `--changed` selection of paths, and the `go test` operands
// and added environment that run it. A selected system-suite package runs as `bench test
// --check system` runs it: under the gate's system-suite operands and environment, because
// a package operand without the tag builds no file of the suite.
func resolveChangedRun(ctx context.Context, root string, paths, env []string) (selection changedSelection, operands, suiteEnv []string, err error) {
	var system string
	selection, err = resolveChangedPackagesWithLoader(ctx, root, paths, func(ctx context.Context, root string) ([]listedPackage, error) {
		packages, err := currentPackagesWithEnvironment(ctx, root, env)
		system = systemSuiteImport(packages)
		return packages, err
	})
	if _, selected := selection[system]; err != nil || !selected {
		return selection, selection.packages(), nil, err
	}
	operands, suiteEnv = gate.SystemSuite(root)
	others := slices.DeleteFunc(selection.packages(), func(importPath string) bool { return importPath == system })
	return selection, append(operands, others...), suiteEnv, nil
}

func resolveChangedPackagesWithLoader(ctx context.Context, root string, paths []string, load func(context.Context, string) ([]listedPackage, error)) (changedSelection, error) {
	inputs := make([]changedPath, 0, len(paths))
	for _, path := range paths {
		input, err := inspectChangedPath(root, path)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}
	if len(inputs) == 0 {
		// A non-nil empty selection keeps the `selected_by` column in the empty report.
		return changedSelection{}, nil
	}
	packages, err := load(ctx, root)
	if err != nil {
		return nil, err
	}
	return selectCurrentPackages(root, packages, inputs)
}

type changedPath struct {
	path   string
	absent bool
}

func inspectChangedPath(root, path string) (changedPath, error) {
	if strings.IndexFunc(path, func(r rune) bool { return r < 0x20 || r == 0x7f }) >= 0 {
		return changedPath{}, fmt.Errorf("changed path is unsafe")
	}
	clean := filepath.Clean(filepath.FromSlash(path))
	if clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return changedPath{}, fmt.Errorf("changed path is outside the repository")
	}
	info, err := os.Lstat(filepath.Join(root, clean))
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			return changedPath{}, fmt.Errorf("changed path is unsafe")
		}
		return changedPath{path: filepath.ToSlash(clean)}, nil
	}
	if !os.IsNotExist(err) {
		return changedPath{}, fmt.Errorf("changed path is unreadable: %w", err)
	}
	if isGoFile(clean) {
		parent, parentErr := os.Stat(filepath.Dir(filepath.Join(root, clean)))
		if parentErr != nil || !parent.IsDir() {
			return changedPath{}, fmt.Errorf("changed Go path is not in a current package")
		}
	}
	return changedPath{path: filepath.ToSlash(clean), absent: true}, nil
}

func currentPackages(ctx context.Context, root string) ([]listedPackage, error) {
	return currentPackagesWithEnvironment(ctx, root, os.Environ())
}

// systemSuitePackage is the package operand among the gate's system-suite operands. The
// other operands are flags.
func systemSuitePackage() string {
	operands, _ := gate.SystemSuite("")
	return operands[slices.IndexFunc(operands, func(operand string) bool { return !strings.HasPrefix(operand, "-") })]
}

// systemSuiteOperands are the gate's system-suite operands where the suite runs and root
// holds its package, and nil otherwise. A linked repository keeps its untagged index, and
// a module without the package keeps a listing that `go list` can answer.
func systemSuiteOperands(root string) []string {
	operands, _ := gate.SystemSuite(root)
	info, err := os.Stat(filepath.Join(root, filepath.FromSlash(systemSuitePackage())))
	if !gate.SystemSuiteRuns(root, testBenchSource(root)) || err != nil || !info.IsDir() {
		return nil
	}
	return operands
}

// systemSuiteImport is the import path of the listed package that `go list` matched by the
// system-suite package operand, or empty when the listing holds no such package.
func systemSuiteImport(packages []listedPackage) string {
	system := systemSuitePackage()
	for _, pkg := range packages {
		if slices.Contains(pkg.Match, system) {
			return pkg.ImportPath
		}
	}
	return ""
}

func currentPackagesWithEnvironment(ctx context.Context, root string, env []string) ([]listedPackage, error) {
	// The tag among the system-suite operands also makes `./...` match the suite, so the
	// match filter below keeps it.
	args := append(append([]string{"list", "-buildvcs=false", "-json", "-test"}, systemSuiteOperands(root)...), currentPackagePattern)
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	cmd.Env = env
	child := runGoChild(ctx, cmd, nil)
	if child.Cancellation.Cause != subprocess.CauseNone {
		return nil, fmt.Errorf("go list interrupted: %s", goChildGroupCancelled)
	}
	if err := cmp.Or(child.Process.StartErr, child.Process.WaitErr, child.Uncertified()); err != nil {
		return nil, fmt.Errorf("go list failed: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(child.Streams.Stdout.Data))
	var packages []listedPackage
	for {
		var pkg listedPackage
		err := decoder.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("go list output malformed: %w", err)
		}
		if pkg.ForTest != "" || !slices.Contains(pkg.Match, currentPackagePattern) {
			continue
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

func selectCurrentPackages(root string, packages []listedPackage, inputs []changedPath) (changedSelection, error) {
	index, err := newPackageIndex(root, packages)
	if err != nil {
		return nil, err
	}
	byImport := index.byImport
	selected := changedSelection{}
	mark := func(importPath, cause string) {
		if current, ok := selected[importPath]; !ok || slices.Index(causePrecedence, cause) < slices.Index(causePrecedence, current) {
			selected[importPath] = cause
		}
	}
	for _, input := range inputs {
		if isGoMetadata(input.path) {
			for importPath := range byImport {
				mark(importPath, causeGoMetadata)
			}
			continue
		}
		goPackage, embedOwner := index.owners(input.path)
		if embedOwner != "" {
			if input.absent && !isGoFile(input.path) {
				return nil, fmt.Errorf("changed embed path is not in the current package graph")
			}
			mark(embedOwner, causeEmbed)
			if goPackage != "" {
				mark(goPackage, causeChanged)
			}
			continue
		}
		if isGoFile(input.path) {
			if goPackage == "" {
				return nil, fmt.Errorf("changed Go path is not in a current package")
			}
			mark(goPackage, causeChanged)
		}
	}
	reverse := make(map[string][]string)
	for _, pkg := range packages {
		for _, dependency := range directDependencies(pkg) {
			reverse[dependency] = append(reverse[dependency], pkg.ImportPath)
		}
	}
	for changed := true; changed; {
		changed = false
		for selectedImport := range selected {
			for _, dependent := range reverse[selectedImport] {
				if _, ok := selected[dependent]; !ok {
					selected[dependent] = ""
					changed = true
				}
			}
		}
	}
	// A package that the closure added names its first selected direct dependency, so the
	// cause does not depend on the order in which the closure reached the package. The
	// package's own path, which its external tests can import, is not a cause.
	for importPath, cause := range selected {
		if cause != "" {
			continue
		}
		dependencies := directDependencies(byImport[importPath])
		sort.Strings(dependencies)
		for _, dependency := range dependencies {
			if dependency == importPath {
				continue
			}
			if _, ok := selected[dependency]; ok {
				selected[importPath] = causeImports + dependency
				break
			}
		}
	}
	return selected, nil
}

// directDependencies returns the imports, the test imports, and the external test imports
// of pkg in one new slice.
func directDependencies(pkg listedPackage) []string {
	return append(append(append([]string{}, pkg.Imports...), pkg.TestImports...), pkg.XTestImports...)
}

func isGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

func isGoMetadata(path string) bool {
	return path == "go.mod" || path == "go.sum" || path == "go.work" || path == "go.work.sum"
}
