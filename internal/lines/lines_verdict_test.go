package lines

import "testing"

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

func TestResolveModelVerdictNamesTheBindingPath(t *testing.T) {
	src := Source{Path: "custom/path.env", Exists: true, Content: []byte(fullBinding)}
	_, _, stderr := ResolveModelVerdict("opencode", "top", true, src)
	if !contains(stderr, "in custom/path.env") {
		t.Errorf("stderr = %q, want to name the path", stderr)
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

// TestCellFault pins the per-harness cell grammar: every cell is an opaque safe token.
// opencode's provider-qualified namespace is a rule on ITS cells rather than a filter on
// a resolved value.
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
		// The rule reads the record's provider binding, not a harness name. The row that
		// binds any provider demands the qualification; a row that binds one provider
		// takes a bare id.
		{"opencode", "gpt-5", "is not provider-qualified"},
		{"codex", "gpt-5", ""},
	} {
		got := CellFault(tt.harness, tt.value)
		if (tt.want == "" && got != "") || (tt.want != "" && !contains(got, tt.want)) {
			t.Errorf("CellFault(%q, %q) = %q, want %q", tt.harness, tt.value, got, tt.want)
		}
	}
}
