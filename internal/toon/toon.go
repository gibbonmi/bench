// Package toon renders flat-table TOON blocks, the machine-queryable output shape
// shared by every AXI query command. It is a thin adapter over the official
// github.com/toon-format/toon-go encoder. The library owns spec-TOON cell escaping and
// quoting. This package adds the block contract the library does not provide. It
// supplies the `name[N]{fields}:` count header, the `{fields}` schema on an empty
// table, and a trailing newline on every line.
//
// Only the flat-table form is exposed: a `name[N]{fields}:` header followed by one
// escaped, comma-joined, two-space-indented row per record. The library's general
// (nested) TOON is not surfaced. The kit emits only flat tables.
package toon

import (
	"fmt"
	"strings"

	toonlib "github.com/toon-format/toon-go"

	"github.com/gibbonmi/bench/internal/bounds"
)

// table is the shared emission core behind Table and TableTyped. It marshals the rows
// through the library with the kit's pinned options and reattaches the block contract.
// A row carries per-cell Go values, so a genuinely numeric cell (a typed int) stays
// bare while a numeric-looking string is quoted per spec. An empty table is
// hand-rendered, because the library drops `{fields}` when it has no elements to
// infer a schema from.
//
// A cell byte that spec-TOON cannot represent — a control character other than the
// escapable tab, newline, or return — makes the library refuse. The refusal surfaces
// as an error, so the caller emits the AXI error contract instead of a crash or a
// lossy block.
func table(name string, fields []string, rows [][]any) (string, error) {
	if len(rows) == 0 {
		// There are no cells, so no escaping is at stake. This branch joins the schema
		// identifiers verbatim to keep an empty result self-describing.
		return fmt.Sprintf("%s[0]{%s}:\n", name, strings.Join(fields, ",")), nil
	}
	objs := make([]toonlib.Object, len(rows))
	for i, row := range rows {
		cells := make([]toonlib.Field, len(fields))
		for j, f := range fields {
			cells[j] = toonlib.Field{Key: f, Value: row[j]}
		}
		objs[i] = toonlib.NewObject(cells...)
	}
	// Length markers stay off, so the kit gets `name[N]{...}` instead of `name[#N]{...}`.
	// The 2-space indent and the default comma delimiter match the AXI block shape.
	out, err := toonlib.MarshalString(map[string]any{name: objs}, toonlib.WithIndent(2))
	if err != nil {
		return "", err
	}
	// The library omits the final newline. Callers concatenate blocks, and the runtime
	// contracts expect a newline on every line, header included.
	return out + "\n", nil
}

// Table renders a flat-table TOON block from all-string rows: a `name[len(rows)]{fields}:`
// header followed by one two-space-indented, comma-joined, spec-escaped row per record.
// Empty rows yield the definitive empty table `name[0]{fields}:`. Field names are
// emitted verbatim (they are schema identifiers, not data). The returned string ends
// with a trailing newline on every line, header included. It returns an error when a
// cell holds a byte spec-TOON cannot represent. The caller then emits the AXI error
// contract instead of a crash or forged output.
func Table(name string, fields []string, rows [][]string) (string, error) {
	typed := make([][]any, len(rows))
	for i, row := range rows {
		cells := make([]any, len(row))
		for j, c := range row {
			cells[j] = c
		}
		typed[i] = cells
	}
	return table(name, fields, typed)
}

// TableTyped renders a flat-table TOON block from per-cell-typed rows, for a column
// that is genuinely mixed. An int cell emits bare, so integer semantics survive a
// round-trip. A string cell is spec-escaped like any other. TableTyped shares the
// same block contract and the same error condition as Table.
func TableTyped(name string, fields []string, rows [][]any) (string, error) {
	return table(name, fields, rows)
}

// Representable reports whether s is a cell value spec-TOON can carry. It rejects
// exactly what the encoder refuses: a control character below U+0020 other than the
// escapable tab, newline, and return. DEL and higher controls pass. A caller uses this
// predicate to filter rows one at a time, so one bad cell does not fail the whole
// table. The package test pins it to the encoder's observed refusal behavior, so a
// library change turns the gate red here instead of a silent drift.
func Representable(s string) bool {
	for _, r := range s {
		if r < 0x20 && r != '\n' && r != '\r' && r != '\t' {
			return false
		}
	}
	return true
}

// IsSpace matches POSIX [[:space:]] in the C locale. It is the one source of the AXI
// whitespace class, and every parser uses it to trim fields. A UTF-8 continuation byte
// or lead byte is >= 0x80 and never matches, so a multibyte rune is never mistaken for
// whitespace.
func IsSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}
	return false
}

// Errorf renders a structured AXI error line: `error: <kind> — <hint>`. The hybrid
// contract prints this line to stdout with exit 1. The caller adds the newline and
// the exit.
func Errorf(kind, hint string) string {
	return fmt.Sprintf("error: %s — %s", kind, hint)
}

// RenderError is the one source of the AXI error line every command emits when Table or
// TableTyped rejects a cell spec-TOON cannot represent. Callers surface it with exit 1.
func RenderError(err error) string {
	return Errorf("unrepresentable TOON cell", err.Error())
}

// Usage renders an AXI usage line for an unknown argument: `usage: <cmd> (unknown
// argument: <arg>)`. Per the hybrid contract this prints to stdout with exit 2.
func Usage(cmd, arg string) string {
	return fmt.Sprintf("usage: %s (unknown argument: %s)", cmd, arg)
}

// MissingArg renders an AXI usage line for a missing required argument: `usage:
// <cmd> (missing argument: <what>)`. It has the same shape as Usage but a distinct
// label. Every command shares this label for a missing positional argument instead of
// the wrong "unknown argument" template. The hybrid contract prints this line to
// stdout with exit 2.
func MissingArg(cmd, what string) string {
	return fmt.Sprintf("usage: %s (missing argument: %s)", cmd, what)
}

// NotInRepo is the structured error every AXI command prints (with exit 1) when the
// cwd is outside a git repository. It is the one source for the shared phrasing.
func NotInRepo() string {
	return Errorf("not in a git repository", "run inside a Bench-linked repo")
}

// RecordError is the AXI error line every command emits when a control record did not
// read as something it may trust: `error: <path> is <state> — <reason>`. path is
// repo-relative, so the line names what an agent would type. Callers add the newline and
// the exit-1. Every fail-closed surface composes this call instead of writing its own
// grammar, so the phrasing an agent parses cannot drift between commands.
func RecordError(path string, state bounds.FileState, reason string) string {
	return Errorf(path+" is "+string(state), reason)
}

// UnknownCell is RecordError's sibling for a surface that degrades instead of exiting.
// It renders the `unknown (<path> is <state>)` detail cell the dashboard prints in
// place of a count it could not derive. The parenthetical carries the same
// `<path> is <state>` clause as the error line. One record's failure reads the same,
// whether the command failed on it or the dashboard reported around it.
func UnknownCell(path string, state bounds.FileState) string {
	return fmt.Sprintf("unknown (%s is %s)", path, state)
}
