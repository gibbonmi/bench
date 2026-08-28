package consumers

import (
	"path/filepath"
	"testing"
)

// This is the package's one subprocess site. It drives the real go/packages `go list`
// path over the minimal fixture module under testdata, so the loader seam's contract is
// observed and not assumed. It skips nothing: a missing go tool is a failure here.
func TestLoadPackagesDrivesTheRealGoList(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "fixturemod"))
	if err != nil {
		t.Fatalf("resolve fixture module: %v", err)
	}
	pkgs, err := load(root, "./...")
	if err != nil {
		t.Fatalf("load fixture module: %v", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("loader returned no packages")
	}
	// One declaration must resolve to one match even though the loader delivers the
	// target package plainly and again as its test variant.
	matches, err := Resolve(pkgs, "target.Symbol")
	if err != nil {
		t.Fatalf("Resolve over loaded packages: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("declaration identity across test variants: got %d matches, want 1", len(matches))
	}
	rows, err := Find(pkgs, "target.Symbol", root)
	if err != nil {
		t.Fatalf("Find over loaded packages: %v", err)
	}
	got := summary(rows)
	want := []string{
		"consumer/consumer.go:6 via=reference enclosing=Direct",
		"consumer/consumer.go:9 via=reference enclosing=registry",
		"consumer/consumer.go:15 via=reference enclosing=Holder.Use",
		"target/target_test.go:7 via=reference enclosing=TestSymbol",
	}
	if len(got) != len(want) {
		t.Fatalf("loaded reference rows: got %d %v, want %d %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("row %d: got %q, want %q", i, got[i], want[i])
		}
	}
}
