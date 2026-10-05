package landing

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing/published"
)

// stagedSpecMatches proves provenance, not agreement: the bytes the landing will
// transition must be the reviewed source tip's committed spec. The function never
// reads the destination's copy for comparison; the source's bytes win. A stale,
// amended, or absent destination spec is not a landing question.
func stagedSpecMatches(root, source, path string, want []byte) error {
	got, err := benchgit.Raw("-C", root, "show", source+":"+path)
	if err != nil || !bytes.Equal(got, want) {
		return errors.New("staged spec bytes are not the reviewed source tip's committed spec")
	}
	return nil
}

// specNeutralizedDestination returns a commit to compose against whose tree already
// carries the source's spec bytes. The spec path is then a one-sided change no merge
// can conflict on. Its parent is the real destination, so the merge base stays
// unchanged. The returned commit is a composition input only; it never becomes a
// published parent.
func specNeutralizedDestination(root, destination, path string, want []byte, mode os.FileMode) (string, error) {
	baseTree, err := output(root, "rev-parse", destination+"^{tree}")
	if err != nil {
		return "", fmt.Errorf("read destination tree: %w", err)
	}
	tree, err := replaceTreeFile(root, baseTree, path, want, mode)
	if err != nil {
		return "", fmt.Errorf("neutralize spec path: %w", err)
	}
	if tree == baseTree {
		return destination, nil
	}
	commit, err := output(root, "commit-tree", tree, "-p", destination, "-m", "compose against the reviewed source spec")
	if err != nil {
		return "", fmt.Errorf("neutralize spec path: %w", err)
	}
	return commit, nil
}

func replaceTreeFile(root, baseTree, path string, content []byte, mode os.FileMode) (string, error) {
	return published.Edit(root, baseTree, func(idx string) error {
		return published.WriteFile(root, idx, path, content, gitRegularFileMode(mode))
	})
}
