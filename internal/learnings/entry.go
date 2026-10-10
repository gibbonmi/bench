package learnings

import (
	"errors"
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
	values := []string{what, right, rule}
	out := "## " + date + " — " + title + "  [open]\n"
	for i, bullet := range entryBullets {
		out += "- **" + bullet.label + ":** " + values[i] + "\n"
	}
	return out
}

// resolvedBulletPrefix opens the one body bullet that links an open entry to the commit
// that fixed it. MarkResolved writes it and Parse reads it, so the two share this prefix.
const resolvedBulletPrefix = "- **Resolved by:** "

// resolvedBy returns the commit a body line names when the line is the resolved bullet.
func resolvedBy(line string) (string, bool) {
	value, ok := strings.CutPrefix(line, resolvedBulletPrefix)
	return strings.TrimSpace(value), ok
}

// Refusals MarkResolved returns when the title does not name exactly one unlinked open
// entry.
var (
	ErrNoOpenEntry     = errors.New("no open entry has that title")
	ErrAmbiguousEntry  = errors.New("more than one open entry has that title")
	ErrAlreadyResolved = errors.New("the entry already names a resolving commit")
)

// MarkResolved links the open entry whose title equals title to commit. It inserts the
// resolved bullet directly below the entry's heading, so the heading keeps its `[open]`
// marker and every other byte of content is unchanged. The drain still owns closing the
// entry. It returns a nil slice and one of the Err values above when no change is made.
func MarkResolved(content []byte, title, commit string) ([]byte, error) {
	var found []Entry
	for _, entry := range Entries(content) {
		if entry.Title == title {
			found = append(found, entry)
		}
	}
	switch {
	case len(found) == 0:
		return nil, ErrNoOpenEntry
	case len(found) > 1:
		return nil, ErrAmbiguousEntry
	case found[0].ResolvedBy != "":
		return nil, ErrAlreadyResolved
	}
	lines := strings.SplitAfter(string(content), "\n")
	heading := found[0].Line - 1
	terminator := "\n"
	if strings.HasSuffix(lines[heading], "\r\n") {
		terminator = "\r\n"
	}
	if !strings.HasSuffix(lines[heading], "\n") {
		lines[heading] += terminator
	}
	bullet := resolvedBulletPrefix + sanitizeField(commit) + terminator
	out := append(append(append([]string(nil), lines[:heading+1]...), bullet), lines[heading+1:]...)
	return []byte(strings.Join(out, "")), nil
}

// entryBullets is the body order FormatEntry renders, so EntryField reads the same order
// the writer uses.
var entryBullets = []struct{ field, label string }{
	{"what", "What happened"},
	{"right", "Right behavior"},
	{"rule", "Proposed rule change"},
}

// EntryField names the FormatEntry argument rendered on the entry's 1-based line:
// "title" for the heading, then "what", "right", or "rule". It returns "" for any other
// line.
func EntryField(line int) string {
	if line == 1 {
		return "title"
	}
	if i := line - 2; i >= 0 && i < len(entryBullets) {
		return entryBullets[i].field
	}
	return ""
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
