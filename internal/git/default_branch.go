package git

import "strings"

// mainCandidate is the branch name the code probes when origin/HEAD is unset. It is a
// candidate, and not an answer, until the object database confirms it. A repository with
// no branch named main falls through to the sole-local-branch fallback and then to
// ok=false. The code discards the candidate the moment the object database does not
// confirm it.
const mainCandidate = "main"

// ResolvedDefault is the single owner of the default-branch fact. It returns a branch
// only when the branch resolves to a commit. The candidates, in order, are origin/HEAD's
// short name with the `origin/` prefix stripped, the mainCandidate probe, and a lone
// local branch. The local branch is the fallback for a master-only repository. When
// nothing resolves, it returns ("", false), so no caller can put a branch this
// repository does not have into a ref or a message.
func ResolvedDefault(root string) (string, bool) {
	def := mainCandidate
	if out, err := Output("-C", root, "symbolic-ref", "--quiet", "--short", "refs/remotes/origin/HEAD"); err == nil && out != "" {
		def = strings.TrimPrefix(out, "origin/")
	}
	if OK("-C", root, "rev-parse", "--verify", "--quiet", def+"^{commit}") {
		return def, true
	}
	branches, err := LocalBranches(root)
	if err == nil && len(branches) == 1 && OK("-C", root, "rev-parse", "--verify", "--quiet", branches[0]+"^{commit}") {
		return branches[0], true
	}
	return "", false
}

// RemoteDefaultRef returns `origin/<default branch>`, or "" when the default does not
// resolve. Callers already fall through on the empty value; the function never fabricates
// a ref.
func RemoteDefaultRef(root string) string {
	if def, ok := ResolvedDefault(root); ok {
		return "origin/" + def
	}
	return ""
}
