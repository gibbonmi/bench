// Fenced blocks in the handoff document. One file holds the whole rule, because a
// second reader with its own idea of a fence would disagree with the parser about
// which bytes are prose. Parse, the section split, and every writer that holds text
// bound for a section read the rule from here.

package handoffdoc

import (
	"iter"
	"strings"
)

// isFence reports whether a line opens or closes a fenced block. Both markdown
// fence characters count, and leading indentation does not disqualify one.
func isFence(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}

// OpenFenceRepair is what a refusal tells a writer holding text this package will not
// read. Nothing derives the fix, so the refusal states it.
const OpenFenceRepair = "close the fence, so the headings below it read as headings"

// OpenFence reports the line that opens a fenced block the text never closes, counting
// from one, and false when every fence closes. A writer holding text bound for a section
// calls it before the render: such a fence absorbs every heading below it, so the written
// document is a shape Parse refuses on the next run. The last unmatched opener wins,
// because a later pair closes the one before it.
func OpenFence(text string) (int, bool) {
	opened := 0
	fenced := false
	for i, raw := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		if !isFence(strings.TrimRight(raw, " \t")) {
			continue
		}
		if fenced {
			opened = 0
		} else {
			opened = i + 1
		}
		fenced = !fenced
	}
	return opened, opened != 0
}

// UnfencedLines yields the State lines that sit outside a fenced block, with each
// line's trailing whitespace removed. A line inside a fence is an example the
// writer pasted, not a claim about this repository, so no reader of State judges it.
//
// A State that breaks the document grammar is refused by Parse itself, which reads
// the file before any writer reaches this walk.
func UnfencedLines(state string) iter.Seq[string] {
	return func(yield func(string) bool) {
		fenced := false
		for _, raw := range strings.Split(state, "\n") {
			trimmed := strings.TrimRight(raw, " \t")
			if isFence(trimmed) {
				fenced = !fenced
				continue
			}
			if fenced {
				continue
			}
			if !yield(trimmed) {
				return
			}
		}
	}
}
