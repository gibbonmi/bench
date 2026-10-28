package status

import (
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/git"
)

// backgroundSignal names the background row of the board.
const backgroundSignal = "background"

// The background row sits on the ladder below the gate, Git, and worktree rows: a red ranks
// beside the guards, because the author acts on it, and each other outcome ranks last, because
// it only reports.
const (
	backgroundRedSeverity  = 3
	backgroundInfoSeverity = 13
)

// appendBackground adds the latest background outcome of the worktree root and its log path.
// The row is a signal and never a verdict: the gate row alone reads the gate cache, so a
// background outcome neither clears nor raises a gate row. A record that does not read is a
// red-rank row, so a broken record hides no red. A root with no git directory or no record
// adds no row.
func appendBackground(rows []row, root string) []row {
	gitdir, err := git.AdminDir(root)
	if err != nil {
		return rows
	}
	record, err := backgroundstate.ReadRecord(gitdir)
	if err != nil {
		return append(rows, row{backgroundRedSeverity, backgroundSignal, "record unreadable", advisoryAction("see " + backgroundstate.RecordPath(gitdir))})
	}
	if record.Outcome == "" {
		return rows
	}
	severity := backgroundInfoSeverity
	if record.Outcome == backgroundstate.Red {
		severity = backgroundRedSeverity
	}
	detail := string(record.Outcome) + " at " + Short(record.Tip) + ", not a verdict"
	return append(rows, row{severity, backgroundSignal, detail, advisoryAction("see " + record.Log)})
}
