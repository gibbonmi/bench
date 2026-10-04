package landing

import (
	"fmt"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/spec"
)

// Deliverable is the reviewed deliverable that the landing closes: the staged spec, the
// tickets-only folder, or nothing on a spec-less landing.
func (r ReviewedRequest) Deliverable() string {
	if r.SpecPath != "" {
		return r.SpecPath
	}
	return r.ClosePath
}

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

// completeTickets consumes the tickets-only folder from the composed tree by index removal
// and applies its verified closure. A folder the destination already removed lists no
// entries, so the removal writes the composed tree back unchanged. The closure reads the
// folder from the reviewed source, so the removal does not hide its evidence.
func completeTickets(r ReviewedRequest, tree, source string) (string, error) {
	tree, err := removeTreeFolder(r.Root, tree, r.ClosePath)
	if err != nil {
		return "", fmt.Errorf("close tickets-only folder: %w", err)
	}
	if tree, err = closeDelivery(r.Root, tree, r.ClosePath, source); err != nil {
		return "", fmt.Errorf("close verified delivery: %w", err)
	}
	return tree, nil
}

// closeDelivery applies the verified closure of the reviewed deliverable to tree. The
// commitment owner derives the exact edits, and publication admission and the spec
// completion oracle consume the same derivation, so the published tree carries the
// delivery and its closure together or carries neither.
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
