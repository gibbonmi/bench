package repository

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
)

// AuthorizeCandidate grades the complete proposed tree against current default-branch authority.
func (store Store) AuthorizeCandidate(tree string) error {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return err
	}
	owner, owned := activeOwner(ledger, store.Root)
	if !owned {
		return fmt.Errorf("commit requires an active owned assignment; run bench worktree create")
	}
	return store.authorizeCandidate(ledger, owner, tree, nil, false)
}

// authorizeCandidate grades tree for owner against one ledger snapshot. The caller supplies
// the owner, so a publication decides for its frozen source assignment rather than for
// whichever checkout holds the store root. Only a publication supplies a delivery, so only
// a publication can carry the verified closure of its reviewed deliverable, and only when
// owner has the authority to close that deliverable. A publication with no delivery is a
// spec-less landing.
func (store Store) authorizeCandidate(ledger intent.Ledger, owner intent.Assignment, tree string, delivery *Delivery, publication bool) error {
	revision, err := store.sourceRevision()
	if err != nil {
		return err
	}
	current, predecessor, err := store.policyAt(revision)
	if err != nil {
		return err
	}
	candidate, identity, err := store.policyAt(tree)
	if err != nil {
		return err
	}
	reference := revision
	if !publication {
		reference = store.branchReference(revision, identity)
	}
	// A branch that keeps its merge base's policy proposes no transition, whatever main
	// approved since then.
	transition := reference == revision && predecessor != identity
	if transition {
		closed, err := store.closedPolicy(revision, candidate, delivery)
		if err != nil {
			return err
		}
		if closed {
			err = store.closureAuthority(ledger, owner, delivery.Spec)
		} else {
			err = store.approvedTransition(ledger, current, candidate, revision)
		}
		if err != nil {
			return err
		}
	}
	// A branch graded from its merge base keeps that base's protection, which is the policy
	// it holds. Main's later protection applies when the landing grades the composed tree.
	protected := current
	if transition || reference != revision {
		protected = candidate
	}
	if err := store.protectedCandidate(protected, reference, tree, transition); err != nil {
		return fmt.Errorf("candidate changes protected commitment: %w; run bench commitment plan --input <file>", err)
	}
	changes, err := git.TreeChangesIncludingSubmodules(store.Root, reference, tree)
	if err != nil {
		return err
	}
	promotions, err := store.planningPromotions(tree, changes)
	if err != nil {
		return err
	}
	var production []string
	for _, change := range changes {
		if store.retiredRecord(tree, change) {
			continue
		}
		if !commitment.PlanningPath(change.Path, change.DstMode, promotions) || !commitment.PlanningPath(change.Path, change.SrcMode, promotions) {
			production = append(production, change.Path)
		}
	}
	if len(production) == 0 {
		return nil
	}
	// A listed legacy run finishes only its approved scope; it never falls back to a binding.
	if scope, listed := continuationScope(ledger, owner); listed {
		for _, path := range production {
			if !inScope(scope, path) {
				return scopeRefusal(path)
			}
		}
		return nil
	}
	// Each guard above already decided, so a light-path ticket waives only the binding.
	err = store.readyFor(ledger, owner, "", "")
	if !errors.Is(err, errUnbound) {
		return err
	}
	if publication {
		return store.lightPathPublication(tree, current, production, delivery, err)
	}
	return store.lightPath(tree, current, production, err)
}

// branchReference names the commit that a commit candidate's own change is read from. The
// candidate composes onto the checkout's HEAD, so the branch's change starts at the merge
// base of HEAD and the default-branch tip, and each commit main made since then is main's.
// A branch that changes the policy proposes a transition from the tip, so it grades against
// the tip. A merge base that does not resolve, or a base policy that does not read, also
// leaves the tip as the reference, which is the stricter grading.
func (store Store) branchReference(revision, identity string) string {
	out, err := git.Output("-C", store.Root, "merge-base", "HEAD", revision)
	base := strings.TrimSpace(out)
	if err != nil || base == "" || base == revision {
		return revision
	}
	if _, inherited, err := store.policyAt(base); err != nil || inherited != identity {
		return revision
	}
	return base
}

// retiredRecord reports whether change deletes the review record of a spec whose folder
// tree no longer holds. That deletion is the planning half of a spec retirement. Any other
// review record change stays a production change. The record path owner decides which
// path is a review record, so a path that only resembles one stays production.
func (store Store) retiredRecord(tree string, change git.TreeChange) bool {
	if change.Status != "D" || change.SrcMode != commitment.PlanningMode {
		return false
	}
	slug := spec.SlugOf(change.Path)
	record, err := reviewrecord.RecordPath(spec.LiveSpecPath(slug))
	return err == nil && record == change.Path && !spec.CommitTree(store.Root, tree).FolderIsDirectory(slug)
}

func (store Store) approvedTransition(ledger intent.Ledger, current, candidate *commitment.Policy, revision string) error {
	if candidate == nil {
		return fmt.Errorf("commitment policy deletion is not approved; run bench commitment plan --input <file>")
	}
	for _, receipt := range ledger.CommitmentReceipts {
		if !receipt.Approved || receipt.Decision == "" {
			continue
		}
		var retained commitment.Plan
		if err := jsonfile.DecodeDocument([]byte(receipt.Payload), &retained); err != nil {
			return err
		}
		// The plan identity binds the runs that the approval listed, and the receipt retains them.
		expected, err := commitment.BuildPlan(current, commitment.Proposal{Policy: *candidate, Continuations: retained.Continuations})
		if err != nil {
			return err
		}
		if receipt.Plan != expected.ID || retained.ID != expected.ID || retained.Predecessor != expected.Predecessor || !bytes.Equal(retained.Proposed, expected.Proposed) {
			continue
		}
		for _, source := range expected.Sources {
			if err := store.validateSourceAt(revision, source); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("candidate policy has no exact approval; run bench commitment plan --input <file>")
}

// closedPolicy reports whether candidate is exactly the policy edit of the verified closure
// that delivery makes to revision.
func (store Store) closedPolicy(revision string, candidate *commitment.Policy, delivery *Delivery) (bool, error) {
	if candidate == nil || delivery == nil {
		return false, nil
	}
	edits, err := store.Closure(revision, *delivery)
	if err != nil {
		return false, err
	}
	for _, edit := range edits {
		if edit.Path == commitment.PolicyPath {
			got, err := commitment.Bytes(*candidate)
			return err == nil && bytes.Equal(edit.Data, got), err
		}
	}
	return false, nil
}

func (store Store) protectedCandidate(policy *commitment.Policy, revision, tree string, transition bool) error {
	if policy == nil {
		return nil
	}
	before, err := roadmap.RevisionDocument(store.Root, revision)
	if err != nil {
		return err
	}
	after, err := roadmap.RevisionDocument(store.Root, tree)
	if err != nil {
		return err
	}
	// A revision with no board has no sequence to project, so its candidate keeps none.
	_, board, err := roadmap.RevisionIndex(store.Root, revision)
	if err != nil {
		return err
	}
	sequence := before.SequenceText
	if transition && board {
		projected, err := roadmap.ProjectSequence([]byte(before.Text), commitment.Remaining(*policy))
		if err != nil {
			return err
		}
		// Parse the sequence through its owner without interpreting rendered commands here.
		sequence = roadmap.SequenceText(projected)
	}
	if after.SequenceText != sequence {
		return fmt.Errorf("protected recommended sequence changed; run bench commitment plan --input <file>")
	}
	rows := map[string]bool{}
	for _, row := range after.Rows {
		rows[row.ID] = true
	}
	sources, _ := commitment.Unsettled(*policy)
	for _, source := range sources {
		if err := store.validateSourceAt(tree, source); err != nil {
			return err
		}
		if roadmap.RowOwner(source.ID, source.Path) && !rows[source.ID] {
			return fmt.Errorf("protected roadmap owner %q was removed", source.ID)
		}
	}
	return nil
}

func (store Store) planningPromotions(tree string, changes []git.TreeChange) ([]string, error) {
	var artifacts [][]byte
	var promotions []string
	for _, change := range changes {
		if change.DstMode != commitment.PlanningMode || !commitment.PlanningPath(change.Path, change.DstMode, nil) || !strings.HasSuffix(change.Path, ".md") {
			continue
		}
		if !strings.HasPrefix(change.Path, "specs/") && !strings.HasPrefix(change.Path, "decisions/") {
			continue
		}
		data, err := git.ReadTreeFile(store.Root, tree, change.Path)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, data)
	}
	for _, change := range changes {
		if !commitment.PlanningPromotionPath(change.Path) {
			continue
		}
		for _, artifact := range artifacts {
			if bytes.Contains(artifact, []byte("`"+change.Path+"`")) || bytes.Contains(artifact, []byte("]("+change.Path+")")) {
				promotions = append(promotions, change.Path)
				break
			}
		}
	}
	return promotions, nil
}

// inScope reports whether a scope entry covers path.
func inScope(scope []string, path string) bool {
	for _, entry := range scope {
		if tickets.Covers(entry, path) {
			return true
		}
	}
	return false
}
