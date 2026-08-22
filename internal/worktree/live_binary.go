package worktree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/freshness"
)

// resolveRunningBinary answers which executable this process is.
// bin/bench.sh resolves exactly one binary and execs it, so the running process is that
// resolution's answer. Reading it here reuses the wrapper's own choice instead of
// restating the search order, and it stays correct for the inherited-BENCH_RUN_BINARY
// branch too.
var resolveRunningBinary = os.Executable

// liveBinaryWarnings is where the residue guard announces a removal that would take the
// CLI down with it. The guard runs below every command's writer. The warning must land
// while the removal is proposed, not in a plan the caller reads afterward.
var liveBinaryWarnings io.Writer = os.Stderr

// warnBeforeRemovingLiveBinary tells the caller that removing candidate out of root
// disables the binary answering bench, and names the command that rebuilds it. It stays
// silent for any other candidate.
func warnBeforeRemovingLiveBinary(root, candidate string) {
	if !isRunningBinary(candidate) {
		return
	}
	fmt.Fprintf(liveBinaryWarnings,
		"bench: %s is the binary currently answering bench; removing it disables the CLI, the git guard it backs, and the gate's BENCH_RUN_BINARY at once. Rebuild it with: %s\n",
		candidate, freshness.RebuildAction(root))
}

// isRunningBinary reports whether candidate is the very file this process is running.
//
// The predicate is identity, not path. A dist/bench in an unrelated checkout is ordinary
// residue. A path-shaped predicate would warn on every checkout ever built. That trains
// the warning away before the one removal that matters arrives.
//
// An executable this process cannot resolve is unknown, never absent. Nothing has been
// shown about the candidate, so the guard warns rather than remove in silence. A candidate
// that cannot be stat'd is the other way round — there is no live file there to lose.
func isRunningBinary(candidate string) bool {
	running, err := resolveRunningBinary()
	if err != nil {
		return true
	}
	if resolved, err := filepath.EvalSymlinks(running); err == nil {
		running = resolved
	}
	runningInfo, err := os.Stat(running)
	if err != nil {
		return true
	}
	candidateInfo, err := os.Stat(candidate)
	if err != nil {
		return false
	}
	return os.SameFile(runningInfo, candidateInfo)
}
