package anchors

import (
	"slices"
	"unicode"
)

// Locate answers the 1-based physical line, in data as stored, of the first
// character of an anchor's first match, or 0 when the needle has no match,
// the section is absent or duplicated, or data is empty.
//
// The evaluator resolves an anchor's presence over comment-stripped and
// whitespace-collapsed text. For ForbidCaseFoldedEmphasis, it also removes
// ordinary emphasis and case-folds whole-file text. Locate uses that same
// resolution, so the two stay in lock step. It maps a match in the transformed
// text back to data by rune index, not by byte offset. A case fold that changes
// a rune's byte length cannot shift the reported line. A forbid kind is located
// the same way as a require kind: a non-zero line locates the violation, the
// presence of the forbidden needle.
func Locate(kind Kind, section, needle, data string) int {
	return locate(kind, section, 0, needle, data)
}

// locate is the package's one location walk. Locate is its step-free projection, and the
// evaluator calls it with the anchor's own step, so a step-scoped anchor reports the line
// of the match inside its step and nothing else.
func locate(kind Kind, section string, step int, needle, data string) int {
	stripped, origin := stripCommentsMapped(data)
	text, textOrigin := stripped, origin
	if kind.sectionScoped() {
		body, bodyOrigin, count := sectionRunesMapped(stripped, origin, section)
		// An absent section and a duplicated one both refuse: the evaluator resolves a
		// scoped anchor against exactly one owning heading.
		if count != 1 {
			return 0
		}
		text, textOrigin = body, bodyOrigin
	}
	if kind.stepScoped() {
		body, bodyOrigin, count := stepRunesMapped(text, textOrigin, step)
		// An absent step and a duplicated one refuse for the same reason a section does.
		// The kind alone decides that this narrowing runs, so the locator and the evaluator
		// cannot disagree about which anchors are step-scoped.
		if count != 1 {
			return 0
		}
		text, textOrigin = body, bodyOrigin
	}
	searchText, searchOrigin := normalizeMatchMapped(kind, text, textOrigin)
	searchNeedleRunes := []rune(needle)
	searchNeedle, _ := normalizeMatchMapped(kind, searchNeedleRunes, identityOrigin(len(searchNeedleRunes)))
	at := indexRunes(searchText, searchNeedle)
	if at < 0 || at >= len(searchOrigin) {
		return 0
	}
	return lineAtRune(data, searchOrigin[at])
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

// scopeRunesMapped is the package's one narrowing walk: it walks the lines once and answers
// the first region opens accepts, that region's origin mapping into runes, and the number of
// lines opens accepted. A line inside a fenced block neither opens nor closes a region. The
// body and the count come from the same walk, so a caller that refuses a duplicated region
// reads the count the body came from. keepOpener keeps the opening line inside the body, for
// a region whose own first line carries text.
func scopeRunesMapped(runes []rune, origin []int, opens, closes func(line []rune) bool, keepOpener bool) (body []rune, bodyOrigin []int, count int) {
	lines := splitRuneLines(runes)
	fenced := false
	start, end := -1, -1
	for i, line := range lines {
		if hasPrefixRunes(trimSpaceRunes(line), fenceMark) {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		if opens(line) {
			count++
			if count == 1 {
				start = i
				if !keepOpener {
					start = i + 1
				}
			}
			continue
		}
		if start >= 0 && end < 0 && closes(line) {
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

// sectionRunesMapped is the package's one H2 section resolution: title's first owning
// section, the number of owning headings, and the next H2 heading as the close. A heading
// inside a fenced block neither delimits a section nor counts as one.
// MarkdownH2Sections is this function's string projection.
func sectionRunesMapped(runes []rune, origin []int, title string) (body []rune, bodyOrigin []int, count int) {
	heading := []rune("## " + title)
	return scopeRunesMapped(runes, origin,
		func(line []rune) bool { return slices.Equal(trimSpaceRunes(line), heading) },
		func(line []rune) bool { return hasPrefixRunes(line, headingMark) },
		false)
}

// stepRunesMapped is the package's one numbered-step resolution: the first body of the step
// the reader sees as step, and the number of lines that open it. The next opener of any step
// closes the body, so an indented continuation line stays inside its step. The opener's own
// line joins the body, because that line carries the step's first words.
// MarkdownNumberedSteps is this function's string projection.
func stepRunesMapped(runes []rune, origin []int, step int) (body []rune, bodyOrigin []int, count int) {
	return scopeRunesMapped(runes, origin,
		// The reader sees the literal digits, so the match is on the written number and
		// never on the opener's ordinal position in the section.
		func(line []rune) bool { number, opens := stepOpener(line); return opens && number == step },
		func(line []rune) bool { _, opens := stepOpener(line); return opens },
		true)
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

func normalizeMatchMapped(kind Kind, runes []rune, origin []int) ([]rune, []int) {
	if kind == ForbidCaseFoldedEmphasis {
		runes, origin = stripMarkdownEmphasisMapped(runes, origin)
	}
	runes, origin = collapseSpaceMapped(runes, origin)
	if kind.sectionScoped() || kind == ForbidCaseFoldedEmphasis {
		runes = toLowerRunes(runes)
	}
	return runes, origin
}

// stripMarkdownEmphasisMapped removes paired emphasis markers at word boundaries.
// Intraword underscores remain visible, so identifiers keep their exact spelling.
// Each delimiter run enters and leaves the stack at most once.
func stripMarkdownEmphasisMapped(runes []rune, origin []int) (out []rune, outOrigin []int) {
	type opener struct {
		at, width int
		marker    rune
		previous  int
	}
	var stack []opener
	top := make(map[rune]int)
	removed := make([]bool, len(runes))
	pop := func() {
		last := stack[len(stack)-1]
		top[last.marker] = last.previous
		stack = stack[:len(stack)-1]
	}
	for i := 0; i < len(runes); {
		width := emphasisMarkerWidth(runes, i)
		if width == 0 {
			i++
			continue
		}
		marker := runes[i]
		if _, seen := top[marker]; !seen {
			top[marker] = -1
		}
		remaining := width
		canClose := i > 0 && !unicode.IsSpace(runes[i-1]) &&
			(i+width == len(runes) || unicode.IsSpace(runes[i+width]) || unicode.IsPunct(runes[i+width]))
		for canClose && remaining > 0 && top[marker] >= 0 {
			open := top[marker]
			for len(stack)-1 > open {
				pop()
			}
			paired := min(remaining, stack[open].width)
			for offset := 0; offset < paired; offset++ {
				removed[stack[open].at+stack[open].width-1-offset] = true
				removed[i+width-remaining+offset] = true
			}
			stack[open].width -= paired
			remaining -= paired
			if stack[open].width == 0 {
				pop()
			}
		}
		if remaining > 0 && emphasisCanOpen(runes, i, width) {
			stack = append(stack, opener{i + width - remaining, remaining, marker, top[marker]})
			top[marker] = len(stack) - 1
		}
		i += width
	}
	for i, r := range runes {
		if !removed[i] {
			out = append(out, r)
			outOrigin = append(outOrigin, origin[i])
		}
	}
	return out, outOrigin
}

func emphasisMarkerWidth(runes []rune, at int) int {
	if at >= len(runes) || runes[at] != '*' && runes[at] != '_' {
		return 0
	}
	if escapedAt(runes, at) {
		return 0
	}
	if at > 0 && runes[at-1] == runes[at] {
		return 0
	}
	end := at
	for end < len(runes) && runes[end] == runes[at] {
		end++
	}
	if width := end - at; width >= 1 && width <= 3 {
		return width
	}
	return 0
}

func escapedAt(runes []rune, at int) bool {
	backslashes := 0
	for i := at - 1; i >= 0 && runes[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 1
}

func emphasisCanOpen(runes []rune, at, width int) bool {
	if at+width >= len(runes) || unicode.IsSpace(runes[at+width]) {
		return false
	}
	return at == 0 || unicode.IsSpace(runes[at-1]) || unicode.IsPunct(runes[at-1])
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
