package testreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// changedCommandOverOnePackage runs `Command --changed` over one changed Go file of the
// package `changedcommand/changed`, with a canned `go list` and `go test`.
func changedCommandOverOnePackage(t *testing.T) (string, int) {
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
	return Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
}

// TestChangedRowsCarrySelectedBy grades that a `--changed` result prints the cause cell
// on each packages row.
func TestChangedRowsCarrySelectedBy(t *testing.T) {
	output, code := changedCommandOverOnePackage(t)
	want := "packages[1]{package,status,elapsed_ms,tests_run,selected_by}:\n  changedcommand/changed,pass,0,0,changed\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}
