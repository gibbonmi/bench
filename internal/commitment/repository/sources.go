package repository

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/spec"
)

// SourceIdentity binds a regular file's bytes or a complete tickets-only tree.
func SourceIdentity(root, revision, path string) (string, error) {
	identity, _, err := sourceAt(root, revision, path)
	return identity, err
}

func sourceAt(root, revision, path string) (string, []byte, error) {
	if ticketsOnlyAt(root, revision, path) {
		object, err := git.Output("-C", root, "rev-parse", "--verify", revision+":"+path)
		if err != nil {
			return "", nil, err
		}
		return "git-tree:" + strings.TrimSpace(object), nil, nil
	}
	data, err := git.ReadTreeFile(root, revision, path)
	if err != nil {
		return "", nil, err
	}
	normalized, err := roadmap.RequirementBytes(path, data)
	if err != nil {
		return "", nil, err
	}
	return commitment.Identity(normalized), data, nil
}

func ticketsOnlyAt(root, revision, path string) bool {
	name := filepath.Base(path)
	return path == spec.ClosedFolderPath(name) && spec.TicketsOnly(spec.CommitTree(root, revision), name)
}

func (store Store) validateDeliverable(source commitment.SourceBinding) error {
	revision, err := store.sourceRevision()
	if err != nil {
		return err
	}
	isSpec := spec.IsLiveSpecPath(source.Path)
	if !isSpec && !ticketsOnlyAt(store.Root, revision, source.Path) {
		return fmt.Errorf("deliverable %q is not an approved spec or tickets-only folder", source.Path)
	}
	identity, data, err := sourceAt(store.Root, revision, source.Path)
	if err != nil {
		return fmt.Errorf("commitment source %q refused: %w", source.ID, err)
	}
	if isSpec {
		if _, err := spec.Implemented(data); err != nil {
			return fmt.Errorf("deliverable %q is not staged: %w", source.Path, err)
		}
	}
	return validateIdentity(source, identity)
}

func (store Store) validateDeliverables(policy commitment.Policy) error {
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			for _, binding := range outcome.Deliverables {
				if err := store.validateDeliverable(binding.Source); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (store Store) sourceRevision() (string, error) {
	branch, ok := git.ResolvedDefault(store.Root)
	if !ok {
		return "", fmt.Errorf("commitment source: default branch is unresolved")
	}
	revision, err := git.Output("-C", store.Root, "rev-parse", branch)
	return strings.TrimSpace(revision), err
}

func (store Store) validateSourceAt(revision string, source commitment.SourceBinding) error {
	got, err := SourceIdentity(store.Root, revision, source.Path)
	if err != nil {
		return fmt.Errorf("commitment source %q refused: %w", source.ID, err)
	}
	return validateIdentity(source, got)
}

func validateIdentity(source commitment.SourceBinding, got string) error {
	if got != source.Identity {
		return fmt.Errorf("commitment source %q changed from %s to %s", source.ID, source.Identity, got)
	}
	return nil
}
