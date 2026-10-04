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
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return err
	}
	owners := intent.AssignmentsOwning(ledger.Assignments, store.Root)
	if len(owners) != 1 || owners[0].State != intent.StateActive {
		return fmt.Errorf("commit requires an active owned assignment; run bench worktree create")
	}
	approved := predecessor != identity
	if approved {
		if err := store.approvedTransition(ledger, current, candidate, revision); err != nil {
			return err
		}
	}
	protected := current
	if approved {
		protected = candidate
	}
	if err := store.protectedCandidate(protected, revision, tree, approved); err != nil {
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
	production := false
	for _, change := range changes {
		if !commitment.PlanningPath(change.Path, change.DstMode, promotions) || !commitment.PlanningPath(change.Path, change.SrcMode, promotions) {
			production = true
			break
		}
	}
	if production {
		return store.ready(ledger, "")
	}
	return nil
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
	sequence := before.SequenceText
	if transition {
		projected, err := roadmap.ProjectSequence([]byte(before.Text), commitment.Selection(*policy).Outcomes)
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
	for _, milestone := range policy.Milestones {
		for _, outcome := range milestone.Outcomes {
			for _, source := range outcome.Sources {
				if err := store.validateSourceAt(tree, source); err != nil {
					return err
				}
				if source.Path == roadmap.RowPath(source.ID) && !rows[source.ID] {
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
		if change.DstMode != "100644" || !commitment.PlanningPath(change.Path, change.DstMode, nil) || !strings.HasSuffix(change.Path, ".md") {
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
		if change.Path != "CONTEXT.md" && !strings.HasPrefix(change.Path, "docs/adr/") {
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
