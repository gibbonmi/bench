package testreport

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// liveKitRoot answers the resolved checkout that holds this package. It declares the
// conformance entry test, and these tests only read it.
func liveKitRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if root, err = filepath.EvalSymlinks(root); err != nil {
		t.Fatal(err)
	}
	return root
}

// installOwnedSelection makes kit the BENCH_KIT of the test, and each focused run builds
// its run binary from kit. The kit is set after the factory install, because that install
// clears BENCH_KIT.
func installOwnedSelection(t *testing.T, kit string) {
	t.Helper()
	installTestSelectionFactory(t, runbinary.Factory{})
	t.Setenv("BENCH_KIT", kit)
}

// rootConformanceRun runs the root-conformance check over root with BENCH_KIT set to
// kit, or unset when kit is empty. A canned run binary and a marker-writing go stand in
// for the child. It answers the output, the exit code, whether the check selected a run
// binary, and the child's recorded environment, which is empty when no Go child started.
func rootConformanceRun(t *testing.T, root, kit string) (output string, code int, selected bool, child string) {
	t.Helper()
	selection := installCannedSelection(t)
	if kit != "" {
		t.Setenv("BENCH_KIT", kit)
	}
	goDir := t.TempDir()
	marker := filepath.Join(goDir, "go-ran")
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	output, code = Command(root, []string{"--check", RootConformanceCheckName})
	data, err := os.ReadFile(marker)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	return output, code, *selection, string(data)
}

func TestRootConformanceCheckRunsUnscoped(t *testing.T) {
	root := liveKitRoot(t)
	output, code, selected, child := rootConformanceRun(t, root, root)
	if code != 0 || !selected || child == "" {
		t.Fatalf("root-conformance over the kit = %d, selected a run binary %t, child started %t\n%s", code, selected, child != "", output)
	}
	lines := strings.Split(child, "\n")
	for _, want := range []string{
		"argv=test -trimpath -count=1 -json ./internal/conformance -run " + namedCheckRunPattern(),
		registry.ConformanceRootEnv + "=" + root,
		registry.ConformanceTierEnv + "=" + string(registry.Dev),
		registry.ConformanceScopeEnv + "=",
	} {
		if !slices.Contains(lines, want) {
			t.Errorf("root-conformance child lacks the line %q:\n%s", want, child)
		}
	}
}

func TestRootConformanceCheckPrintsFailures(t *testing.T) {
	root := gittest.KitCopy(t, liveKitRoot(t))
	const drain, row = ".agents/commands/bench-drain.md", "| `kit-edit` | a `craft-synthesis` kit edit |\n"
	body := readTestReportFile(t, filepath.Join(root, drain))
	if strings.Count(body, row) != 1 {
		t.Fatalf("%s holds the row %q %d times, want 1", drain, row, strings.Count(body, row))
	}
	writeProseCheckFile(t, root, drain, strings.Replace(body, row, "", 1))
	const diagnostic = "row-next-grammar drift: " + drain + " token table lacks token 'kit-edit'"
	installOwnedSelection(t, root)

	output, code := Command(root, []string{"--full", "--check", RootConformanceCheckName})
	if code != 1 || !strings.Contains(output, diagnostic) {
		t.Fatalf("root-conformance over a planted defect = %d, want 1 and the diagnostic %q\n%s", code, diagnostic, output)
	}
}

// The kit case grades the kit clause alone: its root declares the entry test, and only
// the kit names another directory.
func TestRootConformanceCheckRefusesAForeignRoot(t *testing.T) {
	kit := liveKitRoot(t)
	for name, pair := range map[string][2]string{
		"foreign root": {canonicalTestDir(t), kit},
		"foreign kit":  {kit, canonicalTestDir(t)},
	} {
		t.Run(name, func(t *testing.T) {
			root, benchKit := pair[0], pair[1]
			if gate.RootConformanceRuns(root, benchKit) {
				t.Fatalf("gate predicate answers true for root %s under kit %s", root, benchKit)
			}
			output, code, selected, child := rootConformanceRun(t, root, benchKit)
			if code != 1 || !strings.Contains(output, rootConformanceUnavailable) {
				t.Fatalf("root-conformance = %d, %q; want the refusal at exit 1", code, output)
			}
			if selected || child != "" {
				t.Fatalf("refusal selected a run binary %t or started a Go child %t", selected, child != "")
			}
		})
	}
}

func TestChecksFaceListsRootConformance(t *testing.T) {
	_, rows := checksFace(t, t.TempDir())
	if row := rows[RootConformanceCheckName]; len(row) != 3 || row[1] != "conformance" {
		t.Fatalf("--checks row %q; want %s with the kind conformance", row, RootConformanceCheckName)
	}
}

func TestRootConformanceCheckFollowsTheGatePredicate(t *testing.T) {
	kit := liveKitRoot(t)
	withoutEntryTest := canonicalTestDir(t)
	answers := map[bool]int{}
	for name, pair := range map[string][2]string{
		"kit root":             {kit, kit},
		"root without entry":   {withoutEntryTest, withoutEntryTest},
		"foreign root":         {canonicalTestDir(t), kit},
		"kit root with no kit": {kit, ""},
	} {
		root, benchKit := pair[0], pair[1]
		t.Setenv("BENCH_KIT", benchKit)
		want := gate.RootConformanceRuns(root, gate.KitRoot(root))
		answers[want]++
		_, _, _, child := rootConformanceRun(t, root, benchKit)
		if ran := child != ""; ran != want {
			t.Errorf("%s: root-conformance ran %t; the gate predicate answers %t", name, ran, want)
		}
	}
	if answers[true] == 0 || answers[false] == 0 {
		t.Fatalf("gate predicate answers %v; the fixtures must reach both answers", answers)
	}
}
