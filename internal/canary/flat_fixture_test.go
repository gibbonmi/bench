package canary

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// flatFixtures is the harness's own statement of which fixtures run the full inner gate
// from directly under tests/canary/: the one whose EXPECT the gate emits before any phase
// runs, which no phase and so no package can own. It is written out rather than derived
// because flattening is the cheapest wrong migration — a full inner run keeps bite and
// vacuity green while removing the entire scoping win, so the escape has to be named to
// be caught.
var flatFixtures = []string{
	"phase-manifest-defect-admitted",
}

// TestFlatFixtureSetIsExactlyTheUnownedOnes grades the real harness tree: every other
// fixture lives under the family, and for a behavior-owned fixture the package, that owns
// its diagnostic.
func TestFlatFixtureSetIsExactlyTheUnownedOnes(t *testing.T) {
	canaryDir := filepath.Join(kitRoot(t), "tests", "canary")
	entries, err := os.ReadDir(canaryDir)
	if err != nil {
		t.Fatal(err)
	}
	var flat []string
	for _, entry := range entries {
		if entry.IsDir() && holdsExpect(filepath.Join(canaryDir, entry.Name())) {
			flat = append(flat, entry.Name())
		}
	}
	slices.Sort(flat)

	want := slices.Clone(flatFixtures)
	slices.Sort(want)
	if !slices.Equal(flat, want) {
		t.Fatalf("flat fixtures = %v, want %v", flat, want)
	}
}

// kitRoot is the checkout under test, reached from this package's own directory, so a
// test grading the real harness tree reads the tree it was built from.
func kitRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}
