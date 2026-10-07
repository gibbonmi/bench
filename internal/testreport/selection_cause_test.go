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
// selected packages names the first in sorted import-path order. Each case holds the
// first one last in a different dependency list, so an order that the lists or a map
// give fails here.
func TestCauseImportsNamesFirstSortedDependency(t *testing.T) {
	names := []string{"f", "d", "b", "e", "c", "a"}
	packages := []listedPackage{}
	paths := []string{}
	for _, name := range names {
		packages = append(packages, causePackage(name))
		paths = append(paths, name+"/"+name+".go")
	}
	for _, lists := range [][3][]string{
		{{"example/f", "example/d", "example/b"}, {"example/e", "example/c"}, {"example/a"}},
		{{"example/f", "example/d", "example/b"}, {"example/e", "example/a"}, {"example/c"}},
		{{"example/f", "example/d", "example/a"}, {"example/e", "example/c"}, {"example/b"}},
	} {
		dependent := causePackage("z", lists[0]...)
		dependent.TestImports = lists[1]
		dependent.XTestImports = lists[2]
		for range 20 {
			if got := selectCauses(t, append(packages, dependent), paths...)["example/z"]; got != "imports example/a" {
				t.Fatalf("cause of example/z with lists %v = %q, want %q", lists, got, "imports example/a")
			}
		}
	}
}

// TestCauseImportsSkipsUnselectedDependency grades that a package names a selected
// dependency and not an unselected one that sorts first.
func TestCauseImportsSkipsUnselectedDependency(t *testing.T) {
	got := selectCauses(t, []listedPackage{causePackage("a"), causePackage("b", "bytes", "example/a")}, "a/a.go")
	assertCauses(t, got, changedSelection{"example/a": "changed", "example/b": "imports example/a"})
}

// TestCauseImportsThroughTestImport grades that a package that only its tests connect to
// the change names that test import.
func TestCauseImportsThroughTestImport(t *testing.T) {
	dependent := causePackage("b")
	dependent.TestImports = []string{"example/a"}
	got := selectCauses(t, []listedPackage{causePackage("a"), dependent}, "a/a.go")
	assertCauses(t, got, changedSelection{"example/a": "changed", "example/b": "imports example/a"})
}

// TestCauseImportsSkipsOwnPath grades that a package whose external tests import the
// package itself names its real dependency and not its own path.
func TestCauseImportsSkipsOwnPath(t *testing.T) {
	dependent := causePackage("m", "example/n")
	dependent.XTestImports = []string{"example/m"}
	got := selectCauses(t, []listedPackage{causePackage("n"), dependent}, "n/n.go")
	assertCauses(t, got, changedSelection{"example/n": "changed", "example/m": "imports example/n"})
}

// TestCauseChangedBeatsEmbed grades that a package with a changed Go file and a changed
// embed file names the change, in each path order, and also when the changed Go file is
// one of the package's own embed files.
func TestCauseChangedBeatsEmbed(t *testing.T) {
	embed := causePackage("e")
	embed.EmbedFiles = []string{"data.txt", "e.go"}
	want := changedSelection{"example/e": "changed"}
	assertCauses(t, selectCauses(t, []listedPackage{embed}, "e/e.go"), want)
	assertCauses(t, selectCauses(t, []listedPackage{embed}, "e/data.txt", "e/f.go"), want)
	assertCauses(t, selectCauses(t, []listedPackage{embed}, "e/f.go", "e/data.txt"), want)
}

// TestCauseEmbeddedGoFileSelectsItsPackage grades that a changed Go file that another
// package embeds selects the embedder by the embed and its own package by the change, and
// that the closure then reaches the dependents of its own package.
func TestCauseEmbeddedGoFileSelectsItsPackage(t *testing.T) {
	embedder := causePackage("p")
	embedder.EmbedFiles = []string{"sub/s.go"}
	got := selectCauses(t, []listedPackage{embedder, causePackage("p/sub"), causePackage("q", "example/p/sub")}, "p/sub/s.go")
	assertCauses(t, got, changedSelection{"example/p": "embed", "example/p/sub": "changed", "example/q": "imports example/p/sub"})
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
