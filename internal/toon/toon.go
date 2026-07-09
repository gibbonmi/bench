// Package toon renders flat-table TOON blocks — the machine-queryable output shape
// shared by every AXI query command. It is a thin adapter over the official
// github.com/toon-format/toon-go encoder: the library owns spec-TOON cell escaping
// and quoting, and this package preserves the kit's block contract the library does
// not give for free — the `name[N]{fields}:` count header, the `{fields}` schema on
// an empty table, and a trailing newline on every line.
//
// Only the flat-table form is exposed (a `name[N]{fields}:` header followed by one
// escaped, comma-joined, two-space-indented row per record); the library's general
// (nested) TOON is intentionally not surfaced — the kit emits only flat tables.
package toon

import (
	"fmt"
	"strings"

	toonlib "github.com/toon-format/toon-go"
)

// table is the shared emission core behind Table and TableTyped: it marshals the rows
// through the library with the kit's pinned options and reattaches the block contract.
// Rows carry per-cell Go values so a genuinely-numeric cell (a typed int) stays bare
// while a numeric-looking string is quoted per spec. An empty table is hand-rendered
// because the library, having no elements to infer a schema from, drops `{fields}`.
//
// A cell whose bytes spec-TOON cannot represent — a control character other than the
// escapable tab/newline/return — makes the library refuse; that surfaces as an error
// so the caller can emit the AXI error contract, never a crash or a lossy block.
func table(name string, fields []string, rows [][]any) (string, error) {
	if len(rows) == 0 {
		// No cells, so no escaping is at stake — join the schema identifiers verbatim
		// to keep an empty result self-describing.
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
	// Length markers off gives the kit's `name[N]{...}` (on would emit `name[#N]{...}`);
	// 2-space indent and the default comma delimiter match the AXI block shape.
	out, err := toonlib.MarshalString(map[string]any{name: objs}, toonlib.WithIndent(2))
	if err != nil {
		return "", err
	}
	// The library omits the final newline; callers concatenate blocks and the runtime
	// contracts expect a newline per line, header included.
	return out + "\n", nil
}

// Table renders a flat-table TOON block from all-string rows: a `name[len(rows)]{fields}:`
// header followed by one two-space-indented, comma-joined, spec-escaped row per record.
// Empty rows yields the definitive empty table `name[0]{fields}:`. Field names are
// emitted verbatim (they are schema identifiers, not data). The returned string ends
// with a trailing newline on every line, header included. It returns an error when a
// cell holds a byte spec-TOON cannot represent, so the caller emits the AXI error
// contract rather than crashing or forging output.
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
// that is genuinely mixed — an int cell emits bare (integer semantics survive a
// round-trip), a string cell is spec-escaped like any other. Same block contract and
// same error condition as Table.
func TableTyped(name string, fields []string, rows [][]any) (string, error) {
	return table(name, fields, rows)
}

// Representable reports whether s is a cell value spec-TOON can carry: it rejects
// exactly what the encoder refuses — a control character below U+0020 other than the
// escapable tab, newline, and return (DEL and higher controls the library accepts).
// This is the one predicate for callers that pre-filter rows instead of failing a
// whole table on one cell; the package test pins it to the encoder's observed refusal
// behavior, so a library change turns the gate red here rather than drifting silently.
func Representable(s string) bool {
	for _, r := range s {
		if r < 0x20 && r != '\n' && r != '\r' && r != '\t' {
			return false
		}
	}
	return true
}

// IsSpace matches POSIX [[:space:]] in the C locale — the one source of the AXI
// whitespace class, shared by every parser's field trimming. A UTF-8 continuation or
// lead byte is >= 0x80 and never matches, so a multibyte rune is never mistaken for
// whitespace.
func IsSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}
	return false
}

// Errorf renders a structured AXI error line: `error: <kind> — <hint>`. Per the hybrid
// contract this prints to stdout with exit 1; the caller adds the newline and the exit.
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
// <cmd> (missing argument: <what>)`. Same shape as Usage, distinct label — the one
// source every command shares for "you forgot the required positional" instead of
// the wrong "unknown argument" template. Per the hybrid contract this prints to
// stdout with exit 2.
func MissingArg(cmd, what string) string {
	return fmt.Sprintf("usage: %s (missing argument: %s)", cmd, what)
}

// NotInRepo is the structured error every AXI command prints (with exit 1) when the
// cwd is outside a git repository — one source for the shared phrasing.
func NotInRepo() string {
	return Errorf("not in a git repository", "run inside a Bench-linked repo")
}
