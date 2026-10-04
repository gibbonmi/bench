package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

// Verify checks the evidence document input for milestone against the default-branch
// policy at its current revision. A successful check records the receipt that a later
// completion proposal consumes. An identical replay records nothing new.
func (store Store) Verify(milestone string, input []byte) (commitment.Verification, error) {
	evidence, err := commitment.ParseEvidence(input)
	if err != nil {
		return commitment.Verification{}, err
	}
	revision, err := store.sourceRevision()
	if err != nil {
		return commitment.Verification{}, err
	}
	policy, identity, err := store.policyAt(revision)
	if err != nil {
		return commitment.Verification{}, err
	}
	if policy == nil {
		return commitment.Verification{}, errors.New("commitment policy is absent: adoption-required")
	}
	resolve := func(reference string) (string, error) {
		object, err := git.Output("-C", store.Root, "rev-parse", "--verify", "--quiet", revision+":"+reference)
		return strings.TrimSpace(object), err
	}
	verification, err := commitment.VerifyMilestone(*policy, identity, revision, milestone, evidence, resolve)
	if err != nil {
		return commitment.Verification{}, err
	}
	payload, err := json.Marshal(verification)
	if err != nil {
		return commitment.Verification{}, fmt.Errorf("encode milestone verification: %w", err)
	}
	err = intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		for _, receipt := range ledger.MilestoneReceipts {
			if receipt.ID != verification.ID {
				continue
			}
			if receipt.Payload != string(payload) {
				return ledger, false, errors.New("milestone verification identity collision")
			}
			return ledger, false, nil
		}
		ledger.MilestoneReceipts = append(ledger.MilestoneReceipts, intent.MilestoneReceipt{ID: verification.ID, Payload: string(payload)})
		return ledger, true, nil
	}, nil)
	return verification, err
}

// verifiedCompletions checks each milestone completion that proposed adds to current,
// whose identity is predecessor, against its recorded verification receipt.
func verifiedCompletions(ledger intent.Ledger, current *commitment.Policy, predecessor string, proposed commitment.Policy) error {
	added, err := commitment.Completions(current, proposed)
	if err != nil {
		return err
	}
	for _, completion := range added {
		index := slices.IndexFunc(ledger.MilestoneReceipts, func(receipt intent.MilestoneReceipt) bool { return receipt.ID == completion.Verification })
		if index < 0 {
			return fmt.Errorf("completion of milestone %q has no verification receipt %q; run bench commitment verify --milestone %s --evidence <file>", completion.Milestone, completion.Verification, completion.Milestone)
		}
		var verification commitment.Verification
		if err := jsonfile.DecodeDocument([]byte(ledger.MilestoneReceipts[index].Payload), &verification); err != nil {
			return fmt.Errorf("milestone receipt %q: %w", completion.Verification, err)
		}
		if err := verification.Completes(completion, predecessor, proposed); err != nil {
			return err
		}
	}
	return nil
}
