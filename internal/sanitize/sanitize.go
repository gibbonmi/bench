// Package sanitize is the one policy for rendering operator-influenced text into a sink
// that could misread it. Two sinks, two duties. For a terminal render, Controls and its
// siblings escape control runes, so a control byte in a commit subject, an objective, or
// a git-sourced string is escaped one way everywhere. For a command line a reader will
// paste, ShellQuote renders a value as exactly one shell argument.
//
// It is deliberately distinct from internal/toon's cell policy: toon *refuses* a
// control-bearing cell (a closed AXI decision), while this package escapes or quotes.
// The two duties here do not substitute for each other — quoting leaves a control byte
// intact, and escaping produces text no longer naming what it came from.
package sanitize

import (
	"fmt"
	"strings"
	"unicode"
)

// Controls escapes every control rune in value — newline, carriage return, and tab to
// their backslash forms, every other control rune to \uXXXX, and a literal backslash to
// \\ — with no length cap. The result carries no raw control byte, so it is safe to hand
// to a terminal or to embed in html/template output.
func Controls(value string) string {
	var b strings.Builder
	writeEscaped(&b, []rune(value), false)
	return b.String()
}

// Preformatted is Controls for text that lands inside an html/template <pre> block: html/
// template already neutralizes markup there, so the only threat left is a raw control
// byte, and flattening the layout to escape it is unnecessary collateral. It leaves
// newline and tab verbatim so multi-line and tab-aligned content renders as authored,
// while every other control rune — including carriage return — still escapes through the
// same writeEscaped mechanism Controls uses, so no raw control byte reaches the output.
func Preformatted(value string) string {
	var b strings.Builder
	writeEscaped(&b, []rune(value), true)
	return b.String()
}

// Preview is Controls plus a 120-code-point cap and a byte-count suffix: it escapes the
// first 120 runes and, when value was longer, appends "… (N bytes)" naming the original
// byte length. The cap is counted in code points so a multibyte string is not cut at a
// fraction of its apparent length.
func Preview(value string) string {
	runes := []rune(value)
	truncated := len(runes) > 120
	if truncated {
		runes = runes[:120]
	}
	var b strings.Builder
	writeEscaped(&b, runes, false)
	if truncated {
		fmt.Fprintf(&b, "… (%d bytes)", len(value))
	}
	return b.String()
}

// ShellQuote renders value as exactly one POSIX shell argument: the whole string is
// single-quoted, and each embedded quote closes, escapes, and reopens it, so a space, a
// glob character, or a metacharacter inside can neither split the argument nor expand.
// An empty value becomes an empty quoted pair, so it survives as an argument rather
// than vanishing from the command line.
//
// Quoting is all it does. A control byte passes through as a literal byte — single
// quotes make a newline literal but still emit it — so a caller writing into a
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
// preserveLayout is true (Preformatted), newline and tab pass through verbatim and only
// carriage return joins the general control-rune case, so the \u%04x emission stays
// single-sourced regardless of which caller reaches it.
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
