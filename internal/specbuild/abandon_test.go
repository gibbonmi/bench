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
	run, found, err := fixture.service.load("build demo")
	if err != nil || !found || !run.Terminal || len(run.Assignments) != 1 {
		t.Fatalf("terminal evidence = %#v, found=%t err=%v", run, found, err)
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
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R26] second abandonment assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second abandonment ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatal(err)
	}
	first, _, err := service.Assign(t.Context(), "build demo", "one.md", "first")
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := service.Assign(t.Context(), "build demo", "two.md", "second")
	if err != nil {
		t.Fatal(err)
	}
	for _, assigned := range []Assignment{first, second} {
		write(t, filepath.Join(assigned.Path, "internal", "specbuild", "unlanded.go"), "package specbuild\n")
	}
	owner := &abandonOwner{}
	service.worktrees = owner
	plan, err := service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	service.fault = func(point string) error {
		if point == "abandon/release" {
			return errors.New("interrupt after cleanup")
		}
		return nil
	}
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
	run, found, err := service.load("build demo")
	if err != nil || !found || !run.Terminal {
		t.Fatalf("resumed terminal state=%#v found=%t err=%v", run, found, err)
	}
	if _, err := json.Marshal(run); err != nil {
		t.Fatalf("retained terminal evidence does not encode: %v", err)
	}
}

func TestApplyAbandonRecoversAfterOwnerApplyBeforeState(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R26] owner apply recovery assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second owner apply recovery ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatal(err)
	}
	first, _, err := service.Assign(t.Context(), "build demo", "one.md", "first")
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := service.Assign(t.Context(), "build demo", "two.md", "second")
	if err != nil {
		t.Fatal(err)
	}
	for _, assigned := range []Assignment{first, second} {
		write(t, filepath.Join(assigned.Path, "internal", "specbuild", "unlanded.go"), "package specbuild\n")
	}
	owner := &abandonOwner{}
	service.worktrees = owner
	plan, err := service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	service.fault = func(point string) error {
		if point == "abandon/owner-apply" {
			return errors.New("interrupt after owner apply")
		}
		return nil
	}
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
	run, found, err := service.load("build demo")
	if err != nil || !found || !run.Terminal {
		t.Fatalf("owner-apply terminal state=%#v found=%t err=%v", run, found, err)
	}
	op, found := service.operation(run, "abandon", "apply")
	journal := abandonmentJournal{}
	if !found || op.State != "completed" || json.Unmarshal([]byte(op.Result), &journal) != nil || journal.Original.Fingerprint != plan.Fingerprint || len(journal.Current.Worktrees) != 0 {
		t.Fatalf("terminal abandonment journal = %#v, found=%t decoded=%#v", op, found, journal)
	}
}

func TestApplyAbandonResumesCleanReleasedAssignment(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R26] clean abandonment assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second clean abandonment ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(t.Context(), "build demo"); err != nil {
		t.Fatal(err)
	}
	for _, ticket := range []string{"one.md", "two.md"} {
		if _, _, err := service.Assign(t.Context(), "build demo", ticket, ticket); err != nil {
			t.Fatal(err)
		}
	}
	owner := &abandonOwner{}
	service.worktrees = owner
	plan, err := service.Abandon(t.Context(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	service.fault = func(point string) error {
		if point == "abandon/release" {
			return errors.New("interrupt after clean cleanup")
		}
		return nil
	}
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

type abandonOwner struct{ plans, applies int }

func (*abandonOwner) Create(context.Context, string, string, string, string) (OwnedWorktree, error) {
	return OwnedWorktree{}, errors.New("unexpected create")
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

type abandonSnapshot struct{ state, refs, worktrees string }

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
