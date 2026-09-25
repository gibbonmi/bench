package responsebound

import "bytes"

// lineTracker counts the lines of a response and keeps the lines its projection prints:
// the first head lines, and a ring of the last tail lines. It keeps nothing else, so the
// memory of a spilled response does not grow with the response's line count.
//
// A line ends at a newline byte, and a final line without one still counts. The tracker
// copies bytes and never decodes them, so binary bytes stay exact.
type lineTracker struct {
	count int
	// open reports that the last line has no newline yet.
	open bool
	head [][]byte
	ring [][]byte
}

func (t *lineTracker) add(p []byte) {
	for len(p) > 0 {
		if !t.open {
			t.startLine()
		}
		chunk := p
		if end := bytes.IndexByte(p, '\n'); end >= 0 {
			chunk, t.open = p[:end+1], false
		}
		if t.count <= headLines {
			last := len(t.head) - 1
			t.head[last] = append(t.head[last], chunk...)
		}
		last := len(t.ring) - 1
		t.ring[last] = append(t.ring[last], chunk...)
		p = p[len(chunk):]
	}
}

// startLine opens the next line. A full ring reuses its oldest buffer for the new line.
func (t *lineTracker) startLine() {
	t.count++
	t.open = true
	if t.count <= headLines {
		t.head = append(t.head, nil)
	}
	if len(t.ring) < tailLines {
		t.ring = append(t.ring, nil)
		return
	}
	oldest := t.ring[0][:0]
	copy(t.ring, t.ring[1:])
	t.ring[len(t.ring)-1] = oldest
}

// tailCount is how many last lines the projection prints: the ring, less any line the
// head already holds.
func (t *lineTracker) tailCount() int {
	return min(len(t.ring), t.count-len(t.head))
}

// tail answers the bytes of the printed last lines.
func (t *lineTracker) tail() []byte {
	var out []byte
	for _, line := range t.ring[len(t.ring)-t.tailCount():] {
		out = append(out, line...)
	}
	return out
}

// omitted is how many lines the projection leaves to the spill file alone.
func (t *lineTracker) omitted() int {
	return t.count - len(t.head) - t.tailCount()
}
