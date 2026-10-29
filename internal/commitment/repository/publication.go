package repository

import (
	"fmt"
	"slices"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// Publication is the source assignment identity a landing froze before its gate, with
// the reviewed source commit that it publishes and the reviewed deliverable that it
// closes, if any: a staged spec or a tickets-only folder. Destination is the ref that the
// landing publishes, and Lease is the landing's lease on it, if it holds one.
type Publication struct {
	Assignment, Request, Worktree, Source, Deliverable, Destination string
	Lease                                                           *Lease
}

// AdmitPublication grades the exact composed tree for the frozen source assignment.
func (store Store) AdmitPublication(source Publication, tree string) error {
	ledger, err := store.ledger().Read()
	if err != nil {
		return err
	}
	return store.admitPublication(ledger, source, tree)
}

// PublishAdmitted decides the landing of source over tree as a stack of that one member.
func (store Store) PublishAdmitted(source Publication, tree string, publish func() error) error {
	return store.PublishStackAdmitted([]Publication{source}, []string{tree}, publish)
}

// PublishStackAdmitted decides the build window and the tree of each member again under the
// intent lock, and runs publish once before that lock releases. The tree of members[i] is
// trees[i]. No competing claim, blocker, receipt, or gap change can alter the decision between
// that final check and the ref update that publish performs. The first refused member stops
// the decision before publish with its own refusal. The same decision consumes the override
// of each member, so an override admits one publication, and an assignment that the stack
// does not publish keeps its override. It also drops the lease that the first member holds,
// so the next landing can take it once this lock releases.
func (store Store) PublishStackAdmitted(members []Publication, trees []string, publish func() error) error {
	if len(members) == 0 || len(trees) != len(members) {
		return fmt.Errorf("stack publication has %d members and %d trees; want one tree for each member of a non-empty stack", len(members), len(trees))
	}
	published := false
	err := store.ledger().Transact(intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		for i, member := range members {
			if err := store.window(ledger, member); err != nil {
				return ledger, false, err
			}
			if err := store.admitPublication(ledger, member, trees[i]); err != nil {
				return ledger, false, err
			}
		}
		if err := publish(); err != nil {
			return ledger, false, err
		}
		published = true
		state := storedState(ledger)
		for _, member := range members {
			state.Overrides = withoutOverride(state.Overrides, member.Assignment)
		}
		consumed, overridden, _ := withRuntime(ledger, state)
		next, released, err := members[0].Lease.drop(consumed)
		return next, overridden || released, err
	}, nil)
	// The ref update is the commit point. A write that fails after it leaves the lease
	// record with no lock holder, and the next acquire takes that record over. It also
	// leaves each override record, which admits one more landing of its assignment.
	if published {
		return nil
	}
	return err
}

// HoldsDeliveryBinding reports whether the ledger holds a delivery binding for assignment,
// whatever deliverable it binds. It reads the bindings and not the open builds, because a
// started tickets-only deliverable binds its assignment but opens no build.
func (store Store) HoldsDeliveryBinding(assignment string) (bool, error) {
	ledger, err := store.ledger().Read()
	if err != nil {
		return false, err
	}
	return slices.ContainsFunc(storedState(ledger).Bindings, func(binding intent.DeliveryBinding) bool {
		return binding.Assignment == assignment
	}), nil
}

func (store Store) admitPublication(ledger intent.Ledger, source Publication, tree string) error {
	for _, owner := range ledger.Assignments {
		if owner.ID != source.Assignment {
			continue
		}
		if owner.State != intent.StateActive || owner.Request != source.Request || len(intent.AssignmentsOwning([]intent.Assignment{owner}, source.Worktree)) != 1 {
			break
		}
		var delivery *Delivery
		if source.Deliverable != "" {
			delivery = &Delivery{Spec: source.Deliverable, Source: source.Source}
		}
		return store.authorizeCandidate(ledger, owner, tree, delivery, true)
	}
	return refusalroute.Raised{Name: refusalroute.CommitmentNeedsAssignment, Err: fmt.Errorf("publication assignment %q is not active with its presented request and worktree", source.Assignment)}
}

// closureAuthority refuses the verified closure of the deliverable at path unless owner
// may close it. A listed legacy run closes only a deliverable that its scope lists, and
// the policy decides whether that deliverable is approved. Any other owner closes only
// its current bound deliverable, so an unbound run receives the start guidance.
func (store Store) closureAuthority(ledger intent.Ledger, owner intent.Assignment, path string) error {
	if scope, listed := continuationScope(ledger, owner); listed {
		if !slices.Contains(scope, path) {
			return scopeRefusal(path)
		}
		return nil
	}
	return store.readyFor(ledger, owner, path, "")
}

// scopeRefusal is the refusal of a path that a listed legacy scope does not authorize.
func scopeRefusal(path string) error {
	return refusalroute.Raised{Name: refusalroute.CommitmentDecision, Err: fmt.Errorf("legacy continuation scope excludes %q", path)}
}
