package gate

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/landing/published"
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

// gradeTicketsClose commits a reviewed source whose policy approves the tickets-only folder
// as the complete delivery of FT1. It publishes that source through the landing's own
// transform, applies change to the published tree, and grades the result through the
// completion that the landing binds for the folder.
func gradeTicketsClose(t *testing.T, change func(root, source string)) error {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, root, recordtest.Spec, "# example\n\nStatus: staged\n")
	commitmenttest.SeedTicketsOnly(t, root, recordtest.Spec)
	commitmenttest.Commit(t, root, "approve the tickets-only delivery")
	source := gittest.Output(t, root, "rev-parse", "HEAD")
	tree, err := published.Tree(root, gittest.Output(t, root, "rev-parse", "HEAD^{tree}"), commitmenttest.TicketsFolder, source)
	if err != nil {
		t.Fatal(err)
	}
	gittest.Output(t, root, "read-tree", "--reset", "-u", tree)
	change(root, source)
	gittest.Output(t, root, "add", "-A")
	graded, err := captureProspectiveTree(root, gittest.Output(t, root, "write-tree"))
	if err != nil {
		t.Fatal(err)
	}
	e := checkpointEvaluation(WithCompletion(context.Background(), commitmenttest.TicketsFolder, source), newGateEvaluation(root))
	e.prospective = true
	_, err = e.applyCheckpoint(graded, subject{})
	return err
}

// The oracle accepts the exact tickets-only close. It refuses a close that keeps the
// closed folder, the satisfied row's detail owner, or the delivered sequence entry. Each
// refusal is the only difference from the exact close.
func TestCommitmentTicketsOnlyTransform(t *testing.T) {
	t.Parallel()
	if err := gradeTicketsClose(t, func(string, string) {}); err != nil {
		t.Fatalf("exact tickets-only close refused: %v", err)
	}
	for _, row := range []struct {
		name, want string
		change     func(t *testing.T, root, source string)
	}{
		{name: "kept-folder", want: "completion keeps closed " + commitmenttest.TicketsFolder + "/", change: func(t *testing.T, root, source string) {
			gittest.Output(t, root, "checkout", source, "--", commitmenttest.TicketsFolder)
		}},
		{name: "kept-detail-owner", want: "completion keeps closed roadmap/FT1.md", change: func(t *testing.T, root, source string) {
			gittest.Output(t, root, "checkout", source, "--", "roadmap/FT1.md")
		}},
		{name: "kept-sequence-entry", want: "ROADMAP.md differs from the exact closure transform", change: func(t *testing.T, root, _ string) {
			commitmenttest.Write(t, root, "ROADMAP.md", strings.Replace(closedIndex, "1. B\n", "1. "+commitmenttest.DeliveryOutcome+"\n2. B\n", 1))
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			err := gradeTicketsClose(t, func(root, source string) { row.change(t, root, source) })
			if err == nil || !strings.Contains(err.Error(), row.want) {
				t.Fatalf("%s = %v, want a refusal naming %q", row.name, err, row.want)
			}
		})
	}
}

// On each delivery route, the oracle refuses a completion tree that changes one byte of
// the unrelated detail owner FT2. FT2 lies in the directory that the closure edits, so a
// path or directory allowlist would accept the change.
func TestCommitmentUnrelatedByte(t *testing.T) {
	t.Parallel()
	const want = "completion composition changes roadmap/FT2.md"
	flip := func(t *testing.T, root string) {
		path := filepath.Join(root, "roadmap", "FT2.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		changed := strings.Replace(string(data), "Keep", "Keap", 1)
		if len(changed) != len(data) || changed == string(data) {
			t.Fatalf("one-byte change of %q failed", data)
		}
		commitmenttest.Write(t, root, "roadmap/FT2.md", changed)
	}
	for _, row := range []struct {
		name  string
		grade func(t *testing.T) error
	}{
		{name: "spec", grade: func(t *testing.T) error {
			return gradeClosure(t, func(f *recordtest.Fixture) { flip(t, f.Root) })
		}},
		{name: "tickets-only", grade: func(t *testing.T) error {
			return gradeTicketsClose(t, func(root, _ string) { flip(t, root) })
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			if err := row.grade(t); err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("%s close with a changed FT2 byte = %v, want a refusal naming %q", row.name, err, want)
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
