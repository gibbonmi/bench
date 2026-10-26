// The marker step of a first-run landing.
package worktree

import (
	"context"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/landing"
)

// advanceLandedMarker advances the green marker of branch from prior to the published
// commit of result. A publication on a green reader set leaves the marker at prior: the
// reader set graded the readers of the changed families alone, so the marker stays at the
// last whole green.
func advanceLandedMarker(root, branch string, result landing.ReviewedResult, prior string) error {
	if result.Kind == authorization.ReaderGreen {
		return nil
	}
	return authorization.AdvanceMarker(context.Background(), root, branch, result.Commit, prior)
}
