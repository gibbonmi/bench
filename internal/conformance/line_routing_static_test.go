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
	diags = append(diags, checkGuidanceTokens(root)...)
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
	rendered, tabled := profileLinesCells(section)
	if anchored && !tabled {
		diags = append(diags, "profile Lines table missing: projects/benchkit.md renders no '| tier | <harness> |' matrix table under its 'Lines' heading, so no cell's placement can be checked")
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
			case anchored && tabled && rendered[tier][harness] != value:
				// Reached only when the token IS in the section, so this is placement and
				// nothing else: the table binds it to some other harness or tier.
				diags = append(diags, fmt.Sprintf("profile Lines cell misbound: projects/benchkit.md renders the %s %s cell as '%s', but %s='%s' in lines.env", harness, tier, rendered[tier][harness], key, value))
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

// profileLinesCells parses the profile's binding matrix table into cells[tier][harness],
// reporting whether such a table was found. Placement has to be read rather than membership:
// a matrix whose columns are swapped still carries every bound token somewhere in the Lines
// section, so a substring search over the section stays green while the table tells a reader
// the wrong binding for every cell it names. The table is found by its `tier` corner cell,
// its header row names the harness columns, and each row's first cell names the tier — the
// check reads the rendering the way the human it exists for does.
func profileLinesCells(section string) (map[string]map[string]string, bool) {
	var harnesses []string
	cells := map[string]map[string]string{}
	for _, line := range strings.Split(section, "\n") {
		row, ok := markdownRow(line)
		if !ok {
			continue
		}
		if harnesses == nil {
			// Rows before the matrix belong to some other table; only the one cornered
			// `tier` renders the binding.
			if strings.ToLower(row[0]) == "tier" {
				harnesses = make([]string, 0, len(row)-1)
				for _, name := range row[1:] {
					harnesses = append(harnesses, strings.ToLower(name))
				}
			}
			continue
		}
		tier := strings.ToLower(row[0])
		if tier == "" || isRuleRow(row) {
			continue
		}
		if cells[tier] == nil {
			cells[tier] = map[string]string{}
		}
		for i, value := range row[1:] {
			if i < len(harnesses) {
				cells[tier][harnesses[i]] = value
			}
		}
	}
	return cells, harnesses != nil
}

// markdownRow splits one pipe-delimited table row into its cells, stripping the code-span
// backticks the profile renders a model id inside so the cell compares as the bare token.
func markdownRow(line string) ([]string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") {
		return nil, false
	}
	trimmed = strings.TrimSuffix(strings.TrimPrefix(trimmed, "|"), "|")
	cells := strings.Split(trimmed, "|")
	for i, cell := range cells {
		cell = strings.TrimSpace(cell)
		cell = strings.TrimSuffix(strings.TrimPrefix(cell, "`"), "`")
		cells[i] = strings.TrimSpace(cell)
	}
	return cells, true
}

// isRuleRow reports whether row is the `|---|---|` separator under a header rather than data.
func isRuleRow(row []string) bool {
	for _, cell := range row {
		if cell == "" || strings.Trim(cell, "-:") != "" {
			return false
		}
	}
	return true
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

// linesMatrixFixtureEnv binds a full codex and claude matrix in distinct, mutually
// non-containing tokens: a substring collision would let one cell's rendering satisfy
// another's cross-check and hide a missed mutation. The tier name inside each token is what
// makes a swapped rendering readable in a failure message.
const linesMatrixFixtureEnv = "BENCH_CODEX_TOP=alpha-top\nBENCH_CODEX_MID=alpha-mid\n" +
	"BENCH_CODEX_CHEAP=alpha-cheap\nBENCH_CLAUDE_TOP=beta-top\nBENCH_CLAUDE_MID=beta-mid\n" +
	"BENCH_CLAUDE_CHEAP=beta-cheap\n"

// renderLinesProfile builds a profile whose Lines section carries the binding matrix as the
// markdown table the real profile renders it with, taking each cell's text from cell and
// appending trailer after the section. The table shape is load-bearing, not decoration: the
// cross-check reads a tier per row and a harness per column, so a fixture that listed the
// same tokens some other way would grade a rendering the profile does not have.
func renderLinesProfile(cell func(harness, tier string) string, trailer string) string {
	var b strings.Builder
	b.WriteString("# benchkit\n\n## Lines (model + effort routing)\n\n| tier")
	for _, harness := range lines.Harnesses {
		b.WriteString(" | " + harness)
	}
	b.WriteString(" |\n|---")
	for range lines.Harnesses {
		b.WriteString("|---")
	}
	b.WriteString("|\n")
	for _, tier := range lines.Tiers {
		b.WriteString("| " + tier)
		for _, harness := range lines.Harnesses {
			b.WriteString(" | " + cell(harness, tier))
		}
		b.WriteString(" |\n")
	}
	b.WriteString("\n## Notes for cold sessions\n\nThe routing rubric lives in craft-line.\n")
	b.WriteString(trailer)
	return b.String()
}

// writeLinesRoot plants a binding and a profile in a throwaway root the check can read.
func writeLinesRoot(t *testing.T, env, profile string) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{".bench", "projects"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(env), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "projects", "benchkit.md"), []byte(profile), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestLineBindingCatchesSwappedProfileCells is the placement half of the cross-check, and the
// mutation the per-cell absence rows structurally cannot make: a swap moves no token out of
// the Lines section, so every membership test still passes while the table tells a reader
// the wrong binding. One case per unordered pair of declared cells covers both axes — two
// harnesses at one tier and two tiers in one harness — because a check that reads only the
// row or only the column satisfies one axis while the other rots, which is the same
// single-sample failure the per-cell rows exist to kill.
func TestLineBindingCatchesSwappedProfileCells(t *testing.T) {
	binding := lines.ParseBinding([]byte(linesMatrixFixtureEnv))
	type ref struct{ harness, tier string }
	var declared []ref
	for _, harness := range lines.Harnesses {
		for _, tier := range lines.Tiers {
			if binding.Cell(harness, tier) != "" {
				declared = append(declared, ref{harness, tier})
			}
		}
	}
	for i, a := range declared {
		for _, b := range declared[i+1:] {
			a, b := a, b
			t.Run(lines.Key(a.harness, a.tier)+" swapped with "+lines.Key(b.harness, b.tier), func(t *testing.T) {
				profile := renderLinesProfile(func(harness, tier string) string {
					switch (ref{harness, tier}) {
					case a:
						harness, tier = b.harness, b.tier
					case b:
						harness, tier = a.harness, a.tier
					}
					if value := binding.Cell(harness, tier); value != "" {
						return "`" + value + "`"
					}
					return "unbound"
				}, "")
				diags := checkLineBinding(writeLinesRoot(t, linesMatrixFixtureEnv, profile))
				// A stale diagnostic would mean a token left the section, which would make
				// the swap provable by membership alone and this case prove nothing.
				if containsDiagnostic(diags, "profile Lines prose stale") {
					t.Fatalf("the swap moved a token out of the section, so placement was not what bit:\n%s", strings.Join(diags, "\n"))
				}
				named := false
				for _, cell := range []ref{a, b} {
					if containsDiagnostic(diags, "profile Lines cell misbound") &&
						containsDiagnostic(diags, lines.Key(cell.harness, cell.tier)) {
						named = true
					}
				}
				if !named {
					t.Fatalf("swapping %s with %s was accepted:\n%s",
						lines.Key(a.harness, a.tier), lines.Key(b.harness, b.tier), strings.Join(diags, "\n"))
				}
			})
		}
	}
}

// TestLineBindingRequiresTheProfileToRenderTheMatrixAsATable pins the placement check's own
// precondition. Without it, deleting the table and scattering the six ids through the
// section's prose would silently retire the placement arm while every membership test
// stayed green — the check would lose its teeth with nothing turning red.
func TestLineBindingRequiresTheProfileToRenderTheMatrixAsATable(t *testing.T) {
	binding := lines.ParseBinding([]byte(linesMatrixFixtureEnv))
	var b strings.Builder
	b.WriteString("# benchkit\n\n## Lines (model + effort routing)\n\n")
	for _, harness := range lines.Harnesses {
		for _, tier := range lines.Tiers {
			if value := binding.Cell(harness, tier); value != "" {
				b.WriteString("- " + harness + " " + tier + ": `" + value + "`\n")
			}
		}
	}
	diags := checkLineBinding(writeLinesRoot(t, linesMatrixFixtureEnv, b.String()))
	if !containsDiagnostic(diags, "profile Lines table missing") {
		t.Fatalf("a Lines section rendering the matrix as prose was accepted:\n%s", strings.Join(diags, "\n"))
	}
}

// TestLineBindingCrossChecksEveryCellAgainstTheLinesSection proves the quantifier and the
// anchor in one pass, one mutation per declared cell: a checker that samples a single cell
// passes a one-cell test while the other five rot, and an unanchored substring search over
// the whole profile accepts a cell that survives only in an unrelated paragraph. The fixture
// renders its profile from the same parse the check reads, so the cells have one author here
// rather than a second hand-written list.
func TestLineBindingCrossChecksEveryCellAgainstTheLinesSection(t *testing.T) {
	binding := lines.ParseBinding([]byte(linesMatrixFixtureEnv))
	// profile renders every declared cell inside the Lines section except omit, whose cell
	// reads `unbound` the way the real profile spells an unadopted column. When elsewhere is
	// set, omit is still named after the section, which is the only shape that separates an
	// anchored check from a whole-file search.
	profile := func(omit string, elsewhere bool) string {
		trailer := ""
		if elsewhere {
			trailer = "\nHistorical note: this repo once bound `" + omit + "`.\n"
		}
		return renderLinesProfile(func(harness, tier string) string {
			if value := binding.Cell(harness, tier); value != "" && value != omit {
				return "`" + value + "`"
			}
			return "unbound"
		}, trailer)
	}
	write := func(t *testing.T, profile string) string {
		t.Helper()
		return writeLinesRoot(t, linesMatrixFixtureEnv, profile)
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
