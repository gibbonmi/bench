package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/spec"
)

// closedIndex is the board after the delivery of FT1 alone: B's row and B's sequence
// entry stay, and nothing names the delivered outcome.
const closedIndex = "# Roadmap\n\n## Parked\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. B\n"

// gradeClosure commits a reviewed source whose policy approves the example spec as the
// complete delivery of FT1. It then commits that source's exact closure, applies change
// to the closure, and grades the result through the completion oracle.
func gradeClosure(t *testing.T, change func(f *recordtest.Fixture)) error {
	t.Helper()
	f := recordtest.New(t, 1)
	commitmenttest.SeedClosure(t, f.Root, recordtest.Spec)
	record, err := reviewrecord.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(record, "retained completion record\n")
	f.Commit("approve the example delivery")
	source := f.Tip()
	policy, _, err := commitrepo.Store{Root: f.Root}.Policy()
	if err != nil {
		t.Fatal(err)
	}
	binding := policy.Milestones[0].Outcomes[0].Deliverables[0].Source
	policy.Deliveries = append(policy.Deliveries, commitment.DeliveryFact{Milestone: "M", Outcome: commitmenttest.DeliveryOutcome, Binding: binding.ID, Identity: binding.Identity, Source: source, Evidence: f.Git("rev-parse", source+":"+record)})
	commitmenttest.WritePolicy(t, f.Root, policy)
	staged, err := os.ReadFile(filepath.Join(f.Root, recordtest.Spec))
	if err != nil {
		t.Fatal(err)
	}
	implemented, err := spec.Implemented(staged)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(recordtest.Spec, string(implemented))
	f.Write("ROADMAP.md", closedIndex)
	if err := os.Remove(filepath.Join(f.Root, "roadmap", "FT1.md")); err != nil {
		t.Fatal(err)
	}
	change(f)
	f.Commit("close the example delivery")
	graded, err := captureProspectiveTree(f.Root, f.Tree())
	if err != nil {
		t.Fatal(err)
	}
	e := newGateEvaluation(f.Root)
	e.checkpoint, e.completionSource, e.prospective = Checkpoint{Spec: recordtest.Spec, Complete: true}, source, true
	_, err = e.completionTree(graded)
	return err
}

// The oracle accepts the exact closure and refuses a closure that keeps the delivered
// outcome's sequence entry. A path allowlist would accept both.
func TestCommitmentExactTransform(t *testing.T) {
	t.Parallel()
	err := gradeClosure(t, func(f *recordtest.Fixture) {
		f.Write("ROADMAP.md", strings.Replace(closedIndex, "1. B\n", "1. "+commitmenttest.DeliveryOutcome+"\n2. B\n", 1))
	})
	if err == nil || !strings.Contains(err.Error(), "ROADMAP.md differs from the exact closure transform") {
		t.Fatalf("kept sequence entry = %v, want the exact-transform refusal", err)
	}
	if err := gradeClosure(t, func(*recordtest.Fixture) {}); err != nil {
		t.Fatalf("exact closure refused: %v", err)
	}
}

// The oracle refuses a closure that keeps the closed detail owner, and a policy edit
// whose delivery fact names another source or evidence, omits the fact, or makes the
// policy executable. Each refusal comes from the exact closure comparison.
func TestCommitmentClosureNegatives(t *testing.T) {
	t.Parallel()
	const policyRefusal = "completion " + commitment.PolicyPath + " differs from the exact closure transform"
	for _, row := range []struct {
		name, want string
		change     func(*testing.T, *recordtest.Fixture)
	}{
		{name: "retained-detail", want: "completion keeps closed roadmap/FT1.md", change: func(t *testing.T, f *recordtest.Fixture) {
			f.Git("checkout", "HEAD", "--", "roadmap/FT1.md")
		}},
		{name: "wrong-source", want: policyRefusal, change: func(t *testing.T, f *recordtest.Fixture) {
			commitmenttest.EditPolicy(t, f.Root, func(policy *commitment.Policy) { policy.Deliveries[0].Source = f.Git("rev-parse", "HEAD~1") })
		}},
		{name: "wrong-evidence", want: policyRefusal, change: func(t *testing.T, f *recordtest.Fixture) {
			commitmenttest.EditPolicy(t, f.Root, func(policy *commitment.Policy) { policy.Deliveries[0].Evidence = policy.Deliveries[0].Identity })
		}},
		{name: "omitted-fact", want: policyRefusal, change: func(t *testing.T, f *recordtest.Fixture) {
			commitmenttest.EditPolicy(t, f.Root, func(policy *commitment.Policy) { policy.Deliveries = nil })
		}},
		{name: "executable-policy", want: policyRefusal, change: func(t *testing.T, f *recordtest.Fixture) {
			if err := os.Chmod(filepath.Join(f.Root, commitment.PolicyPath), 0o755); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			err := gradeClosure(t, func(f *recordtest.Fixture) { row.change(t, f) })
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("%s = %v, want a refusal naming %q", row.name, err, row.want)
			}
		})
	}
}

// The oracle refuses a closure that also deletes the unrelated obligation B.
func TestCommitmentExtraClosure(t *testing.T) {
	t.Parallel()
	err := gradeClosure(t, func(f *recordtest.Fixture) {
		if err := os.Remove(filepath.Join(f.Root, "roadmap", "FT2.md")); err != nil {
			t.Fatal(err)
		}
	})
	if err == nil || !strings.Contains(err.Error(), "completion composition changes roadmap/FT2.md") {
		t.Fatalf("extra deletion of B = %v, want the destination-delta refusal", err)
	}
}
