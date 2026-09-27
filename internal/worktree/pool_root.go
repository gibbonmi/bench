package worktree

import (
	"errors"
	"fmt"
	"os"
)

// poolRootMode is the one mode a Bench-selected pool root may carry: the checkouts,
// leases, and markers under it are the operator's private working state.
const poolRootMode os.FileMode = 0o700

// The pool-root refusal details, one per refusing branch, so a test can tell which
// branch refused.
const (
	poolRootSymlink      = "worktree pool root is a symlink"
	poolRootNotDirectory = "worktree pool root is not a directory"
	poolRootModeUnset    = "worktree pool root mode cannot be set"
)

// preparePoolRoot creates the pool root when it is absent and sets it to poolRootMode on
// every call, so a root a crash left permissive is tightened on the next entry and an
// already tight root passes unchanged. A root that is a symlink or not a directory is
// refused without being followed. A mode Bench cannot set is a refusal, so the operation
// that needs the root stops before it puts a checkout in it.
func preparePoolRoot(j joins, pool string) error {
	info, err := os.Lstat(pool)
	if errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(pool, poolRootMode); err == nil {
			info, err = os.Lstat(pool)
		}
	}
	if err != nil {
		return err
	}
	wanted := fmt.Sprintf("directory at mode %04o", poolRootMode)
	if info.Mode()&os.ModeSymlink != 0 {
		return refusalError{refusal{detail: poolRootSymlink, observed: pool, wanted: wanted}}
	}
	if !info.IsDir() {
		return refusalError{refusal{detail: poolRootNotDirectory, observed: pool, wanted: wanted}}
	}
	if err := j.chmodPool(pool, poolRootMode); err != nil {
		return refusalError{refusal{detail: poolRootModeUnset + ": " + err.Error(), observed: pool, wanted: wanted}}
	}
	return nil
}
