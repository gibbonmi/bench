package git

import "strings"

// mainCandidate is the branch name probed when origin/HEAD is unset. It is a candidate to
// verify against the object database, never an answer: a repository with no branch named
// main falls through to the sole-local-branch fallback and then to ok=false, and the
// candidate is discarded the moment the object database does not confirm it.
const mainCandidate = "main"

// ResolvedDefault is the single owner of the default-branch fact. It returns the branch
// only when one resolves to a commit: origin/HEAD's short name with the `origin/` prefix
// stripped, else the mainCandidate probe, else a lone local branch — the fallback that
// makes a master-only repository resolve. When nothing resolves it returns ("", false),
// so no caller can put a branch this repository does not have into a ref or a message.
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

// RemoteDefaultRef is `origin/<default branch>`, or "" when the default does not resolve —
// the empty candidate its callers already fall through on rather than a fabricated ref.
func RemoteDefaultRef(root string) string {
	if def, ok := ResolvedDefault(root); ok {
		return "origin/" + def
	}
	return ""
}
