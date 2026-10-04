package worktree

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
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
