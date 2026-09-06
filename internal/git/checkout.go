package git

import "path/filepath"

// IsPrimaryCheckout reports whether root owns the repository's common Git directory.
// Linked worktrees have a private Git directory beneath that common directory.
func IsPrimaryCheckout(root string) (bool, error) {
	gitDir, err := AdminDir(root)
	if err != nil {
		return false, err
	}
	common, err := CommonDir(root)
	if err != nil {
		return false, err
	}
	return filepath.Clean(gitDir) == filepath.Clean(common), nil
}
