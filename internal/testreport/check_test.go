package testreport

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gocache"
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
		registry.ConsumerOnlyEnv:         "ambient-consumer",
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
	for _, name := range []string{registry.ConformanceChecksEnv, registry.ConformanceInheritedEnv, registry.ConsumerOnlyEnv, capability.LogEnv} {
		if strings.Contains(environment, name+"=") {
			t.Errorf("child environment retained %s:\n%s", name, environment)
		}
	}
}

func TestProseNamedCheckRunsCanonicalGraderAndPrintsEveryFinding(t *testing.T) {
	root := t.TempDir()
	writeProseCheckFile(t, root, ".bench/prose-exclusions", "")
	writeProseCheckFile(t, root, "first.md", strings.Repeat("word ", 27)+".\n")
	writeProseCheckFile(t, root, "second.md", strings.Repeat("word ", 27)+".\n")

	output, code := Command(root, []string{"--check", "prose"})
	if code != 1 {
		t.Fatalf("prose named check = %d, want 1\n%s", code, output)
	}
	for _, finding := range []string{"prose: \"first.md\" line 1: sentence of 27 words", "prose: \"second.md\" line 1: sentence of 27 words"} {
		if !strings.Contains(output, finding) {
			t.Errorf("prose named check output = %q, want finding %q", output, finding)
		}
	}
	help, helpCode := Command(root, []string{"--help"})
	if helpCode != 0 || !strings.Contains(help, "\n  prose\n") {
		t.Errorf("named-check help = %d, %q; want prose in the inventory", helpCode, help)
	}
}

func writeProseCheckFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestNamedCheckRunsFromKitAgainstLinkedConsumer(t *testing.T) {
	consumer := t.TempDir()
	kit := t.TempDir()
	marker := filepath.Join(t.TempDir(), "environment")
	goDir := t.TempDir()
	goPath := filepath.Join(goDir, "go")
	source := "#!/usr/bin/env bash\npwd > " + sanitize.ShellQuote(marker) + "\nenv >> " + sanitize.ShellQuote(marker) + "\nprintf '%s\\n' '{\"Action\":\"pass\",\"Package\":\"checkfixture\"}'\n"
	if err := os.WriteFile(goPath, []byte(source), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("BENCH_KIT", kit)
	t.Setenv(registry.ConsumerOnlyEnv, "1")

	selected := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(selected, []byte("selected"), 0o755); err != nil {
		t.Fatal(err)
	}
	previous := selectRunBinary
	selectRunBinary = func(context.Context, string) (*runbinary.Selection, error) {
		return &runbinary.Selection{Path: selected, SourceRoot: kit}, nil
	}
	t.Cleanup(func() { selectRunBinary = previous })

	output, code := Command(consumer, []string{"--check", "load-validity-metadata"})
	if code != 0 {
		t.Fatalf("linked named check = %d, want 0\n%s", code, output)
	}
	environment := readTestReportFile(t, marker)
	if !strings.HasPrefix(environment, kit+"\n") {
		t.Fatalf("named check working directory = %q, want kit %q", strings.SplitN(environment, "\n", 2)[0], kit)
	}
	for name, want := range map[string]string{
		registry.ConformanceRootEnv: consumer,
		registry.ConsumerOnlyEnv:    "1",
		"BENCH_KIT":                 kit,
	} {
		if !strings.Contains(environment, name+"="+want+"\n") {
			t.Errorf("child environment missing %s=%q:\n%s", name, want, environment)
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

func TestNamedCheckHelpListsEverySupportedCheck(t *testing.T) {
	output, code := Command(t.TempDir(), []string{"--help"})
	if code != 0 {
		t.Fatalf("named-check help = %d, want 0\n%s", code, output)
	}
	for _, check := range append(registry.Names(registry.Dev), gate.SystemPhaseName) {
		if !strings.Contains(output, "\n  "+check+"\n") {
			t.Errorf("named-check help missing %q:\n%s", check, output)
		}
	}
}

func TestUnknownNamedCheckReportsOperandAndInventory(t *testing.T) {
	inventory := "\nchecks:\n  " + strings.Join(namedChecks(), "\n  ") + "\n"
	for _, unknown := range []string{"not-registered", "release-evidence-probe"} {
		output, code := Command(t.TempDir(), []string{"--check", unknown})
		if code != 2 {
			t.Errorf("unknown named check %q = %d, want 2\n%s", unknown, code, output)
			continue
		}
		if want := "unknown check: " + unknown + inventory; output != want {
			t.Errorf("unknown named check %q = %q, want %q", unknown, output, want)
		}
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
	if want := "packages[1]{package,status,elapsed_ms}:\n  checkfixture,fail,0\n"; !strings.Contains(output, want) {
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

// LQ17: an unwritable cache directory refuses the focused run before the Go child starts,
// and the refusal names the hold error and the path. A focused run that compiled with the
// lock unheld would write archives a clean is removing.
func TestFocusedRunRefusesAnUnheldCache(t *testing.T) {
	home := unwritableCacheHome(t)
	dir, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}

	out, code := runGoTest(context.Background(), t.TempDir(), focusedRequest{}, []string{"true"}, []string{"HOME=" + home})

	if code != 1 {
		t.Fatalf("exit = %d, want 1; output=%q", code, out)
	}
	if !strings.Contains(out, "cache lock unavailable") || !strings.Contains(out, dir) {
		t.Fatalf("output = %q, want the cache refusal naming %q", out, dir)
	}
}

// unwritableCacheHome answers a HOME whose derived build cache directory exists and denies
// a write, which is the state that fails a holder's lock-file open. The directory is the
// derivation's own answer rather than a second spelling of it, and its mode is restored
// before the temporary tree is removed.
func unwritableCacheHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	dir, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(dir, 0o700); err != nil {
			t.Errorf("restore mode %s: %v", dir, err)
		}
	})
	if err := os.Chmod(dir, 0o500); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot remove directory write permission: %v", err))
	}
	probe, err := os.Create(filepath.Join(dir, "write-probe"))
	if err == nil {
		probe.Close()
		capability.Capability(t, capability.Privilege, "mode 0500 directory remains writable")
	}
	return home
}

// WF46: the conformance registry must never name a check equal to the system phase.
// The --check parser routes that one name past the registry lookup, so a registry entry
// of the same name would run the system suite instead, and refuse nothing.
func TestSystemCheckNameIsReserved(t *testing.T) {
	if check, found := registry.Find(gate.SystemPhaseName); found {
		t.Fatalf("registry names %q (implementation %s); the --check parser shadows it with the system suite", check.Name, check.Implementation)
	}
}
