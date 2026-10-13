package probe

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// find is one requested edit before it matches: the form, the old string it selects, the
// text a swap or an insertion puts in, and the flag that gave a swap its text. An
// insertion's old string is its anchor.
type find struct {
	form     string
	old      string
	text     string
	textFlag string
}

// edit is one change to the start bytes: the byte range [start, end) and the bytes that take
// its place. An insertion is an empty range at the first byte of its anchor.
type edit struct {
	start, end  int
	replacement []byte
}

// mutationResult carries the mutant bytes and the miss classification together, so callers
// cannot recount the subject to tell a miss from a match. remove marks the --delete form,
// which has no mutant bytes.
type mutationResult struct {
	mutated []byte
	miss    bool
	remove  bool
}

// mutate resolves m to one edit list over the start bytes and applies the list once, so
// every form shares the match, overlap, and empty rules. Each old string matches the start
// bytes, never the result of an earlier edit. Several matches refuse unless --occurrence
// selects one, because the first match is a site the caller did not name.
func mutate(start []byte, m mutation, display string) (mutationResult, *refusal) {
	if m.remove {
		return mutationResult{remove: true}, nil
	}
	for _, f := range m.finds {
		if f.form == "swap" && f.text == f.old {
			return mutationResult{}, &refusal{title: "probe mutation empty", detail: f.textFlag + " equals --swap"}
		}
	}
	edits := make([]edit, 0, len(m.finds))
	for _, f := range m.finds {
		e, found, refused := f.locate(start, m.occurrence)
		if refused != nil {
			return mutationResult{}, refused
		}
		if !found {
			return mutationResult{miss: true}, nil
		}
		edits = append(edits, e)
	}
	mutated, refused := apply(start, edits)
	if refused != nil {
		return mutationResult{}, refused
	}
	// An empty file is not an absent one, so a mutant with zero bytes routes to --delete.
	if len(mutated) == 0 {
		return mutationResult{}, &refusal{title: "probe mutation empties the subject", detail: "use --delete to remove " + display + " for the run"}
	}
	return mutationResult{mutated: mutated}, nil
}

// locate answers the edit of the selected match of f in start: the one match, or the nth
// under --occurrence. found is false for an old string that start does not hold.
func (f find) locate(start []byte, occurrence int) (edit, bool, *refusal) {
	old := []byte(f.old)
	count := bytes.Count(start, old)
	switch {
	case count == 0:
		return edit{}, false, nil
	case occurrence == 0 && count > 1:
		hint := fmt.Sprintf("the old string matches %d times, want exactly 1", count)
		return edit{}, false, &refusal{title: "probe mutation ambiguous", detail: hint}
	case occurrence > count:
		hint := fmt.Sprintf("the old string matches %d times, --occurrence is %d", count, occurrence)
		return edit{}, false, &refusal{title: "probe mutation occurrence out of range", detail: hint}
	}
	at := 0
	for n := max(occurrence, 1); n > 1; n-- {
		at += bytes.Index(start[at:], old) + len(old)
	}
	at += bytes.Index(start[at:], old)
	switch f.form {
	case "insert":
		return edit{start: at, end: at, replacement: []byte(f.text)}, true, nil
	case "unwrap":
		inner, ok := unwrapCall(f.old)
		if !ok {
			return edit{}, false, &refusal{title: "probe unwrap invalid", detail: "--unwrap must name one outer call"}
		}
		return edit{start: at, end: at + len(old), replacement: []byte(inner)}, true, nil
	}
	return edit{start: at, end: at + len(old), replacement: []byte(f.text)}, true, nil
}

// apply writes the edits over start in byte order. Two edits whose ranges overlap refuse,
// because one would rewrite the bytes of the other.
func apply(start []byte, edits []edit) ([]byte, *refusal) {
	sort.Slice(edits, func(i, j int) bool { return edits[i].start < edits[j].start })
	var out []byte
	at := 0
	for _, e := range edits {
		if e.start < at {
			return nil, &refusal{title: "probe mutations overlap", detail: "two edits touch the same bytes"}
		}
		out = append(out, start[at:e.start]...)
		out = append(out, e.replacement...)
		at = e.end
	}
	return append(out, start[at:]...), nil
}

func unwrapCall(call string) (string, bool) {
	open := strings.IndexByte(call, '(')
	if open <= 0 || !strings.HasSuffix(call, ")") {
		return "", false
	}
	depth := 0
	for i := open; i < len(call); i++ {
		switch call[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(call)-1 {
				return "", false
			}
			if depth < 0 {
				return "", false
			}
		}
	}
	return call[open+1 : len(call)-1], depth == 0
}
