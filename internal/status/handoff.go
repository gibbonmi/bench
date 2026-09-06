package status

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// HandoffFile is the phase-close continuation artifact `AGENTS.md` names: the document a
// cold session reads to resume without conversation history. Its shape is documented
// inside the file itself so the template cannot drift from the artifact it describes. It
// is the name this package's callers read, and it derives from internal/handoffdoc, which
// owns the document. This way the staleness signal and the writer can never watch
// different files.
const HandoffFile = handoffdoc.DocumentPath

// appendHandoff adds the handoff-staleness signal (sev 12): the commits that landed since
// the handoff was last written. It turns "trust git over the handoff" from a rule the next
// session must recall into an ambient fact the SessionStart hook already prints.
//
// The handoff's age is read from git history rather than from a line inside the document.
// A self-reported date is the same remembered-not-computed defect this signal exists to
// close. A session that rewrites the body and forgets the date leaves the check
// confidently calling a stale document current. It cannot name its own commit anyway,
// since the handoff lands in the commit it describes.
//
// It ranks last deliberately. The value is at cold pickup, where the board is otherwise
// quiet and this row leads. On a busy board a red gate or a dirty tree is the more urgent
// read, and must not be displaced by a document's age.
func appendHandoff(rows []row, root string) []row {
	return append(rows, handoffRows(root, false)...)
}

// handoffRows builds the whole handoff family: the document's own age, then one entry per
// section whose assignment branch has moved past the tip that section records. Both the
// board and its --all expansion read this, so the summary and the expanded list can never
// state a different set.
//
// The document's clock and a section's clock answer different questions. The file is one
// artifact with one write time, so its age covers `main`, the section a primary checkout
// owns and the one no branch backs. A sibling worktree's section is dated against its own
// branch instead: a session that rewrites its own section resets the file's write time,
// and a file-level clock would then read every other section current on that write.
func handoffRows(root string, all bool) []row {
	return append(handoffAgeRows(root), handoffSectionRows(root, all)...)
}

// handoffAgeRows dates the document as a whole. A tracked handoff is dated by the commit
// that last wrote it; an ignored one by its write time.
func handoffAgeRows(root string) []row {
	if _, err := git.Output("-C", root, "check-ignore", "-q", HandoffFile); err == nil {
		return ignoredHandoffAgeRows(root)
	}
	written, ok := handoffWrittenAt(root)
	if !ok {
		return nil
	}
	behind, ok := commitsSince(root, written)
	if !ok || behind == 0 {
		return nil
	}
	detail := fmt.Sprintf("written at %s, %s behind", Short(written), Plural(behind, "commit", "commits"))
	return []row{{12, "handoff", detail, commandAction(handoffAction)}}
}

// ignoredHandoffAgeRows is the age signal for a git-ignored handoff: a local file with
// no commit of its own. The age comes from the file's own write time, the one computed
// fact an ignored file still carries, and the distance counts the commits whose commit
// date is after that write. A missing file reports nothing: keeping no handoff is a
// choice rather than a defect.
func ignoredHandoffAgeRows(root string) []row {
	info, err := os.Stat(filepath.Join(root, filepath.FromSlash(HandoffFile)))
	if err != nil {
		return nil
	}
	written := info.ModTime()
	// The bound starts one second after the write. Commit dates carry one-second
	// granularity, so the commit a handoff rewrite itself lands in must not count
	// against the document it carries.
	out, err := git.Output("-C", root, "rev-list", "--count", "--since="+written.Add(time.Second).Format(time.RFC3339), "HEAD")
	if err != nil {
		return nil
	}
	behind, err := strconv.Atoi(out)
	if err != nil || behind == 0 {
		return nil
	}
	detail := fmt.Sprintf("written at %s, %s behind", written.Format("2006-01-02 15:04"), Plural(behind, "commit", "commits"))
	return []row{{12, "handoff", detail, commandAction(handoffAction)}}
}

// handoffWrittenAt returns the commit that last wrote the handoff, and whether the age is
// worth checking at all. Three states report nothing rather than a distance:
//
//   - an in-flight edit: the handoff is modified, staged, or untracked in the work tree.
//     A session is writing it right now, and its age is about to reset;
//   - no handoff in history. This covers a repo that keeps none, a choice rather than a
//     defect, and one whose handoff has never been committed;
//   - a failed git query, tolerated the way other advisory housekeeping rows tolerate one.
//     A broken query is not evidence of a stale document.
func handoffWrittenAt(root string) (string, bool) {
	dirty, err := git.Output("-C", root, "status", "--porcelain", "--", HandoffFile)
	if err != nil || dirty != "" {
		return "", false
	}
	written, err := git.Output("-C", root, "log", "-1", "--format=%H", "HEAD", "--", HandoffFile)
	if err != nil || written == "" {
		return "", false
	}
	return written, true
}

// handoffSection is one request section joined to the assignment its key names: the
// commits that assignment's branch holds past the tip the section records.
type handoffSection struct {
	key    string
	behind int
}

// handoffSectionRows states the sections a reader must act on. With all false the list
// reduces to the section furthest behind; with all true, under `bench status --all`, every
// section states itself. That is the census row's summary-and-expansion shape, and it keeps
// the default board's five-row budget bounded however many worktrees are live.
func handoffSectionRows(root string, all bool) []row {
	sections := handoffSectionsBehind(root)
	if len(sections) == 0 {
		return nil
	}
	if all {
		out := make([]row, 0, len(sections))
		for _, section := range sections {
			out = append(out, handoffSectionRow(section))
		}
		return out
	}
	// The list is ordered by distance, so the first section is the one furthest behind.
	return []row{handoffSectionRow(sections[0])}
}

// handoffSectionRow renders one section's distance.
func handoffSectionRow(section handoffSection) row {
	key := Short(sanitize.Controls(section.key))
	detail := "request " + key + " " + Plural(section.behind, "commit", "commits") + " behind"
	return row{12, "handoff", detail, commandAction(handoffAction)}
}

// handoffSectionsBehind dates every request section against its own assignment branch. It
// keeps a section only when it states a distance, and orders the distances largest before
// smallest, so the caller's lead is the section furthest behind. An unreadable document
// reports nothing: the document's grammar is the leaf package's to refuse, and this row is
// advisory housekeeping.
func handoffSectionsBehind(root string) []handoffSection {
	doc, err := handoffdoc.Read(filepath.Join(root, filepath.FromSlash(HandoffFile)))
	if err != nil {
		return nil
	}
	branches := activeAssignmentBranches(root)
	var out []handoffSection
	for _, section := range doc.Sections {
		if section.Key == handoffdoc.MainKey {
			continue
		}
		behind, ok := sectionBehind(root, section, branches)
		if !ok || behind == 0 {
			continue
		}
		out = append(out, handoffSection{key: section.Key, behind: behind})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].behind != out[j].behind {
			return out[i].behind > out[j].behind
		}
		return out[i].key < out[j].key
	})
	return out
}

// sectionBehind counts the commits one assignment branch holds past the tip its section
// records, and reports whether the join produced a distance at all. A section the ledger
// has no active assignment for is residue a landing or a release removes, and a section
// whose tip does not resolve names a commit the repository has lost. Neither is a distance,
// and neither is a defect this advisory row can act on, so both report nothing.
func sectionBehind(root string, section handoffdoc.Section, branches map[string]string) (int, bool) {
	branch, found := branches[section.Key]
	if !found {
		return 0, false
	}
	tip := sectionField(section, handoffdoc.LabelWorktreeTip)
	if tip == "" {
		return 0, false
	}
	out, err := git.Output("-C", root, "rev-list", "--count", tip+".."+branch)
	if err != nil {
		return 0, false
	}
	behind, err := strconv.Atoi(out)
	if err != nil {
		return 0, false
	}
	return behind, true
}

// sectionField reads one label line's value.
func sectionField(section handoffdoc.Section, label string) string {
	for _, field := range section.Fields {
		if field.Label == label {
			return field.Value
		}
	}
	return ""
}

// activeAssignmentBranches maps each active assignment's request digest to its branch. The
// digest is the section key, so the map is the join between the document and the ledger. A
// ledger that cannot be read joins no section, so every section reports nothing rather than
// a distance the join never computed.
func activeAssignmentBranches(root string) map[string]string {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return nil
	}
	branches := make(map[string]string, len(assignments))
	for _, assignment := range assignments {
		if assignment.State != intent.StateActive {
			continue
		}
		branches[assignment.Request] = assignment.Branch
	}
	return branches
}

// expandHandoffSignals replaces the summarized handoff rows with the per-section list. It
// is the sibling of expandCensusSignals: it rebuilds the whole family from the same
// source, so the summary it drops and the rows it adds cannot disagree.
func expandHandoffSignals(root string, signals []Signal) []Signal {
	rows := handoffRows(root, true)
	if len(rows) == 0 {
		return signals
	}
	out := make([]Signal, 0, len(signals)+len(rows))
	for _, signal := range signals {
		if signal.Name != "handoff" {
			out = append(out, signal)
		}
	}
	for _, r := range rows {
		out = append(out, newSignal(r.sev, r.signal, r.detail, r.action))
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Severity < out[j].Severity })
	return out
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
