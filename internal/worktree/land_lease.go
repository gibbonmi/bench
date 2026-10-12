// The first run's landing lease: one lease on the destination ref from composition
// through publish, so a second landing refuses before it composes.
package worktree

import (
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

// leaseLanding takes the destination's lease for the assignment after the admission
// proofs and before the composition. A refusal carries the caller's own re-run. The
// publish decision or the composition's refusal releases the lease, and a resume takes none.
func leaseLanding(root, destination string, assignment intent.Assignment, request, base, tip, specArg, path string) (*commitrepo.Lease, error) {
	lease, err := commitrepo.Store{Root: root}.AcquireLease(destination, assignment.ID)
	if err != nil {
		return nil, landingSourceRoute(err, request, base, tip, specArg, path, assignment.ID)
	}
	return lease, nil
}
