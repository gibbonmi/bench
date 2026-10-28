package commit

import (
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/git"
)

// backgroundOutcomePrefix opens the line that states the latest background outcome.
const backgroundOutcomePrefix = "background grade: latest "

// backgroundOutcomeLine states record: its outcome, the tip that it graded, and its log. The
// outcome is never a verdict, and the line says so.
func backgroundOutcomeLine(record backgroundstate.Record) string {
	return fmt.Sprintf("%s%s at %s, not a verdict (log %s)", backgroundOutcomePrefix, record.Outcome, record.Tip, record.Log)
}

// startBackgroundGrade starts the background grade of published, the commit that this call
// published in root. A call that published nothing starts none: a dry run and a refusal each
// leave published empty. Before the start, the call prints the latest background outcome of
// root, so the author sees the outcome of the last grade where the author already looks.
func startBackgroundGrade(root, published string, stderr io.Writer) {
	if published != "" {
		printBackgroundOutcome(root, stderr)
		gate.StartBackground(benchhome.Dir(), gate.KitValue(), root, published, stderr)
	}
}

// printBackgroundOutcome prints the latest background outcome of root to stderr. A root with
// no record prints nothing, and a record that does not read prints its error, so a broken
// record hides no red.
func printBackgroundOutcome(root string, stderr io.Writer) {
	gitdir, err := git.AdminDir(root)
	if err != nil {
		return
	}
	record, err := backgroundstate.ReadRecord(gitdir)
	if err != nil {
		fmt.Fprintf(stderr, "background grade: record unreadable: %v\n", err)
		return
	}
	if record.Outcome != "" {
		fmt.Fprintln(stderr, backgroundOutcomeLine(record))
	}
}
