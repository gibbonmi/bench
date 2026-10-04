package repository

import (
	"errors"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// errDefaultUnresolved refuses a policy read when no default branch resolves.
var errDefaultUnresolved = errors.New("commitment source: default branch is unresolved")

// Outlook projects the published policy and the local runtime state. A repository with
// no resolved default branch has published no commitment, so it reads as adoption-required.
func (store Store) Outlook() (commitment.Outlook, error) {
	policy, exists, err := store.Policy()
	if errors.Is(err, errDefaultUnresolved) {
		return commitment.Project(nil, intent.CommitmentState{}), nil
	}
	if err != nil {
		return commitment.Outlook{}, err
	}
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return commitment.Outlook{}, err
	}
	var published *commitment.Policy
	if exists {
		published = &policy
	}
	return commitment.Project(published, runtimeState(ledger)), nil
}
