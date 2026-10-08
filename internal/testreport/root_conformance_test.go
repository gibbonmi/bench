package testreport

import (
	"context"
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

// installOwnedSelection builds a Bench executable from root and makes it the run binary
// of each focused run in the test.
func installOwnedSelection(t *testing.T, root string) {
	t.Helper()
	selection, err := runbinary.Own(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = selection.Close() })
	previous := selectRunBinary
	selectRunBinary = func(context.Context, string) (*runbinary.Selection, error) { return selection, nil }
	t.Cleanup(func() { selectRunBinary = previous })
}

// rootConformanceRun runs the root-conformance check over root with BENCH_KIT set to
// kit, or unset when kit is empty. A canned run binary and a marker-writing go stand in
// for the child. It answers the output, the exit code, whether the check selected a run
// binary, and the child's recorded environment, which is empty when no Go child started.
func rootConformanceRun(t *testing.T, root, kit string) (output string, code int, selected bool, child string) {
	t.Helper()
	t.Setenv("BENCH_KIT", kit)
	if kit == "" {
		if err := os.Unsetenv("BENCH_KIT"); err != nil {
			t.Fatal(err)
		}
	}
	binary := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(binary, []byte("selected"), 0o755); err != nil {
		t.Fatal(err)
	}
	previous := selectRunBinary
	selectRunBinary = func(_ context.Context, source string) (*runbinary.Selection, error) {
		selected = true
		return &runbinary.Selection{Path: binary, SourceRoot: source}, nil
	}
	t.Cleanup(func() { selectRunBinary = previous })
	goDir := t.TempDir()
	marker := filepath.Join(goDir, "go-ran")
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	output, code = Command(root, []string{"--check", rootConformanceCheckName})
	data, err := os.ReadFile(marker)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	return output, code, selected, string(data)
}

func TestRootConformanceCheckRunsUnscoped(t *testing.T) {
	root := liveKitRoot(t)
	output, code, _, child := rootConformanceRun(t, root, root)
	if code != 0 || child == "" {
		t.Fatalf("root-conformance over the kit = %d, child started %t\n%s", code, child != "", output)
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
	root := gittest.KitCopy(t, filepath.Join("..", ".."))
	const drain, row = ".agents/commands/bench-drain.md", "| `kit-edit` | a `craft-synthesis` kit edit |\n"
	body := readTestReportFile(t, filepath.Join(root, drain))
	if strings.Count(body, row) != 1 {
		t.Fatalf("%s holds the row %q %d times, want 1", drain, row, strings.Count(body, row))
	}
	writeProseCheckFile(t, root, drain, strings.Replace(body, row, "", 1))
	const diagnostic = "row-next-grammar drift: " + drain + " token table lacks token 'kit-edit'"
	t.Setenv("BENCH_KIT", root)
	installOwnedSelection(t, root)

	output, code := Command(root, []string{"--full", "--check", rootConformanceCheckName})
	if code != 1 || !strings.Contains(output, diagnostic) {
		t.Fatalf("root-conformance over a planted defect = %d, want 1 and the diagnostic %q\n%s", code, diagnostic, output)
	}
}

func TestRootConformanceCheckRefusesAForeignRoot(t *testing.T) {
	root := canonicalTestDir(t)
	if gate.RootConformanceRuns(root, liveKitRoot(t)) {
		t.Fatalf("gate predicate answers true for the foreign root %s", root)
	}
	output, code, selected, child := rootConformanceRun(t, root, liveKitRoot(t))
	if code != 1 || !strings.Contains(output, rootConformanceUnavailable) {
		t.Fatalf("root-conformance over a foreign root = %d, %q; want the refusal at exit 1", code, output)
	}
	if selected || child != "" {
		t.Fatalf("foreign-root refusal selected a run binary %t or started a Go child %t", selected, child != "")
	}
}

func TestChecksFaceListsRootConformance(t *testing.T) {
	_, rows := checksFace(t, t.TempDir())
	if row := rows[rootConformanceCheckName]; len(row) != 3 || row[1] != "conformance" {
		t.Fatalf("--checks row %q; want %s with the kind conformance", row, rootConformanceCheckName)
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
