package spec

import (
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// Landing returns the landing commit of slug in the repository at root. The landing is the
// newest first-parent commit of the default branch that changes the spec to `Status:
// implemented`: its tree holds that status, and the tree of its first parent does not. A
// later commit that deletes the folder or lands another spec is not a transition, so the
// answer survives the retirement. found is false when no commit makes the transition.
func Landing(root, slug string) (commit string, found bool, err error) {
	branch, ok := git.ResolvedDefault(root)
	if !ok {
		return "", false, errors.New("default branch is unresolved")
	}
	path := LiveSpecPath(slug)
	raw, err := git.Output("-C", root, "log", "--first-parent", "--format=%H %P", branch, "--", ":(literal,top)"+path)
	if err != nil {
		return "", false, err
	}
	for _, line := range strings.Split(raw, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || !implementedAt(root, fields[0], path) {
			continue
		}
		if len(fields) == 1 || !implementedAt(root, fields[1], path) {
			return fields[0], true, nil
		}
	}
	return "", false, nil
}

// implementedAt reports whether the tree of commit holds path with the retirement marker
// that AwaitsRetirement reads. An absent path holds no status.
func implementedAt(root, commit, path string) bool {
	data, err := git.ReadTreeFile(root, commit, path)
	return err == nil && AwaitsRetirement(data)
}
