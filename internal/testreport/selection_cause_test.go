package testreport

import (
	"path/filepath"
	"reflect"
	"testing"
)

const causeRoot = "/repo"

// causePackage is one canned `go list` package in directory name under causeRoot, with
// the import path `example/<name>` and the given production imports.
func causePackage(name string, imports ...string) listedPackage {
	return listedPackage{Dir: filepath.Join(causeRoot, name), ImportPath: "example/" + name, Imports: imports}
}

// selectCauses runs the selector over the canned packages and the changed paths.
func selectCauses(t *testing.T, packages []listedPackage, paths ...string) changedSelection {
	t.Helper()
	inputs := make([]changedPath, 0, len(paths))
	for _, path := range paths {
		inputs = append(inputs, changedPath{path: path})
	}
	got, err := selectCurrentPackages(causeRoot, packages, inputs)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func assertCauses(t *testing.T, got, want changedSelection) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("causes = %v, want %v", got, want)
	}
}

// TestCauseChangedBeatsImports grades that a package that changed and that imports a
// changed package names its own change.
func TestCauseChangedBeatsImports(t *testing.T) {
	got := selectCauses(t, []listedPackage{causePackage("a"), causePackage("b", "example/a")}, "a/a.go", "b/b.go")
	assertCauses(t, got, changedSelection{"example/a": "changed", "example/b": "changed"})
}

// TestCauseEmbedOnly grades that a package with only a changed embed file names the embed.
func TestCauseEmbedOnly(t *testing.T) {
	embed := causePackage("e")
	embed.EmbedFiles = []string{"data.txt"}
	assertCauses(t, selectCauses(t, []listedPackage{embed}, "e/data.txt"), changedSelection{"example/e": "embed"})
}

// TestCauseImportsNamesFirstSortedDependency grades that a package that imports several
// selected packages names the first in sorted import-path order. Each of the dependency
// lists holds the first one late, so an order that the lists or a map give fails here.
func TestCauseImportsNamesFirstSortedDependency(t *testing.T) {
	names := []string{"f", "d", "b", "e", "c", "a"}
	packages := []listedPackage{}
	paths := []string{}
	for _, name := range names {
		packages = append(packages, causePackage(name))
		paths = append(paths, name+"/"+name+".go")
	}
	dependent := causePackage("z", "example/f", "example/d", "example/b")
	dependent.TestImports = []string{"example/e", "example/c"}
	dependent.XTestImports = []string{"example/a"}
	for range 20 {
		if got := selectCauses(t, append(packages, dependent), paths...)["example/z"]; got != "imports example/a" {
			t.Fatalf("cause of example/z = %q, want %q", got, "imports example/a")
		}
	}
}

// TestCauseImportsNamesDirectDependency grades that a package two import steps from the
// change names its direct dependency, not the changed package.
func TestCauseImportsNamesDirectDependency(t *testing.T) {
	got := selectCauses(t, []listedPackage{causePackage("a"), causePackage("c", "example/a"), causePackage("d", "example/c")}, "a/a.go")
	assertCauses(t, got, changedSelection{"example/a": "changed", "example/c": "imports example/a", "example/d": "imports example/c"})
}

// TestCauseGoMetadataWins grades that a changed `go.mod` names `go-metadata` on each row,
// the changed package included, whatever the order of the changed paths.
func TestCauseGoMetadataWins(t *testing.T) {
	packages := []listedPackage{causePackage("a"), causePackage("b", "example/a")}
	want := changedSelection{"example/a": "go-metadata", "example/b": "go-metadata"}
	assertCauses(t, selectCauses(t, packages, "go.mod", "a/a.go"), want)
	assertCauses(t, selectCauses(t, packages, "a/a.go", "go.mod"), want)
}
