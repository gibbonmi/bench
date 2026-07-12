package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/modelid"
	"github.com/gibbonmi/bench/internal/modelid/modelidtest"
)

func checkLineRouting(root string) []string {
	var diags []string
	diags = append(diags, checkLineBinding(root)...)
	diags = append(diags, checkClaudeHookWiring(root)...)
	diags = append(diags, checkAgentHookBehavior(root)...)
	diags = append(diags, checkAdapterLineGuards(root)...)
	return diags
}

func checkLineBinding(root string) []string {
	path := filepath.Join(root, ".bench", "lines.env")
	if !exists(path) {
		return []string{"lines.env missing: .bench/lines.env is the tier binding enforcement reads"}
	}
	content := []byte(readIfExists(path))
	binding := lines.ParseBinding(content)
	var diags []string
	tiers := []struct {
		key   string
		value string
	}{
		{"BENCH_TIER_TOP", binding.Top},
		{"BENCH_TIER_MID", binding.Mid},
		{"BENCH_TIER_CHEAP", binding.Cheap},
	}
	for _, tier := range tiers {
		if tier.value == "" {
			diags = append(diags, fmt.Sprintf("lines.env tier unset: %s has no value in .bench/lines.env (%s='')", tier.key, tier.key))
		} else if !modelid.SafeToken(tier.value) {
			diags = append(diags, fmt.Sprintf("lines.env tier malformed: %s='%s' is not a safe model token", tier.key, tier.value))
		}
	}
	aliases := []struct {
		key   string
		value string
	}{
		{"BENCH_ALIAS_TOP", binding.AliasTop},
		{"BENCH_ALIAS_MID", binding.AliasMid},
		{"BENCH_ALIAS_CHEAP", binding.AliasCheap},
	}
	aliasRe := regexp.MustCompile(`^[a-z0-9-]+$`)
	for _, alias := range aliases {
		if !regexp.MustCompile(`(?m)^[ \t]*` + regexp.QuoteMeta(alias.key) + `=`).Match(content) {
			continue
		}
		if !aliasRe.MatchString(alias.value) {
			diags = append(diags, fmt.Sprintf("lines.env alias malformed: %s='%s' is not a bare alias", alias.key, alias.value))
		}
	}

	profile := readIfExists(filepath.Join(root, "projects", "benchkit.md"))
	if profile != "" {
		for _, tier := range tiers {
			if tier.value == "" {
				continue
			}
			if !strings.Contains(profile, tier.value) {
				diags = append(diags, fmt.Sprintf("profile Lines prose stale: projects/benchkit.md does not name bound model id '%s' (%s in lines.env)", tier.value, tier.key))
			}
		}
		for _, alias := range aliases {
			if alias.value == "" {
				continue
			}
			want := alias.key + "=" + alias.value
			if !strings.Contains(profile, want) {
				diags = append(diags, fmt.Sprintf("profile Lines prose stale: projects/benchkit.md does not carry alias declaration %s", want))
			}
		}
	}
	return diags
}

func TestLineBindingAcceptsOpaqueSafeModelTokens(t *testing.T) {
	for _, value := range modelidtest.AcceptedTokens {
		t.Run(value, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			content := "BENCH_TIER_TOP=gpt-5.4\n" +
				"BENCH_TIER_MID=" + value + "\n" +
				"BENCH_TIER_CHEAP=openai/gpt-5\n"
			if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkLineBinding(root); len(diags) != 0 {
				t.Fatalf("safe non-Claude tier token got diagnostics:\n%s", strings.Join(diags, "\n"))
			}
		})
	}
}

func TestLineBindingRejectsUnsafeModelTokens(t *testing.T) {
	for _, token := range modelidtest.RejectedTokens {
		t.Run(token.Name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			content := "BENCH_TIER_TOP=gpt-5.4\n" +
				"BENCH_TIER_MID=" + token.Value + "\n" +
				"BENCH_TIER_CHEAP=openai/gpt-5\n"
			if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			diags := checkLineBinding(root)
			want := "BENCH_TIER_MID='" + token.Value + "'"
			if token.Value == "" {
				want = "BENCH_TIER_MID=''"
			}
			if !containsDiagnostic(diags, want) {
				t.Fatalf("want diagnostic containing %q, got:\n%s", want, strings.Join(diags, "\n"))
			}
			if token.Value != "" && !containsDiagnostic(diags, "is not a safe model token") {
				t.Fatalf("want safe-token diagnostic, got:\n%s", strings.Join(diags, "\n"))
			}
		})
	}
}

// checkClaudeHookWiring holds .claude/settings.json to at least the Codex hook
// standard: one parse feeds every wiring assertion. Absent file skips (parity
// with checkCodexHooks — the kit always ships the file and canary fixtures hide
// dot-dirs, so fail-closed-on-absent would misfire); malformed JSON is the
// JSON-validity check's fact, so unmarshal failure returns nothing. Matching is
// by the .bench/hooks/<name>.sh command substring, the stable token under
// Claude's $CLAUDE_PROJECT_DIR prefix. Stop and SessionStart are event-wide (an
// empty matcher filter accepts any group); PreToolUse Bash and Agent filter by
// matcher.
func checkClaudeHookWiring(root string) []string {
	path := filepath.Join(root, ".claude", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cfg struct {
		Hooks map[string][]struct {
			Matcher *string `json:"matcher"`
			Hooks   []struct {
				Type, Command string
				Args          []string
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	var diags []string
	for _, w := range []struct {
		event, matcher, needle, action, diag string
		matcherFree                          bool
		exactCommand                         bool
	}{
		{event: "Stop", needle: ".bench/hooks/stop.sh", diag: "claude settings.json Stop event does not run .bench/hooks/stop.sh"},
		{event: "SessionStart", needle: ".bench/hooks/session-start.sh", diag: "claude settings.json SessionStart event does not run .bench/hooks/session-start.sh"},
		{event: "PreToolUse", matcher: "Bash", needle: ".bench/hooks/block-dangerous-git.sh", diag: "claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-dangerous-git.sh"},
		{event: "PreToolUse", matcher: "Agent", needle: ".bench/hooks/check-agent-line.sh", diag: "claude settings.json PreToolUse Agent matcher missing or does not run .bench/hooks/check-agent-line.sh"},
		{event: "WorktreeCreate", needle: "${CLAUDE_PROJECT_DIR}/.bench/hooks/worktree-lifecycle.sh", action: "create", diag: "claude settings.json WorktreeCreate event must use the brace project-dir placeholder, be matcher-free, and run worktree-lifecycle.sh with args [create]", matcherFree: true, exactCommand: true},
		{event: "WorktreeRemove", needle: "${CLAUDE_PROJECT_DIR}/.bench/hooks/worktree-lifecycle.sh", action: "remove", diag: "claude settings.json WorktreeRemove event must use the brace project-dir placeholder, be matcher-free, and run worktree-lifecycle.sh with args [remove]", matcherFree: true, exactCommand: true},
	} {
		wired := false
		for _, group := range cfg.Hooks[w.event] {
			if w.matcherFree && group.Matcher != nil {
				continue
			}
			if w.matcher != "" && (group.Matcher == nil || *group.Matcher != w.matcher) {
				continue
			}
			for _, hook := range group.Hooks {
				argsMatch := w.action == "" || (len(hook.Args) == 1 && hook.Args[0] == w.action)
				commandMatch := strings.Contains(hook.Command, w.needle)
				if w.exactCommand {
					commandMatch = hook.Command == w.needle
				}
				if hook.Type == "command" && commandMatch && argsMatch {
					wired = true
				}
			}
		}
		if !wired {
			diags = append(diags, w.diag)
		}
	}
	return diags
}

// TestClaudeHookWiringBites is the recorded bite proof for checkClaudeHookWiring
// (per craft-gate): the intact kit-shaped settings.json — two PreToolUse groups
// (Bash, Agent) plus Stop and SessionStart, each carrying the real
// $CLAUDE_PROJECT_DIR command prefix — passes clean; dropping any one wiring fires
// exactly its diagnostic and no sibling. The absent-file case pins the skip-on-
// absent posture (parity with checkCodexHooks) as a regression guard: it passes on
// day one and must keep yielding nothing so a later fail-closed change is a
// deliberate edit, not drift.
func TestClaudeHookWiringBites(t *testing.T) {
	type hook struct {
		Type    string   `json:"type"`
		Command string   `json:"command"`
		Args    []string `json:"args,omitempty"`
	}
	type group struct {
		Matcher *string `json:"matcher,omitempty"`
		Hooks   []hook  `json:"hooks"`
	}
	matcher := func(value string) *string { return &value }
	command := func(name string, args ...string) []hook {
		return []hook{{Type: "command", Command: "$CLAUDE_PROJECT_DIR/.bench/hooks/" + name, Args: args}}
	}
	worktreeCommand := func(action string) []hook {
		return []hook{{Type: "command", Command: "${CLAUDE_PROJECT_DIR}/.bench/hooks/worktree-lifecycle.sh", Args: []string{action}}}
	}
	// intact returns a fresh healthy wiring shape each call so a case can mutate
	// its own copy without disturbing the others.
	intact := func() map[string][]group {
		return map[string][]group{
			"SessionStart":   {{Matcher: matcher(""), Hooks: command("session-start.sh")}},
			"Stop":           {{Matcher: matcher("*"), Hooks: command("stop.sh")}},
			"WorktreeCreate": {{Hooks: worktreeCommand("create")}},
			"WorktreeRemove": {{Hooks: worktreeCommand("remove")}},
			"PreToolUse": {
				{Matcher: matcher("Bash"), Hooks: command("block-dangerous-git.sh")},
				{Matcher: matcher("Agent"), Hooks: command("check-agent-line.sh")},
			},
		}
	}
	write := func(t *testing.T, hooks map[string][]group) string {
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

	if diags := checkClaudeHookWiring(write(t, intact())); len(diags) != 0 {
		t.Fatalf("intact kit-shaped settings.json: want no diagnostics, got %v", diags)
	}

	cases := []struct {
		name   string
		remove func(map[string][]group)
		want   string
	}{
		{"stop dropped", func(h map[string][]group) { delete(h, "Stop") }, "claude settings.json Stop event does not run .bench/hooks/stop.sh"},
		{"session-start dropped", func(h map[string][]group) { delete(h, "SessionStart") }, "claude settings.json SessionStart event does not run .bench/hooks/session-start.sh"},
		{"bash group dropped", func(h map[string][]group) { h["PreToolUse"] = h["PreToolUse"][1:] }, "claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-dangerous-git.sh"},
		{"agent group dropped", func(h map[string][]group) { h["PreToolUse"] = h["PreToolUse"][:1] }, "claude settings.json PreToolUse Agent matcher missing or does not run .bench/hooks/check-agent-line.sh"},
		{"worktree create dropped", func(h map[string][]group) { delete(h, "WorktreeCreate") }, "claude settings.json WorktreeCreate event must use the brace project-dir placeholder, be matcher-free, and run worktree-lifecycle.sh with args [create]"},
		{"worktree remove dropped", func(h map[string][]group) { delete(h, "WorktreeRemove") }, "claude settings.json WorktreeRemove event must use the brace project-dir placeholder, be matcher-free, and run worktree-lifecycle.sh with args [remove]"},
	}
	for _, event := range []string{"WorktreeCreate", "WorktreeRemove"} {
		hooks := intact()
		hooks[event][0].Matcher = matcher("ignored")
		diags := checkClaudeHookWiring(write(t, hooks))
		if !containsDiagnostic(diags, "claude settings.json "+event+" event must use the brace project-dir placeholder") {
			t.Fatalf("%s matcher was accepted: %v", event, diags)
		}
	}
	for event, action := range map[string]string{"WorktreeCreate": "create", "WorktreeRemove": "remove"} {
		hooks := intact()
		hooks[event][0].Hooks[0].Command = "$CLAUDE_PROJECT_DIR/.bench/hooks/worktree-lifecycle.sh"
		diags := checkClaudeHookWiring(write(t, hooks))
		if !containsDiagnostic(diags, "claude settings.json "+event+" event must use the brace project-dir placeholder") {
			t.Fatalf("%s bare project-dir placeholder was accepted for %s: %v", event, action, diags)
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hooks := intact()
			tc.remove(hooks)
			diags := checkClaudeHookWiring(write(t, hooks))
			if !containsDiagnostic(diags, tc.want) {
				t.Fatalf("want %q in diagnostics, got %v", tc.want, diags)
			}
			for _, other := range cases {
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
