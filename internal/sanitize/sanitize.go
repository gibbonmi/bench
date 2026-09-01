// Package sanitize is the one policy for rendering operator-influenced text into a sink
// that could misread it. It covers three sinks and three duties. For a terminal render,
// Controls and its siblings escape control runes, so a control byte in a commit
// subject, an objective, or a git-sourced string is escaped one way everywhere. For a
// command line a reader will paste, ShellQuote renders a value as exactly one shell
// argument. For a table cell that must read as the text it came from, Strip removes the
// control bytes and rewrites nothing else.
//
// The duties here stay distinct from internal/toon's cell policy. toon *refuses* a
// control-bearing cell, a closed AXI decision, while this package escapes, quotes, or
// strips. Strip does compose toon.Representable, so which byte the encoder refuses has
// one source; what to do about that byte stays this package's call. The three duties
// here do not substitute for each other. Quoting leaves a control byte intact; escaping
// produces text that no longer names what it came from; stripping keeps the remaining
// text verbatim and so cannot report what it removed.
package sanitize

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/toon"
)

// Controls escapes every control rune in value, with no length cap. Newline, carriage
// return, and tab escape to their backslash forms; every other control rune escapes to
// \uXXXX; a literal backslash escapes to \\. The result carries no raw control byte, so
// it is safe to hand to a terminal or to embed in html/template output.
func Controls(value string) string {
	var b strings.Builder
	writeEscaped(&b, []rune(value), false)
	return b.String()
}

// Preformatted is Controls for text that lands inside an html/template <pre> block.
// html/template already neutralizes markup there, so the only threat left is a raw
// control byte, and flattening the layout to escape it is unnecessary collateral. It
// leaves newline and tab verbatim, so multi-line and tab-aligned content renders as
// authored. Every other control rune, including carriage return, still escapes through
// the same writeEscaped mechanism Controls uses, so no raw control byte reaches the
// output.
func Preformatted(value string) string {
	var b strings.Builder
	writeEscaped(&b, []rune(value), true)
	return b.String()
}

// Preview is Controls plus the bounds.PreviewRuneLimit cap and a byte-count suffix: it
// escapes the first bounds.PreviewRuneLimit runes and, when value was longer, appends
// "… (N bytes)" naming the original byte length. The cap is counted in code points so a
// multibyte string is not cut at a fraction of its apparent length.
func Preview(value string) string {
	runes := []rune(value)
	truncated := len(runes) > bounds.PreviewRuneLimit
	if truncated {
		runes = runes[:bounds.PreviewRuneLimit]
	}
	var b strings.Builder
	writeEscaped(&b, runes, false)
	if truncated {
		fmt.Fprintf(&b, "… (%d bytes)", len(value))
	}
	return b.String()
}

// Strip removes every rune the TOON cell encoder refuses, and newline and carriage
// return with them, and escapes nothing. It is the filter for a sink whose own encoder
// escapes, where escaping first would reach the reader doubled: a backslash in a path
// would arrive as four. Stripping runs on runes, so a multi-byte sequence is never cut
// into a fragment.
//
// The refusal half is not derived here. Strip asks toon.Representable one rune at a
// time, because only that predicate is pinned to the encoder's observed behavior. A
// second copy of the threshold in this file would drift the day the encoder's refusal
// set moves, and this filter feeds that encoder. The per-rune string is the deliberate
// price of the composition, on text no longer than one phase's stream line.
//
// Newline and carriage return are Strip's own exclusion on top. toon.Representable
// passes both, because the encoder escapes them, but Strip feeds a table of one line
// per row, where a literal newline in a cell forges a row. Every rune from U+0020 up
// still passes verbatim, including DEL and the C1 controls, because Representable
// passes those too.
//
// Removal is all it does. The result no longer says that anything was removed, so a
// caller that must account for the original text escapes it instead.
func Strip(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	for _, r := range value {
		if !toon.Representable(string(r)) || r == '\n' || r == '\r' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// ShellQuote renders value as exactly one POSIX shell argument. The whole string is
// single-quoted, and each embedded quote closes, escapes, and reopens it, so a space, a
// glob character, or a metacharacter inside can neither split the argument nor expand.
// An empty value becomes an empty quoted pair, so it survives as an argument rather
// than vanishing from the command line.
//
// Quoting is all it does. A control byte passes through as a literal byte; single
// quotes make a newline literal but still emit it. A caller that writes into a
// line-structured sink must refuse or escape one itself.
func ShellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

// writeEscaped is the single per-rune escaping rule Controls, Preview, and Preformatted
// all share. When preserveLayout is false (Controls, Preview), newline, carriage return,
// and tab escape to their backslash forms like every other control rune. When
// preserveLayout is true (Preformatted), newline and tab pass through verbatim. Only
// carriage return joins the general control-rune case then, so the \u%04x emission
// stays single-sourced regardless of which caller reaches it.
func writeEscaped(b *strings.Builder, runes []rune, preserveLayout bool) {
	for _, r := range runes {
		switch {
		case preserveLayout && (r == '\n' || r == '\t'):
			b.WriteRune(r)
		case !preserveLayout && r == '\n':
			b.WriteString(`\n`)
		case !preserveLayout && r == '\r':
			b.WriteString(`\r`)
		case !preserveLayout && r == '\t':
			b.WriteString(`\t`)
		case unicode.IsControl(r):
			fmt.Fprintf(b, "\\u%04x", r)
		case r == '\\':
			b.WriteString(`\\`)
		default:
			b.WriteRune(r)
		}
	}
}

// LineSafe reports whether value carries no control rune. It is the predicate a sink
// that writes raw lines uses before it embeds an operator-influenced value: a newline
// there forges a line and an ESC drives the terminal that prints it. Quoting does not
// substitute, because single quotes make a newline literal but still emit the byte, and
// escaping does not substitute either, because an escaped path names a tree that does
// not exist. A caller that fails this predicate emits a pointer instead of the value.
//
// Display-hostile runes outside the control categories — a bidi override, U+2028,
// invalid UTF-8 — pass. This guards line structure, not how a terminal renders one line.
func LineSafe(value string) bool { return !strings.ContainsFunc(value, unicode.IsControl) }
