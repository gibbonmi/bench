// Shared fixtures for the commitment delivery landing tests: the delivery routes, the
// seeded closure fixture, and the delivery state that main and the intent record hold.
package worktree

import (
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// closureLandingFixture is the graded public landing fixture whose policy approves
// closureSpec as the complete delivery of FT1. The delivery outcome also owns each
// residual row, and outcome B owns FT2. A non-nil step adds one prospective gate line.
func closureLandingFixture(t *testing.T, request string, step func(*testrepo.GateFixture, string) string, residual ...string) landingFixture {
	t.Helper()
	return seededClosureFixture(t, request, step, func(t testing.TB, root string) {
		commitmenttest.SeedClosure(t, root, closureSpec, residual...)
	})
}

// seededClosureFixture is the graded public landing fixture whose policy seed writes. Each
// fence path joins the spec's ownership fence, so the source may change that path.
func seededClosureFixture(t *testing.T, request string, step func(*testrepo.GateFixture, string) string, seed func(testing.TB, string), fence ...string) landingFixture {
	t.Helper()
	return seededLandingFixture(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), true, step, seed, fence...)
}

// deliveryRoute is one real delivery route through the landing verb. seed writes its
// approved policy and board, deliverable is what the source assignment binds, and index is
// the board before the delivery.
type deliveryRoute struct {
	name, deliverable, index string
	seed                     func(testing.TB, string)
}

// deliveryRoutes holds the full spec route, the light-path tickets-only route, and the
// route of an outcome that owns no roadmap row.
var deliveryRoutes = []deliveryRoute{
	{name: "spec", deliverable: closureSpec, index: commitmenttest.ClosureIndex(), seed: func(t testing.TB, root string) {
		commitmenttest.SeedClosure(t, root, closureSpec)
	}},
	{name: "tickets-only", deliverable: commitmenttest.TicketsFolder, index: commitmenttest.ClosureIndex(), seed: func(t testing.TB, root string) {
		commitmenttest.SeedTicketsOnly(t, root, closureSpec)
	}},
	{name: "rowless", deliverable: closureSpec, index: commitmenttest.RowlessIndex, seed: func(t testing.TB, root string) {
		commitmenttest.SeedRowless(t, root, closureSpec)
	}},
}

// obligationFreeRoute delivers the spec that the tickets-only seed approves with no
// obligation of its outcome, which owns sources: the legacy binding that the completion
// landing refuses.
var obligationFreeRoute = deliveryRoute{name: "obligation-free", deliverable: closureSpec, index: commitmenttest.ClosureIndex(), seed: deliveryRoutes[1].seed}

func (route deliveryRoute) tickets() bool { return route.deliverable == commitmenttest.TicketsFolder }

// fixture builds the route's public landing fixture and binds the source assignment to the
// route's deliverable. Only a spec landing grades the implemented status. A non-nil step
// adds one prospective gate line, and each fence path joins the spec's ownership fence.
func (route deliveryRoute) fixture(t *testing.T, request string, step func(*testrepo.GateFixture, string) string, fence ...string) landingFixture {
	t.Helper()
	f := seededLandingFixture(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), !route.tickets(), step, route.seed, fence...)
	if route.tickets() {
		commitmenttest.Rebind(t, f.creation.Path, request, route.deliverable)
	}
	return f
}

// args are the landing arguments that deliver the route's deliverable from tip.
func (route deliveryRoute) args(request string, f landingFixture, tip string) []string {
	if route.tickets() {
		return ticketsOnlyLandArgs(request, f.base, tip, commitmenttest.TicketsSlug, f.creation.Path)
	}
	return landArgs(request, f.base, tip, f.creation.Path)
}

// resumeArgs are the arguments that resume the route's publication at published.
func (route deliveryRoute) resumeArgs(published, request string, f landingFixture) []string {
	args := resumeLandArgs(published, request, f.base, f.tip, f.creation.Path)
	if route.tickets() {
		args[slices.Index(args, "--spec")+1] = commitmenttest.TicketsSlug
	}
	return args
}

// closureState is what main and the local intent record say about the delivery outcome.
type closureState struct {
	main, index string
	rows        map[string]bool
	deliveries  []commitment.DeliveryFact
	claimed     bool
}

func readClosureState(t *testing.T, root string) closureState {
	t.Helper()
	state := closureState{main: gitOutput(t, root, "rev-parse", "main"), index: gitOutput(t, root, "show", "main:ROADMAP.md"), rows: map[string]bool{}}
	for _, row := range []string{"FT1", "FT2", "FT3"} {
		state.rows[row] = git.OK("-C", root, "cat-file", "-e", "main:roadmap/"+row+".md")
	}
	policy, _, err := commitrepo.Store{Root: root}.Policy()
	mustNoError(t, err)
	state.deliveries = policy.Deliveries
	state.claimed = deliveryClaimed(t, root)
	return state
}

// deliveryClaimed reports whether the local intent record claims the delivery outcome.
func deliveryClaimed(t *testing.T, root string) bool {
	t.Helper()
	ledger, err := intent.Read(root)
	mustNoError(t, err)
	if ledger.Commitment == nil {
		return false
	}
	return slices.ContainsFunc(ledger.Commitment.Claims, func(claim intent.OutcomeClaim) bool {
		return claim.Outcome == commitmenttest.DeliveryOutcome
	})
}

// requireOpenObligation fails unless main and the intent record still hold open, the
// state before the landing: the route's base board, no delivery fact, and the claim.
func requireOpenObligation(t *testing.T, f landingFixture, route deliveryRoute, open closureState) {
	t.Helper()
	got := readClosureState(t, f.root)
	if !reflect.DeepEqual(got, open) || got.main != f.base || got.index != strings.TrimSpace(route.index) || len(got.deliveries) != 0 || !got.claimed {
		t.Fatalf("unpublished delivery state = %+v, want main %s with the open obligation and its claim", got, f.base)
	}
}
