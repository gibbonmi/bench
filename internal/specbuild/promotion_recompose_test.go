package specbuild

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestPromoteRecomposesAWorkingAdvanceBeforeGate(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	before := loadRun(t, fixture.service)
	owner := &recompositionGate{}
	fixture.service.gate = owner
	advanceWorking(t, fixture.root)
	working := git(t, fixture.root, "rev-parse", "HEAD")
	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.Next != "bench spec build review build demo" || owner.calls != 1 {
		t.Fatalf("recomposition = %#v, %v; gate calls=%d", status, err, owner.calls)
	}
	if owner.branch != before.Branch || owner.tip != working || owner.expected != before.Base {
		t.Fatalf("bootstrap subject = branch %q tip %q expected %q", owner.branch, owner.tip, owner.expected)
	}
	if got := git(t, fixture.root, "rev-parse", "refs/bench/green/"+before.Branch); got != working {
		t.Fatalf("green marker = %s, want recomposed base %s", got, working)
	}
	after := loadRun(t, fixture.service)
	if after.Base != working || after.CandidateTip == before.CandidateTip || after.Review != nil || git(t, fixture.root, "rev-parse", "HEAD") != working || git(t, fixture.root, "rev-parse", after.Candidate) != after.CandidateTip {
		t.Fatalf("recomposed run = %#v", after)
	}
	for _, path := range []string{"advanced.txt", "internal/specbuild/checkpoint-change.go"} {
		if git(t, fixture.root, "show", after.CandidateTip+":"+path) == "" {
			t.Fatalf("recomposed candidate omitted %s", path)
		}
	}
}

func TestPromoteRecompositionRefusesBootstrapFailureWithoutMutation(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	advanceWorking(t, fixture.root)
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	beforeState, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	before := loadRun(t, fixture.service)
	working := git(t, fixture.root, "rev-parse", "HEAD")
	candidate := git(t, fixture.root, "rev-parse", before.Candidate)
	owner := &recompositionGate{err: errors.New("injected exact-evidence failure")}
	fixture.service.gate = owner

	status, promoteErr := fixture.service.Promote(t.Context(), "build demo")
	if promoteErr == nil || !strings.Contains(promoteErr.Error(), "run bench gate, then retry promote") {
		t.Fatalf("Promote error = %v", promoteErr)
	}
	if status.Next != "bench spec build promote build demo" {
		t.Fatalf("Promote next = %q", status.Next)
	}
	afterState, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(afterState) != string(beforeState) || git(t, fixture.root, "rev-parse", before.Candidate) != candidate || git(t, fixture.root, "rev-parse", "HEAD") != working {
		t.Fatal("bootstrap refusal mutated run, candidate, or working branch")
	}
	if got := git(t, fixture.root, "rev-parse", "refs/bench/green/"+before.Branch); got != before.Base {
		t.Fatalf("green marker = %s after bootstrap refusal, want %s", got, before.Base)
	}
}

func TestPromoteCompletesAfterReviewedWorkingAdvanceRecomposition(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	before := loadRun(t, fixture.service)
	advanceWorking(t, fixture.root)
	owner := &promotionGate{accept: true}
	fixture.service.gate = owner
	if status, err := fixture.service.Promote(t.Context(), "build demo"); err != nil || status.Next != "bench spec build review build demo" {
		t.Fatalf("recompose Promote = %#v, %v", status, err)
	}
	run := loadRun(t, fixture.service)
	for key, assigned := range before.Assignments {
		if assigned.Integrated != "" {
			assigned.Integrated = run.CandidateTip
		}
		if got, ok := run.Assignments[key]; !ok || !reflect.DeepEqual(got, assigned) {
			t.Fatalf("recomposed assignment provenance = %#v, want %#v", got, assigned)
		}
	}
	receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatalf("fresh Review: %v", err)
	}
	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.State != "terminal" {
		t.Fatalf("second Promote = %#v, %v", status, err)
	}
	promoted := git(t, fixture.root, "rev-parse", "HEAD")
	if got := git(t, fixture.root, "rev-parse", "refs/bench/green/"+run.Branch); got != promoted {
		t.Fatalf("green marker = %s, want promotion commit %s", got, promoted)
	}
}

func (g *promotionGate) Bootstrap(ctx context.Context, root, branch, tip, expected string) error {
	return greenGate{}.Bootstrap(ctx, root, branch, tip, expected)
}

type recompositionGate struct {
	calls                 int
	branch, tip, expected string
	err                   error
}

func (g *recompositionGate) Bootstrap(ctx context.Context, root, branch, tip, expected string) error {
	g.calls++
	g.branch, g.tip, g.expected = branch, tip, expected
	if g.err != nil {
		return g.err
	}
	return greenGate{}.Bootstrap(ctx, root, branch, tip, expected)
}
