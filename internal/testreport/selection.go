package testreport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"syscall"
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

func resolveChangedPackagesWithEnvironment(ctx context.Context, root string, paths, env []string) (changedSelection, error) {
	return resolveChangedPackagesWithLoader(ctx, root, paths, func(ctx context.Context, root string) ([]listedPackage, error) {
		return currentPackagesWithEnvironment(ctx, root, env)
	})
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
	if strings.HasSuffix(clean, ".go") {
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

func currentPackagesWithEnvironment(ctx context.Context, root string, env []string) ([]listedPackage, error) {
	cmd := exec.Command("go", "list", "-buildvcs=false", "-json", "-test", currentPackagePattern)
	cmd.Dir = root
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("go list failed: %w", err)
	}
	completed := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		completed <- cmd.Wait()
		close(done)
	}()
	var err error
	select {
	case err = <-completed:
	case <-ctx.Done():
		cancelGoProcessGroup(cmd, done)
		err = <-completed
		return nil, fmt.Errorf("go list interrupted: %s", goChildGroupCancelled)
	}
	if err != nil {
		return nil, fmt.Errorf("go list failed: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(output.Bytes()))
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
	byDirectory := make(map[string]string, len(packages))
	byEmbed := make(map[string]string)
	byImport := make(map[string]listedPackage, len(packages))
	for _, pkg := range packages {
		directory, err := filepath.Rel(root, pkg.Dir)
		if err != nil || directory == ".." || strings.HasPrefix(directory, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("go list returned a package outside the repository")
		}
		directory = filepath.ToSlash(directory)
		byDirectory[directory] = pkg.ImportPath
		byImport[pkg.ImportPath] = pkg
		for _, name := range append(append(append([]string{}, pkg.EmbedFiles...), pkg.TestEmbedFiles...), pkg.XTestEmbedFiles...) {
			byEmbed[filepath.ToSlash(filepath.Join(directory, name))] = pkg.ImportPath
		}
	}
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
		if input.absent && !strings.HasSuffix(input.path, ".go") {
			if _, ok := byEmbed[input.path]; ok {
				return nil, fmt.Errorf("changed embed path is not in the current package graph")
			}
			continue
		}
		isGo := strings.HasSuffix(input.path, ".go")
		goPackage, inPackage := "", false
		if isGo {
			goPackage, inPackage = byDirectory[filepath.ToSlash(filepath.Dir(input.path))]
		}
		if importPath, ok := byEmbed[input.path]; ok {
			mark(importPath, causeEmbed)
			if inPackage {
				mark(goPackage, causeChanged)
			}
			continue
		}
		if isGo {
			if !inPackage {
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

func isGoMetadata(path string) bool {
	return path == "go.mod" || path == "go.sum" || path == "go.work" || path == "go.work.sum"
}
