//go:build system

package systemtest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/freshness"
)

func TestSessionStartTE5DiagnosesPartialEnvironment(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	result := fixture.run(t, nil)
	if result.code != 0 {
		t.Fatalf("TE5 partial environment exit = %d, want 0; stdout=%q stderr=%q", result.code, result.stdout, result.stderr)
	}
	for _, want := range []string{"environment closure is partial", fixture.goExecutable} {
		if !strings.Contains(result.stdout, want) {
			t.Fatalf("TE5 partial environment stdout missing %q: %q", want, result.stdout)
		}
	}
}

func TestSessionStartTE5DiscoversGoWithMarkerAbsentOrEmpty(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	for _, marker := range []struct {
		name, override string
	}{
		{name: "absent", override: "ENVMAN_LOAD"},
		{name: "empty", override: "ENVMAN_LOAD="},
	} {
		t.Run(marker.name, func(t *testing.T) {
			result := fixture.run(t, []string{marker.override})
			for _, want := range []string{"environment closure is partial", fixture.goExecutable} {
				if !strings.Contains(result.stdout, want) {
					t.Fatalf("TE5 %s-marker stdout missing %q: %q", marker.name, want, result.stdout)
				}
			}
		})
	}
}

func TestSessionStartTE6PrependsTheDiscoveredDirectory(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	result := fixture.run(t, nil)
	want := "export PATH='" + filepath.Dir(fixture.goExecutable) + "':\"$PATH\""
	if !strings.Contains(result.stdout, want) {
		t.Fatalf("TE6 recovery command missing %q: %q", want, result.stdout)
	}
}

func TestSessionStartTE7DoesNotExecuteDiscoveredGo(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	sentinel := filepath.Join(t.TempDir(), "go-executed")
	result := fixture.run(t, []string{"GO_EXECUTION_SENTINEL=" + sentinel})
	if result.code != 0 {
		t.Fatalf("TE7 hook exit = %d, want 0", result.code)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("TE7 discovered Go execution sentinel exists: %v", err)
	}
}

func TestSessionStartTE7AcceptsSymlinkWithoutExecutingTarget(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	target := fixture.goExecutable + "-target"
	if err := os.Rename(fixture.goExecutable, target); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Base(target), fixture.goExecutable); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(t.TempDir(), "symlink-target-executed")
	result := fixture.run(t, []string{"GO_EXECUTION_SENTINEL=" + sentinel})
	if result.code != 0 || !strings.Contains(result.stdout, fixture.goExecutable) {
		t.Fatalf("TE7 symlink discovery = (%d, %q, %q), want accepted path %q", result.code, result.stdout, result.stderr, fixture.goExecutable)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("TE7 symlink target execution sentinel exists: %v", err)
	}
}

func TestSessionStartTE8MissingAndPartialCasesExitZero(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	for _, test := range []struct {
		name  string
		extra []string
	}{
		{name: "partial"},
		{name: "missing", extra: []string{"HOME=" + t.TempDir()}},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := fixture.run(t, test.extra)
			if result.code != 0 {
				t.Fatalf("TE8 %s exit = %d, want 0; stdout=%q stderr=%q", test.name, result.code, result.stdout, result.stderr)
			}
		})
	}
}

func TestSessionStartTE9HealthyHarnessPrintsNoRecovery(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	healthyPath := filepath.Dir(fixture.goExecutable) + string(os.PathListSeparator) + fixture.harnessPath
	result := fixture.run(t, []string{"PATH=" + healthyPath})
	for _, absent := range []string{"environment closure is partial", "recover without replacing harness tools"} {
		if strings.Contains(result.stdout, absent) {
			t.Fatalf("TE9 healthy stdout contains %q: %q", absent, result.stdout)
		}
	}
}

func TestSessionStartTE10MissingGoPrintsNoPathAssignment(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	result := fixture.run(t, []string{"HOME=" + t.TempDir()})
	if !strings.Contains(result.stdout, "Go is absent from PATH") {
		t.Fatalf("TE10 missing-Go diagnosis absent: %q", result.stdout)
	}
	if strings.Contains(result.stdout, "export PATH=") {
		t.Fatalf("TE10 missing-Go output contains a PATH assignment: %q", result.stdout)
	}
}

func TestSessionStartTE13QuotesOnePathElement(t *testing.T) {
	fixture := newSessionEnvironmentFixtureWithGoDir(t, "Go SDK [stable]*")
	result := fixture.run(t, nil)
	want := "export PATH='" + filepath.Dir(fixture.goExecutable) + "':\"$PATH\""
	if !strings.Contains(result.stdout, want) {
		t.Fatalf("TE13 quoted recovery missing %q: %q", want, result.stdout)
	}
}

func TestSessionStartTE14RejectsUnsafeDiscoveryOutput(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	shellPath := discoveryShellPath(t)
	valid := fixture.goExecutable
	for _, test := range []struct {
		name, output string
	}{
		{name: "relative", output: "relative/go"},
		{name: "nonexistent", output: filepath.Join(t.TempDir(), "missing-go")},
		{name: "multiline", output: valid + "\n" + valid},
		{name: "control-bearing", output: valid + "\x1bunsafe"},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := shellPath + string(os.PathListSeparator) + fixture.harnessPath
			result := fixture.run(t, []string{"PATH=" + path, "DISCOVERY_OUTPUT=" + test.output})
			if !strings.Contains(result.stdout, "Go is absent from PATH") || strings.Contains(result.stdout, "export PATH=") {
				t.Fatalf("TE14 %s output entered a recovery command: %q", test.name, result.stdout)
			}
		})
	}
}

func TestSessionStartNonzeroDiscoveryPrintsNoRecovery(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	path := discoveryShellPath(t) + string(os.PathListSeparator) + fixture.harnessPath
	result := fixture.run(t, []string{
		"PATH=" + path,
		"DISCOVERY_OUTPUT=" + fixture.goExecutable,
		"DISCOVERY_EXIT=23",
	})
	if result.code != 0 || !strings.Contains(result.stdout, "Go is absent from PATH") || strings.Contains(result.stdout, "export PATH=") {
		t.Fatalf("nonzero discovery result = (%d, %q, %q), want zero exit and no recovery assignment", result.code, result.stdout, result.stderr)
	}
}

func TestSessionStartTE15BoundsDiscoveryAndContinues(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	path := hangingDiscoveryShellPath(t) + string(os.PathListSeparator) + fixture.harnessPath
	started := time.Now()
	result := fixture.run(t, []string{"PATH=" + path})
	elapsed := time.Since(started)
	if elapsed >= bounds.EnvironmentDiscoveryTimeout+time.Second {
		t.Fatalf("TE15 hook elapsed = %s, want discovery stopped near %s; result=(%d, %q, %q)", elapsed, bounds.EnvironmentDiscoveryTimeout, result.code, result.stdout, result.stderr)
	}
	if result.code != 0 || strings.Contains(result.stdout, "export PATH=") || !strings.Contains(result.stdout, "guard_scan:") {
		t.Fatalf("TE15 bounded continuation = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
}

func TestSessionStartTE16KillsDiscoveryDescendants(t *testing.T) {
	fixture := newSessionEnvironmentFixture(t)
	shellPath, child := descendantDiscoveryShellPath(t)
	pidFile := filepath.Join(t.TempDir(), "descendant.pid")
	sentinel := filepath.Join(t.TempDir(), "descendant-survived")
	path := shellPath + string(os.PathListSeparator) + fixture.harnessPath
	result := fixture.run(t, []string{
		"PATH=" + path,
		"DISCOVERY_CHILD=" + child,
		"DESCENDANT_PID_FILE=" + pidFile,
		"DESCENDANT_SENTINEL=" + sentinel,
	})
	if result.code != 0 {
		t.Fatalf("TE16 hook exit = %d, want 0; stdout=%q stderr=%q", result.code, result.stdout, result.stderr)
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatal(err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	requireProcessGone(t, pid, "TE16 discovery descendant")
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("TE16 descendant sentinel exists: %v", err)
	}
}

func discoveryShellPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\nif [ \"${1:-}\" != -c ]; then exec /bin/bash \"$@\"; fi\nprintf '%s' \"$DISCOVERY_OUTPUT\"\nexit \"${DISCOVERY_EXIT:-0}\"\n"
	if err := os.WriteFile(filepath.Join(dir, "bash"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func hangingDiscoveryShellPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\nif [ \"${1:-}\" != -c ]; then exec /bin/bash \"$@\"; fi\n/bin/sleep 30\n"
	if err := os.WriteFile(filepath.Join(dir, "bash"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func descendantDiscoveryShellPath(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	shell := "#!/bin/sh\nif [ \"${1:-}\" != -c ]; then exec /bin/bash \"$@\"; fi\n\"$DISCOVERY_CHILD\" &\nwait\n"
	child := "#!/bin/sh\nprintf '%s' \"$$\" > \"$DESCENDANT_PID_FILE\"\ntrap 'printf survived > \"$DESCENDANT_SENTINEL\"' EXIT TERM\n/bin/sleep 4\n"
	if err := os.WriteFile(filepath.Join(dir, "bash"), []byte(shell), 0o755); err != nil {
		t.Fatal(err)
	}
	childPath := filepath.Join(dir, "discovery-child")
	if err := os.WriteFile(childPath, []byte(child), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir, childPath
}

type sessionEnvironmentFixture struct {
	repo, hook, home, goExecutable, harnessPath string
}

func newSessionEnvironmentFixture(t *testing.T) sessionEnvironmentFixture {
	return newSessionEnvironmentFixtureWithGoDir(t, "clean-login-go")
}

func newSessionEnvironmentFixtureWithGoDir(t *testing.T, directory string) sessionEnvironmentFixture {
	t.Helper()
	repo := owner.repos[2]
	installWrapper(t, repo)
	goMod := filepath.Join(repo, "go.mod")
	if err := os.WriteFile(goMod, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(goMod) })
	home := t.TempDir()
	goDir := filepath.Join(t.TempDir(), directory, "bin")
	if err := os.MkdirAll(goDir, 0o755); err != nil {
		t.Fatal(err)
	}
	goExecutable := filepath.Join(goDir, "go")
	if err := os.WriteFile(goExecutable, []byte("#!/bin/sh\nprintf invoked > \"$GO_EXECUTION_SENTINEL\"\nexit 99\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	profile := fmt.Sprintf("if [ -z \"${ENVMAN_LOAD+x}\" ]; then export PATH='%s':\"$PATH\"; fi\n", goDir)
	if err := os.WriteFile(filepath.Join(home, ".bash_profile"), []byte(profile), 0o644); err != nil {
		t.Fatal(err)
	}
	return sessionEnvironmentFixture{
		repo:         repo,
		hook:         sessionStartHook(t),
		home:         home,
		goExecutable: goExecutable,
		harnessPath:  privateToolPath(t, "git", "bash", "uname", "dirname", "basename", "readlink", "tr"),
	}
}

func (f sessionEnvironmentFixture) run(t *testing.T, extra []string) processResult {
	t.Helper()
	overrides := []string{
		"BENCH_KIT=" + owner.kit,
		"BENCH_RUN_BINARY=" + owner.selected.path,
		"BENCH_HOME=" + filepath.Join(f.home, ".bench"),
		"HOME=" + f.home,
		"PATH=" + f.harnessPath,
		"ENVMAN_LOAD=loaded",
	}
	overrides = mergeEnvironment(overrides, extra)
	return runHookWithDeadline(t, f.repo, overrides, f.hook, bounds.TestDeadline(bounds.EnvironmentDiscoveryTimeout))
}

func runHookWithDeadline(t *testing.T, dir string, overrides []string, hook string, deadline time.Duration) processResult {
	t.Helper()
	shell, err := exec.LookPath("bash")
	if err != nil {
		t.Fatal(err)
	}
	timeout, err := exec.LookPath("timeout")
	if err != nil {
		t.Fatal(err)
	}
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runAt(dir, overrides, timeout, "--signal=KILL", deadline.String(), shell, hook)
}

// TestSessionStartNamesTheBuildCommandWhenTheCoreIsUnreachable drives
// .bench/hooks/session-start.sh as a real subprocess in the state the hint exists for:
// a session opening in a repository whose Bench core cannot be reached, which the
// wrapper answers with exit 127 and no dashboard. Every launch goes through the owner's
// process ledger. The expected hint is read from freshness.RebuildAction — the one
// source of the rebuild invocation — so the shell copy in the hook is pinned to the Go
// source rather than asserted against a second hand-written spelling of it.
func TestSessionStartNamesTheBuildCommandWhenTheCoreIsUnreachable(t *testing.T) {
	hook := sessionStartHook(t)
	repo := owner.repos[2]
	installWrapper(t, repo)
	root, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	rebuild := freshness.RebuildAction(root)

	// H23 — the unreachable half. BENCH_KIT names a directory holding no binary and
	// BENCH_RUN_BINARY is removed outright (a bare name in the overrides), so the
	// wrapper runs its real no-binary resolution instead of the ambient one the test
	// process was launched with. PATH carries only the tools the hook and the wrapper
	// themselves run, so `bench` cannot be reached that way either.
	missing := []string{
		"BENCH_KIT=" + filepath.Join(t.TempDir(), "missing-kit"),
		"BENCH_HOME=" + filepath.Join(t.TempDir(), "home"),
		"PATH=" + privateToolPath(t, "git", "bash", "uname", "dirname", "basename", "readlink", "tr"),
		"BENCH_RUN_BINARY",
	}
	cold := runHook(t, repo, missing, hook)
	if cold.code != 0 {
		t.Errorf("H23 unreachable core exit = %d (%q, %q), want 0 — the hint must never block a session", cold.code, cold.stdout, cold.stderr)
	}
	if !strings.Contains(cold.stdout, rebuild) {
		t.Errorf("H23 unreachable core stdout = %q, want the rebuild invocation %q", cold.stdout, rebuild)
	}

	// The hint replaces silence, not the dashboard: a reachable core still routes to
	// `bench session-inspect` and says nothing about rebuilding.
	warm := runHook(t, repo, []string{
		"BENCH_KIT=" + owner.kit,
		"BENCH_RUN_BINARY=" + owner.selected.path,
		"BENCH_HOME=" + filepath.Join(t.TempDir(), "home"),
	}, hook)
	if warm.code != 0 {
		t.Errorf("H23 reachable core exit = %d (%q, %q), want 0", warm.code, warm.stdout, warm.stderr)
	}
	if strings.Contains(warm.stdout, "scripts/go-build.sh") {
		t.Errorf("H23 reachable core printed the rebuild hint: %q", warm.stdout)
	}
	if !strings.Contains(warm.stdout, "bench CLI: ") {
		t.Errorf("H23 reachable core stdout = %q, want the unchanged CLI advertisement", warm.stdout)
	}
}

// TestSessionStartIsSilentOutsideARepository pins H24. The launch deliberately makes a
// wrapper reachable — `bench` resolves on PATH — so the assertion is that the hook
// declines to say anything because the working directory is not a repository, not that
// it had nothing to say. Without that distinction the row would pass on a hook whose
// repository guard had been deleted.
func TestSessionStartIsSilentOutsideARepository(t *testing.T) {
	hook := sessionStartHook(t)
	outside := t.TempDir()
	path := privateToolPath(t, "git", "bash", "uname", "dirname", "basename", "readlink", "tr")
	reachable := filepath.Join(t.TempDir(), "path")
	if err := os.MkdirAll(reachable, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(owner.kit, "bin", "bench.sh"), filepath.Join(reachable, "bench")); err != nil {
		t.Fatal(err)
	}

	result := runHook(t, outside, []string{
		"BENCH_KIT=" + filepath.Join(t.TempDir(), "missing-kit"),
		"BENCH_HOME=" + filepath.Join(t.TempDir(), "home"),
		"PATH=" + path + string(os.PathListSeparator) + reachable,
		"BENCH_RUN_BINARY",
	}, hook)
	if result.code != 0 || result.stdout != "" || result.stderr != "" {
		t.Errorf("H24 outside a repository = (%d, %q, %q), want exit 0 and no output at all", result.code, result.stdout, result.stderr)
	}
}

func sessionStartHook(t *testing.T) string {
	t.Helper()
	hook := filepath.Join(owner.kit, ".bench", "hooks", "session-start.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Fatal(err)
	}
	return hook
}

// runHook launches one hook through the owner's process ledger, the way every other
// system launch in this package does.
func runHook(t *testing.T, dir string, overrides []string, hook string) processResult {
	t.Helper()
	shell, err := exec.LookPath("bash")
	if err != nil {
		t.Fatal(err)
	}
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runAt(dir, overrides, shell, hook)
}

// installWrapper puts the kit's real wrapper where the shared resolver looks for it in
// repo, so both halves of H23 resolve a wrapper regardless of what other journeys in
// this package have installed there.
func installWrapper(t *testing.T, repo string) {
	t.Helper()
	source, err := os.ReadFile(filepath.Join(owner.kit, "bin", "bench.sh"))
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(repo, "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bench.sh"), source, 0o755); err != nil {
		t.Fatal(err)
	}
}
