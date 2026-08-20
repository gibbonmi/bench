//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestDegradedGuardRimDecidesFromTheCommandField drives
// .bench/hooks/block-dangerous-git.sh as a real subprocess in the one state the rim
// exists for: no wrapper reachable on disk and no `bench` on PATH, so the shim cannot
// hand the envelope to internal/gitguard and has to decide alone. The working
// directory is deliberately not a repository and PATH carries only the tools the shim
// itself runs, which is exactly the "required tool missing from PATH" edge class. Every
// launch goes through the owner's process ledger; the assertions read the exit code the
// PreToolUse contract defines (2 blocks, 0 allows) plus the degraded message, which only
// the rim emits, so a verdict that came from a reachable core could not pass as one.
func TestDegradedGuardRimDecidesFromTheCommandField(t *testing.T) {
	hook := filepath.Join(owner.kit, ".bench", "hooks", "block-dangerous-git.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Fatal(err)
	}
	shell, err := exec.LookPath("bash")
	if err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	path := privateToolPath(t, "git", "cat")
	run := func(envelope string) processResult {
		if err := owner.observeSelected(); err != nil {
			t.Fatal(err)
		}
		return owner.runWithInput(outside, []string{"PATH=" + path}, envelope, shell, hook)
	}

	// H20 — the fail-closed half. A destructive git command is refused even though
	// nothing can classify it, including through the one-level wrapper the core
	// documents and behind the honest-mistake prefixes an agent reflexively types.
	for _, blocked := range []struct{ name, envelope string }{
		{"bare destructive verb", `{"tool_name":"Bash","tool_input":{"command":"git reset --hard HEAD~1"}}`},
		{"push", `{"tool_name":"Bash","tool_input":{"command":"cd /srv && git push --force origin main"}}`},
		{"one-level wrapper", `{"tool_name":"Bash","tool_input":{"command":"bash -lc 'git clean -fd'"}}`},
		{"timeout prefix", `{"tool_name":"Bash","tool_input":{"command":"timeout 30 git rebase -i main"}}`},
		{"absolute path to git", `{"tool_name":"Bash","tool_input":{"command":"/usr/bin/git branch -D topic"}}`},
	} {
		result := run(blocked.envelope)
		if result.code != 2 || !strings.Contains(result.stderr, "BLOCKED: guard degraded") {
			t.Errorf("H20 %s = (%d, %q, %q), want exit 2 and the degraded refusal", blocked.name, result.code, result.stdout, result.stderr)
		}
	}

	// H21 — the defect. `git` in a path, in an argument, or in an unrelated envelope
	// field is not a git invocation, and the rim runs during exactly the cold session
	// that needs to read those files.
	for _, allowed := range []struct{ name, envelope string }{
		{"path under .github", `{"session_id":"git-recovery","transcript_path":"/home/u/.claude/git.jsonl","tool_name":"Bash","tool_input":{"command":"cat .github/workflows/gate.yml"}}`},
		{"git as an argument", `{"tool_name":"Bash","tool_input":{"command":"rg git AGENTS.md"}}`},
		{"git only in a sibling field", `{"tool_name":"Bash","tool_input":{"description":"read the git guard","command":"sed -n '1,40p' hooks.md"}}`},
		{"git only in a redirect target", `{"tool_name":"Bash","tool_input":{"command":"printf hi > git"}}`},
		{"empty command", `{"tool_name":"Bash","tool_input":{"command":""}}`},
	} {
		result := run(allowed.envelope)
		if result.code != 0 {
			t.Errorf("H21 %s = (%d, %q, %q), want exit 0", allowed.name, result.code, result.stdout, result.stderr)
		}
	}

	// The rim's wrapper recursion stops where internal/gitguard's does — exactly one
	// level — so this pins the shell copy's depth to the core's rather than letting it
	// drift. It is a depth statement, not a safety claim: the threat model is honest
	// mistakes, and the pre-push hook is the backstop.
	if nested := run(`{"tool_name":"Bash","tool_input":{"command":"bash -c 'sh -c \"git push\"'"}}`); nested.code != 0 {
		t.Errorf("nested wrapper = (%d, %q, %q), want the core's one-level depth (exit 0)", nested.code, nested.stdout, nested.stderr)
	}

	// H22 — the narrowing must not turn an unreadable envelope into an empty command.
	// Each of these carries no command the rim can read, and each must refuse.
	for _, unreadable := range []struct{ name, envelope string }{
		{"empty stdin", ``},
		{"not json", `<!DOCTYPE html>`},
		{"no tool_input", `{"tool_name":"Bash"}`},
		{"no command field", `{"tool_name":"Bash","tool_input":{"description":"read a file"}}`},
		{"non-string command", `{"tool_name":"Bash","tool_input":{"command":["ls"]}}`},
		{"unterminated command string", `{"tool_name":"Bash","tool_input":{"command":"ls -la`},
	} {
		result := run(unreadable.envelope)
		if result.code != 2 || !strings.Contains(result.stderr, "BLOCKED: guard degraded") {
			t.Errorf("H22 %s = (%d, %q, %q), want exit 2 and the degraded refusal", unreadable.name, result.code, result.stdout, result.stderr)
		}
	}
}

// privateToolPath builds a PATH directory holding exactly the named tools, so a launch
// under it reaches the shim's own dependencies and nothing else. It fails when `bench`
// turns out to be reachable through it, because that is the precondition the whole
// degraded-rim assertion rests on.
func privateToolPath(t *testing.T, tools ...string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "path")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, tool := range tools {
		resolved, err := exec.LookPath(tool)
		if err != nil {
			t.Fatalf("%s is required by the shim and is not on PATH: %v", tool, err)
		}
		if err := os.Symlink(resolved, filepath.Join(dir, tool)); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"bench", "bench.sh"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			t.Fatalf("the private tool path reaches %s, so the core would not be missing", name)
		}
	}
	return dir
}
