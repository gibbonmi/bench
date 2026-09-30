// Package canonicalpath owns the one canonical-path derivation for Bench.
package canonicalpath

import "path/filepath"

// Operand answers what a path operand on a Bench command line names: the cleaned path the
// verb acts on, and the repo-relative spelling with forward slashes that every row and
// refusal prints, which is what an agent would type. An absolute operand stays, and a
// relative one joins onto cwd, so the same operand from two working directories names one
// file. A root that cannot relativize the path leaves nothing root-relative to print, so
// the display keeps the cleaned operand itself.
//
// The caller reads its own working directory and passes it in, because this package stays
// pure: it answers from its arguments alone and touches no ambient state.
func Operand(root, cwd, arg string) (path, display string) {
	path = arg
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	path = filepath.Clean(path)
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return path, filepath.ToSlash(filepath.Clean(arg))
	}
	return path, filepath.ToSlash(relative)
}

// Resolve returns one cleaned absolute spelling of path. A symbolic link resolves to its
// target, so two spellings of one directory compare equal. A path that does not exist yet
// carries no link to follow, so it keeps its absolute spelling. The only refusal is the
// working-directory failure that leaves a relative path with no absolute form.
//
// A relative path joins the working directory before any link resolves. A process that
// entered its directory through a link reads that link spelling back as its working
// directory, so a link step that ran first would leave the link spelling in the answer.
// The join is plain concatenation, not filepath.Join, because a lexical clean would pop a
// ".." off the unresolved text.
//
// filepath.EvalSymlinks then runs before any lexical cleaning. A ".." that follows a
// symbolic link component must pop off the link's resolved target, not off the unresolved
// literal text, or a root shaped "<base>/jump/.." — jump a symlink to
// "<base>/physical/child" — resolves to "<base>" instead of the physical "<base>/physical"
// the OS actually reaches. When path does not exist, EvalSymlinks fails, and the fallback
// keeps the absolute, lexically cleaned spelling of path itself, because nothing under an
// absent path resolves.
func Resolve(path string) (string, error) {
	if !filepath.IsAbs(path) {
		cwd, err := filepath.Abs(".")
		if err != nil {
			return "", err
		}
		path = cwd + string(filepath.Separator) + path
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return filepath.Clean(path), nil
	}
	return filepath.Clean(resolved), nil
}
