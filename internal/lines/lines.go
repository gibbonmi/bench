// Package lines is the tier-binding parser and verdict engine behind invariant #2, the
// declared line. Tier is the only shared identity. .bench/lines.env binds a closed
// BENCH_<HARNESS>_<TIER> matrix over Harnesses. Every runtime caller names its own
// harness, and no family is canonical.
//
// It is the single source that both the agent-line PreToolUse guard, check-agent-line,
// and the headless-shift adapters consult. This keeps the model-tier enforcement and its
// advertised binding from drifting. The cmd layer reads .bench/lines.env and the
// PreToolUse envelope from disk/stdin, and passes the bytes in. Everything here is pure,
// so the load-bearing message wording and exit codes are unit testable without a repo.
//
// This is a faithful port of .bench/lib/lines-env.sh: last assignment wins. The value is
// the text after the first '='. One trailing CR is stripped, and whitespace is trimmed,
// trailing then leading. One surrounding double-quote pair is removed, then one
// single-quote pair, independently.
package lines

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/harnesses"
	"github.com/gibbonmi/bench/internal/modelid"
)

// TierValue is the pure port of bench_tier_value, operating on file content rather than a
// path. It returns the value bound to key, the LAST line matching `^[[:space:]]*<key>=`,
// or "" when the key is absent or its value is empty. Ordering is load-bearing: last
// assignment wins. Value is the text after the first '='. Then one trailing CR is
// stripped, and whitespace is trimmed, trailing then leading. One surrounding
// double-quote pair, then one surrounding single-quote pair, are stripped independently,
// not as a matched pair.
func TierValue(key string, content []byte) string {
	value, _ := tierValue(key, content)
	return value
}

// tierValue is TierValue plus the presence bit the matrix parser needs. A key assigned
// an empty value is present but unbound, which separates a declared-but-incomplete
// column from an undeclared one. Both readings come from one scan so the key-matching
// rule has a single source.
func tierValue(key string, content []byte) (string, bool) {
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
		return "", false
	}
	// Step 3: strip exactly ONE trailing carriage return (the shell's `%$'\r'`).
	value = strings.TrimSuffix(value, "\r")
	// Step 4: trim trailing then leading whitespace. The shell's `[[:space:]]`
	// includes CR/FF/VT, so a value like `x\r ` — a CR followed by a space — collapses
	// to `x` there. The trailing-trim must match that cutset, not just space+tab, or
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
	return value, true
}

// ModelFromEnvelope parses the Claude Code Agent PreToolUse envelope, returning
// tool_input.resolvedModel if non-empty, else tool_input.model if non-empty, else "". A
// PreToolUse envelope carries tool_input.model; resolvedModel is a PostToolUse
// tool_response field. Reading it first is a defensive fallback against an envelope
// shape this event is not documented to send, not the contract the guard relies on. It
// returns a non-nil error ONLY when data is not valid JSON, mirroring the Python shim
// raising on unparseable stdin. Valid JSON with no matching field returns ("", nil).
func ModelFromEnvelope(data []byte) (model string, err error) {
	e, err := parseDelegation(data)
	if err != nil {
		return "", err
	}
	return e.model, nil
}

// forkSubagentType is the one delegation type the harness runs on the invoking session's
// model, discarding any declared one. It is an experimental, feature-gated value. Where
// fork mode is not deployed no envelope carries it. If it is renamed upstream, the
// comparison simply stops matching, and the guard keeps its other verdicts.
const forkSubagentType = "fork"

// delegation is the one reading of an Agent PreToolUse envelope the guard branches on: the
// declared model and whether the delegation is a fork.
type delegation struct {
	model  string
	isFork bool
}

// parseDelegation reads tool_input once. subagent_type is decoded separately from the
// raw bytes because it is untrusted text of an unpromised type. A number, object, or
// null there must leave the model reading intact rather than failing the whole envelope.
// Only a JSON string exactly equal to forkSubagentType may select the fork branch.
func parseDelegation(data []byte) (delegation, error) {
	var e struct {
		ToolInput struct {
			ResolvedModel string          `json:"resolvedModel"`
			Model         string          `json:"model"`
			SubagentType  json.RawMessage `json:"subagent_type"`
		} `json:"tool_input"`
	}
	if err := json.Unmarshal(data, &e); err != nil {
		return delegation{}, err
	}
	var subagentType string
	_ = json.Unmarshal(e.ToolInput.SubagentType, &subagentType)
	out := delegation{isFork: subagentType == forkSubagentType}
	// An all-whitespace field is blank for the non-empty test: a whitespace resolvedModel
	// falls back to model, and an all-blank envelope yields "". This closes the whitespace
	// model value at the parse boundary and folds the omitted/empty/whitespace cases into
	// one branch of AgentLineVerdict.
	switch {
	case strings.TrimSpace(e.ToolInput.ResolvedModel) != "":
		out.model = e.ToolInput.ResolvedModel
	case strings.TrimSpace(e.ToolInput.Model) != "":
		out.model = e.ToolInput.Model
	}
	return out, nil
}

// Harnesses is the closed set of harnesses the binding matrix covers, in the order
// diagnostics list them. A harness outside this set names no column: the matrix is closed
// so a typo cannot quietly create a phantom column that nothing grades.
//
// The set derives from the harness record: a row that binds no provider takes no cell, so
// the model-free row stays out. The record's order is therefore the diagnostic order.
var Harnesses = boundHarnesses()

// boundHarnesses names each record row that binds a provider, in record order.
func boundHarnesses() []string {
	names := make([]string, 0, len(harnesses.Rows))
	for _, row := range harnesses.Rows {
		if row.Providers == harnesses.NoProvider {
			continue
		}
		names = append(names, row.Harness)
	}
	return names
}

// harnessBinding returns the provider binding of the row named harness. The second result
// is false for a name the record does not hold.
func harnessBinding(harness string) (harnesses.Provider, bool) {
	row, ok := harnesses.Lookup(harness)
	if !ok {
		return "", false
	}
	return row.Providers, true
}

// harnessOf names the one record row that binds provider. A retired key family migrates to
// that row's column.
func harnessOf(provider harnesses.Provider) string {
	for _, row := range harnesses.Rows {
		if row.Providers == provider {
			return row.Harness
		}
	}
	return ""
}

// Tiers is the closed set of tier names BENCH_MODEL and the matrix share, in descending
// capability order.
var Tiers = []string{"top", "mid", "cheap"}

// Key returns the .bench/lines.env key that binds one harness's tier cell.
func Key(harness, tier string) string {
	return "BENCH_" + strings.ToUpper(harness) + "_" + strings.ToUpper(tier)
}

// KnownHarness reports whether name is one of Harnesses.
func KnownHarness(name string) bool {
	for _, h := range Harnesses {
		if h == name {
			return true
		}
	}
	return false
}

func knownTier(name string) bool {
	for _, t := range Tiers {
		if t == name {
			return true
		}
	}
	return false
}

// foreignKeyRe matches any BENCH_<WORD>_<TIER> assignment, whatever the harness segment
// says. Anything it matches that is not one of the keys the cell reader looks up is a
// foreign key. This includes the retired BENCH_TIER_* / BENCH_ALIAS_* schema. A binding
// carrying only retired keys reads as no binding, rather than as a legacy one.
var foreignKeyRe = regexp.MustCompile(`(?m)^[ \t]*(BENCH_[A-Za-z0-9]+_(?:TOP|MID|CHEAP))=`)

// matrixKey reports whether key is one of the keys ParseBinding reads cells from. The
// test is exact-string against Key's own output, rather than a case-folded reading of
// the harness segment. A segment matched case-insensitively would let a non-canonical
// spelling like BENCH_claude_TOP count as naming a known harness, while no cell lookup
// ever reads it. This leaves the key ignored by the reader and unreported by the
// foreign-key arm alike.
func matrixKey(key string) bool {
	for _, harness := range Harnesses {
		for _, tier := range Tiers {
			if Key(harness, tier) == key {
				return true
			}
		}
	}
	return false
}

// Binding is the parsed BENCH_<HARNESS>_<TIER> matrix: one column per harness in
// Harnesses, plus every harness-shaped key naming a harness outside that closed set.
type Binding struct {
	cells   map[string]cell
	foreign []string
}

// cell separates a key that is absent from one assigned an empty value: an assigned-empty
// cell declares its harness while leaving the column incomplete.
type cell struct {
	present bool
	value   string
}

// ParseBinding reads every matrix cell out of content and records any harness-shaped key
// naming a harness outside Harnesses.
func ParseBinding(content []byte) Binding {
	b := Binding{cells: make(map[string]cell, len(Harnesses)*len(Tiers))}
	for _, harness := range Harnesses {
		for _, tier := range Tiers {
			key := Key(harness, tier)
			value, present := tierValue(key, content)
			b.cells[key] = cell{present: present, value: value}
		}
	}
	seen := make(map[string]bool)
	for _, match := range foreignKeyRe.FindAllStringSubmatch(string(content), -1) {
		key := match[1]
		if matrixKey(key) || seen[key] {
			continue
		}
		seen[key] = true
		b.foreign = append(b.foreign, key)
	}
	return b
}

// Cell returns the model token bound to one harness's tier, or "" when that cell is unset.
func (b Binding) Cell(harness, tier string) string {
	return b.cells[Key(harness, tier)].value
}

// Column returns harness's three cells in Tiers order.
func (b Binding) Column(harness string) []string {
	column := make([]string, 0, len(Tiers))
	for _, tier := range Tiers {
		column = append(column, b.Cell(harness, tier))
	}
	return column
}

// Declared reports whether the binding mentions harness at all. Only a declared harness
// owes all three cells, so an unadopted harness leaves the matrix complete.
func (b Binding) Declared(harness string) bool {
	for _, tier := range Tiers {
		if b.cells[Key(harness, tier)].present {
			return true
		}
	}
	return false
}

// Complete reports whether every one of harness's three cells carries a value.
func (b Binding) Complete(harness string) bool {
	for _, tier := range Tiers {
		if b.Cell(harness, tier) == "" {
			return false
		}
	}
	return true
}

// Empty reports whether no harness binds any cell.
func (b Binding) Empty() bool {
	for _, harness := range Harnesses {
		for _, tier := range Tiers {
			if b.Cell(harness, tier) != "" {
				return false
			}
		}
	}
	return true
}

// ForeignKeys returns the harness-shaped keys naming a harness outside Harnesses, in the
// order they appear.
func (b Binding) ForeignKeys() []string {
	return b.foreign
}

// UnboundKeys returns the keys of harness's cells that bind no token, in Tiers order.
func (b Binding) UnboundKeys(harness string) []string {
	var keys []string
	for _, tier := range Tiers {
		if b.Cell(harness, tier) == "" {
			keys = append(keys, Key(harness, tier))
		}
	}
	return keys
}

// retiredFamilies pairs each retired key prefix with its replacement harness column.
// BENCH_TIER_* held one family's concrete ids,
// which is the codex column. BENCH_ALIAS_* held the tokens the dissolved alias concept
// projected, which are the claude column. These keys bind nothing: the matrix is a
// hard cut with no dual read. So the pairing is migration advice, and never a second
// reading of a binding.
var retiredFamilies = []struct{ prefix, harness string }{
	{"BENCH_TIER_", harnessOf(harnesses.OpenAI)},
	{"BENCH_ALIAS_", harnessOf(harnesses.Anthropic)},
}

// RetiredKeyPrefixes returns the retired schema's key stems in schema order. Consumers
// outside this package read the retired families off retiredFamilies above, rather than
// declaring their own copy. The doctor migration report walks that declaration for the
// rewrites it offers. The guidance sweep builds its prose matcher from this list, so the
// two cannot come to disagree about which schema is retired. Only the stems are
// exported: the harness each family migrates to is the report's business alone.
func RetiredKeyPrefixes() []string {
	prefixes := make([]string, 0, len(retiredFamilies))
	for _, family := range retiredFamilies {
		prefixes = append(prefixes, family.prefix)
	}
	return prefixes
}

// RetiredKeyRewrite is one retired assignment and the matrix assignment replacing it,
// both carrying Value. The migration moves a token between keys and changes no model
// choice.
type RetiredKeyRewrite struct {
	Retired     string
	Replacement string
	Value       string
}

// RetiredKeyRewrites returns every retired schema key assigned in content, in schema order,
// paired with the matrix key that replaces it.
func RetiredKeyRewrites(content []byte) []RetiredKeyRewrite {
	var out []RetiredKeyRewrite
	for _, family := range retiredFamilies {
		for _, tier := range Tiers {
			retired := family.prefix + strings.ToUpper(tier)
			value, present := tierValue(retired, content)
			if !present {
				continue
			}
			out = append(out, RetiredKeyRewrite{
				Retired:     retired,
				Replacement: Key(family.harness, tier),
				Value:       value,
			})
		}
	}
	return out
}

// CellFault reports why value cannot serve as harness's bound cell, or "" when it can.
// Every cell is an opaque safe token. A harness that binds any provider carries a
// provider-qualified namespace, so that requirement is a rule on that row's own cells
// rather than a filter applied to whatever a resolution returns. A row that binds one
// provider names the provider already, so a bare id serves.
func CellFault(harness, value string) string {
	if !modelid.SafeToken(value) {
		return "is not a safe model token"
	}
	if binding, ok := harnessBinding(harness); ok && binding == harnesses.AnyProvider && !isProviderModel(value) {
		return "is not provider-qualified (opencode model ids are provider/model)"
	}
	return ""
}

// Source is the read state of .bench/lines.env. Unreadable is distinct from absent: an
// absent binding means the repo is unrouted. A binding that fails to read is a corrupt
// oracle, and folding the two together would silently disable enforcement.
type Source struct {
	Path       string
	Exists     bool
	Unreadable bool
	Content    []byte
}

// state is the one reading of a Source both verdicts branch on. So the fail-open guard
// and the fail-closed resolver never disagree about what the binding says.
type state int

const (
	stateAbsent     state = iota // no .bench/lines.env: the repo is unrouted
	stateUnreadable              // present but unreadable
	stateForeign                 // carries a harness-shaped key naming no known harness
	stateUnbound                 // well-formed but binds no cell at all
	stateBound                   // at least one cell is bound
)

func classify(src Source) (state, Binding) {
	// Unreadable is tested first: a binding whose bytes are unavailable must never fall
	// through to the unrouted branch, whatever else the caller reports about it.
	if src.Unreadable {
		return stateUnreadable, Binding{}
	}
	if !src.Exists {
		return stateAbsent, Binding{}
	}
	b := ParseBinding(src.Content)
	switch {
	case len(b.ForeignKeys()) > 0:
		return stateForeign, b
	case b.Empty():
		return stateUnbound, b
	}
	return stateBound, b
}

// enforceable reports whether the asking harness has a column the guard can hold a
// delegation to. It is the single condition every deny in AgentLineVerdict is gated on.
// Outside it there is no bound tier, so nothing can be named as the line the delegation
// should have carried. The guard fails open rather than bricking a repo that never
// opted into line enforcement.
func enforceable(st state, b Binding, harness string) bool {
	return st == stateBound && b.Complete(harness)
}

func warn(s string) string {
	return "WARNING: check-agent-line: " + s + " — allowing delegation."
}

func unknownHarness(cmd, harness string) string {
	return cmd + ": unknown harness '" + harness + "'; the binding matrix covers " +
		strings.Join(Harnesses, ", ")
}

func foreignKeyText(b Binding) string {
	return "binds unknown harness key " + strings.Join(b.ForeignKeys(), ", ") +
		"; the binding matrix covers " + strings.Join(Harnesses, ", ")
}

// AgentLineVerdict is the pure agent-line verdict for a delegation asked from harness.
// Every degraded branch is fail-OPEN, exit 0, with a one-line stderr warning. Only a
// present model bound nowhere in the matrix denies, exit 2. An unknown harness is a
// wiring error, exit 1, the shim treats as a core error.
//
// Enforcement is permissive across every bound cell in the whole matrix. A Claude
// session may legitimately name a Codex delegate's tier. The denial's advice names
// only the asking harness's own tokens. The returned stderr carries no trailing
// newline; the caller adds one.
func AgentLineVerdict(stdin []byte, harness string, src Source) (exitCode int, stderr string) {
	if !KnownHarness(harness) {
		return 1, unknownHarness("check-agent-line", harness)
	}
	e, err := parseDelegation(stdin)
	if err != nil {
		return 0, warn("stdin is not parseable as JSON")
	}
	model := e.model
	st, b := classify(src)
	if e.isFork {
		// A fork runs on this session's model whatever it declares. The binding is not
		// what settles the declaration: it is a claim the harness discards and the guard
		// cannot check. An omitted model is the honest signal for behavior no delegation
		// can avoid.
		//
		// Neither verdict can name the session's own model. Only SessionStart is
		// documented to receive one, and it is not guaranteed there. A mid-session switch
		// is re-reported by no hook event.
		//
		// The deny waits for a column to enforce: with no bound tier there is nothing to
		// escalate off. A repo that never opted into line enforcement keeps its
		// delegations. The warning does not wait; it is a warning either way. Withholding
		// it would cost the operator a true statement about inheritance.
		if model != "" && enforceable(st, b, harness) {
			return 2, "DENIED: delegation model '" + model + "' is declared on a fork, which runs on this session's " +
				"model and ignores the declaration — the guard cannot verify a line the harness will not honor. " +
				"Re-delegate the fork with no model field to inherit this session's line, or spawn a non-fork " +
				"delegate on a bound tier token."
		}
		return 0, warn("a fork delegation inherits this session's model, which no hook event reports")
	}
	if model == "" {
		// The one degraded branch that is also the attack path the guard exists for is an
		// omitted or empty model. It inherits the invoking session's model, the silent
		// escalation invariant #2 forbids. Deny it. Every other missing-model branch keeps
		// the fail-open rim — there is no column to enforce, and a broken guard must never
		// brick delegation.
		if enforceable(st, b, harness) {
			return 2, "DENIED: the delegation envelope has a missing or empty model field — an " +
				"omitted model silently inherits this session's model, which invariant #2 forbids. " +
				"Pass a bound tier token from .bench/lines.env; " + describeColumn(harness, b) +
				" (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier."
		}
		return 0, warn("no resolvedModel/model field in tool_input")
	}
	switch st {
	case stateAbsent:
		return 0, warn("no .bench/lines.env at repo root")
	case stateUnreadable:
		return 0, warn(".bench/lines.env is present but unreadable")
	case stateForeign:
		return 0, warn(".bench/lines.env " + foreignKeyText(b))
	case stateUnbound:
		return 0, warn("no BENCH_<HARNESS>_<TIER> cell is bound in .bench/lines.env")
	}
	if !b.Complete(harness) {
		return 0, warn("the " + harness + " column is incomplete in .bench/lines.env")
	}
	for _, known := range Harnesses {
		for _, tier := range Tiers {
			if value := b.Cell(known, tier); value != "" && model == value {
				return 0, ""
			}
		}
	}
	return 2, "DENIED: delegation model '" + model + "' is not a bound tier; " + describeColumn(harness, b) +
		" (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding."
}

// describeColumn formats the asking harness's own three bound tokens, the only tokens a
// denial advertises. The rest of the matrix stays out of the message: enforcement is
// permissive across it. A recovery instruction is executable only where it names what
// this harness can actually pass. The whole matrix is human-readable in the profile.
func describeColumn(harness string, b Binding) string {
	column := b.Column(harness)
	parts := make([]string, 0, len(Tiers))
	for i, tier := range Tiers {
		parts = append(parts, tier+"="+column[i])
	}
	return "harness " + harness + " binds " + strings.Join(parts, " ")
}

// ResolveModelVerdict is the shift adapters' model resolution for harness: BENCH_MODEL
// names a tier and the harness's own column names the model. Unlike the fail-open agent
// guard, an unusable binding refuses, exit 1, so an adapter never launches unguarded.
// The one exception is a binding that declares no cell at all, which degrades to the
// BENCH_MODEL passthrough an unadopted repo relies on. The returned stderr carries no
// trailing newline.
func ResolveModelVerdict(harness, benchModel string, benchModelSet bool, src Source) (model string, exitCode int, stderr string) {
	if !KnownHarness(harness) {
		return "", 1, unknownHarness("bench resolve-model", harness)
	}
	st, b := classify(src)
	switch st {
	case stateAbsent:
		// Unrouted: explicit BENCH_MODEL wins; absent -> "".
		return benchModel, 0, ""
	case stateUnreadable:
		return "", 1, "bench shift: cannot read the tier binding at " + src.Path + " — refusing to run unguarded"
	case stateForeign:
		return "", 1, "bench shift: " + src.Path + " " + foreignKeyText(b)
	case stateUnbound:
		return benchModel, 0, "WARNING: bench adapter: no BENCH_<HARNESS>_<TIER> cell is bound in " +
			src.Path + " — ignoring the binding and falling back to BENCH_MODEL."
	}
	if !b.Complete(harness) {
		keys := make([]string, 0, len(Tiers))
		for _, tier := range Tiers {
			keys = append(keys, Key(harness, tier))
		}
		return "", 1, "bench shift: the " + harness + " column is unbound in " + src.Path + "; bind " +
			strings.Join(keys, ", ") + " before running " + harness + "."
	}
	for _, tier := range Tiers {
		value := b.Cell(harness, tier)
		if fault := CellFault(harness, value); fault != "" {
			return "", 1, "bench shift: " + Key(harness, tier) + "='" + value + "' " + fault + " in " + src.Path
		}
	}
	if !benchModelSet || benchModel == "" {
		return "", 1, "bench shift in a routed repo requires a declared line: set BENCH_MODEL to one of " +
			strings.Join(Tiers, ", ")
	}
	if !knownTier(benchModel) {
		return "", 1, "bench shift: BENCH_MODEL='" + benchModel + "' is not a tier; set it to one of " +
			strings.Join(Tiers, ", ")
	}
	return b.Cell(harness, benchModel), 0, ""
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
