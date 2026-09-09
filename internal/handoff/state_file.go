package handoff

import (
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/toon"
)

// The three ways a drafted State never reaches the document, each named by what the writer
// has to fix. The read faults share one headline because the repair is the same for all
// four file kinds: name a file this command can read. The two text faults are separate,
// because a byte to remove and a line to reword are different edits.
const (
	faultStateFileUnreadable = "the handoff state file is not a readable regular file"
	faultStateFileByte       = "the handoff state file carries a control byte"
	faultStateFileHeading    = "the handoff state file opens a section heading"
)

const (
	stateFileRepair = "name a readable regular file, then rerun bench handoff: "
	stateByteRepair = "remove the control byte, then rerun bench handoff: "
)

// sectionOpener is the prefix that opens a level-two block, derived from the document's
// own main heading rather than spelled again here. handoffdoc splits the file on this
// prefix, and any spelling below it that the grammar has no key for makes every later run
// refuse the document. So the test is the prefix, not the two headings this command writes.
var sectionOpener = strings.TrimSuffix(handoffdoc.MainHeading, handoffdoc.MainKey)

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
	for line := range handoffdoc.UnfencedLines(state) {
		if strings.HasPrefix(line, sectionOpener) {
			return "", refusal{faultStateFileHeading, stateRepair + strings.TrimSpace(line)}
		}
	}
	return state, nil
}
