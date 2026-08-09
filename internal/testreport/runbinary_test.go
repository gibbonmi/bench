package testreport

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestCommandOwnsOneSelectionAndPropagatesItToGoTest(t *testing.T) {
	root, marker := selectedPathTestModule(t)
	t.Setenv(runbinary.Env, "")
	if err := os.Unsetenv(runbinary.Env); err != nil {
		t.Fatal(err)
	}

	builds := 0
	var selectedPath, selectedDir string
	factory := runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			builds++
			selectedPath, selectedDir = output, filepath.Dir(output)
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
	installTestSelectionFactory(t, factory)

	if output, code := Command(root, nil); code != 0 {
		t.Fatalf("Command = %d\n%s", code, output)
	}
	if builds != 1 {
		t.Fatalf("builds = %d, want one", builds)
	}
	if got := strings.TrimSpace(readTestReportFile(t, marker)); got != selectedPath {
		t.Fatalf("child selection = %q, want %q", got, selectedPath)
	}
	if _, err := os.Stat(selectedDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("owned run directory after command = %v, want removed", err)
	}
}

func TestCommandReusesInheritedSelectionWithoutBuild(t *testing.T) {
	root, marker := selectedPathTestModule(t)
	dir := t.TempDir()
	selected := filepath.Join(dir, "bench")
	if err := os.WriteFile(selected, []byte("selected"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, selected)

	builds := 0
	factory := runbinary.Factory{
		Build:  func(context.Context, string, string) error { builds++; return nil },
		Verify: func(string, string) error { return nil },
	}
	installTestSelectionFactory(t, factory)

	if output, code := Command(root, nil); code != 0 {
		t.Fatalf("Command = %d\n%s", code, output)
	}
	if builds != 0 {
		t.Fatalf("builds = %d, want zero", builds)
	}
	if got := strings.TrimSpace(readTestReportFile(t, marker)); got != selected {
		t.Fatalf("child selection = %q, want %q", got, selected)
	}
}

func TestSeparateTopLevelCommandsSelectDifferentPrivatePaths(t *testing.T) {
	root, marker := selectedPathTestModule(t)
	t.Setenv(runbinary.Env, "")
	if err := os.Unsetenv(runbinary.Env); err != nil {
		t.Fatal(err)
	}
	factory := runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
	installTestSelectionFactory(t, factory)

	if output, code := Command(root, nil); code != 0 {
		t.Fatalf("first Command = %d\n%s", code, output)
	}
	first := strings.TrimSpace(readTestReportFile(t, marker))
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	if output, code := Command(root, nil); code != 0 {
		t.Fatalf("second Command = %d\n%s", code, output)
	}
	second := strings.TrimSpace(readTestReportFile(t, marker))
	if first == second {
		t.Fatalf("separate commands selected one path %q", first)
	}
}

func installTestSelectionFactory(t *testing.T, factory runbinary.Factory) {
	t.Helper()
	previous := selectRunBinary
	selectRunBinary = factory.ReuseOrOwn
	t.Cleanup(func() { selectRunBinary = previous })
}

func selectedPathTestModule(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	marker := filepath.Join(root, "selected-path")
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module selectedpathtest\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	source := `package selectedpathtest

import (
	"os"
	"testing"
)

func TestSelectedPath(t *testing.T) {
	if err := os.WriteFile("selected-path", []byte(os.Getenv("BENCH_RUN_BINARY")), 0o644); err != nil {
		t.Fatal(err)
	}
}
`
	if err := os.WriteFile(filepath.Join(root, "selected_test.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, marker
}

func readTestReportFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
