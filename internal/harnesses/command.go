package harnesses

import (
	"fmt"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/harnesstranscript"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// The record view's two flags. Neither is Required, because the compiled views take no flag
// at all; the pair is enforced together below so a half-typed record call names the flag it
// lacks instead of silently printing the compiled view.
const (
	flagRecord = "--record"
	flagFormat = "--format"
)

// grammar is the declared argument shape usage.Parse enforces for `bench harnesses`. The
// verb takes at most one harness name, so a second positional is a usage error rather than
// a silently ignored argument.
var grammar = usage.Grammar{
	Cmd:  "bench harnesses",
	Help: "usage: bench harnesses [<harness>] | bench harnesses <harness> --record <path> --format <source-id>",
	Flags: []usage.Flag{
		{Name: flagRecord, HasValue: true, NoEmptyValue: true},
		{Name: flagFormat, HasValue: true, NoEmptyValue: true},
	},
	MinArgs: 0,
	MaxArgs: 1,
}

// overviewFields is the bare projection's schema. One row names one harness, and the row
// carries the facts a reader needs to choose a harness: what it binds, how a phase is
// invoked, where its hooks live, whether its delegation guard denies, and where its
// headless adapter is.
var overviewFields = []string{"harness", "provider", "phase_form", "hooks", "delegation_guard", "headless", "checked"}

// detailFields is the one-harness projection's schema. Every graded claim on the row
// renders as one cell with the source that was read and the date of that read.
var detailFields = []string{"field", "value", "source", "checked"}

// detailDelegationGuard is the field name the delegation-guard cell renders under. The
// cell is a Row field rather than a mechanic, so the detail view names it explicitly and
// prints it after the twelve mechanics.
const detailDelegationGuard = "delegation_guard"

// measureFields is the measure projection's schema. A measure cell names its supplier
// rather than a source that was read, so the measures render as their own table instead of
// borrowing the cells table's source column.
var measureFields = []string{"measure", "value", "supplier"}

// recordFields identifies the inspected record. A reader who wants to reproduce an
// observation needs the exact bytes, the mapping that read them, and the span they cover.
var recordFields = []string{"digest", "source_id", "harness", "interval"}

// observationFields is the record view's schema. Availability and boundary travel with every
// value, so an unknown measure and a counted one can never read alike.
var observationFields = []string{"metric", "value", "unit", "availability", "source", "boundary"}

// Command implements `bench harnesses [<harness>]`, the record's AXI projection. Bare, it
// prints the schema version and one row per harness, the model-free `none` row included,
// because a projection that hides the degraded path hides the case a reader most needs to
// see. With one harness name it prints that row's cells, each with its source and date.
//
// An unknown name is a usage error at exit 2, not an empty table. The record is closed, so
// an empty table would read as a definitive "this harness has no cells" rather than "no
// such harness".
//
// The verb reads no disk and needs no repository, because the record compiles in. So it
// has no not-in-a-repo refusal.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	record, wantsRecord := parsed.Flags[flagRecord]
	format, wantsFormat := parsed.Flags[flagFormat]
	if wantsRecord != wantsFormat {
		missing := flagFormat
		if wantsFormat {
			missing = flagRecord
		}
		return toon.MissingArg(grammar.Cmd, missing) + "\n", 2
	}
	if wantsRecord {
		return commandRecord(parsed.Positionals, record, format)
	}
	if len(parsed.Positionals) == 0 {
		return render("harnesses", overviewFields, overviewRows())
	}
	row, ok := Lookup(parsed.Positionals[0])
	if !ok {
		return toon.Usage(grammar.Cmd, parsed.Positionals[0]) + "\n", 2
	}
	return renderDetail(row)
}

// commandRecord answers the opt-in record view. The named harness is required, because an
// observation belongs to the harness that produced the record and an unnamed one would
// publish a measure with no owner.
//
// An unpinned source ID is refused before the path is touched: no mapping exists to read the
// bytes with, so the view reports every measure unknown at exit 1 rather than reading a file
// it cannot interpret.
func commandRecord(positionals []string, path, format string) (string, int) {
	if len(positionals) == 0 {
		return toon.MissingArg(grammar.Cmd, "harness") + "\n", 2
	}
	row, ok := Lookup(positionals[0])
	if !ok {
		return toon.Usage(grammar.Cmd, positionals[0]) + "\n", 2
	}
	if !harnesstranscript.Supported(format) {
		return renderRecord(row, harnesstranscript.Unsupported(format), 1)
	}
	observed, failure := harnesstranscript.Read(path, format)
	if failure.State != "" {
		return toon.RecordError(path, bounds.FileState(failure.State), failure.Reason) + "\n", 1
	}
	return renderRecord(row, observed, 0)
}

// renderRecord emits the record identity, then the observations. code is the exit the view
// carries when the tables themselves render: a refused source still prints its unknown rows,
// because an agent needs to see which dimensions it did not get.
func renderRecord(row Row, observed harnesstranscript.Record, code int) (string, int) {
	identity, err := toon.Table("record", recordFields, [][]string{
		{observed.Digest, observed.SourceID, row.Harness, observed.Interval},
	})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	// The value column is genuinely mixed: a counted measure renders as a bare integer, and
	// an unknown one renders as an empty cell. A typed row keeps the integer bare instead of
	// quoting a numeric-looking string.
	rows := make([][]any, 0, len(observed.Observations))
	for _, observation := range observed.Observations {
		rows = append(rows, []any{
			observation.Metric,
			observation.Cell(),
			observation.Unit,
			string(observation.Availability),
			observation.Source,
			observation.Boundary,
		})
	}
	observations, err := toon.TableTyped("observations", observationFields, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out, envelopeCode := envelope(identity + observations)
	if envelopeCode != 0 {
		return out, envelopeCode
	}
	return out, code
}

// render emits the schema line, the table, and the terminal help envelope. Both projections
// are terminal reads, so the disclosure block is honestly empty: the record holds no state
// a follow-on command could repair.
func render(name string, fields []string, rows [][]string) (string, int) {
	table, err := toon.Table(name, fields, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return envelope(table)
}

// renderDetail emits one harness's two tables: the graded cells, then the declared
// measures. The measures follow the cells because a measure is a promise about a future
// read, and the graded facts come first.
func renderDetail(row Row) (string, int) {
	cells, err := toon.Table("cells", detailFields, detailRows(row))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	measures, err := toon.Table("measures", measureFields, measureRows(row))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return envelope(cells + measures)
}

// envelope wraps rendered tables in the schema line and the terminal help block.
func envelope(tables string) (string, int) {
	help, err := axi.RenderHelp(nil)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return fmt.Sprintf("schema: %d\n", Schema) + tables + help, 0
}

// overviewRows projects every record row in record order. The checked date is the
// delegation-guard cell's, because that verdict is the row's one graded claim in this
// projection.
func overviewRows() [][]string {
	rows := make([][]string, 0, len(Rows))
	for _, row := range Rows {
		rows = append(rows, []string{
			row.Harness,
			string(row.Providers),
			row.PhaseForm,
			row.HookConfig,
			string(row.DelegationGuard.Value),
			row.Headless,
			row.DelegationGuard.Checked,
		})
	}
	return rows
}

// detailRows projects one row's cells: the twelve mechanics in Mechanics order, then the
// delegation guard. An unknown cell leaves its source and date empty, so a reader can tell
// a recorded "no" from an ungraded mechanic.
func detailRows(row Row) [][]string {
	rows := make([][]string, 0, len(Mechanics)+1)
	for _, name := range Mechanics {
		rows = append(rows, cellRow(name, row.Mechanics[name]))
	}
	return append(rows, cellRow(detailDelegationGuard, row.DelegationGuard))
}

// measureRows projects one row's measure cells in Measures order. Every value reads unknown
// until the named supplier ships, so the supplier is the row's one live fact.
func measureRows(row Row) [][]string {
	rows := make([][]string, 0, len(Measures))
	for _, name := range Measures {
		cell := row.Measures[name]
		rows = append(rows, []string{name, string(cell.Value), cell.Supplier})
	}
	return rows
}

func cellRow(field string, cell Cell) []string {
	return []string{field, string(cell.Value), cell.Source, cell.Checked}
}
