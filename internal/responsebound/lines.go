package responsebound

import (
	"bytes"
	"unicode/utf8"
)

// lineTracker counts the lines of a response and keeps the lines its projection prints:
// the first head lines, and a ring of the last tail lines. It keeps nothing else, and it
// keeps at most lineCut bytes of each line, so the memory of a spilled response grows
// neither with its line count nor with the length of one line.
//
// A line ends at a newline byte, and a final line without one still counts. The tracker
// copies bytes and decodes them only to find the cut, so binary bytes stay exact.
type lineTracker struct {
	count int
	// open reports that the last line has no newline yet.
	open bool
	head []line
	ring []line
}

// line is one kept line: its first content bytes, without the newline, and the full
// count of its content bytes.
type line struct {
	kept []byte
	size int
}

func (t *lineTracker) add(p []byte) {
	for len(p) > 0 {
		if !t.open {
			t.startLine()
		}
		content, rest := p, []byte(nil)
		if end := bytes.IndexByte(p, '\n'); end >= 0 {
			content, rest, t.open = p[:end], p[end+1:], false
		}
		if t.count <= headLines {
			t.head[len(t.head)-1].append(content)
		}
		t.ring[len(t.ring)-1].append(content)
		p = rest
	}
}

// append keeps the part of content that fits the cut, and counts all of it. A buffer grows
// to the cut at most, so a line never holds more than lineCut bytes.
func (l *line) append(content []byte) {
	l.size += len(content)
	take := content[:min(len(content), lineCut-len(l.kept))]
	if need := len(l.kept) + len(take); need > cap(l.kept) {
		grown := make([]byte, len(l.kept), min(lineCut, max(need, 2*cap(l.kept))))
		copy(grown, l.kept)
		l.kept = grown
	}
	l.kept = append(l.kept, take...)
}

// cut reports that the line is longer than the cut.
func (l line) cut() bool { return l.size > lineCut }

// printed answers the bytes the projection prints for the line, without its newline. A
// cut line whose kept bytes end inside the encoding of one rune ends before that rune
// instead, so the cut never splits a rune. Bytes that encode no rune cut at lineCut.
func (l line) printed() []byte {
	if !l.cut() {
		return l.kept
	}
	for i := len(l.kept) - 1; i >= max(0, len(l.kept)-utf8.UTFMax+1); i-- {
		if utf8.RuneStart(l.kept[i]) {
			if !utf8.FullRune(l.kept[i:]) {
				return l.kept[:i]
			}
			break
		}
	}
	return l.kept
}

// startLine opens the next line. A full ring reuses its oldest buffer for the new line.
func (t *lineTracker) startLine() {
	t.count++
	t.open = true
	if t.count <= headLines {
		t.head = append(t.head, line{})
	}
	if len(t.ring) < tailLines {
		t.ring = append(t.ring, line{})
		return
	}
	oldest := line{kept: t.ring[0].kept[:0]}
	copy(t.ring, t.ring[1:])
	t.ring[len(t.ring)-1] = oldest
}

// tail answers the printed last lines: the ring, less any line the head already holds.
func (t *lineTracker) tail() []line {
	return t.ring[len(t.ring)-min(len(t.ring), t.count-len(t.head)):]
}

// headBytes answers the printed head. Each head line ends with its newline, so the spill
// line that follows the head starts a line of its own.
func (t *lineTracker) headBytes() []byte {
	var out []byte
	for _, l := range t.head {
		out = append(append(out, l.printed()...), '\n')
	}
	return out
}

// tailBytes answers the printed tail. The last line keeps no newline when the response
// ends without one.
func (t *lineTracker) tailBytes() []byte {
	var out []byte
	tail := t.tail()
	for _, l := range tail {
		out = append(append(out, l.printed()...), '\n')
	}
	if t.open && len(tail) > 0 {
		out = out[:len(out)-1]
	}
	return out
}

// omitted is how many lines the projection leaves to the spill file alone.
func (t *lineTracker) omitted() int {
	return t.count - len(t.head) - len(t.tail())
}

// cutLines is how many printed lines the projection cuts.
func (t *lineTracker) cutLines() int {
	cut := 0
	for _, printed := range [][]line{t.head, t.tail()} {
		for _, l := range printed {
			if l.cut() {
				cut++
			}
		}
	}
	return cut
}
