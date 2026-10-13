package probe

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

// find is one requested edit before it matches: the form, the old string it selects, and
// the text a swap or an insertion puts in. An insertion's old string is its anchor.
type find struct {
	form string
	old  string
	text string
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
func mutate(start []byte, m mutation, display string) (mutationResult, string) {
	if m.remove {
		return mutationResult{remove: true}, ""
	}
	for _, f := range m.finds {
		if f.form == "swap" && f.text == f.old {
			return mutationResult{}, toon.Errorf("probe mutation empty", "--with equals --swap") + "\n"
		}
	}
	edits := make([]edit, 0, len(m.finds))
	for _, f := range m.finds {
		e, found, line := f.locate(start, m.occurrence)
		if line != "" {
			return mutationResult{}, line
		}
		if !found {
			return mutationResult{miss: true}, ""
		}
		edits = append(edits, e)
	}
	mutated, line := apply(start, edits)
	if line != "" {
		return mutationResult{}, line
	}
	// An empty file is not an absent one, so a mutant with zero bytes routes to --delete.
	if len(mutated) == 0 {
		return mutationResult{}, toon.Errorf("probe mutation empties the subject", "use --delete to remove "+display+" for the run") + "\n"
	}
	return mutationResult{mutated: mutated}, ""
}

// locate answers the edit of the selected match of f in start: the one match, or the nth
// under --occurrence. found is false for an old string that start does not hold.
func (f find) locate(start []byte, occurrence int) (edit, bool, string) {
	old := []byte(f.old)
	count := bytes.Count(start, old)
	switch {
	case count == 0:
		return edit{}, false, ""
	case occurrence == 0 && count > 1:
		hint := fmt.Sprintf("the old string matches %d times, want exactly 1", count)
		return edit{}, false, toon.Errorf("probe mutation ambiguous", hint) + "\n"
	case occurrence > count:
		hint := fmt.Sprintf("the old string matches %d times, --occurrence is %d", count, occurrence)
		return edit{}, false, toon.Errorf("probe mutation occurrence out of range", hint) + "\n"
	}
	at := 0
	for n := max(occurrence, 1); n > 1; n-- {
		at += bytes.Index(start[at:], old) + len(old)
	}
	at += bytes.Index(start[at:], old)
	switch f.form {
	case "insert":
		return edit{start: at, end: at, replacement: []byte(f.text)}, true, ""
	case "unwrap":
		inner, ok := unwrapCall(f.old)
		if !ok {
			return edit{}, false, toon.Errorf("probe unwrap invalid", "--unwrap must name one outer call") + "\n"
		}
		return edit{start: at, end: at + len(old), replacement: []byte(inner)}, true, ""
	}
	return edit{start: at, end: at + len(old), replacement: []byte(f.text)}, true, ""
}

// apply writes the edits over start in byte order. Two edits whose ranges overlap refuse,
// because one would rewrite the bytes of the other.
func apply(start []byte, edits []edit) ([]byte, string) {
	sort.Slice(edits, func(i, j int) bool { return edits[i].start < edits[j].start })
	var out []byte
	at := 0
	for _, e := range edits {
		if e.start < at {
			return nil, toon.Errorf("probe mutations overlap", "two edits touch the same bytes") + "\n"
		}
		out = append(out, start[at:e.start]...)
		out = append(out, e.replacement...)
		at = e.end
	}
	return append(out, start[at:]...), ""
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
