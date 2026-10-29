// The reader green publication: which landing publishes on a green reader set.
package landing

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestOnlyTheReviewedLandingPublishesOnAReaderGreen is NG57. The reviewed landing onto main
// publishes on a green reader set and names the kind in its receipt. The worktree-commit
// landing answers to its own lane, so it refuses the kind and publishes nothing.
func TestOnlyTheReviewedLandingPublishesOnAReaderGreen(t *testing.T) {
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.ReaderGreen}
	}
	f := newReviewedLanding(t, reviewedBaseSpec, "", reviewedAmendedSpec, 0)
	result, err := o.LandReviewed(context.Background(), f.request(t, reviewedAmendedSpec, "reader landing"))
	if err != nil || result.Kind != authorization.ReaderGreen || gittest.Output(t, f.root, "rev-parse", "refs/heads/main") != result.Commit {
		t.Fatalf("reviewed landing = (%+v, %v), want a publication on %q", result, err, authorization.ReaderGreen)
	}

	root := fixture(t)
	base := gittest.Output(t, root, "rev-parse", "HEAD")
	write(t, root, "named", "changed")
	_, err = o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "named", Paths: []string{"named"}})
	var refused AuthorizationRefusal
	if !errors.As(err, &refused) || refused.Result.Kind != authorization.ReaderGreen {
		t.Fatalf("worktree-commit landing error = %v, want the refusal of %q", err, authorization.ReaderGreen)
	}
	if got := gittest.Output(t, root, "rev-parse", "refs/heads/main"); got != base {
		t.Fatalf("worktree-commit landing moved main to %s, want %s", got, base)
	}
}
