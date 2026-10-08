package repository

import (
	"fmt"
	"slices"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// Publication is the source assignment identity a landing froze before its gate, with
// the reviewed source commit that it publishes and the reviewed deliverable that it
// closes, if any: a staged spec or a tickets-only folder.
type Publication struct {
	Assignment, Request, Worktree, Source, Deliverable string
}

// AdmitPublication grades the exact composed tree for the frozen source assignment.
func (store Store) AdmitPublication(source Publication, tree string) error {
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return err
	}
	return store.admitPublication(ledger, source, tree)
}

// PublishAdmitted decides tree again under the intent lock and runs publish before that
// lock releases. No competing claim, blocker, or receipt can change the decision between
// that final check and the ref update that publish performs.
func (store Store) PublishAdmitted(source Publication, tree string, publish func() error) error {
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		if err := store.admitPublication(ledger, source, tree); err != nil {
			return ledger, false, err
		}
		return ledger, false, publish()
	}, nil)
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
