package conformance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// claudeHook and claudeHookGroup are the settings.json shapes that the wiring fixtures
// marshal.
type claudeHook struct {
	Type    string   `json:"type"`
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

type claudeHookGroup struct {
	Matcher *string      `json:"matcher,omitempty"`
	Hooks   []claudeHook `json:"hooks"`
}

func claudeMatcher(value string) *string { return &value }

func claudeCommand(name string, args ...string) []claudeHook {
	return []claudeHook{{Type: "command", Command: "$CLAUDE_PROJECT_DIR/.bench/hooks/" + name, Args: args}}
}

func claudeWorktreeCommand(action string) []claudeHook {
	return []claudeHook{{Type: "command", Command: "${CLAUDE_PROJECT_DIR}/.bench/hooks/worktree-lifecycle.sh", Args: []string{action}}}
}

// intactClaudeHooks returns a fresh healthy wiring shape each call. A case can then mutate
// its own copy without disturbing the others.
func intactClaudeHooks() map[string][]claudeHookGroup {
	bash := append(claudeCommand("block-dangerous-git.sh"), claudeCommand("block-bench-follow-on.sh")...)
	return map[string][]claudeHookGroup{
		"SessionStart":   {{Matcher: claudeMatcher(""), Hooks: claudeCommand("session-start.sh")}},
		"Stop":           {{Matcher: claudeMatcher("*"), Hooks: claudeCommand("stop.sh")}},
		"WorktreeCreate": {{Hooks: claudeWorktreeCommand("create")}},
		"WorktreeRemove": {{Hooks: claudeWorktreeCommand("remove")}},
		"PreToolUse": {
			{Matcher: claudeMatcher("Bash"), Hooks: append(bash, claudeCommand("block-unsafe-rg.sh")...)},
			{Matcher: claudeMatcher("Agent"), Hooks: claudeCommand("check-agent-line.sh")},
		},
	}
}

// dropBashHook removes the hook that runs script from the first PreToolUse group, which
// the intact shape makes the Bash group.
func dropBashHook(h map[string][]claudeHookGroup, script string) {
	var kept []claudeHook
	for _, hook := range h["PreToolUse"][0].Hooks {
		if !strings.HasSuffix(hook.Command, "/"+script) {
			kept = append(kept, hook)
		}
	}
	h["PreToolUse"][0].Hooks = kept
}

func writeClaudeHooks(t *testing.T, hooks map[string][]claudeHookGroup) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, ".claude")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(map[string]any{"hooks": hooks})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "settings.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestClaudeHookWiringBites is the recorded bite proof for checkClaudeHookWiring. The
// intact kit-shaped settings.json passes clean: two PreToolUse groups (Bash, Agent) plus
// Stop and SessionStart, each carrying the real $CLAUDE_PROJECT_DIR command prefix.
// Dropping any one wiring fires exactly its diagnostic and no sibling. The absent-file
// case pins the skip-on-absent posture, for parity with checkCodexHooks, as a regression
// guard. It passes on day one and must keep yielding nothing. A later fail-closed change
// must then be a deliberate edit, not drift.
func TestClaudeHookWiringBites(t *testing.T) {
	if diags := checkClaudeHookWiring(writeClaudeHooks(t, intactClaudeHooks())); len(diags) != 0 {
		t.Fatalf("intact kit-shaped settings.json: want no diagnostics, got %v", diags)
	}

	cases := []struct {
		name          string
		remove        func(map[string][]claudeHookGroup)
		want          string
		allowSiblings bool
	}{
		{"stop dropped", func(h map[string][]claudeHookGroup) { delete(h, "Stop") }, "claude settings.json Stop event does not run .bench/hooks/stop.sh", false},
		{"session-start dropped", func(h map[string][]claudeHookGroup) { delete(h, "SessionStart") }, "claude settings.json SessionStart event does not run .bench/hooks/session-start.sh", false},
		{"bash group dropped", func(h map[string][]claudeHookGroup) { h["PreToolUse"] = h["PreToolUse"][1:] }, "claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-dangerous-git.sh", true},
		{"follow-on hook dropped", func(h map[string][]claudeHookGroup) { dropBashHook(h, "block-bench-follow-on.sh") }, "claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-bench-follow-on.sh", false},
		{"agent group dropped", func(h map[string][]claudeHookGroup) { h["PreToolUse"] = h["PreToolUse"][:1] }, "claude settings.json PreToolUse Agent matcher missing or does not run .bench/hooks/check-agent-line.sh", false},
		{"worktree create dropped", func(h map[string][]claudeHookGroup) { delete(h, "WorktreeCreate") }, "claude settings.json WorktreeCreate event must use the brace project-dir placeholder, be matcher-free, and run worktree-lifecycle.sh with args [create]", false},
		{"worktree remove dropped", func(h map[string][]claudeHookGroup) { delete(h, "WorktreeRemove") }, "claude settings.json WorktreeRemove event must use the brace project-dir placeholder, be matcher-free, and run worktree-lifecycle.sh with args [remove]", false},
	}
	for _, event := range []string{"WorktreeCreate", "WorktreeRemove"} {
		hooks := intactClaudeHooks()
		hooks[event][0].Matcher = claudeMatcher("ignored")
		diags := checkClaudeHookWiring(writeClaudeHooks(t, hooks))
		if !containsDiagnostic(diags, "claude settings.json "+event+" event must use the brace project-dir placeholder") {
			t.Fatalf("%s matcher was accepted: %v", event, diags)
		}
	}
	for event, action := range map[string]string{"WorktreeCreate": "create", "WorktreeRemove": "remove"} {
		hooks := intactClaudeHooks()
		hooks[event][0].Hooks[0].Command = "$CLAUDE_PROJECT_DIR/.bench/hooks/worktree-lifecycle.sh"
		diags := checkClaudeHookWiring(writeClaudeHooks(t, hooks))
		if !containsDiagnostic(diags, "claude settings.json "+event+" event must use the brace project-dir placeholder") {
			t.Fatalf("%s bare project-dir placeholder was accepted for %s: %v", event, action, diags)
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hooks := intactClaudeHooks()
			tc.remove(hooks)
			diags := checkClaudeHookWiring(writeClaudeHooks(t, hooks))
			if !containsDiagnostic(diags, tc.want) {
				t.Fatalf("want %q in diagnostics, got %v", tc.want, diags)
			}
			for _, other := range cases {
				if tc.allowSiblings {
					break
				}
				if other.want != tc.want && containsDiagnostic(diags, other.want) {
					t.Fatalf("%s also fired sibling diagnostic %q: %v", tc.name, other.want, diags)
				}
			}
		})
	}

	if diags := checkClaudeHookWiring(t.TempDir()); len(diags) != 0 {
		t.Fatalf("absent settings.json: want no diagnostics (skip posture), got %v", diags)
	}
}

// TestClaudeHookWiringNeedsTheRgGuard is R39. A Bash group that runs every other hook but
// not the rg shim fires the shim's diagnostic alone.
func TestClaudeHookWiringNeedsTheRgGuard(t *testing.T) {
	const want = "claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-unsafe-rg.sh"
	hooks := intactClaudeHooks()
	dropBashHook(hooks, "block-unsafe-rg.sh")
	diags := checkClaudeHookWiring(writeClaudeHooks(t, hooks))
	if len(diags) != 1 || diags[0] != want {
		t.Fatalf("Bash group without the rg shim: diagnostics = %v, want only %q", diags, want)
	}
}
