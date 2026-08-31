// Package harnesses is the one record of every harness Bench knows. One row names one
// harness. A row carries the providers the harness binds, the phase invocation form, the
// hook config with its wired events, the deny-capable delegation verdict, the headless
// adapter, the mechanics as cells, and the run measures as cells.
//
// Every reader derives from this record: the tier-binding matrix, the route's prefix
// table, the guard wiring reader, and the conformance harness loops. A new harness
// therefore lands as one row.
//
// The file imports only the standard library, so a low package can compose it without an
// import cycle.
//
// A cell holds yes, no, or unknown. A yes or no cell names the source that was read and
// the date of that read. An unknown cell names neither. The record holds only a fact the
// tree itself records today. Every other cell is unknown, and a reviewed edit fills it.
package harnesses

// Schema is the record's shape version. A consumer prints it, so it can name the shape it
// read.
const Schema = 1

// Value is the closed set of cell values.
type Value string

// The three cell values. Unknown means the tree records no fact, not that the harness
// lacks the mechanic.
const (
	Yes     Value = "yes"
	No      Value = "no"
	Unknown Value = "unknown"
)

// Values is the closed enum in its canonical order.
var Values = []Value{Yes, No, Unknown}

// Provider is the closed set of provider bindings a row can hold. A harness that accepts
// any provider binds Any. The model-free row binds NoProvider.
type Provider string

// The four provider bindings.
const (
	OpenAI      Provider = "openai"
	Anthropic   Provider = "anthropic"
	AnyProvider Provider = "any"
	NoProvider  Provider = "none"
)

// Providers is the closed provider enum in its canonical order.
var Providers = []Provider{OpenAI, Anthropic, AnyProvider, NoProvider}

// Cell is one graded claim. Checked is an ISO date (YYYY-MM-DD) and names the day the
// source was read. Source and Checked are empty exactly when Value is Unknown.
type Cell struct {
	Value   Value
	Source  string
	Checked string
}

// Measure is one declared measure of a harness run. Supplier names the source that will
// read the measure. Value stays Unknown until that supplier ships, because the record holds
// only a fact the tree records today.
type Measure struct {
	Value    Value
	Supplier string
}

// Row is one harness. HookConfig and Headless are repository-relative paths, and each is
// empty when the harness has none. PhaseForm is the phase invocation prefix, and it is
// empty when the harness has no phase surface. Mechanics holds one cell for every name in
// Mechanics order, and Measures holds one cell for every name in Measures order, so a
// reader can print both sets in a stable order.
type Row struct {
	Harness         string
	Providers       Provider
	PhaseForm       string
	HookConfig      string
	HookEvents      []string
	DelegationGuard Cell
	Headless        string
	Mechanics       map[string]Cell
	Measures        map[string]Measure
}

// The twelve mechanic names. Each row holds a cell for each one.
const (
	MechanicSteering       = "steering during an active turn"
	MechanicQuestions      = "structured user questions"
	MechanicPermissions    = "tool-permission controls"
	MechanicHooks          = "hooks"
	MechanicMCP            = "MCP support"
	MechanicSubagents      = "subagent support"
	MechanicIsolation      = "subagent isolation"
	MechanicEffort         = "effort selection"
	MechanicPersistence    = "persistent tasks"
	MechanicResume         = "resume and recovery"
	MechanicStructuredExit = "structured output and exit status"
	MechanicHeadless       = "headless execution"
)

// Mechanics is the mechanic order the detail view prints.
var Mechanics = []string{
	MechanicSteering,
	MechanicQuestions,
	MechanicPermissions,
	MechanicHooks,
	MechanicMCP,
	MechanicSubagents,
	MechanicIsolation,
	MechanicEffort,
	MechanicPersistence,
	MechanicResume,
	MechanicStructuredExit,
	MechanicHeadless,
}

// The four measure names. Each row holds a measure cell for each one.
const (
	MeasureTokens    = "tokens"
	MeasureToolCalls = "tool calls"
	MeasureReadPaths = "Read paths"
	MeasureTurns     = "turns"
)

// Measures is the measure order the detail view prints.
var Measures = []string{
	MeasureTokens,
	MeasureToolCalls,
	MeasureReadPaths,
	MeasureTurns,
}

// measureSupplier names the source that will read every harness-owned measure. One reader
// serves all four measures, so the name is spelled once. The cells stay unknown until that
// reader ships.
const measureSupplier = "FT204 harness transcript reader"

// The sources the initial cells name. Each one is a file in this tree.
const (
	srcClaudeHooks = ".claude/settings.json"
	srcCodexHooks  = ".codex/hooks.json"
	srcAdapters    = ".bench/adapters/"
	srcReference   = ".bench/BENCH-reference.md"
)

// adapterPath builds the repository-relative path of one harness's headless adapter. Every
// adapter path and every headless cell source comes from here, so the directory is spelled
// once.
func adapterPath(name string) string {
	return srcAdapters + name
}

// The dates the initial cells carry. codexGuardChecked is the Hook Layers verdict's own
// date, because that bullet records the read of the upstream Codex hooks docs.
const (
	treeChecked       = "2026-08-26"
	codexGuardChecked = "2026-07-11"
)

// unknownCells builds the twelve cells with no fact recorded. A caller then overwrites the
// cells the tree does record.
func unknownCells() map[string]Cell {
	cells := make(map[string]Cell, len(Mechanics))
	for _, name := range Mechanics {
		cells[name] = Cell{Value: Unknown}
	}
	return cells
}

// unknownMeasures builds the four measure cells with no value recorded. Each cell names its
// supplier, because the supplier is decided today and the value is not.
func unknownMeasures() map[string]Measure {
	cells := make(map[string]Measure, len(Measures))
	for _, name := range Measures {
		cells[name] = Measure{Value: Unknown, Supplier: measureSupplier}
	}
	return cells
}

// mechanics builds one row's cells from the unknown baseline plus the recorded facts.
func mechanics(recorded map[string]Cell) map[string]Cell {
	cells := unknownCells()
	for name, cell := range recorded {
		cells[name] = cell
	}
	return cells
}

// Rows is the record. The order is codex, claude, opencode, none, so every diagnostic that
// walks the record keeps its wording.
//
// HookEvents names each wired hook group. A config that wires one event under two
// matchers therefore names that event twice, qualified by its matcher, because the two
// groups run different scripts.
//
// The effort cells stay unknown. The reference's adapter rule states Bench's own entry
// surface, not the facility a harness offers, so it grades no row.
var Rows = []Row{
	{
		Harness:    "codex",
		Providers:  OpenAI,
		PhaseForm:  "$bench-",
		HookConfig: srcCodexHooks,
		HookEvents: []string{"SessionStart", "Stop", "PreToolUse:Bash"},
		DelegationGuard: Cell{
			Value: No,
			// The Hook Layers bullet records the read: a spawn never surfaces as a
			// matchable tool name on a deny-capable event, and SubagentStart cannot
			// stop the spawn.
			Source:  srcReference + " Hook Layers, the agent-line bullet (Codex hooks docs)",
			Checked: codexGuardChecked,
		},
		Headless: adapterPath("codex"),
		Mechanics: mechanics(map[string]Cell{
			MechanicHooks:    {Value: Yes, Source: srcCodexHooks, Checked: treeChecked},
			MechanicHeadless: {Value: Yes, Source: adapterPath("codex"), Checked: treeChecked},
		}),
		Measures: unknownMeasures(),
	},
	{
		Harness:    "claude",
		Providers:  Anthropic,
		PhaseForm:  "/bench-",
		HookConfig: srcClaudeHooks,
		HookEvents: []string{
			"WorktreeCreate",
			"WorktreeRemove",
			"SessionStart",
			"Stop",
			"PreToolUse:Bash",
			"PreToolUse:Agent",
		},
		DelegationGuard: Cell{
			Value:   Yes,
			Source:  srcClaudeHooks + " PreToolUse Agent matcher runs .bench/hooks/check-agent-line.sh",
			Checked: treeChecked,
		},
		Headless: adapterPath("claude"),
		Mechanics: mechanics(map[string]Cell{
			MechanicHooks:    {Value: Yes, Source: srcClaudeHooks, Checked: treeChecked},
			MechanicHeadless: {Value: Yes, Source: adapterPath("claude"), Checked: treeChecked},
		}),
		Measures: unknownMeasures(),
	},
	{
		Harness:         "opencode",
		Providers:       AnyProvider,
		PhaseForm:       "",
		HookConfig:      "",
		HookEvents:      nil,
		DelegationGuard: Cell{Value: Unknown},
		Headless:        adapterPath("opencode"),
		Mechanics: mechanics(map[string]Cell{
			MechanicHeadless: {Value: Yes, Source: adapterPath("opencode"), Checked: treeChecked},
		}),
		Measures: unknownMeasures(),
	},
	{
		Harness:    "none",
		Providers:  NoProvider,
		PhaseForm:  "",
		HookConfig: "",
		HookEvents: nil,
		DelegationGuard: Cell{
			Value:   No,
			Source:  srcAdapters + " names no none entry, so the model-free path runs no agent",
			Checked: treeChecked,
		},
		Headless: "",
		Mechanics: mechanics(map[string]Cell{
			MechanicHooks:    {Value: No, Source: srcAdapters + " names no none entry, and no config names none", Checked: treeChecked},
			MechanicHeadless: {Value: No, Source: srcAdapters + " names no none entry", Checked: treeChecked},
		}),
		Measures: unknownMeasures(),
	},
}

// Lookup returns the row named name. The second result is false for a name the record does
// not hold, so a caller can refuse an unknown harness rather than print an empty state.
func Lookup(name string) (Row, bool) {
	for _, row := range Rows {
		if row.Harness == name {
			return row, true
		}
	}
	return Row{}, false
}
