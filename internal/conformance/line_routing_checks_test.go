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
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Command string `json:"command"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	var diags []string
	for _, w := range []struct {
		event, matcher, needle, diag string
	}{
		{"Stop", "", ".bench/hooks/stop.sh", "claude settings.json Stop event does not run .bench/hooks/stop.sh"},
		{"SessionStart", "", ".bench/hooks/session-start.sh", "claude settings.json SessionStart event does not run .bench/hooks/session-start.sh"},
		{"PreToolUse", "Bash", ".bench/hooks/block-dangerous-git.sh", "claude settings.json PreToolUse Bash matcher missing or does not run .bench/hooks/block-dangerous-git.sh"},
		{"PreToolUse", "Agent", ".bench/hooks/check-agent-line.sh", "claude settings.json PreToolUse Agent matcher missing or does not run .bench/hooks/check-agent-line.sh"},
	} {
		wired := false
		for _, group := range cfg.Hooks[w.event] {
			if w.matcher != "" && group.Matcher != w.matcher {
				continue
			}
			for _, hook := range group.Hooks {
				if strings.Contains(hook.Command, w.needle) {
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
		Type    string `json:"type"`
		Command string `json:"command"`
	}
	type group struct {
		Matcher string `json:"matcher"`
		Hooks   []hook `json:"hooks"`
	}
	command := func(name string) []hook {
		return []hook{{Type: "command", Command: "$CLAUDE_PROJECT_DIR/.bench/hooks/" + name}}
	}
	// intact returns a fresh healthy wiring shape each call so a case can mutate
	// its own copy without disturbing the others.
	intact := func() map[string][]group {
		return map[string][]group{
			"SessionStart": {{Matcher: "", Hooks: command("session-start.sh")}},
			"Stop":         {{Matcher: "*", Hooks: command("stop.sh")}},
			"PreToolUse": {
				{Matcher: "Bash", Hooks: command("block-dangerous-git.sh")},
				{Matcher: "Agent", Hooks: command("check-agent-line.sh")},
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

func checkAgentHookBehavior(root string) []string {
	hook := filepath.Join(root, ".bench", "hooks", "check-agent-line.sh")
	realBench := filepath.Join(root, "bin", "bench.sh")
	if !exists(hook) {
		return nil
	}
	if !exists(realBench) {
		probe := runWithInput(root, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "bash", hook)
		if probe != nil && probe.ExitCode == 0 {
			return []string{"check-agent-line.sh does not deny an unbound model"}
		}
		return nil
	}

	bindir, cleanup, err := wrapperStubDir(realBench)
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanup()
	env := append(conformanceSubprocessEnv(), "PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"))
	var diags []string

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=gpt-5.3-codex-spark\nBENCH_TIER_CHEAP=openai/gpt-5\nBENCH_ALIAS_MID=opus\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupRouted()

	hookCase := func(label, cwd, input, wantErr string, wantExit int) {
		probe := runWithInputEnv(cwd, env, input, "bash", hook)
		if probe == nil || probe.ExitCode != wantExit {
			got := -1
			if probe != nil {
				got = probe.ExitCode
			}
			diags = append(diags, fmt.Sprintf("check-agent-line.sh %s exit %d (want %d)", label, got, wantExit))
			return
		}
		if wantErr != "" && !strings.Contains(probe.Stderr, wantErr) && !strings.Contains(probe.Stdout, wantErr) {
			diags = append(diags, fmt.Sprintf("check-agent-line.sh %s did not warn with %q", label, wantErr))
		}
	}
	hookCase("allows a bound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-5.3-codex-spark"}}`, "", 0)
	hookCase("allows a declared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus"}}`, "", 0)
	hookCase("denies an undeclared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"sonnet"}}`, "", 2)
	hookCase("denies an unbound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-9"}}`, "", 2)
	hookCase("does not fail open on malformed stdin", routed, `not json at all`, "not parseable as JSON", 0)
	// Ratified posture flip (enforcement-verification): in a routed, completely-bound repo a
	// missing model DENIES (exit 2) rather than warning — an omitted model inherits the
	// session's model, the silent-escalation path invariant #2 exists to stop.
	hookCase("denies a missing model field in a routed repo", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "bound alias", 2)

	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	hookCase("does not fail open without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-9"}}`, "no .bench/lines.env", 0)
	// The missing-model deny is gated on routing: with no binding to enforce, a missing
	// model keeps the fail-open rim (the residual the decision deliberately preserves).
	hookCase("does not fail open on a missing model without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "no resolvedModel/model field", 0)

	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=openai/gpt-5\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupPartial()
	hookCase("does not fail open on an incomplete binding", partial, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-9"}}`, "unset or empty", 0)
	return diags
}

func checkAdapterLineGuards(root string) []string {
	if !exists(filepath.Join(root, ".bench", "adapters")) {
		return nil
	}
	var diags []string
	realBench := filepath.Join(root, "bin", "bench.sh")
	bindir := ""
	cleanup := func() {}
	if exists(realBench) {
		var err error
		bindir, cleanup, err = adapterStubDir(realBench)
		if err != nil {
			return []string{"adapter line guard setup failed: " + err.Error()}
		}
		defer cleanup()
	}
	for _, name := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", name)
		if !exists(path) {
			diags = append(diags, "adapter missing from .bench/adapters: "+name)
			continue
		}
		text := readIfExists(path)
		hasResolveCall := regexp.MustCompile(`(?m)^[[:space:]]*model="\$\([^)]*resolve-model`).MatchString(text)
		if !hasResolveCall {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
	}
	if bindir == "" {
		return diags
	}

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=gpt-5.3-codex-spark\nBENCH_TIER_CHEAP=openai/gpt-5\n")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupRouted()
	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=openai/gpt-5\n")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupPartial()

	for _, name := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", name)
		if !exists(path) {
			continue
		}
		envBase := append(conformanceSubprocessEnv(), "PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"))
		bound := runAtEnv(routed, append(envBase, "BENCH_MODEL=gpt-5.3-codex-spark"), "bash", path, "--line probe prompt")
		if bound == nil || bound.ExitCode != 0 || !strings.Contains(bound.Stdout, "gpt-5.3-codex-spark") || !strings.Contains(bound.Stdout, "--line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass BENCH_MODEL and a dash-leading prompt to the harness in a routed repo", name))
		}
		unset := runAtEnv(routed, envBase, "bash", path, "line probe prompt")
		if unset != nil && unset.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
		}
		unbound := runAtEnv(routed, append(envBase, "BENCH_MODEL=gpt-9"), "bash", path, "line probe prompt")
		if unbound != nil && unbound.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
		pass := runAtEnv(unrouted, envBase, "bash", path, "line probe prompt")
		if pass == nil || pass.ExitCode != 0 || !strings.Contains(pass.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass through in an unrouted repo", name))
		}
		explicit := runAtEnv(unrouted, append(envBase, "BENCH_MODEL=gpt-anything-7"), "bash", path, "line probe prompt")
		if explicit == nil || explicit.ExitCode != 0 || !strings.Contains(explicit.Stdout, "gpt-anything-7") || !strings.Contains(explicit.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass an explicit BENCH_MODEL through in an unrouted repo", name))
		}
		partialProbe := runAtEnv(partial, append(envBase, "BENCH_MODEL=gpt-anything-7"), "bash", path, "line probe prompt")
		if partialProbe == nil || partialProbe.ExitCode != 0 || !strings.Contains(partialProbe.Stdout, "gpt-anything-7") || !strings.Contains(partialProbe.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not fall back to passthrough on an incomplete binding", name))
		}
	}
	return diags
}
