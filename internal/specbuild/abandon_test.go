package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/worktree"
)

const ownershipRefusal = "spec build assignment ownership does not match durable state"

// absentWorktree deletes the fixture's unreleased assignment checkout while the run
// record and the worktree registration still name it, and counts the owner calls the
// resulting cleanup makes.
func absentWorktree(t *testing.T, fixture checkpointFixture) (checkpointFixture, *abandonOwner) {
	t.Helper()
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	if err := os.RemoveAll(fixture.assigned.Path); err != nil {
		t.Fatalf("remove assignment worktree: %v", err)
	}
	return fixture, owner
}

func TestAbandonPlansForRemovedWorktree(t *testing.T) {
	fixture, _ := absentWorktree(t, newCheckpointFixture(t))
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	if len(plan.Worktrees) != 1 || plan.Worktrees[0].ID != fixture.assigned.ID || plan.Worktrees[0].Path != fixture.assigned.Path {
		t.Fatalf("planned worktrees = %#v", plan.Worktrees)
	}
}

func TestAbandonAppliesForRemovedWorktree(t *testing.T) {
	fixture, owner := absentWorktree(t, newCheckpointFixture(t))
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil || status.State != "terminal" || owner.applies != 1 {
		t.Fatalf("ApplyAbandon status=%#v err=%v apply calls=%d", status, err, owner.applies)
	}
	if run := loadRun(t, fixture.service); !run.Terminal {
		t.Fatalf("removed-worktree terminal evidence = %#v", run)
	}
}

func TestRemovedWorktreeRecoveryRefsSurvive(t *testing.T) {
	fixture := newCheckpointFixture(t)
	fixture.service.worktrees = &abandonOwner{}
	run := loadRun(t, fixture.service)
	_, stored, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	owned, found, err := intent.FindAssignmentByRequest(fixture.root, stored.OwnerRequest)
	if err != nil || !found {
		t.Fatalf("owned assignment: found=%v err=%v", found, err)
	}
	ref := intent.RecoveryRefPrefix(owned.OwnerID, owned.ID) + "payload"
	git(t, fixture.root, "update-ref", ref, "HEAD")
	object := git(t, fixture.root, "rev-parse", ref)
	if err := os.RemoveAll(fixture.assigned.Path); err != nil {
		t.Fatalf("remove assignment worktree: %v", err)
	}
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	if want := (AbandonmentRef{Name: ref, Object: object}); len(plan.RecoveryRefs) != 1 || plan.RecoveryRefs[0] != want {
		t.Fatalf("planned recovery refs = %#v, want %#v", plan.RecoveryRefs, want)
	}
	if _, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err != nil {
		t.Fatalf("ApplyAbandon: %v", err)
	}
	if got := git(t, fixture.root, "rev-parse", ref); got != object {
		t.Fatalf("recovery ref %s = %q after abandonment, want %q", ref, got, object)
	}
}

func TestAbandonStillRefusesForgedAssignmentIdentity(t *testing.T) {
	for _, test := range []struct {
		name  string
		forge func(first, second *assignment)
	}{
		{"duplicate id", func(first, second *assignment) { second.ID = first.ID }},
		{"duplicate path", func(first, second *assignment) { second.Path = first.Path }},
		{"duplicate owner request", func(first, second *assignment) { second.OwnerRequest = first.OwnerRequest }},
		{"owner request digest", func(first, _ *assignment) { first.OwnerRequest = digest("forged request") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, service, owner, plan := twoAssignmentAbandonFixture(t, false)
			run := loadRun(t, service)
			keys := make([]string, 0, len(run.Assignments))
			for key := range run.Assignments {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			if len(keys) != 2 {
				t.Fatalf("assignment keys = %#v", keys)
			}
			first, second := run.Assignments[keys[0]], run.Assignments[keys[1]]
			test.forge(&first, &second)
			run.Assignments[keys[0]], run.Assignments[keys[1]] = first, second
			saveRun(t, service, run)
			if _, err := service.Abandon(t.Context(), "build demo"); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
				t.Fatalf("Abandon on forged %s = %v, want the ownership refusal", test.name, err)
			}
			if _, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
				t.Fatalf("ApplyAbandon on forged %s = %v, want the ownership refusal", test.name, err)
			}
			if owner.applies != 0 {
				t.Fatalf("forged %s reached apply: %d", test.name, owner.applies)
			}
		})
	}
}

func TestAbandonRefusesForeignAssignmentCheckout(t *testing.T) {
	fixture, owner := absentWorktree(t, newCheckpointFixture(t))
	stranger := repo(t)
	if err := os.Rename(stranger, fixture.assigned.Path); err != nil {
		t.Fatalf("plant a stranger checkout: %v", err)
	}
	if _, err := fixture.service.Abandon(t.Context(), "build demo"); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
		t.Fatalf("Abandon over a foreign checkout = %v, want the ownership refusal", err)
	}
	if owner.plans != 0 {
		t.Fatalf("foreign checkout reached the owner plan: %d", owner.plans)
	}
}

// startInvocation, assignInvocation, checkpointInvocation, integrateInvocation,
// promoteInvocation, and writeCurrentReviewReceipt drive each precondition-gated
// mutation through the exact service call its production code path takes. The
// absent-worktree and moved-tip regression suites below share them so every lifecycle
// mutation's call shape has one source rather than two independently wired copies.
func startInvocation(t *testing.T, service *Service) error {
	_, err := service.Start(t.Context(), "build demo")
	return err
}

func assignInvocation(t *testing.T, service *Service, request string) error {
	_, _, err := service.Assign(t.Context(), "build demo", "one.md", request)
	return err
}

func checkpointInvocation(t *testing.T, fixture checkpointFixture) error {
	_, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, writeCheckpointReceipt(t, fixture.receipt, "\n"))
	return err
}

func integrateInvocation(t *testing.T, fixture checkpointFixture) error {
	_, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID)
	return err
}

func promoteInvocation(t *testing.T, fixture checkpointFixture) error {
	_, err := fixture.service.Promote(t.Context(), "build demo")
	return err
}

func writeCurrentReviewReceipt(t *testing.T, fixture checkpointFixture) string {
	t.Helper()
	run := loadRun(t, fixture.service)
	receipt := reviewReceipt{Version: 1, Run: run.Run, Candidate: run.CandidateTip, Axes: []reviewAxis{{Axis: "Standards"}, {Axis: "Spec"}, {Axis: "Coverage"}}}
	return writeReviewReceipt(t, receipt)
}

// TestNonAbandonMutationsStillRefuseAbsentWorktree pins that the liveness exemption is
// scoped to abandon: it walks the production mutation list itself, so a mutation added
// later without a matching invocation here fails loudly instead of silently inheriting
// the right to write into a checkout that is gone.
func TestNonAbandonMutationsStillRefuseAbsentWorktree(t *testing.T) {
	invoke := map[mutation]func(t *testing.T) error{
		mutationStart: func(t *testing.T) error {
			fixture, _ := absentWorktree(t, newCheckpointFixture(t))
			return startInvocation(t, fixture.service)
		},
		mutationAssign: func(t *testing.T) error {
			fixture, _ := absentWorktree(t, newCheckpointFixture(t))
			return assignInvocation(t, fixture.service, "absent worktree request")
		},
		mutationCheckpoint: func(t *testing.T) error {
			fixture, _ := absentWorktree(t, newCheckpointFixture(t))
			return checkpointInvocation(t, fixture)
		},
		mutationIntegrate: func(t *testing.T) error {
			fixture, _ := absentWorktree(t, checkpointedReleaseFixture(t))
			return integrateInvocation(t, fixture)
		},
		mutationReview: func(t *testing.T) error {
			fixture, _ := absentWorktree(t, newCheckpointFixture(t))
			_, err := fixture.service.Review(t.Context(), "build demo", writeCurrentReviewReceipt(t, fixture))
			return err
		},
		mutationPromote: func(t *testing.T) error {
			fixture, _ := absentWorktree(t, newCheckpointFixture(t))
			return promoteInvocation(t, fixture)
		},
	}
	for _, op := range lifecycleMutations {
		if op == mutationAbandon {
			continue
		}
		fn, ok := invoke[op]
		if !ok {
			t.Fatalf("no absent-worktree invocation wired for mutation %q", op)
		}
		t.Run(string(op), func(t *testing.T) {
			if err := fn(t); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
				t.Fatalf("%s with an absent assignment worktree = %v, want the ownership refusal", op, err)
			}
		})
	}
}

func TestRemovedWorktreeWithHostilePathIsPlannedAndApplied(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", filepath.Join(home, "a b*c"))
	sibling := filepath.Join(home, "a bXc", "worktrees", "neighbour")
	write(t, filepath.Join(sibling, "keep.txt"), "sibling payload\n")
	fixture, owner := absentWorktree(t, newCheckpointFixture(t))
	if !strings.Contains(fixture.assigned.Path, "a b*c") {
		t.Fatalf("assignment path = %q, want the hostile pool component", fixture.assigned.Path)
	}
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	if len(plan.Worktrees) != 1 || plan.Worktrees[0].Path != fixture.assigned.Path {
		t.Fatalf("planned worktrees = %#v", plan.Worktrees)
	}
	status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil || status.State != "terminal" || owner.applies != 1 {
		t.Fatalf("ApplyAbandon status=%#v err=%v apply calls=%d", status, err, owner.applies)
	}
	if _, err := os.Stat(filepath.Join(sibling, "keep.txt")); err != nil {
		t.Fatalf("glob-sibling worktree was touched: %v", err)
	}
}

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
func TestApplyAbandonSucceedsOnMovedTip(t *testing.T) {
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	advanceWorking(t, fixture.root)
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil {
		t.Fatalf("ApplyAbandon on moved tip: %v", err)
	}
	if status.State != "terminal" || owner.applies != 1 {
		t.Fatalf("moved-tip apply status=%#v apply calls=%d", status, owner.applies)
	}
	if run := loadRun(t, fixture.service); !run.Terminal {
		t.Fatalf("moved-tip terminal evidence = %#v", run)
	}
}
func TestApplyAbandonRefusesDriftedFingerprintOnMovedTip(t *testing.T) {
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	git(t, fixture.root, "update-ref", "refs/bench/specbuild/drift", "HEAD")
	advanceWorking(t, fixture.root)
	if _, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil || !strings.Contains(err.Error(), "plan drifted") {
		t.Fatalf("ApplyAbandon drifted fingerprint on moved tip = %v", err)
	}
	if owner.applies != 0 {
		t.Fatalf("apply calls after drift = %d, want 0", owner.applies)
	}
}
func TestApplyAbandonStillRefusesIdentityDriftOnMovedTip(t *testing.T) {
	for _, test := range []struct {
		name, want string
		drift      func(t *testing.T, fixture checkpointFixture)
	}{
		{"branch", "working checkout does not match recorded subject", func(t *testing.T, fixture checkpointFixture) {
			run := loadRun(t, fixture.service)
			run.Branch = "other"
			saveRun(t, fixture.service, run)
		}},
		{"spec identity", "staged spec no longer matches recorded subject", func(t *testing.T, fixture checkpointFixture) {
			run := loadRun(t, fixture.service)
			run.SpecTip = "changed"
			saveRun(t, fixture.service, run)
		}},
		{"candidate ref", "candidate no longer matches durable tip", func(t *testing.T, fixture checkpointFixture) {
			run := loadRun(t, fixture.service)
			tree := git(t, fixture.root, "rev-parse", run.Candidate+"^{tree}")
			commit := git(t, fixture.root, "commit-tree", tree, "-p", run.Candidate, "-m", "candidate drift")
			git(t, fixture.root, "update-ref", run.Candidate, commit, run.CandidateTip)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			owner := &abandonOwner{}
			fixture.service.worktrees = owner
			plan, err := fixture.service.Abandon(t.Context(), "build demo")
			if err != nil {
				t.Fatalf("Abandon: %v", err)
			}
			advanceWorking(t, fixture.root)
			test.drift(t, fixture)
			if _, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ApplyAbandon %s drift on moved tip = %v, want %q", test.name, err, test.want)
			}
			if owner.applies != 0 {
				t.Fatalf("%s drift reached apply: %d", test.name, owner.applies)
			}
		})
	}
}
func TestApplyAbandonRefusesUnrecognizedHeadMove(t *testing.T) {
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon: %v", err)
	}
	rewriteWorkingHead(t, fixture.root)
	_, err = fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err == nil || !strings.Contains(err.Error(), "does not match recorded subject") || strings.Contains(err.Error(), "bench spec build promote") {
		t.Fatalf("ApplyAbandon unrecognized head move = %v, want the subject-mismatch refusal rather than recomposition", err)
	}
	if owner.applies != 0 {
		t.Fatalf("unrecognized head move reached apply: %d", owner.applies)
	}
}

// TestNonAbandonMutationsStillRecomposeOnMovedTip pins how narrowly a moved tip is
// softened: abandon is exempt at the recomposition return site, checkpoint and start
// fast-forward an empty run onto the moved tip, and every other mutation still refuses.
// It walks the production mutation list itself, so a mutation added later without a
// matching invocation here fails loudly instead of silently joining either set.
func TestNonAbandonMutationsStillRecomposeOnMovedTip(t *testing.T) {
	fastForwards := map[mutation]bool{mutationStart: true, mutationCheckpoint: true}
	invoke := map[mutation]func(t *testing.T) error{
		mutationStart: func(t *testing.T) error {
			fixture := newPreconditionFixture(t, true)
			advanceWorking(t, fixture.root)
			return startInvocation(t, fixture.service)
		},
		mutationAssign: func(t *testing.T) error {
			fixture := newPreconditionFixture(t, true)
			advanceWorking(t, fixture.root)
			return assignInvocation(t, fixture.service, "moved-tip request")
		},
		mutationCheckpoint: func(t *testing.T) error {
			fixture := newCheckpointFixture(t)
			advanceWorking(t, fixture.root)
			return checkpointInvocation(t, fixture)
		},
		mutationIntegrate: func(t *testing.T) error {
			fixture := checkpointedReleaseFixture(t)
			advanceWorking(t, fixture.root)
			return integrateInvocation(t, fixture)
		},
		mutationReview: func(t *testing.T) error {
			fixture := checkpointedReleaseFixture(t)
			if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
				t.Fatalf("Integrate: %v", err)
			}
			evidence := writeCurrentReviewReceipt(t, fixture)
			advanceWorking(t, fixture.root)
			_, err := fixture.service.Review(t.Context(), "build demo", evidence)
			return err
		},
	}
	for _, op := range lifecycleMutations {
		if op == mutationAbandon || op == mutationPromote {
			continue
		}
		fn, ok := invoke[op]
		if !ok {
			t.Fatalf("no moved-tip invocation wired for mutation %q", op)
		}
		t.Run(string(op), func(t *testing.T) {
			err := fn(t)
			if fastForwards[op] {
				if err != nil {
					t.Fatalf("%s on an empty moved tip = %v, want the fast-forward to let it proceed", op, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), "bench spec build promote") {
				t.Fatalf("%s on moved tip = %v, want the recomposition refusal naming promote", op, err)
			}
		})
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

func TestPromoteRecomposesUnreleasedRunOnMovedTip(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	before := setAssignmentsReleased(t, fixture.service, false)
	write(t, filepath.Join(fixture.root, "working-advance.go"), "package advance\n")
	git(t, fixture.root, "add", ".")
	git(t, fixture.root, "commit", "-qm", "working advance")
	working := git(t, fixture.root, "rev-parse", "HEAD")
	if _, err := fixture.service.Promote(t.Context(), "build demo"); err != nil {
		t.Fatalf("unreleased run refused recomposition: %v", err)
	}
	after := loadRun(t, fixture.service)
	if after.CandidateTip == before.CandidateTip || after.Base != working || git(t, fixture.root, "rev-parse", after.Candidate) != after.CandidateTip {
		t.Fatalf("unreleased run recomposition = %#v", after)
	}
}

func reviewedPromotionFixture(t *testing.T, configure ...func(string)) checkpointFixture {
	t.Helper()
	fixture := checkpointedReleaseFixture(t, configure...)
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

// decayedOwner answers abandon planning with a stable synthetic fingerprint instead of
// the real planner, so these fixtures grade only what the precondition classification
// admits and the owner is the counter that proves the release was reached. The composed
// contract — the same shapes carried through worktree.PlanAbandon/ApplyAbandon — is held
// by TestAbandonAppliesForDecayedShapesThroughRealPlanner.
type decayedOwner struct {
	realOwner
	plans, applies int
}

func (o *decayedOwner) PlanAbandon(_ context.Context, _, request, path string) (string, error) {
	o.plans++
	return digest(request + "\x00" + path), nil
}
func (o *decayedOwner) ApplyAbandon(_ context.Context, _, _, _, _ string) error {
	o.applies++
	return nil
}

// decayedFixture builds a checkpoint fixture whose assignment checkout has decayed into
// shape, with an owner that counts the release the decay is supposed to still permit.
func decayedFixture(t *testing.T, shape func(t *testing.T, path string)) (checkpointFixture, *decayedOwner) {
	t.Helper()
	fixture := newCheckpointFixture(t)
	owner := &decayedOwner{}
	fixture.service.worktrees = owner
	shape(t, fixture.assigned.Path)
	return fixture, owner
}

// huskCheckout strips the checkout's git metadata while leaving its bytes in place — the
// decay abandon exists to clean up, and the shape that must never be force-removed.
func huskCheckout(t *testing.T, path string) {
	t.Helper()
	if err := os.RemoveAll(filepath.Join(path, ".git")); err != nil {
		t.Fatalf("strip checkout metadata: %v", err)
	}
	write(t, filepath.Join(path, "husk-payload.txt"), "unlanded husk bytes\n")
}

func danglingSymlinkCheckout(t *testing.T, path string) {
	t.Helper()
	replaceCheckout(t, path)
	if err := os.Symlink(filepath.Join(path+"-target", "gone"), path); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("cannot create a symlink: %v", err))
	}
}

func fifoCheckout(t *testing.T, path string) {
	t.Helper()
	replaceCheckout(t, path)
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("cannot create a FIFO: %v", err))
	}
}

func regularFileCheckout(t *testing.T, path string) {
	t.Helper()
	replaceCheckout(t, path)
	write(t, path, "not a checkout\n")
}

func replaceCheckout(t *testing.T, path string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("clear assignment checkout: %v", err)
	}
}

func TestAbandonAppliesForHuskCheckout(t *testing.T) {
	fixture, owner := decayedFixture(t, huskCheckout)
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon over a husk: %v", err)
	}
	if len(plan.Worktrees) != 1 || plan.Worktrees[0].Path != fixture.assigned.Path {
		t.Fatalf("planned worktrees = %#v", plan.Worktrees)
	}
	status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil || status.State != "terminal" || owner.applies != 1 {
		t.Fatalf("ApplyAbandon over a husk status=%#v err=%v apply calls=%d", status, err, owner.applies)
	}
	run := loadRun(t, fixture.service)
	_, released, ok := assignmentFor(run, fixture.assigned.ID)
	if !run.Terminal || !ok || !released.Released {
		t.Fatalf("husk abandonment evidence = %#v", run)
	}
}

// TestDecayedCheckoutShapesAbandonAndRefuseCheckpoint walks the decayed shapes through
// both directions of the classification at once, so a classifier that special-cases one
// shape — or that softens the shared probe for every operation — cannot pass. The FIFO
// case has no writer, so an implementation that opens the assignment path hangs here
// rather than returning a verdict.
func TestDecayedCheckoutShapesAbandonAndRefuseCheckpoint(t *testing.T) {
	for _, test := range []struct {
		name  string
		shape func(t *testing.T, path string)
	}{
		{"husk", huskCheckout},
		{"dangling symlink", danglingSymlinkCheckout},
		{"fifo", fifoCheckout},
		{"regular file", regularFileCheckout},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture, owner := decayedFixture(t, test.shape)
			plan, err := fixture.service.Abandon(t.Context(), "build demo")
			if err != nil {
				t.Fatalf("Abandon over a %s: %v", test.name, err)
			}
			status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
			if err != nil || status.State != "terminal" || owner.applies != 1 {
				t.Fatalf("ApplyAbandon over a %s status=%#v err=%v apply calls=%d", test.name, status, err, owner.applies)
			}
			checkpointed, _ := decayedFixture(t, test.shape)
			if err := checkpointInvocation(t, checkpointed); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
				t.Fatalf("Checkpoint over a %s = %v, want the ownership refusal", test.name, err)
			}
		})
	}
}

// requireUnreadableMetadata strips the checkout's git metadata entry and proves the strip
// bit, since root ignores the mode entirely and would read straight through the
// assertion. The restore is registered before the chmod so it runs ahead of TempDir's own
// removal, which cannot descend into an entry it cannot enter.
func requireUnreadableMetadata(t *testing.T, path string) {
	t.Helper()
	metadata := filepath.Join(path, ".git")
	t.Cleanup(func() { _ = os.Chmod(metadata, 0o700) })
	if err := os.Chmod(metadata, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip permissions: %v", err))
	}
	if entry, err := os.Open(metadata); err == nil {
		entry.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}
}

// TestAbandonRefusesUnreadableCheckoutMetadata pins the classification's fatal side: a
// path that still carries git metadata may be a live checkout this process merely cannot
// read, so no operation — abandon included — may release it.
func TestAbandonRefusesUnreadableCheckoutMetadata(t *testing.T) {
	fixture, owner := decayedFixture(t, func(t *testing.T, path string) { requireUnreadableMetadata(t, path) })
	if _, err := fixture.service.Abandon(t.Context(), "build demo"); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
		t.Fatalf("Abandon over unreadable metadata = %v, want the ownership refusal", err)
	}
	if err := checkpointInvocation(t, fixture); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
		t.Fatalf("Checkpoint over unreadable metadata = %v, want the ownership refusal", err)
	}
	if owner.plans != 0 || owner.applies != 0 {
		t.Fatalf("unreadable metadata reached the owner: plans=%d applies=%d", owner.plans, owner.applies)
	}
}

// realAbandonFixture stages a decayed checkout under the production planner: the owner
// delegates to internal/worktree rather than answering with a synthetic fingerprint, so
// the shape policy the fixture stages is the one the planner itself enforces.
func realAbandonFixture(t *testing.T, shape func(t *testing.T, path string)) (checkpointFixture, *abandonOwner) {
	t.Helper()
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	shape(t, fixture.assigned.Path)
	return fixture, owner
}

// TestAbandonAppliesForDecayedShapesThroughRealPlanner composes the decayed shapes with
// the real planner, so a state the synthetic fingerprint admits but internal/worktree
// refuses cannot pass unseen. The shapes are the two the classifier decides without
// privilege: a directory that has lost its git metadata, and a path that is no directory
// at all.
func TestAbandonAppliesForDecayedShapesThroughRealPlanner(t *testing.T) {
	for _, test := range []struct {
		name  string
		shape func(t *testing.T, path string)
	}{
		{"husk", huskCheckout},
		{"regular file", regularFileCheckout},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture, owner := realAbandonFixture(t, test.shape)
			plan, err := fixture.service.Abandon(t.Context(), "build demo")
			if err != nil {
				t.Fatalf("Abandon over a %s: %v", test.name, err)
			}
			if len(plan.Worktrees) != 1 || plan.Worktrees[0].Path != fixture.assigned.Path || owner.plans == 0 {
				t.Fatalf("planned worktrees = %#v, owner plan calls = %d", plan.Worktrees, owner.plans)
			}
			status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
			if err != nil || status.State != "terminal" || owner.applies != 1 {
				t.Fatalf("ApplyAbandon over a %s status=%#v err=%v apply calls=%d", test.name, status, err, owner.applies)
			}
			run := loadRun(t, fixture.service)
			_, released, ok := assignmentFor(run, fixture.assigned.ID)
			if !run.Terminal || !ok || !released.Released {
				t.Fatalf("%s abandonment evidence = %#v", test.name, run)
			}
		})
	}
}

// TestAbandonAppliesWithMissingIntentRegistration pins the reclassification that makes
// the blanket exemption deletable: a run record naming an assignment the ledger no longer
// holds has no ownership left to write into, and abandon is what clears it. The owner is
// the production composition — a checkout built by worktree.Create, planned and released
// through internal/worktree itself — because the real planner is the half that refuses a
// registration it cannot resolve, and a stubbed one would hide that refusal.
func TestAbandonAppliesWithMissingIntentRegistration(t *testing.T) {
	fixture := newCheckpointFixture(t)
	owner := &abandonOwner{}
	fixture.service.worktrees = owner
	run := loadRun(t, fixture.service)
	_, stored, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	owned, found, err := intent.FindAssignmentByRequest(fixture.root, stored.OwnerRequest)
	if err != nil || !found {
		t.Fatalf("owned assignment: found=%v err=%v", found, err)
	}
	if err := intent.DeleteAssignment(fixture.root, owned.ID); err != nil {
		t.Fatalf("delete intent registration: %v", err)
	}
	plan, err := fixture.service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatalf("Abandon without an intent registration: %v", err)
	}
	// The ledger entry is the only thing the owner reconciles against, so with it gone
	// there is nothing left to plan or release; the residual bytes stay on disk for the
	// ordinary clean surface to collect.
	if len(plan.Worktrees) != 0 || owner.plans != 0 {
		t.Fatalf("planned worktrees = %#v, owner plan calls = %d", plan.Worktrees, owner.plans)
	}
	status, err := fixture.service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint)
	if err != nil || status.State != "terminal" || owner.applies != 0 {
		t.Fatalf("ApplyAbandon without an intent registration status=%#v err=%v apply calls=%d", status, err, owner.applies)
	}
	final := loadRun(t, fixture.service)
	_, released, ok := assignmentFor(final, fixture.assigned.ID)
	if !final.Terminal || !ok || !released.Released {
		t.Fatalf("terminal evidence without an intent registration = %#v", final)
	}
}

// preparedAbandonFixture drives an abandon apply into its interrupted mid-release state:
// the abandon operation is prepared with a recorded journal and exactly one assignment is
// still unreleased. That is the exact state a blanket prepared-abandon exemption re-admits
// an identity fault through.
func preparedAbandonFixture(t *testing.T) (*Service, *abandonOwner, AbandonmentPlan) {
	t.Helper()
	_, service, owner, plan := twoAssignmentAbandonFixture(t, false)
	service.fault = injectFault("abandon/release")
	if _, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil {
		t.Fatal("ApplyAbandon completed through injected interruption")
	}
	service.fault = nil
	op, found := service.operation(loadRun(t, service), "abandon", "apply")
	if !found || op.State != "prepared" || op.Result == "" || owner.applies != 1 {
		t.Fatalf("prepared abandon operation = %#v found=%t apply calls=%d", op, found, owner.applies)
	}
	return service, owner, plan
}

// preparedAssignments splits the interrupted run's assignments into the one already
// released and the one still owned, since only the owned assignment's identity is still
// read by the ownership check.
func preparedAssignments(t *testing.T, run record) (releasedKey, ownedKey string) {
	t.Helper()
	for key, assigned := range run.Assignments {
		if assigned.Released {
			releasedKey = key
			continue
		}
		ownedKey = key
	}
	if releasedKey == "" || ownedKey == "" {
		t.Fatalf("interrupted assignments = %#v", run.Assignments)
	}
	return releasedKey, ownedKey
}

// TestAbandonRefusesIdentityFaultsDuringPreparedApply enumerates the fatal identity
// classes against the resumed-abandon state, because one representative forgery would let
// an implementation reject it while still swallowing another class.
func TestAbandonRefusesIdentityFaultsDuringPreparedApply(t *testing.T) {
	for _, test := range []struct {
		name  string
		forge func(t *testing.T, released, owned *assignment)
	}{
		{"duplicate id", func(_ *testing.T, released, owned *assignment) { owned.ID = released.ID }},
		{"duplicate path", func(_ *testing.T, released, owned *assignment) { owned.Path = released.Path }},
		{"duplicate owner request", func(_ *testing.T, released, owned *assignment) { owned.OwnerRequest = released.OwnerRequest }},
		{"owner request digest", func(_ *testing.T, _, owned *assignment) { owned.OwnerRequest = digest("forged request") }},
		{"registration assignment id", func(_ *testing.T, _, owned *assignment) { owned.ID = "forged-assignment-id" }},
		{"registration worktree path", func(t *testing.T, _, owned *assignment) {
			owned.Path = filepath.Join(t.TempDir(), "elsewhere")
		}},
		{"foreign checkout", func(t *testing.T, _, owned *assignment) {
			if err := os.RemoveAll(owned.Path); err != nil {
				t.Fatalf("clear assignment checkout: %v", err)
			}
			if err := os.Rename(repo(t), owned.Path); err != nil {
				t.Fatalf("plant a stranger checkout: %v", err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			service, owner, plan := preparedAbandonFixture(t)
			run := loadRun(t, service)
			releasedKey, ownedKey := preparedAssignments(t, run)
			released, owned := run.Assignments[releasedKey], run.Assignments[ownedKey]
			test.forge(t, &released, &owned)
			run.Assignments[releasedKey], run.Assignments[ownedKey] = released, owned
			saveRun(t, service, run)
			if _, err := service.Abandon(t.Context(), "build demo"); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
				t.Fatalf("Abandon on forged %s = %v, want the ownership refusal", test.name, err)
			}
			if _, err := service.ApplyAbandon(t.Context(), "build demo", plan.Fingerprint); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
				t.Fatalf("ApplyAbandon on forged %s = %v, want the ownership refusal", test.name, err)
			}
			if owner.applies != 1 {
				t.Fatalf("forged %s reached a further apply: %d", test.name, owner.applies)
			}
		})
	}
}
