// Package toon renders flat-table TOON blocks — the machine-queryable output shape
// shared by every AXI query command. It is the first internal package of the Go
// port and the single source of the TOON escaping rule that five commands compose,
// ending the two-language duplication the migration exists to close.
//
// Only the flat-table form is implemented (a `name[N]{fields}:` header followed by
// one escaped, comma-joined, two-space-indented row per record); the general TOON
// format is intentionally absent (decision #5).
package toon

import (
	"fmt"
	"strings"
)

// Escape renders one field value for a TOON row. A value that carries a comma, a
// double-quote, or a newline — or one with leading or trailing whitespace, which a
// bare field would lose on parse — is double-quoted with inner quotes doubled;
// anything else is emitted verbatim. This is the one escaping rule; every command's
// row cells pass through here.
func Escape(v string) string {
	if !needsQuote(v) {
		return v
	}
	return `"` + strings.ReplaceAll(v, `"`, `""`) + `"`
}

func needsQuote(v string) bool {
	if strings.ContainsAny(v, ",\"\n") {
		return true
	}
	if v == "" {
		return false
	}
	return IsSpace(v[0]) || IsSpace(v[len(v)-1])
}

// IsSpace matches POSIX [[:space:]] in the C locale — the one source of the AXI
// whitespace class, shared by the emitter's leading/trailing quoting rule and every
// parser's field trimming. A UTF-8 continuation or lead byte is >= 0x80 and never
// matches, so a multibyte rune is never mistaken for whitespace.
func IsSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}
	return false
}

// Table renders a flat-table TOON block: a `name[len(rows)]{fields}:` header followed
// by one two-space-indented, comma-joined, escaped row per record. Empty rows yields
// the definitive empty table `name[0]{fields}:`. Field names are emitted verbatim
// (they are schema identifiers, not data); row cells are Escape'd. The returned string
// ends with a trailing newline on every line, header included.
func Table(name string, fields []string, rows [][]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s[%d]{%s}:\n", name, len(rows), strings.Join(fields, ","))
	for _, row := range rows {
		cells := make([]string, len(row))
		for i, c := range row {
			cells[i] = Escape(c)
		}
		b.WriteString("  ")
		b.WriteString(strings.Join(cells, ","))
		b.WriteByte('\n')
	}
	return b.String()
}

// Errorf renders a structured AXI error line: `error: <kind> — <hint>`. Per the hybrid
// contract this prints to stdout with exit 1; the caller adds the newline and the exit.
func Errorf(kind, hint string) string {
	return fmt.Sprintf("error: %s — %s", kind, hint)
}

// Usage renders an AXI usage line for an unknown argument: `usage: <cmd> (unknown
// argument: <arg>)`. Per the hybrid contract this prints to stdout with exit 2.
func Usage(cmd, arg string) string {
	return fmt.Sprintf("usage: %s (unknown argument: %s)", cmd, arg)
}

// NotInRepo is the structured error every AXI command prints (with exit 1) when the
// cwd is outside a git repository — one source for the shared phrasing.
func NotInRepo() string {
	return Errorf("not in a git repository", "run inside a Bench-linked repo")
}
