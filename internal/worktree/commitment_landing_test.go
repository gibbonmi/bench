package worktree

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// Verified delivery of the outcome publishes its spec flip, its delivery fact, and the
// removal of FT1, its detail owner, and its sequence entry in one commit.
func TestCommitmentDeliveryClosure(t *testing.T) {
	t.Parallel()
	request := "land-commitment-delivery-closure"
	f := closureLandingFixture(t, request, nil)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
		t.Fatalf("delivery landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	got := readClosureState(t, f.root)
	if parents := gitOutput(t, f.root, "rev-list", "--parents", "-n", "1", got.main); parents != got.main+" "+f.base+" "+f.tip {
		t.Fatalf("published parents = %q, want one publication over %s and %s", parents, f.base, f.tip)
	}
	if status := gitOutput(t, f.root, "show", "main:specs/x/spec.md"); !strings.Contains(status, "Status: implemented") {
		t.Fatalf("published spec = %q, want the implemented status", status)
	}
	index := strings.TrimSpace(commitmenttest.DeliveredIndex)
	if got.index != index || got.rows["FT1"] || !got.rows["FT2"] {
		t.Fatalf("published board = %q rows=%v, want FT1 and its sequence entry closed", got.index, got.rows)
	}
	if len(got.deliveries) != 1 || got.deliveries[0].Outcome != commitmenttest.DeliveryOutcome || got.deliveries[0].Source != f.tip {
		t.Fatalf("published delivery facts = %+v, want one fact naming the reviewed source %s", got.deliveries, f.tip)
	}
	if got.claimed {
		t.Fatal("delivered outcome kept its local claim")
	}
}

// A delivery that satisfies FT1 alone keeps the outcome's residual FT3, its sequence
// entry, its claim, and unrelated B.
func TestCommitmentPartialDelivery(t *testing.T) {
	t.Parallel()
	request := "land-commitment-partial-delivery"
	f := closureLandingFixture(t, request, nil, "FT3")
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
		t.Fatalf("partial landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	got := readClosureState(t, f.root)
	index := strings.TrimSpace(commitmenttest.ResidualIndex)
	if got.index != index || got.rows["FT1"] || !got.rows["FT2"] || !got.rows["FT3"] {
		t.Fatalf("partial board = %q rows=%v, want only FT1 closed", got.index, got.rows)
	}
	if len(got.deliveries) != 1 || !got.claimed {
		t.Fatalf("partial delivery = facts %+v claimed=%t, want one fact and the retained claim", got.deliveries, got.claimed)
	}
}

// A listed legacy run that delivers an approved spec in its scope closes that delivery
// with no binding. Reconciliation releases the continuation only when every approved
// deliverable in the scope is delivered, so a partly delivered scope stays open.
func TestCommitmentLegacyClosure(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name    string
		pending string
		open    bool
	}{
		{name: "delivered-scope"},
		{name: "partly-delivered-scope", pending: "specs/y/spec.md", open: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-legacy-closure-" + row.name
			scope := []string{"owned.txt", "reviews/x.md", closureSpec}
			f := seededClosureFixture(t, request, nil, func(t testing.TB, root string) {
				if row.pending == "" {
					commitmenttest.SeedClosure(t, root, closureSpec)
					return
				}
				commitmenttest.SeedClosure(t, root, closureSpec, "FT3")
				commitmenttest.ApprovePending(t, root, row.pending, "FT3")
				scope = append(scope, row.pending)
			})
			mustNoError(t, intent.Transact(f.root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
				ledger.Commitment = &intent.CommitmentState{Continuations: []intent.LegacyContinuation{{Assignment: f.creation.Assignment.ID, Request: f.creation.Assignment.Request, Scope: scope}}}
				return ledger, true, nil
			}, nil))
			r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
			if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
				t.Fatalf("legacy delivery landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			got := readClosureState(t, f.root)
			if got.rows["FT1"] || !got.rows["FT2"] || got.rows["FT3"] != row.open || len(got.deliveries) != 1 || got.deliveries[0].Source != f.tip {
				t.Fatalf("legacy delivery = rows %v facts %+v, want FT1 closed by one fact naming %s", got.rows, got.deliveries, f.tip)
			}
			ledger, err := intent.Read(f.root)
			mustNoError(t, err)
			if open := ledger.Commitment != nil && len(ledger.Commitment.Continuations) == 1; open != row.open {
				t.Fatalf("legacy continuation open = %t, want %t: %+v", open, row.open, ledger.Commitment)
			}
		})
	}
}

// A delivering source that also closes a residual or an unrelated obligation itself, by
// removing its row and its detail owner, refuses at admission on every delivery route.
// The refusal comes before the gate runs and before main moves.
func TestCommitmentClosureAdmission(t *testing.T) {
	t.Parallel()
	residual := deliveryRoute{name: "spec-residual", deliverable: closureSpec, index: commitmenttest.ClosureIndex("FT3"), seed: func(t testing.TB, root string) {
		commitmenttest.SeedClosure(t, root, closureSpec, "FT3")
	}}
	type closing struct {
		route deliveryRoute
		id    string
	}
	rows := []closing{{residual, "FT3"}, {residual, "FT2"}}
	for _, route := range deliveryRoutes[1:] {
		rows = append(rows, closing{route, "FT2"})
	}
	for _, row := range rows {
		t.Run(row.route.name+"-"+row.id, func(t *testing.T) {
			t.Parallel()
			id := row.id
			request := "land-commitment-closure-admission-" + row.route.name + "-" + id
			detail := "roadmap/" + id + ".md"
			f := row.route.fixture(t, request, nil, roadmap.RoadmapFile, detail)
			index, err := roadmap.Close([]byte(row.route.index), []string{id}, []string{commitmenttest.DeliveryOutcome, "B"})
			mustNoError(t, err)
			mustWrite(t, filepath.Join(f.creation.Path, roadmap.RoadmapFile), index, 0o644)
			gitRun(t, f.creation.Path, "rm", "-q", detail)
			gitRun(t, f.creation.Path, "add", roadmap.RoadmapFile)
			gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "close "+id)
			refreshLandingEvidence(t, f.creation.Path, f.base)
			tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
			r := runVerb(t, verbLand, f.call(row.route.args(request, f, tip)...))
			if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=commitment: ") || !strings.Contains(r.stdout, "commitment source \""+id+"\" refused") {
				t.Fatalf("closing %s = (%d, %q, %q), want the protected-source refusal", id, r.exit, r.stdout, r.stderr)
			}
			if tally, err := os.ReadFile(f.tally); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("gate tally = %q, %v; want no gate run", tally, err)
			}
			if main := gitOutput(t, f.root, "rev-parse", "main"); main != f.base {
				t.Fatalf("refused landing moved main to %s, want %s", main, f.base)
			}
		})
	}
}

// A red prospective gate publishes neither the delivery nor its closure on any delivery
// route. The gate ran once, so the refusal is its verdict on the admitted closure.
func TestCommitmentClosureGateRed(t *testing.T) {
	t.Parallel()
	for _, route := range deliveryRoutes {
		t.Run(route.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-closure-red-" + route.name
			f := route.fixture(t, request, func(*testrepo.GateFixture, string) string { return "exit 1\n" })
			open := readClosureState(t, f.root)
			r := runVerb(t, verbLand, f.call(route.args(request, f, f.tip)...))
			if r.exit != 1 || !strings.Contains(r.stdout, "prospective authorization refused") {
				t.Fatalf("red closure landing = (%d, %q, %q), want the gate refusal", r.exit, r.stdout, r.stderr)
			}
			if tally, err := os.ReadFile(f.tally); err != nil || string(tally) != "g" {
				t.Fatalf("gate tally = %q, %v; want the one prospective run", tally, err)
			}
			requireOpenObligation(t, f, route, open)
		})
	}
}

// A gate that a signal ends is not delivery evidence, so the obligation stays open.
func TestCommitmentInterruptedGate(t *testing.T) {
	t.Parallel()
	request := "land-commitment-closure-interrupted"
	f := closureLandingFixture(t, request, func(*testrepo.GateFixture, string) string { return "kill -TERM $$\n" })
	open := readClosureState(t, f.root)
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "prospective authorization refused") {
		t.Fatalf("interrupted closure landing = (%d, %q, %q), want the gate refusal", r.exit, r.stdout, r.stderr)
	}
	requireOpenObligation(t, f, deliveryRoutes[0], open)
}

// A partially written intent record after the identity proofs refuses before the gate.
// It leaves the record, the policy, and the destination as they were.
func TestCommitmentPersistenceBeforeGate(t *testing.T) {
	t.Parallel()
	request := "land-commitment-closure-torn-intent"
	f := closureLandingFixture(t, request, nil)
	address, err := intent.Address(f.root)
	mustNoError(t, err)
	var torn []byte
	j := defaultJoins()
	land := j.landReviewed
	j.landReviewed = func(ctx context.Context, request landing.ReviewedRequest, admission landing.Admission) (landing.ReviewedResult, error) {
		data, err := os.ReadFile(address)
		if err != nil {
			return landing.ReviewedResult{}, err
		}
		torn = data[:len(data)/2]
		if err := os.WriteFile(address, torn, 0o600); err != nil {
			return landing.ReviewedResult{}, err
		}
		return land(ctx, request, admission)
	}
	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=commitment: ") {
		t.Fatalf("torn-intent landing = (%d, %q, %q), want a commitment refusal", r.exit, r.stdout, r.stderr)
	}
	if data, err := os.ReadFile(address); err != nil || string(data) != string(torn) {
		t.Fatalf("torn intent record changed: %v", err)
	}
	if tally, err := os.ReadFile(f.tally); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("gate tally = %q, %v; want no gate run", tally, err)
	}
	if main := gitOutput(t, f.root, "rev-parse", "main"); main != f.base {
		t.Fatalf("torn-intent landing moved main to %s, want %s", main, f.base)
	}
	if index := gitOutput(t, f.root, "show", "main:ROADMAP.md"); index != strings.TrimSpace(commitmenttest.ClosureIndex()) {
		t.Fatalf("torn-intent landing changed the board to %q", index)
	}
}

// A local claim write that fails after the publication leaves the delivery published on
// every delivery route. The real reconciliation fails against an unwritable intent ledger
// directory. The resume finishes that write and publishes nothing a second time.
func TestCommitmentClosureResume(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot deny the intent ledger write")
	}
	for _, route := range deliveryRoutes {
		t.Run(route.name, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-closure-resume-" + route.name
			f := route.fixture(t, request, nil)
			address, err := intent.Address(f.root)
			mustNoError(t, err)
			dir := filepath.Dir(address)
			info, err := os.Stat(dir)
			mustNoError(t, err)
			t.Cleanup(func() { _ = os.Chmod(dir, info.Mode().Perm()) })
			broken := defaultJoins()
			reconcile := broken.reconcileCommitment
			broken.reconcileCommitment = func(root string) error {
				mustNoError(t, os.Chmod(dir, 0o500))
				err := reconcile(root)
				mustNoError(t, os.Chmod(dir, info.Mode().Perm()))
				if err == nil {
					t.Error("reconciliation wrote an unwritable intent ledger")
				}
				return err
			}
			r := runVerb(t, verbLand, f.callWith(broken, route.args(request, f, f.tip)...))
			r.mustViaJoins(t)
			if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:commitment") {
				t.Fatalf("interrupted delivery = (%d, %q, %q), want the commitment step named", r.exit, r.stdout, r.stderr)
			}
			published := readClosureState(t, f.root)
			if published.rows["FT1"] || published.index == strings.TrimSpace(route.index) || len(published.deliveries) != 1 || !published.claimed {
				t.Fatalf("published delivery = %+v, want the closure published and the claim still local", published)
			}
			r = runVerb(t, verbLand, f.call(route.resumeArgs(published.main, request, f)...))
			if r.exit != 0 || !strings.Contains(r.stdout, "published_commit="+published.main) {
				t.Fatalf("resume = (%d, %q, %q), want the original publication completed", r.exit, r.stdout, r.stderr)
			}
			resumed := readClosureState(t, f.root)
			if resumed.main != published.main || resumed.claimed || len(resumed.deliveries) != 1 {
				t.Fatalf("resumed state = %+v, want main %s and the claim released", resumed, published.main)
			}
			if tally, err := os.ReadFile(f.tally); err != nil || string(tally) != "g" {
				t.Fatalf("gate tally = %q, %v; want the one first-run gate", tally, err)
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
