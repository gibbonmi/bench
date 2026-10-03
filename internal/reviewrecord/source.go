package reviewrecord

import (
	"fmt"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// atSource reads the tree of commit, a full commit ID, and returns its source digest and
// the plan in it.
func atSource(root, spec, commit string) (string, Plan, error) {
	tree, err := sourceTree(root, commit)
	if err != nil {
		return "", Plan{}, err
	}
	digest, err := SourceDigest(root, tree, spec)
	if err != nil {
		return "", Plan{}, err
	}
	plan, err := ReadPlan(root, tree, spec)
	return digest, plan, err
}

func sourceTree(root, commit string) (string, error) {
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", "--quiet", commit+"^{tree}")
	if err != nil {
		return "", fmt.Errorf("unreadable tree of commit %s", commit)
	}
	return tree, nil
}
