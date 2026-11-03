package landing

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/gittest"
)

// checkedLanding is a fixture checkout with a changed named path, the owner that lands it
// on a green authorization, and the request that names it. The owner records the tree it
// authorized.
func checkedLanding(t *testing.T) (root string, o Owner, r Request, authorized *string) {
	t.Helper()
	root = fixture(t)
	write(t, root, "named", "changed")
	base := gittest.Output(t, root, "rev-parse", "HEAD")
	authorized = new(string)
	o = New()
	o.authorize = func(_ context.Context, _ string, tree string, _ io.Writer, _ io.Writer) authorization.Result {
		*authorized = tree
		return authorization.Result{Kind: authorization.Green}
	}
	return root, o, Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: "named", Paths: []string{"named"}}, authorized
}

// PP10 and PP11: Land calls the check once, after the candidate commit exists and before any
// ref names it. The candidate holds the authorized composed tree over the expected base, and
// the commit that Land then publishes is that candidate.
func TestLandRunsTheCheckBeforePublication(t *testing.T) {
	root, o, r, authorized := checkedLanding(t)
	var seen []string
	r.Check = func(candidate string) error {
		seen = append(seen, candidate)
		if ref := gittest.Output(t, root, "rev-parse", r.Destination); ref != r.Expected {
			t.Errorf("destination inside the check = %s, want the expected base %s", ref, r.Expected)
		}
		if kind := gittest.Output(t, root, "cat-file", "-t", candidate); kind != "commit" {
			t.Errorf("candidate %s object type = %q, want commit", candidate, kind)
		}
		if tree := gittest.Output(t, root, "rev-parse", candidate+"^{tree}"); tree != *authorized {
			t.Errorf("candidate tree = %s, want the composed tree %s", tree, *authorized)
		}
		if parent := gittest.Output(t, root, "rev-parse", candidate+"^1"); parent != r.Expected {
			t.Errorf("candidate first parent = %s, want the expected base %s", parent, r.Expected)
		}
		return nil
	}
	got, err := o.Land(t.Context(), r)
	if err != nil || len(seen) != 1 || got.Commit != seen[0] || gittest.Output(t, root, "rev-parse", r.Destination) != got.Commit {
		t.Fatalf("Land = (%+v, %v) with check calls %v, want one call whose candidate the destination then names", got, err, seen)
	}
}

// PP12 and PP16: a check error refuses as CheckRefusal and leaves the destination and the
// checkout as they were, and DryRun never calls the check. A green check that moves the
// destination leaves the winner in place, because the update keeps its compare-and-swap.
func TestLandCheckRefusalPublishesNothing(t *testing.T) {
	t.Run("red check", func(t *testing.T) {
		root, o, r, _ := checkedLanding(t)
		before, err := CheckoutFingerprint(root)
		if err != nil {
			t.Fatal(err)
		}
		red := errors.New("chain step red")
		var seen []string
		r.Check = func(candidate string) error {
			seen = append(seen, candidate)
			return red
		}
		o.reconcile = func(Request, []string, composedSnapshot) error {
			t.Fatal("a refused check ran the reconcile")
			return nil
		}
		if err := o.DryRun(t.Context(), r); err != nil || len(seen) != 0 {
			t.Fatalf("DryRun = %v with check calls %v, want green and no call", err, seen)
		}
		_, err = o.Land(t.Context(), r)
		var refusal CheckRefusal
		if !errors.As(err, &refusal) || !errors.Is(err, red) || len(seen) != 1 || refusal.Candidate != seen[0] {
			t.Fatalf("Land err = %v with check calls %v, want CheckRefusal of the one candidate", err, seen)
		}
		after, fingerprintErr := CheckoutFingerprint(root)
		if ref := gittest.Output(t, root, "rev-parse", r.Destination); ref != r.Expected || fingerprintErr != nil || after != before {
			t.Fatalf("after the refusal the destination = %s and the fingerprint = (%s, %v), want %s and %s", ref, after, fingerprintErr, r.Expected, before)
		}
	})
	t.Run("check moves the destination", func(t *testing.T) {
		root, o, r, _ := checkedLanding(t)
		var winner string
		r.Check = func(string) error {
			winner = gittest.Output(t, root, "commit-tree", r.Expected+"^{tree}", "-p", r.Expected, "-m", "winner")
			gittest.MustRun(t, root, "update-ref", r.Destination, winner)
			return nil
		}
		_, err := o.Land(t.Context(), r)
		if !errors.As(err, new(TipMovedError)) {
			t.Fatalf("Land err = %v, want the tip-moved refusal", err)
		}
		if ref := gittest.Output(t, root, "rev-parse", r.Destination); ref != winner {
			t.Fatalf("destination = %s, want the winner %s", ref, winner)
		}
	})
}
