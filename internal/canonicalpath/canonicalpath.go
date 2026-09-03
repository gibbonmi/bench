// Package canonicalpath owns the one canonical-path derivation for Bench.
package canonicalpath

import "path/filepath"

// Resolve returns one cleaned absolute spelling of path. A symbolic link resolves to its
// target, so two spellings of one directory compare equal. A path that does not exist yet
// carries no link to follow, so it keeps its absolute spelling. The only refusal is the
// working-directory failure that leaves a relative path with no absolute form.
func Resolve(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return filepath.Clean(abs), nil
	}
	return filepath.Clean(resolved), nil
}
