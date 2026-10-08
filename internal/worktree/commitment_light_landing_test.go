// Light-path, rowless, and legacy delivery landing tests: each closes its delivery through
// the same broker transform as the spec route.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
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
			acceptTicket(t, f.creation.Path)
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

// An unbound assignment lands its light-path change with --spec naming the change's folder.
// The landing publishes the change, closes the folder, and records no delivery fact. A
// spec-less landing of the same change refuses before publication and names the --spec
// route. A landing whose change has a path outside the ticket follows its printed route.
func TestCommitmentLightPathLanding(t *testing.T) {
	t.Parallel()
	t.Run("outside-route", func(t *testing.T) {
		t.Parallel()
		followLightPathOutsideLanding(t)
	})
	request := "land-commitment-light-path"
	f := lightLandingFixture(t, request)
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, f.tip, f.creation.Path)...))
	if main := gitOutput(t, f.root, "rev-parse", "main"); r.exit == 0 || !strings.Contains(r.stdout, "--spec") || main != f.base {
		t.Fatalf("spec-less light-path landing = (%d, %q, %q) with main %s, want a refusal that names --spec and leaves main at %s", r.exit, r.stdout, r.stderr, main, f.base)
	}
	r = runVerb(t, verbLand, f.call(ticketsOnlyLandArgs(request, f.base, f.tip, lightLandingSlug, f.creation.Path)...))
	published := gitOutput(t, f.root, "rev-parse", "main")
	if r.exit != 0 || published == f.base || !strings.Contains(r.stdout, "published_commit="+published+",") {
		t.Fatalf("light-path landing = (%d, %q, %q) with main %s, want main moved to the published commit", r.exit, r.stdout, r.stderr, published)
	}
	if folder := landing.ClosedFolderPath(lightLandingSlug); git.OK("-C", f.root, "cat-file", "-e", published+":"+folder) {
		t.Fatalf("published tree kept %s", folder)
	}
	if before, after := gitOutput(t, f.root, "show", f.base+":"+commitment.PolicyPath), gitOutput(t, f.root, "show", published+":"+commitment.PolicyPath); after != before {
		t.Fatalf("published policy = %q, want the base policy %q", after, before)
	}
}

// lightPathLandingFixture produces one commitment face that a light-path landing prints.
// build makes the light-path landing fixture and returns the arguments of the landing that
// the face refuses. carry carries out a printed instruction, keyed by the step's index in
// the face's route.
type lightPathLandingFixture struct {
	face  string
	build func(t *testing.T, request string) (landingFixture, []string)
	carry map[int]func(t *testing.T, f landingFixture)
}

// lightPathLandingFixtures are the fixtures of the publication faces, which the landing walk
// proves.
func lightPathLandingFixtures() []lightPathLandingFixture {
	return []lightPathLandingFixture{{
		// The ticket covers the change, and the caller lands it with no --spec.
		face: refusalroute.CommitmentLightPathSpec,
		build: func(t *testing.T, request string) (landingFixture, []string) {
			f := lightLandingFixture(t, request)
			return f, specLessLandArgs(request, f.base, f.tip, f.creation.Path)
		},
	}}
}

// lightPathOutsideLanding is the landing fixture of commitment-light-path-outside: the source
// also commits extra.txt, which the ticket's Writes line does not list. The walk widens the
// line and gives the printed commit of the ticket a lane to pass.
func lightPathOutsideLanding() lightPathLandingFixture {
	return lightPathLandingFixture{
		face: refusalroute.CommitmentLightPathOutside,
		build: func(t *testing.T, request string) (landingFixture, []string) {
			f := lightLandingFixture(t, request)
			commitInWorktree(t, f.creation.Path, "extra.txt", "outside the ticket\n", "outside the ticket")
			f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
			return f, ticketsOnlyLandArgs(request, f.base, f.tip, lightLandingSlug, f.creation.Path)
		},
		carry: map[int]func(*testing.T, landingFixture){0: func(t *testing.T, f landingFixture) {
			commitmenttest.WriteLightTicket(t, f.creation.Path, lightPathTicket, "owned.txt", "extra.txt")
			plantSourceCommitLane(t, f.creation)
		}},
	}
}

// followLightPathOutsideLanding is the landing arm of RR63. The landing prints the
// light-path route in the shape its own grammar takes: the ticket commits first, and the
// re-run names the repaired source tip. The walk follows the route verbatim to a clean
// landing.
func followLightPathOutsideLanding(t *testing.T) {
	for _, face := range refusalroute.Faces(refusalroute.Commitment) {
		if face.Name == refusalroute.CommitmentLightPathOutside {
			followLightPathLanding(t, installedWrapper(t, testRunBinary(t)), face, lightPathOutsideLanding())
			return
		}
	}
	t.Fatalf("the registry declares no face %q", refusalroute.CommitmentLightPathOutside)
}

// followLightPathLanding drives one light-path fixture: the landing prints one route and a
// sentence that holds no route, the walk carries out that route, and its last step lands.
func followLightPathLanding(t *testing.T, wrapper string, face refusalroute.Face, fixture lightPathLandingFixture) {
	request := "light-path-" + face.Name
	f, args := fixture.build(t, request)
	r := runVerb(t, verbLand, f.call(args...))
	detail, _ := recordField(r.stdout, "refused{", "detail", refusalroute.NextField)
	next, printed := recordField(r.stdout, "refused{", refusalroute.NextField)
	if r.exit != 1 || !printed || strings.Count(r.stdout, refusalroute.NextField+"=") != 1 || strings.Contains(detail, "bench ") || strings.Contains(detail, " --") {
		t.Fatalf("%s landing = (%d, %q, %q), want exit 1, one next= route, and a sentence with no route", face.Name, r.exit, r.stdout, r.stderr)
	}
	steps := routetest.Steps(t, face, next, "", nil)
	fill := operatorFill(producedFace{f: f, request: request})
	carry := func(index int) (func(), bool) {
		step, carried := fixture.carry[index]
		return func() { step(t, f) }, carried
	}
	last := routetest.Follow(t, face, steps, carry, func(step string) verbResult {
		ran := runPrintedStep(t, wrapper, f.repoHome, fill(t, step))
		// The commit lane is the walk's scaffold, so it leaves the source before the landing.
		mustRemove(t, filepath.Join(f.creation.Path, ".bench", "phases.json"))
		return ran
	})
	if _, landed := landedField(last.stdout, "worktree"); last.exit != 0 || !landed {
		t.Fatalf("%s re-run = (%d, %q, %q), want exit 0 and a landed record", face.Name, last.exit, last.stdout, last.stderr)
	}
}

// A delivery of an outcome that owns no roadmap row records that outcome's fact alone and
// creates and deletes no row. A project with no board keeps none, and a board keeps every
// row and detail owner. The outcome is delivered only when every approved deliverable is
// delivered: with a spec and a tickets-only folder approved, the first delivery in either
// order keeps the sequence entry and the claim, and the second removes them.
func TestCommitmentNoRoadmapOwner(t *testing.T) {
	t.Parallel()
	pair := func(t testing.TB, root string) { commitmenttest.SeedRowlessPair(t, root, closureSpec) }
	open, delivered := strings.TrimSpace(commitmenttest.RowlessIndex), strings.TrimSpace(commitmenttest.DeliveredIndex)
	for _, row := range []struct {
		name         string
		seed         func(testing.TB, string)
		deliverables []string
		index        string
		remaining    []string
	}{
		{name: "spec-without-board", seed: func(t testing.TB, root string) { commitmenttest.SeedAdmission(t, root, closureSpec) }, deliverables: []string{closureSpec}, remaining: []string{}},
		{name: "spec-with-board", seed: deliveryRoutes[2].seed, deliverables: []string{closureSpec}, index: delivered, remaining: []string{"B"}},
		{name: "spec-then-tickets", seed: pair, deliverables: []string{closureSpec, commitmenttest.TicketsFolder}, index: delivered, remaining: []string{"B"}},
		{name: "tickets-then-spec", seed: pair, deliverables: []string{commitmenttest.TicketsFolder, closureSpec}, index: delivered, remaining: []string{"B"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-no-roadmap-owner-" + row.name
			first := deliveryRoute{deliverable: row.deliverables[0], seed: row.seed}
			f := first.fixture(t, request, nil)
			board := func(revision string) string {
				return gitOutput(t, f.root, "ls-tree", "-r", "--name-only", revision, "--", roadmap.RoadmapFile, roadmap.RoadmapDir)
			}
			for i, deliverable := range row.deliverables {
				var r verbResult
				tip := f.tip
				if i == 0 {
					r = runVerb(t, verbLand, f.call(first.args(request, f, f.tip)...))
				} else {
					r, tip = landNextDelivery(t, f, request+"-next", deliverable)
				}
				if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
					t.Fatalf("rowless delivery of %s = (%d, %q, %q)", deliverable, r.exit, r.stdout, r.stderr)
				}
				if before, after := board(f.base), board("main"); after != before {
					t.Fatalf("published board files = %q, want the base files %q", after, before)
				}
				final := i == len(row.deliverables)-1
				index, remaining := open, []string{commitmenttest.DeliveryOutcome, "B"}
				if final {
					index, remaining = row.index, row.remaining
				}
				published := ""
				if git.OK("-C", f.root, "cat-file", "-e", "main:"+roadmap.RoadmapFile) {
					published = gitOutput(t, f.root, "show", "main:"+roadmap.RoadmapFile)
				}
				if published != index {
					t.Fatalf("published board after %s = %q, want %q", deliverable, published, index)
				}
				binding := "spec"
				if (deliveryRoute{deliverable: deliverable}).tickets() {
					binding = "tickets"
				}
				policy, _, err := commitrepo.Store{Root: f.root}.Policy()
				mustNoError(t, err)
				if len(policy.Deliveries) != i+1 || policy.Deliveries[i].Outcome != commitmenttest.DeliveryOutcome || policy.Deliveries[i].Binding != binding || policy.Deliveries[i].Source != tip {
					t.Fatalf("published delivery facts = %+v, want %d facts, the last a %s fact naming %s", policy.Deliveries, i+1, binding, tip)
				}
				if got := commitment.Remaining(policy); !slices.Equal(got, remaining) {
					t.Fatalf("remaining outcomes after %s = %v, want %v", deliverable, got, remaining)
				}
				if claimed := deliveryClaimed(t, f.root); claimed == final {
					t.Fatalf("rowless outcome claim after %s = %t, want %t", deliverable, claimed, !final)
				}
			}
		})
	}
}

// The completion landing of a legacy binding that names no obligation refuses before the
// gate runs and before main moves, and its route hands the repair plan to the reviewer.
func TestCommitmentLandingRefusesObligationFreeDelivery(t *testing.T) {
	t.Parallel()
	request := "land-commitment-obligation-free"
	route := obligationFreeRoute
	f := route.fixture(t, request, nil)
	open := readClosureState(t, f.root)
	r := runVerb(t, verbLand, f.call(route.args(request, f, f.tip)...))
	if r.exit == 0 || !strings.Contains(r.stdout, ",next=reviewer: bench commitment plan --input <file>") {
		t.Fatalf("obligation-free landing = (%d, %q, %q), want the reviewer route of the plan command", r.exit, r.stdout, r.stderr)
	}
	if tally, err := os.ReadFile(f.tally); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("gate tally = %q, %v; want no gate run", tally, err)
	}
	requireOpenObligation(t, f, route, open)
}

// acceptTicket commits the accepted ticket of the tickets-only folder in worktree.
func acceptTicket(t *testing.T, worktree string) {
	t.Helper()
	commitInWorktree(t, worktree, commitmenttest.TicketsFolder+"/tickets/one.md", "Light path ticket.\n\n- [x] Accepted.\n", "accept the ticket")
}

// landNextDelivery lands deliverable from a new assignment that request names on the
// current main of f. The assignment joins the open delivery outcome, binds deliverable,
// and commits one reviewed change. It returns the landing result and the reviewed tip.
func landNextDelivery(t *testing.T, f landingFixture, request, deliverable string) (verbResult, string) {
	t.Helper()
	route := deliveryRoute{deliverable: deliverable}
	next := f
	next.base = gitOutput(t, f.root, "rev-parse", "main")
	next.creation = mustCreate(t, f.root, f.home, request, "next delivery")
	commitmenttest.Admit(t, next.creation.Path, request, deliverable)
	if route.tickets() {
		acceptTicket(t, next.creation.Path)
	} else {
		commitInWorktree(t, next.creation.Path, "owned.txt", "next reviewed bytes\n", "next reviewed source")
		refreshLandingEvidence(t, next.creation.Path, next.base)
	}
	next.tip = gitOutput(t, next.creation.Path, "rev-parse", "HEAD")
	return runVerb(t, verbLand, next.call(route.args(request, next, next.tip)...)), next.tip
}
