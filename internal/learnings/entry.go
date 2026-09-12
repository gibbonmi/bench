package learnings

import (
	"strings"
	"unicode"
)

// FormatEntry renders one open journal entry in the shape Parse reads back: the dated
// heading with its `[open]` state marker, then the three body bullets. It is the one
// writer of that shape, so the `bench learning` verb, the adopt scaffold's worked
// example, and the parser cannot disagree about what an entry looks like. An empty
// rule renders as "none", the scaffold's own word for "no rule change proposed".
//
// Every field is sanitized first. A raw newline or another control byte in title, what,
// right, or rule would otherwise reach the heading line or a body bullet unescaped: an
// embedded newline in title splits the heading in two, and the surviving half then
// misses the `[open]` marker Parse requires, so the verb's own entry reads back as
// malformed. Sanitizing here, in the one writer, keeps this function's own guarantee
// instead of asking every caller to repeat the rule.
func FormatEntry(date, title, what, right, rule string) string {
	title, what, right, rule = sanitizeField(title), sanitizeField(what), sanitizeField(right), sanitizeField(rule)
	if rule == "" {
		rule = "none"
	}
	return "## " + date + " — " + title + "  [open]\n" +
		"- **What happened:** " + what + "\n" +
		"- **Right behavior:** " + right + "\n" +
		"- **Proposed rule change:** " + rule + "\n"
}

// sanitizeField collapses every control byte in s, a newline, a carriage return, a tab,
// or any other C0/C1 control character, to a single space, then collapses the resulting
// whitespace runs to one space each and trims the ends. FormatEntry is the journal's
// one-line-per-field grammar; this is the one place that grammar is enforced against
// whatever text a caller supplies.
func sanitizeField(s string) string {
	clean := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
	return strings.Join(strings.Fields(clean), " ")
}
