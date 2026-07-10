package gate

// The PhasesCommand surface: signal handling across the process boundary, the
// benchkit phase table, and shellcheck file expansion. The runner engine's own
// tests (concurrency, output shape, exit codes, cancel) live in phases_test.go.

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

func TestPhasesCommandSignalCancelsRunningPhaseGroups(t *testing.T) {
	root := t.TempDir()
	pidfile := filepath.Join(root, "sleep.pid")
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve test executable: %v", err)
	}
	cmd := exec.Command(exe, "-test.run=^TestPhasesCommandSignalHelper$", "--", root)
	cmd.Env = append(os.Environ(),
		"BENCH_TEST_PHASES_SIGNAL_HELPER=1",
		"BENCH_TEST_PHASES_PIDFILE="+pidfile,
		"BENCH_TEST_PHASES_ROOT="+root,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})

	pid := waitForPIDFile(t, pidfile)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("send SIGINT to helper: %v", err)
	}
	err = cmd.Wait()
	if err == nil {
		t.Fatalf("helper exited 0, want 130; stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 130 {
		t.Fatalf("helper exit = %v, want code 130; stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
	waitForProcessExit(t, pid)
}

func TestPhasesCommandSignalHelper(t *testing.T) {
	if os.Getenv("BENCH_TEST_PHASES_SIGNAL_HELPER") != "1" {
		return
	}
	root := os.Getenv("BENCH_TEST_PHASES_ROOT")
	pidfile := os.Getenv("BENCH_TEST_PHASES_PIDFILE")
	benchkitPhasesForCommand = func(root, kit string) []Phase {
		return []Phase{{
			Name: "slow",
			Argv: []string{"bash", "-c", `sleep 30 & echo $! > "$1"; wait`, "bash", pidfile},
		}}
	}
	os.Exit(PhasesCommand([]string{root}, os.Stdout, os.Stderr))
}

func TestPhaseTable(t *testing.T) {
	root := "/tmp/root with spaces"
	kit := "/tmp/kit"
	phases := BenchkitPhases(root, kit)
	if len(phases) != 4 {
		t.Fatalf("BenchkitPhases len = %d, want 4: %#v", len(phases), phases)
	}
	if got, want := phaseNames(phases), []string{"conformance", "contract", "shellcheck", "canary"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("phase names = %#v, want %#v", got, want)
	}
	contractCount := 0
	for _, phase := range phases {
		for _, arg := range phase.Argv {
			if arg == "./internal/contract/..." {
				contractCount++
			}
		}
	}
	if contractCount != 1 {
		t.Fatalf("contract subtree argv count = %d, want exactly 1", contractCount)
	}
	if got, want := phaseNames(phasesForMode(phases, innerMode)), []string{"conformance", "contract", "shellcheck"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("inner phases = %#v, want %#v", got, want)
	}
}

func TestPhasesCommandRoutesCanaryToOwningPhase(t *testing.T) {
	original := benchkitPhasesForCommand
	t.Cleanup(func() { benchkitPhasesForCommand = original })
	benchkitPhasesForCommand = func(root, kit string) []Phase {
		return []Phase{
			{Name: "conformance", Argv: []string{"bash", "-c", "printf 'conformance\\n'"}},
			{Name: "contract", Argv: []string{"bash", "-c", "printf 'contract\\n'"}},
			{Name: "shellcheck", Argv: []string{"bash", "-c", "printf 'shellcheck\\n'"}},
			{Name: "canary", Argv: []string{"bash", "-c", "printf 'canary\\n'"}},
		}
	}
	for _, tc := range []struct {
		name  string
		phase string
	}{
		{name: "behavior-owned", phase: "contract"},
		{name: "load-validity-metadata", phase: "conformance"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("BENCH_CANARY_INNER", "1")
			t.Setenv("BENCH_CANARY_PHASE", tc.phase)

			var stdout, stderr bytes.Buffer
			if code := PhasesCommand([]string{t.TempDir()}, &stdout, &stderr); code != 0 {
				t.Fatalf("PhasesCommand = %d, want 0; stderr=%q", code, stderr.String())
			}
			want := tc.phase + "\ngate: green\n"
			if got := stdout.String(); got != want {
				t.Fatalf("stdout = %q, want %q", got, want)
			}
		})
	}
}

func TestPhaseTableBuildPhase(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "scripts", "go-build.sh"), "#!/usr/bin/env bash\n")
	writeFile(t, filepath.Join(root, "go.mod"), "module fixture\n")

	phases := BenchkitPhases(root, "/tmp/kit")
	if len(phases) != 5 {
		t.Fatalf("BenchkitPhases len = %d, want build + 4: %#v", len(phases), phases)
	}
	build := phases[0]
	if build.Name != "build" || !build.Serial {
		t.Fatalf("first phase = %#v, want serial build phase", build)
	}
	want := []string{"bash", filepath.Join(root, "scripts", "go-build.sh"), root, filepath.Join(root, "dist", "bench")}
	if !reflect.DeepEqual(build.Argv, want) {
		t.Fatalf("build argv = %#v, want %#v", build.Argv, want)
	}
	if got := phaseNames(phasesForMode(phases, innerMode)); got[0] != "build" {
		t.Fatalf("inner phases dropped the build phase: %#v", got)
	}

	// Without the Go build surface the table keeps its four-phase shape — and a
	// half-present surface (either file alone) counts as absent, not buildable.
	if bare := BenchkitPhases(t.TempDir(), "/tmp/kit"); len(bare) != 4 {
		t.Fatalf("bare root BenchkitPhases len = %d, want 4: %#v", len(bare), phaseNames(bare))
	}
	modOnly := t.TempDir()
	writeFile(t, filepath.Join(modOnly, "go.mod"), "module fixture\n")
	if got := BenchkitPhases(modOnly, "/tmp/kit"); len(got) != 4 {
		t.Fatalf("go.mod-only root BenchkitPhases len = %d, want 4: %#v", len(got), phaseNames(got))
	}
	helperOnly := t.TempDir()
	writeFile(t, filepath.Join(helperOnly, "scripts", "go-build.sh"), "#!/usr/bin/env bash\n")
	if got := BenchkitPhases(helperOnly, "/tmp/kit"); len(got) != 4 {
		t.Fatalf("helper-only root BenchkitPhases len = %d, want 4: %#v", len(got), phaseNames(got))
	}
}

func TestShellcheckPhaseExpandsHookAndLibShellFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "bin", "bench.sh"), "#!/usr/bin/env bash\n")
	writeFile(t, filepath.Join(root, ".bench", "hooks", "z.sh"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(root, ".bench", "hooks", "a.sh"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(root, ".bench", "hooks", "README"), "not shell\n")
	writeFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), "#!/bin/sh\n")
	writeFile(t, filepath.Join(root, ".bench", "lib", "notes.txt"), "not shell\n")

	var shellcheck Phase
	for _, phase := range BenchkitPhases("/tmp/root-under-grade", root) {
		if phase.Name == "shellcheck" {
			shellcheck = phase
			break
		}
	}
	want := []string{
		"shellcheck",
		"-S",
		"warning",
		"bin/bench.sh",
		".bench/hooks/a.sh",
		".bench/hooks/z.sh",
		".bench/lib/resolve-bench.sh",
	}
	if !reflect.DeepEqual(shellcheck.Argv, want) {
		t.Fatalf("shellcheck argv = %#v, want %#v", shellcheck.Argv, want)
	}
}

// TestShellcheckPhaseLintsNamedEnforcementShell is the bite-proof for the extended
// lint set: a lint error planted in a shift adapter, in .bench/gate.sh, and in the
// embedded pre-push asset must each turn the shellcheck phase red and be cited by
// path. Green today (adapters/gate.sh/asset were not in the argv, and the hook body
// hid in a Go string no linter read); red once each path joins the linted set. An
// argv typo that drops any of them silently un-lints it and fails here, not in review.
func TestShellcheckPhaseLintsNamedEnforcementShell(t *testing.T) {
	if _, err := exec.LookPath("shellcheck"); err != nil {
		t.Skip("shellcheck not installed")
	}
	kit := t.TempDir()
	const clean = "#!/usr/bin/env bash\ntrue\n"
	const badLint = "#!/usr/bin/env bash\ncd /tmp\n" // SC2164 (warning): cd without || exit
	writeFile(t, filepath.Join(kit, "bin", "bench.sh"), clean)
	writeFile(t, filepath.Join(kit, ".bench", "adapters", "claude"), badLint)
	writeFile(t, filepath.Join(kit, ".bench", "adapters", "codex"), clean)
	writeFile(t, filepath.Join(kit, ".bench", "adapters", "opencode"), clean)
	writeFile(t, filepath.Join(kit, ".bench", "gate.sh"), badLint)
	writeFile(t, filepath.Join(kit, "internal", "adopt", "prepush.sh"), badLint)

	var shellcheck Phase
	for _, p := range BenchkitPhases("/tmp/graded-root", kit) {
		if p.Name == "shellcheck" {
			shellcheck = p
			break
		}
	}

	var stdout, stderr bytes.Buffer
	rc := runPhases(context.Background(), kit, []Phase{shellcheck}, outerMode, &stdout, &stderr)
	out := stdout.String() + stderr.String()
	if rc == 0 {
		t.Fatalf("shellcheck phase stayed green with lint errors planted in the extended set:\n%s", out)
	}
	for _, cited := range []string{".bench/adapters/claude", ".bench/gate.sh", "internal/adopt/prepush.sh"} {
		if !strings.Contains(out, cited) {
			t.Fatalf("shellcheck phase did not lint %q (dropped from the argv set?):\n%s", cited, out)
		}
	}
}

func phaseNames(phases []Phase) []string {
	names := make([]string, 0, len(phases))
	for _, phase := range phases {
		names = append(names, phase.Name)
	}
	return names
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
