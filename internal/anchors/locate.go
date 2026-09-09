package anchors

import (
	"slices"
	"unicode"
)

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
		body, bodyOrigin, count := sectionRunesMapped(stripped, origin, section)
		// An absent section and a duplicated one both refuse: the evaluator resolves a
		// scoped anchor against exactly one owning heading.
		if count != 1 {
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

// commentOpen and commentClose delimit an HTML comment, fenceMark opens a fenced block,
// and headingMark opens an H2 heading. They are declared once here rather than allocated
// fresh on every call.
var (
	commentOpen, commentClose = []rune("<!--"), []rune("-->")
	fenceMark                 = []rune("```")
	headingMark               = []rune("## ")
)

// stripCommentsMapped is the package's one comment strip: a complete comment is removed,
// an unterminated one truncates the text, and the scan restarts on the rewritten text. It
// works over runes and returns each surviving rune's index in data alongside it, so a
// caller can map a match back to its source position; StripHTMLComments is this function's
// string projection. The restart is what makes a comment that the removal itself creates,
// from text either side of a removed comment, hide its needle from Locate as it hides it
// from the evaluator. A single pass would locate a needle the evaluator reads as absent.
func stripCommentsMapped(data string) (text []rune, origin []int) {
	text = []rune(data)
	origin = identityOrigin(len(text))
	for {
		start := indexRunes(text, commentOpen)
		if start < 0 {
			return text, origin
		}
		end := indexRunes(text[start+len(commentOpen):], commentClose)
		if end < 0 {
			return text[:start], origin[:start]
		}
		cut := start + len(commentOpen) + end + len(commentClose)
		text = append(text[:start:start], text[cut:]...)
		origin = append(origin[:start:start], origin[cut:]...)
	}
}

// sectionRunesMapped is the package's one H2 section resolution: it walks the lines once
// and answers title's first owning section as a rune slice, its origin mapping into runes,
// and the number of owning headings. A heading inside a fenced block neither delimits a
// section nor counts as one. The body and the count come from the same walk, so a caller
// that refuses a duplicated section reads the count the body came from.
// MarkdownH2Sections is this function's string projection.
func sectionRunesMapped(runes []rune, origin []int, title string) (body []rune, bodyOrigin []int, count int) {
	lines := splitRuneLines(runes)
	heading := []rune("## " + title)
	fenced := false
	start, end := -1, -1
	for i, line := range lines {
		trimmed := trimSpaceRunes(line)
		if hasPrefixRunes(trimmed, fenceMark) {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		if slices.Equal(trimmed, heading) {
			count++
			if count == 1 {
				start = i + 1
			}
			continue
		}
		if start >= 0 && end < 0 && hasPrefixRunes(line, headingMark) {
			end = i
		}
	}
	if count == 0 {
		return nil, nil, 0
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
	return runes[startOffset:endOffset], origin[startOffset:endOffset], count
}

// collapseSpaceMapped is the package's one whitespace collapse: each whitespace run becomes
// one ASCII space, and a leading or trailing run drops. It returns each output rune's origin
// in runes' own index space, which is itself already an origin into data for a composed
// pipeline. CollapseSpace is this function's string projection.
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

// identityOrigin answers the origin map of a rune slice that has had nothing removed yet:
// rune i comes from index i. A string-form caller starts a mapped transform with it and
// then discards the mapping the transform returns.
func identityOrigin(n int) []int {
	origin := make([]int, n)
	for i := range origin {
		origin[i] = i
	}
	return origin
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
		if slices.Equal(haystack[i:i+len(needle)], needle) {
			return i
		}
	}
	return -1
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

// hasPrefixRunes reports whether runes starts with prefix.
func hasPrefixRunes(runes, prefix []rune) bool {
	if len(prefix) > len(runes) {
		return false
	}
	return slices.Equal(runes[:len(prefix)], prefix)
}
