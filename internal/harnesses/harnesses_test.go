package harnesses

import (
	"fmt"
	"testing"
	"time"
)

// want holds the facts the walk grades per row. The walk takes the rows as an argument, so
// the mutation test can run the same grader over a changed copy.
type want struct {
	harness   string
	providers Provider
	phaseForm string
	hooks     string
	events    int
	guard     Value
	headless  string
	// headlessCell is the row's headless-execution cell value. The `none` row records no
	// adapter, so its cell must stay no.
	headlessCell Value
}

var expected = []want{
	{"codex", OpenAI, "$bench-", ".codex/hooks.json", 3, No, ".bench/adapters/codex", Yes},
	{"claude", Anthropic, "/bench-", ".claude/settings.json", 6, Yes, ".bench/adapters/claude", Yes},
	{"opencode", AnyProvider, "", "", 0, Unknown, ".bench/adapters/opencode", Yes},
	{"none", NoProvider, "", "", 0, No, "", No},
}

// wantMechanics is an independent expectation of the twelve mechanic names and their
// order. The record's own Mechanics slice cannot grade itself: a dropped name shrinks the
// slice and every row's cell map together, and a count taken from the slice stays equal.
// This list is therefore the only thing that reds a dropped mechanic.
var wantMechanics = []string{
	"steering during an active turn",
	"structured user questions",
	"tool-permission controls",
	"hooks",
	"MCP support",
	"subagent support",
	"subagent isolation",
	"effort selection",
	"persistent tasks",
	"resume and recovery",
	"structured output and exit status",
	"headless execution",
}

// wantMeasures is an independent expectation of the four measure names and their order.
// The record's own Measures slice cannot grade itself: a dropped name shrinks the slice and
// every row's measure map together. This list is therefore the only thing that reds a
// dropped measure.
var wantMeasures = []string{
	"tokens",
	"tool calls",
	"Read paths",
	"turns",
}

// walk grades one record against the expected rows and the cell rules. It returns one
// message per fault, so a mutation names what it broke.
func walk(rows []Row, today time.Time) []string {
	var faults []string
	if len(Mechanics) != len(wantMechanics) {
		faults = append(faults, fmt.Sprintf("the record names %d mechanics, want %d", len(Mechanics), len(wantMechanics)))
	}
	for i, name := range wantMechanics {
		if i >= len(Mechanics) {
			faults = append(faults, fmt.Sprintf("the record names no mechanic %q", name))
			continue
		}
		if Mechanics[i] != name {
			faults = append(faults, fmt.Sprintf("mechanic %d is %q, want %q", i, Mechanics[i], name))
		}
	}
	if len(Measures) != len(wantMeasures) {
		faults = append(faults, fmt.Sprintf("the record names %d measures, want %d", len(Measures), len(wantMeasures)))
	}
	for i, name := range wantMeasures {
		if i >= len(Measures) {
			faults = append(faults, fmt.Sprintf("the record names no measure %q", name))
			continue
		}
		if Measures[i] != name {
			faults = append(faults, fmt.Sprintf("measure %d is %q, want %q", i, Measures[i], name))
		}
	}
	if len(rows) != len(expected) {
		return append(faults, fmt.Sprintf("record holds %d rows, want %d", len(rows), len(expected)))
	}
	for i, row := range rows {
		w := expected[i]
		if row.Harness != w.harness {
			faults = append(faults, fmt.Sprintf("row %d is %q, want %q", i, row.Harness, w.harness))
			continue
		}
		if row.Providers != w.providers {
			faults = append(faults, fmt.Sprintf("%s binds %q, want %q", row.Harness, row.Providers, w.providers))
		}
		if !validProvider(row.Providers) {
			faults = append(faults, fmt.Sprintf("%s providers %q is outside the enum", row.Harness, row.Providers))
		}
		if row.PhaseForm != w.phaseForm {
			faults = append(faults, fmt.Sprintf("%s phase form is %q, want %q", row.Harness, row.PhaseForm, w.phaseForm))
		}
		if row.HookConfig != w.hooks {
			faults = append(faults, fmt.Sprintf("%s hook config is %q, want %q", row.Harness, row.HookConfig, w.hooks))
		}
		if len(row.HookEvents) != w.events {
			faults = append(faults, fmt.Sprintf("%s wires %d events, want %d", row.Harness, len(row.HookEvents), w.events))
		}
		for _, event := range row.HookEvents {
			if event == "" {
				faults = append(faults, fmt.Sprintf("%s names an empty hook event", row.Harness))
			}
		}
		if row.HookConfig == "" && len(row.HookEvents) != 0 {
			faults = append(faults, fmt.Sprintf("%s names events with no hook config", row.Harness))
		}
		if row.DelegationGuard.Value != w.guard {
			faults = append(faults, fmt.Sprintf("%s delegation guard is %q, want %q", row.Harness, row.DelegationGuard.Value, w.guard))
		}
		faults = append(faults, cellFaults(row.Harness, "delegation_guard", row.DelegationGuard, today)...)
		if row.Headless != w.headless {
			faults = append(faults, fmt.Sprintf("%s headless adapter is %q, want %q", row.Harness, row.Headless, w.headless))
		}
		if cell := row.Mechanics[MechanicHeadless]; cell.Value != w.headlessCell {
			faults = append(faults, fmt.Sprintf("%s headless execution is %q, want %q", row.Harness, cell.Value, w.headlessCell))
		}
		if len(row.Mechanics) != len(Mechanics) {
			faults = append(faults, fmt.Sprintf("%s holds %d mechanics, want %d", row.Harness, len(row.Mechanics), len(Mechanics)))
		}
		for _, name := range Mechanics {
			cell, ok := row.Mechanics[name]
			if !ok {
				faults = append(faults, fmt.Sprintf("%s holds no %q cell", row.Harness, name))
				continue
			}
			faults = append(faults, cellFaults(row.Harness, name, cell, today)...)
		}
		if len(row.Measures) != len(wantMeasures) {
			faults = append(faults, fmt.Sprintf("%s holds %d measures, want %d", row.Harness, len(row.Measures), len(wantMeasures)))
		}
		for _, name := range wantMeasures {
			cell, ok := row.Measures[name]
			if !ok {
				faults = append(faults, fmt.Sprintf("%s holds no %q measure", row.Harness, name))
				continue
			}
			if !validValue(cell.Value) {
				faults = append(faults, fmt.Sprintf("%s measure %s value %q is outside the enum", row.Harness, name, cell.Value))
			}
			if cell.Value != Unknown {
				faults = append(faults, fmt.Sprintf("%s measure %s is %q, want unknown until a supplier ships", row.Harness, name, cell.Value))
			}
			if cell.Supplier == "" {
				faults = append(faults, fmt.Sprintf("%s measure %s names no supplier", row.Harness, name))
			}
		}
	}
	return faults
}

// cellFaults grades one cell: the value is inside the enum, a yes or no cell names a
// source and an ISO date that is not later than today, and an unknown cell names neither.
func cellFaults(harness, field string, cell Cell, today time.Time) []string {
	var faults []string
	if !validValue(cell.Value) {
		return []string{fmt.Sprintf("%s %s value %q is outside the enum", harness, field, cell.Value)}
	}
	if cell.Value == Unknown {
		if cell.Source != "" || cell.Checked != "" {
			faults = append(faults, fmt.Sprintf("%s %s is unknown but names a source or a date", harness, field))
		}
		return faults
	}
	if cell.Source == "" {
		faults = append(faults, fmt.Sprintf("%s %s is %q with no source", harness, field, cell.Value))
	}
	checked, err := time.Parse("2006-01-02", cell.Checked)
	if err != nil {
		faults = append(faults, fmt.Sprintf("%s %s checked %q is not an ISO date", harness, field, cell.Checked))
		return faults
	}
	if checked.After(today) {
		faults = append(faults, fmt.Sprintf("%s %s checked %q is later than today", harness, field, cell.Checked))
	}
	return faults
}

func validValue(v Value) bool {
	for _, known := range Values {
		if v == known {
			return true
		}
	}
	return false
}

func validProvider(p Provider) bool {
	for _, known := range Providers {
		if p == known {
			return true
		}
	}
	return false
}

// copyRows deep-copies the record so a mutation never leaks into another test.
func copyRows() []Row {
	out := make([]Row, len(Rows))
	for i, row := range Rows {
		row.HookEvents = append([]string(nil), row.HookEvents...)
		cells := make(map[string]Cell, len(row.Mechanics))
		for name, cell := range row.Mechanics {
			cells[name] = cell
		}
		row.Mechanics = cells
		measures := make(map[string]Measure, len(row.Measures))
		for name, cell := range row.Measures {
			measures[name] = cell
		}
		row.Measures = measures
		out[i] = row
	}
	return out
}

func TestRecordWalk(t *testing.T) {
	if Schema != 1 {
		t.Fatalf("Schema is %d, want 1", Schema)
	}
	if faults := walk(Rows, time.Now()); len(faults) != 0 {
		t.Fatalf("record walk: want no faults, got %v", faults)
	}
	for _, w := range expected {
		row, ok := Lookup(w.harness)
		if !ok {
			t.Fatalf("Lookup(%q): want a row", w.harness)
		}
		if row.Harness != w.harness {
			t.Fatalf("Lookup(%q) returned %q", w.harness, row.Harness)
		}
	}
	if _, ok := Lookup("cursor"); ok {
		t.Fatal("Lookup(\"cursor\"): want no row")
	}
}

// TestRecordWalkBites is the recorded bite proof. Each mutation is one class of rot the
// walk must red: an empty cell, a value outside the enum, and a check date in the future.
func TestRecordWalkBites(t *testing.T) {
	today := time.Now()
	cases := []struct {
		name   string
		mutate func(rows []Row)
	}{
		{"empty cell", func(rows []Row) {
			rows[0].Mechanics[MechanicHooks] = Cell{}
		}},
		{"dropped cell", func(rows []Row) {
			delete(rows[1].Mechanics, MechanicHeadless)
		}},
		{"value outside the enum", func(rows []Row) {
			rows[1].Mechanics[MechanicHooks] = Cell{Value: "maybe", Source: ".claude/settings.json", Checked: "2026-08-26"}
		}},
		{"future checked date", func(rows []Row) {
			cell := rows[0].Mechanics[MechanicHeadless]
			cell.Checked = today.AddDate(1, 0, 0).Format("2006-01-02")
			rows[0].Mechanics[MechanicHeadless] = cell
		}},
		{"unknown cell with a source", func(rows []Row) {
			rows[2].Mechanics[MechanicMCP] = Cell{Value: Unknown, Source: "a guess", Checked: "2026-08-26"}
		}},
		{"swapped rows", func(rows []Row) {
			rows[0], rows[1] = rows[1], rows[0]
		}},
		{"dropped hook event", func(rows []Row) {
			rows[1].HookEvents = rows[1].HookEvents[1:]
		}},
		{"guard on codex", func(rows []Row) {
			rows[0].DelegationGuard = Cell{Value: Yes, Source: ".codex/hooks.json", Checked: "2026-08-26"}
		}},
		{"none row claims headless execution", func(rows []Row) {
			rows[3].Mechanics[MechanicHeadless] = Cell{Value: Yes, Source: ".bench/adapters/none", Checked: "2026-08-26"}
		}},
		{"dropped measure cell", func(rows []Row) {
			delete(rows[1].Measures, MeasureReadPaths)
		}},
		{"measure with no supplier", func(rows []Row) {
			rows[0].Measures[MeasureTokens] = Measure{Value: Unknown}
		}},
		{"measure claiming a value with no supplier shipped", func(rows []Row) {
			rows[2].Measures[MeasureTurns] = Measure{Value: Yes, Supplier: "a guess"}
		}},
		{"none row names an adapter", func(rows []Row) {
			rows[3].Headless = ".bench/adapters/none"
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rows := copyRows()
			tc.mutate(rows)
			if faults := walk(rows, today); len(faults) == 0 {
				t.Fatalf("%s: want a fault, got none", tc.name)
			}
		})
	}
}
