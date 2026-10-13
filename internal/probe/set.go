package probe

import (
	"context"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// setSchema is the schema cell of a mutant-set answer. It names the field set of the set
// tables, so it changes with them, and a single-form reader refuses it.
const setSchema = "bench-probe-set-v1"

// mutantFieldsOfSet is the schema of the mutants table, one row for each entry.
var mutantFieldsOfSet = []string{"mutant", "verdict", "mutation", "cause", "failed_tests", "first_failure", "ran", "focused_exit_code", "exit_code", "first_build_line", "restored"}

// setSelectionFields is the schema of the one selection row that every mutant shares.
var setSelectionFields = []string{"form", "target", "run", "baseline", "builds"}

// setExecutionFields is the schema of the execution row. trailing_bytes counts every byte after
// the row, so a reader can tell a complete answer from a cut one.
var setExecutionFields = []string{"schema", "command_exit_code", "trailing_bytes"}

// setTables names the three tables of a mutant-set answer in their printed order. They stay
// apart from resultTables, so FieldSource never names a set field.
var setTables = []struct {
	name   string
	fields []string
}{{"mutants", mutantFieldsOfSet}, {"selection", setSelectionFields}, {"probe_set_execution", setExecutionFields}}

// mutantReportFields is the schema of the table that heads each started mutant's report.
var mutantReportFields = []string{"mutant"}

// The causes of a row that the set stopped before its write.
const (
	setStoppedCause      = "set-stopped"
	gateHeldCause        = "gate-held"
	gateUnavailableCause = "gate-unavailable"
	mutationFailedCause  = "mutation-failed"
	preserveFailedCause  = "preserve-failed"
)

// mutantRow is one mutants row and its evidence. focused is nil when the mutated run never
// started, and report is that run's report.
type mutantRow struct {
	mutation string
	cells    verdictCells
	focused  *int
	exit     int
	report   string
}

// setEnvelope is one mutant-set answer's three leading tables.
type setEnvelope struct {
	rows                        []mutantRow
	form, target, run, baseline string
	builds, command, trailing   int
}

// head renders the three tables, so one cell that the encoder cannot carry refuses them all.
func (e setEnvelope) head() (string, error) {
	rows := make([][]any, len(e.rows))
	for i, row := range e.rows {
		var first, focused, buildLine any
		if row.cells.first != nil {
			first = *row.cells.first
		}
		if row.focused != nil {
			focused = *row.focused
		}
		if row.cells.buildLine != nil {
			buildLine = *row.cells.buildLine
		}
		rows[i] = []any{i + 1, row.cells.verdict, row.mutation, row.cells.cause, row.cells.failed, first, row.cells.ran, focused, row.exit, buildLine, row.cells.restored}
	}
	cells := [][][]any{rows, {{e.form, e.target, e.run, e.baseline, e.builds}}, {{setSchema, e.command, e.trailing}}}
	var out string
	for i, table := range setTables {
		block, err := toon.TableTyped(table.name, table.fields, cells[i])
		if err != nil {
			return "", err
		}
		out += block
	}
	return out, nil
}

// setRun is one mutant set in flight: the subject, the held runs, the entries and their rows
// in entry order, and the preserved table of a failed restore.
type setRun struct {
	replacer  durablefile.Replacer
	root      string
	subject   subject
	ctx       context.Context
	runs      *testreport.HeldRuns
	entries   []setEntry
	rows      []mutantRow
	preserved string
}

// runSet grades a mutant set in the single form's refusal order. One subject lock, one reread,
// one baseline, and one cancel context cover every mutant. Each matching entry then preserves,
// writes, runs, and restores in file order.
func runSet(replacer durablefile.Replacer, root string, parsed usage.Result, subject subject) (string, int) {
	entries, line, code := readMutantSet(parsed, subject)
	if line != "" {
		return line, code
	}
	request, line, code := prepareSelection(root, parsed, subject)
	if line != "" {
		return line, code
	}
	s := &setRun{replacer: replacer, root: root, subject: subject, entries: entries, rows: make([]mutantRow, len(entries))}
	e := setEnvelope{form: request.Form(), target: request.Target(), run: request.Run()}
	matched := false
	for i, entry := range entries {
		s.rows[i] = unwritten(entry.input.cell, untouchedCells(substringMissCause))
		matched = matched || !entry.mutant.miss
	}
	if !matched {
		return s.answer(e, "")
	}
	held, line := holdSubject(root, parsed.Positionals[0], subject)
	if line != "" {
		return line, 1
	}
	defer held.release()
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	s.ctx, s.runs = ctx, testreport.Hold(ctx, root, request)
	defer func() { _ = s.runs.Close() }()
	baseline, report, _ := s.runs.Run()
	e.baseline = string(baseline.Kind)
	if baseline.Kind != testreport.OutcomePassed {
		s.stopFrom(0, baselineCells(baseline))
		e.builds = s.runs.Builds()
		return s.answer(e, report)
	}
	for i, entry := range entries {
		if !entry.mutant.miss && !s.grade(i, entry) {
			break
		}
	}
	e.builds = s.runs.Builds()
	return s.answer(e, "")
}

// grade checks the stop conditions before entry i's write, and then preserves, writes, runs,
// and restores that mutant. It answers false when the set stops. A run that a signal ended
// stops the set at the next write's context check.
func (s *setRun) grade(i int, entry setEntry) bool {
	if cause := s.stopCause(); cause != "" {
		s.stopFrom(i, untouchedCells(cause))
		return false
	}
	preserved, line := preserve(s.root, s.subject)
	if line != "" {
		s.rows[i] = unwritten(entry.input.cell, untouchedCells(preserveFailedCause))
		s.stopFrom(i+1, untouchedCells(setStoppedCause))
		return false
	}
	if published, err := entry.mutant.write(s.replacer, s.subject); err != nil {
		s.rows[i] = unwritten(entry.input.cell, untouchedCells(mutationFailedCause))
		if published {
			s.rows[i].cells.restored = restoredCell
		}
		if kept := recoverMutation(s.replacer, s.subject, preserved, err, published); kept != "" {
			s.rows[i].exit = s.rows[i].cells.failRestore()
			s.preserved = keptCopy(preserved, kept)
		}
		s.stopFrom(i+1, untouchedCells(setStoppedCause))
		return false
	}
	run := runRestored(s.replacer, s.subject, preserved, s.runs.Run)
	cells, code := mutatedCells(run.outcome, entry.expect, testreport.OutcomePassed)
	s.rows[i] = mutantRow{mutation: entry.input.cell, cells: cells, focused: &run.focused, exit: code, report: run.report}
	if !run.restored {
		s.rows[i].exit = s.rows[i].cells.failRestore()
		s.preserved = keptCopy(preserved, run.reason)
		s.stopFrom(i+1, untouchedCells(setStoppedCause))
		return false
	}
	return true
}

// stopCause answers why the set stops before the next write, or an empty cause. A cancelled
// context stops it, and so does a gate lock that is held or unreadable.
func (s *setRun) stopCause() string {
	if s.ctx.Err() != nil {
		return setStoppedCause
	}
	held, err := gate.ExecutionInProgress(s.root)
	switch {
	case err != nil:
		return gateUnavailableCause
	case held:
		return gateHeldCause
	}
	return ""
}

// stopFrom gives row i and every later matching row the cells of a stop before its write. A
// missed row keeps its substring-miss cause.
func (s *setRun) stopFrom(i int, cells verdictCells) {
	for j := i; j < len(s.rows); j++ {
		if !s.entries[j].mutant.miss {
			s.rows[j] = unwritten(s.rows[j].mutation, cells)
		}
	}
}

// unwritten answers the row of a mutant that ended with cells before its write, at exit 1.
func unwritten(mutation string, cells verdictCells) mutantRow {
	return mutantRow{mutation: mutation, cells: cells, exit: 1}
}

// answer renders the three tables and then the suffix: the preserved table of a failed
// restore, the baseline report, and each started mutant's report. The command exit is 2 when a
// restore failed, 0 when every row bit, and 1 otherwise.
func (s *setRun) answer(e setEnvelope, suffix string) (string, int) {
	e.rows = s.rows
	for _, row := range s.rows {
		if row.exit == 2 {
			e.command = 2
			break
		}
		e.command = max(e.command, row.exit)
	}
	reports := ""
	for i, row := range s.rows {
		if row.focused == nil {
			continue
		}
		heading, err := toon.TableTyped("mutant_report", mutantReportFields, [][]any{{i + 1}})
		if err != nil {
			return s.refuse(err, e.command)
		}
		reports += heading + row.report
	}
	suffix = s.preserved + suffix + reports
	e.trailing = len(suffix)
	head, err := e.head()
	if err != nil {
		return s.refuse(err, e.command)
	}
	return head + suffix, e.command
}

// refuse answers a set whose tables the encoder refused. A kept copy still prints, and its
// exit 2 overrides the render refusal's exit 1.
func (s *setRun) refuse(err error, command int) (string, int) {
	return toon.RenderError(err) + "\n" + s.preserved, max(command, 1)
}
