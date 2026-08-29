package consumers

import (
	"os"
	"path/filepath"
	"strings"
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
		"consumer/consumer.go:6 via=call enclosing=Direct",
		"consumer/consumer.go:9 via=reference enclosing=registry",
		"consumer/consumer.go:15 via=call enclosing=Holder.Use",
		"target/target_test.go:7 via=call enclosing=TestSymbol",
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

// stubLoadError swaps the loader seam for one that reports err, so a refusal test drives
// the command's whole path over an error the real go tool produced.
func stubLoadError(t *testing.T, err error) {
	t.Helper()
	original := load
	load = func(string, ...string) ([]*Package, error) { return nil, err }
	t.Cleanup(func() { load = original })
}

// CS8 (story 7): an ill-typed tree refuses, and the refusal carries the first error
// position the real loader reported. The position comes from the go tool here, so a
// tolerant loader or a reshaped message is observable rather than assumed.
func TestIllTypedTreeRefusesNamingTheFirstErrorPosition(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "illtyped"))
	if err != nil {
		t.Fatalf("resolve fixture module: %v", err)
	}
	_, loadErr := load(root, "./...")
	if loadErr == nil {
		t.Fatal("ill-typed fixture module loaded without an error")
	}
	const position = "bad.go:6:14"
	if !strings.Contains(loadErr.Error(), position) {
		t.Fatalf("loader error = %q, want the first error position %q", loadErr, position)
	}
	stubLoadError(t, loadErr)
	out, code := run(t, "target.Symbol")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; out=%q", code, out)
	}
	if !strings.HasPrefix(out, "error: ") {
		t.Fatalf("stdout = %q, want a structured error line", out)
	}
	if !strings.Contains(out, position) {
		t.Fatalf("stdout = %q, want the first error position %q", out, position)
	}
	if !strings.Contains(out, "type-check") {
		t.Fatalf("stdout = %q, want the type-check hint", out)
	}
	if strings.Contains(out, "citation[") || strings.Contains(out, "consumers[") {
		t.Fatalf("stdout = %q, want no citation row and no result block", out)
	}
}

// CS9 (story 8): with the go tool absent from PATH the refusal names the binary and the
// remedy, so an agent reads an action instead of an exec stack.
func TestMissingGoBinaryRefusesNamingTheBinaryAndRemedy(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "fixturemod"))
	if err != nil {
		t.Fatalf("resolve fixture module: %v", err)
	}
	// PATH is emptied for the load alone and restored before the command runs, because an
	// emptied PATH also hides git, and the command resolves its repository root first.
	restore := os.Getenv("PATH")
	t.Setenv("PATH", t.TempDir())
	_, loadErr := load(root, "./...")
	if err := os.Setenv("PATH", restore); err != nil {
		t.Fatalf("restore PATH: %v", err)
	}
	if loadErr == nil {
		t.Fatal("loader found a go tool on an emptied PATH")
	}
	stubLoadError(t, loadErr)
	out, code := run(t, "target.Symbol")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; out=%q", code, out)
	}
	if !strings.Contains(out, "go binary not found on PATH") {
		t.Fatalf("stdout = %q, want the missing binary named", out)
	}
	if !strings.Contains(out, "install Go") {
		t.Fatalf("stdout = %q, want the remedy", out)
	}
	if strings.Contains(out, "citation[") {
		t.Fatalf("stdout = %q, want no citation row", out)
	}
}

// A query of a test function reaches the file go/packages synthesizes for the test
// binary, which sits in the build cache rather than in the tree. The loader must not
// deliver it: a row there names no file a reviewer can open, and the citation promises a
// replay of the checkout alone.
func TestLoadedRowsNeverEscapeTheRoot(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "fixturemod"))
	if err != nil {
		t.Fatalf("resolve fixture module: %v", err)
	}
	pkgs, err := load(root, "./...")
	if err != nil {
		t.Fatalf("load fixture module: %v", err)
	}
	rows, err := Find(pkgs, "target.TestSymbol", root)
	if err != nil {
		t.Fatalf("Find over loaded packages: %v", err)
	}
	for _, r := range rows {
		if strings.HasPrefix(r.File, "..") || filepath.IsAbs(r.File) {
			t.Errorf("row %q escapes the root; every enumerated file must sit inside it", r.File)
		}
	}
	got := summary(rows)
	if len(got) != 0 {
		t.Fatalf("loaded reference rows for a test function: got %d %v, want none in the tree", len(got), got)
	}
}
