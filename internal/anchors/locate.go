package anchors

import (
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/markdown"
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
	var text []rune
	var textOrigin []int
	if kind.sectionScoped() {
		var count int
		text, textOrigin, count = commentFreeSectionMapped(data, section)
		// An absent section and a duplicated one both refuse: the evaluator resolves a
		// scoped anchor against exactly one owning heading.
		if count != 1 {
			return 0
		}
	} else {
		text, textOrigin = stripCommentsMapped(data)
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

// stripCommentsMapped is the package's one comment strip. It keeps each rune of data outside
// the comment spans that markdown.Read removes from data as a whole file, and each kept rune's
// index in data, so a caller can map a match back to its source position. The spans are byte
// ranges, but the origin is a rune index, because lineAtRune and the case fold count runes.
// StripHTMLComments is this function's string projection, so Locate and the evaluator read
// the same comment-free text.
func stripCommentsMapped(data string) (text []rune, origin []int) {
	cut := markdown.Read(data).Comments
	index := 0
	for at, r := range data {
		for len(cut) > 0 && at >= cut[0].End {
			cut = cut[1:]
		}
		if len(cut) == 0 || at < cut[0].Start {
			text = append(text, r)
			origin = append(origin, index)
		}
		index++
	}
	return text, origin
}

// commentFreeSectionMapped answers the comment-free runes of title's first section in data,
// each rune's index in data, and the number of owning headings. Locate and the evaluator read
// a section through it. The walk reads data itself, the document that the comment strip reads,
// so the strip and the walk take one frontmatter state. The body is each comment-free rune
// inside the section's range of data.
func commentFreeSectionMapped(data, title string) (body []rune, bodyOrigin []int, count int) {
	text, origin := stripCommentsMapped(data)
	runes := []rune(data)
	_, span, count := sectionRunesMapped(runes, identityOrigin(len(runes)), title)
	// An absent section and an empty one have no range of data.
	if len(span) == 0 {
		return nil, nil, count
	}
	from, _ := slices.BinarySearch(origin, span[0])
	to, _ := slices.BinarySearch(origin, span[len(span)-1]+1)
	return text[from:to], origin[from:to], count
}

// scopeRunesMapped is the package's one narrowing walk: it reads runes through the block
// reader once and answers the first region opens accepts, that region's origin mapping into
// runes, and the number of lines opens accepted. opens and closes read the reader's classified
// lines. The body and the count come from the same walk, so a caller that refuses a duplicated
// region reads the count the body came from. The body starts on the line after the opener.
// The walk reads runes as a whole document, so a comment opener or a fence marker inside the
// frontmatter hides no section. A caller passes the document itself, not its comment-free text,
// because a removed comment can turn a line of that text into a frontmatter marker.
func scopeRunesMapped(runes []rune, origin []int, opens, closes func(line markdown.Line) bool) (body []rune, bodyOrigin []int, count int) {
	text := string(runes)
	lines := markdown.Read(text).Lines
	// lineStart answers the index in runes where line i starts. A line offset is a byte offset.
	lineStart := func(i int) int { return utf8.RuneCountInString(text[:lines[i].Offset]) }
	start, end := -1, -1
	for i, line := range lines {
		if opens(line) {
			count++
			if count == 1 {
				start = i + 1
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
	startOffset := len(runes)
	if start < len(lines) {
		startOffset = lineStart(start)
	}
	endOffset := len(runes)
	if end < len(lines) {
		endOffset = lineStart(end)
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
// section, the number of owning headings, and the next H2 heading as the close. Only a line
// that the block reader reads as an H2 heading opens or closes a section, so a fenced or an
// indented heading line does neither. MarkdownH2Sections is this function's string projection.
func sectionRunesMapped(runes []rune, origin []int, title string) (body []rune, bodyOrigin []int, count int) {
	return scopeRunesMapped(runes, origin,
		func(line markdown.Line) bool { return line.Heading && line.Section == title },
		func(line markdown.Line) bool { return line.Heading })
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
