package testreport

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
)

func TestNamedCheckOwnsConformanceEnvironment(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(t.TempDir(), "environment")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("BENCH_TEST_MARKER", marker)
	for name, value := range map[string]string{
		registry.ConformanceRootEnv:      "/ambient/root",
		registry.ConformanceTierEnv:      "ship",
		registry.ConformanceScopeEnv:     "ambient-scope",
		registry.ConformanceChecksEnv:    "ambient-checks",
		registry.ConformanceInheritedEnv: "ambient-inherited",
		capability.LogEnv:                filepath.Join(t.TempDir(), "capability-log"),
	} {
		t.Setenv(name, value)
	}

	var selected string
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			selected = output
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})

	output, code := Command(root, []string{"--check", "line-routing"})
	if code != 0 {
		t.Fatalf("named check = %d, want 0\n%s", code, output)
	}
	environment := readTestReportFile(t, marker)
	if want := "argv=test -trimpath -count=1 -json ./internal/conformance -run ^" + registry.RootConformanceTest + "$\n"; !strings.Contains(environment, want) {
		t.Errorf("named-check Go argv missing %q:\n%s", strings.TrimSpace(want), environment)
	}
	for name, want := range map[string]string{
		registry.ConformanceRootEnv:  root,
		registry.ConformanceTierEnv:  string(registry.Dev),
		registry.ConformanceScopeEnv: "line-routing",
		"BENCH_KIT":                  root,
		runbinary.Env:                selected,
	} {
		if !strings.Contains(environment, name+"="+want+"\n") {
			t.Errorf("child environment missing %s=%q:\n%s", name, want, environment)
		}
	}
	for _, name := range []string{registry.ConformanceChecksEnv, registry.ConformanceInheritedEnv, capability.LogEnv} {
		if strings.Contains(environment, name+"=") {
			t.Errorf("child environment retained %s:\n%s", name, environment)
		}
	}
}

func TestFocusedRequestGrammarRefusals(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(t.TempDir(), "go-ran")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	for _, args := range [][]string{
		{"--check", "line-routing", "./internal/conformance"},
		{"--check", "line-routing", "--package", "./internal/conformance"},
		{"--check", "line-routing", "--changed"},
		{"--check", "line-routing", "--run", "^TestRootConformance$"},
		{"--check", "line-routing", "--base", "HEAD"},
		{"--check", "line-routing", "--source-tip", "HEAD"},
		{"--base", "HEAD"},
		{"--source-tip", "HEAD"},
	} {
		output, code := Command(root, args)
		if code != 2 || !strings.HasPrefix(output, "usage: bench test") {
			t.Errorf("Command(%v) = %d, %q; want usage exit 2", args, code, output)
		}
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("grammar refusal started go child: %v", err)
	}
}

func TestNamedCheckRefusalMatrix(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(t.TempDir(), "go-ran")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	for _, args := range [][]string{
		{"--check", "not-registered"},
		{"--check", "release-evidence-probe"},
		{"--check", "line-routing", "--package", "./internal/conformance"},
		{"--check", "line-routing", "--changed"},
		{"--check", "line-routing", "--run", "^TestRootConformance$"},
	} {
		if output, code := Command(root, args); code != 2 || !strings.HasPrefix(output, "usage: bench test") {
			t.Errorf("Command(%v) = %d, %q; want usage exit 2", args, code, output)
		}
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("named-check refusal started go child: %v", err)
	}
}

func TestNamedCheckRefusesCorruptInheritedSelection(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(t.TempDir(), "go-ran")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv(runbinary.Env, filepath.Join(t.TempDir(), "missing-bench"))

	output, code := Command(root, []string{"--check", "line-routing"})
	if code != 1 || !strings.Contains(output, "Bench executable selection failed") {
		t.Fatalf("corrupt inherited selection = %d, %q; want selection refusal", code, output)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("corrupt inherited selection started go child: %v", err)
	}
}

func TestNamedCheckRunsOnlyRegisteredDevScope(t *testing.T) {
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	sourceTiming := registry.TimingPath(root)
	before, beforeErr := os.ReadFile(sourceTiming)
	existed := beforeErr == nil
	if beforeErr != nil && !os.IsNotExist(beforeErr) {
		t.Fatal(beforeErr)
	}
	realGo, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	selection, err := runbinary.Own(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = selection.Close() })
	previous := selectRunBinary
	selectRunBinary = func(context.Context, string) (*runbinary.Selection, error) { return selection, nil }
	t.Cleanup(func() { selectRunBinary = previous })

	timingRoot := t.TempDir()
	runChangedGit(t, timingRoot, "init", "-q")
	goDir := t.TempDir()
	writeGoBoundaryWrapper(t, filepath.Join(goDir, "go"), realGo, timingRoot)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	output, code := Command(root, []string{"--check", "ordinary-build-census"})
	if code != 0 {
		t.Fatalf("named check = %d\n%s", code, output)
	}
	after, afterErr := os.ReadFile(sourceTiming)
	afterExists := afterErr == nil
	if afterErr != nil && !os.IsNotExist(afterErr) {
		t.Fatal(afterErr)
	}
	if afterExists != existed || !bytes.Equal(after, before) {
		t.Fatalf("source timing changed during named check: before=%q after=%q", before, after)
	}
	lines := registry.ReadTimingLines(timingRoot)
	if len(lines) != 1 || !strings.Contains(lines[0], "ordinary-build-census") {
		t.Fatalf("named-check timing = %v, want only ordinary-build-census", lines)
	}
}

func TestNamedChecksWriteNoGateOwnedRecords(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(t.TempDir(), "environment")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	records := map[string]string{
		".git/bench-last-gate":            "cache",
		".git/bench-last-lane":            "lane",
		".git/bench-gate-owner":           "owner",
		".git/bench-gate-pin":             "pin",
		".git/bench-gate-evidence/record": "evidence",
		".git/refs/bench/green/main":      "green",
	}
	for path, want := range records {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(want), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if output, code := Command(root, []string{"--check", "line-routing"}); code != 0 {
		t.Fatalf("named check = %d\n%s", code, output)
	}
	for path, want := range records {
		if got := readTestReportFile(t, filepath.Join(root, path)); got != want {
			t.Errorf("gate-owned record %s = %q, want %q", path, got, want)
		}
	}
}

func writeCheckGo(t *testing.T, path, marker string) {
	t.Helper()
	source := "#!/usr/bin/env bash\nenv > \"" + marker + "\"\nprintf 'argv=%s\\n' \"$*\" >> \"" + marker + "\"\nprintf '%s\\n' '{\"Action\":\"pass\",\"Package\":\"checkfixture\"}'\n"
	if err := os.WriteFile(path, []byte(source), 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeGoBoundaryWrapper(t *testing.T, path, realGo, conformanceRoot string) {
	t.Helper()
	source := "#!/usr/bin/env bash\nexport BENCH_CONFORMANCE_ROOT=" + sanitize.ShellQuote(conformanceRoot) + "\nexec " + sanitize.ShellQuote(realGo) + " \"$@\"\n"
	if err := os.WriteFile(path, []byte(source), 0o755); err != nil {
		t.Fatal(err)
	}
}

// WF16, WF17, WF19: the system check owns the gate's system operands and environment.
// It carries the selected run binary, the kit, and the suite's root, and it carries no
// conformance variable, because the system suite is not a conformance scope.
func TestSystemCheckOwnsTheGateEnvironment(t *testing.T) {
	root := canonicalTestDir(t)
	marker := filepath.Join(t.TempDir(), "environment")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	for name, value := range map[string]string{
		registry.ConformanceRootEnv:      "/ambient/root",
		registry.ConformanceTierEnv:      "ship",
		registry.ConformanceScopeEnv:     "ambient-scope",
		registry.ConformanceChecksEnv:    "ambient-checks",
		registry.ConformanceInheritedEnv: "ambient-inherited",
	} {
		t.Setenv(name, value)
	}

	var selected string
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			selected = output
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	t.Setenv("BENCH_KIT", root)

	output, code := Command(root, []string{"--check", "system"})
	if code != 0 {
		t.Fatalf("system check = %d, want 0\n%s", code, output)
	}
	environment := readTestReportFile(t, marker)
	if want := "argv=test -trimpath -count=1 -json -tags=system ./internal/systemtest\n"; !strings.Contains(environment, want) {
		t.Errorf("system-check Go argv missing %q:\n%s", strings.TrimSpace(want), environment)
	}
	for name, want := range map[string]string{
		runbinary.Env:      selected,
		"BENCH_KIT":        root,
		gate.SystemRootEnv: root,
	} {
		if !strings.Contains(environment, name+"="+want+"\n") {
			t.Errorf("child environment missing %s=%q:\n%s", name, want, environment)
		}
	}
	for _, name := range []string{registry.ConformanceScopeEnv, registry.ConformanceRootEnv, registry.ConformanceTierEnv} {
		if strings.Contains(environment, name+"=") {
			t.Errorf("child environment retained %s:\n%s", name, environment)
		}
	}
}

// WF20: the system suite drives the kit's own wrapper and pool, so a graded root that
// is not the kit refuses before any Go starts.
func TestSystemCheckRefusesAForeignRoot(t *testing.T) {
	root := canonicalTestDir(t)
	marker := filepath.Join(t.TempDir(), "go-ran")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("BENCH_KIT", canonicalTestDir(t))

	output, code := Command(root, []string{"--check", "system"})
	if code != 1 || !strings.Contains(output, "system check unavailable") {
		t.Fatalf("foreign root = %d, %q; want the system-check refusal at exit 1", code, output)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("foreign-root refusal started go child: %v", err)
	}
}

// WF21: a red suite reads as red. The check renders the child's package terminal rather
// than its own opinion of the run.
func TestSystemCheckReportsAFailingSuite(t *testing.T) {
	root := canonicalTestDir(t)
	goDir := t.TempDir()
	writeFailingCheckGo(t, filepath.Join(goDir, "go"))
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	t.Setenv("BENCH_KIT", root)

	output, code := Command(root, []string{"--check", "system"})
	if code != 1 {
		t.Fatalf("failing system check = %d, want 1\n%s", code, output)
	}
	if want := "packages[1]{package,status}:\n  checkfixture,fail\n"; !strings.Contains(output, want) {
		t.Fatalf("failing system check output = %q, want %q", output, want)
	}
}

func writeFailingCheckGo(t *testing.T, path string) {
	t.Helper()
	source := "#!/usr/bin/env bash\nprintf '%s\\n' '{\"Action\":\"fail\",\"Package\":\"checkfixture\"}'\nexit 1\n"
	if err := os.WriteFile(path, []byte(source), 0o755); err != nil {
		t.Fatal(err)
	}
}

// canonicalTestDir returns a temporary directory with its symlinks already resolved.
// The run owner canonicalizes the source root it reports, so an unresolved fixture path
// would not compare equal to the environment the child carries.
func canonicalTestDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
