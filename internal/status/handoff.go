package status

import (
	"fmt"
	"strconv"

	"github.com/gibbonmi/bench/internal/git"
)

// handoffFile is the phase-close continuation artifact `AGENTS.md` names: the document a
// cold session reads to resume without conversation history. Its shape is documented
// inside the file itself so the template cannot drift from the artifact it describes.
const handoffFile = "session-handoff.md"

// appendHandoff adds the handoff-staleness signal (sev 11): the commits that landed since
// the handoff was last written. It turns "trust git over the handoff" from a rule the next
// session must recall into an ambient fact the SessionStart hook already prints.
//
// The handoff's age is read from git history rather than from a line inside the document.
// A self-reported date is the same remembered-not-computed defect this signal exists to
// close — a session that rewrites the body and forgets the date leaves the check
// confidently calling a stale document current — and it cannot name its own commit anyway,
// since the handoff lands in the commit it describes.
//
// It ranks last deliberately. The value is at cold pickup, where the board is otherwise
// quiet and this row leads; on a busy board a red gate or a dirty tree is the more urgent
// read and must not be displaced by a document's age.
func appendHandoff(rows []row, root string) []row {
	written, ok := handoffWrittenAt(root)
	if !ok {
		return rows
	}
	behind, ok := commitsSince(root, written)
	if !ok || behind == 0 {
		return rows
	}
	detail := fmt.Sprintf("written at %s, %s behind", short(written), plural(behind, "commit", "commits"))
	return append(rows, row{11, "handoff", detail, "rewrite " + handoffFile + " at HEAD"})
}

// handoffWrittenAt returns the commit that last wrote the handoff, and whether the age is
// worth checking at all. Three states report nothing rather than a distance:
//
//   - an in-flight edit — the handoff is modified, staged, or untracked in the work tree,
//     which means a session is writing it right now and its age is about to be reset;
//   - no handoff in history, which covers both a repo that keeps none (absence is a
//     choice, never a defect) and one whose handoff has never been committed;
//   - a failed git query, tolerated the way the other advisory housekeeping rows tolerate
//     one, because a broken query is not evidence of a stale document.
func handoffWrittenAt(root string) (string, bool) {
	dirty, err := git.Output("-C", root, "status", "--porcelain", "--", handoffFile)
	if err != nil || dirty != "" {
		return "", false
	}
	written, err := git.Output("-C", root, "log", "-1", "--format=%H", "HEAD", "--", handoffFile)
	if err != nil || written == "" {
		return "", false
	}
	return written, true
}

// commitsSince counts the commits between a commit and HEAD. The subject always comes from
// `git log HEAD`, so it is an ancestor of HEAD by construction and the count is meaningful.
func commitsSince(root, commit string) (int, bool) {
	out, err := git.Output("-C", root, "rev-list", "--count", commit+"..HEAD")
	if err != nil {
		return 0, false
	}
	n, err := strconv.Atoi(out)
	if err != nil {
		return 0, false
	}
	return n, true
}
