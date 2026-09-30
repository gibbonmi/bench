package worktree

import (
	"errors"
	"io"

	"github.com/gibbonmi/bench/internal/intent"
)

// ErrTreeTargetPath is the outcome of a tree-target value that names no label and that the
// worktree target grammar reads as a path. A tree target is never a path, so the caller
// answers this outcome as a grammar refusal.
var ErrTreeTargetPath = errors.New("tree target is a path")

// TreeTarget answers the worktree of the one assignment whose label is exactly label. The
// lookup takes no id, no prefix, and no path, so a tree target is always a name the
// operator chose.
//
// One active match wins over any inactive match of the same label. With no active match,
// one inactive match gives its state refusal. Two or more matches of one class give the
// ambiguity refusal. The label lookup runs before the path-shape test, so a label with a
// separator still names its worktree. The selected assignment then passes the state,
// missing-tree, and creation-bundle checks of every other target-taking verb.
func TreeTarget(root, label string) (string, error) {
	if !lineSafe(label) {
		return "", errTargetControls
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return "", err
	}
	labeled := matchingAssignments(assignments, func(a intent.Assignment) bool { return a.Label == label })
	if active := matchingAssignments(labeled, func(a intent.Assignment) bool { return landingActiveState(a.State) }); len(active) > 0 {
		labeled = active
	}
	switch {
	case len(labeled) == 0 && pathShaped(label):
		return "", ErrTreeTargetPath
	case len(labeled) == 0:
		return "", errTargetUnassigned
	case len(labeled) > 1:
		return "", ambiguousAssignments(labeled)
	}
	// The id is the one address that names exactly this record, so the shared resolver
	// applies its checks to the record that the label selected.
	selected, err := resolveAssignmentIn(root, labeled[0].ID, landingActiveState)
	if err != nil {
		return "", err
	}
	return selected.Worktree, nil
}

// pathShaped reports whether targetPath, the owner of the worktree target grammar, reads
// value as a path. A shape that the grammar refuses, such as `~user` or a relative path,
// is a path shape too.
func pathShaped(value string) bool {
	_, isPath, err := targetPath(value)
	return isPath || err != nil
}

// PrintTreeTargetRefusal prints err as the worktree target refusal of verb and answers
// exit 1.
func PrintTreeTargetRefusal(stderr io.Writer, verb string, err error) int {
	return printTargetRefusal(stderr, verb, err)
}
