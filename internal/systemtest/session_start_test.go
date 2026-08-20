//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
)

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
