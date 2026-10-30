// Fenced blocks and comments in the handoff document. The block reader owns their rules, so
// the parser, the section split, and every writer that holds text bound for a section agree
// about which bytes are prose. This file reads a State body through that reader.

package handoffdoc

import (
	"fmt"
	"iter"
	"strings"

	"github.com/gibbonmi/bench/internal/markdown"
)

// OpenFenceRepair is what a refusal tells a writer holding text this package will not
// read. Nothing derives the fix, so the refusal states it.
const OpenFenceRepair = "close the fence, so the headings below it read as headings"

// OpenCommentRepair is the fix for a comment that the text never closes, for the same reason.
const OpenCommentRepair = "close the comment, so the headings below it read as headings"

// FaultRepair is the repair that a refusal states for a reader fault of kind: the comment
// repair for a comment fault, and the fence repair for every other kind.
func FaultRepair(kind markdown.FaultKind) string {
	if kind == markdown.CommentFault {
		return OpenCommentRepair
	}
	return OpenFenceRepair
}

// ReadState is the one read of a State body. The text is a fragment, so a leading marker line
// opens no frontmatter, and the fault is the comment or the fence that the body leaves open.
// Each line carries the heading rule that Parse applies to the document.
func ReadState(state string) markdown.Document {
	return markdown.ReadFragment(state)
}

// OpenFence reports the line that opens a fenced block the text never closes, counting
// from one, and false when every fence closes or an unterminated comment comes first. The
// text is a State body, so a leading marker line opens no frontmatter. A closer uses the
// opener's character and a run at least as long, so a shorter run or the other character
// leaves the block open.
func OpenFence(text string) (int, bool) {
	return fenceFault(ReadState(text))
}

// blockFault returns the comment or fence fault that Parse refuses in a reading of content.
// The reader reports only a frontmatter fault for an unclosed first-line marker, and that
// reading holds the marker line as body, so the fragment reading names the fault it hides.
func blockFault(content string, read markdown.Document) *markdown.Fault {
	if read.Fault != nil && read.Fault.Kind == markdown.FrontmatterFault {
		return markdown.ReadFragment(content).Fault
	}
	return read.Fault
}

// faultRefusal is the Parse refusal of a comment or fence fault, at the line of its opener.
func faultRefusal(path string, lines []markdown.Line, fault *markdown.Fault) *ParseError {
	opener := strings.TrimRight(lines[fault.Line-1].Raw, " \t")
	opens := "a fenced block"
	if fault.Kind == markdown.CommentFault {
		opens = "an HTML comment"
	}
	return &ParseError{path, fault.Line, fmt.Sprintf("line %q opens %s that the file never closes; %s", opener, opens, FaultRepair(fault.Kind))}
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
		for _, line := range ReadState(state).Lines {
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
