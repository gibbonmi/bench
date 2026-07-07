package lines

import (
	"strings"
	"testing"
)

// pinContent is the retiring gate pin a0 fixture: every load-bearing TierValue case in one
// file. The final line has NO trailing newline and must read intact.
const pinContent = "BENCH_TIER_TOP=\"claude-quoted-1\"\n" +
	"BENCH_TIER_MID='claude-squoted-2'\n" +
	"BENCH_TIER_CHEAP=claude-crlf-3\r\n" +
	"BENCH_ALIAS_TOP=claude-trail-4   \n" +
	"   BENCH_ALIAS_MID=claude-indent-5\n" +
	"BENCH_ALIAS_CHEAP=claude-first-6\n" +
	"BENCH_ALIAS_CHEAP=claude-last-7\n" +
	"BENCH_EMPTY=\n" +
	"BENCH_NONEWLINE=claude-nonl-8"

func TestTierValue(t *testing.T) {
	tests := []struct {
		name, key, want string
	}{
		{"double-quoted", "BENCH_TIER_TOP", "claude-quoted-1"},
		{"single-quoted", "BENCH_TIER_MID", "claude-squoted-2"},
		{"trailing-cr", "BENCH_TIER_CHEAP", "claude-crlf-3"},
		{"trailing-spaces", "BENCH_ALIAS_TOP", "claude-trail-4"},
		{"indented-key", "BENCH_ALIAS_MID", "claude-indent-5"},
		{"last-wins", "BENCH_ALIAS_CHEAP", "claude-last-7"},
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

func TestTierValueEdgeCases(t *testing.T) {
	tests := []struct {
		name, key, content, want string
	}{
		// The '=' anchors the key; KEY must not match KEYX.
		{"prefix-not-matched", "BENCH_TIER_TOP", "BENCH_TIER_TOPX=nope\n", ""},
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

func TestParseBinding(t *testing.T) {
	content := "BENCH_TIER_TOP=fable-5\n" +
		"BENCH_TIER_MID=opus-4-8\n" +
		"BENCH_TIER_CHEAP=sonnet-5\n" +
		"BENCH_ALIAS_TOP=top\n" +
		"BENCH_ALIAS_MID=mid\n" +
		"BENCH_ALIAS_CHEAP=cheap\n"
	got := ParseBinding([]byte(content))
	want := Binding{
		Top: "fable-5", Mid: "opus-4-8", Cheap: "sonnet-5",
		AliasTop: "top", AliasMid: "mid", AliasCheap: "cheap",
	}
	if got != want {
		t.Errorf("ParseBinding = %+v, want %+v", got, want)
	}
}

// fullBinding is a complete binding with aliases, for the verdict tables.
const fullBinding = "BENCH_TIER_TOP=fable-5\n" +
	"BENCH_TIER_MID=opus-4-8\n" +
	"BENCH_TIER_CHEAP=sonnet-5\n" +
	"BENCH_ALIAS_TOP=top-alias\n" +
	"BENCH_ALIAS_MID=mid-alias\n" +
	"BENCH_ALIAS_CHEAP=cheap-alias\n"

func envelope(model string) []byte {
	return []byte(`{"tool_input":{"resolvedModel":"` + model + `"}}`)
}

func TestAgentLineVerdict(t *testing.T) {
	tests := []struct {
		name         string
		stdin        []byte
		exists       bool
		content      string
		wantExit     int
		wantContains string // substring the stderr must contain ("" means stderr must be empty)
	}{
		{"bound-tier-allows", envelope("opus-4-8"), true, fullBinding, 0, ""},
		{"declared-alias-allows", envelope("mid-alias"), true, fullBinding, 0, ""},
		{"undeclared-alias-denies", envelope("cheap"), true, fullBinding, 2, "is not a bound tier"},
		{"unbound-model-denies", envelope("gpt-9"), true, fullBinding, 2, "is not a bound tier"},
		{"malformed-stdin", []byte(`not json`), true, fullBinding, 0, "not parseable as JSON"},
		// Routed + complete binding: an omitted or whitespace-only model is the attack
		// path the guard exists for (a silent inherit of the session's model) and denies.
		{"missing-model-routed-complete-denies", []byte(`{"tool_input":{}}`), true, fullBinding, 2, "missing or empty model field"},
		{"whitespace-model-routed-complete-denies", envelope("   "), true, fullBinding, 2, "missing or empty model field"},
		// Regression rims: a missing model with no binding to enforce stays fail-open.
		{"missing-model-unrouted-fails-open", []byte(`{"tool_input":{}}`), false, "", 0, "no resolvedModel/model field"},
		{"missing-model-incomplete-fails-open", []byte(`{"tool_input":{}}`), true, "BENCH_TIER_TOP=t\nBENCH_TIER_CHEAP=c\n", 0, "no resolvedModel/model field"},
		{"absent-lines-env", envelope("opus-4-8"), false, "", 0, "no .bench/lines.env"},
		{"incomplete-binding", envelope("opus-4-8"), true, "BENCH_TIER_TOP=t\nBENCH_TIER_CHEAP=c\n", 0, "unset or empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit, stderr := AgentLineVerdict(tt.stdin, tt.exists, []byte(tt.content))
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

func TestAgentLineVerdictDenyMessage(t *testing.T) {
	exit, stderr := AgentLineVerdict(envelope("gpt-9"), true, []byte(fullBinding))
	if exit != 2 {
		t.Fatalf("exit = %d, want 2", exit)
	}
	want := "DENIED: delegation model 'gpt-9' is not a bound tier; bound: top=fable-5 mid=opus-4-8 " +
		"cheap=sonnet-5 aliases: top=top-alias mid=mid-alias cheap=cheap-alias " +
		"(see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding."
	if stderr != want {
		t.Errorf("stderr =\n%q\nwant\n%q", stderr, want)
	}
}

func TestAgentLineVerdictMissingModelDenyMessage(t *testing.T) {
	// The routed missing-model deny is distinct from the unbound deny: it names the
	// missing/empty field, instructs re-delegating on a bound alias, and lists the tiers.
	exit, stderr := AgentLineVerdict([]byte(`{"tool_input":{}}`), true, []byte(fullBinding))
	if exit != 2 {
		t.Fatalf("exit = %d, want 2 (stderr=%q)", exit, stderr)
	}
	for _, want := range []string{
		"missing or empty model field",
		"bound alias",
		"top=fable-5 mid=opus-4-8 cheap=sonnet-5",
	} {
		if !contains(stderr, want) {
			t.Errorf("stderr = %q, want to contain %q", stderr, want)
		}
	}
	// It must not reuse the unbound wording, which offers no alias-fix guidance.
	if contains(stderr, "is not a bound tier") {
		t.Errorf("missing-model deny reused the unbound wording: %q", stderr)
	}
}

func TestAgentLineVerdictDashesAbsentAliases(t *testing.T) {
	// No aliases declared: the deny message must show dashes for each.
	binding := "BENCH_TIER_TOP=t\nBENCH_TIER_MID=m\nBENCH_TIER_CHEAP=c\n"
	_, stderr := AgentLineVerdict(envelope("x"), true, []byte(binding))
	want := "aliases: top=- mid=- cheap=-"
	if !contains(stderr, want) {
		t.Errorf("stderr = %q, want to contain %q", stderr, want)
	}
}

func TestResolveModelVerdict(t *testing.T) {
	tests := []struct {
		name         string
		benchModel   string
		benchSet     bool
		exists       bool
		content      string
		wantModel    string
		wantExit     int
		wantContains string
	}{
		{"unrouted-unset", "", false, false, "", "", 0, ""},
		{"unrouted-explicit", "anything", true, false, "", "anything", 0, ""},
		{"incomplete-binding", "opus-4-8", true, true, "BENCH_TIER_TOP=t\n", "opus-4-8", 0, "unset or empty"},
		{"routed-unset", "", false, true, fullBinding, "", 1, "requires a declared line"},
		{"routed-empty-set", "", true, true, fullBinding, "", 1, "requires a declared line"},
		{"routed-unbound", "gpt-9", true, true, fullBinding, "", 1, "is not a bound model"},
		{"routed-bound", "opus-4-8", true, true, fullBinding, "opus-4-8", 0, ""},
		// Aliases do NOT apply here: an alias id is unbound at the adapter.
		{"routed-alias-not-bound", "mid-alias", true, true, fullBinding, "", 1, "is not a bound model"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, exit, stderr := ResolveModelVerdict(tt.benchModel, tt.benchSet, tt.exists, ".bench/lines.env", []byte(tt.content))
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

func TestResolveModelVerdictIncompleteUsesPath(t *testing.T) {
	_, _, stderr := ResolveModelVerdict("m", true, true, "custom/path.env", []byte("BENCH_TIER_TOP=t\n"))
	if !contains(stderr, "in custom/path.env") {
		t.Errorf("stderr = %q, want to name the path", stderr)
	}
}

func TestDescribeBinding(t *testing.T) {
	tests := []struct {
		name    string
		exists  bool
		content string
		want    string
	}{
		{"unrouted", false, "", "unrouted (no .bench/lines.env binding)"},
		{"routed-full", true, fullBinding, "Agent delegation off the bound line (top=fable-5 mid=opus-4-8 cheap=sonnet-5)"},
		{"routed-empty-tiers", true, "", "Agent delegation off the bound line (top=- mid=- cheap=-)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescribeBinding(tt.exists, []byte(tt.content)); got != tt.want {
				t.Errorf("DescribeBinding = %q, want %q", got, tt.want)
			}
		})
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
