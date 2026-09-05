// Package canonicalpath owns the one canonical-path derivation for Bench.
package canonicalpath

import "path/filepath"

// Resolve returns one cleaned absolute spelling of path. A symbolic link resolves to its
// target, so two spellings of one directory compare equal. A path that does not exist yet
// carries no link to follow, so it keeps its absolute spelling. The only refusal is the
// working-directory failure that leaves a relative path with no absolute form.
//
// filepath.EvalSymlinks runs on path first, before any lexical cleaning. A ".." that
// follows a symbolic link component must pop off the link's resolved target, not off the
// unresolved literal text, or a root shaped "<base>/jump/.." — jump a symlink to
// "<base>/physical/child" — resolves to "<base>" instead of the physical "<base>/physical"
// the OS actually reaches. filepath.Abs then only makes the already-resolved path absolute;
// it sees no ".." left to mis-clean. When path does not exist, EvalSymlinks fails, and the
// fallback keeps today's absolute, lexically cleaned spelling of path itself, because
// nothing under an absent path resolves.
func Resolve(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		abs, absErr := filepath.Abs(path)
		if absErr != nil {
			return "", absErr
		}
		return filepath.Clean(abs), nil
	}
	abs, err := filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}
