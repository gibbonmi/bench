package canary

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestSweepAbsolutizesRelativeRootFromDifferentCwd reproduces the field defect:
// `bench canary <relative-root>` run from a cwd other than the repo root must
// still hand the runner an absolute gate path, so the inner gate resolves
// against the repo root rather than each fixture's own temp cwd.
func TestSweepAbsolutizesRelativeRootFromDifferentCwd(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "myroot")
	mkdir(t, root)
	fixture := canaryFixture(root, "test-family", "relroot")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "targeted diagnostic\n")

	realGate := filepath.Join(root, ".bench", "gate.sh")
	write(t, realGate, "#!/bin/sh\n")

	t.Chdir(parent)
	runner := newResolvingRunner(realGate)
	if err := Sweep("myroot", runner.Run); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	if len(runner.calls) != 2 {
		t.Fatalf("runner calls = %d, want baseline + fixture", len(runner.calls))
	}
	for _, call := range runner.calls {
		if !filepath.IsAbs(call.Gate) {
			t.Fatalf("recorded gate %q is not absolute", call.Gate)
		}
		if call.Gate != realGate {
			t.Fatalf("recorded gate = %q, want %q", call.Gate, realGate)
		}
	}
}

// TestSweepDoesNotLetFixtureShadowTheRealGate reproduces the worse failure
// mode: a fixture whose files/ tree materializes its own .bench/gate.sh must
// not have that file win once the gate path is relative — an absolute gate
// path pins every run to the real repo gate regardless of what a fixture
// materializes into its temp cwd.
func TestSweepDoesNotLetFixtureShadowTheRealGate(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "myroot")
	mkdir(t, root)
	fixture := canaryFixture(root, "test-family", "shadow-gate")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "targeted diagnostic\n")
	// Materializes to work/myroot/.bench/gate.sh — the exact relative path the
	// unfixed gate join would resolve to against the fixture's own temp cwd.
	write(t, filepath.Join(fixture, "files", "myroot", "dot-bench", "gate.sh"), "echo shadow\n")

	realGate := filepath.Join(root, ".bench", "gate.sh")
	write(t, realGate, "#!/bin/sh\n")

	t.Chdir(parent)
	runner := newResolvingRunner(realGate)
	if err := Sweep("myroot", runner.Run); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	fixtureCall := runner.calls[len(runner.calls)-1]
	if fixtureCall.Gate != realGate {
		t.Fatalf("fixture ran gate %q, want the real repo gate %q (shadow gate must not win)", fixtureCall.Gate, realGate)
	}
}

// TestSweepCleansTrailingSlashInRoot covers the trailing-slash edge: the
// absolutize step must clean the argument so the recorded gate path carries
// no doubled or trailing separator.
func TestSweepCleansTrailingSlashInRoot(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "myroot")
	mkdir(t, root)
	fixture := canaryFixture(root, "test-family", "trailing-slash")
	mkdir(t, filepath.Join(fixture, "files"))
	write(t, filepath.Join(fixture, "EXPECT"), "targeted diagnostic\n")

	realGate := filepath.Join(root, ".bench", "gate.sh")
	write(t, realGate, "#!/bin/sh\n")

	t.Chdir(parent)
	runner := newResolvingRunner(realGate)
	if err := Sweep("myroot"+string(filepath.Separator), runner.Run); err != nil {
		t.Fatalf("Sweep err = %v", err)
	}
	fixtureCall := runner.calls[len(runner.calls)-1]
	if fixtureCall.Gate != realGate {
		t.Fatalf("recorded gate = %q, want cleaned absolute path %q", fixtureCall.Gate, realGate)
	}
	if strings.Contains(fixtureCall.Gate, "//") {
		t.Fatalf("recorded gate %q has a doubled separator", fixtureCall.Gate)
	}
}

// resolvingRunner reproduces defaultRunner's cwd-relative gate resolution
// (bash resolves a relative script path against cmd.Dir) without spawning
// bash, so the field symptom (exit 127, or a shadowed gate) becomes a
// deterministic unit assertion. realGate is the absolute path the caller
// considers "the real repo gate"; any other resolved path that exists is
// treated as a shadow gate, and one that doesn't exist reproduces exit 127.
type resolvingRunner struct {
	mu       sync.Mutex
	calls    []RunCall
	realGate string
}

func newResolvingRunner(realGate string) *resolvingRunner {
	return &resolvingRunner{realGate: realGate}
}

func (r *resolvingRunner) Run(call RunCall) RunResult {
	r.mu.Lock()
	r.calls = append(r.calls, call)
	r.mu.Unlock()

	resolved := call.Gate
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(call.Cwd, resolved)
	}
	switch {
	case resolved == r.realGate:
		if call.FixtureDir == "" {
			return RunResult{ExitCode: 1, Output: "baseline noise\n"}
		}
		return RunResult{ExitCode: 1, Output: "targeted diagnostic\n"}
	default:
		if info, err := os.Stat(resolved); err == nil && !info.IsDir() {
			return RunResult{ExitCode: 0, Output: "shadow gate ran\n"}
		}
		return RunResult{ExitCode: 127, Output: "bash: " + resolved + ": No such file or directory\n"}
	}
}
