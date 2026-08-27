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

func resolveChangedPackages(ctx context.Context, root string, paths []string) ([]string, error) {
	inputs := make([]changedPath, 0, len(paths))
	for _, path := range paths {
		input, err := inspectChangedPath(root, path)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}
	if len(inputs) == 0 {
		return nil, nil
	}
	packages, err := listCurrentPackages(ctx, root)
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
	cmd := exec.CommandContext(ctx, "go", "list", "-buildvcs=false", "-json", "-test", currentPackagePattern)
	cmd.Dir = root
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list failed: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(output))
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

func selectCurrentPackages(root string, packages []listedPackage, inputs []changedPath) ([]string, error) {
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
	selected := make(map[string]bool)
	for _, input := range inputs {
		if isGoMetadata(input.path) {
			for importPath := range byImport {
				selected[importPath] = true
			}
			continue
		}
		if importPath, ok := byEmbed[input.path]; ok {
			selected[importPath] = true
			continue
		}
		if strings.HasSuffix(input.path, ".go") {
			importPath, ok := byDirectory[filepath.ToSlash(filepath.Dir(input.path))]
			if !ok {
				return nil, fmt.Errorf("changed Go path is not in a current package")
			}
			selected[importPath] = true
			continue
		}
		if input.absent {
			return nil, fmt.Errorf("changed path is not in the current package graph")
		}
	}
	reverse := make(map[string][]string)
	for _, pkg := range packages {
		for _, dependency := range append(append(append([]string{}, pkg.Imports...), pkg.TestImports...), pkg.XTestImports...) {
			reverse[dependency] = append(reverse[dependency], pkg.ImportPath)
		}
	}
	for changed := true; changed; {
		changed = false
		for selectedImport := range selected {
			for _, dependent := range reverse[selectedImport] {
				if !selected[dependent] {
					selected[dependent] = true
					changed = true
				}
			}
		}
	}
	result := make([]string, 0, len(selected))
	for importPath := range selected {
		result = append(result, importPath)
	}
	sort.Strings(result)
	return result, nil
}

func isGoMetadata(path string) bool {
	return path == "go.mod" || path == "go.sum" || path == "go.work" || path == "go.work.sum"
}
