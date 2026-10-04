package repository

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/roadmap"
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
	return store.authorizeCandidate(ledger, owner, tree, nil)
}

// authorizeCandidate grades tree for owner against one ledger snapshot. The caller supplies
// the owner, so a publication decides for its frozen source assignment rather than for
// whichever checkout holds the store root. Only a publication supplies a delivery, so only
// a publication can carry the verified closure of its own approved deliverable.
func (store Store) authorizeCandidate(ledger intent.Ledger, owner intent.Assignment, tree string, delivery *Delivery) error {
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
	transition := predecessor != identity
	if transition {
		closed, err := store.closedPolicy(revision, candidate, delivery)
		if err != nil {
			return err
		}
		if !closed {
			if err := store.approvedTransition(ledger, current, candidate, revision); err != nil {
				return err
			}
		}
	}
	protected := current
	if transition {
		protected = candidate
	}
	if err := store.protectedCandidate(protected, revision, tree, transition); err != nil {
		return fmt.Errorf("candidate changes protected commitment: %w; run bench commitment plan --input <file>", err)
	}
	changes, err := git.TreeChangesIncludingSubmodules(store.Root, revision, tree)
	if err != nil {
		return err
	}
	promotions, err := store.planningPromotions(tree, changes)
	if err != nil {
		return err
	}
	var production []string
	for _, change := range changes {
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
				return fmt.Errorf("legacy continuation scope excludes %q; run bench commitment plan --input <file>", path)
			}
		}
		return nil
	}
	return store.readyFor(ledger, owner, "", "")
}

func (store Store) approvedTransition(ledger intent.Ledger, current, candidate *commitment.Policy, revision string) error {
	if candidate == nil {
		return fmt.Errorf("commitment policy deletion is not approved; run bench commitment plan --input <file>")
	}
	expected, err := commitment.BuildPlan(current, *candidate)
	if err != nil {
		return err
	}
	for _, receipt := range ledger.CommitmentReceipts {
		if !receipt.Approved || receipt.Decision == "" || receipt.Plan != expected.ID {
			continue
		}
		var retained commitment.Plan
		if err := jsonfile.DecodeDocument([]byte(receipt.Payload), &retained); err != nil {
			return err
		}
		if retained.ID != expected.ID || retained.Predecessor != expected.Predecessor || !bytes.Equal(retained.Proposed, expected.Proposed) {
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
	satisfied := commitment.Satisfied(*policy)
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			for _, source := range outcome.Sources {
				if satisfied[source.ID] {
					continue
				}
				if err := store.validateSourceAt(tree, source); err != nil {
					return err
				}
				if roadmap.RowOwner(source.ID, source.Path) && !rows[source.ID] {
					return fmt.Errorf("protected roadmap owner %q was removed", source.ID)
				}
			}
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

// inScope reports whether path is one scope entry or lies below a scope directory.
func inScope(scope []string, path string) bool {
	for _, entry := range scope {
		if path == entry || strings.HasPrefix(path, entry+"/") {
			return true
		}
	}
	return false
}
