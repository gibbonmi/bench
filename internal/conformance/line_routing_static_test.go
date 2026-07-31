package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/modelid/modelidtest"
)

func checkLineRouting(root string) []string {
	var diags []string
	diags = append(diags, checkLineBinding(root)...)
	diags = append(diags, checkClaudeHookWiring(root)...)
	diags = append(diags, checkAgentHookBehavior(root)...)
	diags = append(diags, checkAdapterLineGuards(root)...)
	diags = append(diags, checkLineHarnessSurfaces(root)...)
	return diags
}

// checkLineBinding grades the reviewer-owned matrix and the profile's rendering of it from
// the same parse the runtime reads. Only a DECLARED harness owes all three cells, so an
// unadopted harness leaves the matrix complete instead of reddening the gate for a harness
// nobody runs here.
func checkLineBinding(root string) []string {
	path := filepath.Join(root, ".bench", "lines.env")
	if !exists(path) {
		return []string{"lines.env missing: .bench/lines.env is the tier binding enforcement reads"}
	}
	binding := lines.ParseBinding([]byte(readIfExists(path)))
	var diags []string
	for _, key := range binding.ForeignKeys() {
		diags = append(diags, fmt.Sprintf("lines.env key unknown: %s names no harness in the binding matrix (%s)", key, strings.Join(lines.Harnesses, ", ")))
	}
	profile := readIfExists(filepath.Join(root, "projects", "benchkit.md"))
	section, anchored := profileLinesSection(profile)
	if profile != "" && !anchored {
		diags = append(diags, "profile Lines section missing: projects/benchkit.md has no 'Lines' heading rendering the binding matrix")
	}
	declared := false
	for _, harness := range lines.Harnesses {
		if !binding.Declared(harness) {
			continue
		}
		declared = true
		for _, tier := range lines.Tiers {
			key := lines.Key(harness, tier)
			value := binding.Cell(harness, tier)
			switch {
			case value == "":
				diags = append(diags, fmt.Sprintf("lines.env cell unset: %s has no value in .bench/lines.env (%s='')", key, key))
			case lines.CellFault(harness, value) != "":
				diags = append(diags, fmt.Sprintf("lines.env cell malformed: %s='%s' %s", key, value, lines.CellFault(harness, value)))
			case anchored && !strings.Contains(section, value):
				diags = append(diags, fmt.Sprintf("profile Lines prose stale: projects/benchkit.md does not name bound model id '%s' (%s in lines.env)", value, key))
			}
		}
	}
	if !declared {
		// A binding that declares no harness at all is unusable. The diagnostic names the
		// first cell so an operator receives a key to write rather than a category.
		key := lines.Key(lines.Harnesses[0], lines.Tiers[0])
		diags = append(diags, fmt.Sprintf("lines.env cell unset: %s has no value in .bench/lines.env (%s='')", key, key))
	}
	return diags
}

// profileLinesSection returns the body under the profile's `Lines` heading and whether that
// heading exists. The cross-check is anchored here rather than run over the whole file
// because an unanchored search lets any passing mention keep the check green while the table
// that renders the binding for a human rots. The section runs to the next heading of the
// same or shallower depth, so its own subsections stay inside it.
func profileLinesSection(profile string) (string, bool) {
	depth := 0
	var body []string
	for _, line := range strings.Split(profile, "\n") {
		hashes := len(line) - len(strings.TrimLeft(line, "#"))
		title := ""
		if hashes > 0 && strings.HasPrefix(line[hashes:], " ") {
			title = strings.TrimSpace(line[hashes:])
		}
		if depth == 0 {
			if title == "Lines" || strings.HasPrefix(title, "Lines ") {
				depth = hashes
			}
			continue
		}
		if title != "" && hashes <= depth {
			break
		}
		body = append(body, line)
	}
	return strings.Join(body, "\n"), depth > 0
}

// TestLineBindingGradesDeclaredHarnessesOnly pins the declared-versus-known distinction in
// both directions: an absent opencode column is silent while codex and claude are complete,
// and a declared column missing one cell is not.
func TestLineBindingGradesDeclaredHarnessesOnly(t *testing.T) {
	write := func(t *testing.T, content string) string {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}
	complete := "BENCH_CODEX_TOP=gpt-5.4\nBENCH_CODEX_MID=gpt-5.3\nBENCH_CODEX_CHEAP=gpt-5.2\n" +
		"BENCH_CLAUDE_TOP=fable\nBENCH_CLAUDE_MID=opus\nBENCH_CLAUDE_CHEAP=sonnet\n"
	if diags := checkLineBinding(write(t, complete)); len(diags) != 0 {
		t.Fatalf("an unadopted opencode column reddened the gate:\n%s", strings.Join(diags, "\n"))
	}
	if diags := checkLineBinding(write(t, complete+"BENCH_OPENCODE_TOP=openai/gpt-5\n")); !containsDiagnostic(diags, "BENCH_OPENCODE_MID has no value") {
		t.Fatalf("a declared-but-incomplete opencode column was accepted:\n%s", strings.Join(diags, "\n"))
	}
	partialClaude := "BENCH_CODEX_TOP=gpt-5.4\nBENCH_CODEX_MID=gpt-5.3\nBENCH_CODEX_CHEAP=gpt-5.2\n" +
		"BENCH_CLAUDE_TOP=fable\nBENCH_CLAUDE_CHEAP=sonnet\n"
	if diags := checkLineBinding(write(t, partialClaude)); !containsDiagnostic(diags, "BENCH_CLAUDE_MID has no value") {
		t.Fatalf("a declared-but-incomplete claude column was accepted:\n%s", strings.Join(diags, "\n"))
	}
	// opencode's namespace is provider-qualified, and the rule lives on its own cells.
	if diags := checkLineBinding(write(t, complete+"BENCH_OPENCODE_TOP=gpt-5\nBENCH_OPENCODE_MID=gpt-4\nBENCH_OPENCODE_CHEAP=gpt-3\n")); !containsDiagnostic(diags, "is not provider-qualified") {
		t.Fatalf("a bare opencode column was accepted:\n%s", strings.Join(diags, "\n"))
	}
	// The retired schema names no harness in the matrix, so it is reported rather than read.
	retired := checkLineBinding(write(t, "BENCH_TIER_TOP=gpt-5.4\nBENCH_ALIAS_TOP=fable\n"))
	for _, want := range []string{"BENCH_TIER_TOP names no harness", "BENCH_ALIAS_TOP names no harness"} {
		if !containsDiagnostic(retired, want) {
			t.Fatalf("want %q, got:\n%s", want, strings.Join(retired, "\n"))
		}
	}
}

// TestLineBindingCrossChecksEveryCellAgainstTheLinesSection proves the quantifier and the
// anchor in one pass, one mutation per declared cell: a checker that samples a single cell
// passes a one-cell test while the other five rot, and an unanchored substring search over
// the whole profile accepts a cell that survives only in an unrelated paragraph. The fixture
// renders its profile from the same parse the check reads, so the cells have one author here
// rather than a second hand-written list.
func TestLineBindingCrossChecksEveryCellAgainstTheLinesSection(t *testing.T) {
	// Distinct, mutually non-containing tokens: a substring collision would let one cell's
	// rendering satisfy another's cross-check and hide a missed mutation.
	env := "BENCH_CODEX_TOP=alpha-top\nBENCH_CODEX_MID=alpha-mid\nBENCH_CODEX_CHEAP=alpha-cheap\n" +
		"BENCH_CLAUDE_TOP=beta-top\nBENCH_CLAUDE_MID=beta-mid\nBENCH_CLAUDE_CHEAP=beta-cheap\n"
	binding := lines.ParseBinding([]byte(env))
	// profile renders every declared cell inside the Lines section except omit. When
	// elsewhere is set, omit is still named after the section, which is the only shape that
	// separates an anchored check from a whole-file search.
	profile := func(omit string, elsewhere bool) string {
		var b strings.Builder
		b.WriteString("# benchkit\n\n## Lines (model + effort routing)\n\n")
		for _, harness := range lines.Harnesses {
			for _, tier := range lines.Tiers {
				value := binding.Cell(harness, tier)
				if value == "" || value == omit {
					continue
				}
				b.WriteString("- " + harness + " " + tier + ": `" + value + "`\n")
			}
		}
		b.WriteString("\n## Notes for cold sessions\n\nThe routing rubric lives in craft-line.\n")
		if elsewhere {
			b.WriteString("\nHistorical note: this repo once bound `" + omit + "`.\n")
		}
		return b.String()
	}
	write := func(t *testing.T, profile string) string {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(root, "projects"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(env), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "projects", "benchkit.md"), []byte(profile), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}

	if diags := checkLineBinding(write(t, profile("", false))); len(diags) != 0 {
		t.Fatalf("a profile rendering every declared cell got diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	for _, harness := range lines.Harnesses {
		for _, tier := range lines.Tiers {
			value := binding.Cell(harness, tier)
			if value == "" {
				continue
			}
			for _, mutation := range []struct {
				name      string
				elsewhere bool
			}{
				{"absent from the profile", false},
				{"named only outside the Lines section", true},
			} {
				t.Run(lines.Key(harness, tier)+" "+mutation.name, func(t *testing.T) {
					diags := checkLineBinding(write(t, profile(value, mutation.elsewhere)))
					if !containsDiagnostic(diags, "profile Lines prose stale") || !containsDiagnostic(diags, "'"+value+"'") {
						t.Fatalf("%s %s was accepted:\n%s", lines.Key(harness, tier), mutation.name, strings.Join(diags, "\n"))
					}
					for _, other := range lines.Harnesses {
						for _, otherTier := range lines.Tiers {
							sibling := binding.Cell(other, otherTier)
							if sibling == "" || sibling == value {
								continue
							}
							if containsDiagnostic(diags, "'"+sibling+"'") {
								t.Fatalf("withholding %s also fired for %s:\n%s", value, sibling, strings.Join(diags, "\n"))
							}
						}
					}
				})
			}
		}
	}

	// A profile carrying no Lines heading has nothing to anchor to; the check names that
	// rather than reporting six cells stale for one cause.
	noSection := checkLineBinding(write(t, "# benchkit\n\n## Notes for cold sessions\n\nalpha-top beta-top\n"))
	if !containsDiagnostic(noSection, "profile Lines section missing") {
		t.Fatalf("a profile with no Lines heading was accepted:\n%s", strings.Join(noSection, "\n"))
	}
}

func TestLineBindingAcceptsOpaqueSafeModelTokens(t *testing.T) {
	for _, value := range modelidtest.AcceptedTokens {
		t.Run(value, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			content := "BENCH_CODEX_TOP=gpt-5.4\n" +
				"BENCH_CODEX_MID=" + value + "\n" +
				"BENCH_CODEX_CHEAP=openai/gpt-5\n"
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
			content := "BENCH_CODEX_TOP=gpt-5.4\n" +
				"BENCH_CODEX_MID=" + token.Value + "\n" +
				"BENCH_CODEX_CHEAP=openai/gpt-5\n"
			if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			diags := checkLineBinding(root)
			want := "BENCH_CODEX_MID='" + token.Value + "'"
			if token.Value == "" {
				want = "BENCH_CODEX_MID=''"
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
