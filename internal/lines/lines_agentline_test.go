package lines

import "testing"

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
		// model is the attack path the guard exists for. It silently inherits the session's
		// model and denies.
		{"missing-model-routed-complete-denies", []byte(`{"tool_input":{}}`), "claude", bound(fullBinding), 2, "missing or empty model field"},
		{"whitespace-model-routed-complete-denies", envelope("   "), "claude", bound(fullBinding), 2, "missing or empty model field"},
		// Regression rims: a missing model with no column to enforce stays fail-open.
		// A non-fork delegation type is not the fork branch. The routed missing-model deny
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

// TestAgentLineVerdictDenyMessage pins the denial in full for each asking harness. The
// advice is the asking harness's own column AS PARSED, and no other family's tokens
// ride along. The fixture's claude cells differ from this repo's, so a renderer
// hard-coded to the live binding fails here.
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

// TestAgentLineVerdictForkDelegation pins both fork verdicts, which are opposite
// policies on the same delegation type. A declared model is a claim the harness
// discards, so it denies. An omitted one is the honest signal for behavior nobody can
// avoid, so it allows. Collapsing them into a single deny-all-forks or allow-all-forks
// branch fails one half. The warning states the inheritance without naming a model —
// the invoking session's tier is unknowable at this hook event.
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

// TestAgentLineVerdictForkDenyNeedsARoutedCompleteBinding drives both fork verdicts across
// all three binding states. Without a complete column there is no bound tier, so a
// declared model has nothing to escalate off. A repo that never opted into line
// enforcement must not have its delegations blocked. The fork deny therefore rides the
// same routed-and-complete condition as the missing-model deny.
//
// Every other state takes the fail-open rim. The inheritance warning stays
// unconditional: it is a warning either way. Gating it would cost the operator a true
// statement about what a fork does.
func TestAgentLineVerdictForkDenyNeedsARoutedCompleteBinding(t *testing.T) {
	const inheritanceWarning = "WARNING: check-agent-line: a fork delegation inherits this session's model, " +
		"which no hook event reports — allowing delegation."
	for _, tt := range []struct {
		name     string
		src      Source
		wantDeny bool
	}{
		{"unrouted", Source{Path: ".bench/lines.env"}, false},
		{"incomplete-column", bound("BENCH_CLAUDE_TOP=fable-5\n"), false},
		{"routed-complete", bound(fullBinding), true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			exit, stderr := AgentLineVerdict(forkEnvelope(`"sonnet-5"`), "claude", tt.src)
			switch {
			case tt.wantDeny && (exit != 2 || !contains(stderr, "declared on a fork")):
				t.Errorf("fork declaring a model = (%d, %q), want the fork deny", exit, stderr)
			case !tt.wantDeny && (exit != 0 || stderr != inheritanceWarning):
				t.Errorf("fork declaring a model = (%d, %q), want (0, %q)", exit, stderr, inheritanceWarning)
			}
			if exit, stderr := AgentLineVerdict(forkEnvelope(""), "claude", tt.src); exit != 0 || stderr != inheritanceWarning {
				t.Errorf("fork declaring no model = (%d, %q), want (0, %q)", exit, stderr, inheritanceWarning)
			}
		})
	}
}

// TestAgentLineVerdictAllowsEveryBoundCell pins permissive enforcement across the whole
// matrix through ONE harness's guard. A Claude session may legitimately name the tier a
// Codex delegate will run on. So narrowing enforcement to the asking harness's column
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

// TestSubagentTypeNeverImpersonatesAFork pins the exact-string discriminator. tool_input
// is attacker-shaped text, so only the literal selects the fork branch. Anything else,
// absent, blank, the wrong JSON type, padded, recased, or merely containing the word,
// keeps the non-fork postures. Each value is driven twice, because the two fork
// branches sit on opposite sides of the non-fork verdicts. A bound model must still
// allow, and an omitted one must still hit the routed deny rather than the fork
// warning.
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
