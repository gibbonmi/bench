package lines

import (
	"strings"
	"testing"
)

// pinContent is the retiring gate pin a0 fixture: every load-bearing TierValue case in one
// file. The final line has NO trailing newline and must read intact.
const pinContent = "BENCH_CODEX_TOP=\"claude-quoted-1\"\n" +
	"BENCH_CODEX_MID='claude-squoted-2'\n" +
	"BENCH_CODEX_CHEAP=claude-crlf-3\r\n" +
	"BENCH_CLAUDE_TOP=claude-trail-4   \n" +
	"   BENCH_CLAUDE_MID=claude-indent-5\n" +
	"BENCH_CLAUDE_CHEAP=claude-first-6\n" +
	"BENCH_CLAUDE_CHEAP=claude-last-7\n" +
	"BENCH_EMPTY=\n" +
	"BENCH_NONEWLINE=claude-nonl-8"

func TestTierValue(t *testing.T) {
	tests := []struct {
		name, key, want string
	}{
		{"double-quoted", "BENCH_CODEX_TOP", "claude-quoted-1"},
		{"single-quoted", "BENCH_CODEX_MID", "claude-squoted-2"},
		{"trailing-cr", "BENCH_CODEX_CHEAP", "claude-crlf-3"},
		{"trailing-spaces", "BENCH_CLAUDE_TOP", "claude-trail-4"},
		{"indented-key", "BENCH_CLAUDE_MID", "claude-indent-5"},
		{"last-wins", "BENCH_CLAUDE_CHEAP", "claude-last-7"},
		{"empty-value", "BENCH_EMPTY", ""},
		{"final-line-no-newline", "BENCH_NONEWLINE", "claude-nonl-8"},
		{"absent-key", "BENCH_ABSENT", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TierValue(tt.key, []byte(pinContent)); got != tt.want {
				t.Errorf("TierValue(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

// TestTierValueEdgeCasesReachTheMatrix drives the same shell-compatible parsing through
// ParseBinding, so a schema migration that dropped quoted, CRLF, last-wins, or
// no-final-newline handling fails at the cell a caller actually reads.
func TestTierValueEdgeCasesReachTheMatrix(t *testing.T) {
	b := ParseBinding([]byte(pinContent))
	for _, tt := range []struct {
		harness, tier, want string
	}{
		{"codex", "top", "claude-quoted-1"},
		{"codex", "mid", "claude-squoted-2"},
		{"codex", "cheap", "claude-crlf-3"},
		{"claude", "top", "claude-trail-4"},
		{"claude", "mid", "claude-indent-5"},
		{"claude", "cheap", "claude-last-7"},
	} {
		if got := b.Cell(tt.harness, tt.tier); got != tt.want {
			t.Errorf("Cell(%q, %q) = %q, want %q", tt.harness, tt.tier, got, tt.want)
		}
	}
	if b.Declared("opencode") {
		t.Error("opencode is declared by a binding that names no opencode key")
	}
}

func TestTierValueEdgeCases(t *testing.T) {
	tests := []struct {
		name, key, content, want string
	}{
		// The '=' anchors the key; KEY must not match KEYX.
		{"prefix-not-matched", "BENCH_CODEX_TOP", "BENCH_CODEX_TOPX=nope\n", ""},
		// Value is everything after the FIRST '='.
		{"equals-in-value", "K", "K=a=b=c\n", "a=b=c"},
		// [[:space:]] includes CR, so a trailing run of CRs is trimmed like any
		// whitespace — matching the shell, whose trailing-whitespace strip removes
		// everything after the last non-space char.
		{"double-cr", "K", "K=v\r\r\n", "v"},
		// A CR followed by a space: the whole trailing whitespace run goes.
		{"cr-then-space", "K", "K=x\r \n", "x"},
		// Quote stripping is independent, not paired: one leading OR one trailing.
		{"unbalanced-leading-dquote", "K", "K=\"v\n", "v"},
		{"unbalanced-trailing-dquote", "K", "K=v\"\n", "v"},
		// Double-quote pair stripped BEFORE single-quote pair.
		{"nested-quotes", "K", "K=\"'v'\"\n", "v"},
		// A single unmatched quote char.
		{"only-single-quote-left", "K", "K='v\n", "v"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TierValue(tt.key, []byte(tt.content)); got != tt.want {
				t.Errorf("TierValue(%q, %q) = %q, want %q", tt.key, tt.content, got, tt.want)
			}
		})
	}
}

func TestModelFromEnvelope(t *testing.T) {
	tests := []struct {
		name, data, want string
		wantErr          bool
	}{
		{"resolved-wins", `{"tool_input":{"resolvedModel":"r","model":"m"}}`, "r", false},
		{"model-fallback", `{"tool_input":{"model":"m"}}`, "m", false},
		{"empty-resolved-falls-to-model", `{"tool_input":{"resolvedModel":"","model":"m"}}`, "m", false},
		{"both-absent", `{"tool_input":{}}`, "", false},
		{"no-tool-input", `{}`, "", false},
		// A whitespace-only field is blank for the non-empty test: a whitespace
		// resolvedModel falls back to model, and an all-whitespace envelope yields "".
		{"whitespace-resolved-falls-to-model", `{"tool_input":{"resolvedModel":"   ","model":"m"}}`, "m", false},
		{"all-whitespace-yields-blank", `{"tool_input":{"resolvedModel":"  ","model":"\t"}}`, "", false},
		{"not-json", `not json`, "", true},
		{"empty-input", ``, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ModelFromEnvelope([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("ModelFromEnvelope(%q) err = %v, wantErr %v", tt.data, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ModelFromEnvelope(%q) = %q, want %q", tt.data, got, tt.want)
			}
		})
	}
}

// fullBinding binds codex and claude and leaves opencode unadopted. Its claude cells are
// deliberately NOT this repo's own (`fable`/`opus`/`sonnet`), so a renderer that hard-codes
// the kit's live binding fails every message assertion below.
const fullBinding = "BENCH_CODEX_TOP=gpt-5.6-sol\n" +
	"BENCH_CODEX_MID=gpt-5.6-terra\n" +
	"BENCH_CODEX_CHEAP=gpt-5.6-luna\n" +
	"BENCH_CLAUDE_TOP=fable-5\n" +
	"BENCH_CLAUDE_MID=opus-4-8\n" +
	"BENCH_CLAUDE_CHEAP=sonnet-5\n"

// retiredBinding carries only the retired BENCH_TIER_* / BENCH_ALIAS_* schema. The cut is
// hard: it binds nothing.
const retiredBinding = "BENCH_TIER_TOP=gpt-5.6-sol\n" +
	"BENCH_TIER_MID=gpt-5.6-terra\n" +
	"BENCH_TIER_CHEAP=gpt-5.6-luna\n" +
	"BENCH_ALIAS_TOP=fable\n" +
	"BENCH_ALIAS_MID=opus\n" +
	"BENCH_ALIAS_CHEAP=sonnet\n"

func bound(content string) Source {
	return Source{Path: ".bench/lines.env", Exists: true, Content: []byte(content)}
}

func envelope(model string) []byte {
	return []byte(`{"tool_input":{"model":"` + model + `"}}`)
}

// agentEnvelope is the captured Agent PreToolUse shape: tool_input carries description,
// prompt, subagent_type, and an optional model. subagentType and model are raw JSON, so a
// case can drive a non-string subagent_type or omit either field entirely.
func agentEnvelope(subagentType, model string) []byte {
	fields := `"description":"d","prompt":"p"`
	if subagentType != "" {
		fields += `,"subagent_type":` + subagentType
	}
	if model != "" {
		fields += `,"model":` + model
	}
	return []byte(`{"tool_name":"Agent","tool_input":{` + fields + `}}`)
}

func forkEnvelope(model string) []byte {
	return agentEnvelope(`"fork"`, model)
}

func TestParseBinding(t *testing.T) {
	b := ParseBinding([]byte(fullBinding))
	for _, tt := range []struct {
		harness, tier, want string
	}{
		{"codex", "top", "gpt-5.6-sol"},
		{"codex", "mid", "gpt-5.6-terra"},
		{"codex", "cheap", "gpt-5.6-luna"},
		{"claude", "top", "fable-5"},
		{"claude", "mid", "opus-4-8"},
		{"claude", "cheap", "sonnet-5"},
		{"opencode", "top", ""},
	} {
		if got := b.Cell(tt.harness, tt.tier); got != tt.want {
			t.Errorf("Cell(%q, %q) = %q, want %q", tt.harness, tt.tier, got, tt.want)
		}
	}
	for _, harness := range []string{"codex", "claude"} {
		if !b.Declared(harness) || !b.Complete(harness) {
			t.Errorf("%s column: declared=%v complete=%v, want both true", harness, b.Declared(harness), b.Complete(harness))
		}
	}
	if b.Declared("opencode") || b.Complete("opencode") {
		t.Error("an absent opencode column reads as declared or complete")
	}
	if keys := b.ForeignKeys(); len(keys) != 0 {
		t.Errorf("ForeignKeys = %v, want none", keys)
	}
}

// TestParseBindingRejectsForeignHarnessKeys pins the closed harness set: a
// BENCH_<HARNESS>_<TIER>-shaped key naming a harness outside it creates no column and is
// reported, and the retired schema is exactly such a key — so a lines.env carrying only
// retired keys binds nothing at all rather than reading as a legacy binding.
func TestParseBindingRejectsForeignHarnessKeys(t *testing.T) {
	for _, tt := range []struct {
		name, content string
		want          []string
	}{
		{"unknown harness", "BENCH_CODEX_TOP=a\nBENCH_CODEX_MID=b\nBENCH_CODEX_CHEAP=c\nBENCH_GEMINI_TOP=g\n", []string{"BENCH_GEMINI_TOP"}},
		{"retired schema", retiredBinding, []string{"BENCH_TIER_TOP", "BENCH_TIER_MID", "BENCH_TIER_CHEAP", "BENCH_ALIAS_TOP", "BENCH_ALIAS_MID", "BENCH_ALIAS_CHEAP"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b := ParseBinding([]byte(tt.content))
			if got := strings.Join(b.ForeignKeys(), ","); got != strings.Join(tt.want, ",") {
				t.Errorf("ForeignKeys = %q, want %q", got, strings.Join(tt.want, ","))
			}
		})
	}
	if b := ParseBinding([]byte(retiredBinding)); !b.Empty() {
		t.Error("a retired-schema binding is not empty: the hard cut kept a legacy column alive")
	}
}

// TestRetiredSchemaBindsNothing drives the hard cut through both verdicts: a lines.env
// carrying only retired keys must resolve nothing and enforce nothing, exactly as if it
// bound no cell.
func TestRetiredSchemaBindsNothing(t *testing.T) {
	src := bound(retiredBinding)
	if model, exit, stderr := ResolveModelVerdict("codex", "top", true, src); model != "" || exit != 1 {
		t.Errorf("ResolveModelVerdict on retired keys = (%q, %d, %q), want no model and exit 1", model, exit, stderr)
	}
	// The retired ids and aliases are no longer bound tokens, so enforcement is not run
	// against them; the guard keeps its fail-open rim instead of denying on a stale column.
	for _, model := range []string{"gpt-5.6-sol", "fable"} {
		exit, stderr := AgentLineVerdict(envelope(model), "claude", src)
		if exit != 0 || !contains(stderr, "unknown harness key") {
			t.Errorf("AgentLineVerdict(%q) on retired keys = (%d, %q), want fail-open with the unknown-key warning", model, exit, stderr)
		}
	}
}

// TestUnreadableSourceIsNotAbsent pins the distinction the locator used to fold away: an
// absent binding is an unrouted repo, an unreadable one is a corrupt oracle.
func TestUnreadableSourceIsNotAbsent(t *testing.T) {
	absent := Source{Path: ".bench/lines.env"}
	unreadable := Source{Path: ".bench/lines.env", Exists: true, Unreadable: true}

	if model, exit, stderr := ResolveModelVerdict("codex", "gpt-9", true, absent); model != "gpt-9" || exit != 0 || stderr != "" {
		t.Errorf("absent binding = (%q, %d, %q), want unrouted passthrough", model, exit, stderr)
	}
	model, exit, stderr := ResolveModelVerdict("codex", "top", true, unreadable)
	if model != "" || exit != 1 || !contains(stderr, "cannot read the tier binding") {
		t.Errorf("unreadable binding = (%q, %d, %q), want a fail-closed read error", model, exit, stderr)
	}

	if exit, stderr := AgentLineVerdict(envelope("gpt-9"), "claude", absent); exit != 0 || !contains(stderr, "no .bench/lines.env") {
		t.Errorf("absent binding guard = (%d, %q), want the unrouted warning", exit, stderr)
	}
	if exit, stderr := AgentLineVerdict(envelope("gpt-9"), "claude", unreadable); exit != 0 || !contains(stderr, "unreadable") {
		t.Errorf("unreadable binding guard = (%d, %q), want the unreadable warning", exit, stderr)
	}
}

// TestResolveModelVerdictResolvesEveryCell is the six-cell table: each harness resolves a
// tier through its OWN column, so neither a one-tier nor a one-harness resolver survives.
func TestResolveModelVerdictResolvesEveryCell(t *testing.T) {
	for _, tt := range []struct {
		harness, tier, want string
	}{
		{"codex", "top", "gpt-5.6-sol"},
		{"codex", "mid", "gpt-5.6-terra"},
		{"codex", "cheap", "gpt-5.6-luna"},
		{"claude", "top", "fable-5"},
		{"claude", "mid", "opus-4-8"},
		{"claude", "cheap", "sonnet-5"},
	} {
		t.Run(tt.harness+"-"+tt.tier, func(t *testing.T) {
			model, exit, stderr := ResolveModelVerdict(tt.harness, tt.tier, true, bound(fullBinding))
			if model != tt.want || exit != 0 || stderr != "" {
				t.Errorf("ResolveModelVerdict(%q, %q) = (%q, %d, %q), want %q at exit 0", tt.harness, tt.tier, model, exit, stderr, tt.want)
			}
		})
	}
}

func TestResolveModelVerdict(t *testing.T) {
	tests := []struct {
		name         string
		harness      string
		benchModel   string
		benchSet     bool
		src          Source
		wantModel    string
		wantExit     int
		wantContains string
	}{
		{"unrouted-unset", "codex", "", false, Source{Path: ".bench/lines.env"}, "", 0, ""},
		{"unrouted-explicit", "codex", "anything", true, Source{Path: ".bench/lines.env"}, "anything", 0, ""},
		// A binding that names no cell at all degrades to the BENCH_MODEL passthrough an
		// unadopted repo relies on, rather than bricking every shift.
		{"unbound-matrix", "codex", "gpt-9", true, bound("# nothing bound\n"), "gpt-9", 0, "no BENCH_<HARNESS>_<TIER> cell is bound"},
		{"routed-unset", "codex", "", false, bound(fullBinding), "", 1, "requires a declared line"},
		{"routed-empty-set", "codex", "", true, bound(fullBinding), "", 1, "requires a declared line"},
		// BENCH_MODEL names a tier and never a concrete id, so a bound model id is as
		// unroutable as an unknown token.
		{"routed-model-id-not-a-tier", "codex", "gpt-5.6-sol", true, bound(fullBinding), "", 1, "is not a tier"},
		{"routed-unknown-tier", "codex", "gpt-9", true, bound(fullBinding), "", 1, "is not a tier"},
		{"unknown-harness", "gemini", "top", true, bound(fullBinding), "", 1, "unknown harness 'gemini'"},
		// An unadopted harness fails closed: no fallback to another harness's column.
		{"unbound-column", "opencode", "top", true, bound(fullBinding), "", 1, "the opencode column is unbound"},
		{"partial-column", "claude", "top", true, bound("BENCH_CODEX_TOP=a\nBENCH_CODEX_MID=b\nBENCH_CODEX_CHEAP=c\nBENCH_CLAUDE_TOP=x\n"), "", 1, "the claude column is unbound"},
		{"foreign-key", "codex", "top", true, bound(fullBinding + "BENCH_GEMINI_TOP=g\n"), "", 1, "unknown harness key BENCH_GEMINI_TOP"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, exit, stderr := ResolveModelVerdict(tt.harness, tt.benchModel, tt.benchSet, tt.src)
			if model != tt.wantModel {
				t.Errorf("model = %q, want %q", model, tt.wantModel)
			}
			if exit != tt.wantExit {
				t.Errorf("exit = %d, want %d (stderr=%q)", exit, tt.wantExit, stderr)
			}
			if tt.wantContains == "" {
				if stderr != "" {
					t.Errorf("stderr = %q, want empty", stderr)
				}
			} else if !contains(stderr, tt.wantContains) {
				t.Errorf("stderr = %q, want to contain %q", stderr, tt.wantContains)
			}
		})
	}
}

func TestResolveModelVerdictNamesTheBindingPath(t *testing.T) {
	src := Source{Path: "custom/path.env", Exists: true, Content: []byte(fullBinding)}
	_, _, stderr := ResolveModelVerdict("opencode", "top", true, src)
	if !contains(stderr, "in custom/path.env") {
		t.Errorf("stderr = %q, want to name the path", stderr)
	}
}

// TestCellFault pins the per-harness cell grammar: every cell is an opaque safe token, and
// opencode's provider-qualified namespace is a rule on ITS cells rather than a filter on a
// resolved value.
func TestCellFault(t *testing.T) {
	for _, tt := range []struct {
		harness, value, want string
	}{
		{"codex", "gpt-5.6-sol", ""},
		{"claude", "opus", ""},
		{"codex", "gpt 5", "is not a safe model token"},
		{"opencode", "", "is not a safe model token"},
		{"opencode", "openai/gpt-5.6-luna", ""},
		{"opencode", "openrouter/google/gemini-2.5-flash", ""},
		{"opencode", "gpt-5.6-luna", "is not provider-qualified"},
		{"opencode", "openai//gpt-5", "is not provider-qualified"},
		{"opencode", "openai/gpt-5/", "is not provider-qualified"},
	} {
		got := CellFault(tt.harness, tt.value)
		if (tt.want == "" && got != "") || (tt.want != "" && !contains(got, tt.want)) {
			t.Errorf("CellFault(%q, %q) = %q, want %q", tt.harness, tt.value, got, tt.want)
		}
	}
}

func TestResolveModelVerdictRejectsMalformedColumn(t *testing.T) {
	bare := "BENCH_OPENCODE_TOP=gpt-5.6-sol\nBENCH_OPENCODE_MID=gpt-5.6-terra\nBENCH_OPENCODE_CHEAP=gpt-5.6-luna\n"
	model, exit, stderr := ResolveModelVerdict("opencode", "cheap", true, bound(bare))
	if model != "" || exit != 1 || !contains(stderr, "is not provider-qualified") {
		t.Fatalf("ResolveModelVerdict on a bare opencode column = (%q, %d, %q), want a fail-closed grammar error", model, exit, stderr)
	}
	provider := "BENCH_OPENCODE_TOP=openai/gpt-5.6-sol\nBENCH_OPENCODE_MID=openrouter/google/gemini-2.5-flash\nBENCH_OPENCODE_CHEAP=openai/gpt-5.6-luna\n"
	for _, tt := range []struct{ tier, want string }{
		{"cheap", "openai/gpt-5.6-luna"},
		{"mid", "openrouter/google/gemini-2.5-flash"},
	} {
		model, exit, stderr := ResolveModelVerdict("opencode", tt.tier, true, bound(provider))
		if model != tt.want || exit != 0 || stderr != "" {
			t.Errorf("ResolveModelVerdict(opencode, %q) = (%q, %d, %q), want %q at exit 0", tt.tier, model, exit, stderr, tt.want)
		}
	}
}

func TestAgentLineVerdict(t *testing.T) {
	tests := []struct {
		name         string
		stdin        []byte
		harness      string
		src          Source
		wantExit     int
		wantContains string // substring the stderr must contain ("" means stderr must be empty)
	}{
		{"own-column-allows", envelope("opus-4-8"), "claude", bound(fullBinding), 0, ""},
		// Enforcement is permissive across the whole matrix: a Claude session may name the
		// tier a Codex delegate will run on.
		{"other-column-allows", envelope("gpt-5.6-terra"), "claude", bound(fullBinding), 0, ""},
		{"tier-name-denies", envelope("mid"), "claude", bound(fullBinding), 2, "is not a bound tier"},
		{"unbound-model-denies", envelope("gpt-9"), "claude", bound(fullBinding), 2, "is not a bound tier"},
		{"malformed-stdin", []byte(`not json`), "claude", bound(fullBinding), 0, "not parseable as JSON"},
		// Routed + a complete column for the asking harness: an omitted or whitespace-only
		// model is the attack path the guard exists for (a silent inherit of the session's
		// model) and denies.
		{"missing-model-routed-complete-denies", []byte(`{"tool_input":{}}`), "claude", bound(fullBinding), 2, "missing or empty model field"},
		{"whitespace-model-routed-complete-denies", envelope("   "), "claude", bound(fullBinding), 2, "missing or empty model field"},
		// Regression rims: a missing model with no column to enforce stays fail-open.
		// A non-fork delegation type is not the fork branch: the routed missing-model deny
		// and the fail-open rims below own every envelope that is not exactly a fork.
		{"missing-model-non-fork-denies", agentEnvelope(`"general-purpose"`, ""), "claude", bound(fullBinding), 2, "missing or empty model field"},
		{"missing-model-unrouted-fails-open", []byte(`{"tool_input":{}}`), "claude", Source{Path: ".bench/lines.env"}, 0, "no resolvedModel/model field"},
		{"missing-model-unbound-column-fails-open", []byte(`{"tool_input":{}}`), "opencode", bound(fullBinding), 0, "no resolvedModel/model field"},
		{"absent-lines-env", envelope("opus-4-8"), "claude", Source{Path: ".bench/lines.env"}, 0, "no .bench/lines.env"},
		{"unbound-matrix", envelope("opus-4-8"), "claude", bound("# nothing bound\n"), 0, "no BENCH_<HARNESS>_<TIER> cell is bound"},
		{"incomplete-column", envelope("gpt-9"), "claude", bound("BENCH_CODEX_TOP=a\nBENCH_CODEX_MID=b\nBENCH_CODEX_CHEAP=c\nBENCH_CLAUDE_TOP=x\n"), 0, "the claude column is incomplete"},
		// An unadopted harness has no column to advise in, so the guard degrades rather
		// than denying in tokens nobody can pass.
		{"unbound-column-fails-open", envelope("gpt-9"), "opencode", bound(fullBinding), 0, "the opencode column is incomplete"},
		{"unknown-harness", envelope("gpt-9"), "gemini", bound(fullBinding), 1, "unknown harness 'gemini'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit, stderr := AgentLineVerdict(tt.stdin, tt.harness, tt.src)
			if exit != tt.wantExit {
				t.Errorf("exit = %d, want %d (stderr=%q)", exit, tt.wantExit, stderr)
			}
			if tt.wantContains == "" {
				if stderr != "" {
					t.Errorf("stderr = %q, want empty", stderr)
				}
			} else if !contains(stderr, tt.wantContains) {
				t.Errorf("stderr = %q, want to contain %q", stderr, tt.wantContains)
			}
		})
	}
}

// TestAgentLineVerdictDenyMessage pins the denial in full for each asking harness: the
// advice is the asking harness's own column AS PARSED, and no other family's tokens ride
// along. The fixture's claude cells differ from this repo's, so a renderer hard-coded to
// the live binding fails here.
func TestAgentLineVerdictDenyMessage(t *testing.T) {
	for _, tt := range []struct {
		harness, want, absent string
	}{
		{
			harness: "claude",
			want: "DENIED: delegation model 'gpt-9' is not a bound tier; harness claude binds top=fable-5 mid=opus-4-8 " +
				"cheap=sonnet-5 (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding.",
			absent: "gpt-5.6-sol",
		},
		{
			harness: "codex",
			want: "DENIED: delegation model 'gpt-9' is not a bound tier; harness codex binds top=gpt-5.6-sol mid=gpt-5.6-terra " +
				"cheap=gpt-5.6-luna (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding.",
			absent: "opus-4-8",
		},
	} {
		t.Run(tt.harness, func(t *testing.T) {
			exit, stderr := AgentLineVerdict(envelope("gpt-9"), tt.harness, bound(fullBinding))
			if exit != 2 {
				t.Fatalf("exit = %d, want 2 (stderr=%q)", exit, stderr)
			}
			if stderr != tt.want {
				t.Errorf("stderr =\n%q\nwant\n%q", stderr, tt.want)
			}
			if contains(stderr, tt.absent) {
				t.Errorf("stderr names a token from another harness's column (%q): %q", tt.absent, stderr)
			}
		})
	}
}

func TestAgentLineVerdictMissingModelDenyMessage(t *testing.T) {
	// The routed missing-model deny is distinct from the unbound deny: it names the
	// missing/empty field and advises in the asking harness's own column.
	exit, stderr := AgentLineVerdict([]byte(`{"tool_input":{}}`), "claude", bound(fullBinding))
	if exit != 2 {
		t.Fatalf("exit = %d, want 2 (stderr=%q)", exit, stderr)
	}
	for _, want := range []string{
		"missing or empty model field",
		"bound tier token",
		"harness claude binds top=fable-5 mid=opus-4-8 cheap=sonnet-5",
	} {
		if !contains(stderr, want) {
			t.Errorf("stderr = %q, want to contain %q", stderr, want)
		}
	}
	if contains(stderr, "gpt-5.6-sol") {
		t.Errorf("missing-model deny names another harness's column: %q", stderr)
	}
	// It must not reuse the unbound wording, whose recovery advice is different.
	if contains(stderr, "is not a bound tier") {
		t.Errorf("missing-model deny reused the unbound wording: %q", stderr)
	}
}

// TestAgentLineVerdictForkDelegation pins both fork verdicts, which are opposite policies
// on the same delegation type: a declared model is a claim the harness discards, so it
// denies, while an omitted one is the honest signal for behavior nobody can avoid, so it
// allows. Collapsing them into a single deny-all-forks or allow-all-forks branch fails one
// half. The warning states the inheritance without naming a model — the invoking session's
// tier is unknowable at this hook event.
func TestAgentLineVerdictForkDelegation(t *testing.T) {
	exit, stderr := AgentLineVerdict(forkEnvelope(`"sonnet-5"`), "claude", bound(fullBinding))
	if exit != 2 {
		t.Fatalf("fork declaring a bound model = exit %d, want 2 (stderr=%q)", exit, stderr)
	}
	if stderr != "DENIED: delegation model 'sonnet-5' is declared on a fork, which runs on this session's model and "+
		"ignores the declaration — the guard cannot verify a line the harness will not honor. Re-delegate the fork "+
		"with no model field to inherit this session's line, or spawn a non-fork delegate on a bound tier token." {
		t.Errorf("fork deny stderr = %q", stderr)
	}

	exit, stderr = AgentLineVerdict(forkEnvelope(""), "claude", bound(fullBinding))
	if exit != 0 {
		t.Fatalf("fork declaring no model = exit %d, want 0 (stderr=%q)", exit, stderr)
	}
	if stderr != "WARNING: check-agent-line: a fork delegation inherits this session's model, which no hook event "+
		"reports — allowing delegation." {
		t.Errorf("fork allow stderr = %q", stderr)
	}
}

// TestAgentLineVerdictAllowsEveryBoundCell pins permissive enforcement across the whole
// matrix through ONE harness's guard: a Claude session may legitimately name the tier a
// Codex delegate will run on, so narrowing enforcement to the asking harness's column
// would deny half of these.
func TestAgentLineVerdictAllowsEveryBoundCell(t *testing.T) {
	for _, model := range []string{
		"gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna",
		"fable-5", "opus-4-8", "sonnet-5",
	} {
		t.Run(model, func(t *testing.T) {
			if exit, stderr := AgentLineVerdict(envelope(model), "claude", bound(fullBinding)); exit != 0 || stderr != "" {
				t.Errorf("AgentLineVerdict(%q) = (%d, %q), want a silent allow", model, exit, stderr)
			}
		})
	}
}

// TestSubagentTypeNeverImpersonatesAFork pins the exact-string discriminator. tool_input is
// attacker-shaped text, so only the literal selects the fork branch; anything else — absent,
// blank, the wrong JSON type, padded, recased, or merely containing the word — keeps the
// non-fork postures. Each value is driven twice, because the two fork branches sit on
// opposite sides of the non-fork verdicts: a bound model must still allow, and an omitted
// one must still hit the routed deny rather than the fork warning.
func TestSubagentTypeNeverImpersonatesAFork(t *testing.T) {
	for _, subagentType := range []string{
		``, `""`, `5`, `{}`, `[]`, `null`, `true`,
		`" fork"`, `"fork "`, `"fork\n"`, `"Fork"`, `"FORK"`, `"forked"`, `"my-fork"`, `"general-purpose"`,
	} {
		t.Run(subagentType, func(t *testing.T) {
			if exit, stderr := AgentLineVerdict(agentEnvelope(subagentType, `"opus-4-8"`), "claude", bound(fullBinding)); exit != 0 || stderr != "" {
				t.Errorf("subagent_type %s with a bound model = (%d, %q), want a silent allow", subagentType, exit, stderr)
			}
			exit, stderr := AgentLineVerdict(agentEnvelope(subagentType, ""), "claude", bound(fullBinding))
			if exit != 2 || !contains(stderr, "missing or empty model field") {
				t.Errorf("subagent_type %s with no model = (%d, %q), want the routed missing-model deny", subagentType, exit, stderr)
			}
		})
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
