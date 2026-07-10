// Package lines is the tier-binding parser and verdict engine behind invariant #2
// (the declared line). It is the single source that both the agent-line PreToolUse
// guard (check-agent-line) and the headless-shift adapter consult, so the model-tier
// enforcement and its advertised binding cannot drift. The cmd layer reads
// .bench/lines.env and the PreToolUse envelope from disk/stdin and passes the bytes in;
// everything here is pure, so the load-bearing message wording and exit codes are unit
// testable without a repo.
//
// This is a faithful port of .bench/lib/lines-env.sh: last assignment wins, the value is
// the text after the first '=', with one trailing CR stripped, whitespace trimmed
// (trailing then leading), and one surrounding double-quote pair then one single-quote
// pair removed independently.
package lines

import (
	"encoding/json"
	"strings"

	"github.com/gibbonmi/bench/internal/modelid"
)

// TierValue is the pure port of bench_tier_value, operating on file content rather than a
// path. It returns the value bound to key — the LAST line matching `^[[:space:]]*<key>=` —
// or "" when the key is absent or its value is empty. Ordering is load-bearing: last
// assignment wins; value is the text after the first '='; then one trailing CR, whitespace
// trim (trailing then leading), one surrounding double-quote pair, one surrounding
// single-quote pair — the last two stripped independently, not as a matched pair.
func TierValue(key string, content []byte) string {
	prefix := key + "="
	found := false
	var value string
	for _, line := range strings.Split(string(content), "\n") {
		// Match `^[[:space:]]*<key>=`: skip leading whitespace, then require the literal
		// key followed by '='. The '=' anchors the key so KEY never matches KEYX.
		rest := strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(rest, prefix) {
			continue
		}
		found = true
		// Value = everything after the FIRST '=' on the ORIGINAL line (the shell strips
		// from the first '=', and leading whitespace is trimmed later anyway).
		value = line[strings.Index(line, "=")+1:]
	}
	if !found {
		return ""
	}
	// Step 3: strip exactly ONE trailing carriage return (the shell's `%$'\r'`).
	value = strings.TrimSuffix(value, "\r")
	// Step 4: trim trailing then leading whitespace. The shell's `[[:space:]]`
	// includes CR/FF/VT, so a value like `x\r ` (a CR followed by a space) collapses
	// to `x` there — the trailing-trim must match that cutset, not just space+tab, or
	// a hand-edited binding with interleaved CRs would flip an allow/deny.
	const ws = " \t\r\v\f"
	value = strings.TrimRight(value, ws)
	value = strings.TrimLeft(value, ws)
	// Step 5: strip ONE leading and ONE trailing double-quote (independent).
	value = strings.TrimPrefix(value, "\"")
	value = strings.TrimSuffix(value, "\"")
	// Step 6: strip ONE leading and ONE trailing single-quote (independent).
	value = strings.TrimPrefix(value, "'")
	value = strings.TrimSuffix(value, "'")
	return value
}

// ModelFromEnvelope parses the Claude Code Agent PreToolUse envelope, returning
// tool_input.resolvedModel if non-empty, else tool_input.model if non-empty, else "". It
// returns a non-nil error ONLY when data is not valid JSON, mirroring the Python shim
// raising on unparseable stdin; valid JSON with no matching field returns ("", nil).
func ModelFromEnvelope(data []byte) (model string, err error) {
	var e struct {
		ToolInput struct {
			ResolvedModel string `json:"resolvedModel"`
			Model         string `json:"model"`
		} `json:"tool_input"`
	}
	if err := json.Unmarshal(data, &e); err != nil {
		return "", err
	}
	// An all-whitespace field is blank for the non-empty test: a whitespace resolvedModel
	// falls back to model, and an all-blank envelope yields "". This closes the whitespace
	// model value at the parse boundary and folds the omitted/empty/whitespace cases into
	// one branch of AgentLineVerdict.
	if strings.TrimSpace(e.ToolInput.ResolvedModel) != "" {
		return e.ToolInput.ResolvedModel, nil
	}
	if strings.TrimSpace(e.ToolInput.Model) != "" {
		return e.ToolInput.Model, nil
	}
	return "", nil
}

// Binding is the resolved tier binding from .bench/lines.env: the three model tiers and
// their optional aliases.
type Binding struct {
	Top, Mid, Cheap                string
	AliasTop, AliasMid, AliasCheap string
}

type tierBinding struct {
	model string
	alias string
}

func (b Binding) tiers() [3]tierBinding {
	return [3]tierBinding{
		{model: b.Top, alias: b.AliasTop},
		{model: b.Mid, alias: b.AliasMid},
		{model: b.Cheap, alias: b.AliasCheap},
	}
}

// ParseBinding fills each field of a Binding via TierValue with the corresponding
// BENCH_TIER_* / BENCH_ALIAS_* key.
func ParseBinding(content []byte) Binding {
	return Binding{
		Top:        TierValue("BENCH_TIER_TOP", content),
		Mid:        TierValue("BENCH_TIER_MID", content),
		Cheap:      TierValue("BENCH_TIER_CHEAP", content),
		AliasTop:   TierValue("BENCH_ALIAS_TOP", content),
		AliasMid:   TierValue("BENCH_ALIAS_MID", content),
		AliasCheap: TierValue("BENCH_ALIAS_CHEAP", content),
	}
}

func warn(s string) string {
	return "WARNING: check-agent-line: " + s + " — allowing delegation."
}

func dash(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

// AgentLineVerdict is the pure agent-line verdict. Every degraded branch is fail-OPEN
// (exit 0) with a one-line stderr warning; ONLY a present model that matches no bound tier
// or alias denies (exit 2). The returned stderr carries no trailing newline — the caller
// adds one.
func AgentLineVerdict(stdin []byte, linesEnvExists bool, linesEnvContent []byte) (exitCode int, stderr string) {
	model, err := ModelFromEnvelope(stdin)
	if err != nil {
		return 0, warn("stdin is not parseable as JSON")
	}
	if model == "" {
		// A routed repo with a complete binding is the one degraded branch that is also
		// the attack path the guard exists for: an omitted or empty model inherits the
		// invoking session's model, the silent escalation invariant #2 forbids. Deny it.
		// Every other missing-model branch (unrouted, incomplete binding) keeps the
		// fail-open rim — there is no binding to enforce, and a broken guard must never
		// brick delegation.
		if linesEnvExists {
			if b := ParseBinding(linesEnvContent); b.Top != "" && b.Mid != "" && b.Cheap != "" {
				return 2, "DENIED: the delegation envelope has a missing or empty model field — an " +
					"omitted model silently inherits this session's model, which invariant #2 forbids. " +
					"Pass a bound alias from .bench/lines.env; " + describeBoundTiers(b) +
					" (see .bench/lines.env and the craft-line skill). Re-delegate on a bound alias."
			}
		}
		return 0, warn("no resolvedModel/model field in tool_input")
	}
	if !linesEnvExists {
		return 0, warn("no .bench/lines.env at repo root")
	}
	b := ParseBinding(linesEnvContent)
	if b.Top == "" || b.Mid == "" || b.Cheap == "" {
		return 0, warn("a BENCH_TIER_* value is unset or empty in .bench/lines.env")
	}
	for _, tier := range b.tiers() {
		if model == tier.model || (tier.alias != "" && model == tier.alias) {
			return 0, ""
		}
	}
	return 2, "DENIED: delegation model '" + model + "' is not a bound tier; " + describeBoundTiers(b) +
		" (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding."
}

// describeBoundTiers formats the tier-and-alias listing both deny messages carry, so the
// bound-tiers fact has one source: `bound: top=… mid=… cheap=… aliases: top=… mid=… cheap=…`,
// with a dash for each unset alias.
func describeBoundTiers(b Binding) string {
	return "bound: top=" + b.Top + " mid=" + b.Mid + " cheap=" + b.Cheap +
		" aliases: top=" + dash(b.AliasTop) + " mid=" + dash(b.AliasMid) + " cheap=" + dash(b.AliasCheap)
}

// ResolveModelVerdict is the headless-shift adapter's model resolution. Unlike the
// agent-line verdict, aliases do NOT apply here — only the three tier ids are bound
// targets. The returned stderr carries no trailing newline.
func ResolveModelVerdict(benchModel string, benchModelSet bool, linesEnvExists bool, linesEnvPath string, content []byte) (model string, exitCode int, stderr string) {
	return resolveModelVerdict(benchModel, benchModelSet, linesEnvExists, linesEnvPath, content, false)
}

// ResolveModelAliasVerdict applies the same tier-id validation as
// ResolveModelVerdict, then projects the matched tier to its corresponding alias.
func ResolveModelAliasVerdict(benchModel string, benchModelSet bool, linesEnvExists bool, linesEnvPath string, content []byte) (model string, exitCode int, stderr string) {
	return resolveModelVerdict(benchModel, benchModelSet, linesEnvExists, linesEnvPath, content, true)
}

// ResolveProviderModelVerdict applies the exact tier-id validation, then requires the
// provider/model shape used by harnesses whose model namespace is provider-qualified.
func ResolveProviderModelVerdict(benchModel string, benchModelSet bool, linesEnvExists bool, linesEnvPath string, content []byte) (model string, exitCode int, stderr string) {
	model, exitCode, stderr = ResolveModelVerdict(benchModel, benchModelSet, linesEnvExists, linesEnvPath, content)
	if exitCode != 0 || model == "" {
		return model, exitCode, stderr
	}
	if !isProviderModel(model) {
		return "", 1, "bench shift: BENCH_MODEL='" + model + "' is incompatible with a provider/model harness; use a value listed by 'opencode models' in BENCH_MODEL (and in the matching BENCH_TIER_* binding when routed)"
	}
	return model, 0, stderr
}

func isProviderModel(model string) bool {
	parts := strings.Split(model, "/")
	if len(parts) < 2 || !modelid.SafeToken(model) {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
	}
	return true
}

func resolveModelVerdict(benchModel string, benchModelSet bool, linesEnvExists bool, linesEnvPath string, content []byte, alias bool) (model string, exitCode int, stderr string) {
	if !linesEnvExists {
		// Unrouted: explicit BENCH_MODEL wins; absent -> "".
		return benchModel, 0, ""
	}
	b := ParseBinding(content)
	if b.Top == "" || b.Mid == "" || b.Cheap == "" {
		return benchModel, 0, "WARNING: bench adapter: a BENCH_TIER_* value is unset or empty in " +
			linesEnvPath + " — ignoring the binding and falling back to BENCH_MODEL."
	}
	if !benchModelSet || benchModel == "" {
		return "", 1, "bench shift in a routed repo requires a declared line: set BENCH_MODEL to one of top=" +
			b.Top + " mid=" + b.Mid + " cheap=" + b.Cheap
	}
	for _, tier := range b.tiers() {
		if benchModel != tier.model {
			continue
		}
		if !alias {
			return tier.model, 0, ""
		}
		if tier.alias == "" {
			return "", 1, "bench shift: BENCH_MODEL='" + benchModel + "' has no bound alias; set the corresponding BENCH_ALIAS_* value in " + linesEnvPath
		}
		return tier.alias, 0, ""
	}
	return "", 1, "bench shift: BENCH_MODEL='" + benchModel + "' is not a bound model; set it to one of top=" +
		b.Top + " mid=" + b.Mid + " cheap=" + b.Cheap
}

// DescribeBinding returns the denies-clause body the shim prints after "denies: ".
func DescribeBinding(linesEnvExists bool, content []byte) string {
	if !linesEnvExists {
		return "unrouted (no .bench/lines.env binding)"
	}
	b := ParseBinding(content)
	return "Agent delegation off the bound line (top=" + dash(b.Top) + " mid=" + dash(b.Mid) +
		" cheap=" + dash(b.Cheap) + ")"
}
