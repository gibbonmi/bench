package landing

import (
	"fmt"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/spec"
)

// completeSpec applies the staged spec's implemented status and its verified closure to
// the composed tree.
func completeSpec(r ReviewedRequest, tree, source string) (string, error) {
	implemented, err := spec.Implemented(r.SpecBytes)
	if err != nil {
		return "", err
	}
	if tree, err = replaceTreeFile(r.Root, tree, r.SpecPath, implemented, r.SpecMode); err != nil {
		return "", fmt.Errorf("transition staged spec: %w", err)
	}
	if tree, err = closeDelivery(r.Root, tree, r.SpecPath, source); err != nil {
		return "", fmt.Errorf("close verified delivery: %w", err)
	}
	return tree, nil
}

// closeDelivery applies the verified closure of the reviewed spec to tree. The commitment
// owner derives the exact edits, and the completion oracle grades the same derivation, so
// the published tree carries the delivery and its closure together or carries neither.
func closeDelivery(root, tree, path, source string) (string, error) {
	edits, err := commitrepo.Store{Root: root}.Closure(tree, commitrepo.Delivery{Spec: path, Source: source})
	if err != nil || len(edits) == 0 {
		return tree, err
	}
	return editTree(root, tree, func(idx string) error {
		for _, edit := range edits {
			if edit.Delete {
				if err := indexRun(root, idx, "update-index", "--force-remove", "--", edit.Path); err != nil {
					return err
				}
				continue
			}
			if err := writeIndexFile(root, idx, edit.Path, edit.Data, edit.Mode); err != nil {
				return err
			}
		}
		return nil
	})
}
