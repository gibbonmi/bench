package probe

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"

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

// MutantResult is the part of one mutants row of a complete mutant-set answer that a consumer
// may import. Like Result, it names what the producer reported.
type MutantResult struct {
	// Verdict is the row's verdict cell.
	Verdict string
	// Restored reports that the producer proved the subject's start bytes back after this mutant.
	Restored bool
	// RowExit is the row's exit_code, the single form's exit for the row's verdict.
	RowExit int
	// FocusedExit is the normalized exit of the row's mutated run, or nil when none started.
	FocusedExit *int
	// CommandExit is the exit of the whole set, which no row import binds.
	CommandExit int
}

// ErrMutantSet is ReadResult's refusal of a mutant-set answer, which holds no single probe.
var ErrMutantSet = errors.New("the probe result is a mutant-set answer")

// ErrSingleForm is ReadMutant's refusal of a single-form answer, which holds no mutants table.
var ErrSingleForm = errors.New("the probe result is a single-form answer")

// RowRangeError is ReadMutant's refusal of a row number that names no row of the answer.
type RowRangeError struct {
	// Rows is the answer's row count.
	Rows int
}

func (e *RowRangeError) Error() string {
	return fmt.Sprintf("the mutant-set answer holds %d rows", e.Rows)
}

// ResultSchema is the schema cell of a single-form answer.
const ResultSchema = resultSchema

// answerForm is one answer form that a reader takes: its tables, its schema cell, and the
// sentinel that a reader of the other form answers when this form's first table opens the body.
type answerForm struct {
	tables   []answerTable
	schema   string
	sentinel error
}

var (
	singleForm = answerForm{tables: resultTables, schema: resultSchema, sentinel: ErrSingleForm}
	setForm    = answerForm{tables: setTables, schema: setSchema, sentinel: ErrMutantSet}
)

// answerHead is one decoded envelope: it renders back to its canonical bytes, and it declares
// how many bytes follow it.
type answerHead interface {
	head() (string, error)
	trailingBytes() int
}

// ReadResult reads one complete probe answer: the complete stdout with its optional tree
// row, or the complete spill body the caller selected. readForm states the shared rules.
func ReadResult(data []byte) (Result, error) {
	e, err := readForm(data, singleForm, setForm, decodeEnvelope)
	if err != nil {
		return Result{}, err
	}
	return Result{Verdict: e.cells.verdict, Restored: e.cells.restored == restoredCell, CommandExit: e.command, FocusedExit: e.focused}, nil
}

// ReadMutant reads one complete mutant-set answer under the rules of ReadResult and answers its
// row n, counted from 1.
func ReadMutant(data []byte, n int) (MutantResult, error) {
	e, err := readForm(data, setForm, singleForm, decodeSetEnvelope)
	if err != nil {
		return MutantResult{}, err
	}
	if n < 1 || n > len(e.rows) {
		return MutantResult{}, &RowRangeError{Rows: len(e.rows)}
	}
	row := e.rows[n-1]
	return MutantResult{Verdict: row.cells.verdict, Restored: row.cells.restored == restoredCell, RowExit: row.exit, FocusedExit: row.focused, CommandExit: e.command}, nil
}

// readForm reads one complete answer of form. A projection or spill-failed frame refuses
// before any parse, and no path that a frame names is opened. The envelope must start the
// body after an optional tree row, decode strictly, render back to its exact bytes, and appear
// once. trailing_bytes must equal the exact length of everything after the envelope. A body
// that opens with other's first table refuses with other's sentinel.
func readForm[E answerHead](data []byte, form, other answerForm, decode func([][]cellReader) (E, error)) (E, error) {
	var zero E
	if err := responsebound.RequireComplete(data); err != nil {
		return zero, err
	}
	_, body, _, err := treetarget.ReadLead(data)
	if err != nil {
		return zero, err
	}
	head, suffix, err := splitEnvelope(body, form, other)
	if err != nil {
		return zero, err
	}
	tables, err := decodeTables(head, form)
	var e E
	if err == nil {
		e, err = decode(tables)
	}
	if err == nil && e.trailingBytes() < 0 {
		err = fmt.Errorf("trailing_bytes %d is negative", e.trailingBytes())
	}
	if err != nil {
		return zero, fmt.Errorf("the probe result envelope is malformed: %w", err)
	}
	if canonical, err := e.head(); err != nil || canonical != string(head) {
		return zero, fmt.Errorf("the probe result envelope is not canonical")
	}
	if len(suffix) != e.trailingBytes() {
		return zero, fmt.Errorf("the probe result is incomplete: %d bytes follow the envelope, and trailing_bytes declares %d", len(suffix), e.trailingBytes())
	}
	return e, nil
}

// splitEnvelope splits the body after the envelope's physical lines: each table's header and
// the rows that the header counts. The suffix must start no further envelope table.
func splitEnvelope(body []byte, form, other answerForm) ([]byte, []byte, error) {
	// lines counts the physical lines before end.
	end, lines := 0, 0
	for _, table := range form.tables {
		prefix := []byte(table.name + "[")
		if !bytes.HasPrefix(body[end:], prefix) {
			if end == 0 && bytes.HasPrefix(body, []byte(other.tables[0].name+"[")) {
				return nil, nil, other.sentinel
			}
			return nil, nil, fmt.Errorf("the probe result has no %s table at line %d; the supported schema is %s", table.name, lines+1, form.schema)
		}
		count, _, _ := bytes.Cut(body[end+len(prefix):], []byte("]"))
		rows, err := strconv.Atoi(string(count))
		if err != nil || rows < 0 {
			return nil, nil, fmt.Errorf("the probe result has no row count in its %s header at line %d", table.name, lines+1)
		}
		for range rows + 1 {
			next := bytes.IndexByte(body[end:], '\n')
			if next < 0 {
				return nil, nil, fmt.Errorf("the probe result envelope is incomplete at line %d", lines+1)
			}
			end += next + 1
			lines++
		}
	}
	head, suffix := body[:end], body[end:]
	for _, line := range bytes.SplitAfter(suffix, []byte("\n")) {
		for _, table := range form.tables {
			if bytes.HasPrefix(line, []byte(table.name+"[")) {
				return nil, nil, fmt.Errorf("the probe result repeats its %s table", table.name)
			}
		}
	}
	return head, suffix, nil
}

// decodeTables decodes the envelope through the pinned strict TOON decoder and checks the
// table set, the schema, each row count, and each field set. It answers a cell reader for each
// row of each table in printed order. The strict decoder checks counts and widths but keeps the
// last of two equal keys, so the caller also renders the envelope back.
func decodeTables(head []byte, form answerForm) ([][]cellReader, error) {
	decoded, err := toonlib.Decode(head)
	if err != nil {
		return nil, err
	}
	document, ok := decoded.(map[string]any)
	if !ok || len(document) != len(form.tables) {
		return nil, fmt.Errorf("decodes as %T with the wrong table set", decoded)
	}
	if err := requireSchema(document, form); err != nil {
		return nil, err
	}
	tables := make([][]cellReader, len(form.tables))
	for i, table := range form.tables {
		list, ok := document[table.name].([]any)
		if !ok || len(list) == 0 || !table.several && len(list) != 1 {
			return nil, fmt.Errorf("want one %s row", table.name)
		}
		for _, item := range list {
			row, ok := item.(map[string]any)
			if !ok || len(row) != len(table.fields) {
				return nil, fmt.Errorf("want the %s fields %v", table.name, table.fields)
			}
			tables[i] = append(tables[i], cellReader{table: table.name, row: row})
		}
	}
	return tables, nil
}

// decodeEnvelope reads the cells of a single-form answer's decoded tables.
func decodeEnvelope(tables [][]cellReader) (envelope, error) {
	// Each index is the cell's position in envelope.head, so each cell reads the field that
	// it renders under.
	verdict, selection, execution := &tables[0][0], &tables[1][0], &tables[2][0]
	e := envelope{
		cells: verdictCells{verdict: verdict.text(probeFields[0]), cause: verdict.text(probeFields[3]), failed: verdict.integer(probeFields[4]),
			first: verdict.optionalText(probeFields[5]), restored: verdict.text(probeFields[6]), baseline: selection.text(selectionFields[3]), ran: selection.integer(selectionFields[4])},
		subject: verdict.text(probeFields[1]), mutation: verdict.text(probeFields[2]),
		form: selection.text(selectionFields[0]), target: selection.text(selectionFields[1]), run: selection.text(selectionFields[2]),
		command: execution.integer(executionFields[1]), focused: execution.optionalInteger(executionFields[2]), trailing: execution.integer(executionFields[3]),
	}
	return e, readError(tables)
}

// requireSchema refuses an envelope whose schema cell is not the form's schema. It runs before
// the field sets are checked, because another schema carries another field set, and a
// field error would hide that the answer comes from another schema.
func requireSchema(document map[string]any, form answerForm) error {
	execution := form.tables[len(form.tables)-1]
	var schema any
	if list, ok := document[execution.name].([]any); ok && len(list) == 1 {
		if row, ok := list[0].(map[string]any); ok {
			schema = row[execution.fields[0]]
		}
	}
	if schema != form.schema {
		return fmt.Errorf("the schema %q is unsupported; the supported schema is %s", fmt.Sprint(schema), form.schema)
	}
	return nil
}

// FieldSource names one producer field as `table.field`, so a consumer's help names the
// field this package renders. A name outside the envelope is a defect in the caller.
func FieldSource(field string) string {
	return fieldSource(resultTables, field)
}

// SetFieldSource names one mutant-set field as `table.field`, as FieldSource does for the
// single form.
func SetFieldSource(field string) string {
	return fieldSource(setTables, field)
}

func fieldSource(tables []answerTable, field string) string {
	for _, table := range tables {
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

// optionalText reads a string cell that may be null, and nil stands for null.
func (r *cellReader) optionalText(field string) *string {
	if r.row[field] == nil {
		return nil
	}
	value := r.text(field)
	return &value
}

// optionalInteger reads an integer cell that may be null, and nil stands for null.
func (r *cellReader) optionalInteger(field string) *int {
	if r.row[field] == nil {
		return nil
	}
	value := r.integer(field)
	return &value
}

func (r *cellReader) fail(ok bool, field, want string) {
	if !ok && r.err == nil {
		r.err = fmt.Errorf("%s.%s is %T, want %s", r.table, field, r.row[field], want)
	}
}

// readError answers the first type error of the readers in printed order.
func readError(tables [][]cellReader) error {
	for _, rows := range tables {
		for _, reader := range rows {
			if reader.err != nil {
				return reader.err
			}
		}
	}
	return nil
}
