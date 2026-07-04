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
	if e.ToolInput.ResolvedModel != "" {
		return e.ToolInput.ResolvedModel, nil
	}
	return e.ToolInput.Model, nil
}

// Binding is the resolved tier binding from .bench/lines.env: the three model tiers and
// their optional aliases.
type Binding struct {
	Top, Mid, Cheap                string
	AliasTop, AliasMid, AliasCheap string
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
		return 0, warn("no resolvedModel/model field in tool_input")
	}
	if !linesEnvExists {
		return 0, warn("no .bench/lines.env at repo root")
	}
	b := ParseBinding(linesEnvContent)
	if b.Top == "" || b.Mid == "" || b.Cheap == "" {
		return 0, warn("a BENCH_TIER_* value is unset or empty in .bench/lines.env")
	}
	for _, bound := range []string{b.Top, b.Mid, b.Cheap, b.AliasTop, b.AliasMid, b.AliasCheap} {
		if bound != "" && model == bound {
			return 0, ""
		}
	}
	return 2, "DENIED: delegation model '" + model + "' is not a bound tier; bound: top=" +
		b.Top + " mid=" + b.Mid + " cheap=" + b.Cheap + " aliases: top=" + dash(b.AliasTop) +
		" mid=" + dash(b.AliasMid) + " cheap=" + dash(b.AliasCheap) +
		" (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding."
}

// ResolveModelVerdict is the headless-shift adapter's model resolution. Unlike the
// agent-line verdict, aliases do NOT apply here — only the three tier ids are bound
// targets. The returned stderr carries no trailing newline.
func ResolveModelVerdict(benchModel string, benchModelSet bool, linesEnvExists bool, linesEnvPath string, content []byte) (model string, exitCode int, stderr string) {
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
	if benchModel != b.Top && benchModel != b.Mid && benchModel != b.Cheap {
		return "", 1, "bench shift: BENCH_MODEL='" + benchModel + "' is not a bound model; set it to one of top=" +
			b.Top + " mid=" + b.Mid + " cheap=" + b.Cheap
	}
	return benchModel, 0, ""
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
