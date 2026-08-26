// Package census owns the record of raw shell calls that touch a Bench worktree.
package census

import (
	"path/filepath"

	"github.com/gibbonmi/bench/internal/poolkey"
)

// Dir returns the directory that holds one repository's census records. The
// directory is a sibling of the worktree pool and never a child of it, because
// `bench worktree reclaim` enumerates `<home>/worktrees` and a foreign entry
// there changes its plan.
func Dir(home, root string) string {
	return filepath.Join(home, "census", poolkey.Key(root))
}
