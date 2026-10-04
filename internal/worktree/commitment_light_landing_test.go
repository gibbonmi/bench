// Light-path, rowless, and legacy delivery landing tests: each closes its delivery through
// the same broker transform as the spec route.
package worktree

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/roadmap"
)

// A listed legacy run that delivers an approved spec or tickets-only folder in its scope
// closes that delivery with no binding. Reconciliation releases the continuation only
// when every approved deliverable in the scope is delivered, so a partly delivered scope
// stays open.
func TestCommitmentLegacyClosure(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name    string
		route   deliveryRoute
		pending string
		open    bool
	}{
		{name: "delivered-scope", route: deliveryRoutes[0]},
		{name: "partly-delivered-scope", route: deliveryRoutes[0], pending: "specs/y/spec.md", open: true},
		{name: "tickets-only-scope", route: deliveryRoutes[1]},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-legacy-closure-" + row.name
			scope := []string{"owned.txt", "reviews/x.md", row.route.deliverable}
			route := row.route
			if row.pending != "" {
				route.seed = func(t testing.TB, root string) {
					commitmenttest.SeedClosure(t, root, closureSpec, "FT3")
					commitmenttest.ApprovePending(t, root, row.pending, "FT3")
				}
				scope = append(scope, row.pending)
			}
			f := route.fixture(t, request, nil)
			mustNoError(t, intent.Transact(f.root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
				ledger.Commitment = &intent.CommitmentState{Continuations: []intent.LegacyContinuation{{Assignment: f.creation.Assignment.ID, Request: f.creation.Assignment.Request, Scope: scope}}}
				return ledger, true, nil
			}, nil))
			r := runVerb(t, verbLand, f.call(route.args(request, f, f.tip)...))
			if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
				t.Fatalf("legacy delivery landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			got := readClosureState(t, f.root)
			if got.rows["FT1"] || !got.rows["FT2"] || got.rows["FT3"] != row.open || len(got.deliveries) != 1 || got.deliveries[0].Source != f.tip {
				t.Fatalf("legacy delivery = rows %v facts %+v, want FT1 closed by one fact naming %s", got.rows, got.deliveries, f.tip)
			}
			if route.tickets() && git.OK("-C", f.root, "cat-file", "-e", "main:"+commitmenttest.TicketsFolder) {
				t.Fatal("legacy delivery kept the closed tickets-only folder")
			}
			ledger, err := intent.Read(f.root)
			mustNoError(t, err)
			if open := ledger.Commitment != nil && len(ledger.Commitment.Continuations) == 1; open != row.open {
				t.Fatalf("legacy continuation open = %t, want %t: %+v", open, row.open, ledger.Commitment)
			}
		})
	}
}

// A tickets-only delivery that completely satisfies FT1 publishes, in one commit, the
// folder's close, its delivery fact, and the removal of FT1, its detail owner, and its
// sequence entry. It flips no spec. The fact's evidence is the reviewed folder with the
// acceptance that the source records. A delivery that leaves the residual FT3 keeps that
// row, its sequence entry, and the claim.
func TestCommitmentTicketsOnlyClosure(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, index string
		residual    []string
	}{
		{name: "complete", index: commitmenttest.DeliveredIndex},
		{name: "partial", index: commitmenttest.ResidualIndex, residual: []string{"FT3"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-tickets-only-" + row.name
			route := deliveryRoute{deliverable: commitmenttest.TicketsFolder, seed: func(t testing.TB, root string) {
				commitmenttest.SeedTicketsOnly(t, root, closureSpec, row.residual...)
			}}
			f := route.fixture(t, request, nil)
			commitInWorktree(t, f.creation.Path, commitmenttest.TicketsFolder+"/tickets/one.md", "Light path ticket.\n\n- [x] Accepted.\n", "accept the ticket")
			tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
			r := runVerb(t, verbLand, f.call(route.args(request, f, tip)...))
			if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
				t.Fatalf("tickets-only delivery = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			got := readClosureState(t, f.root)
			if parents := gitOutput(t, f.root, "rev-list", "--parents", "-n", "1", got.main); parents != got.main+" "+f.base+" "+tip {
				t.Fatalf("published parents = %q, want one publication over %s and %s", parents, f.base, tip)
			}
			if git.OK("-C", f.root, "cat-file", "-e", "main:"+commitmenttest.TicketsFolder) {
				t.Fatal("published tree kept the closed tickets-only folder")
			}
			if status := gitOutput(t, f.root, "show", "main:"+closureSpec); !strings.Contains(status, "Status: staged") {
				t.Fatalf("published spec = %q, want the staged spec unchanged", status)
			}
			partial := len(row.residual) > 0
			if got.index != strings.TrimSpace(row.index) || got.rows["FT1"] || !got.rows["FT2"] || got.rows["FT3"] != partial {
				t.Fatalf("published board = %q rows=%v, want only FT1 closed", got.index, got.rows)
			}
			evidence := commitrepo.TreeIdentity(gitOutput(t, f.root, "rev-parse", tip+":"+commitmenttest.TicketsFolder))
			if len(got.deliveries) != 1 || got.deliveries[0].Binding != "tickets" || got.deliveries[0].Source != tip || got.deliveries[0].Evidence != evidence {
				t.Fatalf("published delivery facts = %+v, want one tickets fact naming %s with evidence %s", got.deliveries, tip, evidence)
			}
			if got.claimed != partial {
				t.Fatalf("delivery claim retained = %t, want %t", got.claimed, partial)
			}
		})
	}
}

// A delivery of an outcome that owns no roadmap row records that outcome's fact alone and
// creates and deletes no row. A project with no board keeps none. A board keeps every row
// and detail owner, and its sequence drops only the delivered outcome.
func TestCommitmentNoRoadmapOwner(t *testing.T) {
	t.Parallel()
	rowless := deliveryRoutes[2]
	index := strings.TrimSpace(commitmenttest.DeliveredIndex)
	for _, row := range []struct {
		route          deliveryRoute
		binding, index string
		remaining      []string
	}{
		{route: deliveryRoute{name: "spec-without-board", deliverable: closureSpec, seed: func(t testing.TB, root string) {
			commitmenttest.SeedAdmission(t, root, closureSpec)
		}}, binding: "spec", remaining: []string{}},
		{route: rowless, binding: "spec", index: index, remaining: []string{"B"}},
		{route: deliveryRoute{name: "tickets-with-board", deliverable: commitmenttest.TicketsFolder, seed: rowless.seed}, binding: "tickets", index: index, remaining: []string{"B"}},
	} {
		t.Run(row.route.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-no-roadmap-owner-" + row.route.name
			f := row.route.fixture(t, request, nil)
			board := func(revision string) string {
				return gitOutput(t, f.root, "ls-tree", "-r", "--name-only", revision, "--", roadmap.RoadmapFile, roadmap.RoadmapDir)
			}
			r := runVerb(t, verbLand, f.call(row.route.args(request, f, f.tip)...))
			if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
				t.Fatalf("rowless delivery = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			if before, after := board(f.base), board("main"); after != before {
				t.Fatalf("published board files = %q, want the base files %q", after, before)
			}
			published := ""
			if git.OK("-C", f.root, "cat-file", "-e", "main:"+roadmap.RoadmapFile) {
				published = gitOutput(t, f.root, "show", "main:"+roadmap.RoadmapFile)
			}
			if published != row.index {
				t.Fatalf("published board = %q, want %q", published, row.index)
			}
			policy, _, err := commitrepo.Store{Root: f.root}.Policy()
			mustNoError(t, err)
			if len(policy.Deliveries) != 1 || policy.Deliveries[0].Outcome != commitmenttest.DeliveryOutcome || policy.Deliveries[0].Binding != row.binding || policy.Deliveries[0].Source != f.tip {
				t.Fatalf("published delivery facts = %+v, want one %s fact naming %s", policy.Deliveries, row.binding, f.tip)
			}
			if remaining := commitment.Remaining(policy); !slices.Equal(remaining, row.remaining) {
				t.Fatalf("remaining outcomes = %v, want %v", remaining, row.remaining)
			}
			if deliveryClaimed(t, f.root) {
				t.Fatal("delivered rowless outcome kept its local claim")
			}
		})
	}
}
