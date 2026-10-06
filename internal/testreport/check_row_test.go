package testreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const checkHeader = "check[1]{name,kind,tests_run,subjects}:\n"

// namedCheckOverEvents runs `Command` for one named check over one canned stream.
func namedCheckOverEvents(t *testing.T, root, check string, events []string, exit int) (string, int) {
	t.Helper()
	installCannedGo(t, cannedSet{name: t.Name(), events: events, exit: exit})
	installCannedSelection(t)
	return Command(root, []string{"--check", check})
}

// TestNamedCheckPrintsCheckRowFirst grades that the check row is the first block of the
// result and that the packages table follows it. (Coverage row TP3.)
func TestNamedCheckPrintsCheckRowFirst(t *testing.T) {
	output, code := namedCheckOverEvents(t, t.TempDir(), "line-routing", cannedSetNamed(t, "passing").events, 0)
	want := checkHeader + "  line-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestSystemCheckRowKind grades the kind cell of the system check. (Coverage row TP4.)
func TestSystemCheckRowKind(t *testing.T) {
	root := canonicalTestDir(t)
	t.Setenv("BENCH_KIT", root)
	output, code := namedCheckOverEvents(t, root, "system", cannedSetNamed(t, "passing").events, 0)
	want := checkHeader + "  system,system,1,0\n"
	if !strings.HasPrefix(output, want) || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, want)
	}
}

// TestNamedCheckRanNothingExitsOne grades the zero rule: a package pass with no run
// event is not evidence that the check ran. (Coverage row TP5.)
func TestNamedCheckRanNothingExitsOne(t *testing.T) {
	output, code := namedCheckOverEvents(t, t.TempDir(), "line-routing", []string{`{"Action":"pass","Package":"canned","Elapsed":0.01}`}, 0)
	want := checkHeader + "  line-routing,conformance,0,0\nerror: named check ran nothing"
	if !strings.HasPrefix(output, want) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, prefix %q)", code, output, want)
	}
}

// TestNamedCheckBuildFailureWinsOverZeroRule grades that a compile failure with no run
// event prints its diagnostic and not the zero-rule title. (Coverage row TP51.)
func TestNamedCheckBuildFailureWinsOverZeroRule(t *testing.T) {
	set := cannedSetNamed(t, "build-fail")
	output, code := namedCheckOverEvents(t, t.TempDir(), "line-routing", set.events, set.exit)
	want := checkHeader + "  line-routing,conformance,0,0\n"
	if !strings.HasPrefix(output, want) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, prefix %q)", code, output, want)
	}
	if !strings.Contains(output, "syntax error: unexpected }") || strings.Contains(output, "named check ran nothing") {
		t.Fatalf("Command output = %q, want the compile diagnostic and no zero-rule title", output)
	}
}

// TestPackageRunWithNoTestKeepsExitZero grades that the zero rule is a named-check rule
// only. (Coverage row TP6.)
func TestPackageRunWithNoTestKeepsExitZero(t *testing.T) {
	installCannedGo(t, cannedSet{name: t.Name(), events: []string{`{"Action":"pass","Package":"focusedfixture/chosen","Elapsed":0.01}`}})
	installCannedSelection(t)
	if output, code := Command(focusedTestModule(t), []string{"--package", "chosen"}); code != 0 {
		t.Fatalf("Command = (%d, %q), want exit 0", code, output)
	}
}

// TestChangedRunWithNoTestKeepsExitZero grades that a `--changed` run with no test keeps
// exit 0, because the merge reads that exit as its retry signal. (Coverage row TP55.)
func TestChangedRunWithNoTestKeepsExitZero(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
	goDir := t.TempDir()
	writeChangedSubjectGo(t, filepath.Join(goDir, "go"), filepath.Join(t.TempDir(), "list-environment"), filepath.Join(t.TempDir(), "test-environment"), []listedPackage{{
		Dir:        filepath.Join(root, "changed"),
		ImportPath: "changedcommand/changed",
		Match:      []string{currentPackagePattern},
	}}, "changedcommand/changed")
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installCannedSelection(t)
	if output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip}); code != 0 {
		t.Fatalf("Command = (%d, %q), want exit 0", code, output)
	}
}
