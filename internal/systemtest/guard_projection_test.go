//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// projectionWriteHookPath is the file-write hook in the kit under test.
func projectionWriteHookPath() string {
	return filepath.Join(owner.kit, ".bench", "hooks", "block-primary-file-write.sh")
}

// projectionRunHook feeds hook one raw envelope through the owner's process ledger.
func projectionRunHook(t *testing.T, hook, dir string, overrides []string, envelope string) processResult {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runWithInput(dir, overrides, envelope, shellPath(t), hook)
}

// TestGuardProjectionWriteHookDegraded is GP38. The file-write hook without its library or
// its core warns and allows, so Git's fail-closed refusal never reaches a file write.
func TestGuardProjectionWriteHookDegraded(t *testing.T) {
	envelope := `{"cwd":"/repo","tool_input":{"file_path":"README.md"}}`
	path := []string{"PATH=" + privateToolPath(t, "git", "cat", "dirname")}
	for _, tc := range []struct{ name, hook, reason string }{
		{"GP38 missing library", hookWithoutLibrary(t, projectionWriteHookPath()), "wrapper resolver missing"},
		{"GP38 missing core", projectionWriteHookPath(), "bench core not found"},
	} {
		result := projectionRunHook(t, tc.hook, t.TempDir(), path, envelope)
		if result.code != 0 || result.stdout != "" || !strings.HasPrefix(result.stderr, "WARNING: block-primary-file-write: "+tc.reason) {
			t.Errorf("%s = (%d, %q, %q), want exit 0 and the %q warning", tc.name, result.code, result.stdout, result.stderr, tc.reason)
		}
		if strings.Contains(result.stderr, "BLOCKED:") {
			t.Errorf("%s stderr = %q, carries a refusal", tc.name, result.stderr)
		}
	}
}

// TestGuardProjectionNeverExecutes is GP45. Each expansion-shaped command passes through
// the actual Bash hooks, with and without a core, and its sentinel file stays absent.
func TestGuardProjectionNeverExecutes(t *testing.T) {
	repo := owner.repos[0]
	installWrapper(t, repo)
	live := []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
	degraded := []string{"PATH=" + privateToolPath(t, "git", "cat")}
	followOn := filepath.Join(owner.kit, ".bench", "hooks", "block-bench-follow-on.sh")
	entries := []struct {
		name, hook, dir string
		env             []string
	}{
		{"live Git hook", gitGuardHookPath(), repo, live},
		{"missing-core Git hook", gitGuardHookPath(), t.TempDir(), degraded},
		{"live follow-on hook", followOn, repo, live},
	}
	// The absolute touch path still runs under the missing-core entry's private PATH.
	program, err := exec.LookPath("touch")
	if err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(t.TempDir(), "sentinel")
	touch := program + " '" + sentinel + "'"
	for _, command := range []string{
		"echo $(" + touch + ")",
		"echo `" + touch + "`",
		"cat <(" + touch + ")",
		"echo ${X:-$(" + touch + ")}",
		"git status $(" + touch + ")",
		"bench status $(" + touch + ")",
		"bash -c \"echo $(" + touch + ")\"",
	} {
		for _, entry := range entries {
			result := projectionRunHook(t, entry.hook, entry.dir, entry.env, hookEnvelope(t, "Bash", command))
			if result.code != 0 && result.code != 2 {
				t.Errorf("GP45 %s %q exit = %d (stderr %q), want a verdict", entry.name, command, result.code, result.stderr)
			}
			if _, err := os.Lstat(sentinel); !os.IsNotExist(err) {
				t.Fatalf("GP45 %s %q left the sentinel present (stat error %v)", entry.name, command, err)
			}
		}
	}
}
