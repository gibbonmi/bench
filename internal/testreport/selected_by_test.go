package testreport

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

// changedCommandOverOnePackage runs `Command --changed` over one changed Go file of the
// package `changedcommand/changed`, with a canned `go list` and `go test`. When empty is
// true, the run takes the changed commit as its base too, so the diff has no changed path.
func changedCommandOverOnePackage(t *testing.T, empty bool) (string, int) {
	t.Helper()
	root, base, tip := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
	goDir := t.TempDir()
	writeChangedSubjectGo(t, filepath.Join(goDir, "go"), filepath.Join(t.TempDir(), "list-environment"), filepath.Join(t.TempDir(), "test-environment"), []listedPackage{{
		Dir:        filepath.Join(root, "changed"),
		ImportPath: "changedcommand/changed",
		Match:      []string{currentPackagePattern},
	}}, "changedcommand/changed")
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installCannedSelection(t)
	if empty {
		base = tip
	}
	return Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
}

// TestChangedRowsCarrySelectedBy grades that a `--changed` result prints the cause cell
// on each packages row.
func TestChangedRowsCarrySelectedBy(t *testing.T) {
	output, code := changedCommandOverOnePackage(t, false)
	want := "packages[1]{package,status,elapsed_ms,tests_run,selected_by}:\n  changedcommand/changed,pass,0,0,changed\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestChangedEmptyDiffCarriesSelectedBy grades that a `--changed` run over a diff with no
// changed path prints the empty packages table with the cause column.
func TestChangedEmptyDiffCarriesSelectedBy(t *testing.T) {
	output, code := changedCommandOverOnePackage(t, true)
	want := "packages[0]{package,status,elapsed_ms,tests_run,selected_by}:"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

func TestChangedNonGoSubjectRendersExplicitEmpty(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "README.md", "docs\n", "README.md", "")
	installTestSelectionFactory(t, runbinary.Factory{TempRoot: t.TempDir(), Build: func(_ context.Context, _, output string) error {
		return os.WriteFile(output, []byte("selected"), 0o755)
	}, Verify: func(string, string) error { return nil }})
	output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
	if code != 0 {
		t.Fatalf("emptyReport code = %d\n%s", code, output)
	}
	for _, table := range []string{"packages[0]{package,status,elapsed_ms,tests_run,selected_by}:\n", "failures[0]", "skips[0]"} {
		if !strings.Contains(output, table) {
			t.Fatalf("empty output = %q, want %q", output, table)
		}
	}
}
