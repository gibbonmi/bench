package worktree

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/gate"
)

// stopGradesOf stops the background grades of each worktree at dir or under it in the machine
// queue under the Bench home home, before a verb removes that worktree or lands it as a stack
// member. A stop that fails prints one line, because a background grade never blocks the verb.
// A verb with no Bench home has no machine queue.
func stopGradesOf(home, dir string, stderr io.Writer) {
	if home == "" {
		return
	}
	if err := gate.StopWorktreeGrades(home, dir); err != nil {
		fmt.Fprintf(stderr, "background grade: not stopped: %v\n", err)
	}
}

// stopReclaimGrades stops the background grades under each pool key that the reclaim plan
// removes, before the apply removes the first key.
func stopReclaimGrades(plan poolReclaimPlan, home string, stderr io.Writer) {
	for _, verdict := range plan.verdicts {
		if verdict.Reclaimable() {
			stopGradesOf(home, filepath.Join(poolKeysDirAt(home), verdict.Key), stderr)
		}
	}
}
