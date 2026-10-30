package markdown

import "strings"

// CommentOpen and CommentClose delimit an HTML comment.
const (
	CommentOpen  = "<!--"
	CommentClose = "-->"
)

// cutComments returns the spans of src[from:] that the comment rule removes, in order, and the
// offset of an unterminated opener, or -1. An unterminated comment removes the rest of src.
//
// The result equals a rescan of the rewritten text after each removal. A removal can join the
// bytes on its two sides into a new opener, and that opener opens a comment too, so text that a
// removal rejoins never reads as visible. The earliest opener wins, and the search for its
// closer starts after the opener's four bytes. Only a removal can form an opener that src does
// not hold, so after each removal the pass tests only the bytes on its two sides, and then it
// searches src forward. The work and the allocation therefore stay linear in the input.
func cutComments(src string, from int) (cut cuts, open int) {
	at := from
	for {
		start, after := cut.joinedOpener(src, from, at)
		if start < 0 {
			k := strings.Index(src[at:], CommentOpen)
			if k < 0 {
				return cut, -1
			}
			start, after = at+k, at+k+len(CommentOpen)
		}
		k := strings.Index(src[after:], CommentClose)
		if k < 0 {
			cut.push(Span{Start: start, End: len(src)})
			return cut, start
		}
		at = after + k + len(CommentClose)
		cut.push(Span{Start: start, End: at})
	}
}

// cuts is the spans that the comment rule removes, in order. first[i] is the index of the
// earliest span of the run of adjacent spans that ends with spans[i]. A walk back crosses a whole
// run in one step, so the length of a run does not change the work of a walk. steps counts the
// steps of every walk, so a test can bound the work of the pass.
type cuts struct {
	spans []Span
	first []int
	steps int
}

// joinedOpener returns the offset of an opener whose first bytes are the kept bytes just before
// at and whose last bytes start src[at:], and the offset after that opener. It returns -1 when
// no such opener exists. The kept bytes before at are the bytes of src[from:at] outside c.
func (c *cuts) joinedOpener(src string, from, at int) (start, after int) {
	// kept holds the offsets of up to three kept bytes before at, the nearest first.
	var kept [len(CommentOpen) - 1]int
	n := 0
	for pos, top := at-1, len(c.spans)-1; n < len(kept) && pos >= from; c.steps++ {
		if top >= 0 && pos >= c.spans[top].Start && pos < c.spans[top].End {
			top = c.first[top]
			pos = c.spans[top].Start - 1
			top--
			continue
		}
		kept[n] = pos
		n++
		pos--
	}
	// window holds the kept bytes, the earliest first, and then the bytes that start src[at:].
	var window [2*len(CommentOpen) - 2]byte
	m := 0
	for i := n - 1; i >= 0; i-- {
		window[m] = src[kept[i]]
		m++
	}
	for i := at; i < len(src) && m < len(window); i++ {
		window[m] = src[i]
		m++
	}
	for p := 0; p < n && p+len(CommentOpen) <= m; p++ {
		if string(window[p:p+len(CommentOpen)]) == CommentOpen {
			return kept[n-1-p], at + len(CommentOpen) - (n - p)
		}
	}
	return -1, 0
}

// push adds s. Each earlier span that starts inside s is a removal that s rejoins, so s replaces
// it.
func (c *cuts) push(s Span) {
	for len(c.spans) > 0 && c.spans[len(c.spans)-1].Start >= s.Start {
		c.spans = c.spans[:len(c.spans)-1]
	}
	i := len(c.spans)
	c.first = c.first[:i]
	first := i
	if i > 0 && c.spans[i-1].End == s.Start {
		first = c.first[i-1]
	}
	c.spans = append(c.spans, s)
	c.first = append(c.first, first)
}

// keptText returns src[start:end] without the bytes of cut. cut is in order, and no span of it
// ends at or before start.
func keptText(src string, start, end int, cut []Span) string {
	var text strings.Builder
	for _, s := range cut {
		if s.Start >= end {
			break
		}
		if s.Start > start {
			text.WriteString(src[start:s.Start])
		}
		start = s.End
	}
	if start < end {
		text.WriteString(src[start:end])
	}
	return text.String()
}
