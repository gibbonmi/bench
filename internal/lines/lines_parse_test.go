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
// ParseBinding. A schema migration that dropped quoted, CRLF, last-wins, or
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
		// whitespace. This matches the shell, whose trailing-whitespace strip removes
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

// TestParseBindingRejectsForeignHarnessKeys pins the closed harness set. A
// BENCH_<HARNESS>_<TIER>-shaped key naming a harness outside it creates no column and
// is reported. The retired schema is exactly such a key. So a lines.env carrying only
// retired keys binds nothing at all, rather than reading as a legacy binding.
func TestParseBindingRejectsForeignHarnessKeys(t *testing.T) {
	for _, tt := range []struct {
		name, content string
		want          []string
	}{
		{"unknown harness", "BENCH_CODEX_TOP=a\nBENCH_CODEX_MID=b\nBENCH_CODEX_CHEAP=c\nBENCH_GEMINI_TOP=g\n", []string{"BENCH_GEMINI_TOP"}},
		{"retired schema", retiredBinding, []string{"BENCH_TIER_TOP", "BENCH_TIER_MID", "BENCH_TIER_CHEAP", "BENCH_ALIAS_TOP", "BENCH_ALIAS_MID", "BENCH_ALIAS_CHEAP"}},
		// A key is a cell only under the canonical spelling Key renders, so a lowercased
		// harness segment is foreign. The two readings have to agree: Cell looks up the
		// canonical key and nothing else. A spelling counted as naming a known harness
		// here would be a key neither path ever reads. It would be ignored by the reader
		// and unreported by the gate at once.
		{"non-canonical harness case", "BENCH_CODEX_TOP=a\nBENCH_CODEX_MID=b\nBENCH_CODEX_CHEAP=c\nBENCH_claude_TOP=f\n", []string{"BENCH_claude_TOP"}},
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

// TestRetiredSchemaBindsNothing drives the hard cut through both verdicts. A lines.env
// carrying only retired keys must resolve nothing and enforce nothing, exactly as if it
// bound no cell.
func TestRetiredSchemaBindsNothing(t *testing.T) {
	src := bound(retiredBinding)
	if model, exit, stderr := ResolveModelVerdict("codex", "top", true, src); model != "" || exit != 1 {
		t.Errorf("ResolveModelVerdict on retired keys = (%q, %d, %q), want no model and exit 1", model, exit, stderr)
	}
	// The retired ids and aliases are no longer bound tokens, so enforcement is not run
	// against them. The guard keeps its fail-open rim instead of denying on a stale column.
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
