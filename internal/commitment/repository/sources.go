package repository

import (
	"fmt"
	"path/filepath"
	"slices"
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

// TreeIdentity is the identity of a tickets-only folder whose Git tree is object.
func TreeIdentity(object string) string { return "git-tree:" + object }

func sourceAt(root, revision, path string) (string, []byte, error) {
	if ticketsOnlyAt(root, revision, path) {
		object, err := git.Output("-C", root, "rev-parse", "--verify", revision+":"+path)
		if err != nil {
			return "", nil, err
		}
		return TreeIdentity(strings.TrimSpace(object)), nil, nil
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

// stagedSpecs lists each staged spec that revision holds, in slug order. The status, the
// scope, and the identity of a row come from the same bytes, so plan binds the row that the
// inventory prints.
func (store Store) stagedSpecs(revision string) ([]InventoryItem, error) {
	held, err := git.Output("-C", store.Root, "ls-tree", "-r", "-z", "--name-only", revision, "--", spec.SpecsDir)
	if err != nil {
		return nil, err
	}
	paths := slices.DeleteFunc(strings.Split(held, "\x00"), func(path string) bool { return !spec.IsLiveSpecPath(path) })
	// Git lists specs/a-b/ before specs/a/, so the sort restores slug order.
	slices.SortFunc(paths, func(left, right string) int {
		return strings.Compare(spec.LiveSpecSlug(left), spec.LiveSpecSlug(right))
	})
	var items []InventoryItem
	for _, path := range paths {
		identity, data, err := sourceAt(store.Root, revision, path)
		if err != nil {
			return nil, fmt.Errorf("commitment inventory refused: deliverable %s: %w", path, err)
		}
		if status, scope := spec.Metadata(data); status == "staged" {
			items = append(items, InventoryItem{Kind: "deliverable", ID: path, State: status, Identity: identity, Scope: scope})
		}
	}
	return items, nil
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
	_, bindings := commitment.Unsettled(policy)
	for _, binding := range bindings {
		if err := store.validateDeliverable(binding.Source); err != nil {
			return err
		}
	}
	return nil
}

func (store Store) sourceRevision() (string, error) {
	_, revision, err := store.defaultRevision()
	return revision, err
}

// defaultRevision resolves the default branch and its current revision.
func (store Store) defaultRevision() (string, string, error) {
	branch, ok := git.ResolvedDefault(store.Root)
	if !ok {
		return "", "", errDefaultUnresolved
	}
	revision, err := git.Output("-C", store.Root, "rev-parse", branch)
	return branch, strings.TrimSpace(revision), err
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
