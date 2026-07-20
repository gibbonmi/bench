// Package sanitize is the one control-sequence policy for every terminal render of
// operator-influenced text. It replaces the near-copies that used to live in
// internal/intent and internal/dashboard, so a control byte in a commit subject, an
// objective, or a git-sourced string is escaped one way everywhere. It is deliberately
// distinct from internal/toon's cell policy: toon *refuses* a control-bearing cell (a
// closed AXI decision), while this package *escapes* control runes into a readable form
// for a human-facing string.
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
	writeEscaped(&b, []rune(value))
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
	writeEscaped(&b, runes)
	if truncated {
		fmt.Fprintf(&b, "… (%d bytes)", len(value))
	}
	return b.String()
}

// writeEscaped is the single per-rune escaping rule both Controls and Preview share.
func writeEscaped(b *strings.Builder, runes []rune) {
	for _, r := range runes {
		switch {
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
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
