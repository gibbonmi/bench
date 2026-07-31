package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/worktree"
)

func TestAbandonPlansWithoutMutationAndRecoversDirtyAssignment(t *testing.T) {
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	if _, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n")); err != nil {
		t.Fatalf("Checkpoint: %v", err)
	}
	git(t, fixture.root, "update-ref", "refs/bench/recovery/preexisting", "HEAD")
	dirty := filepath.Join(fixture.assigned.Path, "internal", "specbuild", "checkpoint-change.go")
	write(t, dirty, "package specbuild\n\nvar abandoned = true\n")
	before := abandonmentSnapshot(t, fixture)
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	if len(plan.Worktrees) != 1 || plan.Worktrees[0].ID != fixture.assigned.ID {
		t.Fatalf("planned worktrees = %#v", plan.Worktrees)
	}
	if len(plan.ProvisionalRefs) != 2 || len(plan.UnintegratedCheckpoints) != 1 || len(plan.RecoveryRefs) != 1 {
		t.Fatalf("abandonment inventory = %#v", plan)
	}
	if plan.Fingerprint == "" {
		t.Fatal("Abandon returned no fingerprint")
	}
	if after := abandonmentSnapshot(t, fixture); after != before {
		t.Fatalf("read-only plan mutated state: before=%#v after=%#v", before, after)
	}
	status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil {
		t.Fatalf("ApplyAbandon: %v", err)
	}
	if status.State != "terminal" || owner.applies != 1 {
		t.Fatalf("apply status=%#v apply calls=%d", status, owner.applies)
	}
	if replay, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err != nil || replay != status || owner.applies != 1 {
		t.Fatalf("abandon replay status=%#v err=%v apply calls=%d", replay, err, owner.applies)
	}
	if _, err := os.Stat(fixture.assigned.Path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("owned worktree remains after abandonment: %v", err)
	}
	if refs := git(t, fixture.root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/"); refs == "" {
		t.Fatal("dirty assignment was released without a recovery ref")
	}
	run := loadRun(t, fixture.service)
	if !run.Terminal || len(run.Assignments) != 1 {
		t.Fatalf("terminal evidence = %#v", run)
	}
}
func TestApplyAbandonRefusesInventoryDriftBeforeRelease(t *testing.T) {
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	git(t, fixture.root, "update-ref", "refs/bench/specbuild/drift", "HEAD")
	before := abandonmentSnapshot(t, fixture)
	if _, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil || !strings.Contains(err.Error(), "plan drifted") {
		t.Fatalf("ApplyAbandon drift error = %v", err)
	}
	if owner.applies != 0 {
		t.Fatalf("apply calls after drift = %d, want 0", owner.applies)
	}
	if after := abandonmentSnapshot(t, fixture); after != before {
		t.Fatalf("drift refusal mutated state: before=%#v after=%#v", before, after)
	}
}
func TestApplyAbandonRefusesOwnerPlanDriftBeforeRelease(t *testing.T) {
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	before := git(t, fixture.assigned.Path, "status", "--porcelain", "--untracked-files=all")
	write(t, filepath.Join(fixture.assigned.Path, "dist", "owner-plan-output"), "ignored build output\n")
	if after := git(t, fixture.assigned.Path, "status", "--porcelain", "--untracked-files=all"); after != before {
		t.Fatalf("owner-only drift changed raw status: before=%q after=%q", before, after)
	}
	if _, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil || !strings.Contains(err.Error(), "plan drifted") {
		t.Fatalf("ApplyAbandon owner drift error = %v", err)
	}
	if owner.applies != 0 {
		t.Fatalf("owner drift reached apply: %d", owner.applies)
	}
}
func TestApplyAbandonResumesAfterFirstCleanup(t *testing.T) {
	root, service, owner, plan := twoAssignmentAbandonFixture(t, true)
	service.fault = injectFault("abandon/release")
	if _, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil {
		t.Fatal("ApplyAbandon completed through injected interruption")
	}
	service.fault = nil
	if owner.applies != 1 {
		t.Fatalf("apply calls before resume = %d", owner.applies)
	}
	git(t, root, "update-ref", "refs/bench/specbuild/unplanned-drift", "HEAD")
	if _, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil || !strings.Contains(err.Error(), "plan drifted") || owner.applies != 1 {
		t.Fatalf("drifted resume error=%v apply calls=%d", err, owner.applies)
	}
	git(t, root, "update-ref", "-d", "refs/bench/specbuild/unplanned-drift")
	if status, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err != nil || status.State != "terminal" || owner.applies != 2 {
		t.Fatalf("resume status=%#v err=%v apply calls=%d", status, err, owner.applies)
	}
	run := loadRun(t, service)
	if !run.Terminal {
		t.Fatalf("resumed terminal state=%#v", run)
	}
	if _, err := json.Marshal(run); err != nil {
		t.Fatalf("retained terminal evidence does not encode: %v", err)
	}
}
func TestApplyAbandonRecoversAfterOwnerApplyBeforeState(t *testing.T) {
	root, service, owner, plan := twoAssignmentAbandonFixture(t, true)
	service.fault = injectFault("abandon/owner-apply")
	if _, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil {
		t.Fatal("ApplyAbandon completed through owner-apply interruption")
	}
	if owner.applies != 1 {
		t.Fatalf("apply calls before recovery = %d", owner.applies)
	}
	service.fault = nil
	status, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil || status.State != "terminal" || owner.applies != 3 {
		t.Fatalf("owner-apply recovery status=%#v err=%v apply calls=%d", status, err, owner.applies)
	}
	if refs := strings.Fields(git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/recovery/")); len(refs) != 2 {
		t.Fatalf("recovery refs = %#v, want one per dirty assignment", refs)
	}
	run := loadRun(t, service)
	if !run.Terminal {
		t.Fatalf("owner-apply terminal state=%#v", run)
	}
	op, found := service.operation(run, "abandon", "apply")
	journal := abandonmentJournal{}
	if !found || op.State != "completed" || json.Unmarshal([]byte(op.Result), &journal) != nil || journal.Original.Fingerprint != plan.Fingerprint || len(journal.Current.Worktrees) != 0 {
		t.Fatalf("terminal abandonment journal = %#v, found=%t decoded=%#v", op, found, journal)
	}
}
func TestApplyAbandonResumesCleanReleasedAssignment(t *testing.T) {
	_, service, owner, plan := twoAssignmentAbandonFixture(t, false)
	service.fault = injectFault("abandon/release")
	if _, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil {
		t.Fatal("ApplyAbandon completed through injected interruption")
	}
	service.fault = nil
	if status, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err != nil || status.State != "terminal" || owner.applies != 2 {
		t.Fatalf("clean resume status=%#v err=%v apply calls=%d", status, err, owner.applies)
	}
}
func TestAbandonOwnerRejectsMismatchedRequest(t *testing.T) {
	fixture := newCheckpointFixture(t)
	if _, err := worktree.PlanAbandon(fixture.root, "other request", fixture.assigned.Path); err == nil {
		t.Fatal("PlanAbandon accepted a mismatched request")
	}
	if _, err := os.Stat(fixture.assigned.Path); err != nil {
		t.Fatalf("mismatched request removed assignment: %v", err)
	}
}

func twoAssignmentAbandonFixture(t *testing.T, dirty bool) (string, *Service, *abandonOwner, AbandonmentPlan) {
	t.Helper()
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R26] second abandonment assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second abandonment ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatal(err)
	}
	for _, ticket := range []string{"one.md", "two.md"} {
		assigned, _, err := service.Assign(t.Context(), "build demo", ticket, ticket)
		if err != nil {
			t.Fatal(err)
		}
		if dirty {
			write(t, filepath.Join(assigned.Path, "internal", "specbuild", "unlanded.go"), "package specbuild\n")
		}
	}
	owner := &abandonOwner{}
	service.worktrees = owner
	plan, err := service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	return root, service, owner, plan
}
func TestPromotePublishesOnlyAuthorizedProspectiveSquash(t *testing.T) {
	root, service, firstAssignment, secondAssignment, _ := siblingCheckpoints(t,
		"internal/specbuild/promotion-one.go", "package specbuild\nvar promotionOne = true\n",
		"internal/specbuild/promotion-two.go", "package specbuild\nvar promotionTwo = true\n")
	if _, err := service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("active assignment promoted")
	}
	for _, assigned := range []Assignment{firstAssignment, secondAssignment} {
		if _, err := service.Integrate(t.Context(), "build demo", assigned.ID); err != nil {
			t.Fatal(err)
		}
	}
	run := loadRun(t, service)
	if _, err := service.Review(t.Context(), "build demo", writeReviewReceipt(t, reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}})); err != nil {
		t.Fatal(err)
	}
	setAssignmentsReleased(t, service, false)
	if _, err := service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("unreleased assignment promoted")
	}
	setAssignmentsReleased(t, service, true)
	owner := &promotionGate{accept: true}
	service.gate = owner
	base, candidate := git(t, root, "rev-parse", "HEAD"), run.CandidateTip
	owner.inspect = func(root, tree string) {
		if got := git(t, root, "show", "HEAD:specs/build demo/spec.md"); !strings.Contains(got, "Status: staged") || !strings.Contains(git(t, root, "show", tree+":specs/build demo/spec.md"), "Status: implemented") {
			t.Fatal("gate did not receive an unpublished implemented tree")
		}
	}
	first, firstErr := service.prospectiveTree(t.Context(), run)
	second, secondErr := service.prospectiveTree(t.Context(), run)
	if firstErr != nil || secondErr != nil || first != second {
		t.Fatalf("prospective tree construction = %q/%v then %q/%v", first, firstErr, second, secondErr)
	}
	for disposition, next := range map[GateDisposition]string{GateCandidate: "delegate candidate gate repair", GateInherited: "diagnose inherited gate", GateInfrastructure: "retry promote", GateCapExhausted: "implementation cap exhausted"} {
		owner.red, owner.disposition = true, disposition
		status, err := service.Promote(t.Context(), "build demo")
		if err == nil || status.Next != next || git(t, root, "rev-parse", "HEAD") != base || git(t, root, "rev-parse", run.Candidate) != candidate {
			t.Fatalf("%s red did not preserve or route", disposition)
		}
	}
	owner.err = errors.New("gate unavailable")
	if _, err := service.Promote(t.Context(), "build demo"); err == nil || git(t, root, "rev-parse", "HEAD") != base || git(t, root, "rev-parse", run.Candidate) != candidate {
		t.Fatal("operational gate failure published state")
	}
	owner.err, owner.red, owner.disposition = nil, false, GateCandidate
	owner.contradictory = true
	if _, err := service.Promote(t.Context(), "build demo"); err == nil || git(t, root, "rev-parse", "HEAD") != base || git(t, root, "rev-parse", run.Candidate) != candidate {
		t.Fatal("contradictory green gate outcome published state")
	}
	owner.contradictory, service.fault = false, injectFault("promote/branch")
	if _, err := service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("branch-before-marker interruption did not stop")
	}
	prepared := loadRun(t, service)
	if !refAt(root, "refs/heads/"+prepared.Branch, prepared.PromotionCommit) {
		t.Fatalf("branch interruption did not retain a resumable commit: %#v", prepared)
	}
	name := git(t, root, "symbolic-ref", "--short", "HEAD")
	if git(t, root, "rev-parse", "refs/bench/green/"+name) != base {
		t.Fatal("marker advanced before recovery")
	}
	git(t, root, "update-ref", "-d", prepared.Candidate)
	if _, err := service.Promote(t.Context(), "build demo"); err == nil || git(t, root, "rev-parse", "refs/bench/green/"+name) != base || loadRun(t, service).Terminal {
		t.Fatal("missing retained candidate evidence advanced recovery")
	}
	git(t, root, "update-ref", prepared.Candidate, prepared.CandidateTip)
	write(t, filepath.Join(root, "tracked.txt"), "dirty recovery\n")
	if _, err := service.Promote(t.Context(), "build demo"); err == nil || git(t, root, "rev-parse", "refs/bench/green/"+name) != base || loadRun(t, service).Terminal {
		t.Fatal("tracked recovery dirt was overwritten")
	}
	write(t, filepath.Join(root, "tracked.txt"), "base\n")
	untracked := filepath.Join(root, "recovery-dirt.txt")
	write(t, untracked, "untracked recovery\n")
	if _, err := service.Promote(t.Context(), "build demo"); err == nil || git(t, root, "rev-parse", "refs/bench/green/"+name) != base || loadRun(t, service).Terminal {
		t.Fatal("untracked recovery dirt was overwritten")
	}
	_ = os.Remove(untracked)
	owner.accept = false
	if _, err := service.Promote(t.Context(), "build demo"); err == nil || git(t, root, "rev-parse", "refs/bench/green/"+name) != base {
		t.Fatal("invalid retained evidence advanced marker")
	}
	owner.accept = true
	prepared.PromotionEvidence = "swapped-proof"
	saveRun(t, service, prepared)
	if _, err := service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("swapped retained evidence advanced marker")
	}
	prepared.PromotionEvidence = "owner-proof"
	saveRun(t, service, prepared)
	service.fault = nil
	status, err := service.Promote(t.Context(), "build demo")
	if err != nil || status.State != "terminal" || owner.tree == "" || owner.validations != 5 {
		t.Fatalf("Promote = %#v, %v; gate tree=%q", status, err, owner.tree)
	}
	branch := git(t, root, "rev-parse", "HEAD")
	if git(t, root, "rev-parse", branch+"^") != base || git(t, root, "rev-parse", branch+"^{tree}") != owner.tree || git(t, root, "rev-parse", "refs/bench/green/"+name) != branch {
		t.Fatal("promotion did not publish its authorized squash before the green marker")
	}
	if strings.Contains("\n"+git(t, root, "rev-list", branch)+"\n", "\n"+candidate+"\n") {
		t.Fatal("provisional candidate entered working ancestry")
	}
	if dirty := git(t, root, "status", "--porcelain", "--untracked-files=all"); dirty != "" {
		t.Fatalf("promotion left the working checkout dirty: %q", dirty)
	}
	full, fullErr := service.FullStatus("build demo")
	if fullErr != nil || full.Review == nil || len(full.Assignments) != 2 || full.Assignments[0].Checkpoint == "" || full.Assignments[1].Checkpoint == "" || full.Assignments[0].Integrated == "" || full.Assignments[1].Integrated == "" || full.Assignments[0].Cleanup != "released" || full.Assignments[1].Cleanup != "released" {
		t.Fatalf("retained full status = %#v, %v", full, fullErr)
	}
	if got := git(t, root, "show", branch+":specs/build demo/spec.md"); !strings.Contains(got, "Status: implemented") {
		t.Fatalf("published spec = %q", got)
	}
	if got := git(t, root, "show", candidate+":specs/build demo/spec.md"); !strings.Contains(got, "Status: staged") {
		t.Fatalf("candidate spec = %q", got)
	}
	final := loadRun(t, service)
	op, found := service.operation(final, "promote", final.CandidateTip)
	executions, validations := owner.executions, owner.validations
	if !found || !final.Terminal || op.State != "completed" || op.Result != branch {
		t.Fatalf("terminal promotion journal = %#v", final)
	}
	if replay, err := service.Promote(t.Context(), "build demo"); err != nil || replay != status || owner.executions != executions || owner.validations != validations {
		t.Fatalf("terminal replay = %#v, %v", replay, err)
	}
}

func TestPromoteRefusesUnreleasedRunBeforeRecomposition(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	before := setAssignmentsReleased(t, fixture.service, false)
	write(t, filepath.Join(fixture.root, "working-advance.go"), "package advance\n")
	git(t, fixture.root, "add", ".")
	git(t, fixture.root, "commit", "-qm", "working advance")
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err == nil {
		t.Fatal("active run recomposed")
	}
	after := loadRun(t, fixture.service)
	if after.CandidateTip != before.CandidateTip || after.Base != before.Base || git(t, fixture.root, "rev-parse", after.Candidate) != before.CandidateTip {
		t.Fatal("active run mutated before promotion eligibility")
	}
}

func reviewedPromotionFixture(t *testing.T) checkpointFixture {
	t.Helper()
	fixture := checkpointedReleaseFixture(t)
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatal(err)
	}
	run := loadRun(t, fixture.service)
	receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	if _, err := fixture.service.Review(t.Context(), "build demo", writeReviewReceipt(t, receipt)); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func (g *promotionGate) Execute(_ context.Context, root, tree string) (GateOutcome, error) {
	g.executions++
	g.tree = tree
	if g.inspect != nil {
		g.inspect(root, tree)
	}
	if g.err != nil {
		return GateOutcome{}, g.err
	}
	outcome := GateOutcome{Green: !g.red, Disposition: g.disposition, Evidence: "owner-proof"}
	if outcome.Green && !g.contradictory {
		outcome.Disposition = ""
	}
	return outcome, nil
}

func (o *abandonOwner) PlanAbandon(_ context.Context, root, request, path string) (string, error) {
	o.plans++
	return worktree.PlanAbandon(root, request, path)
}
func (o *abandonOwner) ApplyAbandon(_ context.Context, root, request, path, fingerprint string) error {
	o.applies++
	_, err := worktree.ApplyAbandon(root, request, path, fingerprint)
	return err
}

func abandonmentSnapshot(t *testing.T, fixture checkpointFixture) abandonSnapshot {
	t.Helper()
	statePath, err := fixture.service.statePath("build demo")
	if err != nil {
		t.Fatal(err)
	}
	state, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	return abandonSnapshot{
		state:     string(state),
		refs:      git(t, fixture.root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/"),
		worktrees: git(t, fixture.root, "worktree", "list", "--porcelain"),
	}
}
