package repository

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/spec"
)

// SourceIdentity binds a regular file's bytes or a complete tickets-only tree.
func SourceIdentity(root, revision, path string) (string, error) {
	identity, _, err := sourceAt(nil, root, revision, path)
	return identity, err
}

// TreeIdentity is the identity of a tickets-only folder whose Git tree is object.
func TreeIdentity(object string) string { return "git-tree:" + object }

// sourceAt is the identity and the bytes of the source at path in revision. It reads a
// regular file through run.
func sourceAt(run *git.Run, root, revision, path string) (string, []byte, error) {
	if ticketsOnlyAt(root, revision, path) {
		object, err := git.Output("-C", root, "rev-parse", "--verify", revision+":"+path)
		if err != nil {
			return "", nil, err
		}
		return TreeIdentity(strings.TrimSpace(object)), nil, nil
	}
	data, err := run.ReadTreeFile(root, revision, path)
	if err != nil {
		return "", nil, err
	}
	identity, err := requirementIdentity(path, data)
	if err != nil {
		return "", nil, err
	}
	return identity, data, nil
}

// requirementIdentity is the identity of the source bytes data at path. It drops the
// occurrence ledger of a row file, so the tree read and the working read compute one
// identity for the same bytes.
func requirementIdentity(path string, data []byte) (string, error) {
	normalized, err := roadmap.RequirementBytes(path, data)
	if err != nil {
		return "", err
	}
	return commitment.Identity(normalized), nil
}

// workingIdentity is the identity of the working file at path under root. A linked component
// below root refuses, and only a parsed regular file passes, so a hostile checkout cannot
// redirect or block the read. A link above root passes, because a checkout can sit below a
// linked directory.
func workingIdentity(root, path string) (string, error) {
	if err := bounds.RefuseLinksBelow(root, path); err != nil {
		return "", err
	}
	read := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(path)))
	if read.State != bounds.StateParsed {
		if read.Reason != "" {
			return "", fmt.Errorf("working file %s is %s: %s", path, read.State, read.Reason)
		}
		return "", fmt.Errorf("working file %s is %s", path, read.State)
	}
	return requirementIdentity(path, read.Data)
}

// sourceRead names where one source reads: revision, or the working file under the store
// root when working is true.
type sourceRead struct {
	revision string
	working  bool
}

// identityAt is the identity of the source at path where at names.
func (store Store) identityAt(at sourceRead, path string) (string, error) {
	if at.working {
		return workingIdentity(store.Root, path)
	}
	identity, _, err := sourceAt(store.Run, store.Root, at.revision, path)
	return identity, err
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
		identity, data, err := sourceAt(store.Run, store.Root, revision, path)
		if err != nil {
			return nil, fmt.Errorf("commitment inventory refused: deliverable %s: %w", path, err)
		}
		if status, scope := spec.Metadata(data); status == "staged" {
			items = append(items, InventoryItem{Kind: "deliverable", ID: path, State: status, Identity: identity, Scope: scope})
		}
	}
	return items, nil
}

// DefaultDeclaration reads the `Persisted fields:` declaration of the bound deliverable at
// path from the default-branch tip, never from a checkout. A reviewed spec-stage landing
// wrote that copy, so a candidate tree cannot declare its own fields. A spec without the
// line, and a tickets-only folder, declare none. A malformed line refuses with
// spec.ErrPersistedFieldsInvalid.
func (store Store) DefaultDeclaration(path string) ([]spec.PersistedField, error) {
	revision, err := store.sourceRevision()
	if err != nil {
		return nil, err
	}
	_, data, err := sourceAt(store.Run, store.Root, revision, path)
	if err != nil {
		return nil, fmt.Errorf("default-branch copy of %s: %w", path, err)
	}
	return spec.PersistedFields(data)
}

func ticketsOnlyAt(root, revision, path string) bool {
	name := filepath.Base(path)
	return path == spec.ClosedFolderPath(name) && spec.TicketsOnly(spec.CommitTree(root, revision), name)
}

// readTree names where source reads in a transition from current to proposed. A new
// deliverable source reads tree, which publishes it. A new outcome source whose path revision
// does not hold reads tree too, so one landing publishes a row and its admission. A changed
// outcome source, whose id and path current binds at another identity, reads the working file
// under the store root when checkout is true, at plan and approve, so the row edit can stay
// uncommitted until the approval. At commit and landing it reads tree, so the committed bytes
// must equal the approved identity. Every other source reads revision, so a worktree edit of a
// row whose identity the proposal keeps cannot pass as current.
func (store Store) readTree(current *commitment.Policy, proposed commitment.Policy, source commitment.SourceBinding, tree, revision string, checkout bool) sourceRead {
	if commitment.NewDeliverable(current, proposed, source) {
		return sourceRead{revision: tree}
	}
	if commitment.ChangedSource(current, proposed, source) {
		return sourceRead{revision: tree, working: checkout}
	}
	if commitment.NewSource(current, proposed, source) && store.absentAt(revision, source.Path) {
		return sourceRead{revision: tree}
	}
	return sourceRead{revision: revision}
}

// absentAt reports that revision holds no entry at path. A read error is no proof of absence,
// so the source then reads revision and refuses there.
func (store Store) absentAt(revision, path string) bool {
	listing, err := store.Run.ListTree(store.Root, revision, "--", path)
	return err == nil && len(listing) == 0
}

// checkoutHead resolves the commit that HEAD names in the store root.
func (store Store) checkoutHead() (string, error) {
	head, err := git.Output("-C", store.Root, "rev-parse", "HEAD")
	return strings.TrimSpace(head), err
}

// validateDeliverable validates source at the default-branch tip. A valid source is a staged
// spec or a tickets-only folder, and staged reports which one.
func (store Store) validateDeliverable(source commitment.SourceBinding) (staged bool, err error) {
	revision, err := store.sourceRevision()
	if err != nil {
		return false, err
	}
	return store.validateDeliverableAt(revision, source)
}

func (store Store) validateDeliverableAt(revision string, source commitment.SourceBinding) (bool, error) {
	isSpec := spec.IsLiveSpecPath(source.Path)
	if !isSpec && !ticketsOnlyAt(store.Root, revision, source.Path) {
		return false, fmt.Errorf("deliverable %q is not an approved spec or tickets-only folder", source.Path)
	}
	identity, data, err := sourceAt(store.Run, store.Root, revision, source.Path)
	if err != nil {
		return false, fmt.Errorf("commitment source %q refused: %w", source.ID, err)
	}
	if isSpec {
		if _, err := spec.Implemented(data); err != nil {
			return false, fmt.Errorf("deliverable %q is not staged: %w", source.Path, err)
		}
	}
	return isSpec, validateIdentity(source, identity)
}

func (store Store) validateDeliverables(current *commitment.Policy, proposed commitment.Policy, head string) error {
	revision, err := store.sourceRevision()
	if err != nil {
		return err
	}
	_, bindings := commitment.Unsettled(proposed)
	for _, binding := range bindings {
		// A deliverable reads a revision, never the working file.
		if _, err := store.validateDeliverableAt(store.readTree(current, proposed, binding.Source, head, revision, false).revision, binding.Source); err != nil {
			return err
		}
	}
	return nil
}

func (store Store) validateSources(current *commitment.Policy, proposed commitment.Policy, head string, sources []commitment.SourceBinding) error {
	revision, err := store.sourceRevision()
	if err != nil {
		return err
	}
	for _, source := range sources {
		if err := store.validateSourceAt(store.readTree(current, proposed, source, head, revision, true), source); err != nil {
			return err
		}
	}
	return nil
}

func (store Store) sourceRevision() (string, error) {
	_, revision, err := store.defaultRevision()
	return revision, err
}

// defaultRevision resolves the default branch and its current revision. git.DefaultTip owns
// both, and an unresolved branch maps to errDefaultUnresolved.
func (store Store) defaultRevision() (string, string, error) {
	branch, revision, err := git.DefaultTip(store.Root)
	if errors.Is(err, git.ErrDefaultUnresolved) {
		return "", "", errDefaultUnresolved
	}
	return branch, revision, err
}

func (store Store) validateSourceAt(at sourceRead, source commitment.SourceBinding) error {
	got, err := store.identityAt(at, source.Path)
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
