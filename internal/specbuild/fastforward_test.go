package specbuild

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type markerRaceGate struct {
	AuthorizationGate
	beforeAdvance func()
}

func (g markerRaceGate) AdvanceMarker(ctx context.Context, root, branch, destination, expected string) error {
	g.beforeAdvance()
	return g.AuthorizationGate.AdvanceMarker(ctx, root, branch, destination, expected)
}

func realPromotionFixture(t *testing.T) checkpointFixture {
	t.Helper()
	return assignRealPromotion(t, realPromotionStartFixture(t))
}

func realPromotionStartFixture(t *testing.T) checkpointFixture {
	t.Helper()
	root := specEditStartFixture(t)
	service := New(root, realPromotionGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	return checkpointFixture{root: root, service: service, run: loadRun(t, service)}
}

func assignRealPromotion(t *testing.T, fixture checkpointFixture) checkpointFixture {
	t.Helper()
	service := fixture.service
	assigned, _, err := service.Assign(t.Context(), "build demo", "one.md", "real promotion")
	if err != nil {
		t.Fatalf("Assign: %v", err)
	}
	write(t, filepath.Join(assigned.Path, "internal", "specbuild", "checkpoint-change.go"), "package specbuild\n")
	return checkpointFixture{
		root: fixture.root, service: service, assigned: assigned, run: loadRun(t, service),
		receipt: checkpointReceiptFor(t, service, assigned, []string{"internal/specbuild/checkpoint-change.go"}),
	}
}

func reviewRealPromotion(t *testing.T, fixture checkpointFixture) record {
	t.Helper()
	if err := checkpointInvocation(t, fixture); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	if err := integrateInvocation(t, fixture); err != nil {
		t.Fatalf("Integrate: %v", err)
	}
	run := loadRun(t, fixture.service)
	receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatalf("Review: %v", err)
	}
	return loadRun(t, fixture.service)
}

func TestFastForwardedRunPublishesThroughAuthorizationOwner(t *testing.T) {
	fixture := realPromotionStartFixture(t)
	marker := greenMarkerRef(t, fixture.root)
	commitAdvance(t, fixture.root, "later-advance.txt")
	recordGreenVerdict(t, fixture.root)
	if got := git(t, fixture.root, "rev-parse", marker); got != fixture.run.Base {
		t.Fatalf("ordinary gate moved marker to %s, want start base %s", got, fixture.run.Base)
	}
	if status, err := fixture.service.Start(t.Context(), "build demo"); err != nil || status.State != "active" {
		t.Fatalf("Start re-entry = %#v, %v", status, err)
	}
	fixture = assignRealPromotion(t, fixture)
	reviewRealPromotion(t, fixture)
	status, err := fixture.service.Promote(t.Context(), "build demo")
	final := loadRun(t, fixture.service)
	if err != nil || status.State != "terminal" || !refAt(fixture.root, "refs/heads/"+final.Branch, final.PromotionCommit) || !refAt(fixture.root, marker, final.PromotionCommit) {
		t.Fatalf("Promote = %#v, %v; run=%#v", status, err, final)
	}
	if reloaded, err := New(fixture.root, realPromotionGate{}, realOwner{}).Status("build demo"); err != nil || reloaded.State != "terminal" {
		t.Fatalf("fresh status = %#v, %v", reloaded, err)
	}
}

func TestRecoveryPublishesLaggingMarkerThroughAuthorizationOwner(t *testing.T) {
	fixture := realPromotionFixture(t)
	advanceWorking(t, fixture.root)
	recordGreenVerdict(t, fixture.root)
	reviewRealPromotion(t, fixture)
	fixture.service.fault = injectFault("promote/branch")
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("Promote completed across promote/branch fault")
	}
	fresh := New(fixture.root, realPromotionGate{}, realOwner{})
	status, err := fresh.Promote(t.Context(), "build demo")
	final := loadRun(t, fresh)
	if err != nil || status.State != "terminal" || !refAt(fixture.root, greenMarkerRef(t, fixture.root), final.PromotionCommit) {
		t.Fatalf("recovered Promote = %#v, %v; run=%#v", status, err, final)
	}
}

func TestRecoveryRefusesDivergentMarkerAndCompletesAfterRestore(t *testing.T) {
	fixture := realPromotionFixture(t)
	advanceWorking(t, fixture.root)
	recordGreenVerdict(t, fixture.root)
	reviewRealPromotion(t, fixture)
	marker := greenMarkerRef(t, fixture.root)
	lagging := git(t, fixture.root, "rev-parse", marker)
	fixture.service.fault = injectFault("promote/branch")
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("Promote completed across promote/branch fault")
	}
	prepared := loadRun(t, fixture.service)
	if prepared.Terminal || !refAt(fixture.root, "refs/heads/"+prepared.Branch, prepared.PromotionCommit) {
		t.Fatalf("branch interruption did not retain a prepared nonterminal run: %#v", prepared)
	}
	divergent := git(t, fixture.root, "commit-tree", git(t, fixture.root, "rev-parse", prepared.Base+"^{tree}"), "-p", fixture.run.Base, "-m", "divergent recovery marker")
	git(t, fixture.root, "update-ref", marker, divergent, lagging)
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	stateBefore, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	candidate := git(t, fixture.root, "rev-parse", prepared.Candidate)
	fresh := New(fixture.root, realPromotionGate{}, realOwner{})
	if _, err := fresh.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("recovered Promote accepted a divergent marker")
	}
	stateAfter, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	refused := loadRun(t, fresh)
	if refused.Terminal || git(t, fixture.root, "rev-parse", marker) != divergent || !refAt(fixture.root, "refs/heads/"+prepared.Branch, prepared.PromotionCommit) || git(t, fixture.root, "rev-parse", prepared.Candidate) != candidate || string(stateAfter) != string(stateBefore) {
		t.Fatalf("divergent-marker recovery refusal mutated prepared publication state: %#v", refused)
	}
	git(t, fixture.root, "update-ref", marker, lagging, divergent)
	reloaded := New(fixture.root, realPromotionGate{}, realOwner{})
	status, err := reloaded.Promote(t.Context(), "build demo")
	final := loadRun(t, reloaded)
	if err != nil || status.State != "terminal" || !refAt(fixture.root, "refs/heads/"+final.Branch, final.PromotionCommit) || !refAt(fixture.root, marker, final.PromotionCommit) {
		t.Fatalf("Promote after restoring lagging marker = %#v, %v; run=%#v", status, err, final)
	}
}

func TestCheckpointedFastForwardedRunRecomposesThroughAuthorizationOwner(t *testing.T) {
	fixture := realPromotionFixture(t)
	advanceWorking(t, fixture.root)
	recordGreenVerdict(t, fixture.root)
	reviewRealPromotion(t, fixture)
	commitAdvance(t, fixture.root, "later-advance.txt")
	recordGreenVerdict(t, fixture.root)
	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.Next != "bench spec build review build demo" {
		t.Fatalf("recomposing Promote = %#v, %v", status, err)
	}
	run := loadRun(t, fixture.service)
	receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatalf("fresh Review: %v", err)
	}
	if status, err := fixture.service.Promote(t.Context(), "build demo"); err != nil || status.State != "terminal" {
		t.Fatalf("Promote after recomposition = %#v, %v", status, err)
	}
}

func TestPromotionResumesAfterCommitFaultWithMarkerAtBase(t *testing.T) {
	fixture := realPromotionFixture(t)
	advanceWorking(t, fixture.root)
	recordGreenVerdict(t, fixture.root)
	reviewRealPromotion(t, fixture)
	base := git(t, fixture.root, "rev-parse", "HEAD")
	marker := greenMarkerRef(t, fixture.root)
	markerBase := git(t, fixture.root, "rev-parse", marker)
	fixture.service.fault = injectFault("promote/commit")
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("Promote completed across promote/commit fault")
	}
	if branch, marked := git(t, fixture.root, "rev-parse", "refs/heads/main"), git(t, fixture.root, "rev-parse", marker); branch != base || marked != markerBase {
		t.Fatalf("commit interruption moved publication refs: branch=%s marker=%s want=%s/%s", branch, marked, base, markerBase)
	}
	status, err := New(fixture.root, realPromotionGate{}, realOwner{}).Promote(t.Context(), "build demo")
	if err != nil || status.State != "terminal" {
		t.Fatalf("Promote after commit fault = %#v, %v", status, err)
	}
}

func TestPromotionRefusesDivergentMarkerAndRecoversAfterRestore(t *testing.T) {
	fixture := realPromotionFixture(t)
	advanceWorking(t, fixture.root)
	recordGreenVerdict(t, fixture.root)
	reviewRealPromotion(t, fixture)
	run := loadRun(t, fixture.service)
	marker := greenMarkerRef(t, fixture.root)
	divergent := git(t, fixture.root, "commit-tree", git(t, fixture.root, "rev-parse", run.Base+"^{tree}"), "-p", fixture.run.Base, "-m", "divergent marker")
	git(t, fixture.root, "update-ref", marker, divergent)
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	stateBefore, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	working := git(t, fixture.root, "rev-parse", "refs/heads/"+run.Branch)
	candidate := git(t, fixture.root, "rev-parse", run.Candidate)
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("Promote accepted a divergent marker")
	}
	refused := loadRun(t, fixture.service)
	stateAfter, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if refused.Terminal || git(t, fixture.root, "rev-parse", marker) != divergent || git(t, fixture.root, "rev-parse", "refs/heads/"+run.Branch) != working || git(t, fixture.root, "rev-parse", run.Candidate) != candidate || string(stateAfter) != string(stateBefore) {
		t.Fatalf("divergent-marker refusal did not preserve recovery state: %#v", refused)
	}
	git(t, fixture.root, "update-ref", marker, run.Base, divergent)
	status, err := New(fixture.root, realPromotionGate{}, realOwner{}).Promote(t.Context(), "build demo")
	if err != nil || status.State != "terminal" {
		t.Fatalf("Promote after restoring marker = %#v, %v", status, err)
	}
}

func TestPromotionRestoresStateWhenMarkerChangesAfterPreflight(t *testing.T) {
	fixture := realPromotionFixture(t)
	advanceWorking(t, fixture.root)
	recordGreenVerdict(t, fixture.root)
	reviewRealPromotion(t, fixture)
	run := loadRun(t, fixture.service)
	marker := greenMarkerRef(t, fixture.root)
	divergent := git(t, fixture.root, "commit-tree", git(t, fixture.root, "rev-parse", run.Base+"^{tree}"), "-p", fixture.run.Base, "-m", "racing divergent marker")
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	stateBefore, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	working := git(t, fixture.root, "rev-parse", "refs/heads/"+run.Branch)
	candidate := git(t, fixture.root, "rev-parse", run.Candidate)
	fixture.service.gate = markerRaceGate{beforeAdvance: func() { git(t, fixture.root, "update-ref", marker, divergent) }}
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("Promote accepted a marker changed after preflight")
	}
	stateAfter, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if git(t, fixture.root, "rev-parse", marker) != divergent || git(t, fixture.root, "rev-parse", "refs/heads/"+run.Branch) != working || git(t, fixture.root, "rev-parse", run.Candidate) != candidate || string(stateAfter) != string(stateBefore) || git(t, fixture.root, "status", "--porcelain", "--untracked-files=all") != "" {
		t.Fatal("racing marker refusal mutated publication state")
	}
}

// advancedEmptyRunFixture is the exact state the fast-forward owns: a run whose only
// assignment has never checkpointed, whose candidate still sits on the recorded base,
// and whose working branch has since advanced by a recognized descendant commit.
func advancedEmptyRunFixture(t *testing.T) (checkpointFixture, string) {
	t.Helper()
	fixture := newCheckpointFixture(t)
	advanceWorking(t, fixture.root)
	return fixture, git(t, fixture.root, "rev-parse", "HEAD")
}

func TestCheckpointFastForwardsEmptyRunOntoAdvancedTip(t *testing.T) {
	fixture, tip := advancedEmptyRunFixture(t)
	if err := checkpointInvocation(t, fixture); err != nil {
		t.Fatalf("Checkpoint on an advanced empty run = %v, want the fast-forward to let it proceed", err)
	}
	run := loadRun(t, fixture.service)
	if run.Base != tip || run.CandidateTip != tip || !refAt(fixture.root, run.Candidate, tip) {
		t.Fatalf("fast-forward left base=%s candidate tip=%s ref=%s, want %s", run.Base, run.CandidateTip, git(t, fixture.root, "rev-parse", run.Candidate), tip)
	}
	_, assigned, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok || assigned.Checkpoint == "" {
		t.Fatalf("checkpoint did not proceed after the fast-forward: %#v", assigned)
	}
}

func TestStartFastForwardsEmptyRunAndReportsStatus(t *testing.T) {
	fixture := newPreconditionFixture(t, true)
	advanceWorking(t, fixture.root)
	tip := git(t, fixture.root, "rev-parse", "HEAD")
	status, err := fixture.service.Start(t.Context(), "build demo")
	if err != nil || status.State != "active" {
		t.Fatalf("Start on an advanced empty run = %#v, %v; want the run's status", status, err)
	}
	run := loadRun(t, fixture.service)
	if run.Base != tip || run.CandidateTip != tip || !refAt(fixture.root, run.Candidate, tip) {
		t.Fatalf("fast-forward left base=%s candidate tip=%s ref=%s, want %s", run.Base, run.CandidateTip, git(t, fixture.root, "rev-parse", run.Candidate), tip)
	}
}

func TestFastForwardedRunReloadsInFreshService(t *testing.T) {
	fixture, tip := advancedEmptyRunFixture(t)
	before := loadRun(t, fixture.service)
	if err := checkpointInvocation(t, fixture); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	fresh := New(fixture.root, fixture.gate, realOwner{})
	after, found, err := fresh.load("build demo")
	if err != nil || !found {
		t.Fatalf("fresh load: found=%t err=%v", found, err)
	}
	if after.Run != before.Run || after.Candidate != before.Candidate || after.Base != tip || after.CandidateTip != tip || !refAt(fixture.root, after.Candidate, after.CandidateTip) {
		t.Fatalf("reloaded identity before=%#v after=%#v want tip %s", before, after, tip)
	}
}

// The durable candidate ref is the run's identity, so a fast-forward that overwrote a
// ref another writer had already moved would silently discard that writer's commit.
func TestFastForwardRefusesCandidateRefMovedExternally(t *testing.T) {
	fixture := newPreconditionFixture(t, true)
	advanceWorking(t, fixture.root)
	moveCandidate(t, &fixture)
	before := snapshotPrecondition(t, fixture)
	if _, err := fixture.service.Start(t.Context(), "build demo"); err == nil {
		t.Fatal("Start fast-forwarded onto a candidate ref that moved externally")
	}
	if after := snapshotPrecondition(t, fixture); after != before {
		t.Fatalf("drifted candidate mutated:\n before=%#v\n after=%#v", before, after)
	}
	run := loadRun(t, fixture.service)
	if err := fixture.service.fastForwardEmptyRun(&run, git(t, fixture.root, "rev-parse", "HEAD")); err == nil {
		t.Fatal("fastForwardEmptyRun overwrote a candidate ref that moved externally")
	}
	if after := snapshotPrecondition(t, fixture); after != before {
		t.Fatalf("refused fast-forward mutated:\n before=%#v\n after=%#v", before, after)
	}
}

// A run holding real checkpoint evidence has work a recomposition must replay, so it
// keeps the refusal that routes it to promote no matter how its tip moved.
func TestCheckpointedRunWithMovedTipStillRoutesToPromote(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R17] second assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	checkpointed, _, err := service.Assign(t.Context(), "build demo", "one.md", "checkpointed sibling")
	if err != nil {
		t.Fatalf("Assign checkpointed: %v", err)
	}
	pending, _, err := service.Assign(t.Context(), "build demo", "two.md", "pending sibling")
	if err != nil {
		t.Fatalf("Assign pending: %v", err)
	}
	write(t, filepath.Join(checkpointed.Path, "internal", "specbuild", "checkpointed-change.go"), "package specbuild\n")
	before := checkpointAssignment(t, root, service, checkpointed, []string{"internal/specbuild/checkpointed-change.go"})
	write(t, filepath.Join(pending.Path, "internal", "specbuild", "pending-change.go"), "package specbuild\n")
	rec := checkpointReceiptFor(t, service, pending, []string{"internal/specbuild/pending-change.go"})
	advanceWorking(t, root)
	if _, err := service.Checkpoint(t.Context(), "build demo", pending.ID, writeCheckpointReceipt(t, rec, "\n")); err == nil || !strings.Contains(err.Error(), "bench spec build promote") {
		t.Fatalf("Checkpoint on a checkpointed run with a moved tip = %v, want the recomposition refusal naming promote", err)
	}
	if after := loadRun(t, service); after.Base != before.Base || after.CandidateTip != before.CandidateTip {
		t.Fatalf("refused checkpoint moved the run: before base=%s tip=%s, after base=%s tip=%s", before.Base, before.CandidateTip, after.Base, after.CandidateTip)
	}
}

// Review validates its receipt only after the preconditions, so a fast-forward reachable
// from review would move base and candidate for a receipt that then fails.
func TestReviewOnAdvancedEmptyRunRefusesWithoutMutation(t *testing.T) {
	fixture := newPreconditionFixture(t, true)
	run := loadRun(t, fixture.service)
	evidence := writeReviewReceipt(t, reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}})
	advanceWorking(t, fixture.root)
	before := snapshotPrecondition(t, fixture)
	if _, err := fixture.service.Review(t.Context(), "build demo", evidence); err == nil || !strings.Contains(err.Error(), "bench spec build promote") {
		t.Fatalf("Review on an advanced empty run = %v, want the recomposition refusal naming promote", err)
	}
	if after := snapshotPrecondition(t, fixture); after != before {
		t.Fatalf("review refusal mutated:\n before=%#v\n after=%#v", before, after)
	}
}

func TestNonAncestorTipRefusesFastForwardOnEmptyRun(t *testing.T) {
	fixture := newCheckpointFixture(t)
	rewriteWorkingHead(t, fixture.root)
	before := checkpointSnapshotFor(t, fixture)
	err := checkpointInvocation(t, fixture)
	if err == nil || !strings.Contains(err.Error(), "does not match recorded subject") || strings.Contains(err.Error(), "bench spec build promote") {
		t.Fatalf("Checkpoint on a rewritten head = %v, want the subject-mismatch refusal", err)
	}
	if after := checkpointSnapshotFor(t, fixture); after != before {
		t.Fatalf("rewritten-head refusal mutated: before=%#v after=%#v", before, after)
	}
}

// Every pre-promote transition is provisional, so the fast-forward must not spend the
// gate that only promote is allowed to consult.
func TestFastForwardRunsNoGate(t *testing.T) {
	fixture, _ := advancedEmptyRunFixture(t)
	before := fixture.gate.calls
	marker := git(t, fixture.root, "rev-parse", greenMarkerRef(t, fixture.root))
	if err := checkpointInvocation(t, fixture); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	if fixture.gate.calls != before {
		t.Fatalf("fast-forward ran the gate %d extra time(s)", fixture.gate.calls-before)
	}
	if got := git(t, fixture.root, "rev-parse", greenMarkerRef(t, fixture.root)); got != marker {
		t.Fatalf("fast-forward moved marker to %s, want %s", got, marker)
	}
}

func TestFastForwardIsIdempotentAcrossLaterMutations(t *testing.T) {
	fixture := newCheckpointFixture(t)
	git(t, fixture.root, "config", "core.logAllRefUpdates", "always")
	advanceWorking(t, fixture.root)
	if err := checkpointInvocation(t, fixture); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	moves := git(t, fixture.root, "reflog", "show", fixture.run.Candidate)
	before := checkpointSnapshotFor(t, fixture)
	if _, err := fixture.service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("Start after the fast-forward: %v", err)
	}
	if got := git(t, fixture.root, "reflog", "show", fixture.run.Candidate); got != moves || len(stringsSplitLines(got)) != 1 {
		t.Fatalf("candidate ref moves = %q, want the single fast-forward %q", got, moves)
	}
	if after := checkpointSnapshotFor(t, fixture); after != before {
		t.Fatalf("second mutation rewrote state:\n before=%#v\n after=%#v", before, after)
	}
}
