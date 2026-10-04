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

// Delivery names one verified spec publication: the approved spec path and the reviewed
// source commit whose completion record verifies it.
type Delivery struct{ Spec, Source string }

// Edit is one exact file effect of a verified closure. A deleted path carries no data
// and no mode.
type Edit struct {
	Path, Mode string
	Data       []byte
	Delete     bool
}

// closureMode is the one mode a closure writes. The policy and the board index are
// planning documents, and that class admits only a regular non-executable file.
const closureMode = "100644"

// Delivered returns policy after delivery and the sources that the delivery completely
// satisfies. The completion evidence is the record that the reviewed source retains.
func (store Store) Delivered(policy commitment.Policy, delivery Delivery) (commitment.Policy, []commitment.SourceBinding, error) {
	record, err := reviewrecord.RecordPath(delivery.Spec)
	if err != nil {
		return policy, nil, err
	}
	evidence, err := git.Output("-C", store.Root, "rev-parse", "--verify", delivery.Source+":"+record)
	if err != nil {
		return policy, nil, fmt.Errorf("read completion evidence %s: %w", record, err)
	}
	return commitment.Deliver(policy, delivery.Spec, delivery.Source, evidence)
}

// Closure derives the exact closure that delivery makes to tree. It records the delivery
// in the policy, removes each satisfied row and its detail owner, and projects the
// recommended sequence from the outcomes that remain. A tree whose policy approves no
// obligation for the spec needs no edit.
func (store Store) Closure(tree string, delivery Delivery) ([]Edit, error) {
	current, _, err := store.policyAt(tree)
	if err != nil || current == nil {
		return nil, err
	}
	next, closed, err := store.Delivered(*current, delivery)
	if err != nil || len(closed) == 0 {
		return nil, err
	}
	policy, err := commitment.Bytes(next)
	if err != nil {
		return nil, err
	}
	edits := []Edit{{Path: commitment.PolicyPath, Mode: closureMode, Data: policy}}
	var rows []string
	for _, source := range closed {
		if roadmap.RowOwner(source.ID, source.Path) {
			rows = append(rows, source.ID)
			edits = append(edits, Edit{Path: source.Path, Delete: true})
		}
	}
	index, err := git.ReadTreeFile(store.Root, tree, roadmap.RoadmapFile)
	if err != nil {
		return nil, err
	}
	closedIndex, err := roadmap.Close(index, rows, commitment.Remaining(next))
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(index, closedIndex) {
		edits = append(edits, Edit{Path: roadmap.RoadmapFile, Mode: closureMode, Data: closedIndex})
	}
	return edits, nil
}

// ReconcileDelivered releases each local claim and binding whose outcome the published
// default-branch policy records as delivered. The published fact is the authority, so a
// failure here leaves the delivered obligation closed, and a retry changes nothing more.
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
		return withRuntime(ledger, next)
	}, nil)
}
