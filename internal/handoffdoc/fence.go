// Fenced blocks in the handoff document. The block reader owns the fence rule, so the
// parser, the section split, and every writer that holds text bound for a section agree
// about which bytes are prose. This file reads a State body through that reader.

package handoffdoc

import (
	"iter"
	"strings"

	"github.com/gibbonmi/bench/internal/markdown"
)

// OpenFenceRepair is what a refusal tells a writer holding text this package will not
// read. Nothing derives the fix, so the refusal states it.
const OpenFenceRepair = "close the fence, so the headings below it read as headings"

// OpenFence reports the line that opens a fenced block the text never closes, counting
// from one, and false when every fence closes. A writer holding text bound for a section
// calls it before the render: such a fence absorbs every heading below it, so the written
// document is a shape Parse refuses on the next run. The text is a State body, so a
// leading marker line opens no frontmatter. A closer uses the opener's character and a
// run at least as long, so a shorter run or the other character leaves the block open.
func OpenFence(text string) (int, bool) {
	return fenceFault(markdown.ReadFragment(text))
}

// fenceFault reports the opener line of a reading's fence fault. The reader reports a
// comment fault ahead of a fence fault, so a text with an unterminated comment reports no
// open fence here; the State file read refuses that comment on its own.
func fenceFault(read markdown.Document) (int, bool) {
	if read.Fault == nil || read.Fault.Kind != markdown.FenceFault {
		return 0, false
	}
	return read.Fault.Line, true
}

// UnfencedLines yields the State lines that sit outside a fenced block, with each
// line's trailing whitespace removed. A line inside a fence is an example the
// writer pasted, not a claim about this repository, so no reader of State judges it.
//
// A State that breaks the document grammar is refused by Parse itself, which reads
// the file before any writer reaches this walk.
func UnfencedLines(state string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, line := range markdown.ReadFragment(state).Lines {
			if line.Class != markdown.Body {
				continue
			}
			if !yield(strings.TrimRight(line.Raw, " \t")) {
				return
			}
		}
	}
}

// bodyText joins the raw text of lines and trims the empty lines at each end.
func bodyText(lines []markdown.Line) string {
	raw := make([]string, len(lines))
	for i, line := range lines {
		raw[i] = line.Raw
	}
	return strings.Trim(strings.Join(raw, "\n"), "\n")
}
