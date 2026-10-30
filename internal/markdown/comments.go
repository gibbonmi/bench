package markdown

import "strings"

// commentOpen and commentClose delimit an HTML comment.
const (
	commentOpen  = "<!--"
	commentClose = "-->"
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
func cutComments(src string, from int) (cut []Span, open int) {
	at := from
	for {
		start, after := joinedOpener(src, from, at, cut)
		if start < 0 {
			k := strings.Index(src[at:], commentOpen)
			if k < 0 {
				return cut, -1
			}
			start, after = at+k, at+k+len(commentOpen)
		}
		k := strings.Index(src[after:], commentClose)
		if k < 0 {
			return pushCut(cut, Span{Start: start, End: len(src)}), start
		}
		at = after + k + len(commentClose)
		cut = pushCut(cut, Span{Start: start, End: at})
	}
}

// joinedOpener returns the offset of an opener whose first bytes are the kept bytes just before
// at and whose last bytes start src[at:], and the offset after that opener. It returns -1 when
// no such opener exists. The kept bytes before at are the bytes of src[from:at] outside cut.
func joinedOpener(src string, from, at int, cut []Span) (start, after int) {
	// kept holds the offsets of up to three kept bytes before at, the nearest first.
	var kept [len(commentOpen) - 1]int
	n := 0
	for pos, top := at-1, len(cut)-1; n < len(kept) && pos >= from; {
		if top >= 0 && pos >= cut[top].Start && pos < cut[top].End {
			pos = cut[top].Start - 1
			top--
			continue
		}
		kept[n] = pos
		n++
		pos--
	}
	// window holds the kept bytes, the earliest first, and then the bytes that start src[at:].
	var window [2*len(commentOpen) - 2]byte
	m := 0
	for i := n - 1; i >= 0; i-- {
		window[m] = src[kept[i]]
		m++
	}
	for i := at; i < len(src) && m < len(window); i++ {
		window[m] = src[i]
		m++
	}
	for p := 0; p < n && p+len(commentOpen) <= m; p++ {
		if string(window[p:p+len(commentOpen)]) == commentOpen {
			return kept[n-1-p], at + len(commentOpen) - (n - p)
		}
	}
	return -1, 0
}

// pushCut adds s to cut. Each earlier span that starts inside s is a removal that s rejoins, so
// s replaces it.
func pushCut(cut []Span, s Span) []Span {
	for len(cut) > 0 && cut[len(cut)-1].Start >= s.Start {
		cut = cut[:len(cut)-1]
	}
	return append(cut, s)
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
