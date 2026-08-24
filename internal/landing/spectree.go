package landing

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	benchgit "github.com/gibbonmi/bench/internal/git"
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
	return editTree(root, baseTree, func(idx string) error {
		blob, err := outputInput(root, content, "hash-object", "-w", "--stdin")
		if err != nil {
			return err
		}
		return indexRun(root, idx, "update-index", "--add", "--cacheinfo", gitRegularFileMode(mode)+","+blob+","+path)
	})
}

// removeTreeFolder writes baseTree without every entry beneath rel, through
// removeIndexTree's one spelling of the removal.
func removeTreeFolder(root, baseTree, rel string) (string, error) {
	return editTree(root, baseTree, func(idx string) error { return removeIndexTree(root, idx, rel) })
}

// editTree reads baseTree into a private index, applies edit to that index, and
// writes the resulting tree. No checkout or repository index is touched.
func editTree(root, baseTree string, edit func(idx string) error) (string, error) {
	dir, err := os.MkdirTemp("", "bench-reviewed-landing-index-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if err := indexRun(root, idx, "read-tree", baseTree); err != nil {
		return "", err
	}
	if err := edit(idx); err != nil {
		return "", err
	}
	return indexOutput(root, idx, "write-tree")
}
