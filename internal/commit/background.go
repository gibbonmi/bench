package commit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// worktreeBuild is the build that a publishing commit runs to refresh its checkout's own
// published executable. The tests replace it.
var worktreeBuild = runbinary.BuildSubject

// refreshWorktreeBuild makes the published executable of root match the tree after a commit
// published published, so the background grade that follows resolves a build of that tree.
// A root that declares no build inputs takes no build. An absent executable takes a build
// only when the published commit changes a build input, so a Markdown commit costs no
// build. Any other executable takes the refresh, which builds only when the seal does not
// verify, and its second seal read decides the outcome. The step writes only to stderr and
// never changes the commit's exit, because the commit has already published.
func refreshWorktreeBuild(root, published, home string, build func(context.Context, string, string) error, stderr io.Writer) {
	if !freshness.DeclaresBuildInputs(root) {
		return
	}
	if _, err := os.Lstat(freshness.PublishedExecutable(root)); errors.Is(err, fs.ErrNotExist) && !freshness.DiffChangesBuildInputs(root, published+"^", published) {
		return
	}
	rebuilt, buildErr, err := freshness.Refresh(context.Background(), root, build)
	switch {
	case rebuilt:
		fmt.Fprintln(stderr, "worktree build: rebuilt dist/bench")
	case err != nil:
		fmt.Fprintf(stderr, "worktree build: failed: %s; rebuild with %s\n", sanitize.Controls(rebuildCause(buildErr, err)), freshness.RebuildIn(home, root).Action)
	}
}

// rebuildCause names why a refresh failed. The build's own error names it best. Otherwise
// the seal refusal names it without its remedy, because the failure line prints the remedy
// once itself.
func rebuildCause(buildErr, err error) string {
	if buildErr != nil {
		return buildErr.Error()
	}
	var refusal *freshness.Refusal
	if errors.As(err, &refusal) {
		return refusal.Reason()
	}
	return err.Error()
}

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
