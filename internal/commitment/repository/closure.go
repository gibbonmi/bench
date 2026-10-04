package repository

import (
	"bytes"
	"fmt"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/roadmap"
)

// Delivery names one verified publication: the approved spec path or tickets-only folder,
// and the reviewed source commit whose completion evidence verifies it.
type Delivery struct{ Spec, Source string }

// Edit is one exact file effect of a verified closure. A deleted path carries no data
// and no mode. A written path is a planning document and carries the planning mode.
type Edit struct {
	Path, Mode string
	Data       []byte
	Delete     bool
}

// Delivered returns policy after delivery and the sources that the delivery completely
// satisfies. The completion evidence of a spec is the record that the reviewed source
// retains. A tickets-only folder has no record: its evidence is the folder's tree in the
// reviewed source, which holds the approved ticket acceptance that the gate graded.
func (store Store) Delivered(policy commitment.Policy, delivery Delivery) (commitment.Policy, []commitment.SourceBinding, error) {
	evidence, err := store.completionEvidence(delivery)
	if err != nil {
		return policy, nil, err
	}
	return commitment.Deliver(policy, delivery.Spec, delivery.Source, evidence)
}

func (store Store) completionEvidence(delivery Delivery) (string, error) {
	if ticketsOnlyAt(store.Root, delivery.Source, delivery.Spec) {
		return SourceIdentity(store.Root, delivery.Source, delivery.Spec)
	}
	record, err := reviewrecord.RecordPath(delivery.Spec)
	if err != nil {
		return "", err
	}
	evidence, err := git.Output("-C", store.Root, "rev-parse", "--verify", delivery.Source+":"+record)
	if err != nil {
		return "", fmt.Errorf("read completion evidence %s: %w", record, err)
	}
	return evidence, nil
}

// Closure derives the exact closure that delivery makes to tree. It records the delivery
// in the policy, removes each satisfied row and its detail owner, and projects the
// recommended sequence from the outcomes that remain. A tree whose policy records no new
// fact for the deliverable needs no edit. A rowless delivery closes no row, so a tree
// with no board needs the policy edit alone.
func (store Store) Closure(tree string, delivery Delivery) ([]Edit, error) {
	current, _, err := store.policyAt(tree)
	if err != nil || current == nil {
		return nil, err
	}
	next, closed, err := store.Delivered(*current, delivery)
	if err != nil || len(next.Deliveries) == len(current.Deliveries) {
		return nil, err
	}
	policy, err := commitment.Bytes(next)
	if err != nil {
		return nil, err
	}
	edits := []Edit{{Path: commitment.PolicyPath, Mode: commitment.PlanningMode, Data: policy}}
	var rows []string
	for _, source := range closed {
		if roadmap.RowOwner(source.ID, source.Path) {
			rows = append(rows, source.ID)
			edits = append(edits, Edit{Path: source.Path, Delete: true})
		}
	}
	index, present, err := roadmap.RevisionIndex(store.Root, tree)
	if err != nil {
		return nil, err
	}
	if !present && len(rows) == 0 {
		return edits, nil
	}
	closedIndex, err := roadmap.Close(index, rows, commitment.Remaining(next))
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(index, closedIndex) {
		edits = append(edits, Edit{Path: roadmap.RoadmapFile, Mode: commitment.PlanningMode, Data: closedIndex})
	}
	return edits, nil
}

// ReconcileDelivered releases each local claim and binding whose outcome the published
// default-branch policy records as delivered, and each legacy continuation whose scope
// that policy records as delivered. The published fact is the authority, so a failure
// here leaves the delivered obligation closed, and a retry changes nothing more.
func (store Store) ReconcileDelivered() error {
	policy, exists, err := store.Policy()
	if err != nil || !exists {
		return err
	}
	remaining := map[string]bool{}
	for _, outcome := range commitment.Remaining(policy) {
		remaining[outcome] = true
	}
	open := func(milestone, outcome string) bool { return milestone != policy.ActiveMilestone || remaining[outcome] }
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		state := runtimeState(ledger)
		next := state
		next.Claims = nil
		for _, claim := range state.Claims {
			if open(claim.Milestone, claim.Outcome) {
				next.Claims = append(next.Claims, claim)
			}
		}
		next.Bindings = nil
		for _, binding := range state.Bindings {
			if open(binding.Milestone, binding.Outcome) {
				next.Bindings = append(next.Bindings, binding)
			}
		}
		next.Continuations = nil
		for _, continuation := range state.Continuations {
			if !commitment.ScopeDelivered(policy, continuation.Scope) {
				next.Continuations = append(next.Continuations, continuation)
			}
		}
		return withRuntime(ledger, next)
	}, nil)
}
