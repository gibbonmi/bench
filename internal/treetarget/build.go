package treetarget

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/freshness"
)

// childExecutable answers the executable that the child runs in the tree at dir, or the
// reason that a kit worktree target refuses to start its build. A label target whose root
// declares Bench build inputs is a kit worktree target, and it runs its own worktree build,
// so the registries compiled into the child match the tree. The primary checkout never runs
// a worktree build. Every other target runs the invoking wrapper, which resolves the kit of
// the tree that it runs in. A direct run of the executable has no wrapper, and the child
// runs the same executable.
func childExecutable(call Call, value, dir string) (executable, refusal string) {
	if value != primaryTarget && freshness.DeclaresBuildInputs(dir) {
		return worktreeBuild(dir)
	}
	if call.Wrapper != "" {
		return call.Wrapper, ""
	}
	return call.Running, ""
}

// worktreeBuild answers the published build of the kit worktree at root, after
// freshness.Verify accepts it against the sources of root. A build that is absent is
// missing. Any other refusal of the seal, a changed build input or a changed executable,
// means that the build does not match the tree. The candidate controls the seal, so the
// check detects a stale or changed build and not a forged one.
func worktreeBuild(root string) (executable, refusal string) {
	executable = freshness.PublishedExecutable(root)
	if freshness.Verify(root, executable) == nil {
		return executable, ""
	}
	if _, err := os.Lstat(executable); errors.Is(err, fs.ErrNotExist) {
		return "", "worktree build is missing"
	}
	return "", "worktree build does not match the tree"
}

// printBuildRefusal prints the refusal of a kit worktree target whose build cannot start,
// and answers exit 1. The repair names the label as one shell word, quoted only when it
// needs quoting. Neither line names the executable path.
func printBuildRefusal(stderr io.Writer, command, label, refusal string) int {
	fmt.Fprintln(stderr, command+": "+refusal)
	fmt.Fprintln(stderr, "next=bench worktree build "+axi.ShellQuote(label))
	return 1
}
