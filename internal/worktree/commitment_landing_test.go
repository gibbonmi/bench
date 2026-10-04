package worktree

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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

// closureLandingFixture is the graded public landing fixture whose policy approves
// specs/x/spec.md as the complete delivery of FT1. The delivery outcome also owns each
// residual row, and outcome B owns FT2. A non-nil step adds one prospective gate line.
func closureLandingFixture(t *testing.T, request string, step func(*testrepo.GateFixture, string) string, residual ...string) landingFixture {
	t.Helper()
	return seededClosureFixture(t, request, step, func(t testing.TB, root string) {
		commitmenttest.SeedClosure(t, root, "specs/x/spec.md", residual...)
	})
}

// seededClosureFixture is the graded public landing fixture whose policy seed writes. Each
// fence path joins the spec's ownership fence, so the source may change that path.
func seededClosureFixture(t *testing.T, request string, step func(*testrepo.GateFixture, string) string, seed func(testing.TB, string), fence ...string) landingFixture {
	t.Helper()
	return seededLandingFixture(t, request, "", "", filepath.Join(t.TempDir(), "bench-home"), true, step, seed, fence...)
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
	ledger, err := intent.Read(root)
	mustNoError(t, err)
	if ledger.Commitment != nil {
		for _, claim := range ledger.Commitment.Claims {
			state.claimed = state.claimed || claim.Outcome == commitmenttest.DeliveryOutcome
		}
	}
	return state
}

// requireOpenObligation fails unless main still carries the base board, no delivery
// fact, and the retained delivery claim.
func requireOpenObligation(t *testing.T, f landingFixture) {
	t.Helper()
	got := readClosureState(t, f.root)
	if got.main != f.base || got.index != strings.TrimSpace(commitmenttest.ClosureIndex()) || !got.rows["FT1"] || !got.rows["FT2"] || len(got.deliveries) != 0 || !got.claimed {
		t.Fatalf("unpublished delivery state = %+v, want main %s with the open obligation and its claim", got, f.base)
	}
}

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
	const index = "# Roadmap\n\n## Parked\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. B"
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
	const index = "# Roadmap\n\n## Parked\n\n**FT3 — delivery**\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. delivery\n2. B"
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
			scope := []string{"owned.txt", "reviews/x.md", "specs/x/spec.md"}
			f := seededClosureFixture(t, request, nil, func(t testing.TB, root string) {
				if row.pending == "" {
					commitmenttest.SeedClosure(t, root, "specs/x/spec.md")
					return
				}
				commitmenttest.SeedClosure(t, root, "specs/x/spec.md", "FT3")
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
// removing its row and its detail owner, refuses at admission. The refusal comes before
// the gate runs and before main moves.
func TestCommitmentClosureAdmission(t *testing.T) {
	t.Parallel()
	for _, id := range []string{"FT3", "FT2"} {
		t.Run(id, func(t *testing.T) {
			t.Parallel()
			request := "land-commitment-closure-admission-" + id
			detail := "roadmap/" + id + ".md"
			f := seededClosureFixture(t, request, nil, func(t testing.TB, root string) {
				commitmenttest.SeedClosure(t, root, "specs/x/spec.md", "FT3")
			}, roadmap.RoadmapFile, detail)
			index, err := roadmap.Close([]byte(commitmenttest.ClosureIndex("FT3")), []string{id}, []string{commitmenttest.DeliveryOutcome, "B"})
			mustNoError(t, err)
			mustWrite(t, filepath.Join(f.creation.Path, roadmap.RoadmapFile), index, 0o644)
			gitRun(t, f.creation.Path, "rm", "-q", detail)
			gitRun(t, f.creation.Path, "add", roadmap.RoadmapFile)
			gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "close "+id)
			refreshLandingEvidence(t, f.creation.Path, f.base)
			tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
			r := runVerb(t, verbLand, f.call(landArgs(request, f.base, tip, f.creation.Path)...))
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

// A red prospective gate publishes neither the delivery nor its closure.
func TestCommitmentClosureGateRed(t *testing.T) {
	t.Parallel()
	request := "land-commitment-closure-red"
	f := closureLandingFixture(t, request, func(*testrepo.GateFixture, string) string { return "exit 1\n" })
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "prospective authorization refused") {
		t.Fatalf("red closure landing = (%d, %q, %q), want the gate refusal", r.exit, r.stdout, r.stderr)
	}
	requireOpenObligation(t, f)
}

// A gate that a signal ends is not delivery evidence, so the obligation stays open.
func TestCommitmentInterruptedGate(t *testing.T) {
	t.Parallel()
	request := "land-commitment-closure-interrupted"
	f := closureLandingFixture(t, request, func(*testrepo.GateFixture, string) string { return "kill -TERM $$\n" })
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "prospective authorization refused") {
		t.Fatalf("interrupted closure landing = (%d, %q, %q), want the gate refusal", r.exit, r.stdout, r.stderr)
	}
	requireOpenObligation(t, f)
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

// A local claim write that fails after the publication leaves the delivery published.
// The real reconciliation fails against an unwritable intent ledger directory. The resume
// finishes that write and publishes nothing a second time.
func TestCommitmentClosureResume(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot deny the intent ledger write")
	}
	request := "land-commitment-closure-resume"
	f := closureLandingFixture(t, request, nil)
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
	r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:commitment") {
		t.Fatalf("interrupted delivery = (%d, %q, %q), want the commitment step named", r.exit, r.stdout, r.stderr)
	}
	published := readClosureState(t, f.root)
	if published.rows["FT1"] || len(published.deliveries) != 1 || !published.claimed {
		t.Fatalf("published delivery = %+v, want the closure published and the claim still local", published)
	}
	r = runVerb(t, verbLand, f.call(resumeLandArgs(published.main, request, f.base, f.tip, f.creation.Path)...))
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
}
