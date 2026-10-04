package repository

import (
	"fmt"

	"github.com/gibbonmi/bench/internal/intent"
)

// Publication is the source assignment identity a landing froze before its gate, with
// the reviewed source commit that it publishes.
type Publication struct {
	Assignment, Request, Worktree, Source string
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
		return store.authorizeCandidate(ledger, owner, tree, publishedDelivery(ledger, owner, source.Source))
	}
	return fmt.Errorf("publication assignment %q is not active with its presented request and worktree; run bench worktree create", source.Assignment)
}

// publishedDelivery is the delivery that owner's current binding permits the publication
// to close. An unbound owner closes nothing.
func publishedDelivery(ledger intent.Ledger, owner intent.Assignment, source string) *Delivery {
	for _, binding := range runtimeState(ledger).Bindings {
		if binding.Assignment == owner.ID && binding.Request == owner.Request {
			return &Delivery{Spec: binding.Deliverable, Source: source}
		}
	}
	return nil
}
