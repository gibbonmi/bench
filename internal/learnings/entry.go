package learnings

// FormatEntry renders one open journal entry in the shape Parse reads back: the dated
// heading with its `[open]` state marker, then the three body bullets. It is the one
// writer of that shape, so the `bench learning` verb, the adopt scaffold's worked
// example, and the parser cannot disagree about what an entry looks like. An empty
// rule renders as "none", the scaffold's own word for "no rule change proposed".
func FormatEntry(date, title, what, right, rule string) string {
	if rule == "" {
		rule = "none"
	}
	return "## " + date + " — " + title + "  [open]\n" +
		"- **What happened:** " + what + "\n" +
		"- **Right behavior:** " + right + "\n" +
		"- **Proposed rule change:** " + rule + "\n"
}
