package probe

import (
	"bytes"
	"fmt"
	"math"
	"slices"

	toonlib "github.com/toon-format/toon-go"

	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/treetarget"
)

// Result is the part of one complete probe answer that a consumer may import. It names
// what the producer reported; it does not prove the origin of the bytes or what ran.
type Result struct {
	// Verdict is the verdict cell: bit, silent, invalid, or restore-failed.
	Verdict string
	// Restored reports that the producer proved the subject's start bytes back.
	Restored bool
	// CommandExit is the exit of the probe command itself.
	CommandExit int
	// FocusedExit is the normalized exit of the mutated run, or nil when none started.
	FocusedExit *int
}

// ReadResult reads one complete probe answer: the complete stdout with its optional tree
// row, or the complete spill body the caller selected. A projection or spill-failed frame
// refuses before any parse, and no path that a frame names is opened. The envelope must
// start the body, decode strictly, render back to its exact bytes, and appear once.
// trailing_bytes must equal the exact length of everything after the envelope.
func ReadResult(data []byte) (Result, error) {
	if err := responsebound.RequireComplete(data); err != nil {
		return Result{}, err
	}
	_, body, _, err := treetarget.ReadLead(data)
	if err != nil {
		return Result{}, err
	}
	e, suffix, err := readEnvelope(body)
	if err != nil {
		return Result{}, err
	}
	if len(suffix) != e.trailing {
		return Result{}, fmt.Errorf("the probe result is incomplete: %d bytes follow the envelope, and trailing_bytes declares %d", len(suffix), e.trailing)
	}
	return Result{Verdict: e.cells.verdict, Restored: e.cells.restored == restoredCell, CommandExit: e.command, FocusedExit: e.focused}, nil
}

// readEnvelope splits the body after the envelope's physical lines, a header and one row
// for each table, and decodes those lines. The suffix must start no further envelope table.
func readEnvelope(body []byte) (envelope, []byte, error) {
	end := 0
	for i := range 2 * len(resultTables) {
		if i%2 == 0 && !bytes.HasPrefix(body[end:], []byte(resultTables[i/2].name+"[")) {
			return envelope{}, nil, fmt.Errorf("the probe result has no %s table at line %d; the supported schema is %s", resultTables[i/2].name, i+1, resultSchema)
		}
		next := bytes.IndexByte(body[end:], '\n')
		if next < 0 {
			return envelope{}, nil, fmt.Errorf("the probe result envelope is incomplete at line %d", i+1)
		}
		end += next + 1
	}
	head, suffix := body[:end], body[end:]
	for _, line := range bytes.SplitAfter(suffix, []byte("\n")) {
		for _, table := range resultTables {
			if bytes.HasPrefix(line, []byte(table.name+"[")) {
				return envelope{}, nil, fmt.Errorf("the probe result repeats its %s table", table.name)
			}
		}
	}
	e, err := decodeEnvelope(head)
	if err != nil {
		return envelope{}, nil, fmt.Errorf("the probe result envelope is malformed: %w", err)
	}
	if canonical, err := e.rows(); err != nil || canonical != string(head) {
		return envelope{}, nil, fmt.Errorf("the probe result envelope is not canonical")
	}
	return e, suffix, nil
}

// decodeEnvelope decodes the envelope through the pinned strict TOON decoder and checks
// the table set, each field set, and each cell type. The strict decoder checks counts and
// widths but keeps the last of two equal keys, so the caller also renders the envelope back.
func decodeEnvelope(head []byte) (envelope, error) {
	decoded, err := toonlib.Decode(head)
	if err != nil {
		return envelope{}, err
	}
	document, ok := decoded.(map[string]any)
	if !ok || len(document) != len(resultTables) {
		return envelope{}, fmt.Errorf("decodes as %T with the wrong table set", decoded)
	}
	rows := make([]cellReader, len(resultTables))
	for i, table := range resultTables {
		list, ok := document[table.name].([]any)
		if !ok || len(list) != 1 {
			return envelope{}, fmt.Errorf("want one %s row", table.name)
		}
		row, ok := list[0].(map[string]any)
		if !ok || len(row) != len(table.fields) {
			return envelope{}, fmt.Errorf("want the %s fields %v", table.name, table.fields)
		}
		rows[i] = cellReader{table: table.name, row: row}
	}
	// Each index is the cell's position in envelope.rows, so each cell reads the field that
	// it renders under.
	verdict, selection, execution := &rows[0], &rows[1], &rows[2]
	e := envelope{
		cells: verdictCells{verdict: verdict.text(probeFields[0]), cause: verdict.text(probeFields[3]), failed: verdict.integer(probeFields[4]),
			restored: verdict.text(probeFields[5]), baseline: selection.text(selectionFields[3]), ran: selection.integer(selectionFields[4])},
		subject: verdict.text(probeFields[1]), mutation: verdict.text(probeFields[2]),
		form: selection.text(selectionFields[0]), target: selection.text(selectionFields[1]), run: selection.text(selectionFields[2]),
		command: execution.integer(executionFields[1]), trailing: execution.integer(executionFields[3]),
	}
	if execution.row[executionFields[2]] != nil {
		focused := execution.integer(executionFields[2])
		e.focused = &focused
	}
	schema := execution.text(executionFields[0])
	for _, reader := range rows {
		if reader.err != nil {
			return envelope{}, reader.err
		}
	}
	if schema != resultSchema {
		return envelope{}, fmt.Errorf("the schema %q is unsupported; the supported schema is %s", schema, resultSchema)
	}
	if e.trailing < 0 {
		return envelope{}, fmt.Errorf("trailing_bytes %d is negative", e.trailing)
	}
	return e, nil
}

// FieldSource names one producer field as `table.field`, so a consumer's help names the
// field this package renders. A name outside the envelope is a defect in the caller.
func FieldSource(field string) string {
	for _, table := range resultTables {
		if slices.Contains(table.fields, field) {
			return table.name + "." + field
		}
	}
	panic("probe result has no field " + field)
}

// cellReader reads typed cells from one decoded row and keeps the first type error.
type cellReader struct {
	table string
	row   map[string]any
	err   error
}

func (r *cellReader) text(field string) string {
	value, ok := r.row[field].(string)
	r.fail(ok, field, "a string")
	return value
}

func (r *cellReader) integer(field string) int {
	value, ok := r.row[field].(float64)
	ok = ok && value == math.Trunc(value) && math.Abs(value) <= 1<<53
	r.fail(ok, field, "an integer")
	return int(value)
}

func (r *cellReader) fail(ok bool, field, want string) {
	if !ok && r.err == nil {
		r.err = fmt.Errorf("%s.%s is %T, want %s", r.table, field, r.row[field], want)
	}
}
