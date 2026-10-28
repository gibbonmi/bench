package treetarget

import (
	"errors"
	"io/fs"
	"os"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
)

// childExecutable answers the executable that the child runs in the tree at dir, or the
// reason that a kit worktree target refuses to start its build. A label target whose root
// declares Bench build inputs is a kit worktree target, and it runs a sealed build that
// matches the tree, so the registries compiled into the child match the tree. Its candidates
// are its own published build, the primary checkout's published build, and the running
// executable, in that order. The primary checkout never runs a worktree build. Every other
// target runs the invoking wrapper, which resolves the kit of the tree that it runs in. A
// direct run of the executable has no wrapper, and the child runs the same executable.
func childExecutable(call Call, value, dir string) (executable, refusal string) {
	if value != primaryTarget && freshness.DeclaresBuildInputs(dir) {
		candidates := []string{freshness.PublishedExecutable(dir)}
		if primary, err := git.PrimaryCheckout(call.Root); err == nil {
			candidates = append(candidates, freshness.PublishedExecutable(primary))
		}
		return worktreeBuild(dir, append(candidates, call.Running))
	}
	if call.Wrapper != "" {
		return call.Wrapper, ""
	}
	return call.Running, ""
}

// worktreeBuild answers the first of candidates that freshness.VerifyAt accepts against the
// sources of the kit worktree at root, in the order that childExecutable gives. When no
// candidate passes, an absent build of root is missing. Any other case, a changed build
// input or a changed executable, means that the build does not match the tree. The
// operator already trusts the primary build and the running executable, and the build of
// root controls its own seal, so the check detects a stale or changed build and not a
// forged one.
func worktreeBuild(root string, candidates []string) (executable, refusal string) {
	for _, candidate := range candidates {
		if freshness.VerifyAt(root, root, candidate) == nil {
			return candidate, ""
		}
	}
	if _, err := os.Lstat(freshness.PublishedExecutable(root)); errors.Is(err, fs.ErrNotExist) {
		return "", "worktree build is missing"
	}
	return "", "worktree build does not match the tree"
}
