package testreport

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/toon"
)

// The table bytes here are authored apart from the face, so a wrong owner, an absolute
// path, or a lost row reds these rows.

const emptyFixturesTable = "fixtures[0]{family,fixture,path}:\n"

// inventoryFace runs Command with a counting run binary factory and a marker-writing `go`
// on PATH, and it fails the test when the face builds a run binary or starts a Go child.
func inventoryFace(t *testing.T, root string, args ...string) (string, int) {
	t.Helper()
	builds := 0
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build:    func(context.Context, string, string) error { builds++; return nil },
		Verify:   func(string, string) error { return nil },
	})
	goDir := t.TempDir()
	marker := filepath.Join(goDir, "go-ran")
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	output, code := Command(root, args)
	if builds != 0 {
		t.Errorf("Command(%q) built %d run binaries, want none", args, builds)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Command(%q) started the canned go: %v", args, err)
	}
	return output, code
}

// plantFixture writes one fixture of family under root's tests/canary, in the shape of
// the canary inventory tests. A non-empty check writes the CHECK marker.
func plantFixture(t *testing.T, root, family, name, check string) {
	t.Helper()
	dir := "tests/canary/" + family + "/" + name
	writeProseCheckFile(t, root, dir+"/EXPECT", "planted diagnostic\n")
	if check != "" {
		writeProseCheckFile(t, root, dir+"/CHECK", check+"\n")
	}
}

func TestFixturesFaceListsOwnedFixtures(t *testing.T) {
	root := t.TempDir()
	// The marked line-routing fixture sorts first by path and last by name. With ten owned
	// rows, an unsorted table matches only by a negligible chance.
	plantFixture(t, root, "line-routing", "zz", "package-core-guard")
	plantFixture(t, root, "line-routing", "other", "")
	want := "fixtures[10]{family,fixture,path}:\n  line-routing,zz,tests/canary/line-routing/zz\n"
	for _, name := range []string{"a", "c", "f", "h", "m", "q", "t", "x", "z"} {
		plantFixture(t, root, "package-core-guard", name, "")
		want += "  package-core-guard," + name + ",tests/canary/package-core-guard/" + name + "\n"
	}
	output, code := inventoryFace(t, root, "--check", "package-core-guard", "--fixtures")
	if code != 0 || output != want {
		t.Fatalf("owned fixtures = %d, %q; want 0 and %q", code, output, want)
	}
}

// reassignedFixtureTree holds one package-core-guard fixture whose CHECK marker names
// default-branch-single-source.
func reassignedFixtureTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "reassigned", "default-branch-single-source")
	return root
}

func TestFixturesFaceHonorsCheckMarker(t *testing.T) {
	output, code := inventoryFace(t, reassignedFixtureTree(t), "--check", "default-branch-single-source", "--fixtures")
	want := "fixtures[1]{family,fixture,path}:\n  package-core-guard,reassigned,tests/canary/package-core-guard/reassigned\n"
	if code != 0 || output != want {
		t.Fatalf("marked fixture = %d, %q; want 0 and %q", code, output, want)
	}
}

func TestFixturesFaceOmitsReassignedFixture(t *testing.T) {
	output, code := inventoryFace(t, reassignedFixtureTree(t), "--check", "package-core-guard", "--fixtures")
	if code != 0 || output != emptyFixturesTable {
		t.Fatalf("family owner of a marked fixture = %d, %q; want 0 and %q", code, output, emptyFixturesTable)
	}
}

func TestFixturesFaceEmptyForSystem(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "a", "")
	output, code := inventoryFace(t, root, "--check", "system", "--fixtures")
	if code != 0 || output != emptyFixturesTable {
		t.Fatalf("system fixtures = %d, %q; want 0 and %q", code, output, emptyFixturesTable)
	}
}

func TestFixturesFaceEmptyCanaryDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "tests", "canary"), 0o755); err != nil {
		t.Fatal(err)
	}
	output, code := inventoryFace(t, root, "--check", "package-core-guard", "--fixtures")
	if code != 0 || output != emptyFixturesTable {
		t.Fatalf("present and empty tests/canary = %d, %q; want 0 and %q", code, output, emptyFixturesTable)
	}
}

func TestFixturesFaceRefusesInvalidInventory(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "a", "no-such-check")
	output, code := inventoryFace(t, root, "--check", "package-core-guard", "--fixtures")
	if code != 1 || !strings.Contains(output, "names unknown check") {
		t.Fatalf("invalid inventory = %d, %q; want 1 and names unknown check", code, output)
	}
}

func TestFixturesFaceRefusesUnprintableName(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "bad\x01name", "")
	output, code := inventoryFace(t, root, "--check", "package-core-guard", "--fixtures")
	if code != 1 || !strings.HasPrefix(output, toon.RenderError(errors.New(""))) {
		t.Fatalf("unprintable fixture name = %d, %q; want 1 and the render error", code, output)
	}
}

func TestFixturesFaceUnknownCheck(t *testing.T) {
	output, code := inventoryFace(t, t.TempDir(), "--check", "not-registered", "--fixtures")
	if code != 2 || !strings.HasPrefix(output, "unknown check: not-registered\n") || !strings.HasSuffix(output, checkInventory()) {
		t.Fatalf("unknown check fixtures = %d, %q; want 2 and the unknown-check refusal", code, output)
	}
}

func TestInventoryGrammarRefusals(t *testing.T) {
	for _, args := range [][]string{
		{"--checks", "--full"},
		{"--checks", "--check", "prose"},
		{"--checks", "internal/testreport"},
		{"--checks", "--changed"},
		{"--checks", "--package", "./internal/testreport"},
		{"--checks", "--base", "HEAD"},
		{"--checks", "--run", "^X$"},
		{"--fixtures"},
		{"--fixtures", "--full"},
		{"--check", "prose", "--fixtures", "--full"},
		{"--check", "system", "--fixtures", "--run", "^X$"},
	} {
		output, code := inventoryFace(t, t.TempDir(), args...)
		if code != 2 || !strings.HasPrefix(output, "usage: ") {
			t.Errorf("Command(%q) = %d, %q; want 2 and usage", args, code, output)
		}
	}
}
