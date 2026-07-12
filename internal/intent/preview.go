package intent

import (
	"fmt"
	"strings"
	"unicode"
)

// Preview safely encodes terminal context and caps it at 120 Unicode code points.
func Preview(value string) string {
	runes := []rune(value)
	truncated := len(runes) > 120
	if truncated {
		runes = runes[:120]
	}
	var b strings.Builder
	for _, r := range runes {
		switch {
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case unicode.IsControl(r):
			fmt.Fprintf(&b, "\\u%04x", r)
		case r == '\\':
			b.WriteString(`\\`)
		default:
			b.WriteRune(r)
		}
	}
	if truncated {
		fmt.Fprintf(&b, "… (%d bytes)", len(value))
	}
	return b.String()
}
