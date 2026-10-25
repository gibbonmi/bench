package testreport

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// inheritedRunBinary captures the operator's BENCH_RUN_BINARY before any test clears it.
// Package variables initialize before TestMain and before every t.Setenv.
var inheritedRunBinary, inheritedRunBinaryPresent = os.LookupEnv(runbinary.Env)

// packageRunBinary is the one run binary that the kit-copy tests of a package run share.
// Each focused run validates the shared binary against its own kit copy's sources.
var packageRunBinary = &runbinary.Shared{Raw: inheritedRunBinary, Present: inheritedRunBinaryPresent}

func TestMain(m *testing.M) {
	code := m.Run()
	_ = packageRunBinary.Close()
	os.Exit(code)
}

func TestOwnedSelectionBuildsOnce(t *testing.T) {
	builds := 0
	previous := packageRunBinary
	packageRunBinary = &runbinary.Shared{
		Factory: runbinary.Factory{
			TempRoot: t.TempDir(),
			Build: func(_ context.Context, _, output string) error {
				builds++
				return os.WriteFile(output, []byte("selected"), 0o755)
			},
			Verify: func(string, string) error { return nil },
		},
	}
	t.Cleanup(func() {
		_ = packageRunBinary.Close()
		packageRunBinary = previous
	})
	selected := map[string]bool{}
	for _, name := range []string{"first kit copy", "second kit copy"} {
		t.Run(name, func(t *testing.T) {
			installOwnedSelection(t, t.TempDir())
			selected[os.Getenv(runbinary.Env)] = true
		})
	}
	if builds != 1 || len(selected) != 1 || selected[""] {
		t.Fatalf("two kit copies ran the builder %d times and selected %v; want one build and one selection", builds, selected)
	}
}

func TestChangedGraphUsesOneClosedSelectionEnvironment(t *testing.T) {
	root, base, tip := changedCommandRepository(t, "", "", "changed/changed.go", "package changed\n")
	listEnvironment := filepath.Join(t.TempDir(), "list-environment")
	testEnvironment := filepath.Join(t.TempDir(), "test-environment")
	goDir := t.TempDir()
	writeChangedSubjectGo(t, filepath.Join(goDir, "go"), listEnvironment, testEnvironment, []listedPackage{{
		Dir:        filepath.Join(root, "changed"),
		ImportPath: "changedcommand/changed",
		Match:      []string{currentPackagePattern},
	}}, "changedcommand/changed")
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
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

	builds := 0
	var selected, source string
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, sourceRoot, output string) error {
			builds++
			source, selected = sourceRoot, output
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})

	output, code := Command(root, []string{"--changed", "--base", base, "--source-tip", tip})
	if code != 0 {
		t.Fatalf("changed command = %d\n%s", code, output)
	}
	if builds != 1 {
		t.Fatalf("selection builds = %d, want one", builds)
	}
	for _, environment := range []string{listEnvironment, testEnvironment} {
		got := readTestReportFile(t, environment)
		for _, name := range []string{
			registry.ConformanceRootEnv,
			registry.ConformanceTierEnv,
			registry.ConformanceScopeEnv,
			registry.ConformanceChecksEnv,
			registry.ConformanceInheritedEnv,
			capability.LogEnv,
		} {
			if strings.Contains(got, name+"=") {
				t.Errorf("%s retained ambient %s:\n%s", environment, name, got)
			}
		}
		for name, want := range map[string]string{"BENCH_KIT": source, runbinary.Env: selected} {
			if !strings.Contains(got, name+"="+want+"\n") {
				t.Errorf("%s missing %s=%q:\n%s", environment, name, want, got)
			}
		}
	}
}
