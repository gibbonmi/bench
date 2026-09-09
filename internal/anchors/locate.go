package anchors

import "unicode"

// Locate answers the 1-based physical line, in data as stored, of the first
// character of an anchor's first match, or 0 when the needle has no match,
// the section is absent or duplicated, or data is empty.
//
// The evaluator resolves an anchor's presence over comment-stripped,
// whitespace-collapsed, and (for a section kind) case-folded text; Locate
// answers a position over that same resolution, so the two stay in lock
// step. It maps a match in the transformed text back to data by rune index,
// not by byte offset, so a case fold that changes a rune's byte length
// cannot shift the reported line. A forbid kind is located the same way as
// a require kind: a non-zero line locates the violation, the presence of
// the forbidden needle.
func Locate(kind Kind, section, needle, data string) int {
	stripped, origin := stripCommentsMapped(data)
	text, textOrigin := stripped, origin
	if kind == RequireInSection || kind == ForbidInSection {
		body, bodyOrigin, ok := sectionRunesMapped(stripped, origin, section)
		if !ok {
			return 0
		}
		text, textOrigin = body, bodyOrigin
	}
	collapsed, collapsedOrigin := collapseSpaceMapped(text, textOrigin)
	searchText := collapsed
	searchNeedle := []rune(CollapseSpace(needle))
	if kind == RequireInSection || kind == ForbidInSection {
		searchText = toLowerRunes(collapsed)
		searchNeedle = toLowerRunes(searchNeedle)
	}
	at := indexRunes(searchText, searchNeedle)
	if at < 0 || at >= len(collapsedOrigin) {
		return 0
	}
	return lineAtRune(data, collapsedOrigin[at])
}

// commentOpen and commentClose delimit an HTML comment for stripCommentsMapped, the
// same two literals StripHTMLComments matches against. They are declared once here
// rather than allocated fresh on every call.
var commentOpen, commentClose = []rune("<!--"), []rune("-->")

// stripCommentsMapped strips HTML comments the way StripHTMLComments does — a complete
// comment is removed, and an unterminated one truncates the text — but over runes, and
// it returns each surviving rune's index in data alongside it.
func stripCommentsMapped(data string) (text []rune, origin []int) {
	runes := []rune(data)
	for i := 0; i < len(runes); {
		if runesEqual(runeWindow(runes, i, len(commentOpen)), commentOpen) {
			end := indexRunes(runeWindow(runes, i+len(commentOpen), len(runes)-i-len(commentOpen)), commentClose)
			if end < 0 {
				break
			}
			i = i + len(commentOpen) + end + len(commentClose)
			continue
		}
		text = append(text, runes[i])
		origin = append(origin, i)
		i++
	}
	return text, origin
}

// sectionRunesMapped resolves title's owning section within runes the same way
// MarkdownH2Sections does — headings inside a fenced block neither delimit nor count —
// and returns the body as a rune slice with its origin mapping into runes. ok is false
// when the section is absent or carries more than one owning heading, matching
// MarkdownH2Sections's own duplicate refusal.
func sectionRunesMapped(runes []rune, origin []int, title string) (body []rune, bodyOrigin []int, ok bool) {
	if _, count := MarkdownH2Sections(string(runes), title); count != 1 {
		return nil, nil, false
	}
	lines := splitRuneLines(runes)
	heading := "## " + title
	fenced := false
	start, end := -1, -1
	for i, line := range lines {
		trimmed := trimSpaceRunes(line)
		if hasPrefixRunes(trimmed, []rune("```")) {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		if string(trimmed) == heading {
			if start < 0 {
				start = i + 1
			}
			continue
		}
		if start >= 0 && end < 0 && hasPrefixRunes(line, []rune("## ")) {
			end = i
			break
		}
	}
	if start < 0 {
		return nil, nil, false
	}
	if end < 0 {
		end = len(lines)
	}
	lineStarts := runeLineStarts(runes)
	startOffset := len(runes)
	if start < len(lineStarts) {
		startOffset = lineStarts[start]
	}
	endOffset := len(runes)
	if end < len(lineStarts) {
		endOffset = lineStarts[end]
		if endOffset > startOffset {
			// Drop the boundary line's own leading newline from the body.
			endOffset--
		}
	}
	if endOffset < startOffset {
		endOffset = startOffset
	}
	return runes[startOffset:endOffset], origin[startOffset:endOffset], true
}

// collapseSpaceMapped collapses whitespace like CollapseSpace — each run becomes one
// ASCII space, and a leading or trailing run drops — and returns each output rune's
// origin in runes' own index space (already itself an origin into data, for a composed
// pipeline).
func collapseSpaceMapped(runes []rune, origin []int) (out []rune, outOrigin []int) {
	inField := false
	for i, r := range runes {
		if unicode.IsSpace(r) {
			inField = false
			continue
		}
		if !inField && len(out) > 0 {
			out = append(out, ' ')
			outOrigin = append(outOrigin, origin[i])
		}
		out = append(out, r)
		outOrigin = append(outOrigin, origin[i])
		inField = true
	}
	return out, outOrigin
}

func toLowerRunes(runes []rune) []rune {
	out := make([]rune, len(runes))
	for i, r := range runes {
		out[i] = unicode.ToLower(r)
	}
	return out
}

// lineAtRune answers the 1-based physical line of the rune at runeIdx in data, counting
// newlines strictly before it.
func lineAtRune(data string, runeIdx int) int {
	line := 1
	i := 0
	for _, r := range data {
		if i == runeIdx {
			break
		}
		if r == '\n' {
			line++
		}
		i++
	}
	return line
}

// indexRunes returns the first index in haystack where needle occurs, or -1. It
// searches by rune, never by byte, so a match position is stable across a case fold
// that changes a rune's byte length.
func indexRunes(haystack, needle []rune) int {
	if len(needle) == 0 {
		return 0
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if runesEqual(haystack[i:i+len(needle)], needle) {
			return i
		}
	}
	return -1
}

func runesEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// runeWindow returns runes[start:start+n], or nil when that window runs past the end.
func runeWindow(runes []rune, start, n int) []rune {
	if start < 0 || n < 0 || start+n > len(runes) {
		return nil
	}
	return runes[start : start+n]
}

func splitRuneLines(runes []rune) [][]rune {
	var lines [][]rune
	start := 0
	for i, r := range runes {
		if r == '\n' {
			lines = append(lines, runes[start:i])
			start = i + 1
		}
	}
	lines = append(lines, runes[start:])
	return lines
}

// runeLineStarts answers, for each line splitRuneLines would produce, that line's
// starting rune index in runes.
func runeLineStarts(runes []rune) []int {
	starts := []int{0}
	for i, r := range runes {
		if r == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

func trimSpaceRunes(runes []rune) []rune {
	start := 0
	for start < len(runes) && unicode.IsSpace(runes[start]) {
		start++
	}
	end := len(runes)
	for end > start && unicode.IsSpace(runes[end-1]) {
		end--
	}
	return runes[start:end]
}

func hasPrefixRunes(runes, prefix []rune) bool {
	if len(prefix) > len(runes) {
		return false
	}
	return runesEqual(runes[:len(prefix)], prefix)
}
