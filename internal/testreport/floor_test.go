package testreport

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

// floorFixture returns the changed-selection module and its canned graph. Each extra
// directory becomes one more canned package under the module root, so a test can hold
// nested or sibling packages that the graph does not hold.
func floorFixture(t *testing.T, extra ...string) (string, []listedPackage) {
	t.Helper()
	root := changedSelectionModule(t)
	packages := changedSelectionGraph(root)
	for _, directory := range extra {
		if err := os.MkdirAll(filepath.Join(root, directory), 0o755); err != nil {
			t.Fatal(err)
		}
		packages = append(packages, listedPackage{Dir: filepath.Join(root, directory), ImportPath: "changedfixture/" + directory})
	}
	return root, packages
}

// cannedLoader answers the canned packages in place of `go list`.
func cannedLoader(packages []listedPackage) func(context.Context, string) ([]listedPackage, error) {
	return func(context.Context, string) ([]listedPackage, error) {
		return packages, nil
	}
}

// floorOf derives the floor of the entries over the canned packages.
func floorOf(t *testing.T, root string, packages []listedPackage, entries ...string) Floor {
	t.Helper()
	got, err := writesFloorWithLoader(context.Background(), root, entries, cannedLoader(packages))
	if err != nil {
		t.Fatalf("floor of %q: %v", entries, err)
	}
	return got
}

func assertFloor(t *testing.T, got Floor, packages []string, embed bool) {
	t.Helper()
	if !slices.Equal(got.Packages, packages) || got.Embed != embed {
		t.Fatalf("floor = %+v, want packages %q and embed %t", got, packages, embed)
	}
}

// TestWritesFloorAddsTheFilePackage grades that a present Go file adds the package of
// its directory, and that a Go file at the root adds the root package.
func TestWritesFloorAddsTheFilePackage(t *testing.T) {
	root, packages := floorFixture(t)
	assertFloor(t, floorOf(t, root, packages, "direct/direct.go"), []string{"./direct"}, false)
	if err := os.WriteFile(filepath.Join(root, "root.go"), []byte("package changedfixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rootPackage := listedPackage{Dir: root, ImportPath: "changedfixture"}
	assertFloor(t, floorOf(t, root, append(packages, rootPackage), "root.go"), []string{"./."}, false)
}

// TestWritesFloorAddsTheNewFilePackage grades that a `(new)` Go file in a current
// package directory adds that package, although the file does not exist yet.
func TestWritesFloorAddsTheNewFilePackage(t *testing.T) {
	root, packages := floorFixture(t)
	assertFloor(t, floorOf(t, root, packages, "production/later.go (new)"), []string{"./production"}, false)
}

// TestWritesFloorCoversADirectoryEntry grades that a directory entry adds the package
// at the directory and the package below it, and no package outside it.
func TestWritesFloorCoversADirectoryEntry(t *testing.T) {
	root, packages := floorFixture(t, "nest", "nest/inner")
	assertFloor(t, floorOf(t, root, packages, "nest/"), []string{"./nest", "./nest/inner"}, false)
}

// TestWritesFloorKeepsTheSegmentBoundary grades that a directory entry covers a package
// at a `/` boundary only, so a sibling that shares its leading bytes stays out.
func TestWritesFloorKeepsTheSegmentBoundary(t *testing.T) {
	root, packages := floorFixture(t, "internal/git", "internal/git2")
	assertFloor(t, floorOf(t, root, packages, "internal/git"), []string{"./internal/git"}, false)
}

// TestWritesFloorAddsTheEmbedOwner grades that an embed file adds its owner and sets the
// embed flag, and that a Go file sets the flag only in a package that declares an embed
// pattern.
func TestWritesFloorAddsTheEmbedOwner(t *testing.T) {
	root, packages := floorFixture(t)
	assertFloor(t, floorOf(t, root, packages, "embedprod/input.txt"), []string{"./embedprod"}, true)
	assertFloor(t, floorOf(t, root, packages, "embedprod/embed.go"), []string{"./embedprod"}, false)
	// Each canned package now declares one pattern for each embed file it holds.
	for index := range packages {
		packages[index].EmbedPatterns = slices.Clone(packages[index].EmbedFiles)
	}
	assertFloor(t, floorOf(t, root, packages, "embedprod/embed.go"), []string{"./embedprod"}, true)
	// `go list` itself reports the embed patterns of a module whose root package embeds.
	got, err := WritesFloor(context.Background(), currentPackageGraphModule(t), []string{"embed.go"})
	if err != nil {
		t.Fatal(err)
	}
	assertFloor(t, got, []string{"./."}, true)
}

// TestWritesFloorAddsEveryPackageForModuleMetadata grades that a `go.mod` entry adds the
// pattern that selects every package of the module.
func TestWritesFloorAddsEveryPackageForModuleMetadata(t *testing.T) {
	root, packages := floorFixture(t)
	assertFloor(t, floorOf(t, root, packages, "go.mod"), []string{"./..."}, false)
}

// TestWritesFloorIgnoresAGuidancePath grades that a Markdown file that no package embeds
// adds no package, also inside a package directory.
func TestWritesFloorIgnoresAGuidancePath(t *testing.T) {
	root, packages := floorFixture(t)
	assertFloor(t, floorOf(t, root, packages, "direct/README.md", "guidance.md"), nil, false)
}

// TestWritesFloorSkipsAnUnbornPackage grades that a `(new)` Go file under no current
// package directory adds no package and returns no error.
func TestWritesFloorSkipsAnUnbornPackage(t *testing.T) {
	root, packages := floorFixture(t)
	assertFloor(t, floorOf(t, root, packages, "unborn/unborn.go (new)"), nil, false)
}

// TestWritesFloorSkipsAnUnbornEmbedFile grades that an absent `(new)` file in an embed
// owner's directory adds no package and returns no error, because `go list` names only
// the embed files that exist.
func TestWritesFloorSkipsAnUnbornEmbedFile(t *testing.T) {
	root, packages := floorFixture(t)
	assertFloor(t, floorOf(t, root, packages, "embedprod/later.txt (new)"), nil, false)
}

// TestWritesFloorSkipsANonGoRoot grades that a root without `go.mod` gives an empty floor
// and starts no `go list`.
func TestWritesFloorSkipsANonGoRoot(t *testing.T) {
	calls := 0
	got, err := writesFloorWithLoader(context.Background(), t.TempDir(), []string{"go.mod", "main.go"}, func(context.Context, string) ([]listedPackage, error) {
		calls++
		return nil, nil
	})
	if err != nil || calls != 0 {
		t.Fatalf("floor error = %v and loader calls = %d, want no error and no call", err, calls)
	}
	assertFloor(t, got, nil, false)
}

// TestWritesFloorAgreesWithChangedSelection grades that the floor and the `--changed`
// selection give the same direct packages for each canned path, across the Go file and
// the three embed classes.
func TestWritesFloorAgreesWithChangedSelection(t *testing.T) {
	root, packages := floorFixture(t)
	patterns := map[string]string{}
	for _, pkg := range packages {
		directory, err := filepath.Rel(root, pkg.Dir)
		if err != nil {
			t.Fatal(err)
		}
		patterns[pkg.ImportPath] = packagePattern(root, filepath.ToSlash(directory))
	}
	for _, path := range []string{"direct/direct.go", "embedprod/input.txt", "embedtest/input.txt", "embedxtest/input.txt", "space [glob]/embed input[*].txt"} {
		selection, err := resolveChangedPackagesWithLoader(context.Background(), root, []string{path}, cannedLoader(packages))
		if err != nil {
			t.Fatal(err)
		}
		var direct []string
		for importPath, cause := range selection {
			if slices.Contains(causePrecedence, cause) {
				direct = append(direct, patterns[importPath])
			}
		}
		slices.Sort(direct)
		if got := floorOf(t, root, packages, path).Packages; len(direct) == 0 || !reflect.DeepEqual(got, direct) {
			t.Errorf("floor of %q = %q, want the selection's direct packages %q", path, got, direct)
		}
	}
}
