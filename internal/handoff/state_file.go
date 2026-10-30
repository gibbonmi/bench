package handoff

import (
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/markdown"
	"github.com/gibbonmi/bench/internal/toon"
)

// The five ways a drafted State never reaches the document, each named by what the writer
// has to fix. The read faults share one headline because the repair is the same for all
// four file kinds: name a file this command can read. The four text faults are separate,
// because a byte to remove, a fence or a comment to close, and a line to reword are
// different edits.
const (
	faultStateFileUnreadable = "the handoff state file is not a readable regular file"
	faultStateFileByte       = "the handoff state file carries a control byte"
	faultStateFileFence      = "the handoff state file opens a fence it never closes"
	faultStateFileComment    = "the handoff state file opens an HTML comment it never closes"
	faultStateFileHeading    = "the handoff state file opens a section heading"
)

// Each repair ends in stateRerun. The fence and comment repairs take the fix from the
// document parser, which owns the words. The draft and the document refuse the same fence
// or comment for the same reason, so a reworded instruction there carries here.
const (
	stateRerun      = ", then rerun bench handoff: "
	stateFileRepair = "name a readable regular file" + stateRerun
	stateByteRepair = "remove the control byte" + stateRerun
)

// readStateFile answers the State body a `--state-file` run writes, or the refusal that
// ends the run before anything is written.
//
// The classification is ClassifyNoFollow, never Classify. The following form resolves a
// live link, so a linked path would read bytes outside the file the caller named. The
// type check there also precedes every open, so a FIFO cannot block the command in open(2),
// and a file past the control-record limit answers unreadable rather than a truncated body.
//
// Failed covers the unreadable, wrong-type, and malformed states alone: absence is an
// authoritative answer to that classifier, and a missing draft is a typo rather than a
// deliberate empty State, so it is named beside them. Emptiness stays outside the refusal,
// because an empty draft is the deliberate reset that brings the scaffold guidance back.
func readStateFile(path string) (string, error) {
	file := bounds.ClassifyNoFollow(path)
	if file.State == bounds.StateAbsent || file.State.Failed() {
		return "", refusal{faultStateFileUnreadable, stateFileRepair + path}
	}

	// One trailing newline goes, the way the document's own writer trims one. A draft an
	// editor terminated and one it did not then write the same bytes.
	state := strings.TrimSuffix(string(file.Data), "\n")

	// The predicate the derived fields answer to, widened to the body. State is a block
	// rather than a label line, so the newline and the tab it permits are content here;
	// an escape byte is what would ride into every downstream reader of the artifact.
	if !toon.Representable(state) {
		return "", refusal{faultStateFileByte, stateByteRepair + path}
	}
	// The heading walk below reads a fence and a comment that close. An open one hides every
	// heading under it from that walk and from the document the next run parses, so it
	// refuses first, through the parser's own read. The reader reports a comment ahead of a
	// fence, so the comment refusal also covers a fence that the comment hides. A State opens
	// no frontmatter, so every other fault is a fence.
	read := handoffdoc.ReadState(state)
	if fault := read.Fault; fault != nil {
		kind := faultStateFileFence
		if fault.Kind == markdown.CommentFault {
			kind = faultStateFileComment
		}
		return "", refusal{kind, handoffdoc.FaultRepair(fault.Kind) + stateRerun + strings.TrimSpace(read.Lines[fault.Line-1].Raw)}
	}
	// handoffdoc splits the file on each line that the reader reads as an H2 heading, and any
	// spelling that the grammar has no key for makes every later run refuse the document. So
	// the test is the reader's heading, not the two headings this command writes.
	for _, line := range read.Lines {
		if line.Heading {
			return "", refusal{faultStateFileHeading, stateRepair + strings.TrimSpace(line.Raw)}
		}
	}
	return state, nil
}
