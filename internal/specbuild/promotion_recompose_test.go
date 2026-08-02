package specbuild

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

const samePathFixture = `package specbuild

var candidateValue = "base"

var spacer01 = 1
var spacer02 = 2
var spacer03 = 3
var spacer04 = 4
var spacer05 = 5
var spacer06 = 6
var spacer07 = 7
var spacer08 = 8

var workingValue = "base"
`

func applyCheckpointFixtureConfiguration(root string, configure []func(string)) {
	for _, apply := range configure {
		apply(root)
	}
}

func TestPromoteRecomposesCompatibleSamePathChanges(t *testing.T) {
	candidateContent := strings.Replace(samePathFixture, `candidateValue = "base"`, `candidateValue = "candidate"`, 1)
	fixture := reviewedSamePathPromotionFixture(t, candidateContent)
	before := loadRun(t, fixture.service)
	owner := &recompositionGate{}
	fixture.service.gate = owner
	workingContent := strings.Replace(samePathFixture, `workingValue = "base"`, `workingValue = "working"`, 1)
	write(t, filepath.Join(fixture.root, "internal", "specbuild", "same-path.go"), workingContent)
	git(t, fixture.root, "add", ".")
	git(t, fixture.root, "commit", "-qm", "compatible same-path advance")
	working := git(t, fixture.root, "rev-parse", "HEAD")

	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.Next != "bench spec build review build demo" {
		t.Fatalf("Promote = %#v, %v", status, err)
	}
	after := loadRun(t, fixture.service)
	content := git(t, fixture.root, "show", after.CandidateTip+":internal/specbuild/same-path.go")
	if !strings.Contains(content, `candidateValue = "candidate"`) || !strings.Contains(content, `workingValue = "working"`) {
		t.Fatalf("recomposed same-path content = %q", content)
	}
	if after.Base != working || after.Review != nil || owner.branch != before.Branch || owner.tip != working || owner.expected != before.Base || git(t, fixture.root, "rev-parse", "refs/bench/green/"+before.Branch) != working {
		t.Fatalf("recomposed state = %#v, bootstrap=%#v", after, owner)
	}
}

func TestPromoteRefusesConflictingSamePathChangesWithoutMutation(t *testing.T) {
	candidateContent := strings.Replace(samePathFixture, `candidateValue = "base"`, `candidateValue = "candidate"`, 1)
	fixture := reviewedSamePathPromotionFixture(t, candidateContent)
	workingContent := strings.Replace(samePathFixture, `candidateValue = "base"`, `candidateValue = "working"`, 1)
	write(t, filepath.Join(fixture.root, "internal", "specbuild", "same-path.go"), workingContent)
	git(t, fixture.root, "add", ".")
	git(t, fixture.root, "commit", "-qm", "conflicting same-path advance")
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
	owner := &recompositionGate{}
	fixture.service.gate = owner

	_, promoteErr := fixture.service.Promote(t.Context(), "build demo")
	if promoteErr == nil || !strings.Contains(promoteErr.Error(), "checkpoint patch conflicts with the candidate") {
		t.Fatalf("Promote error = %v", promoteErr)
	}
	afterState, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(afterState) != string(beforeState) || git(t, fixture.root, "rev-parse", before.Candidate) != candidate || git(t, fixture.root, "rev-parse", "HEAD") != working || owner.calls != 0 {
		t.Fatalf("conflict mutated protected state or reached bootstrap: calls=%d", owner.calls)
	}
	if got := git(t, fixture.root, "rev-parse", "refs/bench/green/"+before.Branch); got != before.Base {
		t.Fatalf("green marker = %s after conflict, want %s", got, before.Base)
	}
}

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

func TestRecomposedAttemptReloadsInFreshService(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	abandon := &abandonOwner{}
	fixture.service.worktrees = abandon
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	if _, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err != nil {
		t.Fatalf("ApplyAbandon: %v", err)
	}
	fixture.service.worktrees = realOwner{}
	advanceWorking(t, fixture.root)
	if _, err := fixture.service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("restart: %v", err)
	}
	assigned, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "attempt checkpoint")
	if err != nil {
		t.Fatalf("Assign: %v", err)
	}
	write(t, filepath.Join(assigned.Path, "internal", "specbuild", "attempt-change.go"), "package specbuild\n")
	checkpointAssignment(t, fixture.root, fixture.service, assigned, []string{"internal/specbuild/attempt-change.go"})
	if _, err := fixture.service.Integrate(t.Context(), "build demo", assigned.ID); err != nil {
		t.Fatalf("Integrate: %v", err)
	}
	before := loadRun(t, fixture.service)
	review := reviewReceipt{Version: 1, Run: before.Run, Candidate: before.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, review)); err != nil {
		t.Fatalf("Review: %v", err)
	}
	write(t, filepath.Join(fixture.root, "second-advance.txt"), "advance\n")
	git(t, fixture.root, "add", ".")
	git(t, fixture.root, "commit", "-qm", "second advance")
	owner := &recompositionGate{}
	fixture.service.gate = owner
	recomposed, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || recomposed.Next != "bench spec build review build demo" {
		t.Fatalf("Promote = %#v, %v", recomposed, err)
	}

	fresh := New(fixture.root, owner, realOwner{})
	status, err := fresh.Status("build demo")
	if err != nil || status.Subject != recomposed.Subject || status.Next != "bench spec build review build demo" {
		t.Fatalf("fresh Status = %#v, %v", status, err)
	}
	after, found, err := fresh.load("build demo")
	if err != nil || !found {
		t.Fatalf("fresh load: found=%t err=%v", found, err)
	}
	if after.Run != before.Run || after.Candidate != before.Candidate || after.Base == before.Base || after.CandidateTip == before.CandidateTip || !refAt(fixture.root, before.Candidate, after.CandidateTip) {
		t.Fatalf("reloaded identity before=%#v after=%#v", before, after)
	}
	review.Run, review.Candidate = after.Run, after.CandidateTip
	if _, err := fresh.Review(t.Context(), "build demo", writeReviewReceipt(t, review)); err != nil {
		t.Fatalf("fresh Review: %v", err)
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

func TestPromoteRecomposesMidRepairRun(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	run := loadRun(t, fixture.service)
	run.Review.Axes[0].Findings = []reviewFinding{{ID: "F1", Disposition: "accepted"}}
	saveRun(t, fixture.service, run)
	repair, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "repair finding F1")
	if err != nil {
		t.Fatalf("Assign repair: %v", err)
	}
	advanceWorking(t, fixture.root)
	working := git(t, fixture.root, "rev-parse", "HEAD")

	if _, err := fixture.service.Promote(t.Context(), "build demo"); err != nil {
		t.Fatalf("mid-repair Promote: %v", err)
	}
	after := loadRun(t, fixture.service)
	_, assigned, ok := assignmentFor(after, repair.ID)
	if !ok || assigned.Released || after.Base != working || after.Review != nil {
		t.Fatalf("mid-repair recomposition = %#v, repair assignment %#v", after, assigned)
	}
}

func TestPromoteTwiceOnMovedTipRecomposesOnce(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	advanceWorking(t, fixture.root)
	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.Next != "bench spec build review build demo" {
		t.Fatalf("first Promote = %#v, %v", status, err)
	}
	before := promotionSnapshotFor(t, fixture)

	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil || !strings.Contains(err.Error(), "requires a current clean review") {
		t.Fatalf("second Promote = %v, want the clean-review refusal", err)
	}
	if after := promotionSnapshotFor(t, fixture); after != before {
		t.Fatalf("second Promote replayed: before=%#v after=%#v", before, after)
	}
	if next, err := fixture.service.Status("build demo"); err != nil || next.Next != "bench spec build review build demo" {
		t.Fatalf("recomposed next = %#v, %v", next, err)
	}
}

func TestUnrecognizedHeadMoveStillRefusesPromote(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	rewriteWorkingHead(t, fixture.root)
	requirePromotionRefusal(t, fixture, "does not match recorded subject")
}

func TestPromoteStillRefusesUnreadyRunOnUnmovedTip(t *testing.T) {
	for _, test := range []struct {
		name, want string
		unready    func(*testing.T, checkpointFixture)
	}{
		{"absent review", "requires a current clean review", func(t *testing.T, fixture checkpointFixture) {
			updatePromotionRun(t, fixture, func(run *record) { run.Review = nil })
		}},
		{"stale candidate review", "requires a current clean review", func(t *testing.T, fixture checkpointFixture) {
			updatePromotionRun(t, fixture, func(run *record) { run.Review.Candidate = run.Base })
		}},
		{"accepted finding", "requires a current clean review", func(t *testing.T, fixture checkpointFixture) {
			updatePromotionRun(t, fixture, func(run *record) {
				run.Review.Axes[0].Findings = []reviewFinding{{ID: "F1", Disposition: "accepted"}}
			})
		}},
		{"unreleased assignment", "requires every assignment integrated and released", func(t *testing.T, fixture checkpointFixture) {
			if _, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "second sibling"); err != nil {
				t.Fatalf("Assign second: %v", err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := reviewedPromotionFixture(t)
			test.unready(t, fixture)
			requirePromotionRefusal(t, fixture, test.want)
		})
	}
}

func TestPromoteStillRefusesEachEvidenceFault(t *testing.T) {
	for _, test := range []struct {
		name, want string
		corrupt    func(*testing.T, checkpointFixture)
	}{
		{"drifted candidate ref", "spec build candidate no longer matches durable tip", func(t *testing.T, fixture checkpointFixture) {
			run := loadRun(t, fixture.service)
			driftRef(t, fixture.root, run.Candidate, run.CandidateTip)
		}},
		{"incomplete checkpoint fields", "retained checkpoint evidence is incomplete", func(t *testing.T, fixture checkpointFixture) {
			updatePromotionRun(t, fixture, func(run *record) {
				for key, assigned := range run.Assignments {
					assigned.ReceiptDigest = ""
					run.Assignments[key] = assigned
				}
			})
		}},
		{"drifted checkpoint ref", "retained checkpoint reference drifted", func(t *testing.T, fixture checkpointFixture) {
			run := loadRun(t, fixture.service)
			for _, assigned := range run.Assignments {
				driftRef(t, fixture.root, assigned.CheckpointRef, assigned.Checkpoint)
			}
		}},
		{"integration outside candidate ancestry", "retained integration left candidate ancestry", func(t *testing.T, fixture checkpointFixture) {
			unrelated := git(t, fixture.root, "commit-tree", git(t, fixture.root, "rev-parse", "HEAD^{tree}"), "-m", "unrelated root")
			updatePromotionRun(t, fixture, func(run *record) {
				for key, assigned := range run.Assignments {
					assigned.Integrated = unrelated
					run.Assignments[key] = assigned
				}
			})
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := reviewedPromotionFixture(t)
			test.corrupt(t, fixture)
			requirePromotionRefusal(t, fixture, test.want)
		})
	}
}

type promotionSnapshot struct{ state, refs, head string }

func promotionSnapshotFor(t *testing.T, fixture checkpointFixture) promotionSnapshot {
	t.Helper()
	path, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	state, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return promotionSnapshot{state: string(state), refs: git(t, fixture.root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/"), head: git(t, fixture.root, "rev-parse", "HEAD")}
}

func requirePromotionRefusal(t *testing.T, fixture checkpointFixture, want string) {
	t.Helper()
	owner := &promotionGate{accept: true}
	fixture.service.gate = owner
	before := promotionSnapshotFor(t, fixture)
	_, err := fixture.service.Promote(t.Context(), "build demo")
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("Promote = %v, want a refusal naming %q", err, want)
	}
	if owner.executions != 0 {
		t.Fatalf("refusal reached the prospective gate: %d executions", owner.executions)
	}
	if after := promotionSnapshotFor(t, fixture); after != before {
		t.Fatalf("refusal mutated state: before=%#v after=%#v", before, after)
	}
}

func updatePromotionRun(t *testing.T, fixture checkpointFixture, change func(*record)) {
	t.Helper()
	run := loadRun(t, fixture.service)
	change(&run)
	saveRun(t, fixture.service, run)
}

func driftRef(t *testing.T, root, name, current string) {
	t.Helper()
	commit := git(t, root, "commit-tree", git(t, root, "rev-parse", current+"^{tree}"), "-p", current, "-m", "ref drift")
	git(t, root, "update-ref", name, commit, current)
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

func reviewedSamePathPromotionFixture(t *testing.T, candidateContent string) checkpointFixture {
	t.Helper()
	fixture := newCheckpointFixture(t, func(root string) {
		write(t, filepath.Join(root, "internal", "specbuild", "same-path.go"), samePathFixture)
	})
	path := filepath.Join(fixture.assigned.Path, "internal", "specbuild", "same-path.go")
	write(t, path, candidateContent)
	tree := benchgit.TreeHash(fixture.assigned.Path)
	fixture.receipt.Tree, fixture.receipt.Probe.Tree = tree, tree
	fixture.receipt.Ownership = append(fixture.receipt.Ownership, "internal/specbuild/same-path.go")
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate: %v", err)
	}
	run := loadRun(t, fixture.service)
	receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatalf("Review: %v", err)
	}
	return fixture
}
