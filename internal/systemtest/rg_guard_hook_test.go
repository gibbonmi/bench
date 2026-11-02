//go:build system

package systemtest

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/rgguard"
)

// rgGuardHook is the shim under test, read from the owner kit.
func rgGuardHook() string { return filepath.Join(owner.kit, ".bench", "hooks", "block-unsafe-rg.sh") }

// TestRgGuardHookFailsOpen covers R42 and R43. Each degraded core leaves the shell usable:
// a missing library, a missing core, a stale core, and a core error each warn and exit 0.
func TestRgGuardHookFailsOpen(t *testing.T) {
	t.Parallel()
	hook := rgGuardHook()
	stale := plantWrapperAnswer(t, `bench: unknown subcommand: "guard-rg"`, 2)
	broken := plantWrapperAnswer(t, "core failure", 7)
	for _, tc := range []struct {
		row, hook, dir string
		overrides      []string
		warning        string
	}{
		{"R42 missing library", hookWithoutLibrary(t, hook), t.TempDir(), []string{"PATH=" + privateToolPath(t, "git", "cat")}, "wrapper resolver missing"},
		{"R42 missing core", hook, t.TempDir(), []string{"PATH=" + privateToolPath(t, "git", "cat", "dirname")}, "bench core not found"},
		{"R43 stale core", hook, stale, []string{"BENCH_KIT=" + owner.kit}, "the Bench binary is stale, so the rg guard is off"},
		{"core error", hook, broken, []string{"BENCH_KIT=" + owner.kit}, "bench core errored (exit 7)"},
	} {
		result := runGitGuardHook(t, tc.hook, tc.dir, tc.overrides, "rg foo")
		if result.code != 0 || !strings.Contains(result.stderr, "WARNING: block-unsafe-rg: "+tc.warning) {
			t.Errorf("%s = (%d, %q, %q), want exit 0 and the warning %q", tc.row, result.code, result.stdout, result.stderr, tc.warning)
		}
	}
}

// TestRgGuardHookRefuses covers R45 through a real envelope pipe: the shim forwards the
// core's refusal and its exit 2, and a call with a path passes.
func TestRgGuardHookRefuses(t *testing.T) {
	t.Parallel()
	repo := journeyRepo(t, "rg-guard [journey]-")
	installWrapper(t, repo)
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	overrides := []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
	result := runGitGuardHook(t, rgGuardHook(), repo, overrides, "rg foo")
	if want := rgguard.Classify("rg foo").Message(); result.code != 2 || strings.TrimRight(result.stderr, "\n") != want {
		t.Errorf("R45 rg foo = (%d, %q, %q), want exit 2 and %q", result.code, result.stdout, result.stderr, want)
	}
	if result := runGitGuardHook(t, rgGuardHook(), repo, overrides, "rg foo ."); result.code != 0 || result.stderr != "" {
		t.Errorf("rg foo . = (%d, %q, %q), want exit 0 and no stderr", result.code, result.stdout, result.stderr)
	}
}
