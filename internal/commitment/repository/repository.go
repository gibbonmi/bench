// Package repository adapts commitment decisions to repository state.
package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// Store addresses one repository's tracked policy and local receipts.
type Store struct{ Root string }

// InventoryItem is one source available to an adoption proposal. Identity is the identity that
// plan binds, or a run's request. Scope is a row's detail owner, a spec's row, or a run's scope.
type InventoryItem struct {
	Kind, ID, State, Identity, Scope string
}

// Policy returns the current default-branch policy and whether it exists.
func (store Store) Policy() (commitment.Policy, bool, error) {
	policy, _, err := store.defaultPolicy()
	if err != nil {
		return commitment.Policy{}, false, err
	}
	if policy == nil {
		return commitment.Policy{}, false, nil
	}
	return *policy, true, nil
}

// Inventory reads each roadmap row at the default-branch tip and each staged spec at the checkout
// HEAD. Plan reads a row and a new deliverable there too, so each identity feeds plan unchanged.
func (store Store) Inventory() ([]InventoryItem, error) {
	revision, err := store.sourceRevision()
	if err != nil {
		return nil, err
	}
	document, err := roadmap.RevisionDocument(store.Root, revision)
	if err != nil {
		return nil, fmt.Errorf("commitment inventory refused: %w", err)
	}
	items := make([]InventoryItem, 0, len(document.Rows))
	for _, row := range document.Rows {
		identity, err := SourceIdentity(store.Root, revision, roadmap.RowFilePath(row.ID))
		if err != nil {
			return nil, fmt.Errorf("commitment inventory refused: roadmap row %s: %w", row.ID, err)
		}
		items = append(items, InventoryItem{Kind: "roadmap", ID: row.ID, State: row.SpecStatus, Identity: identity, Scope: roadmap.RowFilePath(row.ID)})
	}
	head, err := store.checkoutHead()
	if err != nil {
		return nil, err
	}
	specs, err := store.stagedSpecs(head)
	if err != nil {
		return nil, err
	}
	items = append(items, specs...)
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return nil, err
	}
	for _, assignment := range ledger.Assignments {
		items = append(items, InventoryItem{Kind: "run", ID: assignment.ID, State: string(assignment.State), Identity: assignment.Request, Scope: runScope(ledger.Commitment, assignment.ID)})
	}
	return items, nil
}

// runScope names the deliverable that binds the assignment, or else its continuation scope.
func runScope(state *intent.CommitmentState, assignment string) string {
	if state == nil {
		return ""
	}
	for _, binding := range state.Bindings {
		if binding.Assignment == assignment {
			return binding.Deliverable
		}
	}
	for _, continuation := range state.Continuations {
		if continuation.Assignment == assignment {
			return strings.Join(continuation.Scope, " ")
		}
	}
	return ""
}

// Plan validates and records one proposed policy.
func (store Store) Plan(input []byte) (commitment.Plan, error) {
	proposal, err := commitment.ParseProposal(input)
	if err != nil {
		return commitment.Plan{}, err
	}
	proposed := proposal.Policy
	current, _, err := store.defaultPolicy()
	if err != nil {
		return commitment.Plan{}, err
	}
	head, err := store.checkoutHead()
	if err != nil {
		return commitment.Plan{}, err
	}
	if err := store.validateDeliverables(current, proposed, head); err != nil {
		return commitment.Plan{}, err
	}
	plan, err := commitment.BuildPlan(current, proposal)
	if err != nil {
		return commitment.Plan{}, err
	}
	if err := store.validateSources(current, proposed, head, plan.Sources); err != nil {
		return commitment.Plan{}, err
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return commitment.Plan{}, fmt.Errorf("encode commitment plan: %w", err)
	}
	err = intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		if err := verifiedCompletions(ledger, current, plan.Predecessor, proposed); err != nil {
			return ledger, false, err
		}
		if err := store.listedRuns(ledger, current, plan.Continuations); err != nil {
			return ledger, false, err
		}
		for _, receipt := range ledger.CommitmentReceipts {
			if receipt.Plan != plan.ID {
				continue
			}
			if receipt.Payload != string(payload) {
				return ledger, false, errors.New("commitment plan identity collision")
			}
			return ledger, false, nil
		}
		ledger.CommitmentReceipts = append(ledger.CommitmentReceipts, intent.CommitmentReceipt{Plan: plan.ID, Payload: string(payload)})
		return ledger, true, nil
	}, nil)
	return plan, err
}

// Approve stages the exact policy transition named by planID.
func (store Store) Approve(planID, decision string, delayed, removed []string) (bool, error) {
	return store.approve(durablefile.Native(), planID, decision, delayed, removed)
}

// approve runs the policy, board, and ledger replacements through one replacer, so a test
// can fail one destination under the real stage and ledger transaction.
func (store Store) approve(replacer durablefile.Replacer, planID, decision string, delayed, removed []string) (bool, error) {
	if strings.TrimSpace(decision) == "" || !sanitize.LineSafe(decision) {
		return false, errors.New("commitment approval refused: decision must be one nonempty line")
	}
	if hasDuplicates(delayed) || hasDuplicates(removed) {
		return false, errors.New("commitment approval refused: effect operands contain duplicates")
	}
	var changed bool
	var restore func() error
	// The Compensation callback returns nothing, so its restore failure travels here.
	var restoreErr error
	err := intent.TransactWithReplacer(replacer, store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		for i := range ledger.CommitmentReceipts {
			receipt := &ledger.CommitmentReceipts[i]
			if receipt.Plan != planID {
				continue
			}
			var plan commitment.Plan
			if err := jsonfile.DecodeDocument([]byte(receipt.Payload), &plan); err != nil {
				return ledger, false, fmt.Errorf("commitment receipt %q: %w", planID, err)
			}
			for otherIndex := range ledger.CommitmentReceipts {
				other := ledger.CommitmentReceipts[otherIndex]
				if otherIndex == i || !other.Approved {
					continue
				}
				var approved commitment.Plan
				if err := jsonfile.DecodeDocument([]byte(other.Payload), &approved); err != nil {
					return ledger, false, fmt.Errorf("commitment receipt %q: %w", other.Plan, err)
				}
				if approved.Predecessor == plan.Predecessor && approved.ProposalIdentity != plan.ProposalIdentity {
					return ledger, false, replan(fmt.Errorf("commitment approval refused: predecessor changed from %s to %s", plan.Predecessor, approved.ProposalIdentity))
				}
			}
			current, predecessor, err := store.defaultPolicy()
			if err != nil {
				return ledger, false, err
			}
			if predecessor != plan.Predecessor && !(receipt.Approved && predecessor == plan.ProposalIdentity) {
				return ledger, false, replan(fmt.Errorf("commitment approval refused: predecessor changed from %s to %s", plan.Predecessor, predecessor))
			}
			if !sameSet(delayed, plan.Effects.Delayed) || !sameSet(removed, plan.Effects.Removed) {
				return ledger, false, errors.New("commitment approval refused: delayed or removed operands do not match the plan")
			}
			proposed, err := commitment.Parse(plan.Proposed)
			if err != nil {
				return ledger, false, err
			}
			head, err := store.checkoutHead()
			if err != nil {
				return ledger, false, err
			}
			if err := store.validateSources(current, proposed, head, plan.Sources); err != nil {
				return ledger, false, replan(err)
			}
			if receipt.Approved {
				if receipt.Decision != decision {
					return ledger, false, errors.New("commitment approval refused: replay decision differs")
				}
				return ledger, false, nil
			}
			if err := verifiedCompletions(ledger, current, plan.Predecessor, proposed); err != nil {
				return ledger, false, err
			}
			if err := store.listedRuns(ledger, current, plan.Continuations); err != nil {
				return ledger, false, replan(err)
			}
			restore, err = store.stage(replacer, plan)
			if err != nil {
				return ledger, false, err
			}
			receipt.Decision = decision
			receipt.Approved = true
			changed = true
			return withContinuations(ledger, plan.Continuations), true, nil
		}
		return ledger, false, fmt.Errorf("commitment approval refused: plan %q is unknown", planID)
	}, func() {
		if restore != nil {
			restoreErr = restore()
		}
	})
	return changed, errors.Join(err, restoreErr)
}

func (store Store) defaultPolicy() (*commitment.Policy, string, error) {
	revision, err := store.sourceRevision()
	if err != nil {
		return nil, "", err
	}
	return store.policyAt(revision)
}

func (store Store) policyAt(branch string) (*commitment.Policy, string, error) {
	const path = commitment.PolicyPath
	listing, err := git.Output("-C", store.Root, "ls-tree", "-z", branch, "--", path)
	if err != nil {
		return nil, "", fmt.Errorf("read commitment policy: %w", err)
	}
	if strings.TrimSuffix(listing, "\x00") == "" {
		return nil, "absent", nil
	}
	data, err := git.ReadTreeFile(store.Root, branch, path)
	if err != nil {
		return nil, "", fmt.Errorf("read commitment policy: %w", err)
	}
	policy, err := commitment.Parse(data)
	if err != nil {
		return nil, "", err
	}
	canonical, err := commitment.Bytes(policy)
	if err != nil {
		return nil, "", err
	}
	return &policy, commitment.Identity(canonical), nil
}

func sameSet(left, right []string) bool {
	left = append([]string(nil), left...)
	right = append([]string(nil), right...)
	slices.Sort(left)
	slices.Sort(right)
	return slices.Equal(left, right)
}

func hasDuplicates(values []string) bool {
	seen := map[string]bool{}
	for _, value := range values {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}
