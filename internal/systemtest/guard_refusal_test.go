//go:build system

package systemtest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// occurrenceCommand runs git after an attached separator and leaves a quote open on its
// next line. A whitespace split joins `hi;git` into one word and so hides the git
// invocation.
const occurrenceCommand = "printf hi;git push --force\necho 'x"

// parseRepairWant is the repair sentence of a malformed-command refusal, written here
// independently of the hook and the core that each emit it.
const parseRepairWant = "Close every quote and escape in the command, then run it again."

// coreMissingMarker is the degraded rim's own refusal text. Only the hook emits it, and
// only when no core answers, so a verdict that carries it did not come from the core.
const coreMissingMarker = "BLOCKED: guard degraded (bench core missing)"

// TestDegradedGuardRefusesMalformedCommands drives the real hook with no core on disk or
// PATH. A command that the rim cannot parse refuses with a parse diagnostic, a malformed
// supported wrapper child refuses the same way, and well-formed git keeps the coarse
// refusal.
func TestDegradedGuardRefusesMalformedCommands(t *testing.T) {
	outside := t.TempDir()
	path := []string{"PATH=" + privateToolPath(t, "git", "cat")}
	// The private PATH holds no Bench executable, and the working directory is in no
	// repository, so no repository-local wrapper exists either.
	if result := owner.runAt(outside, path, "git", "rev-parse", "--show-toplevel"); result.code == 0 {
		t.Fatalf("the missing-core directory is inside the repository %q", result.stdout)
	}
	for _, tc := range []struct{ name, command string }{
		{"GR01 open single quote after an attached separator", occurrenceCommand},
		{"GR02 open double quote after an attached separator", "printf hi;git push --force\necho \"x"},
		{"GR03 trailing escape after an attached separator", "printf hi;git push --force\nprintf x\\"},
		{"GR04 open quote with no git word", "echo 'x"},
		{"GR05 open quote in a supported wrapper child", `bash -c "echo 'x"`},
	} {
		result := runGitGuardHook(t, gitGuardHookPath(), outside, path, tc.command)
		assertParseRefusal(t, tc.name, result)
		if !strings.HasPrefix(result.stderr, coreMissingMarker) {
			t.Errorf("%s stderr = %q, want the degraded rim's refusal", tc.name, result.stderr)
		}
	}
	for _, tc := range []struct{ name, command string }{
		{"GR14 well-formed read-only git", "git status"},
		{"GR31 git after an attached separator", "printf hi;git status"},
	} {
		result := runGitGuardHook(t, gitGuardHookPath(), outside, path, tc.command)
		if result.code != 2 || !strings.HasPrefix(result.stderr, coreMissingMarker) {
			t.Errorf("%s = (%d, %q, %q), want exit 2 and the degraded rim's refusal", tc.name, result.code, result.stdout, result.stderr)
		}
	}
}

// TestGitGuardMissingLibraryRefuses runs a copy of the hook whose adjacent shared library
// is absent. The hook cannot read the envelope, so it refuses.
func TestGitGuardMissingLibraryRefuses(t *testing.T) {
	layout, err := os.MkdirTemp(owner.root, "hook layout [no library]-")
	if err != nil {
		t.Fatal(err)
	}
	hook := filepath.Join(layout, "hooks", "block-dangerous-git.sh")
	data, err := os.ReadFile(gitGuardHookPath())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hook, data, 0o755); err != nil {
		t.Fatal(err)
	}
	result := runGitGuardHook(t, hook, t.TempDir(), []string{"PATH=" + privateToolPath(t, "git", "cat")}, occurrenceCommand)
	if result.code != 2 || !strings.Contains(result.stderr, "BLOCKED: guard degraded (shared library missing)") {
		t.Errorf("GR20 missing library = (%d, %q, %q), want exit 2 and the missing-library refusal", result.code, result.stdout, result.stderr)
	}
}

// TestGitGuardCoreErrorRefuses gives the hook a core that exits 3. The hook keeps its
// analyzer-error refusal and does not read the failure as a parse result or an allow.
func TestGitGuardCoreErrorRefuses(t *testing.T) {
	repo := plantWrapperAnswer(t, "core analyzer failed", 3)
	result := runGitGuardHook(t, gitGuardHookPath(), repo, nil, occurrenceCommand)
	if result.code != 2 {
		t.Errorf("GR21 core error exit = %d (stderr %q), want 2", result.code, result.stderr)
	}
	if !strings.Contains(result.stderr, "BLOCKED: guard analyzer error") {
		t.Errorf("GR22 core error stderr = %q, want the analyzer-error refusal", result.stderr)
	}
}

// TestGitGuardMissingBinaryPreservesRecovery gives the hook a wrapper that exits 127, the
// wrapper's answer when this platform has no installed binary. The degraded rim then
// decides: a malformed command refuses, and a recovery read passes.
func TestGitGuardMissingBinaryPreservesRecovery(t *testing.T) {
	repo := plantWrapperAnswer(t, "no installed binary for this platform", 127)
	assertParseRefusal(t, "GR30 occurrence at the missing-binary rim", runGitGuardHook(t, gitGuardHookPath(), repo, nil, occurrenceCommand))
	if result := runGitGuardHook(t, gitGuardHookPath(), repo, nil, "cat .github/workflows/gate.yml"); result.code != 0 {
		t.Errorf("GR35 recovery read = (%d, %q, %q), want exit 0", result.code, result.stdout, result.stderr)
	}
}

// TestGitGuardHookAndCoreParseRefusal sends the malformed command to the selected binary
// through the unchanged hook and through its guard-git entry. Both refuse with the parse
// diagnostic, and the hook passes the core's refusal through unchanged.
func TestGitGuardHookAndCoreParseRefusal(t *testing.T) {
	repo := owner.repos[0]
	installWrapper(t, repo)
	live := []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
	hooked := runGitGuardHook(t, gitGuardHookPath(), repo, live, occurrenceCommand)
	assertParseRefusal(t, "GR26 live-core hook", hooked)
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	direct := owner.runWithInput(repo, live, bashEnvelope(t, occurrenceCommand), owner.selected.path, "guard-git")
	assertParseRefusal(t, "GR27 guard-git entry", direct)
	if hooked.stderr != direct.stderr {
		t.Errorf("GR26 hook stderr = %q, want the core's refusal %q", hooked.stderr, direct.stderr)
	}
}

func gitGuardHookPath() string {
	return filepath.Join(owner.kit, ".bench", "hooks", "block-dangerous-git.sh")
}

// runGitGuardHook feeds hook one Bash envelope for command through the owner's process
// ledger.
func runGitGuardHook(t *testing.T, hook, dir string, overrides []string, command string) processResult {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runWithInput(dir, overrides, bashEnvelope(t, command), shellPath(t), hook)
}

func bashEnvelope(t *testing.T, command string) string {
	t.Helper()
	data, err := json.Marshal(map[string]any{"tool_name": "Bash", "tool_input": map[string]string{"command": command}})
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// assertParseRefusal requires exit 2 and a refusal that names the parse failure after
// BLOCKED:, gives the repair sentence, and carries none of the destructive-command
// authority advice.
func assertParseRefusal(t *testing.T, name string, result processResult) {
	t.Helper()
	if result.code != 2 {
		t.Errorf("%s exit = %d (stderr %q), want 2", name, result.code, result.stderr)
	}
	if _, refusal, found := strings.Cut(result.stderr, "BLOCKED: "); !found || !strings.Contains(refusal, "cannot parse") {
		t.Errorf("%s stderr = %q, want BLOCKED: and the parse failure", name, result.stderr)
	}
	if !strings.Contains(result.stderr, parseRepairWant) {
		t.Errorf("%s stderr = %q, want the repair sentence", name, result.stderr)
	}
	for _, advice := range []string{"you don't have authority", "Stop and hand back", "history rewrite"} {
		if strings.Contains(result.stderr, advice) {
			t.Errorf("%s stderr = %q, carries the authority advice %q", name, result.stderr, advice)
		}
	}
}
